package cmd

import (
	"fmt"
	"os"
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

	PrintBlockStart()
	PrintField("KEY", key)
	PrintField("ID", t.ID)
	PrintField("TITLE", t.Title)
	PrintField("STATUS", t.Status)
	PrintField("PRIORITY", t.Priority)
	PrintFieldIfNotEmpty("ESTIMATED_TIME", t.EstimatedTime)
	PrintFieldIfNotEmptySlice("DEPENDENCIES", t.Dependencies)
	PrintField("FILE", t.File)
	PrintField("RECORD", t.Record)
	PrintBlockEnd()
}
