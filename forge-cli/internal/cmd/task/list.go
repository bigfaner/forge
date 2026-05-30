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

// fallbackSortPriority is assigned to task IDs that cannot be parsed,
// ensuring they sort after all valid business IDs.
const fallbackSortPriority = 99999

// ANSI escape sequences for cycle marker display.
const (
	colorCycleMarker = "\033[33m"
	colorReset       = "\033[0m"
)

// listIsTerminalFunc detects whether stdout is a terminal (TTY).
// Overridable for testing.
var listIsTerminalFunc = defaultListIsTerminal

func defaultListIsTerminal() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

var listCmd = &cobra.Command{
	Use:   "list [slug]",
	Short: "List all tasks for the current feature",
	Long: `List all tasks for the current feature in a table format.

Displays task ID, type, title (truncated), status, breaking, and mainSession.
Tasks are sorted in topological order by default (dependencies before dependents).
Use --sort id to restore natural ID ordering.

When a slug is provided, lists tasks for that specific feature, reading from
the worktree if one exists for that slug. Use --local to read from the main
repository's index.json regardless of worktree existence.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVar(&listLocal, "local", false, "Read from main repo's index.json (ignore worktree)")
	listCmd.Flags().String("sort", "topo", "Sort order: topo (topological) or id (natural ID)")
	listCmd.Flags().Bool("tree", false, "Display tasks as interactive dependency tree (TUI)")
}

const titleMaxWidth = 50

func runList(cmd *cobra.Command, args []string) error {
	sortMode := "topo"
	if cmd != nil {
		sortMode = cmd.Flags().Lookup("sort").Value.String()
	}
	if sortMode != "topo" && sortMode != "id" {
		return fmt.Errorf("invalid --sort value %q: must be \"topo\" or \"id\"", sortMode)
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	var featureSlug string
	var indexPath string

	if len(args) == 1 {
		// Slug provided: bypass RequireFeature, construct index path directly
		featureSlug = args[0]

		// Validate: feature directory must exist (in main repo or worktree)
		if !featureDirExists(projectRoot, featureSlug) {
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

	// Check for legacy scope fields
	var allTasks []task.Task
	for _, t := range index.TasksMap() {
		allTasks = append(allTasks, t)
	}
	if legacyErr := task.CheckLegacyScope(allTasks); legacyErr != nil {
		scopeErr, ok := legacyErr.(*task.LegacyScopeError)
		if ok {
			base.Exit(base.ErrLegacyScope(scopeErr.Count))
		}
		return legacyErr
	}

	// --tree mode: build dependency tree and render
	treeMode := false
	if cmd != nil {
		treeMode, _ = cmd.Flags().GetBool("tree")
	}
	if treeMode {
		// Terminal capability detection before entering TUI
		isTerminal := listIsTerminalFunc()
		termEnv := os.Getenv("TERM")

		if canUseTUI(isTerminal, termEnv) {
			sortByID := sortMode == "id"
			roots := buildForest(index, withSortByID(sortByID))
			return runTreeTUI(roots)
		}
		// Non-TTY: silently fall back to table mode (continue below)
	}

	// Collect and sort task IDs
	tasks := index.TasksMap()

	var sortedIDs []string
	var cycleSet map[string]bool
	var missingForTask map[string][]string // task ID -> its missing deps

	if sortMode == "topo" {
		ordered, cycles, missing := task.TopologicalSort(index)
		sortedIDs = ordered
		// Append cycle nodes at the end so they still appear in the table
		sortedIDs = append(sortedIDs, cycles...)

		cycleSet = make(map[string]bool, len(cycles))
		for _, c := range cycles {
			cycleSet[c] = true
		}

		// Build missing-deps lookup per task
		missingForTask = buildMissingPerTask(index, missing)
	} else {
		ids := make([]string, 0, len(tasks))
		for _, t := range tasks {
			ids = append(ids, t.ID)
		}
		sortedIDs = naturalSortTaskIDs(ids)
	}

	// Determine if color should be used (TTY only)
	useColor := listIsTerminalFunc()

	// Build display ID with optional marker for each task
	displayID := func(id string) string {
		var markers []string
		if cycleSet != nil && cycleSet[id] {
			markers = append(markers, "cycle")
		}
		if missingForTask != nil {
			if missIDs := missingForTask[id]; len(missIDs) > 0 {
				for _, mid := range missIDs {
					markers = append(markers, "missing: "+mid)
				}
			}
		}
		if len(markers) == 0 {
			return id
		}
		markerText := " [" + strings.Join(markers, ", ") + "]"
		if useColor {
			return id + colorCycleMarker + markerText + colorReset
		}
		return id + markerText
	}

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
		t, _ := index.ByID(id)
		// Use display ID width (includes marker text, but exclude ANSI codes)
		displayW := displayWidthPlain(id, cycleSet, missingForTask)
		if displayW > idCol {
			idCol = displayW
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
		if base.DisplayWidth(string(t.Status)) > statusCol {
			statusCol = base.DisplayWidth(string(t.Status))
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
		t, _ := index.ByID(id)
		title := base.TruncateSlug(t.Title, titleCol)
		idDisplay := displayID(id)
		idDisplayPadded := padRightPlain(idDisplay, idCol)
		fmt.Printf("%s  %s  %s  %s  %s  %s\n",
			idDisplayPadded,
			base.PadRight(t.Type, typeCol),
			base.PadRight(title, titleCol),
			base.PadRight(string(t.Status), statusCol),
			base.PadRight(boolStr(t.Breaking), breakingCol),
			base.PadRight(boolStr(t.MainSession), mainSessCol),
		)
	}

	return nil
}

// displayWidthPlain returns the display width of the ID column for a task,
// accounting for marker text (e.g. " [cycle]", " [missing: 99]") but
// excluding ANSI escape codes.
func displayWidthPlain(id string, cycleSet map[string]bool, missingForTask map[string][]string) int {
	w := base.DisplayWidth(id)
	if cycleSet != nil && cycleSet[id] {
		w += base.DisplayWidth(" [cycle]")
	}
	if missingForTask != nil {
		for _, mid := range missingForTask[id] {
			w += base.DisplayWidth(" [missing: " + mid + "]")
		}
	}
	return w
}

// padRightPlain pads the display ID string to exactly n visible columns.
// ANSI escape codes are not counted as visible width.
func padRightPlain(displayID string, n int) string {
	// Calculate visible width (excluding ANSI codes)
	visibleWidth := 0
	inEscape := false
	for _, r := range displayID {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		visibleWidth += base.DisplayWidth(string(r))
	}
	if visibleWidth >= n {
		return displayID
	}
	return displayID + strings.Repeat(" ", n-visibleWidth)
}

// buildMissingPerTask builds a map from task ID to the list of its missing deps
// that are in the global missing list.
func buildMissingPerTask(idx *task.TaskIndex, globalMissing []string) map[string][]string {
	missingSet := make(map[string]bool, len(globalMissing))
	for _, m := range globalMissing {
		missingSet[m] = true
	}

	result := make(map[string][]string)
	tasks := idx.TasksMap()
	for id, t := range tasks {
		for _, dep := range t.Dependencies {
			switch {
			case missingSet[dep]:
				result[id] = append(result[id], dep)
			case strings.HasSuffix(dep, ".x"):
				// Wildcard: check if it resolved to nothing
				matches, isWildcard := task.ResolveWildcardDep(idx, dep)
				if isWildcard && len(matches) == 0 {
					result[id] = append(result[id], "unresolved: "+dep)
				}
			default:
				// Exact dep: check if it exists
				if _, found := idx.ByID(dep); !found {
					result[id] = append(result[id], dep)
				}
			}
		}
	}
	return result
}

// featureDirExists checks whether a feature directory exists in the main repo
// or inside a worktree (.forge/worktrees/<slug>/docs/features/<slug>/).
func featureDirExists(projectRoot, slug string) bool {
	mainDir := filepath.Join(projectRoot, feature.GetFeatureDir(slug))
	if _, err := os.Stat(mainDir); err == nil {
		return true
	}
	if !listLocal {
		wtDir := filepath.Join(projectRoot, ".forge", "worktrees", slug, feature.GetFeatureDir(slug))
		if _, err := os.Stat(wtDir); err == nil {
			return true
		}
	}
	return false
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
	return idSortKey{numPrefix: fallbackSortPriority}
}
