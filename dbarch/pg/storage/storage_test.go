package storage

import (
	"os"
	"pg/tuple"
	"testing"
)

const filename = "data"

func TestStorage(t *testing.T) {
	var err error
	f, err := os.CreateTemp(os.TempDir(), "data")
	if err != nil {
		t.Fatal("failed to create temp file", err)
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	writer, err := NewFileWriter(f.Name())
	if err != nil {
		t.Fatal("failed to create writer", err)
	}
	defer writer.Close()

	writeTup := tuple.Tuple{Columns: []tuple.Column{
		{
			"id",
			"1",
		},
		{
			"name",
			"Gregory",
		},
	}}
	err = writer.WriteRow(&writeTup)
	if err != nil {
		t.Fatal("failed to write row", err)
	}
	reader, err := NewFileReader(f.Name())
	if err != nil {
		t.Fatal("failed to create reader", err)
	}
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
