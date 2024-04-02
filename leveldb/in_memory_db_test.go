package leveldb

import (
	"bytes"
	"testing"
)

func populatedDb() DB {
	return &inMemoryDb{data: []DataEntry{
		{
			Key:   Key("eggs"),
			Value: Value("scrambled"),
		},
		{
			Key:   Key("spam"),
			Value: Value("ham"),
		},
	}}
}

func emptyDb() DB { return &inMemoryDb{data: make([]DataEntry, 0)} }

func TestInMemoryDb_Get_NoEntry(t *testing.T) {
	var db = emptyDb()
	val, err := db.Get(Key("foo"))
	if err == nil || val != nil {
		t.Errorf("expected error when calling db.Get() for non-existent entry")
	}
}

func TestInMemoryDb_Get_EntryExists(t *testing.T) {
	var db = populatedDb()
	val, err := db.Get(Key("eggs"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	var readableVal = string(val)
	if readableVal != "scrambled" {
		t.Fatalf("expected, value for 'eggs' to be 'scrambled', got %s", readableVal)
	}
}

func TestInMemoryDb_Has_True(t *testing.T) {
	var db = populatedDb()
	var keyExists, err = db.Has(Key("eggs"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if !keyExists {
		t.Fatal("expected key to exist, return value was false")
	}
}

func TestInMemoryDb_Has_False(t *testing.T) {
	var db = populatedDb()
	var keyExists, err = db.Has(Key("foo"))
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if keyExists {
		t.Fatal("expected key not to exist, return value was true")
	}
}

func TestInMemoryDb_Put_NewEntry(t *testing.T) {
	db := populatedDb()
	err := db.Put(Key("foo"), Value("bar"))
	if err != nil {
		t.Fatal("unexpected error when putting value", err)
	}
	val, err := db.Get(Key("foo"))
	if err != nil {
		t.Fatal("unexpected error when getting value", err)
	}
	stringVal := string(val)
	if stringVal != "bar" {
		t.Fatalf("expected value 'bar', got %q\n", stringVal)
	}
}

func TestInMemoryDb_Put_UpdateEntry(t *testing.T) {
	db := populatedDb()
	err := db.Put(Key("eggs"), Value("poached"))
	if err != nil {
		t.Fatal("unexpected error when putting value", err)
	}
	val, err := db.Get(Key("eggs"))
	if err != nil {
		t.Fatal("unexpected error when getting value", err)
	}
	stringVal := string(val)
	if stringVal != "poached" {
		t.Fatalf("expected value 'poached', got %q\n", stringVal)
	}
}

func TestInMemoryDb_Delete_Success(t *testing.T) {
	db := populatedDb()
	err := db.Delete(Key("eggs"))
	if err != nil {
		t.Fatal("unexpected error while deleting key", err)
	}
	_, err = db.Get(Key("eggs"))
	if err == nil {
		t.Fatal("expected an error getting deleted key, but did not get one")
	}
	err = db.Delete(Key("eggs"))
	if err == nil {
		t.Fatal("expected an error deleting non-existent key, but did not get one")
	}
}

func TestInMemoryDb_RangeScan(t *testing.T) {
	data := []struct {
		key   Key
		value Value
	}{
		{Key("abc"), Value("ABC")},
		{Key("abd"), Value("ABD")},
		{Key("abe"), Value("ABE")},
		{Key("abf"), Value("ABF")},
		{Key("abg"), Value("ABG")},
	}
	db := populatedDb()
	var err error
	for _, datum := range data {
		err = db.Put(datum.key, datum.value)
		if err != nil {
			t.Fatal("unexpected error executing Put()", err)
		}
	}
	results, err := db.RangeScan(Key("abc"), Key("abf"))
	if err != nil {
		t.Fatal("unexpected error executing RangeScan()", err)
	}
	expectedResults := data[0:4] // `results` should include matches to `limit` parameter.
	for j, datum := range expectedResults {
		hasNext := results.Next()
		if !hasNext {
			t.Fatalf("expected more results, got %d", j+1)
		}
		if !bytes.Equal(datum.key, results.Key()) {
			t.Fatalf("expected key %q, got %q", datum.key, results.Key())
		}
		if !bytes.Equal(datum.value, results.Value()) {
			t.Fatalf("expected value %q, got %q", datum.value, results.Value())
		}
	}
	if results.Error() != nil {
		t.Fatal("iterator generated unexpected error", err)
	}
	if results.Next() {
		t.Fatal("got more results than expected")
	}
}
