package task

import (
	"fmt"
	"forge-cli/internal/cmd/base"
	"os"
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate index.json: infer type fields and migrate legacy scope to surface-key/surface-type",
	Long: `Migrate the current feature's index.json by:
1. Inferring the type field for every task using the task ID pattern rules.
2. Migrating legacy 'scope' fields to 'surface-key' and 'surface-type' by
   calling 'forge surfaces --json <path>' to resolve the correct values.

Migration is idempotent: tasks that already have the correct fields are left unchanged.

Aborts with exit code 1 if any task is in_progress, leaving index.json unchanged.`,
	Args: cobra.NoArgs,
	RunE: runMigrate,
}

func runMigrate(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		base.Exit(base.ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	var taskCount int
	var scopeMigrated int
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

		// Phase 1: Infer and set type for every task.
		tasks := index.TasksMap()
		taskCount = len(tasks)
		for key, t := range tasks {
			inferred := inferTypeForTask(t)
			t.Type = inferred
			index.SetTask(key, t)
		}

		// Phase 2: Migrate legacy scope -> surface-key + surface-type
		scopeMigrated = migrateScopeToSurface(index, tasksDir, projectRoot)

		return indexPkg.SaveIndexAtomic(indexPath, index)
	}); err != nil {
		base.Exit(fmt.Errorf("migrate: %w", err))
	}

	fmt.Printf("Migrated %d tasks (type inference), %d tasks (scope->surface). Run task validate to verify.\n", taskCount, scopeMigrated)
	return nil
}

// inferTypeForTask returns the inferred type for a task, using the task's
// existing type if already set and non-empty, otherwise using InferType.
func inferTypeForTask(t task.Task) string {
	inferred := task.InferType(t.ID)
	if inferred == "" {
		inferred = task.TypeCodingFeature
	}
	return inferred
}

// migrateScopeToSurface scans for tasks with a legacy 'scope' field but no
// 'surface-key', resolves their surface via forgeconfig, and updates both
// index.json entries and the corresponding frontmatter .md files.
func migrateScopeToSurface(index *task.TaskIndex, tasksDir string, projectRoot string) int {
	surfaces, _ := forgeconfig.ReadSurfaces(projectRoot)

	var count int
	for key, t := range index.TasksMap() {
		if t.Scope == "" || t.SurfaceKey != "" {
			continue
		}

		// Resolve surface via config using the task's file path as the query
		match, err := forgeconfig.MatchSurface(surfaces, t.File)
		if err != nil {
			// If no surfaces configured, log warning and skip
			fmt.Fprintf(os.Stderr, "WARNING: task %s: could not resolve surface for scope %q: %v\n", t.ID, t.Scope, err)
			continue
		}

		t.SurfaceKey = match.Key
		t.SurfaceType = match.Type
		t.Scope = "" // Clear legacy field
		index.SetTask(key, t)
		count++

		// Update the frontmatter .md file
		if t.File != "" {
			mdPath := filepath.Join(tasksDir, t.File)
			if err := updateFrontmatterSurface(mdPath, match.Key, match.Type); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: task %s: failed to update frontmatter: %v\n", t.ID, err)
			}
		}
	}
	return count
}

// updateFrontmatterSurface updates the frontmatter of a task .md file to
// replace the legacy scope field with surface-key and surface-type.
func updateFrontmatterSurface(mdPath string, surfaceKey, surfaceType string) error {
	content, err := os.ReadFile(mdPath)
	if err != nil {
		return err
	}

	fm, body, err := task.ParseFrontmatter(content)
	if err != nil {
		return err
	}

	fm.SurfaceKey = surfaceKey
	fm.SurfaceType = surfaceType
	fm.Scope = ""

	return task.WriteFrontmatter(mdPath, fm, body)
}
