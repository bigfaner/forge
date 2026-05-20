//go:build e2e

package e2e

import (
	"os"
	"testing"
)

// forgeBinary is an alias for ForgeBinary for backward compatibility within this package.
// Set in TestMain (after init()) to capture the value populated by forge_binary.go init().
var forgeBinary string

func TestMain(m *testing.M) {
	// Binary is already built by forge_binary.go init().
	// Copy the value now so all test functions see the correct path.
	forgeBinary = ForgeBinary

	code := m.Run()
	os.Exit(code)
}
