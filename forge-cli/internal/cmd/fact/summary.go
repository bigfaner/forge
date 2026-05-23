package fact

import (
	"io"
	"sort"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/facttable"
	"forge-cli/pkg/project"

	"github.com/spf13/cobra"
)

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show fact table statistics",
	Long: `Show statistics about the Fact Table grouped by source, confidence, and kind.

Displays counts for each group and the total number of facts.

Examples:
  forge fact summary`,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          runSummary,
}

func runSummary(cmd *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	table, err := facttable.Load(projectRoot)
	if err != nil {
		return err
	}

	stats := table.Summary()
	w := cmd.OutOrStdout()

	write(w, "TOTAL: %d facts\n\n", stats.Total)

	printGroup(w, "BY SOURCE", stats.BySource, facttable.ValidSources)
	printGroup(w, "BY CONFIDENCE", stats.ByConfidence, facttable.ValidConfidences)
	printGroup(w, "BY KIND", stats.ByKind, facttable.ValidKinds)

	return nil
}

func printGroup(w io.Writer, title string, counts map[string]int, validValues []string) {
	write(w, "[%s]\n", title)

	if len(counts) == 0 {
		write(w, "\n")
		return
	}

	keys := make([]string, 0, len(counts))
	seen := make(map[string]bool)
	for _, v := range validValues {
		if _, ok := counts[v]; ok {
			keys = append(keys, v)
		}
		seen[v] = true
	}
	extras := make([]string, 0)
	for k := range counts {
		if !seen[k] {
			extras = append(extras, k)
		}
	}
	sort.Strings(extras)
	keys = append(keys, extras...)

	maxKeyLen := 0
	for _, k := range keys {
		if len(k) > maxKeyLen {
			maxKeyLen = len(k)
		}
	}

	for _, k := range keys {
		write(w, "  %s  %d\n", base.PadRight(strings.ToUpper(k), maxKeyLen), counts[k])
	}
	write(w, "\n")
}
