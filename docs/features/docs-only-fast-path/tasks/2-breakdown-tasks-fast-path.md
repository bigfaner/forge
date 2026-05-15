---
id: "2"
title: "Add Docs-Only Fast Path to breakdown-tasks/SKILL.md"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 2: Add Docs-Only Fast Path to breakdown-tasks/SKILL.md

## Description

Add a `## Docs-Only Fast Path` section near the top of `breakdown-tasks/SKILL.md` (after Prerequisites, before Step 0) that explicitly documents which steps to skip when all business tasks are documentation type. Mirrors the same pattern as task 1 for the full-pipeline skill.

## Reference Files
- `docs/proposals/docs-only-fast-path/proposal.md` — Source proposal

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Add `## Docs-Only Fast Path` section after Prerequisites, before Step 0 |

## Acceptance Criteria

- [ ] `breakdown-tasks/SKILL.md` has a `## Docs-Only Fast Path` section positioned after Prerequisites and before Step 0
- [ ] The section lists Step 0 (Resolve Profile) and Step 4b (Standard Test Tasks) as skippable for docs-only features
- [ ] The section defines how to detect docs-only: all business tasks use `templates/task-doc.md` (type: `"documentation"`)
- [ ] An agent reading only this file can determine the complete docs-only workflow

## Implementation Notes

- Use the same section structure as `quick-tasks/SKILL.md` (task 1) for consistency
- Reference Step 4b (not Step 4) since breakdown-tasks uses sub-step numbering
- No runtime behavior changes — purely documenting existing skip behavior
