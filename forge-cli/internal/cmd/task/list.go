package task

import (
	"fmt"
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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks for the current feature",
	Long: `List all tasks for the current feature in a table format.

Displays task ID, type, title (truncated), and status.
Tasks are sorted by ID in natural order: numeric IDs first, then test/gate IDs.`,
	Args: cobra.NoArgs,
	RunE: runList,
}

const titleMaxWidth = 50

func runList(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		base.Exit(base.ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		base.Exit(base.ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

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

	// Calculate dynamic column widths based on actual data
	idCol := len("ID")
	typeCol := len("TYPE")
	titleCol := len("TITLE")
	statusCol := len("STATUS")

	for _, id := range sortedIDs {
		t := tasks[id]
		if len(t.ID) > idCol {
			idCol = len(t.ID)
		}
		if len(t.Type) > typeCol {
			typeCol = len(t.Type)
		}
		titleLen := len(t.Title)
		if titleLen > titleMaxWidth {
			titleLen = titleMaxWidth
		}
		if titleLen > titleCol {
			titleCol = titleLen
		}
		if len(t.Status) > statusCol {
			statusCol = len(t.Status)
		}
	}

	// Print column headers
	fmt.Printf("%s  %s  %s  %s\n",
		base.PadRight("ID", idCol),
		base.PadRight("TYPE", typeCol),
		base.PadRight("TITLE", titleCol),
		base.PadRight("STATUS", statusCol),
	)

	// Print separator
	sep := fmt.Sprintf("%s  %s  %s  %s",
		strings.Repeat("-", idCol),
		strings.Repeat("-", typeCol),
		strings.Repeat("-", titleCol),
		strings.Repeat("-", statusCol),
	)
	fmt.Println(sep)

	// Print task rows
	for _, id := range sortedIDs {
		t := tasks[id]
		title := base.TruncateSlug(t.Title, titleCol)
		fmt.Printf("%s  %s  %s  %s\n",
			base.PadRight(t.ID, idCol),
			base.PadRight(t.Type, typeCol),
			base.PadRight(title, titleCol),
			base.PadRight(t.Status, statusCol),
		)
	}

	return nil
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
