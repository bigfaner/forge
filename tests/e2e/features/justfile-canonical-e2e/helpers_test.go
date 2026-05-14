//go:build e2e

package justfile_canonical_e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// forgeBinary returns the path to the forge binary, built if needed.
// It caches the path in a package variable for reuse across tests.
var forgeBinaryPath string

func forgeBinary(t *testing.T) string {
	t.Helper()
	if forgeBinaryPath != "" {
		if _, err := os.Stat(forgeBinaryPath); err == nil {
			return forgeBinaryPath
		}
	}
	binName := "forge"
	if runtime.GOOS == "windows" {
		binName = "forge.exe"
	}
	binPath := filepath.Join("..", "..", "..", "..", "forge-cli", "bin", binName)
	if _, err := os.Stat(binPath); err != nil {
		// Build the binary
		buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/forge/")
		buildCmd.Dir = filepath.Join("..", "..", "..", "..", "forge-cli")
		if out, err := buildCmd.CombinedOutput(); err != nil {
			t.Fatalf("failed to build forge binary: %s: %s", err, out)
		}
	}
	absPath, err := filepath.Abs(binPath)
	if err != nil {
		t.Fatalf("failed to resolve binary path: %s", err)
	}
	forgeBinaryPath = absPath
	return absPath
}

// runForge executes the forge binary with given args and returns combined output.
func runForge(t *testing.T, args ...string) ([]byte, error) {
	t.Helper()
	bin := forgeBinary(t)
	cmd := exec.Command(bin, args...)
	return cmd.CombinedOutput()
}

// runForgeInDir executes the forge binary in a specific directory.
func runForgeInDir(t *testing.T, dir string, args ...string) ([]byte, error) {
	t.Helper()
	bin := forgeBinary(t)
	cmd := exec.Command(bin, args...)
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
	configContent := "test-profiles:\n  - " + profileName + "\n"
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
