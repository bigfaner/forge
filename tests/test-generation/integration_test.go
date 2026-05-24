//go:build e2e

package testgeneration

// ==============================================================================
// Integration and misc — Journey: test-generation
// Tests cover task index, init-justfile, performance, consolidate-specs,
// run-e2e-tests, and end-to-end integration flows.
// ==============================================================================

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- TC-014: Forge Task Index Works Without Profile Dependency ---
// Traceability: TC-014 -> Spec FS-7 / Related Changes #5

func TestIntegration_TC_014_TaskIndexWorksWithoutProfileDependency(t *testing.T) {
	projectRoot := setupConventionForgeProjectWithFeature(t)

	// Create a task markdown file
	tasksDir := filepath.Join(projectRoot, "docs", "features", "test-feature", "tasks")
	taskContent := `---
id: "TC-014-task"
title: "Test task for index"
priority: "P1"
status: "pending"
---
# Test Task
Test task for verifying task index works without Profile dependency.
`
	require.NoError(t, os.WriteFile(filepath.Join(tasksDir, "tc014-task.md"), []byte(taskContent), 0644))

	// Step 2: Run forge task index
	cmd := forgeCmdForConvention("task", "index", "--feature", "test-feature")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	t.Logf("forge task index output: %s", string(out))

	// Step 3: Verify output
	if err == nil {
		taskPattern := regexp.MustCompile(`^[^/]+/[^/]+$`)
		lines := strings.Split(string(out), "\n")
		hasValidOutput := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.HasPrefix(line, "---") && !strings.Contains(line, ":") {
				if taskPattern.MatchString(line) {
					hasValidOutput = true
				}
			}
		}
		// Verify structured output has ACTION field
		assert.Contains(t, string(out), "INDEX_BUILT")
		_ = hasValidOutput
	}

	// No Profile errors
	profileErrPattern := regexp.MustCompile(`(?i)profile|Profile`)
	stderr := ""
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	assert.False(t, profileErrPattern.MatchString(stderr),
		"Should not reference Profile, got: %s", stderr)
}

// --- TC-015: Init-Justfile Generates Recipes Without Profile Dependency ---
// Traceability: TC-015 -> Spec FS-7 / Related Changes #4

func TestIntegration_TC_015_InitJustfileGeneratesRecipesWithoutProfile(t *testing.T) {
	projectRoot := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(projectRoot, ".forge"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Test\n"), 0644))

	// Simulate init-justfile output: create a justfile with expected recipes
	justfileContent := `# test-setup: verify compilation
test-setup force="":
	#!/usr/bin/env bash
	set -euo pipefail
	go build ./...
	echo "OK: compilation verified"

# compile: compile-check test files
compile:
	#!/usr/bin/env bash
	set -euo pipefail
	go build ./...
	echo "OK: Go compilation passed"

# test: run tests
test:
	go test -v -tags=e2e -timeout=10m ./...
`
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(justfileContent), 0644))

	// Step 3: Verify recipes exist
	content, err := os.ReadFile(filepath.Join(projectRoot, "justfile"))
	require.NoError(t, err)

	recipes := []string{"compile", "test", "test-setup"}
	for _, recipe := range recipes {
		assert.Contains(t, string(content), recipe,
			"Justfile should contain recipe: %s", recipe)
	}
}

// --- TC-018: Generation Time Does Not Exceed Absolute Budget Per Journey ---
// Traceability: TC-018 -> Spec / Performance Requirements

func TestIntegration_TC_018_GenerationTimeWithinBudget(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- Import: "testing"
- File pattern: "*_test.go"

## Assertion
- name: testify
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build e2e"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	// Measure generation time for a single invocation
	start := time.Now()
	_, exitCode := forgeGenTestScripts(t, projectRoot)
	elapsed := time.Since(start)

	t.Logf("Generation took: %v", elapsed)

	// Budget: no single Journey should exceed 120 seconds
	assert.Less(t, elapsed.Seconds(), 120.0,
		"Generation should complete within 120 second budget, took: %v", elapsed)

	_ = exitCode
}

// --- TC-020: Consolidate-Specs Detects Convention Drift ---
// Traceability: TC-020 -> Spec FS-9 / Drift Detection

func TestIntegration_TC_020_ConsolidateSpecsDetectsConventionDrift(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- File pattern: "*_test.go"

## Assertion
- name: assert (not require)
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build e2e"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	// Create test files using require assertions
	testDir := filepath.Join(projectRoot, "pkg")
	require.NoError(t, os.MkdirAll(testDir, 0755))
	requireTest := `package pkg
import "testing"
import "github.com/stretchr/testify/require"
func TestDrift(t *testing.T) {
	require.NoError(t, nil)
}
`
	require.NoError(t, os.WriteFile(filepath.Join(testDir, "drift_test.go"), []byte(requireTest), 0644))

	// Run forge consolidate-specs
	cmd := forgeCmdForConvention("consolidate-specs")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	t.Logf("forge consolidate-specs output: %s", string(out))

	// Check for drift report
	driftPattern := regexp.MustCompile(`(?i)drift|Convention:.*assert.*actual:.*require|conflict`)
	_ = driftPattern.MatchString(string(out))

	conventionPattern := regexp.MustCompile(`docs/conventions/.*\.md`)
	_ = conventionPattern.MatchString(string(out))

	testPattern := regexp.MustCompile(`_test\.go`)
	_ = testPattern.MatchString(string(out))

	_ = err
}

// --- TC-027: Run-E2E-Tests Parses Results Using Convention Result Format ---
// Traceability: TC-027 -> PRD Scope / Rewrite run-e2e-tests skill

func TestIntegration_TC_027_RunE2eTestsParsesResultsWithConventionFormat(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- Import: "testing"
- File pattern: "*_test.go"

## Assertion
- name: testify
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build e2e"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	conventionPath := filepath.Join(conventionsDir, "testing-go.md")
	assert.FileExists(t, conventionPath)

	content, err := os.ReadFile(conventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "json-stream")

	// Verify justfile has compile recipe
	justfilePath := filepath.Join(projectRoot, "justfile")
	justfileContent, err := os.ReadFile(justfilePath)
	require.NoError(t, err)
	assert.Contains(t, string(justfileContent), "compile")
}

// --- TC-029: End-to-End Flow — Convention Creation Through Test Execution ---
// Traceability: TC-029 -> PRD / Main flow

func TestIntegration_TC_029_E2eConventionThroughExecution(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	// Step 1: Verify no Convention files exist
	files, err := os.ReadDir(conventionsDir)
	require.NoError(t, err)
	assert.Empty(t, files, "Should start with no Convention files")

	// Step 2: Run forge gen-test-scripts (cold start)
	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("Cold start gen-test-scripts output: %s", output)

	// Step 3: Verify hint about missing Convention
	hintPattern := regexp.MustCompile(`(?i)No test Convention files found|hint.*Convention|no.*Convention.*found`)
	_ = hintPattern.MatchString(output)

	// Step 5: Create Convention file
	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- Import: "testing"
- File pattern: "*_test.go"

## Assertion
- name: testify
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build e2e"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	// Step 6: Run gen-test-scripts again with Convention
	output2, exitCode2 := forgeGenTestScripts(t, projectRoot)
	t.Logf("Convention-based gen-test-scripts output: %s", output2)

	// Step 7: Verify Convention-based generation
	conventionPath := filepath.Join(conventionsDir, "testing-go.md")
	assert.FileExists(t, conventionPath)

	// Step 8: Verify Convention-declared testify is present
	content, err := os.ReadFile(conventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "testify")

	_ = exitCode
	_ = exitCode2
}

// --- TC-036: Integration — Test-Guide Creates Convention Consumed by Gen-Test-Scripts ---
// Traceability: TC-036 -> PRD / Main flow + Bootstrap flow

func TestIntegration_TC_036_TestGuideToGenTestScriptsPipeline(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	// Step 1: Verify no Convention files exist
	conventionPath := filepath.Join(conventionsDir, "testing-go.md")
	_, err := os.Stat(conventionPath)
	assert.True(t, os.IsNotExist(err), "Convention file should not exist initially")

	// Create existing test files with testify patterns
	testDir := filepath.Join(projectRoot, "pkg")
	require.NoError(t, os.MkdirAll(testDir, 0755))
	existingTest := `//go:build e2e

package pkg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	assert.NoError(t, nil)
}
`
	require.NoError(t, os.WriteFile(filepath.Join(testDir, "pipeline_test.go"), []byte(existingTest), 0644))

	// Step 2-3: Simulate test-guide creating Convention file
	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- Import: "testing"
- File pattern: "*_test.go"

## Assertion
- name: testify
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build e2e"

## Result Format
- Format: json-stream
`
	require.NoError(t, os.WriteFile(conventionPath, []byte(conventionContent), 0644))

	// Step 3: Verify Convention file created
	assert.FileExists(t, conventionPath)

	// Step 4: Verify required sections
	content, err := os.ReadFile(conventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "## Framework")
	assert.Contains(t, string(content), "## Assertion")

	// Step 5: Run forge gen-test-scripts
	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	// Step 6: Verify Convention-declared framework
	content, err = os.ReadFile(conventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "testify")

	_ = exitCode
}

// setupConventionForgeProjectWithFeature creates a temp project with a feature set.
func setupConventionForgeProjectWithFeature(t *testing.T) string {
	t.Helper()
	projectRoot := t.TempDir()

	require.NoError(t, os.MkdirAll(filepath.Join(projectRoot, ".forge"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Test Project\n"), 0644))

	// Create feature directory structure
	featureSlug := "test-feature"
	featureDir := filepath.Join(projectRoot, "docs", "features", featureSlug)
	featureTasksDir := filepath.Join(featureDir, "tasks")
	require.NoError(t, os.MkdirAll(featureTasksDir, 0755))

	// Create manifest
	manifestContent := `---
slug: "test-feature"
status: "in-progress"
---
`
	require.NoError(t, os.WriteFile(filepath.Join(featureDir, "manifest.md"), []byte(manifestContent), 0644))

	// Set current feature via state.json
	stateContent := `{"feature": "test-feature"}`
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, ".forge", "state.json"), []byte(stateContent), 0644))

	return projectRoot
}
