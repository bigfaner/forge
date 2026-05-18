---
id: "4"
title: "Delete old agent definitions and update distribution docs"
priority: "P1"
estimated_time: "30m"
dependencies: ["3"]
type: "cleanup"
mainSession: false
---

# 4: Delete old agent definitions and update distribution docs

## Description

Delete `doc-scorer.md` and `doc-reviser.md` agent definitions now that eval uses protocol+expert composition. Update `forge-distribution.md` to reflect the new file structure.

## Reference Files
- `docs/proposals/expert-template-eval/proposal.md` — Source proposal
- `docs/conventions/forge-distribution.md` — Distribution convention doc

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `docs/conventions/forge-distribution.md` | Remove doc-scorer/doc-reviser from agent listing, add experts/ to distribution tree |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/agents/doc-scorer.md` | Replaced by protocol + expert files |
| `plugins/forge/agents/doc-reviser.md` | Replaced by protocol + expert files |

## Acceptance Criteria

- [ ] `doc-scorer.md` deleted from `plugins/forge/agents/`
- [ ] `doc-reviser.md` deleted from `plugins/forge/agents/`
- [ ] `forge-distribution.md` distribution tree updated: `agents/` section shows `experts/` subdirectory with protocol/ and scorer/
- [ ] `forge-distribution.md` removes `doc-scorer.md` and `doc-reviser.md` from the agent listing
- [ ] `forge-distribution.md` "Core Dependencies" section updated: `doc-scorer` / `doc-reviser` references replaced with `agents/experts/` description
- [ ] No other files reference the deleted agent definitions (verify with grep)

## Hard Rules

- Do NOT delete until Task 3 is complete and eval SKILL.md no longer references the old agent types
- Grep for `doc-scorer` and `doc-reviser` across the entire codebase before deletion to confirm no remaining references

## Implementation Notes

- Simple cleanup task but critical ordering: must come after Task 3
- The `forge-distribution.md` updates are in the "分发包内容" section (lines 19-42) and "核心依赖" section (lines 69-75)
