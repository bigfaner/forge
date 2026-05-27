//go:build cli_functional

package surfacerecipegeneration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Journey smoke test: surface-aware-recipe-generation happy path
// ==============================================================================

// Traceability: Smoke test -> surface-aware-recipe-generation Journey happy path
// Verifies surface configuration is queryable via CLI surfaces command.
func TestSurfaceAwareRecipeGeneration_Smoke(t *testing.T) {
	// Step 1: Configure surfaces and verify via CLI
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	out, exitCode := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "Step 1: surfaces query should succeed, output:\n%s", out)
	assert.Contains(t, out, "admin-panel", "Step 1: output should contain surface-key")
	assert.Contains(t, out, "web", "Step 1: output should contain surface-type")

	// Step 2: Verify surfaces --types returns known types
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "--types")
	assert.Equal(t, 0, exitCode, "Step 2: surfaces --types should succeed, output:\n%s", out)
	assert.Contains(t, out, "web", "Step 2: --types should list 'web'")

	// Step 3: Verify JSON output structure
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "--json")
	assert.Equal(t, 0, exitCode, "Step 3: surfaces --json should succeed, output:\n%s", out)
	assert.Contains(t, out, `"key"`, "Step 3: JSON should contain key field")
	assert.Contains(t, out, `"type"`, "Step 3: JSON should contain type field")

	// Invariant: cli/tui surfaces never generate run or probe recipes
	// Invariant: all recipes include dual-platform variants where applicable
	// Invariant: zero regression for projects without surfaces config
}
