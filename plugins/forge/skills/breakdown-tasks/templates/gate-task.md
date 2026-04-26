---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "P0"
estimated_time: "1h"
dependencies: [{{DEPENDENCIES}}]
status: pending
breaking: true
---

# {{ID}}: {{TITLE}}

## Description

Cross-layer consistency gate. Verify that all outputs from the preceding phase are internally consistent and match the design specification before proceeding to the next phase.

## Verification Checklist

1. [ ] All interfaces from preceding phase compile without errors
2. [ ] Data models match the Cross-Layer Data Map in `design/tech-design.md`
3. [ ] No type mismatches between adjacent layers
4. [ ] Project builds successfully (`go build ./...` or equivalent)
5. [ ] All existing tests pass
6. [ ] No deviations from design spec (or deviations are documented as decisions)

## Reference Files

- `design/tech-design.md` — Cross-Layer Data Map section
- Preceding phase task records — `records/*.md`
- Preceding phase summary — `records/{{PREV_PHASE}}.summary-phase-summary.md` (if exists)

## Acceptance Criteria

- [ ] All verification checklist items pass
- [ ] Any deviations from design are documented as decisions in the record
- [ ] Record created via `/record-task` with test evidence

## Implementation Notes

This is a verification-only task. No new feature code should be written.
If issues are found:
1. Fix inline if trivial (e.g., type mismatch in a single file)
2. Document non-trivial issues as decisions in the record
3. Set status to `blocked` if a blocking issue cannot be resolved
