package task

import (
	"fmt"
	"forge-cli/internal/cmd/base"
	"path/filepath"
	"sort"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var checkDepsCmd = &cobra.Command{
	Use:   "check-deps",
	Short: "Check task dependencies",
	Long: `Check all task dependencies in the current feature.

Validates:
  - All dependencies reference existing tasks
  - Wildcard dependencies match at least one task`,
	Args: cobra.NoArgs,
	RunE: runCheckDeps,
}

type depInfo struct {
	taskKey    string
	taskID     string
	dependency string
	matches    []string
	isWildcard bool
}

func runCheckDeps(_ *cobra.Command, _ []string) error {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return base.ErrFeatureNotSet()
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return base.ErrFileNotFound(indexPath)
	}

	// Collect all task IDs
	taskIDs := make(map[string]bool)
	for _, t := range index.TasksMap() {
		taskIDs[t.ID] = true
	}

	var errors []string
	var depInfos []depInfo

	for key, t := range index.TasksMap() {
		for _, dep := range t.Dependencies {
			matches, isWildcard := task.ResolveWildcardDep(index, dep)

			if isWildcard {
				if len(matches) == 0 {
					errors = append(errors, fmt.Sprintf("Task %s (%s): wildcard '%s' matches NO tasks",
						key, t.ID, dep))
				}
				depInfos = append(depInfos, depInfo{
					taskKey: key, taskID: t.ID, dependency: dep,
					matches: matches, isWildcard: true,
				})
			} else {
				if !taskIDs[dep] {
					errors = append(errors, fmt.Sprintf("Task %s (%s): dependency '%s' does NOT exist",
						key, t.ID, dep))
				}
				depInfos = append(depInfos, depInfo{
					taskKey: key, taskID: t.ID, dependency: dep,
					matches: []string{dep}, isWildcard: false,
				})
			}
		}
	}

	// Output results
	base.PrintSection("TASKS")
	var sortedIDs []string
	for id := range taskIDs {
		sortedIDs = append(sortedIDs, id)
	}
	sort.Slice(sortedIDs, func(i, j int) bool {
		partsI := strings.Split(sortedIDs[i], ".")
		partsJ := strings.Split(sortedIDs[j], ".")
		for k := 0; k < len(partsI) && k < len(partsJ); k++ {
			if partsI[k] != partsJ[k] {
				return partsI[k] < partsJ[k]
			}
		}
		return len(partsI) < len(partsJ)
	})
	for _, id := range sortedIDs {
		base.PrintListItem(id)
	}

	// Output dependencies section
	base.PrintSection("DEPENDENCIES")
	for _, di := range depInfos {
		if di.isWildcard {
			base.PrintListItem(fmt.Sprintf("%s -> [%s] (wildcard)", di.taskID, di.dependency))
		} else {
			base.PrintListItem(fmt.Sprintf("%s -> %s", di.taskID, di.dependency))
		}
	}

	if len(errors) > 0 {
		base.PrintSection("ERRORS")
		for _, e := range errors {
			base.PrintListItem(e)
		}
		base.PrintResult("FAIL", fmt.Sprintf("%d error(s)", len(errors)))
		return fmt.Errorf("%d dependency error(s)", len(errors))
	}

	base.PrintResult("PASS", fmt.Sprintf("%d tasks checked", len(taskIDs)))
	return nil
}
