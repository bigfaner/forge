---
id: "T-test-2"
title: "Generate Test Scripts (go-test)"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["T-test-1b"]
type: "test-pipeline.gen-scripts"
scope: "all"
profile: "go-test"
---

# Generate Test Scripts (go-test)

Profile: **go-test**

# Go Test Generate Strategy

Profile-specific test generation rules for the `gen-test-scripts` skill.

## Test Runner & Imports

| Test type | Runner | Assertion | HTTP | Process |
|-----------|--------|-----------|------|---------|
| CLI | `testing` package | `assert` from testify or `t.Errorf` | -- | `os/exec` via `runCLI()` |
| API | `testing` package | `assert` from testify or `t.Errorf` | `net/http` Client or `net/http/httptest` | -- |
| TUI | `testing` package | Golden file comparison or snapshot | -- | `os/exec` |

All test files import `"testing"`. Prefer `github.com/stretchr/testify/assert` for assertions; standard `t.Errorf`/`t.Fatalf` is acceptable for minimal tests.

## Test File Template

| Test type | Template file | Output filename |
|-----------|--------------|-----------------|
| CLI | `templates/test-file.go` | `<feature>_cli_test.go` |
| API | `templates/test-file.go` | `<feature>_api_test.go` |
| TUI | `templates/test-file.go` | `<feature>_tui_test.go` |

Templates are available via `task profile get go-test --template <filename>`.

## Test Function Naming

All test functions must follow the `TestXxx` pattern where `Xxx` is a capitalized descriptor. Include the TC ID in the name for traceability:

```go
func TestTC_001_LoginWithValidCredentials(t *testing.T) { ... }
```

Pattern: `TestTC_NNN_Description` where `NNN` is the zero-padded test case number.

## Table-Driven Tests

When multiple test cases share the same logic, use table-driven tests:

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

## CLI Testing

Use `os/exec` package to invoke CLI binaries:

```go
func TestTC_NNN_CliCommand(t *testing.T) {
    cmd := exec.Command("myapp", "subcommand", "--flag", "value")
    out, err := cmd.CombinedOutput()
    assert.NoError(t, err)
    assert.Contains(t, string(out), "expected output")
}
```

## API Testing

Use `net/http/httptest` for handler testing or `net/http` Client for integration testing:

```go
func TestTC_NNN_ApiEndpoint(t *testing.T) {
    req, _ := http.NewRequest("GET", "http://localhost:8080/api/resource", nil)
    req.Header.Set("Authorization", "Bearer "+getToken())
    resp, err := http.DefaultClient.Do(req)
    assert.NoError(t, err)
    defer resp.Body.Close()
    assert.Equal(t, http.StatusOK, resp.StatusCode)
}
```

## TUI Testing

Use golden file comparison or snapshot testing:

```go
func TestTC_NNN_TuiOutput(t *testing.T) {
    out := runCLI(t, "myapp", "render")
    golden := filepath.Join("testdata", "tui_output.golden")
    if *update {
        os.WriteFile(golden, []byte(out), 0644)
    }
    expected, _ := os.ReadFile(golden)
    assert.Equal(t, string(expected), out)
}
```

## Build Tags

All test files must include the `e2e` build tag at the top:

```go
//go:build e2e
```

This ensures e2e tests are only compiled when explicitly requested via `-tags=e2e`.

## Import Conventions

Imports use the module path relative to the module root (as declared in `go.mod`). No relative import paths.

```go
import (
    "testing"
    "os/exec"
    "net/http"
    "github.com/stretchr/testify/assert"
)
```

## Auth

API authentication via custom headers or environment variables:

```go
// Header-based auth
req.Header.Set("Authorization", "Bearer "+os.Getenv("E2E_API_TOKEN"))

// Or via helper function
token := getAPIToken(t)
req.Header.Set("X-API-Key", token)
```

Credential caching at package level: `getAPIToken()` fetches once, stores in a package variable, returns cached value on subsequent calls.

## Anti-Patterns (Forbidden)

**Fixed delays**: No `time.Sleep` for synchronization. Use channels, `sync` primitives, or retry loops:

| Forbidden | Replacement |
|-----------|-------------|
| `time.Sleep(2 * time.Second)` | `withRetry(t, func() error { ... }, 10, 100*time.Millisecond)` |
| `time.Sleep` after server start | Health check with retry via `withRetry()` |
| Hardcoded ports | `net.Listen("tcp", ":0")` for ephemeral ports or env var configuration |

**Hardcoded ports**: Never hardcode `localhost:8080`. Use environment variables or ephemeral port allocation.

## Shared Infrastructure

`tests/e2e/main_test.go` provides `TestMain` for suite-level setup and teardown:

```go
//go:build e2e

package e2e

import (
    "os"
    "testing"
)

func TestMain(m *testing.M) {
    // Setup: start servers, seed data
    code := m.Run()
    // Teardown: stop servers, clean data
    os.Exit(code)
}
```

This file is created once and must NOT be overwritten during feature generation.

## Traceability

Each test function must include a traceability comment:

```go
// Traceability: TC-NNN -> {PRD Source}
```

## Compilation Check

After generating all test files:

```bash
just e2e-compile
```

## VERIFY Markers

Resolve all `// VERIFY:` comments using Fact Table values. Post-generation check:

```bash
just e2e-verify --feature <slug>
```

or: `grep -rn '// VERIFY:' tests/e2e/features/<slug>/ --include='*_test.go'`
