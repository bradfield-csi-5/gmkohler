package storage

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
)

func NewPersistentStorage(filename string) (Storage, error) {
	var (
		err  error
		file *os.File
	)

	// TODO: separate fd's for encode/decode?
	if _, err = os.Stat(filename); err == nil {
		if file, err = os.OpenFile(filename, os.O_RDWR, os.ModePerm); err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}
	} else if os.IsNotExist(err) {
		if file, err = os.Create(filename); err != nil {
			return nil, fmt.Errorf("error creating fiel: %w", err)
		}
	} else {
		return nil, fmt.Errorf("error checking file stats: %w", err)
	}

	return &persistentStorage{
		file:    file,
		encoder: gob.NewEncoder(file),
		decoder: gob.NewDecoder(file),
	}, nil
}

type persistentStorage struct {
	file    *os.File
	encoder *gob.Encoder
	decoder *gob.Decoder
}

func (s *persistentStorage) Close() error {
	return s.file.Close()
}

func (s *persistentStorage) Get(key Key) (Value, error) {
	var entries = make(map[Key]Value)
	if err := s.decode(&entries); err != nil {
		return "", err
	}

	value, exists := entries[key]
	if !exists {
		return "", fmt.Errorf("%q: %w", key, ErrNotFound)
	}
	return value, nil
}

func (s *persistentStorage) Put(key Key, value Value) (Value, error) {
	var entries = make(map[Key]Value)
	if err := s.decode(&entries); err != nil {
		return "", err
	}
	entries[key] = value
	if err := s.resetFile(); err != nil {
		return "", err
	}
	if err := s.encoder.Encode(entries); err != nil {
		return "", fmt.Errorf("error encoding value: %w", err)
	}
	return value, nil
}

func (s *persistentStorage) decode(entriesPtr *map[Key]Value) error {
	if err := s.resetFile(); err != nil {
		return err
	}
	if err := s.decoder.Decode(entriesPtr); err != nil {
		if !errors.Is(err, io.EOF) { // assume this means empty file
			return fmt.Errorf("error decoding file: %w", err)
		}
	}
	return nil
}

func (s *persistentStorage) resetFile() error {
	_, err := s.file.Seek(0, io.SeekStart)
	return err
}
