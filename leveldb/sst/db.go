package sst

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"leveldb"
	"leveldb/encoding"
	"leveldb/skiplist"
	"os"
	"slices"
)

const (
	sparseIndexThreshold = 0x400 // add new index every 1K bytes written
	dataOffset           = 0x10
)

type SSTableDB struct {
	readSeeker      io.ReadSeeker
	endOfDataOffset int64
	dir             Directory
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
func BuildSSTable(f *os.File, skipList *skiplist.SkipList) (*SSTableDB, error) {
	// LevelDB’s approach is to flush the mem-table to disk once it reaches the mem-table once it reaches some threshold
	// size, and then truncate the write-ahead log to remove any entries involving flushed data. The data is persisted
	// in an immutable format called an “SSTable” (or “sorted string table”).
	// do writing here

	/* TODO: build and write directory and its metadata */
	if _, err := f.Seek(dataOffset, io.SeekStart); err != nil { // start writing from start
		return nil, err
	}

	// sort tombstones
	var tombstonedKeys = make([]leveldb.Key, len(skipList.Tombstones))
	var tombstoneIdx int
	for tombstonedKey := range skipList.Tombstones {
		tombstonedKeys[tombstoneIdx] = leveldb.Key(tombstonedKey)
		tombstoneIdx++
	}
	tombstoneIdx = 0
	slices.SortFunc(tombstonedKeys, func(a, b leveldb.Key) int {
		return a.Compare(b)
	})

	memTableNode, err := skipList.TraverseUntil(nil, nil)
	if err != nil {
		return nil, err
	}

	// merging loop
	var (
		cumulativeBytesWritten int
		sparseKeys             []leveldb.Key
		keyOffsets             []offset
	)
	for memTableNode != skiplist.NilNode && tombstoneIdx < len(tombstonedKeys) {
		var entryToEncode encoding.Entry

		// pick which source has next smallest key, increment the winner accordingly.
		if tombstoneIdx >= len(tombstonedKeys) {
			entryToEncode = encoding.Entry{
				Key:   encoding.Key(memTableNode.Key()),
				Value: encoding.Value(memTableNode.Value()),
			}
			memTableNode = memTableNode.Next()
		} else {
			var tombstonedKey = tombstonedKeys[tombstoneIdx]
			var comparison = bytes.Compare(tombstonedKey, memTableNode.Key())
			switch {
			case comparison < 0:
				entryToEncode = encoding.Entry{
					Key:   encoding.Key(memTableNode.Key()),
					Value: encoding.Value(memTableNode.Value()),
				}
				memTableNode = memTableNode.Next()
			case comparison > 0:
				entryToEncode = encoding.Entry{Key: encoding.Key(tombstonedKey)}
				tombstoneIdx++
			default:
				return nil, fmt.Errorf(
					"found key %q in both the mem-table and deleted keyset",
					tombstonedKey,
				)
			}
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
		if cumulativeBytesWritten > sparseIndexThreshold {
			var currentOffset, err = f.Seek(0, io.SeekCurrent)
			if err != nil {
				return nil, err
			}
			sparseKeys = append(sparseKeys, leveldb.Key(entryToEncode.Key))
			keyOffsets = append(keyOffsets, offset(currentOffset)) // FIXME: point directly to value
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
	if err := encoding.WriteInt64(f, directoryOffset); err != nil {
		return nil, err
	}
	if err := encoding.WriteInt64(f, int64(len(encodedDirectory))); err != nil {
		return nil, err
	}

	if err := skipList.Reset(); err != nil {
		return nil, err
	}

	return NewSSTableDBFromFile(f)
}

func NewSSTableDBFromFile(readSeeker io.ReadSeeker) (*SSTableDB, error) {
	var (
		bufReader = bufio.NewReader(readSeeker)
		err       error
	)
	// START: read directory metadata
	endOfDataOffset, err := binary.ReadVarint(bufReader)
	if err != nil {
		return nil, err
	}
	dirLen, err := binary.ReadVarint(bufReader)
	if err != nil {
		return nil, err
	}
	// END: read directory metadata
	// START: read directory
	if _, err := readSeeker.Seek(endOfDataOffset, io.SeekStart); err != nil {
		return nil, err
	}
	directoryBuf := make([]byte, dirLen)
	bytesRead, err := readSeeker.Read(directoryBuf)
	if err != nil {
		return nil, fmt.Errorf("sst.NewSSTableDBFromFile: %v", err)
	}
	if int64(bytesRead) != dirLen {
		return nil, fmt.Errorf("sst.NewSSTableDBFromFile: failure to read entire directory.  expected %d bytes, read %d", dirLen, bytesRead)

	}
	// END: read directory
	// reset to start of data
	if _, err = readSeeker.Seek(16, io.SeekStart); err != nil { // 8 == 2 * size(int64)
		return nil, err
	}
	return &SSTableDB{
		readSeeker:      readSeeker,
		endOfDataOffset: endOfDataOffset,
	}, nil
}

func (db *SSTableDB) Get(searchKey leveldb.Key) (leveldb.Value, error) {
	startIndex, err := db.dir.offsetFor(searchKey)
	if err != nil {
		return nil, err
	}
	if _, err := db.readSeeker.Seek(int64(startIndex), io.SeekStart); err != nil {
		return nil, err
	}

	var (
		keyLen        uint64
		currentOffset int64
	)
	currentOffset, err = db.readSeeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	for !db.isAtEndOfData(currentOffset) {
		keyLen, err = encoding.ReadUint64(db.readSeeker)
		if err != nil {
			return nil, err
		}
		var key = make(leveldb.Key, keyLen)
		if _, err = io.ReadFull(db.readSeeker, key); err != nil {
			return nil, err
		}
		var comparison = bytes.Compare(searchKey, key)
		if comparison > 0 {
			return nil, leveldb.NewNotFoundError(searchKey)
		}
		valLen, err := encoding.ReadUint64(db.readSeeker)
		if err != nil {
			return nil, err
		}
		if comparison == 0 {
			value, err := encoding.ReadByteSlice(db.readSeeker, valLen)
			if err != nil {
				return nil, err
			}
			return value, nil
		} else {
			// skip past value dangerous to cast but shrug
			currentOffset, err = db.readSeeker.Seek(int64(valLen), io.SeekCurrent)
			if err != nil {
				return nil, err
			}
		}
	}
	panic("not done")
}

// isAtEndOfData is a substitute for checking for EOF errors because our file has directory data
// at the end of it.  The data offset is inferred from the encoding of where the directory starts (that encoding lives
// at the beginning of the file)
func (db *SSTableDB) isAtEndOfData(offset int64) bool { return db.endOfDataOffset >= offset }

func (db *SSTableDB) Has(key leveldb.Key) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (db *SSTableDB) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	//TODO implement me
	panic("implement me")
}
