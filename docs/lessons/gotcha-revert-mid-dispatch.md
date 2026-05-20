# Debugging with git checkout Leaves HEAD on Wrong Branch

**Date:** 2026-05-20
**Feature:** test-knowledge-convention-driven
**Type:** gotcha

## Symptom

After an active `/run-tasks` dispatch, the working tree shows pre-refactoring code. Tasks that previously completed appear pending with no records. Subsequent dispatches fail because prerequisites are missing.

## Root Cause

1. **Immediate cause**: HEAD was on `v3.0.0` instead of `test-knowledge-convention-driven`, so all subsequent commits (including subagent block-task commits) landed on the wrong branch.
2. **Contributing factor**: During e2e debugging, ran `git checkout v3.0.0 -- .` to test clean-state e2e, then `git checkout test-knowledge-convention-driven -- .` to restore working directory — but this only restores files, not the branch itself.
3. **Root mistake**: `git checkout <branch> -- .` restores working tree files from a branch without switching HEAD to that branch. The correct command is `git checkout <branch>` (without `-- .`).

## Fix

When debugging by temporarily checking out other branches:
1. Use `git checkout <branch>` (full switch) instead of `git checkout <branch> -- .` (file-only)
2. After debugging, switch back with `git checkout <original-branch>`
3. Verify with `git branch --show-current` before resuming dispatch

Recovery: `git stash` → `git checkout <correct-branch>` → `git stash pop` → `forge task index --feature <slug>` → resume `/run-tasks`.

## How to Apply

If `/run-tasks` shows unexpected task state after debugging, run `git branch --show-current` first. If on wrong branch, switch and re-index. The `gotcha-revert-mid-dispatch` pattern (stop → revert → re-index → resume) applies, but the actual trigger is branch confusion, not intentional revert.
