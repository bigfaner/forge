//go:build e2e

package e2e

import (
	"strings"
	"testing"

	"forge-cli/tests/e2e/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Task Types (TC-021 to TC-022) ---

// Traceability: TC-021 -> Story 6 / AC-1
func TestTC_021_ListTypesOutputsAllWithDescriptions(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("task", "list-types")

	assert.Equal(t, 0, exitCode, "task list-types should exit 0")
	assert.True(t, len(strings.TrimSpace(out)) > 0,
		"list-types output should not be empty")

	knownTypes := []string{
		"feature", "fix", "gate",
		"doc-generation", "test-pipeline",
	}
	for _, typ := range knownTypes {
		assert.True(t, strings.Contains(out, typ),
			"list-types output should contain type: %s", typ)
	}

	lines := strings.Split(out, "\n")
	nonEmptyLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		nonEmptyLines++
		parts := strings.SplitN(trimmed, "  ", 2)
		assert.GreaterOrEqual(t, len(parts), 2,
			"each type line should have a description: %q", trimmed)
		if len(parts) >= 2 {
			assert.LessOrEqual(t, len(strings.TrimSpace(parts[1])), 60,
				"description should be <= 60 chars: %q", parts[1])
		}
	}
	assert.True(t, nonEmptyLines >= 5,
		"should list at least 5 task types, got %d", nonEmptyLines)
}

// Traceability: TC-022 -> Story 6 / AC-2
func TestTC_022_ListTypesEmptyRegistryReturnsEmpty(t *testing.T) {
	t.Skip("requires manual setup: empty task type registry")
}
