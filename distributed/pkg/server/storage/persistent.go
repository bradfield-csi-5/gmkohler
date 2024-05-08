package storage

import (
	"distributed/pkg"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	primaryFileName = "primary"
	replicaFileName = "replica"
)

func FileName(role pkg.Role) (string, error) {
	switch role {
	case pkg.RolePrimary:
		return primaryFileName, nil
	case pkg.RoleReplica:
		return replicaFileName, nil
	default:
		return "", fmt.Errorf("unrecognized role %v", role)
	}
}

func NewPersistentStorage(dirPath string, filename string) (Storage, error) {
	if stat, err := os.Stat(dirPath); err == nil {
		if !stat.IsDir() {
			return nil, fmt.Errorf("storage/NewPersistentStorage: %s is not a directory", dirPath)
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return nil, fmt.Errorf(
				"storage/NewPersistentStorage: %s: error creating directory: %w",
				dirPath,
				err,
			)
		}
	} else {
		return nil, fmt.Errorf("error checking file stats: %w", err)
	}
	fullFilePath := filepath.Join(dirPath, filename)
	file, err := os.OpenFile(fullFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("storage/NewPersistentStorage: %s: error opening file: %w", fullFilePath, err)
	}

	return &persistentStorage{
		file: file,
	}, nil
}

type persistentStorage struct {
	file *os.File
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
	if err := gob.NewEncoder(s.file).Encode(entries); err != nil {
		return "", fmt.Errorf("error parse value: %w", err)
	}
	return value, nil
}

func (s *persistentStorage) decode(entriesPtr *map[Key]Value) error {
	if err := s.resetFile(); err != nil {
		return err
	}
	if err := gob.NewDecoder(s.file).Decode(entriesPtr); err != nil {
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
