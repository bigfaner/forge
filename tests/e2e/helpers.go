//go:build e2e

package e2e

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// runCLI executes a CLI command and returns combined output.
func runCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI command failed: %s: %s", err, out)
	}
	return string(out)
}

// runCLIRaw executes a CLI command and returns output, exit code, and error.
// Unlike runCLI, it does not fatalf on non-zero exit — useful for negative tests.
func runCLIRaw(t *testing.T, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(args[0], args[1:]...)
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

// parseBlock extracts lines between "---" separators from raw CLI output.
// Returns the inner lines (without separators) or fails the test.
func parseBlock(t *testing.T, raw string) []string {
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

// hasField checks that a parsed block contains a "KEY: value" line.
func hasField(lines []string, key, value string) bool {
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

// hasNoField checks that a parsed block does NOT contain any line starting with key.
func hasNoField(lines []string, key string) bool {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return false
		}
	}
	return true
}

// fieldIndex returns the index of the line starting with key+": ", or -1.
func fieldIndex(lines []string, key string) int {
	prefix := key + ": "
	for i, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return i
		}
	}
	return -1
}

// fieldValue returns the value for the given key, or "" if not found.
func fieldValue(lines []string, key string) string {
	prefix := key + ": "
	for _, l := range lines {
		if strings.HasPrefix(l, prefix) {
			return strings.TrimPrefix(l, prefix)
		}
	}
	return ""
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
