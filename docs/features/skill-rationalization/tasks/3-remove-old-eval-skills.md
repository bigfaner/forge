---
id: "3"
title: "Remove old eval skill directories, record-task skill, and simplify-skill command"
priority: "P1"
estimated_time: "30m"
dependencies: ["2"]
type: "documentation"
mainSession: false
---

# 3: Remove old eval skill directories, record-task skill, and simplify-skill command

## Description
Delete the 7 old eval skill directories (now superseded by the generic eval skill + command wrappers), the record-task skill (superseded by submit-task), and the simplify-skill command (rarely-used meta-tool).

## Reference Files
- `docs/proposals/skill-rationalization/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| (none) | |

### Delete
| File | Reason |
|------|--------|
| `plugins/forge/skills/eval-proposal/` | Replaced by eval skill + eval-proposal command wrapper |
| `plugins/forge/skills/eval-prd/` | Replaced by eval skill + eval-prd command wrapper |
| `plugins/forge/skills/eval-design/` | Replaced by eval skill + eval-design command wrapper |
| `plugins/forge/skills/eval-ui/` | Replaced by eval skill + eval-ui command wrapper |
| `plugins/forge/skills/eval-test-cases/` | Replaced by eval skill + eval-test-cases command wrapper |
| `plugins/forge/skills/eval-consistency/` | Replaced by eval skill + eval-consistency command wrapper |
| `plugins/forge/skills/eval-harness/` | Replaced by eval skill + eval-harness command wrapper |
| `plugins/forge/skills/record-task/` | Superseded by submit-task |
| `plugins/forge/commands/simplify-skill.md` | Rarely-used meta-tool |

## Acceptance Criteria
- [ ] All 7 eval skill directories fully removed
- [ ] `plugins/forge/skills/record-task/` no longer exists
- [ ] `plugins/forge/commands/simplify-skill.md` no longer exists
- [ ] No dangling references to deleted skills/commands in remaining files (verified via grep for `eval-proposal`, `eval-prd`, `eval-design`, `eval-ui`, `eval-test-cases`, `eval-consistency`, `eval-harness`, `record-task`, `simplify-skill`)

## Hard Rules
- Verify no other file references the deleted skills by name before deletion
- Do NOT delete `skills/submit-task/` — it is the canonical task recording skill

## Implementation Notes
- After deletion, grep for `eval-proposal`, `eval-prd`, etc. across remaining files to catch any dangling references
- The `eval-forge` command (in `.claude/skills/`) may reference these skill paths — update if needed
- Check `forge-cli/pkg/prompt/data/` for any CLI prompt templates that reference deleted eval skill paths
