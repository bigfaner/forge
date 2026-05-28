package task

import (
	"fmt"
	"strings"

	"forge-cli/pkg/types"
)

// TransitionRole represents the role performing a state transition.
type TransitionRole string

const (
	// RoleSubmit represents forge task submit.
	RoleSubmit TransitionRole = "submit"
	// RoleClaim represents forge task claim.
	RoleClaim TransitionRole = "claim"
	// RoleReopen represents forge task reopen.
	RoleReopen TransitionRole = "reopen"
	// RoleAuto represents auto-downgrade or auto-unblock.
	RoleAuto TransitionRole = "auto"
	// RoleManual represents manual operator override (forge task transition).
	RoleManual TransitionRole = "manual"
)

// TransitionError is returned when a state transition is not allowed.
type TransitionError struct {
	From types.Status
	To   types.Status
	Role TransitionRole
	Msg  string
}

// Error implements the error interface.
func (e *TransitionError) Error() string {
	return fmt.Sprintf("invalid transition %s -> %s (role=%s): %s", e.From, e.To, e.Role, e.Msg)
}

// TransitionRule defines a single entry in the transition rule table.
// The table is the single authority for state validation.
type TransitionRule struct {
	From     types.Status   // current status; "*" matches any
	To       types.Status   // target status; "*" matches any
	Role     TransitionRole // required role; "" matches any
	Allowed  bool           // whether the transition is permitted
	GuardMsg string         // human-readable reason when blocked
}

// transitionTable is the single authority for transition validation.
// Rules are evaluated in order; first match wins.
var transitionTable = []TransitionRule{
	// Terminal state: completed is irreversible
	{From: types.StatusCompleted, To: types.Status("*"), Role: "", Allowed: false, GuardMsg: "task already completed, create a subtask if re-work needed"},

	// Terminal state: rejected can only go to pending via reopen
	{From: types.StatusRejected, To: types.StatusPending, Role: RoleReopen, Allowed: true, GuardMsg: ""},
	{From: types.StatusRejected, To: types.Status("*"), Role: "", Allowed: false, GuardMsg: "task rejected, use forge task reopen"},

	// Terminal state: skipped can only go to pending via reopen
	{From: types.StatusSkipped, To: types.StatusPending, Role: RoleReopen, Allowed: true, GuardMsg: ""},
	{From: types.StatusSkipped, To: types.Status("*"), Role: "", Allowed: false, GuardMsg: "task skipped, use forge task reopen"},

	// Suspended cannot directly reach completed (must resume first).
	// Placed before the general "submit -> completed" rule so it matches first.
	{From: types.StatusSuspended, To: types.StatusCompleted, Role: "", Allowed: false, GuardMsg: "use forge task transition to resume task first"},

	// Only submit can reach completed
	{From: types.Status("*"), To: types.StatusCompleted, Role: RoleSubmit, Allowed: true, GuardMsg: ""},
	{From: types.Status("*"), To: types.StatusCompleted, Role: "", Allowed: false, GuardMsg: "use forge task submit"},

	// Submit can auto-downgrade in_progress to blocked
	{From: types.StatusInProgress, To: types.StatusBlocked, Role: RoleSubmit, Allowed: true, GuardMsg: ""},

	// Manual override: operator can unblock or resolve any non-completed task
	{From: types.StatusBlocked, To: types.StatusPending, Role: RoleManual, Allowed: true, GuardMsg: ""},
	{From: types.StatusBlocked, To: types.StatusInProgress, Role: RoleManual, Allowed: true, GuardMsg: ""},

	// blocked -> pending/in_progress requires dependency check (phase 2)
	{From: types.StatusBlocked, To: types.StatusPending, Role: "", Allowed: false, GuardMsg: "dependencies must be checked first"},
	{From: types.StatusBlocked, To: types.StatusInProgress, Role: "", Allowed: false, GuardMsg: "dependencies must be checked first"},

	// pending -> blocked is always allowed (block-source, dependency wait)
	{From: types.StatusPending, To: types.StatusBlocked, Role: "", Allowed: true, GuardMsg: ""},

	// --- suspended: operator manual hold ---
	// Only RoleManual can enter suspended (from any non-terminal state).
	{From: types.Status("*"), To: types.StatusSuspended, Role: RoleManual, Allowed: true, GuardMsg: ""},
	{From: types.Status("*"), To: types.StatusSuspended, Role: "", Allowed: false, GuardMsg: "use forge task transition to suspend tasks"},
	// Manual resume from suspended.
	{From: types.StatusSuspended, To: types.StatusPending, Role: RoleManual, Allowed: true, GuardMsg: ""},
	{From: types.StatusSuspended, To: types.StatusInProgress, Role: RoleManual, Allowed: true, GuardMsg: ""},
	// Manual terminal decisions from suspended.
	{From: types.StatusSuspended, To: types.StatusSkipped, Role: RoleManual, Allowed: true, GuardMsg: ""},
	{From: types.StatusSuspended, To: types.StatusRejected, Role: RoleManual, Allowed: true, GuardMsg: ""},
	// Block system transitions from suspended to blocked (must resume first).
	{From: types.StatusSuspended, To: types.StatusBlocked, Role: "", Allowed: false, GuardMsg: "use forge task transition to resume suspended task"},

	// RoleReopen is only valid for rejected/skipped -> pending (handled above).
	// Using reopen on non-terminal states is invalid.
	{From: types.Status("*"), To: types.Status("*"), Role: RoleReopen, Allowed: false, GuardMsg: "reopen is only for rejected or skipped tasks"},

	// Same-state transition: no-op, always allowed (except terminal states handled above)
	{From: types.Status("*"), To: types.Status("*"), Role: "", Allowed: true, GuardMsg: ""},
}

// ValidateTransition validates a state transition (pure, no data lookup).
// Phase 1 of validation: checks terminal state protection and role-based rules.
// Returns nil if the transition is allowed, or a *TransitionError if not.
// For blocked -> pending/in_progress transitions, returns an error indicating
// that dependency checking (phase 2) is required.
func ValidateTransition(current, target types.Status, role TransitionRole) error {
	for _, rule := range transitionTable {
		if matchRule(rule, current, target, role) {
			if rule.Allowed {
				return nil
			}
			return &TransitionError{
				From: current,
				To:   target,
				Role: role,
				Msg:  rule.GuardMsg,
			}
		}
	}
	// No rule matched: allow by default (open transition)
	return nil
}

// matchRule checks if a transition rule matches the given parameters.
func matchRule(rule TransitionRule, from, to types.Status, role TransitionRole) bool {
	if rule.From != "*" && rule.From != from {
		return false
	}
	if rule.To != "*" && rule.To != to {
		return false
	}
	if rule.Role != "" && rule.Role != role {
		return false
	}
	return true
}

// CheckTransitionDeps validates dependency satisfaction for blocked -> pending/in_progress.
// Phase 2 of validation: call after ValidateTransition indicates dependency check needed.
// Returns unmet dependency IDs, or nil if all deps are met (including canAutoUnblock check).
// Delegates to GetUnmetDeps for wildcard-aware dependency resolution.
func CheckTransitionDeps(index *TaskIndex, taskID string) ([]string, error) {
	t, found := index.ByID(taskID)
	if !found {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	unmet := GetUnmetDeps(index, taskID, t.Dependencies)

	if len(unmet) > 0 {
		return unmet, nil
	}

	// All deps met — check canAutoUnblock (active fix-task check)
	if !canAutoUnblock(index, taskID) {
		return []string{"active fix-task exists"}, nil
	}

	return nil, nil
}

// canAutoUnblock checks whether a blocked task can be unblocked by verifying
// no active fix-tasks are targeting it. A fix-task is considered active if it
// has a SourceTaskID matching the task and is in a non-terminal state.
// This is an unexported helper — not a standalone interface.
func canAutoUnblock(index *TaskIndex, taskID string) bool {
	for _, t := range index.TasksMap() {
		if t.SourceTaskID == taskID && isActiveFixTask(t) {
			return false
		}
	}
	return true
}

// isActiveFixTask returns true if the task is a fix-type task in a non-terminal state.
func isActiveFixTask(t Task) bool {
	if !isFixType(t.Type) {
		return false
	}
	return !isTerminalStatus(t.Status)
}

// isFixType checks if a task type indicates a fix task.
func isFixType(typ string) bool {
	return typ == TypeCodingFix || strings.HasPrefix(typ, "coding.fix")
}

// isTerminalStatus checks if a status is terminal (completed, skipped, rejected).
func isTerminalStatus(status types.Status) bool {
	return types.IsTerminalStatus(status)
}
