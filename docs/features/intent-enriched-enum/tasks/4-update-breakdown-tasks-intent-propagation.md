---
id: "4"
title: "Update breakdown-tasks intent propagation to 1:1 mapping"
priority: "P1"
estimated_time: "1h"
dependencies: [1]
type: "doc"
complexity: "low"
mainSession: false
---

# 4: Update breakdown-tasks intent propagation to 1:1 mapping

## Description
Update breakdown-tasks's Intent Propagation to strict 1:1 mapping for the 6 intent values. Update the Type Assignment table to reflect `coding.fix`'s new constraint (can be mapped from fix intent, but not manually created via CLI).

## Reference Files
- `docs/proposals/intent-enriched-enum/proposal.md` — Proposed Solution, Architecture Decision, Success Criteria
- plugins/forge/skills/breakdown-tasks/SKILL.md: Update Intent Propagation to 1:1 mapping; update Type Assignment table for coding.fix constraint (ref: Scope > In Scope)

## Affected Files

### Create

| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| plugins/forge/skills/breakdown-tasks/SKILL.md | Update Intent Propagation to 1:1 mapping (6 values); update Type Assignment table for coding.fix constraint |

### Delete

| File | Reason |
|------|--------|

## Acceptance Criteria
- [ ] Intent Propagation uses strict 1:1 mapping: new-feature→coding.feature, enhancement→coding.enhancement, refactor→coding.refactor, cleanup→coding.cleanup, fix→coding.fix, doc→doc
- [ ] Type Assignment table entry for `coding.fix` updated to: "可由 fix intent 自动映射，但不可通过 `forge task add` CLI 手动创建"
- [ ] `doc` intent resolves to `doc` task type without sub-type distinction (doc.consolidate/doc.drift unified under doc umbrella)

## Implementation Notes
- Per Architecture Decision: fix intent is allowed to map to coding.fix because it represents explicit user intent declaration in brainstorm, distinct from CLI manual creation
- doc.consolidate and doc.drift are skill-auto-generated types, not user-triggerable — they fall under the doc umbrella
