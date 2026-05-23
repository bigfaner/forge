package cmd

import (
	"fmt"
	"os"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/project"
	"forge-cli/pkg/research"

	"github.com/spf13/cobra"
)

var researchCmd = &cobra.Command{
	Use:   "research [slug]",
	Short: "List or show research report details",
	Long: `List all research reports or show details for a specific report.

Without arguments: lists all research reports in table format.
With a slug argument: shows detailed information for that report.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runResearch,
}

func runResearch(_ *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	if len(args) == 0 {
		return runResearchList(projectRoot)
	}
	return runResearchDetail(projectRoot, args[0])
}

func runResearchList(projectRoot string) error {
	reports, err := research.Discover(projectRoot)
	if err != nil {
		return newErrResearchDiscovery(err)
	}

	if len(reports) == 0 {
		fmt.Fprintln(os.Stderr, "no research found")
		return nil
	}

	slugWidth := base.CalcSlugColWidth(mapReportsToSlugLens(reports))

	base.PrintBlockStart()
	base.PrintField("RESEARCH", fmt.Sprintf("%d found", len(reports)))
	fmt.Println()

	// Table header
	fmt.Printf("  %-s %-12s %-10s %s\n", base.PadRight("SLUG", slugWidth), "CREATED", "TOPIC", "MODE")
	fmt.Printf("  %-s %-12s %-10s %s\n",
		strings.Repeat("-", slugWidth),
		strings.Repeat("-", 10),
		strings.Repeat("-", 10),
		strings.Repeat("-", 10))

	for _, r := range reports {
		topic := "-"
		if r.Topic != "" {
			topic = r.Topic
		}
		mode := "-"
		if r.Mode != "" {
			mode = r.Mode
		}
		fmt.Printf("  %-s %-12s %-10s %s\n",
			base.PadRight(base.TruncateSlug(r.Slug, slugWidth), slugWidth),
			r.Created,
			topic,
			mode)
	}

	fmt.Println()
	base.PrintBlockEnd()
	return nil
}

func runResearchDetail(projectRoot, slug string) error {
	r, err := research.FindBySlug(projectRoot, slug)
	if err != nil {
		return newErrResearchNotFound(slug)
	}

	base.PrintBlockStart()
	base.PrintField("SLUG", r.Slug)
	base.PrintFieldIfNotEmpty("TOPIC", r.Topic)
	base.PrintFieldIfNotEmpty("CREATED", r.Created)
	base.PrintFieldIfNotEmpty("MODE", r.Mode)
	if len(r.Dimensions) > 0 {
		base.PrintField("DIMENSIONS", strings.Join(r.Dimensions, ", "))
	}
	base.PrintField("FILE", r.FilePath)
	base.PrintBlockEnd()
	return nil
}

// mapReportsToSlugLens extracts slug lengths from report list.
func mapReportsToSlugLens(reports []research.Report) []int {
	lens := make([]int, len(reports))
	for i, r := range reports {
		lens[i] = len(r.Slug)
	}
	return lens
}

func newErrResearchDiscovery(err error) *base.AIError {
	return base.NewAIError(
		base.ErrNotFound,
		"Failed to discover research reports",
		err.Error(),
		"Ensure docs/research/ directory exists",
		"ls docs/research/",
	)
}

func newErrResearchNotFound(slug string) *base.AIError {
	return base.NewAIError(
		base.ErrNotFound,
		fmt.Sprintf("Research report not found: %s", slug),
		"No research report .md file found with this slug",
		"Check the slug is correct (without .md extension)",
		"ls docs/research/",
	)
}
