package db

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"leveldb"
	"leveldb/encoding"
	"leveldb/skiplist"
	"leveldb/sst"
	"leveldb/wal"
	"os"
)

type ReadWriteSeeker interface {
	io.Reader
	io.Writer
	io.Seeker
}

type db struct {
	sl       *skiplist.SkipList
	ssTables []leveldb.ReadOnlyDB
	wal      *wal.Log
}

func NewDbFromWal(rw io.ReadWriter) (leveldb.DB, error) {
	reader := bufio.NewReader(rw)
	entries, err := encoding.DecodeLogFile(reader)
	if err != nil {
		return nil, err
	}
	db := NewDb(rw)
	for _, entry := range entries {
		switch entry.Operation {
		case encoding.OpPut:
			if err := db.Put(
				leveldb.Key(entry.Key),
				leveldb.Value(entry.Value),
			); err != nil {
				return nil, err
			}
		case encoding.OpDelete:
			if err := db.Delete(leveldb.Key(entry.Key)); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unrecognized opcode %s", entry.Operation)
		}
	}
	return db, nil
}

func NewDb(readWriter io.ReadWriter) leveldb.DB {
	var log *wal.Log
	if readWriter != nil { // hacky, think of nicer way
		log = wal.NewLog(readWriter)
	}

	return &db{
		sl:  skiplist.NewSkipList(),
		wal: log,
	}
}

func (db *db) Get(key leveldb.Key) (leveldb.Value, error) {
	return db.sl.Search(key)
}

func (db *db) Has(key leveldb.Key) (bool, error) {
	val, err := db.sl.Search(key)
	if err != nil { // FIXME: slow because of reflection
		var notFoundError *leveldb.NotFoundError
		if errors.As(err, &notFoundError) {
			return false, nil
		}
		return false, err
	}
	return val != nil, nil
}

func (db *db) Put(key leveldb.Key, value leveldb.Value) error {
	if err := db.wal.Put(key, value); err != nil {
		return err
	}
	return db.sl.Insert(key, value)
}

func (db *db) Delete(key leveldb.Key) error {
	if err := db.wal.Delete(key); err != nil {
		return err
	}
	return db.sl.Delete(key)
}

func (db *db) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	precedingNode, err := db.sl.TraverseUntil(start, nil)
	if err != nil {
		return nil, err
	}
	return NewSkipListIterator(precedingNode, limit), nil
}

func (db *db) flushSSTable(f *os.File) (*sst.SSTableDB, error) {
	return sst.BuildSSTable(f, db.sl)
}

// consider converting to memory instead of this
type skipListIterator struct {
	limit   leveldb.Key
	current skiplist.Node
	err     error
}

func NewSkipListIterator(precedingNode skiplist.Node, limit leveldb.Key) leveldb.Iterator {
	return &skipListIterator{
		limit:   limit,
		current: precedingNode,
	}
}

func (s *skipListIterator) Next() bool {
	if nextNode := s.current.Next(); nextNode.CompareKey(s.limit) <= 0 {
		s.current = nextNode
		return true
	}
	return false
}

func (s *skipListIterator) Error() error {
	return s.err
}

func (s *skipListIterator) Key() leveldb.Key {
	return s.current.Key()
}

func (s *skipListIterator) Value() leveldb.Value {
	return s.current.Value()
}
