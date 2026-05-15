package task

import (
	"strings"
	"testing"
)

func TestGetBreakdownTestTasks_SingleProfile(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"go-test"})

	// Shared: gen-cases, eval-cases + per-profile: gen-scripts, run, graduate + shared: verify-regression, consolidate = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	// No suffix for single profile
	wantIDs := []string{"T-test-1", "T-test-1b", "T-test-2", "T-test-3", "T-test-4", "T-test-4.5", "T-test-5"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Dependency chain
	if tasks[1].Dependencies[0] != "T-test-1" {
		t.Errorf("eval-cases should depend on gen-cases, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-test-1b" {
		t.Errorf("gen-scripts should depend on eval-cases, got %v", tasks[2].Dependencies)
	}
	if tasks[3].Dependencies[0] != "T-test-2" {
		t.Errorf("run should depend on gen-scripts, got %v", tasks[3].Dependencies)
	}
	if tasks[4].Dependencies[0] != "T-test-3" {
		t.Errorf("graduate should depend on run, got %v", tasks[4].Dependencies)
	}
	if tasks[5].Dependencies[0] != "T-test-4" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[5].Dependencies)
	}
	if tasks[6].Dependencies[0] != "T-test-4.5" {
		t.Errorf("consolidate should depend on verify-regression, got %v", tasks[6].Dependencies)
	}

	// Per-profile tasks have ProfileName
	if tasks[2].ProfileName != "go-test" {
		t.Errorf("gen-scripts ProfileName = %q, want go-test", tasks[2].ProfileName)
	}
}

func TestGetBreakdownTestTasks_MultiProfile(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"web-playwright", "go-test"})

	// 2 shared + 3*2 per-profile + 2 shared = 10
	if len(tasks) != 10 {
		t.Fatalf("expected 10 tasks, got %d", len(tasks))
	}

	// Profile-suffixed IDs
	wantIDs := []string{
		"T-test-1", "T-test-1b",
		"T-test-2a", "T-test-3a", "T-test-4a",
		"T-test-2b", "T-test-3b", "T-test-4b",
		"T-test-4.5", "T-test-5",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// verify-regression depends on both graduates
	if len(tasks[8].Dependencies) != 2 {
		t.Errorf("verify-regression should depend on 2 graduates, got %v", tasks[8].Dependencies)
	}
}

func TestGetQuickTestTasks_SingleProfile(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"go-test"})

	// 4 per-profile + 2 shared = 6
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-quick-1", "T-quick-2", "T-quick-3", "T-quick-4", "T-quick-5", "T-quick-6"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Chain: 2->1, 3->2, 4->3
	if tasks[1].Dependencies[0] != "T-quick-1" {
		t.Errorf("gen-scripts should depend on gen-cases, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-quick-2" {
		t.Errorf("run should depend on gen-scripts, got %v", tasks[2].Dependencies)
	}

	// T-quick-5 depends on T-quick-4
	if tasks[4].Dependencies[0] != "T-quick-4" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[4].Dependencies)
	}

	// T-quick-6 depends on T-quick-5
	if tasks[5].Dependencies[0] != "T-quick-5" {
		t.Errorf("drift detection should depend on verify-regression, got %v", tasks[5].Dependencies)
	}
	// T-quick-6 type and NoTest
	if tasks[5].Type != TypeDocGenerationDrift {
		t.Errorf("T-quick-6 Type = %q, want %q", tasks[5].Type, TypeDocGenerationDrift)
	}
	if !tasks[5].NoTest {
		t.Error("T-quick-6 NoTest should be true")
	}
	if tasks[5].Scope != "all" {
		t.Errorf("T-quick-6 Scope = %q, want %q", tasks[5].Scope, "all")
	}
}

func TestGetQuickTestTasks_MultiProfile(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"web-playwright", "go-test"})

	// 4*2 per-profile + 2 shared = 10
	if len(tasks) != 10 {
		t.Fatalf("expected 10 tasks, got %d", len(tasks))
	}

	// Profile-suffixed
	if tasks[0].ID != "T-quick-1a" {
		t.Errorf("tasks[0].ID = %q, want T-quick-1a", tasks[0].ID)
	}
	if tasks[4].ID != "T-quick-1b" {
		t.Errorf("tasks[4].ID = %q, want T-quick-1b", tasks[4].ID)
	}

	// T-quick-5 depends on both graduates
	if len(tasks[8].Dependencies) != 2 {
		t.Errorf("verify-regression should depend on 2 graduates, got %v", tasks[8].Dependencies)
	}

	// T-quick-6 depends on T-quick-5
	if tasks[9].Dependencies[0] != "T-quick-5" {
		t.Errorf("drift detection should depend on verify-regression, got %v", tasks[9].Dependencies)
	}
	if tasks[9].Type != TypeDocGenerationDrift {
		t.Errorf("T-quick-6 Type = %q, want %q", tasks[9].Type, TypeDocGenerationDrift)
	}
}

func TestGenerateTestTaskMD(t *testing.T) {
	def := TestTaskDef{
		ID: "T-test-2a", Key: "gen-test-scripts-go-test",
		Title: "Generate Test Scripts (go-test)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{"T-test-1b"},
		Type: TypeTestPipelineGenScripts, Scope: "all",
		ProfileName: "go-test", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	// Check frontmatter
	if !strings.Contains(s, `id: "T-test-2a"`) {
		t.Error("missing id in frontmatter")
	}
	if !strings.Contains(s, `type: "test-pipeline.gen-scripts"`) {
		t.Error("missing type in frontmatter")
	}
	if !strings.Contains(s, `"T-test-1b"`) {
		t.Error("missing dependency in frontmatter")
	}

	// Check profile strategy content loaded
	if !strings.Contains(s, "go-test") {
		t.Error("missing profile name in body")
	}
}

func TestGenerateTestTaskMD_SharedTask(t *testing.T) {
	def := TestTaskDef{
		ID: "T-test-1", Key: "gen-test-cases",
		Title: "Generate Test Cases", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{},
		Type: TypeTestPipelineGenCases, Scope: "all", NoTest: true,
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "noTest: true") {
		t.Error("missing noTest in frontmatter")
	}
}

func TestResolveFirstTestDep(t *testing.T) {
	t.Run("breakdown finds gate", func(t *testing.T) {
		existing := map[string]Task{
			"1-gate":  {ID: "1.gate"},
			"2-gate":  {ID: "2.gate"},
			"1.1-foo": {ID: "1.1"},
		}
		tasks := GetBreakdownTestTasks([]string{"go-test"})
		ResolveFirstTestDep(tasks, existing, "breakdown")
		if tasks[0].Dependencies[0] != "2.gate" {
			t.Errorf("T-test-1 should depend on highest gate, got %v", tasks[0].Dependencies)
		}
	})

	t.Run("quick finds max business task", func(t *testing.T) {
		existing := map[string]Task{
			"1-foo": {ID: "1"},
			"2-bar": {ID: "2"},
			"3-baz": {ID: "3"},
		}
		tasks := GetQuickTestTasks([]string{"go-test"})
		ResolveFirstTestDep(tasks, existing, "quick")
		if tasks[0].Dependencies[0] != "3" {
			t.Errorf("T-quick-1 should depend on max business task, got %v", tasks[0].Dependencies)
		}
	})
}

func TestGetDocEvalTask(t *testing.T) {
	task := GetDocEvalTask()

	if task.ID != "T-eval-doc" {
		t.Errorf("ID = %q, want T-eval-doc", task.ID)
	}
	if task.Key != "eval-doc" {
		t.Errorf("Key = %q, want eval-doc", task.Key)
	}
	if task.Type != TypeDocEvaluation {
		t.Errorf("Type = %q, want %q", task.Type, TypeDocEvaluation)
	}
	if !task.NoTest {
		t.Error("NoTest should be true")
	}
	if task.Title == "" {
		t.Error("Title should not be empty")
	}
	if task.Scope == "" {
		t.Error("Scope should not be empty")
	}
	// Dependencies are resolved later by ResolveDocEvalDep
	if len(task.Dependencies) != 0 {
		t.Errorf("Dependencies should be empty (resolved later), got %v", task.Dependencies)
	}
}

func TestResolveDocEvalDep(t *testing.T) {
	t.Run("sets dependency on last business task", func(t *testing.T) {
		existing := map[string]Task{
			"1-doc":    {ID: "1.1", Type: TypeDocumentation},
			"2-doc":    {ID: "1.2", Type: TypeDocumentation},
			"T-test-1": {ID: "T-test-1", Type: TypeTestPipelineGenCases},
		}
		task := GetDocEvalTask()
		ResolveDocEvalDep(&task, existing)

		if len(task.Dependencies) != 1 {
			t.Fatalf("Dependencies = %v, want exactly 1", task.Dependencies)
		}
		if task.Dependencies[0] != "1.2" {
			t.Errorf("dep = %q, want 1.2", task.Dependencies[0])
		}
	})

	t.Run("empty tasks", func(t *testing.T) {
		existing := map[string]Task{}
		task := GetDocEvalTask()
		ResolveDocEvalDep(&task, existing)

		if len(task.Dependencies) != 0 {
			t.Errorf("Dependencies = %v, want empty for no tasks", task.Dependencies)
		}
	})
}
