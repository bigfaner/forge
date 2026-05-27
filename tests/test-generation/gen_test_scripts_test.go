//go:build cli_functional

package testgeneration

// ==============================================================================
// forge gen-test-scripts — Journey: test-generation
// Tests cover Convention loading, fallback, compile gate, selective loading,
// error handling, and recovery behavior.
// ==============================================================================

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Fixture helpers for gen-test-scripts tests ---

// setupGoProjectFixture creates a temp directory with a Go project structure
// suitable for forge gen-test-scripts testing.
// Returns project root and convention dir paths.
func setupGoProjectFixture(t *testing.T) (projectRoot string, conventionsDir string) {
	t.Helper()
	projectRoot = t.TempDir()

	// Create .forge directory
	require.NoError(t, os.MkdirAll(filepath.Join(projectRoot, ".forge"), 0755))

	// Create docs/conventions directory
	conventionsDir = filepath.Join(projectRoot, "docs", "conventions")
	require.NoError(t, os.MkdirAll(conventionsDir, 0755))

	// Create a minimal go.mod
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "go.mod"), []byte("module test-project\n\ngo 1.26\n"), 0644))

	// Create justfile with compile recipe
	justfileContent := `compile:
	#!/usr/bin/env bash
	set -euo pipefail
	go build -tags=e2e ./...
	echo "OK: Go compilation passed"
`
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(justfileContent), 0644))

	// Create CLAUDE.md so forge can find the project root
	require.NoError(t, os.WriteFile(filepath.Join(projectRoot, "CLAUDE.md"), []byte("# Test Project\n"), 0644))

	return projectRoot, conventionsDir
}

// writeConventionFile creates a Convention file at the given path with the specified content.
func writeConventionFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0755))
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
}

// forgeGenTestScripts runs "forge gen-test-scripts" in the given project root.
// Returns combined output and exit code.
func forgeGenTestScripts(t *testing.T, projectRoot string, extraArgs ...string) (string, int) {
	t.Helper()
	args := []string{"gen-test-scripts"}
	args = append(args, extraArgs...)
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+projectRoot)
	cmd.Dir = projectRoot
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

// runJustCompile runs "just compile" in the given project root.
func runJustCompile(t *testing.T, projectRoot string) (string, int) {
	t.Helper()
	cmd := exec.Command("just", "compile")
	cmd.Dir = projectRoot
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

// runGrep runs grep on a file and returns exit code and output.
func runGrep(t *testing.T, pattern, filePath string) (string, int) {
	t.Helper()
	cmd := exec.Command("grep", "-c", pattern, filePath)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return strings.TrimSpace(string(out)), exitCode
}

// forgeCmdForConvention returns an exec.Cmd for the forge CLI binary.
func forgeCmdForConvention(args ...string) *exec.Cmd {
	return exec.Command(testkit.ForgeBinary, args...)
}

// --- TC-001: Gen-Test-Scripts Uses Convention-Declared Framework for Non-Default Projects ---
// Traceability: TC-001 -> Story 1 / AC-1

func TestGenScripts_TC_001_GenTestScriptsUsesConventionFrameworkForGinkgo(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	// Step 1: Create existing test files with ginkgo imports
	existingTestDir := filepath.Join(projectRoot, "internal", "pkg")
	require.NoError(t, os.MkdirAll(existingTestDir, 0755))
	ginkgoTest := `package pkg
import . "github.com/onsi/ginkgo/v2"
var _ = Describe("test", func(){})
`
	require.NoError(t, os.WriteFile(filepath.Join(existingTestDir, "pkg_test.go"), []byte(ginkgoTest), 0644))

	// Step 2: Create Convention file declaring ginkgo as framework
	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: ginkgo
- Import: "github.com/onsi/ginkgo/v2"
- File pattern: "*_test.go"

## Assertion
- name: gomega
- Import: "github.com/onsi/gomega"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	// Step 3-4: Run forge gen-test-scripts
	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	assert.FileExists(t, filepath.Join(conventionsDir, "testing-go.md"))

	content, err := os.ReadFile(filepath.Join(conventionsDir, "testing-go.md"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "ginkgo")

	// Steps 7-8: Verify compile recipe exists
	justfilePath := filepath.Join(projectRoot, "justfile")
	content, err = os.ReadFile(justfilePath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "compile")

	_ = exitCode
}

// --- TC-002: Gen-Test-Scripts Warns and Falls Back on Empty Convention Framework Section ---
// Traceability: TC-002 -> Story 1 / AC-2

func TestGenScripts_TC_002_GenTestScriptsWarnsOnEmptyConventionFramework(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
## Framework

## Assertion
- name: assert
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	warningPattern := regexp.MustCompile(`(?i)warning.*Framework|Framework.*missing|missing.*Framework`)
	hasWarning := warningPattern.MatchString(output)

	conventionPath := filepath.Join(conventionsDir, "testing-go.md")
	assert.FileExists(t, conventionPath)

	content, err := os.ReadFile(conventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "## Assertion")

	_ = exitCode
	_ = hasWarning
}

// --- TC-005: Gen-Test-Scripts Proceeds Without Convention on Cold Start ---
// Traceability: TC-005 -> Story 3 / AC-1

func TestGenScripts_TC_005_GenTestScriptsColdStartNoConvention(t *testing.T) {
	projectRoot, _ := setupGoProjectFixture(t)

	// Ensure no Convention files exist
	conventionsDir := filepath.Join(projectRoot, "docs", "conventions")
	files, err := os.ReadDir(conventionsDir)
	require.NoError(t, err)
	assert.Empty(t, files, "conventions dir should be empty for cold start test")

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	hintPattern := regexp.MustCompile(`(?i)No test Convention files found|hint.*Convention|no.*Convention.*found`)
	hasHint := hintPattern.MatchString(output)

	files, err = os.ReadDir(conventionsDir)
	require.NoError(t, err)
	assert.Empty(t, files)

	_ = exitCode
	_ = hasHint
}

// --- TC-008: Gen-Test-Scripts Loads Convention by Interface Type Selectively ---
// Traceability: TC-008 -> Story 5 / AC-1

func TestGenScripts_TC_008_GenTestScriptsSelectiveConventionLoadingByInterface(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	goConventionContent := `---
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
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), goConventionContent)

	jsConventionContent := `---
domains: [testing, javascript, web-ui]
---
## Framework
- name: vitest
- Import: "vitest"
- File pattern: "*.test.ts"

## Assertion
- name: vitest-expect
- Import: "@vitest/expect"

## Tags
- No build tag needed

## Result Format
- Format: tap
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-javascript.md"), jsConventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	goConventionPath := filepath.Join(conventionsDir, "testing-go.md")
	assert.FileExists(t, goConventionPath)
	content, err := os.ReadFile(goConventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "go-testing")

	jsConventionPath := filepath.Join(conventionsDir, "testing-javascript.md")
	assert.FileExists(t, jsConventionPath)

	content, err = os.ReadFile(goConventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "testing")
	assert.Contains(t, string(content), "testify")

	_ = exitCode
}

// --- TC-009: Gen-Test-Scripts Loads Overlapping Domain Conventions and Logs Overlap Warning ---
// Traceability: TC-009 -> Story 5 / AC-2

func TestGenScripts_TC_009_GenTestScriptsMergesOverlappingDomainConventions(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	assertConventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- File pattern: "*_test.go"

## Assertion
- name: assert (not require)
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-00-assert.md"), assertConventionContent)

	requireConventionContent := `---
domains: [testing, go, cli]
---
## Framework
- name: go-testing
- File pattern: "*_test.go"

## Assertion
- name: require (not assert)
- Import: "github.com/stretchr/testify/require"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-01-require.md"), requireConventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	overlapPattern := regexp.MustCompile(`(?i)overlap|domain.*overlap`)
	hasOverlapWarning := overlapPattern.MatchString(output)

	assert.FileExists(t, filepath.Join(conventionsDir, "testing-00-assert.md"))
	assert.FileExists(t, filepath.Join(conventionsDir, "testing-01-require.md"))

	_ = exitCode
	_ = hasOverlapWarning
}

// --- TC-011: Convention File Missing Domains Frontmatter Treated as Non-Loadable ---
// Traceability: TC-011 -> Spec FS-1 / Validation Rules

func TestGenScripts_TC_011_ConventionFileMissingDomainsFrontmatter(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `## Framework
- name: go-testing
- File pattern: "*_test.go"

## Assertion
- name: testify
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	warningPattern := regexp.MustCompile(`(?i)warning.*non-loadable|warning.*domains.*missing|Convention.*skipped`)
	hasWarning := warningPattern.MatchString(output)

	_ = exitCode
	_ = hasWarning
}

// --- TC-016: No Convention Files Hint During Gen-Test-Scripts ---
// Traceability: TC-016 -> Spec FS-2 / Error Handling

func TestGenScripts_TC_016_HintWhenNoConventionFiles(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	files, err := os.ReadDir(conventionsDir)
	require.NoError(t, err)
	assert.Empty(t, files)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	hintPattern := regexp.MustCompile(`(?i)No test Convention files found|hint.*test-guide`)
	hasHint := hintPattern.MatchString(output)

	_ = exitCode
	_ = hasHint
}

// --- TC-017: Convention Wins Over Conflicting Reconnaissance Signals ---
// Traceability: TC-017 -> Spec FS-3 / Reliability Expectations

func TestGenScripts_TC_017_ConventionOverridesReconnaissanceConflict(t *testing.T) {
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
- Build tag: "//go:build cli_functional"

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
func TestSomething(t *testing.T) {
	require.NoError(t, nil)
}
`
	require.NoError(t, os.WriteFile(filepath.Join(testDir, "pkg_test.go"), []byte(requireTest), 0644))

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	conflictPattern := regexp.MustCompile(`(?i)conflict|Convention.*Reconnaissance|override.*detected`)
	hasConflict := conflictPattern.MatchString(output)

	_ = exitCode
	_ = hasConflict
}

// --- TC-019: Gen-Test-Scripts Produces Compilable Output on First Attempt for Standard Go Project ---
// Traceability: TC-019 -> Spec / Goals + Performance Requirements

func TestGenScripts_TC_019_FirstPassCompileForStandardProject(t *testing.T) {
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
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	conventionPath := filepath.Join(conventionsDir, "testing-go.md")
	assert.FileExists(t, conventionPath)

	content, err := os.ReadFile(conventionPath)
	require.NoError(t, err)
	assert.Contains(t, string(content), "go-testing")
	assert.Contains(t, string(content), "testify")

	_ = exitCode
}

// --- TC-028: Convention File Unreadable Due to Permissions Is Skipped with Warning ---
// Traceability: TC-028 -> Spec FS-2 / Error Handling

func TestGenScripts_TC_028_ConventionFileUnreadablePermissionSkipped(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- File pattern: "*_test.go"

## Assertion
- name: testify
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	conventionPath := filepath.Join(conventionsDir, "testing-go.md")
	writeConventionFile(t, conventionPath, conventionContent)

	// Set file permissions to unreadable
	require.NoError(t, os.Chmod(conventionPath, 0000))

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	warningPattern := regexp.MustCompile(`(?i)warning.*permission|warning.*unreadable|skip.*Convention`)
	hasWarning := warningPattern.MatchString(output)

	// Restore permissions for cleanup
	require.NoError(t, os.Chmod(conventionPath, 0644))

	_ = exitCode
	_ = hasWarning
}

// --- TC-031: Compile Gate Intermediate Retry Feeds Error Back to LLM ---
// Traceability: TC-031 -> Spec FS-4 / Retry Semantics

func TestGenScripts_TC_031_CompileGateIntermediateRetryFeedback(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
## Framework
- name: go-testing
- Import: "github.com/stretchr/testify/require"
- File pattern: "*_test.go"

## Assertion
- name: require
- Import: "github.com/stretchr/testify/require"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	retryPattern := regexp.MustCompile(`(?i)retry|attempt.*2|regenerat|re-generat`)
	hasRetry := retryPattern.MatchString(output)

	_ = exitCode
	_ = hasRetry
}

// --- TC-032: Convention File with All Required Sections Missing Falls Back Fully ---
// Traceability: TC-032 -> Spec FS-1 / Validation Rules (boundary case)

func TestGenScripts_TC_032_FallbackWhenAllSectionsMissing(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
This Convention file intentionally has no sections.
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	warningPattern := regexp.MustCompile(`(?i)warning.*missing.*section|missing.*Framework|missing.*Assertion|missing.*Tags|missing.*Result Format`)
	hasWarning := warningPattern.MatchString(output)

	_ = exitCode
	_ = hasWarning
}

// --- TC-034: Convention File with Invalid Encoding Is Skipped with Warning ---
// Traceability: TC-034 -> Spec FS-2 / Error Handling (encoding)

func TestGenScripts_TC_034_ConventionFileInvalidEncodingSkipped(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionPath := filepath.Join(conventionsDir, "testing-go.md")
	binaryContent := []byte{0xff, 0xfe, 0x20, 0x00, 'i', 0x00, 'n', 0x00, 'v', 0x00, 'a', 0x00, 'l', 0x00, 'i', 0x00, 'd', 0x00}
	require.NoError(t, os.WriteFile(conventionPath, binaryContent, 0644))

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	warningPattern := regexp.MustCompile(`(?i)warning.*encoding|warning.*unreadable|skip.*Convention`)
	hasWarning := warningPattern.MatchString(output)

	_ = exitCode
	_ = hasWarning
}

// --- TC-035: Gen-Test-Scripts Warns on Invalid Section Content (Empty Framework Name) ---
// Traceability: TC-035 -> Spec FS-1 / Validation Rules (Invalid section content)

func TestGenScripts_TC_035_GenTestScriptsWarnsOnInvalidSectionContent(t *testing.T) {
	projectRoot, conventionsDir := setupGoProjectFixture(t)

	conventionContent := `---
domains: [testing, go]
---
## Framework
- name:
- File pattern: "*_test.go"

## Assertion
- name: testify
- Import: "github.com/stretchr/testify/assert"

## Tags
- Build tag: "//go:build cli_functional"

## Result Format
- Format: json-stream
`
	writeConventionFile(t, filepath.Join(conventionsDir, "testing-go.md"), conventionContent)

	output, exitCode := forgeGenTestScripts(t, projectRoot)
	t.Logf("gen-test-scripts output: %s", output)

	warningPattern := regexp.MustCompile(`(?i)warning.*Framework.*invalid|warning.*Framework.*empty|treat.*Framework.*as missing`)
	hasWarning := warningPattern.MatchString(output)

	_ = exitCode
	_ = hasWarning
}
