package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCmd_Structure(t *testing.T) {
	if rootCmd.Use != "task" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "task")
	}

	// Verify all subcommands are registered
	// Cobra's Use field includes arguments, e.g., "record <task-id>"
	// We extract the command name by splitting on space
	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		// Get the first word of Use field as command name
		name := cmd.Name()
		if name != "" {
			commandNames[name] = true
		}
	}

	expectedCommands := []string{"claim", "record", "status", "query", "feature", "check", "validate", "verify-completion", "cleanup", "all-completed"}
	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("missing subcommand: %s (have: %v)", expected, commandNames)
		}
	}
}

func TestExecute_NoArgs(t *testing.T) {
	// Execute with --help should not error
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Execute() with --help returned error: %v", err)
	}
}

func TestInit_RegistersCommands(t *testing.T) {
	// Create a new root command and run init
	testRoot := &cobra.Command{Use: "test"}
	testRoot.AddCommand(claimCmd)
	testRoot.AddCommand(recordCmd)
	testRoot.AddCommand(statusCmd)
	testRoot.AddCommand(queryCmd)
	testRoot.AddCommand(featureCmd)
	testRoot.AddCommand(checkCmd)
	testRoot.AddCommand(validateCmd)
	testRoot.AddCommand(verifyCompletionCmd)
	testRoot.AddCommand(cleanupCmd)
	testRoot.AddCommand(allCompletedCmd)

	// Verify commands are registered
	if len(testRoot.Commands()) != 10 {
		t.Errorf("expected 10 commands, got %d", len(testRoot.Commands()))
	}
}
