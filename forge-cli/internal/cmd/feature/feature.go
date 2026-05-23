// Package feature contains all forge feature subcommand implementations.
//
// Commands are registered into the CLI tree via Register(), called from
// the parent cmd package during initialization.
package feature

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"forge-cli/internal/cmd/base"
	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var verbose bool

// Cmd is the parent feature command, exported for use by the cmd package.
var Cmd = &cobra.Command{
	Use:   "feature [slug]",
	Short: "Set or display the current feature",
	Long: `Set or display the current feature context.

Without arguments: displays the current feature.
With a slug argument: sets the current feature.

Subcommands:
  list            List all features
  status <slug>   Show feature status detail`,
	Args: cobra.MaximumNArgs(1),
	RunE: runFeature,
}

var featureListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all features",
	Long:  "List all features with status, progress, and scores.",
	Args:  cobra.NoArgs,
	RunE:  runFeatureList,
}

var featureStatusCmd = &cobra.Command{
	Use:   "status <slug>",
	Short: "Show feature status detail",
	Long:  "Show detailed status for a feature including manifest, task counts, and artifact scores.",
	Args:  cobra.ExactArgs(1),
	RunE:  runFeatureStatus,
}

var featureSetCmd = &cobra.Command{
	Use:   "set <slug>",
	Short: "Explicitly set the current feature",
	Long: `Set the current feature context by writing to .forge/state.json
and ensuring the feature directory structure exists.

This provides an explicit override for feature resolution,
complementing the existing implicit resolution from git context.`,
	Args: exactArgsNonEmpty(1),
	RunE: runFeatureSet,
}

// Register adds all feature subcommands to Cmd.
func Register() {
	Cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show resolution source")
	Cmd.AddCommand(featureListCmd)
	Cmd.AddCommand(featureStatusCmd)
	Cmd.AddCommand(featureSetCmd)
}

// exactArgsNonEmpty returns a cobra.Args validator that requires exactly n arguments,
// each non-empty. Returns an error for empty strings instead of calling os.Exit.
func exactArgsNonEmpty(n int) cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) != n {
			return fmt.Errorf("requires exactly %d arg(s), got %d", n, len(args))
		}
		for i, arg := range args {
			if strings.TrimSpace(arg) == "" {
				return fmt.Errorf("argument %d must not be empty", i+1)
			}
		}
		return nil
	}
}

func runFeature(_ *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	if len(args) == 0 {
		// Display current feature
		if verbose {
			slug, source, err := feature.GetCurrentFeatureWithSource(projectRoot)
			if err != nil {
				base.PrintBlockStart()
				base.PrintField("FEATURE", "(none)")
				base.PrintBlockEnd()
				return nil
			}
			base.PrintBlockStart()
			base.PrintField("FEATURE", fmt.Sprintf("%s (from: %s)", slug, source))
			base.PrintBlockEnd()
			return nil
		}

		slug, err := feature.GetCurrentFeature(projectRoot)
		if err != nil {
			base.PrintBlockStart()
			base.PrintField("FEATURE", "(none)")
			base.PrintBlockEnd()
			return nil
		}
		base.PrintBlockStart()
		base.PrintField("FEATURE", slug)
		base.PrintBlockEnd()
		return nil
	}

	// Set feature
	slug := args[0]
	if err := feature.EnsureFeatureDir(projectRoot, slug); err != nil {
		return base.ErrFeatureNotFound(slug)
	}
	base.PrintBlockStart()
	base.PrintField("FEATURE", slug)
	base.PrintBlockEnd()
	return nil
}

func runFeatureSet(_ *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	slug := args[0]
	if slug == "" {
		return base.ErrNoInput("feature slug is required")
	}

	if err := feature.EnsureFeatureDir(projectRoot, slug); err != nil {
		return base.NewAIError(
			base.ErrNotFound,
			fmt.Sprintf("Failed to create feature directory for: %s", slug),
			err.Error(),
			"Check filesystem permissions",
			"ls docs/features/",
		)
	}

	if err := feature.EnsureForgeState(projectRoot, slug); err != nil {
		return base.NewAIError(
			base.ErrNotFound,
			fmt.Sprintf("Failed to write state for feature: %s", slug),
			err.Error(),
			"Check .forge/ directory permissions",
			"ls -la .forge/",
		)
	}

	base.PrintBlockStart()
	base.PrintField("FEATURE", slug)
	base.PrintBlockEnd()
	return nil
}

// featureInfo holds information about a discovered feature.
type featureInfo struct {
	Slug          string
	Status        string
	PRDScore      string
	DesignScore   string
	UIScore       string
	TestScore     string
	Completed     int
	Total         int
	Created       string // created field from manifest frontmatter (YYYY-MM-DD); empty if missing
	ManifestMtime int64  // unix seconds of manifest.md mod time; 0 if missing/unreadable
}

func runFeatureList(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	features, err := discoverFeatures(projectRoot)
	if err != nil {
		return newErrFeatureDiscovery(err)
	}

	// Sort by created date descending (newest first).
	// Created is stored as "YYYY-MM-DD" which sorts correctly lexicographically.
	// Features without created field fall back to manifest mtime.
	sort.Slice(features, func(i, j int) bool {
		ci, cj := features[i].Created, features[j].Created
		// Items with created field sort before items without (descending).
		if ci != "" && cj != "" {
			return ci > cj
		}
		// If only one has created, it sorts first.
		if ci != "" {
			return true
		}
		if cj != "" {
			return false
		}
		// Both missing created: fall back to mtime descending.
		return features[i].ManifestMtime > features[j].ManifestMtime
	})

	if len(features) == 0 {
		fmt.Fprintln(os.Stderr, "no features found")
		return nil
	}

	// Calculate dynamic slug column width.
	slugWidth := base.CalcSlugColWidth(mapFeaturesToSlugLens(features))

	base.PrintBlockStart()
	base.PrintField("FEATURES", fmt.Sprintf("%d found", len(features)))
	fmt.Println()

	// Table header
	fmt.Printf("  %-s %-12s %-10s %-10s %-10s %-10s %-10s\n",
		base.PadRight("SLUG", slugWidth), "STATUS", "PROGRESS", "PRD", "DESIGN", "UI", "TESTS")
	fmt.Printf("  %-s %-12s %-10s %-10s %-10s %-10s %-10s\n",
		strings.Repeat("-", slugWidth),
		strings.Repeat("-", 10),
		strings.Repeat("-", 8),
		strings.Repeat("-", 5),
		strings.Repeat("-", 5),
		strings.Repeat("-", 5),
		strings.Repeat("-", 5))

	for _, f := range features {
		progress := fmt.Sprintf("%d/%d", f.Completed, f.Total)
		fmt.Printf("  %-s %-12s %-10s %-10s %-10s %-10s %-10s\n",
			base.PadRight(base.TruncateSlug(f.Slug, slugWidth), slugWidth),
			f.Status,
			progress,
			scoreDisplay(f.PRDScore),
			scoreDisplay(f.DesignScore),
			scoreDisplay(f.UIScore),
			scoreDisplay(f.TestScore))
	}

	fmt.Println()
	base.PrintBlockEnd()
	return nil
}

func runFeatureStatus(_ *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	slug := args[0]
	featureDir := filepath.Join(projectRoot, feature.FeaturesDir, slug)
	if _, err := os.Stat(featureDir); os.IsNotExist(err) {
		return base.ErrFeatureNotFound(slug)
	}

	// Read manifest
	manifestPath := filepath.Join(featureDir, feature.ManifestFileName)
	manifestStatus := ""
	if data, err := os.ReadFile(manifestPath); err == nil {
		var meta struct {
			Status string `yaml:"status"`
		}
		if err := parseYAMLFrontmatter(data, &meta); err == nil {
			manifestStatus = meta.Status
		}
	}

	// Read task index
	indexPath := filepath.Join(featureDir, feature.TasksDirName, feature.IndexFileName)
	var taskStats map[string]int
	total := 0
	if data, err := os.ReadFile(indexPath); err == nil {
		var idx task.TaskIndex
		if err := json.Unmarshal(data, &idx); err == nil {
			taskStats = make(map[string]int)
			for _, t := range idx.TasksMap() {
				taskStats[t.Status]++
				total++
			}
		}
	}

	// Read scores
	prdScore := readScoreFromFrontmatter(filepath.Join(featureDir, feature.PRDDirName, feature.PRDSpecFile))
	designScore := readScoreFromFrontmatter(filepath.Join(featureDir, feature.DesignDirName, feature.TechDesignFile))
	uiScore := readScoreFromFrontmatter(filepath.Join(featureDir, feature.UIDirName, feature.UIDesignFile))

	base.PrintBlockStart()
	base.PrintField("SLUG", slug)
	base.PrintField("STATUS", manifestStatus)
	base.PrintFieldIfNotEmpty("FILE", filepath.Join(feature.FeaturesDir, slug, feature.ManifestFileName))
	fmt.Println()

	base.PrintSection("TASKS")
	if taskStats != nil {
		for _, status := range []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"} {
			if count, ok := taskStats[status]; ok {
				base.PrintListItem(fmt.Sprintf("%s: %d", status, count))
			}
		}
		base.PrintField("TOTAL", fmt.Sprintf("%d", total))
	} else {
		base.PrintField("TOTAL", "0 (no index.json)")
	}

	fmt.Println()
	base.PrintSection("ARTIFACTS")
	base.PrintField("PRD", scoreDisplay(prdScore))
	base.PrintField("DESIGN", scoreDisplay(designScore))
	base.PrintField("UI", scoreDisplay(uiScore))

	base.PrintBlockEnd()
	return nil
}

// discoverFeatures walks docs/features/*/manifest.md and collects feature info.
func discoverFeatures(projectRoot string) ([]featureInfo, error) {
	featuresDir := filepath.Join(projectRoot, feature.FeaturesDir)
	entries, err := os.ReadDir(featuresDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read features directory: %w", err)
	}

	var features []featureInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		slug := entry.Name()
		featureDir := filepath.Join(featuresDir, slug)

		// Read manifest status, created, and mtime
		status := ""
		created := ""
		var manifestMtime int64
		manifestPath := filepath.Join(featureDir, feature.ManifestFileName)
		if data, err := os.ReadFile(manifestPath); err == nil {
			var meta struct {
				Status  string `yaml:"status"`
				Created string `yaml:"created"`
			}
			if err := parseYAMLFrontmatter(data, &meta); err == nil {
				status = meta.Status
				created = meta.Created
			}
		}
		if info, err := os.Stat(manifestPath); err == nil {
			manifestMtime = info.ModTime().Unix()
		}

		// Read task progress
		completed, total := readTaskProgress(filepath.Join(featureDir, feature.TasksDirName, feature.IndexFileName))

		// Read scores
		prdScore := readScoreFromFrontmatter(filepath.Join(featureDir, feature.PRDDirName, feature.PRDSpecFile))
		designScore := readScoreFromFrontmatter(filepath.Join(featureDir, feature.DesignDirName, feature.TechDesignFile))
		uiScore := readScoreFromFrontmatter(filepath.Join(featureDir, feature.UIDirName, feature.UIDesignFile))
		testScore := readScoreFromFrontmatter(filepath.Join(featureDir, feature.TestingResultsDirName, "results.json"))

		features = append(features, featureInfo{
			Slug:          slug,
			Status:        status,
			PRDScore:      prdScore,
			DesignScore:   designScore,
			UIScore:       uiScore,
			TestScore:     testScore,
			Completed:     completed,
			Total:         total,
			ManifestMtime: manifestMtime,
			Created:       created,
		})
	}

	return features, nil
}

// readTaskProgress reads index.json and returns completed/total counts.
func readTaskProgress(indexPath string) (completed, total int) {
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return 0, 0
	}

	var idx task.TaskIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return 0, 0
	}

	for _, t := range idx.TasksMap() {
		total++
		if t.Status == "completed" {
			completed++
		}
	}
	return completed, total
}

// readScoreFromFrontmatter reads the score field from a file's YAML frontmatter.
func readScoreFromFrontmatter(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	var meta struct {
		Score string `yaml:"score"`
	}
	if err := parseYAMLFrontmatter(data, &meta); err != nil {
		return ""
	}
	return meta.Score
}

// scoreDisplay returns the score or em-dash when missing.
func scoreDisplay(score string) string {
	if score == "" {
		return "—" // em dash
	}
	return score
}

// parseYAMLFrontmatter extracts YAML frontmatter from markdown content.
func parseYAMLFrontmatter(content []byte, target any) error {
	text := string(content)

	if !strings.HasPrefix(text, "---") {
		return nil
	}
	text = text[3:]

	closeIdx := strings.Index(text, "\n---")
	if closeIdx < 0 {
		return nil
	}

	yamlContent := text[:closeIdx]
	return parseYAML([]byte(yamlContent), target)
}

// parseYAML is a thin wrapper around yaml.Unmarshal for testability.
var parseYAML = defaultParseYAML

func defaultParseYAML(data []byte, target any) error {
	return yaml.Unmarshal(data, target)
}

func newErrFeatureDiscovery(err error) *base.AIError {
	return base.NewAIError(
		base.ErrNotFound,
		"Failed to discover features",
		err.Error(),
		"Ensure docs/features/ directory exists",
		"ls docs/features/",
	)
}

// mapFeaturesToSlugLens extracts slug lengths from feature list.
func mapFeaturesToSlugLens(features []featureInfo) []int {
	lens := make([]int, len(features))
	for i, f := range features {
		lens[i] = len(f.Slug)
	}
	return lens
}
