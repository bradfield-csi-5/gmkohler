package test

import (
	"bytes"
	"leveldb"
	"leveldb/skiplist"
	"leveldb/wal"
	"os"
	"testing"
)

func TestRecovery(t *testing.T) {
	writeFile, err := os.CreateTemp(os.TempDir(), "wal")
	if err != nil {
		t.Fatal("failed to create temp WAL file:", err)
	}
	log := wal.NewLog(writeFile)
	for _, entry := range []leveldb.DataEntry{
		{leveldb.Key("genre"), leveldb.Value("ambient")},
		{leveldb.Key("artist"), leveldb.Value("Khotin")},
		{leveldb.Key("artist"), leveldb.Value("Eno")},
	} {
		if err := log.Put(entry.Key, entry.Value); err != nil {
			t.Fatal("failed to put entry in log:", err)
		}
	}
	if err := writeFile.Close(); err != nil {
		t.Error("failed to close file descriptor used for writing WAL:", err)
	}
	walFile, err := os.OpenFile(writeFile.Name(), os.O_RDWR, os.ModeType)
	if err != nil {
		t.Fatal("error opening WAL file:", err)
	}
	db, err := skiplist.NewSkipListDbFromWal(walFile)
	if err != nil {
		t.Fatal("error initializing DB from WAL:", err)
	}
	val, err := db.Get(leveldb.Key("artist"))
	if err != nil {
		t.Fatal("error getting value from database:", err)
	}
	if !bytes.Equal(val, leveldb.Value("Eno")) {
		t.Errorf("expected value to be %q, got %q", "Eno", val)
	}
}
