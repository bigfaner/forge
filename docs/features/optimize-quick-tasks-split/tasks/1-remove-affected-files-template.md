---
id: "1"
title: "Remove Affected Files from implementation task template"
priority: "P1"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Remove Affected Files from implementation task template

## Description
Remove the `## Affected Files` section (Create/Modify/Delete tables) from `quick-tasks/templates/task.md` (the implementation task template). This section requires file-level precision that proposals cannot provide, causing agents to over-research downstream skill files to fill the table.

The doc task template (`task-doc.md`) must remain unchanged — doc tasks produce files as deliverables, so target paths are knowable at creation time.

## Reference Files
- `docs/proposals/optimize-quick-tasks-split/proposal.md` — Source proposal
- `plugins/forge/skills/quick-tasks/templates/task.md` — Template to modify
- `plugins/forge/skills/quick-tasks/templates/task-doc.md` — Template to verify unchanged

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/templates/task.md` | Remove `## Affected Files` section (L21-36, the Create/Modify/Delete tables) |

## Acceptance Criteria
- [ ] `templates/task.md` no longer contains `## Affected Files` section
- [ ] `templates/task.md` retains all other sections (Description, Reference Files, Acceptance Criteria, Hard Rules, Implementation Notes)
- [ ] `templates/task-doc.md` still has `## Affected Files` section (unchanged)

## Implementation Notes
- Only remove the Affected Files section from the implementation template
- The frontmatter and all other sections stay intact
- Verify task-doc.md is not accidentally modified
