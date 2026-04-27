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
