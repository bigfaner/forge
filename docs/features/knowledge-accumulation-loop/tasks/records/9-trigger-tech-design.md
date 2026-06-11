---
status: "completed"
started: "2026-05-17 22:09"
completed: "2026-05-17 22:18"
time_spent: "~9m"
---

# Task Record: 9 Add auto-extract trigger to tech-design

## Summary
Added Step 11 (Auto-Extract Knowledge) to tech-design skill. After design completion, reads shared knowledge-extraction.md routine by reference to scan design document for lessons, conventions, business rules, and non-archived decisions. Complements existing Step 7 decision archiving without duplicating it.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/tech-design/SKILL.md

### Key Decisions
- Added as Step 11 (after Step 10 eval prompt) rather than inside Step 7 to avoid modifying existing archiving flow
- References shared extraction routine by path, not by copying content
- Explicitly coordinates with Step 7 via a Coordination section to prevent duplicate decision entries

## Test Results
- **Tests Executed**: Yes
- **Passed**: 1033
- **Failed**: 0
- **Coverage**: 83.4%

## Acceptance Criteria
- [x] After tech-design completion, a knowledge review step runs
- [x] Step reads plugins/forge/references/shared/knowledge-extraction.md for extraction logic
- [x] Scans design document for architecture decisions, dependency choices, data model decisions
- [x] Silent when design contains no notable architectural knowledge
- [x] Presents extracted knowledge via AskUserQuestion for user confirmation
- [x] Writes confirmed knowledge to appropriate directories using shared formats

## Notes
Markdown-only change to SKILL.md. No code changes. Step 7 (decision archiving) remains untouched. Process flow diagram updated to include Step 11.
