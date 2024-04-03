package inmem

import (
	"bytes"
	"fmt"
	"leveldb"
	"slices"
)

type inMemoryDb struct {
	data []leveldb.DataEntry
}

// NewInMemoryDb copies its input slice into a new instance to prevent the database from deleting entries in a reference
// shared by other properties.
func NewInMemoryDb(data []leveldb.DataEntry) leveldb.DB {
	var copied = make([]leveldb.DataEntry, len(data))
	copy(copied, data)
	return &inMemoryDb{data: copied}
}

type inMemoryIterator struct {
	data []leveldb.DataEntry
	curr int
	err  error
}

func NewInMemoryIterator(data []leveldb.DataEntry) leveldb.Iterator {
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

func (i *inMemoryIterator) Key() leveldb.Key {
	if i.curr >= len(i.data) {
		return nil
	}
	return i.data[i.curr].Key
}

func (i *inMemoryIterator) Value() leveldb.Value {
	if i.curr >= len(i.data) {
		return nil
	}
	return i.data[i.curr].Value
}

func (db *inMemoryDb) Get(key leveldb.Key) (leveldb.Value, error) {
	var idx, found = slices.BinarySearchFunc(
		db.data,
		key,
		func(datum leveldb.DataEntry, target leveldb.Key) int { return bytes.Compare(datum.Key, target) },
	)
	if !found {
		return nil, leveldb.NewNotFoundError(key)
	}
	return db.data[idx].Value, nil
}

func (db *inMemoryDb) Has(key leveldb.Key) (bool, error) {
	var _, found = db.findEntryByKey(key)
	return found, nil
}

func (db *inMemoryDb) Put(key leveldb.Key, value leveldb.Value) error {
	var idx, keyExists = db.findEntryByKey(key)

	if keyExists {
		db.data[idx].Value = value
	} else {
		db.data = append(db.data, leveldb.DataEntry{
			Key:   key,
			Value: value,
		})
		db.sortData()
	}

	return nil
}

func (db *inMemoryDb) Delete(key leveldb.Key) error {
	var idx, keyExists = db.findEntryByKey(key)
	if !keyExists {
		return fmt.Errorf("key %q not found\n", key)
	}

	db.data = slices.Delete(db.data, idx, idx+1)
	db.sortData()

	return nil
}

func (db *inMemoryDb) RangeScan(start leveldb.Key, limit leveldb.Key) (leveldb.Iterator, error) {
	firstIdx, _ := db.findEntryByKey(start)
	lastIdx, lastIsInDataset := db.findEntryByKey(limit)
	if lastIsInDataset {
		lastIdx++ // we want to include matching entry in the return Iterator
	}
	return NewInMemoryIterator(db.data[firstIdx:lastIdx]), nil
}

func (db *inMemoryDb) findEntryByKey(key leveldb.Key) (int, bool) {
	return slices.BinarySearchFunc(
		db.data,
		key,
		func(datum leveldb.DataEntry, targetKey leveldb.Key) int { return bytes.Compare(datum.Key, targetKey) },
	)
}

func (db *inMemoryDb) sortData() {
	slices.SortFunc(db.data, func(d1, d2 leveldb.DataEntry) int {
		return bytes.Compare(d1.Key, d2.Key)
	})
}
