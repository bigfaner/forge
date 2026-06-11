// Package forensic contains all forge forensic subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
package forensic

import "github.com/spf13/cobra"

// Cmd is the parent forensic command, exported for use by the cmd package.
var Cmd = &cobra.Command{
	Use:   "forensic",
	Short: "Analyze session transcripts for agent deviation forensics",
	Long: `Extract and analyze evidence from Claude Code session transcripts.

Subcommands:
  search    Find sessions in ~/.claude/history.jsonl
  extract   Extract thinking/tool chains from a session JSONL
  subagents List subagent transcripts for a session`,
	Args: cobra.NoArgs,
}

// Register adds all forensic subcommands to Cmd.
func Register() {
	Cmd.AddCommand(
		searchCmd,
		extractCmd,
		subagentsCmd,
	)
}
