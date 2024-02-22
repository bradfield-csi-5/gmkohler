package iter

import "pg/tuple"

// Iterator lazily yields a single Tuple each time its Next method is called. Many Iterators will need to maintain
// some state in order to do this. For example, Sort may need to accumulate rows before performing
// a sort, and Limit may need to keep track of how many rows it has returned.
type Iterator interface {
	Init()
	Next() *tuple.Tuple // consider (Tuple, err) for EOF?
	Close()
}
