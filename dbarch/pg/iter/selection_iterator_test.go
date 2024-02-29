package iter

import (
	"github.com/google/go-cmp/cmp"
	"pg/expr"
	"pg/tuple"
	"testing"
)

const mariners = "Seattle Mariners"

var isASeattleMariner = expr.NewEqualityExpression("team", mariners)

func TestSelectionIterator_Next(t *testing.T) {
	var tuples = []*tuple.Tuple{
		{
			Columns: []tuple.Column{
				{"name", "Johnny Bench"},
				{"position", "C"},
				{"team", "Cincinnati Reds"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Julio Rodríguez"},
				{"position", "CF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Derek Jeter"},
				{"position", "SS"},
				{"team", "New York Yankees"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Ichiro Suzuki"},
				{"position", "RF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Edgar Martínez"},
				{"position", "DH"},
				{"team", "Seattle Mariners"},
			},
		},
	}
	var expectedTuples = []*tuple.Tuple{
		{
			Columns: []tuple.Column{
				{"name", "Julio Rodríguez"},
				{"position", "CF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Ichiro Suzuki"},
				{"position", "RF"},
				{"team", "Seattle Mariners"},
			},
		},
		{
			Columns: []tuple.Column{
				{"name", "Edgar Martínez"},
				{"position", "DH"},
				{"team", "Seattle Mariners"},
			},
		},
	}

	si := NewSelectionIterator(NewScanIterator(tuples), isASeattleMariner)
	si.Init()

	var results []*tuple.Tuple
	for tup := si.Next(); tup != nil; tup = si.Next() {
		results = append(results, tup)
	}

	for j, expected := range expectedTuples {
		if !cmp.Equal(*expected, *results[j]) {
			t.Errorf("expected %+v, got %+v\n", *expected, *results[j])
		}
	}
}
