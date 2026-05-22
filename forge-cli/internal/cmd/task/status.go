package task

import (
	"forge-cli/internal/cmd/base"
	"path/filepath"
	"strings"

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
// Supports both exact IDs and wildcard patterns (e.g. "1.x").
// selfID is excluded from wildcard matches to avoid self-matching.
func checkUnmetDeps(index *task.TaskIndex, t *task.Task) []string {
	var unmet []string
	for _, dep := range t.Dependencies {
		if strings.HasSuffix(dep, task.IDSuffixWildcard) {
			prefix := strings.TrimSuffix(dep, task.IDSuffixWildcard)
			prefixWithDot := prefix + "."
			found := false
			for _, other := range index.TasksMap() {
				if other.ID == t.ID {
					continue
				}
				if strings.HasPrefix(other.ID, prefixWithDot) && task.IsBusinessTask(other.ID) && other.Status != "completed" && other.Status != "skipped" {
					unmet = append(unmet, other.ID)
					found = true
				}
			}
			if !found {
				continue
			}
		} else {
			_, depTask, err := task.FindTask(index, dep)
			if err != nil || (depTask.Status != "completed" && depTask.Status != "skipped") {
				unmet = append(unmet, dep)
			}
		}
	}
	return unmet
}
