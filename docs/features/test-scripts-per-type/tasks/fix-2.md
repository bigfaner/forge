---
id: "fix-2"
title: "Fix: rewrite e2e tests to use actual CLI interfaces"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# Fix: rewrite e2e tests to use actual CLI interfaces

## Root Cause

E2E tests call non-existent CLI commands (forge gen-test-scripts, forge breakdown-tasks, forge quick-tasks). Must rewrite to test actual interfaces: forge task index (per-type task generation), forge prompt get-by-task-id (type inference + --type pass-through), and Go unit tests for DetectTypesFromTestCases.

## Reference Files

- Source: tests/e2e/features/test-scripts-per-type/test_scripts_per_type_cli_test.go
- Test script: just test-e2e --feature test-scripts-per-type
- Test results: tests/e2e/features/test-scripts-per-type/results/latest.md

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
