package main

import "pg/iter"

/**
 * Our goal will be to implement an iter-style query executor supporting a few basic queries
 * with selection (filtering rows), projection (filtering columns), and aggregation (such as COUNT
 * or SUM).
 *
 * At this point, we recommend using dummy data in memory, rather than attempting to read from a
 * database file on disk. Aim at least to implement a Scan node that yields a single record each
 * time its Next method is called, as well as a Selection node initialized with a predicate
 * function (one which returns true or false) which yields the next record for which the predicate
 * function returns true whenever its own next method is called.
 */
func main() {

}

type Executor struct {
	root *iter.Iterator
}
type instructionType int

const (
	scan instructionType = iota
	limit
	projection
	selection
	sort
)

type instruction struct {
	iType instructionType
	vals  []string
}

func (e *Executor) Execute() []string { return nil }

func NewExecutor(root *iter.Iterator) *Executor {
	return &Executor{root: root}
}
