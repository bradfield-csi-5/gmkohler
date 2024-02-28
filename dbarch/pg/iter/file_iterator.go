package iter

import (
	"errors"
	"fmt"
	"io"
	"pg/storage"
	"pg/tuple"
)

type fileIterator struct {
	reader storage.Reader
	isDone bool
}

func (f fileIterator) Init() {
	fmt.Println("init fileIterator")
}

func (f fileIterator) Next() *tuple.Tuple {
	if f.isDone {
		return nil
	}
	tup, err := f.reader.ReadRow()
	if err != nil {
		if errors.Is(err, io.EOF) {
			f.isDone = true
			return nil
		}
		panic(err)
	}
	return tup
}

func (f fileIterator) Close() {
	f.reader.Close()
}

func NewFileIterator(reader storage.Reader) Iterator {
	return &fileIterator{reader: reader}
}
