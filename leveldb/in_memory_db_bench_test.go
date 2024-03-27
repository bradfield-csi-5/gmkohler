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
	data := make([]dataEntry, dataSize)
	for j := range dataSize {
		if _, err = rand.Read(keyBuf); err != nil {
			panic(err)
		}
		if _, err = rand.Read(valBuf); err != nil {
			panic(err)
		}
		data[j] = dataEntry{
			Key:   keyBuf,
			Value: valBuf,
		}
	}
	db = &inMemoryDb{data: data}
}

func BenchmarkInMemoryDb_Get(b *testing.B) {
	for range b.N {
		_, err := rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		_, _ = db.Get(keyBuf)
	}
}

func BenchmarkInMemoryDb_Has(b *testing.B) {
	for range b.N {
		_, err := rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		_, _ = db.Has(keyBuf)
	}
}

func BenchmarkInMemoryDb_Put(b *testing.B) {
	var err error
	for range b.N {
		_, err = rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		_, err = rand.Read(valBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		_ = db.Put(keyBuf, valBuf)
	}
}

func BenchmarkInMemoryDb_Delete(b *testing.B) {
	var err error
	for range b.N {
		_, err = rand.Read(keyBuf)
		if err != nil {
			b.Error("error reading random bytes", err)
		}
		_ = db.Delete(keyBuf)
	}
}

func BenchmarkInMemoryDb_RangeScan(b *testing.B) {
	for range b.N {
		start := make([]byte, keySize)
		rand.Read(start)
		end := make([]byte, keySize)
		rand.Read(end)
		if bytes.Compare(start, end) > 0 {
			start, end = end, start
		}
		_, _ = db.RangeScan(start, end)
	}
}
