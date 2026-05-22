package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"forge-cli/pkg/feature"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <task-id> [<status>]",
	Short: "Query or set task status",
	Long: `Query the status of a task, or set it to a new status.

With one argument, displays the current status of the task.
With two arguments, sets the task to the given status (e.g., "blocked").

Use "forge task submit" to complete a task or "forge task reopen" to re-activate rejected/skipped tasks.`,
	Args: cobra.RangeArgs(1, 2),
	Run:  runStatus,
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

	if len(args) > 1 {
		// Mutation mode: set task status
		statusArg := args[1]
		if err := indexPkg.WithLock(indexPath, func() error {
			return doSetStatus(indexPath, taskIDArg, statusArg)
		}); err != nil {
			if errors.Is(err, indexPkg.ErrLockConflict) {
				fmt.Fprintln(os.Stderr, "concurrent write conflict, retry")
				os.Exit(1)
			}
			if aiErr, ok := err.(*AIError); ok {
				Exit(aiErr)
			}
			fmt.Fprintf(os.Stderr, "failed to update status: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Query mode: display task status
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		Exit(ErrFileNotFound(indexPath))
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		Exit(ErrTaskNotFound(taskIDArg))
	}

	_ = key

	PrintBlockStart()
	PrintField("TASK_ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintBlockEnd()
}

// doSetStatus mutates a task's status after validating the transition.
func doSetStatus(indexPath, taskIDArg, statusArg string) error {
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return ErrFileNotFound(indexPath)
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		return ErrTaskNotFound(taskIDArg)
	}

	// Validate transition using state machine
	if transitionErr := task.ValidateTransition(t.Status, statusArg, ""); transitionErr != nil {
		te := transitionErr.(*task.TransitionError)
		return NewErrInvalidTransition(t.Status, statusArg, te.Msg)
	}

	t.Status = statusArg
	index.SetTask(key, *t)

	if err := task.SaveIndex(indexPath, index); err != nil {
		return NewAIError(ErrConflict, "Failed to save index", err.Error(), "Check index.json is writable", "cat "+indexPath)
	}

	PrintBlockStart()
	PrintField("TASK_ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintBlockEnd()
	return nil
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
		if strings.HasSuffix(dep, task.IDSuffixWildcard) {
			prefix := strings.TrimSuffix(dep, task.IDSuffixWildcard)
			prefixWithDot := prefix + "."
			found := false
			for _, other := range index.TasksMap() {
				if other.ID == t.ID {
					continue
				}
				if strings.HasPrefix(other.ID, prefixWithDot) && task.IsBusinessTask(other.ID) && other.Status != "completed" && other.Status != "skipped" {
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
