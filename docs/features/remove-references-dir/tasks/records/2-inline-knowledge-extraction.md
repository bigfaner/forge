---
status: "completed"
started: "2026-05-19 01:20"
completed: "2026-05-19 01:23"
time_spent: "~3m"
---

# Task Record: 2 Inline knowledge-extraction.md into consuming skills and commands

## Summary
Inlined full knowledge-extraction.md routine into write-prd/SKILL.md, tech-design/SKILL.md, run-tasks.md, and fix-bug.md. Each file now contains the complete extraction flow (6 steps), knowledge types table, notable knowledge heuristics, deduplication logic, and rules — with trigger-specific parameters preserved. All references/shared/knowledge-extraction path references removed.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md
- plugins/forge/skills/tech-design/SKILL.md
- plugins/forge/commands/run-tasks.md
- plugins/forge/commands/fix-bug.md

### Key Decisions
- Preserved trigger-specific parameter context (trigger name, artifacts list, scanning scope notes) around the inlined content rather than using a generic template
- Kept consumer-specific sections like tech-design's Coordination with Step 7 and run-tasks's Do NOT run on 3 consecutive failures guard

## Test Results
- **Tests Executed**: Yes
- **Passed**: 0
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] No occurrence of references/shared/knowledge-extraction in any of the 4 modified files
- [x] Each file contains the full knowledge extraction routine verbatim
- [x] Path references using ${CLAUDE_SKILL_DIR} to knowledge-extraction.md are fully removed

## Notes
Documentation-only task. No test execution required.
