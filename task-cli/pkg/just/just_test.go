package just

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeJustfile creates a justfile in dir with the given content.
func writeJustfile(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "justfile"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestExtractConciseError(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		got := ExtractConciseError("", 10)
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("fewer than maxLines", func(t *testing.T) {
		input := "line1\nline2\nline3"
		got := ExtractConciseError(input, 10)
		if got != input {
			t.Errorf("expected unchanged, got %q", got)
		}
	})

	t.Run("exactly maxLines non-empty", func(t *testing.T) {
		input := "a\nb\nc"
		got := ExtractConciseError(input, 3)
		if got != input {
			t.Errorf("expected unchanged, got %q", got)
		}
	})

	t.Run("more than maxLines", func(t *testing.T) {
		input := "line1\nline2\nline3\nline4\nline5"
		got := ExtractConciseError(input, 2)
		expected := "...\nline4\nline5"
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

	t.Run("skips empty lines", func(t *testing.T) {
		input := "a\n\nb\n\n\nc\nd"
		got := ExtractConciseError(input, 2)
		expected := "...\nc\nd"
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})
}

func TestResolveScope(t *testing.T) {
	t.Run("empty scope returns empty", func(t *testing.T) {
		got := ResolveScope("/nonexistent", "")
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("all scope returns empty", func(t *testing.T) {
		got := ResolveScope("/nonexistent", "all")
		if got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})

	t.Run("nonexistent dir returns empty", func(t *testing.T) {
		got := ResolveScope("/nonexistent/path/12345", "frontend")
		if got != "" {
			t.Errorf("expected empty for nonexistent dir, got %q", got)
		}
	})

	t.Run("no project-type recipe returns empty", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo hi\n")
		got := ResolveScope(dir, "frontend")
		if got != "" {
			t.Errorf("expected empty without project-type recipe, got %q", got)
		}
	})

	t.Run("mixed project-type returns scope", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "project-type:\n  @echo mixed\n")
		got := ResolveScope(dir, "frontend")
		if got != "frontend" {
			t.Errorf("expected frontend, got %q", got)
		}
	})

	t.Run("backend project-type returns empty", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "project-type:\n  @echo backend\n")
		got := ResolveScope(dir, "frontend")
		if got != "" {
			t.Errorf("expected empty for backend project, got %q", got)
		}
	})

	t.Run("frontend project-type returns empty", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "project-type:\n  @echo frontend\n")
		got := ResolveScope(dir, "backend")
		if got != "" {
			t.Errorf("expected empty for frontend project, got %q", got)
		}
	})

	t.Run("unknown project-type returns empty with warning", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "project-type:\n  @echo unknown\n")
		got := ResolveScope(dir, "frontend")
		if got != "" {
			t.Errorf("expected empty for unknown project type, got %q", got)
		}
	})
}

func TestDefaultGateSequence(t *testing.T) {
	steps := DefaultGateSequence()
	if len(steps) != 4 {
		t.Fatalf("expected 4 steps, got %d", len(steps))
	}

	names := []string{"compile", "fmt", "lint", "test"}
	optional := []bool{false, true, true, false}
	blocking := []bool{true, false, true, true}

	for i, step := range steps {
		if step.Name != names[i] {
			t.Errorf("step %d: expected name %q, got %q", i, names[i], step.Name)
		}
		if step.Optional != optional[i] {
			t.Errorf("step %d: expected optional %v, got %v", i, optional[i], step.Optional)
		}
		if step.Blocking != blocking[i] {
			t.Errorf("step %d: expected blocking %v, got %v", i, blocking[i], step.Blocking)
		}
	}
}

func TestLintGateSequence(t *testing.T) {
	steps := LintGateSequence()
	if len(steps) != 3 {
		t.Fatalf("expected 3 steps, got %d", len(steps))
	}

	names := []string{"compile", "fmt", "lint"}
	optional := []bool{false, true, true}
	blocking := []bool{true, false, true}

	for i, step := range steps {
		if step.Name != names[i] {
			t.Errorf("step %d: expected name %q, got %q", i, names[i], step.Name)
		}
		if step.Optional != optional[i] {
			t.Errorf("step %d: expected optional %v, got %v", i, optional[i], step.Optional)
		}
		if step.Blocking != blocking[i] {
			t.Errorf("step %d: expected blocking %v, got %v", i, blocking[i], step.Blocking)
		}
	}
}

func TestRunCapture(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		output, success := RunCapture(t.TempDir(), "echo", "hello")
		if !success {
			t.Error("RunCapture() success = false, want true")
		}
		if !strings.Contains(output, "hello") {
			t.Errorf("RunCapture() output = %q, want contain hello", output)
		}
	})

	t.Run("failure", func(t *testing.T) {
		_, success := RunCapture(t.TempDir(), "false")
		if success {
			t.Error("RunCapture() success = true, want false for failing command")
		}
	})
}

func TestFileExists(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		dir := t.TempDir()
		p := dir + "/exists.txt"
		os.WriteFile(p, []byte("x"), 0644)
		if !FileExists(p) {
			t.Error("expected true for existing file")
		}
	})

	t.Run("non-existing file", func(t *testing.T) {
		if FileExists(t.TempDir() + "/nope.txt") {
			t.Error("expected false for non-existing file")
		}
	})
}

func TestHasJustfile(t *testing.T) {
	t.Run("no justfile returns false", func(t *testing.T) {
		if HasJustfile(t.TempDir()) {
			t.Error("expected false without justfile")
		}
	})

	t.Run("lowercase justfile detected", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "foo:\n  echo hi\n")
		if !HasJustfile(dir) {
			t.Error("expected true with lowercase justfile")
		}
	})

	t.Run("uppercase Justfile detected", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "Justfile"), []byte("foo:\n  echo hi\n"), 0644); err != nil {
			t.Fatal(err)
		}
		if !HasJustfile(dir) {
			t.Error("expected true with uppercase Justfile")
		}
	})
}

func TestHasRecipe(t *testing.T) {
	t.Run("existing recipe returns true", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo hi\n")
		if !HasRecipe(dir, "compile") {
			t.Error("expected true for existing recipe")
		}
	})

	t.Run("missing recipe returns false", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo hi\n")
		if HasRecipe(dir, "nonexistent") {
			t.Error("expected false for missing recipe")
		}
	})

	t.Run("no justfile returns false", func(t *testing.T) {
		if HasRecipe(t.TempDir(), "compile") {
			t.Error("expected false without justfile")
		}
	})
}

func TestRunGate(t *testing.T) {
	t.Run("no justfile returns true", func(t *testing.T) {
		passed := RunGate(t.TempDir(), "", DefaultGateSequence(), nil)
		if !passed {
			t.Error("expected true without justfile")
		}
	})

	t.Run("all recipes pass returns true", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo compile-ok\ntest:\n  echo test-ok\n")
		// Only test required recipes (compile, test) — fmt and lint are optional
		steps := []GateRecipe{
			{Name: "compile", Optional: false, Blocking: true},
			{Name: "test", Optional: false, Blocking: true},
		}
		passed := RunGate(dir, "", steps, nil)
		if !passed {
			t.Error("expected true when all required recipes pass")
		}
	})

	t.Run("optional missing recipe skipped", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo ok\n")
		steps := []GateRecipe{
			{Name: "compile", Optional: false, Blocking: true},
			{Name: "lint", Optional: true, Blocking: true},
		}
		passed := RunGate(dir, "", steps, nil)
		if !passed {
			t.Error("expected true when optional recipe is missing")
		}
	})

	t.Run("blocking failure calls onFail and returns false", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo compile-ok\nfail:\n  exit 1\n")
		var failStep, failOutput string
		steps := []GateRecipe{
			{Name: "compile", Optional: false, Blocking: true},
			{Name: "fail", Optional: false, Blocking: true},
		}
		passed := RunGate(dir, "", steps, func(step, output string) {
			failStep = step
			failOutput = output
		})
		if passed {
			t.Error("expected false on blocking failure")
		}
		if failStep != "fail" {
			t.Errorf("expected onFail called with step 'fail', got %q", failStep)
		}
		if failOutput == "" {
			t.Error("expected non-empty output from failed step")
		}
	})

	t.Run("non-blocking failure continues", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo ok\nwarn:\n  exit 1\n")
		steps := []GateRecipe{
			{Name: "compile", Optional: false, Blocking: true},
			{Name: "warn", Optional: false, Blocking: false},
		}
		passed := RunGate(dir, "", steps, func(step, output string) {
			t.Errorf("onFail should not be called for non-blocking step, got called with %q", step)
		})
		if !passed {
			t.Error("expected true when non-blocking step fails")
		}
	})

	t.Run("required missing recipe prints warning and skips", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "compile:\n  echo ok\n")
		steps := []GateRecipe{
			{Name: "compile", Optional: false, Blocking: true},
			{Name: "test", Optional: false, Blocking: true},
		}
		passed := RunGate(dir, "", steps, nil)
		if !passed {
			t.Error("expected true when required recipe is missing (graceful skip)")
		}
	})

	t.Run("scope passed for mixed project", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "project-type:\n  @echo mixed\ncompile frontend:\n  echo ok\n")
		steps := []GateRecipe{
			{Name: "compile", Optional: false, Blocking: true},
		}
		passed := RunGate(dir, "frontend", steps, nil)
		if !passed {
			t.Error("expected true with scope resolution for mixed project")
		}
	})

	t.Run("nil onFail with blocking failure does not panic", func(t *testing.T) {
		dir := t.TempDir()
		writeJustfile(t, dir, "fail:\n  exit 1\n")
		steps := []GateRecipe{
			{Name: "fail", Optional: false, Blocking: true},
		}
		passed := RunGate(dir, "", steps, nil)
		if passed {
			t.Error("expected false on blocking failure")
		}
	})
}
