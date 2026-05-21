package task

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	tmpl "forge-cli/pkg/template"
)

// AddTaskOpts holds options for adding a new task.
type AddTaskOpts struct {
	ID            string            // Auto-generated as prefix-N if empty
	Title         string            // Required
	Priority      string            // Default P1 if empty
	EstimatedTime string            // Optional
	Dependencies  []string          // Optional, validated against existing tasks
	Breaking      bool              // Optional
	Description   string            // Optional, becomes markdown body
	Status        string            // Default pending if empty
	Template      string            // Template name (e.g. "fix-task")
	Vars          map[string]string // Variable substitutions for template placeholders
	SourceTaskID  string            // Source task ID: auto-injects {{SOURCE_TASK_ID}} and adds this task as source dependency
	BlockSource   bool              // Block source task before resolution (preserves fix-chain model)
	IDPrefix      string            // Auto-generate ID as prefix-N; empty defaults to "disc"
	Type          string            // Task type (e.g. TypeCodingFix, TypeCodingCleanup). Empty = no type set.
}

// terminalStatuses are task statuses that indicate the task is done.
var terminalStatuses = map[string]bool{
	"completed": true,
	"skipped":   true,
	"rejected":  true,
}

// placeholderRe matches remaining {{KEY}} patterns after variable substitution.
var placeholderRe = regexp.MustCompile(`\{\{(\w+)\}\}`)

// ActiveFixExistsError is returned by AddTask when active fix tasks already exist
// for the specified source task, making the new addition redundant.
type ActiveFixExistsError struct {
	SourceTaskID string
	ActiveFixIDs []string
}

func (e *ActiveFixExistsError) Error() string {
	return fmt.Sprintf("active fix tasks already exist for source %s: %s", e.SourceTaskID, strings.Join(e.ActiveFixIDs, ", "))
}

// hasActiveFixTasks returns IDs of fix tasks targeting sourceTaskID that are not in a terminal state.
func hasActiveFixTasks(index *TaskIndex, sourceTaskID string) []string {
	var active []string
	for _, t := range index.tasks {
		if t.SourceTaskID == sourceTaskID && !terminalStatuses[t.Status] {
			active = append(active, t.ID)
		}
	}
	return active
}

// ResolveSourceTask traces SourceTaskID chains to find the root ancestor.
// If sourceID points to a fix-task that itself has a SourceTaskID, this follows
// the chain until it reaches a task without a SourceTaskID (the original blocked task).
// Returns the original sourceID if no chain exists or the source task is not found.
func ResolveSourceTask(index *TaskIndex, sourceID string) string {
	visited := make(map[string]bool)
	current := sourceID
	for !visited[current] {
		visited[current] = true
		_, t, err := FindTask(index, current)
		if err != nil || t.SourceTaskID == "" {
			break
		}
		current = t.SourceTaskID
	}
	return current
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
		prefix := opts.IDPrefix
		if prefix == "" {
			prefix = "disc"
		}
		opts.ID = generateAutoID(prefix, index)
	}

	// Validate ID uniqueness
	if _, exists := index.ByID(opts.ID); exists {
		return "", fmt.Errorf("task ID already exists: %s", opts.ID)
	}

	// Validate priority
	if !slices.Contains([]string{"P0", "P1", "P2"}, opts.Priority) {
		return "", fmt.Errorf("invalid priority: %s (must be P0, P1, or P2)", opts.Priority)
	}

	// Validate dependencies exist
	for _, dep := range opts.Dependencies {
		if strings.HasSuffix(dep, IDSuffixWildcard) {
			prefix := strings.TrimSuffix(dep, IDSuffixWildcard)
			prefixWithDot := prefix + "."
			found := false
			for _, t := range index.tasks {
				if strings.HasPrefix(t.ID, prefixWithDot) && IsBusinessTask(t.ID) {
					found = true
					break
				}
			}
			if !found {
				return "", fmt.Errorf("wildcard dependency %q matches no business tasks", dep)
			}
		} else {
			if _, exists := index.ByID(dep); !exists {
				return "", fmt.Errorf("dependency not found: %s", dep)
			}
		}
	}

	// Derive file and record paths
	fileName := opts.ID + ".md"
	recordPath := "records/" + opts.ID + ".md"

	// Source handling: dedup → block → resolve.
	// Dedup is a pure read (no mutation), so it must come first.
	// Block before resolution preserves the fix-chain model.
	var srcKey string
	var srcTask *Task
	if opts.SourceTaskID != "" {
		var k string
		var t *Task
		if foundK, foundT, err := FindTask(index, opts.SourceTaskID); err == nil {
			k, t = foundK, foundT
		}

		// Dedup check (pure read): if active fix tasks already exist, skip — no mutation needed.
		if activeFixes := hasActiveFixTasks(index, opts.SourceTaskID); len(activeFixes) > 0 {
			return "", &ActiveFixExistsError{
				SourceTaskID: opts.SourceTaskID,
				ActiveFixIDs: activeFixes,
			}
		}

		// Only mutate after dedup passes.
		if t != nil {
			// Block source before resolution (--block-source flag).
			// Prevents auto-resolve from flattening to root, preserving the chain.
			if opts.BlockSource {
				t.Status = "blocked"
			}

			// Source auto-resolution: when --source-task-id points to a COMPLETED/SKIPPED
			// fix-task, trace the chain to find the root blocked task.
			if t.Status == "completed" || t.Status == "skipped" {
				resolved := ResolveSourceTask(index, opts.SourceTaskID)
				if resolved != opts.SourceTaskID {
					fmt.Fprintf(os.Stderr, "SOURCE-RESOLVE: %s \xe2\x86\x92 %s (source completed, resolving to root)\n", opts.SourceTaskID, resolved)
					opts.SourceTaskID = resolved
					if k2, t2, err2 := FindTask(index, opts.SourceTaskID); err2 == nil {
						k, t = k2, t2
					}
				}
			}
			srcKey, srcTask = k, t
		}
	}

	index.SetTask(opts.ID, Task{
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
		Type:          opts.Type,
	})

	if srcTask != nil && !slices.Contains(srcTask.Dependencies, opts.ID) {
		srcTask.Dependencies = append(srcTask.Dependencies, opts.ID)
		index.SetTask(srcKey, *srcTask)
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
		content, err = ApplyVars(tmpl, opts)
		if err != nil {
			return err
		}
	} else {
		content = buildTaskMarkdown(opts)
	}

	return os.WriteFile(filepath.Join(tasksDir, filename), []byte(content), 0644)
}

// ApplyVars replaces {{KEY}} placeholders in tmpl with values from opts.Vars
// and built-in variables (ID, TITLE, PRIORITY, DESCRIPTION).
// User-provided variables take precedence over builtins.
// Returns an error if any {{...}} placeholders remain unfilled after substitution.
func ApplyVars(tmpl string, opts AddTaskOpts) (string, error) {
	result := tmpl

	// Build merged variable map: user vars override builtins
	vars := map[string]string{
		"ID":             opts.ID,
		"TITLE":          opts.Title,
		"PRIORITY":       opts.Priority,
		"DESCRIPTION":    opts.Description,
		"SOURCE_TASK_ID": opts.SourceTaskID,
	}
	for key, val := range opts.Vars {
		vars[key] = val
	}

	for key, val := range vars {
		result = strings.ReplaceAll(result, "{{"+key+"}}", val)
	}

	// Check for remaining unfilled placeholders
	var unfilled []string
	for _, match := range placeholderRe.FindAllStringSubmatch(result, -1) {
		if len(match) > 1 {
			unfilled = append(unfilled, match[1])
		}
	}
	if len(unfilled) > 0 {
		return "", fmt.Errorf("unfilled template variables: %s", strings.Join(unfilled, ", "))
	}

	return result, nil
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
	if opts.Type != "" {
		fmt.Fprintf(&buf, "type: %q\n", opts.Type)
	}
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

// generateAutoID generates the next prefix-N ID as max(existing) + 1.
func generateAutoID(prefix string, index *TaskIndex) string {
	maxN := 0
	prefixWithDash := prefix + "-"
	for key := range index.tasks {
		numStr, ok := strings.CutPrefix(key, prefixWithDash)
		if ok {
			if n, err := strconv.Atoi(numStr); err == nil && n > maxN {
				maxN = n
			}
		}
	}
	return fmt.Sprintf("%s-%d", prefix, maxN+1)
}

// AddDependency adds depID to the specified task's Dependencies in the index.
// If depID is already listed, this is a no-op.
func AddDependency(indexPath string, taskID string, depID string) error {
	index, err := LoadIndex(indexPath)
	if err != nil {
		return fmt.Errorf("load index: %w", err)
	}

	taskKey, foundTask, err := FindTask(index, taskID)
	if err != nil {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if slices.Contains(foundTask.Dependencies, depID) {
		return nil
	}

	foundTask.Dependencies = append(foundTask.Dependencies, depID)
	index.SetTask(taskKey, *foundTask)

	return SaveIndex(indexPath, index)
}

// GetUnmetDependencies returns the list of dependency IDs that are not "completed" or "skipped".
// Missing task IDs are treated as unmet.
func GetUnmetDependencies(indexPath string, taskID string) ([]string, error) {
	index, err := LoadIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("load index: %w", err)
	}

	_, foundTask, err := FindTask(index, taskID)
	if err != nil {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	var unmet []string
	for _, dep := range foundTask.Dependencies {
		if strings.HasSuffix(dep, IDSuffixWildcard) {
			prefix := strings.TrimSuffix(dep, IDSuffixWildcard)
			prefixWithDot := prefix + "."
			for _, other := range index.tasks {
				if other.ID == foundTask.ID {
					continue
				}
				if strings.HasPrefix(other.ID, prefixWithDot) && IsBusinessTask(other.ID) && other.Status != "completed" && other.Status != "skipped" {
					unmet = append(unmet, other.ID)
				}
			}
			continue
		}
		if depTask, exists := index.ByID(dep); !exists {
			unmet = append(unmet, dep)
		} else if depTask.Status != "completed" && depTask.Status != "skipped" {
			unmet = append(unmet, dep)
		}
	}
	return unmet, nil
}
