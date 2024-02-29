package iter

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"pg/storage"
	"pg/tuple"
	"testing"
)

func TestFileIterator_Next(t *testing.T) {
	tuples := []*tuple.Tuple{
		{
			Columns: []tuple.Column{
				{
					Name:  "name",
					Value: "Bleach",
				},
				{
					Name:  "year",
					Value: "1989",
				},
				{
					Name:  "duration",
					Value: "37:21",
				},
			},
		},
		{
			Columns: []tuple.Column{
				{
					Name:  "name",
					Value: "Nevermind",
				},
				{
					Name:  "year",
					Value: "1991",
				},
				{
					Name:  "duration",
					Value: "42:36",
				},
			},
		},
		{
			Columns: []tuple.Column{
				{
					Name:  "name",
					Value: "In Utero",
				},
				{
					Name:  "year",
					Value: "1993",
				},
				{
					Name:  "duration",
					Value: "41:23",
				},
			},
		},
	}
	f, err := os.CreateTemp(os.TempDir(), "data")
	if err != nil {
		t.Fatal("failed to create file", err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	var colNames []string
	for _, col := range tuples[0].Columns {
		colNames = append(colNames, col.Name)
	}
	writer := storage.NewFileWriter(colNames, len(tuples), f)
	defer func() {
		if err := writer.Close(); err != nil {
			t.Errorf("error closing writer: %v", err)
		}
	}()

	for _, tup := range tuples {
		err := writer.WriteRow(tup)
		if err != nil {
			t.Fatalf("failed to write row %+v: %v\n", *tup, err)
		}
	}
	readFile, err := os.Open(f.Name())
	if err != nil {
		t.Fatal("failed to open file", err)
	}
	reader := storage.NewFileReader(readFile)
	fileIter := NewFileIterator(reader)
	var results []*tuple.Tuple
	for tup := fileIter.Next(); tup != nil; tup = fileIter.Next() {
		results = append(results, tup)
	}
	for j, tup := range results {
		original := tuples[j]
		if !cmp.Equal(*original, *tup) {
			t.Fatalf("expected tuple %+v, got %+v\n", *original, *tup)
		}
	}
}
