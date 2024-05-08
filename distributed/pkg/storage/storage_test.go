package storage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

type impl struct {
	name    string
	factory func() (Storage, error)
}

var impls = []impl{
	{"InMemory", NewInMemoryStorage},
	{"Persistent", persistentStorageFactory},
}

func persistentStorageFactory() (Storage, error) {
	return NewPersistentStorage(os.TempDir())
}

func closeDb(t *testing.T, db Storage) {
	if err := db.Close(); err != nil {
		t.Errorf("error closing storage: %v", err)
	}
	if err := os.Remove(filepath.Join(os.TempDir(), primaryFileName)); err != nil {
		t.Log(fmt.Sprintf("error removing DB file: %v", err))
	}
}

func TestStorage_Get(t *testing.T) {
	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			t.Run("Exists", func(t *testing.T) {
				db, err := impl.factory()
				defer closeDb(t, db)
				if err != nil {
					t.Fatalf("unexpected error initializing storage: %v", err)
				}
				if _, err := db.Put("spam", "ham"); err != nil {
					t.Fatalf("unexpected error putting value: %v", err)
				}
				value, err := db.Get("spam")
				if err != nil {
					t.Errorf("unexpected error getting value: %v", err)
				} else if value != "ham" {
					t.Errorf("expected value %q, got %q", "ham", value)
				}
			})
			t.Run("DoesNotExist", func(t *testing.T) {
				db, err := impl.factory()
				defer closeDb(t, db)
				if err != nil {
					t.Fatalf("unexpected error initializing storage: %v", err)
				}
				value, err := db.Get("spam")
				if err == nil {
					t.Errorf("expected an error, got none (value %q)", value)
				} else if !errors.Is(err, ErrNotFound) {
					t.Errorf("expected ErrNotFound, got %T (%v)", err, err)
				}
			})
		})
	}
}

func TestStorage_Put(t *testing.T) {
	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			db, err := impl.factory()
			defer closeDb(t, db)
			if err != nil {
				t.Fatalf("unexpected error initializing storage: %v", err)
			}
			for _, value := range []Value{"ham", "eggs"} {
				if _, err := db.Put("spam", value); err != nil {
					t.Fatalf("unexpected error putting value: %v", err)
				}
			}
			value, err := db.Get("spam")
			if err != nil {
				t.Fatalf("unexpected error getting value: %v", err)
			}
			if value != "eggs" {
				t.Errorf("expected value %q, got %q", "eggs", value)
			}
		})
	}
}
