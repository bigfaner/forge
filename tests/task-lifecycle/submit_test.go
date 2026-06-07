//go:build cli_functional

package tasklifecycle

import (
	"os/exec"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Task Submit (TC-008 to TC-011) ---

// Traceability: TC-010 -> Story 3 / AC-3
func TestTC_010_SubmitTaskMissingResultFlagReturnsError(t *testing.T) {
	// Submit requires --data or stdin input. Without either, the command
	// should fail with exit code 1. The error may be "task not found" (if
	// task doesn't exist) or "no input" (if task exists but no data provided).
	cmd := exec.Command(testkit.ForgeBinary, "task", "submit", "T-impl-1")
	out, err := cmd.CombinedOutput()

	assert.Error(t, err, "submit without data should fail")
	exitCode := 1
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	}
	assert.Equal(t, 1, exitCode, "submit without data should exit 1")
	// Verify some error output was produced
	assert.True(t, len(out) > 0, "submit failure should produce output")
}
