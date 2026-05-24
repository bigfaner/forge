package base

import (
	"os"
	"os/exec"
)

// RunClaude executes the claude CLI binary with the given arguments,
// connecting stdin/stdout/stderr to the current process.
func RunClaude(args []string) error {
	cmd := exec.Command("claude", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
