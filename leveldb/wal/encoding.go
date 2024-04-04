package wal

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"leveldb" // may need to have our own Key/Value to avoid circular imports
)

const (
	uint8Size  = 1
	uint64Size = 8
)

var byteOrder binary.ByteOrder = binary.LittleEndian

type Opcode uint8

func (o Opcode) String() string {
	switch o {
	case OpPut:
		return "PUT"
	case OpDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func (o Opcode) includeValue() bool {
	switch o {
	case OpPut:
		return true
	default:
		return false
	}
}

const (
	_ Opcode = iota
	OpPut
	OpDelete
)

type DbOperation struct {
	Operation Opcode
	Key       leveldb.Key
	Value     leveldb.Value
}

func DecodeLogFile(reader *bufio.Reader) ([]*DbOperation, error) {
	var (
		err    error
		ops    []*DbOperation
		lenBuf *uint64 = new(uint64)
	)
	for {
		_, err = reader.Peek(uint64Size)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if err = binary.Read(reader, byteOrder, lenBuf); err != nil {
			return nil, err
		}
		var (
			opBytes = make([]byte, *lenBuf)
			opBuf   = new(DbOperation)
		)
		numRead, err := reader.Read(opBytes)
		if err != nil {
			return nil, fmt.Errorf("error reading operation code")
		}
		if numRead != len(opBytes) {
			return nil, fmt.Errorf("expected to read %d bytes, only read %d", *lenBuf, numRead)
		}
		if err = opBuf.decode(opBytes); err != nil {
			return nil, err
		}
		ops = append(ops, opBuf)
	}

	return ops, nil
}

// do not pass in the encoded size of the total packet
func (e *DbOperation) decode(i []byte) error {
	var (
		buf       = bytes.NewBuffer(i)
		err       error
		key       leveldb.Key
		value     leveldb.Value
		lenPtr    = new(uint64)
		opcodePtr = new(Opcode)
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

func (e *DbOperation) encode() ([]byte, error) {
	var (
		buf bytes.Buffer
		w   = io.Writer(&buf)
	)
	var totalLen = uint64(
		uint8Size + // opcode is uint8
			uint64Size + // keyLen encoding is uint64
			len(e.Key)) // length of key (a byte array)

	if e.Operation.includeValue() {
		totalLen += uint64(uint64Size + len(e.Value))
	}

	if err := binary.Write(w, byteOrder, totalLen); err != nil {
		return nil, err
	}
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
