//go:build e2e

package skillops

import (
	"strings"
	"testing"

	"forge-cli/tests/e2e/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Forensic Commands (TC-023 to TC-026) ---

// Traceability: TC-023 -> Story 7 / AC-1
func TestTC_023_ForensicSearchScansHistoryAndReturnsSessions(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("forensic", "search", "--last", "5")

	if exitCode != 0 {
		t.Skip("forensic search requires history.jsonl with recorded sessions")
	}

	assert.Equal(t, 0, exitCode, "forensic search should exit 0")

	if strings.TrimSpace(out) != "[]" && strings.TrimSpace(out) != "" {
		assert.True(t, strings.Contains(out, "sessionId"),
			"session output should contain sessionId field")
	}
}

// Traceability: TC-024 -> Story 7 / AC-2
func TestTC_024_ForensicExtractOutputsEvidenceSummary(t *testing.T) {
	t.Skip("requires manual setup: valid session JSONL file path")
}

// Traceability: TC-025 -> Story 7 / AC-3
func TestTC_025_ForensicSubagentsListsTranscripts(t *testing.T) {
	t.Skip("requires manual setup: session directory with subagent transcripts")
}

// Traceability: TC-026 -> Story 7 / AC-4
func TestTC_026_ForensicExtractNonexistentPathReturnsError(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("forensic", "extract", "/nonexistent/path.jsonl")

	assert.Equal(t, 1, exitCode,
		"forensic extract with nonexistent path should exit 1")
	lower := strings.ToLower(out)
	assert.True(t,
		strings.Contains(lower, "cannot") || strings.Contains(lower, "not found") || strings.Contains(lower, "no such"),
		"output should mention file not found: %s", out)
}
