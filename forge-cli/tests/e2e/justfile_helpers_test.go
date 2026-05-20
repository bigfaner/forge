//go:build e2e

package e2e

import (
	"os/exec"
)

// runJust executes a just recipe with optional args via the system shell,
// returning exit code and combined output.
// Retained in e2e package for scope_resolution_cli_test.go.
func runJust(args ...string) (int, string) {
	cmd := exec.Command("just", args...)
	out, err := cmd.CombinedOutput()
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), string(out)
	}
	if err != nil {
		return 1, err.Error()
	}
	return 0, string(out)
}
