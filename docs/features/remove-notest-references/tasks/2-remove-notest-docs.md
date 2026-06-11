---
id: "2"
title: "Remove noTest references from skill docs, hooks, and templates"
priority: "P1"
estimated_time: "20m"
dependencies: []
type: "doc"
mainSession: false
---

# 2: Remove noTest references from skill docs, hooks, and templates

## Description
Remove all noTest references from skill documentation, hook guides, eval templates, and lesson docs. These references describe the deprecated noTest bypass mechanism that has been replaced by type-prefix-based testability detection.

## Reference Files
- `docs/proposals/remove-notest-references/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/hooks/guide.md` | Remove noTest-related guidance |
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Remove noTest references |
| `.claude/skills/eval-forge/templates/` | Remove noTest audit items from eval templates |
| `docs/lessons/gotcha-docs-only-needs-code-audit.md` | Remove noTest references |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] All noTest references removed from skill docs, hook guides, eval templates, and lesson docs
- [ ] Surrounding context remains coherent after removal (no dangling references or broken sentences)

## Hard Rules
- Do NOT modify files under `docs/proposals/`, `docs/forensics/`, or `docs/self-evolution/`

## Implementation Notes
- Remove entire sentences/paragraphs that are solely about noTest, not just the keyword
- For eval-forge templates, remove audit checklist items related to noTest bypass validation
- Verify no broken cross-references after removal
