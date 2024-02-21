package iter

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

const team = "Seattle Mariners"

func isASeattleMariner(tuple *Tuple) bool {
	val, err := tuple.GetColumnValue("team")
	if err != nil {
		return false
	}
	if val == team {
		return true
	}
	return false
}

func TestSelectionIterator_Next(t *testing.T) {
	var tuples = []*Tuple{
		{
			Columns: []Column{
				{"name", "Johnny Bench"},
				{"position", "C"},
				{"team", "Cincinnati Reds"},
			},
		},
		{
			Columns: []Column{
				{"name", "Julio Rodríguez"},
				{"position", "CF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []Column{
				{"name", "Derek Jeter"},
				{"position", "SS"},
				{"team", "New York Yankees"},
			},
		},
		{
			Columns: []Column{
				{"name", "Ichiro Suzuki"},
				{"position", "RF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []Column{
				{"name", "Edgar Martínez"},
				{"position", "DH"},
				{"team", "Seattle Mariners"},
			},
		},
	}
	var expectedTuples = []*Tuple{
		{
			Columns: []Column{
				{"name", "Julio Rodríguez"},
				{"position", "CF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []Column{
				{"name", "Ichiro Suzuki"},
				{"position", "RF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []Column{
				{"name", "Edgar Martínez"},
				{"position", "DH"},
				{"team", "Seattle Mariners"},
			},
		},
	}

	si := NewSelectionIterator(NewScanIterator(tuples), isASeattleMariner)
	si.Init()

	var results []*Tuple
	for tup := si.Next(); tup != nil; tup = si.Next() {
		results = append(results, tup)
	}

	for j, expected := range expectedTuples {
		if !cmp.Equal(*expected, *results[j]) {
			t.Errorf("expected %+v, got %+v\n", *expected, *results[j])
		}
	}
}
