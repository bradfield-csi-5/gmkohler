package storage

import (
	"fmt"
)

type inMemoryStorage map[Key]Value

func (s inMemoryStorage) Get(key Key) (Value, error) {
	value, exists := s[key]
	if !exists {
		return "", fmt.Errorf("%q: %w", key, ErrNotFound)
	}
	return value, nil
}

func (s inMemoryStorage) Put(key Key, value Value) (Value, error) {
	s[key] = value
	return s[key], nil
}

func (s inMemoryStorage) Close() error {
	for k := range s {
		delete(s, k)
	}
	return nil
}

func NewInMemoryStorage() (Storage, error) {
	return make(inMemoryStorage), nil
}
