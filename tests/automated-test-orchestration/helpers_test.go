//go:build cli_functional

package automatedtestorchestration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	testkit "forge-tests/testkit"

	"github.com/stretchr/testify/assert"
)

// ==============================================================================
// Shared helpers for automated-test-orchestration journey tests
// ==============================================================================

// createProjectWithTask creates a temp project with a task file containing
// the given frontmatter fields.
func createProjectWithTask(t *testing.T, surfaceType string) string {
	t.Helper()
	dir := t.TempDir()

	err := os.WriteFile(filepath.Join(dir, "go.mod"),
		[]byte("module test-project\n\ngo 1.26\n"), 0644)
	assert.NoError(t, err)

	forgeDir := filepath.Join(dir, ".forge")
	err = os.MkdirAll(forgeDir, 0755)
	assert.NoError(t, err)

	surfacesConfig := ""
	if surfaceType != "" {
		surfacesConfig = "\nsurfaces:\n  test-surface: " + surfaceType
	}
	err = os.WriteFile(filepath.Join(forgeDir, "config.yaml"),
		[]byte("version: '1'"+surfacesConfig+"\n"), 0644)
	assert.NoError(t, err)

	// Create feature with task
	tasksDir := filepath.Join(dir, "docs", "features", "test-feature", "tasks")
	err = os.MkdirAll(tasksDir, 0755)
	assert.NoError(t, err)

	fm := "---\nid: T-run-tests\nstatus: in_progress\nfile: task-run-tests.md\n"
	if surfaceType != "" {
		fm += "surface-type: " + surfaceType + "\n"
	}
	fm += "---\n\n# Run tests task\n"
	err = os.WriteFile(filepath.Join(tasksDir, "task-run-tests.md"), []byte(fm), 0644)
	assert.NoError(t, err)

	idx := map[string]interface{}{
		"feature": "test-feature",
		"tasks": map[string]interface{}{
			"T-run-tests": map[string]string{
				"id": "T-run-tests", "status": "in_progress", "file": "task-run-tests.md",
			},
		},
	}
	if surfaceType != "" {
		idx["tasks"].(map[string]interface{})["T-run-tests"].(map[string]string)["surface-type"] = surfaceType
	}
	idxData, err := json.MarshalIndent(idx, "", "  ")
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tasksDir, "index.json"), idxData, 0644)
	assert.NoError(t, err)

	return dir
}

// runForgeRaw runs forge CLI and returns output + exit code.
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

// outputContainsRecoveryHint checks if output contains a recovery hint.
func outputContainsRecoveryHint(output string) bool {
	return strings.Contains(output, "hint") ||
		strings.Contains(output, "verify") ||
		strings.Contains(output, "check") ||
		strings.Contains(output, "configure") ||
		strings.Contains(output, "run")
}
