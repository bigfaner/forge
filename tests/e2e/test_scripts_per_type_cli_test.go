//go:build e2e

package e2e

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// taskIndex represents the index.json structure (minimal fields for testing).
type taskIndex struct {
	Feature string            `json:"feature"`
	Tasks   map[string]taskEntry `json:"tasks"`
}

type taskEntry struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Profile  string `json:"profile,omitempty"`
	Status   string `json:"status"`
	File     string `json:"file,omitempty"`
}

// forgeBinPath caches the built forge binary path.
var forgeBinPath string

// forgeBin returns the path to the forge binary, building it if necessary.
func forgeBin(t *testing.T) string {
	t.Helper()
	if forgeBinPath != "" {
		if _, err := os.Stat(forgeBinPath); err == nil {
			return forgeBinPath
		}
	}
	binName := "forge"
	if runtime.GOOS == "windows" {
		binName = "forge.exe"
	}
	// Resolve forge-cli directory relative to this test file
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// Walk up to project root (contains justfile)
	dir := filepath.Dir(thisFile)
	for dir != "/" && dir != "" {
		if _, err := os.Stat(filepath.Join(dir, "justfile")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	binPath := filepath.Join(dir, "forge-cli", "bin", binName)
	if _, err := os.Stat(binPath); err != nil {
		cliDir := filepath.Join(dir, "forge-cli")
		buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/forge/")
		buildCmd.Dir = cliDir
		if out, err := buildCmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to build forge binary: %s: %s", err, out)
		}
	}
	absPath, err := filepath.Abs(binPath)
	if err != nil {
		t.Fatalf("failed to resolve binary path: %s", err)
	}
	forgeBinPath = absPath
	return absPath
}

// runForge executes the forge binary with given args and returns combined output.
func runForge(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	bin := forgeBin(t)
	cmd := exec.Command(bin, args...)
	return cmd.CombinedOutput()
}

// setupFeatureProject creates a temp project with the given feature structure.
func setupFeatureProject(t *testing.T, slug string, hasPRD bool, testProfiles []string, testCasesContent string) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create feature dirs
	featureDir := filepath.Join(dir, "docs", "features", slug)
	tasksDir := filepath.Join(featureDir, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))

	// Create mode-determining file
	if hasPRD {
		prdDir := filepath.Join(featureDir, "prd")
		require.NoError(t, os.MkdirAll(prdDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(prdDir, "prd-spec.md"), []byte("# PRD"), 0644))
	} else {
		propDir := filepath.Join(dir, "docs", "proposals", slug)
		require.NoError(t, os.MkdirAll(propDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte("# Proposal"), 0644))
	}

	// Create .forge/config.yaml with test profiles
	if len(testProfiles) > 0 {
		forgeDir := filepath.Join(dir, ".forge")
		require.NoError(t, os.MkdirAll(forgeDir, 0755))
		profileLines := "test-profiles:\n"
		for _, p := range testProfiles {
			profileLines += "  - " + p + "\n"
		}
		require.NoError(t, os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(profileLines), 0644))
	}

	// Create test-cases.md if content provided
	if testCasesContent != "" {
		testingDir := filepath.Join(featureDir, "testing")
		require.NoError(t, os.MkdirAll(testingDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(testingDir, "test-cases.md"), []byte(testCasesContent), 0644))
	}

	return dir
}

// readIndexJSON reads and parses the index.json file.
func readIndexJSON(t *testing.T, dir, slug string) taskIndex {
	t.Helper()
	indexPath := filepath.Join(dir, "docs", "features", slug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	require.NoError(t, err, "index.json should exist")
	var idx taskIndex
	require.NoError(t, json.Unmarshal(data, &idx), "index.json should be valid JSON")
	return idx
}

// multiTypeTestCases is a test-cases.md with multiple types (UI, API, CLI).
const multiTypeTestCases = `# Test Cases

## Summary

| Type | Count |
|------|-------|
| UI   | 5     |
| API  | 3     |
| CLI  | 10    |
| **Total** | **18** |
`

// singleTypeTestCases is a test-cases.md with only CLI type.
const singleTypeTestCases = `# Test Cases

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| API  | 0     |
| CLI  | 7     |
| **Total** | **7** |
`

// noTypeTestCases is a test-cases.md with all zero counts.
const noTypeTestCases = `# Test Cases

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| API  | 0     |
| CLI  | 0     |
| **Total** | **0** |
`

// ==============================================================================
// TC-001: forge task index creates per-type gen-scripts tasks for multi-type
// ==============================================================================

// Traceability: TC-001 -> test-scripts-per-type proposal: per-type task generation
func TestTC_001_TaskIndexCreatesPerTypeTasksForMultiType(t *testing.T) {
	dir := setupFeatureProject(t, "multi-type-feat", true, []string{"go-test"}, multiTypeTestCases)

	// Create a business task so index has content
	tasksDir := filepath.Join(dir, "docs", "features", "multi-type-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "multi-type-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "multi-type-feat")

	// Verify per-type gen-scripts tasks exist
	_, hasUI := idx.Tasks["gen-test-scripts-go-test-tui"]
	_, hasAPI := idx.Tasks["gen-test-scripts-go-test-api"]
	_, hasCLI := idx.Tasks["gen-test-scripts-go-test-cli"]

	assert.True(t, hasUI, "index should contain gen-test-scripts-go-test-tui task")
	assert.True(t, hasAPI, "index should contain gen-test-scripts-go-test-api task")
	assert.True(t, hasCLI, "index should contain gen-test-scripts-go-test-cli task")
}

// ==============================================================================
// TC-002: forge task index creates per-type tasks with correct type field
// ==============================================================================

// Traceability: TC-002 -> test-scripts-per-type proposal: task type = test-pipeline.gen-scripts
func TestTC_002_TaskIndexPerTypeTasksHaveCorrectType(t *testing.T) {
	dir := setupFeatureProject(t, "type-check-feat", true, []string{"go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "type-check-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "type-check-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "type-check-feat")

	for _, key := range []string{"gen-test-scripts-go-test-tui", "gen-test-scripts-go-test-api", "gen-test-scripts-go-test-cli"} {
		task, ok := idx.Tasks[key]
		require.True(t, ok, "task %s should exist", key)
		assert.Equal(t, "test-pipeline.gen-scripts", task.Type,
			"task %s should have type test-pipeline.gen-scripts", key)
	}
}

// ==============================================================================
// TC-003: forge task index creates single gen-scripts for single-type project
// ==============================================================================

// Traceability: TC-003 -> test-scripts-per-type proposal: single type = one gen task
func TestTC_003_TaskIndexSingleTypeCreatesOneGenTask(t *testing.T) {
	dir := setupFeatureProject(t, "single-type-feat", true, []string{"go-test"}, singleTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "single-type-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "single-type-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "single-type-feat")

	// Profile capabilities drive per-type tasks regardless of test-cases.md content
	for _, typ := range []string{"tui", "api", "cli"} {
		key := "gen-test-scripts-go-test-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s (from profile capabilities)", key)
	}
}

// ==============================================================================
// TC-004: forge task index without test-cases.md uses profile capabilities for per-type
// ==============================================================================

// Traceability: TC-004 -> per-type tasks from profile capabilities, not test-cases.md
func TestTC_004_TaskIndexWithoutTestCasesUsesProfileCapabilities(t *testing.T) {
	dir := setupFeatureProject(t, "no-types-feat", true, []string{"go-test"}, "")

	tasksDir := filepath.Join(dir, "docs", "features", "no-types-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "no-types-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "no-types-feat")

	// Without test-cases.md, per-type tasks still generated from profile capabilities
	for _, typ := range []string{"tui", "api", "cli"} {
		key := "gen-test-scripts-go-test-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s (from profile capabilities)", key)
	}
}

// ==============================================================================
// TC-005: forge task index with zero-type test-cases uses profile capabilities
// ==============================================================================

// Traceability: TC-005 -> per-type tasks from profile capabilities, not test-cases.md
func TestTC_005_TaskIndexZeroTypeTestCasesUsesProfileCapabilities(t *testing.T) {
	dir := setupFeatureProject(t, "zero-types-feat", true, []string{"go-test"}, noTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "zero-types-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "zero-types-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "zero-types-feat")

	// All-zero test-cases still produces per-type tasks from profile capabilities
	for _, typ := range []string{"tui", "api", "cli"} {
		key := "gen-test-scripts-go-test-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s (from profile capabilities)", key)
	}
}

// ==============================================================================
// TC-006: forge task index creates run task depending on all per-type gen tasks
// ==============================================================================

// Traceability: TC-006 -> test-scripts-per-type proposal: run depends on all gen tasks
func TestTC_006_TaskIndexRunDependsOnAllPerTypeGenTasks(t *testing.T) {
	dir := setupFeatureProject(t, "deps-feat", true, []string{"go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "deps-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "deps-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "deps-feat")

	// Run task should exist
	_, hasRun := idx.Tasks["run-e2e-tests-go-test"]
	require.True(t, hasRun, "index should contain run-e2e-tests-go-test task")

	// Verify the run task's .md file lists per-type gen tasks as dependencies
	runMDPath := filepath.Join(tasksDir, "run-e2e-tests-go-test.md")
	runMDData, err := os.ReadFile(runMDPath)
	require.NoError(t, err, "run task .md file should exist")
	runMDContent := string(runMDData)

	assert.Contains(t, runMDContent, "T-test-2-tui", "run task should depend on T-test-2-tui")
	assert.Contains(t, runMDContent, "T-test-2-api", "run task should depend on T-test-2-api")
	assert.Contains(t, runMDContent, "T-test-2-cli", "run task should depend on T-test-2-cli")
}

// ==============================================================================
// TC-007: forge task index with multi-profile creates per-type per-profile tasks
// ==============================================================================

// Traceability: TC-007 -> test-scripts-per-type proposal: per-profile per-type tasks
func TestTC_007_TaskIndexMultiProfilePerTypeTasks(t *testing.T) {
	dir := setupFeatureProject(t, "multi-prof-feat", true, []string{"web-playwright", "go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "multi-prof-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "multi-prof-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "multi-prof-feat")

	// Profile-a (web-playwright) per-type tasks
	for _, typ := range []string{"web-ui", "api", "cli"} {
		key := "gen-test-scripts-web-playwright-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s for profile web-playwright", key)
	}

	// Profile-b (go-test) per-type tasks
	for _, typ := range []string{"tui", "api", "cli"} {
		key := "gen-test-scripts-go-test-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s for profile go-test", key)
	}

	// Verify profile-suffixed IDs in generated .md files
	for _, key := range []string{
		"gen-test-scripts-web-playwright-web-ui",
		"gen-test-scripts-web-playwright-api",
		"gen-test-scripts-go-test-tui",
		"gen-test-scripts-go-test-cli",
	} {
		mdPath := filepath.Join(tasksDir, key+".md")
		mdData, err := os.ReadFile(mdPath)
		require.NoError(t, err, "%s.md should exist", key)
		content := string(mdData)
		// ID should have profile letter suffix (a for web-playwright, b for go-test)
		if strings.HasPrefix(key, "gen-test-scripts-web-playwright") {
			assert.Contains(t, content, "T-test-2a-", "%s.md should have profile-a suffixed ID", key)
		} else {
			assert.Contains(t, content, "T-test-2b-", "%s.md should have profile-b suffixed ID", key)
		}
	}
}

// ==============================================================================
// TC-008: forge task index quick mode creates per-type gen-and-run tasks
// ==============================================================================

// Traceability: TC-008 -> test-scripts-per-type proposal: quick mode per-type (merged gen+run)
func TestTC_008_TaskIndexQuickModePerTypeTasks(t *testing.T) {
	dir := setupFeatureProject(t, "quick-type-feat", false, []string{"go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "quick-type-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "quick-type-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "quick-type-feat")

	// Quick mode should have per-type gen-and-run tasks with "quick" prefix
	for _, typ := range []string{"tui", "api", "cli"} {
		key := "quick-gen-and-run-go-test-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s for quick mode", key)
	}

	// Quick graduate task should depend on all per-type gen-and-run tasks
	gradMDPath := filepath.Join(tasksDir, "quick-graduate-go-test.md")
	gradMDData, err := os.ReadFile(gradMDPath)
	require.NoError(t, err, "quick graduate task .md should exist")
	gradContent := string(gradMDData)

	assert.Contains(t, gradContent, "T-quick-2-tui", "quick graduate task should depend on T-quick-2-tui")
	assert.Contains(t, gradContent, "T-quick-2-api", "quick graduate task should depend on T-quick-2-api")
	assert.Contains(t, gradContent, "T-quick-2-cli", "quick graduate task should depend on T-quick-2-cli")
}

// ==============================================================================
// TC-009: per-type gen-scripts .md files mention the test type
// ==============================================================================

// Traceability: TC-009 -> test-scripts-per-type proposal: gen-scripts MD contains type info
func TestTC_009_PerTypeGenScriptsMdContainsTestType(t *testing.T) {
	dir := setupFeatureProject(t, "type-md-feat", true, []string{"go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "type-md-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "type-md-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	// Each per-type .md should mention its type in the body
	typeCases := []struct {
		key  string
		typ  string
	}{
		{"gen-test-scripts-go-test-tui", "tui"},
		{"gen-test-scripts-go-test-api", "api"},
		{"gen-test-scripts-go-test-cli", "cli"},
	}

	for _, tc := range typeCases {
		mdPath := filepath.Join(tasksDir, tc.key+".md")
		data, err := os.ReadFile(mdPath)
		require.NoError(t, err, "%s.md should exist", tc.key)
		content := string(data)

		// Body should mention the type
		assert.Contains(t, content, tc.typ, "%s.md should mention type %q", tc.key, tc.typ)
		// Body should mention profile
		assert.Contains(t, content, "go-test", "%s.md should mention profile go-test", tc.key)
	}
}

// ==============================================================================
// TC-010: forge task index idempotent with per-type tasks
// ==============================================================================

// Traceability: TC-010 -> test-scripts-per-type proposal: re-running index is idempotent
func TestTC_010_TaskIndexPerTypeIdempotent(t *testing.T) {
	dir := setupFeatureProject(t, "idempotent-feat", true, []string{"go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "idempotent-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)

	// Run index twice
	for i := 0; i < 2; i++ {
		cmd := exec.Command(bin, "task", "index", "--feature", "idempotent-feat")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "forge task index run %d should succeed: %s", i+1, out)
	}

	idx := readIndexJSON(t, dir, "idempotent-feat")

	// Should have exactly the expected per-type tasks (not duplicated)
	for _, typ := range []string{"tui", "api", "cli"} {
		key := "gen-test-scripts-go-test-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s after idempotent re-run", key)
	}

	// Count total tasks should be reasonable (not doubled)
	assert.LessOrEqual(t, len(idx.Tasks), 15,
		"index should not have excessive tasks after idempotent re-run, got %d", len(idx.Tasks))
}

// ==============================================================================
// TC-011: forge task index per-type gen-scripts .md has correct task IDs
// ==============================================================================

// Traceability: TC-011 -> test-scripts-per-type proposal: task IDs include type suffix
func TestTC_011_PerTypeGenScriptsMdHasCorrectTaskIDs(t *testing.T) {
	dir := setupFeatureProject(t, "tid-feat", true, []string{"go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "tid-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "tid-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "tid-feat")

	// Verify task IDs have type suffix
	expectedIDs := map[string]string{
		"gen-test-scripts-go-test-tui":  "T-test-2-tui",
		"gen-test-scripts-go-test-api": "T-test-2-api",
		"gen-test-scripts-go-test-cli": "T-test-2-cli",
	}
	for key, wantID := range expectedIDs {
		task, ok := idx.Tasks[key]
		require.True(t, ok, "task %s should exist", key)
		assert.Equal(t, wantID, task.ID, "task %s ID mismatch", key)
	}
}

// ==============================================================================
// TC-012: forge task index shared infrastructure tasks are not duplicated
// ==============================================================================

// Traceability: TC-012 -> test-scripts-per-type proposal: shared tasks (gen-cases, eval-cases) not per-type
func TestTC_012_TaskIndexSharedInfrastructureNotDuplicated(t *testing.T) {
	dir := setupFeatureProject(t, "shared-feat", true, []string{"go-test"}, multiTypeTestCases)

	tasksDir := filepath.Join(dir, "docs", "features", "shared-feat", "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"feature\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := forgeBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "shared-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readIndexJSON(t, dir, "shared-feat")

	// gen-test-cases and eval-test-cases should NOT have per-type variants
	_, hasGenCases := idx.Tasks["gen-test-cases"]
	assert.True(t, hasGenCases, "index should contain shared gen-test-cases task")

	_, hasGenCasesUI := idx.Tasks["gen-test-cases-ui"]
	_, hasGenCasesCLI := idx.Tasks["gen-test-cases-cli"]
	assert.False(t, hasGenCasesUI, "index should NOT contain per-type gen-test-cases-ui")
	assert.False(t, hasGenCasesCLI, "index should NOT contain per-type gen-test-cases-cli")

	// Verify shared tasks have correct types
	genCases, ok := idx.Tasks["gen-test-cases"]
	require.True(t, ok)
	assert.Equal(t, "test-pipeline.gen-cases", genCases.Type)
}
