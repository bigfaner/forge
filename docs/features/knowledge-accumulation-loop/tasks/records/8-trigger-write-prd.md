---
status: "completed"
started: "2026-05-17 22:04"
completed: "2026-05-17 22:08"
time_spent: "~4m"
---

# Task Record: 8 Add auto-extract trigger to write-prd

## Summary
Added auto-extract knowledge trigger to write-prd SKILL.md as Step 12 (Knowledge Review), referencing the shared extraction routine by path rather than copying content

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/write-prd/SKILL.md

### Key Decisions
- Placed knowledge review as Step 12 after Step 11 (Adversarial Eval Prompt), consistent with the post-completion trigger pattern used by run-tasks and fix-bug

## Test Results
- **Tests Executed**: No
- **Passed**: 20
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] After PRD completion, a knowledge review step runs
- [x] Step reads plugins/forge/references/shared/knowledge-extraction.md for extraction logic
- [x] Scans PRD content for new business rules and user-facing constraints
- [x] Silent when PRD contains no cross-cutting knowledge
- [x] Presents extracted knowledge via AskUserQuestion for user confirmation
- [x] Writes confirmed knowledge to appropriate directories using shared formats

## Notes
Prompt-only change (SKILL.md). No unit tests applicable. Coverage set to -1.0 since this is a prompt file modification with no testable code paths.
