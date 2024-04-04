package wal

import (
	"bytes"
	"leveldb"
	"testing"
)

func TestDataEntry_Encode(t *testing.T) {
	entries := []DbOperation{
		{Operation: opDelete, Key: leveldb.Key("eggs")},
		{Operation: opPut, Key: leveldb.Key("eggs"), Value: leveldb.Value("over easy")},
	}
	for _, entry := range entries {
		t.Run(entry.Operation.String(), func(t *testing.T) {
			encoded, err := entry.encode()
			if err != nil {
				t.Fatal("error encoding DbOperation:", err)
			}
			payload := encoded[8:] // remove the "total length"
			decodedEntry := new(DbOperation)
			if err := decodedEntry.decode(payload); err != nil {
				t.Fatal("error decoding encoded DbOperation:", err)
			}
			if entry.Operation != decodedEntry.Operation {
				t.Errorf("operations do not match.  expected %s, got %s", entry.Operation, decodedEntry.Operation)
			}
			if !bytes.Equal(entry.Key, decodedEntry.Key) {
				t.Errorf("keys do not match.  expected %q, got %q", entry.Key, decodedEntry.Key)
			}
			if !bytes.Equal(entry.Value, decodedEntry.Value) {
				t.Errorf("values do not match.  expected %q, got %q", entry.Value, decodedEntry.Value)
			}
		})
	}
}
