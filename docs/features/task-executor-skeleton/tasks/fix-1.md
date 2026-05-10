---
id: "fix-1"
title: "Fix: 12 pre-existing e2e regressions expect mixed-project justfile"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: 12 pre-existing e2e regressions expect mixed-project justfile

## Root Cause

12 e2e tests fail because they expect a mixed-project justfile with scope dispatch (frontend/backend branches), but the project is backend-type with scope-free recipes. Affected tests: justfile-e2e-integration (TC-002, TC-016, TC-FJ-001/005/006/007/008/014/015), justfile-execution (TC-017, TC-003), scope-resolution (TC-015). Root cause: tests read the real justfile and assert mixed-project structure, but just project-type returns 'backend'.

## Reference Files

- Source: justfile, tests/e2e/justfile-e2e-integration/, tests/e2e/justfile-execution/, tests/e2e/scope-resolution/
- Test script: just test-e2e
- Test results: tests/e2e/results/

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

When this task is recorded as completed via `task record`, the source task T-test-4.5 is automatically restored to pending if all its dependencies are completed.
