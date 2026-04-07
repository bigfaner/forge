package cmd

import (
	"path/filepath"
	"slices"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"
	"task-cli/pkg/task"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <task-id> [status]",
	Short: "Query or update task status",
	Long: `Query or update the status of a task.

Without status argument: query current status.
With status argument: update to new status.`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	taskIDArg := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		Exit(ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		Exit(ErrFileNotFound(indexPath))
	}

	key, t, err := findTask(index, taskIDArg)
	if err != nil {
		Exit(ErrTaskNotFound(taskIDArg))
	}

	// Query mode: only one argument
	if len(args) == 1 {
		PrintBlockStart()
		PrintField("KEY", key)
		PrintField("ID", t.ID)
		PrintField("STATUS", t.Status)
		PrintField("TITLE", t.Title)
		PrintFieldIfNotEmptySlice("DEPENDENCIES", t.Dependencies)
		PrintBlockEnd()
		return
	}

	// Update mode: two arguments
	newStatus := args[1]

	// Validate status
	if !slices.Contains(index.StatusEnum, newStatus) {
		Exit(ErrInvalidStatus(newStatus, index.StatusEnum))
	}

	t.Status = newStatus
	index.Tasks[key] = *t

	if err := task.SaveIndex(indexPath, index); err != nil {
		Exit(NewAIError(ErrConflict, "Failed to save index", err.Error(), "Check index.json is writable", "cat "+indexPath))
	}

	PrintBlockStart()
	PrintField("KEY", key)
	PrintField("ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintField("TITLE", t.Title)
	PrintFieldIfNotEmptySlice("DEPENDENCIES", t.Dependencies)
	PrintBlockEnd()
}
