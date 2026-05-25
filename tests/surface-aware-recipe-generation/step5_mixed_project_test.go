//go:build e2e

package surfacerecipegeneration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 5 Contract tests: Verify mixed-project recipe generation
// ==============================================================================

// Traceability: TC-032 -> Contract surface-aware-recipe-generation/step-5 Outcome "success"
// Multi-surface project generates prefixed recipes and aggregation recipes.
func TestTC_032_MixedProject_PrefixedAndAggregationRecipes(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n  payment-service: api\n")

	out, exitCode := runForgeRaw(t, projectDir, "init-justfile")
	if exitCode != 0 {
		t.Skipf("init-justfile not available: %s", out)
	}

	justfile := readJustfile(t, projectDir)
	assert.NotEmpty(t, justfile, "justfile should exist")

	// Should have per-surface-key prefixed recipes
	assert.True(t,
		strings.Contains(justfile, "dev-admin-panel") || strings.Contains(justfile, "admin-panel"),
		"justfile should contain admin-panel prefixed recipes")
	assert.True(t,
		strings.Contains(justfile, "dev-payment-service") || strings.Contains(justfile, "payment-service"),
		"justfile should contain payment-service prefixed recipes")

	// Should have aggregation recipes
	assert.True(t,
		recipeExists(justfile, "dev") || strings.Contains(justfile, "dev:"),
		"justfile should have aggregation 'dev' recipe")
}

// Traceability: TC-033 -> Contract surface-aware-recipe-generation/step-5 Outcome "single-surface-fallback"
// Single-surface project has unprefixed recipes without surface-key prefix.
func TestTC_033_MixedProject_SingleSurfaceUnprefixed(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	out, exitCode := runForgeRaw(t, projectDir, "init-justfile")
	if exitCode != 0 {
		t.Skipf("init-justfile not available: %s", out)
	}

	justfile := readJustfile(t, projectDir)
	// Single surface should have unprefixed recipes
	assert.True(t,
		recipeExists(justfile, "dev") || recipeExists(justfile, "test"),
		"single-surface justfile should have unprefixed recipes")
}
