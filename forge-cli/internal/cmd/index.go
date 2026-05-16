package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/profile"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var (
	indexFeatureSlug  string
	indexTestProfiles []string
)

var indexCmd = &cobra.Command{
	Use:   "index --feature <slug> [--test-profiles p1,p2]",
	Short: "Build or rebuild index.json from task markdown files",
	Long: `Scan .md files in the feature's tasks/ directory and generate/update index.json.
Idempotent: re-running with no changes produces the same output.

Test tasks are auto-generated from embedded profiles.
Profiles are read from .forge/config.yaml unless overridden by --test-profiles.`,
	Run: runIndex,
}

func init() {
	indexCmd.Flags().StringVar(&indexFeatureSlug, "feature", "", "Feature slug (required)")
	indexCmd.Flags().StringSliceVar(&indexTestProfiles, "test-profiles", nil, "Override test profiles (comma-separated)")
	_ = indexCmd.MarkFlagRequired("feature")
}

func runIndex(_ *cobra.Command, _ []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	// Validate feature dir exists
	featureDir := filepath.Join(projectRoot, feature.GetFeatureDir(indexFeatureSlug))
	if _, err := os.Stat(featureDir); os.IsNotExist(err) {
		Exit(ErrFeatureNotFound(indexFeatureSlug))
	}

	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(indexFeatureSlug))
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(indexFeatureSlug))

	// Ensure tasks dir exists
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		Exit(fmt.Errorf("create tasks dir: %w", err))
	}

	// Resolve profiles
	profiles := indexTestProfiles
	if len(profiles) == 0 {
		p, err := profile.ReadTestProfiles(projectRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to read profiles: %v\n", err)
		}
		profiles = p
	}

	// Resolve capabilities: config.yaml > UnionCapabilities(profiles)
	var capabilities []string
	cfg, _ := profile.ReadConfig(projectRoot)
	if cfg != nil && len(cfg.Capabilities) > 0 {
		capabilities = cfg.Capabilities
	}
	if len(capabilities) == 0 && len(profiles) > 0 {
		caps, err := profile.UnionCapabilities(profiles)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to resolve capabilities: %v\n", err)
		}
		capabilities = caps
	}

	// Build strategy resolver
	resolveStrategy := func(profileName, kind string) []byte {
		content, err := profile.GetStrategy(profileName, kind)
		if err != nil {
			return nil
		}
		return content
	}

	opts := task.BuildIndexOpts{
		FeatureSlug:      indexFeatureSlug,
		ProjectRoot:      projectRoot,
		TasksDir:         tasksDir,
		IndexPath:        indexPath,
		TestProfiles:     profiles,
		TestCapabilities: capabilities,
		ResolveStrategy:  resolveStrategy,
	}

	result, err := task.BuildIndex(opts)
	if err != nil {
		Exit(fmt.Errorf("build index: %w", err))
	}

	// Print summary
	PrintBlockStart()
	PrintField("ACTION", "INDEX_BUILT")
	PrintField("FEATURE", indexFeatureSlug)
	PrintField("INDEX", indexPath)
	PrintField("NEW", fmt.Sprintf("%d", result.NewCount))
	PrintField("UPDATED", fmt.Sprintf("%d", result.UpdatedCount))
	PrintField("PRESERVED", fmt.Sprintf("%d", result.PreservedCount))
	if len(profiles) > 0 {
		PrintField("TEST_PROFILES", strings.Join(profiles, ", "))
	}
	PrintBlockEnd()

	// Print warnings
	for _, w := range result.Warnings {
		PrintWarning(w)
	}

	// Run validation on the generated index
	v := &validator{filePath: indexPath}
	if err := v.run(); err != nil {
		// Validation errors are printed by v.run(), don't Exit here
		// The index was built successfully, validation is advisory
		fmt.Fprintf(os.Stderr, "NOTE: fix validation errors above and re-run 'forge task index --feature %s'\n", indexFeatureSlug)
	}
}
