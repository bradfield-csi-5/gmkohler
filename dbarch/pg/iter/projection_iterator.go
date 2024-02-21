package iter

import (
	"fmt"
)

func NewProjectionIterator(tuples []*Tuple, columnNames []string) Iterator {
	// Should we be verifying the columnNames exist in the schema?
	return &projectionIterator{
		tuples:      tuples,
		columnNames: columnNames,
		currTuple:   0,
	}
}

type projectionIterator struct {
	tuples      []*Tuple
	columnNames []string
	currTuple   int
}

func (p *projectionIterator) Init() {
	fmt.Println("Init projectionIterator")
}

func (p *projectionIterator) Next() *Tuple {
	if p.currTuple >= len(p.tuples) {
		return nil
	}
	tup := p.tuples[p.currTuple]
	var cols []Column

	for _, colName := range p.columnNames {
		for _, col := range tup.Columns {
			if col.Name == colName {
				cols = append(cols, col)
				break
			}
		}
	}

	p.currTuple++

	return &Tuple{
		Columns: cols,
	}
}

func (p *projectionIterator) Close() {
	fmt.Println("Close projectionIterator")
}

func (p *projectionIterator) Iterators() []Iterator {
	//TODO implement me
	panic("implement me")
}
