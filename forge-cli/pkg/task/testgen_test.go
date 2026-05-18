package task

import (
	"strings"
	"testing"

	"forge-cli/pkg/profile"
)

// defaultAuto is the current default (consolidateSpecs quick=true, e2eTest quick=false).
var defaultAuto = profile.AutoConfigDefaults()

// allEnabledAuto enables all auto-behaviors for tests that need quick + full tasks.
var allEnabledAuto = profile.AutoConfig{
	E2eTest:          profile.ModeToggle{Quick: true, Full: true},
	ConsolidateSpecs: profile.ModeToggle{Quick: true, Full: true},
	CleanCode:        profile.ModeToggle{Quick: false, Full: false},
}

func TestGetBreakdownTestTasks_EmptyInterfaces(t *testing.T) {
	tasks := GetBreakdownTestTasks([]profile.Language{"go"}, nil, defaultAuto)

	// No capabilities -> no test tasks generated
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks with empty capabilities, got %d", len(tasks))
	}
}

func TestGetBreakdownTestTasks_SingleProfile(t *testing.T) {
	tasks := GetBreakdownTestTasks([]profile.Language{"go"}, []string{"cli"}, defaultAuto)

	// Shared: gen-cases, eval-cases + per-type: gen-scripts-cli, run, graduate + shared: verify-regression, consolidate = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	// No suffix for single profile
	wantIDs := []string{"T-test-1", "T-test-1b", "T-test-2-cli", "T-test-3", "T-test-4", "T-test-4.5", "T-specs-1"}
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
	if tasks[3].Dependencies[0] != "T-test-2-cli" {
		t.Errorf("run should depend on gen-scripts-cli, got %v", tasks[3].Dependencies)
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

	// Per-language tasks have Language
	if tasks[2].Language != "go" {
		t.Errorf("gen-scripts Language = %q, want go", tasks[2].Language)
	}
}

func TestGetBreakdownTestTasks_MultiProfile(t *testing.T) {
	tasks := GetBreakdownTestTasks([]profile.Language{"javascript", "go"}, []string{"api"}, defaultAuto)

	// 2 shared + (1 per-type-gen + run + graduate)*2 + 2 shared = 10
	if len(tasks) != 10 {
		t.Fatalf("expected 10 tasks, got %d", len(tasks))
	}

	// Profile-suffixed IDs
	wantIDs := []string{
		"T-test-1", "T-test-1b",
		"T-test-2a-api", "T-test-3a", "T-test-4a",
		"T-test-2b-api", "T-test-3b", "T-test-4b",
		"T-test-4.5", "T-specs-1",
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

func TestGetQuickTestTasks_EmptyInterfaces(t *testing.T) {
	tasks := GetQuickTestTasks([]profile.Language{"go"}, nil, allEnabledAuto)

	// No capabilities -> no test tasks
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks with empty capabilities, got %d", len(tasks))
	}
}

func TestGetQuickTestTasks_SingleProfile(t *testing.T) {
	tasks := GetQuickTestTasks([]profile.Language{"go"}, []string{"cli"}, allEnabledAuto)

	// gen-cases + gen-and-run-cli + graduate + verify-regression + drift = 5
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-quick-1", "T-quick-2-cli", "T-quick-3", "T-quick-4", "T-quick-specs-1"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// T-quick-2-cli type is gen-and-run
	if tasks[1].Type != TypeTestPipelineGenAndRun {
		t.Errorf("T-quick-2-cli Type = %q, want %q", tasks[1].Type, TypeTestPipelineGenAndRun)
	}

	// Chain: 2->1, 3->2
	if tasks[1].Dependencies[0] != "T-quick-1" {
		t.Errorf("gen-and-run should depend on gen-cases, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-quick-2-cli" {
		t.Errorf("graduate should depend on gen-and-run, got %v", tasks[2].Dependencies)
	}

	// T-quick-4 depends on T-quick-3
	if tasks[3].Dependencies[0] != "T-quick-3" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[3].Dependencies)
	}

	// T-quick-specs-1 depends on T-quick-4
	if tasks[4].Dependencies[0] != "T-quick-4" {
		t.Errorf("drift detection should depend on verify-regression, got %v", tasks[4].Dependencies)
	}
	// T-quick-specs-1 type and NoTest
	if tasks[4].Type != TypeDocGenerationDrift {
		t.Errorf("T-quick-specs-1 Type = %q, want %q", tasks[4].Type, TypeDocGenerationDrift)
	}
	if !tasks[4].NoTest {
		t.Error("T-quick-specs-1 NoTest should be true")
	}
	if tasks[4].Scope != "all" {
		t.Errorf("T-quick-specs-1 Scope = %q, want %q", tasks[4].Scope, "all")
	}
}

func TestGetQuickTestTasks_MultiProfile(t *testing.T) {
	tasks := GetQuickTestTasks([]profile.Language{"javascript", "go"}, []string{"api"}, allEnabledAuto)

	// Profile-a: gen-cases + gen-and-run-api + graduate = 3
	// Profile-b: same = 3
	// Shared: verify-regression + drift = 2
	// Total = 8
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

	// T-quick-specs-1 depends on T-quick-4
	if tasks[7].Dependencies[0] != "T-quick-4" {
		t.Errorf("drift detection should depend on verify-regression, got %v", tasks[7].Dependencies)
	}
	if tasks[7].Type != TypeDocGenerationDrift {
		t.Errorf("T-quick-specs-1 Type = %q, want %q", tasks[7].Type, TypeDocGenerationDrift)
	}
}

func TestGenerateTestTaskMD(t *testing.T) {
	def := TestTaskDef{
		ID: "T-test-2a-api", Key: "gen-test-scripts-go-api",
		Title: "Generate Test Scripts (go, api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{"T-test-1b"},
		Type: TypeTestPipelineGenScripts, Scope: "all",
		Language: "go", TestType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	// Check frontmatter
	if !strings.Contains(s, `id: "T-test-2a-api"`) {
		t.Error("missing id in frontmatter")
	}
	if !strings.Contains(s, `type: "test-pipeline.gen-scripts"`) {
		t.Error("missing type in frontmatter")
	}
	if !strings.Contains(s, `"T-test-1b"`) {
		t.Error("missing dependency in frontmatter")
	}

	// Check profile strategy content loaded
	if !strings.Contains(s, "go") {
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
		tasks := GetBreakdownTestTasks([]profile.Language{"go"}, []string{"cli"}, defaultAuto)
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
		tasks := GetQuickTestTasks([]profile.Language{"go"}, []string{"cli"}, allEnabledAuto)
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
	tasks := GetBreakdownTestTasks([]profile.Language{"go"}, []string{"tui", "api"}, defaultAuto)

	// Shared: gen-cases, eval-cases + per-type-gen: 2 (tui, api) + run + graduate + verify-regression + consolidate = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-1", "T-test-1b",
		"T-test-2-tui", "T-test-2-api",
		"T-test-3", "T-test-4",
		"T-test-4.5", "T-specs-1",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Keys include type suffix
	if tasks[2].Key != "gen-test-scripts-go-tui" {
		t.Errorf("tasks[2].Key = %q, want gen-test-scripts-go-tui", tasks[2].Key)
	}
	if tasks[3].Key != "gen-test-scripts-go-api" {
		t.Errorf("tasks[3].Key = %q, want gen-test-scripts-go-api", tasks[3].Key)
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
		[]profile.Language{"javascript", "go"},
		[]string{"tui", "api"},
		defaultAuto,
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
		"T-test-4.5", "T-specs-1",
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
	if tasks[2].Key != "gen-test-scripts-javascript-tui" {
		t.Errorf("tasks[2].Key = %q, want gen-test-scripts-javascript-tui", tasks[2].Key)
	}
	if tasks[6].Key != "gen-test-scripts-go-tui" {
		t.Errorf("tasks[6].Key = %q, want gen-test-scripts-go-tui", tasks[6].Key)
	}

	// verify-regression depends on both graduates
	if len(tasks[10].Dependencies) != 2 {
		t.Errorf("verify-regression should depend on 2 graduates, got %v", tasks[10].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]profile.Language{"go"}, []string{"api"}, defaultAuto)

	// Only api type -> one gen task
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-1", "T-test-1b",
		"T-test-2-api",
		"T-test-3", "T-test-4",
		"T-test-4.5", "T-specs-1",
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
	tasks := GetBreakdownTestTasks([]profile.Language{"go"}, []string{"tui", "api", "cli"}, defaultAuto)

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
		ID: "T-test-2-api", Key: "gen-test-scripts-go-api",
		Title: "Generate Test Scripts (go, api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{"T-test-1b"},
		Type: TypeTestPipelineGenScripts, Scope: "all",
		Language: "go", TestType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	// Check frontmatter has profile
	if !strings.Contains(s, `profile: "go"`) {
		t.Error("missing profile in frontmatter")
	}
	// Check body mentions type
	if !strings.Contains(s, "api") {
		t.Error("missing test type in body")
	}
	if !strings.Contains(s, "go") {
		t.Error("missing profile name in body")
	}
}

func TestGetQuickTestTasks_PerType_SingleProfile(t *testing.T) {
	tasks := GetQuickTestTasks([]profile.Language{"go"}, []string{"tui", "api"}, allEnabledAuto)

	// Per-profile: gen-cases + per-type-gen-and-run(tui,api) + graduate = 4 + shared verify-regression + drift-detection = 6
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-1",
		"T-quick-2-tui", "T-quick-2-api",
		"T-quick-3",
		"T-quick-4",
		"T-quick-specs-1",
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
	if tasks[1].Key != "quick-gen-and-run-go-tui" {
		t.Errorf("tasks[1].Key = %q, want quick-gen-and-run-go-tui", tasks[1].Key)
	}
	if tasks[2].Key != "quick-gen-and-run-go-api" {
		t.Errorf("tasks[2].Key = %q, want quick-gen-and-run-go-api", tasks[2].Key)
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
		[]profile.Language{"javascript", "go"},
		[]string{"tui", "api"},
		allEnabledAuto,
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
		"T-quick-specs-1",
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
	if tasks[1].Key != "quick-gen-and-run-javascript-tui" {
		t.Errorf("tasks[1].Key = %q, want quick-gen-and-run-javascript-tui", tasks[1].Key)
	}
	if tasks[5].Key != "quick-gen-and-run-go-tui" {
		t.Errorf("tasks[5].Key = %q, want quick-gen-and-run-go-tui", tasks[5].Key)
	}

	// T-quick-4 depends on both graduates
	if len(tasks[8].Dependencies) != 2 {
		t.Errorf("verify-regression should depend on 2 graduates, got %v", tasks[8].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks([]profile.Language{"go"}, []string{"api"}, allEnabledAuto)

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
		"T-quick-specs-1",
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
	tasks := GetQuickTestTasks([]profile.Language{"go"}, []string{"tui", "api", "cli"}, allEnabledAuto)

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

func TestGetQuickTestTasks_DefaultAuto_IncludesSpecDrift(t *testing.T) {
	// Verify that default auto config (ConsolidateSpecs.Quick=true) generates T-quick-specs-1.
	// E2eTest.Quick is false by default, so no e2e tasks — only the drift task.
	tasks := GetQuickTestTasks([]profile.Language{"go"}, []string{"cli"}, defaultAuto)

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task (spec drift only), got %d", len(tasks))
	}

	if tasks[0].ID != "T-quick-specs-1" {
		t.Errorf("task ID = %q, want T-quick-specs-1", tasks[0].ID)
	}
	if tasks[0].Type != TypeDocGenerationDrift {
		t.Errorf("task Type = %q, want %q", tasks[0].Type, TypeDocGenerationDrift)
	}
	if !tasks[0].NoTest {
		t.Error("T-quick-specs-1 NoTest should be true")
	}
	// No e2e dependency since E2eTest.Quick is false
	if len(tasks[0].Dependencies) != 0 {
		t.Errorf("T-quick-specs-1 should have no deps without e2e tasks, got %v", tasks[0].Dependencies)
	}
}
