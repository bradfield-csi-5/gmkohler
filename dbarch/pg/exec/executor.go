package exec

import (
	"pg/iter"
	"pg/tuple"
)

type Executor struct {
	root iter.Iterator
}

func (e *Executor) Execute() [][]tuple.ColumnValue {
	e.root.Init()
	var results [][]tuple.ColumnValue
	for tup := e.root.Next(); tup != nil; tup = e.root.Next() {
		var result []tuple.ColumnValue
		for _, col := range tup.Columns {
			result = append(result, col.Value)
		}
		results = append(results, result)
	}
	return results
}

func NewExecutor(root iter.Iterator) *Executor {
	return &Executor{root: root}
}
