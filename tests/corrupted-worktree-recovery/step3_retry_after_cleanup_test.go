//go:build cli_functional

package corruptedworktreerecovery

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 3: Retry start after cleanup
// Contract: corrupted-worktree-recovery/step-3-retry-start-after-cleanup.md
// ==============================================================================

// Traceability: Step 3 / Outcome "success"
// Verifies that a new worktree can be created after cleanup.
func TestStep3_Success_CreatesNewWorktreeAfterCleanup(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "recovered-wt"

	// Phase 1: Create a worktree
	_, stderr1, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")
	assert.True(t, strings.Contains(stderr1, "created new worktree"))

	// Phase 2: Remove the worktree
	stdout2, stderr2, exitCode2 := runForgeRemove(t, projectRoot, slug, "--force")
	assert.Equal(t, 0, exitCode2, "remove should succeed, stdout: %s, stderr: %s", stdout2, stderr2)

	// Verify directory is gone
	assert.NoDirExists(t, worktreeDir(projectRoot, slug))

	// Phase 3: Re-create worktree with same slug
	stdout3, stderr3, exitCode3 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode3, "re-create should succeed, stderr: %s", stderr3)

	// Assert: stderr contains creation message
	assert.True(t, strings.Contains(stderr3, "created new worktree: "+slug),
		"expected creation message after recovery, got: %s", stderr3)

	// Assert: worktree exists with valid .git
	wtDir := worktreeDir(projectRoot, slug)
	assertValidGitFile(t, wtDir)

	// Assert: stdout contains worktree path (behavioral: proves successful creation)
	assert.True(t, strings.Contains(stdout3, "worktree created at"),
		"expected worktree path, got: %s", stdout3)
}

// Traceability: Step 3 / Outcome "worktrees-dir-not-directory"
// Verifies error when .forge/worktrees exists as a file (not directory).
func TestStep3_WorktreesDirNotDirectory_ReturnsError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "file-blocker-wt"

	// Setup: create .forge/worktrees as a regular file (not a directory)
	worktreesPath := projectRoot + "/.forge/worktrees"
	_ = os.WriteFile(worktreesPath, []byte("this is a file, not a directory\n"), 0644)

	// Execute: forge worktree start file-blocker-wt --no-launch
	stdout, stderr, exitCode := runForgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code when worktrees is a file, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains error about creating directory
	assert.True(t, strings.Contains(stderr, "worktrees") || strings.Contains(stderr, "directory"),
		"expected directory creation error, got: %s", stderr)
}

// Traceability: Step 3 / Outcome "source-branch-not-found"
// Verifies error when retrying with non-existent --source-branch.
func TestStep3_SourceBranchNotFound_ReturnsError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "branch-retry-wt"

	// Execute: forge worktree start with non-existent source branch
	stdout, stderr, exitCode := runForgeStartNoLaunch(t, projectRoot, slug, "--source-branch", "nonexistent-branch")

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code for nonexistent source branch, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr mentions source branch not found
	assert.True(t, strings.Contains(stderr, "source branch") && strings.Contains(stderr, "not found"),
		"expected source branch not found error, got: %s", stderr)

	// Assert: hint to verify branch
	assert.True(t, strings.Contains(stderr, "hint") || strings.Contains(stderr, "verify"),
		"expected hint in stderr, got: %s", stderr)

	// Assert: no worktree created
	assert.NoDirExists(t, worktreeDir(projectRoot, slug))
}
