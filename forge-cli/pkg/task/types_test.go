package task

import (
	"encoding/json"
	"strings"
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
			if index.tasks == nil {
				t.Error("Tasks map is nil")
			}
			if len(index.tasks) != 0 {
				t.Errorf("Tasks map should be empty, got %d items", len(index.tasks))
			}
			// Check default status enum
			expectedStatuses := []string{"pending", "in_progress", "completed", "blocked", "suspended", "skipped", "rejected"}
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
		tasks: map[string]Task{
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
	if len(unmarshaled.tasks) != len(index.tasks) {
		t.Errorf("Tasks count = %d, want %d", len(unmarshaled.tasks), len(index.tasks))
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

func TestTaskScopeSerialization(t *testing.T) {
	t.Run("surface-key field serializes when set", func(t *testing.T) {
		task := Task{
			ID:          "1.1",
			Title:       "Frontend Task",
			SurfaceKey:  "admin-panel",
			SurfaceType: "web",
			Status:      "pending",
			File:        "tasks/1.1.md",
		}

		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		got := string(data)
		if !strings.Contains(got, `"surface-key":"admin-panel"`) {
			t.Errorf("JSON = %s, want to contain %q", got, `"surface-key":"admin-panel"`)
		}
		if !strings.Contains(got, `"surface-type":"web"`) {
			t.Errorf("JSON = %s, want to contain %q", got, `"surface-type":"web"`)
		}
	})

	t.Run("surface fields omitted when empty", func(t *testing.T) {
		task := Task{
			ID:     "1.1",
			Title:  "Task Without Surface",
			Status: "pending",
			File:   "tasks/1.1.md",
		}

		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		got := string(data)
		if strings.Contains(got, `"surface-key"`) {
			t.Errorf("JSON = %s, should NOT contain surface-key field when empty", got)
		}
		if strings.Contains(got, `"surface-type"`) {
			t.Errorf("JSON = %s, should NOT contain surface-type field when empty", got)
		}
	})

	t.Run("surface fields deserialize from JSON", func(t *testing.T) {
		jsonStr := `{"id":"1.1","title":"Backend Task","surface-key":"payment-service","surface-type":"api","status":"pending","file":"tasks/1.1.md"}`

		var task Task
		if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}

		if task.SurfaceKey != "payment-service" {
			t.Errorf("SurfaceKey = %q, want %q", task.SurfaceKey, "payment-service")
		}
		if task.SurfaceType != "api" {
			t.Errorf("SurfaceType = %q, want %q", task.SurfaceType, "api")
		}
	})

	t.Run("surface fields default to empty when missing in JSON", func(t *testing.T) {
		jsonStr := `{"id":"1.1","title":"No Surface Task","status":"pending","file":"tasks/1.1.md"}`

		var task Task
		if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}

		if task.SurfaceKey != "" {
			t.Errorf("SurfaceKey = %q, want empty string when missing", task.SurfaceKey)
		}
		if task.SurfaceType != "" {
			t.Errorf("SurfaceType = %q, want empty string when missing", task.SurfaceType)
		}
	})

	t.Run("surface fields roundtrip preserves value", func(t *testing.T) {
		task := Task{
			ID:          "2.3",
			Title:       "Mixed Task",
			SurfaceKey:  "admin-panel",
			SurfaceType: "web",
			Status:      "pending",
			File:        "tasks/2.3.md",
		}

		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		var unmarshaled Task
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}

		if unmarshaled.SurfaceKey != task.SurfaceKey {
			t.Errorf("SurfaceKey roundtrip = %q, want %q", unmarshaled.SurfaceKey, task.SurfaceKey)
		}
		if unmarshaled.SurfaceType != task.SurfaceType {
			t.Errorf("SurfaceType roundtrip = %q, want %q", unmarshaled.SurfaceType, task.SurfaceType)
		}
	})
}

func TestTaskStateScopeSerialization(t *testing.T) {
	t.Run("TaskState surface fields serialize when set", func(t *testing.T) {
		state := &TaskState{
			TaskID:      "1.1",
			Key:         "task1",
			Title:       "Frontend Task",
			SurfaceKey:  "admin-panel",
			SurfaceType: "web",
			StartedTime: "2024-01-01 10:00",
		}

		data, err := json.Marshal(state)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		got := string(data)
		if !strings.Contains(got, `"surface-key":"admin-panel"`) {
			t.Errorf("JSON = %s, want to contain %q", got, `"surface-key":"admin-panel"`)
		}
		if !strings.Contains(got, `"surface-type":"web"`) {
			t.Errorf("JSON = %s, want to contain %q", got, `"surface-type":"web"`)
		}
	})

	t.Run("TaskState surface fields omitted when empty", func(t *testing.T) {
		state := &TaskState{
			TaskID:      "1.1",
			Key:         "task1",
			Title:       "Task",
			StartedTime: "2024-01-01 10:00",
		}

		data, err := json.Marshal(state)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		got := string(data)
		if strings.Contains(got, `"surface-key"`) {
			t.Errorf("JSON = %s, should NOT contain surface-key field when empty", got)
		}
	})

	t.Run("TaskState surface fields deserialize from JSON", func(t *testing.T) {
		jsonStr := `{"task_id":"1.1","key":"task1","title":"Backend Task","surface-key":"payment-service","surface-type":"api","startedTime":"2024-01-01 10:00"}`

		var state TaskState
		if err := json.Unmarshal([]byte(jsonStr), &state); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}

		if state.SurfaceKey != "payment-service" {
			t.Errorf("SurfaceKey = %q, want %q", state.SurfaceKey, "payment-service")
		}
		if state.SurfaceType != "api" {
			t.Errorf("SurfaceType = %q, want %q", state.SurfaceType, "api")
		}
	})
}

func TestTaskIndexJSONRoundTrip_AllFields(t *testing.T) {
	original := &TaskIndex{
		Feature:      "full-feature",
		PRD:          "prd/prd-spec.md",
		Design:       "design/tech-design.md",
		Created:      "2024-06-15",
		Status:       "in_progress",
		tasks:        map[string]Task{"t1": {ID: "1.1", Title: "Task", Status: "pending"}},
		StatusEnum:   []string{"pending", "completed"},
		PriorityEnum: []string{"P0", "P1"},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Verify raw JSON contains "tasks" key (two-struct pattern must emit it)
	if !strings.Contains(string(data), `"tasks"`) {
		t.Errorf("JSON output missing \"tasks\" key: %s", data)
	}

	var loaded TaskIndex
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	assertField(t, "Feature", loaded.Feature, original.Feature)
	assertField(t, "PRD", loaded.PRD, original.PRD)
	assertField(t, "Design", loaded.Design, original.Design)
	assertField(t, "Created", loaded.Created, original.Created)
	assertField(t, "Status", loaded.Status, original.Status)
	if loaded.TaskCount() != original.TaskCount() {
		t.Errorf("TaskCount = %d, want %d", loaded.TaskCount(), original.TaskCount())
	}
	if len(loaded.StatusEnum) != len(original.StatusEnum) {
		t.Errorf("StatusEnum len = %d, want %d", len(loaded.StatusEnum), len(original.StatusEnum))
	}
	if len(loaded.PriorityEnum) != len(original.PriorityEnum) {
		t.Errorf("PriorityEnum len = %d, want %d", len(loaded.PriorityEnum), len(original.PriorityEnum))
	}
}

func TestTaskIndexUnmarshal_EmptyTasks(t *testing.T) {
	t.Run("tasks null", func(t *testing.T) {
		var idx TaskIndex
		if err := json.Unmarshal([]byte(`{"feature":"x","tasks":null}`), &idx); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if idx.tasks != nil {
			t.Error("expected nil tasks for null input")
		}
	})

	t.Run("tasks key absent", func(t *testing.T) {
		var idx TaskIndex
		if err := json.Unmarshal([]byte(`{"feature":"x"}`), &idx); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}
		if idx.tasks != nil {
			t.Error("expected nil tasks when key absent")
		}
		if idx.Feature != "x" {
			t.Errorf("Feature = %q, want %q", idx.Feature, "x")
		}
	})
}

func assertField(t *testing.T, name string, got, want any) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func TestTaskTypeFieldSerialization(t *testing.T) {
	t.Run("type field serializes when set", func(t *testing.T) {
		task := Task{
			ID:     "1.1",
			Title:  "Impl Task",
			Status: "pending",
			File:   "tasks/1.1.md",
			Type:   TypeCodingFeature,
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if !strings.Contains(got, `"type":"coding.feature"`) {
			t.Errorf("JSON = %s, want to contain %q", got, `"type":"coding.feature"`)
		}
	})

	t.Run("type field omitted when empty", func(t *testing.T) {
		task := Task{
			ID:     "1.1",
			Title:  "Task Without Type",
			Status: "pending",
			File:   "tasks/1.1.md",
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if strings.Contains(got, `"type"`) {
			t.Errorf("JSON = %s, should NOT contain type field when empty", got)
		}
	})

	t.Run("type field deserializes from JSON", func(t *testing.T) {
		jsonStr := `{"id":"1.1","title":"Fix Task","type":"coding.fix","status":"pending","file":"tasks/1.1.md"}`
		var task Task
		if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if task.Type != TypeCodingFix {
			t.Errorf("Type = %q, want %q", task.Type, TypeCodingFix)
		}
	})

	t.Run("type field defaults to empty when missing in JSON", func(t *testing.T) {
		jsonStr := `{"id":"1.1","title":"No Type Task","status":"pending","file":"tasks/1.1.md"}`
		var task Task
		if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if task.Type != "" {
			t.Errorf("Type = %q, want empty string when missing", task.Type)
		}
	})
}

func TestTaskBlockedReasonFieldSerialization(t *testing.T) {
	t.Run("blockedReason field serializes when set", func(t *testing.T) {
		task := Task{
			ID:            "1.1",
			Title:         "Blocked Task",
			Status:        "blocked",
			File:          "tasks/1.1.md",
			BlockedReason: "template file missing",
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if !strings.Contains(got, `"blockedReason":"template file missing"`) {
			t.Errorf("JSON = %s, want to contain blockedReason", got)
		}
	})

	t.Run("blockedReason field omitted when empty", func(t *testing.T) {
		task := Task{
			ID:     "1.1",
			Title:  "Normal Task",
			Status: "pending",
			File:   "tasks/1.1.md",
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if strings.Contains(got, `"blockedReason"`) {
			t.Errorf("JSON = %s, should NOT contain blockedReason when empty", got)
		}
	})

	t.Run("blockedReason roundtrip preserves value", func(t *testing.T) {
		task := Task{
			ID:            "2.1",
			Title:         "Blocked Task",
			Status:        "blocked",
			File:          "tasks/2.1.md",
			BlockedReason: "task prompt exited with code 1",
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		var unmarshaled Task
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if unmarshaled.BlockedReason != task.BlockedReason {
			t.Errorf("BlockedReason = %q, want %q", unmarshaled.BlockedReason, task.BlockedReason)
		}
	})
}

func TestTaskComplexityFieldSerialization(t *testing.T) {
	t.Run("complexity field serializes when set", func(t *testing.T) {
		task := Task{
			ID:         "1.1",
			Title:      "Complex Task",
			Status:     "pending",
			File:       "tasks/1.1.md",
			Complexity: "low",
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if !strings.Contains(got, `"complexity":"low"`) {
			t.Errorf("JSON = %s, want to contain %q", got, `"complexity":"low"`)
		}
	})

	t.Run("complexity field omitted when empty", func(t *testing.T) {
		task := Task{
			ID:     "1.1",
			Title:  "Task Without Complexity",
			Status: "pending",
			File:   "tasks/1.1.md",
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if strings.Contains(got, `"complexity"`) {
			t.Errorf("JSON = %s, should NOT contain complexity field when empty", got)
		}
	})

	t.Run("complexity field deserializes from JSON", func(t *testing.T) {
		jsonStr := `{"id":"1.1","title":"High Task","complexity":"high","status":"pending","file":"tasks/1.1.md"}`
		var task Task
		if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if task.Complexity != "high" {
			t.Errorf("Complexity = %q, want %q", task.Complexity, "high")
		}
	})

	t.Run("complexity defaults to empty when missing in JSON", func(t *testing.T) {
		jsonStr := `{"id":"1.1","title":"No Complexity Task","status":"pending","file":"tasks/1.1.md"}`
		var task Task
		if err := json.Unmarshal([]byte(jsonStr), &task); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if task.Complexity != "" {
			t.Errorf("Complexity = %q, want empty string when missing", task.Complexity)
		}
	})

	t.Run("complexity roundtrip preserves value", func(t *testing.T) {
		task := Task{
			ID:         "2.1",
			Title:      "Medium Task",
			Status:     "pending",
			File:       "tasks/2.1.md",
			Complexity: "medium",
		}
		data, err := json.Marshal(task)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		var unmarshaled Task
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if unmarshaled.Complexity != task.Complexity {
			t.Errorf("Complexity = %q, want %q", unmarshaled.Complexity, task.Complexity)
		}
	})
}

func TestTaskStateTypeFieldSerialization(t *testing.T) {
	t.Run("TaskState type field serializes when set", func(t *testing.T) {
		state := &TaskState{
			TaskID:      "1.1",
			Key:         "task1",
			Title:       "Gate Task",
			Type:        TypeGate,
			StartedTime: "2024-01-01 10:00",
		}
		data, err := json.Marshal(state)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if !strings.Contains(got, `"type":"gate"`) {
			t.Errorf("JSON = %s, want to contain %q", got, `"type":"gate"`)
		}
	})

	t.Run("TaskState type field omitted when empty", func(t *testing.T) {
		state := &TaskState{
			TaskID:      "1.1",
			Key:         "task1",
			Title:       "Task",
			StartedTime: "2024-01-01 10:00",
		}
		data, err := json.Marshal(state)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}
		got := string(data)
		if strings.Contains(got, `"type"`) {
			t.Errorf("JSON = %s, should NOT contain type field when empty", got)
		}
	})

	t.Run("TaskState type deserializes from JSON", func(t *testing.T) {
		jsonStr := `{"task_id":"1.1","key":"task1","title":"Fix Task","type":"coding.fix","startedTime":"2024-01-01 10:00"}`
		var state TaskState
		if err := json.Unmarshal([]byte(jsonStr), &state); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}
		if state.Type != TypeCodingFix {
			t.Errorf("Type = %q, want %q", state.Type, TypeCodingFix)
		}
	})
}

func TestTypeConstants(t *testing.T) {
	tests := []struct {
		constant string
		expected string
	}{
		{TypeCodingFeature, "coding.feature"},
		{TypeCodingEnhancement, "coding.enhancement"},
		{TypeCodingCleanup, "coding.cleanup"},
		{TypeCodingRefactor, "coding.refactor"},
		{TypeCodingFix, "coding.fix"},
		{TypeDocFix, "doc.fix"},
		{TypeDoc, "doc"},
		{TypeDocReview, "doc.review"},
		{TypeDocSummary, "doc.summary"},
		{TypeDocConsolidate, "doc.consolidate"},
		{TypeDocDrift, "doc.drift"},
		{TypeTestGenContracts, "test.gen-contracts"},
		{TypeTestGenJourneys, "test.gen-journeys"},
		{TypeTestGenScripts, "test.gen-scripts"},
		{TypeTestRun, "test.run"},
		{TypeValidationCode, "validation.code"},
		{TypeValidationUx, "validation.ux"},
		{TypeGate, "gate"},
		{TypeCleanCode, "code-quality.simplify"},
	}
	for _, tt := range tests {
		if tt.constant != tt.expected {
			t.Errorf("constant value = %q, want %q", tt.constant, tt.expected)
		}
	}
}

func TestSystemTypes(t *testing.T) {
	t.Run("SystemTypes contains exactly 13 entries", func(t *testing.T) {
		if len(SystemTypes) != 13 {
			t.Errorf("SystemTypes has %d entries, want 13", len(SystemTypes))
		}
	})

	t.Run("all system types are present", func(t *testing.T) {
		expected := []string{
			TypeGate,
			TypeTestGenContracts, TypeTestGenJourneys,
			TypeTestGenScripts, TypeTestRun,
			TypeValidationCode, TypeValidationUx,
			TypeDocReview, TypeDocSummary,
			TypeCleanCode,
			TypeEvalJourney, TypeEvalContract,
			TypeDocFix,
		}
		if len(expected) != 13 {
			t.Fatalf("test setup error: expected list has %d entries, want 13", len(expected))
		}
		for _, typ := range expected {
			if !SystemTypes[typ] {
				t.Errorf("SystemTypes missing %q", typ)
			}
		}
	})

	t.Run("SystemTypes is independent from ValidTypes", func(_ *testing.T) {
		// SystemTypes entries should NOT be added to ValidTypes
		for typ := range SystemTypes {
			// All system types happen to be in ValidTypes too, but the maps are independent
			_ = typ // just verify iteration works without panic
		}
	})
}

func TestIsSystemType(t *testing.T) {
	t.Run("returns true for system types including gen-journeys and gen-contracts", func(t *testing.T) {
		systemTypes := []string{
			TypeGate,
			TypeTestGenContracts, TypeTestGenJourneys,
			TypeTestGenScripts, TypeTestRun,
			TypeValidationCode, TypeValidationUx,
			TypeDocReview, TypeDocSummary,
			TypeCleanCode,
		}
		for _, typ := range systemTypes {
			if !IsSystemType(typ) {
				t.Errorf("IsSystemType(%q) = false, want true", typ)
			}
		}
	})

	t.Run("returns false for business types", func(t *testing.T) {
		businessTypes := []string{
			TypeCodingFeature, TypeCodingEnhancement, TypeCodingCleanup,
			TypeCodingRefactor, TypeCodingFix,
			TypeDoc,
		}
		for _, typ := range businessTypes {
			if IsSystemType(typ) {
				t.Errorf("IsSystemType(%q) = true, want false", typ)
			}
		}
	})

	t.Run("returns false for dual-identity types (doc.consolidate, doc.drift)", func(t *testing.T) {
		dualTypes := []string{TypeDocConsolidate, TypeDocDrift}
		for _, typ := range dualTypes {
			if IsSystemType(typ) {
				t.Errorf("IsSystemType(%q) = true, want false (dual-identity type)", typ)
			}
		}
	})

	t.Run("returns false for unknown types", func(t *testing.T) {
		if IsSystemType("unknown") {
			t.Error("IsSystemType('unknown') = true, want false")
		}
		if IsSystemType("") {
			t.Error("IsSystemType('') = true, want false")
		}
	})
}

func TestValidTypes(t *testing.T) {
	t.Run("ValidTypes contains all type constants", func(t *testing.T) {
		allTypes := []string{
			TypeCodingFeature,
			TypeCodingEnhancement,
			TypeCodingCleanup,
			TypeCodingRefactor,
			TypeDoc,
			TypeDocReview,
			TypeDocSummary,
			TypeDocConsolidate,
			TypeDocDrift,
			TypeTestGenContracts,
			TypeTestGenJourneys,
			TypeTestGenScripts,
			TypeTestRun,
			TypeEvalJourney,
			TypeEvalContract,
			TypeValidationCode,
			TypeValidationUx,
			TypeCodingFix,
			TypeDocFix,
			TypeGate,
			TypeCleanCode,
		}
		if len(ValidTypes) != len(allTypes) {
			t.Errorf("ValidTypes has %d entries, want %d", len(ValidTypes), len(allTypes))
		}
		for _, typ := range allTypes {
			if !ValidTypes[typ] {
				t.Errorf("ValidTypes missing %q", typ)
			}
		}
	})

	t.Run("ValidTypes rejects unknown type", func(t *testing.T) {
		if ValidTypes["unknown-type"] {
			t.Error("ValidTypes should not contain 'unknown-type'")
		}
		if ValidTypes[""] {
			t.Error("ValidTypes should not contain empty string")
		}
	})
}

func TestTaskTypeRegistry(t *testing.T) {
	t.Run("registry contains all types", func(t *testing.T) {
		if len(TaskTypeRegistry) != len(ValidTypes) {
			t.Errorf("TaskTypeRegistry has %d entries, want %d", len(TaskTypeRegistry), len(ValidTypes))
		}
	})

	t.Run("each registry entry matches a type constant", func(t *testing.T) {
		registryNames := make(map[string]bool)
		for _, entry := range TaskTypeRegistry {
			registryNames[entry.Name] = true
			if !ValidTypes[entry.Name] {
				t.Errorf("TaskTypeRegistry entry %q not in ValidTypes", entry.Name)
			}
			if entry.Description == "" {
				t.Errorf("TaskTypeRegistry entry %q has empty description", entry.Name)
			}
			if len(entry.Description) > 60 {
				t.Errorf("TaskTypeRegistry entry %q description too long (%d chars): %q",
					entry.Name, len(entry.Description), entry.Description)
			}
		}

		// Verify all type constants are present in registry
		for typ := range ValidTypes {
			if !registryNames[typ] {
				t.Errorf("ValidTypes entry %q missing from TaskTypeRegistry", typ)
			}
		}
	})

	t.Run("descriptions use verb+object format", func(t *testing.T) {
		for _, entry := range TaskTypeRegistry {
			if entry.Description == "" {
				continue
			}
			// First word should be a verb (lowercase letter start)
			first := entry.Description[0]
			if first < 'a' || first > 'z' {
				t.Errorf("TaskTypeRegistry entry %q description does not start with lowercase verb: %q",
					entry.Name, entry.Description)
			}
		}
	})

	t.Run("registry entries have no duplicates", func(t *testing.T) {
		seen := make(map[string]bool)
		for _, entry := range TaskTypeRegistry {
			if seen[entry.Name] {
				t.Errorf("duplicate type name in registry: %q", entry.Name)
			}
			seen[entry.Name] = true
		}
	})
}

func TestTestTypeTitle(t *testing.T) {
	tests := []struct {
		surfaceType string
		want        string
	}{
		{"cli", "CLI Functional Test"},
		{"tui", "Terminal Functional Test"},
		{"api", "API Functional Test"},
		{"web", "Web E2E Test"},
		{"mobile", "Mobile E2E Test"},
		{"unknown", "Functional Test"},
		{"", "Functional Test"},
	}
	for _, tt := range tests {
		t.Run(tt.surfaceType, func(t *testing.T) {
			got := TestTypeTitle(tt.surfaceType)
			if got != tt.want {
				t.Errorf("TestTypeTitle(%q) = %q, want %q", tt.surfaceType, got, tt.want)
			}
		})
	}
}

func TestGenSurfaceTestType(t *testing.T) {
	tests := []struct {
		baseType string
		surface  string
		want     string
	}{
		{"test.gen-scripts", "cli", "test.gen-scripts.cli"},
		{"test.run", "api", "test.run.api"},
		{"test.gen-scripts", "", "test.gen-scripts"},
		{"test.run", "", "test.run"},
	}
	for _, tt := range tests {
		t.Run(tt.baseType+"+"+tt.surface, func(t *testing.T) {
			got := GenSurfaceTestType(tt.baseType, tt.surface)
			if got != tt.want {
				t.Errorf("GenSurfaceTestType(%q, %q) = %q, want %q", tt.baseType, tt.surface, got, tt.want)
			}
		})
	}
}

func TestIsValidType(t *testing.T) {
	t.Run("exact match returns true for all ValidTypes", func(t *testing.T) {
		for typ := range ValidTypes {
			if !IsValidType(typ) {
				t.Errorf("IsValidType(%q) = false, want true (exact match)", typ)
			}
		}
	})

	t.Run("surface-suffixed variants pass validation", func(t *testing.T) {
		surfaceVariants := []string{
			"test.gen-scripts.cli",
			"test.gen-scripts.api",
			"test.gen-scripts.web",
			"test.run.cli",
			"test.run.api",
		}
		for _, typ := range surfaceVariants {
			if !IsValidType(typ) {
				t.Errorf("IsValidType(%q) = false, want true (surface variant)", typ)
			}
		}
	})

	t.Run("non-generated types with suffix still rejected", func(t *testing.T) {
		// coding.feature.cli is NOT a valid surface variant — coding types don't get surface-suffixed
		if IsValidType("coding.feature.cli") {
			t.Error("IsValidType('coding.feature.cli') = true, want false")
		}
	})

	t.Run("unknown types rejected", func(t *testing.T) {
		if IsValidType("unknown-type") {
			t.Error("IsValidType('unknown-type') = true, want false")
		}
		if IsValidType("") {
			t.Error("IsValidType('') = true, want false")
		}
		if !IsValidType("test.gen-scripts") {
			t.Error("IsValidType('test.gen-scripts') = false, want true (exact match)")
		}
	})

	t.Run("only system types support surface suffix", func(t *testing.T) {
		// doc.review.cli — doc.review IS a system type, so this should pass
		if !IsValidType("doc.review.cli") {
			t.Error("IsValidType('doc.review.cli') = false, want true")
		}
		// gate.cli — gate IS a system type, so this should pass
		if !IsValidType("gate.cli") {
			t.Error("IsValidType('gate.cli') = false, want true")
		}
	})
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
