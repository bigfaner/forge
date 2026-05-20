package testrunner

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func captureStdout(f func()) string {
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	f()
	_ = w.Close()
	os.Stdout = old
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func captureStderr(f func()) string {
	var buf bytes.Buffer
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w
	f()
	_ = w.Close()
	os.Stderr = old
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestCapitalize(t *testing.T) {
	tests := []struct{ input, want string }{
		{"compile", "Compile"},
		{"fmt", "Fmt"},
		{"", ""},
		{"a", "A"},
		{"test", "Test"},
	}
	for _, tc := range tests {
		if got := Capitalize(tc.input); got != tc.want {
			t.Errorf("Capitalize(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestPrintHookJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		out := captureStdout(func() {
			PrintHookJSON(map[string]any{"decision": "block", "reason": "test"})
		})
		if !strings.Contains(out, `"decision"`) {
			t.Errorf("expected JSON output, got: %s", out)
		}
	})

	t.Run("marshal error", func(t *testing.T) {
		out := captureStderr(func() {
			PrintHookJSON(map[string]any{"ch": make(chan int)})
		})
		if !strings.Contains(out, "WARNING") {
			t.Errorf("expected warning for marshal error, got: %s", out)
		}
	})
}

func TestHasNpmTestScript(t *testing.T) {
	t.Run("has test script", func(t *testing.T) {
		dir := t.TempDir()
		pkg := `{"scripts": {"test": "jest"}}`
		_ = os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644)
		if !hasNpmTestScript(dir) {
			t.Error("expected true for package with test script")
		}
	})

	t.Run("no test script", func(t *testing.T) {
		dir := t.TempDir()
		pkg := `{"scripts": {"build": "tsc"}}`
		_ = os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644)
		if hasNpmTestScript(dir) {
			t.Error("expected false for package without test script")
		}
	})

	t.Run("no package.json", func(t *testing.T) {
		dir := t.TempDir()
		if hasNpmTestScript(dir) {
			t.Error("expected false when no package.json")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		_ = os.WriteFile(filepath.Join(dir, "package.json"), []byte("not json"), 0644)
		if hasNpmTestScript(dir) {
			t.Error("expected false for invalid JSON")
		}
	})
}

func TestRunProjectTests(t *testing.T) {
	t.Run("no test framework returns true with empty output", func(t *testing.T) {
		dir := t.TempDir()
		output, ok := RunProjectTests(dir)
		if !ok {
			t.Error("expected ok=true for no test command")
		}
		if output != "" {
			t.Errorf("expected empty output, got %q", output)
		}
	})

	t.Run("no test framework prints warning to stdout", func(t *testing.T) {
		dir := t.TempDir()
		out := captureStdout(func() {
			RunProjectTests(dir)
		})
		if !strings.Contains(out, "WARNING") {
			t.Errorf("expected warning when no test command found, got: %s", out)
		}
	})

	t.Run("go.mod uses go test", func(t *testing.T) {
		dir := t.TempDir()
		_ = os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test\n\ngo 1.21\n"), 0644)
		_, ok := RunProjectTests(dir)
		_ = ok // may succeed or fail, just verify no panic
	})

	t.Run("npm test branch", func(t *testing.T) {
		dir := t.TempDir()
		pkg := `{"scripts": {"test": "echo npm-test-pass"}}`
		_ = os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0644)
		out := captureStderr(func() {
			RunProjectTests(dir)
		})
		if !strings.Contains(out, "npm-test-pass") {
			t.Errorf("expected npm test output, got: %s", out)
		}
	})

	t.Run("pytest branch", func(t *testing.T) {
		dir := t.TempDir()
		_ = os.WriteFile(filepath.Join(dir, "pytest.ini"), []byte("[pytest]\n"), 0644)
		out := captureStderr(func() {
			RunProjectTests(dir)
		})
		_ = out // just verify no panic
	})

	t.Run("justfile branch", func(t *testing.T) {
		dir := t.TempDir()
		_ = os.WriteFile(filepath.Join(dir, "justfile"), []byte("test:\n    echo just-test-output\n"), 0644)
		out := captureStderr(func() {
			RunProjectTests(dir)
		})
		_ = out // just verify no panic
	})

	t.Run("Makefile branch", func(t *testing.T) {
		if _, err := exec.LookPath("make"); err != nil {
			t.Skip("make not installed")
		}
		dir := t.TempDir()
		_ = os.WriteFile(filepath.Join(dir, "Makefile"), []byte("test:\n\t@echo make-test-output\n"), 0644)
		out := captureStderr(func() {
			RunProjectTests(dir)
		})
		if !strings.Contains(out, "make-test-output") {
			t.Errorf("expected make test output, got: %s", out)
		}
	})
}
