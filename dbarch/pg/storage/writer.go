package storage

import (
	"encoding/gob"
	"os"
	"pg/tuple"
)

type Writer interface {
	WriteRow(tuple *tuple.Tuple) error
	Close()
}

type fileWriter struct {
	file *os.File
}

func NewFileWriter(filename string) (Writer, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND, os.ModeType)

	if err != nil {
		return nil, err
	}

	return &fileWriter{
		file: file,
	}, nil
}

func (fw *fileWriter) Close() {
	fw.file.Close()
}

func (fw *fileWriter) WriteRow(tuple *tuple.Tuple) error {
	encoder := gob.NewEncoder(fw.file)
	return encoder.Encode(tuple)
}
