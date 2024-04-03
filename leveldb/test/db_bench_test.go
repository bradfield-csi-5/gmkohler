package test

import (
	"bytes"
	"crypto/rand"
	"leveldb"
	"leveldb/inmem"
	"leveldb/skiplist"
	"testing"
)

const (
	keySize   = 128
	valSize   = 128
	benchSize = 1_000 // I peeked in debugger; getting around b.N == 1 with the parent of sub-benchmarks
	dataSize  = 10_000
)

type testImpl struct {
	name string
	db   leveldb.DB
}

var (
	impls  []testImpl
	keyBuf = make(leveldb.Key, keySize)
	valBuf = make(leveldb.Value, valSize)
)

func init() {
	var err error
	data := make([]leveldb.DataEntry, dataSize)
	for j := range dataSize {
		if _, err = rand.Read(keyBuf); err != nil {
			panic(err)
		}
		if _, err = rand.Read(valBuf); err != nil {
			panic(err)
		}
		data[j] = leveldb.DataEntry{
			Key:   keyBuf,
			Value: valBuf,
		}
	}
	slDb := skiplist.NewSkipListDb()
	for _, datum := range data {
		if err := slDb.Put(datum.Key, datum.Value); err != nil {
			panic("failed setup")
		}
	}
	impls = []testImpl{
		{name: "InMemory", db: inmem.NewInMemoryDb(data)},
		{name: "SkipList", db: slDb},
	}
}

func BenchmarkInMemoryDb_Get(b *testing.B) {
	keys := make([]leveldb.Key, benchSize)
	for j := range b.N {
		_, err := rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		keys[j] = keyBuf
	}
	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			for j := range b.N {
				_, _ = impl.db.Get(keys[j%benchSize])
			}
		})
	}
}

func BenchmarkInMemoryDb_Has(b *testing.B) {
	keys := make([]leveldb.Key, benchSize)
	for j := range b.N {
		_, err := rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		keys[j] = keyBuf
	}

	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			for j := range b.N {
				_, _ = impl.db.Has(keys[j%benchSize])
			}
		})
	}
}

func BenchmarkInMemoryDb_Put(b *testing.B) {
	var err error
	var entries = make([]leveldb.DataEntry, benchSize)
	for j := range b.N {
		_, err = rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		_, err = rand.Read(valBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		entries[j] = leveldb.DataEntry{
			Key:   keyBuf,
			Value: valBuf,
		}
	}
	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			for j := range b.N {
				entry := entries[j%benchSize]
				_ = impl.db.Put(entry.Key, entry.Value)
			}
		})
	}
}

func BenchmarkInMemoryDb_Delete(b *testing.B) {
	var err error
	var keys = make([][]byte, benchSize)
	for j := range b.N {
		_, err = rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		keys[j] = keyBuf
	}
	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			for j := range b.N {
				_ = impl.db.Delete(keys[j%benchSize])
			}
		})
	}
}

func BenchmarkInMemoryDb_RangeScan(b *testing.B) {
	type keyRange struct {
		start leveldb.Key
		limit leveldb.Key
	}
	ranges := make([]keyRange, benchSize)
	for j := range b.N {
		start := make(leveldb.Key, keySize)
		if _, err := rand.Read(start); err != nil {
			b.Fatal("error reading random into buffer:", err)
		}
		end := make(leveldb.Key, keySize)
		if _, err := rand.Read(end); err != nil {
			b.Fatal("error reading random into buffer:", err)
		}
		if bytes.Compare(start, end) > 0 {
			start, end = end, start
		}
		ranges[j] = keyRange{start, end}
	}

	for _, impl := range impls {
		b.Run(impl.name, func(b *testing.B) {
			for j := range b.N {
				rng := ranges[j%benchSize]
				_, _ = impl.db.RangeScan(rng.start, rng.limit)
			}
		})
	}
}
