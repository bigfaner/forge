package cmd

import (
	"strings"
	"testing"

	"task-cli/pkg/task"
)

func TestIsTransitionAllowed(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want bool
	}{
		// Same status (idempotent)
		{"pending", "pending", true},
		{"in_progress", "in_progress", true},
		{"completed", "completed", true},
		{"blocked", "blocked", true},
		// Valid forward transitions
		{"pending", "in_progress", true},
		{"pending", "blocked", true},
		{"pending", "skipped", true},
		{"in_progress", "blocked", true},
		{"in_progress", "pending", true},
		{"in_progress", "skipped", true},
		{"blocked", "pending", true},
		{"blocked", "skipped", true},
		{"skipped", "pending", true},
		// Terminal state guards
		{"completed", "pending", false},
		{"completed", "in_progress", false},
		{"completed", "blocked", false},
		{"completed", "skipped", false},
		{"completed", "rejected", false},
		// Rejected is terminal
		{"rejected", "pending", false},
		{"rejected", "in_progress", false},
		{"rejected", "blocked", false},
		{"rejected", "skipped", false},
		{"rejected", "completed", false},
		{"rejected", "rejected", true},
		// Must use task record
		{"in_progress", "completed", false},
		{"pending", "completed", false},
		{"blocked", "completed", false},
		{"skipped", "completed", false},
		// Transitions to rejected allowed from non-terminal states
		{"pending", "rejected", true},
		{"in_progress", "rejected", true},
		{"blocked", "rejected", true},
		{"skipped", "rejected", true},
	}

	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			got := isTransitionAllowed(tt.from, tt.to)
			if got != tt.want {
				t.Errorf("isTransitionAllowed(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestCheckUnmetDeps_Wildcard(t *testing.T) {
	tests := []struct {
		name      string
		subjectID string
		tasks     map[string]task.Task
		deps      []string
		wantUnmet int
	}{
		{
			name: "exact dep completed",
			tasks: map[string]task.Task{
				"a": {ID: "a", Dependencies: []string{"b"}},
				"b": {ID: "b", Status: "completed"},
			},
			deps:      []string{"b"},
			wantUnmet: 0,
		},
		{
			name: "exact dep pending",
			tasks: map[string]task.Task{
				"a": {ID: "a", Dependencies: []string{"b"}},
				"b": {ID: "b", Status: "pending"},
			},
			deps:      []string{"b"},
			wantUnmet: 1,
		},
		{
			name: "wildcard all completed",
			tasks: map[string]task.Task{
				"a":      {ID: "a"},
				"1.1":    {ID: "1.1", Status: "completed"},
				"1.2":    {ID: "1.2", Status: "completed"},
				"1.gate": {ID: "1.gate", Status: "pending"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 0,
		},
		{
			name: "wildcard some pending",
			tasks: map[string]task.Task{
				"a":   {ID: "a"},
				"1.1": {ID: "1.1", Status: "completed"},
				"1.2": {ID: "1.2", Status: "pending"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 1,
		},
		{
			name: "wildcard ignores gate and summary",
			tasks: map[string]task.Task{
				"a":         {ID: "a"},
				"1.1":       {ID: "1.1", Status: "completed"},
				"1.gate":    {ID: "1.gate", Status: "pending"},
				"1.summary": {ID: "1.summary", Status: "pending"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 0,
		},
		{
			name: "wildcard no matches",
			tasks: map[string]task.Task{
				"a":   {ID: "a"},
				"2.1": {ID: "2.1", Status: "pending"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 0,
		},
		{
			name:      "wildcard self-exclusion",
			subjectID: "1.3",
			tasks: map[string]task.Task{
				"1.3": {ID: "1.3", Status: "pending"},
				"1.1": {ID: "1.1", Status: "completed"},
				"1.2": {ID: "1.2", Status: "completed"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := &task.TaskIndex{Feature: "test"}
			index.SetTasks(tt.tasks)
			subjectID := tt.subjectID
			if subjectID == "" {
				subjectID = "a"
			}
			taskWithDeps := &task.Task{ID: subjectID, Dependencies: tt.deps}
			unmet := checkUnmetDeps(index, taskWithDeps)
			if len(unmet) != tt.wantUnmet {
				t.Errorf("got %d unmet (%v), want %d", len(unmet), unmet, tt.wantUnmet)
			}
		})
	}
}

func TestCheckUnmetDeps_KeyDiffersFromID(t *testing.T) {
	t.Run("slug-keyed completed dep resolved by ID", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":           {ID: "src", Dependencies: []string{"T-test-3"}},
			"run-e2e-tests": {ID: "T-test-3", Status: "completed"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"T-test-3"}})
		if len(unmet) != 0 {
			t.Errorf("expected 0 unmet for slug-keyed completed dep, got %v", unmet)
		}
	})

	t.Run("slug-keyed pending dep reported as unmet", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":           {ID: "src", Dependencies: []string{"T-test-3"}},
			"run-e2e-tests": {ID: "T-test-3", Status: "pending"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"T-test-3"}})
		if len(unmet) != 1 || unmet[0] != "T-test-3" {
			t.Errorf("expected [T-test-3] unmet, got %v", unmet)
		}
	})

	t.Run("mixed slug-keyed and dynamic deps", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":           {ID: "src"},
			"run-e2e-tests": {ID: "T-test-3", Status: "completed"},
			"disc-1":        {ID: "disc-1", Status: "completed"},
			"fix-1":         {ID: "fix-1", Status: "pending"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"T-test-3", "disc-1", "fix-1"}})
		if len(unmet) != 1 || unmet[0] != "fix-1" {
			t.Errorf("expected only fix-1 unmet, got %v", unmet)
		}
	})

	t.Run("nonexistent dep treated as unmet", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src": {ID: "src"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"ghost"}})
		if len(unmet) != 1 || unmet[0] != "ghost" {
			t.Errorf("expected [ghost] unmet, got %v", unmet)
		}
	})

	t.Run("slug-keyed skipped dep counts as met", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":           {ID: "src"},
			"run-e2e-tests": {ID: "T-test-3", Status: "skipped"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"T-test-3"}})
		if len(unmet) != 0 {
			t.Errorf("skipped dep should be met, got %v", unmet)
		}
	})
}

func TestCheckUnmetDeps_KeyDiffersFromID_EdgeCases(t *testing.T) {
	t.Run("empty deps returns empty", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src": {ID: "src"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{}})
		if len(unmet) != 0 {
			t.Errorf("empty deps should return 0 unmet, got %v", unmet)
		}
	})

	t.Run("self-referencing dep treated as unmet", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e": {ID: "T-test-3", Status: "blocked"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "T-test-3", Dependencies: []string{"T-test-3"}})
		if len(unmet) != 1 || unmet[0] != "T-test-3" {
			t.Errorf("self-dep should be unmet (blocked), got %v", unmet)
		}
	})

	t.Run("pure slug-keyed deps some met some not", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e":   {ID: "T-test-3", Status: "completed"},
			"run-smoke": {ID: "T-test-7", Status: "pending"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"T-test-3", "T-test-7"}})
		if len(unmet) != 1 || unmet[0] != "T-test-7" {
			t.Errorf("expected only T-test-7 unmet, got %v", unmet)
		}
	})

	t.Run("dep resolved by slug key not by ID", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e": {ID: "T-test-3", Status: "completed"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"run-e2e"}})
		if len(unmet) != 0 {
			t.Errorf("dep by slug key should be found, got %v", unmet)
		}
	})

	t.Run("wildcard combined with slug-keyed exact dep", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e": {ID: "T-test-3", Status: "completed"},
			"1.1":     {ID: "1.1", Status: "completed"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"T-test-3", "1.x"}})
		if len(unmet) != 0 {
			t.Errorf("both slug-keyed exact and wildcard should be met, got %v", unmet)
		}
	})

	t.Run("multiple unmet deps reported correctly", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e":   {ID: "T-test-3", Status: "pending"},
			"run-smoke": {ID: "T-test-7", Status: "pending"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "src", Dependencies: []string{"T-test-3", "T-test-7"}})
		if len(unmet) != 2 {
			t.Errorf("expected 2 unmet, got %d: %v", len(unmet), unmet)
		}
	})
}

func TestGetTransitionHint(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want string
	}{
		{"completed", "pending", "completed is a terminal state"},
		{"completed", "in_progress", "completed is a terminal state"},
		{"rejected", "pending", "rejected is a terminal state"},
		{"rejected", "in_progress", "rejected is a terminal state"},
		{"in_progress", "completed", "use 'task record' to complete a task with quality gate"},
		{"pending", "completed", "use 'task record' to complete a task with quality gate"},
		{"pending", "in_progress", "transition pending -> in_progress is not allowed"},
		{"blocked", "pending", "transition blocked -> pending is not allowed"},
	}
	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			got := getTransitionHint(tt.from, tt.to)
			if !strings.Contains(got, tt.want) {
				t.Errorf("getTransitionHint(%q, %q) = %q, want to contain %q", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestGetTransitionAction(t *testing.T) {
	tests := []struct {
		from string
		to   string
		want string
	}{
		{"completed", "pending", "use --force to override"},
		{"completed", "blocked", "use --force to override"},
		{"rejected", "pending", "use --force to override"},
		{"rejected", "blocked", "use --force to override"},
		{"in_progress", "completed", "task record"},
		{"pending", "completed", "task record"},
		{"pending", "in_progress", "use --force to override"},
	}
	for _, tt := range tests {
		t.Run(tt.from+"->"+tt.to, func(t *testing.T) {
			got := getTransitionAction(tt.from, tt.to)
			if !strings.Contains(got, tt.want) {
				t.Errorf("getTransitionAction(%q, %q) = %q, want to contain %q", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestCheckUnmetDeps_RejectedDepNotSatisfied(t *testing.T) {
	t.Run("rejected exact dep is unmet", func(t *testing.T) {
		index := &task.TaskIndex{Feature: "test"}
		index.SetTasks(map[string]task.Task{
			"a": {ID: "a", Dependencies: []string{"b"}},
			"b": {ID: "b", Status: "rejected"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "a", Dependencies: []string{"b"}})
		if len(unmet) != 1 || unmet[0] != "b" {
			t.Errorf("rejected dep should be unmet, got %v", unmet)
		}
	})

	t.Run("rejected wildcard dep is unmet", func(t *testing.T) {
		index := &task.TaskIndex{Feature: "test"}
		index.SetTasks(map[string]task.Task{
			"a":   {ID: "a"},
			"1.1": {ID: "1.1", Status: "completed"},
			"1.2": {ID: "1.2", Status: "rejected"},
		})
		unmet := checkUnmetDeps(index, &task.Task{ID: "a", Dependencies: []string{"1.x"}})
		if len(unmet) != 1 || unmet[0] != "1.2" {
			t.Errorf("wildcard should report rejected task as unmet, got %v", unmet)
		}
	})
}
