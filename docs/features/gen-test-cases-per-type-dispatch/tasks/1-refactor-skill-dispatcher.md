---
id: "1"
title: "Refactor gen-test-cases SKILL.md into dispatcher"
priority: "P0"
estimated_time: "2h"
dependencies: []
type: "documentation"
mainSession: false
noTest: true
---

# 1: Refactor gen-test-cases SKILL.md into dispatcher

## Description

Refactor the monolithic `plugins/forge/skills/gen-test-cases/SKILL.md` (271 lines) into a slim dispatcher that handles shared Steps 0-2.5 (profile resolution, PRD reading, AC extraction, interface detection) and loops through active types loading per-type instruction files.

The dispatcher keeps all type-agnostic logic and delegates type-specific Steps 3-4 to `types/{type}.md` instruction files (created in Task 2). After the per-type loop, the dispatcher generates `testing/manifest.md` as the aggregator.

## Reference Files
- `docs/proposals/gen-test-cases-per-type-dispatch/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none — types/ files created in Task 2) |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-cases/SKILL.md` | Replace monolithic Steps 3-4 with per-type dispatch loop; add manifest generation step; add convention loading after Step 2.5; add `conventions` frontmatter field |

### Delete
| File | Reason |
|------|--------|
| (none) |

## Acceptance Criteria
- [ ] SKILL.md is under 150 lines
- [ ] Steps 0-2.5 are preserved (profile resolution, PRD reading, AC extraction, interface detection)
- [ ] After Step 2.5, dispatcher loads conventions from per-type instruction frontmatter `conventions` field
- [ ] Dispatcher loops through each active type, loading `types/{type}.md` via Read tool
- [ ] After per-type loop, dispatcher generates `testing/manifest.md` with summary table + cross-type traceability
- [ ] SKILL.md frontmatter includes `conventions: [testing-isolation.md]` for project-wide conventions
- [ } Convention loading: read per-type instruction frontmatter → check `docs/conventions/{filename}` exists → load or skip silently

## Hard Rules
- Do NOT include any type-specific Step 3-4 instructions in SKILL.md — those belong in `types/{type}.md`
- The manifest.md schema must follow the structure defined in the proposal (frontmatter with feature/types/generated, summary table, cross-type traceability table)
- Preserve existing HARD-GATE (only generates test case documents) and HARD-RULES (no invented ACs)

## Implementation Notes
- The split boundary is at Step 2.5/3 — Steps 0-2.5 produce type-agnostic data, Steps 3-4 are type-specific
- The per-type loop should sequentially process each detected type: load instruction → execute Steps 3-4 → write per-type output file
- Route Validation (Step 3.5) stays in the dispatcher since it references all test cases across types
- Integration TC generation (existing-page placements) moves to `types/ui.md` since it only applies to UI/Mobile types
