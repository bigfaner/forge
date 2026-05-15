---
id: "1"
title: "Add Docs-Only Fast Path to quick-tasks/SKILL.md"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Add Docs-Only Fast Path to quick-tasks/SKILL.md

## Description

Add a `## Docs-Only Fast Path` section near the top of `quick-tasks/SKILL.md` (after Prerequisites, before Step 0) that explicitly documents which steps to skip when all business tasks are documentation type. This closes the gap where agents run unnecessary profile resolution and test generation for docs-only features.

## Reference Files
- `docs/proposals/docs-only-fast-path/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Add `## Docs-Only Fast Path` section after Prerequisites, before Step 0 |

## Acceptance Criteria

- [ ] `quick-tasks/SKILL.md` has a `## Docs-Only Fast Path` section positioned after Prerequisites and before Step 0
- [ ] The section lists Step 0 (Resolve Profile) and Step 4 (Test Tasks) as skippable for docs-only features
- [ ] The section defines how to detect docs-only: all business tasks use `templates/task-doc.md` (type: `"documentation"`)
- [ ] An agent reading only this file can determine the complete docs-only workflow

## Implementation Notes

- The fast path section should be a concise, scannable list — agents read this to decide skip behavior before executing steps
- Reference specific step names (Step 0, Step 4) so stale sections are easy to spot if steps are renumbered
- No runtime behavior changes — this is purely documenting existing skip behavior
