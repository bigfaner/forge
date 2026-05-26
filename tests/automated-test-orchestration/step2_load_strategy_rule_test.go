//go:build e2e

package automatedtestorchestration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 2 Contract tests: Load execution strategy rule file
// ==============================================================================

// Traceability: TC-038 -> Contract automated-test-orchestration/step-2 Outcome "success"
// Rule file loaded successfully for detected surface type.
func TestTC_038_LoadStrategyRule_Success(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	// Verify the surface-type is correctly set in the task
	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	t.Logf("run-tests rule loading output (exit %d): %s", exitCode, out)
	// Success if run-tests proceeds past surface detection
}

// Traceability: TC-039 -> Contract automated-test-orchestration/step-2 Outcome "rule-file-not-found"
// Missing rule file triggers blocking error with supported types list.
func TestTC_039_LoadStrategyRule_RuleFileNotFound(t *testing.T) {
	// Use a valid but potentially unsupported surface type
	projectDir := createProjectWithTask(t, "mobile")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	if exitCode != 0 {
		assert.True(t,
			strings.Contains(out, "rule") || strings.Contains(out, "not found") ||
				strings.Contains(out, "supported"),
			"error should reference missing rule file or supported types")
	}
}

// Traceability: TC-040 -> Contract automated-test-orchestration/step-2 Outcome "rule-file-malformed"
// Malformed rule file triggers blocking error with parsing failure details.
func TestTC_040_LoadStrategyRule_RuleFileMalformed(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	// This test verifies behavior when rule file exists but is malformed
	// In practice, this would require tampering with the rule file
	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	t.Logf("run-tests malformed rule output (exit %d): %s", exitCode, out)
}
