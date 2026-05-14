---
name: Task Record Immutability
tags: [architecture, testing]
---

# Task Records Must Be Append-Only, Never Overwritten

## Problem

After task 2 was completed with actual code implementation (commit `b6091ed`), a fix task (fix-1) was needed to resolve test failures. fix-1 was handled correctly — it has its own task record (`records/fix-1.md`) committed together with the fix in `2552563`. However, task 2's original execution record was then overwritten in commit `e036e16` with a "verification only" record claiming no code changes, replacing the original implementation record.

**Impact**: Reading `records/2-info-commands.md` alone gives a misleading picture ("no changes made"). While git history preserves the original, the document-level truth is lost. A reader unaware of the overwrite would conclude task 2 did nothing.

## Root Cause

- **Symptom**: Task 2 record shows "无" (none) for files created/modified, contradicting the actual implementation
- **Direct cause**: `/record-task` was re-invoked for task 2 after fix-1 completed, replacing the original record with a post-fix "verification" perspective
- **Root cause**: No mechanism enforces record immutability — `record-task` acts as upsert rather than append-only
- **Trigger condition**: When a fix task completes and the agent re-verifies the original task, it's tempting to "update" the record to reflect the final verified state, but this overwrites the original execution history

## Solution

Task records should follow an **append-only** model:

1. **First record = immutable**: Once written, a task record file should not be overwritten
2. **Fix tasks have their own records**: fix-1 already has its own `records/fix-1.md` with independent commit — this is the correct pattern
3. **If re-verification is needed**: Append to the existing record (add a `## Post-Fix Verification` section) rather than replacing it, or create a separate file (e.g., `records/2-info-commands-r2.md`)
4. **Cross-reference**: The fix task references the source task via `sourceTaskID` in index.json — the record should similarly link, not replace

## Reusable Pattern

**Task records are audit logs, not mutable state.** When a task produces issues that require a fix:

- Keep the original task record unchanged (documents what was initially done)
- Fix task has its own record with independent commit (documents what went wrong and how it was fixed)
- The pair of records together gives the full picture

If `record-task` detects an existing record, it should warn: "Record for task X already exists. Use a fix task record instead, or append a new section."

## Example

```
# WRONG — overwrite original record
records/2-info-commands.md  ← overwritten to show "no changes"
# (fix-1.md exists with its own commit, but original task 2 context is lost at doc level)

# RIGHT — append-only: keep both records as-is
records/2-info-commands.md  ← original (shows implementation work, 24m)
records/fix-1.md            ← fix (shows test fixes, 17m, independent commit)
```

## Related Files

- `docs/features/forge-info-commands/tasks/records/2-info-commands.md`
- `docs/features/forge-info-commands/tasks/records/fix-1.md`
- `docs/features/forge-info-commands/tasks/index.json`
