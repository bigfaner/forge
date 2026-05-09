package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/just"
	"task-cli/pkg/task"
)

func TestCheckAllCompleted(t *testing.T) {
	tests := []struct {
		name        string
		tasks       map[string]task.Task
		testCommand string
		forgeState  bool
		wantNil     bool
		wantTestCmd string
	}{
		{
			name: "all completed with forge state returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "completed"},
			},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "all skipped with forge state returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "skipped"},
			},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "mixed completed and skipped with forge state returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "skipped"},
			},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "one pending task returns nil even with forge state",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "pending"},
			},
			forgeState: true,
			wantNil:    true,
		},
		{
			name: "in_progress task returns nil",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "in_progress"},
			},
			wantNil: true,
		},
		{
			name: "blocked task returns nil",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "blocked"},
			},
			wantNil: true,
		},
		{
			name:       "empty task list with forge state returns result",
			tasks:      map[string]task.Task{},
			forgeState: true,
			wantNil:    false,
		},
		{
			name: "testCommand from index.json is propagated",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			testCommand: "make test",
			forgeState:  true,
			wantNil:     false,
			wantTestCmd: "make test",
		},
		{
			name: "no forge state returns nil even if all tasks completed",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "completed"},
			},
			forgeState: false,
			wantNil:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("CLAUDE_PROJECT_DIR", dir)

			if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
				t.Fatal(err)
			}

			indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
			index := &task.TaskIndex{
				Feature:     "test",
				StatusEnum:  []string{"pending", "in_progress", "completed", "blocked", "skipped"},
				TestCommand: tc.testCommand,
			}
			index.SetTasks(tc.tasks)
			if err := task.SaveIndex(indexPath, index); err != nil {
				t.Fatal(err)
			}

			if tc.forgeState {
				if err := feature.WriteForgeState(dir, "test"); err != nil {
					t.Fatal(err)
				}
			}

			result, err := checkAllCompleted(false)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tc.wantNil {
				if result != nil {
					t.Errorf("expected nil result, got %+v", result)
				}
				return
			}

			if result == nil {
				t.Fatal("expected non-nil result, got nil")
			}

			if result.FeatureSlug != "test" {
				t.Errorf("FeatureSlug = %q, want %q", result.FeatureSlug, "test")
			}
			if result.ProjectRoot == "" {
				t.Error("ProjectRoot should not be empty")
			}

			if result.TestCommand != tc.wantTestCmd {
				t.Errorf("TestCommand = %q, want %q", result.TestCommand, tc.wantTestCmd)
			}
		})
	}
}

func TestCheckAllCompleted_NoFeature(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := os.MkdirAll(filepath.Join(dir, feature.FeaturesDir), 0755); err != nil {
		t.Fatal(err)
	}

	result, err := checkAllCompleted(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result when no feature set, got %+v", result)
	}
}

func TestCheckAllCompleted_NoProject(t *testing.T) {
	t.Setenv("CLAUDE_PROJECT_DIR", "")

	result, err := checkAllCompleted(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result when no project root, got %+v", result)
	}
}

func TestHasJustfile(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  bool
	}{
		{name: "no justfile", files: []string{}, want: false},
		{name: "lowercase justfile", files: []string{"justfile"}, want: true},
		{name: "capitalized Justfile", files: []string{"Justfile"}, want: true},
		{name: "both present", files: []string{"justfile", "Justfile"}, want: true},
		{name: "unrelated files only", files: []string{"Makefile", "go.mod"}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tc.files {
				if err := os.WriteFile(filepath.Join(dir, f), []byte("test:\n    echo ok\n"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			if got := just.HasJustfile(dir); got != tc.want {
				t.Errorf("just.HasJustfile() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasJustRecipe(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	t.Run("recipe exists", func(t *testing.T) {
		dir := t.TempDir()
		content := "test:\n    echo ok\n"
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if !just.HasRecipe(dir, "test") {
			t.Error("just.HasRecipe() = false, want true for existing recipe")
		}
	})

	t.Run("recipe does not exist", func(t *testing.T) {
		dir := t.TempDir()
		content := "build:\n    echo build\n"
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if just.HasRecipe(dir, "test") {
			t.Error("just.HasRecipe() = true, want false for missing recipe")
		}
	})

	t.Run("no justfile", func(t *testing.T) {
		dir := t.TempDir()
		if just.HasRecipe(dir, "test") {
			t.Error("just.HasRecipe() = true, want false when no justfile")
		}
	})
}

func TestWriteLatestMd(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	t.Run("fail status", func(t *testing.T) {
		stats := TestStats{Fail: 1}

		err := writeLatestMd(dir, "test", stats)
		if err != nil {
			t.Fatalf("writeLatestMd() error = %v", err)
		}

		resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
		data, err := os.ReadFile(filepath.Join(resultsDir, "latest.md"))
		if err != nil {
			t.Fatalf("failed to read latest.md: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "FAIL") {
			t.Error("latest.md should show FAIL status")
		}
		if !strings.Contains(content, "raw-output.txt") {
			t.Error("latest.md should reference raw-output.txt")
		}
	})

	t.Run("pass status", func(t *testing.T) {
		stats := TestStats{Total: 5, Pass: 5}

		err := writeLatestMd(dir, "test", stats)
		if err != nil {
			t.Fatalf("writeLatestMd() error = %v", err)
		}

		resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
		data, err := os.ReadFile(filepath.Join(resultsDir, "latest.md"))
		if err != nil {
			t.Fatalf("failed to read latest.md: %v", err)
		}

		content := string(data)
		if !strings.Contains(content, "PASS") {
			t.Error("latest.md should show PASS status")
		}
	})
}

func TestWriteRawOutput(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	output := "not ok 1 - Login test\n  Error: expected 200, got 404"
	err := writeRawOutput(dir, "test", output)
	if err != nil {
		t.Fatalf("writeRawOutput() error = %v", err)
	}

	resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
	data, err := os.ReadFile(filepath.Join(resultsDir, "raw-output.txt"))
	if err != nil {
		t.Fatalf("failed to read raw-output.txt: %v", err)
	}

	if string(data) != output {
		t.Errorf("raw output mismatch: got %q, want %q", string(data), output)
	}
}




func TestExtractSourceFiles(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		{
			name:   "empty output",
			output: "",
			want:   "See error output for affected files",
		},
		{
			name:   "Go compile error",
			output: "./internal/handler.go:42:2: undefined: foo\n./internal/handler.go:43:1: too many arguments",
			want:   "internal/handler.go",
		},
		{
			name:   "Go lint error deduplicates",
			output: "pkg/service/user.go:108:1: S1000: should use for-range (gosimple)\npkg/service/user.go:200:3: SA1006: printf-style (staticcheck)",
			want:   "pkg/service/user.go",
		},
		{
			name:   "Go test error multiple files",
			output: "--- FAIL: TestHandler (0.00s)\n    handler_test.go:42: Expected 200, got 404\n    service_test.go:10: Error",
			want:   "handler_test.go, service_test.go",
		},
		{
			name:   "TypeScript error",
			output: "src/app.ts:42:5: error TS2304: Cannot find name 'foo'.",
			want:   "src/app.ts",
		},
		{
			name:   "multiple different files",
			output: "a.go:1: error\nb.go:2: error\nc.go:3: error",
			want:   "a.go, b.go, c.go",
		},
		{
			name:   "skips non-source extensions",
			output: "output.txt:10: something\nreport.json:5: error\nhandler.go:42: error",
			want:   "handler.go",
		},
		{
			name:   "no source files found",
			output: "some random output without file paths",
			want:   "See error output for affected files",
		},
		{
			name: "limits to 10 unique files",
			output: func() string {
				var lines []string
				for i := 0; i < 15; i++ {
					lines = append(lines, fmt.Sprintf("file%02d.go:%d: error", i, i+1))
				}
				return strings.Join(lines, "\n")
			}(),
			want: "file00.go, file01.go, file02.go, file03.go, file04.go, file05.go, file06.go, file07.go, file08.go, file09.go",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := extractSourceFiles(tc.output)
			if got != tc.want {
				t.Errorf("extractSourceFiles() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestAddFixTask(t *testing.T) {
	projectRoot := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", projectRoot)
	featureSlug := "test-feature"

	if err := feature.EnsureFeatureDir(projectRoot, featureSlug); err != nil {
		t.Fatal(err)
	}

	// Create index with one completed task so feature is valid
	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index := task.NewTaskIndex(featureSlug)
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed", File: "1.1.md"},
	})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	output := "./internal/handler.go:42:2: undefined: foo\n./internal/handler.go:43:1: too many arguments"
	errorDocPath := "tests/results/unit-raw-output.txt"

	taskID := addFixTask(projectRoot, featureSlug, "compile", output, errorDocPath)
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}

	// Verify task was added to index
	updatedIndex, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	addedTask, exists := updatedIndex.ByID(taskID)
	if !exists {
		t.Fatalf("task %s not found in index", taskID)
	}

	// Verify fix-task defaults applied
	if addedTask.Priority != "P0" {
		t.Errorf("priority = %q, want P0", addedTask.Priority)
	}
	if !addedTask.Breaking {
		t.Error("expected breaking=true")
	}
	if addedTask.Status != "pending" {
		t.Errorf("status = %q, want pending", addedTask.Status)
	}

	// Verify task markdown file was created
	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("task markdown file not found: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, errorDocPath) {
		t.Errorf("task markdown should reference error doc %q", errorDocPath)
	}
	if !strings.Contains(content, "internal/handler.go") {
		t.Error("task markdown should reference extracted source files")
	}
	if !strings.Contains(content, "just compile") {
		t.Error("task markdown should reference the failed command")
	}

	// Verify forge state was reset
	state := feature.ReadForgeState(projectRoot)
	if state == nil {
		t.Fatal("forge state should exist after addFixTask")
	}
	if state.AllCompleted {
		t.Error("allCompleted should be false after adding fix task")
	}
}

func TestAddFixTask_UnitTestStep(t *testing.T) {
	projectRoot := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", projectRoot)
	featureSlug := "test-feature"

	if err := feature.EnsureFeatureDir(projectRoot, featureSlug); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index := task.NewTaskIndex(featureSlug)
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed", File: "1.1.md"},
	})
	task.SaveIndex(indexPath, index)

	taskID := addFixTask(projectRoot, featureSlug, "unit-test", "handler_test.go:10: fail", "tests/results/unit-raw-output.txt")
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}

	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "just test") {
		t.Error("unit-test step should use 'just test' as test script")
	}
}

func TestHandleGateFailure_DistinctReasons(t *testing.T) {
	tests := []struct {
		step string
		want string
	}{
		{"compile", "Project compilation failed in all-completed hook"},
		{"lint", "Lint check failed in all-completed hook"},
		{"unit-test", "Unit tests failed in all-completed hook"},
		{"test-e2e", "E2e regression tests failed in all-completed hook"},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {
			// Capture the hook JSON output from handleGateFailure
			// It calls os.Exit(0), so we run it in a subprocess
			if os.Getenv("TEST_HANDLE_GATE") == "1" {
				handleGateFailure(tc.step, "tests/results/fake.txt", "some error")
				return
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestHandleGateFailure_DistinctReasons/"+tc.step)
			cmd.Env = append(os.Environ(), "TEST_HANDLE_GATE=1")
			output, _ := cmd.CombinedOutput()

			got := string(output)
			if !strings.Contains(got, tc.want) {
				t.Errorf("reason for step %q should contain %q, got:\n%s", tc.step, tc.want, got)
			}
			if !strings.Contains(got, "task claim") {
				t.Errorf("reason for step %q should contain 'task claim'", tc.step)
			}
			// Should NOT contain a specific task ID like disc-N
			if strings.Contains(got, "disc-") {
				t.Errorf("reason for step %q should not contain task ID", tc.step)
			}
		})
	}
}

func TestCheckAllCompleted_RejectedTaskReturnsNil(t *testing.T) {
	projectRoot := t.TempDir()
	featureSlug := "test"
	tasksDir := filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks")
	os.MkdirAll(tasksDir, 0755)

	index := task.NewTaskIndex(featureSlug)
	index.SetTasks(map[string]task.Task{
		"task-a": {ID: "1.1", Status: "completed", File: "1.1.md"},
		"task-b": {ID: "1.2", Status: "rejected", File: "1.2.md"},
	})
	indexPath := filepath.Join(tasksDir, "index.json")
	task.SaveIndex(indexPath, index)
	feature.WriteForgeState(projectRoot, featureSlug)

	result, _ := checkAllCompleted(false)
	if result != nil {
		t.Error("rejected task should prevent all-completed from proceeding")
	}
}
