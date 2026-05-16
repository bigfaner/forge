---
id: "4"
title: "Create per-type rubrics (eval/rubrics/test-cases-*.md)"
priority: "P0"
estimated_time: "2h"
dependencies: ["1"]
type: "documentation"
mainSession: false
noTest: true
---

# 4: Create per-type rubrics (eval/rubrics/test-cases-*.md)

## Description

Decompose the current monolithic `eval/rubrics/test-cases.md` (1000 pts, 6 dimensions) into 5 per-type rubrics. Each per-type rubric keeps 5 shared dimensions (PRD Traceability 200, Step Actionability 250, Completeness 200, Structure & ID 100, Antipattern Prevention 100) and replaces Interface Accuracy (150) with a type-specific dimension using the full 150 pts — no percentage-based splitting.

## Reference Files
- `docs/proposals/gen-test-cases-per-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/eval/rubrics/test-cases.md` — Current monolithic rubric (source of decomposition)

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/eval/rubrics/test-cases-ui.md` | UI rubric: 5 shared dims + Visual State Accuracy (150 pts from web-ui sub-criteria) |
| `plugins/forge/skills/eval/rubrics/test-cases-tui.md` | TUI rubric: 5 shared dims + Output Assertion Accuracy (150 pts from tui sub-criteria) |
| `plugins/forge/skills/eval/rubrics/test-cases-mobile.md` | Mobile rubric: 5 shared dims + Interaction Accuracy (150 pts from mobile-ui sub-criteria) |
| `plugins/forge/skills/eval/rubrics/test-cases-api.md` | API rubric: 5 shared dims + Contract Accuracy (150 pts from api sub-criteria) |
| `plugins/forge/skills/eval/rubrics/test-cases-cli.md` | CLI rubric: 5 shared dims + Command Coverage Accuracy (150 pts from cli sub-criteria) |

### Modify
| File | Changes |
|------|---------|
| (none — keep existing `test-cases.md` rubric for legacy monolithic mode) |

### Delete
| File | Reason |
|------|--------|
| (none — keep legacy rubric for backward compat) |

## Acceptance Criteria
- [ ] 5 per-type rubric files created
- [ ] Each rubric has 6 dimensions totaling 1000 pts (5 shared + 1 type-specific)
- [ ] Shared dimensions preserved with same point values: PRD Traceability 200, Step Actionability 250, Completeness 200, Structure & ID 100, Antipattern Prevention 100
- [ ] Step Actionability blocking threshold (< 200 blocks gen-test-scripts) preserved in each rubric
- [ ] Type-specific dimension (150 pts) derived from current Interface Accuracy sub-criteria for that type — full 150 pts, no percentage splitting
- [ ] UI rubric: Visual State Accuracy covers Route Accuracy (60%) + Route Consistency (40%) from current web-ui criteria
- [ ] TUI rubric: Output Assertion Accuracy covers Output assertions (50%) + Keyboard interaction coverage (50%)
- [ ] Mobile rubric: Interaction Accuracy covers Interaction specificity (50%) + Navigation flow coverage (50%)
- [ ] API rubric: Contract Accuracy covers Contract accuracy (50%) + Error contract coverage (50%)
- [ ] CLI rubric: Command Coverage Accuracy covers Command coverage (50%) + Output assertion specificity (50%)
- [ ] Frontmatter: `scale: 1000`, `target: 900`, `iterations: 6`, `type: test-cases-{type}`

## Hard Rules
- Each per-type rubric must be self-contained — no references to other per-type rubrics or conditional branching
- Do NOT add multi-type filtering logic to per-type rubrics
- The Required Sections check must reflect per-type output format (single type section, not 5 grouped sections)

## Implementation Notes
- Decompose Interface Accuracy dimension: each type's sub-criteria (currently under percentage-based weights) become the full 150-pt type-specific dimension
- The Required Sections in each rubric should reference only that type's template format (e.g., UI rubric expects `ui-test-cases.md` format with frontmatter, UI TC section, traceability, route validation)
- Remove "Active capability filtering" logic from per-type rubrics — each rubric evaluates only its type
