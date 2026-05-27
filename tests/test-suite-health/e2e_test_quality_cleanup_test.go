//go:build cli_functional

package testsuitehealth

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// projectRootQuality returns the forge project root directory.
func projectRootQuality() string {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	// thisFile: .../tests/test-suite-health/e2e_test_quality_cleanup_test.go
	// up 2: test-suite-health -> tests -> project root
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..")
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

// =============================================================================
// TC-001: Deleted test files do not exist
// =============================================================================

// Traceability: TC-001 -> Proposal In Scope #1,#2,#3,#4
func TestTC_001_DeletedTestFilesDoNotExist(t *testing.T) {
	// Old tests/e2e/ directory has been fully removed (cleanup task).
	// All previously tracked deleted files are gone along with the directory.
	root := projectRootQuality()
	assert.NotEmpty(t, root, "project root must be resolved")

	e2eDir := filepath.Join(root, "tests", "e2e")
	_, err := os.Stat(e2eDir)
	assert.True(t, os.IsNotExist(err),
		"tests/e2e/ directory should have been removed during cleanup")
}

// =============================================================================
// TC-002: Deleted test functions do not exist
// =============================================================================

// Traceability: TC-002 -> Proposal In Scope #5,#6,#7
func TestTC_002_DeletedTestFunctionsDoNotExist(t *testing.T) {
	// Old tests/e2e/ directory has been fully removed (cleanup task).
	// All previously tracked deleted functions are gone along with the directory.
	root := projectRootQuality()
	assert.NotEmpty(t, root, "project root must be resolved")

	e2eDir := filepath.Join(root, "tests", "e2e")
	_, err := os.Stat(e2eDir)
	assert.True(t, os.IsNotExist(err),
		"tests/e2e/ directory should have been removed during cleanup")
}

// =============================================================================
// TC-003: E2E test suite compiles successfully
// =============================================================================

// Traceability: TC-003 -> Proposal Success Criterion "just test-e2e compiles and all pass" + Key Risk "compilation failure"
func TestTC_003_E2ETestSuiteCompilesSuccessfully(t *testing.T) {
	root := projectRootQuality()
	assert.NotEmpty(t, root, "project root must be resolved")

	// Tests migrated from tests/e2e/ to Journey directories under tests/.
	cmd := exec.Command("go", "build", "-tags=e2e", "./...")
	cmd.Dir = filepath.Join(root, "tests")
	out, err := cmd.CombinedOutput()
	assert.NoError(t, err,
		"e2e test compilation should succeed, output:\n%s", string(out))
}

// =============================================================================
// TC-004: Zero unconditional t.Skip in test files
// =============================================================================

// Traceability: TC-004 -> Proposal Success Criterion "zero t.Skip unconditional skips"
func TestTC_004_ZeroUnconditionalTSkip(t *testing.T) {
	root := projectRootQuality()
	assert.NotEmpty(t, root, "project root must be resolved")

	testsDir := filepath.Join(root, "tests")

	// Find all _test.go files recursively in Journey directories
	var testFiles []string
	err := filepath.WalkDir(testsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	assert.NoError(t, err, "failed to walk tests directory")

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
	root := projectRootQuality()
	assert.NotEmpty(t, root, "project root must be resolved")

	testsDir := filepath.Join(root, "tests")
	// Construct target dynamically to avoid matching this file's own source.
	target := `exec.Command("` + "go" + `", "` + "test" + `"`

	var testFiles []string
	err := filepath.WalkDir(testsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip testkit/ directory — it's infrastructure, not test files
		if d.IsDir() && d.Name() == "testkit" {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	assert.NoError(t, err, "failed to walk tests directory")

	for _, tf := range testFiles {
		content := fileContent(t, tf)
		assert.NotContains(t, content, target,
			"no exec.Command(\"go\", \"test\") calls should exist in %s", tf)
	}
}

// =============================================================================
// TC-006: No static file text-grep tests remain
// =============================================================================

// Traceability: TC-006 -> Proposal Success Criterion "zero static source file text check tests"
func TestTC_006_NoStaticFileTextGrepTests(t *testing.T) {
	root := projectRootQuality()
	assert.NotEmpty(t, root, "project root must be resolved")

	testsDir := filepath.Join(root, "tests")

	var testFiles []string
	err := filepath.WalkDir(testsDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip testkit/ directory — it's infrastructure
		// Skip test-generation/ — its tests validate generated content which is expected behavior
		if d.IsDir() && (d.Name() == "testkit" || d.Name() == "test-generation") {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), "_test.go") {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	assert.NoError(t, err, "failed to walk tests directory")

	// Pattern: os.ReadFile on a .md or .go file followed by assert.Contains on that content.
	readFilePattern := regexp.MustCompile(`os\.ReadFile\([^)]*(?:\.md|\.go)"?\)`)
	containsPattern := regexp.MustCompile(`assert\.Contains\(t,`)

	for _, tf := range testFiles {
		content := fileContent(t, tf)
		lines := strings.Split(content, "\n")

		for i, line := range lines {
			if readFilePattern.MatchString(line) && !strings.Contains(line, "testdata") {
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
	root := projectRootQuality()
	assert.NotEmpty(t, root, "project root must be resolved")

	e2eDir := filepath.Join(root, "tests", "e2e")

	// Old tests/e2e/ directory has been fully removed (cleanup task).
	// Verify it no longer exists — no duplicates possible.
	_, err := os.Stat(e2eDir)
	assert.True(t, os.IsNotExist(err),
		"tests/e2e/ directory should have been removed during cleanup")
}
