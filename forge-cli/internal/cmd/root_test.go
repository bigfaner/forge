package cmd

import (
	"bytes"
	"testing"

	forensicpkg "forge-cli/internal/cmd/forensic"
	promptpkg "forge-cli/internal/cmd/prompt"
	taskpkg "forge-cli/internal/cmd/task"
	testpkg "forge-cli/internal/cmd/test"
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
	expectedGroups := []string{"task", "forensic", "test", "prompt", "worktree"}
	for _, expected := range expectedGroups {
		if !commandNames[expected] {
			t.Errorf("missing group parent: %s (have: %v)", expected, commandNames)
		}
	}

	// 5 top-level commands
	expectedTopLevel := []string{"cleanup", "quality-gate", "verify-task-done", "feature", "version"}
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
	// 6 groups (task, forensic, test, prompt, worktree, fact) + 4 visible top-level (version is hidden) + config + proposal + lesson + init + claude + research + surfaces = 17 visible
	if visibleCount != 17 {
		t.Errorf("expected 16 visible commands, got %d", visibleCount)
	}
}

func TestRootCmd_VersionIsHidden(t *testing.T) {
	if !versionCmd.Hidden {
		t.Error("versionCmd should be hidden")
	}
}

func TestRootCmd_TaskGroupHasSubcommands(t *testing.T) {
	subcommands := taskpkg.Cmd.Commands()
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

func TestRootCmd_PromptGroupHasGetByTaskId(t *testing.T) {
	subcommands := promptpkg.Cmd.Commands()
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
	subcommands := forensicpkg.Cmd.Commands()
	if len(subcommands) < 3 {
		t.Errorf("forensic group should have at least 3 subcommands, got %d", len(subcommands))
	}
}

func TestRootCmd_TestGroupHasSubcommands(t *testing.T) {
	subcommands := testpkg.Cmd.Commands()
	if len(subcommands) == 0 {
		t.Error("test group should have subcommands")
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
	// Verify rootCmd has the expected explicit commands.
	// Cobra auto-adds "completion" and "help" after first Execute(), so filter those.
	commands := rootCmd.Commands()
	explicit := []string{}
	autoGen := map[string]bool{"completion": true, "help": true}
	for _, cmd := range commands {
		if !autoGen[cmd.Name()] {
			explicit = append(explicit, cmd.Name())
		}
	}

	// 6 groups (task, forensic, test, prompt, worktree, fact) + 5 top-level (cleanup, quality-gate, verify-task-done, feature, version) + config + proposal + lesson + init + claude + research + surfaces = 18
	if len(explicit) != 18 {
		t.Errorf("expected 17 explicit commands, got %d: %v", len(explicit), explicit)
	}
}
