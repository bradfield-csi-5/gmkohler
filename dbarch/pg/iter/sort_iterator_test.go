package iter

import (
	"cmp"
	gocmp "github.com/google/go-cmp/cmp"
	"pg/tuple"
	"testing"
)

const (
	teamColumn     = "team"
	positionColumn = "position"
)

func sortByTeamNameAndPosition(tuple1 *tuple.Tuple, tuple2 *tuple.Tuple) int {
	var err error
	team1, err := tuple1.GetColumnValue(teamColumn)
	if err != nil {
		panic(err)
	}
	team2, err := tuple2.GetColumnValue(teamColumn)
	if err != nil {
		panic(err)
	}
	var teamResult = cmp.Compare(team1, team2)
	if teamResult != 0 {
		return teamResult
	}

	pos1, err := tuple1.GetColumnValue(positionColumn)
	if err != nil {
		panic(err)
	}
	pos2, err := tuple2.GetColumnValue(positionColumn)
	if err != nil {
		panic(err)
	}
	return cmp.Compare(pos1, pos2)
}

func TestSortIterator_Next(t *testing.T) {
	var tuples = []*tuple.Tuple{
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
				{"name", "Johnny Bench"},
				{"position", "C"},
				{"team", "Cincinnati Reds"},
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
				{"name", "Johnny Bench"},
				{"position", "C"},
				{"team", "Cincinnati Reds"},
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
				{"name", "Julio Rodríguez"},
				{"position", "CF"},
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
		{
			Columns: []tuple.Column{
				{"name", "Ichiro Suzuki"},
				{"position", "RF"},
				{"team", "Seattle Mariners"},
			},
		},
	}

	si := NewSortIterator(NewScanIterator(tuples), sortByTeamNameAndPosition)
	si.Init()
	var results []*tuple.Tuple

	for tup := si.Next(); tup != nil; tup = si.Next() {
		results = append(results, tup)
	}

	for j, expected := range expectedTuples {
		if !gocmp.Equal(*expected, *results[j]) {
			t.Errorf("expected %+v, got %+v\n", *expected, *results[j])
		}
	}
}
