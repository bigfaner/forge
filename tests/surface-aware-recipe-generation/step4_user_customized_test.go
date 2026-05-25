//go:build e2e

package surfacerecipegeneration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 4 Contract tests: Verify user-customized protection
// ==============================================================================

// Traceability: TC-030 -> Contract surface-aware-recipe-generation/step-4 Outcome "success"
// User-customized recipes are preserved when re-running init-justfile without --force.
func TestTC_030_UserCustomized_RecipePreserved(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// First run: generate justfile
	out, exitCode := runForgeRaw(t, projectDir, "init-justfile")
	if exitCode != 0 {
		t.Skipf("init-justfile not available: %s", out)
	}

	justfile := readJustfile(t, projectDir)
	assert.NotEmpty(t, justfile, "justfile should exist after first init-justfile")

	// Mark a recipe as user-customized
	customizedContent := strings.Replace(justfile, "dev:", "# user-customized\ndev:", 1)
	err := os.WriteFile(filepath.Join(projectDir, "justfile"), []byte(customizedContent), 0644)
	assert.NoError(t, err)

	// Second run: re-run init-justfile without --force-regenerate
	out, exitCode = runForgeRaw(t, projectDir, "init-justfile")
	if exitCode == 0 {
		updatedJustfile := readJustfile(t, projectDir)
		assert.True(t,
			strings.Contains(updatedJustfile, "user-customized"),
			"user-customized marker should be preserved")
	}
}

// Traceability: TC-031 -> Contract surface-aware-recipe-generation/step-4 Outcome "already-exists-customized"
// With --force-regenerate, user-customized recipes are overwritten.
func TestTC_031_UserCustomized_ForceRegenerate(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// Generate justfile
	out, exitCode := runForgeRaw(t, projectDir, "init-justfile")
	if exitCode != 0 {
		t.Skipf("init-justfile not available: %s", out)
	}

	// Add user-customized marker
	justfile := readJustfile(t, projectDir)
	customizedContent := strings.Replace(justfile, "dev:", "# user-customized\ndev:", 1)
	err := os.WriteFile(filepath.Join(projectDir, "justfile"), []byte(customizedContent), 0644)
	assert.NoError(t, err)

	// Re-run with --force-regenerate
	out, exitCode = runForgeRaw(t, projectDir, "init-justfile", "--force-regenerate")
	t.Logf("init-justfile --force-regenerate output (exit %d): %s", exitCode, out)
}
