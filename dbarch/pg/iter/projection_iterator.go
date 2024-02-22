package iter

import (
	"fmt"
	"pg/tuple"
)

func NewProjectionIterator(source Iterator, columnNames []string) Iterator {
	// Should we be verifying the columnNames exist in the schema?
	return &projectionIterator{
		source:      source,
		columnNames: columnNames,
	}
}

type projectionIterator struct {
	source      Iterator
	columnNames []string
}

func (p *projectionIterator) Init() {
	fmt.Println("Init projectionIterator")
}

func (p *projectionIterator) Next() *tuple.Tuple {
	var tup = p.source.Next()
	if tup == nil {
		return nil
	}

	var cols []tuple.Column

	for _, colName := range p.columnNames {
		for _, col := range tup.Columns {
			if col.Name == colName {
				cols = append(cols, col)
				break
			}
		}
	}

	tup.Columns = cols
	return tup
}

func (p *projectionIterator) Close() {
	p.source.Close()
}
