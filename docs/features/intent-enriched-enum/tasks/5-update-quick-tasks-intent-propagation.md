---
id: "5"
title: "Update quick-tasks intent propagation to 1:1 mapping"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
complexity: "low"
mainSession: false
---

# 5: Update quick-tasks intent propagation to 1:1 mapping

## Description
Update quick-tasks's Intent Propagation section to use the same strict 1:1 mapping as breakdown-tasks for the 6 intent values.

## Reference Files
- `docs/proposals/intent-enriched-enum/proposal.md` — Proposed Solution, Success Criteria
- plugins/forge/skills/quick-tasks/SKILL.md: Update Intent Propagation to 1:1 mapping (ref: Scope > In Scope)

## Affected Files

### Create

| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/quick-tasks/SKILL.md | Update Intent Propagation to 1:1 mapping (6 values) |

### Delete

| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Intent Propagation uses strict 1:1 mapping consistent with breakdown-tasks: new-feature→coding.feature, enhancement→coding.enhancement, refactor→coding.refactor, cleanup→coding.cleanup, fix→coding.fix, doc→doc
- [ ] Mapping table matches breakdown-tasks/SKILL.md exactly

## Implementation Notes
- This task should mirror the changes made in Task 4 (breakdown-tasks) for consistency
- quick-tasks may not have a Type Assignment table equivalent to breakdown-tasks — verify and update only the Intent Propagation section
