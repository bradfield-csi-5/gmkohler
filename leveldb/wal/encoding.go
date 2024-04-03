package wal

import (
	"bytes"
	"encoding/binary"
	"io"
	"leveldb" // may need to have our own Key/Value to avoid circular imports
)

var byteOrder binary.ByteOrder = binary.LittleEndian

type opcode uint8

func (o opcode) String() string {
	switch o {
	case opPut:
		return "PUT"
	case opDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func (o opcode) includeValue() bool {
	switch o {
	case opPut:
		return true
	default:
		return false
	}
}

const (
	opUnknown opcode = iota
	opPut
	opDelete
)

type DbOperation struct {
	Operation opcode
	Key       leveldb.Key
	Value     leveldb.Value
}

func (e *DbOperation) GobDecode(i []byte) error {
	var (
		buf       = bytes.NewBuffer(i)
		err       error
		key       leveldb.Key
		value     leveldb.Value
		lenPtr    = new(uint64)
		opcodePtr = new(opcode)
	)
	if err = binary.Read(buf, byteOrder, opcodePtr); err != nil {
		return err
	}
	// read key length, read key data
	if err = binary.Read(buf, byteOrder, lenPtr); err != nil {
		return err
	}
	key = make(leveldb.Key, *lenPtr)
	err = binary.Read(buf, byteOrder, key)
	if err != nil {
		return err
	}
	if opcodePtr.includeValue() {
		// read value length, read value data
		if err = binary.Read(buf, byteOrder, lenPtr); err != nil {
			return err
		}
		value = make(leveldb.Value, *lenPtr)
		err = binary.Read(buf, byteOrder, value)
		if err != nil {
			return err
		}
	}

	e.Operation = *opcodePtr
	e.Key = key
	e.Value = value
	return nil
}

func (e *DbOperation) GobEncode() ([]byte, error) {
	var (
		buf bytes.Buffer
		w   = io.Writer(&buf)
	)

	if err := binary.Write(w, byteOrder, e.Operation); err != nil {
		return nil, err
	}

	if err := binary.Write(w, byteOrder, uint64(len(e.Key))); err != nil {
		return nil, err
	}
	if err := binary.Write(w, byteOrder, e.Key); err != nil {
		return nil, err
	}

	if e.Operation.includeValue() {
		if err := binary.Write(w, byteOrder, uint64(len(e.Value))); err != nil {
			return nil, err
		}
		if err := binary.Write(w, byteOrder, e.Value); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
