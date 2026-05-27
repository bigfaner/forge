//go:build e2e

package surfacerecipegeneration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 1 Contract tests: Configure surfaces in config.yaml
// Tests verify config validation for surface entries.
// ==============================================================================

// Traceability: TC-021 -> Contract surface-aware-recipe-generation/step-1 Outcome "success"
// Config.yaml accepts valid surface-key to surface-type mappings.
func TestTC_021_ConfigureSurfaces_ValidEntry(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	out, exitCode := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "surfaces should succeed, got output:\n%s", out)
	assert.Contains(t, out, "admin-panel", "output should contain configured surface-key")
	assert.Contains(t, out, "web", "output should contain configured surface-type")
}

// Traceability: TC-022 -> Contract surface-aware-recipe-generation/step-1 Outcome "invalid-surface-type"
// Invalid surface type in config is filtered from --types output (unknown types ignored).
func TestTC_022_ConfigureSurfaces_InvalidSurfaceType(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  my-surface: desktop\n")

	// surfaces --types filters unknown types: empty output means no known types
	out, exitCode := runForgeRaw(t, projectDir, "surfaces", "--types")
	// --types should return 0 but with empty or no "desktop" in output
	assert.NotContains(t, out, "desktop", "unknown surface type should not appear in --types output")

	// Listing still shows the raw entry (surface is stored but type is unrecognized)
	out, exitCode = runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "surfaces listing should succeed")
	assert.Contains(t, out, "desktop", "raw listing should include the configured (but invalid) type")
}

// Traceability: TC-023 -> Contract surface-aware-recipe-generation/step-1 Outcome "config-validation-error"
// Invalid surface config format (no surfaces) returns graceful result.
func TestTC_023_ConfigureSurfaces_InvalidFormat(t *testing.T) {
	// Create project without surfaces field
	dir := createProjectWithoutSurfaces(t)

	// surfaces command with no configured surfaces should exit 0 with empty output
	out, exitCode := runForgeRaw(t, dir, "surfaces")
	assert.Equal(t, 0, exitCode, "surfaces on empty config should succeed (graceful)")
	// No surfaces configured means empty or no output
	t.Logf("surfaces output for empty config (exit %d): %s", exitCode, out)
}
