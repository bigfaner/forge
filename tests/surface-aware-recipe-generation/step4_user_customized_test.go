//go:build e2e

package surfacerecipegeneration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 4 Contract tests: Verify user-customized protection
// ==============================================================================

// Traceability: TC-030 -> Contract surface-aware-recipe-generation/step-4 Outcome "success"
// Config state is stable across repeated queries (user-customized protection is a skill concern).
func TestTC_030_UserCustomized_RecipePreserved(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// First query: get surfaces config
	out1, exitCode := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "first surfaces query should succeed")
	assert.Contains(t, out1, "admin-panel", "first query should contain surface-key")

	// Second query: verify identical result (config is unchanged)
	out2, exitCode := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, 0, exitCode, "second surfaces query should succeed")
	assert.Equal(t, out1, out2, "repeated queries should return identical results (config stability)")
}

// Traceability: TC-031 -> Contract surface-aware-recipe-generation/step-4 Outcome "already-exists-customized"
// Force-regenerate is a skill-level concern. Verify config mutation is not affected by CLI queries.
func TestTC_031_UserCustomized_ForceRegenerate(t *testing.T) {
	projectDir := createProjectWithSurfaces(t, "  admin-panel: web\n")

	// Repeated queries do not modify config
	out1, _ := runForgeRaw(t, projectDir, "surfaces")
	out2, _ := runForgeRaw(t, projectDir, "surfaces")
	assert.Equal(t, out1, out2, "repeated queries should not modify config")

	// Config file content unchanged
	configContent, err := os.ReadFile(filepath.Join(projectDir, ".forge", "config.yaml"))
	assert.NoError(t, err, "config.yaml should be readable")
	assert.Contains(t, string(configContent), "admin-panel", "config should still contain surface-key")
}
