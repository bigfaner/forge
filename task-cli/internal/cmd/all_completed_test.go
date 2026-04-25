package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"task-cli/pkg/feature"
	"task-cli/pkg/task"
)

func TestCheckAllCompleted(t *testing.T) {
	tests := []struct {
		name         string
		tasks        map[string]task.Task
		testCommand  string
		createE2EDir bool
		wantNil      bool
		wantE2EDir   bool
		wantTestCmd  string
	}{
		{
			name: "all completed returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "completed"},
			},
			wantNil: false,
		},
		{
			name: "all skipped returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "skipped"},
			},
			wantNil: false,
		},
		{
			name: "mixed completed and skipped returns result",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "skipped"},
			},
			wantNil: false,
		},
		{
			name: "one pending task returns nil",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
				"t2": {ID: "1.2", Status: "pending"},
			},
			wantNil: true,
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
			name:    "empty task list returns result (vacuously all done)",
			tasks:   map[string]task.Task{},
			wantNil: false,
		},
		{
			name: "e2e scripts dir present is reported",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			createE2EDir: true,
			wantNil:      false,
			wantE2EDir:   true,
		},
		{
			name: "e2e scripts dir absent gives empty field",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			createE2EDir: false,
			wantNil:      false,
			wantE2EDir:   false,
		},
		{
			name: "testCommand from index.json is propagated",
			tasks: map[string]task.Task{
				"t1": {ID: "1.1", Status: "completed"},
			},
			testCommand: "make test",
			wantNil:     false,
			wantTestCmd: "make test",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("CLAUDE_PROJECT_DIR", dir)

			// Create feature directory structure
			if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
				t.Fatal(err)
			}

			// Write index.json
			indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
			index := &task.TaskIndex{
				Feature:     "test",
				StatusEnum:  []string{"pending", "in_progress", "completed", "blocked", "skipped"},
				Tasks:       tc.tasks,
				TestCommand: tc.testCommand,
			}
			if err := task.SaveIndex(indexPath, index); err != nil {
				t.Fatal(err)
			}

			// Optionally create e2e scripts dir
			if tc.createE2EDir {
				e2eDir := filepath.Join(dir, feature.GetFeatureTestingScriptsDir("test"))
				if err := os.MkdirAll(e2eDir, 0755); err != nil {
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

			if tc.wantE2EDir && result.E2EScriptsDir == "" {
				t.Error("expected E2EScriptsDir to be set")
			}
			if !tc.wantE2EDir && result.E2EScriptsDir != "" {
				t.Errorf("expected E2EScriptsDir to be empty, got %q", result.E2EScriptsDir)
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
	// Create features dir but no feature subdirectory
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
		name     string
		files    []string // files to create in temp dir
		want     bool
	}{
		{
			name:  "no justfile",
			files: []string{},
			want:  false,
		},
		{
			name:  "lowercase justfile",
			files: []string{"justfile"},
			want:  true,
		},
		{
			name:  "capitalized Justfile",
			files: []string{"Justfile"},
			want:  true,
		},
		{
			name:  "both present",
			files: []string{"justfile", "Justfile"},
			want:  true,
		},
		{
			name:  "unrelated files only",
			files: []string{"Makefile", "go.mod"},
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, f := range tc.files {
				if err := os.WriteFile(filepath.Join(dir, f), []byte("test:\n    echo ok\n"), 0644); err != nil {
					t.Fatal(err)
				}
			}
			if got := hasJustfile(dir); got != tc.want {
				t.Errorf("hasJustfile() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestHasJustRecipe(t *testing.T) {
	// Skip if just is not installed
	if _, err := exec.LookPath("just"); err != nil {
		t.Skip("just not installed, skipping")
	}

	t.Run("recipe exists", func(t *testing.T) {
		dir := t.TempDir()
		content := "test:\n    echo ok\n"
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if !hasJustRecipe(dir, "test") {
			t.Error("hasJustRecipe() = false, want true for existing recipe")
		}
	})

	t.Run("recipe does not exist", func(t *testing.T) {
		dir := t.TempDir()
		content := "build:\n    echo build\n"
		if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
		if hasJustRecipe(dir, "test") {
			t.Error("hasJustRecipe() = true, want false for missing recipe")
		}
	})

	t.Run("no justfile", func(t *testing.T) {
		dir := t.TempDir()
		if hasJustRecipe(dir, "test") {
			t.Error("hasJustRecipe() = true, want false when no justfile")
		}
	})
}

func TestWriteLatestMd(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create feature directory structure
	if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
		t.Fatal(err)
	}

	stats := TestStats{
		Total:     5,
		Pass:      3,
		Fail:      2,
		Framework: "npm",
	}
	failures := []TestFailure{
		{
			TestName:     "Login with invalid credentials",
			TestCaseID:   "ui-login-login-with-invalid-credentials",
			ErrorMessage: "Expected status 401, got 500",
		},
	}

	err := writeLatestMd(dir, "test", stats, failures)
	if err != nil {
		t.Fatalf("writeLatestMd() error = %v", err)
	}

	resultsDir := filepath.Join(dir, feature.GetFeatureTestingResultsDir("test"))
	latestPath := filepath.Join(resultsDir, "latest.md")

	data, err := os.ReadFile(latestPath)
	if err != nil {
		t.Fatalf("failed to read latest.md: %v", err)
	}

	fileContent := string(data)
	if !strings.Contains(fileContent, "# Test Results: test") {
		t.Error("latest.md missing header")
	}
	if !strings.Contains(fileContent, "FAIL") {
		t.Error("latest.md should show FAIL status")
	}
	if !strings.Contains(fileContent, "failure-ui-login-login-with-invalid-credentials.md") {
		t.Error("latest.md should reference failure file by test case ID")
	}
}

func TestAppendFixTask(t *testing.T) {
	tests := []struct {
		name          string
		existingTasks map[string]task.Task
		failures      []TestFailure
		e2eRound      int
		wantErr       error
		wantAdded     int
		verifyFunc    func(t *testing.T, indexPath string)
	}{
		{
			name: "first failure appends fix-e2e-1-1",
			existingTasks: map[string]task.Task{
				"biz-1": {ID: "1.1", Status: "completed"},
			},
			failures: []TestFailure{
				{
					TestName:     "Login with invalid credentials",
					TestCaseID:   "ui-login-login-with-invalid-credentials",
					ErrorMessage: "Expected status 401, got 500",
				},
			},
			wantErr:   nil,
			wantAdded: 1,
			verifyFunc: func(t *testing.T, indexPath string) {
				index, err := task.LoadIndex(indexPath)
				if err != nil {
					t.Fatal(err)
				}
				if len(index.Tasks) != 2 {
					t.Errorf("expected 2 tasks, got %d", len(index.Tasks))
				}
				fixTask, ok := index.Tasks["fix-e2e-1-1"]
				if !ok {
					t.Error("fix-e2e-1-1 task not found")
				} else {
					if fixTask.ID != "fix-e2e-1-1" {
						t.Errorf("ID = %q, want fix-e2e-1-1", fixTask.ID)
					}
					if fixTask.Priority != "P0" {
						t.Errorf("Priority = %q, want P0", fixTask.Priority)
					}
					if fixTask.Status != "pending" {
						t.Errorf("Status = %q, want pending", fixTask.Status)
					}
					if fixTask.File != "fix-e2e-1-1.md" {
						t.Errorf("File = %q, want fix-e2e-1-1.md", fixTask.File)
					}
				}
			},
		},
		{
			name: "two failures create two tasks",
			existingTasks: map[string]task.Task{
				"biz-1": {ID: "1.1", Status: "completed"},
			},
			failures: []TestFailure{
				{TestName: "Test A", TestCaseID: "tc-a"},
				{TestName: "Test B", TestCaseID: "tc-b"},
			},
			wantErr:   nil,
			wantAdded: 2,
			verifyFunc: func(t *testing.T, indexPath string) {
				index, err := task.LoadIndex(indexPath)
				if err != nil {
					t.Fatal(err)
				}
				if len(index.Tasks) != 3 {
					t.Errorf("expected 3 tasks, got %d", len(index.Tasks))
				}
				if _, ok := index.Tasks["fix-e2e-1-1"]; !ok {
					t.Error("fix-e2e-1-1 not found")
				}
				if _, ok := index.Tasks["fix-e2e-1-2"]; !ok {
					t.Error("fix-e2e-1-2 not found")
				}
			},
		},
		{
			name: "pending fix-e2e exists, skip append",
			existingTasks: map[string]task.Task{
				"biz-1":       {ID: "1.1", Status: "completed"},
				"fix-e2e-1-1": {ID: "fix-e2e-1-1", Status: "pending", Priority: "P0"},
			},
			failures: []TestFailure{
				{TestName: "Login with invalid credentials", TestCaseID: "ui-login"},
			},
			wantErr:   nil,
			wantAdded: 1,
			verifyFunc: func(t *testing.T, indexPath string) {
				index, err := task.LoadIndex(indexPath)
				if err != nil {
					t.Fatal(err)
				}
				if len(index.Tasks) != 2 {
					t.Errorf("expected 2 tasks (no new ones added), got %d", len(index.Tasks))
				}
			},
		},
		{
			name: "fix-e2e round limit (3) reached, returns sentinel",
			existingTasks: map[string]task.Task{
				"biz-1":       {ID: "1.1", Status: "completed"},
				"fix-e2e-1-1": {ID: "fix-e2e-1-1", Status: "completed", Priority: "P0"},
				"fix-e2e-2-1": {ID: "fix-e2e-2-1", Status: "completed", Priority: "P0"},
				"fix-e2e-3-1": {ID: "fix-e2e-3-1", Status: "completed", Priority: "P0"},
			},
			e2eRound: 3,
			failures: []TestFailure{
				{TestName: "Login with invalid credentials", TestCaseID: "ui-login"},
			},
			wantErr:   errFixLimitExceeded,
			wantAdded: 0,
			verifyFunc: func(t *testing.T, indexPath string) {
				index, err := task.LoadIndex(indexPath)
				if err != nil {
					t.Fatal(err)
				}
				if len(index.Tasks) != 4 {
					t.Errorf("expected 4 tasks (unchanged), got %d", len(index.Tasks))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("CLAUDE_PROJECT_DIR", dir)

			// Create feature directory structure
			if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
				t.Fatal(err)
			}

			// Write index.json
			indexPath := filepath.Join(dir, feature.GetFeatureIndexFile("test"))
			index := &task.TaskIndex{
				Feature:    "test",
				StatusEnum: []string{"pending", "in_progress", "completed", "blocked", "skipped"},
				Tasks:      tt.existingTasks,
				E2ERound:   tt.e2eRound,
			}
			if err := task.SaveIndex(indexPath, index); err != nil {
				t.Fatal(err)
			}

			// Run appendFixTask
			added, err := appendFixTask(dir, "test", tt.failures)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("appendFixTask() error = %v, want %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("appendFixTask() unexpected error = %v", err)
			}

			if err == nil && added != tt.wantAdded {
				t.Errorf("appendFixTask() added = %d, want %d", added, tt.wantAdded)
			}

			if tt.verifyFunc != nil {
				tt.verifyFunc(t, indexPath)
			}
		})
	}
}

func TestGraduateTestScripts(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func(t *testing.T, projectRoot, featureSlug string)
		verifyFunc func(t *testing.T, projectRoot, featureSlug string)
	}{
		{
			name: "first success creates marker and copies scripts",
			setupFunc: func(t *testing.T, projectRoot, featureSlug string) {
				// Create test-cases.md with targets
				testCasesPath := filepath.Join(projectRoot, feature.GetFeatureTestCasesFile(featureSlug))
				content := `# Test Cases
## TC-001: Login
- **Target**: ui/login
- **Test ID**: ui/login/login-with-valid-credentials
`
				if err := os.MkdirAll(filepath.Dir(testCasesPath), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(testCasesPath, []byte(content), 0644); err != nil {
					t.Fatal(err)
				}

				// Create scripts directory with ui.spec.ts
				scriptsDir := filepath.Join(projectRoot, feature.GetFeatureTestingScriptsDir(featureSlug))
				if err := os.MkdirAll(scriptsDir, 0755); err != nil {
					t.Fatal(err)
				}
				specContent := "import { test } from 'node:test';\ntest('login', () => {});"
				specPath := filepath.Join(scriptsDir, "ui.spec.ts")
				if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
					t.Fatal(err)
				}
			},
			verifyFunc: func(t *testing.T, projectRoot, featureSlug string) {
				markerPath := feature.GetE2EGraduatedMarker(projectRoot, featureSlug)
				if !fileExists(markerPath) {
					t.Error("graduation marker not created")
				}

				// Check ui/login/ui.spec.ts exists
				uiLoginSpec := filepath.Join(projectRoot, "tests/e2e/ui/login/ui.spec.ts")
				if !fileExists(uiLoginSpec) {
					t.Error("ui/login/ui.spec.ts not copied")
				}
			},
		},
		{
			name: "already graduated skips migration",
			setupFunc: func(t *testing.T, projectRoot, featureSlug string) {
				// Create graduation marker
				markerPath := feature.GetE2EGraduatedMarker(projectRoot, featureSlug)
				if err := os.MkdirAll(filepath.Dir(markerPath), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(markerPath, []byte("2024-01-01T00:00:00Z\n"), 0644); err != nil {
					t.Fatal(err)
				}
			},
			verifyFunc: func(t *testing.T, projectRoot, featureSlug string) {
				// Marker should still exist with original timestamp
				markerPath := feature.GetE2EGraduatedMarker(projectRoot, featureSlug)
				data, err := os.ReadFile(markerPath)
				if err != nil {
					t.Fatal(err)
				}
				if string(data) != "2024-01-01T00:00:00Z\n" {
					t.Error("marker was overwritten")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Setenv("CLAUDE_PROJECT_DIR", dir)

			// Create feature directory structure
			if err := feature.EnsureFeatureDir(dir, "test"); err != nil {
				t.Fatal(err)
			}

			if tt.setupFunc != nil {
				tt.setupFunc(t, dir, "test")
			}

			// Run graduateTestScripts
			err := graduateTestScripts(dir, "test")
			if err != nil {
				t.Fatalf("graduateTestScripts() error = %v", err)
			}

			if tt.verifyFunc != nil {
				tt.verifyFunc(t, dir, "test")
			}
		})
	}
}

func TestSaveIndexAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "index.json")

	index := &task.TaskIndex{
		Feature:    "test",
		StatusEnum: []string{"pending", "completed"},
		Tasks: map[string]task.Task{
			"task1": {ID: "1.1", Status: "pending"},
		},
	}

	err := saveIndexAtomic(path, index)
	if err != nil {
		t.Fatalf("saveIndexAtomic() error = %v", err)
	}

	// Verify file exists
	if !fileExists(path) {
		t.Error("index.json not created")
	}

	// Verify content
	loaded, err := task.LoadIndex(path)
	if err != nil {
		t.Fatalf("failed to load index: %v", err)
	}
	if loaded.Feature != "test" {
		t.Errorf("Feature = %q, want test", loaded.Feature)
	}
}

func TestRunCmdCapture(t *testing.T) {
	dir := t.TempDir()

	output, success := runCmdCapture(dir, "echo", "hello")
	if !success {
		t.Error("runCmdCapture() success = false, want true")
	}
	if !strings.Contains(output, "hello") {
		t.Errorf("runCmdCapture() output = %q, want contain hello", output)
	}
}

func TestCountPassingTests(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   int
	}{
		{
			name:   "TAP ok lines",
			output: "ok 1 - test one\nok 2 - test two\nnot ok 3 - test three",
			want:   2,
		},
		{
			name:   "checkmark lines",
			output: "✓ test one\n✓ test two\n✗ test three",
			want:   2,
		},
		{
			name:   "mixed with noise",
			output: "ok 1\nsome noise ok here\n✓ another\n  ok 2\n",
			want:   3,
		},
		{
			name:   "empty output",
			output: "",
			want:   0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := countPassingTests(tc.output)
			if got != tc.want {
				t.Errorf("countPassingTests() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestRunSpecsIndividually(t *testing.T) {
	dir := t.TempDir()

	// Create a passing spec
	passSpec := `import { test } from 'node:test';
test('pass', () => {});
`
	if err := os.WriteFile(filepath.Join(dir, "a.spec.ts"), []byte(passSpec), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a non-spec file (should be ignored)
	if err := os.WriteFile(filepath.Join(dir, "helper.ts"), []byte("// not a spec"), 0644); err != nil {
		t.Fatal(err)
	}

	output, success := runSpecsIndividually(dir)
	if !success {
		t.Errorf("expected success, got failure. Output: %s", output)
	}
	if !strings.Contains(output, "pass") {
		t.Errorf("output should contain test name, got: %s", output)
	}
}
