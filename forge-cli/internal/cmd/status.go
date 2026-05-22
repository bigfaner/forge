package cmd

import (
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
	Args: cobra.RangeArgs(1, 2),
	RunE: runStatus,
}

func runStatus(_ *cobra.Command, args []string) error {
	taskIDArg := args[0]

	// Status command is now read-only: reject 2-arg mutation calls
	if len(args) > 1 {
		Exit(NewAIError(ErrInvalidInput,
			"task status is read-only. Use forge task submit to complete a task.",
			"", "", ""))
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		Exit(ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	// Query mode: display task status
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		Exit(ErrFileNotFound(indexPath))
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		Exit(ErrTaskNotFound(taskIDArg))
	}

	_ = key

	PrintBlockStart()
	PrintField("TASK_ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintBlockEnd()
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
