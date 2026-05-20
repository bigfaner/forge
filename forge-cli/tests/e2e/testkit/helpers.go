//go:build e2e

package testkit

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// forgeBinaryPath is set by SetForgeBinary. Defaults to "forge" for backward compatibility.
var forgeBinaryPath = "forge"

// SetForgeBinary sets the path to the forge CLI binary used by RunCLI, RunCLIExitCode,
// and RunCLIWithResult. Call this from TestMain after building the binary from source.
func SetForgeBinary(path string) {
	forgeBinaryPath = path
}

// RunCLI executes a forge CLI command and returns combined output.
// Fails the test if the command exits non-zero.
func RunCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command(forgeBinaryPath, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI command failed: %s: %s", err, out)
	}
	return string(out)
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

// RunCLIWithResult executes a forge CLI command and returns stdout, stderr, and exit code.
func RunCLIWithResult(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(forgeBinaryPath, args...)
	var stdoutBuf, stderrBuf strings.Builder
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return stdoutBuf.String(), stderrBuf.String(), exitErr.ExitCode()
	}
	if err != nil {
		return stdoutBuf.String(), err.Error(), 1
	}
	return stdoutBuf.String(), stderrBuf.String(), 0
}

// WithRetry retries a function until it succeeds or max retries exceeded.
func WithRetry(t *testing.T, fn func() error, maxRetries int, delay time.Duration) {
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

// ProjectFileExists returns true if a file exists at the given path relative
// to the project root.
func ProjectFileExists(relPath string) bool {
	// Cannot use t.Helper() here — no *testing.T param. This is intentional:
	// it returns a bool and does not fail the test, matching the helpers.ts signature.
	// Resolve project root without a *testing.T by using runtime.Caller directly.
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		return false
	}
	dir := filepath.Dir(thisFile)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
	path := filepath.Join(dir, relPath)
	_, err := os.Stat(path)
	return err == nil
}

// FileContains asserts that the file at filePath contains the given substring.
// Fails the test if the file cannot be read or does not contain the substring.
func FileContains(t *testing.T, filePath, substring string) {
	t.Helper()
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("cannot read file %q: %v", filePath, err)
	}
	if !strings.Contains(string(data), substring) {
		t.Fatalf("file %q does not contain substring %q", filePath, substring)
	}
}

// FileNotContains asserts that the file at filePath does NOT contain the given
// substring. Fails the test if the file cannot be read or does contain the substring.
func FileNotContains(t *testing.T, filePath, substring string) {
	t.Helper()
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("cannot read file %q: %v", filePath, err)
	}
	if strings.Contains(string(data), substring) {
		t.Fatalf("file %q should not contain substring %q but it does", filePath, substring)
	}
}
