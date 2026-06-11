---
id: "2"
title: "Delete /record-decision and /learn-lesson skills"
priority: "P1"
estimated_time: "30m"
dependencies: ["1"]
type: "cleanup"
mainSession: false
---

# 2: Delete /record-decision and /learn-lesson skills

## Description

Remove the old `/record-decision` and `/learn-lesson` skills now that their functionality has been absorbed into the unified `/learn` skill. The `record-decision` skill directory does not exist in `plugins/forge/skills/` (the decision archiving logic lives in `plugins/forge/references/shared/decision-logging.md` which must be preserved). The `learn-lesson` skill directory must be fully removed.

## Reference Files
- `docs/proposals/knowledge-accumulation-loop/proposal.md` — Source proposal

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
| `plugins/forge/skills/learn-lesson/SKILL.md` | Absorbed by /learn skill |
| `plugins/forge/skills/learn-lesson/templates/template.md` | Absorbed by /learn templates |
| `plugins/forge/skills/learn-lesson/examples/debug-race-condition.md` | Absorbed by /learn examples |

## Acceptance Criteria
- [ ] `plugins/forge/skills/learn-lesson/` directory is fully removed
- [ ] `plugins/forge/skills/record-decision/` directory does not exist (verify, not delete)
- [ ] `plugins/forge/references/shared/decision-logging.md` is preserved unchanged (shared reference used by tech-design and /learn)
- [ ] No other files reference `/record-decision` or `/learn-lesson` as active skills

## Hard Rules
- Do NOT delete `plugins/forge/references/shared/decision-logging.md` — it's a shared reference used by multiple consumers
- Verify `record-decision` skill directory doesn't exist before proceeding (it may have already been removed or never created as a skill directory)

## Implementation Notes
- After deletion, run a grep for `learn-lesson` and `record-decision` across `plugins/forge/` to find any remaining references that need updating (guide.md references are handled in Task 5)
