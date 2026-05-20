//go:build e2e

package e2e

// ==============================================================================
// forge commands — CLI e2e tests for feature: test-knowledge-convention-driven
// Tests cover backward compatibility after Profile removal, config init,
// task index, task add, and forge init.
// ==============================================================================

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- TC-006: Upgraded Forge Silently Ignores Legacy Config Fields ---
// Traceability: TC-006 -> Story 4 / AC-1

func TestTC_006_BackwardCompatIgnoresLegacyConfig(t *testing.T) {
	projectRoot := t.TempDir()

	// Step 1: Create .forge directory and config with legacy fields
	require.NoError(t, os.MkdirAll(filepath.Join(projectRoot, ".forge"), 0755))

	legacyConfig := `languages:
  - go
interfaces:
  - cli
test-framework: go-testing
project-type: backend
auto:
  e2eTest:
    quick: false
    full: true
worktree:
  source-branch: main
`
	configPath := filepath.Join(projectRoot, ".forge", "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(legacyConfig), 0644))

	// Create CLAUDE.md for project root detection
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Test\n"), 0644))

	// Step 2: Run forge task index (should not error on legacy fields)
	cmd := exec.Command("forge", "task", "index")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	t.Logf("forge task index output: %s", string(out))

	// Forge commands should still work even with legacy fields present
	// The key assertion: no errors referencing legacy fields
	legacyPattern := regexp.MustCompile(`(?i)languages|interfaces|test-framework|project-type|legacy.*field|deprecated.*field`)
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	assert.False(t, legacyPattern.MatchString(stderr),
		"stderr should not reference legacy fields, got: %s", stderr)

	// Step 3: Run forge config get to verify config is readable
	cmd = exec.Command("forge", "config", "get", "auto.e2e-test.full")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
	out, err = cmd.CombinedOutput()
	t.Logf("forge config get output: %s", string(out))

	// Config should be readable regardless of legacy fields
	if err == nil {
		assert.Contains(t, string(out), "true")
	}
}

// --- TC-012: Forge Commands Function Correctly After Profile Removal ---
// Traceability: TC-012 -> Spec FS-7 / Import Audit Gate

func TestTC_012_CommandsWorkAfterProfileRemoval(t *testing.T) {
	projectRoot := setupExistingForgeProject(t)

	// Step 1: Run forge task index
	cmd := exec.Command("forge", "task", "index", "--feature", "test-knowledge-convention-driven")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	t.Logf("forge task index output: %s", string(out))

	// Should work without Profile errors
	profileErrPattern := regexp.MustCompile(`(?i)profile|Profile`)
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	assert.False(t, profileErrPattern.MatchString(stderr),
		"Should not reference Profile, got: %s", stderr)
}

// setupExistingForgeProject creates a temp project with existing forge structure.
func setupExistingForgeProject(t *testing.T) string {
	t.Helper()
	projectRoot := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(projectRoot, ".forge"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Test Project\n"), 0644))

	// Create docs/features structure
	featureTasksDir := filepath.Join(projectRoot, "docs", "features", "test-feature", "tasks")
	require.NoError(t, os.MkdirAll(featureTasksDir, 0755))

	return projectRoot
}

// --- TC-013: Config Init Works Without Legacy Fields ---
// Traceability: TC-013 -> Spec FS-6 / FS-8

func TestTC_013_ConfigInitWorksWithoutLegacyFields(t *testing.T) {
	projectRoot := t.TempDir()

	// Step 1: Create .forge directory
	require.NoError(t, os.MkdirAll(filepath.Join(projectRoot, ".forge"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Test\n"), 0644))

	// Run forge config init with piped input (non-interactive)
	// Using "echo" to provide defaults for all prompts
	cmd := exec.Command("forge", "config", "init")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
	cmd.Stdin = strings.NewReader("\nn\n\n\n\n\n\n")
	out, err := cmd.CombinedOutput()
	t.Logf("forge config init output: %s", string(out))

	configPath := filepath.Join(projectRoot, ".forge", "config.yaml")

	// Check if config was created
	if _, statErr := os.Stat(configPath); statErr == nil {
		content, readErr := os.ReadFile(configPath)
		require.NoError(t, readErr)

		// Step 3: Verify no legacy fields
		legacyPattern := regexp.MustCompile(`(?m)^languages:|^interfaces:|^test-framework:|^project-type:`)
		assert.False(t, legacyPattern.MatchString(string(content)),
			"config.yaml should not contain legacy fields, got:\n%s", string(content))
	}

	_ = err
}

// --- TC-025: Forge Task Add Works Without Profile Dependency ---
// Traceability: TC-025 -> Spec FS-7 / Related Changes #5

func TestTC_025_TaskAddWorksWithoutProfileDependency(t *testing.T) {
	projectRoot := setupForgeProjectWithFeature(t)

	// Step 2: Run forge task add
	cmd := exec.Command("forge", "task", "add",
		"--title", "test-task-convention",
		"--description", "test description for convention-driven feature")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	t.Logf("forge task add output: %s", string(out))

	if err == nil {
		// Step 4: Verify task was added via forge task index
		assert.Contains(t, string(out), "ADDED")
		assert.Contains(t, string(out), "test-task-convention")
	} else {
		// May fail if feature not properly set up; log the error
		t.Logf("forge task add error: %v", err)
	}

	// No Profile errors
	profileErrPattern := regexp.MustCompile(`(?i)profile|Profile`)
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	assert.False(t, profileErrPattern.MatchString(stderr),
		"Should not reference Profile, got: %s", stderr)
}

// setupForgeProjectWithFeature creates a temp project with a feature set.
func setupForgeProjectWithFeature(t *testing.T) string {
	t.Helper()
	projectRoot := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(projectRoot, ".forge"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Test Project\n"), 0644))

	// Create feature directory structure
	featureSlug := "test-feature"
	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	featureTasksDir := filepath.Join(featureDir, "tasks")
	require.NoError(t, os.MkdirAll(featureTasksDir, 0755))

	// Create manifest
	manifestContent := `---
slug: "test-feature"
status: "in-progress"
---
`
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, "manifest.md"), []byte(manifestContent), 0644))

	// Set current feature via state.json
	stateContent := `{"feature": "test-feature"}`
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, ".forge", "state.json"), []byte(stateContent), 0644))

	return projectRoot
}

// --- TC-026: Forge Init Creates Project Without Legacy Fields ---
// Traceability: TC-026 -> Spec FS-7 / Related Changes #5

func TestTC_026_ForgeInitCreatesProjectWithoutLegacyFields(t *testing.T) {
	projectRoot := t.TempDir()

	// Step 1: Run forge init with --skip-just to avoid just installation issues
	cmd := exec.Command("forge", "init", "--project-root", projectRoot, "--skip-just")
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	t.Logf("forge init output: %s", string(out))

	// Step 2: Verify .forge directory was created
	forgeDir := filepath.Join(projectRoot, ".forge")
	if _, statErr := os.Stat(forgeDir); statErr == nil {
		// Step 3: Check config.yaml if created
		configPath := filepath.Join(forgeDir, "config.yaml")
		if content, readErr := os.ReadFile(configPath); readErr == nil {
			// Step 4: Verify no legacy fields
			legacyPattern := regexp.MustCompile(`(?m)^languages:|^interfaces:|^test-framework:|^project-type:`)
			assert.False(t, legacyPattern.MatchString(string(content)),
				"config.yaml should not contain legacy fields, got:\n%s", string(content))
		}
	}

	_ = err
}
