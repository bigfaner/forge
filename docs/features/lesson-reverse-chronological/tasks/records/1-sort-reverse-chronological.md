---
status: "completed"
started: "2026-05-19 11:35"
completed: "2026-05-19 11:42"
time_spent: "~7m"
---

# Task Record: 1 Sort lesson list by file modification time in reverse chronological order

## Summary
Added reverse chronological sorting by file modification time to the lesson Discover() function. Lessons now return newest-first, with zero mod-time lessons sorted to the end. Used sort.Slice with no new dependencies.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/lesson/lesson.go
- forge-cli/pkg/lesson/lesson_test.go

### Key Decisions
- Captured ModTime during the existing os.Stat call instead of a separate stat, eliminating redundant IO
- Used a local lessonWithTime struct to pair lessons with their mod times for sorting, then extracted the sorted Lesson slice -- avoids adding modTime field to the exported Lesson type
- Zero mod times sort to end via IsZero() check in the sort comparator

## Test Results
- **Tests Executed**: Yes
- **Passed**: 16
- **Failed**: 0
- **Coverage**: 86.2%

## Acceptance Criteria
- [x] forge lesson output is sorted by file modification time, newest first
- [x] Lessons without valid modification times sort to the end of the list
- [x] forge lesson <name> detail view still works correctly
- [x] New sorting logic has unit test coverage

## Notes
Two new tests added: TestDiscover_SortedByModTimeDescending (3 files with explicit mod times) and TestDiscover_OldestSortsLast (2 files). All 16 tests in the lesson package pass. Pre-existing failures in internal/cmd and pkg/just are unrelated.
