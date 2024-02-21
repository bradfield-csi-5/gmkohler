package iter

import (
	"fmt"
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

func (p *projectionIterator) Next() *Tuple {
	var tup = p.source.Next()
	if tup == nil {
		return nil
	}

	var cols []Column

	for _, colName := range p.columnNames {
		for _, col := range tup.Columns {
			if col.Name == colName {
				cols = append(cols, col)
				break
			}
		}
	}

	return &Tuple{
		Columns: cols,
	}
}

func (p *projectionIterator) Close() {
	p.source.Close()
}
