//go:build e2e

package e2e

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup: verify forge CLI is available in PATH
	// Teardown: no cleanup needed for CLI tests
	code := m.Run()
	os.Exit(code)
}
