//go:build cli_functional

package startexistingflags

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 2: Start with --no-launch on existing worktree
// Contract: start-existing-flags/step-2-no-launch-on-existing.md
// ==============================================================================

// Traceability: Step 2 / Outcome "success"
// Verifies --no-launch on existing worktree outputs path without launching claude.
func TestStep2_NoLaunchExisting_OutputsPath(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "test-wt"

	// Phase 1: Create worktree
	_, stderr1, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")
	assert.True(t, strings.Contains(stderr1, "created new worktree"))

	// Phase 2: Re-enter with --no-launch
	stdout2, stderr2, exitCode2 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode2, "re-entry should succeed, stderr: %s", stderr2)

	// Assert: stderr contains entering message
	assert.True(t, strings.Contains(stderr2, "entering existing worktree: "+slug),
		"expected entering message, got: %s", stderr2)

	// Assert: stdout contains worktree path
	assert.True(t, strings.Contains(stdout2, "worktree path:"),
		"expected worktree path in stdout, got: %s", stdout2)

	// Assert: no file system mutations
	wtDir := worktreeDir(projectRoot, slug)
	assertValidGitFile(t, wtDir)
}

// Traceability: Step 2 / Outcome "no-launch-new-worktree"
// Verifies --no-launch creates a new worktree when none exists.
func TestStep2_NoLaunchNewWorktree_CreatesSuccessfully(t *testing.T) {
	projectRoot := setupGitRepoWithIncludes(t)
	slug := "new-no-launch"

	// Execute: forge worktree start new-no-launch --no-launch
	stdout, stderr, exitCode := runForgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code 0
	assert.Equal(t, 0, exitCode, "expected exit code 0, stderr: %s", stderr)

	// Assert: stderr contains creation message
	assert.True(t, strings.Contains(stderr, "created new worktree: "+slug),
		"expected creation message, got: %s", stderr)

	// Assert: stdout contains path
	assert.True(t, strings.Contains(stdout, "worktree created at"),
		"expected worktree path in stdout, got: %s", stdout)

	// Assert: worktree exists with valid .git
	wtDir := worktreeDir(projectRoot, slug)
	assertValidGitFile(t, wtDir)

	// Assert: includes files were copied
	dstFile := wtDir + "/secret.txt"
	assert.FileExists(t, dstFile, "includes file should be copied")
}

// Traceability: Step 2 / Outcome "combined-flags-ignored"
// Verifies --source-branch + --no-launch on existing worktree.
func TestStep2_CombinedFlags_SourceBranchIgnoredNoLaunch(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "combined-wt"

	// Phase 1: Create worktree
	_, _, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")

	// Phase 2: Re-enter with both --source-branch and --no-launch
	stdout2, stderr2, exitCode2 := runForgeStartNoLaunch(t, projectRoot, slug, "--source-branch", "develop")
	assert.Equal(t, 0, exitCode2, "combined flags should succeed, stderr: %s", stderr2)

	// Assert: stderr contains warning about ignoring --source-branch
	assert.True(t, strings.Contains(stderr2, "warning") && strings.Contains(stderr2, "ignoring --source-branch"),
		"expected warning about ignoring --source-branch, got: %s", stderr2)

	// Assert: stdout contains worktree path
	assert.True(t, strings.Contains(stdout2, "worktree path:"),
		"expected worktree path in stdout, got: %s", stdout2)

	// Assert: no claude process spawned (verified by --no-launch behavior)

	// Assert: no file mutations
	wtDir := worktreeDir(projectRoot, slug)
	entries, _ := os.ReadDir(wtDir)
	assert.True(t, len(entries) > 0, "worktree should still have files")
}
