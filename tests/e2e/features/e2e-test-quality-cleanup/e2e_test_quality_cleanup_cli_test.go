//go:build e2e

package e2e

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// projectRoot returns the forge project root directory.
func projectRoot() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "..")
	abs, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}
	return abs
}

// fileExists checks whether a file exists at the given path.
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// fileContent reads a file and returns its content as a string.
func fileContent(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(data)
}

// fileSHA256 computes the SHA-256 hex digest of a file.
func fileSHA256(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	hash := sha256.Sum256(data)
	return fmt.Sprintf("%x", hash)
}

// =============================================================================
// TC-001: Deleted test files do not exist
// =============================================================================

// Traceability: TC-001 -> Proposal In Scope #1,#2,#3,#4
func TestTC_001_DeletedTestFilesDoNotExist(t *testing.T) {
	root := projectRoot()
	assert.NotEmpty(t, root, "project root must be resolved")

	deletedFiles := []string{
		"extract_design_md_platform_adapters_cli_test.go",
		"cli_list_reverse_chronological_cli_test.go",
		"fix_task_claim_priority_cli_test.go",
		"cli_lean_output_cli_test.go",
	}

	for _, filename := range deletedFiles {
		path := filepath.Join(root, "tests", "e2e", filename)
		assert.False(t, fileExists(path),
			"deleted file should not exist: %s", path)
	}
}

// =============================================================================
// TC-002: Deleted test functions do not exist
// =============================================================================

// Traceability: TC-002 -> Proposal In Scope #5,#6,#7
func TestTC_002_DeletedTestFunctionsDoNotExist(t *testing.T) {
	root := projectRoot()
	assert.NotEmpty(t, root, "project root must be resolved")

	type check struct {
		file     string
		patterns []string
	}

	checks := []check{
		{
			file: filepath.Join(root, "tests", "e2e", "simplify_e2e_tests_cli_test.go"),
			patterns: []string{
				`TestTC_003_`,
				`TestTC_004_`,
			},
		},
		{
			file: filepath.Join(root, "tests", "e2e", "feature_set_command_cli_test.go"),
			patterns: []string{
				`TestTC_016_`,
				`TestTC_017_`,
			},
		},
		{
			file: filepath.Join(root, "tests", "e2e", "quick_test_slim_cli_test.go"),
			patterns: []string{
				`TestTC_003_`,
				`TestTC_009_`,
				`TestTC_010_`,
				`TestTC_013_`,
				`TestTC_016_`,
			},
		},
	}

	for _, c := range checks {
		content := fileContent(t, c.file)
		for _, pattern := range c.patterns {
			assert.NotContains(t, content, pattern,
				"deleted function pattern %s should not exist in %s", pattern, c.file)
		}
	}
}

// =============================================================================
// TC-003: E2E test suite compiles successfully
// =============================================================================

// Traceability: TC-003 -> Proposal Success Criterion "just test-e2e compiles and all pass" + Key Risk "compilation failure"
func TestTC_003_E2ETestSuiteCompilesSuccessfully(t *testing.T) {
	root := projectRoot()
	assert.NotEmpty(t, root, "project root must be resolved")

	cmd := exec.Command("go", "build", "-tags=e2e", "./...")
	cmd.Dir = filepath.Join(root, "tests", "e2e")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err,
		"e2e test compilation should succeed, output:\n%s", string(out))
}

// =============================================================================
// TC-004: Zero unconditional t.Skip in test files
// =============================================================================

// Traceability: TC-004 -> Proposal Success Criterion "zero t.Skip unconditional skips"
func TestTC_004_ZeroUnconditionalTSkip(t *testing.T) {
	root := projectRoot()
	assert.NotEmpty(t, root, "project root must be resolved")

	e2eDir := filepath.Join(root, "tests", "e2e")

	// Find all _test.go files recursively (root + features/ subdirs)
	var testFiles []string
	err := filepath.WalkDir(e2eDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	assert.NoError(t, err, "failed to walk e2e test directory")

	// Pattern to find t.Skip( calls that are NOT inside an if/else conditional.
	// A t.Skip is unconditional if it appears at the top level of a function body
	// without being guarded by a condition.
	skipPattern := regexp.MustCompile(`(?m)^\t+t\.Skip\(`)

	for _, tf := range testFiles {
		content := fileContent(t, tf)
		matches := skipPattern.FindAllString(content, -1)
		assert.Empty(t, matches,
			"unconditional t.Skip() found in %s — all t.Skip calls must be inside conditional branches",
			tf)
	}
}

// =============================================================================
// TC-005: Zero recursive go test invocations
// =============================================================================

// Traceability: TC-005 -> Proposal Success Criterion "zero recursive exec.Command go test calls"
func TestTC_005_ZeroRecursiveGoTestInvocations(t *testing.T) {
	root := projectRoot()
	assert.NotEmpty(t, root, "project root must be resolved")

	e2eDir := filepath.Join(root, "tests", "e2e")

	cmd := exec.Command("grep",
		"-rn", `exec.Command("go", "test"`,
		"--include=*_test.go",
		e2eDir,
	)
	out, _ := cmd.CombinedOutput()
	assert.Empty(t, strings.TrimSpace(string(out)),
		"no exec.Command(\"go\", \"test\") calls should exist in e2e test files")
}

// =============================================================================
// TC-006: No static file text-grep tests remain
// =============================================================================

// Traceability: TC-006 -> Proposal Success Criterion "zero static source file text check tests"
func TestTC_006_NoStaticFileTextGrepTests(t *testing.T) {
	root := projectRoot()
	assert.NotEmpty(t, root, "project root must be resolved")

	e2eDir := filepath.Join(root, "tests", "e2e")

	var testFiles []string
	err := filepath.WalkDir(e2eDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	assert.NoError(t, err, "failed to walk e2e test directory")

	// Pattern: os.ReadFile on a .md or .go file followed by assert.Contains on that content.
	// Look for functions that read static source files and assert text content.
	// A static file text-grep test typically:
	//   1. Calls os.ReadFile on a path ending in .md or .go (not testdata/)
	//   2. Uses assert.Contains or strings.Contains on the result
	//
	// We detect this by finding test functions that have both an os.ReadFile call
	// on a non-testdata path and an assert.Contains call on the result.

	// Simple heuristic: check for os.ReadFile on paths containing ".md" or ".go"
	// that are NOT under testdata/
	readFilePattern := regexp.MustCompile(`os\.ReadFile\([^)]*(?:\.md|\.go)"?\)`)
	containsPattern := regexp.MustCompile(`assert\.Contains\(t,`)

	for _, tf := range testFiles {
		content := fileContent(t, tf)
		lines := strings.Split(content, "\n")

		for i, line := range lines {
			if readFilePattern.MatchString(line) && !strings.Contains(line, "testdata") {
				// Check if there's an assert.Contains within the same function scope
				// (roughly: within 30 lines after the ReadFile call)
				end := i + 30
				if end > len(lines) {
					end = len(lines)
				}
				for j := i; j < end; j++ {
					if containsPattern.MatchString(lines[j]) {
						t.Errorf(
							"suspected static file text-grep test in %s:%d — os.ReadFile on source file with assert.Contains",
							filepath.Base(tf), i+1,
						)
					}
				}
			}
		}
	}
}

// =============================================================================
// TC-007: No duplicate test files between root and features directory
// =============================================================================

// Traceability: TC-007 -> Proposal Success Criterion "no duplicate files between tests/e2e/ and tests/e2e/features/"
func TestTC_007_NoDuplicateTestFilesRootAndFeatures(t *testing.T) {
	root := projectRoot()
	assert.NotEmpty(t, root, "project root must be resolved")

	e2eDir := filepath.Join(root, "tests", "e2e")
	featuresDir := filepath.Join(e2eDir, "features")

	// Collect root-level test files
	rootEntries, err := os.ReadDir(e2eDir)
	assert.NoError(t, err, "failed to read e2e root directory")

	rootFiles := make(map[string]string) // filename -> sha256
	for _, entry := range rootEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(e2eDir, entry.Name())
		rootFiles[entry.Name()] = fileSHA256(t, path)
	}

	// Collect all test files under features/ subdirectories
	var featureFiles []struct {
		name     string
		sha256   string
		fullPath string
	}
	err = filepath.WalkDir(featuresDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			featureFiles = append(featureFiles, struct {
				name     string
				sha256   string
				fullPath string
			}{
				name:     d.Name(),
				sha256:   fileSHA256(t, path),
				fullPath: path,
			})
		}
		return nil
	})
	assert.NoError(t, err, "failed to walk features directory")

	// Compare: root files should not have identical content in features/
	for _, ff := range featureFiles {
		rootHash, exists := rootFiles[ff.name]
		if exists {
			assert.NotEqual(t, rootHash, ff.sha256,
				"duplicate file detected: %s exists in both root and %s with identical content",
				ff.name, ff.fullPath)
		}
	}
}
