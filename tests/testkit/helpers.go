//go:build e2e

package testkit

import (
	"os/exec"
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
