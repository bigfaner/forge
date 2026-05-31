package qualitygate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/just"
	"forge-cli/pkg/task"
	"forge-cli/pkg/types"
)

// CountFixTasks counts active (non-terminal) fix-tasks for a step.
// A fix-task is identified by having a title with the prefix "fix <step>:".
// Terminal statuses (completed, rejected, skipped) are excluded from the count.
// This ensures the fix-task cap reflects work-in-progress only.
func CountFixTasks(index *task.TaskIndex, step string) int {
	count := 0
	prefix := "fix " + step + ":"
	for _, t := range index.TasksMap() {
		if !strings.HasPrefix(t.Title, prefix) {
			continue
		}
		// Exclude terminal statuses
		if t.Status == types.StatusCompleted || t.Status == types.StatusRejected || t.Status == types.StatusSkipped {
			continue
		}
		count++
	}
	return count
}

// fixTypeFromStep returns the deterministic task type for a quality gate failure step.
// compile/test failures -> TypeCodingFix, fmt/lint failures -> TypeCodingCleanup.
func fixTypeFromStep(step string) string {
	switch step {
	case "compile", "unit-test", "test":
		return task.TypeCodingFix
	case "fmt", "lint":
		return task.TypeCodingCleanup
	default:
		return task.TypeCodingFix
	}
}

// AddFixTask creates fix tasks grouped by test suite (directory) using the same
// internal API as `forge task add`. Source files are extracted from the output and
// grouped by directory. Each directory group becomes a separate fix-task, enabling
// parallel execution and bounded scope. Returns the first task ID on success.
// Returns ("", error) on failure: template not found, task add failure, markdown creation failure, or cap exceeded.
func AddFixTask(projectRoot, featureSlug, step, output, errorDocPath string) (string, error) {
	sourceFiles := ExtractSourceFiles(output)

	// Group source files by directory for parallel execution.
	groups := GroupFilesByDir(sourceFiles)

	// If no meaningful groups (e.g. "See error output" fallback), create a single task.
	if len(groups) == 0 {
		return addSingleFixTask(projectRoot, featureSlug, step, sourceFiles, output, errorDocPath)
	}

	// One task per directory group.
	var firstID string
	for _, group := range groups {
		id, err := addSingleFixTask(projectRoot, featureSlug, step, group, output, errorDocPath)
		if err != nil {
			// If a group fails (e.g. cap reached), return the error.
			// Already-created tasks remain in the index.
			return firstID, err
		}
		if firstID == "" {
			firstID = id
		}
	}
	return firstID, nil
}

// GroupFilesByDir splits comma-separated source files into groups by directory.
// Files in the same directory stay in one group. Returns nil if files is empty
// or the fallback message.
func GroupFilesByDir(files string) []string {
	if files == "" || strings.HasPrefix(files, "See error") {
		return nil
	}

	dirMap := make(map[string][]string)
	for _, f := range strings.Split(files, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		dir := filepath.Dir(f)
		dirMap[dir] = append(dirMap[dir], f)
	}

	if len(dirMap) <= 1 {
		// All files in one directory (or no files) — single group.
		return nil
	}

	var groups []string
	for _, dirFiles := range dirMap {
		groups = append(groups, strings.Join(dirFiles, ", "))
	}
	return groups
}

// conciseError extracts all "--- FAIL:" lines from output for test failures.
// Falls back to ExtractConciseError (tail 10 lines) when no "--- FAIL:" lines
// exist (e.g. compile/fmt/lint steps), ensuring the description is never empty.
func conciseError(output string) string {
	if failLines := just.ExtractFailLines(output); failLines != "" {
		return failLines
	}
	return just.ExtractConciseError(output, 10)
}

// fixTaskOverride holds optional overrides for createFixTask.
// When a field is set, it replaces the default value derived from step/output.
// This allows callers (e.g. addRegressionFixTasks) to customize the task
// while reusing the shared creation logic.
type fixTaskOverride struct {
	// Title overrides the default title. When empty, the default
	// "fix <step>: just <step> failure in quality gate" is used.
	Title string
	// Description overrides the default description. When empty, the default
	// description with error doc path and concise error is used.
	Description string
	// ExtraVars are merged into the default Vars map, overwriting keys
	// with the same name (e.g. SOURCE_FILES, TEST_SCRIPT).
	ExtraVars map[string]string
}

// createFixTask is the shared helper for creating a single fix task.
// It encapsulates: surface inference, opts construction, AddTask,
// CreateTaskMarkdown, and EnsureForgeState.
//
// The caller provides the core parameters (projectRoot, featureSlug, step,
// sourceFiles, output, errorDocPath) and optional overrides via fixTaskOverride.
// When no overrides are provided, defaults are derived from step and output
// (matching the original addSingleFixTask behavior).
//
// Returns (taskID, nil) on success.
// Returns ("", error) on failure: template not found, task add failure,
// or markdown creation failure.
func createFixTask(projectRoot, featureSlug, step, sourceFiles, output, errorDocPath string, overrides ...fixTaskOverride) (string, error) {
	var ov fixTaskOverride
	if len(overrides) > 0 {
		ov = overrides[0]
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	// Surface inference with soft-failure policy.
	// When surfaces are not configured or no match is found, fix-task creation
	// proceeds with empty surface key/type.
	surfaceKey, surfaceType := inferSurface(projectRoot, sourceFiles)
	if surfaceKey == "" && surfaceType == "" {
		fmt.Fprintln(os.Stderr, "WARNING: surface inference failed: no surfaces configured or no match for source files. Run 'forge surfaces detect' to configure surfaces")
	}

	testScript := "just " + step

	// Derive title: use override if provided, otherwise default.
	title := ov.Title
	if title == "" {
		title = fmt.Sprintf("fix %s: %s failure in quality gate", step, testScript)
	}

	// Derive description: use override if provided, otherwise default.
	description := ov.Description
	if description == "" {
		description = fmt.Sprintf(
			"Quality gate step `%s` failed during quality-gate hook.\n\n"+
				"Error output saved to: `%s`\n\n"+
				"Concise error:\n```\n%s\n```",
			testScript, errorDocPath, conciseError(output),
		)
	}

	// Build opts — Priority/Breaking/EstimatedTime sourced from template defaults
	// when available (dual-source truth: opts is authoritative, template provides defaults).
	// SourceTaskID is deliberately empty (project-wide gate has no source task).
	// Vars["SOURCE_TASK_ID"] is "N/A (project-wide gate)" for template rendering.
	taskType := fixTypeFromStep(step)

	// Derive Breaking and EstimatedTime from template defaults (dual-source truth).
	// For coding.cleanup (fmt/lint failures): Breaking=false, EstimatedTime="15min".
	// For coding.fix (compile/test failures): Breaking=true, EstimatedTime="30min".
	breaking := true
	estimatedTime := "30min"
	if defs, err := task.GetTaskTemplateDefaults(taskType); err == nil {
		breaking = defs.Breaking
		estimatedTime = defs.EstimatedTime
	}

	// Build default vars, then merge extra vars from overrides.
	vars := map[string]string{
		"SOURCE_FILES":   sourceFiles,
		"TEST_SCRIPT":    testScript,
		"TEST_RESULTS":   errorDocPath,
		"SOURCE_TASK_ID": "N/A (project-wide gate)",
	}
	for k, v := range ov.ExtraVars {
		vars[k] = v
	}

	opts := task.AddTaskOpts{
		Title:         title,
		Priority:      string(types.PriorityP0),
		EstimatedTime: estimatedTime,
		Breaking:      breaking,
		Description:   description,
		Template:      taskType,
		Type:          taskType,
		SurfaceKey:    surfaceKey,
		SurfaceType:   surfaceType,
		Vars:          vars,
	}

	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))

	if _, err := task.GetTaskTemplate(opts.Template); err != nil {
		return "", fmt.Errorf("template %q not found: %w", opts.Template, err)
	}
	if defs, err := task.GetTaskTemplateDefaults(opts.Template); err == nil && defs.IDPrefix != "" {
		opts.IDPrefix = defs.IDPrefix
	}

	id, err := task.AddTask(indexPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to add fix task: %w", err)
	}

	opts.ID = id

	if err := task.CreateTaskMarkdown(tasksDir, id+".md", opts); err != nil {
		return "", fmt.Errorf("failed to create fix task file %s: %w", id+".md", err)
	}

	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to update .forge/state.json: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "Fix task %s added (P0, breaking=%v)\n", id, breaking)
	return id, nil
}

// addRegressionFixTasks creates independent fix tasks for each failing test file
// extracted from the regression output. It uses extractFileLineMap to identify primary
// test files (those with direct --- FAIL: entries) and creates one task per file with
// the relevant filtered output lines.
//
// Soft cap: up to 10 independent tasks + 1 overflow task (files 11-N merged).
// Total tasks <= 11, NOT subject to maxFixTasksPerStep hard cap.
//
// Fallback: when extractFileLineMap returns an empty map (no test files identified),
// falls back to the existing addFixTask behavior with a structured log warning.
//
// Returns the first task ID on success, or ("", error) on failure.
//
//nolint:unparam // errorDocPath is parameterized for API consistency with addFixTask and future callers.
func addRegressionFixTasks(projectRoot, featureSlug, output, errorDocPath string) (string, error) {
	fileLineMap := extractFileLineMap(output)

	// Fallback: no test files identified — use directory-grouped fix task.
	if len(fileLineMap) == 0 {
		fmt.Fprintln(os.Stderr, "WARNING: isTestFile returned zero matches for output, falling back to directory-grouped fix task")
		return AddFixTask(projectRoot, featureSlug, "test", output, errorDocPath)
	}

	// Sort files for deterministic ordering.
	files := make([]string, 0, len(fileLineMap))
	for f := range fileLineMap {
		files = append(files, f)
	}
	sort.Strings(files)

	var firstID string

	// Create tasks for first 10 files (or all if <= 10).
	primaryCount := len(files)
	if primaryCount > 10 {
		primaryCount = 10
	}

	for i := 0; i < primaryCount; i++ {
		file := files[i]
		lines := fileLineMap[file]
		title := fmt.Sprintf("fix test: %s failure in quality gate", file)
		description := fmt.Sprintf(
			"Quality gate test regression failed.\n\n"+
				"Test file: `%s`\n\n"+
				"Error output saved to: `%s`\n\n"+
				"Relevant output lines:\n```\n%s\n```",
			file, errorDocPath, strings.Join(lines, "\n"),
		)

		id, err := createFixTask(
			projectRoot, featureSlug, "test", file, output, errorDocPath,
			fixTaskOverride{
				Title:       title,
				Description: description,
				ExtraVars: map[string]string{
					"SOURCE_FILES": file,
					"TEST_SCRIPT":  "just test",
				},
			},
		)
		if err != nil {
			return firstID, err
		}
		if firstID == "" {
			firstID = id
		}
	}

	// Overflow: merge remaining files into one task.
	if len(files) > 10 {
		overflowFiles := files[10:]
		var overflowLines []string
		for _, f := range overflowFiles {
			overflowLines = append(overflowLines, fileLineMap[f]...)
		}
		overflowCount := len(overflowFiles)
		title := fmt.Sprintf("fix test: regression overflow (%d files)", overflowCount)
		description := fmt.Sprintf(
			"Quality gate test regression failed — overflow group.\n\n"+
				"Files: %s\n\n"+
				"Error output saved to: `%s`\n\n"+
				"Relevant output lines:\n```\n%s\n```",
			strings.Join(overflowFiles, ", "), errorDocPath, strings.Join(overflowLines, "\n"),
		)

		id, err := createFixTask(
			projectRoot, featureSlug, "test",
			strings.Join(overflowFiles, ", "), output, errorDocPath,
			fixTaskOverride{
				Title:       title,
				Description: description,
				ExtraVars: map[string]string{
					"SOURCE_FILES": strings.Join(overflowFiles, ", "),
					"TEST_SCRIPT":  "just test",
				},
			},
		)
		if err != nil {
			return firstID, err
		}
		if firstID == "" {
			firstID = id
		}
	}

	return firstID, nil
}

// addSingleFixTask creates a single fix task using the same internal API as `forge task add`.
// It performs a cap check (countFixTasks + maxFixTasksPerStep) and delegates to createFixTask.
// Returns (taskID, nil) on success.
// Returns ("", error) on failure: template not found, task add failure, markdown creation failure, or cap exceeded.
func addSingleFixTask(projectRoot, featureSlug, step, sourceFiles, output, errorDocPath string) (string, error) {
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	// Check cap before creating a new fix-task.
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to load index for cap check: %v\n", err)
		// Proceed without cap check if index can't be loaded.
	} else {
		active := CountFixTasks(index, step)
		if active >= maxFixTasksPerStep {
			fmt.Fprintf(os.Stderr, "max fix-tasks reached for %s, manual intervention required\n", step)
			return "", ErrMaxFixTasks
		}
	}

	return createFixTask(projectRoot, featureSlug, step, sourceFiles, output, errorDocPath)
}
