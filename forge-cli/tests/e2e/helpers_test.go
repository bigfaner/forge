//go:build e2e

package e2e

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// runCLI executes a forge CLI command and returns combined output.
func runCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("forge", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI command failed: %s: %s", err, out)
	}
	return string(out)
}

// runCLIWithResult executes a forge CLI command and returns stdout, stderr, and exit code.
func runCLIWithResult(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command("forge", args...)
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

// runCLIExitCode executes a forge CLI command and returns exit code and combined output.
func runCLIExitCode(args ...string) (int, string) {
	cmd := exec.Command("forge", args...)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), string(out)
	}
	if err != nil {
		return 1, err.Error()
	}
	return 0, string(out)
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
