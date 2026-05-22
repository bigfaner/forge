package task

import (
	"strings"
	"testing"

	"forge-cli/pkg/forgeconfig"
)

// defaultAuto is the current default (consolidateSpecs quick=true, e2eTest quick=false).
var defaultAuto = forgeconfig.AutoConfigDefaults()

// allEnabledAuto enables all auto-behaviors for tests that need quick + full tasks.
var allEnabledAuto = forgeconfig.AutoConfig{
	E2eTest:          forgeconfig.ModeToggle{Quick: true, Full: true},
	ConsolidateSpecs: forgeconfig.ModeToggle{Quick: true, Full: true},
	CleanCode:        forgeconfig.ModeToggle{Quick: false, Full: false},
}

// validationAuto enables validation + e2e for testing validate-ux gating.
var validationAuto = forgeconfig.AutoConfig{
	E2eTest:          forgeconfig.ModeToggle{Quick: true, Full: true},
	Validation:       forgeconfig.ModeToggle{Quick: true, Full: true},
}

func TestGetBreakdownTestTasks_EmptyInterfaces(t *testing.T) {
	tasks := GetBreakdownTestTasks(nil, defaultAuto)

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks with empty interfaces, got %d", len(tasks))
	}
}

func TestGetBreakdownTestTasks_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)

	// Shared: gen-cases, eval-cases + per-type: gen-scripts-cli + run + graduate + shared: verify-regression, consolidate = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-test-gen-cases", "T-test-eval-cases", "T-test-gen-scripts-cli", "T-test-run", "T-test-graduate", "T-test-verify-regression", "T-specs-consolidate"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Dependency chain
	if tasks[1].Dependencies[0] != "T-test-gen-cases" {
		t.Errorf("eval-cases should depend on gen-cases, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-test-eval-cases" {
		t.Errorf("gen-scripts should depend on eval-cases, got %v", tasks[2].Dependencies)
	}
	if tasks[3].Dependencies[0] != "T-test-gen-scripts-cli" {
		t.Errorf("run should depend on gen-scripts-cli, got %v", tasks[3].Dependencies)
	}
	if tasks[4].Dependencies[0] != "T-test-run" {
		t.Errorf("graduate should depend on run, got %v", tasks[4].Dependencies)
	}
	if tasks[5].Dependencies[0] != "T-test-graduate" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[5].Dependencies)
	}
	if tasks[6].Dependencies[0] != "T-test-verify-regression" {
		t.Errorf("consolidate should depend on verify-regression, got %v", tasks[6].Dependencies)
	}
}

func TestGetQuickTestTasks_EmptyInterfaces(t *testing.T) {
	tasks := GetQuickTestTasks(nil, allEnabledAuto)

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks with empty interfaces, got %d", len(tasks))
	}
}

func TestGetQuickTestTasks_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli"}, allEnabledAuto)

	// gen-cases + gen-and-run-cli + graduate + verify-regression + drift = 5
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-quick-gen-cases", "T-quick-gen-and-run-cli", "T-quick-graduate", "T-quick-verify-regression", "T-quick-doc-drift"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if tasks[1].Type != TypeTestGenAndRun {
		t.Errorf("T-quick-gen-and-run-cli Type = %q, want %q", tasks[1].Type, TypeTestGenAndRun)
	}

	if tasks[1].Dependencies[0] != "T-quick-gen-cases" {
		t.Errorf("gen-and-run should depend on gen-cases, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-quick-gen-and-run-cli" {
		t.Errorf("graduate should depend on gen-and-run, got %v", tasks[2].Dependencies)
	}
	if tasks[3].Dependencies[0] != "T-quick-graduate" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[3].Dependencies)
	}
	if tasks[4].Dependencies[0] != "T-quick-verify-regression" {
		t.Errorf("drift detection should depend on verify-regression, got %v", tasks[4].Dependencies)
	}
	if tasks[4].Type != TypeDocDrift {
		t.Errorf("T-quick-doc-drift Type = %q, want %q", tasks[4].Type, TypeDocDrift)
	}
	if tasks[4].Scope != "all" {
		t.Errorf("T-quick-doc-drift Scope = %q, want %q", tasks[4].Scope, "all")
	}
}

func TestGenerateTestTaskMD(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{"T-test-eval-cases"},
		Type: TypeTestGenScripts, Scope: "all",
		TestType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	if !strings.Contains(s, `id: "T-test-gen-scripts-api"`) {
		t.Error("missing id in frontmatter")
	}
	if !strings.Contains(s, `type: "test.gen-scripts"`) {
		t.Error("missing type in frontmatter")
	}
	if !strings.Contains(s, `"T-test-eval-cases"`) {
		t.Error("missing dependency in frontmatter")
	}
	if !strings.Contains(s, "docs/conventions/testing-") {
		t.Error("body should reference docs/conventions/testing-*")
	}
}

func TestGenerateTestTaskMD_SharedTask(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-cases", Key: "gen-test-cases",
		Title: "Generate Test Cases", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{},
		Type: TypeTestGenCases, Scope: "all",
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	if strings.Contains(s, "noTest") {
		t.Error("noTest should not appear in frontmatter")
	}
}

func TestResolveFirstTestDep(t *testing.T) {
	t.Run("breakdown finds gate", func(t *testing.T) {
		existing := map[string]Task{
			"1-gate":  {ID: "1.gate"},
			"2-gate":  {ID: "2.gate"},
			"1.1-foo": {ID: "1.1"},
		}
		tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)
		ResolveFirstTestDep(tasks, existing, "breakdown")
		if tasks[0].Dependencies[0] != "2.gate" {
			t.Errorf("T-test-gen-cases should depend on highest gate, got %v", tasks[0].Dependencies)
		}
	})

	t.Run("breakdown picks last business task over lower gate", func(t *testing.T) {
		// Phase 3 has 1 task (no gate generated), phase 2 has a gate.
		// Test chain must depend on the last business task (3.1), not the highest gate (2.gate).
		existing := map[string]Task{
			"1-gate":  {ID: "1.gate"},
			"2-gate":  {ID: "2.gate"},
			"1.1-foo": {ID: "1.1"},
			"2.1-bar": {ID: "2.1"},
			"3.1-baz": {ID: "3.1"},
		}
		tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)
		ResolveFirstTestDep(tasks, existing, "breakdown")
		if tasks[0].Dependencies[0] != "3.1" {
			t.Errorf("T-test-gen-cases should depend on last business task (3.1) when it's in a higher phase than highest gate (2.gate), got %v", tasks[0].Dependencies)
		}
	})

	t.Run("quick finds max business task", func(t *testing.T) {
		existing := map[string]Task{
			"1-foo": {ID: "1"},
			"2-bar": {ID: "2"},
			"3-baz": {ID: "3"},
		}
		tasks := GetQuickTestTasks([]string{"cli"}, allEnabledAuto)
		ResolveFirstTestDep(tasks, existing, "quick")
		if tasks[0].Dependencies[0] != "3" {
			t.Errorf("T-quick-gen-cases should depend on max business task, got %v", tasks[0].Dependencies)
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
	if task.Type != TypeDocEval {
		t.Errorf("Type = %q, want %q", task.Type, TypeDocEval)
	}
	if task.Title == "" {
		t.Error("Title should not be empty")
	}
	if task.Scope == "" {
		t.Error("Scope should not be empty")
	}
	if len(task.Dependencies) != 0 {
		t.Errorf("Dependencies should be empty (resolved later), got %v", task.Dependencies)
	}
}

func TestResolveDocEvalDep(t *testing.T) {
	t.Run("sets dependency on last business task", func(t *testing.T) {
		existing := map[string]Task{
			"1-doc":            {ID: "1.1", Type: TypeDoc},
			"2-doc":            {ID: "1.2", Type: TypeDoc},
			"T-test-gen-cases": {ID: "T-test-gen-cases", Type: TypeTestGenCases},
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

func TestGetBreakdownTestTasks_PerType_TwoTypes(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"tui", "api"}, defaultAuto)

	// Shared: gen-cases, eval-cases + per-type-gen: 2 (tui, api) + run + graduate + verify-regression + consolidate = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-cases", "T-test-eval-cases",
		"T-test-gen-scripts-tui", "T-test-gen-scripts-api",
		"T-test-run", "T-test-graduate",
		"T-test-verify-regression", "T-specs-consolidate",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Keys include type suffix (no language)
	if tasks[2].Key != "gen-test-scripts-tui" {
		t.Errorf("tasks[2].Key = %q, want gen-test-scripts-tui", tasks[2].Key)
	}
	if tasks[3].Key != "gen-test-scripts-api" {
		t.Errorf("tasks[3].Key = %q, want gen-test-scripts-api", tasks[3].Key)
	}

	// TestType field set
	if tasks[2].TestType != "tui" {
		t.Errorf("tasks[2].TestType = %q, want tui", tasks[2].TestType)
	}
	if tasks[3].TestType != "api" {
		t.Errorf("tasks[3].TestType = %q, want api", tasks[3].TestType)
	}

	// T-test-run depends on ALL per-type gen-scripts tasks
	if len(tasks[4].Dependencies) != 2 {
		t.Fatalf("T-test-run should depend on 2 gen tasks, got %v", tasks[4].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[4].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-gen-scripts-tui"] || !depSet["T-test-gen-scripts-api"] {
		t.Errorf("T-test-run deps should include T-test-gen-scripts-tui and T-test-gen-scripts-api, got %v", tasks[4].Dependencies)
	}

	// T-test-graduate depends on T-test-run
	if tasks[5].Dependencies[0] != "T-test-run" {
		t.Errorf("graduate should depend on run, got %v", tasks[5].Dependencies)
	}

	// Per-type gen tasks depend on T-test-eval-cases
	if tasks[2].Dependencies[0] != "T-test-eval-cases" {
		t.Errorf("T-test-gen-scripts-tui should depend on T-test-eval-cases, got %v", tasks[2].Dependencies)
	}
	if tasks[3].Dependencies[0] != "T-test-eval-cases" {
		t.Errorf("T-test-gen-scripts-api should depend on T-test-eval-cases, got %v", tasks[3].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"api"}, defaultAuto)

	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-cases", "T-test-eval-cases",
		"T-test-gen-scripts-api",
		"T-test-run", "T-test-graduate",
		"T-test-verify-regression", "T-specs-consolidate",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if len(tasks[3].Dependencies) != 1 || tasks[3].Dependencies[0] != "T-test-gen-scripts-api" {
		t.Errorf("T-test-run should depend on T-test-gen-scripts-api, got %v", tasks[3].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"tui", "api", "cli"}, defaultAuto)

	// 3 types -> 3 gen tasks
	if len(tasks) != 9 {
		t.Fatalf("expected 9 tasks, got %d", len(tasks))
	}

	// T-test-run depends on all 3 gen tasks
	if len(tasks[5].Dependencies) != 3 {
		t.Fatalf("T-test-run should depend on 3 gen tasks, got %v", tasks[5].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[5].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-gen-scripts-tui"] || !depSet["T-test-gen-scripts-api"] || !depSet["T-test-gen-scripts-cli"] {
		t.Errorf("T-test-run missing expected deps, got %v", tasks[5].Dependencies)
	}
}

func TestGenerateTestTaskMD_WithTestType(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{"T-test-eval-cases"},
		Type: TypeTestGenScripts, Scope: "all",
		TestType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, "my-feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "api") {
		t.Error("missing test type in body")
	}
	if !strings.Contains(s, "docs/conventions/testing-") {
		t.Error("body should reference docs/conventions/testing-*")
	}
}

func TestGetQuickTestTasks_PerType_TwoTypes(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"tui", "api"}, allEnabledAuto)

	// gen-cases + per-type-gen-and-run(tui,api) + graduate + verify-regression + drift-detection = 6
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-gen-cases",
		"T-quick-gen-and-run-tui", "T-quick-gen-and-run-api",
		"T-quick-graduate",
		"T-quick-verify-regression",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Keys include type suffix (no language)
	if tasks[1].Key != "quick-gen-and-run-tui" {
		t.Errorf("tasks[1].Key = %q, want quick-gen-and-run-tui", tasks[1].Key)
	}
	if tasks[2].Key != "quick-gen-and-run-api" {
		t.Errorf("tasks[2].Key = %q, want quick-gen-and-run-api", tasks[2].Key)
	}

	// TestType field set
	if tasks[1].TestType != "tui" {
		t.Errorf("tasks[1].TestType = %q, want tui", tasks[1].TestType)
	}
	if tasks[2].TestType != "api" {
		t.Errorf("tasks[2].TestType = %q, want api", tasks[2].TestType)
	}

	// Per-type gen-and-run depend on gen-cases
	if tasks[1].Dependencies[0] != "T-quick-gen-cases" {
		t.Errorf("T-quick-gen-and-run-tui should depend on T-quick-gen-cases, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-quick-gen-cases" {
		t.Errorf("T-quick-gen-and-run-api should depend on T-quick-gen-cases, got %v", tasks[2].Dependencies)
	}

	// T-quick-graduate depends on ALL per-type gen-and-run tasks
	if len(tasks[3].Dependencies) != 2 {
		t.Fatalf("T-quick-graduate should depend on 2 gen-and-run tasks, got %v", tasks[3].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[3].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-quick-gen-and-run-tui"] || !depSet["T-quick-gen-and-run-api"] {
		t.Errorf("T-quick-graduate deps should include T-quick-gen-and-run-tui and T-quick-gen-and-run-api, got %v", tasks[3].Dependencies)
	}

	// T-quick-verify-regression depends on T-quick-graduate
	if tasks[4].Dependencies[0] != "T-quick-graduate" {
		t.Errorf("verify-regression should depend on graduate, got %v", tasks[4].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"api"}, allEnabledAuto)

	// gen-cases + 1 gen-and-run-api + graduate + verify-regression + drift-detection = 5
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-gen-cases",
		"T-quick-gen-and-run-api",
		"T-quick-graduate",
		"T-quick-verify-regression",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if tasks[1].Type != TypeTestGenAndRun {
		t.Errorf("T-quick-gen-and-run-api Type = %q, want %q", tasks[1].Type, TypeTestGenAndRun)
	}

	if len(tasks[2].Dependencies) != 1 || tasks[2].Dependencies[0] != "T-quick-gen-and-run-api" {
		t.Errorf("T-quick-graduate should depend on T-quick-gen-and-run-api, got %v", tasks[2].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"tui", "api", "cli"}, allEnabledAuto)

	// gen-cases + 3 per-type-gen-and-run + graduate + verify-regression + drift-detection = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	// T-quick-graduate depends on all 3 gen-and-run tasks
	if len(tasks[4].Dependencies) != 3 {
		t.Fatalf("T-quick-graduate should depend on 3 gen-and-run tasks, got %v", tasks[4].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[4].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-quick-gen-and-run-tui"] || !depSet["T-quick-gen-and-run-api"] || !depSet["T-quick-gen-and-run-cli"] {
		t.Errorf("T-quick-graduate missing expected deps, got %v", tasks[4].Dependencies)
	}
}

// --- validate-ux should only be generated when interfaces include UI types ---

func TestGetBreakdownTestTasks_ValidateUx_SkippedForCLIOnly(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, validationAuto)

	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			t.Error("validate-ux should not be generated for CLI-only projects (no visual UI)")
		}
	}
}

func TestGetQuickTestTasks_ValidateUx_SkippedForCLIOnly(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli"}, validationAuto)

	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			t.Error("validate-ux should not be generated for CLI-only projects (no visual UI)")
		}
	}
}

func TestGetBreakdownTestTasks_ValidateUx_SkippedForAPIOnly(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"api"}, validationAuto)

	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			t.Error("validate-ux should not be generated for API-only projects (no visual UI)")
		}
	}
}

func TestGetBreakdownTestTasks_ValidateUx_IncludedForTUI(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"tui"}, validationAuto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			found = true
			break
		}
	}
	if !found {
		t.Error("validate-ux should be generated for TUI projects (has visual UI)")
	}
}

func TestGetBreakdownTestTasks_ValidateUx_IncludedForMixed(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli", "tui"}, validationAuto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			found = true
			break
		}
	}
	if !found {
		t.Error("validate-ux should be generated when interfaces include UI types (tui)")
	}
}

func TestGetBreakdownTestTasks_ValidateUx_CodeStillGenerated(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, validationAuto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-validate-code" {
			found = true
			break
		}
	}
	if !found {
		t.Error("validate-code should still be generated for CLI projects")
	}
}

func TestGetQuickTestTasks_DefaultAuto_IncludesSpecDrift(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli"}, defaultAuto)

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task (spec drift only), got %d", len(tasks))
	}

	if tasks[0].ID != "T-quick-doc-drift" {
		t.Errorf("task ID = %q, want T-quick-doc-drift", tasks[0].ID)
	}
	if tasks[0].Type != TypeDocDrift {
		t.Errorf("task Type = %q, want %q", tasks[0].Type, TypeDocDrift)
	}
	if len(tasks[0].Dependencies) != 0 {
		t.Errorf("T-quick-doc-drift should have no deps without e2e tasks, got %v", tasks[0].Dependencies)
	}
}
