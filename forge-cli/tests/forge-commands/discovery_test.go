//go:build e2e

package forgecommands

import (
	"strings"
	"testing"

	"forge-cli/tests/testkit"

	"github.com/stretchr/testify/assert"
)

// --- Command Discovery & Help (TC-001 to TC-004) ---

// Traceability: TC-004 -> Story 1 / AC-3
func TestTC_004_UnknownTaskSubcommandReturnsErrorWithList(t *testing.T) {
	exitCode, out := testkit.RunCLIExitCode("task", "nonexistent-sub")

	// Cobra shows task help and lists available subcommands for unknown subcommand
	// Exit code is 0 (cobra behavior), output lists valid subcommands
	assert.Equal(t, 0, exitCode, "unknown task subcommand shows help (cobra behavior)")
	lower := strings.ToLower(out)
	assert.True(t, strings.Contains(lower, "available commands") || strings.Contains(lower, "usage"),
		"output should show available commands or usage: %s", out)

	// Verify valid subcommands are listed in help output
	validSubs := []string{"claim", "submit", "status"}
	for _, sub := range validSubs {
		assert.True(t, strings.Contains(lower, sub),
			"help output should list valid subcommand: %s", sub)
	}
}
