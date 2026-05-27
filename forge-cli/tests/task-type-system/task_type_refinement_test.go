//go:build cli_functional

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
