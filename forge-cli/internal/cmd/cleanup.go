package cmd

import (
	"os"
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"
	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up task state for completed, blocked, suspended, or rejected tasks",
	Long: `Removes state.json and record.json when the current task has a terminal
or inactive status: completed, blocked, suspended, or rejected.

Called as a Stop hook to clean up:
  - docs/features/<feature>/tasks/process/state.json
  - docs/features/<feature>/tasks/process/record.json (if exists)

Note: .forge/state.json is NOT deleted here — it is only cleared by quality-gate.`,
	Args: cobra.NoArgs,
	RunE: runCleanup,
}

func runCleanup(_ *cobra.Command, _ []string) error {
	cleanupCompletedTaskState()
	return nil
}

// cleanupCompletedTaskState removes state.json when task is completed.
func cleanupCompletedTaskState() {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return // No project context, nothing to clean up
	}

	// Note: .forge/state.json is NOT deleted here — it serves as a workspace marker
	// and is only deleted by `forge quality-gate` after successful execution.

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
	if t.Status == types.StatusCompleted || t.Status == types.StatusBlocked || t.Status == types.StatusSuspended || t.Status == types.StatusRejected {
		_ = os.Remove(statePath)

		// Also delete record.json if exists
		recordPath := feature.GetProcessRecordPath(projectRoot, featureSlug)
		_ = os.Remove(recordPath)
	}
}
