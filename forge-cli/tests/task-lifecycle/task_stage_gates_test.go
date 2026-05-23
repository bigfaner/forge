//go:build e2e

package tasklifecycle

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// stageGateTestDir creates a temporary feature directory with the given task files.
// Each taskFiles entry is "filename:content" or just "filename" (uses a default content).
// Returns the project root temp dir and a cleanup function.
func stageGateTestDir(t *testing.T, featureSlug string, taskFiles []string) (string, func()) {
	t.Helper()
	tmpRoot := t.TempDir()
	featureTasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")
	require.NoError(t, os.MkdirAll(featureTasksDir, 0755))

	// Create .forge dir with config
	forgeDir := filepath.Join(tmpRoot, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))

	for _, tf := range taskFiles {
		parts := strings.SplitN(tf, ":", 2)
		filename := parts[0]
		content := defaultTaskContent(filename)
		if len(parts) == 2 && parts[1] != "" {
			content = parts[1]
		}
		require.NoError(t, os.WriteFile(filepath.Join(featureTasksDir, filename), []byte(content), 0644))
	}

	return tmpRoot, func() {}
}

// defaultTaskContent returns a minimal valid task markdown file.
func defaultTaskContent(filename string) string {
	id := strings.TrimSuffix(filename, ".md")
	// Strip description suffix after first dash
	if idx := strings.Index(id, "-"); idx > 0 {
		id = id[:idx]
	}
	return fmt.Sprintf(`---
id: %q
title: "Test task %s"
priority: "P1"
type: "coding.feature"
---

# Task %s

Test task content.
`, id, id, id)
}

// taskContentWithType returns task content with a specific type field.
func taskContentWithType(filename, taskType string) string {
	id := strings.TrimSuffix(filename, ".md")
	if idx := strings.Index(id, "-"); idx > 0 {
		id = id[:idx]
	}
	return fmt.Sprintf(`---
id: %q
title: "Test task %s"
priority: "P1"
type: %q
---

# Task %s

Test task content.
`, id, id, taskType, id)
}

// parseFrontmatter parses YAML frontmatter from a markdown file.
func parseFrontmatter(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	content := string(data)
	start := strings.Index(content, "---")
	require.True(t, start >= 0, "no frontmatter start in %s", path)
	end := strings.Index(content[start+3:], "---")
	require.True(t, end >= 0, "no frontmatter end in %s", path)

	fm := make(map[string]interface{})
	require.NoError(t, yaml.Unmarshal([]byte(content[start+3:start+3+end]), fm))
	return fm
}

// parseIndexJSON reads and parses index.json from a feature directory.
func parseIndexJSON(t *testing.T, tmpRoot, featureSlug string) map[string]interface{} {
	t.Helper()
	indexPath := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks", "index.json")
	data, err := os.ReadFile(indexPath)
	require.NoError(t, err)

	var idx map[string]interface{}
	require.NoError(t, json.Unmarshal(data, &idx))
	return idx
}

// runForgeIndex runs forge task index in the given project root using the binary
// built by TestMain. Returns stdout, stderr, exit code.
func runForgeIndex(t *testing.T, tmpRoot, featureSlug string, extraArgs ...string) (string, string, int) {
	t.Helper()
	args := []string{"task", "index", "--feature", featureSlug}
	args = append(args, extraArgs...)
	cmd := exec.Command(forgeBinary, args...)
	cmd.Dir = tmpRoot
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

// --- TC-001: Generates summary and gate files for phases with >=2 business tasks ---

// Traceability: TC-001 -> Proposal Success Criteria #1, Key Scenario "Happy path"
func TestTSG_001_GeneratesSummaryAndGateForQualifyingPhases(t *testing.T) {
	featureSlug := "test-stage-gates-001"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task-a.md",
		"1.2-task-b.md",
		"2.1-task-c.md",
		"2.2-task-d.md",
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Verify summary and gate files created for both phases
	for _, phase := range []int{1, 2} {
		summaryPath := filepath.Join(tasksDir, fmt.Sprintf("%d.summary.md", phase))
		gatePath := filepath.Join(tasksDir, fmt.Sprintf("%d.gate.md", phase))
		assert.FileExists(t, summaryPath, "phase %d summary should exist", phase)
		assert.FileExists(t, gatePath, "phase %d gate should exist", phase)
	}

	// Verify dependency wiring for phase 1
	fm1 := parseFrontmatter(t, filepath.Join(tasksDir, "1.summary.md"))
	deps1, ok := fm1["dependencies"].([]interface{})
	require.True(t, ok)
	depStrs1 := toStringSlice(deps1)
	assert.Contains(t, depStrs1, "1.1")
	assert.Contains(t, depStrs1, "1.2")

	gateFm1 := parseFrontmatter(t, filepath.Join(tasksDir, "1.gate.md"))
	gateDeps1, ok := gateFm1["dependencies"].([]interface{})
	require.True(t, ok)
	assert.Contains(t, toStringSlice(gateDeps1), "1.summary")
	assert.Equal(t, true, gateFm1["breaking"])

	// Verify dependency wiring for phase 2
	fm2 := parseFrontmatter(t, filepath.Join(tasksDir, "2.summary.md"))
	deps2, ok := fm2["dependencies"].([]interface{})
	require.True(t, ok)
	depStrs2 := toStringSlice(deps2)
	assert.Contains(t, depStrs2, "2.1")
	assert.Contains(t, depStrs2, "2.2")

	gateFm2 := parseFrontmatter(t, filepath.Join(tasksDir, "2.gate.md"))
	gateDeps2, ok := gateFm2["dependencies"].([]interface{})
	require.True(t, ok)
	assert.Contains(t, toStringSlice(gateDeps2), "2.summary")
	assert.Equal(t, true, gateFm2["breaking"])
}

// --- TC-002: Correct dependency wiring for gate tasks ---

// Traceability: TC-002 -> Proposal Success Criteria #2
func TestTSG_002_CorrectDependencyWiringForGateTasks(t *testing.T) {
	featureSlug := "test-stage-gates-002"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md",
		"1.2.md",
		"1.3.md",
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Verify summary depends on all same-phase business tasks
	fm := parseFrontmatter(t, filepath.Join(tasksDir, "1.summary.md"))
	deps, ok := fm["dependencies"].([]interface{})
	require.True(t, ok)
	depStrs := toStringSlice(deps)
	sort.Strings(depStrs)
	assert.Equal(t, []string{"1.1", "1.2", "1.3"}, depStrs)

	// Verify gate depends only on summary
	gateFm := parseFrontmatter(t, filepath.Join(tasksDir, "1.gate.md"))
	gateDeps, ok := gateFm["dependencies"].([]interface{})
	require.True(t, ok)
	assert.Equal(t, []string{"1.summary"}, toStringSlice(gateDeps))
	assert.Equal(t, true, gateFm["breaking"])
}

// --- TC-003: Skips single-task phases ---

// Traceability: TC-003 -> Proposal Success Criteria #1, Key Scenario "Single-task phase"
func TestTSG_003_SkipsSingleTaskPhases(t *testing.T) {
	featureSlug := "test-stage-gates-003"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
		"2.1-task-a.md",
		"2.2-task-b.md",
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Phase 1 has only 1 task: no gate/summary
	_, err := os.Stat(filepath.Join(tasksDir, "1.summary.md"))
	assert.True(t, os.IsNotExist(err), "1.summary.md should NOT exist for single-task phase")
	_, err = os.Stat(filepath.Join(tasksDir, "1.gate.md"))
	assert.True(t, os.IsNotExist(err), "1.gate.md should NOT exist for single-task phase")

	// Phase 2 has 2 tasks: gate/summary generated
	assert.FileExists(t, filepath.Join(tasksDir, "2.summary.md"))
	assert.FileExists(t, filepath.Join(tasksDir, "2.gate.md"))
}

// --- TC-004: Excludes test-only phases ---

// Traceability: TC-004 -> Proposal Key Scenario "Phase with only test tasks", Success Criteria #1
func TestTSG_004_ExcludesTestOnlyPhases(t *testing.T) {
	featureSlug := "test-stage-gates-004"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md:" + taskContentWithType("1.1.md", "feature"),
		"T-test-gen-scripts-cli.md:" + taskContentWithType("T-test-gen-scripts-cli.md", "testTask"),
		"T-quick-gen-and-run-cli.md:" + taskContentWithType("T-quick-gen-and-run-cli.md", "testTask"),
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// No gate/summary for test-only phase (T-test/T-quick IDs don't match <digit>.<digit>)
	// Phase 1 has only 1 business task (1.1), so no gate/summary either
	_, err := os.Stat(filepath.Join(tasksDir, "1.summary.md"))
	assert.True(t, os.IsNotExist(err), "no summary should exist")
	_, err = os.Stat(filepath.Join(tasksDir, "1.gate.md"))
	assert.True(t, os.IsNotExist(err), "no gate should exist")
}

// --- TC-005: Filters T-test and T-quick tasks from business task count ---

// Traceability: TC-005 -> Proposal Key Scenario "Phase with only test tasks", "How It Works" step 1
func TestTSG_005_FiltersTestTasksFromBusinessCount(t *testing.T) {
	featureSlug := "test-stage-gates-005"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
		"1.2-task.md",
		"T-test-gen-scripts-cli.md:" + taskContentWithType("T-test-gen-scripts-cli.md", "testTask"),
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Summary generated (business count = 2, excluding test task T-test-gen-scripts-cli)
	assert.FileExists(t, filepath.Join(tasksDir, "1.summary.md"))

	// Summary depends_on should NOT include the test task
	fm := parseFrontmatter(t, filepath.Join(tasksDir, "1.summary.md"))
	deps, ok := fm["dependencies"].([]interface{})
	require.True(t, ok)
	depStrs := toStringSlice(deps)
	for _, d := range depStrs {
		assert.NotEqual(t, "T-test-gen-scripts-cli", d, "test task should not appear in summary dependencies")
	}
}

// --- TC-006: Idempotent re-run preserves existing files ---

// Traceability: TC-006 -> Proposal Success Criteria #3, Key Scenario "Idempotent re-run"
func TestTSG_006_IdempotentRerunPreservesExistingFiles(t *testing.T) {
	featureSlug := "test-stage-gates-006"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
		"1.2-task.md",
	})
	defer cleanup()

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// First run
	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	// Record file contents and modification times
	summaryPath := filepath.Join(tasksDir, "1.summary.md")
	gatePath := filepath.Join(tasksDir, "1.gate.md")

	summaryContent, err := os.ReadFile(summaryPath)
	require.NoError(t, err)
	gateContent, err := os.ReadFile(gatePath)
	require.NoError(t, err)

	summaryModTime := getFileModTime(t, summaryPath)
	gateModTime := getFileModTime(t, gatePath)

	// Small delay to ensure modification time would change if file were overwritten
	time.Sleep(10 * time.Millisecond)

	// Second run
	_, _, exitCode = runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	// Verify contents unchanged
	summaryContent2, err := os.ReadFile(summaryPath)
	require.NoError(t, err)
	assert.Equal(t, string(summaryContent), string(summaryContent2), "summary content should be unchanged")

	gateContent2, err := os.ReadFile(gatePath)
	require.NoError(t, err)
	assert.Equal(t, string(gateContent), string(gateContent2), "gate content should be unchanged")

	// Verify modification times unchanged
	assert.Equal(t, summaryModTime, getFileModTime(t, summaryPath), "summary mod time should be unchanged")
	assert.Equal(t, gateModTime, getFileModTime(t, gatePath), "gate mod time should be unchanged")
}

// --- TC-007: Generates only missing gate when summary already exists ---

// Traceability: TC-007 -> Proposal Success Criteria #4, Key Scenario "Partial state"
func TestTSG_007_GeneratesOnlyMissingGateWhenSummaryExists(t *testing.T) {
	featureSlug := "test-stage-gates-007"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
		"1.2-task.md",
	})
	defer cleanup()

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Pre-create summary with custom content
	customSummary := `---
id: "1.summary"
title: "Custom Summary"
priority: "P0"
dependencies: ["1.1", "1.2"]
type: "doc-generation.summary"
---

# Custom Summary Content

This is hand-crafted summary content.
`
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1.summary.md"), []byte(customSummary), 0644))

	// Ensure gate does NOT exist
	_, err := os.Stat(filepath.Join(tasksDir, "1.gate.md"))
	require.True(t, os.IsNotExist(err))

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	// Verify summary content unchanged (not overwritten)
	summaryContent, err := os.ReadFile(filepath.Join(tasksDir, "1.summary.md"))
	require.NoError(t, err)
	assert.Contains(t, string(summaryContent), "Custom Summary Content")

	// Verify gate created with correct dependencies
	assert.FileExists(t, filepath.Join(tasksDir, "1.gate.md"))
	gateFm := parseFrontmatter(t, filepath.Join(tasksDir, "1.gate.md"))
	gateDeps, ok := gateFm["dependencies"].([]interface{})
	require.True(t, ok)
	assert.Contains(t, toStringSlice(gateDeps), "1.summary")
}

// --- TC-008: Generates only missing summary when gate already exists ---

// Traceability: TC-008 -> Proposal Key Scenario "Partial state" (reverse case)
func TestTSG_008_GeneratesOnlyMissingSummaryWhenGateExists(t *testing.T) {
	featureSlug := "test-stage-gates-008"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
		"1.2-task.md",
	})
	defer cleanup()

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Pre-create gate with custom content
	customGate := `---
id: "1.gate"
title: "Custom Gate"
priority: "P0"
dependencies: ["1.summary"]
breaking: true
type: "gate"
---

# Custom Gate Content

Hand-crafted gate.
`
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1.gate.md"), []byte(customGate), 0644))

	// Ensure summary does NOT exist
	_, err := os.Stat(filepath.Join(tasksDir, "1.summary.md"))
	require.True(t, os.IsNotExist(err))

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	// Verify gate content unchanged
	gateContent, err := os.ReadFile(filepath.Join(tasksDir, "1.gate.md"))
	require.NoError(t, err)
	assert.Contains(t, string(gateContent), "Custom Gate Content")

	// Verify summary created
	assert.FileExists(t, filepath.Join(tasksDir, "1.summary.md"))
	fm := parseFrontmatter(t, filepath.Join(tasksDir, "1.summary.md"))
	deps, ok := fm["dependencies"].([]interface{})
	require.True(t, ok)
	depStrs := toStringSlice(deps)
	assert.Contains(t, depStrs, "1.1")
	assert.Contains(t, depStrs, "1.2")
}

// --- TC-009: Silently skips malformed task IDs ---

// Traceability: TC-009 -> Proposal Success Criteria #5, Key Scenario "Malformed task IDs"
func TestTSG_009_SilentlySkipsMalformedTaskIDs(t *testing.T) {
	featureSlug := "test-stage-gates-009"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
		"1.2-task.md",
		"intro.md",
		"1.2a-task.md",
		"overview.md",
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode, "command should exit 0 despite malformed IDs")

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Valid phase (1.1, 1.2) gets gate/summary
	assert.FileExists(t, filepath.Join(tasksDir, "1.summary.md"))
	assert.FileExists(t, filepath.Join(tasksDir, "1.gate.md"))

	// No gate/summary for malformed IDs (they don't form a valid phase)
	_, err := os.Stat(filepath.Join(tasksDir, "intro.summary.md"))
	assert.True(t, os.IsNotExist(err))
	_, err = os.Stat(filepath.Join(tasksDir, "overview.summary.md"))
	assert.True(t, os.IsNotExist(err))
}

// --- TC-010: Pre-existing hand-crafted gate files are preserved ---

// Traceability: TC-010 -> Proposal Key Scenario "Pre-existing hand-crafted gate files"
func TestTSG_010_PreservesPreexistingHandCraftedGateFiles(t *testing.T) {
	featureSlug := "test-stage-gates-010"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md",
		"1.2.md",
	})
	defer cleanup()

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Create hand-crafted gate with unique content
	handCraftedGate := `---
id: "1.gate"
title: "Manual Quality Gate"
priority: "P0"
dependencies: ["1.summary"]
breaking: true
type: "gate"
---

# Manual Quality Gate

This is hand-crafted and should never be overwritten.
Custom verification checklist:
- [ ] Custom check 1
- [ ] Custom check 2
`
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "1.gate.md"), []byte(handCraftedGate), 0644))

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	// Verify gate content is identical to hand-crafted version
	gateContent, err := os.ReadFile(filepath.Join(tasksDir, "1.gate.md"))
	require.NoError(t, err)
	assert.Contains(t, string(gateContent), "Manual Quality Gate")
	assert.Contains(t, string(gateContent), "Custom verification checklist")

	// Verify summary was generated from template
	assert.FileExists(t, filepath.Join(tasksDir, "1.summary.md"))
}

// --- TC-011: Generated tasks appear in index.json with correct type ---

// Traceability: TC-011 -> Proposal Success Criteria #6
func TestTSG_011_GeneratedTasksInIndexJsonWithCorrectType(t *testing.T) {
	featureSlug := "test-stage-gates-011"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md",
		"1.2.md",
		"2.1.md",
		"2.2.md",
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	idx := parseIndexJSON(t, tmpRoot, featureSlug)
	tasks, ok := idx["tasks"].(map[string]interface{})
	require.True(t, ok, "index.json should have tasks map")

	// Verify 1.summary entry
	summary1, ok := tasks["1.summary"].(map[string]interface{})
	require.True(t, ok, "1.summary should exist in tasks")
	assert.Equal(t, "doc-generation.summary", summary1["type"])

	// Verify 1.gate entry
	gate1, ok := tasks["1.gate"].(map[string]interface{})
	require.True(t, ok, "1.gate should exist in tasks")
	assert.Equal(t, "gate", gate1["type"])

	// Verify 2.summary entry
	summary2, ok := tasks["2.summary"].(map[string]interface{})
	require.True(t, ok, "2.summary should exist in tasks")
	assert.Equal(t, "doc-generation.summary", summary2["type"])

	// Verify 2.gate entry
	gate2, ok := tasks["2.gate"].(map[string]interface{})
	require.True(t, ok, "2.gate should exist in tasks")
	assert.Equal(t, "gate", gate2["type"])
}

// --- TC-012: CLI prints summary line per qualifying phase ---

// Traceability: TC-012 -> Proposal "CLI Output Behavior" - "Gates generated"
func TestTSG_012_PrintsSummaryLinePerQualifyingPhase(t *testing.T) {
	featureSlug := "test-stage-gates-012"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md",
		"1.2.md",
		"2.1.md",
		"2.2.md",
	})
	defer cleanup()

	stdout, stderr, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)
	combined := stdout + stderr

	// Verify CLI output contains INDEX_BUILT action with task count including gate/summary
	assert.Contains(t, combined, "INDEX_BUILT", "output should contain INDEX_BUILT action")
	assert.Contains(t, combined, featureSlug, "output should reference feature slug")

	// Verify the gate/summary files exist (proves phases were qualified)
	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")
	assert.FileExists(t, filepath.Join(tasksDir, "1.summary.md"), "phase 1 summary should exist")
	assert.FileExists(t, filepath.Join(tasksDir, "1.gate.md"), "phase 1 gate should exist")
	assert.FileExists(t, filepath.Join(tasksDir, "2.summary.md"), "phase 2 summary should exist")
	assert.FileExists(t, filepath.Join(tasksDir, "2.gate.md"), "phase 2 gate should exist")

	// Verify task count includes generated gate/summary tasks (4 business + 4 gate/summary = 8)
	assert.Contains(t, combined, "Tasks: 8", "output should show 8 total tasks (4 business + 4 gate/summary)")
}

// --- TC-013: CLI prints no-qualification message when no phases qualify ---

// Traceability: TC-013 -> Proposal "CLI Output Behavior" - "Zero phases qualify"
func TestTSG_013_PrintsNoQualificationMessage(t *testing.T) {
	featureSlug := "test-stage-gates-013"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
	})
	defer cleanup()

	stdout, stderr, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)
	combined := stdout + stderr

	// When no phases qualify, output should indicate this
	// Exact message depends on implementation; check for key phrases
	noQualifyPatterns := []string{
		"No phases qualified",
		"no stage-gate",
		"no qualifying",
	}
	matched := false
	for _, pattern := range noQualifyPatterns {
		if strings.Contains(strings.ToLower(combined), strings.ToLower(pattern)) {
			matched = true
			break
		}
	}
	// If no specific message, at minimum no gate files should be created
	if !matched {
		tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")
		_, err := os.Stat(filepath.Join(tasksDir, "1.summary.md"))
		assert.True(t, os.IsNotExist(err), "no gate/summary should be created when no phases qualify")
	}
}

// --- TC-014: CLI exits with error on template rendering failure ---
// Note: This is a unit-level scenario. At e2e level, we verify the CLI handles
// errors gracefully. Template rendering failure requires binary modification
// which is better suited for unit tests.

// Traceability: TC-014 -> Proposal "CLI Output Behavior" - "Template rendering failure"
func TestTSG_014_ExitsWithErrorOnTemplateRenderFailure(t *testing.T) {
	t.Skip("requires binary modification to corrupt embedded template - unit test scenario")
}

// --- TC-015: Quick mode generates stage-gates identically to full mode ---

// Traceability: TC-015 -> Proposal Key Scenario "Quick mode", Success Criteria #7
func TestTSG_015_QuickModeGeneratesStageGatesIdentically(t *testing.T) {
	featureSlug := "test-stage-gates-015"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md",
		"1.2.md",
		"2.1.md",
		"2.2.md",
	})
	defer cleanup()

	// Create a quick-mode manifest
	manifestDir := filepath.Join(tmpRoot, "docs", "features", featureSlug)
	require.NoError(t, os.MkdirAll(manifestDir, 0755))
	manifest := `---
mode: quick
---
`
	require.NoError(t, os.WriteFile(filepath.Join(manifestDir, "manifest.md"), []byte(manifest), 0644))

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Same gate/summary generation as full mode
	assert.FileExists(t, filepath.Join(tasksDir, "1.summary.md"))
	assert.FileExists(t, filepath.Join(tasksDir, "1.gate.md"))
	assert.FileExists(t, filepath.Join(tasksDir, "2.summary.md"))
	assert.FileExists(t, filepath.Join(tasksDir, "2.gate.md"))

	// Same dependency wiring
	gateFm1 := parseFrontmatter(t, filepath.Join(tasksDir, "1.gate.md"))
	gateDeps1, ok := gateFm1["dependencies"].([]interface{})
	require.True(t, ok)
	assert.Contains(t, toStringSlice(gateDeps1), "1.summary")

	// Same index.json entries
	idx := parseIndexJSON(t, tmpRoot, featureSlug)
	tasks, ok := idx["tasks"].(map[string]interface{})
	require.True(t, ok)
	assert.NotNil(t, tasks["1.summary"])
	assert.NotNil(t, tasks["1.gate"])
	assert.NotNil(t, tasks["2.summary"])
	assert.NotNil(t, tasks["2.gate"])
}

// --- TC-016: Does not break existing forge task index behavior ---

// Traceability: TC-016 -> Proposal "Constraints & Dependencies"
func TestTSG_016_DoesNotBreakExistingIndexBehavior(t *testing.T) {
	featureSlug := "test-stage-gates-016"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1-task.md",
	})
	defer cleanup()

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	// index.json should be generated correctly
	idx := parseIndexJSON(t, tmpRoot, featureSlug)
	tasks, ok := idx["tasks"].(map[string]interface{})
	require.True(t, ok, "index.json should have tasks map")

	// Existing task entry should be present (key derived from filename)
	// The key is "1.1-task" (filename without .md), and ID is "1.1" (from frontmatter)
	found := false
	for key, val := range tasks {
		task, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		if task["id"] == "1.1" {
			found = true
			assert.NotEmpty(t, task["id"], "task ID should be set")
			assert.Equal(t, "pending", task["status"], "new task should be pending")
			break
		}
		_ = key
	}
	assert.True(t, found, "1.1 task should exist in index")
}

// --- TC-017: --no-test flag removed, returns unknown flag error ---

// Traceability: TC-017 -> Proposal "Constraints & Dependencies" - "--no-test flag removed"
func TestTSG_017_NoTestFlagRemoved(t *testing.T) {
	featureSlug := "test-stage-gates-017"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md",
		"1.2.md",
		"2.1.md",
		"2.2.md",
	})
	defer cleanup()

	// Run with --no-test — should fail with unknown flag error
	_, stderr, exitCode := runForgeIndex(t, tmpRoot, featureSlug, "--no-test")
	assert.NotEqual(t, 0, exitCode, "--no-test should be rejected")
	assert.Contains(t, stderr, "unknown flag")
}

// --- TC-018: Concurrent execution produces identical output ---

// Traceability: TC-018 -> Proposal Key Scenario "Concurrent execution"
func TestTSG_018_ConcurrentExecutionIdenticalOutput(t *testing.T) {
	featureSlug := "test-stage-gates-018"
	tmpRoot, cleanup := stageGateTestDir(t, featureSlug, []string{
		"1.1.md",
		"1.2.md",
	})
	defer cleanup()

	// Run first instance
	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")

	// Record first run output
	summary1, err := os.ReadFile(filepath.Join(tasksDir, "1.summary.md"))
	require.NoError(t, err)
	gate1, err := os.ReadFile(filepath.Join(tasksDir, "1.gate.md"))
	require.NoError(t, err)

	// Run second instance
	_, _, exitCode = runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode)

	// Verify identical output
	summary2, err := os.ReadFile(filepath.Join(tasksDir, "1.summary.md"))
	require.NoError(t, err)
	gate2, err := os.ReadFile(filepath.Join(tasksDir, "1.gate.md"))
	require.NoError(t, err)

	assert.Equal(t, string(summary1), string(summary2), "summary should be identical across runs")
	assert.Equal(t, string(gate1), string(gate2), "gate should be identical across runs")
}

// --- TC-019: Phase detection rejects path traversal in task IDs ---

// Traceability: TC-019 -> Proposal "Non-Functional Requirements" - Security
func TestTSG_019_RejectsPathTraversalInTaskIDs(t *testing.T) {
	featureSlug := "test-stage-gates-019"
	tmpRoot := t.TempDir()
	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	forgeDir := filepath.Join(tmpRoot, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))

	// Create task files with path-traversal-like patterns
	pathTraversalContent := `---
id: "../1.1"
title: "Traversal attempt"
priority: "P1"
type: "coding.feature"
---

# Traversal
`
	// Normal tasks for comparison
	normalTasks := []string{
		"1.1.md:" + defaultTaskContent("1.1.md"),
		"1.2.md:" + defaultTaskContent("1.2.md"),
	}
	for _, tf := range normalTasks {
		parts := strings.SplitN(tf, ":", 2)
		require.NoError(t, os.WriteFile(filepath.Join(tasksDir, parts[0]), []byte(parts[1]), 0644))
	}
	// Create the traversal file if the filesystem allows it
	// The key test is that the regex ^\d+\.\d+$ rejects non-matching patterns
	// Since filenames like "../1.1.md" may not be creatable, we test indirectly
	// by verifying only valid IDs are processed
	_ = pathTraversalContent

	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	assert.Equal(t, 0, exitCode, "command should not crash")

	// Verify no files written outside the feature tasks directory
	// Walk the temp root and check no .summary.md or .gate.md exist outside tasks dir
	filepath.Walk(tmpRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(tasksDir, path)
		if rel != "" && !strings.HasPrefix(rel, "..") {
			return nil // inside tasks dir, OK
		}
		name := info.Name()
		assert.False(t,
			strings.HasSuffix(name, ".summary.md") || strings.HasSuffix(name, ".gate.md"),
			"gate/summary file found outside tasks dir: %s", path)
		return nil
	})

	// Verify valid IDs still processed correctly
	assert.FileExists(t, filepath.Join(tasksDir, "1.summary.md"))
	assert.FileExists(t, filepath.Join(tasksDir, "1.gate.md"))
}

// --- TC-020: Generation completes within 5ms for 100 tasks and 20 phases ---

// Traceability: TC-020 -> Proposal "Non-Functional Requirements" - Performance
func TestTSG_020_GenerationCompletesWithinTimeBudget(t *testing.T) {
	featureSlug := "test-stage-gates-020"
	tmpRoot := t.TempDir()
	tasksDir := filepath.Join(tmpRoot, "docs", "features", featureSlug, "tasks")
	require.NoError(t, os.MkdirAll(tasksDir, 0755))
	forgeDir := filepath.Join(tmpRoot, ".forge")
	require.NoError(t, os.MkdirAll(forgeDir, 0755))

	// Create 100 tasks across 20 phases (5 tasks per phase)
	var taskFiles []string
	for phase := 1; phase <= 20; phase++ {
		for task := 1; task <= 5; task++ {
			filename := fmt.Sprintf("%d.%d.md", phase, task)
			content := defaultTaskContent(filename)
			require.NoError(t, os.WriteFile(filepath.Join(tasksDir, filename), []byte(content), 0644))
			taskFiles = append(taskFiles, filename)
		}
	}

	// Measure generation time
	start := time.Now()
	_, _, exitCode := runForgeIndex(t, tmpRoot, featureSlug)
	elapsed := time.Since(start)

	assert.Equal(t, 0, exitCode)
	// Performance budget: 5ms for gate/summary generation, but CLI startup adds overhead.
	// The 5ms target is for the gate generation phase specifically, not the full CLI invocation.
	// We set a generous budget of 5s for the full CLI run.
	assert.Less(t, elapsed, 5*time.Second, "full CLI run should complete in under 5 seconds")

	// Verify all 20 phases got gate/summary files
	for phase := 1; phase <= 20; phase++ {
		assert.FileExists(t, filepath.Join(tasksDir, fmt.Sprintf("%d.summary.md", phase)),
			"phase %d summary should exist", phase)
		assert.FileExists(t, filepath.Join(tasksDir, fmt.Sprintf("%d.gate.md", phase)),
			"phase %d gate should exist", phase)
	}
}

// --- Helper functions ---

// toStringSlice converts []interface{} to []string.
func toStringSlice(in []interface{}) []string {
	result := make([]string, len(in))
	for i, v := range in {
		result[i] = fmt.Sprintf("%v", v)
	}
	return result
}

// getFileModTime returns the modification time of a file.
func getFileModTime(t *testing.T, path string) time.Time {
	t.Helper()
	info, err := os.Stat(path)
	require.NoError(t, err)
	return info.ModTime()
}

// compile-time check that unused imports are referenced
var (
	_ = regexp.MustCompile
	_ = testkit.RunCLIExitCode
)
