package main

import (
	"distributed/storage"
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

type operation int

func (o operation) String() string {
	switch o {
	case Get:
		return "get"
	case Put:
		return "put"
	default:
		return "unknown"
	}
}

type Command struct {
	Operation operation
	Key       storage.Key
	Value     storage.Value
}

const (
	unknown operation = iota
	Get
	Put
)

func ParseInput(input string) (*Command, error) {
	results := inputRe.FindStringSubmatch(input)
	if len(results) == 0 {
		return nil, errors.New("input is not valid syntax")
	}

	var (
		op       operation
		opIdx    int    = inputRe.SubexpIndex(labelOperation)
		opString string = results[opIdx]
	)
	switch {
	case opString == "get":
		op = Get
	case opString == "put":
		op = Put
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

	return &Command{
		Operation: op,
		Key:       key,
		Value:     value,
	}, nil
}
