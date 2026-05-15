package prompt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/task"
)

// buildTestIndex writes a minimal index.json to dir and returns its path.
func buildTestIndex(t *testing.T, dir string, tasks map[string]task.Task) string {
	t.Helper()
	index := task.NewTestIndex("test-feature", tasks)
	data, err := json.Marshal(index)
	if err != nil {
		t.Fatalf("marshal index: %v", err)
	}
	path := filepath.Join(dir, "index.json")
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	return path
}

// setupFeatureDir creates the minimal directory structure for a feature.
func setupFeatureDir(t *testing.T, projectRoot string, tasks map[string]task.Task) {
	t.Helper()
	tasksDir := filepath.Join(projectRoot, "docs", "features", "test-feature", "tasks")
	recordsDir := filepath.Join(tasksDir, "records")
	if err := os.MkdirAll(recordsDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	buildTestIndex(t, tasksDir, tasks)
}

// --- Synthesize tests ---

func TestSynthesize_AllTypes(t *testing.T) {
	types := []string{
		task.TypeImplementation,
		task.TypeDocumentation,
		task.TypeDocEvaluation,
		task.TypeDocGenerationSummary,
		task.TypeDocGenerationConsolidate,
		task.TypeDocGenerationDrift,
		task.TypeTestPipelineGenCases,
		task.TypeTestPipelineEvalCases,
		task.TypeTestPipelineGenScripts,
		task.TypeTestPipelineRun,
		task.TypeTestPipelineGraduate,
		task.TypeTestPipelineVerifyRegression,
		task.TypeFix,
		task.TypeGate,
	}

	for _, typ := range types {
		t.Run(typ, func(t *testing.T) {
			dir := t.TempDir()
			tasks := map[string]task.Task{
				"1.1-impl": {
					ID:     "1.1",
					Title:  "Test task",
					Status: "pending",
					File:   "1.1-impl.md",
					Record: "records/1.1-impl.md",
					Type:   typ,
					Scope:  "backend",
				},
			}
			setupFeatureDir(t, dir, tasks)

			opts := SynthesizeOpts{
				ProjectRoot: dir,
				FeatureSlug: "test-feature",
				TaskID:      "1.1",
			}
			result, err := Synthesize(opts)
			if err != nil {
				t.Fatalf("Synthesize(%s): unexpected error: %v", typ, err)
			}
			if result == "" {
				t.Errorf("Synthesize(%s): got empty string", typ)
			}
			if strings.Contains(result, "{{") {
				t.Errorf("Synthesize(%s): result contains unreplaced placeholder: %s", typ, result)
			}
		})
	}
}

func TestSynthesize_FixRecordMissed(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{
		ProjectRoot:     dir,
		FeatureSlug:     "test-feature",
		TaskID:          "1.1",
		FixRecordMissed: true,
	}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("got empty string")
	}
	// fix-record-missed template should mention recovery
	if !strings.Contains(result, "record") {
		t.Error("fix-record-missed template should mention 'record'")
	}
	if strings.Contains(result, "{{") {
		t.Errorf("result contains unreplaced placeholder: %s", result)
	}
}

func TestSynthesize_EmptyType_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   "", // empty
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{
		ProjectRoot: dir,
		FeatureSlug: "test-feature",
		TaskID:      "1.1",
	}
	_, err := Synthesize(opts)
	if err == nil {
		t.Fatal("expected error for empty type, got nil")
	}
	if !strings.Contains(err.Error(), "type field missing") {
		t.Errorf("expected 'type field missing' in error, got: %v", err)
	}
}

func TestSynthesize_UnknownType_ReturnsError(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   "unknown-type-xyz",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{
		ProjectRoot: dir,
		FeatureSlug: "test-feature",
		TaskID:      "1.1",
	}
	_, err := Synthesize(opts)
	if err == nil {
		t.Fatal("expected error for unknown type, got nil")
	}
	if !strings.Contains(err.Error(), "unknown type") {
		t.Errorf("expected 'unknown type' in error, got: %v", err)
	}
}

// --- Empty variable rendering tests ---

func TestSynthesize_EmptyPhaseSummary_NoResidual(t *testing.T) {
	// Phase 1 tasks have no PHASE_SUMMARY; verify no residual text.
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{
		ProjectRoot: dir,
		FeatureSlug: "test-feature",
		TaskID:      "1.1",
	}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The sentence "If `` is non-empty" must NOT appear.
	if strings.Contains(result, "If `` is non-empty") {
		t.Errorf("result contains residual empty-backtick sentence:\n%s", result)
	}

	// No "PHASE_SUMMARY:" label should remain.
	if strings.Contains(result, "PHASE_SUMMARY:") {
		t.Errorf("result contains residual PHASE_SUMMARY label:\n%s", result)
	}

	// Check that consecutive blank lines don't appear (which would indicate
	// a removed placeholder left an extra blank line).
	if strings.Contains(result, "\n\n\n") {
		t.Errorf("result contains triple newlines (likely from removed placeholder):\n%s", result)
	}
}

func TestSynthesize_NonEmptyPhaseSummary_Preserved(t *testing.T) {
	// Phase 2 task with phase 1 completed and a summary file present.
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Phase 1 task",
			Status: "completed",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
			Scope:  "backend",
		},
		"2.1-impl": {
			ID:     "2.1",
			Title:  "Phase 2 task",
			Status: "pending",
			File:   "2.1-impl.md",
			Record: "records/2.1-impl.md",
			Type:   task.TypeImplementation,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	// Create the phase 1 summary.
	summaryPath := filepath.Join(dir, "docs", "features", "test-feature", "tasks", "records", "1-summary.md")
	if err := os.WriteFile(summaryPath, []byte("# Phase 1 Summary"), 0644); err != nil {
		t.Fatalf("write summary: %v", err)
	}

	opts := SynthesizeOpts{
		ProjectRoot: dir,
		FeatureSlug: "test-feature",
		TaskID:      "2.1",
	}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// The conditional sentence should be present with the actual summary path.
	if !strings.Contains(result, "1-summary.md") {
		t.Errorf("result should reference the summary file:\n%s", result)
	}
	// The conditional sentence with backticks should contain the path (not empty backticks).
	if strings.Contains(result, "If `` is non-empty") {
		t.Errorf("result contains residual empty-backtick sentence:\n%s", result)
	}
}

func TestSynthesize_EmptyScope_NoTrailingSpace(t *testing.T) {
	// When scope is "" or "all", the SCOPE variable is empty.
	// Check that "SCOPE: " (with trailing space) is cleaned up.
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
			Scope:  "", // empty scope
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{
		ProjectRoot: dir,
		FeatureSlug: "test-feature",
		TaskID:      "1.1",
	}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "just compile " (trailing space) should not appear.
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		trimmed := strings.TrimRight(line, " \t")
		if line != trimmed && strings.HasPrefix(trimmed, "just ") {
			t.Errorf("line %d has trailing space after 'just' command: %q", i+1, line)
		}
	}
}

func TestSynthesize_EmptyProfile_NoResidual(t *testing.T) {
	// Test pipeline templates with no profile set.
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Test task",
			Status: "pending",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeTestPipelineGenScripts,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{
		ProjectRoot: dir,
		FeatureSlug: "test-feature",
		TaskID:      "1.1",
	}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "PROFILE: " (with trailing space and no value) should not appear.
	if strings.Contains(result, "PROFILE: \n") || strings.Contains(result, "PROFILE: \r\n") {
		t.Errorf("result contains empty PROFILE label:\n%s", result)
	}
}

// --- PhaseDetect tests ---

func TestPhaseDetect_NewPhase(t *testing.T) {
	// currentPhase > maxCompletedPhase AND currentPhase > 1 → inject summary path
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Phase 1 task",
			Status: "completed",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
		},
	}
	setupFeatureDir(t, dir, tasks)

	// Create the summary file that PhaseDetect should find.
	summaryPath := filepath.Join(dir, "docs", "features", "test-feature", "tasks", "records", "1-summary.md")
	if err := os.WriteFile(summaryPath, []byte("# Phase 1 Summary"), 0644); err != nil {
		t.Fatalf("write summary: %v", err)
	}

	// Task 2.1 is in phase 2, maxCompleted is phase 1 → should inject.
	result := PhaseDetect(dir, "test-feature", "2.1")
	if result == "" {
		t.Error("expected non-empty summary path, got empty string")
	}
	if !strings.Contains(result, "1-summary.md") {
		t.Errorf("expected path to contain '1-summary.md', got: %s", result)
	}
}

func TestPhaseDetect_SamePhase(t *testing.T) {
	// currentPhase == maxCompletedPhase → no injection
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Phase 1 task",
			Status: "completed",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
		},
	}
	setupFeatureDir(t, dir, tasks)

	// Task 1.2 is in phase 1, maxCompleted is also phase 1 → no injection.
	result := PhaseDetect(dir, "test-feature", "1.2")
	if result != "" {
		t.Errorf("expected empty string for same-phase task, got: %s", result)
	}
}

func TestPhaseDetect_FirstPhase(t *testing.T) {
	// Phase 1 tasks never get a summary injected (currentPhase > 1 guard).
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Phase 1 task",
			Status: "pending",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
		},
	}
	setupFeatureDir(t, dir, tasks)

	result := PhaseDetect(dir, "test-feature", "1.1")
	if result != "" {
		t.Errorf("expected empty string for phase-1 task, got: %s", result)
	}
}

func TestPhaseDetect_SummaryFileMissing(t *testing.T) {
	// Summary file does not exist → return empty string (not an error).
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1-impl": {
			ID:     "1.1",
			Title:  "Phase 1 task",
			Status: "completed",
			File:   "1.1-impl.md",
			Record: "records/1.1-impl.md",
			Type:   task.TypeImplementation,
		},
	}
	setupFeatureDir(t, dir, tasks)

	// No summary file created — PhaseDetect should return "" gracefully.
	result := PhaseDetect(dir, "test-feature", "2.1")
	if result != "" {
		t.Errorf("expected empty string when summary file missing, got: %s", result)
	}
}

// --- InferType tests ---

func TestInferType(t *testing.T) {
	tests := []struct {
		id       string
		expected string
	}{
		// Summary suffix
		{"1.summary", task.TypeDocGenerationSummary},
		{"2.summary", task.TypeDocGenerationSummary},
		// Gate suffix
		{"1.gate", task.TypeGate},
		{"3.gate", task.TypeGate},
		// T-test exact IDs
		{"T-test-1", task.TypeTestPipelineGenCases},
		{"T-test-1b", task.TypeTestPipelineEvalCases},
		{"T-test-2", task.TypeTestPipelineGenScripts},
		{"T-test-3", task.TypeTestPipelineRun},
		{"T-test-4", task.TypeTestPipelineGraduate},
		{"T-test-4.5", task.TypeTestPipelineVerifyRegression},
		{"T-test-5", task.TypeDocGenerationConsolidate},
		// T-quick-6 drift detection
		{"T-quick-6", task.TypeDocGenerationDrift},
		{"T-quick-6a", task.TypeDocGenerationDrift},
		// Fix prefix
		{"fix-1", task.TypeFix},
		{"fix-auth-bug", task.TypeFix},
		{"disc-1", task.TypeFix},
		{"disc-2", task.TypeFix},
		// Default: empty (no fallback)
		{"1.1", ""},
		{"2.3", ""},
		{"3.1-some-task", ""},
		{"1.1-task-type-fields", ""},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := InferType(tt.id)
			if got != tt.expected {
				t.Errorf("InferType(%q) = %q, want %q", tt.id, got, tt.expected)
			}
		})
	}
}
