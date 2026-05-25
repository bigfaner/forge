---
id: "4"
title: "Add SC internal consistency criterion to proposal rubric D9"
priority: "P1"
estimated_time: "1h"
dependencies: ["3"]
type: "doc"
mainSession: false
---

# 4: Add SC internal consistency criterion to proposal rubric D9

## Description

Modify `plugins/forge/skills/eval/rubrics/proposal.md` Dimension 9 (Success Criteria) to add an "SC internal consistency" criterion (25pts) and adjust existing criterion scores. The total D9 score remains 80pts.

Current D9 structure:
- Criteria are measurable and testable: 0-55 pts
- Coverage is complete: 0-25 pts

New D9 structure:
- Measurable and testable: 0-30 pts (was 0-55)
- Coverage is complete: 0-25 pts (unchanged)
- SC internal consistency: 0-25 pts (new)
- **Total: 80 pts** (unchanged)

## Reference Files

- `proposal.md#In-Scope` — Item 4: modify proposal rubric D9 with new criterion and score redistribution
- `proposal.md#Key-Risks` — risk of D9 score redistribution making existing proposal scores incomparable; SC #10 requires retro-test verifying <5% score difference
- `proposal.md#Success-Criteria` — SC requiring D9 to contain "SC internal consistency" criterion at 25pts with total 80pts unchanged

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rubrics/proposal.md` | Add SC internal consistency criterion to D9, adjust measurable from 55→30, add consistency 25pts |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] D9 contains "SC internal consistency" criterion worth 25pts with clear evaluation guidance (check SC↔SC and SC↔InScope for logical contradictions within clusters)
- [ ] "Criteria are measurable and testable" reduced from 55pts to 30pts
- [ ] "Coverage is complete" reduced from 25pts to 25pts (unchanged, as proposal says 40→25 for coverage but current rubric shows 25 — verify against actual current value)
- [ ] D9 total remains 80pts
- [ ] New criterion description checks SC internal satisfiability (intra-group SC↔SC and SC↔InScope), distinct from D10 which checks SC ↔ Scope/Solution alignment (no overlap)

## Hard Rules

- D9 and D10 must not overlap in responsibility — D9 checks SC internal consistency (within SC set), D10 checks SC ↔ external sections alignment
- Follow `docs/conventions/forge-distribution.md` distribution constraints

## Implementation Notes

- Verify the actual current rubric scores before editing — the proposal says "measurable 40→30, coverage 40→25" but the current rubric shows measurable 55pts and coverage 25pts. The discrepancy may indicate the rubric was already modified since the proposal was written, or the proposal referenced an older version. Use actual current values as the starting point.
- The score redistribution from measurable (55→30, -25pts) funds the new consistency criterion (+25pts), keeping total at 80
- The "What to check" column for the new criterion should reference the clustering + intra-group satisfiability check from scorer-protocol (Task 3)
