//go:build e2e

package testgeneration

// ==============================================================================
// Quick mode test generation — Journey: test-generation
// Tests cover quick mode task index, per-type gen-and-run tasks,
// dependency chains, multi-profile letter suffixes, and InferType mapping.
// ==============================================================================

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// quickSlimTestCases is a test-cases.md with only CLI type (triggers per-type split).
const quickSlimTestCases = `---
feature: "quick-test-slim"
sources:
  - docs/proposals/quick-test-slim/proposal.md
generated: "2026-05-16"
---

# Test Cases: quick-test-slim

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 16  |
| **Total** | **16** |
`

// quickSlimNoTypeTestCases is a test-cases.md with all zero counts (no per-type split).
const quickSlimNoTypeTestCases = `---
feature: "quick-test-slim"
---

# Test Cases: no-types

## Summary

| Type | Count |
|------|-------|
| UI   | 0  |
| API  | 0  |
| CLI  | 0  |
| **Total** | **0** |
`

// quickSlimMultiTypeTestCases is a test-cases.md with multiple types for per-type tests.
const quickSlimMultiTypeTestCases = `---
feature: "quick-test-slim"
---

# Test Cases: multi-type

## Summary

| Type | Count |
|------|-------|
| UI   | 0  |
| API  | 3  |
| CLI  | 10 |
| TUI  | 2  |
| **Total** | **15** |
`

// quickSlimIndex represents the index.json structure (minimal fields for testing).
type quickSlimIndex struct {
	Feature string                       `json:"feature"`
	Tasks   map[string]quickSlimTaskEntry `json:"tasks"`
}

type quickSlimTaskEntry struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Type          string   `json:"type"`
	Profile       string   `json:"profile,omitempty"`
	Status        string   `json:"status"`
	File          string   `json:"file,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
	Scope         string   `json:"scope,omitempty"`
	MainSession   bool     `json:"mainSession,omitempty"`
	EstimatedTime string   `json:"estimatedTime,omitempty"`
}

// quickSlimSetupProject creates a temp project with quick mode structure.
func quickSlimSetupProject(t *testing.T, slug string, testProfiles []string, testCasesContent string) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create proposal directory (triggers quick mode)
	propDir := filepath.Join(dir, "docs", "proposals", slug)
	require.NoError(t, os.MkdirAll(propDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(propDir, "proposal.md"), []byte("# Proposal"), 0644))

	// Create feature dirs
	featureDir := filepath.Join(dir, "docs", "features", slug)
	tasksDir := filepath.Join(featureDir, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))

	// Create .forge/config.yaml with test profiles
	if len(testProfiles) > 0 {
		forgeDir := filepath.Join(dir, ".forge")
		require.NoError(t, os.MkdirAll(forgeDir, 0755))
		profileLines := "languages:\n"
		for _, p := range testProfiles {
			profileLines += "  - " + p + "\n"
		}
		profileLines += "surfaces:\n  backend: api\n  frontend: cli\nauto:\n  test:\n    quick: true\n    full: true\n  consolidateSpecs:\n    quick: true\n    full: true\n"
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

// quickSlimReadIndex reads and parses the index.json file.
func quickSlimReadIndex(t *testing.T, dir, slug string) quickSlimIndex {
	t.Helper()
	indexPath := filepath.Join(dir, "docs", "features", slug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	require.NoError(t, err, "index.json should exist")
	var idx quickSlimIndex
	require.NoError(t, json.Unmarshal(data, &idx), "index.json should be valid JSON")
	return idx
}

// quickSlimAddBusinessTask adds a business task .md file to the feature.
func quickSlimAddBusinessTask(t *testing.T, dir, slug string) {
	t.Helper()
	tasksDir := filepath.Join(dir, "docs", "features", slug, "tasks")
	taskMD := "---\nid: \"1\"\ntitle: \"Implement feature\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"coding.feature\"\nsurface-key: \"\"\nsurface-type: \"api\"\n---\n\n# Implement feature\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-implement.md"), []byte(taskMD), 0644))
}

// quickSlimRunIndex runs forge task index and returns the output.
func quickSlimRunIndex(t *testing.T, dir, slug string) []byte {
	t.Helper()
	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", slug)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)
	return out
}

// =============================================================================
// TC-001: Quick mode single profile generates correct task count
// =============================================================================

// Traceability: TC-001 -> Proposal Success Criteria: quick mode task count
func TestTC_001_QuickModeSingleProfileTaskCount(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-001", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-001")
	quickSlimRunIndex(t, dir, "test-qts-001")

	idx := quickSlimReadIndex(t, dir, "test-qts-001")

	// Count auto-generated pipeline tasks (T-test-* and T-quick-doc-drift)
	// Quick staged topology: gen-journeys(1) + run-test-backend(1) + run-test-frontend(1) + verify-regression(1) = 4 test tasks
	// Plus quick-drift-detection(1) = 5 total auto-gen tasks
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-test-") || strings.HasPrefix(task.ID, "T-quick-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 5, testTaskCount, "quick mode with go profile (2 surfaces) should generate exactly 5 pipeline tasks")
}

// =============================================================================
// TC-002: Quick mode merged task has gen-and-run type
// =============================================================================

// Traceability: TC-002 -> Quick staged topology: run-test tasks have correct type
func TestTC_002_QuickModeRunTestTasksHaveCorrectType(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-002", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-002")
	quickSlimRunIndex(t, dir, "test-qts-002")

	idx := quickSlimReadIndex(t, dir, "test-qts-002")

	// Find run-test tasks (per-surface-key in quick staged topology)
	for _, key := range []string{"backend", "frontend"} {
		taskKey := "run-test-" + key
		task, ok := idx.Tasks[taskKey]
		require.True(t, ok, "task key %s should exist in index", taskKey)
		assert.Equal(t, "test.run", task.Type,
			"run-test-%s should have type test.run", key)
	}
}


// =============================================================================
// TC-004: Quick mode per-type creates independent gen-and-run tasks
// =============================================================================

// Traceability: TC-004 -> Proposal Success Criteria: per-surface run-test tasks are independent
func TestTC_004_QuickModePerSurfaceRunTestTasks(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-004", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-004")
	quickSlimRunIndex(t, dir, "test-qts-004")

	idx := quickSlimReadIndex(t, dir, "test-qts-004")

	// Find per-surface-key run-test tasks
	for _, key := range []string{"backend", "frontend"} {
		taskKey := "run-test-" + key
		task, ok := idx.Tasks[taskKey]
		require.True(t, ok, "task key %s should exist", taskKey)
		assert.Equal(t, "test.run", task.Type,
			"run-test-%s should have type test.run", key)
	}
}

// =============================================================================
// TC-005: Quick mode dependency chain is correct after merge
// =============================================================================

// Traceability: TC-005 -> Task 1 AC: dependency chain gen-journeys -> run-test(s) -> verify-regression -> drift
func TestTC_005_QuickModeDependencyChainCorrect(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-005", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-005")
	quickSlimRunIndex(t, dir, "test-qts-005")

	idx := quickSlimReadIndex(t, dir, "test-qts-005")

	// Build a map from ID to task for easy lookup
	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// Verify dependency chain in quick staged topology
	t.Run("T-test-run-backend depends on T-test-gen-journeys", func(t *testing.T) {
		task, ok := byID["T-test-run-backend"]
		require.True(t, ok, "T-test-run-backend should exist")
		assert.Contains(t, task.Dependencies, "T-test-gen-journeys",
			"T-test-run-backend should depend on T-test-gen-journeys")
	})

	t.Run("T-test-run-frontend depends on T-test-run-backend (serial chain)", func(t *testing.T) {
		task, ok := byID["T-test-run-frontend"]
		require.True(t, ok, "T-test-run-frontend should exist")
		assert.Contains(t, task.Dependencies, "T-test-run-backend",
			"T-test-run-frontend should depend on T-test-run-backend")
	})

	t.Run("T-test-verify-regression depends on last run-test", func(t *testing.T) {
		task, ok := byID["T-test-verify-regression"]
		require.True(t, ok, "T-test-verify-regression should exist")
		assert.Contains(t, task.Dependencies, "T-test-run-frontend",
			"T-test-verify-regression should depend on T-test-run-frontend")
	})

	t.Run("T-quick-doc-drift depends on T-test-verify-regression", func(t *testing.T) {
		task, ok := byID["T-quick-doc-drift"]
		require.True(t, ok, "T-quick-doc-drift should exist")
		assert.Contains(t, task.Dependencies, "T-test-verify-regression",
			"T-quick-doc-drift should depend on T-test-verify-regression")
	})
}

// =============================================================================
// TC-006: Quick mode per-type dependency chain fans in correctly
// =============================================================================

// Traceability: TC-006 -> Proposal Success Criteria: run-test serial chain fans into verify-regression
func TestTC_006_QuickModeRunTestSerialChainFanIn(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-006", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-006")
	quickSlimRunIndex(t, dir, "test-qts-006")

	idx := quickSlimReadIndex(t, dir, "test-qts-006")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// First run-test depends on gen-journeys
	task, ok := byID["T-test-run-backend"]
	require.True(t, ok, "T-test-run-backend should exist")
	assert.Contains(t, task.Dependencies, "T-test-gen-journeys",
		"T-test-run-backend should depend on T-test-gen-journeys")

	// Second run-test depends on first (serial chain)
	task2, ok := byID["T-test-run-frontend"]
	require.True(t, ok, "T-test-run-frontend should exist")
	assert.Contains(t, task2.Dependencies, "T-test-run-backend",
		"T-test-run-frontend should depend on T-test-run-backend")

	// verify-regression depends on last run-test
	verifyTask, ok := byID["T-test-verify-regression"]
	require.True(t, ok, "T-test-verify-regression should exist")
	assert.Contains(t, verifyTask.Dependencies, "T-test-run-frontend",
		"T-test-verify-regression should depend on T-test-run-frontend")
}

// =============================================================================
// TC-007: Breakdown mode is unchanged by quick mode merge
// =============================================================================

// Traceability: TC-007 -> Task 1 Hard Rules: breakdown mode unchanged
func TestTC_007_BreakdownModeUnchangedByQuickMerge(t *testing.T) {
	// Create a breakdown mode project (has PRD, not proposal)
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	featureDir := filepath.Join(dir, "docs", "features", "test-qts-007")
	tasksDir := filepath.Join(featureDir, "tasks")
	prdDir := filepath.Join(featureDir, "prd")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	require.NoError(t, os.MkdirAll(prdDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(prdDir, "prd-spec.md"), []byte("# PRD"), 0644))

	forgeDir := filepath.Join(dir, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(forgeDir, "config.yaml"),
		[]byte("surfaces:\n  backend: api\n  frontend: cli\nauto:\n  test:\n    quick: true\n    full: true\n  consolidateSpecs:\n    quick: true\n    full: true\n"), 0644))

	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"coding.feature\"\nsurface-key: \"\"\nsurface-type: \"api\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := testkit.ForgeBinary
	cmd := exec.Command(bin, "task", "index", "--feature", "test-qts-007")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := quickSlimReadIndex(t, dir, "test-qts-007")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// T-test-gen-scripts-api should have type gen-scripts (NOT gen-and-run)
	task2api, ok := byID["T-test-gen-scripts-api"]
	require.True(t, ok, "T-test-gen-scripts-api should exist in breakdown mode")
	assert.Equal(t, "test.gen-scripts", task2api.Type,
		"breakdown T-test-gen-scripts-api should have type test.gen-scripts, not gen-and-run")

	// Multi-surface: run-test tasks are per-surface-key (T-test-run-backend, T-test-run-frontend)
	taskRun, ok := byID["T-test-run-backend"]
	require.True(t, ok, "T-test-run-backend should exist in breakdown mode")
	assert.Equal(t, "test.run", taskRun.Type,
		"breakdown T-test-run-backend should have type test.run")

	// Total test pipeline tasks:
	// gen-journeys(1) + eval-journey(1) + gen-contracts(1) + eval-contract(1) +
	// 2 per-type gen-scripts + 2 per-surface-key run-test + verify-regression(1) = 10
	// Plus consolidate-specs (T-specs-consolidate) = 11 total T-test-*/T-specs-*/T-eval-* tasks
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-test-") || strings.HasPrefix(task.ID, "T-specs-") || strings.HasPrefix(task.ID, "T-eval-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 10, testTaskCount,
		"breakdown mode with go (2 surfaces, 2 types) should have 10 pipeline tasks")
}

// =============================================================================
// TC-008: Quick mode multi-profile with letter suffixes works
// =============================================================================

// Traceability: TC-008 -> Quick staged topology with multiple profiles
func TestTC_008_QuickModeMultiProfile(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-008", []string{"go", "javascript"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-008")
	quickSlimRunIndex(t, dir, "test-qts-008")

	idx := quickSlimReadIndex(t, dir, "test-qts-008")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// Quick staged topology with multi-profile: surfaces {backend: api, frontend: cli}
	// gen-journeys -> run-test-backend -> run-test-frontend -> verify-regression -> drift
	for _, id := range []string{"T-test-gen-journeys", "T-test-run-backend", "T-test-run-frontend", "T-test-verify-regression"} {
		_, ok := byID[id]
		assert.True(t, ok, "%s should exist", id)
	}

	// Shared tasks T-quick-doc-drift should exist
	_, ok := byID["T-quick-doc-drift"]
	assert.True(t, ok, "T-quick-doc-drift should exist as shared task")
}



// =============================================================================
// TC-011: InferType maps merged IDs correctly
// =============================================================================

// Traceability: TC-011 -> Quick staged topology: task IDs follow correct patterns
func TestTC_011_InferTypeMapsStagedIDsCorrectly(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-011", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-011")
	quickSlimRunIndex(t, dir, "test-qts-011")

	idx := quickSlimReadIndex(t, dir, "test-qts-011")

	// Quick staged topology: gen-journeys, run-test-backend, run-test-frontend, verify-regression
	testCases := []struct {
		key  string
		id   string
		typ  string
	}{
		{"gen-journeys", "T-test-gen-journeys", "test.gen-journeys"},
		{"run-test-backend", "T-test-run-backend", "test.run"},
		{"run-test-frontend", "T-test-run-frontend", "test.run"},
		{"verify-regression", "T-test-verify-regression", "test.verify-regression"},
	}

	for _, tc := range testCases {
		task, ok := idx.Tasks[tc.key]
		require.True(t, ok, "task key %s should exist", tc.key)
		assert.Equal(t, tc.id, task.ID, "key %s should have ID %s", tc.key, tc.id)
		assert.Equal(t, tc.typ, task.Type, "key %s should have type %s", tc.key, tc.typ)
	}

	// Multi-profile quick mode: same staged topology (surfaces drive task keys, not profiles)
	dir2 := quickSlimSetupProject(t, "test-qts-011b", []string{"go", "javascript"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir2, "test-qts-011b")
	quickSlimRunIndex(t, dir2, "test-qts-011b")

	idx2 := quickSlimReadIndex(t, dir2, "test-qts-011b")

	for _, tc := range testCases {
		task, ok := idx2.Tasks[tc.key]
		require.True(t, ok, "multi-profile: task key %s should exist", tc.key)
		assert.Equal(t, tc.typ, task.Type, "multi-profile: key %s should have type %s", tc.key, tc.typ)
	}
}

// =============================================================================
// TC-012: Quick mode single profile produces 5 tasks total
// =============================================================================

// Traceability: TC-012 -> Quick staged: single profile with 2 surfaces produces correct count
func TestTC_012_QuickModeSingleProfileProducesCorrectTaskCount(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-012", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-012")
	quickSlimRunIndex(t, dir, "test-qts-012")

	idx := quickSlimReadIndex(t, dir, "test-qts-012")

	// Quick staged topology with surfaces {backend: api, frontend: cli}:
	// gen-journeys(1) + run-test-backend(1) + run-test-frontend(1) + verify-regression(1) + drift(1) = 5
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-test-") || strings.HasPrefix(task.ID, "T-quick-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 5, testTaskCount,
		"single profile quick mode with 2 surfaces should produce exactly 5 pipeline tasks")
}


// =============================================================================
// TC-014: Merged task generates correct task .md file
// =============================================================================

// Traceability: TC-014 -> Quick staged: run-test task generates valid .md
func TestTC_014_RunTestTaskGeneratesCorrectMD(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-014", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-014")
	quickSlimRunIndex(t, dir, "test-qts-014")

	tasksDir := filepath.Join(dir, "docs", "features", "test-qts-014", "tasks")

	// Find a run-test task .md file
	mdPath := filepath.Join(tasksDir, "run-test-backend.md")
	data, err := os.ReadFile(mdPath)
	require.NoError(t, err, "run-test-backend.md should exist")

	content := string(data)

	// Verify frontmatter fields
	assert.Contains(t, content, `id: "T-test-run-backend"`, "md should contain correct task ID")
	assert.Contains(t, content, `type: "test.run"`, "md should contain correct type")
}

// =============================================================================
// TC-015: DetectTypesFromTestCases correctly parses test-cases.md summary table
// =============================================================================

// Traceability: TC-015 -> Quick staged: surfaces from config drive run-test task generation
func TestTC_015_SurfacesDriveRunTestTasks(t *testing.T) {
	// Quick staged topology: surfaces {backend: api, frontend: cli} drive per-surface run-test tasks.
	dir := quickSlimSetupProject(t, "test-qts-015a", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-015a")
	quickSlimRunIndex(t, dir, "test-qts-015a")

	idx := quickSlimReadIndex(t, dir, "test-qts-015a")

	// With surfaces {backend, frontend}, should have per-surface run-test tasks
	found := false
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-test-run-") {
			found = true
			break
		}
	}
	assert.True(t, found, "per-surface run-test tasks should be generated when surfaces are configured")

	// Even with zero-type test cases, surfaces still drive run-test generation
	dir2 := quickSlimSetupProject(t, "test-qts-015b", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir2, "test-qts-015b")
	quickSlimRunIndex(t, dir2, "test-qts-015b")

	idx2 := quickSlimReadIndex(t, dir2, "test-qts-015b")

	// Per-surface run-test tasks are generated based on config surfaces
	runTestCount := 0
	for _, key := range []string{"backend", "frontend"} {
		foundKey := false
		for _, task := range idx2.Tasks {
			if task.ID == "T-test-run-"+key {
				foundKey = true
				runTestCount++
				break
			}
		}
		assert.True(t, foundKey, "T-test-run-%s should exist (driven by config surfaces)", key)
	}
	assert.Equal(t, 2, runTestCount, "per-surface run-test tasks should be generated from config surfaces (backend, frontend)")
}
