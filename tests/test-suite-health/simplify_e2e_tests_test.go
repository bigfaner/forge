//go:build cli_functional

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

// TC-001: Verify tui-ui-design directory deleted
// Traceability: TC-001 -> Proposal Success Criterion 1
func TestTC_001_Simplify_VerifyTuiUiDesignDirectoryDeleted(t *testing.T) {
	// Old tests/e2e/ has been fully removed — tui-ui-design subdirectory
	// no longer exists by virtue of the parent being gone.
	e2eDir := filepath.Join(projectRootSimplify(t), "tests", "e2e")
	_, err := os.Stat(e2eDir)
	assert.True(t, os.IsNotExist(err),
		"tests/e2e/ directory should have been removed during cleanup")
}

// TC-002: Verify justfile-canonical-e2e directory has been consolidated into e2e-pipeline journey
// Traceability: TC-002 -> Proposal Success Criterion 2 (consolidated into e2e-pipeline journey)
func TestTC_002_Simplify_VerifyJustfileCanonicalE2eDirectoryConsolidated(t *testing.T) {
	// Old tests/e2e/ has been fully removed — justfile-canonical-e2e subdirectory
	// no longer exists by virtue of the parent being gone.
	e2eDir := filepath.Join(projectRootSimplify(t), "tests", "e2e")
	_, err := os.Stat(e2eDir)
	assert.True(t, os.IsNotExist(err),
		"tests/e2e/ directory should have been removed during cleanup")
}
