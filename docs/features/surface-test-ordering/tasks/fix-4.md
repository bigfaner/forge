---
id: "fix-4"
title: "Fix: quick_test_slim old architecture assertions"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
surface-key: ""
surface-type: ""
---

# Fix: quick_test_slim old architecture assertions

## Root Cause

Rewrite quick_test_slim_test.go assertions: T-quick-gen-cases→T-test-gen-journeys, T-quick-gen-and-run-{type}→T-test-run-{key}, remove multi-profile suffixes, update task counts and dependency chains

## Reference Files

- Source: tests/test-generation/quick_test_slim_test.go
- Test script: cd tests && go test -tags e2e ./test-generation/... -count=1
- Test results: 12 tests fail with old T-quick-gen-and-run/T-quick-gen-cases IDs

## Surface Inference

This fix-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `tests/test-generation/quick_test_slim_test.go` to extract the first file path (comma-separated).
2. Run `forge surfaces --json <file-path>` to resolve surface-key/type.
3. Use the resolved surface-type to load the appropriate `rules/surfaces/<type>.md` for test orchestration guidance.

If `forge surfaces --json` fails (no surfaces configured, command not found), proceed without surface information — this does not block the fix.

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

When this task is recorded as completed via `task record`, the source task fix-3 is automatically restored to pending if all its dependencies are completed.
