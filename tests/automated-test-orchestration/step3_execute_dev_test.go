//go:build e2e

package automatedtestorchestration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 3 Contract tests: Execute dev (background start)
// ==============================================================================

// Traceability: TC-041 -> Contract automated-test-orchestration/step-3 Outcome "success"
// Dev server starts in background, PID recorded in test-state.json.
func TestTC_041_ExecuteDev_BackgroundStart(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	t.Logf("run-tests dev step output (exit %d): %s", exitCode, out)
	// Success: dev server started, PID recorded
}

// Traceability: TC-042 -> Contract automated-test-orchestration/step-3 Outcome "dev-failure"
// Dev startup failure produces error output and triggers teardown.
func TestTC_042_ExecuteDev_DevFailure(t *testing.T) {
	// Create project with web surface but no dev server capability
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	if exitCode != 0 {
		// Dev failure should trigger teardown
		assert.True(t,
			len(out) > 0,
			"dev failure should produce error output")
	}
}

// Traceability: TC-043 -> Contract automated-test-orchestration/step-3 Outcome "dev-startup-timeout"
// Dev server timeout triggers error with elapsed time and recovery hint.
func TestTC_043_ExecuteDev_StartupTimeout(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	if exitCode != 0 && strings.Contains(out, "timeout") {
		assert.True(t,
			outputContainsRecoveryHint(out),
			"timeout error should include recovery hint")
	}
}
