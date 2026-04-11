package task

import (
	"encoding/json"
	"testing"
)

func TestNewTaskIndex(t *testing.T) {
	tests := []struct {
		name    string
		feature string
	}{
		{"basic feature", "my-feature"},
		{"empty feature", ""},
		{"feature with special chars", "feature_123-test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := NewTaskIndex(tt.feature)
			if index == nil {
				t.Fatal("NewTaskIndex returned nil")
			}
			if index.Feature != tt.feature {
				t.Errorf("Feature = %q, want %q", index.Feature, tt.feature)
			}
			if index.Tasks == nil {
				t.Error("Tasks map is nil")
			}
			if len(index.Tasks) != 0 {
				t.Errorf("Tasks map should be empty, got %d items", len(index.Tasks))
			}
			// Check default status enum
			expectedStatuses := []string{"pending", "in_progress", "completed", "blocked", "skipped"}
			if len(index.StatusEnum) != len(expectedStatuses) {
				t.Errorf("StatusEnum length = %d, want %d", len(index.StatusEnum), len(expectedStatuses))
			}
			// Check default priority enum
			expectedPriorities := []string{"P0", "P1", "P2"}
			if len(index.PriorityEnum) != len(expectedPriorities) {
				t.Errorf("PriorityEnum length = %d, want %d", len(index.PriorityEnum), len(expectedPriorities))
			}
			// Check Created field format (YYYY-MM-DD)
			if index.Created == "" {
				t.Error("Created field is empty")
			}
			if len(index.Created) != 10 {
				t.Errorf("Created format = %q, want YYYY-MM-DD format", index.Created)
			}
		})
	}
}

func TestTaskJSONRoundTrip(t *testing.T) {
	task := Task{
		ID:            "1.1",
		Title:         "Test Task",
		Priority:      "P0",
		EstimatedTime: "2h",
		Dependencies:  []string{"1.0"},
		Status:        "pending",
		File:          "tasks/1.1.md",
		Record:        "records/1.1.md",
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var unmarshaled Task
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if unmarshaled.ID != task.ID {
		t.Errorf("ID = %q, want %q", unmarshaled.ID, task.ID)
	}
	if unmarshaled.Title != task.Title {
		t.Errorf("Title = %q, want %q", unmarshaled.Title, task.Title)
	}
	if unmarshaled.Priority != task.Priority {
		t.Errorf("Priority = %q, want %q", unmarshaled.Priority, task.Priority)
	}
	if len(unmarshaled.Dependencies) != len(task.Dependencies) {
		t.Errorf("Dependencies length = %d, want %d", len(unmarshaled.Dependencies), len(task.Dependencies))
	}
}

func TestTaskIndexJSONRoundTrip(t *testing.T) {
	index := &TaskIndex{
		Feature: "test-feature",
		PRD:     "prd/prd-spec.md",
		Design:  "design/tech-design.md",
		Created: "2024-01-01",
		Status:  "in_progress",
		Tasks: map[string]Task{
			"task1": {
				ID:       "1.1",
				Title:    "First Task",
				Priority: "P0",
				Status:   "pending",
				File:     "tasks/1.1.md",
				Record:   "records/1.1.md",
			},
		},
		StatusEnum:   []string{"pending", "in_progress", "completed"},
		PriorityEnum: []string{"P0", "P1", "P2"},
	}

	data, err := json.Marshal(index)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var unmarshaled TaskIndex
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if unmarshaled.Feature != index.Feature {
		t.Errorf("Feature = %q, want %q", unmarshaled.Feature, index.Feature)
	}
	if len(unmarshaled.Tasks) != len(index.Tasks) {
		t.Errorf("Tasks count = %d, want %d", len(unmarshaled.Tasks), len(index.Tasks))
	}
}

func TestTaskStateJSONRoundTrip(t *testing.T) {
	state := &TaskState{
		TaskID:        "1.1",
		Key:           "task1",
		Title:         "Test Task",
		Priority:      "P0",
		EstimatedTime: "2h",
		Dependencies:  []string{"1.0"},
		File:          "tasks/1.1.md",
		Record:        "records/1.1.md",
		StartedTime:   "2024-01-01 10:00",
	}

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var unmarshaled TaskState
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if unmarshaled.TaskID != state.TaskID {
		t.Errorf("TaskID = %q, want %q", unmarshaled.TaskID, state.TaskID)
	}
	if unmarshaled.Key != state.Key {
		t.Errorf("Key = %q, want %q", unmarshaled.Key, state.Key)
	}
	if unmarshaled.StartedTime != state.StartedTime {
		t.Errorf("StartedTime = %q, want %q", unmarshaled.StartedTime, state.StartedTime)
	}
}

func TestRecordDataJSONRoundTrip(t *testing.T) {
	rd := &RecordData{
		Status:        "completed",
		Summary:       "Task completed successfully",
		FilesCreated:  []string{"file1.go", "file2.go"},
		FilesModified: []string{"file3.go"},
		KeyDecisions:  []string{"Use pattern X"},
		TestsPassed:   10,
		TestsFailed:   0,
		Coverage:      85.5,
		AcceptanceCriteria: []AcceptanceCriterion{
			{Criterion: "Feature works", Met: true},
			{Criterion: "Tests pass", Met: true},
		},
		Notes: "No issues",
	}

	data, err := json.Marshal(rd)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var unmarshaled RecordData
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if unmarshaled.Status != rd.Status {
		t.Errorf("Status = %q, want %q", unmarshaled.Status, rd.Status)
	}
	if unmarshaled.Summary != rd.Summary {
		t.Errorf("Summary = %q, want %q", unmarshaled.Summary, rd.Summary)
	}
	if unmarshaled.TestsPassed != rd.TestsPassed {
		t.Errorf("TestsPassed = %d, want %d", unmarshaled.TestsPassed, rd.TestsPassed)
	}
	if unmarshaled.Coverage != rd.Coverage {
		t.Errorf("Coverage = %f, want %f", unmarshaled.Coverage, rd.Coverage)
	}
	if len(unmarshaled.AcceptanceCriteria) != len(rd.AcceptanceCriteria) {
		t.Errorf("AcceptanceCriteria count = %d, want %d",
			len(unmarshaled.AcceptanceCriteria), len(rd.AcceptanceCriteria))
	}
}
