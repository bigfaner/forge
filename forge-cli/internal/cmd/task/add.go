package task

import (
	"errors"
	"fmt"
	"forge-cli/internal/cmd/base"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgelog"
	"forge-cli/pkg/project"
	"forge-cli/pkg/proposal"
	"forge-cli/pkg/task"
	"forge-cli/pkg/types"

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
	addVars          []string
	addSourceTaskID  string
	addBlockSource   bool
	addType          string
)

var addCmd = &cobra.Command{
	Use:   "add --title TITLE [flags]",
	Short: "Add a new task to the current feature",
	Long: `Add a new task dynamically. Validates inputs and writes files.
The CLI is a pure tool — the caller decides what to add.`,
	Args: cobra.NoArgs,
	RunE: runAdd,
}

func init() {
	addCmd.Flags().StringVar(&addTitle, "title", "", "Task title (required)")
	addCmd.Flags().StringVar(&addID, "id", "", "Custom task ID (auto-generated as disc-N if omitted)")
	addCmd.Flags().StringVar(&addPriority, "priority", string(types.PriorityP1), "Task priority: P0, P1, or P2")
	addCmd.Flags().StringVar(&addDependsOn, "depends-on", "", "Comma-separated dependency task IDs")
	addCmd.Flags().StringVar(&addEstimatedTime, "estimated-time", "", "Time estimate (e.g. \"1-2h\")")
	addCmd.Flags().BoolVar(&addBreaking, "breaking", false, "Mark as breaking (triggers full test suite)")
	addCmd.Flags().StringVar(&addDescription, "description", "", "Task description (markdown body)")
	addCmd.Flags().StringArrayVar(&addVars, "var", nil, "Template variable in key=value format (repeatable)")
	addCmd.Flags().StringVar(&addSourceTaskID, "source-task-id", "", "Source task ID: auto-resolves to root ancestor, injects {{SOURCE_TASK_ID}}, and adds this task as source dependency")
	addCmd.Flags().BoolVar(&addBlockSource, "block-source", false, "Set source task to blocked (before resolution, preserves fix-chain)")
	addCmd.Flags().StringVar(&addType, "type", "", "Task type (e.g. feature, enhancement, cleanup, refactor, fix, documentation). Validated against known types.")
}

// AddResult holds the result of a successful add operation.
type AddResult struct {
	ID            string
	Title         string
	Priority      string
	Status        string
	File          string
	Record        string
	Breaking      bool
	Dependencies  []string
	FeatureSlug   string
	ProjectRoot   string
	SourceBlocked string // source task ID that was blocked (empty if --block-source not used)
}

func runAdd(cmd *cobra.Command, _ []string) error {
	result, err := executeAdd(cmd)
	if err != nil {
		var dedupErr *task.ActiveFixExistsError
		if errors.As(err, &dedupErr) {
			base.PrintBlockStart()
			base.PrintField("ACTION", "SKIPPED")
			base.PrintField("REASON", dedupErr.Error())
			base.PrintBlockEnd()
			return nil
		}
		base.Exit(err)
	}

	base.PrintBlockStart()
	base.PrintField("ACTION", "ADDED")
	base.PrintField("KEY", result.ID)
	base.PrintField("TASK_ID", result.ID)
	base.PrintField("TITLE", result.Title)
	base.PrintField("PRIORITY", result.Priority)
	base.PrintField("STATUS", result.Status)
	base.PrintField("FILE", result.File)
	base.PrintField("RECORD", result.Record)
	if result.Breaking {
		base.PrintField("BREAKING", "true")
	}
	if result.SourceBlocked != "" {
		base.PrintField("SOURCE_BLOCKED", result.SourceBlocked)
	}
	base.PrintFieldIfNotEmptySlice("DEPENDENCIES", result.Dependencies)
	base.PrintBlockEnd()
	return nil
}

func executeAdd(cmd *cobra.Command) (*AddResult, error) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil, base.ErrProjectNotFound()
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return nil, base.ErrFeatureNotSet()
	}

	// Validate title
	if addTitle == "" {
		return nil, base.ErrNoInput("title is required")
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

	// Validate type
	if addType != "" && !task.IsValidType(addType) {
		return nil, base.ErrNoInput(fmt.Sprintf("invalid type %q. Run 'forge list-types' to see valid types", addType))
	}

	// Parse --var key=value flags
	vars := make(map[string]string)
	for _, v := range addVars {
		parts := strings.SplitN(v, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, base.ErrNoInput(fmt.Sprintf("invalid --var format: %q (expected key=value)", v))
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
		Vars:          vars,
		SourceTaskID:  addSourceTaskID,
		BlockSource:   addBlockSource,
		Type:          addType,
	}

	// When --type is specified, apply template defaults independently of
	// template file existence. Defaults (IDPrefix, Priority, etc.) live in
	// a hardcoded map and should not depend on the embedded FS file being found.
	if addType != "" {
		if defs, err := task.GetTaskTemplateDefaults(addType); err == nil {
			changed := func(name string) bool { return cmd != nil && cmd.Flags().Changed(name) }
			if !changed("priority") {
				opts.Priority = defs.Priority
			}
			if !changed("breaking") {
				opts.Breaking = defs.Breaking
			}
			if !changed("estimated-time") && defs.EstimatedTime != "" {
				opts.EstimatedTime = defs.EstimatedTime
			}
			if !changed("id") && defs.IDPrefix != "" {
				opts.IDPrefix = defs.IDPrefix
			}
		}
		// Template file is only needed for markdown generation
		if _, err := task.GetTaskTemplate(addType); err == nil {
			opts.Template = addType
		}
	}

	id, err := task.AddTask(indexPath, opts)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "task ID already exists") {
			return nil, base.ErrTaskIDConflict(addID)
		}
		if strings.Contains(errMsg, "dependency not found") {
			return nil, base.ErrInvalidDependency(deps)
		}
		return nil, fmt.Errorf("failed to add task: %w", err)
	}

	// Create task markdown file
	opts.ID = id // ensure ID is set after potential auto-generation
	if err := task.CreateTaskMarkdown(tasksDir, id+".md", opts); err != nil {
		return nil, fmt.Errorf("create task file: %w", err)
	}

	// Rebuild index from all .md files (canonical merge)

	// Resolve intent from proposal (defaults to "new-feature" when proposal missing)
	var addIntent string
	if p, err := proposal.FindBySlug(projectRoot, featureSlug); err == nil && p.Intent != "" {
		addIntent = p.Intent
	}

	buildOpts := task.BuildIndexOpts{
		FeatureSlug: featureSlug,
		ProjectRoot: projectRoot,
		TasksDir:    tasksDir,
		IndexPath:   indexPath,
		Intent:      addIntent,
	}
	if _, err := task.BuildIndex(buildOpts); err != nil {
		forgelog.Warn("WARNING: failed to rebuild index: %v\n", err)
	}

	// Reset forge state so claim loop continues
	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		forgelog.Warn("WARNING: failed to update .forge/state.json: %v\n", err)
	}

	// Report which source was blocked (if --block-source was used)
	var sourceBlocked string
	if addBlockSource && addSourceTaskID != "" {
		sourceBlocked = addSourceTaskID
	}

	return &AddResult{
		ID:            id,
		Title:         addTitle,
		Priority:      opts.Priority,
		Status:        string(types.StatusPending),
		File:          filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, id+".md")),
		Record:        filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, "records/"+id+".md")),
		Breaking:      opts.Breaking,
		Dependencies:  deps,
		FeatureSlug:   featureSlug,
		ProjectRoot:   projectRoot,
		SourceBlocked: sourceBlocked,
	}, nil
}
