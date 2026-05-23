---
id: "4"
title: "Update submit-task SKILL.md with type-specific instructions"
priority: "P1"
estimated_time: "30m"
dependencies: ["1", "2", "3"]
type: "doc"
mainSession: false
---

# 4: Update submit-task SKILL.md with type-specific instructions

## Description

Update the `submit-task` skill to reflect type-differentiated record generation. Agents need to know which fields to fill based on task type. Currently the skill shows one uniform JSON format that assumes all tasks have test metrics.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `plugins/forge/skills/submit-task/SKILL.md` — Current skill (155 lines)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/SKILL.md` | Add type-specific record.json format sections for doc vs coding tasks |

## Acceptance Criteria
- [ ] SKILL.md has distinct JSON format examples for `coding.*` and `doc*` type tasks
- [ ] Doc task example shows `referencedDocs`, `reviewStatus`, `docMetrics` fields
- [ ] Doc task example shows `coverage: -1.0` (auto-set) and no test metrics
- [ ] Coding task example unchanged from current format
- [ ] Metrics Collection section clarifies: doc tasks omit test metrics, coding tasks remain mandatory
- [ ] Field table updated with type-specific required/optional annotations

## Hard Rules
- Follow forge plugin distribution model — this file is distributed to user projects
- Do NOT change the CLI command syntax or workflow steps
- Keep the existing structure; add type-specific sections, don't reorganize

## Implementation Notes
- Add a "Type-Specific Fields" section after the existing "Fields" table
- Add a conditional example block showing doc-type record.json vs coding-type record.json
- The existing `typeReclassification` section stays unchanged
