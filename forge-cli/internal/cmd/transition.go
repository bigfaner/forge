package cmd

import (
	"errors"
	"path/filepath"

	"forge-cli/pkg/feature"
	indexPkg "forge-cli/pkg/index"
	"forge-cli/pkg/project"
	"forge-cli/pkg/task"

	"github.com/spf13/cobra"
)

var transitionReason string

var transitionCmd = &cobra.Command{
	Use:   "transition <task-id> <status>",
	Short: "Manually transition a task to a new status",
	Long: `Manually transition a task to a new status with a required reason.

This is an operator override for situations not covered by the standard
workflow (submit, claim, reopen):

  - Unblock a blocked task:       forge task transition 1.2 pending --reason "deps resolved manually"
  - Skip an uncompletable task:   forge task transition 1.2 skipped --reason "superseded by 2.1"
  - Reject a task:                forge task transition 1.2 rejected --reason "out of scope"
  - Block a task manually:        forge task transition 1.2 blocked --reason "waiting on external API"

Terminal state protection: completed, rejected, and skipped tasks can NEVER be transitioned.
Use "forge task reopen" for rejected/skipped -> pending.`,
	Args: cobra.ExactArgs(2),
	RunE: runTransition,
}

func init() {
	transitionCmd.Flags().StringVar(&transitionReason, "reason", "", "Required: reason for the transition")
	_ = transitionCmd.MarkFlagRequired("reason")
}

func runTransition(_ *cobra.Command, args []string) error {
	taskIDArg := args[0]
	targetStatus := args[1]

	projectRoot, err := project.FindProjectRoot()
	if err != nil {
		return ErrProjectNotFound()
	}

	featureSlug, err := feature.RequireFeature(projectRoot)
	if err != nil {
		return ErrFeatureNotSet()
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))

	if lockErr := indexPkg.WithLock(indexPath, func() error {
		return doTransition(indexPath, taskIDArg, targetStatus)
	}); lockErr != nil {
		if errors.Is(lockErr, indexPkg.ErrLockConflict) {
			return NewAIError(ErrConflict, "Concurrent write conflict", "Retry the command", "", "")
		}
		if aiErr, ok := lockErr.(*AIError); ok {
			return aiErr
		}
		return NewAIError(ErrConflict, "Failed to acquire lock", lockErr.Error(), "", "")
	}
	return nil
}

func doTransition(indexPath, taskIDArg, targetStatus string) error {
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		return ErrFileNotFound(indexPath)
	}

	key, t, err := task.FindTask(index, taskIDArg)
	if err != nil {
		return ErrTaskNotFound(taskIDArg)
	}

	// Validate via state machine (RoleManual enables blocked -> pending)
	if transitionErr := task.ValidateTransition(t.Status, targetStatus, task.RoleManual); transitionErr != nil {
		te := transitionErr.(*task.TransitionError)
		return NewErrInvalidTransition(t.Status, targetStatus, te.Msg)
	}

	// Validate target status against index enum
	valid := false
	for _, s := range index.StatusEnum {
		if s == targetStatus {
			valid = true
			break
		}
	}
	if !valid {
		return ErrInvalidStatus(targetStatus, index.StatusEnum)
	}

	t.Status = targetStatus
	if targetStatus == "blocked" {
		t.BlockedReason = transitionReason
		t.ManualBlock = true
	} else {
		t.ManualBlock = false
	}

	index.SetTask(key, *t)

	if err := indexPkg.SaveIndexAtomic(indexPath, index); err != nil {
		return NewAIError(ErrConflict, "Failed to save index", err.Error(), "Check index.json is writable", "cat "+indexPath)
	}

	PrintBlockStart()
	PrintField("TASK_ID", t.ID)
	PrintField("STATUS", t.Status)
	PrintField("REASON", transitionReason)
	PrintBlockEnd()
	return nil
}
