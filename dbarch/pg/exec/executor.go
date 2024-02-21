package exec

import "pg/iter"

type Executor struct {
	root iter.Iterator
}

func (e *Executor) Execute() [][]iter.ColumnValue {
	e.root.Init()
	var results [][]iter.ColumnValue
	for tup := e.root.Next(); tup != nil; tup = e.root.Next() {
		var result []iter.ColumnValue
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
