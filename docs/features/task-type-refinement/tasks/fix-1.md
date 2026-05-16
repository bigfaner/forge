---
id: "fix-1"
title: "Fix: e2e tests failed due to stale forge binary"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# Fix: e2e tests failed due to stale forge binary

## Root Cause

14/20 e2e tests failed because the installed forge binary did not include the new type constants, prompt templates, pipeline logic, and migration changes. Binary has been rebuilt and reinstalled.

## Reference Files

- Source: forge-cli/pkg/task/,forge-cli/pkg/prompt/,forge-cli/pkg/proposal/,forge-cli/internal/cmd/
- Test script: go test ./tests/e2e/features/task-type-refinement/...
- Test results: docs/features/task-type-refinement/testing/

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
