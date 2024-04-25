package leveldb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"leveldb/encoding"
)

type Key []byte

func (k Key) String() string {
	return string(k)
}
func (k Key) Compare(other Key) int { return bytes.Compare(k, other) }

func (k Key) Encode() ([]byte, error) {
	var (
		buf    = bytes.NewBuffer(nil)
		writer = io.Writer(buf)
		err    error
	)
	if err = binary.Write(writer, encoding.ByteOrder, uint64(len(k))); err != nil {
		return nil, err
	}
	if err = binary.Write(writer, encoding.ByteOrder, k); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Value []byte

func (v Value) String() string {
	return string(v)
}

type DataEntry struct {
	Key   Key
	Value Value
}

var ErrKeyNotFound = errors.New("key not found")

func NewNotFoundError(key Key) error {
	return fmt.Errorf("%q: %w", key, ErrKeyNotFound)
}

type ReadOnlyDB interface {
	// Get gets the value for the given key.  It returns an error if the
	// DB does not contain the key.
	Get(key Key) (Value, error)

	// Has returns true if the DB contains the given key.
	Has(key Key) (bool, error)

	// RangeScan returns an Iterator (see below) for scanning through all
	// key-value pairs in the given range, ordered by key ascending.
	RangeScan(start Key, limit Key) (Iterator, error)
}
type DB interface {
	ReadOnlyDB
	// Put sets the value for the given key.  It overwrites any previous value
	// for that key; a DB is not a multi-map.
	Put(key Key, value Value) error

	// Delete deletes the value for the given key.
	Delete(key Key) error
}

type Iterator interface {
	// Next moves the iterator to the next key/value pair.
	// It returns false if the iterator is exhausted
	Next() bool

	// Error returns any accumulated error.  Exhausting all the key/value pairs =
	// is not considered to be an error.
	Error() error

	// Key returns the key of the current key/value pair, or nil if done.
	Key() Key

	// Value returns the value of the current key/value pair, or nil if done.
	Value() Value
}
