//go:build cli_functional

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

// RunCLI executes a forge CLI command and returns combined output.
func RunCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command(ForgeBinary, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI command failed: %s: %s", err, out)
	}
	return string(out)
}

// RunCLIRaw executes a forge CLI command and returns output and exit code.
// Unlike RunCLI, it does not fatalf on non-zero exit -- useful for negative tests.
func RunCLIRaw(t *testing.T, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(ForgeBinary, args...)
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

// ParseBlock extracts lines between "---" separators from raw CLI output.
// Returns the inner lines (without separators) or fails the test.
func ParseBlock(t *testing.T, raw string) []string {
	t.Helper()
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 2 || strings.TrimSpace(lines[0]) != "---" || strings.TrimSpace(lines[len(lines)-1]) != "---" {
		t.Fatalf("output must be wrapped in --- separators, got:\n%s", raw)
	}
	inner := lines[1 : len(lines)-1]
	result := make([]string, 0, len(inner))
	for _, l := range inner {
		result = append(result, strings.TrimSpace(l))
	}
	return result
}

// HasField checks that a parsed block contains a "KEY: value" line.
func HasField(lines []string, key, value string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			if value == "" {
				return true
			}
			return l == prefix+value
		}
	}
	return false
}

// HasNoField checks that a parsed block does NOT contain any line starting with key.
func HasNoField(lines []string, key string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return false
		}
	}
	return true
}

// FieldIndex returns the index of the line starting with key+": ", or -1.
func FieldIndex(lines []string, key string) int {
	prefix := key + ": "
	for i, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return i
		}
	}
	return -1
}

// FieldValue returns the value for the given key, or "" if not found.
func FieldValue(lines []string, key string) string {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return strings.TrimPrefix(l, prefix)
		}
	}
	return ""
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

// RunCLIExitCode executes a forge CLI command and returns exit code and combined output.
// Unlike RunCLIRaw, it does not take a *testing.T and does not fatalf — it simply
// reports the exit code and output, leaving assertion decisions to the caller.
func RunCLIExitCode(args ...string) (int, string) {
	cmd := exec.Command(ForgeBinary, args...)
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
// Since tests/ is an independent Go module (forge-tests), it looks for
// the tests/go.mod file specifically.
func ProjectRoot(t *testing.T) string {
	t.Helper()
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
