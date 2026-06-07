//go:build cli_functional

package surfacerecipegeneration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Shared helpers for surface-aware-recipe-generation journey tests
// ==============================================================================

// createProjectWithSurfaces creates a temp project with go.mod and .forge/config.yaml
// containing the given surfaces configuration.
func createProjectWithSurfaces(t *testing.T, surfacesYAML string) string {
	t.Helper()
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err, "failed to create go.mod")

	forgeDir := filepath.Join(dir, ".forge")
	err = os.MkdirAll(forgeDir, 0755)
	assert.NoError(t, err, "failed to create .forge directory")

	configContent := "version: '1'\nsurfaces:\n" + surfacesYAML
	err = os.WriteFile(filepath.Join(forgeDir, "config.yaml"),
		[]byte(configContent), 0644)
	assert.NoError(t, err, "failed to write config.yaml")

	return dir
}

// createProjectWithoutSurfaces creates a temp project with config.yaml but no surfaces field.
func createProjectWithoutSurfaces(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err)

	forgeDir := filepath.Join(dir, ".forge")
	err = os.MkdirAll(forgeDir, 0755)
	assert.NoError(t, err)

	err = os.WriteFile(filepath.Join(forgeDir, "config.yaml"),
		[]byte("version: '1'\n"), 0644)
	assert.NoError(t, err)

	return dir
}

// runForgeRaw runs forge CLI and returns output + exit code.
func runForgeRaw(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Dir = dir
	cmd.Stdin = strings.NewReader("")
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+dir)
	out, err := cmd.CombinedOutput()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}
	return string(out), exitCode
}
