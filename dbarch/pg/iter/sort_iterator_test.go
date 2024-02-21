package iter

import (
	"cmp"
	gocmp "github.com/google/go-cmp/cmp"
	"testing"
)

const (
	teamColumn     = "team"
	positionColumn = "position"
)

func sortByTeamNameAndPosition(tuple1 *Tuple, tuple2 *Tuple) int {
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
	var tuples = []*Tuple{
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
				{"name", "Johnny Bench"},
				{"position", "C"},
				{"team", "Cincinnati Reds"},
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
				{"name", "Johnny Bench"},
				{"position", "C"},
				{"team", "Cincinnati Reds"},
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
				{"name", "Julio Rodríguez"},
				{"position", "CF"},
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
		{
			Columns: []Column{
				{"name", "Ichiro Suzuki"},
				{"position", "RF"},
				{"team", "Seattle Mariners"},
			},
		},
	}

	si := NewSortIterator(tuples, sortByTeamNameAndPosition)
	si.Init()
	var results []*Tuple

	for tup := si.Next(); tup != nil; tup = si.Next() {
		results = append(results, tup)
	}

	for j, expected := range expectedTuples {
		if !gocmp.Equal(*expected, *results[j]) {
			t.Errorf("expected %+v, got %+v\n", *expected, *results[j])
		}
	}
}
