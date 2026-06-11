---
status: "completed"
started: "2026-05-27 01:37"
completed: "2026-05-27 01:41"
time_spent: "~4m"
---

# Task Record: 12 Complete templates and guidance docs

## Summary
Fixed template and guidance documentation issues (Cluster 8: E1-E14, F1-F16): added template placeholder mapping docs to quick-tasks SKILL.md, fixed stage-gate misclaims, clarified REFERENCE_FILES instruction, fixed breakdown-tasks conditional phase-inventory checklist, added BIZ-error-reporting-001 resolution paths to run-tests, documented HARD-RULE as fourth level in prompt-template-hierarchy, added manual entry point note to execute-task

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/quick-tasks/SKILL.md
- plugins/forge/skills/breakdown-tasks/SKILL.md
- plugins/forge/skills/run-tests/SKILL.md
- docs/conventions/prompt-template-hierarchy.md
- plugins/forge/commands/execute-task.md

### Key Decisions
无

## Document Metrics
5 files modified, 0 created; 10/10 AC passed

## Referenced Documents
- docs/proposals/pipeline-spec-code-alignment/proposal.md
- docs/conventions/forge-distribution.md
- plugins/forge/skills/quick-tasks/templates/manifest-quick.md
- plugins/forge/skills/quick-tasks/templates/task.md
- plugins/forge/skills/quick-tasks/templates/task-doc.md
- plugins/forge/skills/gen-contracts/rules/journey-contract-model.md

## Review Status
completed

## Acceptance Criteria
- [x] manifest-quick.md uses single unified slug placeholder
- [x] manifest-quick.md does not reference non-existent testing/test-cases.md
- [x] quick-tasks/SKILL.md documents all template placeholder mappings
- [x] quick-tasks/SKILL.md Output Checklist is accurate for quick mode (no stage-gate claims)
- [x] breakdown-tasks/SKILL.md has a Commit step
- [x] gen-journeys/SKILL.md error messages match actual CLI output
- [x] run-tests/SKILL.md BIZ-error-reporting-001 has resolvable path
- [x] forge-distribution.md references /learn not /record-decision/learn-lesson
- [x] prompt-template-hierarchy.md documents HARD-RULE as fourth level
- [x] journey-contract-model.md has single canonical copy (not duplicated)

## Notes
Items verified as already correct: E1 (SLUG consistent across templates), E2 (test-cases.md ref already absent), E8 (breakdown-tasks Commit step exists), E12 (gen-journeys error message matches match.go), F9 (clean-code already uses just unit-test), F12 (forge-distribution already uses /learn). Added: template placeholder mapping table (15 placeholders), stage-gate clarification in Step 5, REFERENCE_FILES content replacement clarification, conditional phase-inventory checklist, BIZ-error-reporting-001 file paths, HARD-RULE fourth level documentation, execute-task manual entry point note.
