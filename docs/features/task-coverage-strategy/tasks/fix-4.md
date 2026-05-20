---
id: "fix-4"
title: "fix test-e2e: just test-e2e failure in quality gate"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# fix test-e2e: just test-e2e failure in quality gate

## Root Cause

Quality gate step `just test-e2e` failed during quality-gate hook.

Error output saved to: `tests/e2e/results/raw-output.txt`

Concise error:
```
...
Messages:   	index should contain shared gen-test-cases task
test_scripts_per_type_cli_test.go:602:
Error Trace:	/Users/fanhuifeng/Projects/ai/coding-harness/forge/.forge/worktrees/task-coverage-strategy/tests/e2e/test_scripts_per_type_cli_test.go:602
Error:      	Should be true
Test:       	TestTC_012_TaskIndexSharedInfrastructureNotDuplicated
--- FAIL: TestTC_012_TaskIndexSharedInfrastructureNotDuplicated (0.01s)
FAIL
exit status 1
FAIL	e2e-tests	2.602s
error: recipe `test-e2e` failed with exit code 1
```

## Reference Files

- Source: e2e_test_quality_cleanup_cli_test.go, gen_test_scripts_cli_test.go, test_guide_cli_test.go, quick_test_slim_cli_test.go, Users/fanhuifeng/Projects/ai/coding-harness/forge/.forge/worktrees/task-coverage-strategy/tests/e2e/quick_test_slim_cli_test.go, test_scripts_per_type_cli_test.go, Users/fanhuifeng/Projects/ai/coding-harness/forge/.forge/worktrees/task-coverage-strategy/tests/e2e/test_scripts_per_type_cli_test.go
- Test script: just test-e2e
- Test results: tests/e2e/results/raw-output.txt

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

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
