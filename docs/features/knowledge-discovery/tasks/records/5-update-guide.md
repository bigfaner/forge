---
status: "completed"
started: "2026-05-17 13:47"
completed: "2026-05-17 13:50"
time_spent: "~3m"
---

# Task Record: 5 Update guide.md to reference domains frontmatter

## Summary
Updated the project knowledge note in guide.md (line 30) to mention that convention/business-rule files carry a `domains` frontmatter field with topic keywords, and that agents use it for relevance matching during task execution.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
- Kept the note to 2 sentences by integrating the domains explanation into the existing blockquote rather than adding a separate sentence

## Test Results
- **Tests Executed**: No
- **Passed**: 20
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] The note on line 30 mentions that convention/business-rule files have `domains` frontmatter
- [x] The note explains that domains are used for relevance matching during task execution
- [x] The note is concise (1-2 sentences) — guide.md is session-injected, every line costs tokens

## Notes
Documentation-only change. Coverage set to -1.0 as this task modifies only a markdown file.
