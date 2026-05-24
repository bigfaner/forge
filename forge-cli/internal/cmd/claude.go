package cmd

import (
	"fmt"
	"os/exec"

	"forge-cli/internal/cmd/base"

	"github.com/spf13/cobra"
)

// lookPathFunc resolves a binary name to its full path.
// Variable for testability.
var lookPathFunc = exec.LookPath

// runClaudeFunc executes claude with the given args.
// Variable for testability.
var runClaudeFunc = base.RunClaude

var claudeCmd = &cobra.Command{
	Use:                "claude [flags] [args]",
	Short:              "Launch Claude CLI with permissions skipped",
	Long:               `Launches the Claude CLI binary with --dangerously-skip-permissions always injected. All user args pass through transparently to claude.`,
	DisableFlagParsing: true,
	RunE:               runClaude,
}

func runClaude(cmd *cobra.Command, args []string) error {
	// Pre-flight: verify claude binary exists in PATH
	if _, err := lookPathFunc("claude"); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "error: claude binary not found in PATH\n")
		return fmt.Errorf("claude: %w", err)
	}

	// Always prepend --dangerously-skip-permissions
	allArgs := append([]string{"--dangerously-skip-permissions"}, args...)
	return runClaudeFunc(allArgs)
}
