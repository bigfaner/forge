---
id: "fix-1"
title: "fix unit-test: just unit-test failure in quality gate"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
surface-key: "."
surface-type: "cli"
---

# fix unit-test: just unit-test failure in quality gate

## Root Cause

Quality gate step `just unit-test` failed during quality-gate hook.

Error output saved to: `tests/results/unit-raw-output.txt`

Concise error:
```
--- FAIL: TestPersistentPreRun_InitsWithProjectRoot (0.01s)
--- FAIL: TestPersistentPreRun_InitsWithProjectRoot (0.01s)
```

## Reference Files

- Source: root_test.go
- Test script: just unit-test
- Test results: tests/results/unit-raw-output.txt

## Fix Boundaries

When fixing test failures, observe these boundaries:

**Forbidden:**
- Starting dev server (`npx expo start`, `npm run dev`, etc.)
- Running `npm install` more than 3 times — mark task as blocked if dependency installation fails 3 times
- Running full test suite — regression is verified by the dispatcher after fix completes
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

Full regression is verified by the dispatcher, not by this fix task.

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
