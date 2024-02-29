package tuple

import "errors"

type Tuple struct {
	Columns []Column
}

func (t Tuple) GetColumnValue(name string) (ColumnValue, error) {
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

type ColumnValue string

// Column assumes all values are strings (for now)
type Column struct {
	Name  string
	Value ColumnValue
}
