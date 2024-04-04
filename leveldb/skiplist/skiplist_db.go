package skiplist

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"leveldb"
	"leveldb/wal"
)

type skipListDb struct {
	sl  SkipList
	wal *wal.Log
}

func NewSkipListDbFromWal(rw io.ReadWriter) (leveldb.DB, error) {
	reader := bufio.NewReader(rw)
	entries, err := wal.DecodeLogFile(reader)
	if err != nil {
		return nil, err
	}
	db := NewSkipListDb(rw)
	for _, entry := range entries {
		switch entry.Operation {
		case wal.OpPut:
			if err := db.Put(entry.Key, entry.Value); err != nil {
				return nil, err
			}
		case wal.OpDelete:
			if err := db.Delete(entry.Key); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unrecognized opcode %s", entry.Operation)
		}
	}
	return db, nil
}
func NewSkipListDb(readWriter io.ReadWriter) leveldb.DB {
	var log *wal.Log
	if readWriter != nil { // hacky, think of nicer way
		log = wal.NewLog(readWriter)
	}

	return &skipListDb{
		sl:  NewSkipList(),
		wal: log,
	}
}

func (db *skipListDb) Get(key leveldb.Key) (leveldb.Value, error) {
	return db.sl.Search(key)
}

func (db *skipListDb) Has(key leveldb.Key) (bool, error) {
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

func (db *skipListDb) Put(key leveldb.Key, value leveldb.Value) error {
	if err := db.wal.Put(key, value); err != nil {
		return err
	}
	return db.sl.Insert(key, value)
}

func (db *skipListDb) Delete(key leveldb.Key) error {
	if err := db.wal.Delete(key); err != nil {
		return err
	}
	return db.sl.Delete(key)
}

func (db *skipListDb) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	precedingNode, err := db.sl.traverseUntil(start, nil)
	if err != nil {
		return nil, err
	}
	return NewSkipListIterator(precedingNode, limit), nil
}

// consider converting to memory instead of this
type skipListIterator struct {
	limit   leveldb.Key
	current skipListNode
	err     error
}

// NewSkipListIterator accepts a "preceding node" i.e. one whose level-1 "forward node" is in the desired range,
// and a limit indicating when to stop traversing the skip list.
func NewSkipListIterator(precedingNode skipListNode, limit leveldb.Key) leveldb.Iterator {
	return &skipListIterator{
		limit:   limit,
		current: precedingNode,
	}
}

func (s *skipListIterator) Next() bool {
	if nextNode := s.current.ForwardNodeAtLevel(1); nextNode.CompareKey(s.limit) <= 0 {
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
