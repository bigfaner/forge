package contract

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// --- Stub FactCollector for testing ---

type stubFactCollector struct {
	facts FactTable
	err   error
}

func (s stubFactCollector) Collect(_ string) (FactTable, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.facts, nil
}

// --- Test: Verify with no contracts ---

func TestVerify_NoContracts_ReturnsEmptySummary(t *testing.T) {
	dir := t.TempDir()

	summary, err := Verify(dir, stubFactCollector{facts: FactTable{}})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if summary.Total != 0 {
		t.Errorf("expected Total=0, got %d", summary.Total)
	}
	if summary.Broken != 0 {
		t.Errorf("expected Broken=0, got %d", summary.Broken)
	}
	if summary.OK != 0 {
		t.Errorf("expected OK=0, got %d", summary.OK)
	}
}

// --- Test: Verify with matching Contract (OK) ---

func TestVerify_MatchingContract_ReturnsOK(t *testing.T) {
	dir := t.TempDir()

	// Create contract file
	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "task-lifecycle"
step: 2
step-action: "forge task claim"
---

# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "success"
- Preconditions: "feature exists with slug matching arg; at least one task available"
- Input: (positional: none; flags: none)
- Output: stdout contains "claimed task", exit code 0
- State: tasks status changed to in_progress; index.json updated
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-2-task-claim.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Provide matching facts
	facts := FactTable{
		"forge task claim": {
			Command:  "forge task claim",
			Stdout:   "claimed task T-001 successfully",
			ExitCode: 0,
		},
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if summary.Total != 1 {
		t.Errorf("expected Total=1, got %d", summary.Total)
	}
	if summary.Broken != 0 {
		t.Errorf("expected Broken=0 (no breaks), got %d; results: %+v", summary.Broken, summary.Results)
	}
	if summary.OK != 1 {
		t.Errorf("expected OK=1, got %d", summary.OK)
	}
}

// --- Test: Verify with broken Contract (mismatch) ---

func TestVerify_BrokenContract_DetectsMismatch(t *testing.T) {
	dir := t.TempDir()

	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "task-lifecycle"
step: 2
step-action: "forge task claim"
---

# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "success"
- Preconditions: "feature exists with tasks available"
- Input: (positional: none; flags: none)
- Output: stdout contains "claimed task successfully", exit code 0
- State: tasks status changed to in_progress; index.json updated
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-2-task-claim.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Provide facts where stdout doesn't contain "successfully" (truly different output)
	facts := FactTable{
		"forge task claim": {
			Command:  "forge task claim",
			Stdout:   "Task T-001 claimed with errors",
			ExitCode: 0,
		},
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if summary.Broken == 0 {
		t.Errorf("expected Broken>0 for mismatched contract, got 0; results: %+v", summary.Results)
	}
}

// --- Test: Verify with exit code mismatch ---

func TestVerify_ExitCodeMismatch_DetectsBreak(t *testing.T) {
	dir := t.TempDir()

	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "task-lifecycle"
step: 3
step-action: "forge task submit"
---

# Contract: task-lifecycle / Step 3: forge task submit

## Outcome "success"
- Preconditions: "task status is in_progress"
- Input: --result success
- Output: stdout contains "submitted successfully", exit code 0
- State: status changed to completed
- Side-effect: none

## Journey Invariants
- task_id stable once assigned
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-3-task-submit.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Exit code doesn't match
	facts := FactTable{
		"forge task submit": {
			Command:  "forge task submit",
			Stdout:   "submitted successfully",
			ExitCode: 1, // mismatch: contract expects 0
		},
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if summary.Broken == 0 {
		t.Errorf("expected Broken>0 for exit code mismatch, got 0")
	}

	// Check the break details
	found := false
	for _, r := range summary.Results {
		for _, b := range r.Breaks {
			if b.Dimension == "Output" && strings.Contains(b.Expected, "exit code 0") {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expected Output dimension break with exit code mismatch")
	}
}

// --- Test: Verify with stderr assertion ---

func TestVerify_StderrAssertion_ChecksStderr(t *testing.T) {
	dir := t.TempDir()

	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "task-lifecycle"
step: 2
step-action: "forge task claim"
---

# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "no-tasks-available"
- Preconditions: "feature exists; no tasks available for claiming"
- Input: (positional: none; flags: none)
- Output: stderr contains "no tasks available", exit code 1
- State: unchanged
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-2-task-claim.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Stderr has the expected message
	facts := FactTable{
		"forge task claim": {
			Command:  "forge task claim",
			Stdout:   "",
			Stderr:   "no tasks available for claiming",
			ExitCode: 1,
		},
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if summary.Broken != 0 {
		t.Errorf("expected Broken=0 for matching stderr assertion, got %d; breaks: %+v", summary.Broken, summary.Results)
	}
}

// --- Test: Verify with "unchanged" State always passes ---

func TestVerify_UnchangedState_AlwaysPasses(t *testing.T) {
	dir := t.TempDir()

	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "task-lifecycle"
step: 2
step-action: "forge task claim"
---

# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "no-tasks-available"
- Preconditions: "no tasks available"
- Input: none
- Output: stderr contains "no tasks available", exit code 1
- State: unchanged
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-2-task-claim.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	facts := FactTable{
		"forge task claim": {
			Command:  "forge task claim",
			Stderr:   "no tasks available",
			ExitCode: 1,
		},
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if summary.OK != 1 {
		t.Errorf("expected OK=1 for unchanged state, got OK=%d, Broken=%d", summary.OK, summary.Broken)
	}
}

// --- Test: Verify with no matching fact entry (no false positive) ---

func TestVerify_NoMatchingFact_NoFalsePositive(t *testing.T) {
	dir := t.TempDir()

	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "task-lifecycle"
step: 2
step-action: "forge task claim"
---

# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "success"
- Preconditions: "feature exists with tasks available"
- Input: none
- Output: stdout contains "claimed task", exit code 0
- State: tasks status changed to in_progress
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-2-task-claim.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	// No matching fact entry for "forge task claim"
	facts := FactTable{}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	// No facts means no comparison possible; should report OK (no false positives)
	if summary.OK != 1 {
		t.Errorf("expected OK=1 when no facts available (avoid false positive), got OK=%d, Broken=%d", summary.OK, summary.Broken)
	}
}

// --- Test: Zero false positives on 20+ unchanged contracts ---

func TestVerify_TwentyUnchangedContracts_ZeroFalsePositives(t *testing.T) {
	dir := t.TempDir()

	// Create 22 contracts that should all pass
	for i := 1; i <= 22; i++ {
		journeyDir := filepath.Join(dir, "tests", fmt.Sprintf("journey-%d", i/8+1), "_contracts")
		if err := os.MkdirAll(journeyDir, 0755); err != nil {
			t.Fatal(err)
		}

		cmd := fmt.Sprintf("forge cmd-%d", i)
		contractContent := fmt.Sprintf(`---
journey: "journey-%d"
step: %d
step-action: "%s"
---

# Contract: journey-%d / Step %d: %s

## Outcome "success"
- Preconditions: "feature exists"
- Input: none
- Output: stdout contains "success output %d", exit code 0
- State: status changed to completed
- Side-effect: none

## Journey Invariants
- feature_slug consistent across all steps
`, i/8+1, i, cmd, i/8+1, i, cmd, i)

		filename := fmt.Sprintf("step-%d-action.md", i)
		if err := os.WriteFile(filepath.Join(journeyDir, filename), []byte(contractContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Add matching fact
	}

	// Build fact table with matching facts for all 22 commands
	facts := make(FactTable)
	for i := 1; i <= 22; i++ {
		cmd := fmt.Sprintf("forge cmd-%d", i)
		facts[cmd] = FactEntry{
			Command:  cmd,
			Stdout:   fmt.Sprintf("success output %d", i),
			ExitCode: 0,
		}
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if summary.Total != 22 {
		t.Errorf("expected Total=22, got %d", summary.Total)
	}
	if summary.Broken != 0 {
		t.Errorf("ZERO FALSE POSITIVES REQUIRED: expected Broken=0 for 22 unchanged contracts, got %d", summary.Broken)
		for _, r := range summary.Results {
			if !r.OK {
				t.Errorf("  Broken: %s breaks=%+v", r.ContractPath, r.Breaks)
			}
		}
	}
	if summary.OK != 22 {
		t.Errorf("expected OK=22, got %d", summary.OK)
	}
}

// --- Test: FormatReport output format ---

func TestVerifySummary_FormatReport(t *testing.T) {
	summary := VerifySummary{
		Total:  23,
		Broken: 2,
		OK:     21,
		Results: []VerifyResult{
			{
				ContractPath: "tests/task-lifecycle/_contracts/step-2-task-claim.md",
				Journey:      "task-lifecycle",
				Step:         2,
				OK:           false,
				Breaks: []VerifyBreak{
					{
						Dimension: "Output",
						Outcome:   "success",
						Expected:  "claimed task <task_id>",
						Actual:    "Task <task_id> claimed",
						MatchType: "partial",
					},
				},
			},
			{
				ContractPath: "tests/task-lifecycle/_contracts/step-3-task-submit.md",
				Journey:      "task-lifecycle",
				Step:         3,
				OK:           false,
				Breaks: []VerifyBreak{
					{
						Dimension: "State",
						Outcome:   "success",
						Expected:  "status → completed",
						Actual:    "status → done",
						MatchType: "none",
					},
				},
			},
		},
	}

	report := summary.FormatReport()

	// Verify format matches proposal example
	if !strings.Contains(report, "Scanning 23 Contracts against Fact Table") {
		t.Errorf("report should contain scanning header, got: %s", report)
	}
	if !strings.Contains(report, "BROKEN (2):") {
		t.Errorf("report should contain BROKEN count, got: %s", report)
	}
	if !strings.Contains(report, "OK (21):") {
		t.Errorf("report should contain OK count, got: %s", report)
	}
	if !strings.Contains(report, "Summary: 2 broken, 21 OK, 0 false positives") {
		t.Errorf("report should contain summary line, got: %s", report)
	}
	if !strings.Contains(report, "step-2-task-claim.md") {
		t.Errorf("report should contain broken contract path, got: %s", report)
	}
	if !strings.Contains(report, "Output dimension:") {
		t.Errorf("report should contain dimension name, got: %s", report)
	}
	if !strings.Contains(report, "State dimension:") {
		t.Errorf("report should contain State dimension, got: %s", report)
	}
}

func TestVerifySummary_FormatReport_AllOK(t *testing.T) {
	summary := VerifySummary{
		Total:  3,
		Broken: 0,
		OK:     3,
		Results: []VerifyResult{
			{ContractPath: "tests/j1/_contracts/step-1.md", OK: true},
			{ContractPath: "tests/j2/_contracts/step-1.md", OK: true},
			{ContractPath: "tests/j3/_contracts/step-1.md", OK: true},
		},
	}

	report := summary.FormatReport()

	if strings.Contains(report, "BROKEN") {
		t.Errorf("report should NOT contain BROKEN when all OK, got: %s", report)
	}
	if !strings.Contains(report, "OK (3):") {
		t.Errorf("report should contain OK count, got: %s", report)
	}
	if !strings.Contains(report, "Summary: 0 broken, 3 OK, 0 false positives") {
		t.Errorf("report should contain summary, got: %s", report)
	}
}

// --- Test: Semantic matching ---

func TestSemanticMatch_ExactMatch(t *testing.T) {
	if !semanticMatch("claimed task", "claimed task T-001 successfully") {
		t.Error("expected exact match")
	}
}

func TestSemanticMatch_WordOrderChanged(t *testing.T) {
	// "claimed task" should match "Task T-001 claimed" because both words appear
	if !semanticMatch("claimed task", "Task T-001 claimed") {
		t.Error("expected match with word order changed")
	}
}

func TestSemanticMatch_CaseInsensitive(t *testing.T) {
	if !semanticMatch("SUCCESS OUTPUT", "success output from command") {
		t.Error("expected case-insensitive match")
	}
}

func TestSemanticMatch_PlaceholderIgnored(t *testing.T) {
	if !semanticMatch("claimed task <task_id>", "claimed task T-001") {
		t.Error("expected match with placeholder ignored")
	}
}

func TestSemanticMatch_NoMatch(t *testing.T) {
	if semanticMatch("feature created successfully", "error: feature not found") {
		t.Error("expected no match for unrelated text")
	}
}

func TestSemanticMatch_EmptyDescriptor(t *testing.T) {
	if !semanticMatch("", "anything") {
		t.Error("empty descriptor should always match")
	}
}

func TestSemanticMatch_StopWordsIgnored(t *testing.T) {
	if !semanticMatch("the task is claimed", "task claimed") {
		t.Error("stop words should be ignored in matching")
	}
}

func TestSemanticMatch_PartialMatch(t *testing.T) {
	// "claimed task successfully" vs "Task T-001 claimed" -- "successfully" is missing
	if semanticMatch("claimed task successfully", "Task T-001 claimed") {
		t.Error("partial match (missing key term) should not match")
	}
}

// --- Test: Extract semantic tokens ---

func TestExtractSemanticTokens(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantCount int
		wantToken string // at least one of these should be present
	}{
		{"simple phrase", "claimed task", 2, "claimed"},
		{"with placeholder", "claimed task <task_id>", 2, "claimed"},
		{"with stop words", "the feature is created successfully", 3, "created"},
		{"empty string", "", 0, ""},
		{"punctuation", "success, confirmation containing feature-slug", 4, "confirmation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := extractSemanticTokens(tt.input)
			if len(tokens) != tt.wantCount {
				t.Errorf("expected %d tokens, got %d: %v", tt.wantCount, len(tokens), tokens)
			}
			if tt.wantToken != "" {
				found := false
				for _, tok := range tokens {
					if tok == tt.wantToken {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected token %q in %v", tt.wantToken, tokens)
				}
			}
		})
	}
}

// --- Test: Determine match type ---

func TestDetermineMatchType(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		actual   string
		want     string
	}{
		{"all tokens match", "claimed task", "claimed task T-001", "exact"},
		{"partial match", "claimed task successfully", "Task T-001 claimed", "partial"},
		{"no match", "feature created", "error not found", "none"},
		{"empty expected", "", "anything", "exact"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineMatchType(tt.expected, tt.actual)
			if got != tt.want {
				t.Errorf("determineMatchType(%q, %q) = %q, want %q", tt.expected, tt.actual, got, tt.want)
			}
		})
	}
}

// --- Test: Contract file discovery ---

func TestDiscoverContractFiles_FindsContracts(t *testing.T) {
	dir := t.TempDir()

	// Create contract directory structure
	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create contract files
	for _, name := range []string{"step-1-feature-create.md", "step-2-task-claim.md"} {
		if err := os.WriteFile(filepath.Join(contractDir, name), []byte("# Contract"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create non-contract file (should be ignored)
	if err := os.WriteFile(filepath.Join(contractDir, "README.md"), []byte("# Readme"), 0644); err != nil {
		t.Fatal(err)
	}

	files, err := discoverContractFiles(dir)
	if err != nil {
		t.Fatalf("discoverContractFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 contract files, got %d: %v", len(files), files)
	}

	for _, f := range files {
		base := filepath.Base(f)
		if base != "step-1-feature-create.md" && base != "step-2-task-claim.md" {
			t.Errorf("unexpected file: %s", f)
		}
	}
}

func TestDiscoverContractFiles_NoTestsDir(t *testing.T) {
	dir := t.TempDir()

	files, err := discoverContractFiles(dir)
	if err != nil {
		t.Fatalf("discoverContractFiles failed: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files when no tests dir, got %d", len(files))
	}
}

func TestDiscoverContractFiles_MultipleJourneys(t *testing.T) {
	dir := t.TempDir()

	for _, journey := range []string{"task-lifecycle", "feature-creation"} {
		contractDir := filepath.Join(dir, "tests", journey, "_contracts")
		if err := os.MkdirAll(contractDir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(contractDir, "step-1-action.md"), []byte("# Contract"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := discoverContractFiles(dir)
	if err != nil {
		t.Fatalf("discoverContractFiles failed: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 contract files across journeys, got %d", len(files))
	}
}

// --- Test: ParseContractMarkdown ---

func TestParseContractMarkdown_FullContract(t *testing.T) {
	content := `---
journey: "task-lifecycle"
step: 2
step-action: "forge task claim"
---

# Contract: task-lifecycle / Step 2: forge task claim

## Outcome "success"
- Preconditions: "feature exists with slug matching arg; at least one task available"
- Input: (positional: none; flags: none)
- Output: stdout contains "claimed task <task_id>", exit code 0
- State: tasks/<task_id>/status -> "in_progress"; index.json updated
- Side-effect: none

## Outcome "no-tasks-available"
- Preconditions: "feature exists; no tasks available for claiming"
- Input: (positional: none; flags: none)
- Output: stderr contains "no tasks available", exit code 1
- State: unchanged

## Journey Invariants
- feature_slug consistent across all steps
- task_id stable once assigned
`

	c, err := ParseContractMarkdown(content)
	if err != nil {
		t.Fatalf("ParseContractMarkdown failed: %v", err)
	}

	if c.Journey != "task-lifecycle" {
		t.Errorf("expected journey 'task-lifecycle', got %q", c.Journey)
	}
	if c.Step != 2 {
		t.Errorf("expected step 2, got %d", c.Step)
	}
	if c.Action != "forge task claim" {
		t.Errorf("expected action 'forge task claim', got %q", c.Action)
	}
	if len(c.Outcomes) != 2 {
		t.Fatalf("expected 2 outcomes, got %d", len(c.Outcomes))
	}
	if c.Outcomes[0].Name != "success" {
		t.Errorf("expected first outcome name 'success', got %q", c.Outcomes[0].Name)
	}
	if c.Outcomes[1].Name != "no-tasks-available" {
		t.Errorf("expected second outcome name 'no-tasks-available', got %q", c.Outcomes[1].Name)
	}
	if len(c.Invariants) != 2 {
		t.Errorf("expected 2 invariants, got %d", len(c.Invariants))
	}
}

func TestParseContractMarkdown_PreservesDimensions(t *testing.T) {
	content := `---
journey: "test"
step: 1
step-action: "forge test"
---

# Contract: test / Step 1: forge test

## Outcome "success"
- Preconditions: "feature exists"
- Input: arg1 arg2
- Output: stdout "success", exit code 0
- State: file created

## Journey Invariants
- test invariant
`

	c, err := ParseContractMarkdown(content)
	if err != nil {
		t.Fatalf("ParseContractMarkdown failed: %v", err)
	}

	o := c.Outcomes[0]
	if o.Preconditions == "" {
		t.Error("Preconditions should not be empty")
	}
	if o.Input == "" {
		t.Error("Input should not be empty")
	}
	if o.Output == "" {
		t.Error("Output should not be empty")
	}
	if o.State == "" {
		t.Error("State should not be empty")
	}
}

// --- Test: Extract command ---

func TestExtractCommand(t *testing.T) {
	tests := []struct {
		action string
		want   string
	}{
		{"forge task claim", "forge task claim"},
		{"forge task submit --result success", "forge task submit"},
		{"forge feature my-feature", "forge feature my-feature"},
		{"forge test", "forge test"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			got := extractCommand(tt.action)
			if got != tt.want {
				t.Errorf("extractCommand(%q) = %q, want %q", tt.action, got, tt.want)
			}
		})
	}
}

// --- Test: Extract semantic content ---

func TestExtractSemanticContent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{`stdout contains "claimed task", exit code 0`, `claimed task`},
		{`stderr contains "no tasks available", exit code 1`, `no tasks available`},
		{"success confirmation", "success confirmation"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractSemanticContent(tt.input)
			if got != tt.want {
				t.Errorf("extractSemanticContent(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// --- Test: Verify does not modify files (Hard Rule) ---

func TestVerify_DoesNotModifyFiles(t *testing.T) {
	dir := t.TempDir()

	contractDir := filepath.Join(dir, "tests", "j1", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "j1"
step: 1
step-action: "forge test"
---

# Contract: j1 / Step 1: forge test

## Outcome "success"
- Preconditions: "feature exists"
- Input: none
- Output: stdout contains "success", exit code 0
- State: status changed

## Journey Invariants
- feature_slug consistent
`
	contractFile := filepath.Join(contractDir, "step-1-test.md")
	if err := os.WriteFile(contractFile, []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Record file state before verify
	before, err := os.ReadFile(contractFile)
	if err != nil {
		t.Fatal(err)
	}

	// Run verify
	_, err = Verify(dir, stubFactCollector{facts: FactTable{}})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	// Verify file unchanged
	after, err := os.ReadFile(contractFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("Hard Rule violation: Verify modified a file")
	}

	// Check no new files created
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && path != contractFile {
			t.Errorf("Unexpected file created: %s", path)
		}
		return nil
	})
}

// --- Test: FactTable is freshly collected each run (Hard Rule) ---

func TestVerify_FreshFactTableEachRun(t *testing.T) {
	dir := t.TempDir()

	contractDir := filepath.Join(dir, "tests", "j1", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "j1"
step: 1
step-action: "forge test"
---

# Contract: j1 / Step 1: forge test

## Outcome "success"
- Preconditions: "feature exists"
- Input: none
- Output: stdout contains "success", exit code 0
- State: status changed

## Journey Invariants
- feature_slug consistent
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-1-test.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Track number of Collect calls
	collectCount := 0
	trackingCollector := &trackingFactCollector{
		facts: FactTable{
			"forge test": {Command: "forge test", Stdout: "success", ExitCode: 0},
		},
		onCollect: func() { collectCount++ },
	}

	_, err := Verify(dir, trackingCollector)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if collectCount != 1 {
		t.Errorf("expected exactly 1 Collect call (fresh each run), got %d", collectCount)
	}
}

type trackingFactCollector struct {
	facts     FactTable
	onCollect func()
}

func (t *trackingFactCollector) Collect(_ string) (FactTable, error) {
	t.onCollect()
	return t.facts, nil
}

// --- Test: Multiple contracts across journeys ---

func TestVerify_MultipleJourneys(t *testing.T) {
	dir := t.TempDir()

	// Journey 1: 2 contracts, both OK
	j1Dir := filepath.Join(dir, "tests", "journey-1", "_contracts")
	if err := os.MkdirAll(j1Dir, 0755); err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 2; i++ {
		content := fmt.Sprintf(`---
journey: "journey-1"
step: %d
step-action: "forge cmd-%d"
---

# Contract: journey-1 / Step %d: forge cmd-%d

## Outcome "success"
- Preconditions: "feature exists"
- Input: none
- Output: stdout contains "output %d", exit code 0
- State: status changed

## Journey Invariants
- feature_slug consistent
`, i, i, i, i, i)
		if err := os.WriteFile(filepath.Join(j1Dir, fmt.Sprintf("step-%d-action.md", i)), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Journey 2: 1 contract, broken
	j2Dir := filepath.Join(dir, "tests", "journey-2", "_contracts")
	if err := os.MkdirAll(j2Dir, 0755); err != nil {
		t.Fatal(err)
	}
	content := `---
journey: "journey-2"
step: 1
step-action: "forge cmd-3"
---

# Contract: journey-2 / Step 1: forge cmd-3

## Outcome "success"
- Preconditions: "feature exists"
- Input: none
- Output: stdout contains "expected output", exit code 0
- State: status changed

## Journey Invariants
- feature_slug consistent
`
	if err := os.WriteFile(filepath.Join(j2Dir, "step-1-action.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	facts := FactTable{
		"forge cmd-1": {Command: "forge cmd-1", Stdout: "output 1", ExitCode: 0},
		"forge cmd-2": {Command: "forge cmd-2", Stdout: "output 2", ExitCode: 0},
		"forge cmd-3": {Command: "forge cmd-3", Stdout: "different output", ExitCode: 0},
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if summary.Total != 3 {
		t.Errorf("expected Total=3, got %d", summary.Total)
	}
	if summary.Broken != 1 {
		t.Errorf("expected Broken=1, got %d", summary.Broken)
	}
	if summary.OK != 2 {
		t.Errorf("expected OK=2, got %d", summary.OK)
	}
}

// --- Test: Round-trip Render -> Parse -> Verify ---

func TestVerify_RoundTrip_RenderParseVerify(t *testing.T) {
	dir := t.TempDir()

	c := Contract{
		Journey: "round-trip",
		Step:    1,
		Action:  "forge round-trip",
		Outcomes: []Outcome{
			{
				Name:          "success",
				Preconditions: "feature exists",
				Input:         "none",
				Output:        `stdout contains "round trip success", exit code 0`,
				State:         "completed",
			},
		},
		Invariants: []string{"feature_slug consistent"},
	}

	// Render to markdown
	md := RenderContract(c)

	// Write to file
	contractDir := filepath.Join(dir, "tests", "round-trip", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contractDir, "step-1-round-trip.md"), []byte(md), 0644); err != nil {
		t.Fatal(err)
	}

	// Verify with matching facts
	facts := FactTable{
		"forge round-trip": {
			Command:  "forge round-trip",
			Stdout:   "round trip success",
			ExitCode: 0,
		},
	}

	summary, err := Verify(dir, stubFactCollector{facts: facts})
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}

	if summary.Broken != 0 {
		t.Errorf("round-trip render->parse->verify should produce 0 breaks, got %d", summary.Broken)
		for _, r := range summary.Results {
			if !r.OK {
				t.Errorf("  Broken: %s breaks=%+v", r.ContractPath, r.Breaks)
			}
		}
	}
}

// --- Test: ParseContractsFile ---

func TestParseContractsFile_ValidFile(t *testing.T) {
	dir := t.TempDir()
	content := `---
journey: "test"
step: 1
step-action: "forge test"
---

# Contract: test / Step 1: forge test

## Outcome "success"
- Preconditions: "feature exists"
- Input: none
- Output: stdout "ok", exit code 0
- State: status changed

## Journey Invariants
- invariant 1
`
	filePath := filepath.Join(dir, "step-1-test.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	contracts, err := ParseContractsFile(filePath)
	if err != nil {
		t.Fatalf("ParseContractsFile failed: %v", err)
	}
	if len(contracts) != 1 {
		t.Fatalf("expected 1 contract, got %d", len(contracts))
	}
	if contracts[0].Journey != "test" {
		t.Errorf("expected journey 'test', got %q", contracts[0].Journey)
	}
}

func TestParseContractsFile_NonexistentFile(t *testing.T) {
	_, err := ParseContractsFile("/nonexistent/path.md")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

// --- Test: FormatReport with many OK contracts uses brevity ---

func TestVerifySummary_FormatReport_ManyOK_Brevity(t *testing.T) {
	var results []VerifyResult
	for i := 0; i < 10; i++ {
		results = append(results, VerifyResult{
			ContractPath: fmt.Sprintf("tests/j/_contracts/step-%d.md", i+1),
			OK:           true,
		})
	}

	summary := VerifySummary{
		Total:   10,
		Broken:  0,
		OK:      10,
		Results: results,
	}

	report := summary.FormatReport()
	if !strings.Contains(report, "unchanged contracts omitted for brevity") {
		t.Errorf("expected brevity message for >5 OK contracts, got: %s", report)
	}
}
