package cmd

import (
	"os"
	"path/filepath"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"
	"task-cli/pkg/task"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up completed task state",
	Long: `Removes state.json and record.json when task is completed.

Called as a Stop hook to clean up:
  - docs/features/<feature>/tasks/process/state.json
  - docs/features/<feature>/tasks/process/record.json (if exists)`,
	Run: runCleanup,
}

func runCleanup(cmd *cobra.Command, args []string) {
	cleanupCompletedTaskState()
	os.Exit(0)
}

// cleanupCompletedTaskState removes state.json when task is completed,
// and cleans up .forge/state.json as a fallback.
func cleanupCompletedTaskState() {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return // No project context, nothing to clean up
	}

	// Always clean up .forge/state.json as fallback
	// (normally consumed by all-completed, but safe to clear on session end)
	feature.ClearForgeState(projectRoot)

	// Get current feature
	featureSlug, err := feature.GetCurrentFeature(projectRoot)
	if err != nil {
		return
	}

	statePath := feature.GetTaskStatePath(projectRoot, featureSlug)
	state, err := task.LoadState(statePath)
	if err != nil || state == nil || state.Key == "" {
		return
	}

	// Load index to check task status
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return
	}

	t, exists := index.Tasks[state.Key]
	if !exists {
		return
	}

	// Delete state file if task is completed
	if t.Status == "completed" {
		os.Remove(statePath)

		// Also delete record.json if exists
		recordPath := feature.GetProcessRecordPath(projectRoot, featureSlug)
		os.Remove(recordPath)
	}
}
