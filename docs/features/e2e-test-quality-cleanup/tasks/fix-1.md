---
id: "fix-1"
title: "Fix: TC-005 self-matching grep"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# Fix: TC-005 self-matching grep

## Root Cause

TC-005 greps all _test.go files for exec.Command("go", "test" but its own source file contains that literal string. The test should exclude the features/e2e-test-quality-cleanup/ directory from the grep, or exclude .go files that are the test itself.

## Reference Files

- Source: tests/e2e/features/e2e-test-quality-cleanup/e2e_test_quality_cleanup_cli_test.go
- Test script: TestTC_005_ZeroRecursiveGoTestInvocations
- Test results: tests/e2e/features/e2e-test-quality-cleanup/results/latest.md

## E2E Fix Boundaries

When fixing E2E test failures, observe these boundaries:

**Forbidden:**
- Starting dev server (`npx expo start`, `npm run dev`, etc.)
- Running `npm install` more than 3 times — mark task as blocked if dependency installation fails 3 times
- Running e2e tests (`just test-e2e`) — regression is verified by the dispatcher after fix completes
- Manually opening browser to verify rendering

**Correct workflow:**
1. Read failing test + corresponding component source
2. Compare test's expected testID/selectors vs actual DOM structure
3. Modify component (add testID) or test (adjust selectors/assertions)
4. `just test` — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass

E2e regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task T-quick-3 is automatically restored to pending if all its dependencies are completed.
