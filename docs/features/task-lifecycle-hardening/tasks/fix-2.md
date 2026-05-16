---
id: "fix-2"
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
ERROR: Feature not found: shared-feat
CAUSE: Feature directory does not exist
HINT: Check feature slug is correct
ACTION: ls docs/features/
---
--- FAIL: TestTC_012_TaskIndexSharedInfrastructureNotDuplicated (0.05s)
FAIL
exit status 1
FAIL	e2e-tests	5.399s
error: recipe `test-e2e` failed with exit code 1
```

## Reference Files

- Source: quick_test_slim_cli_test.go, project/ai/forge/tests/e2e/quick_test_slim_cli_test.go, task_lifecycle_hardening_cli_test.go, project/ai/forge/tests/e2e/task_lifecycle_hardening_cli_test.go, test_scripts_per_type_cli_test.go, project/ai/forge/tests/e2e/test_scripts_per_type_cli_test.go
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
4. `just test` — unit tests must pass
5. Record completion

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass

E2e regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
