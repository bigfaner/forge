---
title: "Go Testing Convention"
---

# Go Testing Convention

Convention for generating Go test code using the standard `testing` package with `testify/assert`.

## framework

- **name**: Go testing package + testify/assert
- **version**: go1.18+
- **language**: Go
- **runner_command**: `go test -v -json -tags=e2e`

## discovery

- **test_dir**: `forge-cli/tests/` (forge-cli e2e tests), `tests/` (project-level e2e tests per feature)
- **file_pattern**: `*_test.go`
- **exclude_pattern**: `vendor/`, `node_modules/`

## structure

- **suite_pattern**: `func TestXxx(t *testing.T)` — top-level test functions with `Test` prefix
- **case_pattern**: `TestTC_NNN_Description` — zero-padded test case number with description
- **hook_pattern**: `TestMain(m *testing.M)` — suite-level setup/teardown in `main_test.go`

### Test Function Naming

Pattern: `TestTC_NNN_Description` where `NNN` is zero-padded test case number.

```go
func TestTC_001_LoginWithValidCredentials(t *testing.T) { ... }
```

### Table-Driven Tests

When multiple test cases share the same logic:

```go
func TestTC_NNN_Descriptions(t *testing.T) {
    cases := []struct {
        name     string
        input    string
        expected string
    }{
        {"valid input", "hello", "HELLO"},
        {"empty input", "", ""},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            result := strings.ToUpper(tc.input)
            assert.Equal(t, tc.expected, result)
        })
    }
}
```

### CLI Testing

Use `os/exec` to invoke CLI binaries:

```go
func TestTC_NNN_CliCommand(t *testing.T) {
    cmd := exec.Command("forge", "subcommand", "--flag", "value")
    cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_DIR="+dir)
    out, err := cmd.CombinedOutput()
    assert.NoError(t, err)
    assert.Contains(t, string(out), "expected output")
}
```

### API Testing

Use `net/http` Client for integration testing:

```go
func TestTC_NNN_ApiEndpoint(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:8080/api/resource", nil)
    resp, err := http.DefaultClient.Do(req)
    assert.NoError(t, err)
    defer resp.Body.Close()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### Traceability

Each test function must include a traceability comment:

```go
// Traceability: TC-NNN -> {PRD Source}
```

### Shared Infrastructure

Each test directory (e.g., `forge-cli/tests/<suite>/`, `tests/<suite>/`) provides its own `main_test.go` with `TestMain` for suite-level setup/teardown. This file is created once and must NOT be overwritten during feature generation.

## assertions

- **style**: assert
- **library**: `github.com/stretchr/testify/assert` (NOT `require`)
- **custom_matchers**: none

### Key Functions

- `assert.NoError(t, err, msg...)` — verify no error
- `assert.Equal(t, expected, actual, msg...)` — equality check
- `assert.Contains(t, str, substr, msg...)` — substring check
- `assert.True(t, condition, msg...)` — boolean assertion
- `assert.False(t, condition, msg...)` — boolean assertion
- `assert.Empty(t, collection, msg...)` — empty collection
- `assert.NotNil(t, obj, msg...)` — not nil check
- `assert.Error(t, err, msg...)` — error expected

**Rule**: Always use `assert`, never `require`. `require` stops execution immediately; `assert` collects all failures.

## Tags

- **Build tag**: `//go:build e2e` must be the first line of every test file (before package declaration)
- **Feature tag**: `//go:build feature` used for promoted tests
- **Format**: Pure Go build tag syntax, no comments after the tag

```go
//go:build e2e

package e2e
```

## Result Format

- **Output flags**: `-json -v`
- **Format type**: `json-stream` (one JSON object per line)

### JSON Stream Fields

| Field | Meaning |
|-------|---------|
| `Time` | Timestamp |
| `Action` | `run`, `output`, `pass`, `fail`, `skip` |
| `Package` | Full Go module path |
| `Test` | Test function name |
| `Elapsed` | Duration in seconds |
| `Output` | Captured stdout/stderr line |

### TC ID Extraction

Pattern: `TestTC_NNN_Description` -> `TC-NNN`
Regex: `TC_(\d+)` with separator normalization.

## Import Patterns

Standard imports for e2e tests:

```go
import (
    "os/exec"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)
```

- HTTP tests add: `"net/http"`, `"net/http/httptest"`
- File tests add: `"path/filepath"`
- All test files import `"testing"`

## Anti-patterns (Forbidden)

| Forbidden | Replacement |
|-----------|-------------|
| `time.Sleep` for synchronization | Retry loops with `assert.Eventually` or custom `withRetry()` |
| `require` assertions | `assert` assertions |
| Hardcoded ports | `net.Listen("tcp", ":0")` or environment variables |
| Real secrets/tokens in code | `os.Getenv("E2E_API_TOKEN")` |
| `//go:build` tag missing | Always include `//go:build e2e` as first line |
| Unconditional `t.Skip` | Implement properly or don't generate |
| Recursive test invocation | Recursion guard via env var or `-run` flag |

## Helpers

### runCLI helper

```go
func runCLI(t *testing.T, args ...string) string {
    t.Helper()
    cmd := exec.Command("forge", args...)
    out, err := cmd.CombinedOutput()
    assert.NoError(t, err, "forge %s failed: %s", strings.Join(args, " "), string(out))
    return string(out)
}
```

### withRetry helper

```go
func withRetry(t *testing.T, fn func() error, maxAttempts int, interval time.Duration) {
    t.Helper()
    for i := 0; i < maxAttempts; i++ {
        if err := fn(); err == nil {
            return
        }
        time.Sleep(interval)
    }
    t.Fatal("retry exhausted")
}
```

## VERIFY Markers

Resolve all `// VERIFY:` comments using Fact Table values. Post-generation check:

```bash
just unit-test
grep -rn '// VERIFY:' forge-cli/tests/ tests/ --include='*_test.go'
```
