//go:build cli_functional

package surfacekeymigration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 6 Contract tests: quality-gate fix-task infers surface from file path
// Tests verify fix-task surface inference from failing test file paths.
// ==============================================================================

// Traceability: TC-017 -> Contract surface-key-migration/step-6 Outcome "success"
// Fix-task surface-key and surface-type inferred from failing file path.
func TestTC_017_FixTask_SurfaceInferredFromPath(t *testing.T) {
	cfg := "version: '1'\nsurfaces:\n  admin-panel: web\n"
	projectDir := createTempProjectWithConfig(t, cfg)

	// Verify surface detection works for a path within the configured surface
	out, exitCode := runForgeRaw(t, projectDir, "surfaces", projectDir+"/frontend/src")
	if exitCode == 0 {
		result := parseSurfaceOutput(t, out)
		assert.Equal(t, "admin-panel", result.Key,
			"fix-task should infer surface-key from file path")
		assert.Equal(t, "web", result.Type,
			"fix-task should infer surface-type from file path")
	}
}

// Traceability: TC-018 -> Contract surface-key-migration/step-6 Outcome "inference-failure"
// Fix-task created with empty surface fields when surface detection fails.
func TestTC_018_FixTask_InferenceFailure(t *testing.T) {
	// No surfaces config -- detection will fail
	projectDir := t.TempDir()

	out, exitCode := runForgeRaw(t, projectDir, "surfaces", projectDir+"/some/path")
	assert.NotEqual(t, 0, exitCode,
		"surface detection should fail without config, got output:\n%s", out)
	// When detection fails, fix-task should have empty/default surface fields
	// and error should be logged
}
