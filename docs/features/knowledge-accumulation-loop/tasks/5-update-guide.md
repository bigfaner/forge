---
id: "5"
title: "Update guide.md — replace old skills with /learn"
priority: "P1"
estimated_time: "30m"
dependencies: ["1", "2"]
type: "documentation"
noTest: true
mainSession: false
---

# 5: Update guide.md — replace old skills with /learn

## Description

Update `plugins/forge/hooks/guide.md` to replace references to `/record-decision` and `/learn-lesson` with the unified `/learn` skill. Document the auto-extract flow at pipeline completion points.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal
- `plugins/forge/hooks/guide.md` — Current guide

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Replace old skill references, document /learn and auto-extract |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria
- [ ] `decisions/` directory description updated: references `/learn` instead of `/record-decision`
- [ ] `lessons/` directory description updated: references `/learn` instead of `/learn-lesson`
- [ ] No remaining references to `/record-decision` or `/learn-lesson` as active skills
- [ ] `/learn` documented as the primary knowledge accumulation entry point
- [ ] Auto-extract flow documented: triggers at run-tasks, fix-bug, write-prd, tech-design completion
- [ ] `/consolidate-specs` still documented for bulk extraction + drift detection (unchanged role)

## Hard Rules
- Minimal changes — only update knowledge-related references, don't restructure the guide
- Preserve all other guide sections unchanged

## Implementation Notes
- The guide's Project-Level Documents section has `decisions/` and `lessons/` entries that reference old skills
- The Skill Workflow mermaid diagrams may need `/learn` added as an auxiliary skill
