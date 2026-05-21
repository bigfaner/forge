---
created: "2026-05-20"
tags: [architecture]
---

# Adjacent section over-removal when deleting coupled features

## Problem

Task 2 (remove run-tasks knowledge review) removed the "Commit Remaining Artifacts" section along with the "Knowledge Review" section. This left post-loop artifacts (fix-task records, index.json updates, manifest changes) uncommitted after the dispatcher finishes.

## Root Cause

1. **Surface**: After run-tasks loop completes, uncommitted files remain in the working tree
2. **Direct cause**: The "Commit Remaining Artifacts" section was removed from run-tasks.md along with the Knowledge Review section
3. **Root cause**: The commit section text began with "After knowledge extraction (Step 6)...", creating a false impression that it was part of the knowledge extraction flow. In reality, it committed ALL post-loop artifacts (records, index.json, manifest), not just knowledge files
4. **Deeper**: Adjacent sections with textual coupling (one references the other) create a deletion hazard. Removing the referenced section makes the referencing section appear orphaned, even when it serves an independent purpose

## Solution

When removing a section from a document:
1. Audit all adjacent sections for functional independence vs textual coupling
2. Decouple referencing sections before removing the referenced one
3. Preserve sections that serve independent purposes even if their text references the removed section

For run-tasks.md specifically: the Commit Remaining Artifacts section should have been rewritten to remove the knowledge extraction reference and preserved as a general post-loop cleanup step.

## Reusable Pattern

**Adjacent Section Audit Rule**: Before removing section X from a document, check every section within 1-2 section boundaries of X. For each adjacent section, ask: "Does this section serve a purpose independent of X?" If yes, preserve it (potentially rewriting references to X).

## Example

```
run-tasks.md before removal:
  ### Knowledge Review     ← target for removal
  ### Commit Remaining Artifacts  ← references Knowledge Review but commits ALL artifacts

Correct approach:
1. Remove Knowledge Review section
2. Rewrite Commit Remaining Artifacts to remove knowledge extraction reference
3. Preserve as "Post-loop artifact cleanup"

Incorrect approach (what happened):
1. Remove both sections together
```

## Related Files

- `plugins/forge/commands/run-tasks.md` — dispatcher protocol
- `docs/proposals/auto-knowledge-save/proposal.md` — proposal scope (only specified Knowledge Review removal)
