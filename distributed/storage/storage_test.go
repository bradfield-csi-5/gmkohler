package storage

import (
	"errors"
	"testing"
)

type impl struct {
	name    string
	factory func() Storage
}

var impls = []impl{
	{"InMemory", NewInMemoryStorage},
}

func TestStorage_Get(t *testing.T) {
	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			t.Run("Exists", func(t *testing.T) {
				db := impl.factory()
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
				db := impl.factory()
				_, err := db.Get("spam")
				if !errors.Is(err, ErrNotFound) {
					t.Errorf("expected ErrNotFound, got %T (%v)", err, err)
				}
			})

		})
	}
}

func TestStorage_Put(t *testing.T) {
	for _, impl := range impls {
		t.Run(impl.name, func(t *testing.T) {
			db := impl.factory()
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
