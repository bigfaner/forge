package cmd

import (
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

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

func runQuery(_ *cobra.Command, args []string) {
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

	_, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		Exit(ErrTaskNotFound(taskIDArg))
	}

	PrintBlockStart()
	PrintField("TASK_ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintFieldIfNotEmpty("SCOPE", t.Scope)
	if t.Breaking {
		PrintField("BREAKING", "true")
	}
	PrintBlockEnd()
}
