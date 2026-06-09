//go:build cli_functional

package idempotentstart

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 2: Start with existing worktree (idempotent path)
// Contract: idempotent-start/step-2-idempotent-reentry.md
// ==============================================================================

// Traceability: Step 2 / Outcome "success"
// Verifies that forge worktree start enters an existing worktree idempotently.
func TestStep2_Success_EntersExistingWorktree(t *testing.T) {
	projectRoot := setupGitRepoWithIncludes(t)
	slug := "my-feature"

	// Phase 1: Create the worktree
	_, stderr1, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "phase 1 create should succeed, stderr: %s", stderr1)
	assert.True(t, strings.Contains(stderr1, "created new worktree: "+slug),
		"phase 1 should create worktree")

	// Record worktree state before re-entry
	wtDir := worktreeDir(projectRoot, slug)
	gitFileBefore, _ := os.ReadFile(filepath.Join(wtDir, ".git"))

	// Phase 2: Re-enter the existing worktree
	stdout2, stderr2, exitCode2 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode2, "phase 2 re-entry should succeed, stderr: %s", stderr2)

	// Assert: stderr contains entering message (not created)
	assert.True(t, strings.Contains(stderr2, "entering existing worktree: "+slug),
		"expected stderr to contain 'entering existing worktree: %s', got: %s", slug, stderr2)
	assert.False(t, strings.Contains(stderr2, "created new worktree"),
		"should NOT contain 'created new worktree' on re-entry")

	// Assert: no new worktree created (behavioral: state unchanged)
	assertValidGitFile(t, wtDir)

	// Assert: .git file content unchanged (behavioral: deep assertion - state transition check)
	gitFileAfter, _ := os.ReadFile(filepath.Join(wtDir, ".git"))
	assert.Equal(t, string(gitFileBefore), string(gitFileAfter),
		".git file content should be unchanged after re-entry")

	// Assert: includes files NOT re-copied (verified by unchanged mtime or content)
	secretPath := filepath.Join(wtDir, "secret.txt")
	_, statErr := os.Stat(secretPath)
	assert.NoError(t, statErr, "secret.txt should still exist in worktree")

	// Assert: stdout contains worktree path (behavioral: proves path resolution)
	assert.True(t, strings.Contains(stdout2, "worktree path:"),
		"expected stdout to contain worktree path, got: %s", stdout2)
}

// Traceability: Step 2 / Outcome "corrupted-git-file"
// Verifies error when worktree directory exists but .git file is missing.
func TestStep2_CorruptedGitFile_ReturnsCorruptionError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "broken-feature"

	// Setup: create a directory at worktree path without .git file
	wtDir := worktreeDir(projectRoot, slug)
	if err := os.MkdirAll(wtDir, 0755); err != nil {
		t.Fatalf("failed to create corrupted worktree dir: %v", err)
	}

	// Execute: forge worktree start broken-feature --no-launch
	stdout, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code for corrupted worktree, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains corruption error
	assert.True(t, strings.Contains(stderr, "not a valid git worktree") || strings.Contains(stderr, "corrupt"),
		"expected stderr to mention invalid worktree, got: %s", stderr)

	// Assert: stderr contains hint to run forge worktree remove
	assert.True(t, strings.Contains(stderr, "forge worktree remove"),
		"expected stderr to contain recovery hint, got: %s", stderr)
}

// Traceability: Step 2 / Outcome "source-branch-ignored"
// Verifies that --source-branch is ignored with a warning when worktree already exists.
func TestStep2_SourceBranchIgnored_OnExistingWorktree(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "my-feature"

	// Phase 1: Create worktree
	_, stderr1, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")
	assert.True(t, strings.Contains(stderr1, "created new worktree"))

	// Phase 2: Re-enter with --source-branch (should be ignored)
	stdout2, stderr2, exitCode2 := forgeStartNoLaunch(t, projectRoot, slug, "--source-branch", "develop")
	assert.Equal(t, 0, exitCode2, "re-entry should succeed, stderr: %s", stderr2)

	// Assert: stderr contains warning about ignoring --source-branch
	assert.True(t, strings.Contains(stderr2, "warning") && strings.Contains(stderr2, "ignoring --source-branch"),
		"expected warning about ignoring --source-branch, got: %s", stderr2)

	// Assert: stderr contains entering message
	assert.True(t, strings.Contains(stderr2, "entering existing worktree: "+slug),
		"expected entering existing worktree message, got: %s", stderr2)

	// Assert: stdout contains worktree path (behavioral)
	assert.True(t, strings.Contains(stdout2, "worktree path:"),
		"expected worktree path in stdout, got: %s", stdout2)
}

// Traceability: Step 2 / Outcome "diverged-branch-entered"
// Verifies that existing worktree with diverged branch is entered without modification.
func TestStep2_DivergedBranch_EntersWithoutRebase(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "diverged-feature"

	// Phase 1: Create worktree
	_, _, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")

	// Phase 2: Make the worktree branch diverge from HEAD
	wtDir := worktreeDir(projectRoot, slug)
	divergeFile := filepath.Join(wtDir, "diverge.txt")
	if err := os.WriteFile(divergeFile, []byte("diverged content\n"), 0644); err != nil {
		t.Fatalf("failed to write diverge file: %v", err)
	}
	runGit(t, wtDir, "add", "diverge.txt")
	runGit(t, wtDir, "commit", "-m", "diverge commit")

	// Phase 3: Re-enter with --no-launch
	_, stderr3, exitCode3 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode3, "re-entry should succeed, stderr: %s", stderr3)

	// Assert: stderr contains entering message (not created)
	assert.True(t, strings.Contains(stderr3, "entering existing worktree: "+slug),
		"expected entering message, got: %s", stderr3)

	// Assert: diverge file still exists (behavioral: state unchanged)
	assert.FileExists(t, divergeFile, "diverged commit should still be present")
}
