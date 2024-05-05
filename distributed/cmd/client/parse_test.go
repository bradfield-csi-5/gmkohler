package main

import (
	"distributed/pkg/client"
	"testing"
)

func TestParseInput_Get(t *testing.T) {
	command, err := Input("get foo\n")
	if err != nil {
		t.Errorf("expceted no error, found %v", err)
	} else {
		if command.Operation != client.OpGet {
			t.Errorf("expected command %v, got %v", client.OpGet, command.Operation)
		}
		if command.Key != "foo" {
			t.Errorf("expected key %q, got %q", "foo", command.Key)
		}
		if len(command.Value) > 0 {
			t.Errorf("expected empty value, got %q", command.Value)
		}
	}
}

func TestParseInput_Put(t *testing.T) {
	command, err := Input("put spam=eggs\n")
	if err != nil {
		t.Errorf("expceted no error, found %v", err)
	} else {
		if command.Operation != client.OpPut {
			t.Errorf("expected command %v, got %v", client.OpPut, command.Operation)
		}
		if command.Key != "spam" {
			t.Errorf("expected key %q, got %q", "spam", command.Key)
		}
		if command.Value != "eggs" {
			t.Errorf("expected value %q, got %q", "eggs", command.Value)
		}
	}
}

func TestParseInput_Invalid(t *testing.T) {
	command, err := Input("set spam=eggs\n")
	if err == nil {
		t.Errorf("expceted error, found command %v", command)
	}
}
