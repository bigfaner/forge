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
// Step 4: Verify includes files were copied only once
// Contract: idempotent-start/step-4-verify-includes.md
// ==============================================================================

// Traceability: Step 4 / Outcome "files-match"
// Verifies that includes files match originals and were not modified during re-entry.
func TestStep4_FilesMatch_IncludesCopiedOnce(t *testing.T) {
	projectRoot := setupGitRepoWithIncludes(t)
	slug := "my-feature"

	// Phase 1: Create worktree (copies includes)
	_, stderr1, exitCode1 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode1, "create should succeed, stderr: %s", stderr1)
	assert.True(t, strings.Contains(stderr1, "created new worktree: "+slug))

	// Record file mtime after creation
	wtDir := worktreeDir(projectRoot, slug)
	dstFile := filepath.Join(wtDir, "secret.txt")
	infoAfterCreate, _ := os.Stat(dstFile)

	// Phase 2: Re-enter worktree (should NOT re-copy)
	_, stderr2, exitCode2 := forgeStartNoLaunch(t, projectRoot, slug)
	assert.Equal(t, 0, exitCode2, "re-entry should succeed, stderr: %s", stderr2)

	// Assert: file contents match originals exactly
	srcFile := filepath.Join(projectRoot, "secret.txt")
	assertFileContentsMatch(t, srcFile, dstFile)

	// Assert: file mtime unchanged after re-entry (behavioral: deep - proves no re-copy)
	infoAfterReentry, _ := os.Stat(dstFile)
	assert.Equal(t, infoAfterCreate.ModTime(), infoAfterReentry.ModTime(),
		"includes file should not be modified during re-entry")
}

// Traceability: Step 4 / Outcome "no-includes-config"
// Verifies worktree creation succeeds when no includes are configured.
func TestStep4_NoIncludesConfig_WorktreeCreatedWithoutFiles(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t) // No includes
	slug := "no-includes-feature"

	// Execute: forge worktree start no-includes-feature --no-launch
	_, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code 0
	assert.Equal(t, 0, exitCode, "expected exit code 0, stderr: %s", stderr)

	// Assert: stderr contains creation confirmation
	assert.True(t, strings.Contains(stderr, "created new worktree: "+slug),
		"expected creation message, got: %s", stderr)

	// Assert: worktree exists with valid .git file
	wtDir := worktreeDir(projectRoot, slug)
	assertValidGitFile(t, wtDir)

	// Assert: no extra files in worktree beyond git defaults
	entries, err := os.ReadDir(wtDir)
	assert.NoError(t, err)
	for _, entry := range entries {
		// Only .git file and git default files should exist
		assert.True(t, entry.Name() == ".git" || entry.Name() == "README.md",
			"unexpected file in worktree: %s", entry.Name())
	}
}

// Traceability: Step 4 / Outcome "migrated-config"
// Verifies worktree creation with new 'includes' key (not copy-files).
func TestStep4_MigratedConfig_IncludesCopied(t *testing.T) {
	projectRoot := setupGitRepoWithForge(t)
	slug := "migrated-feature"

	// Create include files
	notesFile := filepath.Join(projectRoot, "notes.txt")
	if err := os.WriteFile(notesFile, []byte("notes content\n"), 0644); err != nil {
		t.Fatalf("failed to write notes.txt: %v", err)
	}
	runGit(t, projectRoot, "add", "notes.txt")
	runGit(t, projectRoot, "commit", "-m", "add notes")

	// Write config with includes key
	writeForgeConfig(t, projectRoot, "notes.txt")

	// Execute: forge worktree start migrated-feature --no-launch
	_, stderr, exitCode := forgeStartNoLaunch(t, projectRoot, slug)

	// Assert: exit code 0
	assert.Equal(t, 0, exitCode, "expected exit code 0, stderr: %s", stderr)

	// Assert: stderr contains creation confirmation
	assert.True(t, strings.Contains(stderr, "created new worktree: "+slug),
		"expected creation message, got: %s", stderr)

	// Assert: includes files were copied
	wtDir := worktreeDir(projectRoot, slug)
	dstFile := filepath.Join(wtDir, "notes.txt")
	assertFileContentsMatch(t, notesFile, dstFile)

	// Assert: no error about missing copy-files key
	assert.False(t, strings.Contains(stderr, "copy-files"),
		"should not mention copy-files when using includes key, got: %s", stderr)
}
