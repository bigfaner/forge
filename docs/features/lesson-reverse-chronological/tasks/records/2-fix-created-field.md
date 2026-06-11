---
status: "completed"
started: "2026-05-19 11:44"
completed: "2026-05-19 11:52"
time_spent: "~8m"
---

# Task Record: 2 Fix Go parser to support created frontmatter field

## Summary
Added `created` field support to the Metadata struct so the Go parser now recognizes both `created` and `date` frontmatter fields. The `created` field takes priority over `date` when both are present, and file modification time remains the final fallback.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/lesson/lesson.go
- forge-cli/pkg/lesson/lesson_test.go

### Key Decisions
- Added `Created string yaml:"created"` field to Metadata struct alongside existing `Date string yaml:"date"` field, rather than using a custom UnmarshalYAML method — simpler and sufficient for two fields
- In Discover(), check `meta.Created` first, then fall back to `meta.Date`, then to file modification time — preserves backward compatibility

## Test Results
- **Tests Executed**: Yes
- **Passed**: 19
- **Failed**: 0
- **Coverage**: 86.7%

## Acceptance Criteria
- [x] Metadata struct parses both created and date frontmatter fields
- [x] created field takes priority over date when both are present
- [x] Existing tests continue to pass
- [x] New parsing logic has unit test coverage

## Notes
无
