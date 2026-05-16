package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/profile"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"
	tmpl "forge-cli/pkg/template"

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
	addBlockSource   bool
	addType          string
)

var addCmd = &cobra.Command{
	Use:   "add --title TITLE [--id ID] [--type TYPE] [--priority P0|P1|P2] [--depends-on ID,ID] [--breaking] [--description TEXT] [--template NAME] [--var KEY=VALUE]",
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

func runAdd(cmd *cobra.Command, _ []string) {
	result, err := executeAdd(cmd)
	if err != nil {
		var dedupErr *task.ActiveFixExistsError
		if errors.As(err, &dedupErr) {
			PrintBlockStart()
			PrintField("ACTION", "SKIPPED")
			PrintField("REASON", dedupErr.Error())
			PrintBlockEnd()
			return
		}
		Exit(err)
	}

	PrintBlockStart()
	PrintField("ACTION", "ADDED")
	PrintField("KEY", result.ID)
	PrintField("TASK_ID", result.ID)
	PrintField("TITLE", result.Title)
	PrintField("PRIORITY", result.Priority)
	PrintField("STATUS", result.Status)
	PrintField("FILE", result.File)
	PrintField("RECORD", result.Record)
	if result.Breaking {
		PrintField("BREAKING", "true")
	}
	if result.SourceBlocked != "" {
		PrintField("SOURCE_BLOCKED", result.SourceBlocked)
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

	// Validate type
	if addType != "" && !task.ValidTypes[addType] {
		return nil, ErrNoInput(fmt.Sprintf("invalid type %q. Run 'forge list-types' to see valid types", addType))
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
		BlockSource:   addBlockSource,
		Type:          addType,
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
			if !cmd.Flags().Changed("id") && defs.IDPrefix != "" {
				opts.IDPrefix = defs.IDPrefix
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

	// Rebuild index from all .md files (canonical merge)
	profiles, _ := profile.ReadTestProfiles(projectRoot)
	resolveStrategy := func(profileName, kind string) []byte {
		content, _ := profile.GetStrategy(profileName, kind)
		return content
	}
	buildOpts := task.BuildIndexOpts{
		FeatureSlug:     featureSlug,
		ProjectRoot:     projectRoot,
		TasksDir:        tasksDir,
		IndexPath:       indexPath,
		TestProfiles:    profiles,
		ResolveStrategy: resolveStrategy,
	}
	if _, err := task.BuildIndex(buildOpts); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to rebuild index: %v\n", err)
	}

	// Reset forge state so claim loop continues
	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to update .forge/state.json: %v\n", err)
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
		Status:        "pending",
		File:          filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, id+".md")),
		Record:        filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, "records/"+id+".md")),
		Breaking:      opts.Breaking,
		Dependencies:  deps,
		FeatureSlug:   featureSlug,
		ProjectRoot:   projectRoot,
		SourceBlocked: sourceBlocked,
	}, nil
}
