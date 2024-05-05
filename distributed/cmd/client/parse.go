package main

import (
	"distributed/pkg/client"
	"distributed/pkg/storage"
	"errors"
	"fmt"
	"regexp"
)

const (
	labelOperation = "operation"
	labelKey       = "key"
	labelValue     = "value"
	commandPattern = `^(?P<operation>\w{3}) (?P<key>\w+)(=(?P<value>.+))?\n$`
)

var inputRe = regexp.MustCompile(commandPattern)

func Input(input string) (*client.Command, error) {
	results := inputRe.FindStringSubmatch(input)
	if len(results) == 0 {
		return nil, errors.New("input is not valid syntax")
	}

	var (
		op       client.Operation
		opIdx    int    = inputRe.SubexpIndex(labelOperation)
		opString string = results[opIdx]
	)
	switch {
	case opString == "get":
		op = client.OpGet
	case opString == "put":
		op = client.OpPut
	default:
		return nil, fmt.Errorf("%q: unrecognized operation", opString)
	}

	var (
		keyIdx = inputRe.SubexpIndex(labelKey)
		key    = storage.Key(results[keyIdx])
		value  storage.Value
	)
	if valueIdx := inputRe.SubexpIndex(labelValue); valueIdx >= 0 {
		value = storage.Value(results[valueIdx])
	}

	return &client.Command{
		Operation: op,
		Key:       key,
		Value:     value,
	}, nil
}
