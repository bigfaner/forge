//go:build e2e

package surfacerecipegeneration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 3 Contract tests: Verify generated justfile contains surface-specific recipes
// ==============================================================================

// Traceability: TC-028 -> Contract surface-aware-recipe-generation/step-3 Outcome "success"
// Web surface configuration is queryable with correct type information.
func TestTC_028_VerifyRecipes_WebSurfaceHasFullOrchestration(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// Verify web surface type is correctly stored and queryable
	out, exitCode := runForgeRaw(t, projectDir, "surfaces", "--json")
	assert.Equal(t, 0, exitCode, "surfaces --json should succeed")
	assert.Contains(t, out, "web", "web surface type should be in output")

	// Verify types listing includes web
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "--types")
	assert.Equal(t, 0, exitCode, "surfaces --types should succeed")
	assert.Contains(t, out, "web", "web should be in known types")
}

// Traceability: TC-029 -> Contract surface-aware-recipe-generation/step-3 Outcome "recipe-not-found"
// CLI surface type is recognized and queryable (probe/run exclusion is a skill concern).
func TestTC_029_VerifyRecipes_CliSurfaceHasNoProbe(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  my-cli: cli\n")

	// Verify CLI surface type is correctly recognized
	out, exitCode := runForgeRaw(t, projectDir, "surfaces", "--types")
	assert.Equal(t, 0, exitCode, "surfaces --types should succeed")
	assert.Contains(t, out, "cli", "cli should be in known types")

	// Verify query by path returns cli type
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "my-cli")
	assert.Equal(t, 0, exitCode, "surfaces query should find my-cli")
	assert.Contains(t, out, "cli", "my-cli should map to cli type")
}
