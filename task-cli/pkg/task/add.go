package task

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	tmpl "task-cli/pkg/template"
)

// AddTaskOpts holds options for adding a new task.
type AddTaskOpts struct {
	ID            string   // Auto-generated as disc-N if empty
	Title         string   // Required
	Priority      string   // Default P1 if empty
	EstimatedTime string   // Optional
	Dependencies  []string // Optional, validated against existing tasks
	Breaking      bool     // Optional
	Description   string   // Optional, becomes markdown body
	Status        string   // Default pending if empty
	Template      string   // Template name (e.g. "fix-task")
	Vars          map[string]string // Variable substitutions for template placeholders
	SourceTaskID  string   // Source task ID: auto-injects {{SOURCE_TASK_ID}} and adds this task as source dependency
}

// AddTask validates opts, adds a task to the index, and saves it.
// Returns the generated or provided task ID.
func AddTask(indexPath string, opts AddTaskOpts) (string, error) {
	if opts.Title == "" {
		return "", fmt.Errorf("title is required")
	}

	index, err := LoadIndex(indexPath)
	if err != nil {
		return "", fmt.Errorf("load index: %w", err)
	}

	// Defaults
	if opts.Status == "" {
		opts.Status = "pending"
	}
	if opts.Priority == "" {
		opts.Priority = "P1"
	}

	// Auto-generate ID if empty
	if opts.ID == "" {
		opts.ID = generateDiscID(index)
	}

	// Validate ID uniqueness
	for key, t := range index.Tasks {
		if t.ID == opts.ID || key == opts.ID {
			return "", fmt.Errorf("task ID already exists: %s", opts.ID)
		}
	}

	// Validate priority
	if !slices.Contains([]string{"P0", "P1", "P2"}, opts.Priority) {
		return "", fmt.Errorf("invalid priority: %s (must be P0, P1, or P2)", opts.Priority)
	}

	// Validate dependencies exist
	for _, dep := range opts.Dependencies {
		if strings.HasSuffix(dep, ".x") {
			prefix := strings.TrimSuffix(dep, ".x")
			prefixWithDot := prefix + "."
			found := false
			for _, t := range index.Tasks {
				if strings.HasPrefix(t.ID, prefixWithDot) && isBusinessTaskID(t.ID) {
					found = true
					break
				}
			}
			if !found {
				return "", fmt.Errorf("wildcard dependency %q matches no business tasks", dep)
			}
		} else {
			found := false
			for _, t := range index.Tasks {
				if t.ID == dep {
					found = true
					break
				}
			}
			if !found {
				return "", fmt.Errorf("dependency not found: %s", dep)
			}
		}
	}

	// Derive file and record paths
	fileName := opts.ID + ".md"
	recordPath := "records/" + opts.ID + ".md"

	key := opts.ID
	if strings.Contains(opts.ID, "-") {
		// Use ID as key directly for disc-*, fix-* style IDs
		key = opts.ID
	}

	index.Tasks[key] = Task{
		ID:            opts.ID,
		Title:         opts.Title,
		Priority:      opts.Priority,
		EstimatedTime: opts.EstimatedTime,
		Dependencies:  opts.Dependencies,
		Status:        opts.Status,
		File:          fileName,
		Record:        recordPath,
		Breaking:      opts.Breaking,
		SourceTaskID:  opts.SourceTaskID,
	}

	// SourceTaskID: add this task as dependency of the source task (in-memory, before single save)
	// Root cause: must iterate by ID or key, not use direct map access with ID, because map keys are slugs
	if opts.SourceTaskID != "" {
		for key, t := range index.Tasks {
			if t.ID == opts.SourceTaskID || key == opts.SourceTaskID {
				if !slices.Contains(t.Dependencies, opts.ID) {
					t.Dependencies = append(t.Dependencies, opts.ID)
					index.Tasks[key] = t
				}
				break
			}
		}
	}

	if err := SaveIndex(indexPath, index); err != nil {
		return "", fmt.Errorf("save index: %w", err)
	}

	return opts.ID, nil
}

// CreateTaskMarkdown writes a task markdown file with YAML frontmatter.
func CreateTaskMarkdown(tasksDir string, filename string, opts AddTaskOpts) error {
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return fmt.Errorf("create tasks dir: %w", err)
	}

	var content string

	if opts.Template != "" {
		tmpl, err := tmpl.Get(opts.Template)
		if err != nil {
			return err
		}
		content = ApplyVars(tmpl, opts)
	} else {
		content = buildTaskMarkdown(opts)
	}

	return os.WriteFile(filepath.Join(tasksDir, filename), []byte(content), 0644)
}

// ApplyVars replaces {{KEY}} placeholders in tmpl with values from opts.Vars
// and built-in variables (ID, TITLE, PRIORITY, DESCRIPTION).
// User-provided variables take precedence over builtins.
func ApplyVars(tmpl string, opts AddTaskOpts) string {
	result := tmpl

	// Build merged variable map: user vars override builtins
	vars := map[string]string{
		"ID":          opts.ID,
		"TITLE":       opts.Title,
		"PRIORITY":    opts.Priority,
		"DESCRIPTION": opts.Description,
		"SOURCE_TASK_ID": opts.SourceTaskID,
	}
	for key, val := range opts.Vars {
		vars[key] = val
	}

	for key, val := range vars {
		result = strings.ReplaceAll(result, "{{"+key+"}}", val)
	}

	return result
}

// buildTaskMarkdown generates task markdown from scratch (non-template mode).
func buildTaskMarkdown(opts AddTaskOpts) string {
	var buf strings.Builder

	// Frontmatter
	buf.WriteString("---\n")
	fmt.Fprintf(&buf, "id: %q\n", opts.ID)
	fmt.Fprintf(&buf, "title: %q\n", opts.Title)
	fmt.Fprintf(&buf, "priority: %q\n", opts.Priority)
	if opts.EstimatedTime != "" {
		fmt.Fprintf(&buf, "estimated_time: %q\n", opts.EstimatedTime)
	}
	if len(opts.Dependencies) > 0 {
		buf.WriteString("dependencies:\n")
		for _, dep := range opts.Dependencies {
			fmt.Fprintf(&buf, "  - %q\n", dep)
		}
	} else {
		buf.WriteString("dependencies: []\n")
	}
	fmt.Fprintf(&buf, "status: %s\n", opts.Status)
	if opts.Breaking {
		buf.WriteString("breaking: true\n")
	}
	buf.WriteString("---\n\n")

	// Title and body
	fmt.Fprintf(&buf, "# %s: %s\n\n", opts.ID, opts.Title)

	if opts.Description != "" {
		buf.WriteString(opts.Description)
		buf.WriteString("\n")
	} else {
		buf.WriteString("## Description\n\n_TBD_\n")
	}

	return buf.String()
}

// generateDiscID generates the next available disc-N ID using gap-filling.
func generateDiscID(index *TaskIndex) string {
	used := make(map[int]bool)
	for key := range index.Tasks {
		numStr, ok := strings.CutPrefix(key, "disc-")
		if ok {
			if n, err := strconv.Atoi(numStr); err == nil && n > 0 {
				used[n] = true
			}
		}
	}
	for i := 1; ; i++ {
		if !used[i] {
			return fmt.Sprintf("disc-%d", i)
		}
	}
}

// AddDependency adds depID to the specified task's Dependencies in the index.
// If depID is already listed, this is a no-op.
func AddDependency(indexPath string, taskID string, depID string) error {
	index, err := LoadIndex(indexPath)
	if err != nil {
		return fmt.Errorf("load index: %w", err)
	}

	t, ok := index.Tasks[taskID]
	if !ok {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if slices.Contains(t.Dependencies, depID) {
		return nil
	}

	t.Dependencies = append(t.Dependencies, depID)
	index.Tasks[taskID] = t

	return SaveIndex(indexPath, index)
}

// GetUnmetDependencies returns the list of dependency IDs that are not "completed" or "skipped".
// Missing task IDs are treated as unmet.
func GetUnmetDependencies(indexPath string, taskID string) ([]string, error) {
	index, err := LoadIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("load index: %w", err)
	}

	t, ok := index.Tasks[taskID]
	if !ok {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	var unmet []string
	for _, dep := range t.Dependencies {
		if strings.HasSuffix(dep, ".x") {
			prefix := strings.TrimSuffix(dep, ".x")
			prefixWithDot := prefix + "."
			for _, other := range index.Tasks {
				if other.ID == t.ID {
					continue
				}
				if strings.HasPrefix(other.ID, prefixWithDot) && isBusinessTaskID(other.ID) && other.Status != "completed" && other.Status != "skipped" {
					unmet = append(unmet, other.ID)
				}
			}
			continue
		}
		depTask, found := index.Tasks[dep]
		if !found || (depTask.Status != "completed" && depTask.Status != "skipped") {
			unmet = append(unmet, dep)
		}
	}
	return unmet, nil
}

// isBusinessTaskID returns true for task IDs that are not gate or summary tasks.
func isBusinessTaskID(id string) bool {
	return !strings.HasSuffix(id, ".gate") && !strings.HasSuffix(id, ".summary")
}
