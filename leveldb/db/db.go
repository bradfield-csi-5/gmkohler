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

var (
	notFoundError *leveldb.NotFoundError
)

type ReadWriteSeeker interface {
	io.Reader
	io.Writer
	io.Seeker
}

type db struct {
	memTable   *skiplist.SkipList
	tombstones *skiplist.SkipList
	ssTables   []leveldb.ReadOnlyDB
	wal        *wal.Log
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

func NewDb(walLog io.ReadWriter) leveldb.DB {
	var log *wal.Log
	if walLog != nil { // hacky, think of nicer way
		log = wal.NewLog(walLog)
	}

	return &db{
		memTable:   skiplist.NewSkipList(),
		tombstones: skiplist.NewSkipList(),
		wal:        log,
	}
}

func (db *db) Get(key leveldb.Key) (leveldb.Value, error) {
	return db.memTable.Search(key)
}

func (db *db) Has(key leveldb.Key) (bool, error) {
	val, err := db.memTable.Search(key)
	if err != nil { // FIXME: slow because of reflection
		if errors.As(err, &notFoundError) {
			return false, nil
		}
		return false, err
	}
	return val != nil, nil
}

func (db *db) Put(key leveldb.Key, value leveldb.Value) error {
	if len(value) == 0 {
		return errors.New("cannot insert blank value")
	}
	if err := db.wal.Put(key, value); err != nil {
		return err
	}
	err := db.memTable.Insert(key, value)
	if err != nil {
		return fmt.Errorf("db.Put: error inserting into memtable: %v", err)
	}
	if err := db.tombstones.Delete(key); err != nil {
		return fmt.Errorf("db.Put: error removing from memtable: %v", err)
	}
	return nil
}

func (db *db) Delete(key leveldb.Key) error {
	if err := db.wal.Delete(key); err != nil {
		return err
	}
	if err := db.memTable.Delete(key); err != nil {
		return fmt.Errorf("db.Delete: error removing from memtable: %v", err)
	}
	if err := db.tombstones.Insert(key, nil); err != nil {
		return fmt.Errorf("db.Delete: error adding to tombstones: %v", err)
	}
	return nil
}

func (db *db) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	precedingNode, err := db.memTable.TraverseUntil(start, nil)
	if err != nil {
		return nil, err
	}
	return NewSkipListIterator(precedingNode, limit), nil
}

func (db *db) flushSSTable(f *os.File) (*sst.SSTableDB, error) {
	sstDb, err := sst.BuildSSTable(f, db.memTable, db.tombstones)
	if err != nil {
		return nil, fmt.Errorf("db.flushSSTable: error building the SSTable: %v", err)
	}

	if err := db.memTable.Reset(); err != nil {
		return nil, fmt.Errorf("db.flushSSTable: error resetting memTable: %v", err)
	}
	if err := db.tombstones.Reset(); err != nil {
		return nil, fmt.Errorf("db.flushSSTable: error resetting tombstones: %v", err)
	}

	return sstDb, nil
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
