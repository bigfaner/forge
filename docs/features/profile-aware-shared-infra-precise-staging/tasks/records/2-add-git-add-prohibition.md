---
status: "completed"
started: "2026-05-17 11:07"
completed: "2026-05-17 11:11"
time_spent: "~4m"
---

# Task Record: 2 Add git add -A prohibition to git-commit.md

## Summary
Added HARD-RULE to git-commit.md prohibiting git add -A, git add ., and git add --all. The rule requires explicit file paths in every git add command and includes a 'Why' section referencing the historical d7f8a13 incident (169-file commit from a 2-file fix). Step 2 in the Steps section was also updated to reinforce explicit file listing.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/commands/git-commit.md

### Key Decisions
- Placed HARD-RULE block inside the Steps section (after step 4) rather than before it, so the agent reads the workflow first and encounters the prohibition at the staging step boundary
- Updated Step 2 text from 'Stage appropriate files with git add' to 'Stage appropriate files with git add — listing each file path explicitly' for reinforcement at the point of action

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] git-commit.md contains a new section with clear prohibition against git add -A, git add ., and git add --all
- [x] The prohibition requires listing explicit file paths in the git add command
- [x] The existing commit workflow (Steps 1-4, Task Completion Template) is otherwise unchanged

## Notes
Documentation-only change. No code or test files modified. Task type is 'documentation'.
