---
id: "fix-5"
title: "Fix: type mismatch after dep check consolidation"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# Fix: type mismatch after dep check consolidation

## Root Cause

Task 6 changed function signatures to accept *TaskIndex but callers in validate_index.go and test files still pass map[string]Task or value TaskIndex. Fix type conversions at all call sites.

## Reference Files

- Source: forge-cli/internal/cmd/task/validate_index.go,forge-cli/internal/cmd/task/check_deps_test.go,forge-cli/internal/cmd/task/validate_index_test.go
- Test script: go build ./forge-cli/... && go test ./forge-cli/internal/cmd/task/...
- Test results: IncompatibleAssign: validateDependencies/validateWildcardSelfDeps/validateLiveness expect *TaskIndex but receive map[string]Task or value TaskIndex

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
4. Run targeted tests on affected packages — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. Run targeted tests on changed packages: `go test -race ./affected/package/...`
2. Replace the path with the actual packages you modified

> **Note:** Full project-wide tests run at CLI submit (`forge task submit`) — agent runs targeted tests only.

E2e regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task 6 is automatically restored to pending if all its dependencies are completed.
