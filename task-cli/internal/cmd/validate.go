package cmd

import (
	"encoding/json"
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

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate index.json file",
	Long: `Validate an index.json file for structural and semantic correctness.

If no file is specified, validates the current feature's index.json.

Validations:
  - JSON syntax
  - Required fields present
  - Dependency references exist
  - No circular dependencies`,
	Run: runValidate,
}

var (
	validStatus   = map[string]bool{"pending": true, "in_progress": true, "completed": true, "blocked": true, "skipped": true}
	validPriority = map[string]bool{"P0": true, "P1": true, "P2": true}
)

func runValidate(cmd *cobra.Command, args []string) {
	var filePath string
	if len(args) > 0 {
		filePath = args[0]
	} else {
		projectRoot, err := project.FindProjectRoot()
		if err != nil {
			Exit(ErrProjectNotFound())
		}
		slug, err := feature.RequireFeature(projectRoot)
		if err != nil {
			Exit(ErrFeatureNotSet())
		}
		filePath = filepath.Join(projectRoot, feature.GetFeatureIndexFile(slug))
	}

	v := &validator{filePath: filePath}
	if err := v.run(); err != nil {
		Exit(err)
	}
}

type validator struct {
	filePath string
	errors   []string
	warnings []string
	info     []string
}

func (v *validator) run() error {
	data, err := os.ReadFile(v.filePath)
	if err != nil {
		return ErrFileNotFound(v.filePath)
	}

	var idx task.TaskIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return ErrInvalidJSON(v.filePath, err.Error())
	}

	if idx.Feature == "" {
		return NewAIError(ErrValidation, "Missing required field", "'feature' field is empty", "Add feature name to index.json", "Add \"feature\": \"<name>\" to index.json")
	}

	v.info = append(v.info, fmt.Sprintf("Feature: %s", idx.Feature))
	v.info = append(v.info, fmt.Sprintf("Tasks: %d", len(idx.Tasks)))

	if idx.PRD == "" {
		v.warnings = append(v.warnings, "Missing 'prd' field")
	}
	if idx.Design == "" {
		v.warnings = append(v.warnings, "Missing 'design' field")
	}
	if len(idx.StatusEnum) == 0 {
		v.warnings = append(v.warnings, "Missing 'statusEnum' field — task record/status commands may fail")
	}

	v.validateTasks(idx.Tasks)
	v.validateDependencies(idx.Tasks)
	v.validateCircularDeps(idx.Tasks)
	v.validateFilesExist(idx.Feature, idx.Tasks)

	if !v.printResults() {
		return NewAIError(ErrValidation, "Validation failed", fmt.Sprintf("%d errors found", len(v.errors)), "Fix errors in index.json", "cat "+v.filePath)
	}
	return nil
}

func (v *validator) validateTasks(tasks map[string]task.Task) {
	for key, t := range tasks {
		if t.ID == "" {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'id'", key))
		}
		if t.Title == "" {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'title'", key))
		}
		if t.File == "" {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': missing 'file'", key))
		}
		if t.Status != "" && !validStatus[t.Status] {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': invalid status '%s'", key, t.Status))
		}
		if t.Priority != "" && !validPriority[t.Priority] {
			v.errors = append(v.errors, fmt.Sprintf("Task '%s': invalid priority '%s'", key, t.Priority))
		}
	}
}

func (v *validator) validateDependencies(tasks map[string]task.Task) {
	taskIDs := make(map[string]bool)
	for _, t := range tasks {
		taskIDs[t.ID] = true
	}

	for key, t := range tasks {
		for _, dep := range t.Dependencies {
			isWildcard := strings.HasSuffix(dep, ".x")

			if isWildcard {
				prefix := strings.TrimSuffix(dep, ".x") + "."
				var matches []string
				for id := range taskIDs {
					if strings.HasPrefix(id, prefix) {
						matches = append(matches, id)
					}
				}
				if len(matches) == 0 {
					v.errors = append(v.errors, fmt.Sprintf("Task '%s': wildcard '%s' matches nothing", key, dep))
				}
			} else {
				if !taskIDs[dep] {
					v.errors = append(v.errors, fmt.Sprintf("Task '%s': dependency '%s' not found", key, dep))
				}
			}
		}
	}
}

func (v *validator) validateCircularDeps(tasks map[string]task.Task) {
	graph := make(map[string][]string)
	taskIDs := make(map[string]bool)
	for _, t := range tasks {
		taskIDs[t.ID] = true
	}
	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if !strings.HasSuffix(dep, ".x") && taskIDs[dep] {
				graph[t.ID] = append(graph[t.ID], dep)
			}
		}
	}

	visited := make(map[string]bool)
	inStack := make(map[string]bool)

	var dfs func(id string, path []string) bool
	dfs = func(id string, path []string) bool {
		visited[id] = true
		inStack[id] = true
		for _, neighbor := range graph[id] {
			if !visited[neighbor] {
				if dfs(neighbor, append(path, neighbor)) {
					return true
				}
			} else if inStack[neighbor] {
				v.errors = append(v.errors, fmt.Sprintf("Circular: %s", strings.Join(append(path, neighbor), " -> ")))
				return true
			}
		}
		inStack[id] = false
		return false
	}

	for id := range taskIDs {
		if !visited[id] && dfs(id, []string{id}) {
			break
		}
	}
}

func (v *validator) validateFilesExist(featureSlug string, tasks map[string]task.Task) {
	// Get project root from filePath (which is .../docs/features/<slug>/tasks/index.json)
	// Go up 4 levels: index.json -> tasks -> <slug> -> features -> docs -> projectRoot
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(v.filePath)))))
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	for key, t := range tasks {
		if t.File == "" {
			continue
		}
		taskFile := filepath.Join(tasksDir, t.File)
		if _, err := os.Stat(taskFile); os.IsNotExist(err) {
			v.warnings = append(v.warnings, fmt.Sprintf("Task '%s': file '%s' missing", key, t.File))
		}
	}
}

func (v *validator) printResults() bool {
	if len(v.info) > 0 {
		PrintSection("INFO")
		sort.Strings(v.info)
		for _, i := range v.info {
			PrintListItem(i)
		}
	}

	if len(v.warnings) > 0 {
		PrintSection("WARNINGS")
		sort.Strings(v.warnings)
		for _, w := range v.warnings {
			PrintListItem(w)
		}
	}

	if len(v.errors) > 0 {
		PrintSection("ERRORS")
		sort.Strings(v.errors)
		for _, e := range v.errors {
			PrintListItem(e)
		}
	}

	if len(v.errors) == 0 {
		PrintResult("PASS", v.filePath)
		return true
	}
	PrintResult("FAIL", fmt.Sprintf("%s (%d errors)", v.filePath, len(v.errors)))
	return false
}
