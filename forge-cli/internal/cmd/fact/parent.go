// Package fact contains all forge fact subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
package fact

import (
	"github.com/spf13/cobra"
)

// Cmd is the parent fact command, exported for use by the cmd package.
var Cmd = &cobra.Command{
	Use:   "fact",
	Short: "Manage the Fact Table",
	Long: `Manage .forge/fact-table.json for structured system facts.

The Fact Table stores facts about the system under test (signatures,
output formats, error codes, side effects, etc.) used by test generation
and run-to-learn pipelines.

Subcommands:
  list     List fact summaries with optional filtering
  get      View a single fact's full content
  summary  Show statistics grouped by source/confidence/kind`,
	Args: cobra.NoArgs,
}
