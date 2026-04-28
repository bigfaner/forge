package task

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
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
		found := false
		for _, t := range index.Tasks {
			if t.ID == dep || strings.HasSuffix(dep, ".x") {
				found = true
				break
			}
		}
		if !found {
			return "", fmt.Errorf("dependency not found: %s", dep)
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

	return os.WriteFile(filepath.Join(tasksDir, filename), []byte(buf.String()), 0644)
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
