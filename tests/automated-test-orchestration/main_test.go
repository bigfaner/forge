//go:build e2e

package automatedtestorchestration

import (
	"testing"

	testkit "forge-tests/testkit"
)

func TestMain(m *testing.M) {
	// Ensure forge binary is built via testkit init
	_ = testkit.ForgeBinary
	m.Run()
}
