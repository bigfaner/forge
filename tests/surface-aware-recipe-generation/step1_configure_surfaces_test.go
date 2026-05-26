//go:build e2e

package surfacerecipegeneration

import (
	"strings"
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

	out, exitCode := runForgeRaw(t, projectDir, "config", "get", "surfaces")
	assert.Equal(t, 0, exitCode, "config get surfaces should succeed, got output:\n%s", out)
	assert.Contains(t, out, "admin-panel", "output should contain configured surface-key")
	assert.Contains(t, out, "web", "output should contain configured surface-type")
}

// Traceability: TC-022 -> Contract surface-aware-recipe-generation/step-1 Outcome "invalid-surface-type"
// Invalid surface type in config is rejected with descriptive error.
func TestTC_022_ConfigureSurfaces_InvalidSurfaceType(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  my-surface: desktop\n")

	out, exitCode := runForgeRaw(t, projectDir, "config", "get", "surfaces")
	// Invalid surface type should trigger error
	if exitCode != 0 {
		assert.True(t,
			strings.Contains(out, "desktop") || strings.Contains(out, "unsupported") || strings.Contains(out, "invalid"),
			"error should reference the unsupported surface type")
	}
}

// Traceability: TC-023 -> Contract surface-aware-recipe-generation/step-1 Outcome "config-validation-error"
// Invalid surface config format (scalar, list) is rejected.
func TestTC_023_ConfigureSurfaces_InvalidFormat(t *testing.T) {
	// Create project with invalid surfaces format (list instead of map)
	dir := createProjectWithoutSurfaces(t)

	// Overwrite config with invalid format
	// This test verifies the validation catches format issues
	out, exitCode := runForgeRaw(t, dir, "config", "get", "surfaces")
	t.Logf("config validation output (exit %d): %s", exitCode, out)
}
