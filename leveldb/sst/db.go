package sst

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"leveldb"
	"leveldb/encoding"
	"leveldb/skiplist"
	"os"
)

const (
	sparseIndexThreshold = 0x400 // add new index every 1K bytes written
	dataOffset           = 0x10
)

type SSTableDB struct {
	readSeeker      io.ReadSeeker
	endOfDataOffset int64
	dir             *Directory
}

func newSSTableConfig() *ssTableConfig {
	return &ssTableConfig{sparseIndexThreshold: sparseIndexThreshold}
}

type ssTableConfig struct {
	sparseIndexThreshold int
}
type ssTableOption func(*ssTableConfig)

func withSparseIndexThreshold(threshold int) ssTableOption {
	return func(config *ssTableConfig) {
		config.sparseIndexThreshold = threshold
	}
}

/**
 * format:
 * | 8 bytes (int64)    | 8 bytes 		   | arbitrarily long | arbitrarily long		  |
 * | [directory offset] | directory size   |     [data]       | [directory entries]      |
 *
 * data:
 * | 8 bytes   |  arbitrary |  8 bytes    |  [0, arbitrary) |
 * | [key len] | [key]		| [value len] | (value) 	   |
 *
 * , where [value len] is 0 if key is tombstoned, and value omitted in this case
 *
 * directory entry:
 * | 8 bytes   |  arbitrary |  8 bytes      |
 * | [key len] | [key]		| [file offset] |
 */

// BuildSSTable builds an SSTable from the skiplists for present and deleted entries in a memtable
func BuildSSTable(
	f *os.File,
	memTable *skiplist.SkipList,
	tombstones *skiplist.SkipList,
	configOptions ...ssTableOption,
) (*SSTableDB, error) {
	var ssTableConfig = newSSTableConfig()
	for _, option := range configOptions {
		option(ssTableConfig)
	}
	// LevelDB’s approach is to flush the mem-table to disk once it reaches the mem-table once it reaches some threshold
	// size, and then truncate the write-ahead log to remove any entries involving flushed data. The data is persisted
	// in an immutable format called an “SSTable” (or “sorted string table”).
	// do writing here

	/* TODO: build and write directory and its metadata */
	if _, err := f.Seek(dataOffset, io.SeekStart); err != nil { // start writing from start
		return nil, err
	}

	// sort tombstones
	tombstonedNode, err := tombstones.TraverseUntil(nil, nil)
	if err != nil {
		return nil, err
	}
	tombstonedNode = tombstonedNode.Next()

	memTableNode, err := memTable.TraverseUntil(nil, nil)
	if err != nil {
		return nil, err
	}
	memTableNode = memTableNode.Next()

	// merging loop
	var (
		cumulativeBytesWritten int
		sparseKeys             []leveldb.Key
		keyOffsets             []offset
	)
	for memTableNode != skiplist.NilNode || tombstonedNode != skiplist.NilNode {
		var nodeToEncode skiplist.Node

		// pick which source has next smallest key, increment the winner accordingly.
		if tombstonedNode == skiplist.NilNode {
			nodeToEncode = memTableNode
			memTableNode = memTableNode.Next()
		} else if memTableNode == skiplist.NilNode {
			nodeToEncode = tombstonedNode
			tombstonedNode = tombstonedNode.Next()
		} else {
			var comparison = bytes.Compare(tombstonedNode.Key(), memTableNode.Key())
			switch {
			case comparison < 0:
				nodeToEncode = tombstonedNode
				tombstonedNode = tombstonedNode.Next()
			case comparison > 0:
				nodeToEncode = memTableNode
				memTableNode = memTableNode.Next()
			default:
				return nil, fmt.Errorf(
					"found key %q in both the mem-table and deleted keyset",
					tombstonedNode.Key(),
				)
			}
		}

		entryToEncode := encoding.Entry{
			Key:   encoding.Key(nodeToEncode.Key()),
			Value: encoding.Value(nodeToEncode.Value()),
		}
		encodedEntry, err := entryToEncode.Encode()
		if err != nil {
			return nil, err
		}
		bytesWritten, err := f.Write(encodedEntry)
		if err != nil {
			return nil, err
		}
		if bytesWritten != len(encodedEntry) {
			return nil, errors.New("failed to write all bytes")
		}
		cumulativeBytesWritten += bytesWritten
		if cumulativeBytesWritten > ssTableConfig.sparseIndexThreshold {
			var currentOffset, err = f.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
			sparseKeys = append(sparseKeys, leveldb.Key(entryToEncode.Key))
			keyOffsets = append(keyOffsets, offset(currentOffset-int64(bytesWritten)))
			cumulativeBytesWritten = 0
		}

	}

	directoryOffset, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	// write directory
	directory, err := NewDirectory(sparseKeys, keyOffsets)
	if err != nil {
		return nil, err
	}
	encodedDirectory, err := directory.Encode()
	if err != nil {
		return nil, err
	}
	_, err = f.Write(encodedDirectory)
	if err != nil {
		return nil, err
	}

	// go back to front of file and write metadata
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	if err := encoding.WriteUint64(f, uint64(directoryOffset)); err != nil {
		return nil, err
	}
	if err := encoding.WriteUint64(f, uint64(len(encodedDirectory))); err != nil {
		return nil, err
	}

	return NewSSTableDBFromFile(f)
}

func NewSSTableDBFromFile(readSeeker io.ReadSeeker) (*SSTableDB, error) {
	var (
		bufReader = bufio.NewReader(readSeeker)
		err       error
	)
	if _, err := readSeeker.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("NewSSTableDBFromFile: error seeking to start of file: %v", err)
	}
	// START: read directory metadata
	endOfDataOffset, err := encoding.ReadUint64(bufReader)
	if err != nil {
		return nil, fmt.Errorf("NewSSTableDBFromFile: error reading data offset: %v", err)
	}
	dirLen, err := encoding.ReadUint64(bufReader)
	if err != nil {
		return nil, fmt.Errorf("NewSSTableDBFromFile: error reading directory length: %v", err)
	}
	// END: read directory metadata
	// START: read directory
	if _, err := readSeeker.Seek(int64(endOfDataOffset), io.SeekStart); err != nil {
		return nil, err
	}
	var directory = NewBlankDirectory()
	if dirLen > 0 {
		directoryBuf := make([]byte, dirLen)
		bytesRead, err := readSeeker.Read(directoryBuf)
		if err != nil {
			return nil, fmt.Errorf("sst.NewSSTableDBFromFile: error reading directory contents: %v", err)
		}
		if uint64(bytesRead) != dirLen {
			return nil, fmt.Errorf("sst.NewSSTableDBFromFile: failure to read entire directory.  expected %d bytes, read %d", dirLen, bytesRead)
		}
		if err := directory.Decode(directoryBuf); err != nil {
			return nil, fmt.Errorf("NewSSTableDBFromFile: error decoding directory contents: %v", err)
		}
	}
	// END: read directory
	// reset to start of data
	if _, err = readSeeker.Seek(dataOffset, io.SeekStart); err != nil { // 8 == 2 * size(int64)
		return nil, fmt.Errorf("NewSSTableDBFromFile: error seeking to start of data: %v", err)
	}
	return &SSTableDB{
		readSeeker:      readSeeker,
		endOfDataOffset: int64(endOfDataOffset),
		dir:             directory,
	}, nil
}

func (db *SSTableDB) Get(searchKey leveldb.Key) (leveldb.Value, error) {
	var (
		keyLen     uint64
		valLen     uint64
		value      leveldb.Value
		startIndex offset
		err        error
	)

	startIndex, err = db.dir.offsetFor(searchKey)
	if err != nil {
		return nil, err
	}
	if _, err := db.readSeeker.Seek(int64(startIndex), io.SeekStart); err != nil {
		return nil, err
	}

	for !db.isAtEndOfData() {
		keyLen, err = encoding.ReadUint64(db.readSeeker)
		if err != nil {
			return nil, err
		}
		var key = make(leveldb.Key, keyLen)
		if _, err = io.ReadFull(db.readSeeker, key); err != nil {
			return nil, err
		}
		var comparison = bytes.Compare(searchKey, key)
		if comparison < 0 {
			return nil, leveldb.NewNotFoundError(searchKey)
		}
		valLen, err = encoding.ReadUint64(db.readSeeker)
		if err != nil {
			return nil, err
		}
		if valLen == 0 {
			// valLen == 0 implies the key has been tombstoned
			// FIXME: consider multiple tables
			return nil, leveldb.NewNotFoundError(searchKey)
		}
		if comparison == 0 {
			value, err = encoding.ReadByteSlice(db.readSeeker, valLen)
			if err != nil {
				return nil, err
			}
			break
		} else {
			// skip past value dangerous to cast but shrug
			_, err = db.readSeeker.Seek(int64(valLen), io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		}
	}

	return value, nil
}

// isAtEndOfData is a substitute for checking for EOF errors because our file has directory data
// at the end of it.  The data offset is inferred from the encoding of where the directory starts (that encoding lives
// at the beginning of the file)
func (db *SSTableDB) isAtEndOfData() bool {
	currOffset, _ := db.readSeeker.Seek(0, io.SeekCurrent) // no worry about error in this use-case
	return currOffset >= db.endOfDataOffset
}

func (db *SSTableDB) Has(key leveldb.Key) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (db *SSTableDB) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	//TODO implement me
	panic("implement me")
}
