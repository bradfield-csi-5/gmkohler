package client

import (
	"distributed/pkg/server/storage"
)

type Operation int

func (o Operation) String() string {
	switch o {
	case OpGet:
		return "get"
	case OpPut:
		return "put"
	default:
		return "unknown"
	}
}

type Command struct {
	Operation Operation
	Key       storage.Key
	Value     storage.Value
}

const (
	OpGet Operation = iota + 1
	OpPut
)
