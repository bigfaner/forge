//go:build e2e

package forgecommands

import (
	"strings"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- E2E Test Commands (TC-017 to TC-020) ---

// Traceability: TC-017 -> Story 5 / AC-1
func TestTC_017_E2ERunWithConfiguredProfileExecutesSuite(t *testing.T) {
	t.Skip("requires manual setup: .forge/config.yaml with valid profile and feature test data")
}

// Traceability: TC-018 -> Story 5 / AC-2
func TestTC_018_E2ERunNoProfileConfiguredReturnsError(t *testing.T) {
	t.Skip("requires manual setup: .forge/config.yaml with no profile field")
}

// Traceability: TC-019 -> Story 5 / AC-3
func TestTC_019_E2ERunUnknownProfileReturnsErrorWithList(t *testing.T) {
	t.Skip("requires manual setup: .forge/config.yaml with unknown profile value")
}

// Traceability: TC-020 -> Story 5 / AC-4
func TestTC_020_E2ERunNonexistentFeatureReturnsError(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("e2e", "run", "--feature", "nonexistent-feature")

	assert.Equal(t, 1, exitCode, "e2e run with nonexistent feature should exit 1")
	lower := strings.ToLower(out)
	assert.True(t,
		strings.Contains(lower, "not found") || strings.Contains(lower, "feature"),
		"output should mention feature not found: %s", out)
}
