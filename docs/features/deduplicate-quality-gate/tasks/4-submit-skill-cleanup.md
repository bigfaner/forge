---
id: "4"
title: "Submit-task SKILL.md cleanup — remove CLI-enforced validation rules"
priority: "P2"
estimated_time: "30min"
dependencies: []
type: "doc"
mainSession: false
---

# 4: Submit-task SKILL.md Cleanup

## Description

The submit-task SKILL.md currently contains detailed descriptions of CLI-enforced validation rules (quality gate pre-check, data validation table, "what submit does" description). These rules are implemented in `submit.go` and duplicated in the SKILL.md. The duplication causes drift — when the CLI logic changes (as in task 1), the SKILL.md must be updated separately.

Remove CLI-enforced validation rules from the SKILL.md. Keep only agent-unique instructions: metrics collection workflow, JSON data format, type reclassification, recovery steps, and forbidden operations.

## Reference Files
- `docs/proposals/deduplicate-quality-gate/proposal.md` — Source proposal (item 2)
- `plugins/forge/skills/submit-task/SKILL.md` — Current SKILL.md

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/SKILL.md` | Remove "Validation Rules (enforced by CLI)" section, "What forge task submit Does" section, quality gate pre-check details. Simplify Metrics Collection to remove `just test [scope]` command. Update coverage rules to reflect new tiered model. |

## Acceptance Criteria

- [ ] "Validation Rules (enforced by CLI)" section removed (lines ~148-183)
- [ ] "What forge task submit Does (One Command = 2 Operations)" section removed or reduced to a single line (lines ~133-145)
- [ ] Quality Gate Pre-check subsection removed (lines ~150-166)
- [ ] Data Validation subsection removed (lines ~168-183)
- [ ] Metrics Collection updated: remove `just test [scope]` command example; instruct agent to capture metrics from targeted test runs
- [ ] Coverage rules updated: `-1.0` auto-set for non-`coding.*` types only (no `noTest` mention)
- [ ] Remaining content: File Locations, JSON Data Format, Fields table, Metrics Collection (updated), Type Reclassification, Usage, Forbidden Operations, Recovery

## Implementation Notes

- The SKILL.md is distributed with the plugin (`~/.claude/plugins/cache/forge/forge/<version>/skills/submit-task/SKILL.md`). Changes here affect all users after plugin update.
- Keep the "Forbidden Operations" section — it contains agent-specific instructions that the CLI cannot enforce.
- Keep the "Recovery" section — it's a workflow guide for agents, not a CLI validation description.
- The `noTest` references in the Fields table should be removed (task 3 removes the flag). Instead, note that non-`coding.*` types automatically skip coverage requirements.
