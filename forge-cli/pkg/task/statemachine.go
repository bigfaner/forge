package task

import (
	"fmt"
	"strings"
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
	From string
	To   string
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
	From     string         // current status; "*" matches any
	To       string         // target status; "*" matches any
	Role     TransitionRole // required role; "" matches any
	Allowed  bool           // whether the transition is permitted
	GuardMsg string         // human-readable reason when blocked
}

// transitionTable is the single authority for transition validation.
// Rules are evaluated in order; first match wins.
var transitionTable = []TransitionRule{
	// Terminal state: completed is irreversible
	{From: "completed", To: "*", Role: "", Allowed: false, GuardMsg: "task already completed, create a subtask if re-work needed"},

	// Terminal state: rejected can only go to pending via reopen
	{From: "rejected", To: "pending", Role: RoleReopen, Allowed: true, GuardMsg: ""},
	{From: "rejected", To: "*", Role: "", Allowed: false, GuardMsg: "task rejected, use forge task reopen"},

	// Terminal state: skipped can only go to pending via reopen
	{From: "skipped", To: "pending", Role: RoleReopen, Allowed: true, GuardMsg: ""},
	{From: "skipped", To: "*", Role: "", Allowed: false, GuardMsg: "task skipped, use forge task reopen"},

	// Only submit can reach completed
	{From: "*", To: "completed", Role: RoleSubmit, Allowed: true, GuardMsg: ""},
	{From: "*", To: "completed", Role: "", Allowed: false, GuardMsg: "use forge task submit"},

	// Submit can auto-downgrade in_progress to blocked
	{From: "in_progress", To: "blocked", Role: RoleSubmit, Allowed: true, GuardMsg: ""},

	// Manual override: operator can unblock or resolve any non-completed task
	{From: "blocked", To: "pending", Role: RoleManual, Allowed: true, GuardMsg: ""},
	{From: "blocked", To: "in_progress", Role: RoleManual, Allowed: true, GuardMsg: ""},

	// blocked -> pending/in_progress requires dependency check (phase 2)
	{From: "blocked", To: "pending", Role: "", Allowed: false, GuardMsg: "dependencies must be checked first"},
	{From: "blocked", To: "in_progress", Role: "", Allowed: false, GuardMsg: "dependencies must be checked first"},

	// pending -> blocked is always allowed (block-source, dependency wait)
	{From: "pending", To: "blocked", Role: "", Allowed: true, GuardMsg: ""},

	// RoleReopen is only valid for rejected/skipped -> pending (handled above).
	// Using reopen on non-terminal states is invalid.
	{From: "*", To: "*", Role: RoleReopen, Allowed: false, GuardMsg: "reopen is only for rejected or skipped tasks"},

	// Same-state transition: no-op, always allowed (except terminal states handled above)
	{From: "*", To: "*", Role: "", Allowed: true, GuardMsg: ""},
}

// ValidateTransition validates a state transition (pure, no data lookup).
// Phase 1 of validation: checks terminal state protection and role-based rules.
// Returns nil if the transition is allowed, or a *TransitionError if not.
// For blocked -> pending/in_progress transitions, returns an error indicating
// that dependency checking (phase 2) is required.
func ValidateTransition(current, target string, role TransitionRole) error {
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
func matchRule(rule TransitionRule, from, to string, role TransitionRole) bool {
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

// depSatisfiedStatuses are the statuses that satisfy dependency checks.
// "completed" and "skipped" satisfy dependencies.
// "rejected", "blocked", "pending", "in_progress" do NOT satisfy.
var depSatisfiedStatuses = map[string]bool{
	"completed": true,
	"skipped":   true,
}

// CheckTransitionDeps validates dependency satisfaction for blocked -> pending/in_progress.
// Phase 2 of validation: call after ValidateTransition indicates dependency check needed.
// Returns unmet dependency IDs, or nil if all deps are met (including canAutoUnblock check).
func CheckTransitionDeps(index *TaskIndex, taskID string) ([]string, error) {
	task, found := index.ByID(taskID)
	if !found {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	// Check basic dependency satisfaction
	var unmet []string
	for _, depID := range task.Dependencies {
		dep, depFound := index.ByID(depID)
		if !depFound {
			unmet = append(unmet, depID)
			continue
		}
		if !depSatisfiedStatuses[dep.Status] {
			unmet = append(unmet, depID)
		}
	}

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

// isTerminalStatus checks if a status is terminal (completed, rejected).
func isTerminalStatus(status string) bool {
	return status == "completed" || status == "rejected"
}
