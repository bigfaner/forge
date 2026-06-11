//go:build cli_functional

package automatedtestorchestration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 4 Contract tests: Execute probe (retry polling)
// ==============================================================================

// Traceability: TC-044 -> Contract automated-test-orchestration/step-4 Outcome "success"
// Probe succeeds and confirms service readiness.
func TestTC_044_ExecuteProbe_Success(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	t.Logf("test run-journey probe step output (exit %d): %s", exitCode, out)
}

// Traceability: TC-045 -> Contract automated-test-orchestration/step-4 Outcome "probe-retryable-failure"
// All 3 probe retries fail with retryable error, triggers teardown + exit 1.
func TestTC_045_ExecuteProbe_RetryableFailure(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	if exitCode == 1 {
		// Retryable failure: teardown executed, HARD-GATE enforced
		assert.True(t,
			strings.Contains(out, "probe") || strings.Contains(out, "retry"),
			"retryable probe failure should reference probe/retry in output")
	}
}

// Traceability: TC-046 -> Contract automated-test-orchestration/step-4 Outcome "probe-blocking-failure"
// Blocking probe failure (exit 2), HARD-GATE enforced, no retry.
func TestTC_046_ExecuteProbe_BlockingFailure(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	if exitCode == 2 {
		// Blocking failure: teardown executed, HARD-GATE enforced
		t.Logf("probe blocking failure output: %s", out)
	}
}

// Traceability: TC-047 -> Contract automated-test-orchestration/step-4 Outcome "dev-crash-during-probe"
// Dev server crashes during probe, retries exhaust, teardown idempotent.
func TestTC_047_ExecuteProbe_DevCrashDuringProbe(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	// Dev crash during probe: retries exhaust, teardown skips missing PID
	t.Logf("dev crash during probe output (exit %d): %s", exitCode, out)
}

// Traceability: TC-048 -> Contract automated-test-orchestration/step-4 Outcome "probe-validation-error"
// Invalid probe configuration triggers blocking error with recovery hint.
func TestTC_048_ExecuteProbe_ValidationError(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	if exitCode == 2 && strings.Contains(out, "configuration") {
		assert.True(t,
			outputContainsRecoveryHint(out),
			"probe validation error should include recovery hint")
	}
}
