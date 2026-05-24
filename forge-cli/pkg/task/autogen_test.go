package task

import (
	"fmt"
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

	// gen-journeys-cli + eval-journey + gen-contracts + eval-contract + gen-scripts-cli + run + verify-regression + consolidate = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-test-gen-journeys-cli", "T-eval-journey", "T-test-gen-contracts", "T-eval-contract", "T-test-gen-scripts-cli", "T-test-run", "T-test-verify-regression", "T-specs-consolidate"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Dependency chain: eval-journey -> gen-journeys, gen-contracts -> eval-journey, eval-contract -> gen-contracts, gen-scripts -> eval-contract, run -> gen-scripts, verify -> run, consolidate -> verify
	if tasks[1].Dependencies[0] != "T-test-gen-journeys-cli" {
		t.Errorf("eval-journey should depend on gen-journeys-cli, got %v", tasks[1].Dependencies)
	}
	if tasks[2].Dependencies[0] != "T-eval-journey" {
		t.Errorf("gen-contracts should depend on eval-journey, got %v", tasks[2].Dependencies)
	}
	if tasks[3].Dependencies[0] != "T-test-gen-contracts" {
		t.Errorf("eval-contract should depend on gen-contracts, got %v", tasks[3].Dependencies)
	}
	if tasks[4].Dependencies[0] != "T-eval-contract" {
		t.Errorf("gen-scripts should depend on eval-contract, got %v", tasks[4].Dependencies)
	}
	if tasks[5].Dependencies[0] != "T-test-gen-scripts-cli" {
		t.Errorf("run should depend on gen-scripts-cli, got %v", tasks[5].Dependencies)
	}
	if tasks[6].Dependencies[0] != "T-test-run" {
		t.Errorf("verify-regression should depend on run, got %v", tasks[6].Dependencies)
	}
	if tasks[7].Dependencies[0] != "T-test-verify-regression" {
		t.Errorf("consolidate should depend on verify-regression, got %v", tasks[7].Dependencies)
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

	// gen-journeys-cli + gen-contracts + gen-scripts-cli + run + verify-regression + drift = 6
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys-cli",
		"T-test-gen-contracts",
		"T-test-gen-scripts-cli",
		"T-test-run",
		"T-test-verify-regression",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// No gen-and-run tasks in Quick mode
	for _, task := range tasks {
		if task.Type == TypeTestGenAndRun {
			t.Errorf("Quick mode should not contain gen-and-run tasks, found %q", task.ID)
		}
	}

	if tasks[0].Type != TypeTestGenJourneys {
		t.Errorf("T-test-gen-journeys-cli Type = %q, want %q", tasks[0].Type, TypeTestGenJourneys)
	}
	if tasks[1].Type != TypeTestGenContracts {
		t.Errorf("T-test-gen-contracts Type = %q, want %q", tasks[1].Type, TypeTestGenContracts)
	}

	if tasks[5].Type != TypeDocDrift {
		t.Errorf("T-quick-doc-drift Type = %q, want %q", tasks[5].Type, TypeDocDrift)
	}
	if tasks[5].Scope != "all" {
		t.Errorf("T-quick-doc-drift Scope = %q, want %q", tasks[5].Scope, "all")
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

	// 2 gen-journeys + eval-journey + gen-contracts + eval-contract + 2 gen-scripts(tui,api) + run + verify-regression + consolidate = 10
	if len(tasks) != 10 {
		t.Fatalf("expected 10 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys-tui", "T-test-gen-journeys-api",
		"T-eval-journey",
		"T-test-gen-contracts",
		"T-eval-contract",
		"T-test-gen-scripts-tui", "T-test-gen-scripts-api",
		"T-test-run",
		"T-test-verify-regression", "T-specs-consolidate",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// gen-journeys keys include type suffix
	if tasks[0].Key != "gen-journeys-tui" {
		t.Errorf("tasks[0].Key = %q, want gen-journeys-tui", tasks[0].Key)
	}
	if tasks[1].Key != "gen-journeys-api" {
		t.Errorf("tasks[1].Key = %q, want gen-journeys-api", tasks[1].Key)
	}

	// gen-scripts keys include type suffix
	if tasks[5].Key != "gen-test-scripts-tui" {
		t.Errorf("tasks[5].Key = %q, want gen-test-scripts-tui", tasks[5].Key)
	}
	if tasks[6].Key != "gen-test-scripts-api" {
		t.Errorf("tasks[6].Key = %q, want gen-test-scripts-api", tasks[6].Key)
	}

	// TestType field set for gen-journeys
	if tasks[0].TestType != "tui" {
		t.Errorf("tasks[0].TestType = %q, want tui", tasks[0].TestType)
	}
	if tasks[1].TestType != "api" {
		t.Errorf("tasks[1].TestType = %q, want api", tasks[1].TestType)
	}

	// TestType field set for gen-scripts
	if tasks[5].TestType != "tui" {
		t.Errorf("tasks[5].TestType = %q, want tui", tasks[5].TestType)
	}
	if tasks[6].TestType != "api" {
		t.Errorf("tasks[6].TestType = %q, want api", tasks[6].TestType)
	}

	// eval-journey depends on ALL gen-journeys tasks
	if len(tasks[2].Dependencies) != 2 {
		t.Fatalf("T-eval-journey should depend on 2 gen-journeys tasks, got %v", tasks[2].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[2].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-gen-journeys-tui"] || !depSet["T-test-gen-journeys-api"] {
		t.Errorf("T-eval-journey deps should include both gen-journeys, got %v", tasks[2].Dependencies)
	}

	// T-test-run depends on ALL per-type gen-scripts tasks
	if len(tasks[7].Dependencies) != 2 {
		t.Fatalf("T-test-run should depend on 2 gen tasks, got %v", tasks[7].Dependencies)
	}
	depSet2 := make(map[string]bool)
	for _, d := range tasks[7].Dependencies {
		depSet2[d] = true
	}
	if !depSet2["T-test-gen-scripts-tui"] || !depSet2["T-test-gen-scripts-api"] {
		t.Errorf("T-test-run deps should include T-test-gen-scripts-tui and T-test-gen-scripts-api, got %v", tasks[7].Dependencies)
	}

	// T-test-verify-regression depends on T-test-run
	if tasks[8].Dependencies[0] != "T-test-run" {
		t.Errorf("verify-regression should depend on run, got %v", tasks[8].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"api"}, defaultAuto)

	// gen-journeys-api + eval-journey + gen-contracts + eval-contract + gen-scripts-api + run + verify-regression + consolidate = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys-api",
		"T-eval-journey",
		"T-test-gen-contracts",
		"T-eval-contract",
		"T-test-gen-scripts-api",
		"T-test-run", "T-test-verify-regression", "T-specs-consolidate",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if len(tasks[5].Dependencies) != 1 || tasks[5].Dependencies[0] != "T-test-gen-scripts-api" {
		t.Errorf("T-test-run should depend on T-test-gen-scripts-api, got %v", tasks[5].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"tui", "api", "cli"}, defaultAuto)

	// 3 gen-journeys + eval-journey + gen-contracts + eval-contract + 3 gen-scripts + run + verify-regression + consolidate = 12
	if len(tasks) != 12 {
		t.Fatalf("expected 12 tasks, got %d", len(tasks))
	}

	// T-test-run depends on all 3 gen tasks
	runIdx := findTaskIndex(tasks, "T-test-run")
	if runIdx < 0 {
		t.Fatalf("T-test-run not found in tasks")
	}
	if len(tasks[runIdx].Dependencies) != 3 {
		t.Fatalf("T-test-run should depend on 3 gen tasks, got %v", tasks[runIdx].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[runIdx].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-gen-scripts-tui"] || !depSet["T-test-gen-scripts-api"] || !depSet["T-test-gen-scripts-cli"] {
		t.Errorf("T-test-run missing expected deps, got %v", tasks[runIdx].Dependencies)
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

	// 2 gen-journeys + gen-contracts + 2 gen-scripts + run + verify-regression + drift = 8
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys-tui", "T-test-gen-journeys-api",
		"T-test-gen-contracts",
		"T-test-gen-scripts-tui", "T-test-gen-scripts-api",
		"T-test-run",
		"T-test-verify-regression",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Keys include type suffix
	if tasks[0].Key != "gen-journeys-tui" {
		t.Errorf("tasks[0].Key = %q, want gen-journeys-tui", tasks[0].Key)
	}
	if tasks[1].Key != "gen-journeys-api" {
		t.Errorf("tasks[1].Key = %q, want gen-journeys-api", tasks[1].Key)
	}

	// TestType field set for gen-journeys
	if tasks[0].TestType != "tui" {
		t.Errorf("tasks[0].TestType = %q, want tui", tasks[0].TestType)
	}
	if tasks[1].TestType != "api" {
		t.Errorf("tasks[1].TestType = %q, want api", tasks[1].TestType)
	}

	// TestType field set for gen-scripts
	if tasks[3].TestType != "tui" {
		t.Errorf("tasks[3].TestType = %q, want tui", tasks[3].TestType)
	}
	if tasks[4].TestType != "api" {
		t.Errorf("tasks[4].TestType = %q, want api", tasks[4].TestType)
	}

	// gen-contracts depends on ALL gen-journeys tasks
	gcIdx := findTaskIndexOrPanic(tasks, "T-test-gen-contracts")
	if len(tasks[gcIdx].Dependencies) != 2 {
		t.Fatalf("gen-contracts should depend on 2 gen-journeys, got %v", tasks[gcIdx].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[gcIdx].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-gen-journeys-tui"] || !depSet["T-test-gen-journeys-api"] {
		t.Errorf("gen-contracts deps should include both gen-journeys, got %v", tasks[gcIdx].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"api"}, allEnabledAuto)

	// gen-journeys-api + gen-contracts + gen-scripts-api + run + verify-regression + drift = 6
	if len(tasks) != 6 {
		t.Fatalf("expected 6 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys-api",
		"T-test-gen-contracts",
		"T-test-gen-scripts-api",
		"T-test-run",
		"T-test-verify-regression",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if tasks[0].Type != TypeTestGenJourneys {
		t.Errorf("T-test-gen-journeys-api Type = %q, want %q", tasks[0].Type, TypeTestGenJourneys)
	}
	if tasks[1].Type != TypeTestGenContracts {
		t.Errorf("T-test-gen-contracts Type = %q, want %q", tasks[1].Type, TypeTestGenContracts)
	}

	// gen-contracts depends on gen-journeys-api
	gcIdx := findTaskIndexOrPanic(tasks, "T-test-gen-contracts")
	if len(tasks[gcIdx].Dependencies) != 1 || tasks[gcIdx].Dependencies[0] != "T-test-gen-journeys-api" {
		t.Errorf("gen-contracts should depend on T-test-gen-journeys-api, got %v", tasks[gcIdx].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"tui", "api", "cli"}, allEnabledAuto)

	// 3 gen-journeys + gen-contracts + 3 gen-scripts + run + verify-regression + drift = 10
	if len(tasks) != 10 {
		t.Fatalf("expected 10 tasks, got %d", len(tasks))
	}

	// gen-contracts depends on all 3 gen-journeys
	gcIdx := findTaskIndexOrPanic(tasks, "T-test-gen-contracts")
	if len(tasks[gcIdx].Dependencies) != 3 {
		t.Fatalf("gen-contracts should depend on 3 gen-journeys, got %v", tasks[gcIdx].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[gcIdx].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-gen-journeys-tui"] || !depSet["T-test-gen-journeys-api"] || !depSet["T-test-gen-journeys-cli"] {
		t.Errorf("gen-contracts missing expected deps, got %v", tasks[gcIdx].Dependencies)
	}

	// run depends on all 3 gen-scripts
	runIdx := findTaskIndexOrPanic(tasks, "T-test-run")
	if len(tasks[runIdx].Dependencies) != 3 {
		t.Fatalf("run should depend on 3 gen-scripts, got %v", tasks[runIdx].Dependencies)
	}
	runDepSet := make(map[string]bool)
	for _, d := range tasks[runIdx].Dependencies {
		runDepSet[d] = true
	}
	if !runDepSet["T-test-gen-scripts-tui"] || !runDepSet["T-test-gen-scripts-api"] || !runDepSet["T-test-gen-scripts-cli"] {
		t.Errorf("run missing expected deps, got %v", tasks[runIdx].Dependencies)
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
Interfaces: {{SURFACES}}
Test Type: {{TEST_TYPE}}
Acceptance:
{{ACCEPTANCE_CRITERIA}}`

	ctx := BodyContext{
		FeatureSlug:        "my-feature",
		Mode:               "quick",
		Scope:              []string{"backend", "frontend"},
		SurfaceTypes:       []string{"api", "cli"},
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
		t.Errorf("SURFACES not substituted, got:\n%s", result)
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
	template := "Interfaces: {{SURFACES}}"
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
	template := "Feature: {{FEATURE_SLUG}}\nMode: {{MODE}}\n## Scope\n{{SCOPE}}\n## Other\nInterfaces: {{SURFACES}}\nType: {{TEST_TYPE}}\nAcceptance:\n{{ACCEPTANCE_CRITERIA}}"
	ctx := BodyContext{}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	knownPlaceholders := []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{SURFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}"}
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
		"{{SURFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}",
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
		template := "## Feature Context\n- Scope: {{SCOPE}}\n- Test interfaces: {{SURFACES}}"
		ctx := BodyContext{
			FeatureSlug:  "feat",
			Scope:        []string{"backend", "frontend"},
			SurfaceTypes: []string{"api", "cli"},
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
		if strings.Contains(result, "{{SURFACES}}") {
			t.Error("SURFACES placeholder should be resolved")
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
		template := "Feature: {{FEATURE_SLUG}}\nMode: {{MODE}}\n## Scope\n{{SCOPE}}\n## Other\nInterfaces: {{SURFACES}}\nType: {{TEST_TYPE}}\n{{ACCEPTANCE_CRITERIA}}"
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
		for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{SURFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}"} {
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
			for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{SURFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}"} {
				if strings.Contains(s, ph) {
					t.Errorf("type %q has unresolved placeholder %s", tt.typ, ph)
				}
			}
		})
	}
}

func TestGenJourneysTemplateContent(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("data/test-gen-journeys.md")
	if err != nil {
		t.Fatalf("cannot read test-gen-journeys.md: %v", err)
	}
	s := string(data)

	// AC1: uses {{FEATURE_SLUG}} and {{MODE}} placeholders
	if !strings.Contains(s, "{{FEATURE_SLUG}}") {
		t.Error("template must contain {{FEATURE_SLUG}} placeholder")
	}
	if !strings.Contains(s, "{{MODE}}") {
		t.Error("template must contain {{MODE}} placeholder")
	}

	// AC3: contains AUTO_COMMIT conditional instruction
	if !strings.Contains(s, "AUTO_COMMIT") {
		t.Error("template must contain AUTO_COMMIT conditional instruction")
	}

	// AC6: template content guides the AI executor
	if !strings.Contains(s, "Journey") {
		t.Error("template must reference Journey generation flow")
	}
}

func TestGenContractsTemplateContent(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("data/test-gen-contracts.md")
	if err != nil {
		t.Fatalf("cannot read test-gen-contracts.md: %v", err)
	}
	s := string(data)

	// AC2: uses {{FEATURE_SLUG}} and {{MODE}} placeholders
	if !strings.Contains(s, "{{FEATURE_SLUG}}") {
		t.Error("template must contain {{FEATURE_SLUG}} placeholder")
	}
	if !strings.Contains(s, "{{MODE}}") {
		t.Error("template must contain {{MODE}} placeholder")
	}

	// AC4: contains SKIP_EVAL_GATE conditional instruction
	if !strings.Contains(s, "SKIP_EVAL_GATE") {
		t.Error("template must contain SKIP_EVAL_GATE conditional instruction")
	}

	// AC6: template content guides the AI executor
	if !strings.Contains(s, "Contract") {
		t.Error("template must reference Contract generation flow")
	}
}

func TestGenJourneysTemplateRendering(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-journeys", Key: "gen-journeys",
		Title: "Generate Test Journeys", Priority: "P1",
		EstimatedTime: "20-30min", Type: TypeTestGenJourneys, Scope: "all",
	}
	ctx := BodyContext{
		FeatureSlug: "my-feature",
		Mode:        "quick",
		Scope:       []string{"backend", "CLI"},
	}

	content, err := GenerateTestTaskMD(def, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(content)

	// Placeholders resolved
	if strings.Contains(s, "{{FEATURE_SLUG}}") {
		t.Error("{{FEATURE_SLUG}} should be resolved")
	}
	if strings.Contains(s, "{{MODE}}") {
		t.Error("{{MODE}} should be resolved")
	}
	// Feature slug present in body
	if !strings.Contains(s, "my-feature") {
		t.Error("rendered body should contain feature slug")
	}
	// Mode present
	if !strings.Contains(s, "quick") {
		t.Error("rendered body should contain mode")
	}
	// No managed placeholders left
	for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{TEST_TYPE}}"} {
		if strings.Contains(s, ph) {
			t.Errorf("placeholder %s should be resolved", ph)
		}
	}
}

func TestGenContractsTemplateRendering(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-contracts", Key: "gen-contracts",
		Title: "Generate Test Contracts", Priority: "P1",
		EstimatedTime: "30-45min", Type: TypeTestGenContracts, Scope: "all",
	}
	ctx := BodyContext{
		FeatureSlug: "my-feature",
		Mode:        "breakdown",
		Scope:       []string{"backend"},
	}

	content, err := GenerateTestTaskMD(def, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(content)

	// Placeholders resolved
	if strings.Contains(s, "{{FEATURE_SLUG}}") {
		t.Error("{{FEATURE_SLUG}} should be resolved")
	}
	if strings.Contains(s, "{{MODE}}") {
		t.Error("{{MODE}} should be resolved")
	}
	// Feature slug present
	if !strings.Contains(s, "my-feature") {
		t.Error("rendered body should contain feature slug")
	}
	// No managed placeholders left
	for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{TEST_TYPE}}"} {
		if strings.Contains(s, ph) {
			t.Errorf("placeholder %s should be resolved", ph)
		}
	}
}

func TestAutogenTypeToFileMapping(t *testing.T) {
	// Verify all auto-gen types have a mapping entry
	wantTypes := []string{
		TypeTestGenScripts, TypeTestGenAndRun, TypeTestRun,
		TypeTestVerifyRegression, TypeEvalJourney, TypeEvalContract,
		TypeValidationCode, TypeValidationUx,
		TypeDocReview, TypeDocConsolidate, TypeDocDrift, TypeCleanCode,
		TypeTestGenJourneys, TypeTestGenContracts,
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

// --- Tests for gen-journeys/gen-contracts in Breakdown mode (Task 3) ---

func TestGetBreakdownTestTasks_GenJourneysPerType(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"tui", "api"}, defaultAuto)

	foundTUI := false
	foundAPI := false
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys-tui" {
			foundTUI = true
			if task.Type != TypeTestGenJourneys {
				t.Errorf("gen-journeys-tui Type = %q, want %q", task.Type, TypeTestGenJourneys)
			}
			if task.TestType != "tui" {
				t.Errorf("gen-journeys-tui TestType = %q, want tui", task.TestType)
			}
			if task.StrategyKind != "interface" {
				t.Errorf("gen-journeys-tui StrategyKind = %q, want interface", task.StrategyKind)
			}
		}
		if task.ID == "T-test-gen-journeys-api" {
			foundAPI = true
			if task.Type != TypeTestGenJourneys {
				t.Errorf("gen-journeys-api Type = %q, want %q", task.Type, TypeTestGenJourneys)
			}
			if task.TestType != "api" {
				t.Errorf("gen-journeys-api TestType = %q, want api", task.TestType)
			}
			if task.StrategyKind != "interface" {
				t.Errorf("gen-journeys-api StrategyKind = %q, want interface", task.StrategyKind)
			}
		}
	}
	if !foundTUI {
		t.Error("missing T-test-gen-journeys-tui task")
	}
	if !foundAPI {
		t.Error("missing T-test-gen-journeys-api task")
	}
}

func TestGetBreakdownTestTasks_GenContracts(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-test-gen-contracts" {
			found = true
			if task.Type != TypeTestGenContracts {
				t.Errorf("gen-contracts Type = %q, want %q", task.Type, TypeTestGenContracts)
			}
		}
	}
	if !found {
		t.Error("missing T-test-gen-contracts task")
	}
}

func TestGetBreakdownTestTasks_NewOrdering(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)

	wantOrder := []string{
		"T-test-gen-journeys-cli",
		"T-eval-journey",
		"T-test-gen-contracts",
		"T-eval-contract",
		"T-test-gen-scripts-cli",
		"T-test-run",
		"T-test-verify-regression",
		"T-specs-consolidate",
	}

	for i, want := range wantOrder {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}
}

func TestGetBreakdownTestTasks_FullDependencyChain(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli", "api"}, defaultAuto)

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// gen-journeys tasks have no deps (pipeline entry points)
	if len(byID["T-test-gen-journeys-cli"].Dependencies) != 0 {
		t.Errorf("gen-journeys-cli should have no deps, got %v", byID["T-test-gen-journeys-cli"].Dependencies)
	}
	if len(byID["T-test-gen-journeys-api"].Dependencies) != 0 {
		t.Errorf("gen-journeys-api should have no deps, got %v", byID["T-test-gen-journeys-api"].Dependencies)
	}

	// eval-journey depends on all gen-journeys
	evalJourneyDeps := byID["T-eval-journey"].Dependencies
	if len(evalJourneyDeps) != 2 {
		t.Fatalf("eval-journey should depend on 2 gen-journeys, got %v", evalJourneyDeps)
	}
	depSet := make(map[string]bool)
	for _, d := range evalJourneyDeps {
		depSet[d] = true
	}
	if !depSet["T-test-gen-journeys-cli"] || !depSet["T-test-gen-journeys-api"] {
		t.Errorf("eval-journey deps should include both gen-journeys, got %v", evalJourneyDeps)
	}

	// gen-contracts depends on eval-journey
	if len(byID["T-test-gen-contracts"].Dependencies) != 1 || byID["T-test-gen-contracts"].Dependencies[0] != "T-eval-journey" {
		t.Errorf("gen-contracts should depend on eval-journey, got %v", byID["T-test-gen-contracts"].Dependencies)
	}

	// eval-contract depends on gen-contracts
	if len(byID["T-eval-contract"].Dependencies) != 1 || byID["T-eval-contract"].Dependencies[0] != "T-test-gen-contracts" {
		t.Errorf("eval-contract should depend on gen-contracts, got %v", byID["T-eval-contract"].Dependencies)
	}

	// gen-scripts depend on eval-contract
	if len(byID["T-test-gen-scripts-cli"].Dependencies) != 1 || byID["T-test-gen-scripts-cli"].Dependencies[0] != "T-eval-contract" {
		t.Errorf("gen-scripts-cli should depend on eval-contract, got %v", byID["T-test-gen-scripts-cli"].Dependencies)
	}
	if len(byID["T-test-gen-scripts-api"].Dependencies) != 1 || byID["T-test-gen-scripts-api"].Dependencies[0] != "T-eval-contract" {
		t.Errorf("gen-scripts-api should depend on eval-contract, got %v", byID["T-test-gen-scripts-api"].Dependencies)
	}

	// run depends on all gen-scripts
	runDeps := byID["T-test-run"].Dependencies
	if len(runDeps) != 2 {
		t.Fatalf("run should depend on 2 gen-scripts, got %v", runDeps)
	}
	runDepSet := make(map[string]bool)
	for _, d := range runDeps {
		runDepSet[d] = true
	}
	if !runDepSet["T-test-gen-scripts-cli"] || !runDepSet["T-test-gen-scripts-api"] {
		t.Errorf("run deps should include both gen-scripts, got %v", runDeps)
	}

	// verify-regression depends on run
	if len(byID["T-test-verify-regression"].Dependencies) != 1 || byID["T-test-verify-regression"].Dependencies[0] != "T-test-run" {
		t.Errorf("verify-regression should depend on run, got %v", byID["T-test-verify-regression"].Dependencies)
	}
}

func TestFindTaskIndexOrPanic_PanicsOnMissing(t *testing.T) {
	tasks := []AutoGenTaskDef{
		{ID: "T-existing-1"},
		{ID: "T-existing-2"},
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when task not found, but did not panic")
		}
		msg := fmt.Sprintf("%v", r)
		if !strings.Contains(msg, "T-missing") {
			t.Errorf("panic message should contain missing task ID %q, got %q", "T-missing", msg)
		}
		if !strings.Contains(msg, "T-existing-1") {
			t.Errorf("panic message should contain task ID T-existing-1, got %q", msg)
		}
		if !strings.Contains(msg, "T-existing-2") {
			t.Errorf("panic message should contain task ID T-existing-2, got %q", msg)
		}
	}()

	findTaskIndexOrPanic(tasks, "T-missing")
}

func TestFindTaskIndexOrPanic_ReturnsIndexWhenFound(t *testing.T) {
	tasks := []AutoGenTaskDef{
		{ID: "T-first"},
		{ID: "T-second"},
		{ID: "T-third"},
	}

	idx := findTaskIndexOrPanic(tasks, "T-second")
	if idx != 1 {
		t.Errorf("expected index 1, got %d", idx)
	}
}

func TestGetBreakdownTestTasks_GenJourneysUsesEmbedTemplate(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-journeys-cli", Key: "gen-journeys-cli",
		Title: "Generate Test Journeys (cli)", Priority: "P1",
		EstimatedTime: "20-30min", Type: TypeTestGenJourneys, Scope: "all",
		TestType: "cli", StrategyKind: "interface",
	}
	ctx := BodyContext{
		FeatureSlug: "test-feature",
		Mode:        "breakdown",
	}

	content, err := GenerateTestTaskMD(def, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "Journey") {
		t.Error("gen-journeys body should contain Journey content from embed template")
	}
	if !strings.Contains(s, "test-feature") {
		t.Error("gen-journeys body should contain feature slug")
	}
	if !strings.Contains(s, "breakdown") {
		t.Error("gen-journeys body should contain mode")
	}
}

func TestGetBreakdownTestTasks_GenContractsUsesEmbedTemplate(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-contracts", Key: "gen-contracts",
		Title: "Generate Test Contracts", Priority: "P1",
		EstimatedTime: "30-45min", Type: TypeTestGenContracts, Scope: "all",
	}
	ctx := BodyContext{
		FeatureSlug: "test-feature",
		Mode:        "breakdown",
	}

	content, err := GenerateTestTaskMD(def, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	if !strings.Contains(s, "Contract") {
		t.Error("gen-contracts body should contain Contract content from embed template")
	}
	if !strings.Contains(s, "test-feature") {
		t.Error("gen-contracts body should contain feature slug")
	}
}

func TestGetBreakdownTestTasks_NoHardcodedIndices(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)

	idSet := make(map[string]bool)
	for _, t := range tasks {
		idSet[t.ID] = true
	}

	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if !idSet[dep] {
				t.Errorf("task %q depends on %q which is not in the task list", task.ID, dep)
			}
		}
	}
}

func TestGetBreakdownTestTasks_RegressionStillValid(t *testing.T) {
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	if byID["T-test-verify-regression"].Dependencies[0] != "T-test-run" {
		t.Errorf("verify-regression should still depend on run, got %v", byID["T-test-verify-regression"].Dependencies)
	}
	if byID["T-test-run"].Dependencies[0] != "T-test-gen-scripts-cli" {
		t.Errorf("run should still depend on gen-scripts-cli, got %v", byID["T-test-run"].Dependencies)
	}
}

// --- Quick mode staged across types topology tests (Task 4) ---

func TestGetQuickTestTasks_NoGenAndRun(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli", "api"}, allEnabledAuto)

	for _, task := range tasks {
		if task.Type == TypeTestGenAndRun {
			t.Errorf("Quick mode should not generate gen-and-run tasks, found %q (type=%q)", task.ID, task.Type)
		}
	}
}

func TestGetQuickTestTasks_StagedAcrossTypesDependencyChain(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli", "api"}, allEnabledAuto)

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Stage 1: gen-journeys have no deps (pipeline entry points)
	if len(byID["T-test-gen-journeys-cli"].Dependencies) != 0 {
		t.Errorf("gen-journeys-cli should have no deps, got %v", byID["T-test-gen-journeys-cli"].Dependencies)
	}
	if len(byID["T-test-gen-journeys-api"].Dependencies) != 0 {
		t.Errorf("gen-journeys-api should have no deps, got %v", byID["T-test-gen-journeys-api"].Dependencies)
	}

	// Stage 2: gen-contracts depends on all gen-journeys
	gcDeps := byID["T-test-gen-contracts"].Dependencies
	if len(gcDeps) != 2 {
		t.Fatalf("gen-contracts should depend on 2 gen-journeys, got %v", gcDeps)
	}
	gcDepSet := make(map[string]bool)
	for _, d := range gcDeps {
		gcDepSet[d] = true
	}
	if !gcDepSet["T-test-gen-journeys-cli"] || !gcDepSet["T-test-gen-journeys-api"] {
		t.Errorf("gen-contracts deps should include both gen-journeys, got %v", gcDeps)
	}

	// Stage 3: gen-scripts depend on gen-contracts
	if len(byID["T-test-gen-scripts-cli"].Dependencies) != 1 || byID["T-test-gen-scripts-cli"].Dependencies[0] != "T-test-gen-contracts" {
		t.Errorf("gen-scripts-cli should depend on gen-contracts, got %v", byID["T-test-gen-scripts-cli"].Dependencies)
	}
	if len(byID["T-test-gen-scripts-api"].Dependencies) != 1 || byID["T-test-gen-scripts-api"].Dependencies[0] != "T-test-gen-contracts" {
		t.Errorf("gen-scripts-api should depend on gen-contracts, got %v", byID["T-test-gen-scripts-api"].Dependencies)
	}

	// Stage 4: run depends on all gen-scripts
	runDeps := byID["T-test-run"].Dependencies
	if len(runDeps) != 2 {
		t.Fatalf("run should depend on 2 gen-scripts, got %v", runDeps)
	}
	runDepSet := make(map[string]bool)
	for _, d := range runDeps {
		runDepSet[d] = true
	}
	if !runDepSet["T-test-gen-scripts-cli"] || !runDepSet["T-test-gen-scripts-api"] {
		t.Errorf("run deps should include both gen-scripts, got %v", runDeps)
	}

	// Stage 5: verify-regression depends on run
	if len(byID["T-test-verify-regression"].Dependencies) != 1 || byID["T-test-verify-regression"].Dependencies[0] != "T-test-run" {
		t.Errorf("verify-regression should depend on run, got %v", byID["T-test-verify-regression"].Dependencies)
	}
}

func TestGetQuickTestTasks_GenJourneysPerType(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"tui", "api"}, allEnabledAuto)

	foundTUI := false
	foundAPI := false
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys-tui" {
			foundTUI = true
			if task.Type != TypeTestGenJourneys {
				t.Errorf("gen-journeys-tui Type = %q, want %q", task.Type, TypeTestGenJourneys)
			}
			if task.TestType != "tui" {
				t.Errorf("gen-journeys-tui TestType = %q, want tui", task.TestType)
			}
		}
		if task.ID == "T-test-gen-journeys-api" {
			foundAPI = true
			if task.Type != TypeTestGenJourneys {
				t.Errorf("gen-journeys-api Type = %q, want %q", task.Type, TypeTestGenJourneys)
			}
			if task.TestType != "api" {
				t.Errorf("gen-journeys-api TestType = %q, want api", task.TestType)
			}
		}
	}
	if !foundTUI {
		t.Error("missing T-test-gen-journeys-tui task in Quick mode")
	}
	if !foundAPI {
		t.Error("missing T-test-gen-journeys-api task in Quick mode")
	}
}

func TestGetQuickTestTasks_GenContracts(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli"}, allEnabledAuto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-test-gen-contracts" {
			found = true
			if task.Type != TypeTestGenContracts {
				t.Errorf("gen-contracts Type = %q, want %q", task.Type, TypeTestGenContracts)
			}
		}
	}
	if !found {
		t.Error("missing T-test-gen-contracts task in Quick mode")
	}
}

func TestGetQuickTestTasks_NoHardcodedIndices(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli"}, allEnabledAuto)

	idSet := make(map[string]bool)
	for _, t := range tasks {
		idSet[t.ID] = true
	}

	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if !idSet[dep] {
				t.Errorf("task %q depends on %q which is not in the task list", task.ID, dep)
			}
		}
	}
}

func TestGetQuickTestTasks_DriftDependsOnVerifyRegression(t *testing.T) {
	tasks := GetQuickTestTasks([]string{"cli"}, allEnabledAuto)

	for _, task := range tasks {
		if task.ID == "T-quick-doc-drift" {
			if len(task.Dependencies) != 1 || task.Dependencies[0] != "T-test-verify-regression" {
				t.Errorf("T-quick-doc-drift should depend on T-test-verify-regression, got %v", task.Dependencies)
			}
			return
		}
	}
	t.Error("T-quick-doc-drift not found")
}

// --- Task 5: ResolveFirstTestDep panic and InferType ordering tests ---

func TestResolveFirstTestDep_BreakdownGracefulOnMissingGenJourneys(_ *testing.T) {
	// When gen-journeys tasks don't exist (no E2E tasks), ResolveFirstTestDep
	// should return gracefully without panicking.
	tasks := []AutoGenTaskDef{
		{ID: "T-eval-journey"},
		{ID: "T-test-gen-contracts"},
		{ID: "T-test-gen-scripts-cli"},
		{ID: "T-test-run"},
	}
	existing := map[string]Task{
		"1-gate": {ID: "1.gate"},
	}
	// Should not panic
	ResolveFirstTestDep(tasks, existing, "breakdown")
}

func TestResolveFirstTestDep_QuickGracefulOnMissingGenJourneys(_ *testing.T) {
	// When gen-journeys tasks don't exist (no E2E tasks), ResolveFirstTestDep
	// should return gracefully without panicking.
	tasks := []AutoGenTaskDef{
		{ID: "T-quick-doc-drift"},
	}
	existing := map[string]Task{
		"1-foo": {ID: "1"},
	}
	// Should not panic
	ResolveFirstTestDep(tasks, existing, "quick")
}

func TestResolveFirstTestDep_BreakdownWithCleanCode(t *testing.T) {
	existing := map[string]Task{
		"1-gate":  {ID: "1.gate"},
		"1.1-foo": {ID: "1.1"},
	}
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)

	// Add a clean-code task
	tasks = append([]AutoGenTaskDef{{ID: "T-clean-code"}}, tasks...)

	ResolveFirstTestDep(tasks, existing, "breakdown")

	cleanIdx := findTaskIndex(tasks, "T-clean-code")
	if cleanIdx < 0 {
		t.Fatal("T-clean-code not found")
	}
	if tasks[cleanIdx].Dependencies[0] != "1.gate" {
		t.Errorf("clean-code should depend on highest gate, got %v", tasks[cleanIdx].Dependencies)
	}

	firstTestIdx := findTaskIndexByPrefix(tasks, "T-test-gen-journeys")
	if firstTestIdx < 0 {
		t.Fatal("gen-journeys not found")
	}
	if tasks[firstTestIdx].Dependencies[0] != "T-clean-code" {
		t.Errorf("first test task should depend on clean-code, got %v", tasks[firstTestIdx].Dependencies)
	}
}

func TestResolveFirstTestDep_QuickWithCleanCode(t *testing.T) {
	existing := map[string]Task{
		"1-foo": {ID: "1"},
		"2-bar": {ID: "2"},
	}
	tasks := GetQuickTestTasks([]string{"cli"}, allEnabledAuto)

	// Add a clean-code task
	tasks = append([]AutoGenTaskDef{{ID: "T-clean-code"}}, tasks...)

	ResolveFirstTestDep(tasks, existing, "quick")

	cleanIdx := findTaskIndex(tasks, "T-clean-code")
	if cleanIdx < 0 {
		t.Fatal("T-clean-code not found")
	}
	if tasks[cleanIdx].Dependencies[0] != "2" {
		t.Errorf("clean-code should depend on max business task, got %v", tasks[cleanIdx].Dependencies)
	}

	firstTestIdx := findTaskIndexByPrefix(tasks, "T-test-gen-journeys")
	if firstTestIdx < 0 {
		t.Fatal("gen-journeys not found")
	}
	if tasks[firstTestIdx].Dependencies[0] != "T-clean-code" {
		t.Errorf("first test task should depend on clean-code, got %v", tasks[firstTestIdx].Dependencies)
	}
}

func TestResolveFirstTestDep_EmptyTasks_NoPanic(_ *testing.T) {
	// Empty tasks should return without panic
	ResolveFirstTestDep(nil, map[string]Task{"1": {ID: "1"}}, "breakdown")
	ResolveFirstTestDep(nil, map[string]Task{"1": {ID: "1"}}, "quick")
}

func TestResolveFirstTestDep_NoDeps_NoPanic(t *testing.T) {
	// No existing business tasks → return without panic
	tasks := GetBreakdownTestTasks([]string{"cli"}, defaultAuto)
	ResolveFirstTestDep(tasks, map[string]Task{}, "breakdown")

	// gen-journeys should have no deps set (no business tasks to depend on)
	firstTestIdx := findTaskIndexByPrefix(tasks, "T-test-gen-journeys")
	if firstTestIdx >= 0 && len(tasks[firstTestIdx].Dependencies) != 0 {
		t.Errorf("gen-journeys should have no deps when no business tasks exist, got %v", tasks[firstTestIdx].Dependencies)
	}
}

// --- {{DOC_TASK_AC}} rendering tests ---

func TestRenderBody_DocTaskAC_Populated(t *testing.T) {
	template := "Feature: {{FEATURE_SLUG}}\n## Acceptance Criteria Summary\n{{DOC_TASK_AC}}\n## End"
	ctx := BodyContext{
		FeatureSlug: "my-feature",
		DocTaskCriteria: map[string]string{
			"1-doc": "- [ ] AC 1\n- [ ] AC 2",
			"2-doc": "- [ ] AC 3",
		},
	}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	if !strings.Contains(result, "### 1-doc") {
		t.Errorf("should contain ### 1-doc sub-section header, got:\n%s", result)
	}
	if !strings.Contains(result, "- [ ] AC 1") {
		t.Errorf("should contain AC content from 1-doc, got:\n%s", result)
	}
	if !strings.Contains(result, "### 2-doc") {
		t.Errorf("should contain ### 2-doc sub-section header, got:\n%s", result)
	}
	if !strings.Contains(result, "- [ ] AC 3") {
		t.Errorf("should contain AC content from 2-doc, got:\n%s", result)
	}
	if strings.Contains(result, "{{DOC_TASK_AC}}") {
		t.Errorf("placeholder should be resolved, got:\n%s", result)
	}
}

func TestRenderBody_DocTaskAC_Empty(t *testing.T) {
	template := "Feature: {{FEATURE_SLUG}}\n{{DOC_TASK_AC}}\n## End"
	ctx := BodyContext{
		FeatureSlug: "my-feature",
	}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	// Empty DocTaskCriteria should produce empty string (placeholder removed)
	if strings.Contains(result, "{{DOC_TASK_AC}}") {
		t.Errorf("placeholder should be resolved to empty when no criteria, got:\n%s", result)
	}
}

func TestRenderBody_DocTaskAC_SortedKeys(t *testing.T) {
	template := "{{DOC_TASK_AC}}"
	ctx := BodyContext{
		FeatureSlug: "feat",
		DocTaskCriteria: map[string]string{
			"3-doc": "AC 3 content",
			"1-doc": "AC 1 content",
			"2-doc": "AC 2 content",
		},
	}
	def := AutoGenTaskDef{}

	result := renderBody(template, def, ctx)

	// Keys should be sorted (1-doc before 2-doc before 3-doc)
	idx1 := strings.Index(result, "### 1-doc")
	idx2 := strings.Index(result, "### 2-doc")
	idx3 := strings.Index(result, "### 3-doc")
	if idx1 >= idx2 || idx2 >= idx3 {
		t.Errorf("keys should be sorted alphabetically, got indices: 1-doc=%d, 2-doc=%d, 3-doc=%d\n%s", idx1, idx2, idx3, result)
	}
}

// --- doc-review autogen template content tests (Task 2) ---

func TestDocReviewAutogenTemplate_ContainsDocTaskACPlaceholder(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("data/doc-review.md")
	if err != nil {
		t.Fatalf("cannot read doc-review.md: %v", err)
	}
	s := string(data)

	if !strings.Contains(s, "{{DOC_TASK_AC}}") {
		t.Error("doc-review autogen template must contain {{DOC_TASK_AC}} placeholder for AC summary injection")
	}
}

func TestDocReviewAutogenTemplate_ContainsACSummarySection(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("data/doc-review.md")
	if err != nil {
		t.Fatalf("cannot read doc-review.md: %v", err)
	}
	s := string(data)

	if !strings.Contains(s, "## Acceptance Criteria Summary") {
		t.Error("doc-review autogen template must contain '## Acceptance Criteria Summary' section header")
	}
	if !strings.Contains(s, "pre-extracted") {
		t.Error("doc-review autogen template should mention that AC is pre-extracted")
	}
}

func TestDocReviewAutogenTemplate_AllowlistDiscoveryStrategy(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("data/doc-review.md")
	if err != nil {
		t.Fatalf("cannot read doc-review.md: %v", err)
	}
	s := string(data)

	// Must use allowlist language
	if !strings.Contains(s, "allowlist") {
		t.Error("doc-review autogen template Discovery Strategy should use allowlist language")
	}
	// Must reference docs/ path
	if !strings.Contains(s, "docs/features/") {
		t.Error("doc-review autogen template should reference docs/features/ path")
	}
}

func TestDocReviewAutogenTemplate_ExcludesTasksAndRecords(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("data/doc-review.md")
	if err != nil {
		t.Fatalf("cannot read doc-review.md: %v", err)
	}
	s := string(data)

	// Must explicitly exclude tasks/, records/, manifest.md, index.json
	if !strings.Contains(s, "tasks/") {
		t.Error("doc-review autogen template should mention tasks/ exclusion")
	}
	if !strings.Contains(s, "records/") {
		t.Error("doc-review autogen template should mention records/ exclusion")
	}
	if !strings.Contains(s, "manifest.md") {
		t.Error("doc-review autogen template should mention manifest.md exclusion")
	}
	if !strings.Contains(s, "index.json") {
		t.Error("doc-review autogen template should mention index.json exclusion")
	}
}

func TestDocReviewAutogenTemplate_NoScanTasksDirective(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("data/doc-review.md")
	if err != nil {
		t.Fatalf("cannot read doc-review.md: %v", err)
	}
	s := string(data)

	// Must NOT contain old-style task scanning instructions
	if strings.Contains(s, "read its acceptance criteria from the task .md file") {
		t.Error("doc-review autogen template should NOT contain old 'read its acceptance criteria from the task .md file' directive")
	}
	if strings.Contains(s, "For each doc task, read its acceptance criteria") {
		t.Error("doc-review autogen template should NOT contain old 'For each doc task, read its acceptance criteria' directive")
	}
}

func TestRenderBody_AllPlaceholdersIncludingDocTaskAC(t *testing.T) {
	template := "Feature: {{FEATURE_SLUG}}\nMode: {{MODE}}\n## Scope\n{{SCOPE}}\n## Other\nInterfaces: {{SURFACES}}\nType: {{TEST_TYPE}}\nAcceptance:\n{{ACCEPTANCE_CRITERIA}}\n{{DOC_TASK_AC}}"
	ctx := BodyContext{
		FeatureSlug:        "feat",
		Mode:               "quick",
		Scope:              []string{"backend"},
		SurfaceTypes:       []string{"cli"},
		AcceptanceCriteria: []string{"AC1"},
		DocTaskCriteria:    map[string]string{"1-doc": "Doc AC"},
	}
	def := AutoGenTaskDef{TestType: "cli"}

	result := renderBody(template, def, ctx)

	allPlaceholders := []string{
		"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}",
		"{{SURFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}",
		"{{DOC_TASK_AC}}",
	}
	for _, ph := range allPlaceholders {
		if strings.Contains(result, ph) {
			t.Errorf("placeholder %s not resolved in output", ph)
		}
	}
	if !strings.Contains(result, "### 1-doc") {
		t.Errorf("should contain DocTaskCriteria sub-section, got:\n%s", result)
	}
}
