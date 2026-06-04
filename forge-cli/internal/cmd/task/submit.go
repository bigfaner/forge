package task

import (
	"encoding/json"
	"errors"
	"fmt"
	"forge-cli/internal/cmd/base"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"forge-cli/pkg/feature"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/just"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"
	"forge-cli/pkg/types"

	"github.com/spf13/cobra"
)

var (
	submitDataPath string
	submitJSON     bool
	submitQuiet    bool
)

var submitCmd = &cobra.Command{
	Use:   "submit <task-id>",
	Short: "Submit task execution result",
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
	RunE: runSubmit,
}

func init() {
	submitCmd.Flags().StringVar(&submitDataPath, "data", "", "Path to JSON data file")
	submitCmd.Flags().BoolVar(&submitJSON, "json", false, "Output result as JSON")
	submitCmd.Flags().BoolVar(&submitQuiet, "quiet", false, "Minimal output")
}

func runSubmit(_ *cobra.Command, args []string) error {
	taskIDArg := args[0]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return base.ErrProjectNotFound()
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return base.ErrFeatureNotSet()
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	// Wrap all index-modifying logic in WithLock for concurrent write safety
	if lockErr := indexPkg.WithLock(indexPath, func() error {
		return doSubmit(projectRoot, featureSlug, indexPath, taskIDArg)
	}); lockErr != nil {
		if errors.Is(lockErr, indexPkg.ErrLockConflict) {
			return base.NewAIError(base.ErrConflict, "Concurrent write conflict", "Retry the command", "", "")
		}
		if aiErr, ok := lockErr.(*base.AIError); ok {
			return aiErr
		}
		return base.NewAIError(base.ErrConflict, "Failed to acquire lock", lockErr.Error(), "", "")
	}
	return nil
}

// doSubmit contains the core submit logic, executed under the advisory lock.
func doSubmit(projectRoot, featureSlug, indexPath, taskIDArg string) error {
	idx, err := task.LoadIndex(indexPath)
	if err != nil {
		return base.ErrFileNotFound(indexPath)
	}

	key, t, err := task.FindTask(idx, taskIDArg)
	if err != nil {
		return base.ErrTaskNotFound(taskIDArg)
	}

	rd, err := readSubmitData(submitDataPath)
	if err != nil {
		return base.ErrNoInput(err.Error())
	}

	if rd.Status == "" {
		rd.Status = string(types.StatusCompleted)
	}

	// Validate taskId matches CLI arg if provided
	if rd.TaskID != "" && rd.TaskID != taskIDArg {
		return base.NewAIError(base.ErrValidation,
			fmt.Sprintf("taskId mismatch: JSON has %q, CLI arg is %q", rd.TaskID, taskIDArg),
			"The taskId in record.json does not match the task being recorded",
			"Either omit taskId from JSON or ensure it matches the CLI argument",
			fmt.Sprintf("Change taskId to %q or remove it from record.json", taskIDArg))
	}

	// Non-testable tasks: auto-set coverage=-1.0 to skip test evidence check
	if !task.IsTestableType(t.Type) {
		if rd.Coverage >= 0 && rd.TestsPassed == 0 && rd.TestsFailed == 0 {
			rd.Coverage = -1.0
		}
	}

	// Capture intended status before validateRecordData may auto-downgrade
	targetStatus := rd.Status

	// Validate required and recommended fields
	if err := validateRecordData(rd, t.Type); err != nil {
		return err
	}

	// State machine validation: check transition before proceeding
	if targetStatus == string(types.StatusCompleted) {
		if transitionErr := task.ValidateTransition(t.Status, types.StatusCompleted, task.RoleSubmit); transitionErr != nil {
			te := transitionErr.(*task.TransitionError)
			return base.NewErrInvalidTransition(string(t.Status), string(types.StatusCompleted), te.Msg)
		}
	}

	// Quality gate pre-check for completed tasks (non-testable types excluded)
	// Tiered model: breaking tasks run full gate (compile+fmt+lint+test),
	// non-breaking coding tasks run static gate (compile+fmt+lint).
	if targetStatus == string(types.StatusCompleted) && task.IsTestableType(t.Type) {
		validateQualityGate(projectRoot, t.SurfaceKey, t.Breaking)
	}

	// After validateRecordData, rd.Status may have been auto-downgraded (completed -> blocked)
	if rd.Status != targetStatus && rd.Status == string(types.StatusBlocked) {
		// Auto-downgrade: validate the blocked transition
		if transitionErr := task.ValidateTransition(t.Status, types.StatusBlocked, task.RoleSubmit); transitionErr != nil {
			te := transitionErr.(*task.TransitionError)
			return base.NewErrInvalidTransition(string(t.Status), string(types.StatusBlocked), te.Msg)
		}
	}

	// Validate status against index statusEnum
	if !slices.Contains(idx.StatusEnum, rd.Status) {
		return base.ErrInvalidStatus(rd.Status, idx.StatusEnum)
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
	if err := os.MkdirAll(filepath.Dir(recordPath), 0o755); err != nil {
		return base.NewAIError(base.ErrValidation, "Failed to create record directory", err.Error(), "Check directory permissions", "mkdir -p "+filepath.Dir(recordPath))
	}

	if err := os.WriteFile(recordPath, []byte(content), 0o644); err != nil {
		return base.NewAIError(base.ErrValidation, "Failed to write record file", err.Error(), "Check file permissions", "cat "+recordPath)
	}

	// Update task status in index
	t.Status = types.Status(rd.Status)

	// Set BlockedReason on auto-downgrade
	if rd.Status == string(types.StatusBlocked) && targetStatus == string(types.StatusCompleted) {
		t.BlockedReason = fmt.Sprintf("auto-downgrade: testsFailed=%d", rd.TestsFailed)
	}

	idx.SetTask(key, *t)

	// Auto-restore: if this fix-task completed or skipped, check if source can be unblocked
	if t.SourceTaskID != "" && (rd.Status == string(types.StatusCompleted) || rd.Status == string(types.StatusSkipped)) {
		autoRestoreSourceTask(idx, t.SourceTaskID)
	}

	if err := saveIndexAndSignalCompletion(indexPath, projectRoot, featureSlug, idx); err != nil {
		return err
	}

	if submitJSON {
		result := map[string]string{
			"recordFile": recordPath,
			"taskId":     t.ID,
			"status":     rd.Status,
		}
		data, _ := json.Marshal(result)
		fmt.Println(string(data))
	} else if !submitQuiet {
		base.PrintBlockStart()
		base.PrintField("STATUS", rd.Status)
		base.PrintBlockEnd()
	}

	return nil
}

// saveIndexAndSignalCompletion saves the index atomically and writes .forge/state.json
// if all tasks are completed or skipped (rejected does not count as done).
func saveIndexAndSignalCompletion(indexPath, projectRoot, featureSlug string, idx *task.TaskIndex) error {
	if err := indexPkg.SaveIndexAtomic(indexPath, idx); err != nil {
		return base.NewAIError(base.ErrConflict, "Failed to update task index", err.Error(), "Check index.json is writable", "cat "+indexPath)
	}

	allDone := true
	for _, t := range idx.TasksMap() {
		if t.Status != types.StatusCompleted && t.Status != types.StatusSkipped {
			allDone = false
			break
		}
	}
	if allDone {
		if err := feature.WriteForgeState(projectRoot, featureSlug); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: failed to write forge state: %v\n", err)
		}
	}
	return nil
}

// autoRestoreSourceTask checks if a blocked source task can be unblocked.
// If the source is blocked and ALL its dependencies are completed or skipped, restores it to pending.
// Root cause: must lookup by ID (iterate), not by direct map key, because map keys are slugs.
// Unknown deps (not in index) are skipped — consistent with claim.go checkDependenciesMet.
func autoRestoreSourceTask(index *task.TaskIndex, sourceTaskID string) {
	srcKey, srcTask, err := task.FindTask(index, sourceTaskID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "AUTO-RESTORE-SKIP: source task %s not found in index\n", sourceTaskID)
		return
	}
	if srcTask.Status != types.StatusBlocked {
		fmt.Fprintf(os.Stderr, "AUTO-RESTORE-SKIP: source task %s is %s (not blocked)\n", sourceTaskID, srcTask.Status)
		return
	}

	unmet := checkUnmetDeps(index, srcTask)
	// Filter out unknown deps — consistent with claim.go checkDependenciesMet which
	// treats unknown deps as vacuously satisfied.
	var knownUnmet []string
	var unknownDeps []string
	for _, id := range unmet {
		if _, found := index.ByID(id); !found {
			unknownDeps = append(unknownDeps, id)
			continue
		}
		knownUnmet = append(knownUnmet, id)
	}
	if len(knownUnmet) > 0 {
		fmt.Fprintf(os.Stderr, "AUTO-RESTORE-SKIP: source task %s has unmet deps: %v\n", sourceTaskID, knownUnmet)
		return
	}

	srcTask.Status = types.StatusPending
	index.SetTask(srcKey, *srcTask)
	fmt.Fprintf(os.Stderr, "AUTO-RESTORE: source task %s restored to pending (all deps completed or skipped)\n", sourceTaskID)
}

// readSubmitData delegates to pkg/task.ReadSubmitData.
// Kept as alias for internal callers and tests.
var readSubmitData = task.ReadSubmitData

// validateRecordData checks required and recommended fields in task.RecordData.
// taskType determines which checks apply based on category:
//   - coding: full validation (test evidence, testsFailed auto-downgrade)
//   - doc/test/validation/gate: skip test evidence and testsFailed checks
//   - eval: accept review fields (summary/findings/severity), reject pure test fields
//
// Hard-required fields (missing = error): summary (all categories)
// Auto-downgrade (coding only): completed + testsFailed > 0 -> blocked
// Hard validation for completed coding tasks:
//   - testsPassed=0 && testsFailed=0 with coverage >= 0
//   - any acceptanceCriteria with met=false
//
// Recommended fields for "completed" status (missing = warning, all categories):
//   - keyDecisions, acceptanceCriteria
func validateRecordData(rd *task.RecordData, taskType string) error {
	category := task.CategoryForType(taskType)
	isCoding := category == task.CategoryCoding
	isEval := category == task.CategoryEval

	var missing []string

	// Hard-required fields (all categories)
	if strings.TrimSpace(rd.Summary) == "" {
		missing = append(missing, "summary")
	}

	if len(missing) > 0 {
		return base.ErrMissingFields(missing)
	}

	// Auto-downgrade (coding only): completed with test failures -> blocked
	if isCoding && rd.Status == string(types.StatusCompleted) && rd.TestsFailed > 0 {
		fmt.Fprintf(os.Stderr, "---\nWARNING: %d test failures detected — auto-downgrading status from 'completed' to 'blocked'\nHINT: Fix test failures, then re-record with status 'completed'\n---\n", rd.TestsFailed)
		rd.Status = string(types.StatusBlocked)
	}

	if rd.Status != string(types.StatusCompleted) {
		return nil
	}

	// Hard validation for completed tasks (coding only)
	if isCoding {
		// Reject completed with no test evidence (unless coverage=-1.0 signals "no tests")
		if rd.Coverage >= 0 && rd.TestsPassed == 0 && rd.TestsFailed == 0 {
			return base.ErrNoTestEvidence()
		}
	}

	// Eval validation: reject completed with only pure test fields (testsPassed/coverage)
	if isEval {
		hasEvalFields := len(rd.Findings) > 0 || rd.Severity != "" || rd.Score > 0
		if !hasEvalFields {
			return base.NewAIError(base.ErrValidation,
				"eval task submission requires eval-specific fields (findings, severity, or score)",
				"Eval tasks are review-type tasks and should not use coding test metrics",
				"Include at least one of: findings, severity, score",
				`{"findings": ["issue1"], "severity": "major", "score": 850}`)
		}
	}

	// Reject completed with unmet acceptance criteria (all categories)
	if len(rd.AcceptanceCriteria) > 0 {
		var unmet []string
		for _, ac := range rd.AcceptanceCriteria {
			if !ac.Met {
				unmet = append(unmet, ac.Criterion)
			}
		}
		if len(unmet) > 0 {
			return base.ErrUnmetAcceptanceCriteria(unmet)
		}
	}

	// Recommended fields for completed tasks (coding only)
	if category == task.CategoryCoding {
		var recommended []string
		if len(rd.KeyDecisions) == 0 {
			recommended = append(recommended, "keyDecisions")
		}
		if len(rd.AcceptanceCriteria) == 0 {
			recommended = append(recommended, "acceptanceCriteria")
		}
		if len(recommended) > 0 {
			base.WarnMissingFields(recommended)
		}
	}

	return nil
}

// fillRecordTemplate delegates to pkg/task.RenderRecord.
// Kept as alias for internal callers and tests.
var fillRecordTemplate = task.RenderRecord

// validateQualityGate runs the quality gate based on the task's breaking flag.
// breaking=true: unit gate (compile -> fmt -> lint -> unit-test) for fast feedback.
// breaking=false: non-breaking gate (compile -> fmt -> lint), skipping tests.
// On failure, exits with base.AIError containing concise error output.
func validateQualityGate(projectRoot, scope string, breaking bool) {
	steps := just.NonBreakingGateSequence()
	if breaking {
		steps = just.UnitGateSequence()
	}
	just.RunGate(projectRoot, scope, steps, func(step, output string) {
		concise := just.ExtractConciseError(output, 10)
		panic(base.NewAIError(base.ErrValidation,
			fmt.Sprintf("Quality gate failed at step: just %s", step),
			concise,
			"Fix the errors above and re-run task record",
			"Or set status to 'blocked' and create a fix task"))
	})
}
