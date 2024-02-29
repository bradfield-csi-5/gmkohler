package iter

import (
	"pg/tuple"
	"slices"
	"testing"
)

func TestScanIterator_Next(t *testing.T) {
	var tuples = []*tuple.Tuple{
		{
			Columns: []tuple.Column{
				{Name: "age", Value: "24"},
				{Name: "name", Value: "Mary Contrary"},
			},
		},
		{
			Columns: []tuple.Column{
				{Name: "age", Value: "22"},
				{Name: "name", Value: "Bob Snob"},
			},
		},
		{
			Columns: []tuple.Column{
				{Name: "age", Value: "30"},
				{Name: "name", Value: "Julia Goulia"},
			},
		},
	}

	var si Iterator = NewScanIterator(tuples)
	si.Init()
	var results []*tuple.Tuple

	for tup := si.Next(); tup != nil; tup = si.Next() {
		results = append(results, tup)
	}

	if !slices.Equal(tuples, results) {
		t.Errorf("Expected tuples %+v, got %+v", tuples, results)
	}
}
