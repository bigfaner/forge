# Reverting Code Mid-Dispatch Breaks Task Dependency Chain

**Date:** 2026-05-20
**Feature:** test-knowledge-convention-driven
**Type:** gotcha

## Symptom

Task dispatcher loops on blocked tasks after partial completion. Tasks that previously completed successfully show as pending, and their acceptance criteria can't be met because the code was reverted.

## Root Cause

During an active `/run-tasks` dispatch, reverting source files to a pre-refactoring state creates a mismatch between task definitions (which describe the refactoring) and the codebase state. The dispatcher re-claims tasks that appear pending (because index.json was also reverted), but the subagent finds the prerequisites missing.

## Fix

Avoid reverting code during active dispatch loops. If a revert is necessary:
1. Stop the dispatcher first
2. Revert the code
3. Re-index tasks to match the new code state (`forge task index`)
4. Resume dispatching

## How to Apply

When `/run-tasks` hits 3 consecutive failures after a revert, check `git log` for recent reverts before assuming a code bug. The fix is administrative (re-sync tasks), not technical.
