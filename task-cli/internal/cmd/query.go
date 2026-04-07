package cmd

import (
	"path/filepath"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"
	"task-cli/pkg/task"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query <task-id-or-key>",
	Short: "Query task information",
	Long: `Query and display information about a specific task.

<task-id-or-key> can be either:
  - Task ID (e.g., "1.2.3")
  - Task Key (e.g., "phase1-1.1.1-project-init")`,
	Args: cobra.ExactArgs(1),
	Run:  runQuery,
}

func runQuery(cmd *cobra.Command, args []string) {
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

	PrintBlockStart()
	PrintField("KEY", key)
	PrintField("ID", t.ID)
	PrintField("TITLE", t.Title)
	PrintField("STATUS", t.Status)
	PrintField("PRIORITY", t.Priority)
	PrintFieldIfNotEmpty("ESTIMATED_TIME", t.EstimatedTime)
	PrintFieldIfNotEmptySlice("DEPENDENCIES", t.Dependencies)
	PrintField("FILE", filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, t.File)))
	PrintField("RECORD", filepath.Join(projectRoot, feature.GetRecordFile(featureSlug, t.Record)))
	PrintBlockEnd()
}
