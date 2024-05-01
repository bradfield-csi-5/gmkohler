package networking

import "distributed/pkg/storage"

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

type ExecuteCommandResponse struct {
	Value storage.Value
	Err   string
}

const (
	unknown Operation = iota
	OpGet
	OpPut
)
