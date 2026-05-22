package task

import (
	"errors"
	"testing"
)

// --- TransitionRole constants ---

func TestTransitionRoleConstants(t *testing.T) {
	tests := []struct {
		role     TransitionRole
		expected string
	}{
		{RoleSubmit, "submit"},
		{RoleClaim, "claim"},
		{RoleReopen, "reopen"},
		{RoleAuto, "auto"},
		{RoleManual, "manual"},
	}
	for _, tt := range tests {
		if string(tt.role) != tt.expected {
			t.Errorf("TransitionRole constant = %q, want %q", tt.role, tt.expected)
		}
	}
}

// --- ValidateTransition: terminal state protection ---

func TestValidateTransition_CompletedIsTerminal(t *testing.T) {
	states := []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"}
	roles := []TransitionRole{RoleSubmit, RoleClaim, RoleReopen, RoleAuto}

	for _, target := range states {
		for _, role := range roles {
			err := ValidateTransition("completed", target, role)
			if err == nil {
				t.Errorf("ValidateTransition(completed, %s, %s) = nil, want error (completed is terminal)", target, role)
			}
		}
	}
}

func TestValidateTransition_CompletedErrorMessage(t *testing.T) {
	err := ValidateTransition("completed", "pending", RoleSubmit)
	if err == nil {
		t.Fatal("expected error for completed -> pending")
	}
	if !containsStr(err.Error(), "completed") {
		t.Errorf("error message should mention 'completed', got: %v", err)
	}
}

// --- ValidateTransition: rejected state ---

func TestValidateTransition_RejectedReopenToPending(t *testing.T) {
	err := ValidateTransition("rejected", "pending", RoleReopen)
	if err != nil {
		t.Errorf("ValidateTransition(rejected, pending, reopen) = %v, want nil", err)
	}
}

func TestValidateTransition_RejectedBlocksNonReopen(t *testing.T) {
	roles := []TransitionRole{RoleSubmit, RoleClaim, RoleAuto}
	targets := []string{"in_progress", "completed", "blocked", "skipped"}

	for _, role := range roles {
		for _, target := range targets {
			err := ValidateTransition("rejected", target, role)
			if err == nil {
				t.Errorf("ValidateTransition(rejected, %s, %s) = nil, want error", target, role)
			}
		}
	}
}

func TestValidateTransition_RejectedNonPendingBlocked(t *testing.T) {
	err := ValidateTransition("rejected", "in_progress", RoleReopen)
	if err == nil {
		t.Error("ValidateTransition(rejected, in_progress, reopen) should fail, reopen only goes to pending")
	}
}

func TestValidateTransition_RejectedErrorMessage(t *testing.T) {
	err := ValidateTransition("rejected", "in_progress", RoleSubmit)
	if err == nil {
		t.Fatal("expected error for rejected -> in_progress")
	}
	if !containsStr(err.Error(), "rejected") {
		t.Errorf("error message should mention 'rejected', got: %v", err)
	}
}

// --- ValidateTransition: skipped state ---

func TestValidateTransition_SkippedReopenToPending(t *testing.T) {
	err := ValidateTransition("skipped", "pending", RoleReopen)
	if err != nil {
		t.Errorf("ValidateTransition(skipped, pending, reopen) = %v, want nil", err)
	}
}

func TestValidateTransition_SkippedBlocksNonReopen(t *testing.T) {
	roles := []TransitionRole{RoleSubmit, RoleClaim, RoleAuto}
	targets := []string{"in_progress", "completed", "blocked", "rejected"}

	for _, role := range roles {
		for _, target := range targets {
			err := ValidateTransition("skipped", target, role)
			if err == nil {
				t.Errorf("ValidateTransition(skipped, %s, %s) = nil, want error", target, role)
			}
		}
	}
}

func TestValidateTransition_SkippedNonPendingBlocked(t *testing.T) {
	err := ValidateTransition("skipped", "in_progress", RoleReopen)
	if err == nil {
		t.Error("ValidateTransition(skipped, in_progress, reopen) should fail, reopen only goes to pending")
	}
}

// --- ValidateTransition: completed target requires RoleSubmit ---

func TestValidateTransition_OnlySubmitCanReachCompleted(t *testing.T) {
	roles := []TransitionRole{RoleClaim, RoleReopen, RoleAuto}
	sources := []string{"pending", "in_progress", "blocked"}

	for _, role := range roles {
		for _, source := range sources {
			err := ValidateTransition(source, "completed", role)
			if err == nil {
				t.Errorf("ValidateTransition(%s, completed, %s) = nil, want error (only submit can reach completed)", source, role)
			}
		}
	}
}

func TestValidateTransition_SubmitCanReachCompleted(t *testing.T) {
	err := ValidateTransition("in_progress", "completed", RoleSubmit)
	if err != nil {
		t.Errorf("ValidateTransition(in_progress, completed, submit) = %v, want nil", err)
	}
}

// --- ValidateTransition: submit auto-downgrade ---

func TestValidateTransition_SubmitCanDowngradeToBlocked(t *testing.T) {
	err := ValidateTransition("in_progress", "blocked", RoleSubmit)
	if err != nil {
		t.Errorf("ValidateTransition(in_progress, blocked, submit) = %v, want nil", err)
	}
}

// --- ValidateTransition: blocked requires dep check (phase 2) ---

func TestValidateTransition_BlockedToPendingNeedsDeps(t *testing.T) {
	// ValidateTransition is phase 1 (pure state check) — it defers blocked→pending to phase 2.
	// Phase 1 should return a special error indicating dependency check is needed.
	err := ValidateTransition("blocked", "pending", RoleClaim)
	if err == nil {
		// This is also acceptable: phase 1 allows it, phase 2 checks deps.
		// The design says "Dep check (phase 2)" for this transition.
		return
	}
	// If phase 1 returns an error, it should indicate deps need checking
	if !containsStr(err.Error(), "depend") {
		t.Errorf("blocked->pending error should mention dependencies, got: %v", err)
	}
}

func TestValidateTransition_BlockedToInProgressNeedsDeps(t *testing.T) {
	err := ValidateTransition("blocked", "in_progress", RoleClaim)
	if err == nil {
		return // acceptable: phase 2 checks deps
	}
	if !containsStr(err.Error(), "depend") {
		t.Errorf("blocked->in_progress error should mention dependencies, got: %v", err)
	}
}

// --- ValidateTransition: pending -> blocked allowed ---

func TestValidateTransition_PendingToBlocked(t *testing.T) {
	roles := []TransitionRole{RoleSubmit, RoleClaim, RoleReopen, RoleAuto}
	for _, role := range roles {
		err := ValidateTransition("pending", "blocked", role)
		if err != nil {
			t.Errorf("ValidateTransition(pending, blocked, %s) = %v, want nil", role, err)
		}
	}
}

// --- ValidateTransition: same state is no-op (allowed) ---

func TestValidateTransition_SameStateNoop(t *testing.T) {
	states := []string{"pending", "in_progress", "blocked", "completed", "skipped", "rejected"}
	roles := []TransitionRole{RoleSubmit, RoleClaim, RoleReopen, RoleAuto}

	for _, state := range states {
		for _, role := range roles {
			err := ValidateTransition(state, state, role)
			if err != nil {
				// completed same-state is terminal, so it should fail
				if state == "completed" {
					continue
				}
				// rejected/skipped same-state: all roles blocked (terminal states)
				if state == "rejected" || state == "skipped" {
					continue
				}
				// reopen role is only for rejected/skipped, not for non-terminal noop
				if role == RoleReopen {
					continue
				}
				t.Errorf("ValidateTransition(%s, %s, %s) = %v, want nil (same state noop)", state, state, role, err)
			}
		}
	}
}

// --- ValidateTransition: general non-terminal transitions ---

func TestValidateTransition_PendingToInProgress(t *testing.T) {
	err := ValidateTransition("pending", "in_progress", RoleClaim)
	if err != nil {
		t.Errorf("ValidateTransition(pending, in_progress, claim) = %v, want nil", err)
	}
}

func TestValidateTransition_InProgressToBlocked_AutoRole(t *testing.T) {
	err := ValidateTransition("in_progress", "blocked", RoleAuto)
	if err != nil {
		t.Errorf("ValidateTransition(in_progress, blocked, auto) = %v, want nil", err)
	}
}

func TestValidateTransition_InProgressToPending(t *testing.T) {
	err := ValidateTransition("in_progress", "pending", RoleClaim)
	if err != nil {
		t.Errorf("ValidateTransition(in_progress, pending, claim) = %v, want nil", err)
	}
}

func TestValidateTransition_BlockedToRejected(t *testing.T) {
	err := ValidateTransition("blocked", "rejected", RoleSubmit)
	if err != nil {
		t.Errorf("ValidateTransition(blocked, rejected, submit) = %v, want nil", err)
	}
}

func TestValidateTransition_BlockedToSkipped(t *testing.T) {
	err := ValidateTransition("blocked", "skipped", RoleSubmit)
	if err != nil {
		t.Errorf("ValidateTransition(blocked, skipped, submit) = %v, want nil", err)
	}
}

func TestValidateTransition_PendingToRejected(t *testing.T) {
	err := ValidateTransition("pending", "rejected", RoleSubmit)
	if err != nil {
		t.Errorf("ValidateTransition(pending, rejected, submit) = %v, want nil", err)
	}
}

func TestValidateTransition_PendingToSkipped(t *testing.T) {
	err := ValidateTransition("pending", "skipped", RoleSubmit)
	if err != nil {
		t.Errorf("ValidateTransition(pending, skipped, submit) = %v, want nil", err)
	}
}

// --- ValidateTransition: RoleAuto behaves like RoleSubmit for non-terminal ---

func TestValidateTransition_AutoRoleLikeSubmit(t *testing.T) {
	// Auto can downgrade to blocked
	err := ValidateTransition("in_progress", "blocked", RoleAuto)
	if err != nil {
		t.Errorf("ValidateTransition(in_progress, blocked, auto) = %v, want nil", err)
	}

	// Auto cannot reach completed (only submit can)
	err = ValidateTransition("in_progress", "completed", RoleAuto)
	if err == nil {
		t.Error("ValidateTransition(in_progress, completed, auto) should fail")
	}

	// Auto cannot escape terminal states
	err = ValidateTransition("completed", "pending", RoleAuto)
	if err == nil {
		t.Error("ValidateTransition(completed, pending, auto) should fail")
	}
}

// --- CheckTransitionDeps ---

func TestCheckTransitionDeps_AllDepsMet(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Status: "completed"},
		"1.2": {ID: "1.2", Status: "completed"},
		"2.1": {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1", "1.2"}},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("unmet deps = %v, want empty (all met)", unmet)
	}
}

func TestCheckTransitionDeps_SomeDepsUnmet(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Status: "completed"},
		"1.2": {ID: "1.2", Status: "pending"},
		"2.1": {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1", "1.2"}},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 1 || unmet[0] != "1.2" {
		t.Errorf("unmet deps = %v, want [1.2]", unmet)
	}
}

func TestCheckTransitionDeps_SkippedSatisfies(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Status: "skipped"},
		"2.1": {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1"}},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("unmet deps = %v, want empty (skipped satisfies deps)", unmet)
	}
}

func TestCheckTransitionDeps_RejectedDoesNotSatisfy(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Status: "rejected"},
		"2.1": {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1"}},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 1 || unmet[0] != "1.1" {
		t.Errorf("unmet deps = %v, want [1.1] (rejected does not satisfy)", unmet)
	}
}

func TestCheckTransitionDeps_NoDeps(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"2.1": {ID: "2.1", Status: "blocked"},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("unmet deps = %v, want empty (no deps)", unmet)
	}
}

func TestCheckTransitionDeps_TaskNotFound(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{})

	_, err := CheckTransitionDeps(idx, "nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent task")
	}
}

func TestCheckTransitionDeps_BlockedDepDoesNotSatisfy(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Status: "blocked"},
		"2.1": {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1"}},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 1 || unmet[0] != "1.1" {
		t.Errorf("unmet deps = %v, want [1.1] (blocked does not satisfy)", unmet)
	}
}

// --- canAutoUnblock (unexported, tested indirectly) ---

func TestCanAutoUnblock_NoActiveFixTasks(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1": {ID: "1.1", Status: "completed"},
		"2.1": {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1"}, SourceTaskID: "1.1"},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("unmet deps = %v, want empty (no active fix tasks)", unmet)
	}
}

func TestCanAutoUnblock_ActiveFixTaskBlocksUnblock(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1":     {ID: "1.1", Status: "completed"},
		"2.1":     {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1"}, SourceTaskID: "1.1"},
		"T-fix-1": {ID: "T-fix-1", Status: "in_progress", SourceTaskID: "2.1", Type: TypeCodingFix},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With an active fix-task pointing to 2.1, unblock should be blocked
	if len(unmet) == 0 {
		t.Error("expected unmet deps when active fix-task exists, got empty")
	}
}

func TestCanAutoUnblock_CompletedFixTaskDoesNotBlock(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1":     {ID: "1.1", Status: "completed"},
		"2.1":     {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1"}, SourceTaskID: "1.1"},
		"T-fix-1": {ID: "T-fix-1", Status: "completed", SourceTaskID: "2.1", Type: TypeCodingFix},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("unmet deps = %v, want empty (completed fix task does not block)", unmet)
	}
}

func TestCanAutoUnblock_RejectedFixTaskDoesNotBlock(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1":     {ID: "1.1", Status: "completed"},
		"2.1":     {ID: "2.1", Status: "blocked", Dependencies: []string{"1.1"}, SourceTaskID: "1.1"},
		"T-fix-1": {ID: "T-fix-1", Status: "rejected", SourceTaskID: "2.1", Type: TypeCodingFix},
	})

	unmet, err := CheckTransitionDeps(idx, "2.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("unmet deps = %v, want empty (rejected fix task does not block)", unmet)
	}
}

// --- Role isolation: RoleReopen only works on rejected/skipped ---

func TestValidateTransition_ReopenOnNonTerminal(t *testing.T) {
	states := []string{"pending", "in_progress", "blocked"}
	for _, state := range states {
		err := ValidateTransition(state, "pending", RoleReopen)
		if err == nil {
			t.Errorf("ValidateTransition(%s, pending, reopen) should fail (reopen only for rejected/skipped)", state)
		}
	}
}

// --- Full transition matrix: all state x role combinations ---

func TestValidateTransition_FullMatrix(t *testing.T) {
	// Test a comprehensive matrix of transitions
	// Format: {from, to, role, shouldPass}
	cases := []struct {
		from string
		to   string
		role TransitionRole
		pass bool
	}{
		// completed -> anything: always blocked
		{"completed", "pending", RoleSubmit, false},
		{"completed", "in_progress", RoleClaim, false},
		{"completed", "blocked", RoleAuto, false},
		{"completed", "completed", RoleSubmit, false},

		// rejected -> pending via reopen only
		{"rejected", "pending", RoleReopen, true},
		{"rejected", "pending", RoleSubmit, false},
		{"rejected", "in_progress", RoleReopen, false},

		// skipped -> pending via reopen only
		{"skipped", "pending", RoleReopen, true},
		{"skipped", "pending", RoleSubmit, false},
		{"skipped", "in_progress", RoleReopen, false},

		// -> completed only via submit
		{"in_progress", "completed", RoleSubmit, true},
		{"pending", "completed", RoleSubmit, true},
		{"pending", "completed", RoleClaim, false},
		{"in_progress", "completed", RoleAuto, false},
		{"blocked", "completed", RoleSubmit, true},

		// submit can downgrade in_progress -> blocked
		{"in_progress", "blocked", RoleSubmit, true},

		// pending -> blocked (any role)
		{"pending", "blocked", RoleSubmit, true},
		{"pending", "blocked", RoleClaim, true},

		// general non-terminal transitions
		{"pending", "in_progress", RoleClaim, true},
		{"pending", "rejected", RoleSubmit, true},
		{"pending", "skipped", RoleSubmit, true},
		{"in_progress", "pending", RoleClaim, true},
		{"in_progress", "rejected", RoleSubmit, true},
		{"blocked", "rejected", RoleSubmit, true},
		{"blocked", "skipped", RoleSubmit, true},


			// manual (RoleManual) overrides
			{"blocked", "pending", RoleManual, true},
			{"blocked", "in_progress", RoleManual, true},
			{"blocked", "skipped", RoleManual, true},
			{"blocked", "rejected", RoleManual, true},
			{"pending", "skipped", RoleManual, true},
			{"pending", "rejected", RoleManual, true},
			{"in_progress", "blocked", RoleManual, true},
			{"in_progress", "skipped", RoleManual, true},
			{"in_progress", "rejected", RoleManual, true},
			{"completed", "pending", RoleManual, false},
			{"rejected", "pending", RoleManual, false},
			{"skipped", "pending", RoleManual, false},
			{"pending", "completed", RoleManual, false},
	}

	for _, tc := range cases {
		err := ValidateTransition(tc.from, tc.to, tc.role)
		if tc.pass && err != nil {
			t.Errorf("ValidateTransition(%s, %s, %s) = %v, want nil", tc.from, tc.to, tc.role, err)
		}
		if !tc.pass && err == nil {
			t.Errorf("ValidateTransition(%s, %s, %s) = nil, want error", tc.from, tc.to, tc.role)
		}
	}
}

// --- Error type check ---

func TestValidateTransition_ErrorType(t *testing.T) {
	err := ValidateTransition("completed", "pending", RoleSubmit)
	if err == nil {
		t.Fatal("expected error")
	}

	// Should be a TransitionError
	var te *TransitionError
	if !errors.As(err, &te) {
		t.Errorf("error should be *TransitionError, got %T", err)
	}
	if te.From != "completed" {
		t.Errorf("TransitionError.From = %q, want %q", te.From, "completed")
	}
	if te.To != "pending" {
		t.Errorf("TransitionError.To = %q, want %q", te.To, "pending")
	}
	if te.Role != RoleSubmit {
		t.Errorf("TransitionError.Role = %q, want %q", te.Role, RoleSubmit)
	}
}

// --- No force parameter anywhere (design constraint) ---

func TestNoForceParameter(_ *testing.T) {
	// This is a compile-time check: ValidateTransition does not accept a force parameter.
	// The function signature is ValidateTransition(current, target string, role TransitionRole) error
	// No test needed, this is enforced by the API.
}

// --- RoleManual: manual override transitions ---

func TestValidateTransition_ManualUnblock(t *testing.T) {
	if err := ValidateTransition("blocked", "pending", RoleManual); err != nil {
		t.Errorf("blocked -> pending (manual) should be allowed, got: %v", err)
	}
}

func TestValidateTransition_ManualUnblockToInProgress(t *testing.T) {
	if err := ValidateTransition("blocked", "in_progress", RoleManual); err != nil {
		t.Errorf("blocked -> in_progress (manual) should be allowed, got: %v", err)
	}
}

func TestValidateTransition_ManualCompletedBlocked(t *testing.T) {
	if err := ValidateTransition("pending", "completed", RoleManual); err == nil {
		t.Error("pending -> completed (manual) should be blocked — use submit")
	}
}

func TestValidateTransition_ManualSkip(t *testing.T) {
	if err := ValidateTransition("blocked", "skipped", RoleManual); err != nil {
		t.Errorf("blocked -> skipped (manual) should be allowed, got: %v", err)
	}
	if err := ValidateTransition("in_progress", "skipped", RoleManual); err != nil {
		t.Errorf("in_progress -> skipped (manual) should be allowed, got: %v", err)
	}
	if err := ValidateTransition("pending", "skipped", RoleManual); err != nil {
		t.Errorf("pending -> skipped (manual) should be allowed, got: %v", err)
	}
}

func TestValidateTransition_ManualReject(t *testing.T) {
	if err := ValidateTransition("blocked", "rejected", RoleManual); err != nil {
		t.Errorf("blocked -> rejected (manual) should be allowed, got: %v", err)
	}
	if err := ValidateTransition("in_progress", "rejected", RoleManual); err != nil {
		t.Errorf("in_progress -> rejected (manual) should be allowed, got: %v", err)
	}
}

func TestValidateTransition_ManualCompletedTerminal(t *testing.T) {
	if err := ValidateTransition("completed", "skipped", RoleManual); err == nil {
		t.Error("completed -> skipped (manual) should be blocked — completed is terminal")
	}
	if err := ValidateTransition("completed", "pending", RoleManual); err == nil {
		t.Error("completed -> pending (manual) should be blocked — completed is terminal")
	}
}

func TestValidateTransition_SuspendedEntry(t *testing.T) {
	if err := ValidateTransition("pending", "suspended", RoleManual); err != nil {
		t.Errorf("pending -> suspended (manual) should be allowed, got: %v", err)
	}
	if err := ValidateTransition("in_progress", "suspended", RoleManual); err != nil {
		t.Errorf("in_progress -> suspended (manual) should be allowed, got: %v", err)
	}
	if err := ValidateTransition("pending", "suspended", RoleAuto); err == nil {
		t.Error("pending -> suspended (auto) should be blocked")
	}
	if err := ValidateTransition("in_progress", "suspended", RoleSubmit); err == nil {
		t.Error("in_progress -> suspended (submit) should be blocked")
	}
}

func TestValidateTransition_SuspendedResume(t *testing.T) {
	if err := ValidateTransition("suspended", "pending", RoleManual); err != nil {
		t.Errorf("suspended -> pending (manual) should be allowed, got: %v", err)
	}
	if err := ValidateTransition("suspended", "in_progress", RoleManual); err != nil {
		t.Errorf("suspended -> in_progress (manual) should be allowed, got: %v", err)
	}
}

func TestValidateTransition_SuspendedTerminal(t *testing.T) {
	if err := ValidateTransition("suspended", "skipped", RoleManual); err != nil {
		t.Errorf("suspended -> skipped (manual) should be allowed, got: %v", err)
	}
	if err := ValidateTransition("suspended", "rejected", RoleManual); err != nil {
		t.Errorf("suspended -> rejected (manual) should be allowed, got: %v", err)
	}
}

func TestValidateTransition_SuspendedBlockedBySystem(t *testing.T) {
	if err := ValidateTransition("suspended", "blocked", RoleAuto); err == nil {
		t.Error("suspended -> blocked (auto) should be blocked")
	}
	if err := ValidateTransition("suspended", "blocked", RoleSubmit); err == nil {
		t.Error("suspended -> blocked (submit) should be blocked")
	}
	if err := ValidateTransition("suspended", "completed", RoleSubmit); err == nil {
		t.Error("suspended -> completed (submit) should be blocked — must resume first")
	}
}

func TestValidateTransition_SuspendedReopen(t *testing.T) {
	if err := ValidateTransition("suspended", "pending", RoleReopen); err == nil {
		t.Error("suspended -> pending (reopen) should be blocked — reopen only for rejected/skipped")
	}
}

func TestValidateTransition_SuspendedSameState(t *testing.T) {
	if err := ValidateTransition("suspended", "suspended", RoleManual); err != nil {
		t.Errorf("suspended -> suspended (same state noop) should be allowed, got: %v", err)
	}
}

// --- Helper ---

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(len(s) > 0 && containsStrHelper(s, sub)))
}

func containsStrHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
