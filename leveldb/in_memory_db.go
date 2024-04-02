package leveldb

import (
	"bytes"
	"fmt"
	"slices"
)

type inMemoryDb struct {
	data []DataEntry
}

type inMemoryIterator struct {
	data []DataEntry
	curr int
	err  error
}

func NewInMemoryIterator(data []DataEntry) Iterator {
	return &inMemoryIterator{
		data: data,
		curr: -1,
	}
}

func (i *inMemoryIterator) Next() bool {
	if i.curr < len(i.data) {
		i.curr++
	}
	return i.curr < len(i.data)
}

func (i *inMemoryIterator) Error() error {
	return i.err
}

func (i *inMemoryIterator) Key() Key {
	if i.curr >= len(i.data) {
		return nil
	}
	return i.data[i.curr].Key
}

func (i *inMemoryIterator) Value() Value {
	if i.curr >= len(i.data) {
		return nil
	}
	return i.data[i.curr].Value
}

type DataEntry struct {
	Key   Key
	Value Value
}

func (db *inMemoryDb) Get(key Key) (Value, error) {
	var idx, found = slices.BinarySearchFunc(
		db.data,
		key,
		func(datum DataEntry, target Key) int { return bytes.Compare(datum.Key, target) },
	)
	if !found {
		return nil, fmt.Errorf("entry not found for key %s\n", key)
	}
	return db.data[idx].Value, nil
}

func (db *inMemoryDb) Has(key Key) (bool, error) {
	var _, found = db.findEntryByKey(key)
	return found, nil
}

func (db *inMemoryDb) Put(key Key, value Value) error {
	var idx, keyExists = db.findEntryByKey(key)

	if keyExists {
		db.data[idx].Value = value
	} else {
		db.data = append(db.data, DataEntry{
			Key:   key,
			Value: value,
		})
		db.sortData()
	}

	return nil
}

func (db *inMemoryDb) Delete(key Key) error {
	var idx, keyExists = db.findEntryByKey(key)
	if !keyExists {
		return fmt.Errorf("key %q not found\n", key)
	}

	db.data = slices.Delete(db.data, idx, idx+1)
	db.sortData()

	return nil
}

func (db *inMemoryDb) RangeScan(start Key, limit Key) (Iterator, error) {
	firstIdx, _ := db.findEntryByKey(start)
	lastIdx, lastIsInDataset := db.findEntryByKey(limit)
	if lastIsInDataset {
		lastIdx++ // we want to include matching entry in the return Iterator
	}
	return NewInMemoryIterator(db.data[firstIdx:lastIdx]), nil
}

func (db *inMemoryDb) findEntryByKey(key Key) (int, bool) {
	return slices.BinarySearchFunc(
		db.data,
		key,
		func(datum DataEntry, targetKey Key) int { return bytes.Compare(datum.Key, targetKey) },
	)
}

func (db *inMemoryDb) sortData() {
	slices.SortFunc(db.data, func(d1, d2 DataEntry) int {
		return bytes.Compare(d1.Key, d2.Key)
	})
}
