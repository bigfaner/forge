package task

import (
	"testing"
)

func TestCheckLegacyScope_NoLegacy(t *testing.T) {
	tasks := []Task{
		{ID: "1.1", SurfaceKey: "admin-panel", SurfaceType: "web"},
		{ID: "1.2", SurfaceKey: "", SurfaceType: ""},
		{ID: "1.3", SurfaceKey: "payment-api", SurfaceType: "api"},
	}
	if err := CheckLegacyScope(tasks); err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}
}

func TestCheckLegacyScope_LegacyScopeDetected(t *testing.T) {
	tasks := []Task{
		{ID: "1.1", Scope: "frontend", SurfaceKey: "", SurfaceType: ""},
		{ID: "1.2", SurfaceKey: "admin-panel", SurfaceType: "web"},
		{ID: "1.3", Scope: "backend", SurfaceKey: "", SurfaceType: ""},
	}
	err := CheckLegacyScope(tasks)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	scopeErr, ok := err.(*LegacyScopeError)
	if !ok {
		t.Fatalf("expected *LegacyScopeError, got %T", err)
	}
	if scopeErr.Count != 2 {
		t.Errorf("expected count 2, got %d", scopeErr.Count)
	}

	expectedMsg := "migration required: found 2 tasks with legacy 'scope' field but no 'surface-key' — run 'forge task migrate' or 'forge breakdown-tasks' to regenerate tasks"
	if scopeErr.Error() != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, scopeErr.Error())
	}
}

func TestCheckLegacyScope_ScopeWithSurfaceKey_NoError(t *testing.T) {
	// If a task has both scope and surface-key, it's already migrated — no error.
	tasks := []Task{
		{ID: "1.1", Scope: "frontend", SurfaceKey: "admin-panel", SurfaceType: "web"},
	}
	if err := CheckLegacyScope(tasks); err != nil {
		t.Errorf("expected nil for task with surface-key set, got error: %v", err)
	}
}

func TestCheckLegacyScope_EmptySlice(t *testing.T) {
	if err := CheckLegacyScope(nil); err != nil {
		t.Errorf("expected nil for empty slice, got error: %v", err)
	}
	if err := CheckLegacyScope([]Task{}); err != nil {
		t.Errorf("expected nil for empty slice, got error: %v", err)
	}
}

func TestLegacyScopeError_MessageFormat(t *testing.T) {
	err := &LegacyScopeError{Count: 5}
	expected := "migration required: found 5 tasks with legacy 'scope' field but no 'surface-key' — run 'forge task migrate' or 'forge breakdown-tasks' to regenerate tasks"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
