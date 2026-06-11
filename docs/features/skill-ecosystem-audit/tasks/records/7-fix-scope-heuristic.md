---
status: "completed"
started: "2026-05-18 01:55"
completed: "2026-05-18 02:00"
time_spent: "~5m"
---

# Task Record: 7 Fix breakdown-tasks scope heuristic for backend projects

## Summary
Fixed scope heuristic in breakdown-tasks SKILL.md: moved src/ out of unconditional frontend classification and added a conditional special case that checks for go.mod/Cargo.toml (backend), package.json (frontend), both (undetermined), or neither (undetermined)

## Changes

### Files Created
无

### Files Modified
- plugins/forge/skills/breakdown-tasks/SKILL.md

### Key Decisions
- Added src/ as an explicit special case block rather than patching the frontend/backend bullet lists, to make the three-way conditional logic clear and unambiguous

## Test Results
- **Tests Executed**: No
- **Passed**: 20
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] src/ is classified as backend when go.mod or Cargo.toml exists at the same level (without requiring package.json to be absent)
- [x] src/ remains frontend for Node.js/React projects (when package.json exists without go.mod/Cargo.toml)
- [x] The heuristic handles mixed projects (both package.json and go.mod) correctly -> scope: 'all'

## Notes
Markdown-only change to skill prompt template. No executable tests apply to prompt text. Project test suite (20 packages) all pass. Coverage set to -1.0 since no test coverage applies to markdown skill files.
