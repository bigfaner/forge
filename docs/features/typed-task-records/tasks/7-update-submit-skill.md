---
id: "7"
title: "Update submit-task SKILL.md with per-type instructions"
priority: "P1"
estimated_time: "30m"
dependencies: ["1", "2", "3", "4", "5", "6"]
type: "doc"
mainSession: false
---

# 7: Update submit-task SKILL.md with per-type instructions

## Description

Update `submit-task` SKILL.md to reflect the 5-category record system. Agents need per-type JSON format examples and field guidance. Currently the skill shows one uniform format assuming all tasks have test metrics.

## Reference Files
- `docs/proposals/typed-task-records/proposal.md` — Source proposal
- `plugins/forge/skills/submit-task/SKILL.md` — Current skill (155 lines)

## Affected Files

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/submit-task/SKILL.md` | Add per-type record.json examples and field tables |

## Acceptance Criteria
- [ ] SKILL.md has JSON format examples for all 5 categories: coding, doc, test, validation, gate
- [ ] Doc example shows `referencedDocs`, `reviewStatus`, `docMetrics` fields (no test fields)
- [ ] Test example shows `casesGenerated`, `casesEvaluated`, `scriptsCreated`, `testResults` fields
- [ ] Validation example shows `validationPassed`, `issuesFound` fields
- [ ] Gate example shows `gatePassed`, `gateChecks` fields
- [ ] Coding example unchanged from current format
- [ ] Field table updated with per-category required/optional annotations
- [ ] Metrics Collection section clarifies per-category expectations

## Hard Rules
- Follow forge plugin distribution model — this file is distributed to user projects
- Do NOT change the CLI command syntax or workflow steps
- Keep existing structure; add type-specific sections, don't reorganize

## Implementation Notes
- Add a "Type-Specific Record Formats" section after the existing "Fields" table
- Each category gets its own JSON example block with comments explaining which fields matter
- The `typeReclassification` section stays unchanged
