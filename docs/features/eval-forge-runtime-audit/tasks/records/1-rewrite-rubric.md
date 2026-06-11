---
status: "completed"
started: "2026-05-15 21:43"
completed: "2026-05-15 21:46"
time_spent: "~3m"
---

# Task Record: 1 Rewrite rubric.md: 6-dimension runtime reliability scoring

## Summary
Complete rewrite of eval-forge rubric from 12-dimension structural consistency (1000 pts) to 6-dimension runtime reliability (1000 pts). New rubric covers: D1 Workflow Completeness (250), D2 Bypass Resistance (250), D3 Instruction Precision (200), D4 Cross-file Dedup (150), D5 Reference Integrity (100), D6 Structural Convention (50). Embeds ground-truth workflow specs, 14 known bypass vectors, CLI-filled variables from prompt.go, and 5 known redundancy instances from the proposal.

## Changes

### Files Created
无

### Files Modified
- .claude/skills/eval-forge/templates/rubric.md

### Key Decisions
- Followed proposal as authoritative source for all dimension specs, scoring criteria, and embedded data
- Deduction tiers updated from flat -5/-15/-25 to graduated -5/-10/-15/-20 (Low/Medium/High/Critical)
- Included 4-phase methodology reference in rubric header for scorer guidance

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] Rubric defines exactly 6 dimensions: Workflow Completeness (250), Bypass Resistance (250), Instruction Precision (200), Cross-file Dedup (150), Reference Integrity (100), Structural Convention (50)
- [x] Each dimension has detailed scoring criteria matching proposal tables (1a-1e, 2a-2e, 3a-3d, 4a-4c, 5a-5d, 6a-6c)
- [x] Dimension 1 embeds the ground-truth workflow specs (Full Mode, Quick Mode, Manifest Status Machine, Per-Skill Precondition/Output Matrix)
- [x] Dimension 2 lists the 5 bypass types with point allocations and the 14 known bypass vectors (BV-2.1 through BV-5.2)
- [x] Dimension 3 lists CLI-filled variables from prompt.go (not marked as undefined)
- [x] Dimension 4 references known redundancy instances (3 categories: content copy, guide overlap, unreasonable inline)
- [x] Deduction tiers updated to reflect new severity model (breakpoints = -20/-15/-10/-5 instead of flat -5/-15/-25)
- [x] Point totals sum to 1000
- [x] Report template reference updated to .claude/skills/eval-forge/templates/report.md

## Notes
Documentation-only task (noTest). This is task 1 of the eval-forge-runtime-audit feature. Only rubric.md was modified; report.md was not touched per hard rules (Task 4 handles that).
