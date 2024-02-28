package storage

import (
	"encoding/gob"
	"os"
	"pg/tuple"
)

type Reader interface {
	ReadRow() (*tuple.Tuple, error)
	Close()
}

func NewFileReader(filename string) (Reader, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return &fileReader{file: file}, nil
}

type fileReader struct {
	file *os.File
}

func (f fileReader) ReadRow() (*tuple.Tuple, error) {
	// scan to right spot?
	decoder := gob.NewDecoder(f.file)
	tup := new(tuple.Tuple)
	if err := decoder.Decode(tup); err != nil {
		return nil, err
	}
	return tup, nil
}

func (f fileReader) Close() {
	f.file.Close()
}
