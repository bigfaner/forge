package task

import (
	"testing"

	"forge-cli/pkg/forgeconfig"
)

// --- AutoConfig gating tests ---

func TestGetBreakdownTestTasks_TestFullFalse(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Full = false

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	// No e2e test tasks, no consolidate (consolidate depends on e2e test chain in breakdown)
	for _, task := range tasks {
		if task.ID == "T-test-gen-scripts-cli" ||
			task.ID == "T-test-run" || task.ID == "T-test-verify-regression" {
			t.Errorf("e2e test task %s should not be generated when test.full=false", task.ID)
		}
	}
}

func TestGetBreakdownTestTasks_ConsolidateSpecsFullFalse(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.ConsolidateSpecs.Full = false

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-specs-consolidate" {
			t.Error("T-specs-consolidate should not be generated when consolidateSpecs.full=false")
		}
	}
}

func TestGetBreakdownTestTasks_CleanCodeFullTrue(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.CleanCode.Full = true

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-clean-code" {
			found = true
			if task.Type != TypeCleanCode {
				t.Errorf("T-clean-code Type = %q, want %q", task.Type, TypeCleanCode)
			}
		}
	}
	if !found {
		t.Error("T-clean-code should be generated when cleanCode.full=true")
	}
}

func TestGetBreakdownTestTasks_CleanCodeFullFalse(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.CleanCode.Full = false

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-clean-code" {
			t.Error("T-clean-code should not be generated when cleanCode.full=false")
		}
	}
}

func TestGetQuickTestTasks_TestQuickFalse(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Quick = false

	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	// No e2e test tasks
	for _, task := range tasks {
		if task.ID == "T-test-gen-journeys-cli" ||
			task.ID == "T-test-gen-contracts" ||
			task.ID == "T-test-gen-scripts-cli" ||
			task.ID == "T-test-run" ||
			task.ID == "T-test-verify-regression" {
			t.Errorf("e2e test task %s should not be generated when test.quick=false", task.ID)
		}
	}
}

func TestGetQuickTestTasks_ConsolidateSpecsQuickFalse(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.ConsolidateSpecs.Quick = false

	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-quick-doc-drift" {
			t.Error("T-quick-doc-drift should not be generated when consolidateSpecs.quick=false")
		}
	}
}

func TestGetQuickTestTasks_CleanCodeQuickTrue(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.CleanCode.Quick = true

	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-clean-code" {
			found = true
			if task.Type != TypeCleanCode {
				t.Errorf("T-clean-code Type = %q, want %q", task.Type, TypeCleanCode)
			}
		}
	}
	if !found {
		t.Error("T-clean-code should be generated when cleanCode.quick=true")
	}
}

func TestGetQuickTestTasks_CleanCodeQuickFalse(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.CleanCode.Quick = false

	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-clean-code" {
			t.Error("T-clean-code should not be generated when cleanCode.quick=false")
		}
	}
}

// --- Backward compat: defaults produce same behavior as before ---

func TestGetBreakdownTestTasks_DefaultsMatchOldBehavior(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	// Should produce 8 tasks (gen-journeys + eval-journey + gen-contracts + eval-contract + gen-scripts + run + verify-regression + consolidate)
	if len(tasks) != 8 {
		t.Fatalf("expected 8 tasks with defaults, got %d", len(tasks))
	}

	wantIDs := []string{"T-test-gen-journeys", "T-eval-journey", "T-test-gen-contracts", "T-eval-contract", "T-test-gen-scripts-cli", "T-test-run", "T-test-verify-regression", "T-specs-consolidate"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}
}

func TestGetQuickTestTasks_DefaultsProduceNoE2ETasks(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	// Defaults: test.quick=false, consolidateSpecs.quick=true -> only spec drift task
	if len(tasks) != 1 {
		t.Fatalf("expected 1 quick task (spec drift) with defaults, got %d", len(tasks))
	}
	if tasks[0].ID != "T-quick-doc-drift" {
		t.Errorf("expected T-quick-doc-drift, got %q", tasks[0].ID)
	}
}

// --- Clean code task dependency wiring ---

func TestGetBreakdownTestTasks_CleanCodeDependsOnVerifyRegression(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.CleanCode.Full = true

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	// Find T-clean-code and verify its dependencies
	for _, task := range tasks {
		if task.ID == "T-clean-code" {
			// In breakdown mode, T-clean-code should not have deps set yet
			// (resolved later by ResolveFirstTestDep)
			return
		}
	}
	t.Error("T-clean-code not found")
}

func TestGetQuickTestTasks_CleanCodeNoE2e(t *testing.T) {
	// When e2e tests and consolidate specs are disabled but clean code is enabled
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Quick = false
	auto.ConsolidateSpecs.Quick = false
	auto.CleanCode.Quick = true

	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	// Should have exactly 1 task (T-clean-code)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-clean-code" {
		t.Errorf("task ID = %q, want T-clean-code", tasks[0].ID)
	}
}

func TestGetBreakdownTestTasks_NoE2eWithCleanCode(t *testing.T) {
	// When e2e tests and consolidate specs are disabled but clean code is enabled
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Full = false
	auto.ConsolidateSpecs.Full = false
	auto.CleanCode.Full = true

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	// Should have T-clean-code only (no e2e test tasks)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-clean-code" {
		t.Errorf("task ID = %q, want T-clean-code", tasks[0].ID)
	}
}

func TestGetBreakdownTestTasks_OnlyConsolidateSpecs(t *testing.T) {
	// When e2e tests are disabled but consolidate specs is enabled
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Full = false
	auto.ConsolidateSpecs.Full = true

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	// Should have T-specs-consolidate only
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-specs-consolidate" {
		t.Errorf("task ID = %q, want T-specs-consolidate", tasks[0].ID)
	}
}

func TestGetQuickTestTasks_OnlyConsolidateSpecs(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Quick = false
	auto.ConsolidateSpecs.Quick = true

	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	// Should have T-quick-doc-drift only
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-quick-doc-drift" {
		t.Errorf("task ID = %q, want T-quick-doc-drift", tasks[0].ID)
	}
}

// --- T-specs-consolidate dependency in breakdown mode when e2e tasks exist ---

func TestGetBreakdownTestTasks_SpecsDependsOnVerifyRegression(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)

	// Find T-specs-consolidate
	for _, task := range tasks {
		if task.ID == "T-specs-consolidate" {
			if len(task.Dependencies) != 1 || task.Dependencies[0] != "T-test-verify-regression" {
				t.Errorf("T-specs-consolidate deps = %v, want [T-test-verify-regression]", task.Dependencies)
			}
			return
		}
	}
	t.Error("T-specs-consolidate not found")
}

func TestGetQuickTestTasks_SpecsDependsOnVerifyRegression(t *testing.T) {
	auto := allEnabledAuto
	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	// Find T-quick-doc-drift
	for _, task := range tasks {
		if task.ID == "T-quick-doc-drift" {
			if len(task.Dependencies) != 1 || task.Dependencies[0] != "T-test-verify-regression" {
				t.Errorf("T-quick-doc-drift deps = %v, want [T-test-verify-regression]", task.Dependencies)
			}
			return
		}
	}
	t.Error("T-quick-doc-drift not found")
}

// --- InferType tests for new IDs ---

func TestInferType_TSpecs1(t *testing.T) {
	got := InferType("T-specs-consolidate", nil)
	if got != TypeDocConsolidate {
		t.Errorf("InferType(T-specs-consolidate) = %q, want %q", got, TypeDocConsolidate)
	}
}

func TestInferType_TQuickSpecs1(t *testing.T) {
	got := InferType("T-quick-doc-drift", nil)
	if got != TypeDocDrift {
		t.Errorf("InferType(T-quick-doc-drift) = %q, want %q", got, TypeDocDrift)
	}
}

func TestInferType_TCleanCode1(t *testing.T) {
	got := InferType("T-clean-code", nil)
	if got != TypeCleanCode {
		t.Errorf("InferType(T-clean-code) = %q, want %q", got, TypeCleanCode)
	}
}

// --- All auto=false produces zero tasks ---

func TestGetBreakdownTestTasks_AllAutoOff(t *testing.T) {
	auto := forgeconfig.AutoConfig{
		Test:             forgeconfig.ModeToggle{Quick: false, Full: false},
		ConsolidateSpecs: forgeconfig.ModeToggle{Quick: false, Full: false},
		CleanCode:        forgeconfig.ModeToggle{Quick: false, Full: false},
	}

	tasks := GetBreakdownTestTasks([]string{"cli"}, auto)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks with all auto off, got %d", len(tasks))
	}
}

func TestGetQuickTestTasks_AllAutoOff(t *testing.T) {
	auto := forgeconfig.AutoConfig{
		Test:             forgeconfig.ModeToggle{Quick: false, Full: false},
		ConsolidateSpecs: forgeconfig.ModeToggle{Quick: false, Full: false},
		CleanCode:        forgeconfig.ModeToggle{Quick: false, Full: false},
	}

	tasks := GetQuickTestTasks([]string{"cli"}, auto)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks with all auto off, got %d", len(tasks))
	}
}

// --- T-clean-code only + consolidate ---

func TestGetQuickTestTasks_CleanCodeAndSpecsNoE2e(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Quick = false
	auto.CleanCode.Quick = true
	auto.ConsolidateSpecs.Quick = true

	tasks := GetQuickTestTasks([]string{"cli"}, auto)

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-quick-doc-drift" {
		t.Errorf("tasks[0].ID = %q, want T-quick-doc-drift", tasks[0].ID)
	}
	if tasks[1].ID != "T-clean-code" {
		t.Errorf("tasks[1].ID = %q, want T-clean-code", tasks[1].ID)
	}
}
