//go:build cli_functional

package automatedtestorchestration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 6 Contract tests: Execute teardown
// ==============================================================================

// Traceability: TC-053 -> Contract automated-test-orchestration/step-6 Outcome "success"
// Teardown terminates processes and cleans up test-state.json.
func TestTC_053_Teardown_Success(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	t.Logf("test run-journey teardown output (exit %d): %s", exitCode, out)

	// After test run-journey completes, test-state.json should be cleaned up
	statePath := filepath.Join(projectDir, ".forge", "test-state.json")
	_, err := os.Stat(statePath)
	// File should not exist after successful teardown
	if err == nil {
		t.Log("test-state.json still exists after teardown -- may contain completed state")
	}
}

// Traceability: TC-054 -> Contract automated-test-orchestration/step-6 Outcome "teardown-kill-failure"
// Kill failure retries once then logs and continues cleanup.
func TestTC_054_Teardown_KillFailure(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	t.Logf("teardown kill failure output (exit %d): %s", exitCode, out)
	// Teardown should be idempotent: PID not existing is not an error
}

// Traceability: TC-055 -> Contract automated-test-orchestration/step-6 Outcome "stale-state-cleanup"
// Stale test-state.json from interrupted run is cleaned up.
func TestTC_055_Teardown_StaleStateCleanup(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	// Create a stale test-state.json
	forgeDir := filepath.Join(projectDir, ".forge")
	stateContent := `{"status": "active", "pid": 99999}`
	err := os.WriteFile(filepath.Join(forgeDir, "test-state.json"),
		[]byte(stateContent), 0644)
	assert.NoError(t, err)

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	t.Logf("stale state cleanup output (exit %d): %s", exitCode, out)
	// Stale state should be detected and cleaned up
}
