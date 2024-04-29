package storage

import "errors"

type Key string
type Value string
type Entry struct {
	Key   Key
	Value Value
}

var ErrNotFound = errors.New("key not found")

type Storage interface {
	Get(key Key) (Value, error)
	Put(key Key, value Value) (Value, error)
	Close() error
}
