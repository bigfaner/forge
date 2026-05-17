---
created: "2026-05-17"
tags: [local-dev-deployment]
---

# Post-Completion Status Commit Missing in run-tasks

## Problem

After all tasks in a feature completed, the `manifest.md` and `proposal.md` status frontmatter were updated (e.g., `status: tasks → completed`) but the changes were left uncommitted in the working tree.

## Root Cause

Causal chain:
1. **Symptom**: manifest.md status edit not committed after task loop ends
2. **Direct cause**: /run-tasks Post-Completion section only instructs "print summary" and "auto git push" — no explicit step to commit status updates
3. **Root cause**: The dispatcher's loop exit path has no "commit closing docs" step. Status transitions are done as file edits but never staged/committed
4. **Trigger**: Any feature where all tasks complete successfully — the status changes persist only in the working tree

## Solution

After the run-tasks loop ends and status updates are applied to manifest.md and proposal.md, commit them:

```bash
git add docs/features/<slug>/manifest.md docs/proposals/<slug>/proposal.md
git commit -m "docs: complete knowledge-discovery feature"
```

## Reusable Pattern

**Rule**: When a run-tasks dispatcher completes all tasks and updates status files, those status changes MUST be committed before printing the summary. The commit is part of the post-completion flow, not a separate manual step.

This applies to any feature lifecycle transition: `tasks → completed`, `tasks → in-progress → completed`.

## Example

```
# run-tasks post-completion flow:
1. Update manifest.md status → completed
2. Update proposal.md status → Completed
3. git add + git commit (status changes)  ← THIS WAS MISSING
4. Print summary
5. Auto git push (if enabled)
```

## Related Files

- `plugins/forge/skills/run-tasks/SKILL.md` — Post-Completion section should include commit step
- `plugins/forge/skills/quick/SKILL.md` — Step 4 references status transitions
