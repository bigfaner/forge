---
status: "completed"
started: "2026-05-15 21:47"
completed: "2026-05-15 21:50"
time_spent: "~3m"
---

# Task Record: 2 Rewrite scorer-prompt.md: 4-phase adversarial scoring methodology

## Summary
Complete rewrite of scorer-prompt.md from flat 12-dimension checklist to 4-phase adversarial scoring methodology: Phase 1 builds workflow graph (D1), Phase 2 per-node adversarial testing (D2), Phase 3 per-file precision review (D3+D4), Phase 4 baseline integrity (D5+D6). Output format updated to 6-dimension scorecard matching new rubric/report structure.

## Changes

### Files Created
无

### Files Modified
- .claude/skills/eval-forge/templates/scorer-prompt.md

### Key Decisions
- Added Phase 0 (Load References) as explicit first step to enforce Hard Rule that scorer must read rubric first
- Preserved Task CLI source code reading section but remapped dimensions from old D7 to new D1/D2/D5
- Output format uses 6-dimension scorecard with D1-D6 prefixes in ATTACKS section
- Phase 3 checks instruction conflicts as highest priority per rubric spec

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Scorer prompt implements 4-phase process
- [x] Phase 1 instructs scorer to read rubric, scan skills, compare specs vs actual, find breakpoints
- [x] Phase 2 instructs scorer to list gates, assume lazy agent, check HARD-RULE enforcement, eval loop independence, quality gate CLI enforcement
- [x] Phase 3 instructs scorer to check instruction conflicts first, then step ambiguity, then incomplete conditionals, then undefined variables, then content redundancy
- [x] Phase 4 instructs scorer to check reference integrity, frontmatter, eval templates, name alignment
- [x] Output format updated to 6-dimension scorecard
- [x] Scorer reads rubric and report template
- [x] Input section specifies all files to scan
- [x] EXTREMELY-IMPORTANT block preserved
- [x] Task CLI source code reading section preserved and adapted

## Notes
Documentation-only task. No test runner applicable. coverage=-1.0 per noTest convention.
