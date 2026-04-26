package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseTestFailures(t *testing.T) {
	tests := []struct {
		name       string
		output     string
		wantCount  int
		wantNames  []string
	}{
		{
			name:      "empty output",
			output:    "",
			wantCount: 0,
		},
		{
			name:      "no failures",
			output:    "ok  task-cli/internal/cmd\nPASS",
			wantCount: 0,
		},
		{
			name: "go test failure",
			output: `--- FAIL: TestClaim (0.00s)
    claim_test.go:42: expected pending, got completed
FAIL`,
			wantCount: 1,
			wantNames: []string{"TestClaim"},
		},
		{
			name: "npm checkmark failure",
			output: `✗ should login successfully
  AssertionError: expected 200
✗ should persist data
  Error: connection refused`,
			wantCount: 2,
			wantNames: []string{"should login successfully", "should persist data"},
		},
		{
			name: "pytest failure",
			output: `FAILED tests/test_api.py::test_endpoint
FAILED tests/test_ui.py::test_login`,
			wantCount: 2,
			wantNames: []string{"test_endpoint", "test_login"},
		},
		{
			name: "deduplication within same pattern",
			output: `✗ should login
✗ should login`,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			failures := parseTestFailures(tt.output)
			if len(failures) != tt.wantCount {
				t.Errorf("got %d failures, want %d", len(failures), tt.wantCount)
			}
			for _, want := range tt.wantNames {
				found := false
				for _, f := range failures {
					if f.TestName == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("missing failure %q in %v", want, failures)
				}
			}
		})
	}
}

func TestExtractErrorContext(t *testing.T) {
	output := `line 1
--- FAIL: TestFoo (0.00s)
Error: expected true, got false
  at TestFoo (foo_test.go:10)
  at runtime (runtime.go:100)
line 6`

	errMsg, stackTrace := extractErrorContext(output, 2)
	if !strings.Contains(errMsg, "Error:") {
		t.Errorf("expected error message, got: %q", errMsg)
	}
	if !strings.Contains(stackTrace, "at ") {
		t.Errorf("expected stack trace, got: %q", stackTrace)
	}
}

func TestExtractErrorContext_Empty(t *testing.T) {
	errMsg, stackTrace := extractErrorContext("hello", 100)
	if errMsg != "" || stackTrace != "" {
		t.Errorf("expected empty for out-of-range, got: %q, %q", errMsg, stackTrace)
	}
}

func TestExtractRelevantOutput(t *testing.T) {
	lines := []string{}
	for i := 0; i < 20; i++ {
		lines = append(lines, "line")
	}
	output := strings.Join(lines, "\n")

	result := extractRelevantOutput(output, 10, 5)
	if result == "" {
		t.Error("expected non-empty output")
	}
}

func TestExtractRelevantOutput_Boundary(t *testing.T) {
	result := extractRelevantOutput("a\nb\nc", 0, 10)
	if !strings.Contains(result, "a") {
		t.Errorf("expected content, got: %q", result)
	}
}

func TestSanitizeTestName(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Login With Valid Credentials", "login-with-valid-credentials"},
		{"test/UI/login-flow", "test-ui-login-flow"},
		{"Special@#$Characters!", "specialcharacters"},
		{"  spaces  ", "--spaces--"},
	}

	for _, tt := range tests {
		got := sanitizeTestName(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeTestName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestMatchTestCaseID(t *testing.T) {
	t.Run("file not found falls back to sanitize", func(t *testing.T) {
		got := matchTestCaseID("Login Test", "/nonexistent/path/test-cases.md")
		want := "login-test"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("matches from test-cases.md", func(t *testing.T) {
		dir := t.TempDir()
		tcPath := filepath.Join(dir, "test-cases.md")
		// Parser reads **Test ID** first, then ## TC- title to associate
		content := `- **Test ID**: ui/login/login-test
## TC-001: Login Test
- **Test ID**: ui/logout/logout-test
## TC-002: Logout Test
`
		os.WriteFile(tcPath, []byte(content), 0644)

		got := matchTestCaseID("Login Test", tcPath)
		if got != "ui/login/login-test" {
			t.Errorf("got %q, want ui/login/login-test", got)
		}

		got = matchTestCaseID("logout test", tcPath)
		if got != "ui/logout/logout-test" {
			t.Errorf("got %q for case-insensitive match, want ui/logout/logout-test", got)
		}

		got = matchTestCaseID("Unknown Test", tcPath)
		if got != "unknown-test" {
			t.Errorf("got %q for unknown test, want unknown-test", got)
		}
	})
}

func TestWriteFailureFiles(t *testing.T) {
	t.Run("no failures returns nil", func(t *testing.T) {
		err := writeFailureFiles("/tmp", "slug", nil)
		if err != nil {
			t.Errorf("expected nil for empty failures, got: %v", err)
		}
	})

	t.Run("writes failure files", func(t *testing.T) {
		dir := t.TempDir()
		failures := []TestFailure{
			{
				TestName:     "TestLogin",
				TestCaseID:   "ui-login",
				File:         "login_test.go",
				Line:         42,
				ErrorMessage: "expected 200",
				Output:       "some output",
				StackTrace:   "at TestLogin (login_test.go:42)",
			},
		}

		err := writeFailureFiles(dir, "test-slug", failures)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		failuresDir := filepath.Join(dir, "docs", "features", "test-slug", "testing", "results", "failures")
		failurePath := filepath.Join(failuresDir, "failure-ui-login.md")
		data, err := os.ReadFile(failurePath)
		if err != nil {
			t.Fatalf("failure file not created: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "TestLogin") {
			t.Errorf("failure file should contain test name, got: %s", content)
		}
		if !strings.Contains(content, "login_test.go:42") {
			t.Errorf("failure file should contain file:line, got: %s", content)
		}
	})

	t.Run("failure without file", func(t *testing.T) {
		dir := t.TempDir()
		failures := []TestFailure{
			{
				TestName:   "TestBasic",
				TestCaseID: "basic",
			},
		}

		err := writeFailureFiles(dir, "slug", failures)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify file was created without file/line info
		failuresDir := filepath.Join(dir, "docs", "features", "slug", "testing", "results", "failures")
		failurePath := filepath.Join(failuresDir, "failure-basic.md")
		data, _ := os.ReadFile(failurePath)
		if strings.Contains(string(data), "File:") {
			t.Error("should not contain File header when empty")
		}
	})
}
