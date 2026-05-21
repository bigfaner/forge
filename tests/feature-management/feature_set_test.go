//go:build e2e

package featuremanagement

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Feature set command tests — Journey: feature-management
// Tests cover: forge feature set subcommand, GetCurrentFeature() priority chain,
// and verbose flag output.
// ==============================================================================

// forgeState mirrors the ForgeState struct from feature package.
type forgeState struct {
	Feature      string `json:"feature"`
	AllCompleted bool   `json:"allCompleted"`
}

// setupProjectDir creates a temp directory acting as a clean project root.
// Returns the temp dir path.
func setupProjectDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// Create .forge dir so CLAUDE_PROJECT_DIR points at a valid project root.
	if err := os.MkdirAll(filepath.Join(dir, ".forge"), 0755); err != nil {
		t.Fatalf("failed to create .forge dir: %v", err)
	}
	return dir
}

// forgeFeatureSet runs "forge feature set <slug>" with CLAUDE_PROJECT_DIR set.
// Returns combined output and exit code. Does NOT fatalf on failure.
func forgeFeatureSet(t *testing.T, projectRoot, slug string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, "feature", "set", slug)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// forgeFeature runs "forge feature" (query) with CLAUDE_PROJECT_DIR set.
// Returns combined output and exit code. Does NOT fatalf on failure.
func forgeFeature(t *testing.T, projectRoot string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, "feature")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// forgeFeatureVerbose runs "forge feature -v" with CLAUDE_PROJECT_DIR set.
// Returns combined output and exit code. Does NOT fatalf on failure.
func forgeFeatureVerbose(t *testing.T, projectRoot string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, "feature", "-v")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// forgeFeaturePositional runs "forge feature <slug>" (positional arg, legacy) with CLAUDE_PROJECT_DIR set.
// Returns combined output and exit code.
func forgeFeaturePositional(t *testing.T, projectRoot, slug string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, "feature", slug)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}

// readForgeState reads .forge/state.json from projectRoot. Returns nil if not found.
func readForgeState(t *testing.T, projectRoot string) *forgeState {
	t.Helper()
	statePath := filepath.Join(projectRoot, ".forge", "state.json")
	data, err := os.ReadFile(statePath)
	if err != nil {
		return nil
	}
	var state forgeState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil
	}
	return &state
}

// writeForgeState writes .forge/state.json with the given feature slug and allCompleted.
func writeForgeState(t *testing.T, projectRoot, featureSlug string, allCompleted bool) {
	t.Helper()
	statePath := filepath.Join(projectRoot, ".forge", "state.json")
	if err := os.MkdirAll(filepath.Dir(statePath), 0755); err != nil {
		t.Fatalf("failed to create .forge dir: %v", err)
	}
	state := forgeState{
		Feature:      featureSlug,
		AllCompleted: allCompleted,
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("failed to marshal state: %v", err)
	}
	if err := os.WriteFile(statePath, data, 0644); err != nil {
		t.Fatalf("failed to write state.json: %v", err)
	}
}

// ensureFeatureDir creates the minimal feature directory structure,
// including a tasks/index.json so that getFeatureFromFeaturesDir can detect it.
func ensureFeatureDir(t *testing.T, projectRoot, slug string) {
	t.Helper()
	featureDir := filepath.Join(projectRoot, "docs", "features", slug)
	subdirs := []string{
		"prd",
		"design",
		"ui",
		"tasks",
		"tasks/process",
		"tasks/records",
	}
	for _, sub := range subdirs {
		if err := os.MkdirAll(filepath.Join(featureDir, sub), 0755); err != nil {
			t.Fatalf("failed to create feature subdir %s: %v", sub, err)
		}
	}
	// Write a minimal tasks/index.json so the features-dir scanner recognizes this directory.
	indexPath := filepath.Join(featureDir, "tasks", "index.json")
	if err := os.WriteFile(indexPath, []byte("{}"), 0644); err != nil {
		t.Fatalf("failed to write index.json: %v", err)
	}
}

// outputContains checks that the raw CLI output contains the expected substring.
func outputContains(t *testing.T, raw, expected string) {
	t.Helper()
	assert.True(t, strings.Contains(raw, expected),
		"expected output to contain %q, got:\n%s", expected, raw)
}

// outputNotContains checks that the raw CLI output does NOT contain the substring.
func outputNotContains(t *testing.T, raw, expected string) {
	t.Helper()
	assert.False(t, strings.Contains(raw, expected),
		"expected output NOT to contain %q, got:\n%s", expected, raw)
}

// =============================================================================
// Task 1: forge feature set subcommand
// =============================================================================

// Traceability: TC-001 -> Proposal SC #1,#2 / Task 1 AC #1,#2
func TestTC_001_SetFeatureCreatesDirectoryAndState(t *testing.T) {
	projectRoot := setupProjectDir(t)

	out, exitCode := forgeFeatureSet(t, projectRoot, "my-feature")
	assert.Equal(t, 0, exitCode, "forge feature set should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: my-feature")

	// Verify .forge/state.json exists with correct values
	state := readForgeState(t, projectRoot)
	assert.NotNil(t, state, ".forge/state.json should exist")
	assert.Equal(t, "my-feature", state.Feature)
	assert.False(t, state.AllCompleted, "allCompleted should be false")

	// Verify feature directory structure exists
	featureDir := filepath.Join(projectRoot, "docs", "features", "my-feature")
	assert.DirExists(t, featureDir, "feature directory should exist")
	assert.DirExists(t, filepath.Join(featureDir, "tasks"), "tasks/ subdir should exist")
	assert.DirExists(t, filepath.Join(featureDir, "tasks", "process"), "tasks/process/ subdir should exist")
}

// Traceability: TC-002 -> Proposal SC #3 / Task 1 AC #3
func TestTC_002_SetFeatureWithEmptySlugReturnsError(t *testing.T) {
	projectRoot := setupProjectDir(t)

	out, exitCode := forgeFeatureSet(t, projectRoot, "")
	assert.NotEqual(t, 0, exitCode, "forge feature set '' should fail, output: %s", out)

	// Verify .forge/state.json does NOT exist
	state := readForgeState(t, projectRoot)
	assert.Nil(t, state, ".forge/state.json should NOT exist when slug is empty")
}

// Traceability: TC-003 -> Task 1 AC #4
func TestTC_003_SetFeaturePrintsSlugToStdout(t *testing.T) {
	projectRoot := setupProjectDir(t)

	out, exitCode := forgeFeatureSet(t, projectRoot, "test-slug")
	assert.Equal(t, 0, exitCode, "forge feature set should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: test-slug")
}

// Traceability: TC-004 -> Task 1 AC #5
func TestTC_004_PositionalArgBackwardCompatibility(t *testing.T) {
	projectRoot := setupProjectDir(t)

	out, exitCode := forgeFeaturePositional(t, projectRoot, "legacy-feature")
	assert.Equal(t, 0, exitCode, "forge feature <slug> should succeed, output: %s", out)

	// Verify directory exists
	featureDir := filepath.Join(projectRoot, "docs", "features", "legacy-feature")
	assert.DirExists(t, featureDir, "feature directory should exist")

	// Verify state.json does NOT exist (positional arg does not write state)
	state := readForgeState(t, projectRoot)
	assert.Nil(t, state, ".forge/state.json should NOT exist for positional arg (old behavior)")
}

// Traceability: TC-005 -> Task 1 Hard Rules (validate slug is non-empty)
func TestTC_005_SetFeatureWithWhitespaceOnlySlugReturnsError(t *testing.T) {
	projectRoot := setupProjectDir(t)

	out, exitCode := forgeFeatureSet(t, projectRoot, "   ")
	assert.NotEqual(t, 0, exitCode, "forge feature set '   ' should fail, output: %s", out)

	// Verify no side effects on filesystem
	state := readForgeState(t, projectRoot)
	assert.Nil(t, state, ".forge/state.json should NOT exist when slug is whitespace-only")
}

// Traceability: TC-006 -> Proposal Key Scenarios (happy path)
func TestTC_006_SetFeatureIdempotentOnRepeatedCalls(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// First call
	out1, exitCode1 := forgeFeatureSet(t, projectRoot, "my-feature")
	assert.Equal(t, 0, exitCode1, "first forge feature set should succeed, output: %s", out1)

	// Second call
	out2, exitCode2 := forgeFeatureSet(t, projectRoot, "my-feature")
	assert.Equal(t, 0, exitCode2, "second forge feature set should succeed, output: %s", out2)

	// Verify state.json remains valid
	state := readForgeState(t, projectRoot)
	assert.NotNil(t, state, ".forge/state.json should exist")
	assert.Equal(t, "my-feature", state.Feature)
}

// Traceability: TC-007 -> Proposal Key Scenarios (worktree mismatch)
func TestTC_007_SetFeatureOverwritesPreviousFeatureInState(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Set first feature
	out1, exitCode1 := forgeFeatureSet(t, projectRoot, "feature-a")
	assert.Equal(t, 0, exitCode1, "first forge feature set should succeed, output: %s", out1)

	// Set second feature
	out2, exitCode2 := forgeFeatureSet(t, projectRoot, "feature-b")
	assert.Equal(t, 0, exitCode2, "second forge feature set should succeed, output: %s", out2)

	// Verify state.json updated to feature-b
	state := readForgeState(t, projectRoot)
	assert.NotNil(t, state, ".forge/state.json should exist")
	assert.Equal(t, "feature-b", state.Feature)

	// Verify feature-b directory exists
	featureBDir := filepath.Join(projectRoot, "docs", "features", "feature-b")
	assert.DirExists(t, featureBDir, "feature-b directory should exist")

	// Verify old feature-a directory still exists (not removed)
	featureADir := filepath.Join(projectRoot, "docs", "features", "feature-a")
	assert.DirExists(t, featureADir, "feature-a directory should remain intact")
}

// =============================================================================
// Task 2: GetCurrentFeature() priority chain
// =============================================================================

// Traceability: TC-008 -> Proposal SC #4 / Task 2 AC #1
func TestTC_008_GetCurrentFeatureReturnsStateJsonFeatureWhenPresent(t *testing.T) {
	projectRoot := setupProjectDir(t)
	ensureFeatureDir(t, projectRoot, "explicit-feature")
	writeForgeState(t, projectRoot, "explicit-feature", false)

	out, exitCode := forgeFeature(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: explicit-feature")
}

// Traceability: TC-009 -> Proposal SC #6 / Task 2 AC #2
func TestTC_009_GetCurrentFeatureFallsBackWhenStateJsonAbsent(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Create a single feature directory (fallback to features-dir)
	ensureFeatureDir(t, projectRoot, "lone-feature")

	// Ensure no state.json
	_ = os.Remove(filepath.Join(projectRoot, ".forge", "state.json"))

	out, exitCode := forgeFeature(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: lone-feature")
}

// Traceability: TC-010 -> Task 2 AC #3
func TestTC_010_GetCurrentFeatureWithSourceReturnsCorrectSourceType(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Scenario 1: state.json set with existing feature dir
	ensureFeatureDir(t, projectRoot, "state-feature")
	writeForgeState(t, projectRoot, "state-feature", false)

	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v (state.json) should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: state-feature (from: state.json)")
}

// Traceability: TC-011 -> Task 2 Hard Rules (feature directory doesn't exist)
func TestTC_011_StateJsonWithNonexistentDirFallsThrough(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Write state.json pointing to a feature with NO directory
	writeForgeState(t, projectRoot, "ghost-feature", false)

	// Create a different feature directory that can be found by features-dir scan
	ensureFeatureDir(t, projectRoot, "fallback-feature")

	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v should succeed, output: %s", out)

	// Should NOT show ghost-feature
	outputNotContains(t, out, "ghost-feature")

	// Should show fallback-feature from features-dir scan
	outputContains(t, out, "FEATURE: fallback-feature")
}

// Traceability: TC-012 -> Task 2 Hard Rules (corrupt state.json silently ignored)
func TestTC_012_CorruptStateJsonFallsThroughSilently(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Write corrupt state.json
	statePath := filepath.Join(projectRoot, ".forge", "state.json")
	if err := os.WriteFile(statePath, []byte("not json at all"), 0644); err != nil {
		t.Fatalf("failed to write corrupt state.json: %v", err)
	}

	// Create a feature directory that can be found by features-dir scan
	ensureFeatureDir(t, projectRoot, "recoverable-feature")

	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v should succeed despite corrupt state, output: %s", out)

	// Should NOT show any error about corrupt state.json
	outputNotContains(t, out, "corrupt")
	outputNotContains(t, out, "invalid")

	// Should show recoverable-feature from features-dir scan
	outputContains(t, out, "FEATURE: recoverable-feature")
}

// Traceability: TC-013 -> Proposal Key Scenarios (worktree mismatch)
func TestTC_013_StateJsonTakesPriorityOverGitWorktree(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Write state.json with one feature
	writeForgeState(t, projectRoot, "oauth-rewrite", false)
	ensureFeatureDir(t, projectRoot, "oauth-rewrite")

	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v should succeed, output: %s", out)

	// state.json should win over any git context
	outputContains(t, out, "FEATURE: oauth-rewrite (from: state.json)")
}

// Traceability: TC-014 -> Task 2 AC #5,#6
func TestTC_014_ExistingCallersUnchangedAfterPriorityChainChange(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Create a single feature directory, no state.json
	ensureFeatureDir(t, projectRoot, "solo-feature")
	_ = os.Remove(filepath.Join(projectRoot, ".forge", "state.json"))

	out, exitCode := forgeFeature(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: solo-feature")
}

// =============================================================================
// Task 3: Verbose flag
// =============================================================================

// Traceability: TC-015 -> Proposal SC #5 / Task 3 AC #1
func TestTC_015_VerboseShowsStateJsonSource(t *testing.T) {
	projectRoot := setupProjectDir(t)
	writeForgeState(t, projectRoot, "my-slug", false)
	ensureFeatureDir(t, projectRoot, "my-slug")

	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: my-slug (from: state.json)")
}

// Traceability: TC-018 -> Task 3 AC #4
func TestTC_016_VerboseShowsFeaturesDirSource(t *testing.T) {
	projectRoot := setupProjectDir(t)
	ensureFeatureDir(t, projectRoot, "detected-feature")
	// No state.json, no git feature context

	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: detected-feature (from: features-dir)")
}

// Traceability: TC-019 -> Task 3 AC #5
func TestTC_017_VerboseShowsNoneWhenNoFeatureSet(t *testing.T) {
	projectRoot := setupProjectDir(t)
	// No feature directories, no state.json

	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v should succeed, output: %s", out)
	outputContains(t, out, "FEATURE: (none)")
}

// Traceability: TC-020 -> Task 3 AC #7 / Task 3 Hard Rules
func TestTC_018_VerboseFlagLocalToFeatureCommandOnly(t *testing.T) {
	projectRoot := setupProjectDir(t)

	// Scenario 1: forge feature -v is recognized
	out, exitCode := forgeFeatureVerbose(t, projectRoot)
	assert.Equal(t, 0, exitCode, "forge feature -v should be recognized, output: %s", out)

	// Scenario 2: forge feature set -v my-feature should NOT be recognized
	cmd := exec.Command(testkit.ForgeBinary, "feature", "set", "-v", "my-feature")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	out2, err := cmd.CombinedOutput()
	assert.Error(t, err, "forge feature set -v should fail or produce error, output: %s", string(out2))

	// Scenario 3: forge feature list -v should NOT be recognized
	cmd3 := exec.Command(testkit.ForgeBinary, "feature", "list", "-v")
	cmd3.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	out3, err3 := cmd3.CombinedOutput()
	assert.Error(t, err3, "forge feature list -v should fail or produce error, output: %s", string(out3))
}
