package sst

import (
	"bytes"
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
	var sl = skiplist.NewSkipList()
	if err := sl.Insert(leveldb.Key("foo"), leveldb.Value("bar")); err != nil {
		t.Fatalf("error calling DB.Put(): %v", err)
	}
	sstDb, err := BuildSSTable(file, sl, skiplist.NewSkipList())
	if err != nil {
		t.Fatal("error building SSTable:", err)
	}
	value, err := sstDb.Get(leveldb.Key("foo"))
	if err != nil {
		t.Fatalf("unexpected error calling sstDb.Get(): %v", err)
	}
	if bytes.Compare(leveldb.Value("bar"), value) != 0 {
		t.Errorf("unexpected returned value.  Expected %q, got %q", "bar", value)
	}
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
