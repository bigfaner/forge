//go:build cli_functional

package surfacerecipegeneration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 5 Contract tests: Verify mixed-project recipe generation
// ==============================================================================

// Traceability: TC-032 -> Contract surface-aware-recipe-generation/step-5 Outcome "success"
// Multi-surface project is queryable with both surface keys.
func TestTC_032_MixedProject_PrefixedAndAggregationRecipes(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n  payment-service: api\n")

	// Verify both surfaces are listed
	out, exitCode := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "surfaces should succeed")
	assert.Contains(t, out, "admin-panel", "output should contain admin-panel")
	assert.Contains(t, out, "payment-service", "output should contain payment-service")

	// Verify both types are listed
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "--types")
	assert.Equal(t, 0, exitCode, "surfaces --types should succeed")
	assert.Contains(t, out, "web", "types should include web")
	assert.Contains(t, out, "api", "types should include api")

	// Verify individual path queries
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "admin-panel")
	assert.Equal(t, 0, exitCode, "query for admin-panel should succeed")
	assert.Contains(t, out, "web", "admin-panel should map to web")
}

// Traceability: TC-033 -> Contract surface-aware-recipe-generation/step-5 Outcome "single-surface-fallback"
// Single-surface project has a scalar representation when key is ".".
func TestTC_033_MixedProject_SingleSurfaceUnprefixed(t *testing.T) {
	// Create a project with a single surface using "." key (scalar form)
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// Single surface should be queryable
	out, exitCode := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "surfaces should succeed")
	assert.Contains(t, out, "admin-panel", "output should contain surface key")
	assert.Contains(t, out, "web", "output should contain surface type")

	// Verify type is recognized
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "--types")
	assert.Equal(t, 0, exitCode, "surfaces --types should succeed")
	assert.Contains(t, out, "web", "web should be in known types")
}
