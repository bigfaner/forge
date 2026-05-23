package task

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// --- Phase Detection Tests ---

func TestDetectPhases_BasicPhases(t *testing.T) {
	taskIDs := []string{"1.1", "1.2", "2.1", "2.2", "3.1"}
	phases := DetectPhases(taskIDs)

	if len(phases) != 3 {
		t.Fatalf("expected 3 phases, got %d: %+v", len(phases), phases)
	}
	if phases[0].Number != 1 || len(phases[0].TaskIDs) != 2 {
		t.Errorf("phase 1: want 2 tasks, got %d", len(phases[0].TaskIDs))
	}
	if phases[1].Number != 2 || len(phases[1].TaskIDs) != 2 {
		t.Errorf("phase 2: want 2 tasks, got %d", len(phases[1].TaskIDs))
	}
	if phases[2].Number != 3 || len(phases[2].TaskIDs) != 1 {
		t.Errorf("phase 3: want 1 task, got %d", len(phases[2].TaskIDs))
	}
}

func TestDetectPhases_ExcludesTestTaskIDs(t *testing.T) {
	taskIDs := []string{"1.1", "1.2", "T-test-gen-scripts-cli", "T-quick-gen-and-run-cli", "2.1"}
	phases := DetectPhases(taskIDs)

	if len(phases) != 2 {
		t.Fatalf("expected 2 phases, got %d", len(phases))
	}
	// Phase 1 should only have 1.1 and 1.2 (T-test/T-quick excluded)
	p1 := phases[0]
	if p1.Number != 1 {
		t.Errorf("first phase number = %d, want 1", p1.Number)
	}
	for _, id := range p1.TaskIDs {
		if strings.HasPrefix(id, "T-test") || strings.HasPrefix(id, "T-quick") {
			t.Errorf("test task ID %q should be excluded from phase tasks", id)
		}
	}
}

func TestDetectPhases_MalformedIDs(t *testing.T) {
	taskIDs := []string{"intro", "1.2a", "abc.def", "1.1", "2.1", "2.2", "1..2"}
	phases := DetectPhases(taskIDs)

	// Only "1.1" (phase 1) and "2.1","2.2" (phase 2) should be detected
	if len(phases) != 2 {
		t.Fatalf("expected 2 phases, got %d: %+v", len(phases), phases)
	}
	if phases[0].Number != 1 || len(phases[0].TaskIDs) != 1 {
		t.Errorf("phase 1: want 1 task, got %d", len(phases[0].TaskIDs))
	}
	if phases[1].Number != 2 || len(phases[1].TaskIDs) != 2 {
		t.Errorf("phase 2: want 2 tasks, got %d", len(phases[1].TaskIDs))
	}
}

func TestDetectPhases_Empty(t *testing.T) {
	phases := DetectPhases(nil)
	if len(phases) != 0 {
		t.Errorf("expected 0 phases for empty input, got %d", len(phases))
	}
}

func TestDetectPhases_SkipsGateSummaryIDs(t *testing.T) {
	taskIDs := []string{"1.1", "1.2", "1.summary", "1.gate", "2.1"}
	phases := DetectPhases(taskIDs)

	if len(phases) != 2 {
		t.Fatalf("expected 2 phases, got %d", len(phases))
	}
	// Phase 1 should only have 1.1 and 1.2
	p1 := phases[0]
	if len(p1.TaskIDs) != 2 {
		t.Errorf("phase 1: want 2 business tasks, got %d: %v", len(p1.TaskIDs), p1.TaskIDs)
	}
	for _, id := range p1.TaskIDs {
		if strings.HasSuffix(id, ".summary") || strings.HasSuffix(id, ".gate") {
			t.Errorf("gate/summary ID %q should be excluded", id)
		}
	}
}

// --- Threshold Tests ---

func TestPhaseInfo_Qualifies(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []string
		expected bool
	}{
		{"2 tasks qualifies", []string{"1.1", "1.2"}, true},
		{"3 tasks qualifies", []string{"1.1", "1.2", "1.3"}, true},
		{"1 task does not qualify", []string{"1.1"}, false},
		{"0 tasks does not qualify", []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PhaseInfo{Number: 1, TaskIDs: tt.tasks}
			if got := p.Qualifies(); got != tt.expected {
				t.Errorf("Qualifies() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// --- Template Generation Tests ---

func TestGenerateSummaryMD(t *testing.T) {
	phase := PhaseInfo{
		Number:  1,
		TaskIDs: []string{"1.1", "1.2"},
	}
	content, err := GenerateSummaryMD(phase, "auth-login")
	if err != nil {
		t.Fatalf("GenerateSummaryMD error: %v", err)
	}
	s := string(content)

	// Check frontmatter
	if !strings.Contains(s, `id: "1.summary"`) {
		t.Error("missing id: 1.summary")
	}
	if !strings.Contains(s, `type: "doc.summary"`) {
		t.Error("missing type: doc.summary")
	}
	if !strings.Contains(s, `"1.1"`) || !strings.Contains(s, `"1.2"`) {
		t.Error("missing dependency on business task IDs")
	}
}

func TestGenerateGateMD(t *testing.T) {
	phase := PhaseInfo{
		Number:  2,
		TaskIDs: []string{"2.1", "2.2", "2.3"},
	}
	content, err := GenerateGateMD(phase, "auth-login")
	if err != nil {
		t.Fatalf("GenerateGateMD error: %v", err)
	}
	s := string(content)

	// Check frontmatter
	if !strings.Contains(s, `id: "2.gate"`) {
		t.Error("missing id: 2.gate")
	}
	if !strings.Contains(s, `type: "gate"`) {
		t.Error("missing type: gate")
	}
	if !strings.Contains(s, `breaking: true`) {
		t.Error("missing breaking: true")
	}
	if !strings.Contains(s, `"2.summary"`) {
		t.Error("gate should depend on 2.summary")
	}
	// Gate should NOT depend directly on business tasks
	if strings.Contains(s, `"2.1"`) {
		t.Error("gate should not depend directly on business task 2.1")
	}
}

func TestGenerateSummaryMD_DependencyOrder(t *testing.T) {
	phase := PhaseInfo{
		Number:  3,
		TaskIDs: []string{"3.3", "3.1", "3.2"},
	}
	content, err := GenerateSummaryMD(phase, "test-feature")
	if err != nil {
		t.Fatalf("GenerateSummaryMD error: %v", err)
	}
	s := string(content)

	// Dependencies should be sorted
	idx31 := strings.Index(s, `"3.1"`)
	idx32 := strings.Index(s, `"3.2"`)
	idx33 := strings.Index(s, `"3.3"`)
	if idx31 >= idx32 || idx32 >= idx33 {
		t.Error("dependencies should be in sorted order")
	}
}

// --- Integration: GenerateStageGates ---

func TestGenerateStageGates_CreatesFiles(t *testing.T) {
	dir := t.TempDir()
	taskIDs := []string{"1.1", "1.2", "2.1"}

	generated, err := GenerateStageGates(taskIDs, dir, "test-feature")
	if err != nil {
		t.Fatalf("GenerateStageGates error: %v", err)
	}
	// Only phase 1 qualifies (2 tasks >= threshold)
	if generated != 2 { // 1 summary + 1 gate
		t.Errorf("generated = %d, want 2", generated)
	}

	// Check files exist
	if _, err := os.Stat(filepath.Join(dir, "1.summary.md")); os.IsNotExist(err) {
		t.Error("1.summary.md not created")
	}
	if _, err := os.Stat(filepath.Join(dir, "1.gate.md")); os.IsNotExist(err) {
		t.Error("1.gate.md not created")
	}
	// Phase 2 should NOT have files (only 1 task)
	if _, err := os.Stat(filepath.Join(dir, "2.summary.md")); err == nil {
		t.Error("2.summary.md should not be created for single-task phase")
	}
}

func TestGenerateStageGates_Idempotent(t *testing.T) {
	dir := t.TempDir()
	taskIDs := []string{"1.1", "1.2"}

	// First run
	gen1, err := GenerateStageGates(taskIDs, dir, "test-feature")
	if err != nil {
		t.Fatalf("first run error: %v", err)
	}
	if gen1 != 2 {
		t.Errorf("first run generated = %d, want 2", gen1)
	}

	// Read content of first run
	summary1, _ := os.ReadFile(filepath.Join(dir, "1.summary.md"))
	gate1, _ := os.ReadFile(filepath.Join(dir, "1.gate.md"))

	// Second run (should skip existing)
	gen2, err := GenerateStageGates(taskIDs, dir, "test-feature")
	if err != nil {
		t.Fatalf("second run error: %v", err)
	}
	if gen2 != 0 {
		t.Errorf("second run generated = %d, want 0 (idempotent)", gen2)
	}

	// Content should be unchanged
	summary2, _ := os.ReadFile(filepath.Join(dir, "1.summary.md"))
	gate2, _ := os.ReadFile(filepath.Join(dir, "1.gate.md"))
	if string(summary1) != string(summary2) {
		t.Error("summary content changed on re-run")
	}
	if string(gate1) != string(gate2) {
		t.Error("gate content changed on re-run")
	}
}

func TestGenerateStageGates_PartialState(t *testing.T) {
	dir := t.TempDir()
	taskIDs := []string{"1.1", "1.2"}

	// Create summary manually
	summaryContent := "---\nid: \"1.summary\"\ntitle: \"Custom Summary\"\n---\n\nCustom content\n"
	if err := os.WriteFile(filepath.Join(dir, "1.summary.md"), []byte(summaryContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Run generation — should only create gate
	generated, err := GenerateStageGates(taskIDs, dir, "test-feature")
	if err != nil {
		t.Fatalf("GenerateStageGates error: %v", err)
	}
	if generated != 1 { // only gate
		t.Errorf("generated = %d, want 1 (partial state)", generated)
	}

	// Summary should be unchanged
	content, _ := os.ReadFile(filepath.Join(dir, "1.summary.md"))
	if string(content) != summaryContent {
		t.Error("existing summary was overwritten")
	}

	// Gate should exist
	if _, err := os.Stat(filepath.Join(dir, "1.gate.md")); os.IsNotExist(err) {
		t.Error("1.gate.md not created")
	}
}

func TestGenerateStageGates_NoQualifyingPhases(t *testing.T) {
	dir := t.TempDir()
	// Single task per phase — no phase qualifies
	taskIDs := []string{"1.1", "2.1", "3.1"}

	generated, err := GenerateStageGates(taskIDs, dir, "test-feature")
	if err != nil {
		t.Fatalf("GenerateStageGates error: %v", err)
	}
	if generated != 0 {
		t.Errorf("generated = %d, want 0 (no qualifying phases)", generated)
	}
}

func TestGenerateStageGates_AllMalformedIDs(t *testing.T) {
	dir := t.TempDir()
	taskIDs := []string{"intro", "setup", "readme"}

	generated, err := GenerateStageGates(taskIDs, dir, "test-feature")
	if err != nil {
		t.Fatalf("GenerateStageGates error: %v", err)
	}
	if generated != 0 {
		t.Errorf("generated = %d, want 0 (all malformed)", generated)
	}
}

// --- PhaseInfo sorted output ---

func TestDetectPhases_SortedByPhaseNumber(t *testing.T) {
	taskIDs := []string{"3.1", "1.1", "2.1"}
	phases := DetectPhases(taskIDs)

	for i := 1; i < len(phases); i++ {
		if phases[i].Number <= phases[i-1].Number {
			t.Errorf("phases not sorted: %d before %d", phases[i-1].Number, phases[i].Number)
		}
	}
}

func TestDetectPhases_TaskIDsSortedWithinPhase(t *testing.T) {
	taskIDs := []string{"1.3", "1.1", "1.2"}
	phases := DetectPhases(taskIDs)

	if len(phases) != 1 {
		t.Fatalf("expected 1 phase, got %d", len(phases))
	}
	p := phases[0]
	sorted := make([]string, len(p.TaskIDs))
	copy(sorted, p.TaskIDs)
	sort.Strings(sorted)
	for i := range p.TaskIDs {
		if p.TaskIDs[i] != sorted[i] {
			t.Errorf("TaskIDs not sorted: got %v, want %v", p.TaskIDs, sorted)
			break
		}
	}
}
