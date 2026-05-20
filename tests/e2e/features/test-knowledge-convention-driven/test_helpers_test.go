//go:build e2e

package e2etestconv

import (
	"os/exec"

	e2etests "e2e-tests"
)

// forgeCmd returns an exec.Cmd for the forge CLI binary built from source.
func forgeCmd(args ...string) *exec.Cmd {
	return exec.Command(e2etests.ForgeBinary, args...)
}
