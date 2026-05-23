package feature

import (
	"os"
	"testing"
)

// TestMain ensures Register() is called so subcommands are available in tests.
func TestMain(m *testing.M) {
	Register()
	os.Exit(m.Run())
}
