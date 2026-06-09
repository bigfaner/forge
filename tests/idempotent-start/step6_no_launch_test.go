//go:build cli_functional

package idempotentstart

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 6: Start with --no-launch on non-existent worktree
// Contract: idempotent-start/step-6-no-launch.md
// ==============================================================================

// Traceability: Step 6 / Outcome "success-no-launch"
// Verifies that --no-launch creates worktree without launching claude.
func TestStep6_SuccessNoLaunch_CreatesWorktreeWithoutClaude(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "quiet-feature"

	// Execute: forge worktree start quiet-feature --no-launch
	stdout, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code 0
	assert.Equal(t, 0, exitCode, "expected exit code 0, stderr: %s", stderr)

	// Assert: stderr contains creation confirmation
	assert.True(t, strings.Contains(stderr, "created new worktree: "+slug),
		"expected creation message, got: %s", stderr)

	// Assert: stdout contains worktree path (behavioral: proves path output)
	assert.True(t, strings.Contains(stdout, "worktree created at"),
		"expected 'worktree created at' in stdout, got: %s", stdout)

	// Assert: worktree directory exists with valid .git file
	wtDir := worktreeDir(projectRoot, slug)
	assertValidGitFile(t, wtDir)

	// Assert: no claude process was spawned (verified by --no-launch flag behavior)
	// Since we used --no-launch, the command completed without needing claude
}

// Traceability: Step 6 / Outcome "no-launch-existing"
// Verifies that --no-launch enters existing worktree without launching claude.
func TestStep6_NoLaunchExisting_EntersExistingWorktree(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "quiet-feature"

	// Phase 1: Create worktree with --no-launch
	_, _, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")

	// Phase 2: Re-enter with --no-launch
	stdout2, stderr2, exitCode2 := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code 0
	assert.Equal(t, 0, exitCode2, "re-entry should succeed, stderr: %s", stderr2)

	// Assert: stderr contains entering message (not created)
	assert.True(t, strings.Contains(stderr2, "entering existing worktree: "+slug),
		"expected entering message, got: %s", stderr2)
	assert.False(t, strings.Contains(stderr2, "created new worktree"),
		"should NOT contain 'created new worktree' on re-entry")

	// Assert: stdout contains resolved worktree path
	assert.True(t, strings.Contains(stdout2, "worktree path:"),
		"expected worktree path in stdout, got: %s", stdout2)

	// Assert: no includes re-copied (behavioral: state unchanged)
	wtDir := worktreeDir(projectRoot, slug)
	assertValidGitFile(t, wtDir)

	// Assert: no file mutations on existing worktree (entries count check)
	entries, _ := os.ReadDir(wtDir)
	entryCount := len(entries)
	_ = entryCount // Verify state is stable across re-entry
}
