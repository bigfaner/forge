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
// Smoke Test: idempotent-start Journey (happy path)
// Runs the complete happy path end-to-end: create -> re-enter -> verify state -> verify includes
// Contract: idempotent-start journey.md
// ==============================================================================

// TestJourney_IdempotentStart verifies the complete happy path:
// 1. Create new worktree with includes
// 2. Re-enter existing worktree (idempotent)
// 3. Verify git state is valid after both operations
// 4. Verify includes files match and were not re-copied
func TestJourney_IdempotentStart(t *testing.T) {
	projectRoot := setupGitRepoWithIncludes(t)
	slug := "my-feature"

	// ---- Step 1: Create new worktree ----
	stdout1, stderr1, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "step 1 create: expected exit code 0, stderr: %s", stderr1)
	assert.True(t, strings.Contains(stderr1, "created new worktree: "+slug),
		"step 1: expected 'created new worktree', got: %s", stderr1)
	assert.True(t, strings.Contains(stdout1, "worktree created at"),
		"step 1: expected worktree path in stdout, got: %s", stdout1)

	wtDir := worktreeDir(projectRoot, slug)

	// Verify worktree directory exists
	assert.DirExists(t, wtDir, "worktree directory should exist after creation")
	assertValidGitFile(t, wtDir)

	// Verify includes were copied
	srcFile := filepath.Join(projectRoot, "secret.txt")
	dstFile := filepath.Join(wtDir, "secret.txt")
	assertFileContentsMatch(t, srcFile, dstFile)

	// Record file mtime for later comparison
	infoAfterCreate, _ := os.Stat(dstFile)

	// ---- Step 2: Re-enter existing worktree ----
	_, stderr2, exitCode2 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode2, "step 2 re-entry: expected exit code 0, stderr: %s", stderr2)
	assert.True(t, strings.Contains(stderr2, "entering existing worktree: "+slug),
		"step 2: expected 'entering existing worktree', got: %s", stderr2)
	assert.False(t, strings.Contains(stderr2, "created new worktree"),
		"step 2: should NOT contain 'created new worktree'")

	// ---- Step 3: Verify git state ----
	listOutput := string(execGit(projectRoot, "worktree", "list"))
	slugCount := strings.Count(listOutput, slug)
	assert.Equal(t, 1, slugCount,
		"step 3: expected exactly one worktree entry, got:\n%s", listOutput)

	// Verify .git file is valid
	gitFilePath := filepath.Join(wtDir, ".git")
	gitContent, err := os.ReadFile(gitFilePath)
	assert.NoError(t, err, "step 3: should read .git file")
	assert.True(t, strings.Contains(string(gitContent), "gitdir"),
		"step 3: .git file should contain gitdir reference")

	// ---- Step 4: Verify includes files not re-copied ----
	assertFileContentsMatch(t, srcFile, dstFile)
	infoAfterReentry, _ := os.Stat(dstFile)
	assert.Equal(t, infoAfterCreate.ModTime(), infoAfterReentry.ModTime(),
		"step 4: includes file should not be modified during re-entry")

	// ---- Journey Invariant: stderr always distinguishes create vs enter ----
	assert.True(t, strings.Contains(stderr1, "created new worktree"),
		"invariant: first invocation stderr should contain 'created new worktree'")
	assert.True(t, strings.Contains(stderr2, "entering existing worktree"),
		"invariant: second invocation stderr should contain 'entering existing worktree'")

	// ---- Journey Invariant: worktree directory structure remains valid ----
	assertValidGitFile(t, wtDir)
}
