package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	forensicpkg "forge-cli/internal/cmd/forensic"
	promptpkg "forge-cli/internal/cmd/prompt"
	taskpkg "forge-cli/internal/cmd/task"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/forgelog"

	"gopkg.in/yaml.v3"
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
	expectedGroups := []string{"task", "forensic", "prompt", "worktree", "fact"}
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
	// 5 groups (task, forensic, prompt, worktree, fact) + 4 visible top-level (version is hidden) + config + proposal + lesson + init + claude + research + surfaces + upgrade + justfile = 18 visible
	if visibleCount != 18 {
		t.Errorf("expected 18 visible commands, got %d", visibleCount)
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

	expectedTaskSubs := []string{"claim", "submit", "status", "query", "check-deps", "validate", "add", "index", "migrate"}
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

	// 5 groups (task, forensic, prompt, worktree, fact) + 5 top-level (cleanup, quality-gate, verify-task-done, feature, version) + config + proposal + lesson + init + claude + research + surfaces + upgrade + justfile = 19
	if len(explicit) != 19 {
		t.Errorf("expected 19 explicit commands, got %d: %v", len(explicit), explicit)
	}
}

// --- forgelog integration tests ---

// ptrBool is a test helper that returns a pointer to the given bool value.
func ptrBool(v bool) *bool { return &v }

func TestPersistentPreRun_InitsWithProjectRoot(t *testing.T) {
	// Clear env vars that override FindProjectRoot (set by CI/hook environments)
	t.Setenv("CLAUDE_PROJECT_DIR", "")
	t.Setenv("PROJECT_ROOT", "")
	// Unset FORGE_NO_LOG (set by TestMain) to allow file backend
	t.Setenv("FORGE_NO_LOG", "")

	// Setup: create a temp directory with project markers and config
	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0o644)
	forgeDir := filepath.Join(dir, feature.ForgeDir)
	_ = os.MkdirAll(forgeDir, 0o755)

	// Write config with logs section
	cfg := forgeconfig.Config{
		Logs: &forgeconfig.LogsConfig{
			Enabled:       ptrBool(true),
			Level:         "info",
			RetentionDays: 7,
		},
	}
	configData, _ := yaml.Marshal(&cfg)
	_ = os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), configData, 0o644)

	// Change to project dir so FindProjectRoot resolves
	oldWd, _ := os.Getwd()
	_ = os.Chdir(dir)
	t.Cleanup(func() { _ = os.Chdir(oldWd) })

	cmd := rootCmd
	err := persistentPreRun(cmd, nil)
	if err != nil {
		t.Fatalf("persistentPreRun() error = %v", err)
	}
	t.Cleanup(func() { forgelog.Close() })

	// Verify log directory was created
	logsDir := filepath.Join(dir, feature.ForgeLogsDir)
	info, err := os.Stat(logsDir)
	if err != nil {
		t.Fatalf("logsDir should exist after persistentPreRun: %v", err)
	}
	if !info.IsDir() {
		t.Error("logsDir should be a directory")
	}
}

func TestPersistentPreRun_FallbackWithoutProject(t *testing.T) {
	// Clear env vars that override FindProjectRoot
	t.Setenv("CLAUDE_PROJECT_DIR", "")
	t.Setenv("PROJECT_ROOT", "")

	// In a temp dir with no project markers, should not fail
	dir := t.TempDir()
	oldWd, _ := os.Getwd()
	_ = os.Chdir(dir)
	t.Cleanup(func() { _ = os.Chdir(oldWd) })

	err := persistentPreRun(rootCmd, nil)
	if err != nil {
		t.Fatalf("persistentPreRun() without project should not fail: %v", err)
	}
	t.Cleanup(func() { forgelog.Close() })
}

func TestPersistentPreRun_WithEnvDisable(t *testing.T) {
	// Clear env vars that override FindProjectRoot
	t.Setenv("CLAUDE_PROJECT_DIR", "")
	t.Setenv("PROJECT_ROOT", "")
	t.Setenv("FORGE_NO_LOG", "1")

	dir := t.TempDir()
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0o644)
	forgeDir := filepath.Join(dir, feature.ForgeDir)
	_ = os.MkdirAll(forgeDir, 0o755)

	oldWd, _ := os.Getwd()
	_ = os.Chdir(dir)
	t.Cleanup(func() { _ = os.Chdir(oldWd) })

	err := persistentPreRun(rootCmd, nil)
	if err != nil {
		t.Fatalf("persistentPreRun() error = %v", err)
	}
	t.Cleanup(func() { forgelog.Close() })

	// Verify log directory was NOT created
	logsDir := filepath.Join(dir, feature.ForgeLogsDir)
	if _, err := os.Stat(logsDir); !os.IsNotExist(err) {
		t.Error("logsDir should NOT be created when FORGE_NO_LOG=1")
	}
}

func TestForgelogClose_ReleasesResources(t *testing.T) {
	// Init logging first
	dir := t.TempDir()
	logsDir := filepath.Join(dir, feature.ForgeLogsDir)
	_ = forgelog.Init(nil, logsDir)

	// Close should not panic and should release file handles
	forgelog.Close()

	// Double close should also be safe
	forgelog.Close()
}

func TestForgelogIntegration_SubmitWritesAutoRestore(t *testing.T) {
	// Unset FORGE_NO_LOG (set by TestMain) to allow file backend
	t.Setenv("FORGE_NO_LOG", "")

	// AC-2: verify that forgelog writes structured log output with correct format
	dir := t.TempDir()
	logsDir := filepath.Join(dir, feature.ForgeLogsDir)

	_ = forgelog.Init(nil, logsDir)
	t.Cleanup(func() { forgelog.Close() })

	// Simulate AUTO-RESTORE diagnostic from submit
	forgelog.Info("AUTO-RESTORE: source task %s restored to pending\n", "3")

	forgelog.Close()

	// Find the log file
	entries, err := os.ReadDir(logsDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one log file")
	}

	data, err := os.ReadFile(filepath.Join(logsDir, entries[0].Name()))
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)

	// Verify structured format: 2006-01-02T15:04:05.000 [LEVEL] message
	if !strings.Contains(content, "[INFO]") {
		t.Errorf("log file should contain [INFO], got: %q", content)
	}
	if !strings.Contains(content, "AUTO-RESTORE: source task 3 restored to pending") {
		t.Errorf("log file should contain AUTO-RESTORE message, got: %q", content)
	}

	// Verify timestamp format prefix (regex: ^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3})
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		t.Fatal("expected at least one line in log file")
	}
	line := lines[0]
	// Check format: YYYY-MM-DDTHH:MM:SS.mmm [LEVEL] msg
	if len(line) < 24 {
		t.Errorf("log line too short: %q", line)
	}
	tsPart := line[:23] // "2026-06-04T17:30:00.123"
	if !strings.Contains(tsPart, "T") || !strings.Contains(tsPart, ".") {
		t.Errorf("log line timestamp format wrong: %q (expected YYYY-MM-DDTHH:MM:SS.mmm)", tsPart)
	}
}
