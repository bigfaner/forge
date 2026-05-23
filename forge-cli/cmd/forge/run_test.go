package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_Help(t *testing.T) {
	// Save original args and stdout
	origArgs := os.Args
	origStdout := os.Stdout
	defer func() {
		os.Args = origArgs
		os.Stdout = origStdout
	}()

	// Create temp dir with go.mod
	tmpDir := t.TempDir()
	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Capture stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Change to temp dir
	origWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(origWd) }()
	_ = os.Chdir(tmpDir)
	os.Args = []string{"forge", "--help"}

	// Run in goroutine since it might call os.Exit
	done := make(chan bool)
	go func() {
		Run()
		done <- true
	}()

	// Wait for completion
	<-done
	_ = w.Close()
	os.Stdout = origStdout

	// Read output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if output == "" {
		t.Error("expected some output from --help")
	}
}
