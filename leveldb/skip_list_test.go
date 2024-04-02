package leveldb

import (
	"bytes"
	"testing"
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
			t.Fatalf("unexpected error calling skipList.Insert() with %+v: %v", datum, err)
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
