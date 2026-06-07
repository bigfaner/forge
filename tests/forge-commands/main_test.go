//go:build cli_functional

package forgecommands

import (
	"testing"

	testkit "forge-tests/testkit"
)

func TestMain(m *testing.M) {
	// Ensure forge binary is built via testkit init
	_ = testkit.ForgeBinary
	m.Run()
}
