package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/contract"

	"github.com/stretchr/testify/assert"
)

// --- Test: forge test verify command registered ---

func TestTestVerify_CommandRegistered(t *testing.T) {
	found := false
	for _, cmd := range testCmd.Commands() {
		if cmd.Name() == "verify" {
			found = true
			break
		}
	}
	if !found {
		t.Error("test group missing 'verify' subcommand")
	}
}

// --- Test: forge test verify with no contracts returns unverifiable (Exit via Total==0) ---

func TestTestVerify_NoContracts_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create minimal project structure
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Verify() returns empty summary without error; runTestVerify handles Total==0 via Exit()
	collector := contract.RealFactCollector{}
	summary, err := contract.Verify(dir, collector)
	assert.NoError(t, err)
	assert.Equal(t, 0, summary.Total)
}

// --- Test: forge test verify with contracts produces report ---

func TestTestVerify_WithContracts_ProducesReport(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	// Create minimal project structure
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a contract file
	contractDir := filepath.Join(dir, "tests", "task-lifecycle", "_contracts")
	if err := os.MkdirAll(contractDir, 0755); err != nil {
		t.Fatal(err)
	}
	contractContent := `---
journey: "task-lifecycle"
step: 1
step-action: "forge feature create"
---

# Contract: task-lifecycle / Step 1: forge feature create

## Outcome "success"
- Preconditions: "no feature with this slug exists"
- Input: feature-slug as positional arg
- Output: stdout contains "success confirmation", exit code 0
- State: feature directory created

## Journey Invariants
- feature_slug consistent across all steps
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-1-feature-create.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	output, err := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "verify"})
		return rootCmd.Execute()
	})
	if err != nil {
		t.Fatalf("test verify failed: %v", err)
	}

	if !strings.Contains(output, "Scanning 1 Contracts") {
		t.Errorf("expected 'Scanning 1 Contracts' in output, got: %s", output)
	}
	if !strings.Contains(output, "Summary:") {
		t.Errorf("expected 'Summary:' in output, got: %s", output)
	}
}

// --- Test: verify does not modify files (Hard Rule) ---

func TestTestVerify_DoesNotModifyContractFiles(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

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
- Output: stdout "ok", exit code 0
- State: status changed

## Journey Invariants
- test invariant
`
	contractFile := filepath.Join(contractDir, "step-1-test.md")
	if err := os.WriteFile(contractFile, []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	before, err := os.ReadFile(contractFile)
	if err != nil {
		t.Fatal(err)
	}

	_, _ = captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "verify"})
		return rootCmd.Execute()
	})

	after, err := os.ReadFile(contractFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Error("Hard Rule violation: verify modified a contract file")
	}
}

// --- Test: verify output format matches proposal ---

func TestTestVerify_OutputFormat(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("CLAUDE_PROJECT_DIR", dir)

	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

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
- Output: stdout contains "ok", exit code 0
- State: status changed

## Journey Invariants
- invariant
`
	if err := os.WriteFile(filepath.Join(contractDir, "step-1-test.md"), []byte(contractContent), 0644); err != nil {
		t.Fatal(err)
	}

	output, _ := captureOutput(func() error {
		rootCmd.SetArgs([]string{"test", "verify"})
		return rootCmd.Execute()
	})

	// Verify output has the canonical format
	if !strings.Contains(output, "Scanning") {
		t.Errorf("expected 'Scanning' in output, got: %s", output)
	}
	if !strings.Contains(output, "Contracts against Fact Table") {
		t.Errorf("expected 'Contracts against Fact Table' in output, got: %s", output)
	}
	if !strings.Contains(output, "OK (") {
		t.Errorf("expected 'OK (' in output, got: %s", output)
	}
	if !strings.Contains(output, "Summary:") {
		t.Errorf("expected 'Summary:' in output, got: %s", output)
	}
	if !strings.Contains(output, "false positives") {
		t.Errorf("expected 'false positives' in output, got: %s", output)
	}
}
