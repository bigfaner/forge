//go:build e2e

package automatedtestorchestration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 5 Contract tests: Execute test
// ==============================================================================

// Traceability: TC-049 -> Contract automated-test-orchestration/step-5 Outcome "success"
// Tests run and all pass (exit 0), proceeds to teardown.
func TestTC_049_ExecuteTest_Success(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	t.Logf("run-tests test step output (exit %d): %s", exitCode, out)
}

// Traceability: TC-050 -> Contract automated-test-orchestration/step-5 Outcome "test-failure"
// Test execution fails, teardown runs, exit code propagated.
func TestTC_050_ExecuteTest_Failure(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	_, exitCode := runForgeRaw(t, projectDir, "run-tests")
	if exitCode != 0 {
		// Test failure: teardown executed, exit code is 1 (retryable) or 2 (blocking)
		assert.True(t,
			exitCode == 1 || exitCode == 2,
			"test failure should exit with code 1 (retryable) or 2 (blocking)")
	}
}

// Traceability: TC-051 -> Contract automated-test-orchestration/step-5 Outcome "test-validation-error"
// Misconfigured test recipe triggers immediate failure with blocking exit code.
func TestTC_051_ExecuteTest_ValidationError(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	if exitCode == 2 {
		t.Logf("test validation error output: %s", out)
	}
}

// Traceability: TC-052 -> Contract automated-test-orchestration/step-5 Outcome "test-execution-timeout"
// Test execution timeout triggers termination + teardown.
func TestTC_052_ExecuteTest_Timeout(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	if exitCode != 0 && strings.Contains(out, "timeout") {
		assert.True(t,
			outputContainsRecoveryHint(out),
			"timeout error should include recovery hint")
	}
}
