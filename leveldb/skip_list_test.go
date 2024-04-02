package leveldb

import (
	"bytes"
	"testing"
)

const (
	insertError = "unexpected error calling skipList.Insert() with %+v: %v"
)

func TestSkipList_Search(t *testing.T) {
	testData := []dataEntry{
		{Key("foo"), Value("bar")},
		{Key("bizz"), Value("buzz")},
		{Key("jamb"), Value("lamb")},
		{Key("ball"), Value("fall")},
		{Key("sun"), Value("moon")},
		{Key("cloud"), Value("sky")},
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
		data         dataEntry
		shouldUpdate bool
		newValue     Value
	}{
		{dataEntry{Key("foo"), Value("bar")}, false, Value("bar")},
		{dataEntry{Key("bizz"), Value("buzz")}, true, Value("updated")},
	}
	sl := NewSkipList()
	for _, datum := range testData {
		if err := sl.Insert(datum.data.Key, datum.data.Value); err != nil {
			t.Fatalf(insertError, datum.data, err)
		}
		if datum.shouldUpdate {
			if err := sl.Insert(datum.data.Key, datum.newValue); err != nil {
				t.Fatalf(insertError, dataEntry{datum.data.Key, datum.newValue}, err)
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
		data            dataEntry
		shouldBeDeleted bool
	}{
		{dataEntry{Key("foo"), Value("bar")}, false},
		{dataEntry{Key("bizz"), Value("buzz")}, true},
		{dataEntry{Key("jamb"), Value("lamb")}, false},
		{dataEntry{Key("ball"), Value("fall")}, true},
		{dataEntry{Key("sun"), Value("moon")}, false},
		{dataEntry{Key("cloud"), Value("sky")}, true},
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
