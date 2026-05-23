package task

import (
	"testing"
)

func TestResolveWildcardDep_Exact(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1-setup": {ID: "1.1-setup", Status: "completed"},
	})

	matches, isWildcard := ResolveWildcardDep(idx, "1.1-setup")
	if isWildcard {
		t.Error("expected isWildcard=false for exact dep")
	}
	if len(matches) != 1 || matches[0] != "1.1-setup" {
		t.Errorf("expected [1.1-setup], got %v", matches)
	}
}

func TestResolveWildcardDep_Wildcard(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1-setup":  {ID: "1.1-setup", Status: "pending"},
		"1.2-build":  {ID: "1.2-build", Status: "completed"},
		"1.gate":     {ID: "1.gate", Status: "pending"},
		"1.summary":  {ID: "1.summary", Status: "pending"},
		"2.1-deploy": {ID: "2.1-deploy", Status: "pending"},
	})

	matches, isWildcard := ResolveWildcardDep(idx, "1.x")
	if !isWildcard {
		t.Error("expected isWildcard=true for wildcard dep")
	}
	// Should only match business tasks (exclude gate/summary)
	expected := map[string]bool{"1.1-setup": true, "1.2-build": true}
	if len(matches) != len(expected) {
		t.Errorf("expected %d matches, got %d: %v", len(expected), len(matches), matches)
	}
	for _, m := range matches {
		if !expected[m] {
			t.Errorf("unexpected match: %s", m)
		}
	}
}

func TestResolveWildcardDep_WildcardNoMatch(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1-setup": {ID: "1.1-setup", Status: "pending"},
	})

	matches, isWildcard := ResolveWildcardDep(idx, "5.x")
	if !isWildcard {
		t.Error("expected isWildcard=true")
	}
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %v", matches)
	}
}

func TestGetUnmetDeps_AllCompleted(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1-setup": {ID: "1.1-setup", Status: "completed"},
		"1.2-build": {ID: "1.2-build", Status: "skipped"},
	})

	unmet := GetUnmetDeps(idx, "2.1-deploy", []string{"1.1-setup", "1.2-build"})
	if len(unmet) != 0 {
		t.Errorf("expected no unmet, got %v", unmet)
	}
}

func TestGetUnmetDeps_SomeUnmet(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1-setup": {ID: "1.1-setup", Status: "completed"},
		"1.2-build": {ID: "1.2-build", Status: "pending"},
	})

	unmet := GetUnmetDeps(idx, "2.1-deploy", []string{"1.1-setup", "1.2-build"})
	if len(unmet) != 1 || unmet[0] != "1.2-build" {
		t.Errorf("expected [1.2-build], got %v", unmet)
	}
}

func TestGetUnmetDeps_WildcardSelfExcluded(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1-setup": {ID: "1.1-setup", Status: "completed"},
		"1.2-build": {ID: "1.2-build", Status: "pending"},
	})

	// selfID matches wildcard but should be excluded
	unmet := GetUnmetDeps(idx, "1.1-setup", []string{"1.x"})
	if len(unmet) != 1 || unmet[0] != "1.2-build" {
		t.Errorf("expected [1.2-build] (self excluded), got %v", unmet)
	}
}

func TestGetUnmetDeps_MissingExactDep(t *testing.T) {
	idx := NewTestIndex("test", map[string]Task{
		"1.1-setup": {ID: "1.1-setup", Status: "completed"},
	})

	unmet := GetUnmetDeps(idx, "2.1-deploy", []string{"1.1-setup", "9.9-missing"})
	if len(unmet) != 1 || unmet[0] != "9.9-missing" {
		t.Errorf("expected [9.9-missing], got %v", unmet)
	}
}

func TestIsDepSatisfied(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{"completed", true},
		{"skipped", true},
		{"pending", false},
		{"in_progress", false},
		{"blocked", false},
		{"rejected", false},
		{"suspended", false},
	}
	for _, tt := range tests {
		if got := IsDepSatisfied(tt.status); got != tt.expected {
			t.Errorf("IsDepSatisfied(%q) = %v, want %v", tt.status, got, tt.expected)
		}
	}
}
