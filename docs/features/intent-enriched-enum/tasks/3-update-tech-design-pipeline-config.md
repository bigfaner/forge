---
id: "3"
title: "Update tech-design pipeline configuration for 6 intents with override signals"
priority: "P1"
estimated_time: "2h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 3: Update tech-design pipeline configuration for 6 intents with override signals

## Description
Replace tech-design's binary pipeline branching with the same 6-row Pipeline Configuration table and Override Signals mechanism as write-prd. Ensure the two copies remain synchronized. Update design quality checks for 6 intent values.

## Reference Files
- `docs/proposals/intent-enriched-enum/proposal.md` — Proposed Solution, Key Scenarios, Success Criteria, Key Risks
- plugins/forge/skills/tech-design/SKILL.md: Replace binary branching with Pipeline Configuration table + Override Signals table; add `<!-- Override: ... -->` comment generation logic (ref: Proposed Solution)
- plugins/forge/skills/tech-design/rules/design-quality-checks.md: Update intent-gated checks from 3 to 6 values (ref: Scope > In Scope)

## Affected Files

### Create

| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/tech-design/SKILL.md | Replace binary branching with Pipeline Configuration table (6 rows × 6 columns) + Override Signals table (5 signal types); add `<!-- Override: ... -->` comment generation instruction |
| plugins/forge/skills/tech-design/rules/design-quality-checks.md | Update all intent-gated checks from 3-value to 6-value enum |

### Delete

| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] tech-design/SKILL.md uses Pipeline Configuration table with 6 rows (identical structure to write-prd)
- [ ] Override Signals table matches write-prd's exactly (5 signal types, same keywords, same override actions)
- [ ] Override trigger generates `<!-- Override: ... -->` comment in tech-design output
- [ ] tech-design/rules/design-quality-checks.md intent-gated checks reference all 6 intent values
- [ ] Existing new-feature, refactor, cleanup pipeline artifacts unchanged from pre-modification behavior

## Hard Rules
- Pipeline Configuration table must be identical to write-prd/SKILL.md's copy (same rows, columns, defaults) — verify by diff

## Implementation Notes
- This task should be implemented by copying the Pipeline Configuration table and Override Signals table from write-prd/SKILL.md (Task 2), then adapting for tech-design context
- Key Risk: Synchronization drift between write-prd and tech-design copies — verify with diff after both tasks complete
