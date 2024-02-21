package iter

import "fmt"

// limitIterator limits the number of output rows.
type limitIterator struct {
	tuples []*Tuple
	limit  int
	sent   int
}

func NewLimitIterator(tuples []*Tuple, limit int) Iterator {
	// TODO: Consider (Iterator, err) return signature, and fail when limit < 0
	return &limitIterator{
		tuples: tuples,
		limit:  limit,
		sent:   0,
	}
}

func (l *limitIterator) Init() {
	fmt.Println("Init limitIterator")
}

func (l *limitIterator) Next() *Tuple {
	if l.sent >= l.limit || l.sent >= len(l.tuples) {
		return nil
	}

	tup := l.tuples[l.sent]
	l.sent++
	return tup
}

func (l *limitIterator) Close() {
	fmt.Println("Close limitIterator")
}

func (l *limitIterator) Iterators() []Iterator {
	//TODO implement me
	panic("implement me")
}
