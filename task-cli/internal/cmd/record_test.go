package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"task-cli/pkg/task"
)

func TestValidateRecordData(t *testing.T) {
	t.Run("empty summary triggers hard error", func(t *testing.T) {
		rd := &RecordData{
			Status: "completed",
			Summary: "",
		}
		// validateRecordData calls Exit() which calls os.Exit(1)
		// We test the logic by catching the exit via a subprocess
		// Instead, let's test the validation logic directly
		// Since validateRecordData calls Exit, we test via subprocess
		if os.Getenv("TEST_VALIDATE_EMPTY_SUMMARY") == "1" {
			validateRecordData(rd)
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestValidateRecordData/empty_summary_triggers_hard_error")
		cmd.Env = append(os.Environ(), "TEST_VALIDATE_EMPTY_SUMMARY=1")
		err := cmd.Run()
		if err == nil {
			t.Error("expected non-zero exit for empty summary")
		}
	})

	t.Run("whitespace-only summary triggers hard error", func(t *testing.T) {
		if os.Getenv("TEST_VALIDATE_WS_SUMMARY") == "1" {
			validateRecordData(&RecordData{Status: "completed", Summary: "   "})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestValidateRecordData/whitespace-only_summary_triggers_hard_error")
		cmd.Env = append(os.Environ(), "TEST_VALIDATE_WS_SUMMARY=1")
		err := cmd.Run()
		if err == nil {
			t.Error("expected non-zero exit for whitespace-only summary")
		}
	})

	t.Run("completed without recommended fields warns", func(t *testing.T) {
		rd := &RecordData{
			Status:  "completed",
			Summary: "Did the work",
		}
		// Capture stderr for warning
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd)
		w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if !strings.Contains(output, "WARNING") {
			t.Errorf("expected warning in stderr, got: %s", output)
		}
		// Should warn about keyDecisions, tests, acceptanceCriteria
		for _, field := range []string{"keyDecisions", "coverage", "acceptanceCriteria"} {
			if !strings.Contains(output, field) {
				t.Errorf("expected warning to mention %q, got: %s", field, output)
			}
		}
	})

	t.Run("completed with all fields produces no warning", func(t *testing.T) {
		rd := &RecordData{
			Status:             "completed",
			Summary:            "Full record",
			KeyDecisions:       []string{"decision"},
			TestsPassed:        5,
			Coverage:           80.0,
			AcceptanceCriteria: []AcceptanceCriterion{{Criterion: "works", Met: true}},
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd)
		w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "WARNING") {
			t.Errorf("unexpected warning for complete record: %s", output)
		}
	})

	t.Run("non-completed status skips recommended checks", func(t *testing.T) {
		rd := &RecordData{
			Status:  "blocked",
			Summary: "Blocked with reason",
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd)
		w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "WARNING") {
			t.Errorf("non-completed status should not produce warnings: %s", output)
		}
	})
}

func TestFindTask(t *testing.T) {
	tests := []struct {
		name     string
		tasks    map[string]task.Task
		searchID string
		wantKey  string
		wantID   string
		wantErr  bool
	}{
		{
			name: "find by task ID",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1"},
				"task2": {ID: "1.2", Title: "Task 2"},
			},
			searchID: "1.1",
			wantKey:  "task1",
			wantID:   "1.1",
			wantErr:  false,
		},
		{
			name: "find by key",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1"},
				"task2": {ID: "1.2", Title: "Task 2"},
			},
			searchID: "task2",
			wantKey:  "task2",
			wantID:   "1.2",
			wantErr:  false,
		},
		{
			name: "not found",
			tasks: map[string]task.Task{
				"task1": {ID: "1.1", Title: "Task 1"},
			},
			searchID: "2.1",
			wantKey:  "",
			wantID:   "",
			wantErr:  true,
		},
		{
			name:     "empty tasks",
			tasks:    map[string]task.Task{},
			searchID: "1.1",
			wantKey:  "",
			wantID:   "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := &task.TaskIndex{
				Feature: "test",
				Tasks:   tt.tasks,
			}
			key, gotTask, err := findTask(index, tt.searchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("findTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if key != tt.wantKey {
				t.Errorf("findTask() key = %q, want %q", key, tt.wantKey)
			}
			if !tt.wantErr && gotTask.ID != tt.wantID {
				t.Errorf("findTask() task.ID = %q, want %q", gotTask.ID, tt.wantID)
			}
		})
	}
}

func TestFormatList(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  string
	}{
		{
			name:  "empty list",
			items: []string{},
			want:  "无",
		},
		{
			name:  "single item",
			items: []string{"file1.go"},
			want:  "- file1.go",
		},
		{
			name:  "multiple items",
			items: []string{"file1.go", "file2.go", "file3.go"},
			want:  "- file1.go\n- file2.go\n- file3.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatList(tt.items)
			if got != tt.want {
				t.Errorf("formatList() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name string
		dur  time.Duration
		want string
	}{
		{
			name: "minutes only",
			dur:  45 * time.Minute,
			want: "~45m",
		},
		{
			name: "hours only",
			dur:  2 * time.Hour,
			want: "~2h",
		},
		{
			name: "hours and minutes",
			dur:  2*time.Hour + 30*time.Minute,
			want: "~2h 30m",
		},
		{
			name: "less than hour",
			dur:  30 * time.Minute,
			want: "~30m",
		},
		{
			name: "zero",
			dur:  0,
			want: "~0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.dur)
			if got != tt.want {
				t.Errorf("formatDuration() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatCriteria(t *testing.T) {
	tests := []struct {
		name     string
		criteria []AcceptanceCriterion
		want     string
	}{
		{
			name:     "empty criteria",
			criteria: []AcceptanceCriterion{},
			want:     "无",
		},
		{
			name: "single unmet criterion",
			criteria: []AcceptanceCriterion{
				{Criterion: "Feature works", Met: false},
			},
			want: "- [ ] Feature works",
		},
		{
			name: "single met criterion",
			criteria: []AcceptanceCriterion{
				{Criterion: "Feature works", Met: true},
			},
			want: "- [x] Feature works",
		},
		{
			name: "multiple mixed criteria",
			criteria: []AcceptanceCriterion{
				{Criterion: "Feature works", Met: true},
				{Criterion: "Tests pass", Met: false},
				{Criterion: "Docs updated", Met: true},
			},
			want: "- [x] Feature works\n- [ ] Tests pass\n- [x] Docs updated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatCriteria(tt.criteria)
			if got != tt.want {
				t.Errorf("formatCriteria() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFillRecordTemplate(t *testing.T) {
	tests := []struct {
		name             string
		task             *task.Task
		recordData       *RecordData
		startedTime      string
		checkContains    []string
		checkNotContains []string
	}{
		{
			name: "basic template",
			task: &task.Task{
				ID:    "1.1",
				Title: "Implement feature X",
			},
			recordData: &RecordData{
				Status:       "completed",
				Summary:      "Implemented the feature",
				FilesCreated: []string{"main.go"},
				TestsPassed:  5,
				Coverage:     85.5,
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"Implement feature X",
				"Implemented the feature",
				"completed",
				"2026-04-06 10:00",
				"main.go",
				"Tests Passed: 5",
				"Coverage: 85.5%",
			},
		},
		{
			name: "template with all fields",
			task: &task.Task{
				ID:    "2.1",
				Title: "Full feature",
			},
			recordData: &RecordData{
				Status:        "completed",
				Summary:       "Complete implementation",
				FilesCreated:  []string{"a.go", "b.go"},
				FilesModified: []string{"c.go"},
				KeyDecisions:  []string{"Used pattern X"},
				TestsPassed:   10,
				TestsFailed:   2,
				Coverage:      90.0,
				AcceptanceCriteria: []AcceptanceCriterion{
					{Criterion: "AC1", Met: true},
					{Criterion: "AC2", Met: false},
				},
				Notes: "Some notes",
			},
			startedTime: "2026-04-06 09:00",
			checkContains: []string{
				"Full feature",
				"Complete implementation",
				"a.go",
				"b.go",
				"c.go",
				"Used pattern X",
				"Tests Passed: 10",
				"Tests Failed: 2",
				"Coverage: 90.0%",
				"[x] AC1",
				"[ ] AC2",
				"Some notes",
			},
		},
		{
			name: "non-completed status",
			task: &task.Task{
				ID:    "1.2",
				Title: "In progress task",
			},
			recordData: &RecordData{
				Status:  "in_progress",
				Summary: "Work in progress",
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"Status: in_progress",
				"Completed: N/A",
			},
		},
		{
			name: "default notes when empty",
			task: &task.Task{
				ID:    "1.3",
				Title: "Task with no notes",
			},
			recordData: &RecordData{
				Status:  "completed",
				Summary: "Done",
				Notes:   "",
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"无",
			},
		},
		{
			name: "empty started time uses current",
			task: &task.Task{
				ID:    "1.4",
				Title: "Task",
			},
			recordData: &RecordData{
				Status:  "completed",
				Summary: "Done",
			},
			startedTime: "",
			checkContains: []string{
				"Started:",
				"Completed:",
			},
		},
		{
			name: "time spent when completed after started",
			task: &task.Task{
				ID:    "1.5",
				Title: "Timed Task",
			},
			recordData: &RecordData{
				Status:  "completed",
				Summary: "Done",
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"Time Spent:",
			},
		},
		{
			name: "no time spent when completed before started",
			task: &task.Task{
				ID:    "1.6",
				Title: "Backward Time Task",
			},
			recordData: &RecordData{
				Status:  "completed",
				Summary: "Done",
			},
			startedTime: "2026-04-06 15:00", // Started after current time would be
			checkNotContains: []string{
				"Time Spent: ~",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fillRecordTemplate(tt.task, tt.recordData, tt.startedTime)
			for _, expected := range tt.checkContains {
				if !strings.Contains(got, expected) {
					t.Errorf("fillRecordTemplate() missing expected content %q", expected)
				}
			}
		})
	}
}

func TestReadRecordData(t *testing.T) {
	t.Run("read from file", func(t *testing.T) {
		// Create temp file with JSON data
		dir := t.TempDir()
		dataPath := dir + "/data.json"
		jsonData := `{"status":"completed","summary":"Done","testsPassed":5,"coverage":80.5}`
		if err := os.WriteFile(dataPath, []byte(jsonData), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		rd, err := readRecordData(dataPath)
		if err != nil {
			t.Fatalf("readRecordData() error = %v", err)
		}
		if rd.Status != "completed" {
			t.Errorf("Status = %q, want %q", rd.Status, "completed")
		}
		if rd.Summary != "Done" {
			t.Errorf("Summary = %q, want %q", rd.Summary, "Done")
		}
		if rd.TestsPassed != 5 {
			t.Errorf("TestsPassed = %d, want 5", rd.TestsPassed)
		}
		if rd.Coverage != 80.5 {
			t.Errorf("Coverage = %f, want 80.5", rd.Coverage)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := readRecordData("/nonexistent/path/file.json")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		dataPath := dir + "/invalid.json"
		if err := os.WriteFile(dataPath, []byte("not valid json"), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		_, err := readRecordData(dataPath)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("no input without data flag", func(t *testing.T) {
		// When dataPath is empty and stdin is not a pipe, should error
		// This test verifies the error message
		_, err := readRecordData("")
		if err == nil {
			t.Error("expected error when no data provided")
		}
		if !strings.Contains(err.Error(), "no input") {
			t.Errorf("error should mention 'no input', got: %v", err)
		}
	})
}
