//go:build e2e

package e2e

import (
	"os"
	"testing"
)

// forgeBinary is an alias for ForgeBinary for backward compatibility within this package.
// The actual build happens in forge_binary.go init().
var forgeBinary = ForgeBinary

func TestMain(m *testing.M) {
	// Binary is already built by forge_binary.go init().
	// TestMain only drives test execution.
	code := m.Run()
	os.Exit(code)
}
