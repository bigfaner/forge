---
created: "2026-05-17"
tags: [testing, local-dev-deployment]
---

# E2E Tests Inherit Project-Root Env Vars from Parent Process

## Problem

E2e tests that create isolated temp directories with `go.mod` and `.forge/state.json` fail because `forge` CLI resolves to the **real** project instead of the test's temp dir. Symptoms:

- Submit tests: record files not created in expected temp paths
- Query tests: output shows real feature data (e.g., `KEY: 1-submit-write-once`) instead of test fixture data (e.g., `KEY: task-1`)
- All 8 tests fail consistently, but unit tests in the same package pass

## Root Cause

Causal chain (4 levels):

1. **Symptom**: E2e tests querying temp-dir fixtures get real project data instead
2. **Direct cause**: `project.FindProjectRoot()` resolves to the real project root, ignoring `go.mod` in the temp dir
3. **Root cause**: `FindProjectRoot()` checks `CLAUDE_PROJECT_DIR` and `PROJECT_ROOT` environment variables **before** walking the filesystem. `exec.Command` inherits the parent process's environment, so these env vars point to the real project
4. **Trigger condition**: Running e2e tests from within Claude Code (which sets `CLAUDE_PROJECT_DIR`) or any environment where project-root env vars are set

## Solution

Strip project-root env vars from the child process environment in test helpers:

```go
var envBlacklist = []string{"CLAUDE_PROJECT_DIR", "PROJECT_ROOT"}

func cleanForgeEnv() []string {
    env := os.Environ()
    filtered := make([]string, 0, len(env))
    for _, e := range env {
        skip := false
        for _, bl := range envBlacklist {
            if strings.HasPrefix(e, bl+"=") {
                skip = true
                break
            }
        }
        if !skip {
            filtered = append(filtered, e)
        }
    }
    return filtered
}

func runForgeInDir(t *testing.T, dir string, args ...string) string {
    cmd := exec.Command("forge", args...)
    cmd.Dir = dir
    cmd.Env = cleanForgeEnv()  // Key fix
    // ...
}
```

## Reusable Pattern

**When writing e2e tests that spawn CLI subprocesses in isolated directories, always strip environment variables that override project root detection.**

`cmd.Dir` only sets the working directory — it does NOT prevent env-var-based resolution from bypassing filesystem markers. The env var takes precedence over filesystem traversal in `FindProjectRoot()`.

Diagnostic checklist for "test sees wrong project data":
1. Check if `cmd.Env` is set — if not, the child inherits parent's env
2. Check if the CLI's project-root logic has env-var overrides (grep for `Getenv`, `ENV`)
3. Strip those vars in test helpers

## Related Files

- `forge-cli/pkg/project/root.go` — `FindProjectRoot()` env-var check (lines 42-49)
- `tests/e2e/task_record_immutability_cli_test.go` — fixed test helpers
- `docs/lessons/gotcha-fix-task-index-test-isolation.md` — related: unit test isolation (different root cause)
