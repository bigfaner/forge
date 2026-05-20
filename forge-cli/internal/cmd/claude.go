package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// lookPathFunc resolves a binary name to its full path.
// Variable for testability.
var lookPathFunc = exec.LookPath

// runClaudeFunc executes claude with the given args.
// Variable for testability.
var runClaudeFunc = defaultRunClaude

// claudeSupportsContinueFlagFunc checks whether the installed claude CLI
// supports the -c / --continue flag. Overridable for testing.
var claudeSupportsContinueFlagFunc = defaultClaudeSupportsContinueFlag

func defaultRunClaude(args []string) error {
	cmd := exec.Command("claude", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// defaultClaudeSupportsContinueFlag probes claude --help for the -c flag.
// Returns true if the flag is present, false otherwise.
func defaultClaudeSupportsContinueFlag() bool {
	cmd := exec.Command("claude", "--help")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	helpText := string(output)
	// Look for "-c" or "--continue" in the help output
	return strings.Contains(helpText, "-c,") || strings.Contains(helpText, "--continue")
}

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
