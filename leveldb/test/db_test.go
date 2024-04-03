package test

import (
	"bytes"
	"leveldb"
	"leveldb/inmem"
	"leveldb/skiplist"
	"testing"
)

var testImpls = []TestSetup{
	{
		Name:    "InMemory",
		NewDb:   inmem.NewInMemoryDb,
		EmptyDb: func() leveldb.DB { return inmem.NewInMemoryDb(nil) },
	},
	{
		Name: "SkipList",
		NewDb: func(entries []leveldb.DataEntry) leveldb.DB {
			sl := skiplist.NewSkipListDb()
			for _, entry := range entries {
				_ = sl.Put(entry.Key, entry.Value)
			}
			return sl
		},
		EmptyDb: skiplist.NewSkipListDb,
	},
}

// consider just importing the impls here and avoiding duplicating tests
type TestSetup struct {
	Name    string
	NewDb   func(entries []leveldb.DataEntry) leveldb.DB
	EmptyDb func() leveldb.DB
}

var data = []leveldb.DataEntry{
	{
		Key:   leveldb.Key("eggs"),
		Value: leveldb.Value("scrambled"),
	},
	{
		Key:   leveldb.Key("spam"),
		Value: leveldb.Value("ham"),
	},
}

func TestDb_Get_NoEntry(t *testing.T) {
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			db := impl.NewDb(data)
			val, err := db.Get(leveldb.Key("foo"))
			if err == nil || val != nil {
				t.Errorf("expected error when calling db.Get() for non-existent entry")
			}
		})
	}
}

func TestDb_Get_EntryExists(t *testing.T) {
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			var db = impl.NewDb(data)
			val, err := db.Get(leveldb.Key("eggs"))
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			var readableVal = string(val)
			if readableVal != "scrambled" {
				t.Fatalf("expected value for \"eggs\" to be \"scrambled\", got %q", readableVal)
			}
		})
	}
}

func TestDb_Has_True(t *testing.T) {
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			var db = impl.NewDb(data)
			var keyExists, err = db.Has(leveldb.Key("eggs"))
			if err != nil {
				t.Fatalf("%s: unexpected error", err)
			}
			if !keyExists {
				t.Fatal("expected key to exist, return value was false")
			}
		})
	}
}

func TestDb_Has_False(t *testing.T) {
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			var db = impl.NewDb(data)
			var keyExists, err = db.Has(leveldb.Key("foo"))
			if err != nil {
				t.Fatal("unexpected error:", err)
			}
			if keyExists {
				t.Fatal("expected key not to exist, return value was true")
			}
		})
	}
}

func TestDb_Put_NewEntry(t *testing.T) {
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			db := impl.NewDb(data)
			err := db.Put(leveldb.Key("foo"), leveldb.Value("bar"))
			if err != nil {
				t.Fatal("unexpected error when putting value", err)
			}
			val, err := db.Get(leveldb.Key("foo"))
			if err != nil {
				t.Fatal("unexpected error when getting value", err)
			}
			stringVal := string(val)
			if stringVal != "bar" {
				t.Fatalf("expected value 'bar', got %q", stringVal)
			}
		})
	}
}

func TestDb_Put_UpdateEntry(t *testing.T) {
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			db := impl.NewDb(data)
			err := db.Put(leveldb.Key("eggs"), leveldb.Value("poached"))
			if err != nil {
				t.Fatal("unexpected error when putting value", err)
			}
			val, err := db.Get(leveldb.Key("eggs"))
			if err != nil {
				t.Fatal("unexpected error when getting value", err)
			}
			stringVal := string(val)
			if stringVal != "poached" {
				t.Fatalf("expected value \"poached\", got %q", stringVal)
			}
		})
	}
}

func TestDb_Delete_Success(t *testing.T) {
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			db := impl.NewDb(data)
			err := db.Delete(leveldb.Key("eggs"))
			if err != nil {
				t.Fatal("unexpected error while deleting key", err)
			}
			_, err = db.Get(leveldb.Key("eggs"))
			if err == nil {
				t.Fatal("expected an error getting deleted key, but did not get one")
			}
			err = db.Delete(leveldb.Key("eggs"))
			if err == nil {
				t.Fatal("expected an error deleting non-existent key, but did not get one")
			}
		})
	}
}

func TestDb_RangeScan(t *testing.T) {
	scanData := []leveldb.DataEntry{
		{leveldb.Key("abc"), leveldb.Value("ABC")},
		{leveldb.Key("abd"), leveldb.Value("ABD")},
		{leveldb.Key("abe"), leveldb.Value("ABE")},
		{leveldb.Key("abf"), leveldb.Value("ABF")},
		{leveldb.Key("abg"), leveldb.Value("ABG")},
	}
	for _, impl := range testImpls {
		t.Run(impl.Name, func(t *testing.T) {
			db := impl.NewDb(data)
			var err error
			for _, datum := range scanData {
				err = db.Put(datum.Key, datum.Value)
				if err != nil {
					t.Fatal("unexpected error executing Put()", err)
				}
			}
			results, err := db.RangeScan(leveldb.Key("abc"), leveldb.Key("abf"))
			if err != nil {
				t.Fatal("unexpected error executing RangeScan()", err)
			}
			expectedResults := scanData[0:4] // `results` should include matches to `limit` parameter.
			for j, datum := range expectedResults {
				hasNext := results.Next()
				if !hasNext {
					t.Fatalf("expected more results, got %d", j+1)
				}
				if !bytes.Equal(datum.Key, results.Key()) {
					t.Fatalf("expected key %q, got %q", datum.Key, results.Key())
				}
				if !bytes.Equal(datum.Value, results.Value()) {
					t.Fatalf("expected value %q, got %q", datum.Value, results.Value())
				}
			}
			if results.Error() != nil {
				t.Fatal("iterator generated unexpected error", err)
			}
			if results.Next() {
				t.Fatal("got more results than expected")
			}
		})
	}
}
