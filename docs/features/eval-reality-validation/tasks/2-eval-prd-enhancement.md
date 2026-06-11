---
id: "2"
title: "eval-prd Rubric Enhancement"
priority: "P0"
estimated_time: "2h"
dependencies: [1]
type: "documentation"
mainSession: false
---

# 2: eval-prd Rubric Enhancement

## Description

Enhance the eval-prd rubric with two new dimensions: "Scenario Completeness" and "Edge Case Coverage". Add `context` frontmatter declaration so the scorer receives injected conventions and business-rules to detect contradictions between PRD content and project reality.

This is Batch 2 from the proposal — the first concrete application of context injection.

## Reference Files
- `docs/proposals/eval-reality-validation/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) | |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/rubrics/prd.md` | Add context frontmatter, add 2 new dimensions, adjust existing dimension points |

### Delete
| File | Reason |
|------|--------|
| (none) | |

## Acceptance Criteria

- [ ] `prd.md` rubric frontmatter includes `context` field declaring required conventions and business-rules categories
- [ ] New dimension "Scenario Completeness" added with clear criteria for evaluating whether PRD scenarios cover real-world usage patterns
- [ ] New dimension "Edge Case Coverage" added with criteria for error paths, boundary conditions, and failure states
- [ ] Existing dimensions adjusted to maintain 1000-point total scale
- [ ] Scoring criteria for new dimensions reference injected context: scorer should check PRD scenarios against loaded business-rules and conventions for contradictions
- [ ] Mode A/B detection still works correctly with new dimensions

## Hard Rules

- Total scale MUST remain 1000 points after adding new dimensions
- Do NOT remove existing dimensions — only adjust point allocations
- Do NOT modify eval SKILL.md — only the rubric file changes

## Implementation Notes

- Current prd.md has 5 dimensions: Background & Goals (150), Flow Diagrams (200), Functional Specs/Flow Completeness (200), User Stories (300), Scope Clarity (150).
- Adding 2 new dimensions requires reducing points from existing dimensions. Suggested: reduce User Stories from 300 to 200, reduce Flow Diagrams from 200 to 150, reduce Background & Goals from 150 to 100. New dimensions get ~150 each.
- "Scenario Completeness" should check: Are all user-facing scenarios described end-to-end? Are there implicit assumptions not stated? Do scenarios match known business rules?
- "Edge Case Coverage" should check: Are error paths documented? Are boundary conditions (empty input, max limits, concurrent access) covered? Do scenarios include failure recovery?
- The context frontmatter should declare `conventions: []` (no specific conventions needed for PRD eval) and `business-rules: auto` (load all business rules for contradiction detection).
