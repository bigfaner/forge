# Task Executor "Already Completed" Detection Ignores Feature Scope

## Problem

`permission-granularity` feature's T-test-2 (gen-test-scripts) never generated e2e test scripts. The task executor detected commit `51a3be7` from the `integration-test-coverage` feature (April 28) and incorrectly declared "task already completed", skipping execution entirely.

Consequence: `docs/features/permission-granularity/testing/scripts/` does not exist. Only `test-cases.md` was generated (by T-test-1). The feature has no executable e2e test scripts.

## Root Cause

**Causal chain (4 levels):**

1. **Symptom:** No scripts generated under `permission-granularity/testing/scripts/`
2. **Direct cause:** Task executor found commit `51a3be7` (from `integration-test-coverage`) via git log search for "T-test-2" and skipped execution
3. **Root cause:** Task executor's "already completed" detection searches git history globally by task ID, without verifying the matched commit belongs to the CURRENT feature
4. **Trigger condition:** Multiple features share standard task IDs (T-test-1, T-test-2) from `/breakdown-tasks`; once any feature's T-test-2 commit is merged to main, all subsequent features with the same task ID will be falsely detected as "already completed"

**Dispatcher also failed:** Record verification (Step 3) should have caught this — the "already completed on 2026-04-28" claim was from before the current session started, yet no new code was generated. The dispatcher did not question this inconsistency.

## Solution

For the task-executor agent template:
- When detecting "already completed", verify the matched commit's file paths contain the current feature slug
- Or: require feature-scoped task IDs (e.g., `pg-T-test-2` instead of `T-test-2`)

For the dispatcher:
- When an agent reports "already completed" with a commit dated before the session started, verify that the commit actually modified files under the current feature's directory
- If no files match the current feature, re-dispatch with explicit instruction to regenerate

## Key Takeaway

Standard task IDs (T-test-1, T-test-2) from `/breakdown-tasks` are not namespaced per feature. The task executor searches git history by task ID globally and will match commits from ANY feature. Always verify that "already completed" detection is scoped to the current feature — both in the task executor and in the dispatcher's record verification step.
