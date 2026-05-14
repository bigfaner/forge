//go:build e2e

package testkit

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// RunCLI executes a forge CLI command and returns combined output.
// Fails the test if the command exits non-zero.
func RunCLI(t *testing.T, args ...string) string {
	t.Helper()
	cmd := exec.Command("forge", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("CLI command failed: %s: %s", err, out)
	}
	return string(out)
}

// RunCLIExitCode executes a forge CLI command and returns exit code and combined output.
func RunCLIExitCode(args ...string) (int, string) {
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

// RunCLIWithResult executes a forge CLI command and returns stdout, stderr, and exit code.
func RunCLIWithResult(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
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
