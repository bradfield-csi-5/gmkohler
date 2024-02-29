package iter

import (
	"pg/tuple"
	"slices"
	"testing"
)

type limitDatum struct {
	specified       int
	expectedResults int
}

func TestLimitIterator_Next(t *testing.T) {
	var limits = []limitDatum{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 3},
		{4, 3}, // can't return more than the number of tuples we have
	}
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

	for _, lim := range limits {
		var li = NewLimitIterator(NewScanIterator(tuples), lim.specified)

		var results []*tuple.Tuple
		for tup := li.Next(); tup != nil; tup = li.Next() {
			results = append(results, tup)
		}

		if len(results) != lim.expectedResults {
			t.Errorf("Exepected %d tuples to be returned, got %d", lim.expectedResults, len(results))
		}
		if !slices.Equal(tuples[:lim.expectedResults], results) {
			t.Errorf("Expected %+v tuples to be returned, got %+v", results, tuples[:lim.expectedResults])
		}
	}
}
