package cmd

import (
	"bytes"
	"testing"
)

func TestRootCmd_Structure(t *testing.T) {
	if rootCmd.Use != "forge" {
		t.Errorf("rootCmd.Use = %q, want %q", rootCmd.Use, "forge")
	}

	// Verify all group parents and top-level commands are registered
	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)
	for _, cmd := range commands {
		name := cmd.Name()
		if name != "" {
			commandNames[name] = true
		}
	}

	// 5 group parents
	expectedGroups := []string{"task", "e2e", "forensic", "profile", "prompt"}
	for _, expected := range expectedGroups {
		if !commandNames[expected] {
			t.Errorf("missing group parent: %s (have: %v)", expected, commandNames)
		}
	}

	// 6 top-level commands
	expectedTopLevel := []string{"cleanup", "probe", "quality-gate", "verify-task-done", "feature", "version"}
	for _, expected := range expectedTopLevel {
		if !commandNames[expected] {
			t.Errorf("missing top-level command: %s (have: %v)", expected, commandNames)
		}
	}
}

func TestRootCmd_HelpShowsTenVisibleEntries(t *testing.T) {
	commands := rootCmd.Commands()
	// Filter out Cobra auto-generated commands (completion, help)
	autoGen := map[string]bool{"completion": true, "help": true}
	visibleCount := 0
	for _, cmd := range commands {
		if !cmd.Hidden && !autoGen[cmd.Name()] {
			visibleCount++
		}
	}
	// 5 groups + 6 visible top-level (version is hidden) + config + proposal + lesson + init = 14 visible
	if visibleCount != 14 {
		t.Errorf("expected 14 visible commands, got %d", visibleCount)
	}
}

func TestRootCmd_VersionIsHidden(t *testing.T) {
	if !versionCmd.Hidden {
		t.Error("versionCmd should be hidden")
	}
}

func TestRootCmd_TaskGroupHasSubcommands(t *testing.T) {
	subcommands := taskCmd.Commands()
	if len(subcommands) == 0 {
		t.Error("task group should have subcommands")
	}

	taskSubNames := make(map[string]bool)
	for _, cmd := range subcommands {
		taskSubNames[cmd.Name()] = true
	}

	expectedTaskSubs := []string{"claim", "submit", "status", "query", "check-deps", "validate-index", "add", "index", "migrate"}
	for _, expected := range expectedTaskSubs {
		if !taskSubNames[expected] {
			t.Errorf("missing task subcommand: %s (have: %v)", expected, taskSubNames)
		}
	}
}

func TestRootCmd_E2eGroupHasValidateSpecs(t *testing.T) {
	subcommands := e2eCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "validate-specs" {
			found = true
			break
		}
	}
	if !found {
		t.Error("e2e group should have validate-specs subcommand")
	}
}

func TestRootCmd_PromptGroupHasGetByTaskId(t *testing.T) {
	subcommands := promptCmd.Commands()
	found := false
	for _, cmd := range subcommands {
		if cmd.Name() == "get-by-task-id" {
			found = true
			break
		}
	}
	if !found {
		t.Error("prompt group should have get-by-task-id subcommand")
	}
}

func TestRootCmd_ForensicGroupHasSubcommands(t *testing.T) {
	subcommands := forensicCmd.Commands()
	if len(subcommands) < 3 {
		t.Errorf("forensic group should have at least 3 subcommands, got %d", len(subcommands))
	}
}

func TestRootCmd_ProfileGroupHasSubcommands(t *testing.T) {
	subcommands := profileCmd.Commands()
	if len(subcommands) == 0 {
		t.Error("profile group should have subcommands")
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
	// Verify rootCmd has the expected 10 explicitly registered commands.
	// Cobra auto-adds "completion" and "help" after first Execute(), so filter those.
	commands := rootCmd.Commands()
	explicit := []string{}
	autoGen := map[string]bool{"completion": true, "help": true}
	for _, cmd := range commands {
		if !autoGen[cmd.Name()] {
			explicit = append(explicit, cmd.Name())
		}
	}

	// 5 groups + 6 top-level + config + proposal + lesson + init = 15
	if len(explicit) != 15 {
		t.Errorf("expected 15 explicit commands, got %d: %v", len(explicit), explicit)
	}
}
