package cmd

import (
	"os"
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

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

func runCleanup(_ *cobra.Command, _ []string) {
	cleanupCompletedTaskState()
	os.Exit(0)
}

// cleanupCompletedTaskState removes state.json when task is completed.
func cleanupCompletedTaskState() {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return // No project context, nothing to clean up
	}

	// Note: .forge/state.json is NOT deleted here — it serves as a workspace marker
	// and is only deleted by `task all-completed` after successful execution.

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

	t, exists := index.ByID(state.Key)
	if !exists {
		return
	}

	// Delete state file if task is completed, blocked, or rejected
	if t.Status == "completed" || t.Status == "blocked" || t.Status == "rejected" {
		_ = os.Remove(statePath)

		// Also delete record.json if exists
		recordPath := feature.GetProcessRecordPath(projectRoot, featureSlug)
		_ = os.Remove(recordPath)
	}
}
