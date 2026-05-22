package cmd

import (
	"fmt"
	"path/filepath"
	"sort"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var queryVerbose bool

var queryCmd = &cobra.Command{
	Use:   "query <task-id-or-key>",
	Short: "Query task information",
	Long: `Query and display information about a specific task.

<task-id-or-key> can be either:
  - Task ID (e.g., "1.2.3")
  - Task Key (e.g., "phase1-1.1.1-project-init")`,
	Args: cobra.ExactArgs(1),
	RunE: runQuery,
}

func init() {
	queryCmd.Flags().BoolVarP(&queryVerbose, "verbose", "v", false, "show all task fields including related fixes")
}

func runQuery(_ *cobra.Command, args []string) error {
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

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		Exit(ErrTaskNotFound(taskIDArg))
	}

	if queryVerbose {
		printVerboseQuery(key, t, featureSlug, index)
	} else {
		printDefaultQuery(t)
	}
	return nil
}

func printDefaultQuery(t *task.Task) {
	PrintBlockStart()
	PrintField("TASK_ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintFieldIfNotEmpty("SCOPE", t.Scope)
	if t.Breaking {
		PrintField("BREAKING", "true")
	}
	PrintBlockEnd()
}

func printVerboseQuery(key string, t *task.Task, featureSlug string, index *task.TaskIndex) {
	PrintBlockStart()
	PrintField("KEY", key)
	PrintField("TASK_ID", t.ID)
	PrintField("TITLE", t.Title)
	PrintField("STATUS", t.Status)
	PrintField("PRIORITY", t.Priority)
	PrintFieldIfNotEmpty("TYPE", t.Type)
	PrintFieldIfNotEmpty("SCOPE", t.Scope)
	if len(t.Dependencies) > 0 {
		PrintField("DEPENDENCIES:", "")
		for _, dep := range t.Dependencies {
			PrintListItem(dep)
		}
	}
	PrintField("TASK_FILE", feature.GetTaskFile(featureSlug, t.File))
	PrintField("RECORD_FILE", feature.GetTaskFile(featureSlug, t.Record))
	PrintBlockEnd()

	// RELATED_FIXES: find tasks whose SourceTaskID matches this task's ID
	var fixes []task.Task
	for _, ft := range index.TasksMap() {
		if ft.SourceTaskID == t.ID {
			fixes = append(fixes, ft)
		}
	}
	if len(fixes) > 0 {
		// Sort by ID for deterministic output
		sort.Slice(fixes, func(i, j int) bool {
			return fixes[i].ID < fixes[j].ID
		})
		PrintField("RELATED_FIXES:", "")
		for _, fix := range fixes {
			PrintListItem(fmt.Sprintf("%s [%s] %s", fix.ID, fix.Status, fix.Title))
		}
	}
}
