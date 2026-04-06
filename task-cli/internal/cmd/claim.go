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
1. Minimum phase with pending tasks
2. All dependencies must be met
3. Priority (P0 > P1 > P2)
4. Task ID (semantic version ordering)`,
	Run: runClaim,
}

func runClaim(cmd *cobra.Command, args []string) {
	result, err := executeClaim()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print result
	if result.Action == "CONTINUE" {
		printContinueTask(result.State, result.Task)
	} else {
		printNewTask(result.Key, result.Task)
	}
}

// ClaimResult represents the result of a claim operation
type ClaimResult struct {
	Action      string // "CLAIMED" or "CONTINUE"
	Key         string
	Task        *task.Task
	State       *task.TaskState
	StartedTime string // for CONTINUE action
}

// executeClaim contains the core logic for the claim command, testable
func executeClaim() (*ClaimResult, error) {
	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return nil, err
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return nil, err
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return nil, err
	}

	// Check for existing task state
	statePath := feature.GetTaskStatePath(projectRoot, featureSlug)
	continueTask, hasIssues, issues := checkExistingTaskState(projectRoot, index, statePath)

	if hasIssues {
		return nil, fmt.Errorf("task data integrity issues: %v", issues)
	}

	if continueTask {
		state, _ := task.LoadState(statePath)
		t := index.Tasks[state.Key]
		return &ClaimResult{
			Action:      "CONTINUE",
			Key:         state.Key,
			Task:        &t,
			State:       state,
			StartedTime: state.StartedTime,
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
	}
	if err := task.SaveState(statePath, state); err != nil {
		return nil, err
	}

	return &ClaimResult{
		Action: "CLAIMED",
		Key:    key,
		Task:   t,
		State:  state,
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

	t, exists := index.Tasks[state.Key]
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

	minPhase := getMinPendingPhase(index)
	if minPhase == -1 {
		return "", nil, fmt.Errorf("no pending tasks available")
	}

	for key, t := range index.Tasks {
		if t.Status == "pending" {
			if met, _ := checkDependenciesMet(index, t); met {
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
	t := index.Tasks[twk.key]
	t.Status = "in_progress"
	index.Tasks[twk.key] = t
	return twk.key, &t, nil
}

func getMinPendingPhase(index *task.TaskIndex) int {
	minPhase := -1
	for _, t := range index.Tasks {
		if t.Status == "pending" {
			phase := getTaskPhase(t.ID)
			if phase != -1 && (minPhase == -1 || phase < minPhase) {
				minPhase = phase
			}
		}
	}
	return minPhase
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

func checkDependenciesMet(index *task.TaskIndex, t task.Task) (bool, []string) {
	var unmet []string
	for _, dep := range t.Dependencies {
		if strings.HasSuffix(dep, ".x") || strings.HasSuffix(dep, "x") {
			prefix := strings.TrimSuffix(strings.TrimSuffix(dep, "x"), ".")
			found := false
			prefixWithDot := prefix + "."
			for _, other := range index.Tasks {
				if strings.HasPrefix(other.ID, prefixWithDot) && other.Status != "completed" {
					unmet = append(unmet, other.ID)
					found = true
				}
			}
			if !found {
				continue
			}
		} else {
			for _, other := range index.Tasks {
				if other.ID == dep {
					if other.Status != "completed" {
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
		var numA, numB int
		if i < len(partsA) {
			numA, _ = strconv.Atoi(partsA[i])
		}
		if i < len(partsB) {
			numB, _ = strconv.Atoi(partsB[i])
		}
		if numA != numB {
			return numA < numB
		}
	}
	return false
}

func printTaskDetails(key string, t *task.Task) {
	PrintField("KEY", key)
	PrintField("ID", t.ID)
	PrintField("TITLE", t.Title)
	PrintField("PRIORITY", t.Priority)
	PrintField("STATUS", t.Status)
	PrintFieldIfNotEmpty("ESTIMATED_TIME", t.EstimatedTime)
	PrintFieldIfNotEmptySlice("DEPENDENCIES", t.Dependencies)
	PrintField("FILE", t.File)
	PrintField("RECORD", t.Record)
}

func printContinueTask(state *task.TaskState, t *task.Task) {
	PrintBlockStart()
	PrintField("ACTION", "CONTINUE")
	printTaskDetails(state.Key, t)
	PrintField("STARTED_AT", state.StartedTime)
	PrintBlockEnd()
}

func printNewTask(key string, t *task.Task) {
	PrintBlockStart()
	PrintField("ACTION", "CLAIMED")
	printTaskDetails(key, t)
	PrintBlockEnd()
}
