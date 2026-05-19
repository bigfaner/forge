package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
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

	// Resolve languages: flag > config.yaml > auto-detect
	languages := indexLanguages
	if len(languages) == 0 {
		langs, err := forgeconfig.ReadLanguages(projectRoot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to read languages: %v\n", err)
		}
		languages = langs
	}

	// Resolve interfaces: config.yaml > union of language capabilities
	var interfaces []string
	cfg, _ := forgeconfig.ReadConfig(projectRoot)
	if cfg != nil && len(cfg.Interfaces) > 0 {
		interfaces = cfg.Interfaces
	}
	if len(interfaces) == 0 && len(languages) > 0 {
		ifaces, err := forgeconfig.UnionLanguageInterfaces(languages)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to resolve interfaces: %v\n", err)
		}
		interfaces = ifaces
	}

	opts := task.BuildIndexOpts{
		FeatureSlug: indexFeatureSlug,
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
	}

	// Read auto-behavior config (returns defaults when missing)
	auto, err := forgeconfig.ReadAutoConfig(projectRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to read auto config: %v\n", err)
	}
	opts.AutoConfig = auto

	result, err := task.BuildIndex(opts)
	if err != nil {
		Exit(fmt.Errorf("build index: %w", err))
	}

	// Generate test pipeline tasks (caller-side assembly)
	if len(languages) > 0 && len(interfaces) > 0 {
		mode := detectIndexMode(projectRoot, indexFeatureSlug)
		if mode != "" {
			testTasks := task.GenerateTestTasks(mode, languages, interfaces, auto)
			if len(testTasks) > 0 {
				task.ResolveFirstTestDep(testTasks, result.Index.TasksMap(), mode)
				generateTestTaskFiles(testTasks, tasksDir, result)
			}
		}
	}

	// Re-save index after test task generation
	if err := task.SaveIndex(indexPath, result.Index); err != nil {
		Exit(fmt.Errorf("save index: %w", err))
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
		fmt.Fprintf(os.Stderr, "NOTE: fix validation errors above and re-run 'forge task index --feature %s'\n", indexFeatureSlug)
	}
}

func detectIndexMode(projectRoot, slug string) string {
	featureDir := filepath.Join(projectRoot, "docs", "features", slug)
	if _, err := os.Stat(filepath.Join(featureDir, "prd", "prd-spec.md")); err == nil {
		return "breakdown"
	}
	if _, err := os.Stat(filepath.Join(projectRoot, "docs", "proposals", slug, "proposal.md")); err == nil {
		return "quick"
	}
	return ""
}

func generateTestTaskFiles(testTasks []task.TestTaskDef, tasksDir string, result *task.BuildIndexResult) {
	for i := range testTasks {
		td := &testTasks[i]
		key := td.Key

		// Generate .md if missing
		mdPath := filepath.Join(tasksDir, key+".md")
		if _, err := os.Stat(mdPath); os.IsNotExist(err) {
			content, err := task.GenerateTestTaskMD(*td, "")
			if err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("generate %s: %v", key, err))
				continue
			}
			if err := os.WriteFile(mdPath, content, 0644); err != nil {
				result.Warnings = append(result.Warnings, fmt.Sprintf("write %s: %v", key, err))
				continue
			}
		}

		t := td.TaskFromFile()

		if existing, found := result.Index.ByID(td.ID); found {
			t.Status = existing.Status
			t.SourceTaskID = existing.SourceTaskID
			t.BlockedReason = existing.BlockedReason
			result.Index.SetTask(key, t)
			result.UpdatedCount++
		} else {
			result.Index.SetTask(key, t)
			result.NewCount++
		}
	}
}
