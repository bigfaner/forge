package task

import (
	"fmt"
	"strings"
	"testing"

	"forge-cli/pkg/forgeconfig"
)

// defaultAuto is the current default (consolidateSpecs quick=true, test quick=false).
var defaultAuto = forgeconfig.AutoConfigDefaults()

// allEnabledAuto enables all auto-behaviors for tests that need quick + full tasks.
var allEnabledAuto = forgeconfig.AutoConfig{
	Test:             forgeconfig.ModeToggle{Quick: true, Full: true},
	ConsolidateSpecs: forgeconfig.ModeToggle{Quick: true, Full: true},
	CleanCode:        forgeconfig.ModeToggle{Quick: false, Full: false},
}

// validationAuto enables validation + e2e for testing validate-ux gating.
var validationAuto = forgeconfig.AutoConfig{
	Test:       forgeconfig.ModeToggle{Quick: true, Full: true},
	Validation: forgeconfig.ModeToggle{Quick: true, Full: true},
}

// scalarSurface creates a single-surface (scalar form) surfaces map for tests.
func scalarSurface(typ string) map[string]string {
	return map[string]string{".": typ}
}

// multiSurface creates a multi-surface map from alternating key-type pairs.
func multiSurface(keyTypes ...string) map[string]string {
	m := make(map[string]string, len(keyTypes)/2)
	for i := 0; i+1 < len(keyTypes); i += 2 {
		m[keyTypes[i]] = keyTypes[i+1]
	}
	return m
}

func TestGetBreakdownTestTasks_EmptyInterfaces(t *testing.T) {
	tasks := GetBreakdownTestTasks(nil, nil, defaultAuto, "")

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks with empty interfaces, got %d", len(tasks))
	}
}

func TestGetBreakdownTestTasks_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

	// gen-journeys + eval-journey + gen-contracts + eval-contract + gen-scripts-cli + run + consolidate = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	wantIDs := []string{"T-test-gen-journeys", "T-eval-journey", "T-test-gen-contracts", "T-eval-contract", "T-test-gen-scripts-cli", "T-test-run", "T-specs-consolidate"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// Dependency chain: eval-journey -> gen-journeys, gen-contracts -> eval-journey, eval-contract -> gen-contracts, gen-scripts -> eval-contract, run -> gen-scripts, consolidate -> run
	if tasks[1].Dependencies[0] != "T-test-gen-journeys" {
		t.Errorf("eval-journey should depend on gen-journeys, got %v", tasks[1].Dependencies)
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
		t.Errorf("consolidate should depend on run, got %v", tasks[6].Dependencies)
	}
}

func TestGetQuickTestTasks_EmptyInterfaces(t *testing.T) {
	tasks := GetQuickTestTasks(nil, nil, allEnabledAuto, "")

	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks with empty interfaces, got %d", len(tasks))
	}
}

func TestGetQuickTestTasks_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "")

	// gen-journeys + run + drift = 3
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys",
		"T-test-run",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if tasks[0].Type != TypeTestGenJourneys {
		t.Errorf("T-test-gen-journeys Type = %q, want %q", tasks[0].Type, TypeTestGenJourneys)
	}

	if tasks[2].Type != TypeDocDrift {
		t.Errorf("T-quick-doc-drift Type = %q, want %q", tasks[2].Type, TypeDocDrift)
	}
}

func TestGenerateTestTaskMD(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{},
		Type:        GenSurfaceTestType(TypeTestGenScripts, "api"),
		SurfaceType: "api", StrategyKind: "generate",
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)

	if !strings.Contains(s, `id: "T-test-gen-scripts-api"`) {
		t.Error("missing id in frontmatter")
	}
	if !strings.Contains(s, `type: "test.gen-scripts.api"`) {
		t.Error("missing type in frontmatter")
	}
	// Body now loaded from embed template, should contain strategy-based content
	if !strings.Contains(s, "executable test scripts") {
		t.Error("body should contain strategy-based content from embed template")
	}
}

func TestGenerateTestTaskMD_SharedTask(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-run", Key: "run-test",
		Title: "Run e2e Tests", Priority: "P1",
		EstimatedTime: "30min-1h", Dependencies: []string{},
		Type: TypeTestRun,
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
		tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")
		ResolveFirstTestDep(tasks, existing, "breakdown", "")
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
		tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")
		ResolveFirstTestDep(tasks, existing, "breakdown", "")
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
		tasks := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "")
		ResolveFirstTestDep(tasks, existing, "quick", "")
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
	if len(task.Dependencies) != 0 {
		t.Errorf("Dependencies should be empty (resolved later), got %v", task.Dependencies)
	}
}

func TestResolveReviewDocDep(t *testing.T) {
	t.Run("depends on all doc-type tasks, not just the last", func(t *testing.T) {
		existing := map[string]Task{
			"1-doc":              {ID: "1.1", Type: TypeDoc},
			"2-doc":              {ID: "1.2", Type: TypeDoc},
			"3-code":             {ID: "2.1", Type: TypeCodingFeature},
			"4-doc":              {ID: "3.1", Type: TypeDoc},
			"T-test-gen-scripts": {ID: "T-test-gen-scripts-cli", Type: TypeTestGenScripts},
		}
		task := GetReviewDocTask()
		ResolveReviewDocDep(&task, existing)

		if len(task.Dependencies) != 3 {
			t.Fatalf("Dependencies = %v, want exactly 3 doc tasks", task.Dependencies)
		}

		depSet := make(map[string]bool)
		for _, dep := range task.Dependencies {
			depSet[dep] = true
		}

		if !depSet["1.1"] || !depSet["1.2"] || !depSet["3.1"] {
			t.Errorf("should depend on all doc tasks {1.1, 1.2, 3.1}, got %v", task.Dependencies)
		}
		if depSet["2.1"] {
			t.Error("should NOT depend on non-doc task 2.1")
		}
	})

	t.Run("excludes non-doc business tasks", func(t *testing.T) {
		existing := map[string]Task{
			"1-code": {ID: "1", Type: TypeCodingFeature},
			"2-code": {ID: "2", Type: TypeCodingEnhancement},
			"3-fix":  {ID: "3", Type: TypeCodingFix},
		}
		task := GetReviewDocTask()
		ResolveReviewDocDep(&task, existing)

		if len(task.Dependencies) != 0 {
			t.Errorf("should have no deps when no doc tasks exist, got %v", task.Dependencies)
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

	t.Run("single doc task", func(t *testing.T) {
		existing := map[string]Task{
			"1-doc": {ID: "1.1", Type: TypeDoc},
		}
		task := GetReviewDocTask()
		ResolveReviewDocDep(&task, existing)

		if len(task.Dependencies) != 1 || task.Dependencies[0] != "1.1" {
			t.Errorf("should depend on the single doc task, got %v", task.Dependencies)
		}
	})
}

func TestGetBreakdownTestTasks_PerType_TwoTypes(t *testing.T) {
	tasks := GetBreakdownTestTasks(multiSurface("tui", "tui", "api", "api", ""), []string{"api", "tui"}, defaultAuto, "")

	// gen-journeys + eval-journey + gen-contracts + eval-contract + 2 gen-scripts(tui,api) + 2 run-test(api,tui) + consolidate = 9
	if len(tasks) != 9 {
		t.Fatalf("expected 9 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys",
		"T-eval-journey",
		"T-test-gen-contracts",
		"T-eval-contract",
		"T-test-gen-scripts-api", "T-test-gen-scripts-tui",
		"T-test-run-api", "T-test-run-tui",
		"T-specs-consolidate",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	// gen-journeys key has no type suffix (single task)
	if tasks[0].Key != "gen-journeys" {
		t.Errorf("tasks[0].Key = %q, want gen-journeys", tasks[0].Key)
	}

	// gen-scripts keys include type suffix (sorted alphabetically)
	if tasks[4].Key != "gen-test-scripts-api" {
		t.Errorf("tasks[4].Key = %q, want gen-test-scripts-api", tasks[4].Key)
	}
	if tasks[5].Key != "gen-test-scripts-tui" {
		t.Errorf("tasks[5].Key = %q, want gen-test-scripts-tui", tasks[5].Key)
	}

	// SurfaceType empty for single gen-journeys task
	if tasks[0].SurfaceType != "" {
		t.Errorf("tasks[0].SurfaceType = %q, want empty (single gen-journeys)", tasks[0].SurfaceType)
	}

	// TestType field set for gen-scripts (sorted alphabetically)
	if tasks[4].SurfaceType != "api" {
		t.Errorf("tasks[4].SurfaceType = %q, want api", tasks[4].SurfaceType)
	}
	if tasks[5].SurfaceType != "tui" {
		t.Errorf("tasks[5].SurfaceType = %q, want tui", tasks[5].SurfaceType)
	}

	// run-test tasks have surface-key and surface-type set
	if tasks[6].SurfaceKey != "api" {
		t.Errorf("tasks[6].SurfaceKey = %q, want api", tasks[6].SurfaceKey)
	}
	if tasks[7].SurfaceKey != "tui" {
		t.Errorf("tasks[7].SurfaceKey = %q, want tui", tasks[7].SurfaceKey)
	}

	// eval-journey depends on single gen-journeys task
	if len(tasks[1].Dependencies) != 1 {
		t.Fatalf("T-eval-journey should depend on 1 gen-journeys task, got %v", tasks[1].Dependencies)
	}
	if tasks[1].Dependencies[0] != "T-test-gen-journeys" {
		t.Errorf("T-eval-journey deps should be T-test-gen-journeys, got %v", tasks[1].Dependencies)
	}

	// T-test-run-api (first in chain) depends on ALL per-type gen-scripts tasks
	if len(tasks[6].Dependencies) != 2 {
		t.Fatalf("T-test-run-api should depend on 2 gen tasks, got %v", tasks[6].Dependencies)
	}
	depSet2 := make(map[string]bool)
	for _, d := range tasks[6].Dependencies {
		depSet2[d] = true
	}
	if !depSet2["T-test-gen-scripts-tui"] || !depSet2["T-test-gen-scripts-api"] {
		t.Errorf("T-test-run-api deps should include both gen-scripts, got %v", tasks[6].Dependencies)
	}

	// T-test-run-tui (second in chain) depends on T-test-run-api (serial)
	if len(tasks[7].Dependencies) != 1 || tasks[7].Dependencies[0] != "T-test-run-api" {
		t.Errorf("T-test-run-tui should depend on T-test-run-api (serial), got %v", tasks[7].Dependencies)
	}

	// T-specs-consolidate depends on last run-test (T-test-run-tui)
	if tasks[8].Dependencies[0] != "T-test-run-tui" {
		t.Errorf("consolidate should depend on T-test-run-tui (chain tail), got %v", tasks[8].Dependencies)
	}
}

func TestGetBreakdownTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("api"), nil, defaultAuto, "")

	// gen-journeys + eval-journey + gen-contracts + eval-contract + gen-scripts-api + run + consolidate = 7
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys",
		"T-eval-journey",
		"T-test-gen-contracts",
		"T-eval-contract",
		"T-test-gen-scripts-api",
		"T-test-run", "T-specs-consolidate",
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
	tasks := GetBreakdownTestTasks(multiSurface("tui", "tui", "api", "api", "cli", "cli", ""), []string{"api", "tui", "cli"}, defaultAuto, "")

	// gen-journeys + eval-journey + gen-contracts + eval-contract + 3 gen-scripts + 3 run-tests + consolidate = 11
	if len(tasks) != 11 {
		t.Fatalf("expected 12 tasks, got %d", len(tasks))
	}

	// T-test-run-api (first in chain) depends on all 3 gen tasks
	runAPIIdx := findTaskIndex(tasks, "T-test-run-api")
	if runAPIIdx < 0 {
		t.Fatalf("T-test-run-api not found in tasks")
	}
	if len(tasks[runAPIIdx].Dependencies) != 3 {
		t.Fatalf("T-test-run-api should depend on 3 gen tasks, got %v", tasks[runAPIIdx].Dependencies)
	}
	depSet := make(map[string]bool)
	for _, d := range tasks[runAPIIdx].Dependencies {
		depSet[d] = true
	}
	if !depSet["T-test-gen-scripts-tui"] || !depSet["T-test-gen-scripts-api"] || !depSet["T-test-gen-scripts-cli"] {
		t.Errorf("T-test-run-api missing expected deps, got %v", tasks[runAPIIdx].Dependencies)
	}

	// Serial chain: api -> tui -> cli
	runTUIIdx := findTaskIndex(tasks, "T-test-run-tui")
	if len(tasks[runTUIIdx].Dependencies) != 1 || tasks[runTUIIdx].Dependencies[0] != "T-test-run-api" {
		t.Errorf("T-test-run-tui should depend on T-test-run-api, got %v", tasks[runTUIIdx].Dependencies)
	}
	runCLIIdx := findTaskIndex(tasks, "T-test-run-cli")
	if len(tasks[runCLIIdx].Dependencies) != 1 || tasks[runCLIIdx].Dependencies[0] != "T-test-run-tui" {
		t.Errorf("T-test-run-cli should depend on T-test-run-tui, got %v", tasks[runCLIIdx].Dependencies)
	}

	// Consolidate-specs depends on last run-test (cli)
	consolidateIdx := findTaskIndex(tasks, "T-specs-consolidate")
	if consolidateIdx < 0 {
		t.Fatalf("T-specs-consolidate not found in tasks")
	}
	if tasks[consolidateIdx].Dependencies[0] != "T-test-run-cli" {
		t.Errorf("consolidate-specs should depend on last run-test (after run-test chain), got %v", tasks[consolidateIdx].Dependencies)
	}
}

func TestGenerateTestTaskMD_WithTestType(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Dependencies: []string{},
		Type:        GenSurfaceTestType(TypeTestGenScripts, "api"),
		SurfaceType: "api", StrategyKind: "generate",
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
	tasks := GetQuickTestTasks(multiSurface("tui", "tui", "api", "api", ""), []string{"api", "tui"}, allEnabledAuto, "")

	// gen-journeys + 2 run-tests(api,tui) + drift = 4
	if len(tasks) != 4 {
		t.Fatalf("expected 4 tasks, got %d", len(tasks))
	}

	// Fixed-position tasks (before map-iteration-dependent ones)
	if tasks[0].ID != "T-test-gen-journeys" {
		t.Errorf("tasks[0].ID = %q, want T-test-gen-journeys", tasks[0].ID)
	}

	// Map-iteration-dependent tasks: run-tests use byID lookups
	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Both run-test tasks exist with correct SurfaceKey
	if _, ok := byID["T-test-run-api"]; !ok {
		t.Error("missing T-test-run-api task")
	}
	if _, ok := byID["T-test-run-tui"]; !ok {
		t.Error("missing T-test-run-tui task")
	}

	// drift at fixed position
	if tasks[3].ID != "T-quick-doc-drift" {
		t.Errorf("tasks[3].ID = %q, want T-quick-doc-drift", tasks[3].ID)
	}

	// gen-journeys key has no type suffix (single task)
	if tasks[0].Key != "gen-journeys" {
		t.Errorf("tasks[0].Key = %q, want gen-journeys", tasks[0].Key)
	}

	// SurfaceType empty for single gen-journeys task
	if tasks[0].SurfaceType != "" {
		t.Errorf("tasks[0].SurfaceType = %q, want empty (single gen-journeys)", tasks[0].SurfaceType)
	}
}

func TestGetQuickTestTasks_PerType_SingleType(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("api"), nil, allEnabledAuto, "")

	// gen-journeys + run + drift = 3
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	wantIDs := []string{
		"T-test-gen-journeys",
		"T-test-run",
		"T-quick-doc-drift",
	}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}

	if tasks[0].Type != TypeTestGenJourneys {
		t.Errorf("T-test-gen-journeys Type = %q, want %q", tasks[0].Type, TypeTestGenJourneys)
	}

	// T-test-run depends on gen-journeys (direct, no gen-contracts/gen-scripts in Quick)
	runIdx := findTaskIndexOrPanic(tasks, "T-test-run")
	if len(tasks[runIdx].Dependencies) != 1 || tasks[runIdx].Dependencies[0] != "T-test-gen-journeys" {
		t.Errorf("T-test-run should depend on T-test-gen-journeys, got %v", tasks[runIdx].Dependencies)
	}
}

func TestGetQuickTestTasks_PerType_ThreeTypes(t *testing.T) {
	tasks := GetQuickTestTasks(multiSurface("tui", "tui", "api", "api", "cli", "cli", ""), []string{"api", "tui", "cli"}, allEnabledAuto, "")

	// gen-journeys + 3 run-tests + drift = 5
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	// T-test-run-api (first in chain) depends on gen-journeys (not gen-scripts)
	runAPIIdx := findTaskIndexOrPanic(tasks, "T-test-run-api")
	if len(tasks[runAPIIdx].Dependencies) != 1 {
		t.Fatalf("T-test-run-api should depend on 1 gen-journeys, got %v", tasks[runAPIIdx].Dependencies)
	}
	if tasks[runAPIIdx].Dependencies[0] != "T-test-gen-journeys" {
		t.Errorf("T-test-run-api should depend on T-test-gen-journeys, got %v", tasks[runAPIIdx].Dependencies)
	}

	// Serial chain: api -> tui -> cli
	runTUIIdx := findTaskIndexOrPanic(tasks, "T-test-run-tui")
	if len(tasks[runTUIIdx].Dependencies) != 1 || tasks[runTUIIdx].Dependencies[0] != "T-test-run-api" {
		t.Errorf("T-test-run-tui should depend on T-test-run-api, got %v", tasks[runTUIIdx].Dependencies)
	}
	runCLIIdx := findTaskIndexOrPanic(tasks, "T-test-run-cli")
	if len(tasks[runCLIIdx].Dependencies) != 1 || tasks[runCLIIdx].Dependencies[0] != "T-test-run-tui" {
		t.Errorf("T-test-run-cli should depend on T-test-run-tui, got %v", tasks[runCLIIdx].Dependencies)
	}
}

// --- validate-ux should only be generated when interfaces include UI types ---

func TestGetBreakdownTestTasks_ValidateUx_SkippedForCLIOnly(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, validationAuto, "")

	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			t.Error("validate-ux should not be generated for CLI-only projects (no visual UI)")
		}
	}
}

func TestGetQuickTestTasks_ValidateUx_SkippedForCLIOnly(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, validationAuto, "")

	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			t.Error("validate-ux should not be generated for CLI-only projects (no visual UI)")
		}
	}
}

func TestGetBreakdownTestTasks_ValidateUx_SkippedForAPIOnly(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("api"), nil, validationAuto, "")

	for _, task := range tasks {
		if task.ID == "T-validate-ux" {
			t.Error("validate-ux should not be generated for API-only projects (no visual UI)")
		}
	}
}

func TestGetBreakdownTestTasks_ValidateUx_IncludedForTUI(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("tui"), nil, validationAuto, "")

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
	tasks := GetBreakdownTestTasks(multiSurface("cli", "cli", "tui", "tui", ""), []string{"tui", "cli"}, validationAuto, "")

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
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, validationAuto, "")

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
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

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
		{"run", TypeTestRun, "staged test scripts"},
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
				EstimatedTime: "1h", Type: tt.typ,
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
		EstimatedTime: "1-2h", Type: TypeTestGenScripts,
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
		EstimatedTime: "1-2h", Type: TypeTestGenScripts,
		SurfaceType: "api",
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

// --- Task 1.2b: TaskFromFile propagation tests ---

func TestTaskFromFile_PropagatesSurfaceKeyAndType(t *testing.T) {
	def := AutoGenTaskDef{
		ID:            "T-test-gen-scripts-api",
		Key:           "gen-test-scripts-api",
		Title:         "Generate Test Scripts (api)",
		Priority:      "P1",
		EstimatedTime: "1-2h",
		Dependencies:  []string{"dep1"},
		Type:          TypeTestGenScripts,
		SurfaceKey:    "admin-panel",
		SurfaceType:   "web",
		Breaking:      true,
		MainSession:   true,
	}

	task := def.TaskFromFile()

	if task.SurfaceKey != "admin-panel" {
		t.Errorf("Task.SurfaceKey = %q, want %q", task.SurfaceKey, "admin-panel")
	}
	if task.SurfaceType != "web" {
		t.Errorf("Task.SurfaceType = %q, want %q", task.SurfaceType, "web")
	}
	if task.ID != "T-test-gen-scripts-api" {
		t.Errorf("Task.ID = %q, want %q", task.ID, "T-test-gen-scripts-api")
	}
	if task.Type != TypeTestGenScripts {
		t.Errorf("Task.Type = %q, want %q", task.Type, TypeTestGenScripts)
	}
}

func TestTaskFromFile_EmptySurfaceFields(t *testing.T) {
	def := AutoGenTaskDef{
		ID:       "T-test-run",
		Key:      "run-test",
		Title:    "Run e2e Tests",
		Priority: "P1",
		Type:     TypeTestRun,
	}

	task := def.TaskFromFile()

	if task.SurfaceKey != "" {
		t.Errorf("Task.SurfaceKey = %q, want empty when not set", task.SurfaceKey)
	}
	if task.SurfaceType != "" {
		t.Errorf("Task.SurfaceType = %q, want empty when not set", task.SurfaceType)
	}
}

func TestGenerateTestTaskMD_FrontmatterUnchanged(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-run", Key: "run-test",
		Title: "Run e2e Tests", Priority: "P1",
		EstimatedTime: "30min-1h", Dependencies: []string{"dep1"},
		Type:        TypeTestRun,
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
	if !strings.Contains(s, `surface-key: ""`) {
		t.Error("missing surface-key in frontmatter")
	}
	if !strings.Contains(s, `surface-type: ""`) {
		t.Error("missing surface-type in frontmatter")
	}
	if !strings.Contains(s, "mainSession: true") {
		t.Error("missing mainSession in frontmatter")
	}
}

// --- autogenTemplateData and renderBody tests ---

func TestRenderBody_SubstitutesAllPlaceholders(t *testing.T) {
	tmpl := "Feature: {{.FeatureSlug}}{{if .Mode}}\nMode: {{.Mode}}{{end}}\nInterfaces: {{.SurfaceTypes}}{{if .SurfaceType}}\nTest Type: {{.SurfaceType}}{{end}}\nAcceptance:\n{{.AcceptanceCriteria}}"

	data := autogenTemplateData{
		FeatureSlug:        "my-feature",
		Mode:               "quick",
		SurfaceTypes:       "- api\n- cli",
		SurfaceType:        "api",
		AcceptanceCriteria: "- [ ] AC1: works\n- [ ] AC2: fast",
	}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	if !strings.Contains(result, "Feature: my-feature") {
		t.Error("FeatureSlug not substituted")
	}
	if !strings.Contains(result, "Mode: quick") {
		t.Error("Mode not substituted")
	}
	if !strings.Contains(result, "- api") {
		t.Errorf("SurfaceTypes not substituted, got:\n%s", result)
	}
	if !strings.Contains(result, "Test Type: api") {
		t.Error("SurfaceType not substituted")
	}
	if !strings.Contains(result, "- [ ] AC1: works") {
		t.Errorf("AcceptanceCriteria not substituted, got:\n%s", result)
	}
}

func TestRenderBody_EmptyMode_OmitsLine(t *testing.T) {
	tmpl := "Feature: {{.FeatureSlug}}{{if .Mode}}\nMode: {{.Mode}}{{end}}\nDone"
	data := autogenTemplateData{FeatureSlug: "test"}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	if strings.Contains(result, "Mode:") {
		t.Errorf("Mode line should be omitted when empty, got:\n%s", result)
	}
	if !strings.Contains(result, "Feature: test") {
		t.Error("FeatureSlug should still be present")
	}
}

func TestRenderBody_EmptySurfaceType_OmitsLine(t *testing.T) {
	tmpl := "Feature: {{.FeatureSlug}}{{if .SurfaceType}}\nType: {{.SurfaceType}}{{end}}\nDone"
	data := autogenTemplateData{FeatureSlug: "test"}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	if strings.Contains(result, "Type:") {
		t.Errorf("SurfaceType line should be omitted when empty, got:\n%s", result)
	}
}

func TestRenderBody_DefaultAcceptanceCriteria(t *testing.T) {
	tmpl := "Acceptance:\n{{.AcceptanceCriteria}}"
	data := autogenTemplateData{AcceptanceCriteria: "- [ ] All acceptance criteria met"}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	if !strings.Contains(result, "- [ ] All acceptance criteria met") {
		t.Errorf("default acceptance criteria not rendered, got:\n%s", result)
	}
}

func TestRenderBody_DefaultSurfaceTypes(t *testing.T) {
	tmpl := "Interfaces: {{.SurfaceTypes}}"
	data := autogenTemplateData{SurfaceTypes: "See .forge/config.yaml"}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	if !strings.Contains(result, "See .forge/config.yaml") {
		t.Errorf("default SurfaceTypes not rendered, got:\n%s", result)
	}
}

func TestRenderBody_NoTemplateMarkersLeft(t *testing.T) {
	tmpl := "Feature: {{.FeatureSlug}}{{if .Mode}}\nMode: {{.Mode}}{{end}}{{if .SurfaceKey}}\nScope: {{.SurfaceKey}}{{end}}\nInterfaces: {{.SurfaceTypes}}{{if .SurfaceType}}\nType: {{.SurfaceType}}{{end}}\n{{.AcceptanceCriteria}}"
	data := autogenTemplateData{
		FeatureSlug:        "",
		AcceptanceCriteria: "- [ ] All acceptance criteria met",
		SurfaceTypes:       "See .forge/config.yaml",
	}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	templateMarkers := []string{"{{.FeatureSlug}}", "{{.Mode}}", "{{.SurfaceKey}}", "{{.SurfaceTypes}}", "{{.SurfaceType}}", "{{.AcceptanceCriteria}}"}
	for _, ph := range templateMarkers {
		if strings.Contains(result, ph) {
			t.Errorf("template marker %s not resolved in output:\n%s", ph, result)
		}
	}
}

func TestGenerateTestTaskMD_WithBodyContext(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Type: TypeTestGenScripts,
		SurfaceType: "api",
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
	if !strings.Contains(s, `id: "T-test-gen-scripts-api"`) {
		t.Error("missing id in frontmatter")
	}
	if !strings.Contains(s, "executable test scripts") {
		t.Error("body should contain template content")
	}
}

func TestGenerateTestTaskMD_BackwardCompat_EmptyBodyContext(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-scripts-api", Key: "gen-test-scripts-api",
		Title: "Generate Test Scripts (api)", Priority: "P1",
		EstimatedTime: "1-2h", Type: TypeTestGenScripts,
	}

	content, err := GenerateTestTaskMD(def, BodyContext{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	s := string(content)
	managed := []string{
		"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}",
		"{{SURFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}",
		"{{.FeatureSlug}}", "{{.Mode}}", "{{.SurfaceKey}}",
		"{{.SurfaceTypes}}", "{{.SurfaceType}}", "{{.AcceptanceCriteria}}",
	}
	for _, ph := range managed {
		if strings.Contains(s, ph) {
			t.Errorf("placeholder %s should be resolved", ph)
		}
	}
	if !strings.Contains(s, `id: "T-test-gen-scripts-api"`) {
		t.Error("frontmatter should be intact")
	}
}

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
			tmpl := "Feature: {{.FeatureSlug}} is ready."
			data := autogenTemplateData{FeatureSlug: tt.featureSlug}

			result, err := renderBody(tmpl, data)
			if err != nil {
				t.Fatalf("renderBody error: %v", err)
			}

			if !strings.Contains(result, tt.want) {
				t.Errorf("expected %q in result, got: %s", tt.want, result)
			}
			if strings.Contains(result, "{{.FeatureSlug}}") {
				t.Error("FeatureSlug placeholder should be resolved")
			}
		})
	}
}

func TestRenderBody_AcceptanceCriteria(t *testing.T) {
	t.Run("criteria filled as checklist", func(t *testing.T) {
		tmpl := "## Validation Criteria\n{{.AcceptanceCriteria}}\n## End"
		data := autogenTemplateData{
			FeatureSlug:        "feat",
			AcceptanceCriteria: "- [ ] AC1: Login works\n- [ ] AC2: Logout works",
		}

		result, err := renderBody(tmpl, data)
		if err != nil {
			t.Fatalf("renderBody error: %v", err)
		}

		if !strings.Contains(result, "- [ ] AC1: Login works") {
			t.Error("first AC should be present")
		}
		if !strings.Contains(result, "- [ ] AC2: Logout works") {
			t.Error("second AC should be present")
		}
	})
}

func TestRenderBody_EmptyFields(t *testing.T) {
	t.Run("all fields empty uses fallbacks", func(t *testing.T) {
		tmpl := "Feature: {{.FeatureSlug}}{{if .Mode}}\nMode: {{.Mode}}{{end}}{{if .SurfaceKey}}\n## Scope\n{{.SurfaceKey}}{{end}}\n## Other\nInterfaces: {{.SurfaceTypes}}{{if .SurfaceType}}\nType: {{.SurfaceType}}{{end}}\n{{.AcceptanceCriteria}}"
		data := autogenTemplateData{
			AcceptanceCriteria: "- [ ] All acceptance criteria met",
			SurfaceTypes:       "See .forge/config.yaml",
		}

		result, err := renderBody(tmpl, data)
		if err != nil {
			t.Fatalf("renderBody error: %v", err)
		}

		if strings.Contains(result, "Mode:") {
			t.Error("Mode line should be omitted when empty")
		}
		if strings.Contains(result, "## Scope") {
			t.Error("Scope section should be omitted when empty")
		}
		if !strings.Contains(result, "See .forge/config.yaml") {
			t.Error("Empty SurfaceTypes should use fallback")
		}
		if strings.Contains(result, "Type:") {
			t.Error("SurfaceType line should be omitted when empty")
		}
		if !strings.Contains(result, "- [ ] All acceptance criteria met") {
			t.Error("Empty AcceptanceCriteria should use fallback")
		}
		for _, ph := range []string{"{{.FeatureSlug}}", "{{.Mode}}", "{{.SurfaceKey}}", "{{.SurfaceTypes}}", "{{.SurfaceType}}", "{{.AcceptanceCriteria}}"} {
			if strings.Contains(result, ph) {
				t.Errorf("template marker %s should be resolved", ph)
			}
		}
	})
}

func TestBodyContentPerStrategy(t *testing.T) {
	tests := []struct {
		name         string
		typ          string
		testType     string
		ctx          BodyContext
		wantContains []string
	}{
		{"gen-scripts has feature context", TypeTestGenScripts, "api", BodyContext{
			FeatureSlug: "feat",
		}, []string{"feat", "api"}},
		{"run has feature context", TypeTestRun, "", BodyContext{
			FeatureSlug: "feat",
		}, []string{"feat"}},
		{"validation-code has criteria", TypeValidationCode, "", BodyContext{
			FeatureSlug: "feat", AcceptanceCriteria: []string{"AC1: works", "AC2: fast"},
		}, []string{"feat", "- [ ] AC1: works", "- [ ] AC2: fast"}},
		{"validation-ux has criteria", TypeValidationUx, "", BodyContext{
			FeatureSlug: "feat", AcceptanceCriteria: []string{"AC1: accessible"},
		}, []string{"feat", "- [ ] AC1: accessible"}},
		{"doc-review has discovery strategy", TypeDocReview, "", BodyContext{
			FeatureSlug: "feat", Mode: "breakdown",
		}, []string{"feat", "breakdown mode"}},
		{"doc-consolidate has discovery strategy", TypeDocConsolidate, "", BodyContext{
			FeatureSlug: "feat",
		}, []string{"feat", "Discovery Strategy"}},
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
				EstimatedTime: "1h", Type: tt.typ,
				SurfaceType: tt.testType,
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

			for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{SURFACES}}", "{{TEST_TYPE}}", "{{ACCEPTANCE_CRITERIA}}"} {
				if strings.Contains(s, ph) {
					t.Errorf("type %q has unresolved old-style placeholder %s", tt.typ, ph)
				}
			}
		})
	}
}

func TestGenJourneysTemplateContent(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("templates/test-gen-journeys.md")
	if err != nil {
		t.Fatalf("cannot read test-gen-journeys.md: %v", err)
	}
	s := string(data)

	// AC1: uses {{.FeatureSlug}} and {{.Mode}} placeholders
	if !strings.Contains(s, "{{.FeatureSlug}}") {
		t.Error("template must contain {{.FeatureSlug}} placeholder")
	}
	if !strings.Contains(s, "{{.Mode}}") {
		t.Error("template must contain {{.Mode}} placeholder")
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
	data, err := autogenTemplateFS.ReadFile("templates/test-gen-contracts.md")
	if err != nil {
		t.Fatalf("cannot read test-gen-contracts.md: %v", err)
	}
	s := string(data)

	// AC2: uses {{.FeatureSlug}} and {{.Mode}} placeholders
	if !strings.Contains(s, "{{.FeatureSlug}}") {
		t.Error("template must contain {{.FeatureSlug}} placeholder")
	}
	if !strings.Contains(s, "{{.Mode}}") {
		t.Error("template must contain {{.Mode}} placeholder")
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
		EstimatedTime: "20-30min", Type: TypeTestGenJourneys,
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

	// Placeholders resolved
	if strings.Contains(s, "{{.FeatureSlug}}") {
		t.Error("{{.FeatureSlug}} should be resolved")
	}
	if strings.Contains(s, "{{.Mode}}") {
		t.Error("{{.Mode}} should be resolved")
	}
	// Feature slug present in body
	if !strings.Contains(s, "my-feature") {
		t.Error("rendered body should contain feature slug")
	}
	// Mode present
	if !strings.Contains(s, "quick") {
		t.Error("rendered body should contain mode")
	}
	// No old-style or new-style placeholders left
	for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{TEST_TYPE}}", "{{.SurfaceKey}}"} {
		if strings.Contains(s, ph) {
			t.Errorf("placeholder %s should be resolved", ph)
		}
	}
}

func TestGenContractsTemplateRendering(t *testing.T) {
	def := AutoGenTaskDef{
		ID: "T-test-gen-contracts", Key: "gen-contracts",
		Title: "Generate Test Contracts", Priority: "P1",
		EstimatedTime: "30-45min", Type: TypeTestGenContracts,
	}
	ctx := BodyContext{
		FeatureSlug: "my-feature",
		Mode:        "breakdown",
	}

	content, err := GenerateTestTaskMD(def, ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(content)

	// Placeholders resolved
	if strings.Contains(s, "{{.FeatureSlug}}") {
		t.Error("{{.FeatureSlug}} should be resolved")
	}
	if strings.Contains(s, "{{.Mode}}") {
		t.Error("{{.Mode}} should be resolved")
	}
	// Feature slug present
	if !strings.Contains(s, "my-feature") {
		t.Error("rendered body should contain feature slug")
	}
	// No old-style or new-style placeholders left
	for _, ph := range []string{"{{FEATURE_SLUG}}", "{{MODE}}", "{{SCOPE}}", "{{TEST_TYPE}}", "{{.SurfaceKey}}"} {
		if strings.Contains(s, ph) {
			t.Errorf("placeholder %s should be resolved", ph)
		}
	}
}

func TestAutogenTemplateDiscovery(t *testing.T) {
	// Verify all auto-gen types resolve to a readable template via naming convention
	wantTypes := []string{
		TypeTestGenScripts, TypeTestRun,
		TypeValidationCode, TypeValidationUx,
		TypeDocReview, TypeDocConsolidate, TypeDocDrift, TypeCleanCode,
		TypeTestGenJourneys, TypeTestGenContracts,
	}

	for _, typ := range wantTypes {
		file := autogenTemplatePath(typ)
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
	tasks := GetBreakdownTestTasks(multiSurface("tui", "tui", "api", "api", ""), []string{"api", "tui"}, defaultAuto, "")

	foundGenJourneys := false
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys" {
			foundGenJourneys = true
			if task.Type != TypeTestGenJourneys {
				t.Errorf("gen-journeys Type = %q, want %q", task.Type, TypeTestGenJourneys)
			}
			if task.SurfaceType != "" {
				t.Errorf("gen-journeys TestType = %q, want empty (single task)", task.SurfaceType)
			}
			if task.StrategyKind != "interface" {
				t.Errorf("gen-journeys StrategyKind = %q, want interface", task.StrategyKind)
			}
		}
	}
	if !foundGenJourneys {
		t.Error("missing T-test-gen-journeys task")
	}
}

func TestGetBreakdownTestTasks_GenContracts(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

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
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

	wantOrder := []string{
		"T-test-gen-journeys",
		"T-eval-journey",
		"T-test-gen-contracts",
		"T-eval-contract",
		"T-test-gen-scripts-cli",
		"T-test-run",
		"T-specs-consolidate",
	}

	for i, want := range wantOrder {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}
}

func TestGetBreakdownTestTasks_FullDependencyChain(t *testing.T) {
	tasks := GetBreakdownTestTasks(multiSurface("cli", "cli", "api", "api", ""), []string{"api", "cli"}, defaultAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// gen-journeys task has no deps (pipeline entry point)
	if len(byID["T-test-gen-journeys"].Dependencies) != 0 {
		t.Errorf("gen-journeys should have no deps, got %v", byID["T-test-gen-journeys"].Dependencies)
	}

	// eval-journey depends on single gen-journeys
	evalJourneyDeps := byID["T-eval-journey"].Dependencies
	if len(evalJourneyDeps) != 1 {
		t.Fatalf("eval-journey should depend on 1 gen-journeys, got %v", evalJourneyDeps)
	}
	if evalJourneyDeps[0] != "T-test-gen-journeys" {
		t.Errorf("eval-journey deps should be T-test-gen-journeys, got %v", evalJourneyDeps)
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

	// T-test-run-api (first in chain) depends on all gen-scripts
	runAPIDeps := byID["T-test-run-api"].Dependencies
	if len(runAPIDeps) != 2 {
		t.Fatalf("T-test-run-api should depend on 2 gen-scripts, got %v", runAPIDeps)
	}
	runDepSet := make(map[string]bool)
	for _, d := range runAPIDeps {
		runDepSet[d] = true
	}
	if !runDepSet["T-test-gen-scripts-cli"] || !runDepSet["T-test-gen-scripts-api"] {
		t.Errorf("T-test-run-api deps should include both gen-scripts, got %v", runAPIDeps)
	}

	// Serial chain: T-test-run-cli depends on T-test-run-api
	if len(byID["T-test-run-cli"].Dependencies) != 1 || byID["T-test-run-cli"].Dependencies[0] != "T-test-run-api" {
		t.Errorf("T-test-run-cli should depend on T-test-run-api (serial), got %v", byID["T-test-run-cli"].Dependencies)
	}

	// consolidate-specs depends on last run-test (T-test-run-cli)
	if len(byID["T-specs-consolidate"].Dependencies) != 1 || byID["T-specs-consolidate"].Dependencies[0] != "T-test-run-cli" {
		t.Errorf("consolidate-specs should depend on T-test-run-cli (chain tail), got %v", byID["T-specs-consolidate"].Dependencies)
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
		ID: "T-test-gen-journeys", Key: "gen-journeys",
		Title: "Generate Test Journeys", Priority: "P1",
		EstimatedTime: "20-30min", Type: TypeTestGenJourneys,
		StrategyKind: "interface",
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
		EstimatedTime: "30-45min", Type: TypeTestGenContracts,
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
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

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
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// consolidate-specs depends on T-test-run (last run-test)
	if byID["T-specs-consolidate"].Dependencies[0] != "T-test-run" {
		t.Errorf("consolidate-specs should depend on T-test-run, got %v", byID["T-specs-consolidate"].Dependencies)
	}
	if byID["T-test-run"].Dependencies[0] != "T-test-gen-scripts-cli" {
		t.Errorf("run should still depend on gen-scripts-cli, got %v", byID["T-test-run"].Dependencies)
	}
}

// --- Quick mode staged across types topology tests (Task 4) ---

func TestGetQuickTestTasks_StagedPipelineTypesOnly(t *testing.T) {
	tasks := GetQuickTestTasks(multiSurface("cli", "cli", "api", "api", ""), []string{"api", "cli"}, allEnabledAuto, "")

	// Quick mode should only generate test pipeline or doc task types
	validPrefixes := []string{"test.", "doc.", "code-quality."}
	for _, task := range tasks {
		valid := false
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(task.Type, prefix) {
				valid = true
				break
			}
		}
		if !valid {
			t.Errorf("Quick mode generated unexpected task type %q (type=%q)", task.ID, task.Type)
		}
	}
}

func TestGetQuickTestTasks_StagedAcrossTypesDependencyChain(t *testing.T) {
	tasks := GetQuickTestTasks(multiSurface("cli", "cli", "api", "api", ""), []string{"api", "cli"}, allEnabledAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Stage 1: gen-journeys has no deps (pipeline entry point)
	if len(byID["T-test-gen-journeys"].Dependencies) != 0 {
		t.Errorf("gen-journeys should have no deps, got %v", byID["T-test-gen-journeys"].Dependencies)
	}

	// Stage 2: first run-test depends on gen-journeys (no gen-contracts/gen-scripts in Quick mode)
	runAPIDeps := byID["T-test-run-api"].Dependencies
	if len(runAPIDeps) != 1 {
		t.Fatalf("T-test-run-api should depend on 1 gen-journeys, got %v", runAPIDeps)
	}
	if runAPIDeps[0] != "T-test-gen-journeys" {
		t.Errorf("T-test-run-api deps should be T-test-gen-journeys, got %v", runAPIDeps)
	}

	// Serial chain: T-test-run-cli depends on T-test-run-api
	if len(byID["T-test-run-cli"].Dependencies) != 1 || byID["T-test-run-cli"].Dependencies[0] != "T-test-run-api" {
		t.Errorf("T-test-run-cli should depend on T-test-run-api (serial), got %v", byID["T-test-run-cli"].Dependencies)
	}

	// Drift depends on last run-test (T-test-run-cli)
	if len(byID["T-quick-doc-drift"].Dependencies) != 1 || byID["T-quick-doc-drift"].Dependencies[0] != "T-test-run-cli" {
		t.Errorf("drift should depend on T-test-run-cli (chain tail), got %v", byID["T-quick-doc-drift"].Dependencies)
	}
}
func TestGetQuickTestTasks_GenJourneysPerType(t *testing.T) {
	tasks := GetQuickTestTasks(multiSurface("tui", "tui", "api", "api", ""), []string{"api", "tui"}, allEnabledAuto, "")

	foundGenJourneys := false
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys" {
			foundGenJourneys = true
			if task.Type != TypeTestGenJourneys {
				t.Errorf("gen-journeys Type = %q, want %q", task.Type, TypeTestGenJourneys)
			}
			if task.SurfaceType != "" {
				t.Errorf("gen-journeys TestType = %q, want empty (single task)", task.SurfaceType)
			}
		}
	}
	if !foundGenJourneys {
		t.Error("missing T-test-gen-journeys task in Quick mode")
	}
}

func TestGetQuickTestTasks_NoGenContractsOrScripts(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "")

	// Quick mode should NOT generate gen-contracts or gen-scripts
	for _, task := range tasks {
		if task.ID == "T-test-gen-contracts" {
			t.Error("Quick mode should not generate T-test-gen-contracts")
		}
		if strings.HasPrefix(task.ID, "T-test-gen-scripts-") {
			t.Errorf("Quick mode should not generate %q", task.ID)
		}
	}
}

func TestGetQuickTestTasks_NoHardcodedIndices(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "")

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

func TestGetQuickTestTasks_DriftDependsOnLastRunTest(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "")

	for _, task := range tasks {
		if task.ID == "T-quick-doc-drift" {
			if len(task.Dependencies) != 1 || task.Dependencies[0] != "T-test-run" {
				t.Errorf("T-quick-doc-drift should depend on last run-test (after run-test), got %v", task.Dependencies)
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
	ResolveFirstTestDep(tasks, existing, "breakdown", "")
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
	ResolveFirstTestDep(tasks, existing, "quick", "")
}

func TestResolveFirstTestDep_BreakdownWithCleanCode(t *testing.T) {
	existing := map[string]Task{
		"1-gate":  {ID: "1.gate"},
		"1.1-foo": {ID: "1.1"},
	}
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

	// Add a clean-code task
	tasks = append([]AutoGenTaskDef{{ID: "T-clean-code"}}, tasks...)

	ResolveFirstTestDep(tasks, existing, "breakdown", "")

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
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "")

	// Add a clean-code task
	tasks = append([]AutoGenTaskDef{{ID: "T-clean-code"}}, tasks...)

	ResolveFirstTestDep(tasks, existing, "quick", "")

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
	ResolveFirstTestDep(nil, map[string]Task{"1": {ID: "1"}}, "breakdown", "")
	ResolveFirstTestDep(nil, map[string]Task{"1": {ID: "1"}}, "quick", "")
}

func TestResolveFirstTestDep_NoDeps_NoPanic(t *testing.T) {
	// No existing business tasks → return without panic
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")
	ResolveFirstTestDep(tasks, map[string]Task{}, "breakdown", "")

	// gen-journeys should have no deps set (no business tasks to depend on)
	firstTestIdx := findTaskIndexByPrefix(tasks, "T-test-gen-journeys")
	if firstTestIdx >= 0 && len(tasks[firstTestIdx].Dependencies) != 0 {
		t.Errorf("gen-journeys should have no deps when no business tasks exist, got %v", tasks[firstTestIdx].Dependencies)
	}
}

func TestResolveDriftFallbackDep(t *testing.T) {
	t.Run("quick drift with no deps falls back to last business task", func(t *testing.T) {
		index := NewTaskIndex("test-feature")
		index.SetTasks(map[string]Task{
			"1-foo": {ID: "1", Title: "Foo", Priority: "P1", Status: "pending", Dependencies: []string{}, Type: TypeCodingFeature},
			"drift": {ID: "T-quick-doc-drift", Title: "Drift", Priority: "P2", Status: "pending", Dependencies: nil, Type: TypeDocDrift},
		})

		ResolveDriftFallbackDep(index)

		driftTask, _ := index.ByID("T-quick-doc-drift")
		if len(driftTask.Dependencies) != 1 || driftTask.Dependencies[0] != "1" {
			t.Errorf("drift deps = %v, want [1]", driftTask.Dependencies)
		}
	})

	t.Run("breakdown consolidate with no deps falls back to last business task", func(t *testing.T) {
		index := NewTaskIndex("test-feature")
		index.SetTasks(map[string]Task{
			"2-bar":       {ID: "2", Title: "Bar", Priority: "P1", Status: "pending", Dependencies: []string{}, Type: TypeCodingFeature},
			"consolidate": {ID: "T-specs-consolidate", Title: "Consolidate", Priority: "P2", Status: "pending", Dependencies: nil, Type: TypeDocConsolidate},
		})

		ResolveDriftFallbackDep(index)

		ct, _ := index.ByID("T-specs-consolidate")
		if len(ct.Dependencies) != 1 || ct.Dependencies[0] != "2" {
			t.Errorf("consolidate deps = %v, want [2]", ct.Dependencies)
		}
	})

	t.Run("no-op when drift already has deps", func(t *testing.T) {
		index := NewTaskIndex("test-feature")
		index.SetTasks(map[string]Task{
			"1-foo": {ID: "1", Title: "Foo", Priority: "P1", Status: "pending", Dependencies: []string{}, Type: TypeCodingFeature},
			"drift": {ID: "T-quick-doc-drift", Title: "Drift", Priority: "P2", Status: "pending", Dependencies: []string{"T-test-run"}, Type: TypeDocDrift},
		})

		ResolveDriftFallbackDep(index)

		dt, _ := index.ByID("T-quick-doc-drift")
		if len(dt.Dependencies) != 1 || dt.Dependencies[0] != "T-test-run" {
			t.Errorf("existing deps should be preserved, got %v", dt.Dependencies)
		}
	})

	t.Run("no-op when no business tasks", func(t *testing.T) {
		index := NewTaskIndex("test-feature")
		index.SetTasks(map[string]Task{
			"drift": {ID: "T-quick-doc-drift", Title: "Drift", Priority: "P2", Status: "pending", Dependencies: nil, Type: TypeDocDrift},
		})

		ResolveDriftFallbackDep(index)

		dt, _ := index.ByID("T-quick-doc-drift")
		if len(dt.Dependencies) != 0 {
			t.Errorf("drift should have no deps when no business tasks, got %v", dt.Dependencies)
		}
	})
}

// --- {{DOC_TASK_AC}} rendering tests ---

func TestRenderBody_DocTaskAC_Populated(t *testing.T) {
	tmpl := "Feature: {{.FeatureSlug}}\n## Acceptance Criteria Summary\n{{.DocTaskCriteria}}\n## End"
	data := autogenTemplateData{
		FeatureSlug:     "my-feature",
		DocTaskCriteria: "### 1-doc\n- [ ] AC 1\n- [ ] AC 2\n\n### 2-doc\n- [ ] AC 3",
	}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

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
	if strings.Contains(result, "{{.DocTaskCriteria}}") {
		t.Errorf("placeholder should be resolved, got:\n%s", result)
	}
}

func TestRenderBody_DocTaskAC_Empty(t *testing.T) {
	tmpl := "Feature: {{.FeatureSlug}}\n{{.DocTaskCriteria}}\n## End"
	data := autogenTemplateData{
		FeatureSlug:     "my-feature",
		DocTaskCriteria: "",
	}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	// Empty DocTaskCriteria should produce empty string
	if strings.Contains(result, "{{.DocTaskCriteria}}") {
		t.Errorf("placeholder should be resolved to empty when no criteria, got:\n%s", result)
	}
}

func TestRenderBody_DocTaskAC_SortedKeys(t *testing.T) {
	tmpl := "{{.DocTaskCriteria}}"
	data := autogenTemplateData{
		FeatureSlug:     "feat",
		DocTaskCriteria: "### 1-doc\nAC 1 content\n\n### 2-doc\nAC 2 content\n\n### 3-doc\nAC 3 content",
	}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

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
	data, err := autogenTemplateFS.ReadFile("templates/doc-review.md")
	if err != nil {
		t.Fatalf("cannot read doc-review.md: %v", err)
	}
	s := string(data)

	if !strings.Contains(s, "{{.DocTaskCriteria}}") {
		t.Error("doc-review autogen template must contain {{.DocTaskCriteria}} placeholder for AC summary injection")
	}
}

func TestDocReviewAutogenTemplate_ContainsACSummarySection(t *testing.T) {
	data, err := autogenTemplateFS.ReadFile("templates/doc-review.md")
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
	data, err := autogenTemplateFS.ReadFile("templates/doc-review.md")
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
	data, err := autogenTemplateFS.ReadFile("templates/doc-review.md")
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
	data, err := autogenTemplateFS.ReadFile("templates/doc-review.md")
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
	tmpl := "Feature: {{.FeatureSlug}}{{if .Mode}}\nMode: {{.Mode}}{{end}}{{if .SurfaceKey}}\n## Scope\n{{.SurfaceKey}}{{end}}\n## Other\nInterfaces: {{.SurfaceTypes}}{{if .SurfaceType}}\nType: {{.SurfaceType}}{{end}}\n{{.AcceptanceCriteria}}\n{{.DocTaskCriteria}}"
	data := autogenTemplateData{
		FeatureSlug:        "feat",
		Mode:               "quick",
		SurfaceTypes:       "- cli",
		SurfaceType:        "cli",
		AcceptanceCriteria: "- [ ] AC1",
		DocTaskCriteria:    "### 1-doc\nDoc AC",
	}

	result, err := renderBody(tmpl, data)
	if err != nil {
		t.Fatalf("renderBody error: %v", err)
	}

	allMarkers := []string{
		"{{.FeatureSlug}}", "{{.Mode}}", "{{.SurfaceKey}}",
		"{{.SurfaceTypes}}", "{{.SurfaceType}}", "{{.AcceptanceCriteria}}",
		"{{.DocTaskCriteria}}",
	}
	for _, ph := range allMarkers {
		if strings.Contains(result, ph) {
			t.Errorf("template marker %s not resolved in output", ph)
		}
	}
	if !strings.Contains(result, "### 1-doc") {
		t.Errorf("should contain DocTaskCriteria sub-section, got:\n%s", result)
	}
}

// --- Task 4: Split run-tests into per-surface-key serial tasks ---

// AC1: surfaces { frontend: web, backend: api } with no execution-order:
// ResolveExecutionOrder defaults to api < web, so backend runs before frontend.
func TestGetBreakdownTestTasks_DefaultExecutionOrder_BackendBeforeFrontend(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web")
	resolved, err := forgeconfig.ResolveExecutionOrder(surfaces, nil)
	if err != nil {
		t.Fatalf("ResolveExecutionOrder: %v", err)
	}

	tasks := GetBreakdownTestTasks(surfaces, resolved, defaultAuto, "")

	backendIdx := findTaskIndex(tasks, "T-test-run-backend")
	frontendIdx := findTaskIndex(tasks, "T-test-run-frontend")
	if backendIdx < 0 {
		t.Fatal("T-test-run-backend not found")
	}
	if frontendIdx < 0 {
		t.Fatal("T-test-run-frontend not found")
	}

	// Backend (api) must come before frontend (web) in task list
	if backendIdx >= frontendIdx {
		t.Errorf("T-test-run-backend (idx=%d) should come before T-test-run-frontend (idx=%d)", backendIdx, frontendIdx)
	}

	// Verify serial chain: frontend depends on backend
	if len(tasks[frontendIdx].Dependencies) != 1 || tasks[frontendIdx].Dependencies[0] != "T-test-run-backend" {
		t.Errorf("T-test-run-frontend should depend on T-test-run-backend (serial chain), got %v", tasks[frontendIdx].Dependencies)
	}
}

// AC2: Failure propagation - serial dependency chain means frontend is blocked when backend fails.
func TestGetBreakdownTestTasks_SerialChain_BlockedOnUpstreamFailure(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetBreakdownTestTasks(surfaces, resolved, defaultAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	backendDeps := byID["T-test-run-backend"].Dependencies
	frontendDeps := byID["T-test-run-frontend"].Dependencies

	if len(backendDeps) == 0 {
		t.Error("T-test-run-backend should have dependencies (on gen-scripts)")
	}

	if len(frontendDeps) != 1 || frontendDeps[0] != "T-test-run-backend" {
		t.Errorf("T-test-run-frontend should depend only on T-test-run-backend, got %v", frontendDeps)
	}
}

// AC3: Single surface project degenerates to no-suffix T-test-run
func TestGetBreakdownTestTasks_SingleSurfaceDegeneration(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("api"), nil, defaultAuto, "")

	for _, task := range tasks {
		if strings.HasPrefix(task.ID, "T-test-run-") {
			t.Errorf("single surface should not generate suffixed run-test, got %q", task.ID)
		}
	}

	runIdx := findTaskIndex(tasks, "T-test-run")
	if runIdx < 0 {
		t.Fatal("single surface should generate T-test-run (no suffix)")
	}

	if tasks[runIdx].Key != "run-test" {
		t.Errorf("single surface run-test Key = %q, want run-test", tasks[runIdx].Key)
	}
}

// AC4: Quick mode - T-test-gen-journeys is upstream of T-test-run-*
func TestGetQuickTestTasks_GenJourneysUpstreamOfRunTests(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetQuickTestTasks(surfaces, resolved, allEnabledAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	genJourneys, ok := byID["T-test-gen-journeys"]
	if !ok {
		t.Fatal("T-test-gen-journeys not found")
	}

	if len(genJourneys.Dependencies) != 0 {
		t.Errorf("gen-journeys should have no deps (pipeline entry), got %v", genJourneys.Dependencies)
	}

	firstRun := byID["T-test-run-backend"]
	if !dependsTransitively(tasks, firstRun, "T-test-gen-journeys") {
		t.Errorf("T-test-run-backend should transitively depend on T-test-gen-journeys")
	}
}

// dependsTransitively checks if a task transitively depends on targetID.
func dependsTransitively(tasks []AutoGenTaskDef, task AutoGenTaskDef, targetID string) bool {
	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	visited := make(map[string]bool)
	var dfs func(id string) bool
	dfs = func(id string) bool {
		if id == targetID {
			return true
		}
		if visited[id] {
			return false
		}
		visited[id] = true
		t, ok := byID[id]
		if !ok {
			return false
		}
		for _, dep := range t.Dependencies {
			if dfs(dep) {
				return true
			}
		}
		return false
	}
	return dfs(task.ID)
}

// Consolidate-specs depends on last run-test in execution order
func TestGetBreakdownTestTasks_ConsolidateDependsOnChainTail(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web", "mobile-app", "tui")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetBreakdownTestTasks(surfaces, resolved, defaultAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Last in execution order: api < web < tui
	lastRunID := "T-test-run-mobile-app"

	consolidateDeps := byID["T-specs-consolidate"].Dependencies
	if len(consolidateDeps) != 1 || consolidateDeps[0] != lastRunID {
		t.Errorf("consolidate-specs should depend on %s (chain tail), got %v", lastRunID, consolidateDeps)
	}
}

// Quick mode variant: drift depends on last run-test in chain
func TestGetQuickTestTasks_DriftDependsOnChainTail(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web", "mobile-app", "tui")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetQuickTestTasks(surfaces, resolved, allEnabledAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	lastRunID := "T-test-run-mobile-app"

	driftDeps := byID["T-quick-doc-drift"].Dependencies
	if len(driftDeps) != 1 || driftDeps[0] != lastRunID {
		t.Errorf("drift should depend on %s (chain tail), got %v", lastRunID, driftDeps)
	}
}

// Verify key != type scenario correctly populates SurfaceKey and SurfaceType
func TestGetBreakdownTestTasks_KeyNotType_SetsSurfaceFields(t *testing.T) {
	surfaces := multiSurface("admin-panel", "web", "payment-service", "api")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetBreakdownTestTasks(surfaces, resolved, defaultAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	psTask, ok := byID["T-test-run-payment-service"]
	if !ok {
		t.Fatal("T-test-run-payment-service not found")
	}
	if psTask.SurfaceKey != "payment-service" {
		t.Errorf("SurfaceKey = %q, want payment-service", psTask.SurfaceKey)
	}
	if psTask.SurfaceType != "api" {
		t.Errorf("SurfaceType = %q, want api", psTask.SurfaceType)
	}

	apTask, ok := byID["T-test-run-admin-panel"]
	if !ok {
		t.Fatal("T-test-run-admin-panel not found")
	}
	if apTask.SurfaceKey != "admin-panel" {
		t.Errorf("SurfaceKey = %q, want admin-panel", apTask.SurfaceKey)
	}
	if apTask.SurfaceType != "web" {
		t.Errorf("SurfaceType = %q, want web", apTask.SurfaceType)
	}
}

// Verify explicit execution-order overrides default type-based ordering
func TestGetBreakdownTestTasks_ExplicitExecutionOrder(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web")
	execOrder := []string{"frontend", "backend"}

	tasks := GetBreakdownTestTasks(surfaces, execOrder, defaultAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	frontendDeps := byID["T-test-run-frontend"].Dependencies
	foundGenScriptDep := false
	for _, dep := range frontendDeps {
		if strings.HasPrefix(dep, "T-test-gen-scripts-") {
			foundGenScriptDep = true
			break
		}
	}
	if !foundGenScriptDep {
		t.Errorf("T-test-run-frontend (first in chain) should depend on gen-scripts, got %v", frontendDeps)
	}

	backendDeps := byID["T-test-run-backend"].Dependencies
	if len(backendDeps) != 1 || backendDeps[0] != "T-test-run-frontend" {
		t.Errorf("T-test-run-backend should depend on T-test-run-frontend (serial), got %v", backendDeps)
	}

	consolidateDeps := byID["T-specs-consolidate"].Dependencies
	if len(consolidateDeps) != 1 || consolidateDeps[0] != "T-test-run-backend" {
		t.Errorf("consolidate-specs should depend on T-test-run-backend (chain tail), got %v", consolidateDeps)
	}
}

// --- Task 5: Dependency resolution chain tests ---

// AC1: Breakdown mode: gen-journeys -> gen-contracts -> gen-scripts-* -> run-test-{keys}
func TestGetBreakdownTestTasks_AC1_FullDAG(t *testing.T) {
	surfaces := multiSurface("auth-service", "api", "admin", "web", "cli", "cli")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetBreakdownTestTasks(surfaces, resolved, defaultAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// gen-journeys -> gen-contracts (via eval-journey)
	if byID["T-test-gen-contracts"].Dependencies[0] != "T-eval-journey" {
		t.Errorf("gen-contracts should depend on eval-journey, got %v", byID["T-test-gen-contracts"].Dependencies)
	}

	// gen-scripts depend on eval-contract
	for _, typ := range []string{"api", "web", "cli"} {
		gsID := "T-test-gen-scripts-" + typ
		if byID[gsID].Dependencies[0] != "T-eval-contract" {
			t.Errorf("%s should depend on eval-contract, got %v", gsID, byID[gsID].Dependencies)
		}
	}

	// run-test chain: first depends on all gen-scripts
	firstRunID := "T-test-run-auth-service"
	firstRunDeps := byID[firstRunID].Dependencies
	if len(firstRunDeps) != 3 {
		t.Fatalf("%s should depend on 3 gen-scripts, got %v", firstRunID, firstRunDeps)
	}

	// consolidate-specs depends on last run-test (chain tail)
	consolidateDeps := byID["T-specs-consolidate"].Dependencies
	lastRunID := "T-test-run-cli"
	if len(consolidateDeps) != 1 || consolidateDeps[0] != lastRunID {
		t.Errorf("consolidate-specs should depend on %s (chain tail), got %v", lastRunID, consolidateDeps)
	}
}

// AC2: Quick mode: gen-journeys -> run-test-{keys} (no gen-contracts/gen-scripts)
func TestGetQuickTestTasks_AC2_DirectDependencyChain(t *testing.T) {
	surfaces := multiSurface("auth-service", "api", "admin", "web", "cli", "cli")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetQuickTestTasks(surfaces, resolved, allEnabledAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// No gen-contracts or gen-scripts should exist
	if _, ok := byID["T-test-gen-contracts"]; ok {
		t.Error("Quick mode should not contain T-test-gen-contracts")
	}
	for _, typ := range []string{"api", "web", "cli"} {
		gsID := "T-test-gen-scripts-" + typ
		if _, ok := byID[gsID]; ok {
			t.Errorf("Quick mode should not contain %q", gsID)
		}
	}

	// First run-test depends directly on gen-journeys
	firstRunID := "T-test-run-auth-service"
	firstRunDeps := byID[firstRunID].Dependencies
	if len(firstRunDeps) != 1 || firstRunDeps[0] != "T-test-gen-journeys" {
		t.Errorf("%s should depend on T-test-gen-journeys, got %v", firstRunID, firstRunDeps)
	}

	// drift depends on last run-test (chain tail)
	driftDeps := byID["T-quick-doc-drift"].Dependencies
	lastRunID := "T-test-run-cli"
	if len(driftDeps) != 1 || driftDeps[0] != lastRunID {
		t.Errorf("drift should depend on %s (chain tail), got %v", lastRunID, driftDeps)
	}
}

// AC3: downstream tasks depend ONLY on the last run-test in execution order
func TestDownstreamDependsOnlyOnLastRunTest(t *testing.T) {
	t.Run("breakdown mode", func(t *testing.T) {
		surfaces := multiSurface("backend", "api", "frontend", "web")
		resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

		tasks := GetBreakdownTestTasks(surfaces, resolved, defaultAuto, "")
		byID := make(map[string]AutoGenTaskDef)
		for _, t := range tasks {
			byID[t.ID] = t
		}

		// Last in execution order is frontend (web comes after api)
		lastRunID := "T-test-run-frontend"
		consolidateDeps := byID["T-specs-consolidate"].Dependencies
		if len(consolidateDeps) != 1 {
			t.Fatalf("consolidate-specs should depend on exactly 1 task, got %v", consolidateDeps)
		}
		if consolidateDeps[0] != lastRunID {
			t.Errorf("consolidate-specs should depend on %s (chain tail), got %q", lastRunID, consolidateDeps[0])
		}
	})

	t.Run("quick mode", func(t *testing.T) {
		surfaces := multiSurface("backend", "api", "frontend", "web")
		resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

		tasks := GetQuickTestTasks(surfaces, resolved, allEnabledAuto, "")
		byID := make(map[string]AutoGenTaskDef)
		for _, t := range tasks {
			byID[t.ID] = t
		}

		lastRunID := "T-test-run-frontend"
		driftDeps := byID["T-quick-doc-drift"].Dependencies
		if len(driftDeps) != 1 {
			t.Fatalf("drift should depend on exactly 1 task, got %v", driftDeps)
		}
		if driftDeps[0] != lastRunID {
			t.Errorf("drift should depend on %s (chain tail), got %q", lastRunID, driftDeps[0])
		}
	})
}

// Quick mode single surface: run depends directly on gen-journeys
func TestGetQuickTestTasks_SingleSurface_RunDependsOnGenJourneys(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("api"), nil, allEnabledAuto, "")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	runDeps := byID["T-test-run"].Dependencies
	if len(runDeps) != 1 || runDeps[0] != "T-test-gen-journeys" {
		t.Errorf("T-test-run should depend on T-test-gen-journeys (direct), got %v", runDeps)
	}
}

// Quick mode multi-surface: full serial chain from gen-journeys
func TestGetQuickTestTasks_MultiSurface_FullSerialChain(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web", "mobile-app", "tui")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetQuickTestTasks(surfaces, resolved, allEnabledAuto, "")
	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Verify serial chain: gen-journeys -> backend -> frontend -> mobile-app -> drift
	if byID["T-test-run-backend"].Dependencies[0] != "T-test-gen-journeys" {
		t.Errorf("T-test-run-backend should depend on gen-journeys, got %v", byID["T-test-run-backend"].Dependencies)
	}
	if byID["T-test-run-frontend"].Dependencies[0] != "T-test-run-backend" {
		t.Errorf("T-test-run-frontend should depend on T-test-run-backend, got %v", byID["T-test-run-frontend"].Dependencies)
	}
	if byID["T-test-run-mobile-app"].Dependencies[0] != "T-test-run-frontend" {
		t.Errorf("T-test-run-mobile-app should depend on T-test-run-frontend, got %v", byID["T-test-run-mobile-app"].Dependencies)
	}
	if byID["T-quick-doc-drift"].Dependencies[0] != "T-test-run-mobile-app" {
		t.Errorf("drift should depend on T-test-run-mobile-app (chain tail), got %v", byID["T-quick-doc-drift"].Dependencies)
	}
}

// --- Intent-driven pipeline branching tests (Task 5) ---

// Helper: auto config with validation + clean-code + consolidate enabled for refactor tests
var refactorAuto = forgeconfig.AutoConfig{
	Test:             forgeconfig.ModeToggle{Quick: false, Full: true},
	ConsolidateSpecs: forgeconfig.ModeToggle{Quick: true, Full: true},
	CleanCode:        forgeconfig.ModeToggle{Quick: true, Full: true},
	Validation:       forgeconfig.ModeToggle{Quick: true, Full: true},
}

// AC1: GetBreakdownTestTasks skips test tasks for refactor/cleanup intent but keeps validation/consolidate/clean-code

func TestGetBreakdownTestTasks_Refactor_SkipsTestTasks(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")

	// Should NOT contain test pipeline tasks
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys" {
			t.Error("refactor should not generate T-test-gen-journeys")
		}
		if task.ID == "T-eval-journey" {
			t.Error("refactor should not generate T-eval-journey")
		}
		if task.ID == "T-test-gen-contracts" {
			t.Error("refactor should not generate T-test-gen-contracts")
		}
		if task.ID == "T-eval-contract" {
			t.Error("refactor should not generate T-eval-contract")
		}
		if strings.HasPrefix(task.ID, "T-test-gen-scripts-") {
			t.Errorf("refactor should not generate %q", task.ID)
		}
		if strings.HasPrefix(task.ID, "T-test-run") {
			t.Errorf("refactor should not generate %q", task.ID)
		}
	}

	// Should still contain validation, consolidate-specs, clean-code
	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}
	if _, ok := byID["T-validate-code"]; !ok {
		t.Error("refactor should still generate T-validate-code")
	}
	if _, ok := byID["T-specs-consolidate"]; !ok {
		t.Error("refactor should still generate T-specs-consolidate")
	}
	if _, ok := byID["T-clean-code"]; !ok {
		t.Error("refactor should still generate T-clean-code")
	}
}

func TestGetBreakdownTestTasks_Cleanup_SkipsTestTasks(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, refactorAuto, "cleanup")

	// Should NOT contain test pipeline tasks
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys" {
			t.Error("cleanup should not generate T-test-gen-journeys")
		}
		if strings.HasPrefix(task.ID, "T-test-run") {
			t.Errorf("cleanup should not generate %q", task.ID)
		}
	}

	// Should still contain validation, consolidate-specs, clean-code
	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}
	if _, ok := byID["T-validate-code"]; !ok {
		t.Error("cleanup should still generate T-validate-code")
	}
	if _, ok := byID["T-specs-consolidate"]; !ok {
		t.Error("cleanup should still generate T-specs-consolidate")
	}
	if _, ok := byID["T-clean-code"]; !ok {
		t.Error("cleanup should still generate T-clean-code")
	}
}

// AC4: new-feature intent produces identical results to current behavior (backward compat)

func TestGetBreakdownTestTasks_NewFeature_IdenticalToCurrent(t *testing.T) {
	// Call with explicit "new-feature" intent -- should produce identical results to empty intent
	tasksNewFeature := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "new-feature")
	tasksCurrent := GetBreakdownTestTasks(scalarSurface("cli"), nil, defaultAuto, "")

	if len(tasksNewFeature) != len(tasksCurrent) {
		t.Fatalf("new-feature intent: got %d tasks, current behavior: %d tasks", len(tasksNewFeature), len(tasksCurrent))
	}

	for i := range tasksNewFeature {
		if tasksNewFeature[i].ID != tasksCurrent[i].ID {
			t.Errorf("tasks[%d]: new-feature=%q, current=%q", i, tasksNewFeature[i].ID, tasksCurrent[i].ID)
		}
		if len(tasksNewFeature[i].Dependencies) != len(tasksCurrent[i].Dependencies) {
			t.Errorf("tasks[%d] deps length: new-feature=%d, current=%d", i, len(tasksNewFeature[i].Dependencies), len(tasksCurrent[i].Dependencies))
		}
		for j := range tasksNewFeature[i].Dependencies {
			if tasksNewFeature[i].Dependencies[j] != tasksCurrent[i].Dependencies[j] {
				t.Errorf("tasks[%d].Dependencies[%d]: new-feature=%q, current=%q", i, j, tasksNewFeature[i].Dependencies[j], tasksCurrent[i].Dependencies[j])
			}
		}
	}
}

func TestGetQuickTestTasks_NewFeature_IdenticalToCurrent(t *testing.T) {
	tasksNewFeature := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "new-feature")
	tasksCurrent := GetQuickTestTasks(scalarSurface("cli"), nil, allEnabledAuto, "")

	if len(tasksNewFeature) != len(tasksCurrent) {
		t.Fatalf("new-feature intent: got %d tasks, current behavior: %d tasks", len(tasksNewFeature), len(tasksCurrent))
	}

	for i := range tasksNewFeature {
		if tasksNewFeature[i].ID != tasksCurrent[i].ID {
			t.Errorf("tasks[%d]: new-feature=%q, current=%q", i, tasksNewFeature[i].ID, tasksCurrent[i].ID)
		}
	}
}

// AC2: resolveBreakdownDeps for refactor wires to last business task (no run-tests lookup)

func TestResolveBreakdownDeps_Refactor_NoRunTestDep(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// For refactor: validate-code should NOT depend on run-tests (no lastRunID lookup)
	validateTask, ok := byID["T-validate-code"]
	if !ok {
		t.Fatal("T-validate-code should exist for refactor")
	}

	// validate-code should NOT have run-test as dependency (run-tests don't exist)
	for _, dep := range validateTask.Dependencies {
		if strings.HasPrefix(dep, "T-test-run") {
			t.Errorf("refactor validate-code should not depend on run-test, got dep %q", dep)
		}
	}
}

// AC2: resolveQuickDeps for refactor/cleanup wires to last business task

func TestResolveQuickDeps_Refactor_SkipsTestTasks(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")

	// Should NOT contain test pipeline tasks
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys" {
			t.Error("refactor quick should not generate T-test-gen-journeys")
		}
		if strings.HasPrefix(task.ID, "T-test-run") {
			t.Errorf("refactor quick should not generate %q", task.ID)
		}
	}

	// Should still contain clean-code and doc-drift
	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}
	if _, ok := byID["T-clean-code"]; !ok {
		t.Error("refactor quick should still generate T-clean-code")
	}
	if _, ok := byID["T-quick-doc-drift"]; !ok {
		t.Error("refactor quick should still generate T-quick-doc-drift")
	}
}

func TestResolveQuickDeps_Cleanup_SkipsTestTasks(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, refactorAuto, "cleanup")

	// Same as refactor quick: no test tasks, but clean-code and doc-drift present
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys" {
			t.Error("cleanup quick should not generate T-test-gen-journeys")
		}
		if strings.HasPrefix(task.ID, "T-test-run") {
			t.Errorf("cleanup quick should not generate %q", task.ID)
		}
	}

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}
	if _, ok := byID["T-clean-code"]; !ok {
		t.Error("cleanup quick should still generate T-clean-code")
	}
	if _, ok := byID["T-quick-doc-drift"]; !ok {
		t.Error("cleanup quick should still generate T-quick-doc-drift")
	}
}

// AC5: Full wiring verification for all 5 scenarios

func TestIntent_RefactorBreakdown_FullWiring(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Expected tasks: validate-code, consolidate-specs, clean-code
	if len(tasks) != 3 {
		t.Fatalf("refactor breakdown should generate 3 tasks, got %d: %v", len(tasks), taskIDs(tasks))
	}

	// Wiring: no task should depend on run-test
	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if strings.HasPrefix(dep, "T-test-run") {
				t.Errorf("%s should not depend on run-test, got %v", task.ID, task.Dependencies)
			}
		}
	}
}

func TestIntent_RefactorQuick_FullWiring(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Expected tasks: validate-code, doc-drift, clean-code (refactorAuto has Validation.Quick=true)
	if len(tasks) != 3 {
		t.Fatalf("refactor quick should generate 3 tasks, got %d: %v", len(tasks), taskIDs(tasks))
	}

	// No task should depend on run-test
	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if strings.HasPrefix(dep, "T-test-run") {
				t.Errorf("%s should not depend on run-test, got %v", task.ID, task.Dependencies)
			}
		}
	}
}

func TestIntent_CleanupQuick_FullWiring(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, refactorAuto, "cleanup")

	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}

	// Expected tasks: validate-code, doc-drift, clean-code (refactorAuto has Validation.Quick=true)
	if len(tasks) != 3 {
		t.Fatalf("cleanup quick should generate 3 tasks, got %d: %v", len(tasks), taskIDs(tasks))
	}

	for _, task := range tasks {
		for _, dep := range task.Dependencies {
			if strings.HasPrefix(dep, "T-test-run") {
				t.Errorf("%s should not depend on run-test, got %v", task.ID, task.Dependencies)
			}
		}
	}
}

// AC3: Zero business task protection -- ResolveFirstTestDep with refactor/cleanup and empty business tasks

func TestResolveFirstTestDep_Refactor_NoBusinessTasks_NoPanic(_ *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")
	existing := map[string]Task{}

	// Should not panic with empty business tasks
	ResolveFirstTestDep(tasks, existing, "breakdown", "refactor")
}

func TestResolveFirstTestDep_Cleanup_NoBusinessTasks_NoPanic(_ *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, refactorAuto, "cleanup")
	existing := map[string]Task{}

	// Should not panic with empty business tasks
	ResolveFirstTestDep(tasks, existing, "quick", "cleanup")
}

// ResolveFirstTestDep with refactor/cleanup wires validate-code/clean-code to last business task

func TestResolveFirstTestDep_Refactor_WiresValidateCodeToLastBiz(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")
	existing := map[string]Task{
		"1-foo": {ID: "1"},
		"2-bar": {ID: "2"},
	}

	ResolveFirstTestDep(tasks, existing, "breakdown", "refactor")

	// For refactor breakdown: validate-code should depend on last business task
	validateIdx := findTaskIndex(tasks, "T-validate-code")
	if validateIdx < 0 {
		t.Fatal("T-validate-code not found")
	}

	if len(tasks[validateIdx].Dependencies) != 1 || tasks[validateIdx].Dependencies[0] != "2" {
		t.Errorf("validate-code should depend on last business task '2', got %v", tasks[validateIdx].Dependencies)
	}
}

// ResolveFirstTestDep for refactor/cleanup with business tasks wires the first validation/clean-code task

func TestResolveFirstTestDep_Refactor_WiresValidateCode(t *testing.T) {
	tasks := GetBreakdownTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")
	existing := map[string]Task{
		"1-foo": {ID: "1"},
		"2-bar": {ID: "2"},
	}

	ResolveFirstTestDep(tasks, existing, "breakdown", "refactor")

	// For refactor breakdown: first downstream task (validate-code) should depend on last business task
	validateIdx := findTaskIndex(tasks, "T-validate-code")
	if validateIdx < 0 {
		t.Fatal("T-validate-code not found")
	}

	if len(tasks[validateIdx].Dependencies) == 0 {
		t.Error("validate-code should have dependencies for refactor with business tasks")
	} else if tasks[validateIdx].Dependencies[0] != "2" {
		t.Errorf("validate-code should depend on last business task '2', got %v", tasks[validateIdx].Dependencies)
	}
}

func TestResolveFirstTestDep_RefactorQuick_WiresCleanCode(t *testing.T) {
	tasks := GetQuickTestTasks(scalarSurface("cli"), nil, refactorAuto, "refactor")
	existing := map[string]Task{
		"1-foo": {ID: "1"},
		"3-bar": {ID: "3"},
	}

	ResolveFirstTestDep(tasks, existing, "quick", "refactor")

	// For refactor quick: first downstream task (clean-code) should depend on last business task
	cleanIdx := findTaskIndex(tasks, "T-clean-code")
	if cleanIdx < 0 {
		t.Fatal("T-clean-code not found")
	}

	if len(tasks[cleanIdx].Dependencies) == 0 {
		t.Error("clean-code should have dependencies for refactor with business tasks")
	} else if tasks[cleanIdx].Dependencies[0] != "3" {
		t.Errorf("clean-code should depend on last business task '3', got %v", tasks[cleanIdx].Dependencies)
	}
}

// Multi-surface refactor: ensure no run-test tasks generated

func TestGetBreakdownTestTasks_Refactor_MultiSurface_NoTestTasks(t *testing.T) {
	surfaces := multiSurface("backend", "api", "frontend", "web")
	resolved, _ := forgeconfig.ResolveExecutionOrder(surfaces, nil)

	tasks := GetBreakdownTestTasks(surfaces, resolved, refactorAuto, "refactor")

	for _, task := range tasks {
		if strings.HasPrefix(task.ID, "T-test-") {
			t.Errorf("refactor multi-surface should not generate test tasks, got %q", task.ID)
		}
	}

	// Should still have validate-code, consolidate-specs, clean-code
	byID := make(map[string]AutoGenTaskDef)
	for _, t := range tasks {
		byID[t.ID] = t
	}
	if _, ok := byID["T-validate-code"]; !ok {
		t.Error("refactor multi-surface should still generate T-validate-code")
	}
}

// Verify GenerateTestTasks passes intent through

func TestGenerateTestTasks_PassesIntent_Breakdown(t *testing.T) {
	// GenerateTestTasks with refactor should produce no test pipeline tasks
	tasks := GenerateTestTasks("breakdown", scalarSurface("cli"), nil, refactorAuto, "refactor")

	for _, task := range tasks {
		if strings.HasPrefix(task.ID, "T-test-") {
			t.Errorf("GenerateTestTasks breakdown refactor should not generate test tasks, got %q", task.ID)
		}
	}
}

func TestGenerateTestTasks_PassesIntent_Quick(t *testing.T) {
	// GenerateTestTasks with cleanup should produce no test pipeline tasks
	tasks := GenerateTestTasks("quick", scalarSurface("cli"), nil, refactorAuto, "cleanup")

	for _, task := range tasks {
		if strings.HasPrefix(task.ID, "T-test-") {
			t.Errorf("GenerateTestTasks quick cleanup should not generate test tasks, got %q", task.ID)
		}
	}
}

// Helper to collect task IDs for error messages
func taskIDs(tasks []AutoGenTaskDef) []string {
	ids := make([]string, len(tasks))
	for i, t := range tasks {
		ids[i] = t.ID
	}
	return ids
}
