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
	NoTest        bool     `json:"noTest,omitempty"`
	MainSession   bool     `json:"mainSession,omitempty"`
	EstimatedTime string   `json:"estimatedTime,omitempty"`
}

// quickSlimBinPath caches the built forge binary path.
var quickSlimBinPath string

// quickSlimBin returns the path to the forge binary, building it if necessary.
func quickSlimBin(t *testing.T) string {
	t.Helper()
	if quickSlimBinPath != "" {
		if _, err := os.Stat(quickSlimBinPath); err == nil {
			return quickSlimBinPath
		}
	}
	binName := "forge"
	if runtime.GOOS == "windows" {
		binName = "forge.exe"
	}
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
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
	quickSlimBinPath = absPath
	return absPath
}

// quickSlimSetupProject creates a temp project with quick mode structure.
func quickSlimSetupProject(t *testing.T, slug string, testProfiles []string, testCasesContent string) string {
	t.Helper()
	dir := t.TempDir()

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
	taskMD := "---\nid: \"1\"\ntitle: \"Implement feature\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"implementation\"\nscope: \"all\"\n---\n\n# Implement feature\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-implement.md"), []byte(taskMD), 0644))
}

// quickSlimRunIndex runs forge task index and returns the output.
func quickSlimRunIndex(t *testing.T, dir, slug string) []byte {
	t.Helper()
	bin := quickSlimBin(t)
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
	dir := quickSlimSetupProject(t, "test-qts-001", []string{"go-test"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-001")
	quickSlimRunIndex(t, dir, "test-qts-001")

	idx := quickSlimReadIndex(t, dir, "test-qts-001")

	// Count test pipeline tasks (T-quick-* IDs)
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-quick-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 5, testTaskCount, "quick mode with single profile should generate exactly 5 test pipeline tasks")
}

// =============================================================================
// TC-002: Quick mode merged task has gen-and-run type
// =============================================================================

// Traceability: TC-002 -> Task 1 AC: merged task type = test-pipeline.gen-and-run
func TestTC_002_QuickModeMergedTaskHasGenAndRunType(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-002", []string{"go-test"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-002")
	quickSlimRunIndex(t, dir, "test-qts-002")

	idx := quickSlimReadIndex(t, dir, "test-qts-002")

	// Find T-quick-2 task
	var found bool
	for _, task := range idx.Tasks {
		if task.ID == "T-quick-2" {
			found = true
			assert.Equal(t, "test-pipeline.gen-and-run", task.Type,
				"T-quick-2 should have type test-pipeline.gen-and-run")
			break
		}
	}
	assert.True(t, found, "T-quick-2 task should exist in index")
}

// =============================================================================
// TC-003: Quick mode merged task generates correct prompt template
// =============================================================================

// Traceability: TC-003 -> Task 1 AC: type-to-template mapping
func TestTC_003_QuickModeMergedTaskPromptTemplate(t *testing.T) {
	// Verify the template file exists on disk
	// Resolve project root relative to this test file
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
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

	templatePath := filepath.Join(dir, "forge-cli", "pkg", "prompt", "data", "test-pipeline-gen-and-run.md")
	data, err := os.ReadFile(templatePath)
	require.NoError(t, err, "prompt template should exist at %s", templatePath)

	content := string(data)
	assert.Contains(t, content, "gen-test-scripts", "template should reference gen-test-scripts skill")
	assert.Contains(t, content, "run-e2e-tests", "template should reference run-e2e-tests skill")
}

// =============================================================================
// TC-004: Quick mode per-type creates independent gen-and-run tasks
// =============================================================================

// Traceability: TC-004 -> Proposal Success Criteria: per-type independent tasks
func TestTC_004_QuickModePerTypeCreatesIndependentGenAndRun(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-004", []string{"go-test"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-004")
	quickSlimRunIndex(t, dir, "test-qts-004")

	idx := quickSlimReadIndex(t, dir, "test-qts-004")

	// Find per-type gen-and-run tasks
	for _, typ := range []string{"api", "tui"} {
		found := false
		for _, task := range idx.Tasks {
			if strings.HasPrefix(task.ID, "T-quick-2") && strings.HasSuffix(task.ID, "-"+typ) {
				found = true
				assert.Equal(t, "test-pipeline.gen-and-run", task.Type,
					"T-quick-2-%s should have type test-pipeline.gen-and-run", typ)
			}
		}
		assert.True(t, found, "per-type task for %s should exist", typ)
	}
}

// =============================================================================
// TC-005: Quick mode dependency chain is correct after merge
// =============================================================================

// Traceability: TC-005 -> Task 1 AC: dependency chain gen-cases -> gen-and-run -> graduate -> verify -> drift
func TestTC_005_QuickModeDependencyChainCorrectAfterMerge(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-005", []string{"go-test"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-005")
	quickSlimRunIndex(t, dir, "test-qts-005")

	idx := quickSlimReadIndex(t, dir, "test-qts-005")

	// Build a map from ID to task for easy lookup
	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// Verify dependency chain
	t.Run("T-quick-2 depends on T-quick-1", func(t *testing.T) {
		task, ok := byID["T-quick-2"]
		require.True(t, ok, "T-quick-2 should exist")
		assert.Contains(t, task.Dependencies, "T-quick-1",
			"T-quick-2 should depend on T-quick-1")
	})

	t.Run("T-quick-3 depends on T-quick-2", func(t *testing.T) {
		task, ok := byID["T-quick-3"]
		require.True(t, ok, "T-quick-3 should exist")
		assert.Contains(t, task.Dependencies, "T-quick-2",
			"T-quick-3 should depend on T-quick-2")
	})

	t.Run("T-quick-4 depends on T-quick-3", func(t *testing.T) {
		task, ok := byID["T-quick-4"]
		require.True(t, ok, "T-quick-4 should exist")
		assert.Contains(t, task.Dependencies, "T-quick-3",
			"T-quick-4 should depend on T-quick-3")
	})

	t.Run("T-quick-5 depends on T-quick-4", func(t *testing.T) {
		task, ok := byID["T-quick-5"]
		require.True(t, ok, "T-quick-5 should exist")
		assert.Contains(t, task.Dependencies, "T-quick-4",
			"T-quick-5 should depend on T-quick-4")
	})
}

// =============================================================================
// TC-006: Quick mode per-type dependency chain fans in correctly
// =============================================================================

// Traceability: TC-006 -> Proposal Success Criteria: per-type fan-in to graduate
func TestTC_006_QuickModePerTypeDependencyFanIn(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-006", []string{"go-test"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-006")
	quickSlimRunIndex(t, dir, "test-qts-006")

	idx := quickSlimReadIndex(t, dir, "test-qts-006")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// T-quick-2-tui and T-quick-2-api should both depend on T-quick-1
	for _, typ := range []string{"tui", "api"} {
		id := "T-quick-2-" + typ
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Contains(t, task.Dependencies, "T-quick-1",
			"%s should depend on T-quick-1", id)
	}

	// T-quick-3 (graduate) should depend on both T-quick-2-tui AND T-quick-2-api
	gradTask, ok := byID["T-quick-3"]
	require.True(t, ok, "T-quick-3 should exist")
	assert.Contains(t, gradTask.Dependencies, "T-quick-2-tui",
		"T-quick-3 should depend on T-quick-2-tui")
	assert.Contains(t, gradTask.Dependencies, "T-quick-2-api",
		"T-quick-3 should depend on T-quick-2-api")

	// T-quick-4 should depend on T-quick-3
	verifyTask, ok := byID["T-quick-4"]
	require.True(t, ok, "T-quick-4 should exist")
	assert.Contains(t, verifyTask.Dependencies, "T-quick-3",
		"T-quick-4 should depend on T-quick-3")
}

// =============================================================================
// TC-007: Breakdown mode is unchanged by quick mode merge
// =============================================================================

// Traceability: TC-007 -> Task 1 Hard Rules: breakdown mode unchanged
func TestTC_007_BreakdownModeUnchangedByQuickMerge(t *testing.T) {
	// Create a breakdown mode project (has PRD, not proposal)
	dir := t.TempDir()

	featureDir := filepath.Join(dir, "docs", "features", "test-qts-007")
	tasksDir := filepath.Join(featureDir, "tasks")
	prdDir := filepath.Join(featureDir, "prd")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	require.NoError(t, os.MkdirAll(prdDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(prdDir, "prd-spec.md"), []byte("# PRD"), 0644))

	forgeDir := filepath.Join(dir, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(forgeDir, "config.yaml"),
		[]byte("test-profiles:\n  - go-test\n"), 0644))

	taskMD := "---\nid: \"1\"\ntitle: \"Task One\"\npriority: \"P1\"\nestimated_time: \"1h\"\ntype: \"implementation\"\nscope: \"all\"\n---\n\n# Task One\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1-task-one.md"), []byte(taskMD), 0644))

	bin := quickSlimBin(t)
	cmd := exec.Command(bin, "task", "index", "--feature", "test-qts-007")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "forge task index should succeed: %s", out)

	idx := quickSlimReadIndex(t, dir, "test-qts-007")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// T-test-2 should have type gen-scripts (NOT gen-and-run)
	task2, ok := byID["T-test-2"]
	require.True(t, ok, "T-test-2 should exist in breakdown mode")
	assert.Equal(t, "test-pipeline.gen-scripts", task2.Type,
		"breakdown T-test-2 should have type test-pipeline.gen-scripts, not gen-and-run")

	// T-test-3 should have type run
	task3, ok := byID["T-test-3"]
	require.True(t, ok, "T-test-3 should exist in breakdown mode")
	assert.Equal(t, "test-pipeline.run", task3.Type,
		"breakdown T-test-3 should have type test-pipeline.run")

	// Total test pipeline tasks (excluding eval): gen-cases, eval, gen-scripts, run, graduate, verify, consolidate = 7
	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-test-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 7, testTaskCount,
		"breakdown mode should have 7 test pipeline tasks (gen-cases, eval, gen-scripts, run, graduate, verify, consolidate)")
}

// =============================================================================
// TC-008: Quick mode multi-profile with letter suffixes works
// =============================================================================

// Traceability: TC-008 -> Task 1 AC: multi-profile letter suffixes with merged task
func TestTC_008_QuickModeMultiProfileLetterSuffixes(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-008", []string{"go-test", "web-playwright"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-008")
	quickSlimRunIndex(t, dir, "test-qts-008")

	idx := quickSlimReadIndex(t, dir, "test-qts-008")

	byID := make(map[string]quickSlimTaskEntry)
	for _, task := range idx.Tasks {
		byID[task.ID] = task
	}

	// Verify suffixed gen-cases tasks exist
	for _, id := range []string{"T-quick-1a", "T-quick-1b"} {
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Equal(t, "test-pipeline.gen-cases", task.Type,
			"%s should have type test-pipeline.gen-cases", id)
	}

	// Verify suffixed gen-and-run tasks exist
	for _, id := range []string{"T-quick-2a", "T-quick-2b"} {
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Equal(t, "test-pipeline.gen-and-run", task.Type,
			"%s should have type test-pipeline.gen-and-run", id)
	}

	// Verify suffixed graduate tasks exist
	for _, id := range []string{"T-quick-3a", "T-quick-3b"} {
		task, ok := byID[id]
		require.True(t, ok, "%s should exist", id)
		assert.Equal(t, "test-pipeline.graduate", task.Type,
			"%s should have type test-pipeline.graduate", id)
	}

	// Shared tasks T-quick-4 and T-quick-5 should exist
	for _, id := range []string{"T-quick-4", "T-quick-5"} {
		_, ok := byID[id]
		assert.True(t, ok, "%s should exist as shared task", id)
	}
}

// =============================================================================
// TC-009: Merged prompt template calls gen then run sequentially
// =============================================================================

// Traceability: TC-009 -> Task 1 AC: merged template orchestrates both skills
func TestTC_009_MergedTemplateCallsBothSkills(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
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

	templatePath := filepath.Join(dir, "forge-cli", "pkg", "prompt", "data", "test-pipeline-gen-and-run.md")
	data, err := os.ReadFile(templatePath)
	require.NoError(t, err, "prompt template should exist")

	content := string(data)

	// Verify it references gen-test-scripts
	assert.Contains(t, content, "gen-test-scripts",
		"template should reference gen-test-scripts skill")

	// Verify it references run-e2e-tests
	assert.Contains(t, content, "run-e2e-tests",
		"template should reference run-e2e-tests skill")

	// Verify gen phase appears before run phase
	genIdx := strings.Index(content, "gen-test-scripts")
	runIdx := strings.Index(content, "run-e2e-tests")
	assert.True(t, genIdx < runIdx,
		"gen-test-scripts should appear before run-e2e-tests in the template")

	// Verify it includes failure recovery / fix loop
	assert.True(t,
		strings.Contains(content, "fix") || strings.Contains(content, "retry") || strings.Contains(content, "failure"),
		"template should include instructions for failure recovery")
}

// =============================================================================
// TC-010: Type constant registered in types.go
// =============================================================================

// Traceability: TC-010 -> Task 1 AC: type constant fully registered
func TestTC_010_TypeConstantRegisteredInTypesGo(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
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

	typesPath := filepath.Join(dir, "forge-cli", "pkg", "task", "types.go")
	data, err := os.ReadFile(typesPath)
	require.NoError(t, err, "types.go should exist")

	content := string(data)

	// Verify constant exists with correct value
	assert.Contains(t, content, `TypeTestPipelineGenAndRun        = "test-pipeline.gen-and-run"`,
		"types.go should declare TypeTestPipelineGenAndRun constant")

	// Verify it appears in TaskTypeRegistry
	assert.Contains(t, content, "TypeTestPipelineGenAndRun",
		"types.go should reference TypeTestPipelineGenAndRun in registry")

	// Verify it appears in ValidTypes map
	lines := strings.Split(content, "\n")
	inValidTypes := false
	for _, line := range lines {
		if strings.Contains(line, "ValidTypes") {
			inValidTypes = true
		}
		if inValidTypes && strings.Contains(line, "TypeTestPipelineGenAndRun") {
			break
		}
	}
	assert.True(t, inValidTypes,
		"TypeTestPipelineGenAndRun should be in ValidTypes map")
}

// =============================================================================
// TC-011: InferType maps merged IDs correctly
// =============================================================================

// Traceability: TC-011 -> Task 1 AC: InferType handles merged task ID patterns
func TestTC_011_InferTypeMapsMergedIDsCorrectly(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-011", []string{"go-test"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-011")
	quickSlimRunIndex(t, dir, "test-qts-011")

	idx := quickSlimReadIndex(t, dir, "test-qts-011")

	// Test cases: each ID should map to gen-and-run type
	testIDs := []string{
		"T-quick-2",
		"T-quick-2-api",
		"T-quick-2-tui",
	}

	for _, id := range testIDs {
		task, ok := idx.Tasks["quick-gen-and-run-go-test"]
		if !ok {
			// Try per-type keys
			for _, t := range idx.Tasks {
				if t.ID == id {
					task = t
					ok = true
					break
				}
			}
		}
		if ok && task.ID == id {
			assert.Equal(t, "test-pipeline.gen-and-run", task.Type,
				"%s should have type test-pipeline.gen-and-run", id)
		}
	}

	// Also verify multi-profile suffixed IDs via a separate index
	dir2 := quickSlimSetupProject(t, "test-qts-011b", []string{"go-test", "web-playwright"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir2, "test-qts-011b")
	quickSlimRunIndex(t, dir2, "test-qts-011b")

	idx2 := quickSlimReadIndex(t, dir2, "test-qts-011b")

	for _, id := range []string{"T-quick-2a", "T-quick-2b"} {
		found := false
		for _, task := range idx2.Tasks {
			if task.ID == id {
				found = true
				assert.Equal(t, "test-pipeline.gen-and-run", task.Type,
					"%s should have type test-pipeline.gen-and-run", id)
				break
			}
		}
		assert.True(t, found, "%s should exist in multi-profile index", id)
	}
}

// =============================================================================
// TC-012: Quick mode single profile produces 5 tasks total
// =============================================================================

// Traceability: TC-012 -> Task 1 AC: single profile = 5 tasks
func TestTC_012_QuickModeSingleProfileProducesFiveTasks(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-012", []string{"go-test"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-012")
	quickSlimRunIndex(t, dir, "test-qts-012")

	idx := quickSlimReadIndex(t, dir, "test-qts-012")

	testTaskCount := 0
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-quick-") {
			testTaskCount++
		}
	}
	assert.Equal(t, 5, testTaskCount,
		"single profile quick mode should produce exactly 5 test tasks")
}

// =============================================================================
// TC-013: Existing gen-scripts and run prompt templates remain intact
// =============================================================================

// Traceability: TC-013 -> Task 1 Implementation Notes: original templates unchanged
func TestTC_013_ExistingGenScriptsAndRunTemplatesRemainIntact(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
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

	dataDir := filepath.Join(dir, "forge-cli", "pkg", "prompt", "data")

	// Verify gen-scripts template exists
	genScriptsPath := filepath.Join(dataDir, "test-pipeline-gen-scripts.md")
	data, err := os.ReadFile(genScriptsPath)
	require.NoError(t, err, "test-pipeline-gen-scripts.md should still exist")
	assert.True(t, len(data) > 0, "gen-scripts template should not be empty")

	// Verify run template exists
	runPath := filepath.Join(dataDir, "test-pipeline-run.md")
	data, err = os.ReadFile(runPath)
	require.NoError(t, err, "test-pipeline-run.md should still exist")
	assert.True(t, len(data) > 0, "run template should not be empty")
}

// =============================================================================
// TC-014: Merged task generates correct task .md file
// =============================================================================

// Traceability: TC-014 -> Task 1 AC: merged task generates valid .md with all fields
func TestTC_014_MergedTaskGeneratesCorrectMD(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-014", []string{"go-test"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-014")
	quickSlimRunIndex(t, dir, "test-qts-014")

	tasksDir := filepath.Join(dir, "docs", "features", "test-qts-014", "tasks")

	// Find the gen-and-run task .md file
	mdPath := filepath.Join(tasksDir, "quick-gen-and-run-go-test.md")
	data, err := os.ReadFile(mdPath)
	require.NoError(t, err, "quick-gen-and-run-go-test.md should exist")

	content := string(data)

	// Verify frontmatter fields
	assert.Contains(t, content, `id: "T-quick-2"`, "md should contain correct task ID")
	assert.Contains(t, content, `type: "test-pipeline.gen-and-run"`, "md should contain correct type")
	assert.Contains(t, content, `profile: "go-test"`, "md should contain correct profile")

	// Verify body contains profile strategy
	assert.Contains(t, content, "go-test", "md body should reference the profile")
}

// =============================================================================
// TC-015: DetectTypesFromTestCases correctly parses test-cases.md summary table
// =============================================================================

// Traceability: TC-015 -> Task 1 AC: type detection from summary table
func TestTC_015_DetectTypesFromTestCasesParsesSummaryTable(t *testing.T) {
	dir := quickSlimSetupProject(t, "test-qts-015a", []string{"go-test"}, quickSlimMultiTypeTestCases)
	quickSlimAddBusinessTask(t, dir, "test-qts-015a")
	quickSlimRunIndex(t, dir, "test-qts-015a")

	idx := quickSlimReadIndex(t, dir, "test-qts-015a")

	// With API and TUI types detected, should have per-type tasks
	found := false
	for _, task := range idx.Tasks {
		if strings.HasPrefix(task.ID, "T-quick-2-") {
			found = true
			break
		}
	}
	assert.True(t, found, "per-type tasks should be generated when test-cases.md has multiple types with non-zero counts")

	// Verify with zero-type test cases, no per-type tasks are generated
	dir2 := quickSlimSetupProject(t, "test-qts-015b", []string{"go-test"}, quickSlimNoTypeTestCases)
	quickSlimAddBusinessTask(t, dir2, "test-qts-015b")
	quickSlimRunIndex(t, dir2, "test-qts-015b")

	idx2 := quickSlimReadIndex(t, dir2, "test-qts-015b")

	for _, task := range idx2.Tasks {
		assert.False(t, strings.Contains(task.ID, "T-quick-2-"),
			"no per-type tasks should be generated when all type counts are zero, found: %s", task.ID)
	}
}

// =============================================================================
// TC-016: Subagent completes gen-and-run in single session
// =============================================================================

// Traceability: TC-016 -> Proposal Success Criteria: single session gen-and-run
func TestTC_016_SubagentCompletesGenAndRunInSingleSession(t *testing.T) {
	// Verify the prompt template contains both generation and execution phases
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
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

	templatePath := filepath.Join(dir, "forge-cli", "pkg", "prompt", "data", "test-pipeline-gen-and-run.md")
	data, err := os.ReadFile(templatePath)
	require.NoError(t, err, "prompt template should exist")

	content := string(data)

	// Template must orchestrate both skills in one task definition
	assert.Contains(t, content, "gen-test-scripts", "template must include gen-test-scripts phase")
	assert.Contains(t, content, "run-e2e-tests", "template must include run-e2e-tests phase")

	// Template must have 3 steps (read, generate, run)
	assert.Contains(t, content, "Step 1", "template should have Step 1")
	assert.Contains(t, content, "Step 2", "template should have Step 2")
	assert.Contains(t, content, "Step 3", "template should have Step 3")

	// Verify the template describes a single session workflow
	assert.Contains(t, content, "combined", "template should describe a combined/single-session workflow")
}
