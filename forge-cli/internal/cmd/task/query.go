package task

import (
	"fmt"
	"forge-cli/internal/cmd/base"
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
		base.Exit(base.ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		base.Exit(base.ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		base.Exit(base.ErrFileNotFound(indexPath))
	}

	// Check for legacy scope fields
	var allTasks []task.Task
	for _, t := range index.TasksMap() {
		allTasks = append(allTasks, t)
	}
	if legacyErr := task.CheckLegacyScope(allTasks); legacyErr != nil {
		scopeErr, ok := legacyErr.(*task.LegacyScopeError)
		if ok {
			base.Exit(base.ErrLegacyScope(scopeErr.Count))
		}
		base.Exit(legacyErr)
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		base.Exit(base.ErrTaskNotFound(taskIDArg))
	}

	if queryVerbose {
		printVerboseQuery(key, t, featureSlug, index)
	} else {
		printDefaultQuery(t)
	}
	return nil
}

func printDefaultQuery(t *task.Task) {
	base.PrintBlockStart()
	base.PrintField("TASK_ID", t.ID)
	base.PrintField("STATUS", string(t.Status))
	base.PrintFieldIfNotEmpty("TYPE", t.Type)
	base.PrintFieldIfNotEmpty("TASK_CATEGORY", task.CategoryForType(t.Type))
	base.PrintFieldIfNotEmpty("SURFACE_KEY", t.SurfaceKey)
	base.PrintFieldIfNotEmpty("SURFACE_TYPE", t.SurfaceType)
	if t.Breaking {
		base.PrintField("BREAKING", "true")
	}
	base.PrintBlockEnd()
}

func printVerboseQuery(key string, t *task.Task, featureSlug string, index *task.TaskIndex) {
	base.PrintBlockStart()
	base.PrintField("KEY", key)
	base.PrintField("TASK_ID", t.ID)
	base.PrintField("TITLE", t.Title)
	base.PrintField("STATUS", string(t.Status))
	base.PrintField("PRIORITY", string(t.Priority))
	base.PrintFieldIfNotEmpty("TYPE", t.Type)
	base.PrintFieldIfNotEmpty("TASK_CATEGORY", task.CategoryForType(t.Type))
	base.PrintFieldIfNotEmpty("SURFACE_KEY", t.SurfaceKey)
	base.PrintFieldIfNotEmpty("SURFACE_TYPE", t.SurfaceType)
	if len(t.Dependencies) > 0 {
		base.PrintField("DEPENDENCIES:", "")
		for _, dep := range t.Dependencies {
			base.PrintListItem(dep)
		}
	}
	base.PrintField("TASK_FILE", feature.GetTaskFile(featureSlug, t.File))
	base.PrintField("RECORD_FILE", feature.GetTaskFile(featureSlug, t.Record))
	base.PrintBlockEnd()

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
		base.PrintField("RELATED_FIXES:", "")
		for _, fix := range fixes {
			base.PrintListItem(fmt.Sprintf("%s [%s] %s", fix.ID, fix.Status, fix.Title))
		}
	}
}
