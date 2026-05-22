package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"forge-cli/pkg/feature"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var reopenCmd = &cobra.Command{
	Use:   "reopen <task-id>",
	Short: "Re-activate a rejected or skipped task to pending",
	Long: `Re-activate a rejected or skipped task back to pending status.

Terminal state protection:
  - completed tasks are NEVER re-openable
  - only rejected and skipped tasks can be reopened
  - reopen target is always pending (cannot specify other states)`,
	Args: cobra.ExactArgs(1),
	Run:  runReopen,
}

func runReopen(_ *cobra.Command, args []string) {
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

	if lockErr := indexPkg.WithLock(indexPath, func() error {
		return doReopen(indexPath, taskIDArg)
	}); lockErr != nil {
		if errors.Is(lockErr, indexPkg.ErrLockConflict) {
			fmt.Fprintln(os.Stderr, "concurrent write conflict, retry")
			os.Exit(1)
		}
		if aiErr, ok := lockErr.(*AIError); ok {
			Exit(aiErr)
		}
		fmt.Fprintf(os.Stderr, "failed to acquire lock: %v\n", lockErr)
		os.Exit(1)
	}
}

func doReopen(indexPath, taskIDArg string) error {
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return ErrFileNotFound(indexPath)
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		return ErrTaskNotFound(taskIDArg)
	}

	// Validate transition using state machine (RoleReopen, target always pending)
	if transitionErr := task.ValidateTransition(t.Status, "pending", task.RoleReopen); transitionErr != nil {
		te := transitionErr.(*task.TransitionError)
		return NewErrInvalidTransition(t.Status, "pending", te.Msg)
	}

	t.Status = "pending"
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
