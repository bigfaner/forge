package task

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func newTestIndex(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	_ = os.MkdirAll(tasksDir, 0755)

	index := NewTaskIndex("test-feature")
	index.tasks["1.1-init"] = Task{
		ID:       "1.1",
		Title:    "Init project",
		Priority: "P0",
		Status:   "completed",
		File:     "1.1-init.md",
		Record:   "records/1.1-init.md",
	}
	index.tasks["1.2-setup"] = Task{
		ID:       "1.2",
		Title:    "Setup config",
		Priority: "P1",
		Status:   "pending",
		File:     "1.2-setup.md",
		Record:   "records/1.2-setup.md",
	}

	indexPath := filepath.Join(dir, "index.json")
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatal(err)
	}
	return indexPath
}

func TestAddTask_Basic(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:    "Fix auth timeout",
		Priority: "P0",
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "disc-1" {
		t.Errorf("expected disc-1, got %s", id)
	}

	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks["disc-1"]
	if task.Title != "Fix auth timeout" {
		t.Errorf("expected title 'Fix auth timeout', got %s", task.Title)
	}
	if task.Priority != "P0" {
		t.Errorf("expected P0, got %s", task.Priority)
	}
	if task.Status != "pending" {
		t.Errorf("expected pending, got %s", task.Status)
	}
	if task.File != "disc-1.md" {
		t.Errorf("expected disc-1.md, got %s", task.File)
	}
}

func TestAddTask_WithProvidedID(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		ID:       "fix-auth-1",
		Title:    "Fix auth",
		Priority: "P0",
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "fix-auth-1" {
		t.Errorf("expected fix-auth-1, got %s", id)
	}
}

func TestAddTask_AutoGenerateID(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{Title: "First disc"})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "disc-1" {
		t.Errorf("expected disc-1, got %s", id)
	}
}

func TestAddTask_AutoGenerateID_Sequential(t *testing.T) {
	indexPath := newTestIndex(t)

	if _, err := AddTask(indexPath, AddTaskOpts{Title: "Disc 1"}); err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}
	if _, err := AddTask(indexPath, AddTaskOpts{Title: "Disc 2"}); err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}

	id, err := AddTask(indexPath, AddTaskOpts{Title: "Disc 3"})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "disc-3" {
		t.Errorf("expected disc-3, got %s", id)
	}
}

func TestAddTask_AutoGenerateID_MaxPlusOne(t *testing.T) {
	indexPath := newTestIndex(t)

	if _, err := AddTask(indexPath, AddTaskOpts{Title: "Disc 1", ID: "disc-1"}); err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}
	if _, err := AddTask(indexPath, AddTaskOpts{Title: "Disc 3", ID: "disc-3"}); err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}

	id, err := AddTask(indexPath, AddTaskOpts{Title: "Max plus one"})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "disc-4" {
		t.Errorf("expected disc-4 (max+1), got %s", id)
	}
}

func TestAddTask_DuplicateID(t *testing.T) {
	indexPath := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{ID: "1.1", Title: "Duplicate"})
	if err == nil {
		t.Fatal("expected error for duplicate ID")
	}
}

func TestAddTask_InvalidPriority(t *testing.T) {
	indexPath := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{Title: "Bad priority", Priority: "P5"})
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
}

func TestAddTask_EmptyTitle(t *testing.T) {
	indexPath := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{Title: ""})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestAddTask_DefaultPriority(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{Title: "Default prio"})
	if err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}
	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks[id]
	if task.Priority != "P1" {
		t.Errorf("expected default P1, got %s", task.Priority)
	}
}

func TestAddTask_DependencyNotFound(t *testing.T) {
	indexPath := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Bad dep",
		Dependencies: []string{"9.9"},
	})
	if err == nil {
		t.Fatal("expected error for missing dependency")
	}
}

func TestAddTask_DependenciesExist(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "With dep",
		Dependencies: []string{"1.1"},
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks[id]
	if len(task.Dependencies) != 1 || task.Dependencies[0] != "1.1" {
		t.Errorf("expected deps [1.1], got %v", task.Dependencies)
	}
}

func TestAddTask_Breaking(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{Title: "Breaking task", Breaking: true})
	if err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}
	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks[id]
	if !task.Breaking {
		t.Error("expected breaking=true")
	}
}

func TestCreateTaskMarkdown_Basic(t *testing.T) {
	dir := t.TempDir()
	opts := AddTaskOpts{
		ID:       "disc-1",
		Title:    "Fix timeout",
		Priority: "P0",
	}

	err := CreateTaskMarkdown(dir, "disc-1.md", opts)
	if err != nil {
		t.Fatalf("CreateTaskMarkdown failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "disc-1.md"))
	if err != nil {
		t.Fatalf("read file failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, `id: "disc-1"`) {
		t.Errorf("missing id in frontmatter: %s", content)
	}
	if !strings.Contains(content, "# disc-1: Fix timeout") {
		t.Errorf("missing title heading: %s", content)
	}
}

func TestCreateTaskMarkdown_WithBody(t *testing.T) {
	dir := t.TempDir()
	opts := AddTaskOpts{
		ID:          "disc-1",
		Title:       "Fix timeout",
		Priority:    "P0",
		Description: "## Steps\n\n1. Read error\n2. Fix\n",
	}

	err := CreateTaskMarkdown(dir, "disc-1.md", opts)
	if err != nil {
		t.Fatalf("CreateTaskMarkdown failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "disc-1.md"))
	if err != nil {
		t.Fatalf("setup: read file failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "## Steps") {
		t.Errorf("missing description body: %s", content)
	}
}

func TestCreateTaskMarkdown_WithDependencies(t *testing.T) {
	dir := t.TempDir()
	opts := AddTaskOpts{
		ID:           "disc-1",
		Title:        "Fix",
		Priority:     "P1",
		Dependencies: []string{"1.1", "1.2"},
	}

	if err := CreateTaskMarkdown(dir, "disc-1.md", opts); err != nil {
		t.Fatalf("setup: CreateTaskMarkdown failed: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "disc-1.md"))
	if err != nil {
		t.Fatalf("setup: read file failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, `"1.1"`) || !strings.Contains(content, `"1.2"`) {
		t.Errorf("missing dependencies: %s", content)
	}
}

func TestCreateTaskMarkdown_Breaking(t *testing.T) {
	dir := t.TempDir()
	opts := AddTaskOpts{
		ID:       "disc-1",
		Title:    "Fix",
		Priority: "P0",
		Breaking: true,
	}

	if err := CreateTaskMarkdown(dir, "disc-1.md", opts); err != nil {
		t.Fatalf("setup: CreateTaskMarkdown failed: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "disc-1.md"))
	if err != nil {
		t.Fatalf("setup: read file failed: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "breaking: true") {
		t.Errorf("missing breaking: %s", content)
	}
}

func TestGenerateAutoID_Empty(t *testing.T) {
	index := NewTaskIndex("test")
	id := generateAutoID("disc", index)
	if id != "disc-1" {
		t.Errorf("expected disc-1, got %s", id)
	}
}

func TestGenerateAutoID_Sequential(t *testing.T) {
	index := NewTaskIndex("test")
	index.tasks["disc-1"] = Task{ID: "disc-1", Title: "D1", Priority: "P1", Status: "completed", File: "disc-1.md", Record: "records/disc-1.md"}
	index.tasks["disc-2"] = Task{ID: "disc-2", Title: "D2", Priority: "P1", Status: "completed", File: "disc-2.md", Record: "records/disc-2.md"}
	id := generateAutoID("disc", index)
	if id != "disc-3" {
		t.Errorf("expected disc-3, got %s", id)
	}
}

func TestGenerateAutoID_NonDiscIgnored(t *testing.T) {
	index := NewTaskIndex("test")
	index.tasks["1.1-init"] = Task{ID: "1.1", Title: "Init", Priority: "P0", Status: "completed", File: "1.1-init.md", Record: "records/1.1-init.md"}
	index.tasks["fix-e2e-1-1"] = Task{ID: "fix-e2e-1-1", Title: "Fix e2e", Priority: "P0", Status: "completed", File: "fix-e2e-1-1.md", Record: "records/fix-e2e-1-1.md"}
	id := generateAutoID("disc", index)
	if id != "disc-1" {
		t.Errorf("expected disc-1, got %s", id)
	}
}

func TestGenerateAutoID_FixPrefix(t *testing.T) {
	tests := []struct {
		name     string
		existing map[string]Task
		prefix   string
		expected string
	}{
		{
			name:     "empty index",
			existing: nil,
			prefix:   "fix",
			expected: "fix-1",
		},
		{
			name: "sequential",
			existing: map[string]Task{
				"fix-1": {ID: "fix-1"},
				"fix-2": {ID: "fix-2"},
			},
			prefix:   "fix",
			expected: "fix-3",
		},
		{
			name: "gap skipped (max+1, not gap-fill)",
			existing: map[string]Task{
				"fix-1": {ID: "fix-1"},
				"fix-3": {ID: "fix-3"},
			},
			prefix:   "fix",
			expected: "fix-4",
		},
		{
			name: "disc tasks ignored for fix prefix",
			existing: map[string]Task{
				"disc-1": {ID: "disc-1"},
				"disc-2": {ID: "disc-2"},
				"fix-1":  {ID: "fix-1"},
			},
			prefix:   "fix",
			expected: "fix-2",
		},
		{
			name: "fix prefix ignores disc",
			existing: map[string]Task{
				"fix-1": {ID: "fix-1"},
			},
			prefix:   "disc",
			expected: "disc-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := NewTaskIndex("test")
			if tt.existing != nil {
				for k, v := range tt.existing {
					index.tasks[k] = v
				}
			}
			id := generateAutoID(tt.prefix, index)
			if id != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, id)
			}
		})
	}
}

func TestApplyVars(t *testing.T) {
	tests := []struct {
		name      string
		tmpl      string
		opts      AddTaskOpts
		expected  string
		wantError bool
	}{
		{
			name:     "builtin ID and TITLE",
			tmpl:     "id: {{ID}}, title: {{TITLE}}",
			opts:     AddTaskOpts{ID: "fix-1", Title: "Fix: login bug"},
			expected: "id: fix-1, title: Fix: login bug",
		},
		{
			name:     "user variable",
			tmpl:     "source: {{SOURCE_TASK_ID}}",
			opts:     AddTaskOpts{Vars: map[string]string{"SOURCE_TASK_ID": "T-test-3"}},
			expected: "source: T-test-3",
		},
		{
			name:     "user var overrides builtin",
			tmpl:     "{{TITLE}}",
			opts:     AddTaskOpts{Title: "original", Vars: map[string]string{"TITLE": "overridden"}},
			expected: "overridden",
		},
		{
			name:      "missing var returns error",
			tmpl:      "keep {{UNKNOWN}} placeholder",
			opts:      AddTaskOpts{},
			wantError: true,
		},
		{
			name:     "no placeholders",
			tmpl:     "plain text",
			opts:     AddTaskOpts{},
			expected: "plain text",
		},
		{
			name:     "multiple same placeholder",
			tmpl:     "{{ID}}-{{ID}}",
			opts:     AddTaskOpts{ID: "fix-1"},
			expected: "fix-1-fix-1",
		},
		{
			name:     "builtin DESCRIPTION",
			tmpl:     "desc: {{DESCRIPTION}}",
			opts:     AddTaskOpts{Description: "root cause"},
			expected: "desc: root cause",
		},
		{
			name:     "builtin PRIORITY",
			tmpl:     "prio: {{PRIORITY}}",
			opts:     AddTaskOpts{Priority: "P0"},
			expected: "prio: P0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ApplyVars(tt.tmpl, tt.opts)
			if tt.wantError {
				if err == nil {
					t.Fatal("ApplyVars() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("ApplyVars() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("ApplyVars() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCreateTaskMarkdown_TemplateMode(t *testing.T) {
	dir := t.TempDir()

	opts := AddTaskOpts{
		ID:          "fix-1",
		Title:       "Fix: login selector mismatch",
		Priority:    "P0",
		Description: "Selector [data-testid='submit-btn'] not found.",
		Template:    "fix-task",
		Vars: map[string]string{
			"SOURCE_TASK_ID": "T-test-3",
			"SOURCE_FILES":   "src/components/Login.tsx",
			"TEST_SCRIPT":    "tests/e2e/features/auth/login.spec.ts",
			"TEST_RESULTS":   "tests/e2e/features/auth/results/latest.md",
		},
	}

	if err := CreateTaskMarkdown(dir, "fix-1.md", opts); err != nil {
		t.Fatalf("CreateTaskMarkdown() error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "fix-1.md"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(data)

	checks := []string{
		`id: "fix-1"`,
		`title: "Fix: login selector mismatch"`,
		`priority: "P0"`,
		"Selector [data-testid='submit-btn'] not found.",
		"automatically restored to pending",
		"## Reference Files",
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, got)
		}
	}

	if strings.Contains(got, "{{") {
		t.Errorf("unsubstituted placeholder remaining in output:\n%s", got)
	}
}

func TestCreateTaskMarkdown_TemplateNotFound(t *testing.T) {
	dir := t.TempDir()
	opts := AddTaskOpts{
		ID:       "fix-1",
		Title:    "Fix: test",
		Template: "nonexistent",
	}
	err := CreateTaskMarkdown(dir, "fix-1.md", opts)
	if err == nil {
		t.Fatal("expected error for missing template")
	}
}

func TestAddDependency(t *testing.T) {
	indexPath := newTestIndex(t)

	if err := AddDependency(indexPath, "1.2-setup", "disc-1"); err != nil {
		t.Fatalf("AddDependency failed: %v", err)
	}

	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks["1.2-setup"]
	if !slices.Contains(task.Dependencies, "disc-1") {
		t.Errorf("expected disc-1 in dependencies, got %v", task.Dependencies)
	}
}

func TestAddDependency_Duplicate(t *testing.T) {
	indexPath := newTestIndex(t)

	if err := AddDependency(indexPath, "1.2-setup", "disc-1"); err != nil {
		t.Fatalf("setup: AddDependency failed: %v", err)
	}
	err := AddDependency(indexPath, "1.2-setup", "disc-1")
	if err != nil {
		t.Errorf("duplicate AddDependency should be no-op, got: %v", err)
	}

	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks["1.2-setup"]
	count := 0
	for _, d := range task.Dependencies {
		if d == "disc-1" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 occurrence of disc-1, got %d", count)
	}
}

func TestAddDependency_TaskNotFound(t *testing.T) {
	indexPath := newTestIndex(t)
	err := AddDependency(indexPath, "nonexistent", "disc-1")
	if err == nil {
		t.Fatal("expected error for nonexistent task")
	}
}

func TestGetUnmetDependencies(t *testing.T) {
	indexPath := newTestIndex(t)

	// 1.1 is completed, 1.2 is pending — depend on both
	if _, err := AddTask(indexPath, AddTaskOpts{Title: "Fix", ID: "fix-1"}); err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}
	if err := AddDependency(indexPath, "1.2-setup", "fix-1"); err != nil {
		t.Fatalf("setup: AddDependency failed: %v", err)
	}

	unmet, err := GetUnmetDependencies(indexPath, "1.2-setup")
	if err != nil {
		t.Fatalf("GetUnmetDependencies failed: %v", err)
	}
	// fix-1 is pending (not completed) → unmet
	if !slices.Contains(unmet, "fix-1") {
		t.Errorf("expected fix-1 in unmet, got %v", unmet)
	}

	// Complete fix-1
	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	fixTask := index.tasks["fix-1"]
	fixTask.Status = "completed"
	index.tasks["fix-1"] = fixTask
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("setup: SaveIndex failed: %v", err)
	}

	unmet2, _ := GetUnmetDependencies(indexPath, "1.2-setup")
	if slices.Contains(unmet2, "fix-1") {
		t.Errorf("fix-1 is completed, should not be unmet, got %v", unmet2)
	}
}

func TestAddTask_SourceTaskID_Persisted(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "1.1-init",
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks[id]
	if task.SourceTaskID != "1.1-init" {
		t.Errorf("expected sourceTaskID '1.1-init', got %q", task.SourceTaskID)
	}
}

func TestAddTask_SourceTaskID_UpdatesSourceDeps(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "1.1-init",
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	srcTask := index.tasks["1.1-init"]
	if !slices.Contains(srcTask.Dependencies, id) {
		t.Errorf("source task should have %s as dependency, got %v", id, srcTask.Dependencies)
	}
}

func TestAddTask_SourceTaskID_SourceNotFound(t *testing.T) {
	indexPath := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "nonexistent",
	})
	if err != nil {
		t.Fatalf("AddTask should succeed even if source not found, got: %v", err)
	}

	index, err := LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	task := index.tasks[id]
	if task.SourceTaskID != "nonexistent" {
		t.Errorf("SourceTaskID should still be persisted, got %q", task.SourceTaskID)
	}
}

func TestAddTask_SourceTaskID_IdempotentDep(t *testing.T) {
	indexPath := newTestIndex(t)

	id1, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1-init",
	})
	if err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}

	// Complete fix-1 so dedup allows fix-2
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[id1]
	fix1.Status = "completed"
	index.tasks[id1] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Add again with same source — source dep should not duplicate
	id2, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: "1.1-init",
	})
	if err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}

	index, err = LoadIndex(indexPath)
	if err != nil {
		t.Fatalf("setup: LoadIndex failed: %v", err)
	}
	srcTask := index.tasks["1.1-init"]
	count := 0
	for _, d := range srcTask.Dependencies {
		if d == id1 || d == id2 {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 distinct deps, got %d in %v", count, srcTask.Dependencies)
	}

	// Verify id1 only appears once
	count1 := 0
	for _, d := range srcTask.Dependencies {
		if d == id1 {
			count1++
		}
	}
	if count1 != 1 {
		t.Errorf("id1 should appear exactly once, got %d", count1)
	}
}

// TestAddTask_SourceTaskID_LookupByID verifies that SourceTaskID lookup works
// when the source task's key differs from its ID (e.g. key="1.1-init", id="1.1").
// This is the core bug: index.tasks[opts.SourceTaskID] fails when SourceTaskID is
// the task ID but the map key is a slug.
func TestAddTask_SourceTaskID_LookupByID(t *testing.T) {
	indexPath := newTestIndex(t)

	// SourceTaskID uses the task ID "1.1", but the map key is "1.1-init"
	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "1.1", // task ID, not map key
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	srcTask := index.tasks["1.1-init"]
	if !slices.Contains(srcTask.Dependencies, id) {
		t.Errorf("source task should have %s as dependency (looked up by ID), got deps %v", id, srcTask.Dependencies)
	}
}

// TestAddTask_SourceTaskID_LookupByID_PreservesExistingDeps verifies that appending
// a new dep by-ID lookup does not clobber the source task's existing dependencies.
func TestAddTask_SourceTaskID_LookupByID_PreservesExistingDeps(t *testing.T) {
	indexPath := newTestIndex(t)

	// Give the source task an existing dependency before the add
	index, _ := LoadIndex(indexPath)
	src := index.tasks["1.1-init"]
	src.Dependencies = []string{"some-other-task"}
	index.tasks["1.1-init"] = src
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "1.1", // task ID, map key is "1.1-init"
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ = LoadIndex(indexPath)
	srcTask := index.tasks["1.1-init"]
	if !slices.Contains(srcTask.Dependencies, "some-other-task") {
		t.Errorf("existing dep 'some-other-task' was lost, got %v", srcTask.Dependencies)
	}
	if !slices.Contains(srcTask.Dependencies, id) {
		t.Errorf("new dep %s missing, got %v", id, srcTask.Dependencies)
	}
}

// TestAddTask_SourceTaskID_DynamicTaskWhereKeyEqualsID verifies lookup works
// for dynamically added tasks where the map key equals the task ID.
func TestAddTask_SourceTaskID_DynamicTaskWhereKeyEqualsID(t *testing.T) {
	indexPath := newTestIndex(t)

	// First add a disc task (key == ID)
	firstID, _ := AddTask(indexPath, AddTaskOpts{
		Title:    "Disc task",
		Priority: "P1",
	})

	// Now add another task with SourceTaskID pointing to the disc task
	secondID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix for disc",
		Priority:     "P0",
		SourceTaskID: firstID, // disc-1: key == ID
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	srcTask := index.tasks[firstID]
	if !slices.Contains(srcTask.Dependencies, secondID) {
		t.Errorf("disc task should have %s as dependency, got %v", secondID, srcTask.Dependencies)
	}
}

// TestAddTask_SourceTaskID_MultipleAddsToSameSourceByID verifies that multiple
// tasks added with the same SourceTaskID (by ID) all appear as dependencies.
func TestAddTask_SourceTaskID_MultipleAddsToSameSourceByID(t *testing.T) {
	indexPath := newTestIndex(t)

	id1, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix A",
		SourceTaskID: "1.1", // ID, not key
	})
	// Complete fix A so dedup allows fix B
	index, _ := LoadIndex(indexPath)
	fixA := index.tasks[id1]
	fixA.Status = "completed"
	index.tasks[id1] = fixA
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	id2, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix B",
		SourceTaskID: "1.1",
	})
	// Complete fix B so dedup allows fix C
	index, _ = LoadIndex(indexPath)
	fixB := index.tasks[id2]
	fixB.Status = "completed"
	index.tasks[id2] = fixB
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	id3, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix C",
		SourceTaskID: "1.1",
	})

	index, _ = LoadIndex(indexPath)
	srcTask := index.tasks["1.1-init"]
	for _, id := range []string{id1, id2, id3} {
		if !slices.Contains(srcTask.Dependencies, id) {
			t.Errorf("source missing dep %s, got %v", id, srcTask.Dependencies)
		}
	}
	if len(srcTask.Dependencies) != 3 {
		t.Errorf("expected 3 deps, got %d: %v", len(srcTask.Dependencies), srcTask.Dependencies)
	}
}

// TestAddTask_SourceTaskID_LookupByID_SourceNotFoundIsNoOp verifies that passing
// a nonexistent ID does not error and does not corrupt the index.
func TestAddTask_SourceTaskID_LookupByID_SourceNotFoundIsNoOp(t *testing.T) {
	indexPath := newTestIndex(t)

	indexBefore, _ := LoadIndex(indexPath)
	taskCountBefore := len(indexBefore.tasks)

	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Orphan fix",
		Priority:     "P0",
		SourceTaskID: "9.9-nonexistent", // no such ID or key
	})
	if err != nil {
		t.Fatalf("AddTask should succeed even for missing source, got: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	if len(index.tasks) != taskCountBefore+1 {
		t.Errorf("expected %d tasks, got %d", taskCountBefore+1, len(index.tasks))
	}
	// Original tasks unchanged
	for key, before := range indexBefore.tasks {
		after := index.tasks[key]
		if len(after.Dependencies) != len(before.Dependencies) {
			t.Errorf("task %s deps changed unexpectedly: before=%v after=%v", key, before.Dependencies, after.Dependencies)
		}
	}
}

// --- Source Resolution Tests (auto-resolve --source-task-id when source is completed) ---

func TestAddTask_SourceResolution_CompletedSourceResolvesToRoot(t *testing.T) {
	indexPath := newTestIndex(t)

	// disc-1 is a fix-task for source "1.1" — mark it completed
	fix1ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		Priority:     "P0",
		SourceTaskID: "1.1", // points to root source
	})
	if err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}
	// Mark disc-1 as completed (simulates it finishing its quality gate)
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "completed"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Now add disc-2 with --source-task-id pointing to disc-1 (a COMPLETED fix-task)
	// Auto-resolve should trace disc-1 → 1.1 (root) because disc-1 is completed
	fix2ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2 (deeper)",
		Priority:     "P0",
		SourceTaskID: fix1ID, // pointing to completed fix-task
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	// disc-2 should have SourceTaskID resolved to "1.1" (root), not disc-1
	index, _ = LoadIndex(indexPath)
	if index.tasks[fix2ID].SourceTaskID != "1.1" {
		t.Errorf("disc-2 SourceTaskID should resolve to '1.1', got %q", index.tasks[fix2ID].SourceTaskID)
	}

	// Source task "1.1" should have BOTH disc-1 AND disc-2 as dependencies
	srcTask := index.tasks["1.1-init"]
	if !slices.Contains(srcTask.Dependencies, fix1ID) {
		t.Errorf("source should have %s as dep, got %v", fix1ID, srcTask.Dependencies)
	}
	if !slices.Contains(srcTask.Dependencies, fix2ID) {
		t.Errorf("source should have %s as dep, got %v", fix2ID, srcTask.Dependencies)
	}

	// disc-1 should NOT have disc-2 as dependency (resolved to root, flat for completed source)
	if slices.Contains(index.tasks[fix1ID].Dependencies, fix2ID) {
		t.Errorf("disc-1 should NOT have disc-2 as dep (resolved to root), got %v", index.tasks[fix1ID].Dependencies)
	}
}

func TestAddTask_SourceResolution_BlockedSourcePreservesChain(t *testing.T) {
	indexPath := newTestIndex(t)

	// disc-1 is a fix-task for source "1.1" — mark it BLOCKED (its quality gate failed)
	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		Priority:     "P0",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "blocked"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Now add disc-2 with --source-task-id pointing to disc-1 (a BLOCKED fix-task)
	// Chain model should be preserved — no resolution
	fix2ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2 (nested)",
		Priority:     "P0",
		SourceTaskID: fix1ID, // pointing to blocked fix-task
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ = LoadIndex(indexPath)
	// disc-2 should keep SourceTaskID as disc-1 (chain preserved)
	if index.tasks[fix2ID].SourceTaskID != fix1ID {
		t.Errorf("disc-2 SourceTaskID should remain %q (chain preserved), got %q", fix1ID, index.tasks[fix2ID].SourceTaskID)
	}

	// disc-1 should have disc-2 as dependency (chain model)
	if !slices.Contains(index.tasks[fix1ID].Dependencies, fix2ID) {
		t.Errorf("disc-1 should have disc-2 as dep (chain model), got %v", index.tasks[fix1ID].Dependencies)
	}

	// Root source should NOT have disc-2 as direct dep
	srcTask := index.tasks["1.1-init"]
	if slices.Contains(srcTask.Dependencies, fix2ID) {
		t.Errorf("root source should NOT have disc-2 as dep (chain model), got %v", srcTask.Dependencies)
	}
}

func TestAddTask_SourceResolution_MultiLevelCompletedChain(t *testing.T) {
	indexPath := newTestIndex(t)

	// Chain: disc-1 → 1.1 (root), mark completed
	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	tmp := index.tasks[fix1ID]
	tmp.Status = "completed"
	index.tasks[fix1ID] = tmp
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}
	fix2ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: fix1ID,
	})
	index, _ = LoadIndex(indexPath)
	tmp2 := index.tasks[fix2ID]
	tmp2.Status = "completed"
	index.tasks[fix2ID] = tmp2
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// disc-3 → disc-2 (completed) → disc-1 (completed) → resolves to 1.1
	fix3ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 3",
		SourceTaskID: fix2ID,
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ = LoadIndex(indexPath)

	// All should resolve to "1.1"
	for _, id := range []string{fix1ID, fix2ID, fix3ID} {
		if index.tasks[id].SourceTaskID != "1.1" {
			t.Errorf("%s SourceTaskID should be '1.1', got %q", id, index.tasks[id].SourceTaskID)
		}
	}

	// Source should have all 3 as deps
	srcTask := index.tasks["1.1-init"]
	if len(srcTask.Dependencies) != 3 {
		t.Errorf("source should have 3 deps, got %d: %v", len(srcTask.Dependencies), srcTask.Dependencies)
	}
}

func TestAddTask_SourceResolution_NoChainPassthrough(t *testing.T) {
	indexPath := newTestIndex(t)

	// 1.1 has no SourceTaskID — direct passthrough, no resolution needed
	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix direct",
		SourceTaskID: "1.1",
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	if index.tasks[id].SourceTaskID != "1.1" {
		t.Errorf("SourceTaskID should remain '1.1', got %q", index.tasks[id].SourceTaskID)
	}
}

func TestAddTask_SourceResolution_CycleDetection(t *testing.T) {
	indexPath := newTestIndex(t)

	// Manually create a cycle: disc-1 → disc-2 → disc-1, both completed
	index, _ := LoadIndex(indexPath)
	index.tasks["disc-1"] = Task{
		ID:           "disc-1",
		Title:        "Fix 1",
		Priority:     "P0",
		Status:       "completed",
		File:         "disc-1.md",
		Record:       "records/disc-1.md",
		SourceTaskID: "disc-2",
	}
	index.tasks["disc-2"] = Task{
		ID:           "disc-2",
		Title:        "Fix 2",
		Priority:     "P0",
		Status:       "completed",
		File:         "disc-2.md",
		Record:       "records/disc-2.md",
		SourceTaskID: "disc-1",
	}
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Adding with --source-task-id disc-1 should not infinite loop
	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 3",
		SourceTaskID: "disc-1",
	})
	if err != nil {
		t.Fatalf("AddTask should handle cycle gracefully, got: %v", err)
	}

	// Should resolve to disc-1 or disc-2 (whichever the cycle stops at)
	index, _ = LoadIndex(indexPath)
	if index.tasks[id].SourceTaskID == "" {
		t.Error("SourceTaskID should be set even with cycle")
	}
}

func TestAddTask_SourceResolution_SkippedSourceResolves(t *testing.T) {
	indexPath := newTestIndex(t)

	// disc-1 is a SKIPPED fix-task for source "1.1"
	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	tmp3 := index.tasks[fix1ID]
	tmp3.Status = "skipped"
	index.tasks[fix1ID] = tmp3
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// disc-2 → disc-1 (skipped) → resolves to 1.1
	fix2ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: fix1ID,
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ = LoadIndex(indexPath)
	if index.tasks[fix2ID].SourceTaskID != "1.1" {
		t.Errorf("disc-2 SourceTaskID should resolve to '1.1' (skipped source), got %q", index.tasks[fix2ID].SourceTaskID)
	}
}

func TestResolveSourceTask_Direct(t *testing.T) {
	indexPath := newTestIndex(t)
	index, _ := LoadIndex(indexPath)

	// 1.1 has no SourceTaskID → direct passthrough
	result := ResolveSourceTask(index, "1.1")
	if result != "1.1" {
		t.Errorf("expected '1.1', got %q", result)
	}
}

func TestResolveSourceTask_NotFound(t *testing.T) {
	indexPath := newTestIndex(t)
	index, _ := LoadIndex(indexPath)

	// Nonexistent task → return as-is
	result := ResolveSourceTask(index, "nonexistent")
	if result != "nonexistent" {
		t.Errorf("expected 'nonexistent', got %q", result)
	}
}

func TestGetUnmetDependencies_Wildcard(t *testing.T) {
	indexPath := newTestIndex(t)

	// Add wildcard dep to 1.2-setup
	if err := AddDependency(indexPath, "1.2-setup", "0.x"); err != nil {
		t.Fatalf("setup: AddDependency failed: %v", err)
	}

	// Add a phase-0 task that's pending
	if _, err := AddTask(indexPath, AddTaskOpts{Title: "Phase 0 task", ID: "0.1", Status: "pending"}); err != nil {
		t.Fatalf("setup: AddTask failed: %v", err)
	}

	unmet, err := GetUnmetDependencies(indexPath, "1.2-setup")
	if err != nil {
		t.Fatalf("GetUnmetDependencies failed: %v", err)
	}
	// 0.1 is pending -> unmet, 1.1 is completed -> met
	found := false
	for _, u := range unmet {
		if u == "0.1" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected 0.1 in unmet for wildcard dep, got %v", unmet)
	}
}

func TestGetUnmetDependencies_WildcardAllCompleted(t *testing.T) {
	indexPath := newTestIndex(t)

	// 1.1 is completed. 1.2 matches wildcard but is the task itself — self-excluded.
	if err := AddDependency(indexPath, "1.2-setup", "1.x"); err != nil {
		t.Fatalf("setup: AddDependency failed: %v", err)
	}

	unmet, err := GetUnmetDependencies(indexPath, "1.2-setup")
	if err != nil {
		t.Fatalf("GetUnmetDependencies failed: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("expected 0 unmet (self-excluded + all others completed), got %v", unmet)
	}
}

func TestAddDependency_LookupByID(t *testing.T) {
	indexPath := newTestIndex(t)

	// "1.2" is the task ID, but map key is "1.2-setup"
	err := AddDependency(indexPath, "1.2", "disc-1")
	if err != nil {
		t.Fatalf("AddDependency by ID failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	task := index.tasks["1.2-setup"]
	if !slices.Contains(task.Dependencies, "disc-1") {
		t.Errorf("expected disc-1 in dependencies, got %v", task.Dependencies)
	}
}

func TestAddDependency_LookupByID_NotFound(t *testing.T) {
	indexPath := newTestIndex(t)

	err := AddDependency(indexPath, "9.9", "disc-1")
	if err == nil {
		t.Fatal("expected error for nonexistent task ID")
	}
}

func TestAddDependency_WriteBackUsesSlugKey(t *testing.T) {
	indexPath := newTestIndex(t)

	err := AddDependency(indexPath, "1.2", "disc-1")
	if err != nil {
		t.Fatalf("AddDependency failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	if _, ok := index.tasks["1.2"]; ok {
		t.Error("should not create duplicate entry under ID key '1.2'")
	}
	if _, ok := index.tasks["1.2-setup"]; !ok {
		t.Error("original slug key '1.2-setup' should still exist")
	}
}

func TestGetUnmetDependencies_SlugKeyDeps(t *testing.T) {
	indexPath := newTestIndex(t)

	// Add a new task that depends on slug-keyed task "1.1" (key="1.1-init", id="1.1")
	_, _ = AddTask(indexPath, AddTaskOpts{Title: "Watcher", Dependencies: []string{"1.1"}})

	// 1.1 is completed → should have 0 unmet
	unmet, err := GetUnmetDependencies(indexPath, "disc-1")
	if err != nil {
		t.Fatalf("GetUnmetDependencies failed: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("expected 0 unmet (1.1 is completed), got %v", unmet)
	}
}

func TestGetUnmetDependencies_SlugKeyDeps_Pending(t *testing.T) {
	indexPath := newTestIndex(t)

	// Depends on "1.2" (key="1.2-setup", status=pending)
	_, _ = AddTask(indexPath, AddTaskOpts{Title: "Watcher", Dependencies: []string{"1.2"}})

	unmet, _ := GetUnmetDependencies(indexPath, "disc-1")
	if !slices.Contains(unmet, "1.2") {
		t.Errorf("expected 1.2 in unmet, got %v", unmet)
	}
}

func TestGetUnmetDependencies_LookupByID(t *testing.T) {
	indexPath := newTestIndex(t)

	// Query by task ID "1.2", not by key "1.2-setup"
	unmet, err := GetUnmetDependencies(indexPath, "1.2")
	if err != nil {
		t.Fatalf("GetUnmetDependencies by ID failed: %v", err)
	}
	// 1.2-setup depends on nothing → 0 unmet
	if len(unmet) != 0 {
		t.Errorf("expected 0 unmet, got %v", unmet)
	}
}

func TestGetUnmetDependencies_LookupByID_NotFound(t *testing.T) {
	indexPath := newTestIndex(t)

	_, err := GetUnmetDependencies(indexPath, "9.9")
	if err == nil {
		t.Fatal("expected error for nonexistent task ID")
	}
}

func TestGetUnmetDependencies_AllSlugKeyedCompleted(t *testing.T) {
	indexPath := newTestIndex(t)

	// Both deps are slug-keyed: "1.1" (key="1.1-init"), "1.2" (key="1.2-setup")
	// 1.1 is completed, make 1.2 completed too
	index, _ := LoadIndex(indexPath)
	t2 := index.tasks["1.2-setup"]
	t2.Status = "completed"
	index.tasks["1.2-setup"] = t2
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	_, _ = AddTask(indexPath, AddTaskOpts{Title: "Watcher", Dependencies: []string{"1.1", "1.2"}})

	unmet, _ := GetUnmetDependencies(indexPath, "disc-1")
	if len(unmet) != 0 {
		t.Errorf("all slug-keyed deps completed, expected 0 unmet, got %v", unmet)
	}
}

func TestGetUnmetDependencies_SkippedDepTreatedAsMet(t *testing.T) {
	indexPath := newTestIndex(t)

	index, _ := LoadIndex(indexPath)
	t1 := index.tasks["1.1-init"]
	t1.Status = "skipped"
	index.tasks["1.1-init"] = t1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	_, _ = AddTask(indexPath, AddTaskOpts{Title: "Watcher", Dependencies: []string{"1.1"}})

	unmet, _ := GetUnmetDependencies(indexPath, "disc-1")
	if len(unmet) != 0 {
		t.Errorf("skipped dep should be met, got %v", unmet)
	}
}

func TestGetUnmetDependencies_NonexistentDepTreatedAsUnmet(t *testing.T) {
	indexPath := newTestIndex(t)

	// Bypass AddTask dependency validation — directly create task with phantom dep
	index, _ := LoadIndex(indexPath)
	index.tasks["disc-1"] = Task{
		ID:           "disc-1",
		Title:        "Watcher",
		Priority:     "P1",
		Status:       "pending",
		File:         "disc-1.md",
		Record:       "records/disc-1.md",
		Dependencies: []string{"9.9"},
	}
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	unmet, _ := GetUnmetDependencies(indexPath, "disc-1")
	if len(unmet) != 1 || unmet[0] != "9.9" {
		t.Errorf("expected [9.9] unmet, got %v", unmet)
	}
}

func TestGetUnmetDependencies_WildcardWithSlugKeyedTasks(t *testing.T) {
	indexPath := newTestIndex(t)

	// "1.x" wildcard should match slug-keyed tasks "1.1" and "1.2"
	// 1.1 is completed, 1.2 is pending
	_, _ = AddTask(indexPath, AddTaskOpts{Title: "Watcher", Dependencies: []string{"1.x"}})

	unmet, _ := GetUnmetDependencies(indexPath, "disc-1")
	if len(unmet) != 1 || unmet[0] != "1.2" {
		t.Errorf("expected [1.2] unmet from wildcard, got %v", unmet)
	}
}

func TestGetUnmetDependencies_MixedWildcardAndExactSlugKeyed(t *testing.T) {
	indexPath := newTestIndex(t)

	// dep on "1.x" (wildcard) + "1.1" (exact slug-keyed, completed)
	_, _ = AddTask(indexPath, AddTaskOpts{Title: "Watcher", Dependencies: []string{"1.x", "1.1"}})

	// 1.1 completed, 1.2 pending → wildcard reports 1.2 as unmet
	unmet, _ := GetUnmetDependencies(indexPath, "disc-1")
	if len(unmet) != 1 || unmet[0] != "1.2" {
		t.Errorf("expected [1.2] unmet, got %v", unmet)
	}
}

// --- Active fix-task dedup tests ---

func TestAddTask_Dedup_ActiveFixBlocksCreation(t *testing.T) {
	indexPath := newTestIndex(t)

	// Add fix-1 targeting source "1.1" — still pending (active)
	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		Priority:     "P0",
		SourceTaskID: "1.1",
	})
	if err != nil {
		t.Fatalf("setup: AddTask fix-1 failed: %v", err)
	}

	// Try adding fix-2 for the same source — should be blocked by dedup
	_, err = AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		Priority:     "P0",
		SourceTaskID: "1.1",
	})
	if err == nil {
		t.Fatal("expected ActiveFixExistsError, got nil")
	}
	var dedupErr *ActiveFixExistsError
	if !errors.As(err, &dedupErr) {
		t.Fatalf("expected *ActiveFixExistsError, got %T: %v", err, err)
	}
	if dedupErr.SourceTaskID != "1.1" {
		t.Errorf("SourceTaskID = %q, want %q", dedupErr.SourceTaskID, "1.1")
	}
	if len(dedupErr.ActiveFixIDs) != 1 || dedupErr.ActiveFixIDs[0] != "disc-1" {
		t.Errorf("ActiveFixIDs = %v, want [disc-1]", dedupErr.ActiveFixIDs)
	}
}

func TestAddTask_Dedup_CompletedFixAllowsCreation(t *testing.T) {
	indexPath := newTestIndex(t)

	// Add fix-1 and mark it completed
	fix1ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		Priority:     "P0",
		SourceTaskID: "1.1",
	})
	if err != nil {
		t.Fatalf("setup: AddTask fix-1 failed: %v", err)
	}
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "completed"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Adding fix-2 for the same source should succeed (fix-1 is completed)
	fix2ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		Priority:     "P0",
		SourceTaskID: "1.1",
	})
	if err != nil {
		t.Fatalf("AddTask fix-2 should succeed when fix-1 is completed, got: %v", err)
	}
	if fix2ID == "" {
		t.Error("expected non-empty fix-2 ID")
	}
}

func TestAddTask_Dedup_SkippedFixAllowsCreation(t *testing.T) {
	indexPath := newTestIndex(t)

	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "skipped"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: "1.1",
	})
	if err != nil {
		t.Fatalf("should succeed when fix-1 is skipped, got: %v", err)
	}
}

func TestAddTask_Dedup_RejectedFixAllowsCreation(t *testing.T) {
	indexPath := newTestIndex(t)

	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "rejected"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: "1.1",
	})
	if err != nil {
		t.Fatalf("should succeed when fix-1 is rejected, got: %v", err)
	}
}

func TestAddTask_Dedup_BlockedFixBlocksCreation(t *testing.T) {
	indexPath := newTestIndex(t)

	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "blocked"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: "1.1",
	})
	if err == nil {
		t.Fatal("expected ActiveFixExistsError for blocked fix-task")
	}
	var dedupErr *ActiveFixExistsError
	if !errors.As(err, &dedupErr) {
		t.Fatalf("expected *ActiveFixExistsError, got %T: %v", err, err)
	}
	if dedupErr.ActiveFixIDs[0] != fix1ID {
		t.Errorf("ActiveFixIDs = %v, want [%s]", dedupErr.ActiveFixIDs, fix1ID)
	}
}

func TestAddTask_Dedup_NoSourceTaskID(t *testing.T) {
	indexPath := newTestIndex(t)

	// No dedup check when SourceTaskID is empty
	_, err := AddTask(indexPath, AddTaskOpts{Title: "Task A"})
	if err != nil {
		t.Fatalf("AddTask without sourceTaskID failed: %v", err)
	}
	_, err = AddTask(indexPath, AddTaskOpts{Title: "Task B"})
	if err != nil {
		t.Fatalf("second AddTask without sourceTaskID failed: %v", err)
	}
}

func TestAddTask_Dedup_ResolvedSourceChecksAgainstRoot(t *testing.T) {
	indexPath := newTestIndex(t)

	// fix-1 targets source "1.1", mark it completed (triggers source resolution)
	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "completed"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// fix-2 targets completed fix-1 → resolves to "1.1" (root), then checks dedup.
	// fix-1 is completed → dedup passes → fix-2 created
	fix2ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: fix1ID,
	})
	if err != nil {
		t.Fatalf("AddTask fix-2 should succeed (fix-1 completed, resolved to 1.1), got: %v", err)
	}

	// Mark fix-2 as pending (active), then try fix-3 targeting "1.1"
	index, _ = LoadIndex(indexPath)
	// fix-2 resolved to SourceTaskID "1.1" — it's active (pending)
	if index.tasks[fix2ID].SourceTaskID != "1.1" {
		t.Fatalf("fix-2 SourceTaskID should resolve to '1.1', got %q", index.tasks[fix2ID].SourceTaskID)
	}

	// fix-3 targets "1.1" — fix-2 is active with source "1.1" → dedup blocks
	_, err = AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 3",
		SourceTaskID: "1.1",
	})
	if err == nil {
		t.Fatal("expected ActiveFixExistsError (fix-2 is active for source 1.1)")
	}
	var dedupErr *ActiveFixExistsError
	if !errors.As(err, &dedupErr) {
		t.Fatalf("expected *ActiveFixExistsError, got %T: %v", err, err)
	}
}

func TestAddTask_Dedup_MixedActiveAndCompleted(t *testing.T) {
	indexPath := newTestIndex(t)

	// fix-1 completed
	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "completed"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// fix-2 pending (active)
	_, _ = AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: "1.1",
	})

	// Try fix-3 → blocked by active fix-2
	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 3",
		SourceTaskID: "1.1",
	})
	if err == nil {
		t.Fatal("expected ActiveFixExistsError")
	}
	var dedupErr *ActiveFixExistsError
	if !errors.As(err, &dedupErr) {
		t.Fatalf("expected *ActiveFixExistsError, got %T: %v", err, err)
	}
	// Should only report fix-2 (active), not fix-1 (completed)
	if len(dedupErr.ActiveFixIDs) != 1 {
		t.Errorf("expected 1 active fix, got %d: %v", len(dedupErr.ActiveFixIDs), dedupErr.ActiveFixIDs)
	}
}

// --- BlockSource tests ---

func TestAddTask_BlockSource_SetsBlocked(t *testing.T) {
	indexPath := newTestIndex(t)

	// 1.2 is pending — block it via --block-source
	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix for 1.2",
		Priority:     "P0",
		SourceTaskID: "1.2",
		BlockSource:  true,
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	srcTask := index.tasks["1.2-setup"]
	if srcTask.Status != "blocked" {
		t.Errorf("source should be blocked, got %q", srcTask.Status)
	}
}

func TestAddTask_BlockSource_CompletedSourcePreservesChain(t *testing.T) {
	indexPath := newTestIndex(t)

	// fix-1 targets source "1.1" — mark completed
	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "completed"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Add fix-2 with --block-source targeting completed fix-1
	// BlockSource sets fix-1 to blocked BEFORE resolution → no resolution → chain preserved
	fix2ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2 (fix-of-fix)",
		SourceTaskID: fix1ID,
		BlockSource:  true,
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ = LoadIndex(indexPath)
	// fix-2 should chain to fix-1 (not resolve to root 1.1)
	if index.tasks[fix2ID].SourceTaskID != fix1ID {
		t.Errorf("fix-2 SourceTaskID should be %q (chain preserved), got %q", fix1ID, index.tasks[fix2ID].SourceTaskID)
	}
	// fix-1 should be blocked
	if index.tasks[fix1ID].Status != "blocked" {
		t.Errorf("fix-1 should be blocked, got %q", index.tasks[fix1ID].Status)
	}
	// fix-2 should be a dependency of fix-1 (not root 1.1)
	if !slices.Contains(index.tasks[fix1ID].Dependencies, fix2ID) {
		t.Errorf("fix-1 should have fix-2 as dep, got %v", index.tasks[fix1ID].Dependencies)
	}
}

func TestAddTask_BlockSource_WithoutFlagFlattensToRoot(t *testing.T) {
	indexPath := newTestIndex(t)

	// fix-1 targets "1.1" — mark completed
	fix1ID, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1",
	})
	index, _ := LoadIndex(indexPath)
	fix1 := index.tasks[fix1ID]
	fix1.Status = "completed"
	index.tasks[fix1ID] = fix1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// Without --block-source, completed fix-1 resolves to root "1.1"
	fix2ID, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: fix1ID,
		BlockSource:  false,
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ = LoadIndex(indexPath)
	// fix-2 should flatten to root "1.1"
	if index.tasks[fix2ID].SourceTaskID != "1.1" {
		t.Errorf("fix-2 SourceTaskID should resolve to '1.1', got %q", index.tasks[fix2ID].SourceTaskID)
	}
	// fix-1 should still be completed (not blocked)
	if index.tasks[fix1ID].Status != "completed" {
		t.Errorf("fix-1 should remain completed, got %q", index.tasks[fix1ID].Status)
	}
}

func TestAddTask_BlockSource_SourceNotFound(t *testing.T) {
	indexPath := newTestIndex(t)

	// --block-source with nonexistent source — should succeed (no source to block)
	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix orphan",
		SourceTaskID: "9.9-nonexistent",
		BlockSource:  true,
	})
	if err != nil {
		t.Fatalf("AddTask should succeed with missing source, got: %v", err)
	}
}

func TestAddTask_BlockSource_NoSourceTaskID(t *testing.T) {
	indexPath := newTestIndex(t)

	// --block-source without --source-task-id — flag is ignored
	_, err := AddTask(indexPath, AddTaskOpts{
		Title:       "Plain task",
		BlockSource: true,
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
}

func TestAddTask_BlockSource_DedupPreventsMutation(t *testing.T) {
	indexPath := newTestIndex(t)

	// Add first fix task for source 1.1 (active/pending)
	_, _ = AddTask(indexPath, AddTaskOpts{
		Title:        "Fix: first attempt",
		SourceTaskID: "1.1",
	})

	// Verify 1.1 is still "completed" before the dedup test
	index, _ := LoadIndex(indexPath)
	if index.tasks["1.1-init"].Status != "completed" {
		t.Fatalf("1.1 should be completed before dedup test, got %q", index.tasks["1.1-init"].Status)
	}

	// Try to add second fix with --block-source — dedup should trigger BEFORE blocking
	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix: second attempt",
		SourceTaskID: "1.1",
		BlockSource:  true,
	})
	var dedupErr *ActiveFixExistsError
	if !errors.As(err, &dedupErr) {
		t.Fatalf("expected ActiveFixExistsError, got %v", err)
	}

	// Key assertion: source 1.1 must NOT have been mutated to "blocked"
	index, _ = LoadIndex(indexPath)
	if index.tasks["1.1-init"].Status != "completed" {
		t.Errorf("source 1.1 should remain completed (dedup prevented blocking mutation), got %q",
			index.tasks["1.1-init"].Status)
	}
}

func TestAddTask_BlockSource_AlreadyBlocked(t *testing.T) {
	indexPath := newTestIndex(t)

	// Mark 1.2 as blocked manually
	index, _ := LoadIndex(indexPath)
	t1 := index.tasks["1.2-setup"]
	t1.Status = "blocked"
	index.tasks["1.2-setup"] = t1
	if err := SaveIndex(indexPath, index); err != nil {
		t.Fatalf("SaveIndex failed: %v", err)
	}

	// --block-source on already-blocked source — idempotent
	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix for blocked 1.2",
		SourceTaskID: "1.2",
		BlockSource:  true,
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ = LoadIndex(indexPath)
	if index.tasks["1.2-setup"].Status != "blocked" {
		t.Errorf("source should remain blocked, got %q", index.tasks["1.2-setup"].Status)
	}
}
