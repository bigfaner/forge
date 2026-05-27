//go:build e2e

package forgecommands

import (
	"strings"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- E2E Test Commands (TC-017 to TC-020) ---

// Traceability: TC-020 -> Story 5 / AC-4
func TestTC_020_E2ERunNonexistentFeatureReturnsError(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("e2e", "run", "--feature", "nonexistent-feature")

	assert.Equal(t, 1, exitCode, "e2e run with nonexistent feature should exit 1")
	lower := strings.ToLower(out)
	assert.True(t,
		strings.Contains(lower, "not found") || strings.Contains(lower, "feature"),
		"output should mention feature not found: %s", out)
}
