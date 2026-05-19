---
id: "3"
title: "Strengthen docs-only classification in quick-tasks and breakdown-tasks"
priority: "P1"
estimated_time: "30min"
dependencies: ["1"]
type: "documentation"
mainSession: false
---

# 3: Strengthen docs-only classification in quick-tasks and breakdown-tasks

## Description
Update both task generator skills to explicitly reference the "classify by output artifact" rule from type-assignment.md. When all In Scope items target non-compilable files (.md, .yaml, .json under skills/docs), agents must assign `type: "documentation"` and use the `task-doc.md` template.

Evidence from past features shows agents assigned `type: "enhancement"` or `type: "implementation"` to tasks that only modify markdown files, because they classified by intent rather than output.

## Reference Files
- `docs/proposals/task-type-code-docs-boundary/proposal.md` — Source proposal
- `plugins/forge/references/shared/type-assignment.md` — Classification rule (task 1 output)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/quick-tasks/SKILL.md` | Add explicit reference to "按产出物分类" rule in Type Assignment section |
| `plugins/forge/skills/breakdown-tasks/SKILL.md` | Add explicit reference to "按产出物分类" rule in Type Assignment section |

## Acceptance Criteria
- [ ] `quick-tasks SKILL.md` Type Assignment section explicitly states "classify by output artifact, not intent" with cross-reference to type-assignment.md classification rule
- [ ] `breakdown-tasks SKILL.md` Type Assignment section explicitly states the same rule
- [ ] Both skills reference the new classification table from type-assignment.md

## Hard Rules
- Must load `docs/conventions/forge-distribution.md` before modifying plugin files

## Implementation Notes
- Both skills already read `type-assignment.md` via `${CLAUDE_SKILL_DIR}/../../references/shared/type-assignment.md` — the change is to make the classification rule more prominent in the skill instructions
- Avoid duplicating the full classification table — reference type-assignment.md as the single source of truth
