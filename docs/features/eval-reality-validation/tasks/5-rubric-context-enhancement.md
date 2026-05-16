---
id: "5"
title: "Remaining Rubric Context Enhancement"
priority: "P2"
estimated_time: "3h"
dependencies: [1]
type: "documentation"
mainSession: false
---

# 5: Remaining Rubric Context Enhancement

## Description

Add `context` frontmatter declarations to all 14 existing rubrics, and enhance select rubrics with additional dimensions that leverage injected context for reality validation.

This is Batch 5 from the proposal — covers the full rubric set.

## Reference Files
- `docs/proposals/eval-reality-validation/proposal.md` — Source proposal (Batch 5 rubric table)
- `plugins/forge/skills/eval/rubrics/` — All 14 rubric files

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rubrics/design.md` | Add context frontmatter; consider adding "Implementation Feasibility" dimension |
| `plugins/forge/skills/eval/rubrics/proposal.md` | Add context frontmatter |
| `plugins/forge/skills/eval/rubrics/ui-web.md` | Add context frontmatter |
| `plugins/forge/skills/eval/rubrics/ui-mobile.md` | Add context frontmatter |
| `plugins/forge/skills/eval/rubrics/ui-tui.md` | Add context frontmatter |
| `plugins/forge/skills/eval/rubrics/test-cases.md` | Add context frontmatter; consider adding "Convention Compliance" dimension |
| `plugins/forge/skills/eval/rubrics/ui-test-cases.md` | Add context frontmatter; enhance D3 Visual State Accuracy |
| `plugins/forge/skills/eval/rubrics/tui-test-cases.md` | Add context frontmatter; enhance D3 Output Assertion Accuracy |
| `plugins/forge/skills/eval/rubrics/mobile-test-cases.md` | Add context frontmatter; enhance D3 Interaction Accuracy |
| `plugins/forge/skills/eval/rubrics/api-test-cases.md` | Add context frontmatter; enhance D3 Contract Accuracy |
| `plugins/forge/skills/eval/rubrics/cli-test-cases.md` | Add context frontmatter; enhance D3 Command Coverage Accuracy |
| `plugins/forge/skills/eval/rubrics/consistency.md` | Add context frontmatter |
| `plugins/forge/skills/eval/rubrics/harness.md` | Add context frontmatter |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] All 14 rubric files have `context` frontmatter with appropriate `conventions` and `business-rules` declarations
- [ ] Each rubric's `conventions` list is tailored to its evaluation focus (e.g., test-cases declares `testing-isolation`, ui-web declares `ux`)
- [ ] `design.md` considers an "Implementation Feasibility" dimension — if added, total remains 1000 pts
- [ ] `test-cases.md` considers a "Convention Compliance" dimension — if added, total remains 1000 pts
- [ ] Per-type test-cases rubrics (ui/tui/mobile/api/cli) have enhanced D3 dimensions that reference injected conventions
- [ ] Rubrics without `context` continue to work (backward compatible — `context` is optional)
- [ ] Total point scale remains 1000 for all affected rubrics (except harness which is 100)

## Hard Rules

- Do NOT change `scale`, `target`, or `iterations` values unless adding new dimensions requires point reallocation
- If new dimensions are added, total MUST still equal the declared `scale`
- Do NOT modify eval SKILL.md — this task only changes rubric files
- The `prd.md` rubric is NOT modified in this task (it was handled in Task 2)

## Implementation Notes

- Context frontmatter convention: each rubric declares which conventions it needs. The mapping should be semantic:
  - `design.md`: `conventions: [api, error-handling]` — design needs API and error handling conventions
  - `test-cases.md`: `conventions: [testing-isolation]` — test cases need isolation conventions
  - `ui-web.md`: `conventions: [ux, frontend]` — UI needs UX and frontend conventions
  - `consistency.md`: `conventions: []`, `business-rules: auto` — consistency checks all business rules
  - `harness.md`: `conventions: []`, `business-rules: []` — harness eval is self-contained
- For per-type test-cases rubrics, the D3 dimension enhancement should add criteria like: "Do test cases comply with project conventions for [type] testing?" and "Are there convention violations in test step descriptions?"
- "Implementation Feasibility" for design rubric: check whether the design references technologies/patterns that exist in the project, whether dependencies are available, whether the architecture fits the project structure.
