//go:build cli_functional

package startexistingflags

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 1: Start with --source-branch on existing worktree
// Contract: start-existing-flags/step-1-source-branch-on-existing.md
// ==============================================================================

// Traceability: Step 1 / Outcome "success-ignored-flag"
// Verifies --source-branch is ignored with a warning when entering existing worktree.
func TestStep1_SourceBranchIgnored_OnExistingWorktree(t *testing.T) {
	projectRoot := setupGitRepoWithIncludes(t)
	slug := "my-feature"

	// Phase 1: Create worktree
	_, stderr1, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")
	assert.True(t, strings.Contains(stderr1, "created new worktree: "+slug))

	// Phase 2: Re-enter with --source-branch develop
	stdout2, stderr2, exitCode2 := runForgeStartNoLaunch(t, projectRoot, slug, "--source-branch", "develop")
	assert.Equal(t, 0, exitCode2, "re-entry with --source-branch should succeed, stderr: %s", stderr2)

	// Assert: stderr contains warning about ignoring --source-branch
	assert.True(t, strings.Contains(stderr2, "warning") && strings.Contains(stderr2, "ignoring --source-branch"),
		"expected warning about ignoring --source-branch, got: %s", stderr2)

	// Assert: stderr contains entering message
	assert.True(t, strings.Contains(stderr2, "entering existing worktree: "+slug),
		"expected entering message, got: %s", stderr2)

	// Assert: stdout contains worktree path
	assert.True(t, strings.Contains(stdout2, "worktree path:"),
		"expected worktree path in stdout, got: %s", stdout2)
}

// Traceability: Step 1 / Outcome "diverged-branch"
// Verifies that existing worktree with diverged branch is entered without rebase when --source-branch is given.
func TestStep1_DivergedBranch_EntersWithoutRebase(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "diverged-wt"

	// Phase 1: Create worktree
	_, stderr1, exitCode1 := runForgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed")
	assert.True(t, strings.Contains(stderr1, "created new worktree"))

	// Phase 2: Diverge the worktree branch
	wtDir := worktreeDir(projectRoot, slug)
	divergeFile := wtDir + "/diverge.txt"
	writeFile(t, divergeFile, "diverged\n")
	runGitIn(t, wtDir, "add", "diverge.txt")
	runGitIn(t, wtDir, "commit", "-m", "diverge")

	// Phase 3: Re-enter with --source-branch
	stdout3, stderr3, exitCode3 := runForgeStartNoLaunch(t, projectRoot, slug, "--source-branch", "develop")
	assert.Equal(t, 0, exitCode3, "re-entry should succeed, stderr: %s", stderr3)

	// Assert: stderr contains warning about ignoring --source-branch
	assert.True(t, strings.Contains(stderr3, "warning") && strings.Contains(stderr3, "ignoring --source-branch"),
		"expected warning, got: %s", stderr3)

	// Assert: stderr contains entering message
	assert.True(t, strings.Contains(stderr3, "entering existing worktree: "+slug),
		"expected entering message, got: %s", stderr3)

	// Assert: diverge file still exists (behavioral: no branch modification)
	assert.FileExists(t, divergeFile, "diverged commit should be preserved")

	_ = stdout3
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := writeFileHelper(path, content); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func writeFileHelper(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
