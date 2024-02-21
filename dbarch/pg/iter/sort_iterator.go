package iter

import "slices"

type SortFunc func(*Tuple, *Tuple) int

type sortIterator struct {
	tuples    []*Tuple
	sortFunc  SortFunc
	currTuple int
}

func NewSortIterator(tuples []*Tuple, sortFunc SortFunc) Iterator {
	return &sortIterator{
		tuples:   tuples,
		sortFunc: sortFunc,
	}

}

func (s *sortIterator) Init() {
	// TODO make this lazier (heap?)
	slices.SortFunc(s.tuples, s.sortFunc)
}

func (s *sortIterator) Next() *Tuple {
	if s.currTuple >= len(s.tuples) {
		return nil
	}
	var tup = s.tuples[s.currTuple]
	s.currTuple++
	return tup
}

func (s *sortIterator) Close() {
	//TODO implement me
	panic("implement me")
}

func (s *sortIterator) Iterators() []Iterator {
	//TODO implement me
	panic("implement me")
}
