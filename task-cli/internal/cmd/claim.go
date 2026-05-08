package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"task-cli/pkg/feature"
	"task-cli/pkg/project"
	"task-cli/pkg/task"

	"github.com/spf13/cobra"
)

var claimCmd = &cobra.Command{
	Use:   "claim",
	Short: "Claim the next available task",
	Long: `Claim the next available task from the current feature's task list.

The task is selected based on:
1. All dependencies must be met
2. Priority (P0 > P1 > P2)
3. Task ID (semantic version ordering)`,
	Run: runClaim,
}

func runClaim(cmd *cobra.Command, args []string) {
	result, err := executeClaim()
	if err != nil {
		Exit(err)
	}

	// Print result
	if result.Action == "CONTINUE" {
		printContinueTask(result.State, result.Task, result.ProjectRoot, result.FeatureSlug)
	} else {
		printNewTask(result.Key, result.Task, result.ProjectRoot, result.FeatureSlug)
	}
}

// ClaimResult represents the result of a claim operation
type ClaimResult struct {
	Action      string // "CLAIMED" or "CONTINUE"
	Key         string
	Task        *task.Task
	State       *task.TaskState
	StartedTime string // for CONTINUE action
	ProjectRoot string
	FeatureSlug string
}

// executeClaim contains the core logic for the claim command, testable
func executeClaim() (*ClaimResult, error) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil, ErrProjectNotFound()
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return nil, ErrFeatureNotSet()
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return nil, NewAIError(ErrNotFound, "Failed to load task index", err.Error(), "Check index.json exists and is valid", "cat "+indexPath)
	}

	// Check for existing task state
	statePath := feature.GetTaskStatePath(projectRoot, featureSlug)
	continueTask, hasIssues, issues := checkExistingTaskState(projectRoot, index, statePath)

	if hasIssues {
		return nil, ErrDataIntegrity(issues)
	}

	if continueTask {
		state, _ := task.LoadState(statePath)
		t, _ := index.ByID(state.Key)
		return &ClaimResult{
			Action:      "CONTINUE",
			Key:         state.Key,
			Task:        &t,
			State:       state,
			StartedTime: state.StartedTime,
			ProjectRoot: projectRoot,
			FeatureSlug: featureSlug,
		}, nil
	}

	// Claim new task
	key, t, err := claimNextTask(index)
	if err != nil {
		return nil, err
	}

	// Update index
	if err := task.SaveIndex(indexPath, index); err != nil {
		return nil, err
	}

	// Bootstrap .forge/state.json as workspace marker so subagents in subdirectories
	// (e.g., backend/) can find the project root via FindProjectRoot().
	// Placed after all validation to avoid creating artifacts in error paths.
	if err := feature.EnsureForgeState(projectRoot, featureSlug); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: failed to write .forge/state.json: %v\n", err)
	}

	// Save state
	state := &task.TaskState{
		TaskID:        t.ID,
		Key:           key,
		Title:         t.Title,
		Priority:      t.Priority,
		EstimatedTime: t.EstimatedTime,
		Dependencies:  t.Dependencies,
		File:          t.File,
		Record:        t.Record,
		StartedTime:   time.Now().Format("2006-01-02 15:04"),
		Breaking:      t.Breaking,
		Scope:         t.Scope,
		MainSession:   t.MainSession,
	}
	if err := task.SaveState(statePath, state); err != nil {
		return nil, err
	}

	return &ClaimResult{
		Action:      "CLAIMED",
		Key:         key,
		Task:        t,
		State:       state,
		ProjectRoot: projectRoot,
		FeatureSlug: featureSlug,
	}, nil
}

func checkExistingTaskState(projectRoot string, index *task.TaskIndex, statePath string) (bool, bool, []string) {
	state, err := task.LoadState(statePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load task state: %v\n", err)
		return false, false, nil
	}
	if state == nil {
		return false, false, nil
	}

	t, exists := index.ByID(state.Key)
	if !exists {
		return false, true, []string{fmt.Sprintf("Task key '%s' not found in index.json", state.Key)}
	}

	switch t.Status {
	case "in_progress":
		return true, false, nil
	case "completed":
		fmt.Printf("Previous task '%s' is completed. Claiming new task...\n", t.Title)
		task.DeleteState(statePath)
		return false, false, nil
	default:
		return false, true, []string{fmt.Sprintf("Task '%s' has unexpected status: %s", t.Title, t.Status)}
	}
}

func claimNextTask(index *task.TaskIndex) (string, *task.Task, error) {
	type taskWithKey struct {
		key string
		t   task.Task
	}
	var eligibleTasks []taskWithKey

	hasPending := false
	for _, t := range index.TasksMap() {
		if t.Status == "pending" {
			hasPending = true
			break
		}
	}
	if !hasPending {
		return "", nil, ErrNoPendingTasks()
	}

	for key, t := range index.TasksMap() {
		if t.Status == "pending" {
			if met, _ := checkDependenciesMet(index, t.ID, t); met {
				eligibleTasks = append(eligibleTasks, taskWithKey{key: key, t: t})
			}
		}
	}

	if len(eligibleTasks) == 0 {
		return "", nil, fmt.Errorf("no task available with met dependencies")
	}

	priorityOrder := map[string]int{"P0": 0, "P1": 1, "P2": 2}
	sort.Slice(eligibleTasks, func(i, j int) bool {
		pi, pj := priorityOrder[eligibleTasks[i].t.Priority], priorityOrder[eligibleTasks[j].t.Priority]
		if pi != pj {
			return pi < pj
		}
		return compareVersionIDs(eligibleTasks[i].t.ID, eligibleTasks[j].t.ID)
	})

	twk := eligibleTasks[0]
	t, _ := index.ByID(twk.key)
	t.Status = "in_progress"
	index.SetTask(twk.key, t)
	return twk.key, &t, nil
}

func getTaskPhase(id string) int {
	parts := strings.Split(id, ".")
	if len(parts) > 0 {
		phase, err := strconv.Atoi(parts[0])
		if err == nil {
			return phase
		}
	}
	return -1
}

func checkDependenciesMet(index *task.TaskIndex, selfID string, t task.Task) (bool, []string) {
	var unmet []string
	for _, dep := range t.Dependencies {
		if strings.HasSuffix(dep, ".x") {
			prefix := strings.TrimSuffix(dep, ".x")
			found := false
			prefixWithDot := prefix + "."
			for _, other := range index.TasksMap() {
				if other.ID == selfID {
					continue
				}
				if strings.HasPrefix(other.ID, prefixWithDot) && isBusinessTask(other.ID) && other.Status != "completed" && other.Status != "skipped" {
					unmet = append(unmet, other.ID)
					found = true
				}
			}
			if !found {
				continue
			}
		} else {
			for _, other := range index.TasksMap() {
				if other.ID == dep {
					if other.Status != "completed" && other.Status != "skipped" {
						unmet = append(unmet, dep)
					}
					break
				}
			}
		}
	}
	return len(unmet) == 0, unmet
}

func compareVersionIDs(a, b string) bool {
	partsA := strings.Split(a, ".")
	partsB := strings.Split(b, ".")
	maxLen := len(partsA)
	if len(partsB) > maxLen {
		maxLen = len(partsB)
	}
	for i := 0; i < maxLen; i++ {
		na, aIsNum := parseSegment(partsA, i)
		nb, bIsNum := parseSegment(partsB, i)
		if aIsNum != bIsNum {
			return aIsNum // numeric sorts before alphabetic
		}
		if na != nb {
			return na < nb
		}
	}
	return false
}

// parseSegment returns the numeric value of a segment and whether it's numeric.
// Numeric segments (e.g., "1", "12") return their int value with true.
// Alphabetic segments (e.g., "summary", "gate") return a lexicographic rank with false.
func parseSegment(parts []string, i int) (int, bool) {
	if i >= len(parts) {
		return -1, true // missing segments sort before everything
	}
	if n, err := strconv.Atoi(parts[i]); err == nil {
		return n, true
	}
	// Alphabetic segments: sort after all numeric, with deterministic order
	switch parts[i] {
	case "gate":
		return 1, false
	case "summary":
		return 2, false
	default:
		return 0, false
	}
}

func printTaskDetails(key string, t *task.Task, projectRoot, featureSlug string) {
	PrintField("KEY", key)
	PrintField("TASK_ID", t.ID)
	PrintField("TITLE", t.Title)
	PrintField("PRIORITY", t.Priority)
	PrintField("STATUS", t.Status)
	PrintFieldIfNotEmpty("ESTIMATED_TIME", t.EstimatedTime)
	PrintFieldIfNotEmptySlice("DEPENDENCIES", t.Dependencies)
	PrintField("BREAKING", strconv.FormatBool(t.Breaking))
	PrintField("MAIN_SESSION", strconv.FormatBool(t.MainSession))
	PrintFieldIfNotEmpty("SCOPE", t.Scope)
	PrintFieldIfNotEmpty("FEATURE", featureSlug)
	PrintField("FILE", filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, t.File)))
	PrintField("RECORD", filepath.Join(projectRoot, feature.GetTaskFile(featureSlug, t.Record)))
}

func printContinueTask(state *task.TaskState, t *task.Task, projectRoot, featureSlug string) {
	PrintBlockStart()
	PrintField("ACTION", "CONTINUE")
	printTaskDetails(state.Key, t, projectRoot, featureSlug)
	PrintField("STARTED_AT", state.StartedTime)
	PrintBlockEnd()
}

func printNewTask(key string, t *task.Task, projectRoot, featureSlug string) {
	PrintBlockStart()
	PrintField("ACTION", "CLAIMED")
	printTaskDetails(key, t, projectRoot, featureSlug)
	PrintBlockEnd()
}
