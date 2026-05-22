package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var indexFeatureSlug string

var indexCmd = &cobra.Command{
	Use:   "index --feature <slug>",
	Short: "Build or rebuild index.json from task markdown files",
	Long: `Scan .md files in the feature's tasks/ directory and generate/update index.json.
Idempotent: re-running with no changes produces the same output.`,
	Args: cobra.NoArgs,
	RunE: runIndex,
}

func init() {
	indexCmd.Flags().StringVar(&indexFeatureSlug, "feature", "", "Feature slug (required)")
	_ = indexCmd.MarkFlagRequired("feature")
}

func runIndex(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	// Validate feature dir exists
	featureDir := filepath.Join(projectRoot, feature.GetFeatureDir(indexFeatureSlug))
	if _, err := os.Stat(featureDir); os.IsNotExist(err) {
		return ErrFeatureNotFound(indexFeatureSlug)
	}

	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(indexFeatureSlug))
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(indexFeatureSlug))

	// Ensure tasks dir exists
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return fmt.Errorf("create tasks dir: %w", err)
	}

	// Read auto-behavior config (returns defaults when missing)
	auto, err := forgeconfig.ReadAutoConfig(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to read auto config: %v\n", err)
	}

	opts := task.BuildIndexOpts{
		FeatureSlug: indexFeatureSlug,
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
		AutoConfig:  auto,
	}

	result, err := task.BuildIndex(opts)
	if err != nil {
		return fmt.Errorf("build index: %w", err)
	}

	// Save index atomically under lock (BuildIndex already saves internally,
	// but the lock ensures exclusive access during the full rebuild cycle).
	if err := indexPkg.WithLock(indexPath, func() error {
		return indexPkg.SaveIndexAtomic(indexPath, result.Index)
	}); err != nil {
		return fmt.Errorf("save index: %w", err)
	}

	// Print summary
	PrintBlockStart()
	PrintField("ACTION", "INDEX_BUILT")
	PrintField("FEATURE", indexFeatureSlug)
	PrintField("INDEX", indexPath)
	PrintField("NEW", fmt.Sprintf("%d", result.NewCount))
	PrintField("UPDATED", fmt.Sprintf("%d", result.UpdatedCount))
	PrintField("PRESERVED", fmt.Sprintf("%d", result.PreservedCount))
	PrintBlockEnd()

	// Print warnings
	for _, w := range result.Warnings {
		PrintWarning(w)
	}

	// Run validation on the generated index
	v := &validator{filePath: indexPath}
	if err := v.run(); err != nil {
		fmt.Fprintf(os.Stderr, "NOTE: fix validation errors above and re-run 'forge task index --feature %s'\n", indexFeatureSlug)
	}
	return nil
}
