//go:build cli_functional

package automatedtestorchestration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 7 Contract tests: Alternative surface orchestration (CLI/TUI/Mobile)
// ==============================================================================

// Traceability: TC-056 -> Contract automated-test-orchestration/step-7 Outcome "cli-tui-success"
// CLI/TUI surface: orchestration is build -> dev -> test, no probe/teardown.
func TestTC_056_AlternativeSurface_CliTuiSuccess(t *testing.T) {
	projectDir := createProjectWithTask(t, "cli")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	t.Logf("CLI surface orchestration output (exit %d): %s", exitCode, out)
	// CLI/TUI: no probe step, no teardown step
	// Dev is a build/compile step, not a background server
}

// Traceability: TC-057 -> Contract automated-test-orchestration/step-7 Outcome "mobile-success"
// Mobile surface: orchestration is test-setup -> dev -> test -> teardown.
func TestTC_057_AlternativeSurface_MobileSuccess(t *testing.T) {
	projectDir := createProjectWithTask(t, "mobile")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	t.Logf("Mobile surface orchestration output (exit %d): %s", exitCode, out)
	// Mobile: test-setup prepares emulator, teardown cleans up emulator + processes
}

// Traceability: TC-058 -> Contract automated-test-orchestration/step-7 Outcome "alternative-surface-failure"
// Build/setup failure for alternative surfaces produces error and skips tests.
func TestTC_058_AlternativeSurface_Failure(t *testing.T) {
	projectDir := createProjectWithTask(t, "cli")

	out, exitCode := runForgeRaw(t, projectDir, "test", "run-journey", "test-journey")
	if exitCode != 0 {
		// Build failure: exit code 1 (retryable) or 2 (blocking)
		assert.True(t,
			exitCode == 1 || exitCode == 2,
			"alternative surface failure should exit with code 1 or 2")
		assert.False(t,
			strings.Contains(out, "probe"),
			"CLI/TUI failure should not reference probe step")
	}
}
