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
	indexFeatureSlug string
	indexLanguages   []string
)

var indexCmd = &cobra.Command{
	Use:   "index --feature <slug> [--languages go,javascript]",
	Short: "Build or rebuild index.json from task markdown files",
	Long: `Scan .md files in the feature's tasks/ directory and generate/update index.json.
Idempotent: re-running with no changes produces the same output.

Test tasks are auto-generated from detected language strategies.
Languages are read from .forge/config.yaml unless overridden by --languages.`,
	Run: runIndex,
}

func init() {
	indexCmd.Flags().StringVar(&indexFeatureSlug, "feature", "", "Feature slug (required)")
	indexCmd.Flags().StringSliceVar(&indexLanguages, "languages", nil, "Override detected languages (comma-separated)")
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

	// Resolve languages
	languages := indexLanguages
	if len(languages) == 0 {
		langs, err := profile.ReadLanguages(projectRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to read languages: %v\n", err)
		}
		languages = langs
	}

	// Resolve interfaces: config.yaml > UnionLanguageInterfaces(languages)
	var interfaces []string
	cfg, _ := profile.ReadConfig(projectRoot)
	if cfg != nil && len(cfg.Interfaces) > 0 {
		interfaces = cfg.Interfaces
	}
	if len(interfaces) == 0 && len(languages) > 0 {
		ifaces, err := profile.UnionLanguageInterfaces(languages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to resolve interfaces: %v\n", err)
		}
		interfaces = ifaces
	}

	// Build strategy resolver
	resolveStrategy := func(language, kind string) []byte {
		content, err := profile.GetStrategy(language, kind)
		if err != nil {
			return nil
		}
		return content
	}

	opts := task.BuildIndexOpts{
		FeatureSlug:     indexFeatureSlug,
		ProjectRoot:     projectRoot,
		TasksDir:        tasksDir,
		IndexPath:       indexPath,
		Languages:       languages,
		TestInterfaces:  interfaces,
		ResolveStrategy: resolveStrategy,
	}

	// Read auto-behavior config (returns defaults when missing)
	auto, err := profile.ReadAutoConfig(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to read auto config: %v\n", err)
	}
	opts.AutoConfig = auto

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
	if len(languages) > 0 {
		PrintField("LANGUAGES", strings.Join(languages, ", "))
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
