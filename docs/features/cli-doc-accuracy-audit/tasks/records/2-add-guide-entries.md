---
status: "completed"
started: "2026-06-07 22:03"
completed: "2026-06-07 22:05"
time_spent: "~2m"
---

# Task Record: 2 Add new commands and flags to guide.md

## Summary
Added 4 missing CLI entries to guide.md: forge task query (G5), forge task check-deps (G6), forge feature list (G7), and --tree flag for forge task list (G8). All descriptions verified against --help output.

## Changes

### Files Created
无

### Files Modified
- plugins/forge/hooks/guide.md

### Key Decisions
无

## Document Metrics
4 entries added, 1 entry updated, all verified against CLI --help output

## Referenced Documents
- docs/proposals/cli-doc-accuracy-audit/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] guide.md has forge task query <id-or-key> matching --help output
- [x] guide.md has forge task check-deps matching --help output
- [x] guide.md has forge feature list matching --help output
- [x] guide.md forge task list description includes --tree flag

## Notes
All new entries follow existing guide.md format (backtick command, em dash, description). Placement follows logical grouping: task commands in Task Management, feature list in Feature Management, check-deps in Pipeline Utilities.
