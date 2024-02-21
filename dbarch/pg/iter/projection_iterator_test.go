package iter

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestProjectionIterator_Next(t *testing.T) {
	tuples := []*Tuple{
		{
			Columns: []Column{
				{"name", "Ada"},
				{"gender", "F"},
				{"department", "computer_science"},
				{"year", "freshman"},
			},
		},
		{
			Columns: []Column{
				{"name", "Malcolm"},
				{"gender", "M"},
				{"department", "sociology"},
				{"year", "sophomore"},
			},
		},
		{
			Columns: []Column{
				{"name", "Richard"},
				{"gender", "M"},
				{"department", "physics"},
				{"year", "junior"},
			},
		},
		{
			Columns: []Column{
				{"name", "Marie"},
				{"gender", "F"},
				{"department", "chemistry"},
				{"year", "senior"},
			},
		},
	}
	expectedTuples := []*Tuple{
		{
			Columns: []Column{
				{"department", "computer_science"},
				{"name", "Ada"},
			},
		},
		{
			Columns: []Column{
				{"department", "sociology"},
				{"name", "Malcolm"},
			},
		},
		{
			Columns: []Column{
				{"department", "physics"},
				{"name", "Richard"},
			},
		},
		{
			Columns: []Column{
				{"department", "chemistry"},
				{"name", "Marie"},
			},
		},
	}
	// tests rearranging columns
	pi := NewProjectionIterator(tuples, []string{"department", "name"})
	var results []*Tuple
	for tup := pi.Next(); tup != nil; tup = pi.Next() {
		results = append(results, tup)
	}

	for j, expectedTup := range expectedTuples {
		if !cmp.Equal(*expectedTup, *results[j]) {
			t.Errorf("expected %+v, got %+v\n", *expectedTup, *results[j])
		}
	}
}
