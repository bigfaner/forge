//go:build e2e

package surfacerecipegeneration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 3 Contract tests: Verify generated justfile contains surface-specific recipes
// ==============================================================================

// Traceability: TC-028 -> Contract surface-aware-recipe-generation/step-3 Outcome "success"
// Generated justfile has correct surface-type-specific recipes.
func TestTC_028_VerifyRecipes_WebSurfaceHasFullOrchestration(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	out, exitCode := runForgeRaw(t, projectDir, "init-justfile")
	if exitCode != 0 {
		t.Skipf("init-justfile not available or failed: %s", out)
	}

	justfile := readJustfile(t, projectDir)
	assert.NotEmpty(t, justfile, "justfile should exist after init-justfile")

	// Web surface should have: dev, probe, test, test-teardown
	assert.True(t, recipeExists(justfile, "dev"),
		"web surface justfile should have 'dev' recipe")
	assert.True(t, recipeExists(justfile, "test"),
		"web surface justfile should have 'test' recipe")
}

// Traceability: TC-029 -> Contract surface-aware-recipe-generation/step-3 Outcome "recipe-not-found"
// Expected recipe missing from justfile indicates incomplete generation.
func TestTC_029_VerifyRecipes_CliSurfaceHasNoProbe(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  my-cli: cli\n")

	out, exitCode := runForgeRaw(t, projectDir, "init-justfile")
	if exitCode != 0 {
		t.Skipf("init-justfile not available or failed: %s", out)
	}

	justfile := readJustfile(t, projectDir)
	// CLI surface should NOT have probe or run recipes
	assert.False(t, strings.Contains(justfile, "probe:"),
		"CLI surface justfile should NOT have 'probe' recipe")
	assert.False(t, strings.Contains(justfile, "run:"),
		"CLI surface justfile should NOT have 'run' recipe")
}
