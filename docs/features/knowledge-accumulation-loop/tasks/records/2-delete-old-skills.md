---
status: "completed"
started: "2026-05-17 21:24"
completed: "2026-05-17 21:29"
time_spent: "~5m"
---

# Task Record: 2 Delete /record-decision and /learn-lesson skills

## Summary
Removed the learn-lesson skill directory (SKILL.md, templates/template.md, examples/debug-race-condition.md) and updated stale references in forensic/SKILL.md and knowledge-extraction.md to point to the unified /learn skill and its template paths.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/forensic/SKILL.md
- plugins/forge/references/shared/knowledge-extraction.md

### Key Decisions
- Updated knowledge-extraction.md template paths to use learn/templates/lesson-entry.md (verified actual filename) instead of the deleted learn-lesson/templates/template.md

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 83.4%

## Acceptance Criteria
- [x] plugins/forge/skills/learn-lesson/ directory is fully removed
- [x] plugins/forge/skills/record-decision/ directory does not exist (verify, not delete)
- [x] plugins/forge/references/shared/decision-logging.md is preserved unchanged
- [x] No other files reference /record-decision or /learn-lesson as active skills

## Notes
guide.md references to learn-lesson/record-decision are deferred to Task 5 per implementation notes. decision-logging.md and commands/record-decision.md reference record-decision as a command flow, not as a skill -- these are preserved. learn/SKILL.md mentions absorbing the old skills as historical context, which is fine.
