---
id: "fix-1"
title: "Fix: update autogen_test.go for interleaved dependency chain"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
---

# Fix: update autogen_test.go for interleaved dependency chain

## Root Cause

Hard Rules for task 1 restrict modifications to pipeline.go/pipeline_validate.go/pipeline_test.go only. The interleaved dependency wiring in pipeline.go changes the dependency chain behavior, causing 11 tests in autogen_test.go to fail because they assert the old serial chain wiring. These tests need updating to reflect the new interleaved dependencies: test-run-{surface} depends on gen-scripts-{surface}, and gen-scripts-{surface-N} (N>0) depends on test-run-{surface-N-1}.

## Reference Files

- Source: forge-cli/pkg/task/autogen_test.go
- Test script: cd forge-cli && go test -race ./pkg/task/ -run 'TestGetBreakdownTestTasks_PerKey_TwoTypes|TestGetBreakdownTestTasks_PerKey_ThreeTypes|TestGetQuickTestTasks_PerKey_ThreeTypes|TestGetBreakdownTestTasks_FullDependencyChain|TestGetQuickTestTasks_StagedAcrossTypesDependencyChain|TestGetBreakdownTestTasks_DefaultExecutionOrder_BackendBeforeFrontend|TestGetBreakdownTestTasks_SerialChain_BlockedOnUpstreamFailure|TestGetBreakdownTestTasks_ExplicitExecutionOrder|TestGetBreakdownTestTasks_AC1_FullDAG|TestGetQuickTestTasks_AC2_DirectDependencyChain|TestGetQuickTestTasks_MultiSurface_FullSerialChain'
- Test results: 11 tests fail due to interleaved dependency chain changes: T-test-run-{surface} now depends on T-test-gen-scripts-{surface} instead of serial chain

## Surface Inference

This fix-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `forge-cli/pkg/task/autogen_test.go` to extract the first file path (comma-separated).
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

When this task is recorded as completed via `task record`, the source task 1 is automatically restored to pending if all its dependencies are completed.
