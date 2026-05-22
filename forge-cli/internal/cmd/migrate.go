package cmd

import (
	"fmt"
	"path/filepath"

	"forge-cli/pkg/feature"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/project"
	"forge-cli/pkg/prompt"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate index.json by inferring type fields for all tasks",
	Long: `Migrate the current feature's index.json by inferring the type field
for every task using the task ID pattern rules.

Migration is idempotent: tasks that already have a type field are overwritten
with the inferred value.

Aborts with exit code 1 if any task is in_progress, leaving index.json unchanged.`,
	Run: runMigrate,
}

func runMigrate(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		Exit(ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	taskCount := 0
	if err := indexPkg.WithLock(indexPath, func() error {
		index, err := task.LoadIndex(indexPath)
		if err != nil {
			return fmt.Errorf("load index: %w", err)
		}

		// Pre-flight check: abort if any task is in_progress.
		for _, t := range index.TasksMap() {
			if t.Status == feature.StatusInProgress {
				return fmt.Errorf("task %q is in_progress — finish or pause it before migrating", t.ID)
			}
		}

		// Infer and set type for every task.
		tasks := index.TasksMap()
		taskCount = len(tasks)
		for key, t := range tasks {
			inferred := prompt.InferType(t.ID)
			if inferred == "" {
				inferred = task.TypeCodingFeature
			}
			t.Type = inferred
			index.SetTask(key, t)
		}

		return indexPkg.SaveIndexAtomic(indexPath, index)
	}); err != nil {
		Exit(fmt.Errorf("migrate: %w", err))
	}

	fmt.Printf("Migrated %d tasks. Run task validate to verify.\n", taskCount)
}
