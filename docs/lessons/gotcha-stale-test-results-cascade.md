---
created: "2026-05-14"
tags: [testing, error-handling]
---

# Stale Test Result Files Cause Cascading False-Positive Fix Tasks

## Problem

E2e test runner writes results to `results/latest.md`. Downstream agents (T-quick-3 retry, T-quick-4 graduation) and hooks read this file as the source of truth. When a fix is applied between runs, the stale file still shows old failures, causing:

- T-quick-4 refuses to graduate (prerequisite: `latest.md` must show PASS)
- Subagent re-runs report blocked based on stale results
- Dispatcher creates fix tasks for failures that no longer exist

## Root Cause

Causal chain (4 levels):

1. **Symptom**: T-quick-3 reports blocked, T-quick-4 refuses to graduate — all downstream tasks stall
2. **Direct cause**: `results/latest.md` contains pre-fix failure data, downstream agents trust it blindly
3. **Root cause**: Test runner writes results to a file but never invalidates it; agents check the file instead of re-running tests
4. **Trigger condition**: Any workflow where fix tasks modify code between test generation and test execution, but the results file isn't regenerated

## Solution

After applying a fix task, always re-run the actual tests and regenerate the results file before marking the original task as resolved. Never trust a cached results file when code has changed since it was written.

In the dispatcher: after a fix task completes and unblocks a test task, the resumed test task should re-run tests from scratch, not check the old `latest.md`.

## Reusable Pattern

**Never read a test results file as truth when code has changed since it was written.** Always re-execute tests.

Checklist for fix-task → test-task handoff:
1. Fix task applies code changes
2. Test task re-runs the actual test command (not reads `latest.md`)
3. Test task writes fresh results
4. Only then does graduation/regression proceed

## Example

```bash
# BAD: trust stale file
cat tests/e2e/features/my-feature/results/latest.md | grep "Result: FAIL"

# GOOD: re-run tests
go test -tags=e2e ./tests/e2e/features/my-feature/... -v
```

## Related Files

- `forge-cli/tests/e2e/features/*/results/latest.md`
- `tests/results/unit-raw-output.txt`
