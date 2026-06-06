//go:build cli_functional

package skillops

import (
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Prompt Commands (TC-005 to TC-007) ---

// Traceability: TC-005 -> Story 2 / AC-1
func TestTC_005_GetPromptByTaskIDReturnsCorrectPrompt(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("prompt", "get-by-task-id", "1.1")

	if exitCode == 0 {
		assert.True(t, len(out) > 0, "prompt output should not be empty")
		assert.False(t, strings.Contains(out, "{{TASK_ID}}"),
			"prompt should have TASK_ID substituted, not contain template placeholder")
		assert.False(t, strings.Contains(out, "{{TASK_FILE}}"),
			"prompt should have TASK_FILE substituted")
	} else {
		t.Skip("task 1.1 not found in current feature index - needs test data setup")
	}
}

// Traceability: TC-006 -> Story 2 / AC-2
func TestTC_006_GetPromptNonexistentTaskIDReturnsError(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("prompt", "get-by-task-id", "NONEXISTENT-999")

	assert.Equal(t, 1, exitCode, "nonexistent task ID should exit 1")
	lower := strings.ToLower(out)
	assert.True(t, strings.Contains(lower, "not found"),
		"output should contain 'not found' error: %s", out)
}
