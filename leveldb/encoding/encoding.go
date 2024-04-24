package encoding

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	uint8Size  = 1
	uint64Size = 8
)

type Key []byte

func (key Key) Compare(other Key) int {
	return bytes.Compare(key, other)
}

type Value []byte

type Opcode interface {
	fmt.Stringer
	IncludeValue() bool
}

type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode([]byte) error
}

var ByteOrder binary.ByteOrder = binary.LittleEndian

type opcode uint8

func (o opcode) String() string {
	switch o {
	case OpPut:
		return "PUT"
	case OpDelete:
		return "DELETE"
	default:
		return "UNKNOWN"
	}
}

func (o opcode) IncludeValue() bool {
	switch o {
	case OpPut:
		return true
	default:
		return false
	}
}

const (
	_ opcode = iota
	OpPut
	OpDelete
)

type Entry struct {
	Key   Key
	Value Value
}

func (e *Entry) Encode() ([]byte, error) {
	// similar to DbOperation.Encode(), not sure if can be DRYed or if it should be
	var (
		buf      = bytes.NewBuffer(nil)
		writer   = io.Writer(buf)
		valueLen = uint64(len(e.Value))
	)

	if err := binary.Write(writer, ByteOrder, uint64(len(e.Key))); err != nil {
		return nil, err
	}
	if err := binary.Write(writer, ByteOrder, e.Key); err != nil {
		return nil, err
	}

	if err := binary.Write(writer, ByteOrder, valueLen); err != nil {
		return nil, err
	}
	if valueLen > 0 {
		if err := binary.Write(writer, ByteOrder, e.Value); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

type DbOperation struct {
	Operation Opcode
	Entry
}

func DecodeLogFile(reader *bufio.Reader) ([]*DbOperation, error) {
	var (
		err    error
		ops    []*DbOperation
		lenBuf = new(uint64)
	)
	for {
		_, err = reader.Peek(uint64Size)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}

		if err = binary.Read(reader, ByteOrder, lenBuf); err != nil {
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
		if err = opBuf.Decode(opBytes); err != nil {
			return nil, err
		}
		ops = append(ops, opBuf)
	}

	return ops, nil
}

func ReadByteSlice(r io.Reader, size uint64) ([]byte, error) {
	var buf = make([]byte, size)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadUint64(r io.Reader) (uint64, error) {
	var buf, err = ReadByteSlice(r, uint64Size)
	if err != nil {
		return 0, err
	}
	return ByteOrder.Uint64(buf), nil
}

func WriteUint64(w io.Writer, v uint64) error {
	var buf = make([]byte, uint64Size)
	ByteOrder.PutUint64(buf, v)
	bytesWritten, err := w.Write(buf)
	if bytesWritten != len(buf) {
		return errors.New("encoding.WriteUint64: not all bytes written")
	} else if err != nil {
		return err
	}
	return nil
}

// do not pass in the encoded size of the total packet
func (e *DbOperation) Decode(i []byte) error {
	var (
		buf       = bytes.NewBuffer(i)
		err       error
		key       Key
		keyLenBuf uint64
		opcodeBuf opcode
		value     Value
		valLenBuf uint64
	)
	if err = binary.Read(buf, ByteOrder, &opcodeBuf); err != nil {
		return err
	}
	// read key length, read key data
	if err = binary.Read(buf, ByteOrder, &keyLenBuf); err != nil {
		return err
	}
	key = make(Key, keyLenBuf)
	err = binary.Read(buf, ByteOrder, key)
	if err != nil {
		return err
	}
	if opcodeBuf.IncludeValue() {
		// read value length, read value data
		if err = binary.Read(buf, ByteOrder, &valLenBuf); err != nil {
			return err
		}
		value = make(Value, valLenBuf)
		err = binary.Read(buf, ByteOrder, value)
		if err != nil {
			return err
		}
	}

	e.Operation = opcodeBuf
	e.Key = key
	e.Value = value
	return nil
}

func (e *DbOperation) Encode() ([]byte, error) {
	var (
		buf bytes.Buffer
		w   = io.Writer(&buf)
	)
	var totalLen = uint64(
		uint8Size + // opcode is uint8
			uint64Size + // keyLen encoding is uint64
			len(e.Key)) // length of key (a byte array)

	if e.Operation.IncludeValue() {
		totalLen += uint64(uint64Size + len(e.Value))
	}

	if err := binary.Write(w, ByteOrder, totalLen); err != nil {
		return nil, err
	}
	if err := binary.Write(w, ByteOrder, e.Operation); err != nil {
		return nil, err
	}

	if err := binary.Write(w, ByteOrder, uint64(len(e.Key))); err != nil {
		return nil, err
	}
	if err := binary.Write(w, ByteOrder, e.Key); err != nil {
		return nil, err
	}

	if e.Operation.IncludeValue() {
		if err := binary.Write(w, ByteOrder, uint64(len(e.Value))); err != nil {
			return nil, err
		}
		if err := binary.Write(w, ByteOrder, e.Value); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
