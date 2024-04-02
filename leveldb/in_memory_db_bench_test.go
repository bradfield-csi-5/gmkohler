package leveldb

import (
	"bytes"
	"crypto/rand"
	"testing"
)

const (
	keySize  = 128
	valSize  = 128
	dataSize = 10_000
)

var (
	db     DB
	keyBuf = make([]byte, keySize)
	valBuf = make([]byte, valSize)
)

func init() {
	var err error
	data := make([]DataEntry, dataSize)
	for j := range dataSize {
		if _, err = rand.Read(keyBuf); err != nil {
			panic(err)
		}
		if _, err = rand.Read(valBuf); err != nil {
			panic(err)
		}
		data[j] = DataEntry{
			Key:   keyBuf,
			Value: valBuf,
		}
	}
	db = &inMemoryDb{data: data}
}

func BenchmarkInMemoryDb_Get(b *testing.B) {
	keys := make([][]byte, b.N)
	for j := range b.N {
		_, err := rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		keys[j] = keyBuf
	}
	b.ResetTimer()
	for j := range b.N {
		_, _ = db.Get(keys[j])
	}
}

func BenchmarkInMemoryDb_Has(b *testing.B) {
	keys := make([][]byte, b.N)
	for j := range b.N {
		_, err := rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		keys[j] = keyBuf
	}
	b.ResetTimer()
	for j := range b.N {
		_, _ = db.Has(keys[j])
	}
}

func BenchmarkInMemoryDb_Put(b *testing.B) {
	var err error
	var entries = make([]DataEntry, b.N)
	for j := range b.N {
		_, err = rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		_, err = rand.Read(valBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		entries[j] = DataEntry{
			Key:   keyBuf,
			Value: valBuf,
		}
	}
	b.ResetTimer()
	for j := range b.N {
		entry := entries[j]
		_ = db.Put(entry.Key, entry.Value)
	}
}

func BenchmarkInMemoryDb_Delete(b *testing.B) {
	var err error
	var keys = make([][]byte, b.N)
	for j := range b.N {
		_, err = rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		keys[j] = keyBuf
	}
	b.ResetTimer()
	for j := range b.N {
		_ = db.Delete(keys[j])
	}
}

func BenchmarkInMemoryDb_RangeScan(b *testing.B) {
	type keyRange struct {
		start []byte
		limit []byte
	}
	ranges := make([]keyRange, b.N)
	for j := range b.N {
		start := make([]byte, keySize)
		rand.Read(start)
		end := make([]byte, keySize)
		rand.Read(end)
		if bytes.Compare(start, end) > 0 {
			start, end = end, start
		}
		ranges[j] = keyRange{start, end}
	}
	b.ResetTimer()
	for j := range b.N {
		rng := ranges[j]
		_, _ = db.RangeScan(rng.start, rng.limit)
	}
}
