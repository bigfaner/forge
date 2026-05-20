package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// stubExec is a hand-rolled mock matching project convention (no testify/gomock).
type stubExec struct {
	responses map[string]execResponse
}

type execResponse struct {
	output []byte
	err    error
}

func (s *stubExec) Run(name string, args ...string) ([]byte, error) {
	key := name + " " + strings.Join(args, " ")
	if r, ok := s.responses[key]; ok {
		return r.output, r.err
	}
	return nil, fmt.Errorf("stubExec: unexpected command: %s", key)
}

func TestStubExec(t *testing.T) {
	t.Run("returns configured response", func(t *testing.T) {
		s := &stubExec{responses: map[string]execResponse{
			"echo hello": {output: []byte("hello\n"), err: nil},
		}}

		out, err := s.Run("echo", "hello")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(out) != "hello\n" {
			t.Fatalf("expected 'hello\\n', got %q", string(out))
		}
	})

	t.Run("returns error for unexpected command", func(t *testing.T) {
		s := &stubExec{responses: map[string]execResponse{}}

		_, err := s.Run("unknown", "cmd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "stubExec: unexpected command") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("returns configured error", func(t *testing.T) {
		s := &stubExec{responses: map[string]execResponse{
			"fail cmd": {output: nil, err: fmt.Errorf("command failed")},
		}}

		_, err := s.Run("fail", "cmd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "command failed" {
			t.Fatalf("expected 'command failed', got %q", err.Error())
		}
	})
}

func TestRealExecImplementsExecRunner(_ *testing.T) {
	// Compile-time interface check
	var _ ExecRunner = RealExec{}
}

func TestRunnerDefault(t *testing.T) {
	// Verify runner is set to RealExec by default
	_, ok := runner.(RealExec)
	if !ok {
		t.Fatal("expected runner to be RealExec by default")
	}
}

// setupE2EDir creates a temp dir with an e2e test directory.
// Used by tests that need the directory structure for file scanning.
func setupE2EDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	e2eDir := filepath.Join(dir, "tests", "e2e")
	if err := os.MkdirAll(e2eDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// Create a minimal test file so discovery can find something
	if err := os.WriteFile(filepath.Join(e2eDir, "example_test.go"), []byte("package e2e\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestRun(t *testing.T) {
	t.Run("delegates to just test-e2e", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just test-e2e": {output: []byte("ok\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir, Feature: ""})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("passes feature as justfile argument", func(t *testing.T) {
		dir := t.TempDir()
		// Create feature directory so Run's existence check passes
		featureDir := filepath.Join(dir, "tests", "e2e", "features", "my-feature")
		if err := os.MkdirAll(featureDir, 0o755); err != nil {
			t.Fatal(err)
		}
		s := &stubExec{responses: map[string]execResponse{
			"just test-e2e feature=my-feature": {output: []byte("ok\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir, Feature: "my-feature"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("just not on PATH returns actionable error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just test-e2e": {output: nil, err: fmt.Errorf("exec: \"just\": executable file not found in $PATH")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error for just not found")
		}
		if !strings.Contains(err.Error(), "'just' is required but not found on PATH") {
			t.Fatalf("expected 'just' not found error, got %q", err.Error())
		}
	})

	t.Run("just failure returns formatted error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just test-e2e": {output: []byte("first line of error\nsecond line"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error for tool failure")
		}
		if !strings.Contains(err.Error(), "just test-e2e failed:") {
			t.Fatalf("expected error to contain 'just test-e2e failed:', got %q", err.Error())
		}
		if !strings.Contains(err.Error(), "first line of error") {
			t.Fatalf("expected error to contain first line of stderr, got %q", err.Error())
		}
	})
}

func TestSetup(t *testing.T) {
	t.Run("delegates to just e2e-setup", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-setup": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Setup(RunOpts{ProjectRoot: dir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("just not on PATH returns actionable error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-setup": {output: nil, err: fmt.Errorf("exec: \"just\": executable file not found in $PATH")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Setup(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error for just not found")
		}
		if !strings.Contains(err.Error(), "'just' is required but not found on PATH") {
			t.Fatalf("expected 'just' not found error, got %q", err.Error())
		}
	})

	t.Run("just failure returns formatted error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-setup": {output: []byte("EACCES: permission denied\n"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Setup(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "just e2e-setup failed:") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestVerify(t *testing.T) {
	t.Run("scans for VERIFY markers", func(t *testing.T) {
		dir := setupE2EDir(t)
		// Write a file without VERIFY markers
		e2eDir := filepath.Join(dir, "tests", "e2e")
		if err := os.WriteFile(filepath.Join(e2eDir, "clean_test.go"), []byte("package e2e\n// no markers\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		oldRunner := runner
		runner = &stubExec{}
		defer func() { runner = oldRunner }()

		err := Verify(RunOpts{ProjectRoot: dir, Feature: ""})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("finds VERIFY markers returns error", func(t *testing.T) {
		dir := setupE2EDir(t)
		e2eDir := filepath.Join(dir, "tests", "e2e")
		if err := os.WriteFile(filepath.Join(e2eDir, "has_verify_test.go"), []byte("// VERIFY: placeholder\npackage e2e\n"), 0o644); err != nil {
			t.Fatal(err)
		}

		oldRunner := runner
		runner = &stubExec{}
		defer func() { runner = oldRunner }()

		err := Verify(RunOpts{ProjectRoot: dir, Feature: ""})
		if err == nil {
			t.Fatal("expected error for VERIFY markers")
		}
		if !strings.Contains(err.Error(), "VERIFY markers found") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("feature not found returns ErrFeatureNotFound", func(t *testing.T) {
		dir := setupE2EDir(t)

		err := Verify(RunOpts{ProjectRoot: dir, Feature: "nonexistent"})
		if !strings.Contains(err.Error(), "feature not found") {
			t.Fatalf("expected feature not found error, got %v", err)
		}
	})
}

func TestCompile(t *testing.T) {
	t.Run("delegates to just e2e-compile", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-compile": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("just not on PATH returns actionable error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-compile": {output: nil, err: fmt.Errorf("exec: \"just\": executable file not found in $PATH")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err == nil {
			t.Fatal("expected error for just not found")
		}
		if !strings.Contains(err.Error(), "'just' is required but not found on PATH") {
			t.Fatalf("expected 'just' not found error, got %q", err.Error())
		}
	})

	t.Run("just failure returns formatted error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-compile": {output: []byte("./tests/e2e/main_test.go:15: undefined: Foo\n"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "just e2e-compile failed:") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestDiscover(t *testing.T) {
	t.Run("delegates to just e2e-discover", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-discover": {output: []byte("TestExample\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("just not on PATH returns actionable error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-discover": {output: nil, err: fmt.Errorf("exec: \"just\": executable file not found in $PATH")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err == nil {
			t.Fatal("expected error for just not found")
		}
		if !strings.Contains(err.Error(), "'just' is required but not found on PATH") {
			t.Fatalf("expected 'just' not found error, got %q", err.Error())
		}
	})

	t.Run("just failure returns formatted error", func(t *testing.T) {
		dir := t.TempDir()
		s := &stubExec{responses: map[string]execResponse{
			"just e2e-discover": {output: []byte("build constraints exclude all tests\n"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "just e2e-discover failed:") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
