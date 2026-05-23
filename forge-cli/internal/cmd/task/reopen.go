package task

import (
	"errors"
	"forge-cli/internal/cmd/base"
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
	RunE: runReopen,
}

func runReopen(_ *cobra.Command, args []string) error {
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

	if lockErr := indexPkg.WithLock(indexPath, func() error {
		return doReopen(indexPath, taskIDArg)
	}); lockErr != nil {
		if errors.Is(lockErr, indexPkg.ErrLockConflict) {
			return base.NewAIError(base.ErrConflict, "Concurrent write conflict", "Retry the command", "Wait a moment and try again", "forge reopen "+taskIDArg)
		}
		if aiErr, ok := lockErr.(*base.AIError); ok {
			return aiErr
		}
		return base.NewAIError(base.ErrConflict, "Failed to acquire lock", lockErr.Error(), "Check index.json is not locked by another process", "cat "+indexPath)
	}
	return nil
}

func doReopen(indexPath, taskIDArg string) error {
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return base.ErrFileNotFound(indexPath)
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		return base.ErrTaskNotFound(taskIDArg)
	}

	// Validate transition using state machine (RoleReopen, target always pending)
	if transitionErr := task.ValidateTransition(t.Status, "pending", task.RoleReopen); transitionErr != nil {
		te := transitionErr.(*task.TransitionError)
		return base.NewErrInvalidTransition(t.Status, "pending", te.Msg)
	}

	t.Status = "pending"
	index.SetTask(key, *t)

	if err := indexPkg.SaveIndexAtomic(indexPath, index); err != nil {
		return base.NewAIError(base.ErrConflict, "Failed to save index", err.Error(), "Check index.json is writable", "cat "+indexPath)
	}

	base.PrintBlockStart()
	base.PrintField("TASK_ID", t.ID)
	base.PrintField("STATUS", t.Status)
	base.PrintBlockEnd()
	return nil
}
