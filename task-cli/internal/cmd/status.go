package cmd

import (
	"fmt"
	"os"
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
With status argument: update to new status.

Valid statuses: pending, in_progress, completed, blocked, skipped`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) {
	taskIDArg := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	key, t, err := findTask(index, taskIDArg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Query mode: only one argument
	if len(args) == 1 {
		PrintBlockStart()
		PrintField("KEY", key)
		PrintField("ID", t.ID)
		PrintField("STATUS", t.Status)
		PrintField("TITLE", t.Title)
		PrintBlockEnd()
		return
	}

	// Update mode: two arguments
	newStatus := args[1]

	// Validate status
	if !slices.Contains(index.StatusEnum, newStatus) {
		fmt.Fprintf(os.Stderr, "Error: invalid status '%s' (valid: %v)\n", newStatus, index.StatusEnum)
		os.Exit(1)
	}

	t.Status = newStatus
	index.Tasks[key] = *t

	if err := task.SaveIndex(indexPath, index); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	PrintBlockStart()
	PrintField("KEY", key)
	PrintField("ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintField("TITLE", t.Title)
	PrintFieldIfNotEmptySlice("DEPENDENCIES", t.Dependencies)
	PrintBlockEnd()
}
