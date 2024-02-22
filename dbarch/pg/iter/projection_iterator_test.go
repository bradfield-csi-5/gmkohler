package iter

import (
	"github.com/google/go-cmp/cmp"
	"pg/tuple"
	"testing"
)

func TestProjectionIterator_Next(t *testing.T) {
	tuples := []*tuple.Tuple{
		{
			Columns: []tuple.Column{
				{"name", "Ada"},
				{"gender", "F"},
				{"department", "computer_science"},
				{"year", "freshman"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Malcolm"},
				{"gender", "M"},
				{"department", "sociology"},
				{"year", "sophomore"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Richard"},
				{"gender", "M"},
				{"department", "physics"},
				{"year", "junior"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Marie"},
				{"gender", "F"},
				{"department", "chemistry"},
				{"year", "senior"},
			},
		},
	}
	expectedTuples := []*tuple.Tuple{
		{
			Columns: []tuple.Column{
				{"department", "computer_science"},
				{"name", "Ada"},
			},
		},
		{
			Columns: []tuple.Column{
				{"department", "sociology"},
				{"name", "Malcolm"},
			},
		},
		{
			Columns: []tuple.Column{
				{"department", "physics"},
				{"name", "Richard"},
			},
		},
		{
			Columns: []tuple.Column{
				{"department", "chemistry"},
				{"name", "Marie"},
			},
		},
	}
	// tests rearranging columns
	pi := NewProjectionIterator(NewScanIterator(tuples), []string{"department", "name"})
	var results []*tuple.Tuple
	for tup := pi.Next(); tup != nil; tup = pi.Next() {
		results = append(results, tup)
	}

	for j, expectedTup := range expectedTuples {
		if !cmp.Equal(*expectedTup, *results[j]) {
			t.Errorf("expected %+v, got %+v\n", *expectedTup, *results[j])
		}
	}
}
