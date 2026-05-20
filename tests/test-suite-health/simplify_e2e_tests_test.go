//go:build e2e

package testsuitehealth

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// projectRootSimplify returns the forge project root directory.
func projectRootSimplify(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot determine test file location")
	}
	// thisFile: .../tests/test-suite-health/simplify_e2e_tests_test.go
	// up 2: test-suite-health -> tests -> project root
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
	return filepath.Join(projectRootSimplify(t), "tests", "e2e")
}

// TC-001: Verify tui-ui-design directory deleted
// Traceability: TC-001 -> Proposal Success Criterion 1
func TestTC_001_Simplify_VerifyTuiUiDesignDirectoryDeleted(t *testing.T) {
	tuiDir := filepath.Join(e2eRoot(t), "tui-ui-design")

	_, err := os.Stat(tuiDir)
	assert.True(t, os.IsNotExist(err),
		"directory tests/e2e/tui-ui-design/ should not exist, but it does")
}

// TC-002: Verify justfile-canonical-e2e directory has been consolidated into e2e-pipeline journey
// Traceability: TC-002 -> Proposal Success Criterion 2 (consolidated into e2e-pipeline journey)
func TestTC_002_Simplify_VerifyJustfileCanonicalE2eDirectoryConsolidated(t *testing.T) {
	oldDir := filepath.Join(e2eRoot(t), "justfile-canonical-e2e")

	_, err := os.Stat(oldDir)
	assert.True(t, os.IsNotExist(err),
		"directory tests/e2e/justfile-canonical-e2e/ should not exist (consolidated into tests/e2e-pipeline/)")
}
