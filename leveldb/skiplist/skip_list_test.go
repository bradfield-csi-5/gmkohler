package skiplist

import (
	"bytes"
	"leveldb"
	"testing"
)

const (
	insertError = "unexpected error calling skipList.Insert() with %+v: %v"
)

func TestSkipList_Search(t *testing.T) {
	testData := []leveldb.DataEntry{
		{leveldb.Key("foo"), leveldb.Value("bar")},
		{leveldb.Key("bizz"), leveldb.Value("buzz")},
		{leveldb.Key("jamb"), leveldb.Value("lamb")},
		{leveldb.Key("ball"), leveldb.Value("fall")},
		{leveldb.Key("sun"), leveldb.Value("moon")},
		{leveldb.Key("cloud"), leveldb.Value("sky")},
	}
	sl := NewSkipList()
	for _, datum := range testData {
		if err := sl.Insert(datum.Key, datum.Value); err != nil {
			t.Fatalf(insertError, datum, err)
		}
	}
	for _, datum := range testData {
		searchResult, err := sl.Search(datum.Key)
		if err != nil {
			t.Fatal("unexpected error calling skipList.Search()", err)
		}
		if !bytes.Equal(searchResult, datum.Value) {
			t.Fatalf("expected search result to be %q, got %q", datum.Value, searchResult)
		}
	}
}

func TestSkipList_Insert(t *testing.T) {
	testData := []struct {
		data         leveldb.DataEntry
		shouldUpdate bool
		newValue     leveldb.Value
	}{
		{data: leveldb.DataEntry{Key: leveldb.Key("foo"), Value: leveldb.Value("bar")}, newValue: leveldb.Value("bar")},
		{data: leveldb.DataEntry{Key: leveldb.Key("bizz"), Value: leveldb.Value("buzz")}, shouldUpdate: true, newValue: leveldb.Value("updated")},
	}
	sl := NewSkipList()
	for _, datum := range testData {
		if err := sl.Insert(datum.data.Key, datum.data.Value); err != nil {
			t.Fatalf(insertError, datum.data, err)
		}
		if datum.shouldUpdate {
			if err := sl.Insert(datum.data.Key, datum.newValue); err != nil {
				t.Fatalf(insertError, leveldb.DataEntry{Key: datum.data.Key, Value: datum.newValue}, err)
			}
		}
	}
	for _, datum := range testData {
		value, err := sl.Search(datum.data.Key)
		if err != nil {
			t.Errorf("unexpected error calling SkipList.Search() with key %q: %v", datum.data.Key, err)
		} else if !bytes.Equal(value, datum.newValue) {
			t.Errorf("expected value of %q to be %q, got %q", datum.data.Key, datum.newValue, value)
		}

	}
}

func TestSkipList_Delete(t *testing.T) {
	testData := []struct {
		data            leveldb.DataEntry
		shouldBeDeleted bool
	}{
		{leveldb.DataEntry{Key: leveldb.Key("foo"), Value: leveldb.Value("bar")}, false},
		{leveldb.DataEntry{Key: leveldb.Key("bizz"), Value: leveldb.Value("buzz")}, true},
		{leveldb.DataEntry{Key: leveldb.Key("jamb"), Value: leveldb.Value("lamb")}, false},
		{leveldb.DataEntry{Key: leveldb.Key("ball"), Value: leveldb.Value("fall")}, true},
		{leveldb.DataEntry{Key: leveldb.Key("sun"), Value: leveldb.Value("moon")}, false},
		{leveldb.DataEntry{Key: leveldb.Key("cloud"), Value: leveldb.Value("sky")}, true},
	}
	sl := NewSkipList()
	for _, datum := range testData {
		if err := sl.Insert(datum.data.Key, datum.data.Value); err != nil {
			t.Fatalf(insertError, datum.data, err)
		}
	}
	for _, datum := range testData {
		if !datum.shouldBeDeleted {
			continue
		}
		err := sl.Delete(datum.data.Key)
		if err != nil {
			t.Fatalf("unexpected error calling SkipList.Delete() with key %q: %v", datum.data.Key, err)
		}
	}
	for _, datum := range testData {
		value, err := sl.Search(datum.data.Key)
		if datum.shouldBeDeleted {
			if err == nil {
				t.Errorf(
					"expected error calling SkipList.Search() for deleted key %q, but did not get one",
					datum.data.Key,
				)
			}
		} else if !datum.shouldBeDeleted {
			if err != nil {
				t.Errorf(
					"unexpected error calling SkipList.Search() for non-deleted key %q: %v",
					datum.data.Key,
					err,
				)
			} else if !bytes.Equal(value, datum.data.Value) {
				t.Errorf(
					"expected value for key %q to be %q, got %q",
					datum.data.Key,
					datum.data.Value,
					value,
				)
			}
		}
	}
}
