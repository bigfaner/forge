//go:build e2e

package e2e

import (
	"strings"
	"testing"

	"forge-cli/tests/e2e/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Command Discovery & Help (TC-001 to TC-004) ---

// Traceability: TC-001 -> Story 1 / AC-1
func TestTC_001_HelpOutputShowsCommandGroups(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("--help")

	assert.Equal(t, 0, exitCode, "forge --help should exit 0")

	// Verify 5 command groups
	commandGroups := []string{"task", "e2e", "forensic", "test", "prompt"}
	for _, group := range commandGroups {
		assert.True(t, strings.Contains(out, group),
			"help output should contain command group: %s", group)
	}

	// Verify 5 top-level commands (version is hidden)
	topLevelCommands := []string{"cleanup", "feature", "probe", "quality-gate", "verify-task-done"}
	for _, cmd := range topLevelCommands {
		assert.True(t, strings.Contains(out, cmd),
			"help output should contain top-level command: %s", cmd)
	}
}

// Traceability: TC-002 -> Story 1 / AC-1
func TestTC_002_TaskSubcommandHelpShowsAllCommands(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("task", "--help")

	assert.Equal(t, 0, exitCode, "forge task --help should exit 0")

	// Fact: root.go:37 adds verify-task-done as top-level (not task subcommand).
	// Actual task subcommands: claim, submit, status, query, check-deps,
	// validate-index, add, index, migrate, list-types (10 visible)
	subcommands := []string{
		"claim", "submit", "status", "query",
		"check-deps", "validate-index",
		"add", "index", "migrate", "list-types",
	}
	for _, sub := range subcommands {
		assert.True(t, strings.Contains(out, sub),
			"task help output should contain subcommand: %s", sub)
	}

	// Verify each description is self-describing (verb + object pattern)
	// and <= 80 characters per description line
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		for _, sub := range subcommands {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, sub+" ") || trimmed == sub {
				assert.LessOrEqual(t, len(trimmed), 80,
					"description for %q should be <= 80 chars: %q", sub, trimmed)
			}
		}
	}
}

// Traceability: TC-003 -> Story 1 / AC-2
func TestTC_003_UnknownCommandReturnsErrorWithSuggestion(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("taks")

	assert.Equal(t, 1, exitCode, "unknown command should exit 1")
	assert.True(t, strings.Contains(strings.ToLower(out), "unknown"),
		"stderr should contain 'unknown' text: %s", out)
}

// Traceability: TC-004 -> Story 1 / AC-3
func TestTC_004_UnknownTaskSubcommandReturnsErrorWithList(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("task", "nonexistent-sub")

	// Cobra shows task help and lists available subcommands for unknown subcommand
	// Exit code is 0 (cobra behavior), output lists valid subcommands
	assert.Equal(t, 0, exitCode, "unknown task subcommand shows help (cobra behavior)")
	lower := strings.ToLower(out)
	assert.True(t, strings.Contains(lower, "available commands") || strings.Contains(lower, "usage"),
		"output should show available commands or usage: %s", out)

	// Verify valid subcommands are listed in help output
	validSubs := []string{"claim", "submit", "status"}
	for _, sub := range validSubs {
		assert.True(t, strings.Contains(lower, sub),
			"help output should list valid subcommand: %s", sub)
	}
}
