package e2e

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"forge-cli/pkg/feature"
)

// setupProfile creates a temp directory with a valid profile config.
func setupProfile(t *testing.T, profileName string) string {
	t.Helper()
	dir := t.TempDir()
	forgeDir := filepath.Join(dir, feature.ForgeDir)
	if err := os.MkdirAll(forgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configContent := fmt.Sprintf("test-profiles:\n  - %s\n", profileName)
	if err := os.WriteFile(filepath.Join(forgeDir, feature.ForgeConfigFileName), []byte(configContent), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

// setupProfileWithE2E creates a temp dir with a valid profile and e2e test directory.
// Used by TestVerify which needs the directory structure for file scanning.
func setupProfileWithE2E(t *testing.T, profileName string) string {
	t.Helper()
	dir := setupProfile(t, profileName)
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
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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

	t.Run("no profile returns ErrNoProfile", func(t *testing.T) {
		dir := t.TempDir()
		err := Run(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error")
		}
		if !errors.Is(err, ErrNoProfile) {
			t.Fatalf("expected ErrNoProfile, got %v", err)
		}
	})
}

func TestSetup(t *testing.T) {
	t.Run("delegates to just e2e-setup", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "web-playwright")
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

	t.Run("no profile returns ErrNoProfile", func(t *testing.T) {
		dir := t.TempDir()
		err := Setup(RunOpts{ProjectRoot: dir})
		if !errors.Is(err, ErrNoProfile) {
			t.Fatalf("expected ErrNoProfile, got %v", err)
		}
	})
}

func TestVerify(t *testing.T) {
	t.Run("go-test profile scans for VERIFY markers", func(t *testing.T) {
		dir := setupProfileWithE2E(t, "go-test")
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
		dir := setupProfileWithE2E(t, "go-test")
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
		dir := setupProfileWithE2E(t, "go-test")

		err := Verify(RunOpts{ProjectRoot: dir, Feature: "nonexistent"})
		if !errors.Is(err, ErrFeatureNotFound) {
			t.Fatalf("expected ErrFeatureNotFound, got %v", err)
		}
	})

	t.Run("no profile returns ErrNoProfile", func(t *testing.T) {
		dir := t.TempDir()
		err := Verify(RunOpts{ProjectRoot: dir})
		if !errors.Is(err, ErrNoProfile) {
			t.Fatalf("expected ErrNoProfile, got %v", err)
		}
	})
}

func TestCompile(t *testing.T) {
	t.Run("delegates to just e2e-compile", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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

	t.Run("no profile returns ErrNoProfile", func(t *testing.T) {
		dir := t.TempDir()
		err := Compile(dir)
		if !errors.Is(err, ErrNoProfile) {
			t.Fatalf("expected ErrNoProfile, got %v", err)
		}
	})
}

func TestDiscover(t *testing.T) {
	t.Run("delegates to just e2e-discover", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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
		dir := setupProfile(t, "go-test")
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

	t.Run("no profile returns ErrNoProfile", func(t *testing.T) {
		dir := t.TempDir()
		err := Discover(dir)
		if !errors.Is(err, ErrNoProfile) {
			t.Fatalf("expected ErrNoProfile, got %v", err)
		}
	})
}
