---
id: "T-test-5"
title: "Consolidate Specs"
priority: "P2"
estimated_time: "20min"
dependencies: ["T-test-4.5"]
status: pending
noTest: true
mainSession: false
---

# Consolidate Specs

## Description

Call `/consolidate-specs` skill to extract business rules from PRD and technical specifications from design into `specs/` directory.

## Reference Files

- `docs/features/task-executor-skeleton/prd/prd-spec.md` — Source for business rules
- `docs/features/task-executor-skeleton/design/tech-design.md` — Source for technical specs

## Acceptance Criteria

- [ ] `docs/features/task-executor-skeleton/specs/biz-specs.md` exists
- [ ] `docs/features/task-executor-skeleton/specs/tech-specs.md` exists
- [ ] If `[CROSS]` items exist: `review-choices.md` exists with user's approved/rejected items
- [ ] `docs/features/task-executor-skeleton/specs/.integrated` marker exists

## User Stories

No direct user story mapping. This is a standard knowledge consolidation task.

## Implementation Notes

1. Run `/consolidate-specs` skill
2. If ALL items are `[LOCAL]`, skip integration
3. If running non-interactively and CROSS items exist, mark blocked
