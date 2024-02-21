package iter

import "fmt"

type Predicate func(*Tuple) bool

type selectionIterator struct {
	tuples    []*Tuple
	currTuple int
	predicate Predicate
}

func NewSelectionIterator(tuples []*Tuple, predicate Predicate) Iterator {
	return &selectionIterator{
		tuples:    tuples,
		currTuple: 0,
		predicate: predicate,
	}
}

func (s *selectionIterator) Init() {
	fmt.Println("Init SelectionIterator")
}

func (s *selectionIterator) Next() *Tuple {
	var tup *Tuple
	for ; s.currTuple < len(s.tuples) && tup == nil; s.currTuple++ {
		var t = s.tuples[s.currTuple]
		if s.predicate(t) {
			tup = t
		}
	}
	return tup
}

func (s *selectionIterator) Close() {
	fmt.Println("Close SelectionIterator")
}

func (s *selectionIterator) Iterators() []Iterator {
	//TODO implement me
	panic("implement me")
}
