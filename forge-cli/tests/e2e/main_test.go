//go:build e2e

package e2e

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup: verify forge binary is available
	code := m.Run()
	// Teardown: clean up any test artifacts
	os.Exit(code)
}
