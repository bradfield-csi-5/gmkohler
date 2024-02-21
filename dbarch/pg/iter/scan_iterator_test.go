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
	var results []*Tuple
	var tup *Tuple = si.Next()

	for tup != nil {
		results = append(results, tup)
		tup = si.Next()
	}

	if !slices.Equal(tuples, results) {
		t.Errorf("Expected tuples %+v, got %+v", tuples, results)
	}
}
