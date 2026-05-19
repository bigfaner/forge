//go:build e2e

package e2e

import (
	"strings"
	"testing"

	"forge-cli/tests/e2e/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Profile Commands (TC-027 to TC-030) ---

// Traceability: TC-027 -> Story 8 / AC-1
func TestTC_027_ProfileDetectScansAndOutputsProfiles(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("profile", "detect")

	assert.Equal(t, 0, exitCode, "profile detect should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"profile detect output should not be empty")
}

// Traceability: TC-028 -> Story 8 / AC-2
func TestTC_028_ProfileSetUpdatesConfigWithValidProfile(t *testing.T) {
	t.Skip("requires manual setup: writable .forge/config.yaml, verify config restoration")
}

// Traceability: TC-029 -> Story 8 / AC-3
func TestTC_029_ProfileGetOutputsStrategyFileContent(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("profile", "get", "go", "--generate")

	assert.Equal(t, 0, exitCode, "profile get go --generate should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"profile get output should not be empty")
}

// Traceability: TC-030 -> Story 8 / AC-4
func TestTC_030_ProfileSetInvalidProfileReturnsErrorWithList(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("profile", "set", "nonexistent-profile")

	assert.Equal(t, 1, exitCode,
		"profile set with invalid name should exit 1")
	lower := strings.ToLower(out)
	assert.True(t, strings.Contains(lower, "unknown"),
		"output should contain 'unknown': %s", out)

	knownLanguages := []string{"javascript", "go"}
	foundAtLeastOne := false
	for _, p := range knownLanguages {
		if strings.Contains(lower, p) {
			foundAtLeastOne = true
			break
		}
	}
	assert.True(t, foundAtLeastOne,
		"output should list at least one known language: %s", out)
}
