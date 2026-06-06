package task

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/internal/cmd/base"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var listLocal bool

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

// colWidths holds computed column widths for the task table.
type colWidths struct {
	id       int
	typ      int
	title    int
	status   int
	breaking int
	mainSess int
}

func runList(cmd *cobra.Command, args []string) error {
	sortMode, err := parseSortMode(cmd)
	if err != nil {
		return err
	}

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	featureSlug, indexPath, err := resolveFeatureArgs(projectRoot, args)
	if err != nil {
		return err
	}

	index, err := loadTaskIndex(indexPath, featureSlug)
	if err != nil {
		if errors.Is(err, errSentinelNoTasks) {
			return nil // message already printed
		}
		return err
	}

	if err := checkLegacyScope(index); err != nil {
		return err
	}

	// --tree mode: build dependency tree and render
	if isTreeMode(cmd) {
		handled, err := handleTreeMode(index, sortMode)
		if handled || err != nil {
			return err
		}
		// Not handled: terminal doesn't support TUI, fall through to table mode
	}

	// Table mode
	sortedIDs, cycleSet, missingForTask := sortTaskIDs(index, sortMode)
	displayID := newDisplayIDFunc(cycleSet, missingForTask)

	fmt.Printf("%d found  (feature: %s)\n\n", len(sortedIDs), featureSlug)

	widths := computeColumnWidths(sortedIDs, index, cycleSet, missingForTask)
	printTaskTable(sortedIDs, index, widths, displayID)

	return nil
}

// parseSortMode extracts and validates the --sort flag value.
func parseSortMode(cmd *cobra.Command) (string, error) {
	sortMode := "topo"
	if cmd != nil {
		sortMode = cmd.Flags().Lookup("sort").Value.String()
	}
	if sortMode != "topo" && sortMode != "id" {
		return "", fmt.Errorf("invalid --sort value %q: must be \"topo\" or \"id\"", sortMode)
	}
	return sortMode, nil
}

// resolveFeatureArgs determines the feature slug and index path from CLI args.
func resolveFeatureArgs(projectRoot string, args []string) (featureSlug, indexPath string, err error) {
	if len(args) == 1 {
		featureSlug = args[0]
		if !featureDirExists(projectRoot, featureSlug) {
			return "", "", base.ErrFeatureNotFound(featureSlug)
		}
		indexPath = resolveListIndexPath(projectRoot, featureSlug)
		return featureSlug, indexPath, nil
	}

	// No slug: use existing auto-detection logic
	var err2 error
	featureSlug, err2 = feature.RequireFeature(projectRoot)
	if err2 != nil {
		base.Exit(base.ErrFeatureNotSet())
	}
	indexPath = filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	return featureSlug, indexPath, nil
}

// loadTaskIndex loads the task index, printing a message and returning a
// sentinel error if the file is missing or contains no tasks.
func loadTaskIndex(indexPath, featureSlug string) (*task.TaskIndex, error) {
	index, err := task.LoadIndex(indexPath)
	if err != nil || index.TaskCount() == 0 {
		fmt.Printf("no tasks found  (feature: %s)\n", featureSlug)
		return nil, errSentinelNoTasks
	}
	return index, nil
}

// errSentinelNoTasks is returned by loadTaskIndex when no tasks are found.
// Callers should return nil for this case (the message is already printed).
var errSentinelNoTasks = errors.New("no tasks found")

// checkLegacyScope verifies no legacy scope fields remain in tasks.
func checkLegacyScope(index *task.TaskIndex) error {
	var allTasks []task.Task
	for _, t := range index.TasksMap() {
		allTasks = append(allTasks, t)
	}
	if legacyErr := task.CheckLegacyScope(allTasks); legacyErr != nil {
		if scopeErr, ok := legacyErr.(*task.LegacyScopeError); ok {
			base.Exit(base.ErrLegacyScope(scopeErr.Count))
		}
		return legacyErr
	}
	return nil
}

// isTreeMode returns true if the --tree flag is set.
func isTreeMode(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	treeMode, _ := cmd.Flags().GetBool("tree")
	return treeMode
}

// handleTreeMode attempts TUI rendering. Returns (handled, error):
//   - (true, err) if TUI was launched (err may be nil on success)
//   - (false, nil) if terminal doesn't support TUI (fall back to table mode)
func handleTreeMode(index *task.TaskIndex, sortMode string) (bool, error) {
	isTerminal := listIsTerminalFunc()
	termEnv := os.Getenv("TERM")

	if !canUseTUI(isTerminal, termEnv) {
		return false, nil // fall back to table mode
	}

	sortByID := sortMode == "id"
	roots := buildForest(index, withSortByID(sortByID))
	return true, runTreeTUI(roots)
}

// sortTaskIDs returns sorted task IDs along with cycle and missing-dep metadata.
func sortTaskIDs(index *task.TaskIndex, sortMode string) (sortedIDs []string, cycleSet map[string]bool, missingForTask map[string][]string) {
	if sortMode != "topo" {
		ids := make([]string, 0, len(index.TasksMap()))
		for _, t := range index.TasksMap() {
			ids = append(ids, t.ID)
		}
		return naturalSortTaskIDs(ids), nil, nil
	}

	ordered, cycles, missing := task.TopologicalSort(index)
	sortedIDs = ordered
	sortedIDs = append(sortedIDs, cycles...)

	cycleSet = make(map[string]bool, len(cycles))
	for _, c := range cycles {
		cycleSet[c] = true
	}

	missingForTask = buildMissingPerTask(index, missing)
	return sortedIDs, cycleSet, missingForTask
}

// newDisplayIDFunc returns a closure that formats a task ID with optional
// cycle/missing markers and ANSI color codes.
func newDisplayIDFunc(cycleSet map[string]bool, missingForTask map[string][]string) func(string) string {
	useColor := listIsTerminalFunc()

	return func(id string) string {
		var markers []string
		if cycleSet != nil && cycleSet[id] {
			markers = append(markers, "cycle")
		}
		if missingForTask != nil {
			for _, mid := range missingForTask[id] {
				markers = append(markers, "missing: "+mid)
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
}

// computeColumnWidths calculates dynamic column widths for the task table.
func computeColumnWidths(sortedIDs []string, index *task.TaskIndex, cycleSet map[string]bool, missingForTask map[string][]string) colWidths {
	w := colWidths{
		id:       base.DisplayWidth("ID"),
		typ:      base.DisplayWidth("TYPE"),
		title:    base.DisplayWidth("TITLE"),
		status:   base.DisplayWidth("STATUS"),
		breaking: base.DisplayWidth("BREAKING"),
		mainSess: base.DisplayWidth("MAIN_SESS"),
	}

	for _, id := range sortedIDs {
		t, _ := index.ByID(id)
		displayW := displayWidthPlain(id, cycleSet, missingForTask)
		if displayW > w.id {
			w.id = displayW
		}
		if dw := base.DisplayWidth(t.Type); dw > w.typ {
			w.typ = dw
		}
		titleLen := base.DisplayWidth(t.Title)
		if titleLen > titleMaxWidth {
			titleLen = titleMaxWidth
		}
		if titleLen > w.title {
			w.title = titleLen
		}
		if dw := base.DisplayWidth(string(t.Status)); dw > w.status {
			w.status = dw
		}
		boolBreaking := boolStr(t.Breaking)
		if dw := base.DisplayWidth(boolBreaking); dw > w.breaking {
			w.breaking = dw
		}
		boolMainSess := boolStr(t.MainSession)
		if dw := base.DisplayWidth(boolMainSess); dw > w.mainSess {
			w.mainSess = dw
		}
	}
	return w
}

// boolStr converts a bool to its lowercase string representation.
func boolStr(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

// printTaskTable renders the task table with headers, separator, and rows.
func printTaskTable(sortedIDs []string, index *task.TaskIndex, w colWidths, displayID func(string) string) {
	// Print column headers
	fmt.Printf("%s  %s  %s  %s  %s  %s\n",
		base.PadRight("ID", w.id),
		base.PadRight("TYPE", w.typ),
		base.PadRight("TITLE", w.title),
		base.PadRight("STATUS", w.status),
		base.PadRight("BREAKING", w.breaking),
		base.PadRight("MAIN_SESS", w.mainSess),
	)

	// Print separator
	sep := fmt.Sprintf("%s  %s  %s  %s  %s  %s",
		strings.Repeat("-", w.id),
		strings.Repeat("-", w.typ),
		strings.Repeat("-", w.title),
		strings.Repeat("-", w.status),
		strings.Repeat("-", w.breaking),
		strings.Repeat("-", w.mainSess),
	)
	fmt.Println(sep)

	// Print task rows
	for _, id := range sortedIDs {
		t, _ := index.ByID(id)
		title := base.TruncateSlug(t.Title, w.title)
		idDisplay := displayID(id)
		idDisplayPadded := padRightPlain(idDisplay, w.id)
		fmt.Printf("%s  %s  %s  %s  %s  %s\n",
			idDisplayPadded,
			base.PadRight(t.Type, w.typ),
			base.PadRight(title, w.title),
			base.PadRight(string(t.Status), w.status),
			base.PadRight(boolStr(t.Breaking), w.breaking),
			base.PadRight(boolStr(t.MainSession), w.mainSess),
		)
	}
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
