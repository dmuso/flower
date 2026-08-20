package main

import (
	"strings"
	"testing"
)

func TestParseCommandStartsServerWhenNoArgs(t *testing.T) {
	command, operand, err := parseCommand([]string{"server"})
	if err != nil {
		t.Fatalf("parseCommand: %v", err)
	}
	if command != "" {
		t.Fatalf("command: got %q, want empty", command)
	}
	if operand != "" {
		t.Fatalf("operand: got %q, want empty", operand)
	}
}

func TestParseCommandRecognisesMigrateAndRollback(t *testing.T) {
	for _, name := range []string{"migrate", "rollback"} {
		command, operand, err := parseCommand([]string{"server", name})
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if command != name {
			t.Fatalf("command: got %q, want %q", command, name)
		}
		if operand != "" {
			t.Fatalf("operand: got %q, want empty", operand)
		}
	}
}

func TestParseCommandRequiresForceVersion(t *testing.T) {
	_, _, err := parseCommand([]string{"server", "force"})
	if err == nil {
		t.Fatal("expected missing force version to fail")
	}
	if !strings.Contains(err.Error(), "force") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseCommandReadsForceVersion(t *testing.T) {
	command, operand, err := parseCommand([]string{"server", "force", "1"})
	if err != nil {
		t.Fatalf("parseCommand: %v", err)
	}
	if command != "force" {
		t.Fatalf("command: got %q, want %q", command, "force")
	}
	if operand != "1" {
		t.Fatalf("operand: got %q, want %q", operand, "1")
	}
}

func TestParseCommandRejectsUnknownCommand(t *testing.T) {
	_, _, err := parseCommand([]string{"server", "serve"})
	if err == nil {
		t.Fatal("expected unknown command to fail")
	}
}
