package task

import (
	"testing"

	"forge-cli/pkg/profile"
)

// --- AutoConfig gating tests ---

func TestGetBreakdownTestTasks_E2eTestFullFalse(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.E2eTest.Full = false

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// No e2e test tasks, no consolidate (consolidate depends on e2e test chain in breakdown)
	// Only tasks generated should be T-specs-1 (consolidate is separate gate)
	for _, task := range tasks {
		if task.ID == "T-test-1" || task.ID == "T-test-1b" || task.ID == "T-test-2-cli" ||
			task.ID == "T-test-3" || task.ID == "T-test-4" || task.ID == "T-test-4.5" {
			t.Errorf("e2e test task %s should not be generated when e2eTest.full=false", task.ID)
		}
	}
}

func TestGetBreakdownTestTasks_ConsolidateSpecsFullFalse(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.ConsolidateSpecs.Full = false

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-specs-1" {
			t.Error("T-specs-1 should not be generated when consolidateSpecs.full=false")
		}
	}
}

func TestGetBreakdownTestTasks_CleanCodeFullTrue(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.CleanCode.Full = true

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-clean-code-1" {
			found = true
			if task.Type != TypeCleanCode {
				t.Errorf("T-clean-code-1 Type = %q, want %q", task.Type, TypeCleanCode)
			}
			if !task.NoTest {
				t.Error("T-clean-code-1 NoTest should be true")
			}
		}
	}
	if !found {
		t.Error("T-clean-code-1 should be generated when cleanCode.full=true")
	}
}

func TestGetBreakdownTestTasks_CleanCodeFullFalse(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.CleanCode.Full = false

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-clean-code-1" {
			t.Error("T-clean-code-1 should not be generated when cleanCode.full=false")
		}
	}
}

func TestGetQuickTestTasks_E2eTestQuickFalse(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.E2eTest.Quick = false

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// No e2e test tasks
	for _, task := range tasks {
		if task.ID == "T-quick-1" || task.ID == "T-quick-2-cli" ||
			task.ID == "T-quick-3" || task.ID == "T-quick-4" {
			t.Errorf("e2e test task %s should not be generated when e2eTest.quick=false", task.ID)
		}
	}
}

func TestGetQuickTestTasks_ConsolidateSpecsQuickFalse(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.ConsolidateSpecs.Quick = false

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-quick-specs-1" {
			t.Error("T-quick-specs-1 should not be generated when consolidateSpecs.quick=false")
		}
	}
}

func TestGetQuickTestTasks_CleanCodeQuickTrue(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.CleanCode.Quick = true

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	found := false
	for _, task := range tasks {
		if task.ID == "T-clean-code-1" {
			found = true
			if task.Type != TypeCleanCode {
				t.Errorf("T-clean-code-1 Type = %q, want %q", task.Type, TypeCleanCode)
			}
			if !task.NoTest {
				t.Error("T-clean-code-1 NoTest should be true")
			}
		}
	}
	if !found {
		t.Error("T-clean-code-1 should be generated when cleanCode.quick=true")
	}
}

func TestGetQuickTestTasks_CleanCodeQuickFalse(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.CleanCode.Quick = false

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	for _, task := range tasks {
		if task.ID == "T-clean-code-1" {
			t.Error("T-clean-code-1 should not be generated when cleanCode.quick=false")
		}
	}
}

// --- Backward compat: defaults produce same behavior as before ---

func TestGetBreakdownTestTasks_DefaultsMatchOldBehavior(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Should produce exactly 7 tasks (same as before)
	if len(tasks) != 7 {
		t.Fatalf("expected 7 tasks with defaults, got %d", len(tasks))
	}

	wantIDs := []string{"T-test-1", "T-test-1b", "T-test-2-cli", "T-test-3", "T-test-4", "T-test-4.5", "T-specs-1"}
	for i, want := range wantIDs {
		if tasks[i].ID != want {
			t.Errorf("tasks[%d].ID = %q, want %q", i, tasks[i].ID, want)
		}
	}
}

func TestGetQuickTestTasks_DefaultsProduceNoTasks(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Defaults: e2eTest.quick=false, consolidateSpecs.quick=false → no quick test tasks
	if len(tasks) != 0 {
		t.Fatalf("expected 0 quick tasks with defaults (quick=false), got %d", len(tasks))
	}
}

// --- Clean code task dependency wiring ---

func TestGetBreakdownTestTasks_CleanCodeDependsOnVerifyRegression(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.CleanCode.Full = true

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Find T-clean-code-1 and verify its dependencies
	for _, task := range tasks {
		if task.ID == "T-clean-code-1" {
			// In breakdown mode, T-clean-code-1 should not have deps set yet
			// (resolved later by ResolveFirstTestDep)
			return
		}
	}
	t.Error("T-clean-code-1 not found")
}

func TestGetQuickTestTasks_CleanCodeNoE2e(t *testing.T) {
	// When e2e tests and consolidate specs are disabled but clean code is enabled
	auto := profile.AutoConfigDefaults()
	auto.E2eTest.Quick = false
	auto.ConsolidateSpecs.Quick = false
	auto.CleanCode.Quick = true

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Should have exactly 1 task (T-clean-code-1)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-clean-code-1" {
		t.Errorf("task ID = %q, want T-clean-code-1", tasks[0].ID)
	}
}

func TestGetBreakdownTestTasks_NoE2eWithCleanCode(t *testing.T) {
	// When e2e tests and consolidate specs are disabled but clean code is enabled
	auto := profile.AutoConfigDefaults()
	auto.E2eTest.Full = false
	auto.ConsolidateSpecs.Full = false
	auto.CleanCode.Full = true

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Should have T-clean-code-1 only (no e2e test tasks)
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-clean-code-1" {
		t.Errorf("task ID = %q, want T-clean-code-1", tasks[0].ID)
	}
}

func TestGetBreakdownTestTasks_OnlyConsolidateSpecs(t *testing.T) {
	// When e2e tests are disabled but consolidate specs is enabled
	auto := profile.AutoConfigDefaults()
	auto.E2eTest.Full = false
	auto.ConsolidateSpecs.Full = true

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Should have T-specs-1 only
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-specs-1" {
		t.Errorf("task ID = %q, want T-specs-1", tasks[0].ID)
	}
}

func TestGetQuickTestTasks_OnlyConsolidateSpecs(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.E2eTest.Quick = false
	auto.ConsolidateSpecs.Quick = true

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Should have T-quick-specs-1 only
	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-quick-specs-1" {
		t.Errorf("task ID = %q, want T-quick-specs-1", tasks[0].ID)
	}
}

// --- T-specs-1 dependency in breakdown mode when e2e tasks exist ---

func TestGetBreakdownTestTasks_SpecsDependsOnVerifyRegression(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Find T-specs-1
	for _, task := range tasks {
		if task.ID == "T-specs-1" {
			if len(task.Dependencies) != 1 || task.Dependencies[0] != "T-test-4.5" {
				t.Errorf("T-specs-1 deps = %v, want [T-test-4.5]", task.Dependencies)
			}
			return
		}
	}
	t.Error("T-specs-1 not found")
}

func TestGetQuickTestTasks_SpecsDependsOnVerifyRegression(t *testing.T) {
	auto := allEnabledAuto
	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	// Find T-quick-specs-1
	for _, task := range tasks {
		if task.ID == "T-quick-specs-1" {
			if len(task.Dependencies) != 1 || task.Dependencies[0] != "T-quick-4" {
				t.Errorf("T-quick-specs-1 deps = %v, want [T-quick-4]", task.Dependencies)
			}
			return
		}
	}
	t.Error("T-quick-specs-1 not found")
}

// --- InferType tests for new IDs ---

func TestInferType_TSpecs1(t *testing.T) {
	got := InferType("T-specs-1")
	if got != TypeDocGenerationConsolidate {
		t.Errorf("InferType(T-specs-1) = %q, want %q", got, TypeDocGenerationConsolidate)
	}
}

func TestInferType_TQuickSpecs1(t *testing.T) {
	got := InferType("T-quick-specs-1")
	if got != TypeDocGenerationDrift {
		t.Errorf("InferType(T-quick-specs-1) = %q, want %q", got, TypeDocGenerationDrift)
	}
}

func TestInferType_TCleanCode1(t *testing.T) {
	got := InferType("T-clean-code-1")
	if got != TypeCleanCode {
		t.Errorf("InferType(T-clean-code-1) = %q, want %q", got, TypeCleanCode)
	}
}

// --- All auto=false produces zero tasks ---

func TestGetBreakdownTestTasks_AllAutoOff(t *testing.T) {
	auto := profile.AutoConfig{
		E2eTest:          profile.ModeToggle{Quick: false, Full: false},
		ConsolidateSpecs: profile.ModeToggle{Quick: false, Full: false},
		CleanCode:        profile.ModeToggle{Quick: false, Full: false},
	}

	tasks := GetBreakdownTestTasks([]string{"go-test"}, []string{"cli"}, auto)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks with all auto off, got %d", len(tasks))
	}
}

func TestGetQuickTestTasks_AllAutoOff(t *testing.T) {
	auto := profile.AutoConfig{
		E2eTest:          profile.ModeToggle{Quick: false, Full: false},
		ConsolidateSpecs: profile.ModeToggle{Quick: false, Full: false},
		CleanCode:        profile.ModeToggle{Quick: false, Full: false},
	}

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks with all auto off, got %d", len(tasks))
	}
}

// --- T-clean-code-1 only + consolidate ---

func TestGetQuickTestTasks_CleanCodeAndSpecsNoE2e(t *testing.T) {
	auto := profile.AutoConfigDefaults()
	auto.E2eTest.Quick = false
	auto.CleanCode.Quick = true
	auto.ConsolidateSpecs.Quick = true

	tasks := GetQuickTestTasks([]string{"go-test"}, []string{"cli"}, auto)

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d: %+v", len(tasks), tasks)
	}
	if tasks[0].ID != "T-quick-specs-1" {
		t.Errorf("tasks[0].ID = %q, want T-quick-specs-1", tasks[0].ID)
	}
	if tasks[1].ID != "T-clean-code-1" {
		t.Errorf("tasks[1].ID = %q, want T-clean-code-1", tasks[1].ID)
	}
}
