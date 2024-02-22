package iter

import (
	"pg/expr"
	"pg/tuple"
)

type Predicate func(*tuple.Tuple) bool

type selectionIterator struct {
	source    Iterator
	predicate expr.Expression
}

func NewSelectionIterator(source Iterator, predicate expr.Expression) Iterator {
	return &selectionIterator{
		source:    source,
		predicate: predicate,
	}
}

func (s *selectionIterator) Init() {
	s.source.Init()
}

func (s *selectionIterator) Next() *tuple.Tuple {
	for tup := s.source.Next(); tup != nil; tup = s.source.Next() {
		if s.predicate.Execute(*tup) {
			return tup
		}
	}
	return nil
}

func (s *selectionIterator) Close() {
	s.source.Close()
}
