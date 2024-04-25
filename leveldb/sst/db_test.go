package sst

import (
	"bytes"
	"errors"
	"leveldb"
	"leveldb/skiplist"
	"os"
	"slices"
	"testing"
)

const indexThreshold = 0x4 // small to force dictionary use
type entry struct {
	key   string
	value string
}

func TestSSTable(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), "sst")
	if err != nil {
		t.Fatal("failed to create SST file:", err)
	}
	var testData = []entry{
		{"alpha", "male"},
		{"bravo", "network"},
		{"charlie", "horse"},
		{"delta", "airlines"},
		{"echo", "echo"},
		{"foible", "foo"},
		{"foo", "bar"},
		{"fox", "hound"},
		{"foxtrot", "dance"},
		{"frolic", "room"},
		{"garage", "band"},
		{"glazed", "donut"},
		{"golf", "ball"},
		{"grape", "soda"},
		{"hostel", "nightmare"},
		{"hotel", "california"},
		{"igloo", "vestibule"},
		{"juliett", "romeo"},
		{"kilo", "gram"},
		{"mike", "drop"},
		{"november", "rain"},
		{"oscar", "meyer"},
		{"papa", "johns"},
		{"quebec", "city"},
		{"romeo", "alfa"},
		{"sierra", "nevada"},
		{"tango", "night"},
		{"uniform", "distribution"},
		{"victor", "spoils"},
		{"whiskey", "neat"},
		{"x-ray", "vision"},
		{"yankee", "stadium"},
		{"zebra", "stripes"},
	}

	var memTable = skiplist.NewSkipList()
	for _, entry := range testData {
		if err := memTable.Insert(leveldb.Key(entry.key), leveldb.Value(entry.value)); err != nil {
			t.Fatalf("error inserting key into memTable skiplist: %v", err)
		}
	}
	var tombstones = skiplist.NewSkipList()
	for _, key := range []string{
		"aardvark",
		"alabaster",
		"ajax",
		"fog",
		"frog",
		"funk",
		"hovel",
		"icicle",
		"meter",
		"spam",
	} {
		if err := tombstones.Insert(leveldb.Key(key), nil); err != nil {
			t.Fatalf("error inserting into tombstone skiplist: %v", err)
		}
	}

	sstDb, err := BuildSSTable(file, memTable, tombstones, withSparseIndexThreshold(indexThreshold))
	if err != nil {
		t.Fatal("error building SSTable:", err)
	}

	t.Run("Get", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			value, err := sstDb.Get(leveldb.Key("foo"))
			if err != nil {
				t.Fatalf("unexpected error calling sstDb.Get(): %v", err)
			}
			if bytes.Compare(leveldb.Value("bar"), value) != 0 {
				t.Errorf("unexpected returned value.  Expected %q, got %q", "bar", value)
			}
		})
		t.Run("NoEntry", func(t *testing.T) {
			_, err := sstDb.Get(leveldb.Key("baseball"))
			if err == nil {
				t.Error("expected error calling sstDb.Get() for non-existent value, did not get one")
			}
			var notFoundError *leveldb.NotFoundError
			if !errors.As(err, &notFoundError) {
				t.Errorf("expected a NotFoundError, got %T: %v", err, err)
			}
		})
		t.Run("Tombstoned", func(t *testing.T) {
			_, err := sstDb.Get(leveldb.Key("spam"))
			if err == nil {
				t.Error("expected error calling sstDb.Get() for tombstoned value, did not get one")
			}
			var notFoundError *leveldb.NotFoundError
			if !errors.As(err, &notFoundError) {
				t.Errorf("expected a NotFoundError, got %T: %v", err, err)
			}
		})
	})

	t.Run("Has", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			exists, err := sstDb.Has(leveldb.Key("foo"))
			if err != nil {
				t.Fatalf("unexpected error calling sstDb.Has(): %v", err)
			}
			if !exists {
				t.Error("unexpected returned value.  Expected true, got false")
			}
		})
		t.Run("NoEntry", func(t *testing.T) {
			var key = leveldb.Key("baseball")
			exists, err := sstDb.Has(key)
			if err != nil {
				t.Errorf("unexpected error calling sstDb.Has() for non-existent value: %v", err)
			}
			if exists {
				t.Errorf("expected %q not to exists, got exists", key)
			}
		})
		t.Run("Tombstoned", func(t *testing.T) {
			var key = leveldb.Key("spam")
			exists, err := sstDb.Has(key)
			if err != nil {
				t.Errorf("unexpected error calling sstDb.Has() for tombstoned value: %v", err)
			}
			if exists {
				t.Errorf("expected %q not to exists, got exists", key)
			}
		})
	})

	t.Run("RangeScan", func(t *testing.T) {
		t.Run("NonEmptyIterator", func(t *testing.T) {
			startIdx := slices.IndexFunc(testData, func(e entry) bool {
				return e.key == "frolic"
			})
			endIdx := slices.IndexFunc(testData, func(e entry) bool {
				return e.key == "whiskey"
			})
			expectedRange := testData[startIdx : endIdx+1]

			results, err := sstDb.RangeScan(
				leveldb.Key("frog"), // deleted entry before "frolic" to expect "frolic" first
				leveldb.Key("whiskey"),
			)
			if err != nil {
				t.Fatalf("error executing RangeScan: %v", err)
			}
			var j int
			for results.Next() {
				expectedResult := expectedRange[j]
				if expectedResult.key != string(results.Key()) || expectedResult.value != string(results.Value()) {
					t.Errorf("expected  %q=%q, got %q=%q", expectedResult.key, expectedResult.value, results.Key(), results.Value())
				}
				j++
			}
		})
		t.Run("EmptyIterator", func(t *testing.T) {
			results, err := sstDb.RangeScan( // two tombstoned keys
				leveldb.Key("aardvark"),
				leveldb.Key("ajax"),
			)
			if err != nil {
				t.Fatalf("unexpected error executing RangeScan(): %v", err)
			}
			if results.Next() {
				t.Errorf("expected zero results, got %q=%q", results.Key(), results.Value())
			}
		})
	})
}
