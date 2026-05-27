//go:build e2e

package testkit

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// forgeBinaryPath is set by SetForgeBinary. Defaults to "forge" for backward compatibility.
var forgeBinaryPath = "forge"

// SetForgeBinary sets the path to the forge CLI binary used by RunCLIExitCode.
// Call this from TestMain after building the binary from source.
func SetForgeBinary(path string) {
	forgeBinaryPath = path
}

// RunCLIExitCode executes a forge CLI command and returns exit code and combined output.
func RunCLIExitCode(args ...string) (int, string) {
	cmd := exec.Command(forgeBinaryPath, args...)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), string(out)
	}
	if err != nil {
		return 1, err.Error()
	}
	return 0, string(out)
}

// ProjectRoot resolves the project root directory by walking up from the
// source file location (via runtime.Caller) to find a go.mod marker.
// This mirrors the helpers.ts approach where __dirname is used as anchor.
func ProjectRoot(t *testing.T) string {
	t.Helper()
	// Use runtime.Caller to get this source file's path, then walk up.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed: cannot determine source file path")
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("cannot find project root (no go.mod found walking up from testkit source)")
		}
		dir = parent
	}
}

// ReadProjectFile reads and returns the content of a file relative to the
// project root. Fails the test if the file cannot be read.
func ReadProjectFile(t *testing.T, relPath string) string {
	t.Helper()
	root := ProjectRoot(t)
	path := filepath.Join(root, relPath)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read project file %q: %v", relPath, err)
	}
	return string(data)
}
