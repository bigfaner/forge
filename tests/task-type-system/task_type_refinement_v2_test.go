//go:build cli_functional

package tasktypesystem

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Shared helpers for task-type-refinement tests (v2) ---

// trCreateFeatureDir creates a temporary feature directory with the given tasks in
// index.json under docs/features/<featureSlug>/tasks/. Returns the project root.
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

// trDefaultTaskContent returns a minimal task file content for the given filename.
func trDefaultTaskContent(filename string) string {
	id := strings.TrimSuffix(filename, ".md")
	return fmt.Sprintf("---\nid: %q\ntitle: %q\npriority: \"P1\"\ntype: \"feature\"\n---\n\n# Task %s\n", id, id, id)
}

// trTaskContentWithType returns task content with a specific type field.
func trTaskContentWithType(filename, taskType string) string {
	id := strings.TrimSuffix(filename, ".md")
	return fmt.Sprintf("---\nid: %q\ntitle: %q\npriority: \"P1\"\ntype: %q\n---\n\n# Task %s\n", id, id, taskType, id)
}

// trParseIndexJSON reads and parses index.json from a feature directory.
func trParseIndexJSON(t *testing.T, tmpRoot, featureSlug string) map[string]interface{} {
	t.Helper()
	indexPath := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	require.NoError(t, err)
	var idx map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &idx))
	return idx
}

// trGetTasksFromIndex extracts the tasks map from a parsed index.json.
func trGetTasksFromIndex(t *testing.T, idx map[string]interface{}) map[string]interface{} {
	t.Helper()
	tasks, ok := idx["tasks"].(map[string]interface{})
	require.True(t, ok, "index.json should have tasks map")
	return tasks
}

// trHasTaskWithIDPrefix checks if any task in the index has an ID matching the given prefix.
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
	cmd := exec.Command(testkit.ForgeBinary, "task", "list-types")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "forge task list-types should succeed")
	assert.Equal(t, 0, cmd.ProcessState.ExitCode(), "should exit 0")

	output := string(out)
	newTypes := []string{"feature", "enhancement", "cleanup", "refactor"}
	for _, typ := range newTypes {
		assert.True(t, strings.Contains(output, typ),
			"task list-types output should contain type: %s", typ)
	}
}

// --- TC-002: forge list-types no longer shows implementation type ---

// Traceability: TC-002 -> Task 1 AC-2, TypeImplementation removed from registry
func TestTC_TypeRefine_002_ListTypesNoLongerShowsImplementation(t *testing.T) {
	cmd := exec.Command(testkit.ForgeBinary, "task", "list-types")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err, "forge task list-types should succeed")

	output := string(out)
	assert.False(t, strings.Contains(output, "implementation"),
		"task list-types output should NOT contain 'implementation' type")
}

// --- TC-007: forge build-index skips test pipeline for cleanup-only feature ---

// Traceability: TC-007 -> Task 2 AC-1, Proposal D2 (needsTestPipeline false for cleanup)
func TestTC_TypeRefine_007_BuildIndexSkipsPipelineForCleanupOnly(t *testing.T) {
	featureSlug := "tr-pipeline-cleanup"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "cleanup"),
	})

	cmd := exec.Command(testkit.ForgeBinary, "task", "index", "--feature", featureSlug)
	cmd.Dir = tmpRoot
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task index should succeed: %s", string(out))

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

	cmd := exec.Command(testkit.ForgeBinary, "task", "index", "--feature", featureSlug)
	cmd.Dir = tmpRoot
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task index should succeed: %s", string(out))

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

	cmd := exec.Command(testkit.ForgeBinary, "task", "index", "--feature", featureSlug)
	cmd.Dir = tmpRoot
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task index should succeed: %s", string(out))

	idx := trParseIndexJSON(t, tmpRoot, featureSlug)
	tasks := trGetTasksFromIndex(t, idx)
	assert.True(t, trHasTaskWithIDPrefix(tasks, "T-review-doc"),
		"documentation-only feature should have T-review-doc auto-gen task")
	assert.False(t, trHasTaskWithIDPrefix(tasks, "T-test-"),
		"documentation-only feature should NOT have test pipeline tasks")
}

// --- TC-010: forge build-index generates neither pipeline nor review-doc for mixed cleanup-refactor ---

// Traceability: TC-010 -> Task 2 AC-5, Proposal D2 table row "only cleanup/refactor"
func TestTC_TypeRefine_010_BuildIndexNoPipelineNoEvalForCleanupRefactor(t *testing.T) {
	featureSlug := "tr-cleanup-refactor"
	tmpRoot := trCreateFeatureDir(t, featureSlug, []string{
		"biz-1.md:" + trTaskContentWithType("biz-1.md", "cleanup"),
		"biz-2.md:" + trTaskContentWithType("biz-2.md", "refactor"),
	})

	cmd := exec.Command(testkit.ForgeBinary, "task", "index", "--feature", featureSlug)
	cmd.Dir = tmpRoot
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task index should succeed: %s", string(out))

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

	cmd := exec.Command(testkit.ForgeBinary, "task", "index", "--feature", featureSlug)
	cmd.Dir = tmpRoot
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task index should succeed: %s", string(out))

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

// --- TC-015: forge creates fix-typed dynamic task on compile failure ---

// Traceability: TC-015 -> Task 4 AC-1, Proposal D4 (compile failure -> TypeCodingFix)
func TestTC_TypeRefine_015_FixTypeFromCompileFailureIsFix(t *testing.T) {
	cmd := exec.Command(testkit.ForgeBinary, "task", "list-types")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task list-types should exit 0")
	assert.True(t, strings.Contains(string(out), "fix"),
		"fix type must exist in registry for compile failure fix tasks")
}

// --- TC-016: forge creates cleanup-typed dynamic task on fmt failure ---

// Traceability: TC-016 -> Task 4 AC-2, Proposal D4 (fmt failure -> TypeCodingCleanup)
func TestTC_TypeRefine_016_FixTypeFromFmtFailureIsCleanup(t *testing.T) {
	cmd := exec.Command(testkit.ForgeBinary, "task", "list-types")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task list-types should exit 0")
	assert.True(t, strings.Contains(string(out), "cleanup"),
		"cleanup type must exist in registry for fmt failure fix tasks")
}

// --- TC-017: forge creates cleanup-typed dynamic task on lint failure ---

// Traceability: TC-017 -> Task 4 AC-2, Proposal D4 (lint failure -> TypeCodingCleanup)
func TestTC_TypeRefine_017_FixTypeFromLintFailureIsCleanup(t *testing.T) {
	cmd := exec.Command(testkit.ForgeBinary, "task", "list-types")
	out, err := cmd.CombinedOutput()
	require.NoError(t, err, "task list-types should exit 0")
	assert.True(t, strings.Contains(string(out), "cleanup"),
		"cleanup type must exist in registry for lint failure fix tasks")
}
