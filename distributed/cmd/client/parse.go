package main

import (
	"distributed/pkg/networking"
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

func Input(input string) (*networking.Command, error) {
	results := inputRe.FindStringSubmatch(input)
	if len(results) == 0 {
		return nil, errors.New("input is not valid syntax")
	}

	var (
		op       networking.Operation
		opIdx    int    = inputRe.SubexpIndex(labelOperation)
		opString string = results[opIdx]
	)
	switch {
	case opString == "get":
		op = networking.OpGet
	case opString == "put":
		op = networking.OpPut
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

	return &networking.Command{
		Operation: op,
		Key:       key,
		Value:     value,
	}, nil
}
