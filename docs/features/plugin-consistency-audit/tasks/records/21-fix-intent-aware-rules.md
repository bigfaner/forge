---
status: "completed"
started: "2026-05-30 06:16"
completed: "2026-05-30 06:18"
time_spent: "~2m"
---

# Task Record: 21 Fix: add intent-aware checks to write-prd + tech-design rule files

## Summary
Added intent-aware conditional checks to write-prd self-check rules and tech-design design-quality-checks rules

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/rules/self-check.md
- plugins/forge/skills/tech-design/rules/design-quality-checks.md

### Key Decisions
无

## Document Metrics
2 rules files updated, 3 AC items verified against SKILL.md intent branches

## Referenced Documents
- docs/features/plugin-consistency-audit/reports/06-consolidated-report.md
- docs/features/plugin-consistency-audit/reports/02-skills-batch-a.md

## Review Status
final

## Acceptance Criteria
- [x] write-prd/rules/self-check.md contains refactor/cleanup intent checks: Change Scope, Constraints, Verification Criteria; skips user stories/flow diagram/UI checks
- [x] tech-design/rules/design-quality-checks.md 5.1 PRD coverage branches by intent: new-feature from user stories, refactor/cleanup from Verification Criteria
- [x] tech-design/rules/design-quality-checks.md 5.5 DB schema check skips for refactor/cleanup intent

## Notes
Both rule files now mirror the intent branching that SKILL.md already implements. No other files were modified per Hard Rules.
