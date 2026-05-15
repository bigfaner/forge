---
status: "completed"
started: "2026-05-15 21:45"
completed: "2026-05-15 21:48"
time_spent: "~3m"
---

# Task Record: 1 Add drift detection + auto-fix steps to consolidate-specs SKILL.md

## Summary
Extended consolidate-specs SKILL.md with drift detection and auto-fix capabilities. Added Steps 9-11 (Detect Drift, Auto-fix Drift, Commit Changes), updated HARD-GATE to allow modifications when drift is detected, updated workflow diagram, and added drift-only mode support for quick workflows without PRD/design files. Renamed old Step 9 (Record Task) to Step 12.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/consolidate-specs/SKILL.md

### Key Decisions
- Drift detection compares rule keywords against actual code rather than simple text matching to mitigate false positives
- Drift-only mode is implicit: if prd/prd-spec.md and design/tech-design.md don't exist, skip Steps 1-8 and run only Steps 9-11
- New implicit rules from code are extracted with [CROSS] classification and presented to user before appending
- Project-global IDs are preserved during auto-fix - only description/behavior text is updated

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] SKILL.md contains new Steps 9-11 after existing Step 8
- [x] Step 9: Detect Drift - reads all docs/business-rules/*.md and docs/conventions/*.md, validates each rule against current code, classifies as current/drifted/orphaned
- [x] Step 10: Auto-fix Drift - updates drifted rules in-place preserving project-global IDs, removes orphaned rules with commit message recording rule ID + reason, detects implicit new rules from code changes
- [x] Step 11: Commit Changes - commits modified spec files with descriptive message listing changed rule IDs
- [x] HARD-GATE updated: second bullet changes to 'unless drift is detected in Step 9'
- [x] Workflow diagram updated to include Steps 9-11
- [x] Skill supports drift-only mode: when no PRD/design exists, skip Steps 1-8 and run only Steps 9-11

## Notes
Documentation-only task. No code changes or tests required.
