---
id: "fix-1"
title: "fix lint: just lint failure in quality gate"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# fix lint: just lint failure in quality gate

## Root Cause

Quality gate step `just lint` failed during quality-gate hook.

Error output saved to: `tests/results/unit-raw-output.txt`

Concise error:
```
...
// TaskTypeInfo describes a single task type for display and discovery.
^
..\..\..\forge\forge-cli\pkg\task\types.go:120:6: exported: type name will be used as task.TaskIndex by other packages, and that stutters; consider calling this Index (revive)
^
..\..\..\forge\forge-cli\pkg\task\types.go:187:6: exported: type name will be used as task.TaskState by other packages, and that stutters; consider calling this State (revive)
}
^
3 issues:
* revive: 3
error: recipe `lint` failed with exit code 1
```

## Reference Files

- Source: types.go
- Test script: just lint
- Test results: tests/results/unit-raw-output.txt

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
