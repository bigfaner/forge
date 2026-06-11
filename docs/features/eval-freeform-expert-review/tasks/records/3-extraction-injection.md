---
status: "completed"
started: "2026-05-23 20:06"
completed: "2026-05-23 20:07"
time_spent: "~1m"
---

# Task Record: 3 Extraction Prompt Template & Injection Mechanism

## Summary
Created extraction prompt template and injection mechanism rules for bridging freeform expert review with rubric scoring

## Changes

### Files Created
- plugins/forge/skills/eval/experts/freeform/extraction-prompt.md
- plugins/forge/skills/eval/rules/freeform-injection.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
2 files created, ~215 lines total

## Referenced Documents
- docs/proposals/eval-freeform-expert-review/proposal.md
- plugins/forge/skills/eval/rules/scorer-composition.md
- plugins/forge/skills/eval/experts/freeform/freeform-review-protocol.md
- plugins/forge/skills/eval/experts/freeform/freeform-reviewer.md

## Review Status
final

## Acceptance Criteria
- [x] Extraction prompt includes System/User roles, JSON array output format, summary/severity/quote fields, explicit-only rules
- [x] Extraction prompt uses {{FREEFORM_REVIEW}} placeholder
- [x] Injection rules define key findings appended as attack points to scorer prompt end
- [x] Injection rules define [beyond-rubric] tag for unmappable findings in ATTACKS list
- [x] Injection rules define contradiction annotation with exact marker text
- [x] JSON validation rules: non-empty, valid JSON, all fields non-empty, severity enum
- [x] Degradation logic: empty or invalid extraction skips injection, falls back to standard rubric
- [x] Partial extraction failure: hit rate < 50% triggers alert + attaches full narrative
- [x] Hit rate coarse heuristic estimation method with limitations documented
- [x] Hard rule: injected content explicitly labeled as from freeform expert review
- [x] Hard rule: extraction prompt only extracts explicitly stated risks, no inference

## Notes
Extraction prompt closely follows the template defined in the proposal. Injection rules integrate with existing scorer-composition.md by appending after protocol+expert+context blocks.
