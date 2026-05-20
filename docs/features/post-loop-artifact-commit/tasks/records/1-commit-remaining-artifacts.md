---
status: "completed"
started: "2026-05-20 18:07"
completed: "2026-05-20 18:08"
time_spent: "~1m"
---

# Task Record: 1 Add post-loop artifact commit step to run-tasks

## Summary
Added 'Commit Remaining Artifacts' section to run-tasks Post-Completion flow after Knowledge Review. The section runs git status --porcelain, filters to feature-scope paths (docs/features/<slug>/, docs/decisions/, docs/lessons/, docs/conventions/, docs/business-rules/), and commits matched artifacts with conventional commit message. Skips silently when nothing to commit. Guards against 3 consecutive failures.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/run-tasks.md

### Key Decisions
- Placed after Knowledge Review Step 6 (Write confirmed knowledge) as Step 7, before end of Post-Completion section
- Used explicit path filter with 5 directories to prevent committing unrelated files
- Slug variable reused from Step 0 resolution, no new variable needed

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] After Knowledge Review section (Step 6: Write confirmed knowledge), a new 'Commit Remaining Artifacts' section exists
- [x] The section unconditionally runs git status --porcelain to detect uncommitted files
- [x] Filters results to feature-scope paths only: docs/features/<slug>/ and knowledge directories
- [x] If filtered results are non-empty: git add the matched paths + git commit with message chore(<slug>): commit post-loop artifacts
- [x] If filtered results are empty: skip silently, no output
- [x] Does NOT run when loop ended due to 3 consecutive failures (matches Knowledge Review guard)

## Notes
无
