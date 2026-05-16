package task

import (
	"strings"
	"testing"
)

func TestDetectTypesFromTestCases(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name: "multiple types with non-zero counts",
			content: `## Summary

| Type | Count |
|------|-------|
| UI   | 5   |
| **Integration** | **2** |
| API  | 3  |
| CLI  | 10  |
| **Total** | **20** |`,
			want: []string{"ui", "integration", "api", "cli"},
		},
		{
			name: "only CLI type present",
			content: `## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 75  |
| **Total** | **75** |`,
			want: []string{"cli"},
		},
		{
			name: "TUI type present",
			content: `## Summary

| Type | Count |
|------|-------|
| TUI  | 0     |
| **Integration** | **0** |
| API  | 0     |
| CLI  | 31    |
| **Total** | **31** |`,
			want: []string{"cli"},
		},
		{
			name: "UI and API only",
			content: `## Summary

| Type | Count |
|------|-------|
| UI   | 8   |
| **Integration** | **0** |
| API  | 4  |
| CLI  | 0  |
| **Total** | **12** |`,
			want: []string{"ui", "api"},
		},
		{
			name:    "empty content returns nil",
			content: ``,
			want:    nil,
		},
		{
			name: "no summary table returns nil",
			content: `# Test Cases

Some text without a table.
`,
			want: nil,
		},
		{
			name: "all zero counts returns nil",
			content: `## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 0  |
| **Total** | **0** |`,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectTypesFromTestCases([]byte(tt.content))
			if len(got) != len(tt.want) {
				t.Fatalf("DetectTypesFromTestCases() = %v, want %v", got, tt.want)
			}
			for i, w := range tt.want {
				if got[i] != w {
					t.Errorf("got[%d] = %q, want %q", i, got[i], w)
				}
			}
		})
	}
}

func TestGetBreakdownTestTasks_SingleProfile(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"go-test"}, nil)

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
	tasks := GetBreakdownTestTasks([]string{"web-playwright", "go-test"}, nil)

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
	tasks := GetQuickTestTasks([]string{"go-test"}, nil)

	// 3 per-profile (gen-cases, gen-and-run, graduate) + 2 shared (verify-regression, drift) = 5
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-quick-1", "T-quick-2", "T-quick-3", "T-quick-4", "T-quick-5"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// T-quick-2 type is gen-and-run
	if tasks[1].Type != TypeTestPipelineGenAndRun {
		t.Errorf("T-quick-2 Type = %q, want %q", tasks[1].Type, TypeTestPipelineGenAndRun)
	}

	// Chain: 2->1, 3->2
	if tasks[1].Dependencies[0] != "T-quick-1" {
		t.Errorf("gen-and-run should depend on gen-cases, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-quick-2" {
		t.Errorf("graduate should depend on gen-and-run, got %v", tasks[2].Dependencies)
	}

	// T-quick-4 depends on T-quick-3
	if tasks[3].Dependencies[0] != "T-quick-3" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[3].Dependencies)
	}

	// T-quick-5 depends on T-quick-4
	if tasks[4].Dependencies[0] != "T-quick-4" {
		t.Errorf("drift detection should depend on verify-regression, got %v", tasks[4].Dependencies)
	}
	// T-quick-5 type and NoTest
	if tasks[4].Type != TypeDocGenerationDrift {
		t.Errorf("T-quick-5 Type = %q, want %q", tasks[4].Type, TypeDocGenerationDrift)
	}
	if !tasks[4].NoTest {
		t.Error("T-quick-5 NoTest should be true")
	}
	if tasks[4].Scope != "all" {
		t.Errorf("T-quick-5 Scope = %q, want %q", tasks[4].Scope, "all")
	}
}

func TestGetQuickTestTasks_MultiProfile(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"web-playwright", "go-test"}, nil)

	// 3*2 per-profile + 2 shared = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
	}

	// Profile-suffixed
	if tasks[0].ID != "T-quick-1a" {
		t.Errorf("tasks[0].ID = %q, want T-quick-1a", tasks[0].ID)
	}
	if tasks[3].ID != "T-quick-1b" {
		t.Errorf("tasks[3].ID = %q, want T-quick-1b", tasks[3].ID)
	}

	// T-quick-4 depends on both graduates
	if len(tasks[6].Dependencies) != 2 {
		t.Errorf("verify-regression should depend on 2 graduates, got %v", tasks[6].Dependencies)
	}

	// T-quick-5 depends on T-quick-4
	if tasks[7].Dependencies[0] != "T-quick-4" {
		t.Errorf("drift detection should depend on verify-regression, got %v", tasks[7].Dependencies)
	}
	if tasks[7].Type != TypeDocGenerationDrift {
		t.Errorf("T-quick-5 Type = %q, want %q", tasks[7].Type, TypeDocGenerationDrift)
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
		tasks := GetBreakdownTestTasks([]string{"go-test"}, nil)
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
		tasks := GetQuickTestTasks([]string{"go-test"}, nil)
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

func TestGetBreakdownTestTasks_PerType_SingleProfile(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"tui", "api"})

	// Shared: gen-cases, eval-cases + per-type-gen: 2 (tui, api) + run + graduate + verify-regression + consolidate = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-1", "T-test-1b",
		"T-test-2-tui", "T-test-2-api",
		"T-test-3", "T-test-4",
		"T-test-4.5", "T-test-5",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Keys include type suffix
	if tasks[2].Key != "gen-test-scripts-go-test-tui" {
		t.Errorf("tasks[2].Key = %q, want gen-test-scripts-go-test-tui", tasks[2].Key)
	}
	if tasks[3].Key != "gen-test-scripts-go-test-api" {
		t.Errorf("tasks[3].Key = %q, want gen-test-scripts-go-test-api", tasks[3].Key)
	}

	// TestType field set
	if tasks[2].TestType != "tui" {
		t.Errorf("tasks[2].TestType = %q, want tui", tasks[2].TestType)
	}
	if tasks[3].TestType != "api" {
		t.Errorf("tasks[3].TestType = %q, want api", tasks[3].TestType)
	}

	// T-test-3 depends on ALL per-type T-test-2-* tasks
	if len(tasks[4].Dependencies) != 2 {
		t.Fatalf("T-test-3 should depend on 2 gen tasks, got %v", tasks[4].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[4].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-2-tui"] || !depSet["T-test-2-api"] {
		t.Errorf("T-test-3 deps should include T-test-2-tui and T-test-2-api, got %v", tasks[4].Dependencies)
	}

	// T-test-4 depends on T-test-3
	if tasks[5].Dependencies[0] != "T-test-3" {
		t.Errorf("graduate should depend on run, got %v", tasks[5].Dependencies)
	}

	// Per-type gen tasks depend on T-test-1b
	if tasks[2].Dependencies[0] != "T-test-1b" {
		t.Errorf("T-test-2-tui should depend on T-test-1b, got %v", tasks[2].Dependencies)
	}
	if tasks[3].Dependencies[0] != "T-test-1b" {
		t.Errorf("T-test-2-api should depend on T-test-1b, got %v", tasks[3].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_MultiProfile(t *testing.T) {
	tasks := GetBreakdownTestTasks(
		[]string{"web-playwright", "go-test"},
		[]string{"tui", "api"},
	)

	// Shared: 2 + profile-a: 2 per-type-gen + run + grad = 4 + profile-b: same = 4 + shared: 2 = 12
	if len(tasks) != 12 {
		t.Fatalf("expected 12 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-1", "T-test-1b",
		"T-test-2a-tui", "T-test-2a-api",
		"T-test-3a", "T-test-4a",
		"T-test-2b-tui", "T-test-2b-api",
		"T-test-3b", "T-test-4b",
		"T-test-4.5", "T-test-5",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// T-test-3a depends on both profile-a gen tasks
	if len(tasks[4].Dependencies) != 2 {
		t.Fatalf("T-test-3a should depend on 2 gen tasks, got %v", tasks[4].Dependencies)
	}
	depSetA := make(map[string]bool)
	for _, d := range tasks[4].Dependencies {
		depSetA[d] = true
	}
	if !depSetA["T-test-2a-tui"] || !depSetA["T-test-2a-api"] {
		t.Errorf("T-test-3a deps should include T-test-2a-tui and T-test-2a-api, got %v", tasks[4].Dependencies)
	}

	// Keys include profile and type
	if tasks[2].Key != "gen-test-scripts-web-playwright-tui" {
		t.Errorf("tasks[2].Key = %q, want gen-test-scripts-web-playwright-tui", tasks[2].Key)
	}
	if tasks[6].Key != "gen-test-scripts-go-test-tui" {
		t.Errorf("tasks[6].Key = %q, want gen-test-scripts-go-test-tui", tasks[6].Key)
	}

	// verify-regression depends on both graduates
	if len(tasks[10].Dependencies) != 2 {
		t.Errorf("verify-regression should depend on 2 graduates, got %v", tasks[10].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"api"})

	// Only api type -> one gen task
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-1", "T-test-1b",
		"T-test-2-api",
		"T-test-3", "T-test-4",
		"T-test-4.5", "T-test-5",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// T-test-3 depends on single gen task
	if len(tasks[3].Dependencies) != 1 || tasks[3].Dependencies[0] != "T-test-2-api" {
		t.Errorf("T-test-3 should depend on T-test-2-api, got %v", tasks[3].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"tui", "api", "cli"})

	// 3 types -> 3 gen tasks
	if len(tasks) != 9 {
		t.Fatalf("expected 9 tasks, got %d", len(tasks))
	}

	// T-test-3 depends on all 3 gen tasks
	if len(tasks[5].Dependencies) != 3 {
		t.Fatalf("T-test-3 should depend on 3 gen tasks, got %v", tasks[5].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[5].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-2-tui"] || !depSet["T-test-2-api"] || !depSet["T-test-2-cli"] {
		t.Errorf("T-test-3 missing expected deps, got %v", tasks[5].Dependencies)
	}
}

func TestGenerateTestTaskMD_WithTestType(t *testing.T) {
	def := TestTaskDef{
		ID: "T-test-2-api", Key: "gen-test-scripts-go-test-api",
		Title: "Generate Test Scripts (go-test, api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{"T-test-1b"},
		Type: TypeTestPipelineGenScripts, Scope: "all",
		ProfileName: "go-test", TestType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	// Check frontmatter has profile
	if !strings.Contains(s, `profile: "go-test"`) {
		t.Error("missing profile in frontmatter")
	}
	// Check body mentions type
	if !strings.Contains(s, "api") {
		t.Error("missing test type in body")
	}
	if !strings.Contains(s, "go-test") {
		t.Error("missing profile name in body")
	}
}

func TestGetQuickTestTasks_PerType_SingleProfile(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"tui", "api"})

	// Per-profile: gen-cases + per-type-gen-and-run(tui,api) + graduate = 4 + shared verify-regression + drift-detection = 6
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-1",
		"T-quick-2-tui", "T-quick-2-api",
		"T-quick-3",
		"T-quick-4",
		"T-quick-5",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// T-quick-2-tui and T-quick-2-api are gen-and-run type
	if tasks[1].Type != TypeTestPipelineGenAndRun {
		t.Errorf("T-quick-2-tui Type = %q, want %q", tasks[1].Type, TypeTestPipelineGenAndRun)
	}
	if tasks[2].Type != TypeTestPipelineGenAndRun {
		t.Errorf("T-quick-2-api Type = %q, want %q", tasks[2].Type, TypeTestPipelineGenAndRun)
	}

	// Keys include type suffix
	if tasks[1].Key != "quick-gen-and-run-go-test-tui" {
		t.Errorf("tasks[1].Key = %q, want quick-gen-and-run-go-test-tui", tasks[1].Key)
	}
	if tasks[2].Key != "quick-gen-and-run-go-test-api" {
		t.Errorf("tasks[2].Key = %q, want quick-gen-and-run-go-test-api", tasks[2].Key)
	}

	// TestType field set
	if tasks[1].TestType != "tui" {
		t.Errorf("tasks[1].TestType = %q, want tui", tasks[1].TestType)
	}
	if tasks[2].TestType != "api" {
		t.Errorf("tasks[2].TestType = %q, want api", tasks[2].TestType)
	}

	// Per-type gen-and-run depend on gen-cases
	if tasks[1].Dependencies[0] != "T-quick-1" {
		t.Errorf("T-quick-2-tui should depend on T-quick-1, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-quick-1" {
		t.Errorf("T-quick-2-api should depend on T-quick-1, got %v", tasks[2].Dependencies)
	}

	// T-quick-3 depends on ALL per-type T-quick-2-* tasks
	if len(tasks[3].Dependencies) != 2 {
		t.Fatalf("T-quick-3 should depend on 2 gen-and-run tasks, got %v", tasks[3].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[3].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-quick-2-tui"] || !depSet["T-quick-2-api"] {
		t.Errorf("T-quick-3 deps should include T-quick-2-tui and T-quick-2-api, got %v", tasks[3].Dependencies)
	}

	// T-quick-4 depends on T-quick-3
	if tasks[4].Dependencies[0] != "T-quick-3" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[4].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_MultiProfile(t *testing.T) {
	tasks := GetQuickTestTasks(
		[]string{"web-playwright", "go-test"},
		[]string{"tui", "api"},
	)

	// Profile-a: gen-cases + 2 per-type-gen-and-run + graduate = 4
	// Profile-b: same = 4
	// Shared: verify-regression + drift-detection = 2
	// Total = 10
	if len(tasks) != 10 {
		t.Fatalf("expected 10 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-1a",
		"T-quick-2a-tui", "T-quick-2a-api",
		"T-quick-3a",
		"T-quick-1b",
		"T-quick-2b-tui", "T-quick-2b-api",
		"T-quick-3b",
		"T-quick-4",
		"T-quick-5",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// T-quick-3a depends on both profile-a gen-and-run tasks
	if len(tasks[3].Dependencies) != 2 {
		t.Fatalf("T-quick-3a should depend on 2 gen-and-run tasks, got %v", tasks[3].Dependencies)
	}
	depSetA := make(map[string]bool)
	for _, d := range tasks[3].Dependencies {
		depSetA[d] = true
	}
	if !depSetA["T-quick-2a-tui"] || !depSetA["T-quick-2a-api"] {
		t.Errorf("T-quick-3a deps should include T-quick-2a-tui and T-quick-2a-api, got %v", tasks[3].Dependencies)
	}

	// Keys include profile and type
	if tasks[1].Key != "quick-gen-and-run-web-playwright-tui" {
		t.Errorf("tasks[1].Key = %q, want quick-gen-and-run-web-playwright-tui", tasks[1].Key)
	}
	if tasks[5].Key != "quick-gen-and-run-go-test-tui" {
		t.Errorf("tasks[5].Key = %q, want quick-gen-and-run-go-test-tui", tasks[5].Key)
	}

	// T-quick-4 depends on both graduates
	if len(tasks[8].Dependencies) != 2 {
		t.Errorf("verify-regression should depend on 2 graduates, got %v", tasks[8].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"api"})

	// Only api type -> one gen-and-run task
	// gen-cases + 1 gen-and-run-api + graduate + verify-regression + drift-detection = 5
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-1",
		"T-quick-2-api",
		"T-quick-3",
		"T-quick-4",
		"T-quick-5",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// T-quick-2-api type is gen-and-run
	if tasks[1].Type != TypeTestPipelineGenAndRun {
		t.Errorf("T-quick-2-api Type = %q, want %q", tasks[1].Type, TypeTestPipelineGenAndRun)
	}

	// T-quick-3 depends on single gen-and-run task
	if len(tasks[2].Dependencies) != 1 || tasks[2].Dependencies[0] != "T-quick-2-api" {
		t.Errorf("T-quick-3 should depend on T-quick-2-api, got %v", tasks[2].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"tui", "api", "cli"})

	// gen-cases + 3 per-type-gen-and-run + graduate + verify-regression + drift-detection = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	// T-quick-3 depends on all 3 gen-and-run tasks
	if len(tasks[4].Dependencies) != 3 {
		t.Fatalf("T-quick-3 should depend on 3 gen-and-run tasks, got %v", tasks[4].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[4].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-quick-2-tui"] || !depSet["T-quick-2-api"] || !depSet["T-quick-2-cli"] {
		t.Errorf("T-quick-3 missing expected deps, got %v", tasks[4].Dependencies)
	}
}
