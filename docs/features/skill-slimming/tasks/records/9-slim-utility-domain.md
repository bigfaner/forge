---
status: "completed"
started: "2026-05-20 14:56"
completed: "2026-05-20 15:01"
time_spent: "~5m"
---

# Task Record: 9 Slim utility domain (clean-code + run-e2e-tests)

## Summary
Slimmed utility domain skills: clean-code (190 lines, already under 350 limit, no changes needed) and run-e2e-tests (299 -> 193 lines, 35.5% reduction). Extracted result-parsing rules and failure-diagnosis rules to rules/ subdirectory. No ambiguity terms found in either skill.

## Changes

### Files Created
- plugins/forge/skills/run-e2e-tests/rules/result-parsing.md
- plugins/forge/skills/run-e2e-tests/rules/failure-diagnosis.md

### Files Modified
- plugins/forge/skills/run-e2e-tests/SKILL.md

### Key Decisions
- clean-code SKILL.md (190 lines) was already well under the 350-line limit with no ambiguity terms; no structural changes needed
- Extracted 3 format parsing strategies (json-stream, json-report, text-verbose) from Step 4 into rules/result-parsing.md (63 lines)
- Extracted Failure Diagnosis section (App Health Gate + General Failure Analysis) into rules/failure-diagnosis.md (51 lines)
- Retained all workflow steps, error handling table, and output format in SKILL.md to preserve self-containment

## Test Results
- **Tests Executed**: No
- **Passed**: -1
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] clean-code SKILL.md <= 350 lines
- [x] run-e2e-tests SKILL.md <= 350 lines
- [x] No ambiguity terms remaining
- [x] All referenced auxiliary file paths exist and are readable

## Notes
Task 9 is the last task in the skill-slimming initiative. clean-code had no changes needed. run-e2e-tests reduced from 299 to 193 lines with 114 lines extracted to 2 rules files.
