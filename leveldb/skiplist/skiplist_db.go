package skiplist

import "leveldb"

type skipListDb struct {
	sl SkipList
}

func NewSkipListDb() leveldb.DB {
	return &skipListDb{
		sl: NewSkipList(),
	}
}

func (s *skipListDb) Get(key leveldb.Key) (leveldb.Value, error) {
	return s.sl.Search(key)
}

func (s *skipListDb) Has(key leveldb.Key) (bool, error) {
	val, err := s.sl.Search(key)
	if err != nil { // FIXME: relay "not found" search to false, nil
		return false, err
	}
	return val != nil, nil
}

func (s *skipListDb) Put(key leveldb.Key, value leveldb.Value) error {
	return s.sl.Insert(key, value)
}

func (s *skipListDb) Delete(key leveldb.Key) error {
	return s.sl.Delete(key)
}

func (s *skipListDb) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	precedingNode, err := s.sl.traverseUntil(start)
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
