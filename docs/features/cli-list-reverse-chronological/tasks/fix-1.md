---
id: "fix-1"
title: "Fix: T-quick-2 e2e test failures"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# Fix: T-quick-2 e2e test failures

## Root Cause

3 e2e tests (TC-001, TC-004, TC-005) fail asserting sort behavior. Implementation tasks 1 and 2 already completed — tests may need adjustment or may now pass.

## Reference Files

- Source: tests/e2e/features/cli-list-reverse-chronological/cli_list_reverse_chronological_cli_test.go
- Test script: go test ./tests/e2e/features/cli-list-reverse-chronological/...
- Test results: {{TEST_RESULTS}}

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

When this task is recorded as completed via `task record`, the source task T-quick-2 is automatically restored to pending if all its dependencies are completed.
