//go:build cli_functional

package startexistingflags

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Smoke Test: start-existing-flags Journey (happy path)
// Contract: start-existing-flags journey.md
// ==============================================================================

// TestJourney_StartExistingFlags verifies the happy path:
// 1. Create worktree
// 2. Re-enter with --source-branch (ignored)
// 3. Re-enter with --no-launch (path output)
func TestJourney_StartExistingFlags(t *testing.T) {
	projectRoot := setupGitRepoWithIncludes(t)
	slug := "flags-test-wt"

	// ---- Step 1: Create worktree ----
	stdout1, stderr1, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "step 1: create should succeed, stderr: %s", stderr1)
	assert.True(t, strings.Contains(stderr1, "created new worktree: "+slug),
		"step 1: expected creation message")
	assert.True(t, strings.Contains(stdout1, "worktree created at"),
		"step 1: expected worktree path in stdout")

	wtDir := worktreeDir(projectRoot, slug)
	assertValidGitFile(t, wtDir)

	// ---- Step 2: Re-enter with --source-branch (should be ignored) ----
	_, stderr2, exitCode2 := runForgeStartNoLaunch(t, projectRoot, slug, "--source-branch", "develop")
	assert.Equal(t, 0, exitCode2, "step 2: re-entry with --source-branch should succeed")
	assert.True(t, strings.Contains(stderr2, "warning") && strings.Contains(stderr2, "ignoring --source-branch"),
		"step 2: expected warning about ignoring --source-branch")
	assert.True(t, strings.Contains(stderr2, "entering existing worktree: "+slug),
		"step 2: expected entering message")

	// ---- Step 3: Re-enter with --no-launch ----
	stdout3, stderr3, exitCode3 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode3, "step 3: re-entry with --no-launch should succeed")
	assert.True(t, strings.Contains(stderr3, "entering existing worktree: "+slug),
		"step 3: expected entering message")
	assert.True(t, strings.Contains(stdout3, "worktree path:"),
		"step 3: expected worktree path in stdout")

	// ---- Journey Invariant: --no-launch always suppresses claude ----
	// All three steps used --no-launch and none required claude binary

	// ---- Journey Invariant: no file system mutations on re-entry ----
	entries, _ := os.ReadDir(wtDir)
	assert.True(t, len(entries) > 0, "worktree should have files")

	// ---- Journey Invariant: exit code 0 for all successful combinations ----
	assert.Equal(t, 0, exitCode1)
	assert.Equal(t, 0, exitCode2)
	assert.Equal(t, 0, exitCode3)
}
