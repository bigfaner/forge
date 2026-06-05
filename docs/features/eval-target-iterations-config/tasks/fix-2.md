---
id: "fix-2"
title: "Fix: Pre-existing test failures in qualitygate and scripts packages"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
---

# Fix: Pre-existing test failures in qualitygate and scripts packages

## Root Cause

Pre-existing test failures unrelated to fix-1 work: quality_gate_test.go line 154 expects error when no project root but gets nil; scripts install test fails. These block forge task submit quality gate.

## Reference Files

- Source: forge-cli/internal/cmd/qualitygate/quality_gate_test.go,forge-cli/scripts/
- Test script: go test -race ./forge-cli/internal/cmd/qualitygate/... ./forge-cli/scripts/...
- Test results: TestCheckAllCompleted_NoProject: expected error got nil; TestInstallSh/Ps1_PrintsVerificationInstructions: FAIL

## Surface Inference

This fix-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `forge-cli/internal/cmd/qualitygate/quality_gate_test.go,forge-cli/scripts/` to extract the first file path (comma-separated).
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

When this task is recorded as completed via `task record`, the source task fix-1 is automatically restored to pending if all its dependencies are completed.
