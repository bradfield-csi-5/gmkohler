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

var (
	notFoundError *leveldb.NotFoundError
)

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
					"found key %q in both the mem-table and tombstones",
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

type SSTableDB struct {
	readSeeker      io.ReadSeeker
	endOfDataOffset int64
	dir             *Directory
}

func (db *SSTableDB) Get(searchKey leveldb.Key) (leveldb.Value, error) {
	var (
		entry = new(encoding.Entry)
		err   error
	)

	if err = db.scanTowards(searchKey); err != nil {
		return nil, err
	}

	for !db.isAtEndOfData() {
		bytesRead, err := readEntry(db.readSeeker, entry)
		if err != nil {
			_, _ = db.readSeeker.Seek(-bytesRead, io.SeekCurrent)
			return nil, err
		}

		var comparison = bytes.Compare(entry.Key, searchKey)
		if comparison > 0 {
			return nil, leveldb.NewNotFoundError(searchKey)
		} else if comparison == 0 {
			if len(entry.Value) == 0 {
				// valLen == 0 implies the key has been tombstoned
				// FIXME: consider multiple tables
				return nil, leveldb.NewNotFoundError(searchKey)
			}

			break
		} else {
			entry = new(encoding.Entry)
			continue
		}
	}
	if entry.IsZeroEntry() {
		return nil, leveldb.NewNotFoundError(searchKey)
	} else {
		return leveldb.Value(entry.Value), nil
	}
}

func (db *SSTableDB) Has(key leveldb.Key) (bool, error) {
	_, err := db.Get(key)
	if err != nil {
		if errors.As(err, &notFoundError) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (db *SSTableDB) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	var (
		encodingStart = encoding.Key(start)
		currentEntry  = new(encoding.Entry)
		err           error
		bytesRead     int64
	)
	if err = db.scanTowards(start); err != nil {
		return nil, err
	}
	for !db.isAtEndOfData() {
		bytesRead, err = readEntry(db.readSeeker, currentEntry)
		if err != nil {
			return nil, err
		}
		if currentEntry.Key.Compare(encodingStart) < 0 {
			continue
		}
		// we have found the entry gte our key, let's rewind so Next() returns it
		_, err := db.readSeeker.Seek(-bytesRead, io.SeekCurrent)
		if err != nil {
			return nil, err
		}
		break
	}

	return NewIterator(
		db.readSeeker,
		encoding.Key(limit),
		db.endOfDataOffset,
	), nil
}

// scanTowards scans to the key in the sparse index that's closest to searchKey (less than or equal to)
func (db *SSTableDB) scanTowards(searchKey leveldb.Key) error {
	startIndex, err := db.dir.offsetFor(searchKey)
	if err != nil {
		return err
	}
	if _, err := db.readSeeker.Seek(int64(startIndex), io.SeekStart); err != nil {
		return err
	}
	return nil
}

// isAtEndOfData is a substitute for checking for EOF errors because our file has directory data
// at the end of it.  The data offset is inferred from the encoding of where the directory starts (that encoding lives
// at the beginning of the file)
func (db *SSTableDB) isAtEndOfData() bool {
	currOffset, _ := db.readSeeker.Seek(0, io.SeekCurrent) // no risk to get EOF with these parameters
	return currOffset >= db.endOfDataOffset
}

func NewIterator(
	readSeeker io.ReadSeeker,
	limit encoding.Key,
	endOfDataOffset int64,
) *Iterator {
	return &Iterator{
		readSeeker:      readSeeker,
		limit:           limit,
		endOfDataOffset: endOfDataOffset,
		currentEntry:    new(encoding.Entry), // call Next() first
		err:             nil,
	}
}

// Iterator is used for satisfying a RangeScan.  It is similar to the read functions.
type Iterator struct {
	readSeeker      io.ReadSeeker
	limit           encoding.Key
	endOfDataOffset int64
	currentEntry    *encoding.Entry // should this start at the preceding entry?
	err             error
}

func (i *Iterator) Next() bool {
	if i.isAtEndOfData() {
		return false
	}
	_, err := readEntry(i.readSeeker, i.currentEntry)
	if i.currentEntry.Key.Compare(i.limit) > 0 {
		return false
	}
	if err != nil {
		i.err = err
		return false
	}
	if len(i.currentEntry.Value) == 0 {
		return i.Next() // don't return tombstoned data
	}
	return true
}

func (i *Iterator) Error() error {
	return i.err
}

func (i *Iterator) Key() leveldb.Key {
	return leveldb.Key(i.currentEntry.Key)
}

func (i *Iterator) Value() leveldb.Value {
	return leveldb.Value(i.currentEntry.Value)
}

func (i *Iterator) isAtEndOfData() bool {
	currOffset, _ := i.readSeeker.Seek(0, io.SeekCurrent) // no risk to get EOF with these parameters
	return currOffset > i.endOfDataOffset
}

// readEntry reads an entry into the supplied pointer and returns how many bytes were read
// in case the caller needs to "peek" (in which case they can seek backwards by that number of bytes)
func readEntry(rs io.ReadSeeker, entry *encoding.Entry) (int64, error) {
	var bytesRead int64
	keyLen, err := encoding.ReadUint64(rs)
	bytesRead += 8
	if err != nil {
		return bytesRead, err
	}
	var key = make(encoding.Key, keyLen)
	keyBytes, err := io.ReadFull(rs, key)
	bytesRead += int64(keyBytes)
	if err != nil {
		return bytesRead, err
	}
	valLen, err := encoding.ReadUint64(rs)
	bytesRead += 8
	if err != nil {
		return bytesRead, err
	}
	if valLen == 0 {
		*entry = encoding.Entry{Key: key, Value: nil}
		return bytesRead, nil
	}
	value, err := encoding.ReadByteSlice(rs, valLen)
	bytesRead += int64(valLen)
	if err != nil {
		return bytesRead, err
	}
	*entry = encoding.Entry{
		Key:   key,
		Value: value,
	}
	return bytesRead, nil
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
