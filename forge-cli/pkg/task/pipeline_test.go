package task

import (
	"strings"
	"testing"

	"forge-cli/pkg/forgeconfig"
)

// ---------------------------------------------------------------------------
// Phase 1: Registry validation tests
// ---------------------------------------------------------------------------

func TestValidatePipelineRegistry_Valid(t *testing.T) {
	// ValidatePipelineRegistry runs at init() and panics on failure.
	// If we reach this test, the registry is structurally valid.
	if err := ValidatePipelineRegistry(); err != nil {
		t.Errorf("ValidatePipelineRegistry returned error: %v", err)
	}
}

func TestValidatePipelineRegistry_AllNodesHaveGenerateCondition(t *testing.T) {
	for i, node := range PipelineRegistry {
		if node.GenerateCondition == nil {
			t.Errorf("node %q (index %d): GenerateCondition is nil", node.ID, i)
		}
	}
}

func TestValidatePipelineRegistry_ExpectedNodeCount(t *testing.T) {
	// PipelineRegistry should have exactly 12 nodes (as of initial implementation).
	if len(PipelineRegistry) != 12 {
		t.Errorf("PipelineRegistry has %d nodes, want 12", len(PipelineRegistry))
	}
}

func TestValidatePipelineRegistry_AllTypesDistinct(t *testing.T) {
	seen := make(map[string]string) // type -> ID
	for _, node := range PipelineRegistry {
		if prev, dup := seen[node.Type]; dup {
			t.Errorf("duplicate type %q in nodes %q and %q", node.Type, prev, node.ID)
		}
		seen[node.Type] = node.ID
	}
}

func TestValidatePipelineRegistry_NoBlankIDs(t *testing.T) {
	for i, node := range PipelineRegistry {
		if node.ID == "" {
			t.Errorf("node at index %d has empty ID", i)
		}
		if node.Type == "" {
			t.Errorf("node %q (index %d) has empty Type", node.ID, i)
		}
		if node.Title == "" {
			t.Errorf("node %q (index %d) has empty Title", node.ID, i)
		}
	}
}

func TestValidatePipelineRegistry_PlaceholderExpansionConsistency(t *testing.T) {
	for i, node := range PipelineRegistry {
		hasSurfaceKey := strings.Contains(node.ID, "{surface-key}") || strings.Contains(node.Key, "{surface-key}")
		hasSurfaceType := strings.Contains(node.ID, "{surface-type}") || strings.Contains(node.Key, "{surface-type}")

		switch node.Expansion {
		case "per-surface-key":
			if !hasSurfaceKey {
				t.Errorf("node %q (index %d): per-surface-key expansion but no {surface-key} in ID/Key", node.ID, i)
			}
		case "per-surface-type":
			if !hasSurfaceType {
				t.Errorf("node %q (index %d): per-surface-type expansion but no {surface-type} in ID/Key", node.ID, i)
			}
		default:
			if hasSurfaceKey || hasSurfaceType {
				t.Errorf("node %q (index %d): no expansion but has placeholders in ID/Key", node.ID, i)
			}
		}
	}
}

func TestValidatePipelineRegistry_DepRefIntegrity(t *testing.T) {
	// All static DepRef.Ref strings must reference existing node IDs.
	knownIDs := make(map[string]bool)
	for _, node := range PipelineRegistry {
		knownIDs[node.ID] = true
	}
	for _, node := range PipelineRegistry {
		for _, dep := range node.DependsOn {
			if dep.Resolve == nil && dep.Ref != "" {
				if !knownIDs[dep.Ref] && !idExistsInRegistry(dep.Ref) {
					t.Errorf("node %q: DependsOn.Ref %q does not match any registry node ID", node.ID, dep.Ref)
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Config gate function tests
// ---------------------------------------------------------------------------

func TestGateTest(t *testing.T) {
	tests := []struct {
		name string
		mode string
		auto forgeconfig.AutoConfig
		want bool
	}{
		{"quick enabled", "quick", forgeconfig.AutoConfig{Test: forgeconfig.ModeToggle{Quick: true}}, true},
		{"quick disabled", "quick", forgeconfig.AutoConfig{Test: forgeconfig.ModeToggle{Quick: false}}, false},
		{"breakdown enabled", "breakdown", forgeconfig.AutoConfig{Test: forgeconfig.ModeToggle{Full: true}}, true},
		{"breakdown disabled", "breakdown", forgeconfig.AutoConfig{Test: forgeconfig.ModeToggle{Full: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GateTest(tt.mode, tt.auto); got != tt.want {
				t.Errorf("GateTest(%q, %+v) = %v, want %v", tt.mode, tt.auto, got, tt.want)
			}
		})
	}
}

func TestGateValidation(t *testing.T) {
	tests := []struct {
		name string
		mode string
		auto forgeconfig.AutoConfig
		want bool
	}{
		{"quick enabled", "quick", forgeconfig.AutoConfig{Validation: forgeconfig.ModeToggle{Quick: true}}, true},
		{"quick disabled", "quick", forgeconfig.AutoConfig{Validation: forgeconfig.ModeToggle{Quick: false}}, false},
		{"breakdown enabled", "breakdown", forgeconfig.AutoConfig{Validation: forgeconfig.ModeToggle{Full: true}}, true},
		{"breakdown disabled", "breakdown", forgeconfig.AutoConfig{Validation: forgeconfig.ModeToggle{Full: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GateValidation(tt.mode, tt.auto); got != tt.want {
				t.Errorf("GateValidation(%q, %+v) = %v, want %v", tt.mode, tt.auto, got, tt.want)
			}
		})
	}
}

func TestGateConsolidateSpecs(t *testing.T) {
	tests := []struct {
		name string
		mode string
		auto forgeconfig.AutoConfig
		want bool
	}{
		{"quick enabled", "quick", forgeconfig.AutoConfig{ConsolidateSpecs: forgeconfig.ModeToggle{Quick: true}}, true},
		{"quick disabled", "quick", forgeconfig.AutoConfig{ConsolidateSpecs: forgeconfig.ModeToggle{Quick: false}}, false},
		{"breakdown enabled", "breakdown", forgeconfig.AutoConfig{ConsolidateSpecs: forgeconfig.ModeToggle{Full: true}}, true},
		{"breakdown disabled", "breakdown", forgeconfig.AutoConfig{ConsolidateSpecs: forgeconfig.ModeToggle{Full: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GateConsolidateSpecs(tt.mode, tt.auto); got != tt.want {
				t.Errorf("GateConsolidateSpecs(%q, %+v) = %v, want %v", tt.mode, tt.auto, got, tt.want)
			}
		})
	}
}

func TestGateCleanCode(t *testing.T) {
	tests := []struct {
		name string
		mode string
		auto forgeconfig.AutoConfig
		want bool
	}{
		{"quick enabled", "quick", forgeconfig.AutoConfig{CleanCode: forgeconfig.ModeToggle{Quick: true}}, true},
		{"quick disabled", "quick", forgeconfig.AutoConfig{CleanCode: forgeconfig.ModeToggle{Quick: false}}, false},
		{"breakdown enabled", "breakdown", forgeconfig.AutoConfig{CleanCode: forgeconfig.ModeToggle{Full: true}}, true},
		{"breakdown disabled", "breakdown", forgeconfig.AutoConfig{CleanCode: forgeconfig.ModeToggle{Full: false}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GateCleanCode(tt.mode, tt.auto); got != tt.want {
				t.Errorf("GateCleanCode(%q, %+v) = %v, want %v", tt.mode, tt.auto, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Intent gate function tests
// ---------------------------------------------------------------------------

func TestGateAllowAll(t *testing.T) {
	intents := []string{"new-feature", "refactor", "cleanup", ""}
	for _, intent := range intents {
		if !GateAllowAll(intent) {
			t.Errorf("GateAllowAll(%q) = false, want true for all intents", intent)
		}
	}
}

func TestGateBlockSkipTest(t *testing.T) {
	tests := []struct {
		intent string
		want   bool
	}{
		{"new-feature", true},
		{"", true},
		{"refactor", false},
		{"cleanup", false},
	}
	for _, tt := range tests {
		t.Run(tt.intent, func(t *testing.T) {
			if got := GateBlockSkipTest(tt.intent); got != tt.want {
				t.Errorf("GateBlockSkipTest(%q) = %v, want %v", tt.intent, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Generate condition function tests
// ---------------------------------------------------------------------------

func TestCondHasTestableTasks(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  bool
	}{
		{"nil tasks returns true (legacy compat)", nil, true},
		{"empty tasks returns false", []Task{}, false},
		{"feature task returns true", []Task{{ID: "1", Type: TypeCodingFeature}}, true},
		{"doc task returns false", []Task{{ID: "1", Type: TypeDoc}}, false},
		{"mixed tasks returns true", []Task{{ID: "1", Type: TypeDoc}, {ID: "2", Type: TypeCodingFeature}}, true},
		{"clean-code type returns true", []Task{{ID: "1", Type: TypeCleanCode}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CondHasTestableTasks(tt.tasks); got != tt.want {
				t.Errorf("CondHasTestableTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCondHasDocTasks(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
		want  bool
	}{
		{"nil tasks returns false", nil, false},
		{"empty tasks returns false", []Task{}, false},
		{"doc task returns true", []Task{{ID: "1", Type: TypeDoc}}, true},
		{"doc.fix returns true", []Task{{ID: "1", Type: TypeDocFix}}, true},
		{"feature task returns false", []Task{{ID: "1", Type: TypeCodingFeature}}, false},
		{"mixed returns true", []Task{{ID: "1", Type: TypeCodingFeature}, {ID: "2", Type: TypeDoc}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CondHasDocTasks(tt.tasks); got != tt.want {
				t.Errorf("CondHasDocTasks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCondAlways(t *testing.T) {
	if !CondAlways(nil) {
		t.Error("CondAlways(nil) = false, want true")
	}
	if !CondAlways([]Task{}) {
		t.Error("CondAlways([]) = false, want true")
	}
	if !CondAlways([]Task{{ID: "1", Type: TypeCodingFeature}}) {
		t.Error("CondAlways(tasks) = false, want true")
	}
}

// ---------------------------------------------------------------------------
// Dependency resolver function tests
// ---------------------------------------------------------------------------

func TestPipeline_ResolveLastRunTest(t *testing.T) {
	tests := []struct {
		name string
		ctx  *GenContext
		want []string
	}{
		{"empty chain", &GenContext{RunTestChain: nil}, nil},
		{"single task", &GenContext{RunTestChain: []string{"T-test-run-api"}}, []string{"T-test-run-api"}},
		{"multiple tasks returns last", &GenContext{RunTestChain: []string{"T-test-run-api", "T-test-run-web"}}, []string{"T-test-run-web"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveLastRunTest(tt.ctx)
			if len(got) != len(tt.want) {
				t.Fatalf("ResolveLastRunTest() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ResolveLastRunTest()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestPipeline_ResolveUpstream(t *testing.T) {
	tests := []struct {
		name string
		ctx  *GenContext
		want []string
	}{
		{"empty upstream", &GenContext{UpstreamIDs: nil}, nil},
		{"single upstream", &GenContext{UpstreamIDs: []string{"T-test-gen-journeys"}}, []string{"T-test-gen-journeys"}},
		{"multiple upstream", &GenContext{UpstreamIDs: []string{"T-test-gen-scripts-api", "T-test-gen-scripts-cli"}}, []string{"T-test-gen-scripts-api", "T-test-gen-scripts-cli"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveUpstream(tt.ctx)
			if len(got) != len(tt.want) {
				t.Fatalf("ResolveUpstream() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ResolveUpstream()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestPipeline_ResolveDocTasks(t *testing.T) {
	tests := []struct {
		name string
		ctx  *GenContext
		want []string
	}{
		{"no business tasks", &GenContext{BusinessTasks: nil}, nil},
		{"no doc tasks", &GenContext{BusinessTasks: []Task{{ID: "1", Type: TypeCodingFeature}}}, nil},
		{"one doc task", &GenContext{BusinessTasks: []Task{{ID: "1", Type: TypeDoc}}}, []string{"1"}},
		{"mixed tasks returns only doc IDs", &GenContext{BusinessTasks: []Task{
			{ID: "1", Type: TypeCodingFeature},
			{ID: "2", Type: TypeDoc},
			{ID: "3", Type: TypeDocFix},
		}}, []string{"2", "3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveDocTasks(tt.ctx)
			if len(got) != len(tt.want) {
				t.Fatalf("ResolveDocTasks() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ResolveDocTasks()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestPipeline_ResolveLastBusinessTask(t *testing.T) {
	tests := []struct {
		name string
		ctx  *GenContext
		want []string
	}{
		{"no business tasks", &GenContext{BusinessTasks: nil}, nil},
		{"single task", &GenContext{BusinessTasks: []Task{{ID: "1", Type: TypeCodingFeature}}}, []string{"1"}},
		{"multiple tasks returns highest numeric", &GenContext{BusinessTasks: []Task{
			{ID: "1.1", Type: TypeCodingFeature},
			{ID: "2.1", Type: TypeCodingFeature},
			{ID: "1.2", Type: TypeDoc},
		}}, []string{"2.1"}},
		{"task with simple numeric ID", &GenContext{BusinessTasks: []Task{
			{ID: "1", Type: TypeCodingFeature},
			{ID: "3", Type: TypeDoc},
			{ID: "2", Type: TypeDocFix},
		}}, []string{"3"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveLastBusinessTask(tt.ctx)
			if len(got) != len(tt.want) {
				t.Fatalf("ResolveLastBusinessTask() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ResolveLastBusinessTask()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestPipeline_ResolveHighestGateOrLastBiz(t *testing.T) {
	tests := []struct {
		name          string
		existingTasks map[string]Task
		businessTasks []Task
		wantContains  string // substring to check in result
		wantNil       bool
	}{
		{
			name:          "empty returns nil",
			existingTasks: map[string]Task{},
			businessTasks: nil,
			wantNil:       true,
		},
		{
			name:          "gate takes priority",
			existingTasks: map[string]Task{"1.gate": {ID: "1.gate"}, "1.summary": {ID: "1.summary"}},
			businessTasks: []Task{{ID: "1.1", Type: TypeCodingFeature}},
			wantContains:  "1.gate",
		},
		{
			name:          "summary when no gate",
			existingTasks: map[string]Task{"1.summary": {ID: "1.summary"}},
			businessTasks: []Task{{ID: "1.1", Type: TypeCodingFeature}},
			wantContains:  "1.summary",
		},
		{
			name:          "higher phase gate wins",
			existingTasks: map[string]Task{"1.gate": {ID: "1.gate"}, "2.gate": {ID: "2.gate"}},
			businessTasks: []Task{{ID: "1.1", Type: TypeCodingFeature}},
			wantContains:  "2.gate",
		},
		{
			name:          "business task wins over lower gate",
			existingTasks: map[string]Task{"1.gate": {ID: "1.gate"}},
			businessTasks: []Task{{ID: "3", Type: TypeCodingFeature}},
			wantContains:  "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &GenContext{ExistingTasks: tt.existingTasks, BusinessTasks: tt.businessTasks}
			got := ResolveHighestGateOrLastBiz(ctx)
			if tt.wantNil {
				if got != nil {
					t.Errorf("ResolveHighestGateOrLastBiz() = %v, want nil", got)
				}
				return
			}
			if len(got) != 1 {
				t.Fatalf("ResolveHighestGateOrLastBiz() = %v, want exactly 1 result", got)
			}
			if !strings.Contains(got[0], tt.wantContains) {
				t.Errorf("ResolveHighestGateOrLastBiz() = %v, want result containing %q", got, tt.wantContains)
			}
		})
	}
}

func TestPipeline_ResolveLastRunTestOrBusiness(t *testing.T) {
	tests := []struct {
		name         string
		runTestChain []string
		bizTasks     []Task
		want         []string
	}{
		{"empty returns nil", nil, nil, nil},
		{"run-test chain takes priority", []string{"T-test-run-api"}, []Task{{ID: "1"}}, []string{"T-test-run-api"}},
		{"falls back to last business task", nil, []Task{{ID: "1"}, {ID: "2"}}, []string{"2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &GenContext{RunTestChain: tt.runTestChain, BusinessTasks: tt.bizTasks}
			got := ResolveLastRunTestOrBusiness(ctx)
			if len(got) != len(tt.want) {
				t.Fatalf("ResolveLastRunTestOrBusiness() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ResolveLastRunTestOrBusiness()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestPipeline_ResolveIfGenerated(t *testing.T) {
	resolver := ResolveIfGenerated("T-review-doc")

	tests := []struct {
		name         string
		allGenerated []string
		want         []string
	}{
		{"found", []string{"T-review-doc", "T-clean-code"}, []string{"T-review-doc"}},
		{"not found", []string{"T-clean-code"}, nil},
		{"empty", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &GenContext{AllGenerated: tt.allGenerated}
			got := resolver(ctx)
			if len(got) != len(tt.want) {
				t.Fatalf("ResolveIfGenerated() = %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ResolveIfGenerated()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GenerateTestTasks edge case tests
// ---------------------------------------------------------------------------

func TestGenerateTestTasks_EmptySurfaces_NonSurfaceTasksStillGenerate(t *testing.T) {
	auto := allEnabledAuto
	tasks := GenerateTestTasks("quick", map[string]string{}, nil, auto, "", nil, nil)

	// Non-surface tasks should still generate
	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}
	// T-quick-doc-drift generates because CondAlways + GateConsolidateSpecs.Quick=true
	if _, ok := byID["T-quick-doc-drift"]; !ok {
		t.Error("T-quick-doc-drift should generate even with empty surfaces")
	}
	// Surface-dependent nodes produce zero tasks
	for _, task := range tasks {
		if strings.HasPrefix(task.ID, "T-test-run") || strings.HasPrefix(task.ID, "T-test-gen-scripts") {
			t.Errorf("surface-dependent task %s should not generate with empty surfaces", task.ID)
		}
	}
}

func TestGenerateTestTasks_RefactorIntent_BlocksTestPipeline(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Full = true
	auto.ConsolidateSpecs.Full = true

	tasks := GenerateTestTasks("breakdown", scalarSurface("api"), nil, auto, "refactor",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)

	// Test pipeline tasks should be blocked by GateBlockSkipTest
	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}
	if _, ok := byID["T-test-gen-journeys"]; ok {
		t.Error("T-test-gen-journeys should NOT generate for refactor intent")
	}
	if _, ok := byID["T-test-run"]; ok {
		t.Error("T-test-run should NOT generate for refactor intent")
	}
	// Non-test nodes with GateAllowAll should still generate
	if _, ok := byID["T-specs-consolidate"]; !ok {
		t.Error("T-specs-consolidate should generate for refactor intent (GateAllowAll)")
	}
}

func TestGenerateTestTasks_CleanupIntent_NonTestNodesStillGenerate(t *testing.T) {
	auto := forgeconfig.AutoConfigDefaults()
	auto.Test.Quick = true
	auto.ConsolidateSpecs.Quick = true
	auto.CleanCode.Quick = true

	tasks := GenerateTestTasks("quick", scalarSurface("cli"), nil, auto, "cleanup", nil, nil)

	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}
	// T-clean-code has GateAllowAll, should generate for cleanup
	if _, ok := byID["T-clean-code"]; !ok {
		t.Error("T-clean-code should generate for cleanup intent")
	}
	// T-quick-doc-drift has GateAllowAll, should generate for cleanup
	if _, ok := byID["T-quick-doc-drift"]; !ok {
		t.Error("T-quick-doc-drift should generate for cleanup intent")
	}
	// Test pipeline nodes should be blocked
	if _, ok := byID["T-test-gen-journeys"]; ok {
		t.Error("T-test-gen-journeys should NOT generate for cleanup intent")
	}
}

func TestGenerateTestTasks_MultiSurface_Expansion(t *testing.T) {
	auto := allEnabledAuto
	surfaces := map[string]string{
		"backend":  "api",
		"frontend": "web",
	}

	tasks := GenerateTestTasks("quick", surfaces, []string{"backend", "frontend"}, auto, "",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)

	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}

	// Run-test should expand per surface key
	if _, ok := byID["T-test-run-backend"]; !ok {
		t.Error("T-test-run-backend should exist for multi-surface")
	}
	if _, ok := byID["T-test-run-frontend"]; !ok {
		t.Error("T-test-run-frontend should exist for multi-surface")
	}

	// Serial chain: frontend depends on backend
	frontend := byID["T-test-run-frontend"]
	foundBackendDep := false
	for _, dep := range frontend.Dependencies {
		if dep == "T-test-run-backend" {
			foundBackendDep = true
		}
	}
	if !foundBackendDep {
		t.Errorf("T-test-run-frontend should depend on T-test-run-backend, deps=%v", frontend.Dependencies)
	}
}

func TestGenerateTestTasks_SingleSurface_NoSuffix(t *testing.T) {
	auto := allEnabledAuto
	tasks := GenerateTestTasks("quick", scalarSurface("api"), nil, auto, "",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)

	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}

	// Single surface: run-test should NOT have surface-key suffix
	if _, ok := byID["T-test-run"]; !ok {
		t.Error("T-test-run should exist without surface-key suffix for single surface")
	}
	// The suffixed version should NOT exist
	if _, ok := byID["T-test-run-."]; ok {
		t.Error("T-test-run-. should NOT exist for scalar single surface")
	}
}

func TestGenerateTestTasks_UISurfaceOnly_Gating(t *testing.T) {
	auto := validationAuto
	auto.ConsolidateSpecs = forgeconfig.ModeToggle{Quick: true}

	// CLI surface — no visual UI, T-validate-ux should NOT generate
	tasks := GenerateTestTasks("quick", scalarSurface("cli"), nil, auto, "",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)

	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}
	if _, ok := byID["T-validate-ux"]; ok {
		t.Error("T-validate-ux should NOT generate for CLI surface (no visual UI)")
	}
	if _, ok := byID["T-validate-code"]; !ok {
		t.Error("T-validate-code should generate for CLI surface")
	}
}

func TestGenerateTestTasks_UISurfaceOnly_WithWeb(t *testing.T) {
	auto := validationAuto
	auto.ConsolidateSpecs = forgeconfig.ModeToggle{Quick: true}

	// Web surface — has visual UI, T-validate-ux should generate
	tasks := GenerateTestTasks("quick", scalarSurface("web"), nil, auto, "",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)

	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}
	if _, ok := byID["T-validate-ux"]; !ok {
		t.Error("T-validate-ux should generate for web surface (visual UI)")
	}
	if _, ok := byID["T-validate-code"]; !ok {
		t.Error("T-validate-code should generate for web surface")
	}
}

func TestGenerateTestTasks_BreakdownMode_SkipsQuickOnlyNodes(t *testing.T) {
	auto := allEnabledAuto
	tasks := GenerateTestTasks("breakdown", scalarSurface("api"), nil, auto, "",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)

	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}
	// T-quick-doc-drift has Mode: "quick", should NOT generate in breakdown
	if _, ok := byID["T-quick-doc-drift"]; ok {
		t.Error("T-quick-doc-drift should NOT generate in breakdown mode")
	}
	// T-specs-consolidate has Mode: "breakdown", should generate
	if _, ok := byID["T-specs-consolidate"]; !ok {
		t.Error("T-specs-consolidate should generate in breakdown mode")
	}
}

func TestGenerateTestTasks_QuickMode_SkipsBreakdownOnlyNodes(t *testing.T) {
	auto := allEnabledAuto
	tasks := GenerateTestTasks("quick", scalarSurface("api"), nil, auto, "",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)

	byID := make(map[string]AutoGenTaskDef)
	for _, task := range tasks {
		byID[task.ID] = task
	}
	// T-eval-journey has Mode: "breakdown", should NOT generate in quick
	if _, ok := byID["T-eval-journey"]; ok {
		t.Error("T-eval-journey should NOT generate in quick mode")
	}
	// T-test-gen-contracts has Mode: "breakdown", should NOT generate in quick
	if _, ok := byID["T-test-gen-contracts"]; ok {
		t.Error("T-test-gen-contracts should NOT generate in quick mode")
	}
	// T-test-gen-journeys has no mode restriction, should generate
	if _, ok := byID["T-test-gen-journeys"]; !ok {
		t.Error("T-test-gen-journeys should generate in quick mode")
	}
}

func TestGenerateTestTasks_DocReviewOnlyForDocTasks(t *testing.T) {
	auto := allEnabledAuto

	// Feature tasks only — T-review-doc should NOT generate
	tasks := GenerateTestTasks("quick", scalarSurface("api"), nil, auto, "",
		[]Task{{ID: "1", Type: TypeCodingFeature}}, nil)
	for _, task := range tasks {
		if task.ID == "T-review-doc" {
			t.Error("T-review-doc should NOT generate when no doc tasks exist")
		}
	}

	// Doc tasks — T-review-doc should generate
	tasks = GenerateTestTasks("quick", scalarSurface("api"), nil, auto, "",
		[]Task{{ID: "1", Type: TypeDoc}}, nil)
	found := false
	for _, task := range tasks {
		if task.ID == "T-review-doc" {
			found = true
		}
	}
	if !found {
		t.Error("T-review-doc should generate when doc tasks exist")
	}
}

func TestGenerateTestTasks_DocReviewDependsOnDocTasks(t *testing.T) {
	auto := allEnabledAuto
	bizTasks := []Task{
		{ID: "1", Type: TypeDoc},
		{ID: "2", Type: TypeDoc},
	}
	tasks := GenerateTestTasks("quick", scalarSurface("api"), nil, auto, "", bizTasks, nil)

	for _, task := range tasks {
		if task.ID == "T-review-doc" {
			if len(task.Dependencies) != 2 {
				t.Errorf("T-review-doc deps count = %d, want 2", len(task.Dependencies))
			}
			found1, found2 := false, false
			for _, dep := range task.Dependencies {
				if dep == "1" {
					found1 = true
				}
				if dep == "2" {
					found2 = true
				}
			}
			if !found1 || !found2 {
				t.Errorf("T-review-doc deps = %v, want [1, 2]", task.Dependencies)
			}
			return
		}
	}
	t.Error("T-review-doc not found")
}

// ---------------------------------------------------------------------------
// Registry-driven InferType matching tests
// ---------------------------------------------------------------------------

func TestMatchRegistryID_ExactMatch(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"T-review-doc", TypeDocReview},
		{"T-clean-code", TypeCleanCode},
		{"T-test-gen-journeys", TypeTestGenJourneys},
		{"T-eval-journey", TypeEvalJourney},
		{"T-test-gen-contracts", TypeTestGenContracts},
		{"T-eval-contract", TypeEvalContract},
		{"T-validate-code", TypeValidationCode},
		{"T-validate-ux", TypeValidationUx},
		{"T-specs-consolidate", TypeDocConsolidate},
		{"T-quick-doc-drift", TypeDocDrift},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := matchRegistryID(tt.id, nil)
			if got != tt.want {
				t.Errorf("matchRegistryID(%q, nil) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestMatchRegistryID_SurfaceKeySuffix_GenScripts(t *testing.T) {
	surfaces := map[string]string{"backend": "api", "frontend": "web"}

	tests := []struct {
		id       string
		surfaces map[string]string
		want     string
	}{
		{"T-test-gen-scripts-backend", surfaces, TypeTestGenScripts},
		{"T-test-gen-scripts-frontend", surfaces, TypeTestGenScripts},
		{"T-test-gen-scripts", nil, TypeTestGenScripts}, // degenerate form (single surface)
		{"T-test-gen-scripts-backend", nil, ""},         // no surfaces map -> no key match
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := matchRegistryID(tt.id, tt.surfaces)
			if got != tt.want {
				t.Errorf("matchRegistryID(%q, ...) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestMatchRegistryID_SurfaceKeySuffix(t *testing.T) {
	surfaces := map[string]string{"backend": "api", "frontend": "web"}

	tests := []struct {
		id   string
		want string
	}{
		{"T-test-run-backend", TypeTestRun},
		{"T-test-run-frontend", TypeTestRun},
		{"T-test-run", TypeTestRun}, // degenerate single-surface form
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := matchRegistryID(tt.id, surfaces)
			if got != tt.want {
				t.Errorf("matchRegistryID(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestMatchRegistryID_UnknownID(t *testing.T) {
	got := matchRegistryID("T-unknown-task", nil)
	if got != "" {
		t.Errorf("matchRegistryID(T-unknown-task) = %q, want empty", got)
	}
}

// ---------------------------------------------------------------------------
// Helper function tests
// ---------------------------------------------------------------------------

func TestDeriveKey(t *testing.T) {
	tests := []struct {
		key, id, want string
	}{
		{"review-doc", "T-review-doc", "review-doc"},
		{"", "T-review-doc", "review-doc"},          // derive from ID
		{"", "no-prefix", "no-prefix"},              // no T- prefix
		{"custom-key", "T-something", "custom-key"}, // explicit key wins
	}
	for _, tt := range tests {
		t.Run(tt.key+"/"+tt.id, func(t *testing.T) {
			if got := deriveKey(tt.key, tt.id); got != tt.want {
				t.Errorf("deriveKey(%q, %q) = %q, want %q", tt.key, tt.id, got, tt.want)
			}
		})
	}
}

func TestHasVisualUI(t *testing.T) {
	tests := []struct {
		name     string
		surfaces map[string]string
		want     bool
	}{
		{"nil surfaces", nil, false},
		{"empty surfaces", map[string]string{}, false},
		{"api only", map[string]string{".": "api"}, false},
		{"cli only", map[string]string{".": "cli"}, false},
		{"web", map[string]string{".": "web"}, true},
		{"tui", map[string]string{".": "tui"}, true},
		{"mobile", map[string]string{".": "mobile"}, true},
		{"mixed with web", map[string]string{"api": "api", "web": "web"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasVisualUI(tt.surfaces); got != tt.want {
				t.Errorf("hasVisualUI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExpandTitle(t *testing.T) {
	got := expandTitle("Generate {test-type-title} Scripts", "api")
	if !strings.Contains(got, "API") {
		t.Errorf("expandTitle() = %q, should contain API title", got)
	}
}

func TestSortedSurfaceKeys(t *testing.T) {
	surfaces := map[string]string{
		"charlie": "api",
		"alpha":   "cli",
		"bravo":   "web",
	}
	got := sortedSurfaceKeys(surfaces)
	if len(got) != 3 || got[0] != "alpha" || got[1] != "bravo" || got[2] != "charlie" {
		t.Errorf("sortedSurfaceKeys() = %v, want [alpha bravo charlie]", got)
	}
}

func TestMatchTypeSuffixedID(t *testing.T) {
	// matchTypeSuffixedID is a generic helper for per-surface-type pattern matching.
	// Note: gen-test-scripts now uses per-surface-key, so the {surface-type} template
	// below is hypothetical — no current registry node uses per-surface-type.
	tests := []struct {
		id       string
		template string
		want     bool
	}{
		{"T-test-gen-scripts-api", "T-test-gen-scripts-{surface-type}", true},
		{"T-test-gen-scripts", "T-test-gen-scripts-{surface-type}", true}, // degenerate
		{"T-test-gen-scripts-", "T-test-gen-scripts-{surface-type}", false},
		{"T-other-api", "T-test-gen-scripts-{surface-type}", false},
		{"T-test-gen-scripts-api123", "T-test-gen-scripts-{surface-type}", false}, // digits invalid
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if got := matchTypeSuffixedID(tt.id, tt.template); got != tt.want {
				t.Errorf("matchTypeSuffixedID(%q, %q) = %v, want %v", tt.id, tt.template, got, tt.want)
			}
		})
	}
}

func TestMatchSurfaceKeyID(t *testing.T) {
	surfaces := map[string]string{"backend": "api", "frontend": "web"}

	tests := []struct {
		id       string
		template string
		want     bool
	}{
		{"T-test-run-backend", "T-test-run-{surface-key}", true},
		{"T-test-run-frontend", "T-test-run-{surface-key}", true},
		{"T-test-run", "T-test-run-{surface-key}", true}, // degenerate
		{"T-test-run-unknown", "T-test-run-{surface-key}", false},
		{"T-other-backend", "T-test-run-{surface-key}", false},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			if got := matchSurfaceKeyID(tt.id, tt.template, surfaces); got != tt.want {
				t.Errorf("matchSurfaceKeyID(%q, %q) = %v, want %v", tt.id, tt.template, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Phase 2 dynamic validation tests
// ---------------------------------------------------------------------------

func TestValidateGeneratedTasks_ValidTasks(t *testing.T) {
	tasks := []AutoGenTaskDef{
		{ID: "T-gen-journeys", Dependencies: nil},
		{ID: "T-run-test", Dependencies: []string{"T-gen-journeys"}},
	}
	if err := validateGeneratedTasks(tasks); err != nil {
		t.Errorf("validateGeneratedTasks() error: %v", err)
	}
}

func TestValidateGeneratedTasks_DanglingDependency(t *testing.T) {
	tasks := []AutoGenTaskDef{
		{ID: "T-run-test", Dependencies: []string{"T-nonexistent"}},
	}
	if err := validateGeneratedTasks(tasks); err == nil {
		t.Error("expected error for dangling T- dependency")
	}
}

func TestValidateGeneratedTasks_BusinessTaskDepOK(t *testing.T) {
	// Dependencies on non-T- IDs (business tasks) should not trigger error
	tasks := []AutoGenTaskDef{
		{ID: "T-review-doc", Dependencies: []string{"1", "2"}},
	}
	if err := validateGeneratedTasks(tasks); err != nil {
		t.Errorf("validateGeneratedTasks() should not error for business task deps: %v", err)
	}
}

func TestCheckNoCycles_Valid(t *testing.T) {
	tasks := []AutoGenTaskDef{
		{ID: "A", Dependencies: nil},
		{ID: "B", Dependencies: []string{"A"}},
		{ID: "C", Dependencies: []string{"B"}},
	}
	if err := checkNoCycles(tasks); err != nil {
		t.Errorf("checkNoCycles() error: %v", err)
	}
}

func TestCheckNoCycles_Cycle(t *testing.T) {
	tasks := []AutoGenTaskDef{
		{ID: "A", Dependencies: []string{"B"}},
		{ID: "B", Dependencies: []string{"A"}},
	}
	if err := checkNoCycles(tasks); err == nil {
		t.Error("expected cycle error")
	}
}

func TestCheckNoCycles_Empty(t *testing.T) {
	if err := checkNoCycles(nil); err != nil {
		t.Errorf("checkNoCycles(nil) error: %v", err)
	}
}
