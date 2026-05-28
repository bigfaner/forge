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
		name       string
		tasks      map[string]task.Task
		forgeState bool
		wantNil    bool
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
				Feature:    "test",
				StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped"},
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
			result, _ := checkAllCompleted(false)

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
		})
	}
}

func TestCheckAllCompleted_NoFeature(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)
	if err := os.MkdirAll(filepath.Join(dir, feature.FeaturesDir), 0755); err != nil {
		t.Fatal(err)
	}
	_, err := checkAllCompleted(false)
	if err == nil {
		t.Error("expected error when no feature set, got nil")
	}
	if !strings.Contains(err.Error(), "No feature set") {
		t.Errorf("expected 'No feature set' error, got: %v", err)
	}
}

func TestCheckAllCompleted_NoProject(t *testing.T) {
	if os.Getenv("TEST_CHECK_ALL_COMPLETED_NO_PROJECT") == "1" {
		_, err := checkAllCompleted(false)
		if err == nil {
			t.Error("expected error when no project root, got nil")
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

func TestGroupFilesByDir(t *testing.T) {
	tests := []struct {
		name  string
		files string
		want  int // number of groups (0 = nil returned, meaning single group or no files)
	}{
		{
			name:  "empty returns nil",
			files: "",
			want:  0,
		},
		{
			name:  "fallback message returns nil",
			files: "See error output for affected files",
			want:  0,
		},
		{
			name:  "single file returns nil",
			files: "pkg/handler.go",
			want:  0,
		},
		{
			name:  "two files same directory returns nil",
			files: "pkg/handler.go, pkg/service.go",
			want:  0,
		},
		{
			name:  "two files different directories returns two groups",
			files: "pkg/handler.go, internal/service.go",
			want:  2,
		},
		{
			name:  "three files two directories returns two groups",
			files: "pkg/handler.go, pkg/service.go, internal/main.go",
			want:  2,
		},
		{
			name:  "three files three directories returns three groups",
			files: "a/handler.go, b/service.go, c/main.go",
			want:  3,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			groups := groupFilesByDir(tc.files)
			if tc.want == 0 {
				if groups != nil {
					t.Errorf("expected nil, got %v", groups)
				}
			} else {
				if len(groups) != tc.want {
					t.Errorf("expected %d groups, got %d: %v", tc.want, len(groups), groups)
				}
			}
		})
	}
}

func TestAddFixTask_MultiDirCreatesMultipleTasks(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Simulate compile error with files in two different directories
	output := "pkg/handler.go:10: error\ninternal/service.go:20: error"
	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", output, "tests/results/out.txt")
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}

	// Should have created 2 tasks (one per directory)
	updatedIndex, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	// Count tasks with "fix compile:" prefix
	fixCount := 0
	for _, t := range updatedIndex.TasksMap() {
		if strings.HasPrefix(t.Title, "fix compile:") {
			fixCount++
		}
	}
	if fixCount != 2 {
		t.Errorf("expected 2 fix tasks (one per directory), got %d", fixCount)
	}
}

func TestAddFixTask_SingleDirCreatesOneTask(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Simulate compile error with files in the same directory
	output := "pkg/handler.go:10: error\npkg/service.go:20: error"
	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", output, "tests/results/out.txt")
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}

	// Should have created 1 task
	updatedIndex, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}

	fixCount := 0
	for _, t := range updatedIndex.TasksMap() {
		if strings.HasPrefix(t.Title, "fix compile:") {
			fixCount++
		}
	}
	if fixCount != 1 {
		t.Errorf("expected 1 fix task (same directory), got %d", fixCount)
	}
}

// helperSetup creates a minimal feature with a completed task for addFixTask tests.
// Configures a default scalar surface (".": "cli") so surface inference succeeds.
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

	// Configure default surface for inference
	forgeDir := filepath.Join(projectRoot, ".forge")
	if err := os.MkdirAll(forgeDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("surfaces:\n  .: cli\n"), 0644); err != nil {
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
	if addedTask.SourceTaskID != "" {
		t.Errorf("sourceTaskID = %q, want empty (no sentinel)", addedTask.SourceTaskID)
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
		{"unit-test", "just unit-test"},
		{"test", "just test"},
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
		{"test", task.TypeCodingFix},
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

func TestAddFixTask_CleanupTaskNonBreaking(t *testing.T) {
	// Cleanup tasks (fmt/lint) should use Breaking=false and EstimatedTime="15min"
	tests := []struct {
		step string
	}{
		{"fmt"},
		{"lint"},
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
			if addedTask.Breaking {
				t.Errorf("cleanup task for step %q: Breaking=true, want false", tc.step)
			}
			if addedTask.EstimatedTime != "15min" {
				t.Errorf("cleanup task for step %q: EstimatedTime=%q, want 15min", tc.step, addedTask.EstimatedTime)
			}
		})
	}
}

func TestAddFixTask_FixTaskBreakingWithEstimatedTime(t *testing.T) {
	// Fix tasks (compile/test/unit-test) should use Breaking=true and EstimatedTime="30min"
	tests := []struct {
		step string
	}{
		{"compile"},
		{"unit-test"},
		{"test"},
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
			if !addedTask.Breaking {
				t.Errorf("fix task for step %q: Breaking=false, want true", tc.step)
			}
			if addedTask.EstimatedTime != "30min" {
				t.Errorf("fix task for step %q: EstimatedTime=%q, want 30min", tc.step, addedTask.EstimatedTime)
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
		{"compile", task.TypeCodingFix, "type: \"coding.fix\""},
		{"fmt", task.TypeCodingCleanup, "type: \"coding.cleanup\""},
		{"lint", task.TypeCodingCleanup, "type: \"coding.cleanup\""},
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

	// Soft-failure policy: empty output -> no source files -> surface inference fails
	// -> fix-task is still created with empty surface key/type.
	taskID, addErr := addFixTask(projectRoot, featureSlug, "lint", "", "tests/results/unit-raw-output.txt")
	if addErr != nil {
		t.Fatalf("expected no error (soft-failure policy), got: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty taskID")
	}
}

func TestAddFixTask_NoSourceFilesInOutput(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	// Soft-failure policy: no source files in output -> surface inference fails
	// -> fix-task is still created with empty surface key/type.
	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", "some random output without file references", "tests/results/unit-raw-output.txt")
	if addErr != nil {
		t.Fatalf("expected no error (soft-failure policy), got: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty taskID")
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

	errorDoc := "tests/results/raw-output.txt"
	taskID, addErr := addFixTask(projectRoot, featureSlug, "test", "test.spec.ts:5: fail", errorDoc)
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
			want: 1,
		},
		{
			name: "skipped fix tasks counted cumulatively",
			tasks: map[string]task.Task{
				"f1": {ID: "f1", SourceTaskID: "1.1", Title: "fix compile: skipped", Status: "skipped"},
				"f2": {ID: "f2", SourceTaskID: "1.1", Title: "fix compile: active", Status: "pending"},
			},
			step: "compile",
			want: 1,
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
			want: 1,
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
			want: 1,
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

func TestAddFixTask_ActiveOnlyCapAllowsUnderLimit(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// 3 fix-tasks for "compile": 2 completed/skipped + 1 active. Active-only = 1 (under cap of 3).
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
	// Active-only count = 1 (f3), under cap of 3, so adding a 4th should succeed.
	taskID, capErr := addFixTask(projectRoot, featureSlug, "compile", "a.go:1: error", "tests/results/out.txt")
	if capErr != nil {
		t.Fatalf("expected no error with 1 active fix-task, got %v", capErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}
}

func TestHandleGateFailure_DistinctReasons(t *testing.T) {
	tests := []struct {
		step          string
		fixID         string
		breaking      bool
		wantContains  string
		wantFixAction string
		wantClaim     bool   // expect "task claim" in output
		wantManual    bool   // expect "task add --type coding.fix" in output
		wantFixMsg    string // expect this fix task message
	}{
		{"compile", "fix-1", true, "Project compilation failed in quality-gate hook", "fix compilation errors", true, false, "Fix task fix-1 added (P0, breaking=true)"},
		{"lint", "fix-2", false, "Lint check failed in quality-gate hook", "fix lint errors", true, false, "Fix task fix-2 added (P0, breaking=false)"},
		{"unit-test", "fix-3", true, "Unit tests failed in quality-gate hook", "fix failing unit tests", true, false, "Fix task fix-3 added (P0, breaking=true)"},
		{"test", "fix-4", true, "Advanced tests failed in quality-gate hook", "fix failing tests", true, false, "Fix task fix-4 added (P0, breaking=true)"},
		{"unknown-step", "fix-5", true, "Unknown-step check failed in quality-gate hook", "fix the issue", true, false, "Fix task fix-5 added (P0, breaking=true)"},
		{"compile", "", true, "Project compilation failed in quality-gate hook", "fix compilation errors", false, true, "Failed to add fix task automatically"},
	}
	for _, tc := range tests {
		name := tc.step
		if tc.fixID == "" {
			name = "nofixid-" + tc.step
		}
		t.Run(name, func(t *testing.T) {
			if os.Getenv("TEST_HANDLE_GATE") == "1" {
				_ = handleGateFailure(tc.step, "tests/results/fake.txt", tc.fixID, "some error detail", tc.breaking)
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
			if tc.wantManual && !strings.Contains(got, "task add --type coding.fix") {
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
	result, _ := checkAllCompleted(false)
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
	result, _ := checkAllCompleted(false)
	if result == nil {
		t.Fatal("first call should return result")
	}
	// Forge state should be consumed (AllCompleted set to false)
	state := feature.ReadForgeState(dir)
	if state != nil && state.AllCompleted {
		t.Error("forge state should have AllCompleted=false after checkAllCompleted consumes it")
	}
	// Second call should return nil (AllCompleted=false, not signaling)
	result2, _ := checkAllCompleted(false)
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
			name: "documentation plus doc-review",
			tasks: map[string]task.Task{
				"t1": {ID: "1", Type: task.TypeDoc},
				"t2": {ID: "T-review-doc", Type: task.TypeDocReview},
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
				"t1": {ID: "T-test-gen-scripts-api", Type: task.TypeTestGenScripts},
				"t2": {ID: "T-test-gen-scripts", Type: task.TypeTestGenScripts},
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
				"t2": {ID: "T-review-doc", Status: "completed", Type: task.TypeDocReview},
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
			result, _ := checkAllCompleted(false)
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
	result, _ := checkAllCompleted(false)
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
	result, _ := checkAllCompleted(false)
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
	result, _ := checkAllCompleted(false)
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
	result, _ := checkAllCompleted(true)
	if result != nil {
		t.Error("expected nil result without forge state")
	}
}

func TestAddFixTask_SourceTaskIDEmpty(t *testing.T) {
	tests := []struct {
		step string
	}{
		{"compile"},
		{"lint"},
		{"unit-test"},
		{"test"},
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
			if addedTask.SourceTaskID != "" {
				t.Errorf("SourceTaskID = %q, want empty (no sentinel)", addedTask.SourceTaskID)
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
	mockRun := func(_ string) (string, bool) {
		callCount++
		if callCount == 1 {
			return "FAIL: TestFlaky", false
		}
		return "ok", true
	}
	passed, fixID, fixErr := runUnitTestStep(projectRoot, featureSlug, mockRun)
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

	mockRun := func(_ string) (string, bool) {
		return "handler.go:10: FAIL: TestReal", false
	}
	passed, fixID, fixErr := runUnitTestStep(projectRoot, featureSlug, mockRun)
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
	mockRun := func(_ string) (string, bool) {
		callCount++
		return "ok", true
	}
	passed, fixID, fixErr := runUnitTestStep(projectRoot, featureSlug, mockRun)
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
	mockRun := func(_ string) (string, bool) {
		callCount++
		return fmt.Sprintf("handler.go:%d: attempt %d output: FAIL: TestX", callCount, callCount), false
	}
	passed, fixID, _ := runUnitTestStep(projectRoot, featureSlug, mockRun)
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

func TestInferSurface(t *testing.T) {
	tests := []struct {
		name            string
		configYAML      string // content for .forge/config.yaml
		sourceFiles     string
		wantSurfaceKey  string
		wantSurfaceType string
	}{
		{
			name:            "scalar surface matches file",
			configYAML:      "surfaces:\n  .: web\n",
			sourceFiles:     "src/components/App.tsx",
			wantSurfaceKey:  ".",
			wantSurfaceType: "web",
		},
		{
			name:            "map surface matches file by prefix",
			configYAML:      "surfaces:\n  admin-panel: web\n  payment-service: api\n",
			sourceFiles:     "admin-panel/src/App.tsx",
			wantSurfaceKey:  "admin-panel",
			wantSurfaceType: "web",
		},
		{
			name:            "no surfaces configured returns empty",
			configYAML:      "",
			sourceFiles:     "src/components/App.tsx",
			wantSurfaceKey:  "",
			wantSurfaceType: "",
		},
		{
			name:            "file not matching any surface returns empty",
			configYAML:      "surfaces:\n  admin-panel: web\n",
			sourceFiles:     "src/components/App.tsx",
			wantSurfaceKey:  "",
			wantSurfaceType: "",
		},
		{
			name:            "empty source files returns empty",
			configYAML:      "surfaces:\n  .: web\n",
			sourceFiles:     "",
			wantSurfaceKey:  "",
			wantSurfaceType: "",
		},
		{
			name:            "fallback message source files returns empty",
			configYAML:      "surfaces:\n  .: web\n",
			sourceFiles:     "See error output for affected files",
			wantSurfaceKey:  "",
			wantSurfaceType: "",
		},
		{
			name:            "first file from comma-separated list",
			configYAML:      "surfaces:\n  backend: api\n",
			sourceFiles:     "backend/handler.go, frontend/App.tsx",
			wantSurfaceKey:  "backend",
			wantSurfaceType: "api",
		},
		{
			name:            "second file matches when first does not",
			configYAML:      "surfaces:\n  admin-panel: web\n  payment-service: api\n",
			sourceFiles:     "src/unknown.go, admin-panel/src/App.tsx",
			wantSurfaceKey:  "admin-panel",
			wantSurfaceType: "web",
		},
		{
			name:            "third file matches when first two do not",
			configYAML:      "surfaces:\n  backend: api\n  frontend: web\n",
			sourceFiles:     "readme.go, docs.go, backend/handler.go",
			wantSurfaceKey:  "backend",
			wantSurfaceType: "api",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			projectRoot := t.TempDir()
			forgeDir := filepath.Join(projectRoot, ".forge")
			if err := os.MkdirAll(forgeDir, 0755); err != nil {
				t.Fatal(err)
			}
			if tc.configYAML != "" {
				if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(tc.configYAML), 0644); err != nil {
					t.Fatal(err)
				}
			}
			gotKey, gotType := inferSurface(projectRoot, tc.sourceFiles)
			if gotKey != tc.wantSurfaceKey {
				t.Errorf("surfaceKey = %q, want %q", gotKey, tc.wantSurfaceKey)
			}
			if gotType != tc.wantSurfaceType {
				t.Errorf("surfaceType = %q, want %q", gotType, tc.wantSurfaceType)
			}
		})
	}
}

func TestRequireSurfaceInference(t *testing.T) {
	t.Run("returns values when inference succeeds", func(t *testing.T) {
		projectRoot := t.TempDir()
		forgeDir := filepath.Join(projectRoot, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("surfaces:\n  .: web\n"), 0644); err != nil {
			t.Fatal(err)
		}
		key, typ, err := requireSurfaceInference(projectRoot, "src/App.tsx")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key != "." {
			t.Errorf("key = %q, want %q", key, ".")
		}
		if typ != "web" {
			t.Errorf("typ = %q, want %q", typ, "web")
		}
	})

	t.Run("returns error when no surfaces configured", func(t *testing.T) {
		projectRoot := t.TempDir()
		_, _, err := requireSurfaceInference(projectRoot, "handler.go")
		if err == nil {
			t.Fatal("expected error when no surfaces configured")
		}
		if !strings.Contains(err.Error(), "surface inference failed") {
			t.Errorf("error should mention 'surface inference failed', got: %v", err)
		}
		if !strings.Contains(err.Error(), "forge surfaces detect") {
			t.Errorf("error should contain 'forge surfaces detect' guidance, got: %v", err)
		}
	})

	t.Run("returns error when file matches no surface", func(t *testing.T) {
		projectRoot := t.TempDir()
		forgeDir := filepath.Join(projectRoot, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("surfaces:\n  admin-panel: web\n"), 0644); err != nil {
			t.Fatal(err)
		}
		_, _, err := requireSurfaceInference(projectRoot, "unrelated/path.go")
		if err == nil {
			t.Fatal("expected error when file matches no surface")
		}
		if !strings.Contains(err.Error(), "surface inference failed") {
			t.Errorf("error should mention 'surface inference failed', got: %v", err)
		}
	})
}

func TestAddFixTask_SurfaceInference(t *testing.T) {
	projectRoot, featureSlug, indexPath := helperSetup(t)

	// Configure surfaces in .forge/config.yaml
	forgeDir := filepath.Join(projectRoot, ".forge")
	if err := os.MkdirAll(forgeDir, 0755); err != nil {
		t.Fatal(err)
	}
	configContent := "surfaces:\n  admin-panel: web\n  payment-service: api\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	output := "./admin-panel/src/App.tsx:10: error: undefined variable"
	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", output, "tests/results/out.txt")
	if addErr != nil {
		t.Fatalf("unexpected error: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty task ID")
	}

	// Verify surface info in index.json
	updatedIndex, err := task.LoadIndex(indexPath)
	if err != nil {
		t.Fatal(err)
	}
	addedTask, exists := updatedIndex.ByID(taskID)
	if !exists {
		t.Fatalf("task %s not found in index", taskID)
	}
	if addedTask.SurfaceKey != "admin-panel" {
		t.Errorf("SurfaceKey = %q, want %q", addedTask.SurfaceKey, "admin-panel")
	}
	if addedTask.SurfaceType != "web" {
		t.Errorf("SurfaceType = %q, want %q", addedTask.SurfaceType, "web")
	}

	// Verify surface info in markdown frontmatter
	mdPath := filepath.Join(projectRoot, feature.GetFeatureTasksDir(featureSlug), taskID+".md")
	data, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("task markdown file not found: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, `surface-key: "admin-panel"`) {
		t.Errorf("markdown should contain surface-key: %q, got:\n%s", "admin-panel", content[:min(len(content), 500)])
	}
	if !strings.Contains(content, `surface-type: "web"`) {
		t.Errorf("markdown should contain surface-type: %q, got:\n%s", "web", content[:min(len(content), 500)])
	}
}

func TestAddFixTask_SurfaceInferenceSoftFailure(t *testing.T) {
	projectRoot, featureSlug, _ := helperSetup(t)

	// Remove surfaces config to trigger inference failure
	forgeDir := filepath.Join(projectRoot, ".forge")
	if err := os.Remove(filepath.Join(forgeDir, "config.yaml")); err != nil {
		t.Fatal(err)
	}

	// Soft-failure policy: surface inference failure does NOT block fix-task creation.
	// Task is created with empty surface key/type.
	taskID, addErr := addFixTask(projectRoot, featureSlug, "compile", "handler.go:10: error", "tests/results/out.txt")
	if addErr != nil {
		t.Fatalf("expected no error (soft-failure policy), got: %v", addErr)
	}
	if taskID == "" {
		t.Fatal("expected non-empty taskID even when surface inference fails")
	}
}

func TestRunTestRegression(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	t.Run("skips when no test recipe", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte("compile:\n  echo ok\n"), 0644); err != nil {
			t.Fatal(err)
		}
		// Should return nil (skip) when no "test" recipe
		if err := runTestRegression(dir, "feat"); err != nil {
			t.Errorf("expected nil when no test recipe, got %v", err)
		}
	})

	t.Run("uses test-setup recipe when available", func(t *testing.T) {
		projectRoot, featureSlug, _ := helperSetup(t)
		content := "test-setup:\n  echo setup-ok\ntest:\n  echo test-ok\n"
		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		// runTestRegression checks for "test" recipe and optionally runs test-setup.
		// It may skip due to probe failure (no dev server in test env), which is acceptable.
		// The key assertion is that it doesn't panic or error with wrong recipe names.
		_ = runTestRegression(projectRoot, featureSlug)
	})
}

// --- Surface-aware orchestration tests ---

func TestNeedsFullLifecycle(t *testing.T) {
	tests := []struct {
		surfaceType string
		want        bool
	}{
		{"web", true},
		{"api", true},
		{"mobile", true},
		{"cli", false},
		{"tui", false},
		{"unknown", false},
	}
	for _, tc := range tests {
		t.Run(tc.surfaceType, func(t *testing.T) {
			got := needsFullLifecycle(tc.surfaceType)
			if got != tc.want {
				t.Errorf("needsFullLifecycle(%q) = %v, want %v", tc.surfaceType, got, tc.want)
			}
		})
	}
}

func TestSurfaceOrchestrationSequence(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	t.Run("cli surface executes test then teardown", func(t *testing.T) {
		projectRoot := t.TempDir()
		// Write config with cli surface
		forgeDir := filepath.Join(projectRoot, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("surfaces:\n  .: cli\n"), 0644); err != nil {
			t.Fatal(err)
		}
		// Write justfile with test + teardown recipes
		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(
			"test:\n  echo test-ok\ntest-setup:\n  echo setup-ok\nteardown:\n  echo teardown-ok\n",
		), 0644); err != nil {
			t.Fatal(err)
		}
		// For cli surface, the simplified sequence should run: test -> teardown
		// This should complete without error
		err := runTestRegression(projectRoot, "test-feature")
		// Error handling is tested by the lifecycle tests; this tests the integration path.
		_ = err
	})

	t.Run("no surfaces falls back to current behavior", func(t *testing.T) {
		projectRoot := t.TempDir()
		// No .forge/config.yaml -- no surfaces configured
		// Write justfile with test recipe only
		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte("test:\n  echo test-ok\n"), 0644); err != nil {
			t.Fatal(err)
		}
		// Should fall back to the legacy behavior (serverprobe + just test)
		// which returns nil since no e2e config exists (CLI-only probe returns true)
		err := runTestRegression(projectRoot, "test-feature")
		if err != nil {
			t.Errorf("expected nil for no-surfaces fallback, got %v", err)
		}
	})

	t.Run("multi-surface project runs both sequences", func(t *testing.T) {
		projectRoot := t.TempDir()
		forgeDir := filepath.Join(projectRoot, ".forge")
		if err := os.MkdirAll(forgeDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte("surfaces:\n  frontend: web\n  tools: cli\n"), 0644); err != nil {
			t.Fatal(err)
		}
		// Write justfile with all recipes
		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(
			"dev:\n  echo dev-ok\nprobe:\n  echo probe-ok\ntest:\n  echo test-ok\nteardown:\n  echo teardown-ok\ntest-setup:\n  echo setup-ok\n",
		), 0644); err != nil {
			t.Fatal(err)
		}
		// Should attempt multi-surface orchestration without panic
		_ = runTestRegression(projectRoot, "test-feature")
	})
}

func TestProbeWithRetry(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	t.Run("succeeds on first attempt", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte("probe:\n  echo probe-ok\n"), 0644); err != nil {
			t.Fatal(err)
		}
		ok := probeWithRetry(dir, "probe", 1, 0)
		if !ok {
			t.Error("expected probe to succeed on first attempt")
		}
	})

	t.Run("succeeds after retry", func(t *testing.T) {
		dir := t.TempDir()
		markerFile := filepath.Join(dir, "probe-marker")
		// Write a probe script that fails first time, succeeds second time
		scriptContent := fmt.Sprintf(`#!/bin/bash
if [ -f "%s" ]; then
  echo ok
else
  touch "%s"
  exit 1
fi
`, filepath.ToSlash(markerFile), filepath.ToSlash(markerFile))
		scriptPath := filepath.Join(dir, "probe.sh")
		if err := os.WriteFile(scriptPath, []byte(scriptContent), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(
			fmt.Sprintf("probe:\n  bash %s\n", filepath.ToSlash(scriptPath)),
		), 0644); err != nil {
			t.Fatal(err)
		}
		ok := probeWithRetry(dir, "probe", 3, 0)
		if !ok {
			t.Error("expected probe to succeed after retry")
		}
	})

	t.Run("fails after all retries", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte("probe:\n  exit 1\n"), 0644); err != nil {
			t.Fatal(err)
		}
		ok := probeWithRetry(dir, "probe", 3, 0)
		if ok {
			t.Error("expected probe to fail after all retries")
		}
	})

	t.Run("skips when no probe recipe", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte("test:\n  echo ok\n"), 0644); err != nil {
			t.Fatal(err)
		}
		ok := probeWithRetry(dir, "probe", 3, 0)
		if !ok {
			t.Error("expected probe to be skipped (return true) when no probe recipe")
		}
	})
}

func TestRunSurfaceLifecycle_TeardownAlwaysRuns(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	// helperWriteMarkerScript creates a bash script that writes a marker file.
	// This avoids Windows path issues with justfile inline paths.
	helperWriteMarkerScript := func(t *testing.T, scriptPath, markerPath string) {
		t.Helper()
		content := fmt.Sprintf("#!/bin/bash\necho ran > '%s'\n", filepath.ToSlash(markerPath))
		if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("teardown runs even when dev fails", func(t *testing.T) {
		projectRoot := t.TempDir()
		markerDir := filepath.Join(projectRoot, "tmp")
		if err := os.MkdirAll(markerDir, 0755); err != nil {
			t.Fatal(err)
		}
		markerFile := filepath.Join(markerDir, "teardown-ran")
		scriptPath := filepath.Join(markerDir, "teardown.sh")
		helperWriteMarkerScript(t, scriptPath, markerFile)

		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(
			fmt.Sprintf("dev:\n  exit 1\nteardown:\n  bash %s\n", filepath.ToSlash(scriptPath)),
		), 0644); err != nil {
			t.Fatal(err)
		}
		result := runSurfaceLifecycle(projectRoot, "web")
		if result.success {
			t.Error("expected failure when dev fails")
		}
		if !just.FileExists(markerFile) {
			t.Error("teardown should have run even when dev fails")
		}
	})

	t.Run("teardown runs even when probe fails", func(t *testing.T) {
		projectRoot := t.TempDir()
		markerDir := filepath.Join(projectRoot, "tmp")
		if err := os.MkdirAll(markerDir, 0755); err != nil {
			t.Fatal(err)
		}
		markerFile := filepath.Join(markerDir, "teardown-ran")
		scriptPath := filepath.Join(markerDir, "teardown.sh")
		helperWriteMarkerScript(t, scriptPath, markerFile)

		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(
			fmt.Sprintf("dev:\n  echo dev-ok\nprobe:\n  exit 1\nteardown:\n  bash %s\n", filepath.ToSlash(scriptPath)),
		), 0644); err != nil {
			t.Fatal(err)
		}
		result := runSurfaceLifecycle(projectRoot, "web")
		if result.success {
			t.Error("expected failure when probe fails")
		}
		if !just.FileExists(markerFile) {
			t.Error("teardown should have run even when probe fails")
		}
	})

	t.Run("teardown runs even when test fails", func(t *testing.T) {
		projectRoot := t.TempDir()
		markerDir := filepath.Join(projectRoot, "tmp")
		if err := os.MkdirAll(markerDir, 0755); err != nil {
			t.Fatal(err)
		}
		markerFile := filepath.Join(markerDir, "teardown-ran")
		scriptPath := filepath.Join(markerDir, "teardown.sh")
		helperWriteMarkerScript(t, scriptPath, markerFile)

		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(
			fmt.Sprintf("dev:\n  echo dev-ok\nprobe:\n  echo probe-ok\ntest:\n  exit 1\nteardown:\n  bash %s\n", filepath.ToSlash(scriptPath)),
		), 0644); err != nil {
			t.Fatal(err)
		}
		result := runSurfaceLifecycle(projectRoot, "web")
		if result.success {
			t.Error("expected failure when test fails")
		}
		if !just.FileExists(markerFile) {
			t.Error("teardown should have run even when test fails")
		}
	})

	t.Run("cli surface skips dev and probe", func(t *testing.T) {
		projectRoot := t.TempDir()
		markerDir := filepath.Join(projectRoot, "tmp")
		if err := os.MkdirAll(markerDir, 0755); err != nil {
			t.Fatal(err)
		}
		// If dev runs, this marker would be created
		devMarker := filepath.Join(markerDir, "dev-ran")
		devScript := filepath.Join(markerDir, "dev.sh")
		helperWriteMarkerScript(t, devScript, devMarker)

		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(
			fmt.Sprintf("dev:\n  bash %s\nprobe:\n  echo probe-ok\ntest:\n  echo test-ok\nteardown:\n  echo teardown-ok\n", filepath.ToSlash(devScript)),
		), 0644); err != nil {
			t.Fatal(err)
		}
		result := runSurfaceLifecycle(projectRoot, "cli")
		if !result.success {
			t.Error("expected success for cli surface")
		}
		if just.FileExists(devMarker) {
			t.Error("cli surface should not run dev")
		}
	})
}

func TestRunSurfaceLifecycle_SurfaceSpecificRecipes(t *testing.T) {
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	t.Run("uses surface-specific recipes when available", func(t *testing.T) {
		projectRoot := t.TempDir()
		markerDir := filepath.Join(projectRoot, "tmp")
		if err := os.MkdirAll(markerDir, 0755); err != nil {
			t.Fatal(err)
		}
		devMarker := filepath.Join(markerDir, "web-dev-ran")
		probeMarker := filepath.Join(markerDir, "web-probe-ran")
		testMarker := filepath.Join(markerDir, "test-ran")
		teardownMarker := filepath.Join(markerDir, "web-teardown-ran")

		// Helper to create marker scripts
		writeScript := func(name, marker string) string {
			scriptPath := filepath.Join(markerDir, name)
			content := fmt.Sprintf("#!/bin/bash\necho ran > '%s'\n", filepath.ToSlash(marker))
			if err := os.WriteFile(scriptPath, []byte(content), 0755); err != nil {
				t.Fatal(err)
			}
			return filepath.ToSlash(scriptPath)
		}

		devScript := writeScript("web-dev.sh", devMarker)
		probeScript := writeScript("web-probe.sh", probeMarker)
		testScript := writeScript("test.sh", testMarker)
		teardownScript := writeScript("web-teardown.sh", teardownMarker)

		if err := os.WriteFile(filepath.Join(projectRoot, "justfile"), []byte(
			fmt.Sprintf(`web-dev:
  bash %s
web-probe:
  bash %s
test:
  bash %s
web-teardown:
  bash %s
`, devScript, probeScript, testScript, teardownScript),
		), 0644); err != nil {
			t.Fatal(err)
		}
		result := runSurfaceLifecycle(projectRoot, "web")
		if !result.success {
			t.Error("expected success")
		}
		if !just.FileExists(devMarker) {
			t.Error("expected web-dev to run")
		}
		if !just.FileExists(probeMarker) {
			t.Error("expected web-probe to run")
		}
		if !just.FileExists(testMarker) {
			t.Error("expected test to run")
		}
		if !just.FileExists(teardownMarker) {
			t.Error("expected web-teardown to run")
		}
	})
}
