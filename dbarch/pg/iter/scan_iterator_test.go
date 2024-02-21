package iter

import (
	"slices"
	"testing"
)

func TestScanIterator_Next(t *testing.T) {
	var tuples = []*Tuple{
		{
			Columns: []Column{
				{Name: "age", Value: "24"},
				{Name: "name", Value: "Mary Contrary"},
			},
		},
		{
			Columns: []Column{
				{Name: "age", Value: "22"},
				{Name: "name", Value: "Bob Snob"},
			},
		},
		{
			Columns: []Column{
				{Name: "age", Value: "30"},
				{Name: "name", Value: "Julia Goulia"},
			},
		},
	}

	var si Iterator = NewScanIterator(tuples)
	si.Init()
	var results []*Tuple

	for tup := si.Next(); tup != nil; tup = si.Next() {
		results = append(results, tup)
	}

	if !slices.Equal(tuples, results) {
		t.Errorf("Expected tuples %+v, got %+v", tuples, results)
	}
}
