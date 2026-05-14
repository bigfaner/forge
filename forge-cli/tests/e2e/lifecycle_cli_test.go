//go:build e2e

package e2e

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// --- Hook-Triggered Lifecycle (TC-012 to TC-016) ---

// Traceability: TC-012 -> Story 4 / AC-1
func TestTC_012_CleanupRemovesTerminalStateFiles(t *testing.T) {
	t.Skip("requires manual setup: feature with terminal-state tasks and state.json")
}

// Traceability: TC-013 -> Story 4 / AC-2
func TestTC_013_QualityGateRunsCompileFmtLintTestSequence(t *testing.T) {
	t.Skip("requires manual setup: all tasks completed with compilable project")
}

// Traceability: TC-014 -> Story 4 / AC-3
func TestTC_014_CleanupNoTerminalTasksOutputsMessage(t *testing.T) {
	t.Skip("requires manual setup: feature with no terminal-state tasks")
}

// Traceability: TC-015 -> Story 4 / AC-4
func TestTC_015_QualityGateCreatesNewFixTaskOnRepeatedFailure(t *testing.T) {
	t.Skip("requires manual setup: failing quality gate with existing fix-task")
}

// Traceability: TC-016 -> Story 4 / AC-5
func TestTC_016_QualityGateStopsCreatingFixTasksAfterMax3(t *testing.T) {
	t.Skip("requires manual setup: 3 existing fix-tasks for same step")
}

// verifyMaxFixTaskError verifies the fix-task cap error message pattern.
func verifyMaxFixTaskError(t *testing.T, output, step string) {
	t.Helper()
	lower := strings.ToLower(output)
	assert.True(t,
		strings.Contains(lower, "max") && strings.Contains(lower, step),
		"output should mention max fix-tasks for step %s: %s", step, output)
}
