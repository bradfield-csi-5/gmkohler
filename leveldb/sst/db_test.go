package sst

import (
	"bytes"
	"errors"
	"leveldb"
	"leveldb/skiplist"
	"os"
	"testing"
)

func TestBuildSSTable(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), "sst")
	if err != nil {
		t.Fatal("failed to create SST file:", err)
	}

	var memTable = skiplist.NewSkipList()
	if err := memTable.Insert(leveldb.Key("foo"), leveldb.Value("bar")); err != nil {
		t.Fatalf("error inserting key into memTable skiplist: %v", err)
	}
	var tombstones = skiplist.NewSkipList()
	if err := tombstones.Insert(leveldb.Key("spam"), nil); err != nil {
		t.Fatalf("error inserting into tombstone skiplist: %v", err)
	}

	sstDb, err := BuildSSTable(file, memTable, tombstones)
	if err != nil {
		t.Fatal("error building SSTable:", err)
	}

	t.Run("GetExists", func(t *testing.T) {
		value, err := sstDb.Get(leveldb.Key("foo"))
		if err != nil {
			t.Fatalf("unexpected error calling sstDb.Get(): %v", err)
		}
		if bytes.Compare(leveldb.Value("bar"), value) != 0 {
			t.Errorf("unexpected returned value.  Expected %q, got %q", "bar", value)
		}
	})

	t.Run("GetTombstoned", func(t *testing.T) {
		_, err := sstDb.Get(leveldb.Key("spam"))
		if err == nil {
			t.Error("expected error calling sstDb.Get() for tombstoned value, did not get one")
		}
		var notFoundError *leveldb.NotFoundError
		if !errors.As(err, &notFoundError) {
			t.Errorf("expected a NotFoundError, got %T: %v", err, err)
		}
	})
}

func TestNewSSTableDB(t *testing.T) {
	t.Skip("not yet implemented")
	//file, err := os.CreateTemp(os.TempDir(), "sst")
	//if err != nil {
	//	t.Fatal("failed to create SST file:", err)
	//}
	//
	//var sstDb leveldb.ReadOnlyDB
	//sstDb, err = NewSSTableDBFromFile(file)
	//fmt.Println(sstDb)
	//t.Fatal("finish test")
}
