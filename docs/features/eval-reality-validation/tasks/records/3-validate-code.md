---
status: "completed"
started: "2026-05-17 01:27"
completed: "2026-05-17 01:30"
time_spent: "~3m"
---

# Task Record: 3 validate-code Implementation

## Summary
Created validate-code eval type: rubric with scenario traceability/path completeness/code-PRD consistency dimensions (1000pt, iterations=1), task templates for both breakdown-tasks and quick-tasks (type=gate, breaking=false), and extended eval SKILL.md with validate-code entries in Prerequisites, Parameters, Pre-Processing, Report path, Default Doc Dir, Next Step, and Rubric Reference tables.

## Changes

### Files Created
- plugins/forge/skills/eval/rubrics/validate-code.md
- plugins/forge/skills/breakdown-tasks/templates/validate-code-task.md
- plugins/forge/skills/quick-tasks/templates/validate-code-task.md

### Files Modified
- plugins/forge/skills/eval/SKILL.md

### Key Decisions
- Rubric uses 3 dimensions (Scenario Traceability 400pts, Path Completeness 350pts, Code-PRD Consistency 250pts) to focus on traceability rather than document quality
- Task templates use type=gate and breaking=false per Hard Rules — produces a verification report, not a blocking gate
- Quick mode template references proposal as primary input (no PRD guaranteed), breakdown-tasks template references PRD

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] validate-code.md rubric exists with scale: 1000, target: 700, iterations: 1, type: validate-code
- [x] Rubric defines dimensions for: scenario traceability, path completeness, code-prd consistency
- [x] Eval SKILL.md Prerequisites table includes validate-code entry requiring PRD + git diff
- [x] Eval SKILL.md Pre-Processing table includes validate-code entry describing PRD + git diff + code file list assembly
- [x] Eval SKILL.md Parameters table lists validate-code as a valid type
- [x] Eval SKILL.md Rubric Reference table includes validate-code row with scale/target/iterations
- [x] Task templates exist in both breakdown-tasks and quick-tasks with correct positioning
- [x] validate-code uses iterations=1, never triggers revise loop

## Notes
Documentation-only task. No code compilation or tests applicable. Did not modify doc-scorer.md or doc-reviser.md per Hard Rules.
