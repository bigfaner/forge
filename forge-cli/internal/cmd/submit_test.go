package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/task"
)

func TestValidateRecordData(t *testing.T) {
	t.Run("empty summary triggers hard error", func(t *testing.T) {
		rd := &task.RecordData{
			Status:  "completed",
			Summary: "",
		}
		if os.Getenv("TEST_VALIDATE_EMPTY_SUMMARY") == "1" {
			validateRecordData(rd, false)
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
			validateRecordData(&task.RecordData{Status: "completed", Summary: "   "}, false)
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestValidateRecordData/whitespace-only_summary_triggers_hard_error")
		cmd.Env = append(os.Environ(), "TEST_VALIDATE_WS_SUMMARY=1")
		err := cmd.Run()
		if err == nil {
			t.Error("expected non-zero exit for whitespace-only summary")
		}
	})

	t.Run("completed with testsFailed auto-downgrades to blocked", func(t *testing.T) {
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Partial pass",
			TestsPassed: 3,
			TestsFailed: 2,
			Coverage:    60.0,
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 2048)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if rd.Status != "blocked" {
			t.Errorf("expected status downgraded to 'blocked', got %q", rd.Status)
		}
		if !strings.Contains(output, "auto-downgrading") {
			t.Errorf("expected auto-downgrade warning in stderr, got: %s", output)
		}
	})

	t.Run("force does NOT prevent auto-downgrade of testsFailed", func(t *testing.T) {
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Partial pass",
			TestsPassed: 3,
			TestsFailed: 2,
			Coverage:    60.0,
		}
		validateRecordData(rd, true)

		if rd.Status != "blocked" {
			t.Errorf("expected status downgraded even with force=true, got %q", rd.Status)
		}
	})

	t.Run("completed without test evidence triggers hard error", func(t *testing.T) {
		if os.Getenv("TEST_VALIDATE_NO_TESTS") == "1" {
			validateRecordData(&task.RecordData{Status: "completed", Summary: "Did the work", TestsPassed: 0, TestsFailed: 0, Coverage: 0}, false)
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestValidateRecordData/completed_without_test_evidence_triggers_hard_error")
		cmd.Env = append(os.Environ(), "TEST_VALIDATE_NO_TESTS=1")
		err := cmd.Run()
		if err == nil {
			t.Error("expected non-zero exit for completed with no test evidence")
		}
	})

	t.Run("completed with coverage=-1 skips test evidence check", func(t *testing.T) {
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Doc task",
			Coverage:    -1.0,
			TestsPassed: 0,
			TestsFailed: 0,
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "ERROR") {
			t.Errorf("coverage=-1.0 should skip test evidence check, got: %s", output)
		}
	})

	t.Run("noTest task with coverage=-1 skips validation", func(t *testing.T) {
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Documentation-only task",
			Coverage:    -1.0,
			TestsPassed: 0,
			TestsFailed: 0,
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "ERROR") {
			t.Errorf("noTest task should skip test evidence check, got: %s", output)
		}
	})

	t.Run("noTest task with testsPassed > 0 passes validation", func(t *testing.T) {
		rd := &task.RecordData{
			Status:       "completed",
			Summary:      "Ran some tests despite noTest flag",
			Coverage:     80.0,
			TestsPassed:  5,
			TestsFailed:  0,
			KeyDecisions: []string{"tested anyway"},
			AcceptanceCriteria: []task.AcceptanceCriterion{
				{Criterion: "Works", Met: true},
			},
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "ERROR") {
			t.Errorf("noTest + testsPassed > 0 should pass, got: %s", output)
		}
	})

	t.Run("completed with tests passes test evidence check", func(t *testing.T) {
		rd := &task.RecordData{
			Status:             "completed",
			Summary:            "Full record",
			KeyDecisions:       []string{"decision"},
			TestsPassed:        5,
			Coverage:           80.0,
			AcceptanceCriteria: []task.AcceptanceCriterion{{Criterion: "works", Met: true}},
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "WARNING") {
			t.Errorf("unexpected warning for complete record: %s", output)
		}
	})

	t.Run("completed with unmet AC triggers hard error", func(t *testing.T) {
		if os.Getenv("TEST_VALIDATE_UNMET_AC") == "1" {
			validateRecordData(&task.RecordData{
				Status:      "completed",
				Summary:     "Partial",
				TestsPassed: 1,
				Coverage:    50.0,
				AcceptanceCriteria: []task.AcceptanceCriterion{
					{Criterion: "works", Met: true},
					{Criterion: "edge case", Met: false},
				},
			}, false)
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestValidateRecordData/completed_with_unmet_AC_triggers_hard_error")
		cmd.Env = append(os.Environ(), "TEST_VALIDATE_UNMET_AC=1")
		err := cmd.Run()
		if err == nil {
			t.Error("expected non-zero exit for completed with unmet AC")
		}
	})

	t.Run("blocked with unmet AC is allowed", func(t *testing.T) {
		rd := &task.RecordData{
			Status:      "blocked",
			Summary:     "Blocked",
			TestsPassed: 0,
			TestsFailed: 0,
			Coverage:    0,
			AcceptanceCriteria: []task.AcceptanceCriterion{
				{Criterion: "works", Met: false},
			},
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "ERROR") {
			t.Errorf("blocked status should allow unmet AC, got: %s", output)
		}
	})

	t.Run("force overrides test evidence check", func(t *testing.T) {
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Force override",
			TestsPassed: 0,
			TestsFailed: 0,
			Coverage:    0,
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, true)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "ERROR") {
			t.Errorf("force should override test evidence check, got: %s", output)
		}
	})

	t.Run("completed without recommended fields warns", func(t *testing.T) {
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Did the work",
			TestsPassed: 1,
			Coverage:    50.0,
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if !strings.Contains(output, "WARNING") {
			t.Errorf("expected warning in stderr, got: %s", output)
		}
		for _, field := range []string{"keyDecisions", "acceptanceCriteria"} {
			if !strings.Contains(output, field) {
				t.Errorf("expected warning to mention %q, got: %s", field, output)
			}
		}
	})

	t.Run("non-completed status skips all checks", func(t *testing.T) {
		rd := &task.RecordData{
			Status:  "blocked",
			Summary: "Blocked with reason",
		}
		old := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w
		validateRecordData(rd, false)
		_ = w.Close()
		os.Stderr = old

		buf := make([]byte, 1024)
		n, _ := r.Read(buf)
		output := string(buf[:n])

		if strings.Contains(output, "WARNING") {
			t.Errorf("non-completed status should not produce warnings: %s", output)
		}
	})
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
		criteria []task.AcceptanceCriterion
		want     string
	}{
		{
			name:     "empty criteria",
			criteria: []task.AcceptanceCriterion{},
			want:     "无",
		},
		{
			name: "single unmet criterion",
			criteria: []task.AcceptanceCriterion{
				{Criterion: "Feature works", Met: false},
			},
			want: "- [ ] Feature works",
		},
		{
			name: "single met criterion",
			criteria: []task.AcceptanceCriterion{
				{Criterion: "Feature works", Met: true},
			},
			want: "- [x] Feature works",
		},
		{
			name: "multiple mixed criteria",
			criteria: []task.AcceptanceCriterion{
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
		recordData       *task.RecordData
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
			recordData: &task.RecordData{
				Status:       "completed",
				Summary:      "Implemented the feature",
				FilesCreated: []string{"main.go"},
				TestsPassed:  5,
				Coverage:     85.5,
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"1.1",
				"Implement feature X",
				"Implemented the feature",
				"completed",
				"2026-04-06 10:00",
				"main.go",
				"**Passed**: 5",
				"**Coverage**: 85.5%",
			},
		},
		{
			name: "template with all fields",
			task: &task.Task{
				ID:    "2.1",
				Title: "Full feature",
			},
			recordData: &task.RecordData{
				Status:        "completed",
				Summary:       "Complete implementation",
				FilesCreated:  []string{"a.go", "b.go"},
				FilesModified: []string{"c.go"},
				KeyDecisions:  []string{"Used pattern X"},
				TestsPassed:   10,
				TestsFailed:   2,
				Coverage:      90.0,
				AcceptanceCriteria: []task.AcceptanceCriterion{
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
				"**Passed**: 10",
				"**Failed**: 2",
				"**Coverage**: 90.0%",
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
			recordData: &task.RecordData{
				Status:  "in_progress",
				Summary: "Work in progress",
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"status: \"in_progress\"",
				"completed: \"N/A\"",
			},
		},
		{
			name: "default notes when empty",
			task: &task.Task{
				ID:    "1.3",
				Title: "Task with no notes",
			},
			recordData: &task.RecordData{
				Status:      "completed",
				Summary:     "Done",
				TestsPassed: 1,
				Coverage:    50.0,
				Notes:       "",
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
			recordData: &task.RecordData{
				Status:      "completed",
				Summary:     "Done",
				TestsPassed: 1,
				Coverage:    50.0,
			},
			startedTime: "",
			checkContains: []string{
				"started:",
				"completed:",
			},
		},
		{
			name: "time spent when completed after started",
			task: &task.Task{
				ID:    "1.5",
				Title: "Timed Task",
			},
			recordData: &task.RecordData{
				Status:      "completed",
				Summary:     "Done",
				TestsPassed: 1,
				Coverage:    50.0,
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"time_spent:",
			},
		},
		{
			name: "no time spent when completed before started",
			task: &task.Task{
				ID:    "1.6",
				Title: "Backward Time Task",
			},
			recordData: &task.RecordData{
				Status:      "completed",
				Summary:     "Done",
				TestsPassed: 1,
				Coverage:    50.0,
			},
			startedTime: "2026-04-06 15:00",
			checkNotContains: []string{
				"time_spent: ~",
			},
		},
		{
			name: "noTest task with coverage=-1",
			task: &task.Task{
				ID:     "1.7",
				Title:  "Write PRD",
				NoTest: true,
			},
			recordData: &task.RecordData{
				Status:   "completed",
				Summary:  "Created PRD",
				Coverage: -1.0,
			},
			startedTime: "2026-04-06 10:00",
			checkContains: []string{
				"Tests Executed**: No (noTest task)",
				"Coverage**: N/A (task has no tests)",
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
		dir := t.TempDir()
		dataPath := dir + "/data.json"
		jsonData := `{"status":"completed","summary":"Done","testsPassed":5,"coverage":80.5}`
		if err := os.WriteFile(dataPath, []byte(jsonData), 0644); err != nil {
			t.Fatalf("failed to write temp file: %v", err)
		}

		rd, err := readSubmitData(dataPath)
		if err != nil {
			t.Fatalf("readSubmitData() error = %v", err)
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
		_, err := readSubmitData("/nonexistent/path/file.json")
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

		_, err := readSubmitData(dataPath)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
	})

	t.Run("no input without data flag", func(t *testing.T) {
		_, err := readSubmitData("")
		if err == nil {
			t.Error("expected error when no data provided")
		}
		if !strings.Contains(err.Error(), "no input") {
			t.Errorf("error should mention 'no input', got: %v", err)
		}
	})
}

func TestFormatCoverage(t *testing.T) {
	tests := []struct {
		input float64
		want  string
	}{
		{-1.0, "N/A (task has no tests)"},
		{85.5, "85.5%"},
		{0.0, "0.0%"},
		{100.0, "100.0%"},
	}
	for _, tt := range tests {
		got := formatCoverage(tt.input)
		if got != tt.want {
			t.Errorf("formatCoverage(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatTestsExecuted(t *testing.T) {
	tests := []struct {
		coverage float64
		noTest   bool
		want     string
	}{
		{-1.0, true, "No (noTest task)"},
		{-1.0, false, "No"},
		{0.0, false, "Yes"},
		{85.5, false, "Yes"},
		{0.0, true, "No (noTest task)"},
		{80.0, true, "No (noTest task)"},
	}
	for _, tt := range tests {
		got := formatTestsExecuted(tt.coverage, tt.noTest)
		if got != tt.want {
			t.Errorf("formatTestsExecuted(%v, %v) = %q, want %q", tt.coverage, tt.noTest, got, tt.want)
		}
	}
}

func TestSaveIndexAndSignalCompletion(t *testing.T) {
	t.Run("all tasks completed writes forge state", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("CLAUDE_PROJECT_DIR", dir)
		featureDir := filepath.Join(dir, "docs", "features", "test-f")
		tasksDir := filepath.Join(featureDir, "tasks")
		_ = os.MkdirAll(tasksDir, 0755)

		indexPath := filepath.Join(tasksDir, "index.json")
		index := &task.TaskIndex{
			Feature: "test-f",
		}
		index.SetTasks(map[string]task.Task{
			"t1": {ID: "1.1", Title: "Done", Status: "completed", Priority: "P0", File: "1.1.md"},
			"t2": {ID: "1.2", Title: "Skipped", Status: "skipped", Priority: "P1", File: "1.2.md"},
		})

		saveIndexAndSignalCompletion(indexPath, dir, "test-f", index)

		// Verify index was saved
		data, err := os.ReadFile(indexPath)
		if err != nil {
			t.Fatalf("index not saved: %v", err)
		}
		if !strings.Contains(string(data), "completed") {
			t.Error("index should contain completed status")
		}

		// Verify forge state was written
		statePath := filepath.Join(dir, ".forge", "state.json")
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			t.Error("forge state should be written when all tasks done")
		}
	})

	t.Run("incomplete tasks does not write forge state", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("CLAUDE_PROJECT_DIR", dir)
		featureDir := filepath.Join(dir, "docs", "features", "test-f")
		tasksDir := filepath.Join(featureDir, "tasks")
		_ = os.MkdirAll(tasksDir, 0755)

		indexPath := filepath.Join(tasksDir, "index.json")
		index := &task.TaskIndex{
			Feature: "test-f",
		}
		index.SetTasks(map[string]task.Task{
			"t1": {ID: "1.1", Title: "Done", Status: "completed", Priority: "P0", File: "1.1.md"},
			"t2": {ID: "1.2", Title: "Pending", Status: "pending", Priority: "P1", File: "1.2.md"},
		})

		saveIndexAndSignalCompletion(indexPath, dir, "test-f", index)

		// Verify forge state was NOT written
		statePath := filepath.Join(dir, ".forge", "state.json")
		if _, err := os.Stat(statePath); err == nil {
			t.Error("forge state should NOT be written when tasks are pending")
		}
	})
}

func TestAutoRestoreSourceTask(t *testing.T) {
	t.Run("restores blocked source when all deps completed", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":   {ID: "src", Status: "blocked", Dependencies: []string{"fix-1"}},
			"fix-1": {ID: "fix-1", Status: "completed", SourceTaskID: "src"},
		})

		autoRestoreSourceTask(index, "src")

		if index.TasksMap()["src"].Status != "pending" {
			t.Errorf("expected pending, got %s", index.TasksMap()["src"].Status)
		}
	})

	t.Run("does not restore when some deps incomplete", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":   {ID: "src", Status: "blocked", Dependencies: []string{"fix-1", "fix-2"}},
			"fix-1": {ID: "fix-1", Status: "completed"},
			"fix-2": {ID: "fix-2", Status: "pending"},
		})

		autoRestoreSourceTask(index, "src")

		if index.TasksMap()["src"].Status != "blocked" {
			t.Errorf("source should stay blocked when deps incomplete, got %s", index.TasksMap()["src"].Status)
		}
	})

	t.Run("no-op when source is not blocked", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":   {ID: "src", Status: "in_progress", Dependencies: []string{"fix-1"}},
			"fix-1": {ID: "fix-1", Status: "completed"},
		})

		autoRestoreSourceTask(index, "src")

		if index.TasksMap()["src"].Status != "in_progress" {
			t.Errorf("source should stay in_progress, got %s", index.TasksMap()["src"].Status)
		}
	})

	t.Run("no-op when source not found", func(_ *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}

		autoRestoreSourceTask(index, "nonexistent")
	})

	t.Run("blocked with no deps restores to pending", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src": {ID: "src", Status: "blocked"},
		})

		autoRestoreSourceTask(index, "src")

		if index.TasksMap()["src"].Status != "pending" {
			t.Errorf("blocked with no deps should restore, got %s", index.TasksMap()["src"].Status)
		}
	})

	t.Run("nested chain: fix-B restores fix-A only", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":   {ID: "src", Status: "blocked", Dependencies: []string{"fix-A"}},
			"fix-A": {ID: "fix-A", Status: "blocked", Dependencies: []string{"fix-B"}, SourceTaskID: "src"},
			"fix-B": {ID: "fix-B", Status: "completed", SourceTaskID: "fix-A"},
		})

		autoRestoreSourceTask(index, "fix-A")

		if index.TasksMap()["fix-A"].Status != "pending" {
			t.Errorf("fix-A should be restored to pending, got %s", index.TasksMap()["fix-A"].Status)
		}
		if index.TasksMap()["src"].Status != "blocked" {
			t.Errorf("src should stay blocked until fix-A completes, got %s", index.TasksMap()["src"].Status)
		}
	})

	t.Run("skipped dep counts as completed (aligned with validate)", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":   {ID: "src", Status: "blocked", Dependencies: []string{"fix-1", "fix-2"}},
			"fix-1": {ID: "fix-1", Status: "completed"},
			"fix-2": {ID: "fix-2", Status: "skipped"},
		})

		autoRestoreSourceTask(index, "src")

		if index.TasksMap()["src"].Status != "pending" {
			t.Errorf("source should be restored when deps are completed/skipped, got %s", index.TasksMap()["src"].Status)
		}
	})
	t.Run("skipped fix-task triggers auto-restore", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":   {ID: "src", Status: "blocked", Dependencies: []string{"fix-1"}},
			"fix-1": {ID: "fix-1", Status: "skipped", SourceTaskID: "src"},
		})

		// Simulate what record.go does: check SourceTaskID with skipped status
		if index.TasksMap()["fix-1"].SourceTaskID != "" {
			autoRestoreSourceTask(index, "src")
		}

		if index.TasksMap()["src"].Status != "pending" {
			t.Errorf("source should be restored when skipped fix-task completes, got %s", index.TasksMap()["src"].Status)
		}
	})
}

func TestAutoRestoreSourceTask_WildcardDeps(t *testing.T) {
	t.Run("restores when wildcard deps all completed", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src":    {ID: "src", Status: "blocked", Dependencies: []string{"1.x", "fix-1"}},
			"1.1":    {ID: "1.1", Status: "completed"},
			"1.2":    {ID: "1.2", Status: "completed"},
			"1.gate": {ID: "1.gate", Status: "pending"},
			"fix-1":  {ID: "fix-1", Status: "completed"},
		})

		autoRestoreSourceTask(index, "src")

		if index.TasksMap()["src"].Status != "pending" {
			t.Errorf("should restore with wildcard deps all completed, got %s", index.TasksMap()["src"].Status)
		}
	})

	t.Run("does not restore when wildcard dep has pending task", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"src": {ID: "src", Status: "blocked", Dependencies: []string{"1.x"}},
			"1.1": {ID: "1.1", Status: "completed"},
			"1.2": {ID: "1.2", Status: "pending"},
		})

		autoRestoreSourceTask(index, "src")

		if index.TasksMap()["src"].Status != "blocked" {
			t.Errorf("should stay blocked when wildcard dep has pending, got %s", index.TasksMap()["src"].Status)
		}
	})
}

func TestAutoRestoreSourceTask_KeyDiffersFromID(t *testing.T) {
	t.Run("restores blocked source by ID when key is slug", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e-tests": {ID: "T-test-3", Status: "blocked", Dependencies: []string{"fix-1"}},
			"fix-1":         {ID: "fix-1", Status: "completed", SourceTaskID: "T-test-3"},
		})

		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e-tests"].Status != "pending" {
			t.Errorf("expected pending, got %s", index.TasksMap()["run-e2e-tests"].Status)
		}
	})

	t.Run("no-op when source not found by ID", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e-tests": {ID: "T-test-3", Status: "blocked", Dependencies: []string{"fix-1"}},
			"fix-1":         {ID: "fix-1", Status: "completed"},
		})

		autoRestoreSourceTask(index, "nonexistent-id")
		if index.TasksMap()["run-e2e-tests"].Status != "blocked" {
			t.Errorf("should stay blocked, got %s", index.TasksMap()["run-e2e-tests"].Status)
		}
	})

	t.Run("stays blocked when some deps incomplete (slug-keyed source)", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e-tests": {ID: "T-test-3", Status: "blocked", Dependencies: []string{"fix-1", "fix-2"}},
			"fix-1":         {ID: "fix-1", Status: "completed"},
			"fix-2":         {ID: "fix-2", Status: "pending"},
		})

		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e-tests"].Status != "blocked" {
			t.Errorf("should stay blocked with incomplete deps, got %s", index.TasksMap()["run-e2e-tests"].Status)
		}
	})

	t.Run("restores by key when key equals ID (dynamic task)", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"disc-1": {ID: "disc-1", Status: "blocked", Dependencies: []string{"fix-1"}},
			"fix-1":  {ID: "fix-1", Status: "completed"},
		})

		autoRestoreSourceTask(index, "disc-1")

		if index.TasksMap()["disc-1"].Status != "pending" {
			t.Errorf("expected pending, got %s", index.TasksMap()["disc-1"].Status)
		}
	})

	t.Run("write-back uses correct slug key, does not create duplicate entry", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e-tests": {ID: "T-test-3", Status: "blocked", Dependencies: []string{"fix-1"}},
			"fix-1":         {ID: "fix-1", Status: "completed"},
		})

		autoRestoreSourceTask(index, "T-test-3")

		_, hasSlugKey := index.TasksMap()["run-e2e-tests"]
		_, hasIDKey := index.TasksMap()["T-test-3"]
		if !hasSlugKey {
			t.Error("slug key 'run-e2e-tests' was lost after restore")
		}
		if hasIDKey {
			t.Error("should not create duplicate entry under ID key 'T-test-3'")
		}
		if index.TasksMap()["run-e2e-tests"].Status != "pending" {
			t.Errorf("expected pending, got %s", index.TasksMap()["run-e2e-tests"].Status)
		}
	})
}

func TestAutoRestoreSourceTask_SlugKeyedFullChain(t *testing.T) {
	t.Run("slug source with slug-keyed deps restores", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e":  {ID: "T-test-3", Status: "blocked", Dependencies: []string{"T-fix-7"}},
			"fix-auth": {ID: "T-fix-7", Status: "completed"},
		})

		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e"].Status != "pending" {
			t.Errorf("should restore with slug-keyed dep completed, got %s", index.TasksMap()["run-e2e"].Status)
		}
	})

	t.Run("slug source with wildcard deps all completed", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e": {ID: "T-test-3", Status: "blocked", Dependencies: []string{"1.x"}},
			"1.1":     {ID: "1.1", Status: "completed"},
			"1.2":     {ID: "1.2", Status: "completed"},
		})

		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e"].Status != "pending" {
			t.Errorf("should restore with wildcard all completed, got %s", index.TasksMap()["run-e2e"].Status)
		}
	})

	t.Run("slug source completed status is no-op", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e":  {ID: "T-test-3", Status: "completed", Dependencies: []string{"T-fix-7"}},
			"fix-auth": {ID: "T-fix-7", Status: "completed"},
		})

		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e"].Status != "completed" {
			t.Error("completed source should not be changed")
		}
	})

	t.Run("slug source skipped status is no-op", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e": {ID: "T-test-3", Status: "skipped", Dependencies: []string{}},
		})

		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e"].Status != "skipped" {
			t.Error("skipped source should not be changed")
		}
	})

	t.Run("slug blocked source with no deps restores", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e": {ID: "T-test-3", Status: "blocked", Dependencies: []string{}},
		})

		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e"].Status != "pending" {
			t.Errorf("blocked with no deps should restore, got %s", index.TasksMap()["run-e2e"].Status)
		}
	})

	t.Run("idempotent: second call is no-op for slug-keyed task", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e":  {ID: "T-test-3", Status: "blocked", Dependencies: []string{"T-fix-7"}},
			"fix-auth": {ID: "T-fix-7", Status: "completed"},
		})

		autoRestoreSourceTask(index, "T-test-3")
		// Second call: status is now "pending", should be no-op
		autoRestoreSourceTask(index, "T-test-3")

		if index.TasksMap()["run-e2e"].Status != "pending" {
			t.Errorf("should stay pending after second call, got %s", index.TasksMap()["run-e2e"].Status)
		}
	})

	t.Run("nested chain: both source and fix-A are slug-keyed", func(t *testing.T) {
		index := &task.TaskIndex{
			Feature: "test",
		}
		index.SetTasks(map[string]task.Task{
			"run-e2e":  {ID: "T-test-3", Status: "blocked", Dependencies: []string{"T-fix-A"}},
			"fix-auth": {ID: "T-fix-A", Status: "blocked", Dependencies: []string{"fix-B"}, SourceTaskID: "T-test-3"},
			"fix-B":    {ID: "fix-B", Status: "completed", SourceTaskID: "T-fix-A"},
		})

		// Restoring fix-A (slug-keyed, has completed dep fix-B)
		autoRestoreSourceTask(index, "T-fix-A")

		if index.TasksMap()["fix-auth"].Status != "pending" {
			t.Errorf("fix-A should be restored, got %s", index.TasksMap()["fix-auth"].Status)
		}
		if index.TasksMap()["run-e2e"].Status != "blocked" {
			t.Errorf("source should stay blocked (fix-A not completed yet), got %s", index.TasksMap()["run-e2e"].Status)
		}
	})
}

func TestFillRecordTemplate_TypeReclassification(t *testing.T) {
	t.Run("nil TypeReclassification omits block", func(t *testing.T) {
		tmpl := &task.Task{ID: "1.1", Title: "Test Task"}
		rd := &task.RecordData{
			Status:               "completed",
			Summary:              "Done",
			TestsPassed:          1,
			Coverage:             50.0,
			TypeReclassification: nil,
		}
		content := fillRecordTemplate(tmpl, rd, "2026-01-01 10:00")
		if strings.Contains(content, "Type Reclassification") {
			t.Error("should not contain Type Reclassification block when nil")
		}
	})

	t.Run("non-nil TypeReclassification renders block", func(t *testing.T) {
		tmpl := &task.Task{ID: "fix-1", Title: "Fix compile error"}
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Was actually a cleanup",
			TestsPassed: 3,
			Coverage:    80.0,
			TypeReclassification: &task.TypeReclassification{
				OriginalType: "fix",
				ActualType:   "cleanup",
				Reason:       "e2e test TestTC_003_Login has race condition in assertion timing",
			},
		}
		content := fillRecordTemplate(tmpl, rd, "2026-01-01 10:00")

		if !strings.Contains(content, "## Type Reclassification") {
			t.Error("should contain Type Reclassification heading")
		}
		if !strings.Contains(content, "- Original: fix") {
			t.Error("should contain original type")
		}
		if !strings.Contains(content, "- Actual: cleanup") {
			t.Error("should contain actual type")
		}
		if !strings.Contains(content, "- Reason: e2e test TestTC_003_Login has race condition") {
			t.Error("should contain reason")
		}
	})

	t.Run("TypeReclassification appears between Summary and Changes", func(t *testing.T) {
		tmpl := &task.Task{ID: "fix-1", Title: "Fix task"}
		rd := &task.RecordData{
			Status:      "completed",
			Summary:     "Changed type",
			TestsPassed: 1,
			Coverage:    50.0,
			TypeReclassification: &task.TypeReclassification{
				OriginalType: "fix",
				ActualType:   "cleanup",
				Reason:       "flaky test",
			},
		}
		content := fillRecordTemplate(tmpl, rd, "2026-01-01 10:00")

		summaryIdx := strings.Index(content, "## Summary")
		reclassIdx := strings.Index(content, "## Type Reclassification")
		changesIdx := strings.Index(content, "## Changes")

		if summaryIdx == -1 || reclassIdx == -1 || changesIdx == -1 {
			t.Fatal("expected all three sections to be present")
		}
		if summaryIdx >= reclassIdx || reclassIdx >= changesIdx {
			t.Error("Type Reclassification should appear between Summary and Changes")
		}
	})
}

func TestFillRecordTemplate_RejectedStatus(t *testing.T) {
	tmpl := &task.Task{ID: "1.1", Title: "Test Task"}
	rd := &task.RecordData{Status: "rejected", Summary: "Did not pass acceptance criteria"}
	content := fillRecordTemplate(tmpl, rd, "2026-01-01 10:00")
	if !strings.Contains(content, `status: "rejected"`) {
		t.Error("template should contain rejected status")
	}
	if !strings.Contains(content, "N/A") {
		t.Error("completed time should be N/A for non-completed status")
	}
}

func TestSaveIndexAndSignalCompletion_RejectedNotDone(t *testing.T) {
	projectRoot := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", projectRoot)
	featureSlug := "test"
	tasksDir := filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks")
	_ = os.MkdirAll(tasksDir, 0755)

	index := task.NewTaskIndex(featureSlug)
	index.SetTasks(map[string]task.Task{
		"task-a": {ID: "1.1", Status: "completed", File: "1.1.md", Record: "1.1-record.md"},
		"task-b": {ID: "1.2", Status: "rejected", File: "1.2.md", Record: "1.2-record.md"},
	})
	indexPath := filepath.Join(tasksDir, "index.json")
	_ = task.SaveIndex(indexPath, index)

	saveIndexAndSignalCompletion(indexPath, projectRoot, featureSlug, index)

	// Should NOT write forge state because rejected != done
	forgeState := feature.ReadForgeState(projectRoot)
	if forgeState != nil && forgeState.AllCompleted {
		t.Error("rejected task should prevent allCompleted signal")
	}
}

func TestAutoRestoreSourceTask_RejectedDepNotMet(t *testing.T) {
	index := &task.TaskIndex{Feature: "test"}
	index.SetTasks(map[string]task.Task{
		"source": {ID: "1.1", Status: "blocked", Dependencies: []string{"1.2"}},
		"dep":    {ID: "1.2", Status: "rejected"},
	})
	// Should not restore because dep is rejected (not completed/skipped)
	autoRestoreSourceTask(index, "1.1")
	src, _ := index.ByID("source")
	if src.Status != "blocked" {
		t.Errorf("source should stay blocked when dep is rejected, got %s", src.Status)
	}
}

func TestValidateRecordData_RejectedSkipsCompletedChecks(_ *testing.T) {
	// Rejected status should skip test evidence and AC checks
	rd := &task.RecordData{
		Status:   "rejected",
		Summary:  "Acceptance criteria not met",
		Coverage: -1.0,
	}
	// Should not exit or error — rejected skips completed validation
	validateRecordData(rd, false)
}

// TestRecordExistsCheck tests the write-once protection for record files.
// Uses the subprocess pattern (like TestValidateRecordData) because runSubmit calls Exit().
func TestRecordExistsCheck(t *testing.T) {
	t.Run("submit fails when record already exists without force", func(t *testing.T) {
		if os.Getenv("TEST_RECORD_EXISTS_NO_FORCE") == "1" {
			setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1", Title: "T1", Status: "pending", File: "1.md", Record: "records/1.md"},
				},
			})

			dir, _ := os.Getwd()
			recordPath := filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.md")
			_ = os.WriteFile(recordPath, []byte("existing record"), 0644)

			dataPath := filepath.Join(dir, "record.json")
			jsonData := `{"status":"completed","summary":"Done","testsPassed":1,"coverage":50.0}`
			_ = os.WriteFile(dataPath, []byte(jsonData), 0644)

			submitDataPath = dataPath
			submitForce = false
			runSubmit(submitCmd, []string{"1"})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestRecordExistsCheck/submit_fails_when_record_already_exists_without_force")
		cmd.Env = append(os.Environ(), "TEST_RECORD_EXISTS_NO_FORCE=1")
		output, _ := cmd.CombinedOutput()
		if !strings.Contains(string(output), "already exists") {
			t.Errorf("expected 'already exists' error, got: %s", string(output))
		}
		if !strings.Contains(string(output), "VALIDATION_ERROR") {
			t.Errorf("expected VALIDATION_ERROR code, got: %s", string(output))
		}
		if !strings.Contains(string(output), "--force") {
			t.Errorf("expected hint to mention --force, got: %s", string(output))
		}
	})

	t.Run("submit with --force overwrites existing record with warning", func(t *testing.T) {
		if os.Getenv("TEST_RECORD_EXISTS_FORCE") == "1" {
			setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1", Title: "T1", Status: "pending", File: "1.md", Record: "records/1.md"},
				},
			})

			dir, _ := os.Getwd()
			recordPath := filepath.Join(dir, "docs", "features", "test", "tasks", "records", "1.md")
			_ = os.WriteFile(recordPath, []byte("existing record"), 0644)

			dataPath := filepath.Join(dir, "record.json")
			jsonData := `{"status":"completed","summary":"Overwritten","testsPassed":2,"coverage":60.0,"keyDecisions":["d1"],"acceptanceCriteria":[{"criterion":"works","met":true}]}`
			_ = os.WriteFile(dataPath, []byte(jsonData), 0644)

			submitDataPath = dataPath
			submitForce = true
			runSubmit(submitCmd, []string{"1"})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestRecordExistsCheck/submit_with_--force_overwrites_existing_record_with_warning")
		cmd.Env = append(os.Environ(), "TEST_RECORD_EXISTS_FORCE=1")
		output, _ := cmd.CombinedOutput()
		out := string(output)
		if !strings.Contains(out, "WARNING") || !strings.Contains(out, "Overwriting") {
			t.Errorf("expected WARNING about overwriting on stderr, got: %s", out)
		}
		// Should succeed (exit code 0)
		if !cmd.ProcessState.Success() {
			t.Errorf("expected success exit, got: %s", out)
		}
	})

	t.Run("submit succeeds when record does not exist", func(t *testing.T) {
		if os.Getenv("TEST_RECORD_NOT_EXISTS") == "1" {
			setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1", Title: "T1", Status: "pending", File: "1.md", Record: "records/1.md"},
				},
			})

			dir, _ := os.Getwd()
			// Ensure record file does NOT exist (it shouldn't by default)

			dataPath := filepath.Join(dir, "record.json")
			jsonData := `{"status":"completed","summary":"New record","testsPassed":3,"coverage":70.0,"keyDecisions":["d1"],"acceptanceCriteria":[{"criterion":"works","met":true}]}`
			_ = os.WriteFile(dataPath, []byte(jsonData), 0644)

			submitDataPath = dataPath
			submitForce = false
			runSubmit(submitCmd, []string{"1"})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestRecordExistsCheck/submit_succeeds_when_record_does_not_exist")
		cmd.Env = append(os.Environ(), "TEST_RECORD_NOT_EXISTS=1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("expected success, got error: %v, output: %s", err, string(output))
		}
	})
}

func TestSubmit_NonTestableTypeSkipsQualityGate(t *testing.T) {
	t.Run("documentation type skips quality gate", func(t *testing.T) {
		if os.Getenv("TEST_SUBMIT_DOC_SKIPS_QG") == "1" {
			setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1", Title: "Doc Task", Status: "pending", File: "1.md", Record: "records/1.md", Type: task.TypeDocumentation},
				},
			})

			dir, _ := os.Getwd()
			dataPath := filepath.Join(dir, "record.json")
			jsonData := `{"status":"completed","summary":"Doc task done","coverage":-1.0}`
			_ = os.WriteFile(dataPath, []byte(jsonData), 0644)

			submitDataPath = dataPath
			submitForce = false
			runSubmit(submitCmd, []string{"1"})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestSubmit_NonTestableTypeSkipsQualityGate/documentation_type_skips_quality_gate")
		cmd.Env = append(os.Environ(), "TEST_SUBMIT_DOC_SKIPS_QG=1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("documentation type should skip quality gate, got error: %v, output: %s", err, string(output))
		}
	})

	t.Run("feature type runs quality gate", func(t *testing.T) {
		if os.Getenv("TEST_SUBMIT_FEAT_RUNS_QG") == "1" {
			setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1", Title: "Feature Task", Status: "pending", File: "1.md", Record: "records/1.md", Type: task.TypeFeature},
				},
			})

			dir, _ := os.Getwd()
			// Create a justfile so RunGate actually attempts execution
			justfile := "compile:\n\t@echo \"compile fails\" && exit 1\nfmt:\n\t@true\nlint:\n\t@true\ntest:\n\t@true\n"
			_ = os.WriteFile(filepath.Join(dir, "justfile"), []byte(justfile), 0644)

			dataPath := filepath.Join(dir, "record.json")
			jsonData := `{"status":"completed","summary":"Feature done","testsPassed":3,"coverage":80.0}`
			_ = os.WriteFile(dataPath, []byte(jsonData), 0644)

			submitDataPath = dataPath
			submitForce = false
			runSubmit(submitCmd, []string{"1"})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestSubmit_NonTestableTypeSkipsQualityGate/feature_type_runs_quality_gate")
		cmd.Env = append(os.Environ(), "TEST_SUBMIT_FEAT_RUNS_QG=1")
		output, _ := cmd.CombinedOutput()
		out := string(output)
		if !strings.Contains(out, "Quality gate failed") {
			t.Errorf("feature type should run quality gate, got: %s", out)
		}
	})
}

func TestSubmit_NonTestableTypeAutoSetCoverage(t *testing.T) {
	t.Run("documentation type auto-sets coverage to -1", func(t *testing.T) {
		if os.Getenv("TEST_SUBMIT_DOC_AUTO_COV") == "1" {
			setupFullProject(t, SetupOpts{
				Tasks: map[string]task.Task{
					"t1": {ID: "1", Title: "Doc Task", Status: "pending", File: "1.md", Record: "records/1.md", Type: task.TypeDocumentation},
				},
			})

			dir, _ := os.Getwd()
			dataPath := filepath.Join(dir, "record.json")
			jsonData := `{"status":"completed","summary":"Doc task done"}`
			_ = os.WriteFile(dataPath, []byte(jsonData), 0644)

			submitDataPath = dataPath
			submitForce = false
			runSubmit(submitCmd, []string{"1"})
			return
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestSubmit_NonTestableTypeAutoSetCoverage/documentation_type_auto-sets_coverage_to_-1")
		cmd.Env = append(os.Environ(), "TEST_SUBMIT_DOC_AUTO_COV=1")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Errorf("documentation type should auto-set coverage and succeed, got error: %v, output: %s", err, string(output))
		}
	})
}
