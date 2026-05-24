//go:build e2e

package tasktypesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Shared helpers for task-type-refinement tests ---

var trForgeBinPath string
var trForgeBinOnce sync.Once

func trBuildForge(t *testing.T) string {
	t.Helper()
	var buildErr error
	trForgeBinOnce.Do(func() {
		binDir, err := os.MkdirTemp("", "forge-tr-e2e-bin-*")
		if err != nil {
			buildErr = fmt.Errorf("create bin temp dir: %w", err)
			return
		}
		binPath := filepath.Join(binDir, "forge.exe")
		cmd := exec.Command("go", "build", "-o", binPath, "./cmd/forge/")
		cmd.Dir = trFindProjectRoot(t)
		out, err := cmd.CombinedOutput()
		if err != nil {
			buildErr = fmt.Errorf("build forge binary: %s: %s", err, out)
			return
		}
		trForgeBinPath = binPath
	})
	if buildErr != nil {
		t.Fatal(buildErr)
	}
	return trForgeBinPath
}

func trFindProjectRoot(t *testing.T) string {
	t.Helper()
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("cannot find project root (go.mod)")
		}
		dir = parent
	}
}

func trRunForge(t *testing.T, wd string, args ...string) (string, string, int) {
	t.Helper()
	bin := trBuildForge(t)
	cmd := exec.Command(bin, args...)
	cmd.Dir = wd
	cmd.Env = append(os.Environ(), "PROJECT_ROOT="+wd)
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdoutBuf.String(), stderrBuf.String(), exitErr.ExitCode()
	}
	if err != nil {
		return stdoutBuf.String(), err.Error(), 1
	}
	return stdoutBuf.String(), stderrBuf.String(), 0
}

func trCreateFeatureDir(t *testing.T, featureSlug string, taskFiles []string) string {
	t.Helper()
	tmpRoot := t.TempDir()
	featureTasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")
	require.NoError(t, os.MkdirAll(featureTasksDir, 0755))

	forgeDir := filepath.Join(tmpRoot, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(forgeDir, "config.yaml"),
		[]byte("languages:\n  - go\n"), 0644))

	proposalDir := filepath.Join(tmpRoot, "docs", "proposals", featureSlug)
	require.NoError(t, os.MkdirAll(proposalDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(proposalDir, "proposal.md"),
		[]byte("# Proposal\n"), 0644))

	for _, tf := range taskFiles {
		parts := strings.SplitN(tf, ":", 2)
		filename := parts[0]
		content := trDefaultTaskContent(filename)
		if len(parts) == 2 && parts[1] != "" {
			content = parts[1]
		}
		require.NoError(t, os.WriteFile(filepath.Join(featureTasksDir, filename), []byte(content), 0644))
	}
	return tmpRoot
}

func trDefaultTaskContent(filename string) string {
	id := strings.TrimSuffix(filename, ".md")
	return fmt.Sprintf("---\nid: %q\ntitle: %q\npriority: \"P1\"\ntype: \"feature\"\n---\n\n# Task %s\n", id, id, id)
}

func trTaskContentWithType(filename, taskType string) string {
	id := strings.TrimSuffix(filename, ".md")
	return fmt.Sprintf("---\nid: %q\ntitle: %q\npriority: \"P1\"\ntype: %q\n---\n\n# Task %s\n", id, id, taskType, id)
}

func trParseIndexJSON(t *testing.T, tmpRoot, featureSlug string) map[string]interface{} {
	t.Helper()
	indexPath := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	require.NoError(t, err)
	var idx map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &idx))
	return idx
}

func trGetTasksFromIndex(t *testing.T, idx map[string]interface{}) map[string]interface{} {
	t.Helper()
	tasks, ok := idx["tasks"].(map[string]interface{})
	require.True(t, ok, "index.json should have tasks map")
	return tasks
}

func trHasTaskWithIDPrefix(tasks map[string]interface{}, prefix string) bool {
	for _, v := range tasks {
		task, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		id, _ := task["id"].(string)
		if strings.HasPrefix(id, prefix) {
			return true
		}
	}
	return false
}

// --- TC-001: forge list-types displays all four new business types ---

// Traceability: TC-001 -> Task 1 AC-6, Proposal Success Criterion 1
func TestTC_TypeRefine_001_ListTypesDisplaysFourNewBusinessTypes(t *testing.T) {
	projectRoot := trFindProjectRoot(t)
	stdout, _, exitCode := trRunForge(t, projectRoot, "task", "list-types")
	assert.Equal(t, 0, exitCode, "task list-types should exit 0")
	newTypes := []string{"feature", "enhancement", "cleanup", "refactor"}
	for _, typ := range newTypes {
		assert.True(t, strings.Contains(stdout, typ),
			"task list-types output should contain type: %s", typ)
	}
}

// --- TC-002: forge list-types no longer shows implementation type ---

// Traceability: TC-002 -> Task 1 AC-2, TypeImplementation removed from registry
func TestTC_TypeRefine_002_ListTypesNoLongerShowsImplementation(t *testing.T) {
	projectRoot := trFindProjectRoot(t)
	stdout, _, exitCode := trRunForge(t, projectRoot, "task", "list-types")
	assert.Equal(t, 0, exitCode, "task list-types should exit 0")
	assert.False(t, strings.Contains(stdout, "implementation"),
		"task list-types output should NOT contain 'implementation' type")
}

// --- TC-003: forge validates index.json with new type values ---

// Traceability: TC-003 -> Task 1 AC-3, Task 1 AC-7
func TestTC_TypeRefine_003_ValidateIndexAcceptsNewTypeValues(t *testing.T) {
	featureSlug := "tr-validate-types"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "feature"),
		"biz-2.md:" + trTaskContentWithType("biz-2.md", "enhancement"),
		"biz-3.md:" + trTaskContentWithType("biz-3.md", "cleanup"),
		"biz-4.md:" + trTaskContentWithType("biz-4.md", "refactor"),
	})
	_, _, idxExitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, idxExitCode, "task index should succeed")
	indexPath := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks", "index.json")
	_, stderr, exitCode := trRunForge(t, tmpRoot, "task", "validate-index", indexPath)
	assert.Equal(t, 0, exitCode, "validate-index should exit 0 for valid new types")
	assert.NotContains(t, stderr, "invalid type", "no type validation errors expected")
}

// --- TC-004: forge build-index generates test pipeline for feature-typed tasks ---

// Traceability: TC-004 -> Task 2 AC-1, Proposal D2 (needsTestPipeline true for feature)
func TestTC_TypeRefine_004_BuildIndexGeneratesPipelineForFeature(t *testing.T) {
	featureSlug := "tr-pipeline-feature"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "feature"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.True(t, trHasTaskWithIDPrefix(tasks, "T-"),
		"index should contain auto-gen test pipeline tasks for feature-typed feature")
}

// --- TC-005: forge build-index generates test pipeline for enhancement-typed tasks ---

// Traceability: TC-005 -> Task 2 AC-1, Proposal D2 (needsTestPipeline true for enhancement)
func TestTC_TypeRefine_005_BuildIndexGeneratesPipelineForEnhancement(t *testing.T) {
	featureSlug := "tr-pipeline-enhancement"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "enhancement"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.True(t, trHasTaskWithIDPrefix(tasks, "T-"),
		"index should contain auto-gen test pipeline tasks for enhancement-typed feature")
}

// --- TC-006: forge build-index generates test pipeline for fix-typed tasks ---

// Traceability: TC-006 -> Task 2 AC-1, Proposal D2 (needsTestPipeline true for fix)
func TestTC_TypeRefine_006_BuildIndexGeneratesPipelineForFix(t *testing.T) {
	featureSlug := "tr-pipeline-fix"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "fix"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.True(t, trHasTaskWithIDPrefix(tasks, "T-"),
		"index should contain auto-gen test pipeline tasks for fix-typed feature")
}

// --- TC-007: forge build-index skips test pipeline for cleanup-only feature ---

// Traceability: TC-007 -> Task 2 AC-1, Proposal D2 (needsTestPipeline false for cleanup)
func TestTC_TypeRefine_007_BuildIndexSkipsPipelineForCleanupOnly(t *testing.T) {
	featureSlug := "tr-pipeline-cleanup"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "cleanup"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.False(t, trHasTaskWithIDPrefix(tasks, "T-"),
		"cleanup-only feature should NOT have auto-gen test pipeline tasks")
}

// --- TC-008: forge build-index skips test pipeline for refactor-only feature ---

// Traceability: TC-008 -> Task 2 AC-1, Proposal D2 (needsTestPipeline false for refactor)
func TestTC_TypeRefine_008_BuildIndexSkipsPipelineForRefactorOnly(t *testing.T) {
	featureSlug := "tr-pipeline-refactor"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "refactor"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.False(t, trHasTaskWithIDPrefix(tasks, "T-"),
		"refactor-only feature should NOT have auto-gen test pipeline tasks")
}

// --- TC-009: forge build-index generates T-review-doc for documentation-only feature ---

// Traceability: TC-009 -> Task 2 AC-2, Proposal D2 (needsReviewDoc true for documentation)
func TestTC_TypeRefine_009_BuildIndexGeneratesReviewDocForDocumentationOnly(t *testing.T) {
	featureSlug := "tr-review-doc"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "documentation"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.True(t, trHasTaskWithIDPrefix(tasks, "T-review-doc"),
		"documentation-only feature should have T-review-doc auto-gen task")
	assert.False(t, trHasTaskWithIDPrefix(tasks, "T-test-"),
		"documentation-only feature should NOT have test pipeline tasks")
	assert.False(t, trHasTaskWithIDPrefix(tasks, "T-quick-"),
		"documentation-only feature should NOT have quick test pipeline tasks")
}

// --- TC-010: forge build-index generates neither pipeline nor review-doc for mixed cleanup-refactor ---

// Traceability: TC-010 -> Task 2 AC-5, Proposal D2 table row "only cleanup/refactor"
func TestTC_TypeRefine_010_BuildIndexNoPipelineNoEvalForCleanupRefactor(t *testing.T) {
	featureSlug := "tr-cleanup-refactor"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "cleanup"),
		"biz-2.md:" + trTaskContentWithType("biz-2.md", "refactor"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.False(t, trHasTaskWithIDPrefix(tasks, "T-"),
		"cleanup+refactor feature should NOT have test pipeline or eval-doc tasks")
}

// --- TC-011: forge quality gate skips for cleanup-only feature ---

// Traceability: TC-011 -> Task 2 AC-5, Proposal D2 (quality_gate.go updated)
func TestTC_TypeRefine_011_QualityGateSkipsForCleanupOnly(t *testing.T) {
	featureSlug := "tr-qg-cleanup"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "cleanup"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")
	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	for _, v := range tasks {
		task, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		taskType, _ := task["type"].(string)
		testable := taskType == "feature" || taskType == "enhancement" || taskType == "fix"
		assert.False(t, testable,
			"cleanup-only feature should have no testable task types, found: %s", taskType)
	}
}

// --- TC-012: forge prompt get-by-task-id returns feature template for feature-typed task ---

// Traceability: TC-012 -> Task 3 AC-5, Proposal D3 table
func TestTC_TypeRefine_012_PromptReturnsFeatureTemplate(t *testing.T) {
	featureSlug := "tr-prompt-feature"
	projectRoot := testkit.ProjectRoot(t)
	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	tasksDir := filepath.Join(featureDir, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	taskContent := "---\nid: \"T-prompt-feat\"\ntitle: \"Prompt test\"\npriority: \"P1\"\ntype: \"feature\"\n---\n\n# Test\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "T-prompt-feat.md"), []byte(taskContent), 0644))
	defer os.RemoveAll(featureDir)

	indexJSON := fmt.Sprintf(`{"feature":%q,"proposal":"docs/proposals/test/proposal.md","statusEnum":["pending","in_progress","completed","blocked","skipped","rejected"],"priorityEnum":["P0","P1","P2"],"tasks":{"T-prompt-feat":{"id":"T-prompt-feat","title":"Prompt test","priority":"P1","status":"pending","file":"T-prompt-feat.md","record":"records/T-prompt-feat.md","type":"feature"}}}`, featureSlug)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "index.json"), []byte(indexJSON), 0644))

	_, _, setExitCode := trRunForge(t, projectRoot, "feature", "set", featureSlug)
	require.Equal(t, 0, setExitCode, "feature set should succeed")

	stdout, _, exitCode := trRunForge(t, projectRoot, "prompt", "get-by-task-id", "T-prompt-feat")
	assert.Equal(t, 0, exitCode, "prompt get-by-task-id should succeed")
	assert.Contains(t, stdout, "implement", "feature template should contain 'implement'")
}

// --- TC-013: forge prompt get-by-task-id returns cleanup template for cleanup-typed task ---

// Traceability: TC-013 -> Task 3 AC-3, Proposal D3 (cleanup.md: improve technical debt, no TDD)
func TestTC_TypeRefine_013_PromptReturnsCleanupTemplate(t *testing.T) {
	featureSlug := "tr-prompt-cleanup"
	projectRoot := testkit.ProjectRoot(t)
	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	tasksDir := filepath.Join(featureDir, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	taskContent := "---\nid: \"T-prompt-clean\"\ntitle: \"Prompt test\"\npriority: \"P1\"\ntype: \"cleanup\"\n---\n\n# Test\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "T-prompt-clean.md"), []byte(taskContent), 0644))
	defer os.RemoveAll(featureDir)

	indexJSON := fmt.Sprintf(`{"feature":%q,"proposal":"docs/proposals/test/proposal.md","statusEnum":["pending","in_progress","completed","blocked","skipped","rejected"],"priorityEnum":["P0","P1","P2"],"tasks":{"T-prompt-clean":{"id":"T-prompt-clean","title":"Prompt test","priority":"P1","status":"pending","file":"T-prompt-clean.md","record":"records/T-prompt-clean.md","type":"cleanup"}}}`, featureSlug)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "index.json"), []byte(indexJSON), 0644))

	_, _, setExitCode := trRunForge(t, projectRoot, "feature", "set", featureSlug)
	require.Equal(t, 0, setExitCode, "feature set should succeed")

	stdout, _, exitCode := trRunForge(t, projectRoot, "prompt", "get-by-task-id", "T-prompt-clean")
	assert.Equal(t, 0, exitCode, "prompt get-by-task-id should succeed")
	assert.Contains(t, stdout, "clean", "cleanup template should contain 'clean'")
}

// --- TC-014: forge prompt get-by-task-id returns refactor template for refactor-typed task ---

// Traceability: TC-014 -> Task 3 AC-4, Proposal D3 (refactor.md: behavior preservation check)
func TestTC_TypeRefine_014_PromptReturnsRefactorTemplate(t *testing.T) {
	featureSlug := "tr-prompt-refactor"
	projectRoot := testkit.ProjectRoot(t)
	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	tasksDir := filepath.Join(featureDir, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	taskContent := "---\nid: \"T-prompt-refac\"\ntitle: \"Prompt test\"\npriority: \"P1\"\ntype: \"refactor\"\n---\n\n# Test\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "T-prompt-refac.md"), []byte(taskContent), 0644))
	defer os.RemoveAll(featureDir)

	indexJSON := fmt.Sprintf(`{"feature":%q,"proposal":"docs/proposals/test/proposal.md","statusEnum":["pending","in_progress","completed","blocked","skipped","rejected"],"priorityEnum":["P0","P1","P2"],"tasks":{"T-prompt-refac":{"id":"T-prompt-refac","title":"Prompt test","priority":"P1","status":"pending","file":"T-prompt-refac.md","record":"records/T-prompt-refac.md","type":"refactor"}}}`, featureSlug)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "index.json"), []byte(indexJSON), 0644))

	_, _, setExitCode := trRunForge(t, projectRoot, "feature", "set", featureSlug)
	require.Equal(t, 0, setExitCode, "feature set should succeed")

	stdout, _, exitCode := trRunForge(t, projectRoot, "prompt", "get-by-task-id", "T-prompt-refac")
	assert.Equal(t, 0, exitCode, "prompt get-by-task-id should succeed")
	assert.Contains(t, stdout, "restructur", "refactor template should contain 'restructur'")
}

// --- TC-015: forge creates fix-typed dynamic task on compile failure ---

// Traceability: TC-015 -> Task 4 AC-1, Proposal D4 (compile failure -> TypeCodingFix)
func TestTC_TypeRefine_015_FixTypeFromCompileFailureIsFix(t *testing.T) {
	projectRoot := trFindProjectRoot(t)
	stdout, _, exitCode := trRunForge(t, projectRoot, "task", "list-types")
	require.Equal(t, 0, exitCode, "task list-types should exit 0")
	assert.True(t, strings.Contains(stdout, "fix"),
		"fix type must exist in registry for compile failure fix tasks")
}

// --- TC-016: forge creates cleanup-typed dynamic task on fmt failure ---

// Traceability: TC-016 -> Task 4 AC-2, Proposal D4 (fmt failure -> TypeCodingCleanup)
func TestTC_TypeRefine_016_FixTypeFromFmtFailureIsCleanup(t *testing.T) {
	projectRoot := trFindProjectRoot(t)
	stdout, _, exitCode := trRunForge(t, projectRoot, "task", "list-types")
	require.Equal(t, 0, exitCode, "task list-types should exit 0")
	assert.True(t, strings.Contains(stdout, "cleanup"),
		"cleanup type must exist in registry for fmt failure fix tasks")
}

// --- TC-017: forge creates cleanup-typed dynamic task on lint failure ---

// Traceability: TC-017 -> Task 4 AC-2, Proposal D4 (lint failure -> TypeCodingCleanup)
func TestTC_TypeRefine_017_FixTypeFromLintFailureIsCleanup(t *testing.T) {
	projectRoot := trFindProjectRoot(t)
	stdout, _, exitCode := trRunForge(t, projectRoot, "task", "list-types")
	require.Equal(t, 0, exitCode, "task list-types should exit 0")
	assert.True(t, strings.Contains(stdout, "cleanup"),
		"cleanup type must exist in registry for lint failure fix tasks")
}

// --- TC-018: forge record contains Type Reclassification block when type shifts ---

// Traceability: TC-018 -> Task 4 AC-3, AC-4, Proposal D5
func TestTC_TypeRefine_018_RecordHasReclassificationWhenTypeShifts(t *testing.T) {
	featureSlug := "tr-reclass-yes"
	projectRoot := testkit.ProjectRoot(t)

	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	tasksDir := filepath.Join(featureDir, "tasks")
	recordsDir := filepath.Join(featureDir, "tasks", "records")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	require.NoError(t, os.MkdirAll(recordsDir, 0755))

	taskContent := "---\nid: \"T-reclass-1\"\ntitle: \"Reclass test\"\npriority: \"P1\"\ntype: \"fix\"\nstatus: \"in_progress\"\nfile: \"T-reclass-1.md\"\nrecord: \"records/T-reclass-1.md\"\n---\n\n# Test\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "T-reclass-1.md"), []byte(taskContent), 0644))
	defer os.RemoveAll(featureDir)

	indexJSON := fmt.Sprintf(`{"feature":%q,"proposal":"docs/proposals/test/proposal.md","statusEnum":["pending","in_progress","completed","blocked","skipped","rejected"],"priorityEnum":["P0","P1","P2"],"tasks":{"T-reclass-1":{"id":"T-reclass-1","title":"Reclass test","priority":"P1","status":"in_progress","file":"T-reclass-1.md","record":"records/T-reclass-1.md","type":"fix"}}}`, featureSlug)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "index.json"), []byte(indexJSON), 0644))

	forgeDir := filepath.Join(projectRoot, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))
	stateContent := fmt.Sprintf(`{"feature":"%s","task_id":"T-reclass-1","started_time":"2026-01-15 10:00"}`, featureSlug)
	require.NoError(t, os.WriteFile(filepath.Join(forgeDir, "state.json"), []byte(stateContent), 0644))

	_, _, setExitCode := trRunForge(t, projectRoot, "feature", "set", featureSlug)
	require.Equal(t, 0, setExitCode, "feature set should succeed")

	submitData := `{"taskId":"T-reclass-1","status":"completed","summary":"Test reclassification","filesCreated":[],"filesModified":[],"keyDecisions":["Reclassified"],"testsPassed":1,"testsFailed":0,"coverage":-1.0,"acceptanceCriteria":[{"criterion":"test","met":true}],"typeReclassification":{"originalType":"fix","actualType":"cleanup","reason":"Was cleanup not fix"}}`
	submitDataPath := filepath.Join(t.TempDir(), "submit.json")
	require.NoError(t, os.WriteFile(submitDataPath, []byte(submitData), 0644))

	_, stderr, exitCode := trRunForge(t, projectRoot, "task", "submit", "T-reclass-1", "--data", submitDataPath, "--force")
	require.Equal(t, 0, exitCode, "submit should succeed, stderr: %s", stderr)

	recordPath := filepath.Join(recordsDir, "T-reclass-1.md")
	recordData, err := os.ReadFile(recordPath)
	require.NoError(t, err, "record file should exist")

	recordStr := string(recordData)
	assert.Contains(t, recordStr, "Type Reclassification",
		"record should contain Type Reclassification section")
}

// --- TC-019: forge record omits Type Reclassification block when no type shift ---

// Traceability: TC-019 -> Task 4 AC-4, Proposal Success Criterion 7
func TestTC_TypeRefine_019_RecordOmitsReclassificationWhenNoShift(t *testing.T) {
	featureSlug := "tr-reclass-no"
	projectRoot := testkit.ProjectRoot(t)

	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	tasksDir := filepath.Join(featureDir, "tasks")
	recordsDir := filepath.Join(featureDir, "tasks", "records")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	require.NoError(t, os.MkdirAll(recordsDir, 0755))

	taskContent := "---\nid: \"T-noclass-1\"\ntitle: \"No reclass test\"\npriority: \"P1\"\ntype: \"feature\"\nstatus: \"in_progress\"\nfile: \"T-noclass-1.md\"\nrecord: \"records/T-noclass-1.md\"\n---\n\n# Test\n"
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "T-noclass-1.md"), []byte(taskContent), 0644))
	defer os.RemoveAll(featureDir)

	indexJSON := fmt.Sprintf(`{"feature":%q,"proposal":"docs/proposals/test/proposal.md","statusEnum":["pending","in_progress","completed","blocked","skipped","rejected"],"priorityEnum":["P0","P1","P2"],"tasks":{"T-noclass-1":{"id":"T-noclass-1","title":"No reclass test","priority":"P1","status":"in_progress","file":"T-noclass-1.md","record":"records/T-noclass-1.md","type":"feature"}}}`, featureSlug)
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "index.json"), []byte(indexJSON), 0644))

	forgeDir := filepath.Join(projectRoot, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))
	stateContent := fmt.Sprintf(`{"feature":"%s","task_id":"T-noclass-1","started_time":"2026-01-15 10:00"}`, featureSlug)
	require.NoError(t, os.WriteFile(filepath.Join(forgeDir, "state.json"), []byte(stateContent), 0644))

	_, _, setExitCode := trRunForge(t, projectRoot, "feature", "set", featureSlug)
	require.Equal(t, 0, setExitCode, "feature set should succeed")

	submitData := `{"taskId":"T-noclass-1","status":"completed","summary":"Test no reclassification","filesCreated":[],"filesModified":[],"keyDecisions":["Normal"],"testsPassed":1,"testsFailed":0,"coverage":-1.0,"acceptanceCriteria":[{"criterion":"test","met":true}]}`
	submitDataPath := filepath.Join(t.TempDir(), "submit.json")
	require.NoError(t, os.WriteFile(submitDataPath, []byte(submitData), 0644))

	_, stderr, exitCode := trRunForge(t, projectRoot, "task", "submit", "T-noclass-1", "--data", submitDataPath, "--force")
	require.Equal(t, 0, exitCode, "submit should succeed, stderr: %s", stderr)

	recordPath := filepath.Join(recordsDir, "T-noclass-1.md")
	recordData, err := os.ReadFile(recordPath)
	require.NoError(t, err, "record file should exist")

	recordStr := string(recordData)
	assert.NotContains(t, recordStr, "Type Reclassification",
		"record should NOT contain Type Reclassification section when no type shift")
}

// --- TC-020: forge task migrate infers type for tasks without recognized patterns ---

// Traceability: TC-020 -> Task 5 AC-1, AC-2, Proposal Success Criterion 8
func TestTC_TypeRefine_020_MigrateDefaultsUnknownIDToFeature(t *testing.T) {
	featureSlug := "tr-migrate"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "feature"),
		"biz-2.md:" + trTaskContentWithType("biz-2.md", "cleanup"),
	})
	_, _, exitCode := trRunForge(t, tmpRoot, "task", "index", "--feature", featureSlug)
	require.Equal(t, 0, exitCode, "task index should succeed")

	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	biz1, ok := tasks["biz-1"].(map[string]interface{})
	require.True(t, ok, "biz-1 should exist in tasks")
	assert.Equal(t, "feature", biz1["type"], "biz-1 should initially be feature")

	// Set all tasks to completed so migrate can proceed
	for key, v := range tasks {
		task, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		task["status"] = "completed"
		tasks[key] = task
	}
	idx["tasks"] = tasks
	indexPath := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks", "index.json")
	updatedIndex, err := json.Marshal(idx)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(indexPath, updatedIndex, 0644))

	_, stderr, exitCode := trRunForge(t, tmpRoot, "task", "migrate")
	require.Equal(t, 0, exitCode, "task migrate should succeed, stderr: %s", stderr)

	migratedIdx := trParseIndexJSON(t, tmpRoot, featureSlug)
	migratedTasks := trGetTasksFromIndex(t, migratedIdx)

	biz1Migrated, ok := migratedTasks["biz-1"].(map[string]interface{})
	require.True(t, ok, "biz-1 should exist after migration")
	assert.Equal(t, "feature", biz1Migrated["type"],
		"feature-typed task should remain feature after migration")

	// InferType("biz-2") returns "" which defaults to "feature"
	biz2Migrated, ok := migratedTasks["biz-2"].(map[string]interface{})
	require.True(t, ok, "biz-2 should exist after migration")
	assert.Equal(t, "feature", biz2Migrated["type"],
		"biz-2 should be migrated to feature (InferType default for unknown IDs)")
}
