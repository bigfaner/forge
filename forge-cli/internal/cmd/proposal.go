package cmd

import (
	"fmt"
	"os"
	"strings"

	"forge-cli/pkg/project"
	"forge-cli/pkg/proposal"

	"github.com/spf13/cobra"
)

var proposalCmd = &cobra.Command{
	Use:   "proposal [slug]",
	Short: "List or show proposal details",
	Long: `List all proposals or show details for a specific proposal.

Without arguments: lists all proposals in table format.
With a slug argument: shows detailed information for that proposal.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProposal,
}

func runProposal(_ *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	if len(args) == 0 {
		runProposalList(projectRoot)
	} else {
		runProposalDetail(projectRoot, args[0])
	}
	return nil
}

func runProposalList(projectRoot string) {
	proposals, err := proposal.Discover(projectRoot)
	if err != nil {
		Exit(newErrProposalDiscovery(err))
	}

	if len(proposals) == 0 {
		fmt.Fprintln(os.Stderr, "no proposals found")
		return
	}

	// Proposals are already sorted by Created descending (newest first)
	// via infocmd.Discover.

	// Calculate dynamic slug column width.
	slugWidth := CalcSlugColWidth(mapProposalsToSlugLens(proposals))

	PrintBlockStart()
	PrintField("PROPOSALS", fmt.Sprintf("%d found", len(proposals)))
	fmt.Println()

	// Table header
	fmt.Printf("  %-s %-12s %-10s %-4s %s\n", PadRight("SLUG", slugWidth), "CREATED", "STATUS", "PRD", "FEATURE")
	fmt.Printf("  %-s %-12s %-10s %-4s %s\n",
		strings.Repeat("-", slugWidth),
		strings.Repeat("-", 10),
		strings.Repeat("-", 8),
		strings.Repeat("-", 3),
		strings.Repeat("-", 10))

	for _, p := range proposals {
		prdMark := "-"
		if p.HasPRD {
			prdMark = "yes"
		}
		featureStatus := "-"
		if p.FeatureStatus != "" {
			featureStatus = p.FeatureStatus
		}
		fmt.Printf("  %-s %-12s %-10s %-4s %s\n",
			PadRight(TruncateSlug(p.Slug, slugWidth), slugWidth),
			p.Created,
			p.Status,
			prdMark,
			featureStatus)
	}

	fmt.Println()
	PrintBlockEnd()
}

func runProposalDetail(projectRoot, slug string) {
	p, err := proposal.FindBySlug(projectRoot, slug)
	if err != nil {
		Exit(newErrProposalNotFound(slug))
	}

	prdMark := "no"
	if p.HasPRD {
		prdMark = "yes"
	}
	featureStatus := "(none)"
	if p.FeatureStatus != "" {
		featureStatus = p.FeatureStatus
	}

	PrintBlockStart()
	PrintField("SLUG", p.Slug)
	PrintField("CREATED", p.Created)
	PrintField("STATUS", p.Status)
	PrintFieldIfNotEmpty("AUTHOR", p.Author)
	PrintField("PRD", prdMark)
	PrintField("FEATURE", featureStatus)
	PrintField("FILE", p.FilePath)
	PrintBlockEnd()
}

// mapProposalsToSlugLens extracts slug lengths from proposal list.
func mapProposalsToSlugLens(proposals []proposal.Proposal) []int {
	lens := make([]int, len(proposals))
	for i, p := range proposals {
		lens[i] = len(p.Slug)
	}
	return lens
}

func newErrProposalDiscovery(err error) *AIError {
	return NewAIError(
		ErrNotFound,
		"Failed to discover proposals",
		err.Error(),
		"Ensure docs/proposals/ directory exists",
		"ls docs/proposals/",
	)
}

func newErrProposalNotFound(slug string) *AIError {
	return NewAIError(
		ErrNotFound,
		fmt.Sprintf("Proposal not found: %s", slug),
		"No proposal.md found for this slug",
		"Check slug is correct",
		"ls docs/proposals/",
	)
}
