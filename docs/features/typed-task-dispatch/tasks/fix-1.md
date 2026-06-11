---
id: "fix-1"
title: "Fix: golangci-lint Go version mismatch (go1.25 vs go1.26.1)"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: golangci-lint Go version mismatch (go1.25 vs go1.26.1)

## Root Cause

golangci-lint was built with go1.25 but project targets go1.26.1. Upgrade golangci-lint or pin Go version to resolve.

## Reference Files

- Source: .
- Test script: just lint
- Test results: golangci-lint: go1.25 < go1.26.1

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

When this task is recorded as completed via `task record`, the source task 3.1 is automatically restored to pending if all its dependencies are completed.
