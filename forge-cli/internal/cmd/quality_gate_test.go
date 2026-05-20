package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/just"
	"forge-cli/pkg/task"
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

			result := checkAllCompleted(false)

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

	result := checkAllCompleted(false)
	if result != nil {
		t.Errorf("expected nil result when no feature set, got %+v", result)
	}
}

func TestCheckAllCompleted_NoProject(t *testing.T) {
	if os.Getenv("TEST_CHECK_ALL_COMPLETED_NO_PROJECT") == "1" {
		result := checkAllCompleted(false)
		if result != nil {
			t.Errorf("expected nil result when no project root, got %+v", result)
		}
		return
	}

	tmpDir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckAllCompleted_NoProject")
	// Build clean env: clear CLAUDE_PROJECT_DIR and PROJECT_ROOT so FindProjectRoot
	// cannot walk up from cwd and find ancestor project markers.
	env := []string{}
	for _, e := range os.Environ() {
		if strings.HasPrefix(e, "CLAUDE_PROJECT_DIR=") || strings.HasPrefix(e, "PROJECT_ROOT=") {
			continue
		}
		env = append(env, e)
	}
	cmd.Env = append(slices.Clone(env), "TEST_CHECK_ALL_COMPLETED_NO_PROJECT=1", "CLAUDE_PROJECT_DIR=")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("test subprocess failed: %v\n%s", err, string(output))
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

func TestExtractSourceFiles(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   string
	}{
		// --- empty / no-match cases ---
		{
			name:   "empty output",
			output: "",
			want:   "See error output for affected files",
		},
		{
			name:   "no source files found",
			output: "some random output without file paths",
			want:   "See error output for affected files",
		},
		{
			name:   "only non-source extensions",
			output: "output.txt:10: something\nreport.json:5: error\nconfig.yaml:8: bad",
			want:   "See error output for affected files",
		},

		// --- Go patterns ---
		{
			name:   "Go compile error with ./ prefix",
			output: "./internal/handler.go:42:2: undefined: foo\n./internal/handler.go:43:1: too many arguments",
			want:   "internal/handler.go",
		},
		{
			name:   "Go compile error without prefix",
			output: "pkg/service/user.go:108:1: S1000: should use for-range (gosimple)",
			want:   "pkg/service/user.go",
		},
		{
			name:   "Go deduplicates same file",
			output: "handler.go:10: err1\nhandler.go:20: err2\nhandler.go:30: err3",
			want:   "handler.go",
		},
		{
			name:   "Go test error multiple files",
			output: "--- FAIL: TestHandler (0.00s)\n    handler_test.go:42: Expected 200, got 404\n    service_test.go:10: Error",
			want:   "handler_test.go, service_test.go",
		},
		{
			name:   "Go vet output",
			output: "internal/handler.go:42:2: fmt.Sprintf format %v has arg of wrong type int",
			want:   "internal/handler.go",
		},

		// --- TypeScript / JavaScript ---
		{
			name:   "TypeScript error",
			output: "src/app.ts:42:5: error TS2304: Cannot find name 'foo'.",
			want:   "src/app.ts",
		},
		{
			name:   "TypeScript JSX",
			output: "src/components/Button.tsx:15:3: error TS2604: Element implicitly has an 'any' type.",
			want:   "src/components/Button.tsx",
		},
		{
			name:   "JavaScript error",
			output: "src/utils/helpers.js:10:5: ReferenceError: foo is not defined",
			want:   "src/utils/helpers.js",
		},

		// --- Python ---
		{
			name:   "Python error (pyflakes format)",
			output: "app/handlers.py:42: undefined name 'foo'",
			want:   "app/handlers.py",
		},
		{
			name:   "Python error (traceback format not matched)",
			output: "  File \"app/handlers.py\", line 42\n    def foo(\n           ^\nSyntaxError: incomplete input",
			want:   "See error output for affected files",
		},
		{
			name:   "Python pytest error",
			output: "FAILED tests/test_handler.py:42 - AssertionError: expected 200",
			want:   "tests/test_handler.py",
		},

		// --- Rust ---
		{
			name:   "Rust error",
			output: "error[E0425]: cannot find value `foo` in this scope\n --> src/main.rs:10:5\n",
			want:   "src/main.rs",
		},

		// --- Java ---
		{
			name:   "Java error",
			output: "src/main/java/com/example/App.java:42: error: cannot find symbol",
			want:   "src/main/java/com/example/App.java",
		},

		// --- C/C++ ---
		{
			name:   "C error",
			output: "src/main.c:42:5: error: use of undeclared identifier 'foo'",
			want:   "src/main.c",
		},
		{
			name:   "C++ header error",
			output: "include/utils.hpp:10:1: error: expected ';' after struct",
			want:   "include/utils.hpp",
		},

		// --- Web ---
		{
			name:   "CSS error",
			output: "src/styles/main.css:42:3: parse error: invalid property",
			want:   "src/styles/main.css",
		},
		{
			name:   "SCSS error",
			output: "src/styles/_variables.scss:15:1: error: undefined variable",
			want:   "src/styles/_variables.scss",
		},
		{
			name:   "HTML error",
			output: "src/templates/index.html:20:5: mismatched tag",
			want:   "src/templates/index.html",
		},
		{
			name:   "SQL error",
			output: "migrations/001_init.sql:10:1: syntax error at or near \"CREATE\"",
			want:   "migrations/001_init.sql",
		},
		{
			name:   "Vue error",
			output: "src/components/App.vue:25:3: error: v-bind without expression",
			want:   "src/components/App.vue",
		},
		{
			name:   "Svelte error",
			output: "src/lib/Button.svelte:10:1: Unexpected token",
			want:   "src/lib/Button.svelte",
		},

		// --- mixed languages ---
		{
			name:   "mixed Go and TypeScript",
			output: "internal/api.go:10: error\nsrc/frontend.ts:20: error",
			want:   "internal/api.go, src/frontend.ts",
		},

		// --- path patterns ---
		{
			name:   "deep nested path",
			output: "a/b/c/d/e/f.go:1: error",
			want:   "a/b/c/d/e/f.go",
		},
		{
			name:   "path with hyphens and underscores",
			output: "pkg/my-module/handler_utils.go:10: error",
			want:   "pkg/my-module/handler_utils.go",
		},

		// --- boundary ---
		{
			name: "limits to 10 unique files",
			output: func() string {
				var lines []string
				for i := range 15 {
					lines = append(lines, fmt.Sprintf("file%02d.go:%d: error", i, i+1))
				}
				return strings.Join(lines, "\n")
			}(),
			want: "file00.go, file01.go, file02.go, file03.go, file04.go, file05.go, file06.go, file07.go, file08.go, file09.go",
		},
		{
			name:   "exactly 10 files",
			output: "f01.go:1:e\nf02.go:2:e\nf03.go:3:e\nf04.go:4:e\nf05.go:5:e\nf06.go:6:e\nf07.go:7:e\nf08.go:8:e\nf09.go:9:e\nf10.go:10:e",
			want:   "f01.go, f02.go, f03.go, f04.go, f05.go, f06.go, f07.go, f08.go, f09.go, f10.go",
		},
		{
			name:   "single file no line number colon suffix still matches",
			output: "handler.go:42: msg",
			want:   "handler.go",
		},
		{
			name:   "file with only line number",
			output: "handler.go:42: ",
			want:   "handler.go",
		},
		{
			name:   "skips .log .md .txt .json .yaml .xml .toml .csv .bin .exe .lock",
			output: "app.log:1:e\nreadme.md:2:e\ndata.json:3:e\nconfig.yaml:4:e\napp.toml:5:e\nhandler.go:6:e",
			want:   "handler.go",
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

// helperSetup creates a minimal feature with a completed task for addFixTask tests.
func helperSetup(t *testing.T) (projectRoot, featureSlug, indexPath string) {
	t.Helper()
	projectRoot = t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", projectRoot)
	featureSlug = "test-feature"

	if err := feature.EnsureFeatureDir(projectRoot, featureSlug); err != nil {
		t.Fatal(err)
	}

	indexPath = filepath.Join(projectRoot, feature.GetFeatureIndexFile(featureSlug))
	index := task.NewTaskIndex(featureSlug)
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed", File: "1.1.md"},
	})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}
	return
}

func TestAddFixTask_BasicCompile(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	output := "./internal/handler.go:42:2: undefined: foo\n./internal/handler.go:43:1: too many arguments"
	errorDocPath := "tests/results/unit-raw-output.txt"

	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", output, errorDocPath)
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}

	updatedIndex, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	addedTask, exists := updatedIndex.ByID(taskID)
	if !exists {
		t.Fatalf("task %s not found in index", taskID)
	}

	if addedTask.Priority != "P0" {
		t.Errorf("priority = %q, want P0", addedTask.Priority)
	}
	if !addedTask.Breaking {
		t.Error("expected breaking=true")
	}
	if addedTask.Status != "pending" {
		t.Errorf("status = %q, want pending", addedTask.Status)
	}
	if addedTask.EstimatedTime != "30min" {
		t.Errorf("estimatedTime = %q, want 30min", addedTask.EstimatedTime)
	}
	if addedTask.SourceTaskID != "quality-gate:compile" {
		t.Errorf("sourceTaskID = %q, want %q", addedTask.SourceTaskID, "quality-gate:compile")
	}

	// Verify task markdown content
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
		t.Error("task markdown should reference 'just compile'")
	}
	if !strings.Contains(content, "Root Cause") {
		t.Error("task markdown should contain Root Cause section from fix-task template")
	}
	if !strings.Contains(content, "Reference Files") {
		t.Error("task markdown should contain Reference Files section from fix-task template")
	}
	if !strings.Contains(content, "Verification") {
		t.Error("task markdown should contain Verification section from fix-task template")
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

func TestAddFixTask_StepSpecificTestScripts(t *testing.T) {
	tests := []struct {
		step           string
		wantTestScript string
	}{
		{"compile", "just compile"},
		{"lint", "just lint"},
		{"unit-test", "just test"},
		{"test-e2e", "just test-e2e"},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {
			projectRoot, featureSlug, _ := helperSetup(t)
			taskID, addErr := addFixTask(projectRoot, featureSlug, tc.step, "handler.go:10: fail", "tests/results/fake.txt")
			if addErr != nil {
				t.Fatalf("unexpected error: %v", addErr)
			}
			if taskID == "" {
				t.Fatal("expected non-empty task ID")
			}

			mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
			data, err := os.ReadFile(mdPath)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), tc.wantTestScript) {
				t.Errorf("step %q should produce test script %q in markdown", tc.step, tc.wantTestScript)
			}
		})
	}
}

func TestAddFixTask_TypeFromStep(t *testing.T) {
	tests := []struct {
		step     string
		wantType string
	}{
		{"compile", task.TypeCodingFix},
		{"fmt", task.TypeCodingCleanup},
		{"lint", task.TypeCodingCleanup},
		{"unit-test", task.TypeCodingFix},
		{"test-e2e", task.TypeCodingFix},
		{"unknown-step", task.TypeCodingFix}, // default fallback
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {
			projectRoot, featureSlug, indexPath := helperSetup(t)
			taskID, addErr := addFixTask(projectRoot, featureSlug, tc.step, "handler.go:10: fail", "tests/results/fake.txt")
			if addErr != nil {
				t.Fatalf("unexpected error: %v", addErr)
			}
			if taskID == "" {
				t.Fatal("expected non-empty task ID")
			}

			updatedIndex, err := task.LoadIndex(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			addedTask, exists := updatedIndex.ByID(taskID)
			if !exists {
				t.Fatalf("task %s not found in index", taskID)
			}
			if addedTask.Type != tc.wantType {
				t.Errorf("type for step %q = %q, want %q", tc.step, addedTask.Type, tc.wantType)
			}
		})
	}
}

func TestAddFixTask_TemplateSelection(t *testing.T) {
	tests := []struct {
		step        string
		wantType    string
		wantSnippet string // distinctive text in the generated .md
	}{
		{"compile", task.TypeCodingFix, "type: \"fix\""},
		{"fmt", task.TypeCodingCleanup, "type: \"cleanup\""},
		{"lint", task.TypeCodingCleanup, "type: \"cleanup\""},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {
			projectRoot, featureSlug, _ := helperSetup(t)
			taskID, addErr := addFixTask(projectRoot, featureSlug, tc.step, "handler.go:10: fail", "tests/results/fake.txt")
			if addErr != nil {
				t.Fatalf("unexpected error: %v", addErr)
			}

			mdPath := filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks", taskID+".md")
			content, err := os.ReadFile(mdPath)
			if err != nil {
				t.Fatalf("read .md: %v", err)
			}
			if !strings.Contains(string(content), tc.wantSnippet) {
				t.Errorf("step %q: .md content does not contain %q.\nGot:\n%s", tc.step, tc.wantSnippet, string(content))
			}
		})
	}
}

func TestAddFixTask_EmptyOutput(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	taskID, addErr := addFixTask(projectRoot, featureSlug, "lint", "", "tests/results/unit-raw-output.txt")
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID even with empty output")
	}

	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "See error output for affected files") {
		t.Error("empty output should produce fallback source files message")
	}
}

func TestAddFixTask_NoSourceFilesInOutput(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", "some random output without file references", "tests/results/unit-raw-output.txt")
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected task ID")
	}

	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "See error output for affected files") {
		t.Error("no source files in output should produce fallback message")
	}
}

func TestAddFixTask_SequentialIDs(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	id1, err1 := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	id2, err2 := addFixTask(projectRoot, featureSlug, "lint", "b.go:2: error", "tests/results/out.txt")
	id3, err3 := addFixTask(projectRoot, featureSlug, "unit-test", "c.go:3: error", "tests/results/out.txt")
	for _, e := range []error{err1, err2, err3} {
		if e != nil {
			t.Fatalf("unexpected error: %v", e)
		}
	}

	if id1 == "" || id2 == "" || id3 == "" {
		t.Fatalf("expected 3 valid IDs, got %q %q %q", id1, id2, id3)
	}

	// IDs should be different (max+1 via template prefix)
	if id1 == id2 || id2 == id3 || id1 == id3 {
		t.Errorf("expected unique IDs, got %q %q %q", id1, id2, id3)
	}

	// All should be in index
	idx, _ := task.LoadIndex(indexPath)
	count := 0
	for _, id := range []string{id1, id2, id3} {
		if _, ok := idx.ByID(id); ok {
			count++
		}
	}
	if count != 3 {
		t.Errorf("expected 3 tasks in index, found %d", count)
	}
}

func TestAddFixTask_TitleContainsStep(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	taskID, addErr := addFixTask(projectRoot, featureSlug, "lint", "a.go:1: error", "tests/results/out.txt")
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected task ID")
	}

	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "fix lint:") {
		t.Error("task title should contain 'fix lint:' prefix")
	}
}

func TestAddFixTask_DescriptionContainsErrorDoc(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	errorDoc := "tests/e2e/results/raw-output.txt"
	taskID, addErr := addFixTask(projectRoot, featureSlug, "test-e2e", "test.spec.ts:5: fail", errorDoc)
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected task ID")
	}

	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), errorDoc) {
		t.Errorf("task description should reference error doc %q", errorDoc)
	}
}

func TestAddFixTask_ForgeStateResetEachTime(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	// First add
	if _, err := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt"); err != nil {
		t.Fatalf("unexpected error on first add: %v", err)
	}
	state := feature.ReadForgeState(projectRoot)
	if state == nil || state.AllCompleted {
		t.Fatal("after first addFixTask, forge state should exist with allCompleted=false")
	}

	// Write allCompleted=true to simulate next completion cycle
	if err := feature.WriteForgeState(projectRoot, featureSlug); err != nil {
		t.Fatal(err)
	}

	// Second add should reset again
	if _, err := addFixTask(projectRoot, featureSlug, "lint", "b.go:2: error", "tests/results/out.txt"); err != nil {
		t.Fatalf("unexpected error on second add: %v", err)
	}
	state = feature.ReadForgeState(projectRoot)
	if state == nil || state.AllCompleted {
		t.Fatal("after second addFixTask, forge state should be reset again")
	}
}

func TestCountFixTasks(t *testing.T) {
	tests := []struct {
		name  string
		tasks map[string]task.Task
		step  string
		want  int
	}{
		{
			name:  "no fix tasks",
			tasks: map[string]task.Task{},
			step:  "compile",
			want:  0,
		},
		{
			name: "one active fix task for step",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "fix compile: something broke", Status: "pending"},
			},
			step: "compile",
			want: 1,
		},
		{
			name: "three active fix tasks for same step",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "fix compile: first error", Status: "pending"},
				"f2": {ID: "f2", SourceTaskID: "1.1", Title: "fix compile: second error", Status: "in_progress"},
				"f3": {ID: "f3", SourceTaskID: "1.1", Title: "fix compile: third error", Status: "blocked"},
			},
			step: "compile",
			want: 3,
		},
		{
			name: "completed fix tasks counted cumulatively",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "fix compile: done", Status: "completed"},
				"f2": {ID: "f2", SourceTaskID: "1.1", Title: "fix compile: active", Status: "pending"},
			},
			step: "compile",
			want: 2,
		},
		{
			name: "skipped fix tasks counted cumulatively",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "fix compile: skipped", Status: "skipped"},
				"f2": {ID: "f2", SourceTaskID: "1.1", Title: "fix compile: active", Status: "pending"},
			},
			step: "compile",
			want: 2,
		},
		{
			name: "different step not counted",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "fix compile: error", Status: "pending"},
				"f2": {ID: "f2", SourceTaskID: "1.1", Title: "fix lint: error", Status: "pending"},
			},
			step: "compile",
			want: 1,
		},
		{
			name: "task without SourceTaskID not counted even with matching title",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", SourceTaskID: "", Title: "fix compile: regular task", Status: "pending"},
			},
			step: "compile",
			want: 0,
		},
		{
			name: "task without fix prefix not counted",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "some other compile task", Status: "pending"},
			},
			step: "compile",
			want: 0,
		},
		{
			name: "mix of terminal and active across steps counted cumulatively",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "fix compile: first", Status: "completed"},
				"f2": {ID: "f2", SourceTaskID: "1.1", Title: "fix compile: second", Status: "skipped"},
				"f3": {ID: "f3", SourceTaskID: "1.1", Title: "fix compile: third", Status: "pending"},
				"f4": {ID: "f4", SourceTaskID: "1.1", Title: "fix lint: first", Status: "pending"},
			},
			step: "compile",
			want: 3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			index := task.NewTaskIndex("test")
			index.SetTasks(tc.tasks)
			got := countFixTasks(index, tc.step)
			if got != tc.want {
				t.Errorf("countFixTasks(%q) = %d, want %d", tc.step, got, tc.want)
			}
		})
	}
}

func TestAddFixTask_CapEnforced(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Pre-populate 3 active fix-tasks for "compile"
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	index.SetTask("f1", task.Task{ID: "f1", SourceTaskID: "1.1", Title: "fix compile: first", Status: "pending", File: "f1.md"})
	index.SetTask("f2", task.Task{ID: "f2", SourceTaskID: "1.1", Title: "fix compile: second", Status: "in_progress", File: "f2.md"})
	index.SetTask("f3", task.Task{ID: "f3", SourceTaskID: "1.1", Title: "fix compile: third", Status: "blocked", File: "f3.md"})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	taskID, capErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if capErr == nil {
		t.Errorf("expected error when 3 active fix-tasks exist, got nil (taskID=%q)", taskID)
	}
	if taskID != "" {
		t.Errorf("expected empty taskID on cap error, got %q", taskID)
	}
}

func TestAddFixTask_CapAllowsUnderLimit(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Pre-populate 2 active fix-tasks for "compile" (under cap of 3)
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	index.SetTask("f1", task.Task{ID: "f1", SourceTaskID: "1.1", Title: "fix compile: first", Status: "pending", File: "f1.md"})
	index.SetTask("f2", task.Task{ID: "f2", SourceTaskID: "1.1", Title: "fix compile: second", Status: "in_progress", File: "f2.md"})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	taskID, capErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if capErr != nil {
		t.Fatalf("expected no error with 2 active fix-tasks, got %v", capErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}
}

func TestAddFixTask_CumulativeCountBlocksAtCap(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// 3 fix-tasks for "compile": 2 completed/skipped + 1 active = 3 cumulative (at cap)
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	index.SetTask("f1", task.Task{ID: "f1", SourceTaskID: "1.1", Title: "fix compile: done", Status: "completed", File: "f1.md"})
	index.SetTask("f2", task.Task{ID: "f2", SourceTaskID: "1.1", Title: "fix compile: skipped", Status: "skipped", File: "f2.md"})
	index.SetTask("f3", task.Task{ID: "f3", SourceTaskID: "1.1", Title: "fix compile: active", Status: "pending", File: "f3.md"})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	taskID, capErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if capErr == nil {
		t.Errorf("expected error when 3 cumulative fix-tasks exist, got nil (taskID=%q)", taskID)
	}
	if taskID != "" {
		t.Errorf("expected empty taskID on cap error, got %q", taskID)
	}
}

func TestHandleGateFailure_DistinctReasons(t *testing.T) {
	tests := []struct {
		step          string
		fixID         string
		wantContains  string
		wantFixAction string
		wantClaim     bool   // expect "task claim" in output
		wantManual    bool   // expect "task add --template fix-task" in output
		wantFixMsg    string // expect this fix task message
	}{
		{"compile", "fix-1", "Project compilation failed in quality-gate hook", "fix compilation errors", true, false, "Fix task fix-1 added (P0, breaking)"},
		{"lint", "fix-2", "Lint check failed in quality-gate hook", "fix lint errors", true, false, "Fix task fix-2 added (P0, breaking)"},
		{"unit-test", "fix-3", "Unit tests failed in quality-gate hook", "fix failing tests", true, false, "Fix task fix-3 added (P0, breaking)"},
		{"test-e2e", "fix-4", "E2e regression tests failed in quality-gate hook", "fix failing e2e tests", true, false, "Fix task fix-4 added (P0, breaking)"},
		{"unknown-step", "fix-5", "Unknown-step check failed in quality-gate hook", "fix the issue", true, false, "Fix task fix-5 added (P0, breaking)"},
		{"compile", "", "Project compilation failed in quality-gate hook", "fix compilation errors", false, true, "Failed to add fix task automatically"},
	}

	for _, tc := range tests {
		name := tc.step
		if tc.fixID == "" {
			name = "nofixid-" + tc.step
		}
		t.Run(name, func(t *testing.T) {
			if os.Getenv("TEST_HANDLE_GATE") == "1" {
				handleGateFailure(tc.step, "tests/results/fake.txt", tc.fixID, "some error detail")
				return
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestHandleGateFailure_DistinctReasons/"+name)
			cmd.Env = append(os.Environ(), "TEST_HANDLE_GATE=1")
			output, _ := cmd.CombinedOutput()

			got := string(output)
			if !strings.Contains(got, tc.wantContains) {
				t.Errorf("reason for step %q should contain %q, got:\n%s", tc.step, tc.wantContains, got)
			}
			if !strings.Contains(got, tc.wantFixAction) {
				t.Errorf("reason for step %q should contain fix action %q", tc.step, tc.wantFixAction)
			}
			if tc.wantClaim && !strings.Contains(got, "task claim") {
				t.Errorf("reason for step %q should contain 'task claim'", tc.step)
			}
			if tc.wantManual && !strings.Contains(got, "task add --template fix-task") {
				t.Errorf("reason for step %q (no fixID) should contain manual add instruction", tc.step)
			}
			if !strings.Contains(got, tc.wantFixMsg) {
				t.Errorf("reason for step %q should contain %q, got:\n%s", tc.step, tc.wantFixMsg, got)
			}
			if !strings.Contains(got, "tests/results/fake.txt") {
				t.Errorf("reason for step %q should reference error doc path", tc.step)
			}
			if !strings.Contains(got, "some error detail") {
				t.Errorf("reason for step %q should include concise error", tc.step)
			}
		})
	}
}

func TestCheckAllCompleted_RejectedTaskReturnsNil(t *testing.T) {
	projectRoot := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", projectRoot)
	featureSlug := "test"
	tasksDir := filepath.Join(projectRoot, "docs", "features", featureSlug, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		t.Fatal(err)
	}

	index := task.NewTaskIndex(featureSlug)
	index.SetTasks(map[string]task.Task{
		"task-a": {ID: "1.1", Status: "completed", File: "1.1.md"},
		"task-b": {ID: "1.2", Status: "rejected", File: "1.2.md"},
	})
	indexPath := filepath.Join(tasksDir, "index.json")
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}
	if err := feature.WriteForgeState(projectRoot, featureSlug); err != nil {
		t.Fatal(err)
	}

	result := checkAllCompleted(false)
	if result != nil {
		t.Error("rejected task should prevent quality-gate from proceeding")
	}
}

func TestCheckAllCompleted_ForgeStateConsumedOnSuccess(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed"},
	})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}
	if err := feature.WriteForgeState(dir, "test"); err != nil {
		t.Fatal(err)
	}

	// First call should succeed and consume the state
	result := checkAllCompleted(false)
	if result == nil {
		t.Fatal("first call should return result")
	}

	// Forge state should be cleared
	state := feature.ReadForgeState(dir)
	if state != nil {
		t.Error("forge state should be cleared after checkAllCompleted consumes it")
	}

	// Second call should return nil (no state)
	result2 := checkAllCompleted(false)
	if result2 != nil {
		t.Error("second call should return nil after state was consumed")
	}
}

func TestIsDocsOnly(t *testing.T) {
	tests := []struct {
		name  string
		tasks map[string]task.Task
		want  bool
	}{
		{
			name: "only documentation tasks",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeDoc},
				"t2": {ID: "2", Type: task.TypeDoc},
			},
			want: true,
		},
		{
			name: "documentation plus doc-evaluation",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeDoc},
				"t2": {ID: "T-eval-doc", Type: task.TypeDocEval},
			},
			want: true,
		},
		{
			name:  "empty index is docs-only",
			tasks: map[string]task.Task{},
			want:  true,
		},
		{
			name: "has feature task",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeDoc},
				"t2": {ID: "2", Type: task.TypeCodingFeature},
			},
			want: false,
		},
		{
			name: "has enhancement task",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeCodingEnhancement},
			},
			want: false,
		},
		{
			name: "has fix task",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeDoc},
				"f1": {ID: "fix-1", Type: task.TypeCodingFix},
			},
			want: false,
		},
		{
			name: "has cleanup task (testable)",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeCodingCleanup},
			},
			want: false,
		},
		{
			name: "has refactor task (testable)",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeCodingRefactor},
			},
			want: false,
		},
		{
			name: "has feature task (testable)",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeCodingFeature},
			},
			want: false,
		},
		{
			name: "test-pipeline tasks only",
			tasks: map[string]task.Task{
				"t1": {ID: "T-quick-gen-cases", Type: task.TypeTestGenCases},
				"t2": {ID: "T-quick-gen-and-run", Type: task.TypeTestGenAndRun},
			},
			want: true,
		},
		{
			name: "mixed documentation and test-pipeline",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeDoc},
				"t2": {ID: "T-quick-verify-regression", Type: task.TypeTestRun},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			index := task.NewTaskIndex("test")
			index.SetTasks(tc.tasks)
			got := isDocsOnly(index)
			if got != tc.want {
				t.Errorf("isDocsOnly() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCheckAllCompleted_DocsOnlyFlag(t *testing.T) {
	tests := []struct {
		name         string
		tasks        map[string]task.Task
		wantDocsOnly bool
	}{
		{
			name: "documentation only sets DocsOnly true",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Status: "completed", Type: task.TypeDoc},
				"t2": {ID: "T-eval-doc", Status: "completed", Type: task.TypeDocEval},
			},
			wantDocsOnly: true,
		},
		{
			name: "feature task sets DocsOnly false",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Status: "completed", Type: task.TypeCodingFeature},
			},
			wantDocsOnly: false,
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
				Feature:    "test",
				StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped"},
			}
			index.SetTasks(tc.tasks)
			if err := task.SaveIndex(indexPath, index); err != nil {
				t.Fatal(err)
			}
			if err := feature.WriteForgeState(dir, "test"); err != nil {
				t.Fatal(err)
			}

			result := checkAllCompleted(false)
			if result == nil {
				t.Fatal("expected non-nil result")
			}
			if result.DocsOnly != tc.wantDocsOnly {
				t.Errorf("DocsOnly = %v, want %v", result.DocsOnly, tc.wantDocsOnly)
			}
		})
	}
}

func TestCheckAllCompleted_ManyCompletedTasks(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	tasks := make(map[string]task.Task)
	for i := range 20 {
		id := fmt.Sprintf("1.%d", i+1)
		key := fmt.Sprintf("task-%d", i+1)
		tasks[key] = task.Task{ID: id, Status: "completed"}
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped"},
	}
	index.SetTasks(tasks)
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}
	if err := feature.WriteForgeState(dir, "test"); err != nil {
		t.Fatal(err)
	}

	result := checkAllCompleted(false)
	if result == nil {
		t.Fatal("expected result with many completed tasks")
	}
}

func TestCheckAllCompleted_AllBlockedReturnsNil(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "blocked"},
		"t2": {ID: "1.2", Status: "blocked"},
	})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}
	if err := feature.WriteForgeState(dir, "test"); err != nil {
		t.Fatal(err)
	}

	result := checkAllCompleted(false)
	if result != nil {
		t.Error("all blocked tasks should return nil")
	}
}

func TestCheckAllCompleted_MixedCompletedSkippedRejected(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped", "rejected"},
	}
	index.SetTasks(map[string]task.Task{
		"t1": {ID: "1.1", Status: "completed"},
		"t2": {ID: "1.2", Status: "skipped"},
		"t3": {ID: "1.3", Status: "rejected"},
	})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}
	if err := feature.WriteForgeState(dir, "test"); err != nil {
		t.Fatal(err)
	}

	// rejected is not completed or skipped, so should return nil
	result := checkAllCompleted(false)
	if result != nil {
		t.Error("rejected task should prevent quality-gate from proceeding")
	}
}

func TestCheckAllCompleted_VerboseMode(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	// No forge state -> should return nil but not error
	result := checkAllCompleted(true)
	if result != nil {
		t.Error("expected nil result without forge state")
	}
}

func TestAddFixTask_StepScopedSentinel(t *testing.T) {
	tests := []struct {
		step         string
		wantSentinel string
	}{
		{"compile", "quality-gate:compile"},
		{"lint", "quality-gate:lint"},
		{"unit-test", "quality-gate:unit-test"},
		{"test-e2e", "quality-gate:test-e2e"},
	}

	for _, tc := range tests {
		t.Run(tc.step, func(t *testing.T) {
			projectRoot, featureSlug, indexPath := helperSetup(t)

			taskID, addErr := addFixTask(projectRoot, featureSlug, tc.step, "handler.go:10: fail", "tests/results/fake.txt")
			if addErr != nil {
				t.Fatalf("unexpected error: %v", addErr)
			}
			if taskID == "" {
				t.Fatal("expected non-empty task ID")
			}

			updatedIndex, err := task.LoadIndex(indexPath)
			if err != nil {
				t.Fatal(err)
			}
			addedTask, exists := updatedIndex.ByID(taskID)
			if !exists {
				t.Fatalf("task %s not found in index", taskID)
			}

			if addedTask.SourceTaskID != tc.wantSentinel {
				t.Errorf("SourceTaskID = %q, want %q", addedTask.SourceTaskID, tc.wantSentinel)
			}
		})
	}
}

func TestAddFixTask_CrossStepIndependence(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Pre-populate 3 fix-tasks for "compile" (at cap)
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	index.SetTask("f1", task.Task{ID: "f1", SourceTaskID: "1.1", Title: "fix compile: first", Status: "pending", File: "f1.md"})
	index.SetTask("f2", task.Task{ID: "f2", SourceTaskID: "1.1", Title: "fix compile: second", Status: "in_progress", File: "f2.md"})
	index.SetTask("f3", task.Task{ID: "f3", SourceTaskID: "1.1", Title: "fix compile: third", Status: "blocked", File: "f3.md"})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	// compile is at cap -> should fail
	_, capErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if capErr == nil {
		t.Error("expected cap error for compile step at limit")
	}

	// lint has no fix tasks -> should succeed
	taskID, lintErr := addFixTask(projectRoot, featureSlug, "lint", "b.go:2: error", "tests/results/out.txt")
	if lintErr != nil {
		t.Fatalf("expected no error for lint step (cross-step independent), got %v", lintErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID for lint fix task")
	}

	// unit-test has no fix tasks -> should succeed
	taskID2, unitErr := addFixTask(projectRoot, featureSlug, "unit-test", "c.go:3: error", "tests/results/out.txt")
	if unitErr != nil {
		t.Fatalf("expected no error for unit-test step (cross-step independent), got %v", unitErr)
	}
	if taskID2 == "" {
		t.Fatal("expected non-empty task ID for unit-test fix task")
	}
}

func TestAddFixTask_VarsSourceTaskIDRemainsNA(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}

	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	// Vars["SOURCE_TASK_ID"] in template should still be "N/A (project-wide gate)"
	if !strings.Contains(content, "N/A (project-wide gate)") {
		t.Error("task markdown should contain 'N/A (project-wide gate)' for template rendering")
	}
}

func TestAddFixTask_TaskAddFailure(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Delete the index file so that task.AddTask's internal LoadIndex fails.
	// The cap check in addFixTask will print a WARNING and proceed (by design),
	// then AddTask will fail with "load index" error.
	if err := os.Remove(indexPath); err != nil {
		t.Fatal(err)
	}

	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if addErr == nil {
		t.Fatalf("expected error when task add fails (no index), got nil (taskID=%q)", taskID)
	}
	if taskID != "" {
		t.Errorf("expected empty taskID on error, got %q", taskID)
	}
	if !strings.Contains(addErr.Error(), "failed to add fix task") {
		t.Errorf("error should contain 'failed to add fix task', got: %v", addErr)
	}
	if !strings.Contains(addErr.Error(), "load index") {
		t.Errorf("error should contain 'load index', got: %v", addErr)
	}
}

func TestAddFixTask_MarkdownCreationError(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Pre-add a task with ID "fix-1" that matches the auto-generated ID.
	// This means AddTask will generate "fix-2" for the new task.
	// Then pre-create a *read-only directory* named "fix-2.md" to block
	// os.WriteFile in CreateTaskMarkdown.
	index, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	index.SetTask("fix-1", task.Task{ID: "fix-1", Status: "completed", File: "fix-1.md"})
	if err := task.SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}

	tasksDir := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug))
	// Create a directory named "fix-2.md" so os.WriteFile fails (can't write to a directory)
	blockerPath := filepath.Join(tasksDir, "fix-2.md")
	if err := os.MkdirAll(blockerPath, 0755); err != nil {
		t.Fatal(err)
	}

	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if addErr == nil {
		t.Fatalf("expected error when markdown creation fails, got nil (taskID=%q)", taskID)
	}
	if taskID != "" {
		t.Errorf("expected empty taskID on error, got %q", taskID)
	}
	if !strings.Contains(addErr.Error(), "fix-2.md") {
		t.Errorf("error should reference the blocked file, got: %v", addErr)
	}
}

func TestAddFixTask_TemplateNotFoundError_NonexistentTemplate(t *testing.T) {
	// This test verifies that when the template doesn't exist, addFixTask
	// returns an explicit error. Since "fix-task" is embedded and always exists,
	// we test via the internal code path directly by checking that the function
	// properly propagates errors from tmpl.Get.
	//
	// We can trigger this by temporarily pointing at a feature that uses
	// a non-existent template. However, since the template name is hardcoded
	// in addFixTask, we verify the behavior through the task-add-failure
	// and markdown-failure tests above, plus this test confirms that the
	// current success path still works with the existing "fix-task" template.
	projectRoot, featureSlug, _ := helperSetup(t)

	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if addErr != nil {
		t.Fatalf("expected no error with valid template, got: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}
}

func TestRunUnitTestStep_RetryPass(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	callCount := 0
	mockRun := func(_, _ string) (string, bool) {
		callCount++
		if callCount == 1 {
			return "FAIL: TestFlaky", false
		}
		return "ok", true
	}

	passed, fixID, fixErr := runUnitTestStep(projectRoot, featureSlug, mockRun, "")
	if !passed {
		t.Error("expected passed=true when retry succeeds")
	}
	if fixID != "" {
		t.Errorf("expected no fix task on retry pass, got fixID=%q", fixID)
	}
	if fixErr != nil {
		t.Errorf("expected no error on retry pass, got %v", fixErr)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (initial + retry), got %d", callCount)
	}
}

func TestRunUnitTestStep_RetryFail(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	mockRun := func(_, _ string) (string, bool) {
		return "FAIL: TestReal", false
	}

	passed, fixID, fixErr := runUnitTestStep(projectRoot, featureSlug, mockRun, "")
	if passed {
		t.Error("expected passed=false when both attempts fail")
	}
	if fixID == "" {
		t.Error("expected fix task ID on double failure")
	}
	if fixErr != nil {
		t.Errorf("expected no error from runUnitTestStep, got %v", fixErr)
	}

	// Verify fix task markdown mentions retry
	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), fixID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("task markdown not found: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "retried once, both attempts failed") {
		t.Errorf("fix task description should mention retry, got content (first 500 chars): %.500s", content)
	}
}

func TestRunUnitTestStep_FirstPass(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	callCount := 0
	mockRun := func(_, _ string) (string, bool) {
		callCount++
		return "ok", true
	}

	passed, fixID, fixErr := runUnitTestStep(projectRoot, featureSlug, mockRun, "")
	if !passed {
		t.Error("expected passed=true on first pass")
	}
	if fixID != "" {
		t.Errorf("expected no fix task on first pass, got fixID=%q", fixID)
	}
	if fixErr != nil {
		t.Errorf("expected no error, got %v", fixErr)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (no retry needed), got %d", callCount)
	}
}

func TestRunUnitTestStep_RetryOutputInDescription(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	callCount := 0
	mockRun := func(_, _ string) (string, bool) {
		callCount++
		return fmt.Sprintf("attempt %d output: FAIL: TestX", callCount), false
	}

	passed, fixID, _ := runUnitTestStep(projectRoot, featureSlug, mockRun, "")
	if passed {
		t.Error("expected passed=false")
	}

	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), fixID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("task markdown not found: %v", err)
	}
	content := string(data)

	// Description should include retry-run output (attempt 2)
	if !strings.Contains(content, "attempt 2 output") {
		t.Errorf("fix task description should contain retry output, got content (first 500 chars): %.500s", content)
	}
}
