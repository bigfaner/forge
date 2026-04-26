package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"
	"task-cli/pkg/task"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check task dependencies",
	Long: `Check all task dependencies in the current feature.

Validates:
  - All dependencies reference existing tasks
  - Wildcard dependencies match at least one task`,
	Run: runCheck,
}

type depInfo struct {
	taskKey    string
	taskID     string
	dependency string
	matches    []string
	isWildcard bool
}

func runCheck(cmd *cobra.Command, args []string) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		Exit(ErrProjectNotFound())
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		Exit(ErrFeatureNotSet())
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		Exit(ErrFileNotFound(indexPath))
	}

	// Collect all task IDs
	taskIDs := make(map[string]bool)
	for _, t := range index.Tasks {
		taskIDs[t.ID] = true
	}

	var errors []string
	var depInfos []depInfo

	for key, t := range index.Tasks {
		for _, dep := range t.Dependencies {
			isWildcard := strings.HasSuffix(dep, ".x") || strings.HasSuffix(dep, "x")

			if isWildcard {
				prefix := strings.TrimSuffix(strings.TrimSuffix(dep, "x"), ".")
				prefixWithDot := prefix + "."

				var matches []string
				for id := range taskIDs {
					if strings.HasPrefix(id, prefixWithDot) && isBusinessTask(id) {
						matches = append(matches, id)
					}
				}

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
	PrintSection("TASKS")
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
		PrintListItem(id)
	}

	// Output dependencies section
	PrintSection("DEPENDENCIES")
	for _, di := range depInfos {
		if di.isWildcard {
			fmt.Printf("  %s -> [%s] (wildcard)\n", di.taskID, di.dependency)
		} else {
			fmt.Printf("  %s -> %s\n", di.taskID, di.dependency)
		}
	}

	if len(errors) > 0 {
		PrintSection("ERRORS")
		for _, e := range errors {
			PrintListItem(e)
		}
		PrintResult("FAIL", fmt.Sprintf("%d error(s)", len(errors)))
		os.Exit(1)
	}

	PrintResult("PASS", fmt.Sprintf("%d tasks checked", len(taskIDs)))
}
