//go:build e2e

package e2epipeline

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	testkit "forge-tests/testkit"
)

// runForgeInDir executes the forge binary in a specific directory.
func runForgeInDir(t *testing.T, dir string, args ...string) ([]byte, error) {
	t.Helper()
	cmd := exec.Command(testkit.ForgeBinary, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// setupTempProject creates a temp directory with .forge/config.yaml containing
// the specified test profile.
func setupTempProject(t *testing.T, profileName string) string {
	t.Helper()
	dir := t.TempDir()
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configContent := "languages:\n  - " + profileName + "\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// setupTempProjectWithE2E creates a temp project with tests/e2e/ directory.
func setupTempProjectWithE2E(t *testing.T, profileName string) string {
	t.Helper()
	dir := setupTempProject(t, profileName)
	e2eDir := filepath.Join(dir, "tests", "e2e")
	if err := os.MkdirAll(e2eDir, 0o755); err != nil {
		t.Fatal(err)
	}
	return dir
}

// withRetry retries a function until it succeeds or max retries exceeded.
func withRetry(t *testing.T, fn func() error, maxRetries int, delay time.Duration) {
	t.Helper()
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = fn(); err == nil {
			return
		}
		time.Sleep(delay)
	}
	t.Fatalf("retry exhausted: %s", err)
}
