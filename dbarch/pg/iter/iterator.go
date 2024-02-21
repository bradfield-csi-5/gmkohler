package iter

import "errors"

// Iterator lazily yields a single Tuple each time its Next method is called. Many Iterators will need to maintain
// some state in order to do this. For example, Sort may need to accumulate rows before performing
// a sort, and Limit may need to keep track of how many rows it has returned.
type Iterator interface {
	Init()
	Next() *Tuple // consider (Tuple, err) for EOF?
	Close()
}

type Tuple struct {
	Columns []Column
}

func (t Tuple) GetColumnValue(name string) (string, error) {
	var col *Column
	for _, c := range t.Columns {
		if c.Name == name {
			col = &c
			break
		}
	}

	if col == nil {
		return "", errors.New("column does not exist")
	}

	return col.Value, nil
}

// Column assumes all values are strings (for now)
type Column struct {
	Name  string
	Value string
}
