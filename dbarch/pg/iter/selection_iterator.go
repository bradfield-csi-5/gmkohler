package iter

type Predicate func(*Tuple) bool

type selectionIterator struct {
	source    Iterator
	predicate Predicate
}

func NewSelectionIterator(source Iterator, predicate Predicate) Iterator {
	return &selectionIterator{
		source:    source,
		predicate: predicate,
	}
}

func (s *selectionIterator) Init() {
	s.source.Init()
}

func (s *selectionIterator) Next() *Tuple {
	for tup := s.source.Next(); tup != nil; tup = s.source.Next() {
		if s.predicate(tup) {
			return tup
		}
	}
	return nil
}

func (s *selectionIterator) Close() {
	s.source.Close()
}
