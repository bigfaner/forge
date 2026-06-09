---
id: "fix-1"
title: "Fix: forge worktree remove fails on corrupted worktrees"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
type: "coding.fix"
---

# Fix: forge worktree remove fails on corrupted worktrees

## Root Cause

Production code defect: forge worktree remove cannot handle corrupted worktrees (directory exists but .git file missing). git worktree remove --force fails validation. Need fallback to manual cleanup (os.RemoveAll + git worktree prune) when git reports .git validation failure. Affects 2 tests in corrupted-worktree-recovery journey.

## Reference Files

- Source: forge-cli/internal/cmd/worktree/cmd_remove.go,tests/corrupted-worktree-recovery/step2_remove_corrupted_test.go,tests/corrupted-worktree-recovery/corrupted_worktree_recovery_smoke_test.go
- Test script: cd tests && go test -v -tags=cli_functional -timeout=10m ./corrupted-worktree-recovery/...
- Test results: TestStep2_Success_RemovesCorruptedWorktree: FAIL - git worktree remove --force fails with 'validation failed, cannot remove working tree: .git does not exist'. TestJourney_CorruptedWorktreeRecovery: FAIL - cascade from step 2 failure. Root cause: cmd_remove.go has no fallback when git worktree remove fails on corrupted worktrees. Fix: detect missing .git validation error and fall back to os.RemoveAll + git worktree prune.

## Surface Inference

This fix-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `forge-cli/internal/cmd/worktree/cmd_remove.go,tests/corrupted-worktree-recovery/step2_remove_corrupted_test.go,tests/corrupted-worktree-recovery/corrupted_worktree_recovery_smoke_test.go` to extract the first file path (comma-separated).
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

When this task is recorded as completed via `task record`, the source task T-test-run is automatically restored to pending if all its dependencies are completed.
