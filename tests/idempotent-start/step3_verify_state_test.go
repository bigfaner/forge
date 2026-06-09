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
// Step 3: Verify worktree state after both invocations
// Contract: idempotent-start/step-3-verify-worktree-state.md
// ==============================================================================

// Traceability: Step 3 / Outcome "single-entry"
// Verifies that git worktree list shows exactly one entry after create + re-entry.
func TestStep3_SingleEntry_AfterCreateAndReentry(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "my-feature"

	// Phase 1: Create worktree
	_, stderr1, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed, stderr: %s", stderr1)

	// Phase 2: Re-enter worktree
	_, stderr2, exitCode2 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode2, "re-entry should succeed, stderr: %s", stderr2)

	// Assert: git worktree list shows exactly one entry for slug
	listOutput := string(execGit(projectRoot, "worktree", "list"))
	slugLines := 0
	for _, line := range strings.Split(listOutput, "\n") {
		if strings.Contains(line, slug) {
			slugLines++
		}
	}
	assert.Equal(t, 1, slugLines,
		"expected exactly one worktree entry for %s, got:\n%s", slug, listOutput)

	// Assert: entry has a valid path (behavioral: proves git worktree is registered)
	assert.True(t, strings.Contains(listOutput, ".forge/worktrees/"+slug),
		"worktree list should contain worktree path, got:\n%s", listOutput)
}

// Traceability: Step 3 / Outcome "git-file-valid"
// Verifies that .git file is a file (not directory) and contains valid gitdir reference.
func TestStep3_GitFileValid_AfterCreateAndReentry(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "my-feature"

	// Phase 1: Create worktree
	_, stderr1, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed, stderr: %s", stderr1)

	// Phase 2: Re-enter worktree
	_, stderr2, exitCode2 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode2, "re-entry should succeed, stderr: %s", stderr2)

	// Assert: .git is a file (not a directory)
	wtDir := worktreeDir(projectRoot, slug)
	gitFilePath := filepath.Join(wtDir, ".git")
	info, err := os.Stat(gitFilePath)
	assert.NoError(t, err, ".git should exist")
	assert.False(t, info.IsDir(), ".git should be a file, not a directory")

	// Assert: .git file contains gitdir reference (behavioral: deep assertion - cross-entity reference)
	content, err := os.ReadFile(gitFilePath)
	assert.NoError(t, err, "should be able to read .git file")
	assert.True(t, strings.Contains(string(content), "gitdir"),
		".git file should contain gitdir reference, got: %s", string(content))

	// Assert: gitdir reference points to a valid path
	gitdirPath := strings.TrimPrefix(strings.TrimSpace(string(content)), "gitdir: ")
	gitdirPath = filepath.Clean(gitdirPath)
	// The gitdir should reference a path under .git/worktrees/ in the main repo
	assert.True(t, strings.Contains(gitdirPath, "worktrees"),
		"gitdir should reference worktrees directory, got: %s", gitdirPath)
}
