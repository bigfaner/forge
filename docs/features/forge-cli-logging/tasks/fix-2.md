---
id: "fix-2"
title: "Fix: pre-existing TestCheckAllCompleted_NoProject test failure"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
---

# Fix: pre-existing TestCheckAllCompleted_NoProject test failure

## Root Cause

Pre-existing test failure: TestCheckAllCompleted_NoProject runs a subprocess that calls CheckAllCompleted(false), which calls base.Exit(err) -> os.Exit(1). The test expects the subprocess to succeed and return a non-nil error, but base.Exit terminates the process with exit code 1. This failure existed before task 5 migration. Fix requires updating the test to expect non-zero exit code from subprocess.

## Reference Files

- Source: forge-cli/internal/cmd/qualitygate/quality_gate_test.go
- Test script: go test ./internal/cmd/qualitygate/... -run TestCheckAllCompleted_NoProject
- Test results: test subprocess failed: exit status 1 -- CheckAllCompleted calls base.Exit(err) which calls os.Exit(1), but test expects subprocess to exit 0 and return error

## Surface Inference

This fix-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `forge-cli/internal/cmd/qualitygate/quality_gate_test.go` to extract the first file path (comma-separated).
2. Run `forge surfaces --json <file-path>` to resolve surface-key/type.
3. Use the resolved surface-type to load the appropriate `rules/surfaces/<type>.md` for test orchestration guidance.

If `forge surfaces --json` fails (no surfaces configured, command not found), proceed without surface information — this does not block the fix.

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

When this task is recorded as completed via `task record`, the source task 5 is automatically restored to pending if all its dependencies are completed.
