package wal

import (
	"bufio"
	"bytes"
	"leveldb"
	"leveldb/encoding"
	"os"
	"testing"
)

func TestLog_Put(t *testing.T) {
	writeFile, err := os.CreateTemp(os.TempDir(), "wal")
	if err != nil {
		t.Fatal("failed to create temp WAL file:", err)
	}
	defer func() {
		err := os.Remove(writeFile.Name())
		if err != nil {
			t.Error("failed to remove file:", err)
		}
	}()

	var operations = []encoding.DbOperation{
		{
			Operation: encoding.OpPut,
			Entry:     encoding.Entry{Key: encoding.Key("howdely"), Value: encoding.Value("doodley")},
		},
		{
			Operation: encoding.OpDelete,
			Entry:     encoding.Entry{Key: encoding.Key("neighbor")},
		},
	}
	wal := NewLog(writeFile)

	for _, op := range operations {
		if op.Operation == encoding.OpPut {
			err := wal.Put(leveldb.Key(op.Key), leveldb.Value(op.Value))
			if err != nil {
				t.Fatal("error writing PUT to log:", err)
			}
		} else if op.Operation == encoding.OpDelete {
			err := wal.Delete(leveldb.Key(op.Key))
			if err != nil {
				t.Fatal("error writing DELETE to log:", err)
			}
		}
	}

	readFile, err := os.Open(writeFile.Name())
	if err != nil {
		t.Fatal("error opening WAL file for reading:", err)
	}
	decodedOps, err := encoding.DecodeLogFile(bufio.NewReader(readFile))
	if err != nil {
		t.Fatal("error decoding WAL file:", err)
	}
	for j, decodedOp := range decodedOps {
		originalOp := operations[j]
		if decodedOp.Operation != originalOp.Operation {
			t.Errorf("expected operation %s, got %s", originalOp.Operation, decodedOp.Operation)
		}
		if !bytes.Equal(decodedOp.Key, originalOp.Key) {
			t.Errorf("expected key %q, got %q", originalOp.Key, decodedOp.Key)
		}
		if !bytes.Equal(decodedOp.Value, originalOp.Value) {
			t.Errorf("expected value %q, got %q", originalOp.Value, decodedOp.Value)
		}
	}
}
