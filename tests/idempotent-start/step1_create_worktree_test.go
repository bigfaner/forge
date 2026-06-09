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
// Step 1: Start with non-existent worktree (create path)
// Contract: idempotent-start/step-1-create-worktree.md
// ==============================================================================

// Traceability: Step 1 / Outcome "success"
// Verifies that forge worktree start creates a new worktree on first invocation.
func TestStep1_Success_CreatesNewWorktree(t *testing.T) {
	projectRoot := setupGitRepoWithIncludes(t)
	slug := "my-feature"

	// Execute: forge worktree start my-feature --no-launch
	stdout, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code 0
	assert.Equal(t, 0, exitCode, "expected exit code 0, stderr: %s", stderr)

	// Assert: stderr contains creation confirmation
	assert.True(t, strings.Contains(stderr, "created new worktree: "+slug),
		"expected stderr to contain 'created new worktree: %s', got stderr: %s", slug, stderr)

	// Assert: worktree directory exists
	wtDir := worktreeDir(projectRoot, slug)
	info, err := os.Stat(wtDir)
	assert.NoError(t, err, "worktree directory should exist")
	assert.True(t, info.IsDir(), "worktree path should be a directory")

	// Assert: .git file is present and is a file (not directory)
	assertValidGitFile(t, wtDir)

	// Assert: includes files were copied
	srcFile := filepath.Join(projectRoot, "secret.txt")
	dstFile := filepath.Join(wtDir, "secret.txt")
	assertFileContentsMatch(t, srcFile, dstFile)

	// Assert: stdout contains worktree path (behavioral: proves path resolution)
	assert.True(t, strings.Contains(stdout, "worktree created at"),
		"expected stdout to contain 'worktree created at', got: %s", stdout)

	// Assert: git branch was created (behavioral: proves git state transition)
	cmd := execGit(projectRoot, "branch", "--list", slug)
	assert.True(t, strings.Contains(string(cmd), slug),
		"expected git branch %s to exist, got: %s", slug, string(cmd))
}

// Traceability: Step 1 / Outcome "slug-not-found"
// Verifies that forge worktree start without a slug argument returns an error.
func TestStep1_SlugNotFound_ReturnsError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)

	// Execute: forge worktree start (no arguments)
	stdout, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, "")

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code when no slug provided, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains error about slug
	assert.True(t, strings.Contains(stderr, "slug") || strings.Contains(stderr, "Slug"),
		"expected stderr to mention slug requirement, got: %s", stderr)

	// Assert: no worktree was created
	entries, _ := os.ReadDir(filepath.Join(projectRoot, ".forge", "worktrees"))
	assert.Equal(t, 0, len(entries), "no worktrees should be created")
}

// Traceability: Step 1 / Outcome "not-git-repository"
// Verifies that forge worktree start in a non-git directory returns an error.
func TestStep1_NotGitRepository_ReturnsError(t *testing.T) {
	// Setup: temp dir without git init
	dir := t.TempDir()
	// Create .forge/config.yaml so project root is found
	forgeDir := filepath.Join(dir, ".forge")
	_ = os.MkdirAll(forgeDir, 0755)
	_ = os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("version: \"1\"\nsurfaces: cli\n"), 0644)
	// Create go.mod for project root detection
	_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n"), 0644)

	// Execute: forge worktree start my-feature --no-launch
	stdout, stderr, exitCode := forgeStartNoLaunch(t, dir, "my-feature")

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code for non-git directory, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains git repository error
	assert.True(t, strings.Contains(stderr, "git") || strings.Contains(stderr, "repository"),
		"expected stderr to mention git repository, got: %s", stderr)
}

// Traceability: Step 1 / Outcome "claude-not-found"
// Verifies that forge worktree start without --no-launch fails when claude is not in PATH.
func TestStep1_ClaudeNotFound_ReturnsError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "claude-test-feature"

	// Execute: forge worktree start my-feature (without --no-launch, claude not in PATH)
	// We set PATH to empty to simulate claude not being available
	args := []string{"worktree", "start", slug}
	cmd := createForgeCommand(args...)
	cmd.Env = []string{
		"CLAUDE_PROJECT_DIR=" + projectRoot,
		"PATH=/nonexistent", // claude binary not in PATH
		"HOME=" + os.Getenv("HOME"),
	}
	out, _ := cmd.CombinedOutput()
	output := string(out)

	// Assert: exit code is non-zero
	assert.False(t, cmd.ProcessState.Success(),
		"expected failure when claude not in PATH, output: %s", output)

	// Assert: output mentions claude binary not found
	assert.True(t, strings.Contains(output, "claude") || strings.Contains(output, "Claude"),
		"expected output to mention claude binary, got: %s", output)
}

// Traceability: Step 1 / Outcome "source-branch-not-found"
// Verifies that forge worktree start with non-existent --source-branch returns an error.
func TestStep1_SourceBranchNotFound_ReturnsError(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "branch-test-feature"

	// Execute: forge worktree start my-feature --source-branch nonexistent-branch --no-launch
	stdout, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug, "--source-branch", "nonexistent-branch")

	// Assert: exit code is non-zero
	assert.NotEqual(t, 0, exitCode,
		"expected non-zero exit code for nonexistent source branch, stdout: %s, stderr: %s", stdout, stderr)

	// Assert: stderr contains source branch not found error
	assert.True(t, strings.Contains(stderr, "source branch") && strings.Contains(stderr, "not found"),
		"expected stderr to mention source branch not found, got: %s", stderr)

	// Assert: stderr contains hint
	assert.True(t, strings.Contains(stderr, "hint") || strings.Contains(stderr, "verify"),
		"expected stderr to contain a hint, got: %s", stderr)

	// Assert: no worktree was created
	wtDir := worktreeDir(projectRoot, slug)
	_, err := os.Stat(wtDir)
	assert.True(t, os.IsNotExist(err), "no worktree should be created for invalid source branch")
}
