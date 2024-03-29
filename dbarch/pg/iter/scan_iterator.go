package iter

import (
	"fmt"
	"pg/tuple"
)

// scanIterator yields each row for the table as needed. In this initial implementation scanIterator returns rows from
// a predefined list in memory.
type scanIterator struct {
	tuples    []*tuple.Tuple
	currTuple int
}

func NewScanIterator(tuples []*tuple.Tuple) Iterator {
	return &scanIterator{
		tuples:    tuples,
		currTuple: 0,
	}
}

func (s *scanIterator) Init() {
	fmt.Println("Init scanIterator")
}

func (s *scanIterator) Next() *tuple.Tuple {
	if s.currTuple >= len(s.tuples) {
		return nil
	}

	tup := s.tuples[s.currTuple]
	s.currTuple++
	return tup
}

func (s *scanIterator) Close() {
	fmt.Println("Close scanIterator")
}
