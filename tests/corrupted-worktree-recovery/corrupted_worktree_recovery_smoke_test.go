//go:build cli_functional

package corruptedworktreerecovery

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Smoke Test: corrupted-worktree-recovery Journey (happy path)
// Contract: corrupted-worktree-recovery journey.md
// ==============================================================================

// TestJourney_CorruptedWorktreeRecovery verifies the complete happy path:
// 1. Attempt start on corrupted worktree (detect corruption)
// 2. Remove corrupted worktree
// 3. Retry start (create new worktree successfully)
func TestJourney_CorruptedWorktreeRecovery(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "recovery-wt"

	// ---- Phase 1: Create a valid worktree, then corrupt it ----
	_, stderr1, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "phase 1 create should succeed")
	assert.True(t, strings.Contains(stderr1, "created new worktree: "+slug))

	// Corrupt: remove .git file
	wtDir := worktreeDir(projectRoot, slug)
	_ = os.Remove(wtDir + "/.git")

	// ---- Step 1: Attempt start on corrupted worktree ----
	_, stderrS1, exitCodeS1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.NotEqual(t, 0, exitCodeS1, "step 1: should fail on corrupted worktree")
	assert.True(t, strings.Contains(stderrS1, "not a valid git worktree") || strings.Contains(stderrS1, "corrupt"),
		"step 1: expected corruption error, got: %s", stderrS1)
	assert.True(t, strings.Contains(stderrS1, "forge worktree remove"),
		"step 1: expected recovery hint")

	// ---- Step 2: Remove corrupted worktree ----
	stdoutS2, stderrS2, exitCodeS2 := runForgeRemove(t, projectRoot, slug, "--force")
	assert.Equal(t, 0, exitCodeS2, "step 2: remove should succeed, stderr: %s", stderrS2)
	assert.True(t, strings.Contains(stdoutS2, "Removed worktree"),
		"step 2: expected removal confirmation")
	assert.NoDirExists(t, wtDir, "step 2: directory should be removed")

	// ---- Step 3: Retry start after cleanup ----
	stdoutS3, stderrS3, exitCodeS3 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCodeS3, "step 3: re-create should succeed, stderr: %s", stderrS3)
	assert.True(t, strings.Contains(stderrS3, "created new worktree: "+slug),
		"step 3: expected creation message")
	assert.True(t, strings.Contains(stdoutS3, "worktree created at"),
		"step 3: expected worktree path in stdout")

	// ---- Verify final state ----
	assertValidGitFile(t, wtDir)
	assert.DirExists(t, wtDir, "worktree should exist after recovery")

	// ---- Journey Invariant: recovery restores valid state ----
	// The worktree was corrupted, removed, and re-created -- ending in a valid state
	gitFile := wtDir + "/.git"
	_, err := os.Stat(gitFile)
	assert.NoError(t, err, ".git file should exist after recovery")

	// ---- Journey Invariant: exit code is non-zero for errors, zero for success ----
	assert.NotEqual(t, 0, exitCodeS1, "step 1 should fail (corruption)")
	assert.Equal(t, 0, exitCodeS2, "step 2 should succeed (remove)")
	assert.Equal(t, 0, exitCodeS3, "step 3 should succeed (re-create)")
}
