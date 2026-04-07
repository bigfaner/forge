package version

import "testing"

func TestGetVersion(t *testing.T) {
	v := GetVersion()
	if v == "" {
		t.Error("expected non-empty version")
	}
}

func TestGetName(t *testing.T) {
	n := GetName()
	if n != "task" {
		t.Errorf("expected name 'task', got %q", n)
	}
}

func TestVersionDefault(t *testing.T) {
	// Default version should be "dev"
	if Version != "dev" {
		t.Errorf("expected default version 'dev', got %q", Version)
	}
}
