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

var verifyCompletionCmd = &cobra.Command{
	Use:   "verify-completion",
	Short: "Verify task completion before git commit",
	Long: `Verifies task completion before allowing git commit.

Checks:
  - Task status is "completed"
  - Record file exists (if specified)

Exit codes:
  0 - All checks passed
  2 - Verification failed`,
	Run: runVerifyCompletion,
}

func runVerifyCompletion(cmd *cobra.Command, args []string) {
	if err := verifyTaskCompletion(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	os.Exit(0)
}

// verifyTaskCompletion verifies task completion before allowing git commit.
// Returns error if task is not completed or record file is missing.
func verifyTaskCompletion() error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil // No project context, allow commit
	}

	// Get current feature
	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		return nil // No feature set, allow commit
	}

	// Read state file
	statePath := feature.GetTaskStatePath(projectRoot, featureSlug)
	state, err := task.LoadState(statePath)
	if err != nil {
		return nil // No state file, allow commit
	}
	if state == nil || state.TaskID == "" {
		return nil // No active task, allow commit
	}

	// Load index
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return fmt.Errorf("cannot read index.json: %w", err)
	}

	// Find task and check status
	var foundTask *task.Task
	for _, t := range index.TasksMap() {
		if t.ID == state.TaskID {
			foundTask = &t
			break
		}
	}

	if foundTask == nil {
		return fmt.Errorf("task %s not found in index.json", state.TaskID)
	}

	if foundTask.Status != "completed" {
		return fmt.Errorf("task %s status is %q, not completed. Run 'task record' before committing",
			state.TaskID, foundTask.Status)
	}

	// Verify record file exists
	if foundTask.Record != "" {
		recordPath := filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, foundTask.Record))
		if _, err := os.Stat(recordPath); os.IsNotExist(err) {
			return fmt.Errorf("record file missing: %s. Create it before committing",
				feature.GetTaskFile(featureSlug, foundTask.Record))
		}
	}

	return nil
}
