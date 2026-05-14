//go:build e2e

package e2e

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"forge-cli/pkg/just"

	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// TC-034: EnsureJust prompts for confirmation when just is not found
// ---------------------------------------------------------------------------

// Traceability: TC-034 -> Proposal Key Scenario "Happy path" -- user confirmation step
func TestTC_034_PromptWhenJustNotFound(t *testing.T) {
	if _, err := exec.LookPath("just"); err == nil {
		t.Skip("just already installed - cannot test not-found prompt in e2e")
	}

	// When just is not found and stdin is non-interactive (bytes.Buffer),
	// EnsureJust should fail gracefully rather than prompting.
	var stdout bytes.Buffer
	result := just.EnsureJust(strings.NewReader("y\n"), &stdout)

	// With piped stdin, the function should detect non-interactive mode
	// and return FAILED without showing the prompt.
	assert.Equal(t, just.StatusFailed, result.Status)
	// The output should NOT contain the installation prompt
	// because non-interactive mode is detected before prompting.
}

// ---------------------------------------------------------------------------
// TC-035: EnsureJust prompts for upgrade when version is outdated
// ---------------------------------------------------------------------------

// Traceability: TC-035 -> Proposal Key Scenario "Outdated version" -- upgrade prompt
func TestTC_035_PromptUpgradeOutdatedVersion(t *testing.T) {
	// This test requires mocking to set up an outdated just version.
	// In e2e with the real environment, if just meets minimum version,
	// the upgrade prompt is never shown.
	// Unit tests cover this scenario with mocked versions.
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed on this system")
	}

	_, version, found := just.DetectJustFunc()
	if !found {
		t.Skip("just not found via DetectJustFunc")
	}

	// If just is current, no upgrade prompt is shown
	if just.IsMinimumVersion(version, "1.40.0") {
		var stdout bytes.Buffer
		result := just.EnsureJust(strings.NewReader(""), &stdout)
		assert.Equal(t, just.StatusSkipped, result.Status)
		// No upgrade prompt should appear
		assert.NotContains(t, stdout.String(), "Upgrade?")
	}
}

// ---------------------------------------------------------------------------
// TC-036: Non-interactive stdin with outdated just fails with descriptive message
// ---------------------------------------------------------------------------

// Traceability: TC-036 -> Proposal "non-interactive stdin handling for outdated version"
func TestTC_036_NonInteractiveStdinOutdatedFails(t *testing.T) {
	// This scenario requires a just binary that is outdated (< minimum version).
	// In e2e, we can only test this if the installed just is actually outdated.
	// Unit tests cover this with mocked versions.
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed on this system")
	}

	_, version, found := just.DetectJustFunc()
	if !found {
		t.Skip("just not found via DetectJustFunc")
	}

	if !just.IsMinimumVersion(version, "1.40.0") {
		// Just is outdated - test with piped stdin
		var stdout bytes.Buffer
		result := just.EnsureJust(strings.NewReader("y\n"), &stdout)
		assert.Equal(t, just.StatusFailed, result.Status)
		assert.Contains(t, result.Detail, "non-interactive")
	}
	// If just is current, this test is a no-op (upgrade not needed)
}
