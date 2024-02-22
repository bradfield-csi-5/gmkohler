package iter

import (
	"pg/tuple"
	"slices"
)

type SortFunc func(*tuple.Tuple, *tuple.Tuple) int

type sortIterator struct {
	sortFunc  SortFunc
	currTuple int
	source    Iterator
	sorted    []*tuple.Tuple
}

func NewSortIterator(source Iterator, sortFunc SortFunc) Iterator {
	return &sortIterator{
		source:   source,
		sortFunc: sortFunc,
	}
}

func (s *sortIterator) Init() {
	s.source.Init()
	// TODO make this lazier (heap?)
	var tuples []*tuple.Tuple
	for tup := s.source.Next(); tup != nil; tup = s.source.Next() {
		tuples = append(tuples, tup)
	}
	slices.SortFunc(tuples, s.sortFunc)
	s.sorted = tuples
}

func (s *sortIterator) Next() *tuple.Tuple {
	if s.currTuple >= len(s.sorted) {
		return nil
	}
	var tup = s.sorted[s.currTuple]
	s.currTuple++
	return tup
}

func (s *sortIterator) Close() {
	s.source.Close()
}
