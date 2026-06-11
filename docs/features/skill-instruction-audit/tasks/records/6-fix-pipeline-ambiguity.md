---
status: "completed"
started: "2026-05-28 23:04"
completed: "2026-05-28 23:07"
time_spent: "~3m"
---

# Task Record: 6 Fix ambiguity and logic issues in pipeline commands/skills

## Summary
Fixed ambiguity and logic issues across 6 pipeline files: quick.md fallback, execute-task.md MAIN_SESSION, run-tasks.md definitions, gen-journeys matching rule, submit-task reclassification criteria, fix-bug.md duplicate+ARGUMENTS

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/quick.md
- plugins/forge/commands/execute-task.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/skills/gen-journeys/SKILL.md
- plugins/forge/skills/submit-task/SKILL.md
- plugins/forge/commands/fix-bug.md

### Key Decisions
无

## Document Metrics
6 files, 6 AC items, all PASS

## Referenced Documents
- docs/features/skill-instruction-audit/tasks/6-fix-pipeline-ambiguity.md

## Review Status
final

## Acceptance Criteria
- [x] quick.md non-zero fallback = show gate
- [x] execute-task.md defines MAIN_SESSION; Step 1.5 self-contained
- [x] run-tasks.md defines successful = STATUS completed; defines T-test-run; has slug failure path
- [x] gen-journeys has case-insensitive scenario matching rule
- [x] submit-task has objective type reclassification criteria
- [x] fix-bug has no duplicate between E-I and HARD-GATE; ARGUMENTS parsing clarified

## Notes
All changes are documentation-only edits clarifying ambiguous or missing definitions. No functional code changes.
