package cmd

import (
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
		// Must use task record
		{"in_progress", "completed", false},
		{"pending", "completed", false},
		{"blocked", "completed", false},
		{"skipped", "completed", false},
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
				"a":   {ID: "a", Dependencies: []string{"b"}},
				"b":   {ID: "b", Status: "completed"},
			},
			deps:      []string{"b"},
			wantUnmet: 0,
		},
		{
			name: "exact dep pending",
			tasks: map[string]task.Task{
				"a":   {ID: "a", Dependencies: []string{"b"}},
				"b":   {ID: "b", Status: "pending"},
			},
			deps:      []string{"b"},
			wantUnmet: 1,
		},
		{
			name: "wildcard all completed",
			tasks: map[string]task.Task{
				"a":     {ID: "a"},
				"1.1":   {ID: "1.1", Status: "completed"},
				"1.2":   {ID: "1.2", Status: "completed"},
				"1.gate": {ID: "1.gate", Status: "pending"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 0,
		},
		{
			name: "wildcard some pending",
			tasks: map[string]task.Task{
				"a":     {ID: "a"},
				"1.1":   {ID: "1.1", Status: "completed"},
				"1.2":   {ID: "1.2", Status: "pending"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 1,
		},
		{
			name: "wildcard ignores gate and summary",
			tasks: map[string]task.Task{
				"a":        {ID: "a"},
				"1.1":      {ID: "1.1", Status: "completed"},
				"1.gate":   {ID: "1.gate", Status: "pending"},
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
			name: "wildcard self-exclusion",
			subjectID: "1.3",
			tasks: map[string]task.Task{
				"1.3":   {ID: "1.3", Status: "pending"},
				"1.1":   {ID: "1.1", Status: "completed"},
				"1.2":   {ID: "1.2", Status: "completed"},
			},
			deps:      []string{"1.x"},
			wantUnmet: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := &task.TaskIndex{Feature: "test", Tasks: tt.tasks}
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
