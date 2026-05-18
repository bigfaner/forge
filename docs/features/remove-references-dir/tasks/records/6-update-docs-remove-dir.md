---
status: "completed"
started: "2026-05-19 01:31"
completed: "2026-05-19 01:38"
time_spent: "~7m"
---

# Task Record: 6 Update forge-distribution.md and remove references/ directory

## Summary
Updated forge-distribution.md to remove all 4 references/ locations (directory tree, component table, path rules table, reference file explanation paragraph) and deleted the entire plugins/forge/references/ directory (8 files). Also fixed pre-existing argument-hint vs argument-hints frontmatter mismatch in extract-design-md command. Verified zero references/shared/ occurrences remain in plugins/forge/.

## Changes

### Files Created
无

### Files Modified
- docs/conventions/forge-distribution.md
- plugins/forge/commands/extract-design-md.md

### Key Decisions
- Preserved docs/reference/ (user project directory) as documented in Hard Rules -- that is a separate concept from plugins/forge/references/
- Fixed pre-existing argument-hint/argument-hints frontmatter mismatch in extract-design-md to unblock quality gate

## Test Results
- **Tests Executed**: Yes
- **Passed**: 24
- **Failed**: 0
- **Coverage**: 0.0%

## Acceptance Criteria
- [x] plugins/forge/references/ directory no longer exists
- [x] Zero occurrences of references/shared/ in any file under plugins/forge/
- [x] docs/conventions/forge-distribution.md has no mention of references/ directory
- [x] docs/conventions/forge-distribution.md directory tree and component table are accurate without references/
- [x] Path resolution section no longer lists references/shared/ cross-skill reference pattern

## Notes
无
