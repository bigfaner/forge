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

// projectRoot returns the forge project root directory.
// This file lives at tests/e2e/.
// Project root is 2 levels up from this file.
func projectRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	// thisFile: .../tests/e2e/simplify_e2e_tests_cli_test.go
	// up 2: e2e -> tests -> project root
	dir := filepath.Join(filepath.Dir(thisFile), "..", "..")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("cannot resolve project root: %s", err)
	}
	return abs
}

// e2eRoot returns the tests/e2e/ directory.
func e2eRoot(t *testing.T) string {
	t.Helper()
	return filepath.Join(projectRoot(t), "tests", "e2e")
}

// TC-001: Verify tui-ui-design directory deleted
// Traceability: TC-001 -> Proposal Success Criterion 1
func TestTC_001_VerifyTuiUiDesignDirectoryDeleted(t *testing.T) {
	// Step 1: Check that the directory tests/e2e/tui-ui-design/ does not exist
	tuiDir := filepath.Join(e2eRoot(t), "tui-ui-design")

	_, err := os.Stat(tuiDir)
	assert.True(t, os.IsNotExist(err),
		"directory tests/e2e/tui-ui-design/ should not exist, but it does")
}

// TC-002: Verify TC-020 removed from justfile-canonical-e2e
// Traceability: TC-002 -> Proposal Success Criterion 2
func TestTC_002_VerifyTC020RemovedFromJustfileCanonicalE2e(t *testing.T) {
	// Step 1: Read the file tests/e2e/justfile-canonical-e2e/justfile_canonical_e2e_cli_test.go
	testFile := filepath.Join(e2eRoot(t), "justfile-canonical-e2e", "justfile_canonical_e2e_cli_test.go")

	content, err := os.ReadFile(testFile)
	assert.NoError(t, err, "failed to read justfile_canonical_e2e_cli_test.go")

	// Step 2: Search for the function name TestTC_020_AllManifestsContainZeroRunAndGraduateFields
	funcName := "TestTC_020_AllManifestsContainZeroRunAndGraduateFields"
	assert.False(t, strings.Contains(string(content), funcName),
		"function %s should not be present in the file, but it was found", funcName)
}

// TC-003: Verify e2e test suite compiles
// Traceability: TC-003 -> Proposal Success Criterion 3
func TestTC_003_VerifyE2eTestSuiteCompiles(t *testing.T) {
	// Step 1: Run go test -tags=e2e ./... -count=1 -run=^$ (compile-only)
	// Must run from tests/e2e/ (the Go module root), not the project root.
	cmd := exec.Command("go", "test", "-tags=e2e", "./...", "-count=1", "-run=^$")
	cmd.Dir = e2eRoot(t)
	out, err := cmd.CombinedOutput()

	// Step 2: Check exit code is 0 and no compilation errors
	assert.NoError(t, err, "e2e test suite compilation failed:\n%s", string(out))
}

// TC-004: Verify remaining CLI behavior tests pass
// Traceability: TC-004 -> Proposal Success Criterion 4
func TestTC_004_VerifyRemainingCliBehaviorTestsPass(t *testing.T) {
	// Step 1: Run go test -tags=e2e ./... -count=1 -timeout 120s
	// Must run from tests/e2e/ (the Go module root), not the project root.
	cmd := exec.Command("go", "test", "-tags=e2e", "./...", "-count=1", "-timeout", "120s")
	cmd.Dir = e2eRoot(t)
	out, err := cmd.CombinedOutput()

	// Step 2: Check exit code is 0 and all tests pass
	assert.NoError(t, err, "e2e test suite execution failed:\n%s", string(out))
}
