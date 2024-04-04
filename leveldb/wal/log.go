package wal

import (
	"bufio"
	"fmt"
	"io"
	"leveldb"
)

type Log struct {
	//operations []DbOperation
	writer *bufio.Writer
}

func NewLog(writer io.Writer) *Log {
	var bufferedWriter = bufio.NewWriter(writer)
	//go func() {
	//	time.Sleep(time.Second)
	//	var err = bufferedWriter.Flush()
	//	if err != nil {
	//		fmt.Println("WARN:", err)
	//	}
	//}()
	//recoveryReader := bufio.NewReader(rw)
	return &Log{
		writer: bufferedWriter,
	}
}

func (log *Log) Put(key leveldb.Key, value leveldb.Value) error {
	if log == nil {
		return nil
	}
	return log.write(DbOperation{
		Operation: OpPut,
		Key:       key,
		Value:     value,
	})
}

func (log *Log) Delete(key leveldb.Key) error {
	if log == nil {
		return nil
	}
	return log.write(DbOperation{Operation: OpDelete, Key: key, Value: nil})
}

func (log *Log) write(dbOp DbOperation) error {
	encoded, err := dbOp.encode()
	if err != nil {
		return err
	}
	bytesWritten, err := log.writer.Write(encoded)
	if err != nil {
		return err
	}
	if bytesWritten != len(encoded) {
		return fmt.Errorf(
			"expected to write %d bytes, only wrote %d",
			len(encoded),
			bytesWritten,
		)
	}
	return log.writer.Flush()
}
