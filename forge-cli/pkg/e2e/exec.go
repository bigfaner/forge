package e2e

import (
	"os/exec"
)

// ExecRunner is an injectable wrapper for testability.
type ExecRunner interface {
	Run(name string, args ...string) ([]byte, error)
}

// RealExec is the production implementation using os/exec.
type RealExec struct{}

// Run executes the named command with args and returns combined output.
func (RealExec) Run(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

// runner is the package-level ExecRunner used by action functions.
// Production code uses RealExec; tests override with stubExec.
var runner ExecRunner = RealExec{}
