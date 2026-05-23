package prompt

import (
	"fmt"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	promptpkg "forge-cli/pkg/prompt"

	"github.com/spf13/cobra"
)

var promptFixRecordMissed bool

var promptGetCmd = &cobra.Command{
	Use:   "get-by-task-id <id>",
	Short: "Synthesize the agent prompt for a task",
	Long: `Synthesize and print the agent prompt for the given task ID.

The prompt is selected based on the task's type field and rendered with
runtime values (task file path, record file path, scope, feature slug).

Use --fix-record-missed to use the fix-record-missed recovery template
regardless of the task's type.`,
	Args: cobra.ExactArgs(1),
	RunE: runPrompt,
}

func init() {
	promptGetCmd.Flags().BoolVar(&promptFixRecordMissed, "fix-record-missed", false, "Use fix-record-missed template")
}

func runPrompt(_ *cobra.Command, args []string) error {
	taskID := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		base.Exit(base.ErrFeatureNotSet())
	}

	opts := promptpkg.SynthesizeOpts{
		ProjectRoot:     projectRoot,
		FeatureSlug:     featureSlug,
		TaskID:          taskID,
		FixRecordMissed: promptFixRecordMissed,
	}

	result, err := promptpkg.Synthesize(opts)
	if err != nil {
		base.Exit(fmt.Errorf("%w", err))
	}

	fmt.Print(result)
	return nil
}
