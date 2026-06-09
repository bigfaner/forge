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
// Step 5: Start when worktree directory exists but .git file is missing
// Contract: idempotent-start/step-5-corrupted-directory.md
// ==============================================================================

// Traceability: Step 5 / Outcome "corrupted-directory"
// Verifies error when worktree directory exists but has no .git file.
func TestStep5_CorruptedDirectory_ReturnsCorruptionError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "broken-feature"

	// Setup: create directory at worktree path without .git file
	wtDir := worktreeDir(projectRoot, slug)
	if err := os.MkdirAll(wtDir, 0755); err != nil {
		t.Fatalf("failed to create corrupted directory: %v", err)
	}

	// Execute: forge worktree start broken-feature --no-launch
	stdout, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code for corrupted directory, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains corruption error message
	assert.True(t, strings.Contains(stderr, "not a valid git worktree") || strings.Contains(stderr, "corrupt"),
		"expected corruption error message, got: %s", stderr)

	// Assert: stderr suggests running forge worktree remove
	assert.True(t, strings.Contains(stderr, "forge worktree remove"),
		"expected recovery hint with 'forge worktree remove', got: %s", stderr)

	// Assert: corrupted directory still exists (behavioral: no auto-cleanup)
	assert.DirExists(t, wtDir, "corrupted directory should not be auto-removed")
}

// Traceability: Step 5 / Outcome "symlink-resolution-failure"
// Verifies error when worktree directory is a broken symlink.
func TestStep5_SymlinkResolutionFailure_ReturnsError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "symlink-feature"

	// Setup: create a broken symlink at the worktree path
	wtDir := worktreeDir(projectRoot, slug)
	parentDir := filepath.Dir(wtDir)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		t.Fatalf("failed to create parent dir: %v", err)
	}
	// Create symlink pointing to non-existent target
	nonExistentTarget := filepath.Join(t.TempDir(), "nonexistent")
	if err := os.Symlink(nonExistentTarget, wtDir); err != nil {
		t.Fatalf("failed to create broken symlink: %v", err)
	}

	// Execute: forge worktree start symlink-feature --no-launch
	stdout, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code for broken symlink, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains error about path resolution
	assert.True(t, strings.Contains(stderr, "resolve") || strings.Contains(stderr, "path"),
		"expected error about path resolution, got: %s", stderr)
}
