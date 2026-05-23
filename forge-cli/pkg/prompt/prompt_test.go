package prompt

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode"

	"forge-cli/pkg/task"
)

// intPtr returns a pointer to the given int value.
func intPtr(v int) *int {
	return &v
}

// containsChinese returns true if the string contains any CJK Unified Ideograph characters.
func containsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

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
		task.TypeCodingFeature,
		task.TypeCodingEnhancement,
		task.TypeCodingCleanup,
		task.TypeCodingRefactor,
		task.TypeDoc,
		task.TypeDocEval,
		task.TypeDocSummary,
		task.TypeDocConsolidate,
		task.TypeDocDrift,
		task.TypeTestGenCases,
		task.TypeTestEvalCases,
		task.TypeTestGenScripts,
		task.TypeTestRun,
		task.TypeTestGraduate,
		task.TypeTestVerifyRegression,
		task.TypeCodingFix,
		task.TypeGate,
		task.TypeCleanCode,
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
			Type:   task.TypeCodingFeature,
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

// --- New type template content tests ---

func TestSynthesize_FeatureTemplate_HasTDDWorkflow(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Feature task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingFeature,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "TDD") {
		t.Error("feature template should mention TDD workflow")
	}
	if !strings.Contains(result, "RED") {
		t.Error("feature template should mention RED step")
	}
}

func TestSynthesize_EnhancementTemplate_HasTDDWorkflow(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Enhancement task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingEnhancement,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "TDD") {
		t.Error("enhancement template should mention TDD workflow")
	}
}

func TestSynthesize_CleanupTemplate_NoTDD(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Cleanup task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingCleanup,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "RED") {
		t.Error("cleanup template should NOT mention RED/TDD step")
	}
	if !strings.Contains(result, "Targeted Tests") {
		t.Error("cleanup template should mention Targeted Tests")
	}
}

func TestSynthesize_RefactorTemplate_NoTDD(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Refactor task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingRefactor,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "RED") {
		t.Error("refactor template should NOT mention RED/TDD step")
	}
	if !strings.Contains(result, "behavior") {
		t.Error("refactor template should mention behavior preservation")
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
			Type:   task.TypeCodingFeature,
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
			Type:   task.TypeCodingFeature,
			Scope:  "backend",
		},
		"2.1-impl": {
			ID:     "2.1",
			Title:  "Phase 2 task",
			Status: "pending",
			File:   "2.1-impl.md",
			Record: "records/2.1-impl.md",
			Type:   task.TypeCodingFeature,
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
			Type:   task.TypeCodingFeature,
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
			Type:   task.TypeCodingFeature,
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
			Type:   task.TypeCodingFeature,
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
			Type:   task.TypeCodingFeature,
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
			Type:   task.TypeCodingFeature,
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
		{"1.summary", task.TypeDocSummary},
		{"2.summary", task.TypeDocSummary},
		// Gate suffix
		{"1.gate", task.TypeGate},
		{"3.gate", task.TypeGate},
		// T-test exact IDs
		{"T-test-gen-cases", task.TypeTestGenCases},
		{"T-test-eval-cases", task.TypeTestEvalCases},
		{"T-test-gen-scripts", task.TypeTestGenScripts},
		{"T-test-run", task.TypeTestRun},
		{"T-test-graduate", task.TypeTestGraduate},
		{"T-test-verify-regression", task.TypeTestVerifyRegression},
		{"T-specs-consolidate", task.TypeDocConsolidate},
		// T-quick-doc-drift drift detection
		{"T-quick-doc-drift", task.TypeDocDrift},
		// Fix prefix
		{"fix-1", task.TypeCodingFix},
		{"fix-auth-bug", task.TypeCodingFix},
		{"disc-1", task.TypeCodingFix},
		{"disc-2", task.TypeCodingFix},
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

func TestSynthesize_CleanCodeTemplate_InvokesSkill(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"T-clean-code": {
			ID:     "T-clean-code",
			Title:  "Clean code task",
			Status: "pending",
			File:   "T-clean-code-1.md",
			Record: "records/T-clean-code-1.md",
			Type:   task.TypeCleanCode,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "T-clean-code"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, `Skill(skill="forge:clean-code")`) {
		t.Error("clean-code template should invoke forge:clean-code skill")
	}
	if !strings.Contains(result, "T-clean-code") {
		t.Error("clean-code template should contain the task ID")
	}
}

// --- Test type suffix in gen-scripts template ---

func TestSynthesize_GenScripts_WithTypeSuffix(t *testing.T) {
	tests := []struct {
		name         string
		taskID       string
		wantContains string
		dontWant     string
	}{
		{
			name:         "T-test-gen-scripts-api includes --type api",
			taskID:       "T-test-gen-scripts-api",
			wantContains: `Skill(skill="forge:gen-test-scripts" --type api)`,
			dontWant:     `{{TEST_TYPE_ARG}}`,
		},
		{
			name:         "T-test-gen-scripts-tui includes --type tui",
			taskID:       "T-test-gen-scripts-tui",
			wantContains: `Skill(skill="forge:gen-test-scripts" --type tui)`,
			dontWant:     `{{TEST_TYPE_ARG}}`,
		},
		{
			name:         "T-quick-gen-and-run-cli includes --type cli",
			taskID:       "T-quick-gen-and-run-cli",
			wantContains: `Skill(skill="forge:gen-test-scripts" --type cli)`,
			dontWant:     `{{TEST_TYPE_ARG}}`,
		},
		{
			name:         "T-quick-gen-and-run-web-ui includes --type web-ui",
			taskID:       "T-quick-gen-and-run-web-ui",
			wantContains: `Skill(skill="forge:gen-test-scripts" --type web-ui)`,
			dontWant:     `{{TEST_TYPE_ARG}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tasks := map[string]task.Task{
				tt.taskID: {
					ID:     tt.taskID,
					Title:  "Gen scripts typed",
					Status: "pending",
					File:   tt.taskID + ".md",
					Record: "records/" + tt.taskID + ".md",
					Type:   task.TypeTestGenScripts,
					Scope:  "backend",
				},
			}
			setupFeatureDir(t, dir, tasks)

			opts := SynthesizeOpts{
				ProjectRoot: dir,
				FeatureSlug: "test-feature",
				TaskID:      tt.taskID,
			}
			result, err := Synthesize(opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(result, tt.wantContains) {
				t.Errorf("result should contain %q, got:\n%s", tt.wantContains, result)
			}
			if strings.Contains(result, tt.dontWant) {
				t.Errorf("result should not contain unreplaced placeholder %q", tt.dontWant)
			}
		})
	}
}

func TestSynthesize_GenScripts_NoTypeSuffix(t *testing.T) {
	// Ensure backward compatibility: no --type when no type suffix.
	tests := []struct {
		name   string
		taskID string
	}{
		{"T-test-gen-scripts", "T-test-gen-scripts"},
		{"T-test-gen-scriptsa", "T-test-gen-scriptsa"},
		{"T-quick-gen-and-run", "T-quick-gen-and-run"},
		{"T-quick-gen-and-runa", "T-quick-gen-and-runa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tasks := map[string]task.Task{
				tt.taskID: {
					ID:     tt.taskID,
					Title:  "Gen scripts",
					Status: "pending",
					File:   tt.taskID + ".md",
					Record: "records/" + tt.taskID + ".md",
					Type:   task.TypeTestGenScripts,
					Scope:  "backend",
				},
			}
			setupFeatureDir(t, dir, tasks)

			opts := SynthesizeOpts{
				ProjectRoot: dir,
				FeatureSlug: "test-feature",
				TaskID:      tt.taskID,
			}
			result, err := Synthesize(opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Contains(result, "--type") {
				t.Errorf("result should not contain --type for non-type-suffixed ID, got:\n%s", result)
			}
			// Should still contain the skill invocation without --type
			if !strings.Contains(result, `Skill(skill="forge:gen-test-scripts")`) {
				t.Errorf("result should contain skill invocation without --type, got:\n%s", result)
			}
		})
	}
}

// --- Consolidate/Drift non-interactive mode tests ---

func TestSynthesize_ConsolidateTemplate_NonInteractive(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"T-specs-consolidate": {
			ID:     "T-specs-consolidate",
			Title:  "Consolidate specs",
			Status: "pending",
			File:   "T-specs-consolidate.md",
			Record: "records/T-specs-consolidate.md",
			Type:   task.TypeDocConsolidate,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "T-specs-consolidate"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Must instruct non-interactive mode for pipeline execution.
	if !strings.Contains(result, "non-interactive") {
		t.Error("consolidate template should mention non-interactive mode")
	}

	// Must NOT instruct the agent to block or wait for user confirmation.
	if strings.Contains(result, "blocked") {
		t.Error("consolidate template should NOT mention 'blocked' — auto mode should proceed")
	}

	// Must invoke the consolidate-specs skill.
	if !strings.Contains(result, `Skill(skill="forge:consolidate-specs"`) {
		t.Error("consolidate template should invoke forge:consolidate-specs skill")
	}
}

func TestSynthesize_DriftTemplate_NonInteractive(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"T-quick-doc-drift": {
			ID:     "T-quick-doc-drift",
			Title:  "Drift detection",
			Status: "pending",
			File:   "T-quick-doc-drift.md",
			Record: "records/T-quick-doc-drift.md",
			Type:   task.TypeDocDrift,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "T-quick-doc-drift"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Must instruct non-interactive mode for pipeline execution.
	if !strings.Contains(result, "non-interactive") {
		t.Error("drift template should mention non-interactive mode")
	}

	// Must NOT instruct the agent to block or wait for user confirmation.
	if strings.Contains(result, "blocked") {
		t.Error("drift template should NOT mention 'blocked' — auto mode should proceed")
	}

	// Must invoke the consolidate-specs skill.
	if !strings.Contains(result, `Skill(skill="forge:consolidate-specs"`) {
		t.Error("drift template should invoke forge:consolidate-specs skill")
	}
}

// --- Coding principles injection tests ---

func TestSynthesize_CodingTemplates_ContainCodingPrinciples(t *testing.T) {
	codingTypes := []struct {
		typ     string
		typeVal string
	}{
		{"feature", task.TypeCodingFeature},
		{"enhancement", task.TypeCodingEnhancement},
		{"fix", task.TypeCodingFix},
		{"refactor", task.TypeCodingRefactor},
		{"cleanup", task.TypeCodingCleanup},
	}

	for _, ct := range codingTypes {
		t.Run(ct.typ, func(t *testing.T) {
			dir := t.TempDir()
			taskID := "1.1"
			tasks := map[string]task.Task{
				taskID: {
					ID:     taskID,
					Title:  "Coding task",
					Status: "pending",
					File:   "1.1.md",
					Record: "records/1.1.md",
					Type:   ct.typeVal,
					Scope:  "backend",
				},
			}
			setupFeatureDir(t, dir, tasks)

			opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: taskID}
			result, err := Synthesize(opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !strings.Contains(result, "<CODING_PRINCIPLES>") {
				t.Errorf("%s template: synthesized prompt missing <CODING_PRINCIPLES> opening tag", ct.typ)
			}
			if !strings.Contains(result, "</CODING_PRINCIPLES>") {
				t.Errorf("%s template: synthesized prompt missing </CODING_PRINCIPLES> closing tag", ct.typ)
			}
		})
	}
}

// --- Coverage target injection tests ---

func TestSynthesize_CodingFeature_DefaultCoverageTarget(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Feature task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingFeature,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "percentage") {
		t.Error("coding.feature prompt should contain COVERAGE_STRATEGY 'percentage'")
	}
	if !strings.Contains(result, "Achieve 80% test coverage") {
		t.Error("coding.feature prompt should contain COVERAGE_TARGET 'Achieve 80% test coverage'")
	}
}

func TestSynthesize_CodingFix_DefaultCoverageTarget(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"fix-1": {
			ID:     "fix-1",
			Title:  "Fix task",
			Status: "pending",
			File:   "fix-1.md",
			Record: "records/fix-1.md",
			Type:   task.TypeCodingFix,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "fix-1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "percentage") {
		t.Error("coding.fix prompt should contain COVERAGE_STRATEGY 'percentage'")
	}
	if !strings.Contains(result, "Achieve 60% test coverage") {
		t.Error("coding.fix prompt should contain COVERAGE_TARGET 'Achieve 60% test coverage'")
	}
}

func TestSynthesize_CodingRefactor_MaintainStrategy(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Refactor task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingRefactor,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "maintain") {
		t.Error("coding.refactor prompt should contain COVERAGE_STRATEGY 'maintain'")
	}
	if !strings.Contains(result, "Maintain existing coverage, no more than 2% decrease") {
		t.Error("coding.refactor prompt should contain COVERAGE_TARGET 'Maintain existing coverage, no more than 2% decrease'")
	}
}

func TestSynthesize_FrontmatterCoverageOverridesConfig(t *testing.T) {
	dir := t.TempDir()
	coverage90 := 90
	tasks := map[string]task.Task{
		"1.1": {
			ID:       "1.1",
			Title:    "Feature task with custom coverage",
			Status:   "pending",
			File:     "1.1.md",
			Record:   "records/1.1.md",
			Type:     task.TypeCodingFeature,
			Scope:    "backend",
			Coverage: &coverage90,
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Frontmatter coverage=90 should override the default 80
	if !strings.Contains(result, "Achieve 90% test coverage") {
		t.Errorf("frontmatter coverage=90 should produce 'Achieve 90%% test coverage', got:\n%s", result)
	}
}

func TestSynthesize_NonTestableType_NoCoverageInjection(t *testing.T) {
	nonTestableTypes := []string{
		task.TypeDoc,
		task.TypeDocEval,
		task.TypeDocSummary,
		task.TypeDocConsolidate,
		task.TypeDocDrift,
		task.TypeGate,
		task.TypeTestGenCases,
		task.TypeTestEvalCases,
		task.TypeTestGenScripts,
		task.TypeTestRun,
		task.TypeCleanCode,
	}

	for _, typ := range nonTestableTypes {
		t.Run(typ, func(t *testing.T) {
			dir := t.TempDir()
			taskID := "1.1"
			tasks := map[string]task.Task{
				taskID: {
					ID:     taskID,
					Title:  "Non-testable task",
					Status: "pending",
					File:   "1.1.md",
					Record: "records/1.1.md",
					Type:   typ,
					Scope:  "backend",
				},
			}
			setupFeatureDir(t, dir, tasks)

			opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: taskID}
			result, err := Synthesize(opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if strings.Contains(result, "COVERAGE_TARGET") {
				t.Errorf("%s template: should not contain unreplaced COVERAGE_TARGET placeholder", typ)
			}
			if strings.Contains(result, "COVERAGE_STRATEGY") {
				t.Errorf("%s template: should not contain unreplaced COVERAGE_STRATEGY placeholder", typ)
			}
		})
	}
}

func TestSynthesize_ConfigCoverageOverridesDefault(t *testing.T) {
	// Create a config.yaml with custom coverage for coding.feature
	dir := t.TempDir()
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configContent := `coverage:
  coding.feature:
    type: percentage
    percentage: 75
`
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}

	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Feature task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingFeature,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Config overrides default: 75 instead of 80
	if !strings.Contains(result, "Achieve 75% test coverage") {
		t.Errorf("config coverage=75 should produce 'Achieve 75%% test coverage', got:\n%s", result)
	}
}

// --- Coverage directive actionability tests ---
// Verify that templates contain actionable instructions guiding agent behavior.

func TestSynthesize_CodingFeature_CoverageDirectiveActionable(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Feature task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingFeature,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Feature template should tell agent to stop adding tests when target is reached
	if !strings.Contains(result, "Stop adding tests once the target is reached") {
		t.Error("coding.feature prompt should contain actionable directive: 'Stop adding tests once the target is reached'")
	}
}

func TestSynthesize_CodingEnhancement_CoverageDirectiveActionable(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Enhancement task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingEnhancement,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Enhancement template should tell agent to stop adding tests when target is reached
	if !strings.Contains(result, "Stop adding tests once the target is reached") {
		t.Error("coding.enhancement prompt should contain actionable directive: 'Stop adding tests once the target is reached'")
	}
}

func TestSynthesize_CodingFix_CoverageDirectiveActionable(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"fix-1": {
			ID:     "fix-1",
			Title:  "Fix task",
			Status: "pending",
			File:   "fix-1.md",
			Record: "records/fix-1.md",
			Type:   task.TypeCodingFix,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "fix-1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Fix template should tell agent to write targeted tests, not chase high coverage
	if !strings.Contains(result, "Write targeted fix tests") {
		t.Error("coding.fix prompt should contain actionable directive: 'Write targeted fix tests'")
	}
	if !strings.Contains(result, "stop adding once the target is reached") {
		t.Error("coding.fix prompt should contain actionable directive: 'stop adding once the target is reached'")
	}
}

func TestSynthesize_CodingRefactor_CoverageDirectiveActionable(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Refactor task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingRefactor,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Refactor template should contain incremental compile strategy
	if !strings.Contains(result, "just compile") {
		t.Error("coding.refactor prompt should contain incremental compile strategy with 'just compile'")
	}
	if !strings.Contains(result, "no new tests") {
		t.Error("coding.refactor prompt should contain directive: 'no new tests'")
	}
}

func TestSynthesize_CodingCleanup_CoverageDirectiveActionable(t *testing.T) {
	dir := t.TempDir()
	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Cleanup task",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingCleanup,
			Scope:  "backend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Cleanup template should say no new tests needed, maintain strategy
	if !strings.Contains(result, "no new tests") {
		t.Error("coding.cleanup prompt should contain directive: 'no new tests'")
	}
}

// --- resolveCoverage language and cleanup/refactor tests ---

func TestResolveCoverage_EnglishText(t *testing.T) {
	tests := []struct {
		name         string
		task         task.Task
		wantStrategy string
		wantTarget   string
	}{
		{
			name: "frontmatter coverage returns English",
			task: task.Task{
				ID:       "1.1",
				Type:     task.TypeCodingFeature,
				Coverage: intPtr(90),
			},
			wantStrategy: "percentage",
			wantTarget:   "Achieve 90% test coverage",
		},
		{
			name: "refactor maintain strategy returns English",
			task: task.Task{
				ID:   "1.1",
				Type: task.TypeCodingRefactor,
			},
			wantStrategy: "maintain",
			wantTarget:   "Maintain existing coverage, no more than 2% decrease",
		},
		{
			name: "cleanup maintain strategy returns English",
			task: task.Task{
				ID:   "1.1",
				Type: task.TypeCodingCleanup,
			},
			wantStrategy: "maintain",
			wantTarget:   "Maintain existing coverage, no more than 2% decrease",
		},
		{
			name: "feature default returns English",
			task: task.Task{
				ID:   "1.1",
				Type: task.TypeCodingFeature,
			},
			wantStrategy: "percentage",
			wantTarget:   "Achieve 80% test coverage",
		},
		{
			name: "fix default returns English",
			task: task.Task{
				ID:   "fix-1",
				Type: task.TypeCodingFix,
			},
			wantStrategy: "percentage",
			wantTarget:   "Achieve 60% test coverage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			strategy, target := resolveCoverage(dir, tt.task)
			if strategy != tt.wantStrategy {
				t.Errorf("strategy = %q, want %q", strategy, tt.wantStrategy)
			}
			if target != tt.wantTarget {
				t.Errorf("target = %q, want %q", target, tt.wantTarget)
			}
			// Verify no Chinese characters in output
			if containsChinese(strategy) {
				t.Errorf("strategy contains Chinese characters: %q", strategy)
			}
			if containsChinese(target) {
				t.Errorf("target contains Chinese characters: %q", target)
			}
		})
	}
}

func TestResolveCoverage_CleanupAndRefactor_SkipPercentageDirective(t *testing.T) {
	// Even when frontmatter specifies a coverage percentage, cleanup/refactor
	// should not inject a percentage target because their templates say "No new tests".
	tests := []struct {
		name       string
		task       task.Task
		wantStrat  string
		wantTarget string
	}{
		{
			name: "cleanup with frontmatter coverage skips percentage",
			task: task.Task{
				ID:       "1.1",
				Type:     task.TypeCodingCleanup,
				Coverage: intPtr(80),
			},
			wantStrat:  "maintain",
			wantTarget: "Maintain existing coverage, no more than 2% decrease",
		},
		{
			name: "refactor with frontmatter coverage skips percentage",
			task: task.Task{
				ID:       "1.1",
				Type:     task.TypeCodingRefactor,
				Coverage: intPtr(90),
			},
			wantStrat:  "maintain",
			wantTarget: "Maintain existing coverage, no more than 2% decrease",
		},
		{
			name: "cleanup default skips percentage",
			task: task.Task{
				ID:   "1.1",
				Type: task.TypeCodingCleanup,
			},
			wantStrat:  "maintain",
			wantTarget: "Maintain existing coverage, no more than 2% decrease",
		},
		{
			name: "refactor default skips percentage",
			task: task.Task{
				ID:   "1.1",
				Type: task.TypeCodingRefactor,
			},
			wantStrat:  "maintain",
			wantTarget: "Maintain existing coverage, no more than 2% decrease",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			strategy, target := resolveCoverage(dir, tt.task)
			if strategy != tt.wantStrat {
				t.Errorf("strategy = %q, want %q", strategy, tt.wantStrat)
			}
			if target != tt.wantTarget {
				t.Errorf("target = %q, want %q", target, tt.wantTarget)
			}
		})
	}
}

// --- Scope resolution fallback tests ---

func TestResolveScope_MismatchFallsBack(t *testing.T) {
	tests := []struct {
		name        string
		projectType string
		taskScope   string
		wantScope   string
	}{
		{
			name:        "backend project with frontend scope falls back to empty",
			projectType: "backend",
			taskScope:   "frontend",
			wantScope:   "",
		},
		{
			name:        "backend project with backend scope keeps scope",
			projectType: "backend",
			taskScope:   "backend",
			wantScope:   "backend",
		},
		{
			name:        "fullstack project with frontend scope keeps scope",
			projectType: "fullstack",
			taskScope:   "frontend",
			wantScope:   "frontend",
		},
		{
			name:        "fullstack project with backend scope keeps scope",
			projectType: "fullstack",
			taskScope:   "backend",
			wantScope:   "backend",
		},
		{
			name:        "no project type keeps scope",
			projectType: "",
			taskScope:   "frontend",
			wantScope:   "frontend",
		},
		{
			name:        "empty scope stays empty",
			projectType: "backend",
			taskScope:   "",
			wantScope:   "",
		},
		{
			name:        "all scope becomes empty",
			projectType: "backend",
			taskScope:   "all",
			wantScope:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			// Create config.yaml if projectType is set
			if tt.projectType != "" {
				forgeDir := filepath.Join(dir, ".forge")
				if err := os.MkdirAll(forgeDir, 0o755); err != nil {
					t.Fatal(err)
				}
				configContent := "project-type: " + tt.projectType + "\n"
				if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
					t.Fatal(err)
				}
			}

			got := resolveScope(dir, tt.taskScope)
			if got != tt.wantScope {
				t.Errorf("resolveScope(%q, %q) = %q, want %q", tt.projectType, tt.taskScope, got, tt.wantScope)
			}
		})
	}
}

func TestSynthesize_ScopeMismatch_GeneratesDefaultCommand(t *testing.T) {
	// Backend project + frontend scope → scope should be empty → no scope suffix on commands
	dir := t.TempDir()
	forgeDir := filepath.Join(dir, ".forge")
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configContent := "project-type: backend\n"
	if err := os.WriteFile(filepath.Join(forgeDir, "config.yaml"), []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}

	tasks := map[string]task.Task{
		"1.1": {
			ID:     "1.1",
			Title:  "Frontend task in backend project",
			Status: "pending",
			File:   "1.1.md",
			Record: "records/1.1.md",
			Type:   task.TypeCodingFeature,
			Scope:  "frontend",
		},
	}
	setupFeatureDir(t, dir, tasks)

	opts := SynthesizeOpts{ProjectRoot: dir, FeatureSlug: "test-feature", TaskID: "1.1"}
	result, err := Synthesize(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should not contain "SCOPE: frontend" since it's a mismatch
	if strings.Contains(result, "SCOPE: frontend") {
		t.Errorf("result should not contain 'SCOPE: frontend' for mismatched scope, got:\n%s", result)
	}

	// Should contain default commands without scope suffix
	if !strings.Contains(result, "just compile") {
		t.Error("result should contain 'just compile' default command")
	}
}

func TestExtractTestTypeArg(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"T-test-gen-scripts-api", " --type api"},
		{"T-test-gen-scripts-tui", " --type tui"},
		{"T-quick-gen-and-run-cli", " --type cli"},
		{"T-quick-gen-and-run-web-ui", " --type web-ui"},
		{"T-test-gen-scripts", ""},
		{"T-test-gen-scripts", ""},
		{"T-quick-gen-and-run", ""},
		{"T-test-3-api", ""}, // not a gen-scripts base
		{"1.1", ""},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := extractTestTypeArg(tt.id)
			if got != tt.want {
				t.Errorf("extractTestTypeArg(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}
