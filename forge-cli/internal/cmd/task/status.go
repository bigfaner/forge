package task

import (
	"forge-cli/internal/cmd/base"
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <task-id>",
	Short: "Query task status",
	Long: `Query the status of a task.

Use "forge task submit" to complete a task or "forge task reopen" to re-activate rejected/skipped tasks.`,
	Args: cobra.ExactArgs(1),
	RunE: runStatus,
}

func runStatus(_ *cobra.Command, args []string) error {
	taskIDArg := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		base.Exit(base.ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	// Query mode: display task status
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		base.Exit(base.ErrFileNotFound(indexPath))
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		base.Exit(base.ErrTaskNotFound(taskIDArg))
	}

	_ = key

	base.PrintBlockStart()
	base.PrintField("TASK_ID", t.ID)
	base.PrintField("STATUS", t.Status)
	base.PrintBlockEnd()
	return nil
}

// checkUnmetDeps returns dependency IDs that are not "completed" or "skipped".
// Delegates to task.GetUnmetDeps which handles both exact IDs and wildcard patterns.
func checkUnmetDeps(index *task.TaskIndex, t *task.Task) []string {
	return task.GetUnmetDeps(index, t.ID, t.Dependencies)
}
