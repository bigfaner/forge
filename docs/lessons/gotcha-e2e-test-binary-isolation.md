---
created: "2026-05-20"
tags: [testing, local-dev-deployment]
---

# E2E Tests Must Not Depend on System-Installed Forge Binary

## Problem

E2e tests for a feature branch failed (8/36) because they invoked the system-installed `forge` binary instead of a branch-specific test binary. The system binary was compiled from `main` branch and lacked the new feature's code changes.

Symptoms:
- Tests for new CLI commands (e.g., `forge gen-test-scripts`) reported "unknown command"
- Tests for removed commands (e.g., `forge test detect`) still succeeded when they should have failed
- Tests for modified behavior (e.g., `forge config init`) produced old-format output

## Root Cause

1. **Surface**: 8 e2e tests failed with "command not found" or unexpected output format
2. **Direct cause**: Tests invoked `exec.Command("forge", ...)` which resolved to the system-installed binary via `$PATH`
3. **Systemic cause**: No convention existed requiring e2e tests to compile and use a dedicated forge binary built from the current source tree. Tests assumed the installed `forge` matched the branch under test.
4. **Design gap**: E2e test infrastructure (helpers, fixtures) focused on filesystem isolation (`CLAUDE_PROJECT_DIR`, `t.TempDir()`) but neglected binary isolation — the forge executable itself was treated as an external dependency rather than part of the system under test.

## Solution

E2e tests must compile a dedicated forge binary from the current source tree and use that binary for all CLI invocations. The binary should be built once per test suite (via `TestMain`) and cleaned up automatically.

## Reusable Pattern

When writing e2e tests for a CLI tool built from the same repository:

1. Build the binary from source in `TestMain` or a test helper, targeting a temp directory
2. Pass the binary path to all test functions via a package-level variable
3. Use the full path in `exec.Command(binaryPath, ...)` instead of relying on `$PATH` resolution
4. Never assume the system-installed binary matches the code under test

## Example

```go
var forgeBinary string

func TestMain(m *testing.M) {
    tmpDir, _ := os.MkdirTemp("", "forge-e2e-*")
    defer os.RemoveAll(tmpDir)

    forgeBinary = filepath.Join(tmpDir, "forge")
    buildCmd := exec.Command("go", "build", "-o", forgeBinary, "./cmd/forge")
    if out, err := buildCmd.CombinedOutput(); err != nil {
        fmt.Fprintf(os.Stderr, "build failed: %s\n%s", err, out)
        os.Exit(1)
    }

    os.Exit(m.Run())
}

func TestMyFeature(t *testing.T) {
    cmd := exec.Command(forgeBinary, "task", "claim")
    // ... assertions ...
}
```

## Related Files

- `docs/conventions/testing-isolation.md` — TEST-isolation-004

## References

- `docs/lessons/gotcha-e2e-env-override-isolation.md` — related filesystem isolation lesson
