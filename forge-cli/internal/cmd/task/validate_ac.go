package task

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

// validateACCount checks that each task .md file has between 1 and 6 acceptance criteria.
func (v *validator) validateACCount(featureSlug string, tasks map[string]task.Task) {
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(v.filePath)))))
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	for _, t := range tasks {
		if t.File == "" {
			continue
		}
		taskFile := filepath.Join(tasksDir, t.File)
		data, err := os.ReadFile(taskFile)
		if err != nil {
			continue // File existence already checked in validateFilesExist
		}

		acCount := countACItems(string(data))
		if acCount == 0 {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s' (%s): has 0 acceptance criteria (must have 1-6)", t.File, t.ID))
		} else if acCount > 6 {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s' (%s): has %d acceptance criteria (max 6)", t.File, t.ID, acCount))
		}
	}
}

// countACItems parses a task .md file content and counts `- [ ]` lines
// under the `## Acceptance Criteria` section.
func countACItems(content string) int {
	lines := strings.Split(content, "\n")
	inACSection := false
	count := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Detect section headers
		if strings.HasPrefix(trimmed, "## ") {
			if trimmed == "## Acceptance Criteria" {
				inACSection = true
			} else {
				inACSection = false
			}
			continue
		}
		if inACSection && strings.HasPrefix(trimmed, "- [ ]") {
			count++
		}
	}
	return count
}
