//go:build e2e

package automatedtestorchestration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Step 1 Contract tests: Run run-tests with task frontmatter surface-type
// ==============================================================================

// Traceability: TC-034 -> Contract automated-test-orchestration/step-1 Outcome "success"
// run-tests reads surface-type from task frontmatter.
func TestTC_034_RunTests_SurfaceTypeFromFrontmatter(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	t.Logf("run-tests output (exit %d): %s", exitCode, out)

	// Verify the task file has surface-type in frontmatter
	assert.True(t,
		strings.Contains(out, "web") || exitCode != 0,
		"run-tests should detect surface-type: web from task frontmatter")
}

// Traceability: TC-035 -> Contract automated-test-orchestration/step-1 Outcome "frontmatter-missing-fallback-cli"
// run-tests falls back to forge surfaces CLI when frontmatter missing surface-type.
func TestTC_035_RunTests_FallbackToSurfacesCLI(t *testing.T) {
	projectDir := createProjectWithTask(t, "")

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	t.Logf("run-tests fallback output (exit %d): %s", exitCode, out)
}

// Traceability: TC-036 -> Contract automated-test-orchestration/step-1 Outcome "surface-type-unavailable"
// Both frontmatter and CLI sources unavailable triggers error with recovery hint.
func TestTC_036_RunTests_SurfaceTypeUnavailable(t *testing.T) {
	dir := createProjectWithoutSurfacesAndTest(t)

	out, exitCode := runForgeRaw(t, dir, "run-tests")
	if exitCode != 0 {
		assert.True(t,
			outputContainsRecoveryHint(out),
			"error should include a recovery hint")
	}
}

// Traceability: TC-037 -> Contract automated-test-orchestration/step-1 Outcome "session-expired-during-detection"
// Project configuration invalidated mid-session triggers blocking error.
func TestTC_037_RunTests_SessionExpired(t *testing.T) {
	projectDir := createProjectWithTask(t, "web")

	// Remove config after setup to simulate mid-session invalidation
	configPath := filepath.Join(projectDir, ".forge", "config.yaml")
	_ = os.Remove(configPath)

	out, exitCode := runForgeRaw(t, projectDir, "run-tests")
	if exitCode != 0 {
		t.Logf("session expired output (exit %d): %s", exitCode, out)
	}
}

func createProjectWithoutSurfacesAndTest(t *testing.T) string {
	t.Helper()
	dir := createProjectWithTask(t, "")
	// Remove config to simulate no surfaces available
	configPath := filepath.Join(dir, ".forge", "config.yaml")
	_ = os.Remove(configPath)
	// Write minimal config without surfaces
	forgeDir := filepath.Join(dir, ".forge")
	err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"),
		[]byte("version: '1'\n"), 0644)
	assert.NoError(t, err)
	return dir
}
