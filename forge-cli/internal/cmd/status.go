package cmd

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var statusForce bool

var statusCmd = &cobra.Command{
	Use:   "status <task-id> [status]",
	Short: "Query or update task status",
	Long: `Query or update the status of a task.

Without status argument: query current status.
With status argument: update to new status.

State machine guards:
  - completed and rejected are terminal (cannot leave without --force)
  - in_progress -> completed is blocked (use "forge task submit" instead)
  - pending/in_progress transitions require all dependencies to be completed or skipped`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runStatus,
}

func init() {
	statusCmd.Flags().BoolVar(&statusForce, "force", false, "Override state machine guards (use with caution)")
}

func runStatus(_ *cobra.Command, args []string) {
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

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		Exit(ErrTaskNotFound(taskIDArg))
	}

	// Query mode: only one argument
	if len(args) == 1 {
		PrintBlockStart()
		PrintField("TASK_ID", t.ID)
		PrintField("STATUS", t.Status)
		PrintBlockEnd()
		return
	}

	// Update mode: two arguments
	newStatus := args[1]

	// Validate status is in enum
	if !slices.Contains(index.StatusEnum, newStatus) {
		Exit(ErrInvalidStatus(newStatus, index.StatusEnum))
	}

	// State machine validation (unless --force)
	if !statusForce && !isTransitionAllowed(t.Status, newStatus) {
		Exit(NewAIError(ErrValidation,
			fmt.Sprintf("Invalid transition: %s -> %s", t.Status, newStatus),
			fmt.Sprintf("Current status is %s", t.Status),
			getTransitionHint(t.Status, newStatus),
			getTransitionAction(t.Status, newStatus),
		))
	}

	// Dependency check for pending and in_progress transitions
	if newStatus == "pending" || newStatus == "in_progress" {
		unmet := checkUnmetDeps(index, t)
		if len(unmet) > 0 {
			PrintBlockStart()
			PrintField("TASK_ID", t.ID)
			PrintField("STATUS", t.Status)
			fmt.Printf("WARNING: %s has unmet dependencies: %s. Status not changed.\n", t.ID, strings.Join(unmet, ", "))
			PrintBlockEnd()
			return
		}
	}

	t.Status = newStatus
	index.SetTask(key, *t)

	if err := task.SaveIndex(indexPath, index); err != nil {
		Exit(NewAIError(ErrConflict, "Failed to save index", err.Error(), "Check index.json is writable", "cat "+indexPath))
	}

	PrintBlockStart()
	PrintField("TASK_ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintBlockEnd()
}

// isTransitionAllowed returns true if the state transition is valid.
// completed and rejected are terminal. in_progress -> completed must go through task record.
// rejected does not satisfy dependency checks — downstream tasks cannot proceed.
func isTransitionAllowed(from, to string) bool {
	if from == to {
		return true
	}
	if from == "completed" || from == "rejected" {
		return false
	}
	if to == "completed" {
		return false
	}
	return true
}

func getTransitionHint(from, to string) string {
	if from == "completed" {
		return "completed is a terminal state"
	}
	if from == "rejected" {
		return "rejected is a terminal state (task ran but acceptance criteria not met)"
	}
	if to == "completed" {
		return "use 'forge task submit' to complete a task with quality gate"
	}
	return fmt.Sprintf("transition %s -> %s is not allowed", from, to)
}

func getTransitionAction(from, to string) string {
	if from == "completed" || from == "rejected" {
		return "use --force to override (may break lifecycle tracking)"
	}
	if to == "completed" {
		return "run 'forge task submit <task-id> --data record.json' or use --force"
	}
	return "use --force to override"
}

// checkUnmetDeps returns dependency IDs that are not "completed" or "skipped".
// Supports both exact IDs and wildcard patterns (e.g. "1.x").
// selfID is excluded from wildcard matches to avoid self-matching.
func checkUnmetDeps(index *task.TaskIndex, t *task.Task) []string {
	var unmet []string
	for _, dep := range t.Dependencies {
		if strings.HasSuffix(dep, ".x") {
			prefix := strings.TrimSuffix(dep, ".x")
			prefixWithDot := prefix + "."
			found := false
			for _, other := range index.TasksMap() {
				if other.ID == t.ID {
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
			_, depTask, err := task.FindTask(index, dep)
			if err != nil || (depTask.Status != "completed" && depTask.Status != "skipped") {
				unmet = append(unmet, dep)
			}
		}
	}
	return unmet
}
