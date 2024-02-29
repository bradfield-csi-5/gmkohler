package storage

import "errors"

type HeaderVersion int

const (
	LatestVersion HeaderVersion = 1
)

type fileHeader struct {
	Version     HeaderVersion
	NumRows     int
	ColumnNames []string
}

func (h *fileHeader) Validate() error {
	if h.Version == 0 {
		return errors.New("fileHeader.Validate(): Version must not be zero")
	}
	if h.NumRows == 0 {
		return errors.New("fileHeader.Validate(): NumRows must not be zero")
	}
	if len(h.ColumnNames) == 0 {
		return errors.New("fileHeader.Validate(): len(ColumnNames) must not be zero")
	}
	return nil
}
