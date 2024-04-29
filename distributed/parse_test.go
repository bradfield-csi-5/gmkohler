package main

import "testing"

func TestParseInput_Get(t *testing.T) {
	command, err := ParseInput("get foo\n")
	if err != nil {
		t.Errorf("expceted no error, found %v", err)
	} else {
		if command.Operation != Get {
			t.Errorf("expected command %v, got %v", Get, command.Operation)
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
	command, err := ParseInput("put spam=eggs\n")
	if err != nil {
		t.Errorf("expceted no error, found %v", err)
	} else {
		if command.Operation != Put {
			t.Errorf("expected command %v, got %v", Put, command.Operation)
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
	command, err := ParseInput("set spam=eggs\n")
	if err == nil {
		t.Errorf("expceted error, found command %v", command)
	}
}
