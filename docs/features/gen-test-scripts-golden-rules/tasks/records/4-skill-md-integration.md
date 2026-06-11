---
status: "completed"
started: "2026-05-21 23:17"
completed: "2026-05-21 23:18"
time_spent: "~1m"
---

# Task Record: 4 Integrate types/ into SKILL.md with loading logic

## Summary
Modified gen-test-scripts SKILL.md to add Step 2.5 (Load Type Rules) between Step 2 (Read Contract Specifications) and Step 3 (Generate Test Code). Added two HARD-RULE declarations for types/ vs Convention priority and Reconnaissance Hints usage. Step 2.5 follows gen-test-cases's per-type dispatch pattern: extract interface types from Contracts, always load _shared.md, load only matching type files, emit WARNING if >3 types detected.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/gen-test-scripts/SKILL.md

### Key Decisions
- Step 2.5 extracts interface types from Contract Actions/Outcomes rather than from project structure (Contracts are the immediate input for generation)
- Per-type dispatch pattern mirrors gen-test-cases: _shared.md always loaded, type files loaded only for detected types
- Token budget warning at >3 types is advisory (not blocking) to match proposal's 'WARNING' language

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] New step added between Step 2 and Step 3: Step 2.5 Load Type Rules
- [x] Step 2.5 reads Contract files, extracts all interface types referenced
- [x] Step 2.5 loads _shared.md (always) + matching type files from types/ (only for detected interface types)
- [x] Step 2.5 emits WARNING if >3 types detected (token budget risk)
- [x] HARD-RULE added declaring types/ vs Convention priority
- [x] HARD-RULE added: Reconnaissance Hints are discovery aids only, must convert to Fact Table values
- [x] Step 2.5 is consistent with gen-test-cases's type loading pattern (per-type dispatch, not bulk load)

## Notes
无
