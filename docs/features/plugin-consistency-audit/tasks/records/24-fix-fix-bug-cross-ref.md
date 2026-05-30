---
status: "completed"
started: "2026-05-30 06:21"
completed: "2026-05-30 06:22"
time_spent: "~1m"
---

# Task Record: 24 Fix: fix-bug cross-reference precision

## Summary
Fixed imprecise cross-reference in fix-bug.md line 260: changed 'domain-to-file mapping from /consolidate-specs skill Step 5' to 'domain-to-decision-file mapping from /consolidate-specs rules/overlap-detection.md'

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/fix-bug.md

### Key Decisions
无

## Document Metrics
1 cross-reference fix, terminology aligned with overlap-detection.md

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/05-commands-agent-hooks.md
- plugins/forge/skills/consolidate-specs/rules/overlap-detection.md

## Review Status
final

## Acceptance Criteria
- [x] fix-bug cross-reference points to correct file path (rules/overlap-detection.md)
- [x] No longer references 'Step 5' for the mapping

## Notes
Minimal single-line change per Hard Rules. The mapping table is in overlap-detection.md line 8, titled 'Domain-to-decision-file mapping'.
