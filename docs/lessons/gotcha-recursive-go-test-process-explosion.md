---
created: "2026-05-16"
tags: [testing, local-dev-deployment]
---

# Recursive `go test` Call Causes Process Explosion on Windows

## Problem

After running e2e tests, 126+ orphaned `e2e-tests.test.exe` and 137+ `go.exe` processes remained running, consuming ~48MB RAM each (~6GB total). The machine became sluggish.

## Root Cause

Causal chain (3 levels deep):

1. **Symptom**: Hundreds of orphaned test processes running long after the test command finished.
2. **Direct cause**: `go.exe` parent processes died (timed out), but their child `e2e-tests.test.exe` processes were never terminated.
3. **Root cause**: `simplify_e2e_tests_cli_test.go` TC-004 calls `exec.Command("go", "test", "-tags=e2e", "./...", ...)` from **within a test that belongs to the same package being tested**. This creates infinite recursion:
   - Outer `go test ./...` compiles and runs all tests
   - TC-004 spawns another `go test ./...`
   - The inner run encounters TC-004 again, spawning yet another level
   - Growth: ~3 packages × 3 subpackages per recursion level, bounded only by `-timeout 300s`
4. **Why orphans persist on Windows**: Unlike Linux (where child processes die with parent via `PR_SET_PDEATHSIG`), Windows has no automatic parent-death signal. When `go test` times out and panics, spawned child processes keep running indefinitely.

## Solution

Add a recursion guard via environment variable. When TC-004 spawns `go test`, it sets a guard variable; the inner test detects it and skips itself:

```go
func TestTC_004_VerifyRemainingCliBehaviorTestsPass(t *testing.T) {
    if os.Getenv("FORGE_E2E_RECURSION_GUARD") != "" {
        t.Skip("skipping to prevent recursive test execution")
    }
    cmd := exec.Command("go", "test", "-tags=e2e", "./...", "-count=1", "-timeout", "300s")
    cmd.Dir = e2eRoot(t)
    cmd.Env = append(os.Environ(), "FORGE_E2E_RECURSION_GUARD=1")
    out, err := cmd.CombinedOutput()
    // ...
}
```

## Reusable Pattern

**Never call `go test ./...` (or any test runner for the full package) from within a test in that same package.** This applies to all languages and test frameworks, not just Go.

When a test must verify "all other tests pass" (a meta-test):
1. Use an environment variable as a recursion guard — set it when spawning the subprocess, check it at the top of the meta-test.
2. Or exclude the meta-test itself from the run via `-run` flag filtering.
3. On Windows, always ensure child process cleanup. `exec.Command` does NOT kill children when the parent process dies. Consider using `cmd.Cancel` / `cmd.WaitDelay` (Go 1.20+) or process group kill via `syscall.SysProcAttr`.

## Example

**Dangerous** — `go test ./...` inside a test that's part of the same package:
```go
func TestAllTestsPass(t *testing.T) {
    cmd := exec.Command("go", "test", "./...", "-timeout", "5m")
    cmd.Dir = "tests/e2e"  // same package as this test!
    cmd.CombinedOutput()   // infinite recursion
}
```

**Safe** — recursion guard prevents re-entry:
```go
func TestAllTestsPass(t *testing.T) {
    if os.Getenv("RECURSION_GUARD") != "" {
        t.Skip("re-entry guard")
    }
    cmd := exec.Command("go", "test", "./...", "-timeout", "5m")
    cmd.Dir = "tests/e2e"
    cmd.Env = append(os.Environ(), "RECURSION_GUARD=1")
    cmd.CombinedOutput()
}
```

## Related Files

- `tests/e2e/simplify_e2e_tests_cli_test.go` (TC-003, TC-004)
- `justfile` (`test-e2e` recipe)
