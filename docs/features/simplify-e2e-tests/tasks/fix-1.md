---
id: "fix-1"
title: "Fix: TC-003/TC-004 go.mod path mismatch"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# Fix: TC-003/TC-004 go.mod path mismatch

## Root Cause

TC-003/TC-004 run 'go test -tags=e2e ./tests/e2e/...' from project root, but e2e has its own go.mod (module e2e-tests). Must run from tests/e2e/ directory with './...' instead.

## Reference Files

- Source: tests/e2e/features/simplify-e2e-tests/simplify_e2e_tests_cli_test.go
- Test script: tests/e2e/features/simplify-e2e-tests/simplify_e2e_tests_cli_test.go
- Test results: tests/e2e/features/simplify-e2e-tests/results/latest.md

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
