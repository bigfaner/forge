package cmd

import (
	"fmt"
	"os"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"

	"github.com/spf13/cobra"
)

var featureCmd = &cobra.Command{
	Use:   "feature [slug]",
	Short: "Set or display the current feature",
	Long: `Set or display the current feature context.

Without arguments: displays the current feature.
With a slug argument: sets the current feature.`,
	Args: cobra.MaximumNArgs(1),
	Run:  runFeature,
}

func runFeature(cmd *cobra.Command, args []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(args) == 0 {
		// Display current feature
		slug, err := feature.GetCurrentFeature(projectRoot)
		if err != nil {
			PrintField("FEATURE", "(none)")
			return
		}
		PrintField("FEATURE", slug)
		return
	}

	// Set feature
	slug := args[0]
	if err := feature.SetFeature(projectRoot, slug); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	PrintField("FEATURE", slug)
}
