package fact

import (
	"encoding/json"
	"fmt"
	"io"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/facttable"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <fact_id>",
	Short: "View a single fact's full content",
	Long: `View the complete content of a single fact entry by its fact_id.

The value field is pretty-printed as JSON.

Examples:
  forge fact get cli.forge-signature-12345`,
	Args:          cobra.ExactArgs(1),
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runGet,
}

// write writes to w, ignoring errors (interactive output is best-effort).
func write(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func runGet(cmd *cobra.Command, args []string) error {
	factID := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	table, err := facttable.Load(projectRoot)
	if err != nil {
		return err
	}

	entry := table.GetByID(factID)
	if entry == nil {
		return base.NewAIError(
			base.ErrNotFound,
			fmt.Sprintf("Fact not found: %s", factID),
			"No fact entry with this fact_id exists in the table",
			"Verify the fact_id is correct",
			"forge fact list",
		)
	}

	w := cmd.OutOrStdout()

	write(w, "FACT_ID:     %s\n", entry.FactID)
	write(w, "SOURCE:      %s\n", entry.Source)
	write(w, "SUBJECT:     %s\n", entry.Subject)
	write(w, "KIND:        %s\n", entry.Kind)
	write(w, "CONFIDENCE:  %s\n", entry.Confidence)
	write(w, "UPDATED_AT:  %s\n", entry.UpdatedAt)

	// Pretty-print value
	var value interface{}
	if err := json.Unmarshal(entry.Value, &value); err == nil {
		pretty, _ := json.MarshalIndent(value, "             ", "  ")
		write(w, "VALUE:       %s\n", string(pretty))
	} else {
		write(w, "VALUE:       %s\n", string(entry.Value))
	}

	return nil
}
