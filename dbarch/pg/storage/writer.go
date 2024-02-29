package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"pg/tuple"
)

type Writer interface {
	WriteRow(tuple *tuple.Tuple) error
	Close() error
}

type fileWriter struct {
	numRows     int
	columnNames []string
	w           io.Writer

	numWritten int
	uvarintBuf []byte
}

func NewFileWriter(columnNames []string, numRows int, w io.Writer) Writer {
	return &fileWriter{
		w:           w,
		columnNames: columnNames,
		numRows:     numRows,
		uvarintBuf:  make([]byte, binary.MaxVarintLen64),
	}
}

func (fw *fileWriter) Close() error {
	if fw.numWritten != fw.numRows {
		return fmt.Errorf(
			"fileWriter.Close(): expected to write %d rows, but wrote: %d",
			fw.numRows,
			fw.numWritten,
		)
	}
	return nil
}

func (fw *fileWriter) WriteRow(tup *tuple.Tuple) error {
	if fw.numWritten == 0 {
		if err := fw.writeHeader(); err != nil {
			return fmt.Errorf("Writer.WriteRow(): error writing fileHeader: %v", err)
		}
	}
	if len(tup.Columns) != len(fw.columnNames) {
		return fmt.Errorf(
			"Writer.WriteRow(): tried to write tuple: %+v with %d values, but writer expects %d",
			tup,
			len(tup.Columns),
			len(fw.columnNames),
		)
	}
	for _, col := range tup.Columns {
		if err := fw.writeUvarint(uint64(len(col.Value))); err != nil {
			return fmt.Errorf("Writer.WriteRow(): error writing string length uvarint: %v", err)
		}

		if _, err := fw.w.Write([]byte(col.Value)); err != nil {
			return fmt.Errorf("Writer.WriteRow(): error writing string %v, err: %v", col, err)
		}
	}
	fw.numWritten++
	return nil
}

func (fw *fileWriter) writeHeader() error {
	var head = fileHeader{
		Version:     LatestVersion,
		NumRows:     fw.numRows,
		ColumnNames: fw.columnNames,
	}

	var buf = bytes.NewBuffer(nil)
	var headerEncoder = gob.NewEncoder(buf)
	if err := headerEncoder.Encode(head); err != nil {
		return fmt.Errorf("Writer.writeHeader(): error encoding fileHeader: %+v, err: %v", head, err)
	}
	var headerBytes = buf.Bytes()
	if err := fw.writeUvarint(uint64(len(headerBytes))); err != nil {
		return fmt.Errorf("Writer.writeHeader(): error writing fileHeader bytes uvarint length: %v", err)
	}
	if _, err := fw.w.Write(headerBytes); err != nil {
		return fmt.Errorf("Writer.writeHeader(): error writing fileHeader bytes: %v", err)
	}

	return nil
}

func (fw *fileWriter) writeUvarint(x uint64) error {
	varintLen := binary.PutUvarint(fw.uvarintBuf, x)
	if _, err := fw.w.Write(fw.uvarintBuf[:varintLen]); err != nil {
		return fmt.Errorf("writeUvarint: error writing uvarint: %v", err)
	}
	return nil
}
