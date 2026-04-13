package cmd

import (
	"testing"
)

func TestErrInvalidStatus_EmptyValidStatuses(t *testing.T) {
	// Should not panic when validStatuses is empty (missing statusEnum in index.json)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("ErrInvalidStatus panicked with empty validStatuses: %v", r)
		}
	}()

	err := ErrInvalidStatus("completed", []string{})
	if err == nil {
		t.Error("expected non-nil error")
	}
}

func TestErrInvalidStatus_WithValidStatuses(t *testing.T) {
	err := ErrInvalidStatus("bad_status", []string{"pending", "completed"})
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Code != ErrValidation {
		t.Errorf("Code = %q, want %q", err.Code, ErrValidation)
	}
	if err.Action != "task status <id> pending" {
		t.Errorf("Action = %q, want %q", err.Action, "task status <id> pending")
	}
}
