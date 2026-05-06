package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"task-cli/pkg/feature"
	"task-cli/pkg/just"
	"task-cli/pkg/project"
	"task-cli/pkg/task"

	"github.com/spf13/cobra"
)

var (
	recordDataPath string
	recordJSON     bool
	recordQuiet    bool
	recordForce    bool
)

var recordCmd = &cobra.Command{
	Use:   "record <task-id>",
	Short: "Generate a task execution record",
	Long: `Generate a task execution record from a Markdown template.

The record data can be provided via:
  - --data flag pointing to a JSON file
  - stdin pipe (may have issues on Windows)

The record is saved to docs/features/<slug>/records/<task.record>
and the task status is updated in index.json.

JSON input format:
  {
    "taskId": "1.1",
    "status": "completed",
    "summary": "Brief description of what was done",
    "filesCreated": ["path/to/new/file.go"],
    "filesModified": ["path/to/modified/file.go"],
    "keyDecisions": ["Decision 1", "Decision 2"],
    "testsPassed": 5,
    "testsFailed": 0,
    "coverage": 85.5,
    "acceptanceCriteria": [
      {"criterion": "Feature works", "met": true}
    ],
    "notes": "Optional notes"
  }

Required fields: summary
Optional fields: taskId (verified against CLI arg if provided), status (default: completed),
                 filesCreated, filesModified, keyDecisions, testsPassed,
                 testsFailed, coverage, acceptanceCriteria, notes`,
	Args: cobra.ExactArgs(1),
	Run:  runRecord,
}

func init() {
	recordCmd.Flags().StringVar(&recordDataPath, "data", "", "Path to JSON data file")
	recordCmd.Flags().BoolVar(&recordJSON, "json", false, "Output result as JSON")
	recordCmd.Flags().BoolVar(&recordQuiet, "quiet", false, "Minimal output")
	recordCmd.Flags().BoolVar(&recordForce, "force", false, "Override validation errors (use with caution)")
}

func runRecord(cmd *cobra.Command, args []string) {
	taskIDArg := args[0]

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

	key, t, err := findTask(index, taskIDArg)
	if err != nil {
		Exit(ErrTaskNotFound(taskIDArg))
	}

	rd, err := readRecordData(recordDataPath)
	if err != nil {
		Exit(ErrNoInput(err.Error()))
	}

	if rd.Status == "" {
		rd.Status = "completed"
	}

	// Validate taskId matches CLI arg if provided
	if rd.TaskID != "" && rd.TaskID != taskIDArg {
		Exit(NewAIError(ErrValidation,
			fmt.Sprintf("taskId mismatch: JSON has %q, CLI arg is %q", rd.TaskID, taskIDArg),
			"The taskId in record.json does not match the task being recorded",
			"Either omit taskId from JSON or ensure it matches the CLI argument",
			fmt.Sprintf("Change taskId to %q or remove it from record.json", taskIDArg)))
	}

	// Validate required and recommended fields
	validateRecordData(rd, recordForce)

	// Quality gate pre-check for completed tasks (unless --force)
	if rd.Status == "completed" && !recordForce {
		validateQualityGate(projectRoot, t.Scope)
	}

	// Validate status
	validStatus := false
	for _, s := range index.StatusEnum {
		if s == rd.Status {
			validStatus = true
			break
		}
	}
	if !validStatus {
		Exit(ErrInvalidStatus(rd.Status, index.StatusEnum))
	}

	// Read startedTime from task-state.json
	statePath := feature.GetTaskStatePath(projectRoot, featureSlug)
	state, _ := task.LoadState(statePath)
	startedTime := ""
	if state != nil && state.TaskID == t.ID {
		startedTime = state.StartedTime
	}

	content := fillRecordTemplate(t, rd, startedTime)

	// Write record file
	recordPath := filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, t.Record))
	if err := os.MkdirAll(filepath.Dir(recordPath), 0755); err != nil {
		Exit(NewAIError(ErrValidation, "Failed to create record directory", err.Error(), "Check directory permissions", "mkdir -p "+filepath.Dir(recordPath)))
	}
	if err := os.WriteFile(recordPath, []byte(content), 0644); err != nil {
		Exit(NewAIError(ErrValidation, "Failed to write record file", err.Error(), "Check file permissions", "cat "+recordPath))
	}

	// Update task status in index
	t.Status = rd.Status
	index.Tasks[key] = *t

	// Auto-restore: if this task has a SourceTaskID and completed or skipped, check if source can be unblocked
	if t.SourceTaskID != "" && (rd.Status == "completed" || rd.Status == "skipped") {
		autoRestoreSourceTask(index, t.SourceTaskID)
	}

	saveIndexAndSignalCompletion(indexPath, projectRoot, featureSlug, index)

	if recordJSON {
		result := map[string]string{
			"recordFile": recordPath,
			"taskId":     t.ID,
			"status":     rd.Status,
		}
		data, _ := json.Marshal(result)
		fmt.Println(string(data))
	} else if !recordQuiet {
		PrintBlockStart()
		PrintField("TASK_ID", t.ID)
		PrintField("RECORD_FILE", recordPath)
		PrintField("STATUS", rd.Status)
		PrintBlockEnd()
	}
}

// saveIndexAndSignalCompletion saves the index and writes .forge/state.json
// if all tasks are completed or skipped.
func saveIndexAndSignalCompletion(indexPath, projectRoot, featureSlug string, index *task.TaskIndex) {
	if err := task.SaveIndex(indexPath, index); err != nil {
		Exit(NewAIError(ErrConflict, "Failed to update task index", err.Error(), "Check index.json is writable", "cat "+indexPath))
	}

	allDone := true
	for _, t := range index.Tasks {
		if t.Status != feature.StatusCompleted && t.Status != feature.StatusSkipped {
			allDone = false
			break
		}
	}
	if allDone {
		if err := feature.WriteForgeState(projectRoot, featureSlug); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to write forge state: %v\n", err)
		}
	}
}

// autoRestoreSourceTask checks if a blocked source task can be unblocked.
// If the source is blocked and ALL its dependencies are completed or skipped, restores it to pending.
func autoRestoreSourceTask(index *task.TaskIndex, sourceTaskID string) {
	srcTask, found := index.Tasks[sourceTaskID]
	if !found || srcTask.Status != "blocked" {
		return
	}

	unmet := checkUnmetDeps(index, &srcTask)
	if len(unmet) > 0 {
		return
	}

	srcTask.Status = "pending"
	index.Tasks[sourceTaskID] = srcTask
	fmt.Fprintf(os.Stderr, "AUTO-RESTORE: source task %s restored to pending (all deps completed or skipped)\n", sourceTaskID)
}

func findTask(index *task.TaskIndex, taskID string) (string, *task.Task, error) {
	for key, t := range index.Tasks {
		if t.ID == taskID || key == taskID {
			return key, &t, nil
		}
	}
	return "", nil, fmt.Errorf("task not found: %s", taskID)
}

func readRecordData(dataPath string) (*task.RecordData, error) {
	var data []byte
	var err error

	if dataPath != "" {
		data, err = os.ReadFile(dataPath)
	} else {
		stat, _ := os.Stdin.Stat()
		if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
			return nil, fmt.Errorf("no input: provide --data flag or pipe JSON to stdin")
		}
		data, err = io.ReadAll(os.Stdin)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read record data: %w", err)
	}

	var rd task.RecordData
	if err := json.Unmarshal(data, &rd); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return &rd, nil
}

// validateRecordData checks required and recommended fields in task.RecordData.
// Hard-required fields (missing = error): summary
// Hard validation for completed tasks (overridable with --force):
//   - testsPassed=0 && testsFailed=0 with coverage >= 0
//   - any acceptanceCriteria with met=false
// Recommended fields for "completed" status (missing = warning):
//   - keyDecisions, acceptanceCriteria
func validateRecordData(rd *task.RecordData, force bool) {
	var missing []string

	// Hard-required fields
	if strings.TrimSpace(rd.Summary) == "" {
		missing = append(missing, "summary")
	}

	if len(missing) > 0 {
		Exit(ErrMissingFields(missing))
	}

	if rd.Status != "completed" {
		return
	}

	// Hard validation for completed tasks (unless --force)
	if !force {
		// Reject completed with no test evidence (unless coverage=-1.0 signals "no tests")
		if rd.Coverage >= 0 && rd.TestsPassed == 0 && rd.TestsFailed == 0 {
			Exit(ErrNoTestEvidence())
		}

		// Reject completed with unmet acceptance criteria
		if len(rd.AcceptanceCriteria) > 0 {
			var unmet []string
			for _, ac := range rd.AcceptanceCriteria {
				if !ac.Met {
					unmet = append(unmet, ac.Criterion)
				}
			}
			if len(unmet) > 0 {
				Exit(ErrUnmetAcceptanceCriteria(unmet))
			}
		}
	}

	// Recommended fields for completed tasks
	var recommended []string
	if len(rd.KeyDecisions) == 0 {
		recommended = append(recommended, "keyDecisions")
	}
	if len(rd.AcceptanceCriteria) == 0 {
		recommended = append(recommended, "acceptanceCriteria")
	}
	if len(recommended) > 0 {
		WarnMissingFields(recommended)
	}
}

func fillRecordTemplate(t *task.Task, rd *task.RecordData, startedTime string) string {
	status := rd.Status
	started := startedTime
	if started == "" {
		started = time.Now().Format("2006-01-02 15:04")
	}
	completed := time.Now().Format("2006-01-02 15:04")
	if status != "completed" {
		completed = "N/A"
	}

	// Calculate time spent
	timeSpent := ""
	startedT, err1 := time.Parse("2006-01-02 15:04", started)
	completedT, err2 := time.Parse("2006-01-02 15:04", completed)
	if err1 == nil && err2 == nil && completedT.After(startedT) {
		dur := completedT.Sub(startedT)
		timeSpent = formatDuration(dur)
	}

	notes := rd.Notes
	if notes == "" {
		notes = "无"
	}

	return fmt.Sprintf(`---
status: "%s"
started: "%s"
completed: "%s"
time_spent: "%s"
---

# Task Record: %s %s

## Summary
%s

## Changes

### Files Created
%s

### Files Modified
%s

### Key Decisions
%s

## Test Results
- **Passed**: %d
- **Failed**: %d
- **Coverage**: %s

## Acceptance Criteria
%s

## Notes
%s
`,
		status, started, completed, timeSpent,
		t.ID, t.Title,
		rd.Summary,
		formatList(rd.FilesCreated),
		formatList(rd.FilesModified),
		formatList(rd.KeyDecisions),
		rd.TestsPassed, rd.TestsFailed, formatCoverage(rd.Coverage),
		formatCriteria(rd.AcceptanceCriteria),
		notes,
	)
}

func formatCoverage(c float64) string {
	if c < 0 {
		return "N/A (task has no tests)"
	}
	return fmt.Sprintf("%.1f%%", c)
}

func formatList(items []string) string {
	if len(items) == 0 {
		return "无"
	}
	lines := make([]string, len(items))
	for i, item := range items {
		lines[i] = "- " + item
	}
	return strings.Join(lines, "\n")
}

func formatDuration(dur time.Duration) string {
	d := int(dur.Hours())
	m := int(dur.Minutes()) % 60
	switch {
	case d > 0 && m > 0:
		return fmt.Sprintf("~%dh %dm", d, m)
	case d > 0:
		return fmt.Sprintf("~%dh", d)
	default:
		return fmt.Sprintf("~%dm", m)
	}
}

func formatCriteria(criteria []task.AcceptanceCriterion) string {
	if len(criteria) == 0 {
		return "无"
	}
	lines := make([]string, len(criteria))
	for i, c := range criteria {
		check := "[ ]"
		if c.Met {
			check = "[x]"
		}
		lines[i] = "- " + check + " " + c.Criterion
	}
	return strings.Join(lines, "\n")
}

// validateQualityGate runs the full quality gate (compile -> fmt -> lint -> test).
// On failure, exits with AIError containing concise error output.
func validateQualityGate(projectRoot, scope string) {
	just.RunGate(projectRoot, scope, just.DefaultGateSequence(), func(step, output string) {
		concise := just.ExtractConciseError(output, 10)
		Exit(NewAIError(ErrValidation,
			fmt.Sprintf("Quality gate failed at step: just %s", step),
			concise,
			"Fix the errors above and re-run task record",
			"Or use --force to bypass the quality gate"))
	})
}
