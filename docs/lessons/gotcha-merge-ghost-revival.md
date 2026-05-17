---
created: "2026-05-17"
tags: [dependencies, architecture]
---

# Merge Ghost Revival: Deleted Files Silently Reintroduced by Feature Branch Merge

## Problem

A completed cleanup refactor (skill-rationalization: 7 eval-* skill directories → 1 generic eval skill) was silently undone. Weeks after merge, all 7 deleted directories reappeared in the codebase as if the cleanup never happened. No merge conflict was raised.

## Root Cause

Causal chain (3 levels):

1. **Symptom**: 7 eval-* directories exist despite rationalization task records showing "completed" and merge PR merged.
2. **Direct cause**: Commit `d7f8a13` on a separate feature branch (`gen-test-cases-per-type-dispatch`) added ~2000 lines restoring all 7 eval-* skill directories. This branch was developed from a base that predated the deletion, so the files existed in its history.
3. **Root cause**: Git merge treats file deletion and file addition as independent operations across branches. When branch A deletes a file and branch B adds the same file from a common ancestor, git sees no conflict — it simply carries forward the addition. There is no "this file was intentionally deleted" sentinel.
4. **Trigger condition**: Long-lived or stale-based feature branches that diverge before a cleanup/delete commit, then merge back carrying the old files as "new additions".

## Solution

Re-delete the resurrected directories. Verify with `git log --diff-filter=A` to identify which commit re-added them.

## Reusable Pattern

**After merging any branch that was created before a significant refactor/cleanup, verify that deleted files haven't been resurrected.**

Concrete checks (adapted for merge-based workflow):

1. **Pre-merge audit**: Before merging a stale-based branch, diff it against target and check for files that were intentionally deleted on target: `git diff --name-only --diff-filter=A target...branch | xargs -I{} sh -c 'git log --oneline target -- {} | head -1 | grep -q delete && echo "GHOST: {}"'`
2. **Post-merge verification**: After merge, run a checklist of known deletions to confirm they're still absent.
3. **Merge PR review**: When reviewing PRs that touch many files, specifically look for files that shouldn't exist anymore. GitHub diff will show them as green (added) — treat them with the same scrutiny as red (deleted).

The key insight: **git merge is additive by default**. It will never drop files that exist in either parent. Any intentional deletion must be actively protected, not passively assumed.

## Example

```bash
# Detect ghost revival after merge
# List files added by a branch that were previously deleted on target
git log --diff-filter=A --name-only --pretty=format: branch..target | sort -u | while read f; do
  if git log --oneline --diff-filter=D -- "$f" | grep -q .; then
    echo "GHOST REVIVAL: $f was deleted but re-added"
  fi
done
```

## Related Files

- `plugins/forge/skills/eval-*/` — resurrected directories
- `plugins/forge/skills/eval/` — correct generic eval skill
- `docs/proposals/skill-rationalization/proposal.md` — original cleanup proposal
- `docs/features/skill-rationalization/tasks/records/3-remove-old-eval-skills.md` — deletion record

## References

- Git merge strategies and rename/deletion handling
- [skill-rationalization proposal](../proposals/skill-rationalization/proposal.md)
