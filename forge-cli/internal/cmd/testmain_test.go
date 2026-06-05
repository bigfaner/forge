package cmd

import (
	"os"
	"testing"

	"forge-cli/pkg/forgelog"
)

// TestMain disables forgelog file backend in tests to prevent file handle
// leaks that cause Windows "file in use" errors during temp directory cleanup.
// Tests that need to verify forgelog file behavior set FORGE_NO_LOG= explicitly
// or use forgelog.Init() directly.
func TestMain(m *testing.M) {
	_ = os.Setenv("FORGE_NO_LOG", "1")
	code := m.Run()
	forgelog.Close()
	os.Exit(code)
}
