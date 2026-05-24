---
id: "fix-6"
title: "Fix: test file broken after re-export cleanup"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "fix"
---

# Fix: test file broken after re-export cleanup

## Root Cause

Task 9 removed re-export layer but output_contract_test.go still references PrintBlock, PrintBlockStart, PrintField through old import path. Update test imports to use base.* directly.

## Reference Files

- Source: forge-cli/internal/cmd/output_contract_test.go
- Test script: go build ./forge-cli/... && go test ./forge-cli/internal/cmd/...
- Test results: PrintBlock/PrintBlockStart/PrintField undefined in output_contract_test.go, unused import base

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

When this task is recorded as completed via `task record`, the source task 9 is automatically restored to pending if all its dependencies are completed.
