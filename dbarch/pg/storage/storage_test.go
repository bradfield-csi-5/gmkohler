package storage

import (
	"os"
	"pg/tuple"
	"testing"
)

const filename = "data"

func TestStorage(t *testing.T) {
	writeTups := []tuple.Tuple{
		{
			Columns: []tuple.Column{
				{
					"id",
					"1",
				},
				{
					"name",
					"Larry David",
				},
			},
		},
		{
			Columns: []tuple.Column{
				{
					"id",
					"2",
				},
				{
					"name",
					"Richard Lewis",
				},
			},
		},
	}
	var colNames []string
	for _, col := range writeTups[0].Columns {
		colNames = append(colNames, col.Name)
	}
	var err error
	f, err := os.CreateTemp(os.TempDir(), "data")
	if err != nil {
		t.Fatal("failed to create temp file", err)
	}
	defer func() {
		os.Remove(f.Name())
	}()
	writer := NewFileWriter(colNames, len(writeTups), f)
	if err != nil {
		t.Fatal("failed to create writer", err)
	}

	for _, writeTup := range writeTups {
		err = writer.WriteRow(&writeTup)
		if err != nil {
			t.Fatal("failed to write row", err)
		}
	}
	writer.Close()
	readFile, err := os.Open(f.Name())
	if err != nil {
		t.Fatal("failed to open file for reading", err)
	}
	defer readFile.Close()
	reader := NewFileReader(readFile)
	if err != nil {
		t.Fatal("failed to create reader", err)
	}
	defer reader.Close()
	for _, writeTup := range writeTups {
		readTup, err := reader.ReadRow()
		if err != nil {
			t.Fatal("failed to read row", err)
		}
		writeCols := writeTup.Columns
		for j, readCol := range readTup.Columns {
			if writeCols[j] != readCol {
				t.Errorf("expected column to match %+v, got %+v\n", writeCols[j], readCol)
			}
		}
	}
}
