package test

import (
	"testing"
)

// --- Test: forge test verify command registered ---

func TestTestVerify_CommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range Cmd.Commands() {
		if cmd.Name() == "verify" {
			found = true
			break
		}
	}
	if !found {
		t.Error("test group missing 'verify' subcommand")
	}
}

// --- Test: forge test run-journey command registered ---

func TestTestingRunJourney_CommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range Cmd.Commands() {
		if cmd.Name() == "run-journey" {
			found = true
			break
		}
	}
	if !found {
		t.Error("testing group missing 'run-journey' subcommand")
	}
}

// --- Test: test group subcommands ---

func TestTestCommand_Subcommands(t *testing.T) {
	subNames := make(map[string]bool)
	for _, cmd := range Cmd.Commands() {
		subNames[cmd.Name()] = true
	}

	// Only these subcommands should exist after simplification
	expected := []string{"promote", "run-journey", "verify"}
	for _, name := range expected {
		if !subNames[name] {
			t.Errorf("test group missing subcommand: %s (have: %v)", name, subNames)
		}
	}

	// These subcommands should NOT exist
	removed := []string{"detect", "get", "interfaces", "framework"}
	for _, name := range removed {
		if subNames[name] {
			t.Errorf("test group should NOT have subcommand: %s", name)
		}
	}
}
