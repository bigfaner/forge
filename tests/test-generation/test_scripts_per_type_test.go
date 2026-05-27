//go:build cli_functional

package testgeneration

// ==============================================================================
// Per-type test scripts generation — Journey: test-generation
// Tests cover per-type gen-scripts tasks, task types, dependency chains,
// multi-profile tasks, quick mode per-type, and idempotent re-runs.
// ==============================================================================

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// perTypeTaskIndex represents the index.json structure (minimal fields for testing).
type perTypeTaskIndex struct {
	Feature string                 `json:"feature"`
	Tasks   map[string]perTypeTask `json:"tasks"`
}

type perTypeTask struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Profile  string `json:"profile,omitempty"`
	Status   string `json:"status"`
	File     string `json:"file,omitempty"`
}

// runForge executes the forge binary with given args and returns combined output.
func runForge(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, args...)
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
		profileLines := "languages:\n"
		for _, p := range testProfiles {
			profileLines += "  - " + p + "\n"
		}
		profileLines += "surfaces:\n  backend: api\n  frontend: cli\nauto:\n  test:\n    quick: true\n    full: true\n"
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

// readPerTypeIndexJSON reads and parses the index.json file.
func readPerTypeIndexJSON(t *testing.T, dir, slug string) perTypeTaskIndex {
	t.Helper()
	indexPath := filepath.Join(dir, "docs", "features", slug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	require.NoError(t, err, "index.json should exist")
	var idx perTypeTaskIndex
	require.NoError(t, json.Unmarshal(data, &idx), "index.json should be valid JSON")
	return idx
}

// multiTypeTestCases is a test-cases.md with multiple types (UI, API, CLI).
const perTypeMultiTypeTestCases = `# Test Cases

## Summary

| Type | Count |
|------|-------|
| UI   | 5     |
| API  | 3     |
| CLI  | 10    |
| **Total** | **18** |
`

// singleTypeTestCases is a test-cases.md with only CLI type.
const perTypeSingleTypeTestCases = `# Test Cases

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| API  | 0     |
| CLI  | 7     |
| **Total** | **7** |
`

// perTypeNoTypeTestCases is a test-cases.md with all zero counts.
const perTypeNoTypeTestCases = `# Test Cases

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| API  | 0     |
| CLI  | 0     |
| **Total** | **0** |
`

// addBusinessTask adds a business task .md file to the feature.
func addBusinessTask(t *testing.T, dir, slug string) {
	t.Helper()
	tasksDir := filepath.Join(dir, "docs", "features", slug, "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"coding.feature\"\nsurface-key: \"\"\nsurface-type: \"api\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))
}

// ==============================================================================
// TC-001: forge task index creates per-type gen-scripts tasks for multi-type
// ==============================================================================

// Traceability: TC-001 -> test-scripts-per-type proposal: per-type task generation
func TestPerType_TC_001_TaskIndexCreatesPerTypeTasksForMultiType(t *testing.T) {
	dir := setupFeatureProject(t, "multi-type-feat", true, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "multi-type-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "multi-type-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "multi-type-feat")

	// Verify per-type gen-scripts tasks exist
	// go profile capabilities: [api, cli]
	_, hasAPI := idx.Tasks["gen-test-scripts-api"]
	_, hasCLI := idx.Tasks["gen-test-scripts-cli"]

	assert.True(t, hasAPI, "index should contain gen-test-scripts-api task")
	assert.True(t, hasCLI, "index should contain gen-test-scripts-cli task")
}

// ==============================================================================
// TC-002: forge task index creates per-type tasks with correct type field
// ==============================================================================

// Traceability: TC-002 -> test-scripts-per-type proposal: task type = test.gen-scripts
func TestPerType_TC_002_TaskIndexPerTypeTasksHaveCorrectType(t *testing.T) {
	dir := setupFeatureProject(t, "type-check-feat", true, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "type-check-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "type-check-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "type-check-feat")

	for _, key := range []string{"gen-test-scripts-api", "gen-test-scripts-cli"} {
		task, ok := idx.Tasks[key]
		require.True(t, ok, "task %s should exist", key)
		assert.Equal(t, "test.gen-scripts", task.Type,
			"task %s should have type test.gen-scripts", key)
	}
}

// ==============================================================================
// TC-003: forge task index creates single gen-scripts for single-type project
// ==============================================================================

// Traceability: TC-003 -> test-scripts-per-type proposal: single type = one gen task
func TestPerType_TC_003_TaskIndexSingleTypeCreatesOneGenTask(t *testing.T) {
	dir := setupFeatureProject(t, "single-type-feat", true, []string{"go"}, perTypeSingleTypeTestCases)
	addBusinessTask(t, dir, "single-type-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "single-type-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "single-type-feat")

	// go profile capabilities [api, cli] are always used regardless of test-cases content
	// All per-type gen-scripts tasks should exist
	_, hasCLI := idx.Tasks["gen-test-scripts-cli"]
	assert.True(t, hasCLI, "index should contain gen-test-scripts-cli task")

	// Should also have api variant (driven by profile capabilities)
	_, hasAPI := idx.Tasks["gen-test-scripts-api"]
	assert.True(t, hasAPI, "index should contain gen-test-scripts-api task")
}

// ==============================================================================
// TC-004: forge task index without test-cases.md uses profile capabilities for per-type
// ==============================================================================

// Traceability: TC-004 -> per-type tasks from profile capabilities, not test-cases.md
func TestPerType_TC_004_TaskIndexWithoutTestCasesUsesProfileCapabilities(t *testing.T) {
	dir := setupFeatureProject(t, "no-types-feat", true, []string{"go"}, "")
	addBusinessTask(t, dir, "no-types-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "no-types-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "no-types-feat")

	// Without test-cases.md, capabilities still come from profile manifest
	// Per-type tasks are generated based on profile capabilities
	for _, typ := range []string{"api", "cli"} {
		key := "gen-test-scripts-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s (capabilities from profile manifest)", key)
	}
}

// ==============================================================================
// TC-005: forge task index with zero-type test-cases uses profile capabilities
// ==============================================================================

// Traceability: TC-005 -> per-type tasks from profile capabilities, not test-cases.md
func TestPerType_TC_005_TaskIndexZeroTypeTestCasesUsesProfileCapabilities(t *testing.T) {
	dir := setupFeatureProject(t, "zero-types-feat", true, []string{"go"}, perTypeNoTypeTestCases)
	addBusinessTask(t, dir, "zero-types-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "zero-types-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "zero-types-feat")

	// All-zero test-cases types do not change per-type generation;
	// capabilities from profile manifest drive the per-type split
	for _, typ := range []string{"api", "cli"} {
		key := "gen-test-scripts-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s for zero-type test-cases", key)
	}
}

// ==============================================================================
// TC-006: forge task index creates run task depending on all per-type gen tasks
// ==============================================================================

// Traceability: TC-006 -> test-scripts-per-type proposal: run depends on all gen tasks
func TestPerType_TC_006_TaskIndexRunDependsOnAllPerTypeGenTasks(t *testing.T) {
	dir := setupFeatureProject(t, "deps-feat", true, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "deps-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "deps-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "deps-feat")

	// Multi-surface: run-test tasks are per-surface-key (serial chain)
	// First run-test (run-test-backend) should exist and depend on all gen-scripts
	_, hasRun := idx.Tasks["run-test-backend"]
	require.True(t, hasRun, "index should contain run-test-backend task")

	// Verify the first run task's .md file lists per-type gen tasks as dependencies
	tasksDir := filepath.Join(dir, "docs", "features", "deps-feat", "tasks")
	runMDPath := filepath.Join(tasksDir, "run-test-backend.md")
	runMDData, err := os.ReadFile(runMDPath)
	require.NoError(t, err, "run task .md file should exist")
	runMDContent := string(runMDData)

	assert.Contains(t, runMDContent, "T-test-gen-scripts-api", "run task should depend on T-test-gen-scripts-api")
	assert.Contains(t, runMDContent, "T-test-gen-scripts-cli", "run task should depend on T-test-gen-scripts-cli")
}

// ==============================================================================
// TC-007: forge task index with multi-profile creates per-type per-profile tasks
// ==============================================================================

// Traceability: TC-007 -> test-scripts-per-type proposal: per-type gen-scripts with multi-surface
func TestPerType_TC_007_TaskIndexMultiProfilePerTypeTasks(t *testing.T) {
	dir := setupFeatureProject(t, "multi-prof-feat", true, []string{"javascript", "go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "multi-prof-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "multi-prof-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "multi-prof-feat")

	// Multi-surface: surfaces {backend: api, frontend: cli}
	// Per-type gen-scripts tasks generated for each surface type
	for _, typ := range []string{"api", "cli"} {
		key := "gen-test-scripts-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s", key)
	}

	// Run-test tasks are per-surface-key (serial chain)
	for _, key := range []string{"run-test-backend", "run-test-frontend"} {
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s", key)
	}

	// Verify per-type gen-scripts .md files have correct task IDs
	tasksDir := filepath.Join(dir, "docs", "features", "multi-prof-feat", "tasks")
	for _, key := range []string{
		"gen-test-scripts-api",
		"gen-test-scripts-cli",
	} {
		mdPath := filepath.Join(tasksDir, key+".md")
		mdData, err := os.ReadFile(mdPath)
		require.NoError(t, err, "%s.md should exist", key)
		content := string(mdData)
		// ID should match the expected pattern
		assert.Contains(t, content, "T-test-gen-scripts-", "%s.md should have T-test-gen-scripts- ID", key)
	}
}

// ==============================================================================
// TC-008: forge task index quick mode creates per-type gen-and-run tasks
// ==============================================================================

// Traceability: TC-008 -> Quick staged topology: run-test tasks per surface key
func TestPerType_TC_008_TaskIndexQuickModePerSurfaceRunTestTasks(t *testing.T) {
	dir := setupFeatureProject(t, "quick-type-feat", false, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "quick-type-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "quick-type-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "quick-type-feat")

	// Quick staged topology: surfaces {backend: api, frontend: cli}
	// Run-test tasks per surface key
	for _, key := range []string{"backend", "frontend"} {
		taskKey := "run-test-" + key
		_, ok := idx.Tasks[taskKey]
		assert.True(t, ok, "index should contain %s for quick mode", taskKey)
	}

	// Verify verify-regression depends on last run-test
	tasksDir := filepath.Join(dir, "docs", "features", "quick-type-feat", "tasks")
	verifyMDPath := filepath.Join(tasksDir, "verify-regression.md")
	verifyMDData, err := os.ReadFile(verifyMDPath)
	require.NoError(t, err, "verify-regression task .md should exist")
	verifyContent := string(verifyMDData)

	assert.Contains(t, verifyContent, "T-test-run-frontend", "verify-regression should depend on T-test-run-frontend (last in serial chain)")
}

// ==============================================================================
// TC-009: per-type gen-scripts .md files mention the test type
// ==============================================================================

// Traceability: TC-009 -> test-scripts-per-type proposal: gen-scripts MD contains type info
func TestPerType_TC_009_PerTypeGenScriptsMdContainsTestType(t *testing.T) {
	dir := setupFeatureProject(t, "type-md-feat", true, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "type-md-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "type-md-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	// Each per-type .md should mention its type in the body
	// surface types: [api, cli]
	typeCases := []struct {
		key string
		typ string
	}{
		{"gen-test-scripts-api", "api"},
		{"gen-test-scripts-cli", "cli"},
	}

	tasksDir := filepath.Join(dir, "docs", "features", "type-md-feat", "tasks")
	for _, tc := range typeCases {
		mdPath := filepath.Join(tasksDir, tc.key+".md")
		data, err := os.ReadFile(mdPath)
		require.NoError(t, err, "%s.md should exist", tc.key)
		content := string(data)

		// Body should mention the type
		assert.Contains(t, content, tc.typ, "%s.md should mention type %q", tc.key, tc.typ)
	}
}

// ==============================================================================
// TC-010: forge task index idempotent with per-type tasks
// ==============================================================================

// Traceability: TC-010 -> test-scripts-per-type proposal: re-running index is idempotent
func TestPerType_TC_010_TaskIndexPerTypeIdempotent(t *testing.T) {
	dir := setupFeatureProject(t, "idempotent-feat", true, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "idempotent-feat")

	bin := testkit.ForgeBinary

	// Run index twice
	for i := 0; i < 2; i++ {
		cmd := exec.Command(bin, "task", "index", "--feature", "idempotent-feat")
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		require.NoError(t, err, "forge task index run %d should succeed: %s", i+1, out)
	}

	idx := readPerTypeIndexJSON(t, dir, "idempotent-feat")

	// Should have exactly the expected per-type tasks (not duplicated)
	// go capabilities: [api, cli]
	for _, typ := range []string{"api", "cli"} {
		key := "gen-test-scripts-" + typ
		_, ok := idx.Tasks[key]
		assert.True(t, ok, "index should contain %s after idempotent re-run", key)
	}

	// Count total tasks should be reasonable (not doubled)
	// breakdown mode with go (2 capabilities): gen-cases + eval + 2 per-type gen + run + graduate + verify + consolidate = 8
	assert.LessOrEqual(t, len(idx.Tasks), 15,
		"index should not have excessive tasks after idempotent re-run, got %d", len(idx.Tasks))
}

// ==============================================================================
// TC-011: forge task index per-type gen-scripts .md has correct task IDs
// ==============================================================================

// Traceability: TC-011 -> test-scripts-per-type proposal: task IDs include type suffix
func TestPerType_TC_011_PerTypeGenScriptsMdHasCorrectTaskIDs(t *testing.T) {
	dir := setupFeatureProject(t, "tid-feat", true, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "tid-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "tid-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "tid-feat")

	// Verify task IDs have type suffix
	// go capabilities: [api, cli]
	expectedIDs := map[string]string{
		"gen-test-scripts-api": "T-test-gen-scripts-api",
		"gen-test-scripts-cli": "T-test-gen-scripts-cli",
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

// Traceability: TC-012 -> test-scripts-per-type proposal: shared tasks not per-type
func TestPerType_TC_012_TaskIndexSharedInfrastructureNotDuplicated(t *testing.T) {
	dir := setupFeatureProject(t, "shared-feat", true, []string{"go"}, perTypeMultiTypeTestCases)
	addBusinessTask(t, dir, "shared-feat")

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "shared-feat")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := readPerTypeIndexJSON(t, dir, "shared-feat")

	// gen-journeys should NOT have per-type variants (shared infrastructure)
	_, hasGenJourneys := idx.Tasks["gen-journeys"]
	assert.True(t, hasGenJourneys, "index should contain shared gen-journeys task")

	_, hasGenJourneysTUI := idx.Tasks["gen-journeys-tui"]
	_, hasGenJourneysCLI := idx.Tasks["gen-journeys-cli"]
	assert.False(t, hasGenJourneysTUI, "index should NOT contain per-type gen-journeys-tui")
	assert.False(t, hasGenJourneysCLI, "index should NOT contain per-type gen-journeys-cli")

	// Verify shared tasks have correct types
	genJourneys, ok := idx.Tasks["gen-journeys"]
	require.True(t, ok)
	assert.Equal(t, "test.gen-journeys", genJourneys.Type)
}
