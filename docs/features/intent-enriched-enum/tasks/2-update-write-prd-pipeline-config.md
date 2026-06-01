---
id: "2"
title: "Update write-prd pipeline configuration for 6 intents with override signals"
priority: "P1"
estimated_time: "2h"
dependencies: [1]
type: "doc"
complexity: "high"
mainSession: false
---

# 2: Update write-prd pipeline configuration for 6 intents with override signals

## Description
Replace write-prd's binary pipeline branching (new-feature → Full PRD; others → Spec-only PRD) with a 6-row Pipeline Configuration table and Override Signals mechanism. Each intent maps to a specific PRD format (Full, Simplified, Spec-only, or Minimal). Override signals from PRD content can enable additional pipeline steps (API handbook, Security Review, etc.) on top of the intent baseline.

## Reference Files
- `docs/proposals/intent-enriched-enum/proposal.md` — Proposed Solution, Key Scenarios, Success Criteria, Key Risks
- plugins/forge/skills/write-prd/SKILL.md: Replace binary branching with Pipeline Configuration table + Override Signals table; add `<!-- Override: ... -->` comment generation logic (ref: Proposed Solution)
- plugins/forge/skills/write-prd/rules/self-check.md: Update intent-gated checks from 3 to 6 values (ref: Scope > In Scope)

## Affected Files

### Create

| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/write-prd/SKILL.md | Replace binary branching with Pipeline Configuration table (6 rows × 6 columns) + Override Signals table (5 signal types); add `<!-- Override: ... -->` comment generation instruction; add enhancement Simplified PRD format |
| plugins/forge/skills/write-prd/rules/self-check.md | Update all intent-gated checks from 3-value to 6-value enum |

### Delete

| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] write-prd/SKILL.md uses Pipeline Configuration table with 6 rows (one per intent) and 6 columns (PRD Format, User Stories, API Handbook, Test Pipeline, Security Review, at minimum)
- [ ] Override Signals table exists with 5 signal types: API 变更, 用户可见行为, 安全相关, 性能相关, 数据迁移
- [ ] Override trigger generates `<!-- Override: ... -->` comment in PRD output (e.g., `<!-- Override: API handbook enabled by signal "接口变更" -->`)
- [ ] Enhancement intent produces Simplified PRD format (Background + Goals + Test Pipeline), skipping User Stories
- [ ] Doc intent produces Minimal PRD format (title + goals + scope only)
- [ ] write-prd/rules/self-check.md intent-gated checks reference all 6 intent values
- [ ] Existing new-feature, refactor, cleanup pipeline artifacts unchanged from pre-modification behavior

## Hard Rules
- Pipeline Configuration table must match the table in tech-design/SKILL.md exactly (same 6 rows, same columns, same default values) — synchronized copy, not divergent

## Implementation Notes
- Override signals are detected during PRD content generation (same LLM call, not sequential scan). This is content generation + signal matching in parallel inference
- Negation handling: LLM should skip signals in negative context (e.g., "不涉及 API 变更"). Relies on LLM context understanding, not keyword matching
- Override only adds steps (开启), never removes. Worst case: unnecessary artifact generated, caught in user review
- For doc intent: Minimal PRD format has no pipeline steps that can be overridden — override signals become no-op by design
- Key Risk: Pipeline Configuration table is a synchronized copy with tech-design — changes must be applied to both files
