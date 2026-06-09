//go:build cli_functional

package startexistingflags

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 3: Start with --interactive on existing worktree
// Contract: start-existing-flags/step-3-interactive-mode.md
// ==============================================================================

// Traceability: Step 3 / Outcome "non-tty"
// Verifies that --interactive in a non-TTY environment returns an error.
func TestStep3_NonTTY_ReturnsError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)

	// Execute: forge worktree start --interactive --no-launch (non-TTY, piped stdin)
	stdout, stderr, exitCode := runForgeStartInteractive(t, projectRoot)

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code for non-TTY, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains TTY error
	assert.True(t, strings.Contains(stderr, "interactive mode requires a terminal") || strings.Contains(stderr, "TTY"),
		"expected TTY error message, got: %s", stderr)
}

// Traceability: Step 3 / Outcome "no-existing-worktrees"
// Verifies --interactive with no proposals or features shows "not found" message.
func TestStep3_NoItems_NoUnfinishedProposals(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)

	// Execute: forge worktree start --interactive --no-launch (no proposals/features)
	stdout, stderr, exitCode := runForgeStartInteractive(t, projectRoot)

	// Assert: exit code is 0 (graceful, not an error)
	assert.Equal(t, 0, exitCode, "expected exit code 0 for no items, stderr: %s", stderr)

	// Assert: stdout contains "No unfinished proposals or features found"
	assert.True(t, strings.Contains(stdout, "No unfinished proposals or features found"),
		"expected no-items message, got stdout: %s", stdout)

	// Assert: stdout contains suggestion to create one
	assert.True(t, strings.Contains(stdout, "forge proposal") || strings.Contains(stdout, "forge feature"),
		"expected suggestion to create proposal/feature, got stdout: %s", stdout)
}

// Traceability: Step 3 / Outcome "success" (partial - TTY required)
// We cannot fully test interactive mode in a non-TTY environment,
// but we can verify the setup works and the TTY check is the only blocker.
func TestStep3_Success_RequiresTTY(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "interactive-wt"

	// Create a proposal directory for interactive selection
	proposalDir := filepath.Join(projectRoot, "docs", "proposals", "test-proposal")
	if err := os.MkdirAll(proposalDir, 0755); err != nil {
		t.Fatalf("failed to create proposal dir: %v", err)
	}
	_ = os.WriteFile(filepath.Join(proposalDir, "proposal.md"), []byte("# Test Proposal\n"), 0644)

	// Create a worktree for the proposal
	_, stderr1, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")
	assert.True(t, strings.Contains(stderr1, "created new worktree"))

	// The interactive mode requires TTY which we can't provide in tests.
	// Verify the non-TTY error is returned when trying interactive mode.
	_, stderr2, exitCode2 := runForgeStartInteractive(t, projectRoot)
	assert.NotEqual(t, 0, exitCode2, "interactive should fail in non-TTY")
	assert.True(t, strings.Contains(stderr2, "interactive mode requires a terminal") || strings.Contains(stderr2, "TTY"),
		"expected TTY requirement error, got: %s", stderr2)
}
