package task

import (
	"testing"
)

func TestPreserveRuntimeFields(t *testing.T) {
	tests := []struct {
		name     string
		existing *Task
		newTask  Task
		expected Task
	}{
		{
			name:     "preserves all runtime fields including Dependencies",
			existing: &Task{Status: "in_progress", SourceTaskID: "source-1", BlockedReason: "waiting for review", Dependencies: []string{"fix-1"}},
			newTask:  Task{ID: "1", Title: "Test", Status: "pending"},
			expected: Task{ID: "1", Title: "Test", Status: "in_progress", SourceTaskID: "source-1", BlockedReason: "waiting for review", Dependencies: []string{"fix-1"}},
		},
		{
			name:     "nil existing does nothing",
			existing: nil,
			newTask:  Task{ID: "1", Status: "pending"},
			expected: Task{ID: "1", Status: "pending"},
		},
		{
			name:     "partial preservation - empty strings should still copy",
			existing: &Task{Status: "completed", SourceTaskID: "", BlockedReason: ""},
			newTask:  Task{ID: "1", Status: "pending"},
			expected: Task{ID: "1", Status: "completed", SourceTaskID: "", BlockedReason: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newCopy := tt.newTask
			PreserveRuntimeFields(tt.existing, &newCopy)
			if newCopy.Status != tt.expected.Status {
				t.Errorf("Status = %q, want %q", newCopy.Status, tt.expected.Status)
			}
			if newCopy.SourceTaskID != tt.expected.SourceTaskID {
				t.Errorf("SourceTaskID = %q, want %q", newCopy.SourceTaskID, tt.expected.SourceTaskID)
			}
			if newCopy.BlockedReason != tt.expected.BlockedReason {
				t.Errorf("BlockedReason = %q, want %q", newCopy.BlockedReason, tt.expected.BlockedReason)
			}
			// Dependencies preservation
			if len(newCopy.Dependencies) != len(tt.expected.Dependencies) {
				t.Errorf("Dependencies = %v, want %v", newCopy.Dependencies, tt.expected.Dependencies)
			} else {
				for i, d := range newCopy.Dependencies {
					if d != tt.expected.Dependencies[i] {
						t.Errorf("Dependencies[%d] = %q, want %q", i, d, tt.expected.Dependencies[i])
					}
				}
			}
			// Ensure fields NOT in PreserveRuntimeFields are not affected
			if newCopy.ID != tt.newTask.ID {
				t.Errorf("ID changed from %q to %q", tt.newTask.ID, newCopy.ID)
			}
			if newCopy.Title != tt.newTask.Title {
				t.Errorf("Title changed from %q to %q", tt.newTask.Title, newCopy.Title)
			}
		})
	}
}
