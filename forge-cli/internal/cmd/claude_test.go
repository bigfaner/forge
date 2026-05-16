package cmd

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestClaudeCmd_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "claude" {
			found = true
			break
		}
	}
	if !found {
		t.Error("claude command should be registered as top-level command")
	}
}

func TestClaudeCmd_DisabledFlagParsing(t *testing.T) {
	if !claudeCmd.DisableFlagParsing {
		t.Error("claudeCmd.DisableFlagParsing should be true for transparent arg passthrough")
	}
}

func TestClaudeCmd_Use(t *testing.T) {
	if !strings.Contains(claudeCmd.Use, "claude") {
		t.Errorf("claudeCmd.Use = %q, should contain 'claude'", claudeCmd.Use)
	}
}

func TestClaudeCmd_FlagPassthrough(t *testing.T) {
	// Verify that --model flag does NOT get consumed by cobra
	// (DisableFlagParsing means cobra ignores all flags)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"claude", "--model", "opus", "-p", "hello"})

	// This will fail because claude binary doesn't exist, but should NOT
	// fail due to "unknown flag" -- that proves flag passthrough works.
	_ = rootCmd.Execute()
	output := buf.String()
	if strings.Contains(output, "unknown flag") || strings.Contains(output, "unknown shorthand") {
		t.Errorf("flags should pass through, but got: %s", output)
	}
}

func TestClaudeCmd_ErrorWhenBinaryNotFound(t *testing.T) {
	// Override lookPathFunc to simulate claude not in PATH
	orig := lookPathFunc
	lookPathFunc = func(_ string) (string, error) {
		return "", &exec.Error{Name: "claude", Err: exec.ErrNotFound}
	}
	defer func() { lookPathFunc = orig }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"claude"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when claude binary not found")
	}
	stderr := buf.String()
	if !strings.Contains(stderr, "claude") {
		t.Errorf("error should mention 'claude', got: %s", stderr)
	}
}

func TestClaudeCmd_PrependsDangerouslySkipPermissions(t *testing.T) {
	// Capture the args passed to runClaudeFunc
	var capturedArgs []string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	// Path validation must succeed
	origLookPath := lookPathFunc
	lookPathFunc = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	defer func() { lookPathFunc = origLookPath }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"claude", "-c"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(capturedArgs) < 2 {
		t.Fatalf("expected at least 2 args, got %d: %v", len(capturedArgs), capturedArgs)
	}
	if capturedArgs[0] != "--dangerously-skip-permissions" {
		t.Errorf("first arg should be --dangerously-skip-permissions, got %q", capturedArgs[0])
	}
	if capturedArgs[1] != "-c" {
		t.Errorf("second arg should be -c, got %q", capturedArgs[1])
	}
}

func TestClaudeCmd_PrependsDangerouslySkipPermissions_WithMultipleFlags(t *testing.T) {
	var capturedArgs []string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	origLookPath := lookPathFunc
	lookPathFunc = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	defer func() { lookPathFunc = origLookPath }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"claude", "--model", "opus", "-p", "hello world"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"--dangerously-skip-permissions", "--model", "opus", "-p", "hello world"}
	if len(capturedArgs) != len(expected) {
		t.Fatalf("expected %d args, got %d: %v", len(expected), len(capturedArgs), capturedArgs)
	}
	for i, want := range expected {
		if capturedArgs[i] != want {
			t.Errorf("arg[%d] = %q, want %q", i, capturedArgs[i], want)
		}
	}
}

func TestClaudeCmd_EmptyArgs(t *testing.T) {
	var capturedArgs []string
	origRunClaude := runClaudeFunc
	runClaudeFunc = func(args []string) error {
		capturedArgs = args
		return nil
	}
	defer func() { runClaudeFunc = origRunClaude }()

	origLookPath := lookPathFunc
	lookPathFunc = func(_ string) (string, error) {
		return "/usr/bin/claude", nil
	}
	defer func() { lookPathFunc = origLookPath }()

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"claude"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(capturedArgs) != 1 || capturedArgs[0] != "--dangerously-skip-permissions" {
		t.Errorf("expected only [--dangerously-skip-permissions], got %v", capturedArgs)
	}
}

func TestClaudeCmd_DefaultsAreNotNil(t *testing.T) {
	if lookPathFunc == nil {
		t.Error("lookPathFunc should not be nil")
	}
	if runClaudeFunc == nil {
		t.Error("runClaudeFunc should not be nil")
	}
}
