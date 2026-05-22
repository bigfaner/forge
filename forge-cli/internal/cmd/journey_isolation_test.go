package cmd

import (
	"testing"
)

// --- Test: run-journey CLI command ---

func TestTestingRunJourney_CommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range testCmd.Commands() {
		if cmd.Name() == "run-journey" {
			found = true
			break
		}
	}
	if !found {
		t.Error("testing group missing 'run-journey' subcommand")
	}
}
