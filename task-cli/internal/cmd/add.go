package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"
	"task-cli/pkg/task"
	tmpl "task-cli/pkg/template"

	"github.com/spf13/cobra"
)

var (
	addTitle         string
	addID            string
	addPriority      string
	addDependsOn     string
	addEstimatedTime string
	addBreaking      bool
	addDescription   string
	addTemplate      string
	addVars          []string
	addSourceTaskID  string
)

var addCmd = &cobra.Command{
	Use:   "add --title TITLE [--id ID] [--priority P0|P1|P2] [--depends-on ID,ID] [--breaking] [--description TEXT] [--template NAME] [--var KEY=VALUE]",
	Short: "Add a new task to the current feature",
	Long: `Add a new task dynamically. Validates inputs and writes files.
The CLI is a pure tool — the caller decides what to add.`,
	Run: runAdd,
}

func init() {
	addCmd.Flags().StringVar(&addTitle, "title", "", "Task title (required)")
	addCmd.Flags().StringVar(&addID, "id", "", "Custom task ID (auto-generated as disc-N if omitted)")
	addCmd.Flags().StringVar(&addPriority, "priority", "P1", "Task priority: P0, P1, or P2")
	addCmd.Flags().StringVar(&addDependsOn, "depends-on", "", "Comma-separated dependency task IDs")
	addCmd.Flags().StringVar(&addEstimatedTime, "estimated-time", "", "Time estimate (e.g. \"1-2h\")")
	addCmd.Flags().BoolVar(&addBreaking, "breaking", false, "Mark as breaking (triggers full test suite)")
	addCmd.Flags().StringVar(&addDescription, "description", "", "Task description (markdown body)")
	addCmd.Flags().StringVar(&addTemplate, "template", "", "Template name (reads from tasks/_templates/<name>.md)")
	addCmd.Flags().StringArrayVar(&addVars, "var", nil, "Template variable in key=value format (repeatable)")
	addCmd.Flags().StringVar(&addSourceTaskID, "source-task-id", "", "Source task ID: auto-injects {{SOURCE_TASK_ID}} and adds this task as source dependency")
}

// AddResult holds the result of a successful add operation.
type AddResult struct {
	ID           string
	Title        string
	Priority     string
	Status       string
	File         string
	Record       string
	Breaking     bool
	Dependencies []string
	FeatureSlug  string
	ProjectRoot  string
}

func runAdd(cmd *cobra.Command, args []string) {
	result, err := executeAdd(cmd)
	if err != nil {
		Exit(err)
	}

	PrintBlockStart()
	PrintField("ACTION", "ADDED")
	PrintField("KEY", result.ID)
	PrintField("ID", result.ID)
	PrintField("TITLE", result.Title)
	PrintField("PRIORITY", result.Priority)
	PrintField("STATUS", result.Status)
	PrintField("FILE", result.File)
	PrintField("RECORD", result.Record)
	if result.Breaking {
		PrintField("BREAKING", "true")
	}
	PrintFieldIfNotEmptySlice("DEPENDENCIES", result.Dependencies)
	PrintBlockEnd()
}

func executeAdd(cmd *cobra.Command) (*AddResult, error) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil, ErrProjectNotFound()
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return nil, ErrFeatureNotSet()
	}

	// Validate title
	if addTitle == "" {
		return nil, ErrNoInput("title is required")
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	// Parse dependencies
	var deps []string
	if addDependsOn != "" {
		for _, d := range strings.Split(addDependsOn, ",") {
			d = strings.TrimSpace(d)
			if d != "" {
				deps = append(deps, d)
			}
		}
	}

	// Parse --var key=value flags
	vars := make(map[string]string)
	for _, v := range addVars {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, ErrNoInput(fmt.Sprintf("invalid --var format: %q (expected key=value)", v))
		}
		vars[parts[0]] = parts[1]
	}

	opts := task.AddTaskOpts{
		ID:            addID,
		Title:         addTitle,
		Priority:      addPriority,
		EstimatedTime: addEstimatedTime,
		Dependencies:  deps,
		Breaking:      addBreaking,
		Description:   addDescription,
		Template:      addTemplate,
		Vars:          vars,
		SourceTaskID:  addSourceTaskID,
	}

	// Apply template defaults for fixed fields
	if addTemplate != "" {
		if _, err := tmpl.Get(addTemplate); err != nil {
			return nil, ErrNoInput(fmt.Sprintf("template %q not found. Available: %v", addTemplate, tmpl.List()))
		}
		defs, err := tmpl.GetDefaults(addTemplate)
		if err == nil {
			if !cmd.Flags().Changed("priority") {
				opts.Priority = defs.Priority
			}
			if !cmd.Flags().Changed("breaking") {
				opts.Breaking = defs.Breaking
			}
			if !cmd.Flags().Changed("estimated-time") && defs.EstimatedTime != "" {
				opts.EstimatedTime = defs.EstimatedTime
			}
		}
	}

	id, err := task.AddTask(indexPath, opts)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "task ID already exists") {
			return nil, ErrTaskIDConflict(addID)
		}
		if strings.Contains(errMsg, "dependency not found") {
			return nil, ErrInvalidDependency(deps)
		}
		return nil, fmt.Errorf("failed to add task: %w", err)
	}

	// Create task markdown file
	opts.ID = id // ensure ID is set after potential auto-generation
	if err := task.CreateTaskMarkdown(tasksDir, id+".md", opts); err != nil {
		return nil, fmt.Errorf("create task file: %w", err)
	}

	// Reset forge state so claim loop continues
	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to update .forge/state.json: %v\n", err)
	}

	return &AddResult{
		ID:           id,
		Title:        addTitle,
		Priority:     opts.Priority,
		Status:       "pending",
		File:         filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, id+".md")),
		Record:       filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, "records/"+id+".md")),
		Breaking:     opts.Breaking,
		Dependencies: deps,
		FeatureSlug:  featureSlug,
		ProjectRoot:  projectRoot,
	}, nil
}
