package task

import (
	"os"
	"path/filepath"
	"testing"
)

func newTestIndex(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()
	tasksDir := filepath.Join(dir, "tasks")
	os.MkdirAll(tasksDir, 0755)

	index := NewTaskIndex("test-feature")
	index.Tasks["1.1-init"] = Task{
		ID:       "1.1",
		Title:    "Init project",
		Priority: "P0",
		Status:   "completed",
		File:     "1.1-init.md",
		Record:   "records/1.1-init.md",
	}
	index.Tasks["1.2-setup"] = Task{
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
	return indexPath, tasksDir
}

func TestAddTask_Basic(t *testing.T) {
	indexPath, _ := newTestIndex(t)

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

	index, _ := LoadIndex(indexPath)
	task := index.Tasks["disc-1"]
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
	indexPath, _ := newTestIndex(t)

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
	indexPath, _ := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{Title: "First disc"})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "disc-1" {
		t.Errorf("expected disc-1, got %s", id)
	}
}

func TestAddTask_AutoGenerateID_Sequential(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	AddTask(indexPath, AddTaskOpts{Title: "Disc 1"})
	AddTask(indexPath, AddTaskOpts{Title: "Disc 2"})

	id, err := AddTask(indexPath, AddTaskOpts{Title: "Disc 3"})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "disc-3" {
		t.Errorf("expected disc-3, got %s", id)
	}
}

func TestAddTask_AutoGenerateID_GapFill(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	AddTask(indexPath, AddTaskOpts{Title: "Disc 1", ID: "disc-1"})
	AddTask(indexPath, AddTaskOpts{Title: "Disc 3", ID: "disc-3"})

	id, err := AddTask(indexPath, AddTaskOpts{Title: "Fill gap"})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	if id != "disc-2" {
		t.Errorf("expected disc-2 (gap fill), got %s", id)
	}
}

func TestAddTask_DuplicateID(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{ID: "1.1", Title: "Duplicate"})
	if err == nil {
		t.Fatal("expected error for duplicate ID")
	}
}

func TestAddTask_InvalidPriority(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{Title: "Bad priority", Priority: "P5"})
	if err == nil {
		t.Fatal("expected error for invalid priority")
	}
}

func TestAddTask_EmptyTitle(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{Title: ""})
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestAddTask_DefaultPriority(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	id, _ := AddTask(indexPath, AddTaskOpts{Title: "Default prio"})
	index, _ := LoadIndex(indexPath)
	task := index.Tasks[id]
	if task.Priority != "P1" {
		t.Errorf("expected default P1, got %s", task.Priority)
	}
}

func TestAddTask_DependencyNotFound(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	_, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Bad dep",
		Dependencies: []string{"9.9"},
	})
	if err == nil {
		t.Fatal("expected error for missing dependency")
	}
}

func TestAddTask_DependenciesExist(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "With dep",
		Dependencies: []string{"1.1"},
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}
	index, _ := LoadIndex(indexPath)
	task := index.Tasks[id]
	if len(task.Dependencies) != 1 || task.Dependencies[0] != "1.1" {
		t.Errorf("expected deps [1.1], got %v", task.Dependencies)
	}
}

func TestAddTask_Breaking(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	id, _ := AddTask(indexPath, AddTaskOpts{Title: "Breaking task", Breaking: true})
	index, _ := LoadIndex(indexPath)
	task := index.Tasks[id]
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
	if !contains(content, `id: "disc-1"`) {
		t.Errorf("missing id in frontmatter: %s", content)
	}
	if !contains(content, "# disc-1: Fix timeout") {
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

	data, _ := os.ReadFile(filepath.Join(dir, "disc-1.md"))
	content := string(data)
	if !contains(content, "## Steps") {
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

	CreateTaskMarkdown(dir, "disc-1.md", opts)
	data, _ := os.ReadFile(filepath.Join(dir, "disc-1.md"))
	content := string(data)
	if !contains(content, `"1.1"`) || !contains(content, `"1.2"`) {
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

	CreateTaskMarkdown(dir, "disc-1.md", opts)
	data, _ := os.ReadFile(filepath.Join(dir, "disc-1.md"))
	content := string(data)
	if !contains(content, "breaking: true") {
		t.Errorf("missing breaking: %s", content)
	}
}

func TestGenerateDiscID_Empty(t *testing.T) {
	index := NewTaskIndex("test")
	id := generateDiscID(index)
	if id != "disc-1" {
		t.Errorf("expected disc-1, got %s", id)
	}
}

func TestGenerateDiscID_Sequential(t *testing.T) {
	index := NewTaskIndex("test")
	index.Tasks["disc-1"] = Task{ID: "disc-1", Title: "D1", Priority: "P1", Status: "completed", File: "disc-1.md", Record: "records/disc-1.md"}
	index.Tasks["disc-2"] = Task{ID: "disc-2", Title: "D2", Priority: "P1", Status: "completed", File: "disc-2.md", Record: "records/disc-2.md"}
	id := generateDiscID(index)
	if id != "disc-3" {
		t.Errorf("expected disc-3, got %s", id)
	}
}

func TestGenerateDiscID_NonDiscIgnored(t *testing.T) {
	index := NewTaskIndex("test")
	index.Tasks["1.1-init"] = Task{ID: "1.1", Title: "Init", Priority: "P0", Status: "completed", File: "1.1-init.md", Record: "records/1.1-init.md"}
	index.Tasks["fix-e2e-1-1"] = Task{ID: "fix-e2e-1-1", Title: "Fix e2e", Priority: "P0", Status: "completed", File: "fix-e2e-1-1.md", Record: "records/fix-e2e-1-1.md"}
	id := generateDiscID(index)
	if id != "disc-1" {
		t.Errorf("expected disc-1, got %s", id)
	}
}

func TestApplyVars(t *testing.T) {
	tests := []struct {
		name     string
		tmpl     string
		opts     AddTaskOpts
		expected string
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
			name:     "missing var left as-is",
			tmpl:     "keep {{UNKNOWN}} placeholder",
			opts:     AddTaskOpts{},
			expected: "keep {{UNKNOWN}} placeholder",
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
			result := ApplyVars(tt.tmpl, tt.opts)
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
		if !contains(got, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, got)
		}
	}

	if contains(got, "{{") {
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
	indexPath, _ := newTestIndex(t)

	if err := AddDependency(indexPath, "1.2-setup", "disc-1"); err != nil {
		t.Fatalf("AddDependency failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	task := index.Tasks["1.2-setup"]
	if !containsSlice(task.Dependencies, "disc-1") {
		t.Errorf("expected disc-1 in dependencies, got %v", task.Dependencies)
	}
}

func TestAddDependency_Duplicate(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	AddDependency(indexPath, "1.2-setup", "disc-1")
	err := AddDependency(indexPath, "1.2-setup", "disc-1")
	if err != nil {
		t.Errorf("duplicate AddDependency should be no-op, got: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	task := index.Tasks["1.2-setup"]
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
	indexPath, _ := newTestIndex(t)
	err := AddDependency(indexPath, "nonexistent", "disc-1")
	if err == nil {
		t.Fatal("expected error for nonexistent task")
	}
}

func TestGetUnmetDependencies(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	// 1.1 is completed, 1.2 is pending — depend on both
	AddTask(indexPath, AddTaskOpts{Title: "Fix", ID: "fix-1"})
	AddDependency(indexPath, "1.2-setup", "fix-1")

	unmet, err := GetUnmetDependencies(indexPath, "1.2-setup")
	if err != nil {
		t.Fatalf("GetUnmetDependencies failed: %v", err)
	}
	// fix-1 is pending (not completed) → unmet
	if !containsSlice(unmet, "fix-1") {
		t.Errorf("expected fix-1 in unmet, got %v", unmet)
	}

	// Complete fix-1
	index, _ := LoadIndex(indexPath)
	fixTask := index.Tasks["fix-1"]
	fixTask.Status = "completed"
	index.Tasks["fix-1"] = fixTask
	SaveIndex(indexPath, index)

	unmet2, _ := GetUnmetDependencies(indexPath, "1.2-setup")
	if containsSlice(unmet2, "fix-1") {
		t.Errorf("fix-1 is completed, should not be unmet, got %v", unmet2)
	}
}

func containsSlice(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstr(s, substr))
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestAddTask_SourceTaskID_Persisted(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "1.1-init",
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	task := index.Tasks[id]
	if task.SourceTaskID != "1.1-init" {
		t.Errorf("expected sourceTaskID '1.1-init', got %q", task.SourceTaskID)
	}
}

func TestAddTask_SourceTaskID_UpdatesSourceDeps(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "1.1-init",
	})
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	srcTask := index.Tasks["1.1-init"]
	if !containsSlice(srcTask.Dependencies, id) {
		t.Errorf("source task should have %s as dependency, got %v", id, srcTask.Dependencies)
	}
}

func TestAddTask_SourceTaskID_SourceNotFound(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	id, err := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix auth",
		Priority:     "P0",
		SourceTaskID: "nonexistent",
	})
	if err != nil {
		t.Fatalf("AddTask should succeed even if source not found, got: %v", err)
	}

	index, _ := LoadIndex(indexPath)
	task := index.Tasks[id]
	if task.SourceTaskID != "nonexistent" {
		t.Errorf("SourceTaskID should still be persisted, got %q", task.SourceTaskID)
	}
}

func TestAddTask_SourceTaskID_IdempotentDep(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	id1, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 1",
		SourceTaskID: "1.1-init",
	})
	// Add again with same source — source dep should not duplicate
	id2, _ := AddTask(indexPath, AddTaskOpts{
		Title:        "Fix 2",
		SourceTaskID: "1.1-init",
	})

	index, _ := LoadIndex(indexPath)
	srcTask := index.Tasks["1.1-init"]
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

func TestGetUnmetDependencies_Wildcard(t *testing.T) {
	indexPath, _ := newTestIndex(t)

	// Add wildcard dep to 1.2-setup
	AddDependency(indexPath, "1.2-setup", "0.x")

	// Add a phase-0 task that's pending
	AddTask(indexPath, AddTaskOpts{Title: "Phase 0 task", ID: "0.1", Status: "pending"})

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
	indexPath, _ := newTestIndex(t)

	// 1.1 is completed. 1.2 matches wildcard but is the task itself — self-excluded.
	AddDependency(indexPath, "1.2-setup", "1.x")

	unmet, err := GetUnmetDependencies(indexPath, "1.2-setup")
	if err != nil {
		t.Fatalf("GetUnmetDependencies failed: %v", err)
	}
	if len(unmet) != 0 {
		t.Errorf("expected 0 unmet (self-excluded + all others completed), got %v", unmet)
	}
}
