---
created: "2026-05-14"
tags: [testing, architecture]
---

# Quality Gate Stop Hook Creates Fix-Task Loop on Stale Failures

## Problem

The `forge quality-gate` stop hook runs `just test` at session end. When it detects a failure, it auto-creates a P0 fix task. But the hook reads a cached `tests/results/unit-raw-output.txt` file instead of re-running tests live. This caused:

- fix-3, fix-4, fix-5, fix-6, fix-7, fix-8 created in sequence
- Each fix task was closed as "transient failure, tests pass now"
- The hook re-read the same stale file and created another fix task
- Infinite loop until manually broken

## Root Cause

Causal chain (4 levels):

1. **Symptom**: Six consecutive fix tasks created for the same transient test failure
2. **Direct cause**: Hook reads `unit-raw-output.txt` which still contains the old failure output
3. **Root cause**: The hook doesn't re-run `just test` before creating a fix task — it trusts the output file from the previous run. Even after fix tasks close and re-run tests successfully, the file isn't updated because the fix tasks completed without code changes (tests were already passing)
4. **Trigger condition**: Any non-deterministic test that fails intermittently, or any test that fails once then passes on retry (race condition, timing, cached results)

## Solution

The quality gate hook should:

1. Re-run `just test` live (not read cached output)
2. Only create a fix task if the live run actually fails
3. Invalidate `unit-raw-output.txt` after a successful re-run

## Reusable Pattern

**Hooks that trigger on test failures must re-run tests live, never read cached result files.** Cached results become stale the moment any code changes.

For the fix-task auto-creation pattern:
- Before creating a fix task, re-execute the failing command
- If the command passes on re-run, don't create a fix task (was transient)
- If it fails again, create the fix task with the fresh failure output

## Example

```bash
# BAD: hook reads stale file
if grep -q "FAIL" tests/results/unit-raw-output.txt; then
  forge task add --template fix-task ...
fi

# GOOD: hook re-runs live
if ! just test > /dev/null 2>&1; then
  just test 2>&1 | tee tests/results/unit-raw-output.txt
  forge task add --template fix-task ...
fi
```

## Related Files

- `tests/results/unit-raw-output.txt` — cached test output
- `plugins/forge/hooks/` — stop hook configuration
- `forge-cli/internal/cmd/all_completed.go` — all-completed hook logic
