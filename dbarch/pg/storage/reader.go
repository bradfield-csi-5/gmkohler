package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"pg/tuple"
)

type Reader interface {
	ReadRow() (*tuple.Tuple, error)
	Close()
}

func NewFileReader(r io.Reader) Reader {
	return &fileReader{
		r: newByteReader(r),
	}
}

type fileReader struct {
	header  *fileHeader
	r       *byteReader
	numRead int
	next    *tuple.Tuple
}

func (f *fileReader) ReadRow() (*tuple.Tuple, error) {
	if f.header == nil {
		err := f.readHeader()
		if err != nil {
			return nil, fmt.Errorf("fileReader.ReadRow(): error reading header: %v", err)
		}
	}
	if f.numRead < f.header.NumRows {
		f.readTuple()
		f.numRead++
		return f.next, nil
	}
	return nil, nil
}
func (f *fileReader) readTuple() {
	var tup = new(tuple.Tuple)
	for _, col := range f.header.ColumnNames {
		valLen, err := binary.ReadUvarint(f.r)
		if err != nil {
			panic(fmt.Sprintf("fileReader.readTuple(): error reading next value length: %v", err))
		}
		var valBytes = make([]byte, valLen)
		if _, err := io.ReadFull(f.r, valBytes); err != nil {
			panic(fmt.Sprintf("fileReader.readTuple(): error reading value bytes: %v", err))
		}
		tup.Columns = append(tup.Columns, tuple.Column{Name: col, Value: tuple.ColumnValue(valBytes)})
	}

	f.next = tup
}

func (f *fileReader) Close() {
	// TODO
}

func (f *fileReader) readHeader() error {
	headerLength, err := binary.ReadUvarint(f.r)
	if err != nil {
		return fmt.Errorf("fileReader.readHeader(): error reading fileHeader length: %v", err)
	}
	headerBytes := make([]byte, headerLength)
	if _, err := io.ReadFull(f.r, headerBytes); err != nil {
		return fmt.Errorf("fileReader.readHeader(): error reading fileHeader bytes: %v", err)
	}
	var header = new(fileHeader)
	if err := gob.NewDecoder(bytes.NewReader(headerBytes)).Decode(header); err != nil {
		return fmt.Errorf("fileReader.readHeader(): error decoding fileHeader: %v", err)
	}
	if err = header.Validate(); err != nil {
		return fmt.Errorf("fileReader.readHeader(): %v", err)
	}
	f.header = header
	return nil
}

type byteReader struct {
	io.Reader
	byteBuf []byte
}

func newByteReader(r io.Reader) *byteReader {
	return &byteReader{
		Reader:  r,
		byteBuf: make([]byte, 1),
	}
}

func (b *byteReader) ReadByte() (byte, error) {
	n, err := b.Reader.Read(b.byteBuf)
	if err != nil {
		return 0, fmt.Errorf("byteReader.ReadByte(): error reading byte: %v", err)
	}
	if n != 1 {
		return 0, fmt.Errorf("byteReader.ReadByte(): expected to read 1 byte, but read %d", n)
	}
	return b.byteBuf[0], nil
}
