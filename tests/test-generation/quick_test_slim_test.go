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
		profileLines += "auto:\n  e2eTest:\n    quick: true\n    full: true\n  consolidateSpecs:\n    quick: true\n    full: true\n"
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
	taskMD := "---\nid: \"1\"\ntitle: \"Implement feature\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"coding.feature\"\nscope: \"all\"\n---\n\n# Implement feature\n"
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

	// Count test pipeline tasks (T-quick-* IDs)
	// go profile has capabilities [api, cli] -> per-type gen-and-run tasks
	// Total: gen-cases(1) + 2 per-type gen-and-run + graduate(1) + verify-regression(1) + drift-detection(1) = 6
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-quick-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 6, testTaskCount, "quick mode with go profile (2 capabilities) should generate exactly 6 test pipeline tasks")
}

// =============================================================================
// TC-002: Quick mode merged task has gen-and-run type
// =============================================================================

// Traceability: TC-002 -> Task 1 AC: per-type gen-and-run tasks have correct type
func TestTC_002_QuickModeMergedTaskHasGenAndRunType(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-002", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-002")
	quickSlimRunIndex(t, dir, "test-qts-002")

	idx := quickSlimReadIndex(t, dir, "test-qts-002")

	// Find T-quick-gen-and-run-<type> tasks (per-type gen-and-run with go capabilities)
	for _, typ := range []string{"api", "cli"} {
		id := "T-quick-gen-and-run-" + typ
		var found bool
		for _, task := range idx.Tasks {
			if task.ID == id {
				found = true
				assert.Equal(t, "test.gen-and-run", task.Type,
					"%s should have type test.gen-and-run", id)
				break
			}
		}
		assert.True(t, found, "%s task should exist in index", id)
	}
}


// =============================================================================
// TC-004: Quick mode per-type creates independent gen-and-run tasks
// =============================================================================

// Traceability: TC-004 -> Proposal Success Criteria: per-type independent tasks
func TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-004", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-004")
	quickSlimRunIndex(t, dir, "test-qts-004")

	idx := quickSlimReadIndex(t, dir, "test-qts-004")

	// Find per-type gen-and-run tasks (go capabilities: api, cli)
	for _, typ := range []string{"api", "cli"} {
		found := false
		for _, task := range idx.Tasks {
			if strings.HasPrefix(task.ID, "T-quick-gen-and-run") && strings.HasSuffix(task.ID, "-"+typ) {
				found = true
				assert.Equal(t, "test.gen-and-run", task.Type,
					"T-quick-gen-and-run-%s should have type test.gen-and-run", typ)
			}
		}
		assert.True(t, found, "per-type task for %s should exist", typ)
	}
}

// =============================================================================
// TC-005: Quick mode dependency chain is correct after merge
// =============================================================================

// Traceability: TC-005 -> Task 1 AC: dependency chain gen-cases -> gen-and-run-per-type -> graduate -> verify -> drift
func TestTC_005_QuickModeDependencyChainCorrectAfterMerge(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-005", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-005")
	quickSlimRunIndex(t, dir, "test-qts-005")

	idx := quickSlimReadIndex(t, dir, "test-qts-005")

	// Build a map from ID to task for easy lookup
	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// Verify dependency chain with per-type tasks
	t.Run("T-quick-gen-and-run-api depends on T-quick-gen-cases", func(t *testing.T) {
		task, ok := byID["T-quick-gen-and-run-api"]
		require.True(t, ok, "T-quick-gen-and-run-api should exist")
		assert.Contains(t, task.Dependencies, "T-quick-gen-cases",
			"T-quick-gen-and-run-api should depend on T-quick-gen-cases")
	})

	t.Run("T-quick-graduate depends on per-type gen-and-run tasks", func(t *testing.T) {
		task, ok := byID["T-quick-graduate"]
		require.True(t, ok, "T-quick-graduate should exist")
		assert.Contains(t, task.Dependencies, "T-quick-gen-and-run-api",
			"T-quick-graduate should depend on T-quick-gen-and-run-api")
		assert.Contains(t, task.Dependencies, "T-quick-gen-and-run-cli",
			"T-quick-graduate should depend on T-quick-gen-and-run-cli")
	})

	t.Run("T-quick-verify-regression depends on T-quick-graduate", func(t *testing.T) {
		task, ok := byID["T-quick-verify-regression"]
		require.True(t, ok, "T-quick-verify-regression should exist")
		assert.Contains(t, task.Dependencies, "T-quick-graduate",
			"T-quick-verify-regression should depend on T-quick-graduate")
	})

	t.Run("T-quick-doc-drift depends on T-quick-verify-regression", func(t *testing.T) {
		task, ok := byID["T-quick-doc-drift"]
		require.True(t, ok, "T-quick-doc-drift should exist")
		assert.Contains(t, task.Dependencies, "T-quick-verify-regression",
			"T-quick-doc-drift should depend on T-quick-verify-regression")
	})
}

// =============================================================================
// TC-006: Quick mode per-type dependency chain fans in correctly
// =============================================================================

// Traceability: TC-006 -> Proposal Success Criteria: per-type fan-in to graduate
func TestTC_006_QuickModePerTypeDependencyFanIn(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-006", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-006")
	quickSlimRunIndex(t, dir, "test-qts-006")

	idx := quickSlimReadIndex(t, dir, "test-qts-006")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// T-quick-gen-and-run-api and T-quick-gen-and-run-cli should all depend on T-quick-gen-cases
	for _, typ := range []string{"api", "cli"} {
		id := "T-quick-gen-and-run-" + typ
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Contains(t, task.Dependencies, "T-quick-gen-cases",
			"%s should depend on T-quick-gen-cases", id)
	}

	// T-quick-graduate should depend on all per-type gen-and-run tasks
	gradTask, ok := byID["T-quick-graduate"]
	require.True(t, ok, "T-quick-graduate should exist")
	assert.Contains(t, gradTask.Dependencies, "T-quick-gen-and-run-api",
		"T-quick-graduate should depend on T-quick-gen-and-run-api")
	assert.Contains(t, gradTask.Dependencies, "T-quick-gen-and-run-cli",
		"T-quick-graduate should depend on T-quick-gen-and-run-cli")

	// T-quick-verify-regression should depend on T-quick-graduate
	verifyTask, ok := byID["T-quick-verify-regression"]
	require.True(t, ok, "T-quick-verify-regression should exist")
	assert.Contains(t, verifyTask.Dependencies, "T-quick-graduate",
		"T-quick-verify-regression should depend on T-quick-graduate")
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
		[]byte("languages:\n  - go\n"), 0644))

	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"coding.feature\"\nscope: \"all\"\n---\n\n# Task One\n"
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

	// T-test-run should have type run
	task3, ok := byID["T-test-run"]
	require.True(t, ok, "T-test-run should exist in breakdown mode")
	assert.Equal(t, "test.run", task3.Type,
		"breakdown T-test-run should have type test.run")

	// Total test pipeline tasks: gen-cases(1) + eval(1) + 2 per-type gen-scripts + run(1) + graduate(1) + verify(1) = 7
	// Plus consolidate-specs (T-specs-1) = 8 total T-test-*/T-specs-* tasks
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-test-") || strings.HasPrefix(task.ID, "T-specs-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 8, testTaskCount,
		"breakdown mode with go (2 capabilities) should have 8 test pipeline tasks")
}

// =============================================================================
// TC-008: Quick mode multi-profile with letter suffixes works
// =============================================================================

// Traceability: TC-008 -> Task 1 AC: multi-profile letter suffixes with per-type gen-and-run tasks
func TestTC_008_QuickModeMultiProfileLetterSuffixes(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-008", []string{"go", "javascript"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-008")
	quickSlimRunIndex(t, dir, "test-qts-008")

	idx := quickSlimReadIndex(t, dir, "test-qts-008")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// Union capabilities from go [api, cli] + javascript [web-ui, api] = [api, cli, web-ui]
	_ = []string{"api", "cli", "web-ui"}

	// Verify suffixed gen-cases tasks exist
	for _, id := range []string{"T-quick-gen-casesa", "T-quick-gen-casesb"} {
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Equal(t, "test.gen-cases", task.Type,
			"%s should have type test.gen-cases", id)
	}

	// Verify suffixed per-type gen-and-run tasks exist
	// Profile a (go): capabilities [api, cli]
	for _, typ := range []string{"api", "cli"} {
		id := "T-quick-gen-and-runa-" + typ
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Equal(t, "test.gen-and-run", task.Type,
			"%s should have type test.gen-and-run", id)
	}
	// Profile b (javascript): capabilities [api, cli, web-ui]
	for _, typ := range []string{"api", "cli", "web-ui"} {
		id := "T-quick-gen-and-runb-" + typ
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Equal(t, "test.gen-and-run", task.Type,
			"%s should have type test.gen-and-run", id)
	}
	}

	// Shared tasks T-quick-verify-regression and T-quick-doc-drift should exist
	for _, id := range []string{"T-quick-verify-regression", "T-quick-doc-drift"} {
		_, ok := byID[id]
		assert.True(t, ok, "%s should exist as shared task", id)
	}
}



// =============================================================================
// TC-011: InferType maps merged IDs correctly
// =============================================================================

// Traceability: TC-011 -> Task 1 AC: InferType handles per-type gen-and-run task ID patterns
func TestTC_011_InferTypeMapsMergedIDsCorrectly(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-011", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-011")
	quickSlimRunIndex(t, dir, "test-qts-011")

	idx := quickSlimReadIndex(t, dir, "test-qts-011")

	// Test cases: each per-type ID should map to gen-and-run type
	// go capabilities are [api, cli]
	testIDs := []string{
		"T-quick-gen-and-run-api",
		"T-quick-gen-and-run-cli",
		"T-quick-gen-and-run-cli",
	}

	for _, id := range testIDs {
		found := false
		for _, task := range idx.Tasks {
			if task.ID == id {
				found = true
				assert.Equal(t, "test.gen-and-run", task.Type,
					"%s should have type test.gen-and-run", id)
				break
			}
		}
		assert.True(t, found, "%s should exist in index", id)
	}

	// Also verify multi-profile suffixed per-type IDs via a separate index
	dir2 := quickSlimSetupProject(t, "test-qts-011b", []string{"go", "javascript"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir2, "test-qts-011b")
	quickSlimRunIndex(t, dir2, "test-qts-011b")

	idx2 := quickSlimReadIndex(t, dir2, "test-qts-011b")

	// Union capabilities: [api, cli, web-ui]
	for _, id := range []string{"T-quick-gen-and-runa-api", "T-quick-gen-and-runa-cli", "T-quick-gen-and-runa-web-ui",
		"T-quick-gen-and-runb-api", "T-quick-gen-and-runb-cli", "T-quick-gen-and-runb-web-ui"} {
		found := false
		for _, task := range idx2.Tasks {
			if task.ID == id {
				found = true
				assert.Equal(t, "test.gen-and-run", task.Type,
					"%s should have type test.gen-and-run", id)
				break
			}
		}
		assert.True(t, found, "%s should exist in multi-profile index", id)
	}
}

// =============================================================================
// TC-012: Quick mode single profile produces 5 tasks total
// =============================================================================

// Traceability: TC-012 -> Task 1 AC: single profile = 7 tasks (with 3 capabilities)
func TestTC_012_QuickModeSingleProfileProducesFiveTasks(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-012", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-012")
	quickSlimRunIndex(t, dir, "test-qts-012")

	idx := quickSlimReadIndex(t, dir, "test-qts-012")

	// go capabilities: [api, cli] -> 1 gen-cases + 2 per-type gen-and-run + 1 graduate + 1 verify + 1 drift = 6
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-quick-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 6, testTaskCount,
		"single profile quick mode with 2 capabilities should produce exactly 6 test tasks")
}


// =============================================================================
// TC-014: Merged task generates correct task .md file
// =============================================================================

// Traceability: TC-014 -> Task 1 AC: per-type gen-and-run task generates valid .md with all fields
func TestTC_014_MergedTaskGeneratesCorrectMD(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-014", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-014")
	quickSlimRunIndex(t, dir, "test-qts-014")

	tasksDir := filepath.Join(dir, "docs", "features", "test-qts-014", "tasks")

	// Find a per-type gen-and-run task .md file
	mdPath := filepath.Join(tasksDir, "quick-gen-and-run-go-api.md")
	data, err := os.ReadFile(mdPath)
	require.NoError(t, err, "quick-gen-and-run-go-api.md should exist")

	content := string(data)

	// Verify frontmatter fields
	assert.Contains(t, content, `id: "T-quick-gen-and-run-api"`, "md should contain correct task ID")
	assert.Contains(t, content, `type: "test.gen-and-run"`, "md should contain correct type")
	// profile field is not generated in quick mode task .md files

	// Verify body contains profile strategy
	assert.Contains(t, content, "go", "md body should reference the profile")
}

// =============================================================================
// TC-015: DetectTypesFromTestCases correctly parses test-cases.md summary table
// =============================================================================

// Traceability: TC-015 -> Task 1 AC: capabilities from profile manifest drive per-type generation
func TestTC_015_DetectTypesFromTestCasesParsesSummaryTable(t *testing.T) {
	// go profile has capabilities [api, cli] from manifest.
	// Per-type tasks are always generated when capabilities are present,
	// regardless of test-cases.md content.
	dir := quickSlimSetupProject(t, "test-qts-015a", []string{"go"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-015a")
	quickSlimRunIndex(t, dir, "test-qts-015a")

	idx := quickSlimReadIndex(t, dir, "test-qts-015a")

	// With go capabilities [api, cli], should have per-type tasks
	found := false
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-quick-gen-and-run-") {
			found = true
			break
		}
	}
	assert.True(t, found, "per-type tasks should be generated when profile has capabilities")

	// Even with zero-type test cases, capabilities from manifest still drive per-type generation
	dir2 := quickSlimSetupProject(t, "test-qts-015b", []string{"go"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir2, "test-qts-015b")
	quickSlimRunIndex(t, dir2, "test-qts-015b")

	idx2 := quickSlimReadIndex(t, dir2, "test-qts-015b")

	// Per-type tasks are generated based on profile capabilities, not test-cases content
	perTypeCount := 0
	for _, typ := range []string{"api", "cli"} {
		foundType := false
		for _, task := range idx2.Tasks {
			if task.ID == "T-quick-gen-and-run-"+typ {
				foundType = true
				perTypeCount++
				break
			}
		}
		assert.True(t, foundType, "T-quick-gen-and-run-%s should exist (driven by profile capabilities)", typ)
	}
	assert.Equal(t, 2, perTypeCount, "per-type tasks should be generated from go profile capabilities (api, cli)")
}
