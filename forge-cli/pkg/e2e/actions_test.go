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

// setupGoTestProfile creates a temp dir with go-test profile and e2e test directory.
func setupGoTestProfile(t *testing.T) string {
	t.Helper()
	dir := setupProfile(t, "go-test")
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
	t.Run("go-test profile dispatches go test", func(t *testing.T) {
		dir := setupGoTestProfile(t)
		s := &stubExec{responses: map[string]execResponse{
			"go test ./tests/e2e/...": {output: []byte("ok\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir, Feature: ""})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("go-test profile with feature flag", func(t *testing.T) {
		dir := setupGoTestProfile(t)
		s := &stubExec{responses: map[string]execResponse{
			"go test ./tests/e2e/features/my-feature/...": {output: []byte("ok\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir, Feature: "my-feature"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("web-playwright profile dispatches npx playwright test", func(t *testing.T) {
		dir := setupProfile(t, "web-playwright")
		s := &stubExec{responses: map[string]execResponse{
			"npx playwright test": {output: []byte("ok\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unsupported profile returns error", func(t *testing.T) {
		dir := setupProfile(t, "maestro")
		oldRunner := runner
		runner = &stubExec{}
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error for unsupported profile")
		}
		if !strings.Contains(err.Error(), "unsupported profile for run") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("external tool failure returns formatted error", func(t *testing.T) {
		dir := setupGoTestProfile(t)
		s := &stubExec{responses: map[string]execResponse{
			"go test ./tests/e2e/...": {output: []byte("first line of error\nsecond line"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Run(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error for tool failure")
		}
		if !strings.Contains(err.Error(), "go test failed:") {
			t.Fatalf("expected error to contain 'go test failed:', got %q", err.Error())
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
	t.Run("go-test profile dispatches go install", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
		s := &stubExec{responses: map[string]execResponse{
			"go install": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Setup(RunOpts{ProjectRoot: dir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("web-playwright profile dispatches npx playwright install", func(t *testing.T) {
		dir := setupProfile(t, "web-playwright")
		s := &stubExec{responses: map[string]execResponse{
			"npx playwright install": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Setup(RunOpts{ProjectRoot: dir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("pytest profile dispatches pip install", func(t *testing.T) {
		dir := setupProfile(t, "pytest")
		s := &stubExec{responses: map[string]execResponse{
			"python -m pip install pytest": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Setup(RunOpts{ProjectRoot: dir})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("external tool failure returns formatted error", func(t *testing.T) {
		dir := setupProfile(t, "web-playwright")
		s := &stubExec{responses: map[string]execResponse{
			"npx playwright install": {output: []byte("EACCES: permission denied\n"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Setup(RunOpts{ProjectRoot: dir})
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "npx playwright install failed:") {
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
		dir := setupGoTestProfile(t)
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
		dir := setupGoTestProfile(t)
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
		dir := setupGoTestProfile(t)

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
	t.Run("go-test profile dispatches go build", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
		s := &stubExec{responses: map[string]execResponse{
			"go build ./tests/e2e/...": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("web-playwright profile dispatches tsc --noEmit", func(t *testing.T) {
		dir := setupProfile(t, "web-playwright")
		s := &stubExec{responses: map[string]execResponse{
			"npx tsc --noEmit": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("pytest profile dispatches compileall", func(t *testing.T) {
		dir := setupProfile(t, "pytest")
		s := &stubExec{responses: map[string]execResponse{
			"python -m compileall tests/e2e/ -q": {output: []byte(""), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unsupported profile returns error", func(t *testing.T) {
		dir := setupProfile(t, "maestro")
		oldRunner := runner
		runner = &stubExec{}
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err == nil {
			t.Fatal("expected error for unsupported profile")
		}
		if !strings.Contains(err.Error(), "unsupported profile for compile") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("external tool failure returns formatted error", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
		s := &stubExec{responses: map[string]execResponse{
			"go build ./tests/e2e/...": {output: []byte("./tests/e2e/main_test.go:15: undefined: Foo\n"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Compile(dir)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "go build failed:") {
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
	t.Run("go-test profile dispatches go test -list", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
		s := &stubExec{responses: map[string]execResponse{
			"go test ./tests/e2e/... -list .* -tags=e2e": {output: []byte("TestExample\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("web-playwright profile dispatches playwright test --list", func(t *testing.T) {
		dir := setupProfile(t, "web-playwright")
		s := &stubExec{responses: map[string]execResponse{
			"npx playwright test --list": {output: []byte("test list\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("pytest profile dispatches pytest --collect-only", func(t *testing.T) {
		dir := setupProfile(t, "pytest")
		s := &stubExec{responses: map[string]execResponse{
			"python -m pytest tests/e2e/ --collect-only -q": {output: []byte("test list\n"), err: nil},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unsupported profile returns error", func(t *testing.T) {
		dir := setupProfile(t, "maestro")
		oldRunner := runner
		runner = &stubExec{}
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err == nil {
			t.Fatal("expected error for unsupported profile")
		}
		if !strings.Contains(err.Error(), "unsupported profile for discover") {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("external tool failure returns formatted error", func(t *testing.T) {
		dir := setupProfile(t, "go-test")
		s := &stubExec{responses: map[string]execResponse{
			"go test ./tests/e2e/... -list .* -tags=e2e": {output: []byte("build constraints exclude all tests\n"), err: fmt.Errorf("exit status 1")},
		}}
		oldRunner := runner
		runner = s
		defer func() { runner = oldRunner }()

		err := Discover(dir)
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "go test -list failed:") {
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
