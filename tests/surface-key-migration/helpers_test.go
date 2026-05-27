//go:build cli_functional

package surfacekeymigration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Shared helpers for surface-key-migration journey tests
// ==============================================================================

// createTempProjectWithConfig creates a temporary directory with go.mod and an
// optional .forge/config.yaml with the given YAML content for the surfaces field.
func createTempProjectWithConfig(t *testing.T, configYAML string) string {
	t.Helper()
	dir := t.TempDir()

	// go.mod for project root detection
	err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err, "failed to create go.mod")

	if configYAML != "" {
		forgeDir := filepath.Join(dir, ".forge")
		err := os.MkdirAll(forgeDir, 0755)
		assert.NoError(t, err, "failed to create .forge directory")
		err = os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configYAML), 0644)
		assert.NoError(t, err, "failed to write config.yaml")
	}

	return dir
}

// configWithSurfaces returns a minimal .forge/config.yaml with the given surfaces map.
func configWithSurfaces(surfaces string) string {
	return "version: '1'\nsurfaces:\n" + surfaces
}

// runForge runs the forge CLI in a given working directory, fatalf on failure.
func runForge(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+dir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("forge %s failed: %s: %s", stringsJoin(args, " "), err, out)
	}
	return string(out)
}

// runForgeRaw runs the forge CLI in a given working directory, returning output and exit code.
func runForgeRaw(t *testing.T, dir string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Dir = dir
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

// stringsJoin joins strings with a separator.
func stringsJoin(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	result := elems[0]
	for i := 1; i < len(elems); i++ {
		result += sep + elems[i]
	}
	return result
}

// surfaceResult represents a single entry in the JSON array output from forge surfaces.
type surfaceResult struct {
	Key string `json:"key"`
	Type string `json:"type"`
}

// parseSurfaceOutput attempts to parse forge surfaces CLI output as a JSON array.
func parseSurfaceOutput(t *testing.T, output string) surfaceResult {
	t.Helper()
	var results []surfaceResult
	err := json.Unmarshal([]byte(output), &results)
	assert.NoError(t, err, "failed to parse surface output as JSON: %s", output)
	if len(results) > 0 {
		return results[0]
	}
	return surfaceResult{}
}

// createFeatureWithTasks creates a feature directory structure with tasks.
func createFeatureWithTasks(t *testing.T, projectDir, featureSlug string, tasks map[string]map[string]string) {
	t.Helper()
	tasksDir := filepath.Join(projectDir, "docs", "features", featureSlug, "tasks")
	err := os.MkdirAll(tasksDir, 0755)
	assert.NoError(t, err, "failed to create tasks directory")

	indexEntries := make(map[string]interface{})
	for taskKey, frontmatterFields := range tasks {
		taskFile := filepath.Join(tasksDir, taskKey+".md")
		fm := "---\n"
		for k, v := range frontmatterFields {
			fm += k + ": " + v + "\n"
		}
		fm += "---\n\n# Task " + taskKey + "\n"
		err := os.WriteFile(taskFile, []byte(fm), 0644)
		assert.NoError(t, err, "failed to write task file %s", taskKey)

		indexEntries[taskKey] = map[string]interface{}{
			"id":     taskKey,
			"status": frontmatterFields["status"],
			"file":   taskKey + ".md",
		}
	}

	idxData, err := json.MarshalIndent(map[string]interface{}{
		"feature": featureSlug,
		"tasks":   indexEntries,
	}, "", "  ")
	assert.NoError(t, err, "failed to marshal index.json")
	err = os.WriteFile(filepath.Join(tasksDir, "index.json"), idxData, 0644)
	assert.NoError(t, err, "failed to write index.json")
}
