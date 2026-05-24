package task

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"forge-cli/internal/cmd/base"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var listLocal bool

var listCmd = &cobra.Command{
	Use:   "list [slug]",
	Short: "List all tasks for the current feature",
	Long: `List all tasks for the current feature in a table format.

Displays task ID, type, title (truncated), status, breaking, and mainSession.
Tasks are sorted by ID in natural order: numeric IDs first, then test/gate IDs.

When a slug is provided, lists tasks for that specific feature, reading from
the worktree if one exists for that slug. Use --local to read from the main
repository's index.json regardless of worktree existence.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVar(&listLocal, "local", false, "Read from main repo's index.json (ignore worktree)")
}

const titleMaxWidth = 50

func runList(_ *cobra.Command, args []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	var featureSlug string
	var indexPath string

	if len(args) == 1 {
		// Slug provided: bypass RequireFeature, construct index path directly
		featureSlug = args[0]

		// Validate: feature directory must exist
		featureDir := filepath.Join(projectRoot, feature.GetFeatureDir(featureSlug))
		if _, err := os.Stat(featureDir); os.IsNotExist(err) {
			return base.ErrFeatureNotFound(featureSlug)
		}

		indexPath = resolveListIndexPath(projectRoot, featureSlug)
	} else {
		// No slug: use existing auto-detection logic
		featureSlug, err = feature.RequireFeature(projectRoot)
		if err != nil {
			base.Exit(base.ErrFeatureNotSet())
		}

		indexPath = filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	}

	index, err := task.LoadIndex(indexPath)
	if err != nil || index.TaskCount() == 0 {
		fmt.Printf("no tasks found  (feature: %s)\n", featureSlug)
		return nil
	}

	// Collect and sort task IDs
	tasks := index.TasksMap()
	ids := make([]string, 0, len(tasks))
	for id := range tasks {
		ids = append(ids, id)
	}
	sortedIDs := naturalSortTaskIDs(ids)

	// Print header
	fmt.Printf("%d found  (feature: %s)\n\n", len(sortedIDs), featureSlug)

	// Calculate dynamic column widths based on display width (CJK-aware)
	idCol := base.DisplayWidth("ID")
	typeCol := base.DisplayWidth("TYPE")
	titleCol := base.DisplayWidth("TITLE")
	statusCol := base.DisplayWidth("STATUS")
	breakingCol := base.DisplayWidth("BREAKING")
	mainSessCol := base.DisplayWidth("MAIN_SESS")

	boolStr := func(v bool) string {
		if v {
			return "true"
		}
		return "false"
	}

	for _, id := range sortedIDs {
		t := tasks[id]
		if base.DisplayWidth(t.ID) > idCol {
			idCol = base.DisplayWidth(t.ID)
		}
		if base.DisplayWidth(t.Type) > typeCol {
			typeCol = base.DisplayWidth(t.Type)
		}
		titleLen := base.DisplayWidth(t.Title)
		if titleLen > titleMaxWidth {
			titleLen = titleMaxWidth
		}
		if titleLen > titleCol {
			titleCol = titleLen
		}
		if base.DisplayWidth(t.Status) > statusCol {
			statusCol = base.DisplayWidth(t.Status)
		}
		if base.DisplayWidth(boolStr(t.Breaking)) > breakingCol {
			breakingCol = base.DisplayWidth(boolStr(t.Breaking))
		}
		if base.DisplayWidth(boolStr(t.MainSession)) > mainSessCol {
			mainSessCol = base.DisplayWidth(boolStr(t.MainSession))
		}
	}

	// Print column headers
	fmt.Printf("%s  %s  %s  %s  %s  %s\n",
		base.PadRight("ID", idCol),
		base.PadRight("TYPE", typeCol),
		base.PadRight("TITLE", titleCol),
		base.PadRight("STATUS", statusCol),
		base.PadRight("BREAKING", breakingCol),
		base.PadRight("MAIN_SESS", mainSessCol),
	)

	// Print separator
	sep := fmt.Sprintf("%s  %s  %s  %s  %s  %s",
		strings.Repeat("-", idCol),
		strings.Repeat("-", typeCol),
		strings.Repeat("-", titleCol),
		strings.Repeat("-", statusCol),
		strings.Repeat("-", breakingCol),
		strings.Repeat("-", mainSessCol),
	)
	fmt.Println(sep)

	// Print task rows
	for _, id := range sortedIDs {
		t := tasks[id]
		title := base.TruncateSlug(t.Title, titleCol)
		fmt.Printf("%s  %s  %s  %s  %s  %s\n",
			base.PadRight(t.ID, idCol),
			base.PadRight(t.Type, typeCol),
			base.PadRight(title, titleCol),
			base.PadRight(t.Status, statusCol),
			base.PadRight(boolStr(t.Breaking), breakingCol),
			base.PadRight(boolStr(t.MainSession), mainSessCol),
		)
	}

	return nil
}

// resolveListIndexPath determines the index.json path for a feature slug.
// If --local is not set and a worktree exists for the slug, reads from the
// worktree's copy of index.json; otherwise reads from the main repo.
func resolveListIndexPath(projectRoot, slug string) string {
	if !listLocal {
		// Check if .forge/worktrees/<slug> directory exists (fast path)
		worktreeDir := filepath.Join(projectRoot, ".forge", "worktrees", slug)
		if info, err := os.Stat(worktreeDir); err == nil && info.IsDir() {
			wtIndex := filepath.Join(worktreeDir, feature.GetFeatureIndexFile(slug))
			if _, err := os.Stat(wtIndex); err == nil {
				return wtIndex
			}
		}
	}

	// Fallback: main repo's index.json
	return filepath.Join(projectRoot, feature.GetFeatureIndexFile(slug))
}

// naturalSortTaskIDs sorts task IDs in natural order:
// business IDs grouped by numeric prefix (1, 1.gate, 2, 2.summary, ...),
// then test pipeline IDs (T-1, T-2, ...).
// Within the same numeric prefix, pure numeric comes before compound.
func naturalSortTaskIDs(ids []string) []string {
	sorted := make([]string, len(ids))
	copy(sorted, ids)

	sort.SliceStable(sorted, func(i, j int) bool {
		ki := sortKey(sorted[i])
		kj := sortKey(sorted[j])

		// Primary: T-prefixed IDs sort after all business IDs
		if ki.isTestPipeline != kj.isTestPipeline {
			return !ki.isTestPipeline
		}

		// Secondary: numeric prefix (groups 1, 1.gate together before 2)
		if ki.numPrefix != kj.numPrefix {
			return ki.numPrefix < kj.numPrefix
		}

		// Tertiary: same prefix — pure numeric before compound (1 < 1.gate)
		if ki.isPureNumeric != kj.isPureNumeric {
			return ki.isPureNumeric
		}

		// Quaternary: string comparison for compound suffixes
		return sorted[i] < sorted[j]
	})

	return sorted
}

type idSortKey struct {
	isTestPipeline bool
	isPureNumeric  bool
	numPrefix      int
}

func sortKey(id string) idSortKey {
	// Test pipeline IDs: T-1, T-2, etc.
	if numStr, ok := strings.CutPrefix(id, task.IDPrefixTestPipeline); ok {
		num, _ := strconv.Atoi(numStr)
		return idSortKey{isTestPipeline: true, numPrefix: num}
	}

	// Try pure numeric: "1", "2", "10"
	if num, err := strconv.Atoi(id); err == nil {
		return idSortKey{isPureNumeric: true, numPrefix: num}
	}

	// Compound IDs: "1.gate", "1.summary", "1.1" etc.
	dotIdx := strings.Index(id, ".")
	if dotIdx > 0 {
		prefix := id[:dotIdx]
		if num, err := strconv.Atoi(prefix); err == nil {
			return idSortKey{numPrefix: num}
		}
	}

	// Fallback: high numPrefix to sort last among business IDs
	return idSortKey{numPrefix: 99999}
}
