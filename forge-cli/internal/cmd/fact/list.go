package fact

import (
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/facttable"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var (
	listSource     string
	listConfidence string
)

var listCmd = &cobra.Command{
	Use:   "list [--source static|runtime|manual] [--confidence confirmed|inferred|assumed]",
	Short: "List fact summaries with optional filtering",
	Long: `List all facts in a summary table format, optionally filtered by source and/or confidence.

Each row shows: fact_id, subject, kind, source, confidence.

Examples:
  forge fact list
  forge fact list --source runtime
  forge fact list --confidence confirmed
  forge fact list --source runtime --confidence confirmed`,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runList,
}

func init() {
	listCmd.Flags().StringVar(&listSource, "source", "", "filter by source (static|runtime|manual)")
	listCmd.Flags().StringVar(&listConfidence, "confidence", "", "filter by confidence (confirmed|inferred|assumed)")
}

func runList(cmd *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	table, err := facttable.Load(projectRoot)
	if err != nil {
		return err
	}

	filtered := table.Filter(listSource, listConfidence)
	w := cmd.OutOrStdout()

	if len(filtered) == 0 {
		write(w, "no facts found\n")
		return nil
	}

	sorted := filtered.SortedEntries()

	// Column widths
	idCol := 35
	subjectCol := 25
	kindCol := 20
	sourceCol := 10
	confCol := 12

	// Header
	write(w, "%d facts found\n\n", len(sorted))
	write(w, "%s  %s  %s  %s  %s\n",
		base.PadRight("FACT_ID", idCol),
		base.PadRight("SUBJECT", subjectCol),
		base.PadRight("KIND", kindCol),
		base.PadRight("SOURCE", sourceCol),
		base.PadRight("CONFIDENCE", confCol),
	)

	// Separator
	write(w, "%s  %s  %s  %s  %s\n",
		strings.Repeat("-", idCol),
		strings.Repeat("-", subjectCol),
		strings.Repeat("-", kindCol),
		strings.Repeat("-", sourceCol),
		strings.Repeat("-", confCol),
	)

	// Rows
	for _, e := range sorted {
		write(w, "%s  %s  %s  %s  %s\n",
			base.PadRight(e.FactID, idCol),
			base.PadRight(base.TruncateSlug(e.Subject, subjectCol), subjectCol),
			base.PadRight(e.Kind, kindCol),
			base.PadRight(e.Source, sourceCol),
			base.PadRight(e.Confidence, confCol),
		)
	}

	return nil
}
