package main

import (
	"distributed/pkg/networking"
	"testing"
)

func TestParseInput_Get(t *testing.T) {
	command, err := Input("get foo\n")
	if err != nil {
		t.Errorf("expceted no error, found %v", err)
	} else {
		if command.Operation != networking.OpGet {
			t.Errorf("expected command %v, got %v", networking.OpGet, command.Operation)
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
		if command.Operation != networking.OpPut {
			t.Errorf("expected command %v, got %v", networking.OpPut, command.Operation)
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
