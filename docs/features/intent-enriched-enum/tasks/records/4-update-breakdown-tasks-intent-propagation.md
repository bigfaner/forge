---
status: "completed"
started: "2026-05-31 15:03"
completed: "2026-05-31 15:05"
time_spent: "~2m"
---

# Task Record: 4 Update breakdown-tasks intent propagation to 1:1 mapping

## Summary
Updated breakdown-tasks SKILL.md: Intent Propagation now uses strict 1:1 mapping table (6 intent values) and Type Assignment table updated coding.fix entry to reflect fix intent auto-mapping constraint

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
无

## Document Metrics
Intent Propagation: 6-row mapping table; Type Assignment: 1 row updated

## Referenced Documents
- docs/proposals/intent-enriched-enum/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] Intent Propagation uses strict 1:1 mapping: new-feature->coding.feature, enhancement->coding.enhancement, refactor->coding.refactor, cleanup->coding.cleanup, fix->coding.fix, doc->doc
- [x] Type Assignment table entry for coding.fix updated to: fix intent auto-mapping but not CLI creatable
- [x] doc intent resolves to doc task type without sub-type distinction (doc.consolidate/doc.drift unified under doc umbrella)

## Notes
Per Architecture Decision in proposal: fix intent maps to coding.fix because it represents explicit user intent declaration in brainstorm, distinct from CLI manual creation. doc.consolidate and doc.drift are skill-auto-generated, unified under doc umbrella.
