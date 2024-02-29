package exec

import (
	"github.com/gocarina/gocsv"
	"os"
	"pg/expr"
	"pg/iter"
	"pg/tuple"
	"slices"
	"strconv"
	"strings"
	"testing"
)

const (
	movieId = "movie_id"
)

type genres []string

func (g *genres) UnmarshalCSV(s string) error {
	*g = strings.Split(s, "|")
	return nil
}

type movie struct {
	MovieId int    `csv:"movieId"`
	Title   string `csv:"title"`
	Genres  genres `csv:"genres"`
}

var tuples []*tuple.Tuple

func tupleFromMovie(m movie) *tuple.Tuple {
	return &tuple.Tuple{
		Columns: []tuple.Column{
			{Name: movieId, Value: tuple.ColumnValue(strconv.Itoa(m.MovieId))},
			{Name: "title", Value: tuple.ColumnValue(m.Title)},
			{Name: "genres", Value: tuple.ColumnValue(strings.Join(m.Genres, "|"))},
		}}
}

func init() {
	moviesFile, err := os.Open("movies.csv")
	if err != nil {
		panic(err)
	}
	defer moviesFile.Close()
	var movies []*movie
	if err := gocsv.UnmarshalFile(moviesFile, &movies); err != nil {
		panic(err)
	}
	for _, m := range movies {
		tuples = append(tuples, tupleFromMovie(*m))
	}
}

func TestExecutor_Execute(t *testing.T) {
	iterator := iter.NewProjectionIterator(
		iter.NewSelectionIterator(
			iter.NewScanIterator(tuples),
			expr.NewEqualityExpression(movieId, "5000"),
		),
		[]string{"title"},
	)
	var executor = NewExecutor(iterator)
	var results = executor.Execute()
	var expected = [][]tuple.ColumnValue{
		{
			"Medium Cool (1969)",
		},
	}
	for j, exp := range expected {
		if !slices.Equal(results[j], exp) {
			t.Errorf("expected %v, got %v\n", exp, results[j])
		}
	}
}
