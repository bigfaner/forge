package cmd

import (
	"strings"
	"testing"
)

func TestProfileCommand_Removed(t *testing.T) {
	// The 'profile' command should not exist on rootCmd
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "profile" {
			t.Error("forge profile command should not exist -- it should be replaced by forge test")
		}
	}
}

func TestTestingCommand_Removed(t *testing.T) {
	// The old 'testing' command should not exist on rootCmd
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "testing" {
			t.Error("forge testing command should not exist -- it is renamed to forge test")
		}
	}
}

func TestTestCommand_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "test" {
			found = true
			break
		}
	}
	if !found {
		t.Error("forge test command should be registered on rootCmd")
	}
}

func TestTestCommand_Subcommands(t *testing.T) {
	subNames := make(map[string]bool)
	for _, cmd := range testCmd.Commands() {
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

func TestTestCommand_DefaultRun_ShowsHelp(t *testing.T) {
	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test command failed: %v", err)
	}

	if !strings.Contains(output, "SUBCOMMANDS:") {
		t.Errorf("expected 'SUBCOMMANDS:' in default output, got: %q", output)
	}
}
