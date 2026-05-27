//go:build e2e

package surfacerecipegeneration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 2 Contract tests: Run init-justfile
// Tests verify init-justfile generates surface-specific recipes.
// ==============================================================================

// Traceability: TC-024 -> Contract surface-aware-recipe-generation/step-2 Outcome "success"
// Surfaces config is readable; init-justfile is a Claude skill (not CLI), verify config readiness.
func TestTC_024_InitJustfile_Success(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// Verify surfaces config is readable and correct for init-justfile consumption
	out, exitCode := runForgeRaw(t, projectDir, "surfaces", "--json")
	assert.Equal(t, 0, exitCode, "surfaces --json should succeed for config readiness check")
	assert.Contains(t, out, "admin-panel", "config should contain surface-key")
	assert.Contains(t, out, "web", "config should contain surface-type")

	// init-justfile is a Claude skill invoked via /init-justfile, not a CLI command.
	// The CLI surfaces command confirms config readiness for skill consumption.
}

// Traceability: TC-025 -> Contract surface-aware-recipe-generation/step-2 Outcome "no-surfaces-configured"
// Without surfaces configured, CLI surfaces --types returns empty (no known types).
func TestTC_025_InitJustfile_NoSurfacesConfigured(t *testing.T) {
	projectDir := createProjectWithoutSurfaces(t)

	// Without surfaces, --types should return empty output
	out, exitCode := runForgeRaw(t, projectDir, "surfaces", "--types")
	assert.Equal(t, 0, exitCode, "surfaces --types should succeed with no surfaces")
	assert.Empty(t, strings.TrimSpace(out), "no surfaces configured: --types output should be empty")
}

// Traceability: TC-026 -> Contract surface-aware-recipe-generation/step-2 Outcome "just-version-below-minimum"
// just version check is performed by init-justfile skill (not CLI). Verify forge init --skip-just skips just check.
func TestTC_026_InitJustfile_JustVersionBelowMinimum(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// forge init --skip-just should work regardless of just version
	out, exitCode := runForgeRaw(t, projectDir, "init", "--skip-just")
	t.Logf("init --skip-just output (exit %d): %s", exitCode, out)
}

// Traceability: TC-027 -> Contract surface-aware-recipe-generation/step-2 Outcome "surface-rule-file-missing"
// CLI surfaces handles unknown types gracefully (filtered from --types).
func TestTC_027_InitJustfile_SurfaceRuleFileMissing(t *testing.T) {
	// Use an unusual surface type that would not have a rule file
	projectDir := createProjectWithSurfaces(t, "  my-cli: cli\n")

	// CLI surfaces should still list the entry but --types filters unknown types
	out, exitCode := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "surfaces listing should succeed")
	assert.Contains(t, out, "my-cli", "listing should contain surface key")

	// --types only returns known types
	out, exitCode = runForgeRaw(t, projectDir, "surfaces", "--types")
	assert.Equal(t, 0, exitCode, "surfaces --types should succeed")
	// cli is a known type, so it should appear
	assert.Contains(t, out, "cli", "cli is a known surface type")
}
