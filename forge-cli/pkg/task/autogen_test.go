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
	E2eTest:    forgeconfig.ModeToggle{Quick: true, Full: true},
	Validation: forgeconfig.ModeToggle{Quick: true, Full: true},
}

func TestGetBreakdownTestTasks_EmptyInterfaces(t *testing.T) {
	tasks := GetBreakdownTestTasks(nil, defaultAuto)

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks with empty interfaces, got %d", len(tasks))
	}
}

func TestGetBreakdownTestTasks_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)

	// eval-journey + eval-contract + gen-scripts-cli + run + verify-regression + consolidate = 6
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-eval-journey", "T-eval-contract", "T-test-gen-scripts-cli", "T-test-run", "T-test-verify-regression", "T-specs-consolidate"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Dependency chain: eval-contract -> eval-journey, gen-scripts -> eval-contract, run -> gen-scripts, verify -> run, consolidate -> verify
	if tasks[1].Dependencies[0] != "T-eval-journey" {
		t.Errorf("eval-contract should depend on eval-journey, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-eval-contract" {
		t.Errorf("gen-scripts should depend on eval-contract, got %v", tasks[2].Dependencies)
	}
	if tasks[3].Dependencies[0] != "T-test-gen-scripts-cli" {
		t.Errorf("run should depend on gen-scripts-cli, got %v", tasks[3].Dependencies)
	}
	if tasks[4].Dependencies[0] != "T-test-run" {
		t.Errorf("verify-regression should depend on run, got %v", tasks[4].Dependencies)
	}
	if tasks[5].Dependencies[0] != "T-test-verify-regression" {
		t.Errorf("consolidate should depend on verify-regression, got %v", tasks[5].Dependencies)
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

	// gen-and-run-cli + verify-regression + drift = 3
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-quick-gen-and-run-cli", "T-quick-verify-regression", "T-quick-doc-drift"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if tasks[0].Type != TypeTestGenAndRun {
		t.Errorf("T-quick-gen-and-run-cli Type = %q, want %q", tasks[0].Type, TypeTestGenAndRun)
	}

	if tasks[1].Dependencies[0] != "T-quick-gen-and-run-cli" {
		t.Errorf("verify-regression should depend on gen-and-run, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Type != TypeDocDrift {
		t.Errorf("T-quick-doc-drift Type = %q, want %q", tasks[2].Type, TypeDocDrift)
	}
	if tasks[2].Scope != "all" {
		t.Errorf("T-quick-doc-drift Scope = %q, want %q", tasks[2].Scope, "all")
	}
}

func TestGenerateTestTaskMD(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{},
		Type: TypeTestGenScripts, Scope: "all",
		TestType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
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
	// Body now loaded from embed template, should contain strategy-based content
	if !strings.Contains(s, "executable test scripts") {
		t.Error("body should contain strategy-based content from embed template")
	}
}

func TestGenerateTestTaskMD_SharedTask(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-run", Key: "run-e2e-tests",
		Title: "Run e2e Tests", Priority: "P1",
		EstimatedTime: "30min-1h", Dependencies: []string{},
		Type: TypeTestRun, Scope: "all",
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
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
			t.Errorf("first test task should depend on highest gate, got %v", tasks[0].Dependencies)
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
			t.Errorf("first test task should depend on last business task (3.1) when it's in a higher phase than highest gate (2.gate), got %v", tasks[0].Dependencies)
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
			t.Errorf("first quick test task should depend on max business task, got %v", tasks[0].Dependencies)
		}
	})
}

func TestGetReviewDocTask(t *testing.T) {
	task := GetReviewDocTask()

	if task.ID != "T-review-doc" {
		t.Errorf("ID = %q, want T-review-doc", task.ID)
	}
	if task.Key != "review-doc" {
		t.Errorf("Key = %q, want review-doc", task.Key)
	}
	if task.Type != TypeDocReview {
		t.Errorf("Type = %q, want %q", task.Type, TypeDocReview)
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

func TestResolveReviewDocDep(t *testing.T) {
	t.Run("sets dependency on last business task", func(t *testing.T) {
		existing := map[string]Task{
			"1-doc":              {ID: "1.1", Type: TypeDoc},
			"2-doc":              {ID: "1.2", Type: TypeDoc},
			"T-test-gen-scripts": {ID: "T-test-gen-scripts-cli", Type: TypeTestGenScripts},
		}
		task := GetReviewDocTask()
		ResolveReviewDocDep(&task, existing)

		if len(task.Dependencies) != 1 {
			t.Fatalf("Dependencies = %v, want exactly 1", task.Dependencies)
		}
		if task.Dependencies[0] != "1.2" {
			t.Errorf("dep = %q, want 1.2", task.Dependencies[0])
		}
	})

	t.Run("empty tasks", func(t *testing.T) {
		existing := map[string]Task{}
		task := GetReviewDocTask()
		ResolveReviewDocDep(&task, existing)

		if len(task.Dependencies) != 0 {
			t.Errorf("Dependencies = %v, want empty for no tasks", task.Dependencies)
		}
	})
}

func TestGetBreakdownTestTasks_PerType_TwoTypes(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"tui", "api"}, defaultAuto)

	// eval-journey + eval-contract + 2 gen-scripts(tui,api) + run + verify-regression + consolidate = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-eval-journey", "T-eval-contract",
		"T-test-gen-scripts-tui", "T-test-gen-scripts-api",
		"T-test-run",
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

	// T-test-verify-regression depends on T-test-run
	if tasks[5].Dependencies[0] != "T-test-run" {
		t.Errorf("verify-regression should depend on run, got %v", tasks[5].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"api"}, defaultAuto)

	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-eval-journey", "T-eval-contract",
		"T-test-gen-scripts-api",
		"T-test-run", "T-test-verify-regression", "T-specs-consolidate",
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

	// eval-journey + eval-contract + 3 gen tasks + run + verify-regression + consolidate = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
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
		EstimatedTime: "1-2h", Dependencies: []string{},
		Type: TypeTestGenScripts, Scope: "all",
		TestType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "api") {
		t.Error("missing test type in body")
	}
	// Body loaded from embed template with TestType appended
	if !strings.Contains(s, "Type: **api**") {
		t.Error("body should contain TestType note")
	}
}

func TestGetQuickTestTasks_PerType_TwoTypes(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"tui", "api"}, allEnabledAuto)

	// per-type-gen-and-run(tui,api) + verify-regression + drift-detection = 4
	if len(tasks) != 4 {
		t.Fatalf("expected 4 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-gen-and-run-tui", "T-quick-gen-and-run-api",
		"T-quick-verify-regression",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Keys include type suffix (no language)
	if tasks[0].Key != "quick-gen-and-run-tui" {
		t.Errorf("tasks[0].Key = %q, want quick-gen-and-run-tui", tasks[0].Key)
	}
	if tasks[1].Key != "quick-gen-and-run-api" {
		t.Errorf("tasks[1].Key = %q, want quick-gen-and-run-api", tasks[1].Key)
	}

	// TestType field set
	if tasks[0].TestType != "tui" {
		t.Errorf("tasks[0].TestType = %q, want tui", tasks[0].TestType)
	}
	if tasks[1].TestType != "api" {
		t.Errorf("tasks[1].TestType = %q, want api", tasks[1].TestType)
	}

	// T-quick-verify-regression depends on ALL per-type gen-and-run tasks
	if len(tasks[2].Dependencies) != 2 {
		t.Fatalf("T-quick-verify-regression should depend on 2 gen-and-run tasks, got %v", tasks[2].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[2].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-quick-gen-and-run-tui"] || !depSet["T-quick-gen-and-run-api"] {
		t.Errorf("T-quick-verify-regression deps should include T-quick-gen-and-run-tui and T-quick-gen-and-run-api, got %v", tasks[2].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"api"}, allEnabledAuto)

	// 1 gen-and-run-api + verify-regression + drift-detection = 3
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-quick-gen-and-run-api",
		"T-quick-verify-regression",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if tasks[0].Type != TypeTestGenAndRun {
		t.Errorf("T-quick-gen-and-run-api Type = %q, want %q", tasks[0].Type, TypeTestGenAndRun)
	}

	if len(tasks[1].Dependencies) != 1 || tasks[1].Dependencies[0] != "T-quick-gen-and-run-api" {
		t.Errorf("T-quick-verify-regression should depend on T-quick-gen-and-run-api, got %v", tasks[1].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"tui", "api", "cli"}, allEnabledAuto)

	// 3 per-type-gen-and-run + verify-regression + drift-detection = 5
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	// T-quick-verify-regression depends on all 3 gen-and-run tasks
	if len(tasks[3].Dependencies) != 3 {
		t.Fatalf("T-quick-verify-regression should depend on 3 gen-and-run tasks, got %v", tasks[3].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[3].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-quick-gen-and-run-tui"] || !depSet["T-quick-gen-and-run-api"] || !depSet["T-quick-gen-and-run-cli"] {
		t.Errorf("T-quick-verify-regression missing expected deps, got %v", tasks[3].Dependencies)
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

// --- Embed template loading tests ---

func TestGenerateTestTaskMD_EmbedTemplate_LoadsContent(t *testing.T) {
	tests := []struct {
		name         string
		typ          string
		wantContains string
	}{
		{"gen-scripts", TypeTestGenScripts, "executable test scripts"},
		{"gen-and-run", TypeTestGenAndRun, "Phase 1"},
		{"run", TypeTestRun, "staged e2e test scripts"},
		{"verify-regression", TypeTestVerifyRegression, "just test-e2e"},
		{"eval-journey", TypeEvalJourney, "6-dimension rubric"},
		{"eval-contract", TypeEvalContract, "6-dimension rubric"},
		{"validation-code", TypeValidationCode, "quality gate"},
		{"validation-ux", TypeValidationUx, "accessibility, usability"},
		{"doc-review", TypeDocReview, "acceptance criteria"},
		{"doc-consolidate", TypeDocConsolidate, "CROSS items"},
		{"doc-drift", TypeDocDrift, "git diff --name-only"},
		{"clean-code", TypeCleanCode, "Simplify and clean"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := AutoGenTaskDef{
				ID: "T-test", Key: "test",
				Title: "Test Task", Priority: "P1",
				EstimatedTime: "1h", Type: tt.typ, Scope: "all",
			}

			content, err := GenerateTestTaskMD(def, BodyContext{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			s := string(content)
			if !strings.Contains(s, tt.wantContains) {
				t.Errorf("body for type %q should contain %q", tt.typ, tt.wantContains)
			}
		})
	}
}

func TestGenerateTestTaskMD_StrategyContentAppendedAfterTemplate(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Type: TypeTestGenScripts, Scope: "all",
		StrategyKind:    "generate",
		StrategyContent: []byte("# Custom Strategy\n\nUse this strategy."),
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	// Should contain template content
	if !strings.Contains(s, "executable test scripts") {
		t.Error("body should contain template content")
	}
	// StrategyContent appended AFTER template
	templateIdx := strings.Index(s, "executable test scripts")
	strategyIdx := strings.Index(s, "Custom Strategy")
	if strategyIdx <= templateIdx {
		t.Error("StrategyContent should appear after template content")
	}
}

func TestGenerateTestTaskMD_TestTypeNotedInBody(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Type: TypeTestGenScripts, Scope: "all",
		TestType: "api",
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "Type: **api**") {
		t.Error("body should contain TestType note")
	}
}

func TestGenerateTestTaskMD_FrontmatterUnchanged(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-run", Key: "run-e2e-tests",
		Title: "Run e2e Tests", Priority: "P1",
		EstimatedTime: "30min-1h", Dependencies: []string{"dep1"},
		Type: TypeTestRun, Scope: "all",
		MainSession: true,
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	// Frontmatter fields unchanged
	if !strings.Contains(s, `id: "T-test-run"`) {
		t.Error("missing id in frontmatter")
	}
	if !strings.Contains(s, `title: "Run e2e Tests"`) {
		t.Error("missing title in frontmatter")
	}
	if !strings.Contains(s, `priority: "P1"`) {
		t.Error("missing priority in frontmatter")
	}
	if !strings.Contains(s, `"dep1"`) {
		t.Error("missing dependency in frontmatter")
	}
	if !strings.Contains(s, `type: "test.run"`) {
		t.Error("missing type in frontmatter")
	}
	if !strings.Contains(s, `scope: "all"`) {
		t.Error("missing scope in frontmatter")
	}
	if !strings.Contains(s, "mainSession: true") {
		t.Error("missing mainSession in frontmatter")
	}
}

// --- BodyContext and renderBody tests ---

func TestRenderBody_SubstitutesAllPlaceholders(t *testing.T) {
	template := `Feature: {{FEATURE_SLUG}}
Mode: {{MODE}}
Scope: {{SCOPE}}
Interfaces: {{INTERFACES}}
Test Type: {{TEST_TYPE}}
Acceptance:
{{ACCEPTANCE_CRITERIA}}`

	ctx := BodyContext{
		FeatureSlug:        "my-feature",
		Mode:               "quick",
		Scope:              []string{"backend", "frontend"},
		Interfaces:         []string{"api", "cli"},
		AcceptanceCriteria: []string{"AC1: works", "AC2: fast"},
	}
	def := AutoGenTaskDef{TestType: "api"}

	result := renderBody(template, def, ctx)

	if !strings.Contains(result, "Feature: my-feature") {
		t.Error("FEATURE_SLUG not substituted")
	}
	if !strings.Contains(result, "Mode: quick") {
		t.Error("MODE not substituted")
	}
	if !strings.Contains(result, "- backend\n- frontend") {
		t.Errorf("SCOPE not substituted, got:\n%s", result)
	}
	if !strings.Contains(result, "- api\n- cli") {
		t.Errorf("INTERFACES not substituted, got:\n%s", result)
	}
	if !strings.Contains(result, "Test Type: api") {
		t.Error("TEST_TYPE not substituted")
	}
	if !strings.Contains(result, "- [ ] AC1: works") {
		t.Errorf("ACCEPTANCE_CRITERIA not substituted, got:\n%s", result)
	}
}

func TestRenderBody_EmptyMode_OmitsLine(t *testing.T) {
	template := "Feature: {{FEATURE_SLUG}}\nMode: {{MODE}}\nDone"
	ctx := BodyContext{FeatureSlug: "test"}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	if strings.Contains(result, "Mode:") {
		t.Errorf("Mode line should be omitted when empty, got:\n%s", result)
	}
	if !strings.Contains(result, "Feature: test") {
		t.Error("FEATURE_SLUG should still be present")
	}
}

func TestRenderBody_EmptyScope_OmitsSection(t *testing.T) {
	template := "Start\n## Scope\n{{SCOPE}}\n## End"
	ctx := BodyContext{FeatureSlug: "test"}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	if strings.Contains(result, "## Scope") {
		t.Errorf("Scope section should be omitted when empty, got:\n%s", result)
	}
	if !strings.Contains(result, "## End") {
		t.Error("Content after scope section should remain")
	}
}

func TestRenderBody_EmptyInterfaces_Default(t *testing.T) {
	template := "Interfaces: {{INTERFACES}}"
	ctx := BodyContext{FeatureSlug: "test"}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	if !strings.Contains(result, "See .forge/config.yaml") {
		t.Errorf("Empty interfaces should default to 'See .forge/config.yaml', got:\n%s", result)
	}
}

func TestRenderBody_EmptyTestType_OmitsLine(t *testing.T) {
	template := "Feature: {{FEATURE_SLUG}}\nType: {{TEST_TYPE}}\nDone"
	ctx := BodyContext{FeatureSlug: "test"}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	if strings.Contains(result, "Type:") {
		t.Errorf("TestType line should be omitted when empty, got:\n%s", result)
	}
}

func TestRenderBody_EmptyAcceptanceCriteria_Default(t *testing.T) {
	template := "Acceptance:\n{{ACCEPTANCE_CRITERIA}}"
	ctx := BodyContext{FeatureSlug: "test"}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	if !strings.Contains(result, "- [ ] All acceptance criteria met") {
		t.Errorf("Empty acceptance criteria should default, got:\n%s", result)
	}
}

func TestRenderBody_EmptyBodyContext_KnownPlaceholdersResolved(t *testing.T) {
	template := "Feature: {{FEATURE_SLUG}}\nMode: {{MODE}}\n## Scope\n{{SCOPE}}\n## Other\nInterfaces: {{INTERFACES}}\nType: {{TEST_TYPE}}\nAcceptance:\n{{ACCEPTANCE_CRITERIA}}"
	ctx := BodyContext{}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	knownPlaceholders := []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{INTERFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}"}
	for _, ph := range knownPlaceholders {
		if strings.Contains(result, ph) {
			t.Errorf("placeholder %s not resolved in output:\n%s", ph, result)
		}
	}
}

func TestGenerateTestTaskMD_WithBodyContext(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Type: TypeTestGenScripts, Scope: "all",
		TestType: "api",
	}
	ctx := BodyContext{
		FeatureSlug: "my-feature",
		Mode:        "quick",
	}

	content, err := GenerateTestTaskMD(def, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	// Frontmatter still works
	if !strings.Contains(s, `id: "T-test-gen-scripts-api"`) {
		t.Error("missing id in frontmatter")
	}
	// Template body loaded (placeholder substitution applied)
	if !strings.Contains(s, "executable test scripts") {
		t.Error("body should contain template content")
	}
}

func TestGenerateTestTaskMD_BackwardCompat_EmptyBodyContext(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Type: TypeTestGenScripts, Scope: "all",
	}

	// Passing empty BodyContext should produce same output as before
	content, err := GenerateTestTaskMD(def, BodyContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	// Only check our 6 managed placeholders are resolved
	managed := []string{
		"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}",
		"{{INTERFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}",
	}
	for _, ph := range managed {
		if strings.Contains(s, ph) {
			t.Errorf("managed placeholder %s should be resolved", ph)
		}
	}
	if !strings.Contains(s, `id: "T-test-gen-scripts-api"`) {
		t.Error("frontmatter should be intact")
	}
}

// --- Focused body content verification tests ---

func TestRenderBody_FeatureSlug(t *testing.T) {
	tests := []struct {
		name        string
		featureSlug string
		want        string
	}{
		{"populated slug", "my-feature", "my-feature"},
		{"empty slug", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := "Feature: {{FEATURE_SLUG}} is ready."
			ctx := BodyContext{FeatureSlug: tt.featureSlug}
			def := AutoGenTaskDef{}

			result := renderBody(template, def, ctx)

			if !strings.Contains(result, tt.want) {
				t.Errorf("expected %q in result, got: %s", tt.want, result)
			}
			if strings.Contains(result, "{{FEATURE_SLUG}}") {
				t.Error("FEATURE_SLUG placeholder should be resolved")
			}
		})
	}
}

func TestRenderBody_ScopeAndInterfaces(t *testing.T) {
	t.Run("populated scope and interfaces", func(t *testing.T) {
		template := "## Feature Context\n- Scope: {{SCOPE}}\n- Test interfaces: {{INTERFACES}}"
		ctx := BodyContext{
			FeatureSlug: "feat",
			Scope:       []string{"backend", "frontend"},
			Interfaces:  []string{"api", "cli"},
		}
		def := AutoGenTaskDef{}

		result := renderBody(template, def, ctx)

		if !strings.Contains(result, "- backend") {
			t.Error("scope item 'backend' should be present")
		}
		if !strings.Contains(result, "- frontend") {
			t.Error("scope item 'frontend' should be present")
		}
		if !strings.Contains(result, "- api") {
			t.Error("interface 'api' should be present")
		}
		if !strings.Contains(result, "- cli") {
			t.Error("interface 'cli' should be present")
		}
		if strings.Contains(result, "{{SCOPE}}") {
			t.Error("SCOPE placeholder should be resolved")
		}
		if strings.Contains(result, "{{INTERFACES}}") {
			t.Error("INTERFACES placeholder should be resolved")
		}
	})
}

func TestRenderBody_AcceptanceCriteria(t *testing.T) {
	t.Run("criteria filled as checklist", func(t *testing.T) {
		template := "## Validation Criteria\n{{ACCEPTANCE_CRITERIA}}\n## End"
		ctx := BodyContext{
			FeatureSlug:        "feat",
			AcceptanceCriteria: []string{"AC1: Login works", "AC2: Logout works"},
		}
		def := AutoGenTaskDef{}

		result := renderBody(template, def, ctx)

		if !strings.Contains(result, "- [ ] AC1: Login works") {
			t.Error("first AC should be filled as unchecked checklist item")
		}
		if !strings.Contains(result, "- [ ] AC2: Logout works") {
			t.Error("second AC should be filled as unchecked checklist item")
		}
		if strings.Contains(result, "{{ACCEPTANCE_CRITERIA}}") {
			t.Error("ACCEPTANCE_CRITERIA placeholder should be resolved")
		}
	})
}

func TestRenderBody_EmptyFields(t *testing.T) {
	t.Run("all fields empty uses fallbacks", func(t *testing.T) {
		template := "Feature: {{FEATURE_SLUG}}\nMode: {{MODE}}\n## Scope\n{{SCOPE}}\n## Other\nInterfaces: {{INTERFACES}}\nType: {{TEST_TYPE}}\n{{ACCEPTANCE_CRITERIA}}"
		ctx := BodyContext{}
		def := AutoGenTaskDef{}

		result := renderBody(template, def, ctx)

		// Mode line omitted
		if strings.Contains(result, "Mode:") {
			t.Error("Mode line should be omitted when empty")
		}
		// Scope section omitted
		if strings.Contains(result, "## Scope") {
			t.Error("Scope section should be omitted when empty")
		}
		// Interfaces fallback
		if !strings.Contains(result, "See .forge/config.yaml") {
			t.Error("Empty interfaces should use fallback")
		}
		// TestType line omitted
		if strings.Contains(result, "Type:") {
			t.Error("TestType line should be omitted when empty")
		}
		// AcceptanceCriteria fallback
		if !strings.Contains(result, "- [ ] All acceptance criteria met") {
			t.Error("Empty acceptance criteria should use fallback")
		}
		// No leftover placeholders
		for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{INTERFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}"} {
			if strings.Contains(result, ph) {
				t.Errorf("placeholder %s should be resolved", ph)
			}
		}
	})
}

// TestBodyContentPerStrategy verifies each task type gets correct body content
// with a populated BodyContext, grouped by strategy (A/B/C).
func TestBodyContentPerStrategy(t *testing.T) {
	tests := []struct {
		name         string
		typ          string
		testType     string // per-type interface suffix (e.g., "api")
		ctx          BodyContext
		wantContains []string
	}{
		// Strategy A: Feature context (slug + scope + interfaces injected)
		{"gen-scripts has feature context", TypeTestGenScripts, "api", BodyContext{
			FeatureSlug: "feat",
		}, []string{"feat", "api"}},
		{"gen-and-run has feature context", TypeTestGenAndRun, "tui", BodyContext{
			FeatureSlug: "feat",
		}, []string{"feat", "tui"}},
		{"run has feature context", TypeTestRun, "", BodyContext{
			FeatureSlug: "feat", Scope: []string{"backend"},
		}, []string{"feat", "- backend"}},

		// Strategy B: Acceptance criteria pre-filled as validation checklist
		{"validation-code has criteria", TypeValidationCode, "", BodyContext{
			FeatureSlug: "feat", AcceptanceCriteria: []string{"AC1: works", "AC2: fast"},
		}, []string{"feat", "- [ ] AC1: works", "- [ ] AC2: fast"}},
		{"validation-ux has criteria", TypeValidationUx, "", BodyContext{
			FeatureSlug: "feat", AcceptanceCriteria: []string{"AC1: accessible"},
		}, []string{"feat", "- [ ] AC1: accessible"}},

		// Strategy C: Discovery strategy steps present (git diff, directory scan)
		{"doc-review has discovery strategy", TypeDocReview, "", BodyContext{
			FeatureSlug: "feat", Mode: "breakdown",
		}, []string{"feat", "breakdown mode"}},
		{"doc-consolidate has discovery strategy", TypeDocConsolidate, "", BodyContext{
			FeatureSlug: "feat", Scope: []string{"backend"},
		}, []string{"feat", "- backend", "Discovery Strategy"}},
		{"doc-drift has git diff strategy", TypeDocDrift, "", BodyContext{
			FeatureSlug: "feat",
		}, []string{"feat", "git diff"}},
		{"clean-code has git diff strategy", TypeCleanCode, "", BodyContext{
			FeatureSlug: "feat",
		}, []string{"feat", "git diff"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := AutoGenTaskDef{
				ID: "T-test", Key: "test",
				Title: "Test Task", Priority: "P1",
				EstimatedTime: "1h", Type: tt.typ, Scope: "all",
				TestType: tt.testType,
			}

			content, err := GenerateTestTaskMD(def, tt.ctx)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			s := string(content)
			for _, want := range tt.wantContains {
				if !strings.Contains(s, want) {
					t.Errorf("type %q body should contain %q, got:\n%s", tt.typ, want, s)
				}
			}

			// Verify no managed placeholders left unresolved
			for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{INTERFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}"} {
				if strings.Contains(s, ph) {
					t.Errorf("type %q has unresolved placeholder %s", tt.typ, ph)
				}
			}
		})
	}
}

func TestAutogenTypeToFileMapping(t *testing.T) {
	// Verify all 12 auto-gen types have a mapping entry
	wantTypes := []string{
		TypeTestGenScripts, TypeTestGenAndRun, TypeTestRun,
		TypeTestVerifyRegression, TypeEvalJourney, TypeEvalContract,
		TypeValidationCode, TypeValidationUx,
		TypeDocReview, TypeDocConsolidate, TypeDocDrift, TypeCleanCode,
	}

	if len(autogenTypeToFile) != len(wantTypes) {
		t.Errorf("autogenTypeToFile has %d entries, want %d", len(autogenTypeToFile), len(wantTypes))
	}

	for _, typ := range wantTypes {
		file, ok := autogenTypeToFile[typ]
		if !ok {
			t.Errorf("type %q missing from autogenTypeToFile", typ)
			continue
		}
		// Verify file can be read from embed FS
		data, err := autogenTemplateFS.ReadFile(file)
		if err != nil {
			t.Errorf("cannot read template file %q for type %q: %v", file, typ, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("template file %q for type %q is empty", file, typ)
		}
	}
}
