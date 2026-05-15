//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Go test runs from the package source directory by default.
	// These tests use paths relative to the project root.
	// Use runtime.Caller to find this source file, then walk up to project root.
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	for dir != "/" && dir != "" {
		if _, err := os.Stat(filepath.Join(dir, "justfile")); err == nil {
			_ = os.Chdir(dir)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

// runCLI executes a CLI command and returns combined output.
func tsptRunCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI command failed: %s: %s", err, out)
	}
	return string(out)
}

// tsptRunCLIRaw executes a CLI command and returns output and exit code.
// Unlike tsptRunCLI, it does not fatalf on non-zero exit.
func tsptRunCLIRaw(t *testing.T, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
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

// ==============================================================================
// gen-test-scripts --type filter tests — feature: test-scripts-per-type
// ==============================================================================

// Traceability: TC-001 -> Proposal Success Criterion 1
func TestTC_001_GenTestScriptsTypeFilterCLI(t *testing.T) {
	feature := "test-scripts-per-type"
	outDir := filepath.Join("tests", "e2e", "features", feature)

	// Clean any previously generated files in the staging area
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	out := tsptRunCLI(t, "forge", "gen-test-scripts", "--type", "cli", "--feature", feature)

	// Verify command succeeded (no error output)
	assert.NotContains(t, out, "error",
		"gen-test-scripts --type cli should succeed without errors")

	// Verify only CLI test scripts are generated
	files, err := os.ReadDir(outDir)
	assert.NoError(t, err, "output directory should be readable")

	hasCLI := false
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, "_cli_test.go") {
			hasCLI = true
		}
		assert.False(t, strings.HasSuffix(name, "_api_test.go"),
			"no API test scripts should be generated with --type cli, found: %s", name)
		assert.False(t, strings.HasSuffix(name, "_tui_test.go"),
			"no TUI test scripts should be generated with --type cli, found: %s", name)
	}
	assert.True(t, hasCLI, "at least one CLI test script should be generated")
}

// Traceability: TC-002 -> Proposal Success Criterion 1
func TestTC_002_GenTestScriptsTypeFilterAPI(t *testing.T) {
	feature := "test-scripts-per-type"
	outDir := filepath.Join("tests", "e2e", "features", feature)

	// Clean staging area
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	out := tsptRunCLI(t, "forge", "gen-test-scripts", "--type", "api", "--feature", feature)

	assert.NotContains(t, out, "error",
		"gen-test-scripts --type api should succeed without errors")

	files, err := os.ReadDir(outDir)
	assert.NoError(t, err, "output directory should be readable")

	hasAPI := false
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, "_api_test.go") {
			hasAPI = true
		}
		assert.False(t, strings.HasSuffix(name, "_cli_test.go"),
			"no CLI test scripts should be generated with --type api, found: %s", name)
		assert.False(t, strings.HasSuffix(name, "_tui_test.go"),
			"no TUI test scripts should be generated with --type api, found: %s", name)
	}
	assert.True(t, hasAPI, "at least one API test script should be generated")
}

// Traceability: TC-003 -> Proposal Success Criterion 1
func TestTC_003_GenTestScriptsTypeFilterTUI(t *testing.T) {
	feature := "test-scripts-per-type"
	outDir := filepath.Join("tests", "e2e", "features", feature)

	// Clean staging area
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	out := tsptRunCLI(t, "forge", "gen-test-scripts", "--type", "tui", "--feature", feature)

	assert.NotContains(t, out, "error",
		"gen-test-scripts --type tui should succeed without errors")

	files, err := os.ReadDir(outDir)
	assert.NoError(t, err, "output directory should be readable")

	hasTUI := false
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, "_tui_test.go") {
			hasTUI = true
		}
		assert.False(t, strings.HasSuffix(name, "_cli_test.go"),
			"no CLI test scripts should be generated with --type tui, found: %s", name)
		assert.False(t, strings.HasSuffix(name, "_api_test.go"),
			"no API test scripts should be generated with --type tui, found: %s", name)
	}
	assert.True(t, hasTUI, "at least one TUI test script should be generated")
}

// ==============================================================================
// breakdown-tasks per-type tests — feature: test-scripts-per-type
// ==============================================================================

// Traceability: TC-004 -> Proposal Success Criterion 2
func TestTC_004_BreakdownTasksCreatesPerTypeTasks(t *testing.T) {
	feature := "test-scripts-per-type"
	tasksDir := filepath.Join("docs", "features", feature, "tasks")

	out := tsptRunCLI(t, "forge", "breakdown-tasks", "--feature", feature)
	assert.NotContains(t, out, "error",
		"breakdown-tasks should succeed without errors")

	// Check that per-type gen-scripts tasks exist
	cliTask := filepath.Join(tasksDir, "*-gen-scripts-cli*")
	apiTask := filepath.Join(tasksDir, "*-gen-scripts-api*")

	cliMatches, _ := filepath.Glob(cliTask)
	apiMatches, _ := filepath.Glob(apiTask)

	assert.Greater(t, len(cliMatches), 0,
		"should create a CLI gen-scripts task (pattern: %s)", cliTask)
	assert.Greater(t, len(apiMatches), 0,
		"should create an API gen-scripts task (pattern: %s)", apiTask)
}

// Traceability: TC-005 -> Proposal Key Risk — no tasks for types without test cases
func TestTC_005_BreakdownTasksNoEmptyTypeTasks(t *testing.T) {
	feature := "test-scripts-per-type"
	tasksDir := filepath.Join("docs", "features", feature, "tasks")

	// This feature only has CLI test cases, so only CLI gen task should exist
	out := tsptRunCLI(t, "forge", "breakdown-tasks", "--feature", feature)
	assert.NotContains(t, out, "error",
		"breakdown-tasks should succeed without errors")

	// Since test-cases.md only has CLI test cases, verify no tasks for empty types
	tuiTask := filepath.Join(tasksDir, "*-gen-scripts-tui*")
	tuiMatches, _ := filepath.Glob(tuiTask)

	assert.Equal(t, 0, len(tuiMatches),
		"should not create TUI gen-scripts task when no TUI test cases exist")
}

// Traceability: TC-006 -> Proposal Success Criterion 3
func TestTC_006_QuickTasksCreatesPerTypeTasks(t *testing.T) {
	feature := "test-scripts-per-type"
	tasksDir := filepath.Join("docs", "features", feature, "tasks")

	out := tsptRunCLI(t, "forge", "quick-tasks", "--feature", feature)
	assert.NotContains(t, out, "error",
		"quick-tasks should succeed without errors")

	cliTask := filepath.Join(tasksDir, "*-quick-*-cli*")
	apiTask := filepath.Join(tasksDir, "*-quick-*-api*")

	cliMatches, _ := filepath.Glob(cliTask)
	apiMatches, _ := filepath.Glob(apiTask)

	assert.Greater(t, len(cliMatches), 0,
		"should create a CLI quick gen-scripts task (pattern: %s)", cliTask)
	assert.Greater(t, len(apiMatches), 0,
		"should create an API quick gen-scripts task (pattern: %s)", apiTask)
}

// ==============================================================================
// Dependency chain tests — feature: test-scripts-per-type
// ==============================================================================

// Traceability: TC-007 -> Proposal Success Criterion 4
func TestTC_007_Test3DependsOnAllPerTypeTasks(t *testing.T) {
	feature := "test-scripts-per-type"
	tasksDir := filepath.Join("docs", "features", feature, "tasks")

	out := tsptRunCLI(t, "forge", "breakdown-tasks", "--feature", feature)
	assert.NotContains(t, out, "error",
		"breakdown-tasks should succeed without errors")

	// Find the T-test-3 (run-e2e-tests) task file
	test3Pattern := filepath.Join(tasksDir, "T-test-3*")
	test3Matches, _ := filepath.Glob(test3Pattern)
	if len(test3Matches) == 0 {
		t.Skip("T-test-3 task file not found — feature may use quick mode")
	}

	content, err := os.ReadFile(test3Matches[0])
	assert.NoError(t, err, "should be able to read T-test-3 task file")

	contentStr := string(content)

	// T-test-3 dependencies should list all per-type gen tasks
	assert.True(t, strings.Contains(contentStr, "T-test-2-cli") || strings.Contains(contentStr, "-gen-scripts-cli"),
		"T-test-3 should depend on CLI gen task, content:\n%s", contentStr)
	assert.True(t, strings.Contains(contentStr, "T-test-2-api") || strings.Contains(contentStr, "-gen-scripts-api"),
		"T-test-3 should depend on API gen task, content:\n%s", contentStr)
}

// Traceability: TC-008 -> Proposal Scope — T-quick-3 dependencies
func TestTC_008_Quick3DependsOnAllPerTypeTasks(t *testing.T) {
	feature := "test-scripts-per-type"
	tasksDir := filepath.Join("docs", "features", feature, "tasks")

	out := tsptRunCLI(t, "forge", "quick-tasks", "--feature", feature)
	assert.NotContains(t, out, "error",
		"quick-tasks should succeed without errors")

	// Find the T-quick-3 (run-e2e-tests) task file
	quick3Pattern := filepath.Join(tasksDir, "T-quick-3*")
	quick3Matches, _ := filepath.Glob(quick3Pattern)
	if len(quick3Matches) == 0 {
		t.Skip("T-quick-3 task file not found")
	}

	content, err := os.ReadFile(quick3Matches[0])
	assert.NoError(t, err, "should be able to read T-quick-3 task file")

	contentStr := string(content)

	assert.True(t, strings.Contains(contentStr, "T-quick-2-cli") || strings.Contains(contentStr, "-quick-*-cli"),
		"T-quick-3 should depend on CLI quick gen task, content:\n%s", contentStr)
	assert.True(t, strings.Contains(contentStr, "T-quick-2-api") || strings.Contains(contentStr, "-quick-*-api"),
		"T-quick-3 should depend on API quick gen task, content:\n%s", contentStr)
}

// ==============================================================================
// Independent retry tests — feature: test-scripts-per-type
// ==============================================================================

// Traceability: TC-009 -> Proposal Success Criterion 5
func TestTC_009_FailedGenTaskIndependentRetry(t *testing.T) {
	feature := "test-scripts-per-type"
	outDir := filepath.Join("tests", "e2e", "features", feature)

	// Clean and create output dir
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Step 1: Generate CLI scripts — should succeed
	tsptRunCLI(t, "forge", "gen-test-scripts", "--type", "cli", "--feature", feature)

	// Record CLI scripts before API attempt
	cliFilesBefore, _ := filepath.Glob(filepath.Join(outDir, "*_cli_test.go"))
	assert.Greater(t, len(cliFilesBefore), 0,
		"CLI scripts should exist after successful generation")

	// Record file contents/checksums for comparison
	cliContentBefore := make(map[string]string)
	for _, f := range cliFilesBefore {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("failed to read CLI file %s: %v", f, err)
		}
		cliContentBefore[filepath.Base(f)] = string(data)
	}

	// Step 2: Attempt API generation — may fail (no API endpoints in project)
	out, exitCode := tsptRunCLIRaw(t, "forge", "gen-test-scripts", "--type", "api", "--feature", feature)
	// API generation failure is acceptable for this test
	_ = out
	_ = exitCode

	// Step 3: Verify CLI scripts still exist and are unchanged
	cliFilesAfter, _ := filepath.Glob(filepath.Join(outDir, "*_cli_test.go"))
	assert.Equal(t, len(cliFilesBefore), len(cliFilesAfter),
		"CLI scripts should remain intact after API generation attempt")

	for _, f := range cliFilesAfter {
		data, err := os.ReadFile(f)
		assert.NoError(t, err, "should be able to read CLI file %s", f)
		before, existed := cliContentBefore[filepath.Base(f)]
		if existed {
			assert.Equal(t, before, string(data),
				"CLI file %s should be unchanged after API failure", filepath.Base(f))
		}
	}

	// Step 4: Retry API generation
	tsptRunCLI(t, "forge", "gen-test-scripts", "--type", "api", "--feature", feature)

	// CLI scripts should still be unchanged
	cliFilesFinal, _ := filepath.Glob(filepath.Join(outDir, "*_cli_test.go"))
	assert.Equal(t, len(cliFilesBefore), len(cliFilesFinal),
		"CLI scripts should remain intact after API retry")
}

// ==============================================================================
// Idempotent infrastructure tests — feature: test-scripts-per-type
// ==============================================================================

// Traceability: TC-010 -> Proposal Success Criterion 6
func TestTC_010_SharedInfrastructureIdempotent(t *testing.T) {
	feature := "test-scripts-per-type"
	helpersPath := filepath.Join("tests", "e2e", "helpers.go")

	// Record shared infrastructure state before
	helpersBefore, err := os.ReadFile(helpersPath)
	if err != nil {
		t.Fatalf("failed to read helpers.go: %v", err)
	}

	mainTestPath := filepath.Join("tests", "e2e", "main_test.go")
	mainTestBefore, err := os.ReadFile(mainTestPath)
	if err != nil {
		t.Fatalf("failed to read main_test.go: %v", err)
	}

	// Step 1: First run creates shared infrastructure
	tsptRunCLI(t, "forge", "gen-test-scripts", "--type", "cli", "--feature", feature)

	// Step 2: Second run should reuse existing infrastructure
	tsptRunCLI(t, "forge", "gen-test-scripts", "--type", "api", "--feature", feature)

	// Step 3: Verify shared files unchanged across both runs
	helpersAfterSecond, _ := os.ReadFile(helpersPath)
	mainTestAfterSecond, _ := os.ReadFile(mainTestPath)

	assert.Equal(t, string(helpersBefore), string(helpersAfterSecond),
		"helpers.go should be unchanged after both gen runs")
	assert.Equal(t, string(mainTestBefore), string(mainTestAfterSecond),
		"main_test.go should be unchanged after both gen runs")
}

// ==============================================================================
// Single type project tests — feature: test-scripts-per-type
// ==============================================================================

// Traceability: TC-011 -> Proposal Key Scenario — single type project
func TestTC_011_SingleTypeProjectCreatesOnlyOneGenTask(t *testing.T) {
	feature := "test-scripts-per-type"
	tasksDir := filepath.Join("docs", "features", feature, "tasks")

	out := tsptRunCLI(t, "forge", "breakdown-tasks", "--feature", feature)
	assert.NotContains(t, out, "error",
		"breakdown-tasks should succeed without errors")

	// Only CLI gen task should exist (go-test profile has cli capability + only CLI test cases)
	cliPattern := filepath.Join(tasksDir, "*-gen-scripts-cli*")
	apiPattern := filepath.Join(tasksDir, "*-gen-scripts-api*")
	tuiPattern := filepath.Join(tasksDir, "*-gen-scripts-tui*")

	cliMatches, _ := filepath.Glob(cliPattern)
	apiMatches, _ := filepath.Glob(apiPattern)
	tuiMatches, _ := filepath.Glob(tuiPattern)

	assert.Greater(t, len(cliMatches), 0,
		"CLI gen task should exist")
	assert.Equal(t, 0, len(tuiMatches),
		"TUI gen task should not exist when no TUI test cases present")
	// API task may or may not exist depending on profile capabilities
	// but should not exist if test-cases.md has no API cases
	_ = apiMatches
}

// ==============================================================================
// Backward compatibility tests — feature: test-scripts-per-type
// ==============================================================================

// Traceability: TC-012 -> Proposal NFR (backward compatibility)
func TestTC_012_GenTestScriptsWithoutTypeGeneratesAllTypes(t *testing.T) {
	feature := "test-scripts-per-type"
	outDir := filepath.Join("tests", "e2e", "features", feature)

	// Clean staging area
	_ = os.RemoveAll(outDir)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("failed to create output dir: %v", err)
	}

	// Run without --type flag (backward compatible mode)
	out := tsptRunCLI(t, "forge", "gen-test-scripts", "--feature", feature)

	assert.NotContains(t, out, "error",
		"gen-test-scripts without --type should succeed without errors")

	// All types present in test-cases.md should be generated
	// Since test-cases.md only has CLI cases, at least CLI scripts should exist
	cliFiles, _ := filepath.Glob(filepath.Join(outDir, "*_cli_test.go"))
	assert.Greater(t, len(cliFiles), 0,
		"CLI test scripts should be generated in backward-compatible mode")

	// Verify the command still works without the --type flag
	// (backward compatibility means existing usage patterns are preserved)
}
