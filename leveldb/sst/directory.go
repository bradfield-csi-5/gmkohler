package sst

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"leveldb"
	"leveldb/encoding"
	"slices"
)

type offset int64
type Directory struct {
	sparseKeys []leveldb.Key
	offsets    []offset
}

func NewBlankDirectory() *Directory {
	return &Directory{
		sparseKeys: make([]leveldb.Key, 0),
		offsets:    make([]offset, 0),
	}
}

func (dir *Directory) Decode(i []byte) error {
	var (
		reader = bytes.NewReader(i)
		keyLen uint64
	)
	if err := binary.Read(reader, encoding.ByteOrder, keyLen); err != nil {
		return err
	}
	panic("uh oh")
}

func (dir *Directory) Encode() ([]byte, error) {
	if len(dir.sparseKeys) != len(dir.offsets) {
		return nil, errors.New("directory does not have equal number of keys and offsets")
	}

	var (
		buf    = bytes.NewBuffer(nil)
		writer = io.Writer(buf)
	)

	for j := range len(dir.sparseKeys) {
		var (
			key    = dir.sparseKeys[j]
			offset = dir.offsets[j]
		)

		encodedKey, err := key.Encode()
		if err != nil {
			return nil, err
		}
		buf.Write(encodedKey)
		if err = binary.Write(writer, encoding.ByteOrder, uint64(offset)); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewDirectory(sparseKeys []leveldb.Key, offsets []offset) (*Directory, error) {
	//if len(sparseKeys) != len(offsets) {
	//	return nil, errors.New("sparse keys and corresponding offsets must be of equal length")
	//}
	//if len(sparseKeys) == 0 {
	//	return nil, errors.New("sparseKeys and offsets cannot be empty")
	//}

	return &Directory{
		sparseKeys: sparseKeys,
		offsets:    offsets,
	}, nil
}

func (dir *Directory) offsetFor(searchKey leveldb.Key) (offset, error) {
	if len(dir.sparseKeys) == 0 {
		return dataOffset, nil
	}
	offsetIndex, _ := slices.BinarySearchFunc(dir.sparseKeys, searchKey, func(key, key2 leveldb.Key) int {
		return bytes.Compare(key, key2)
	})

	return dir.offsets[offsetIndex], nil
}
