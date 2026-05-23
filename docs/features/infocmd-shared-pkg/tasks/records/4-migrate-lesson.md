---
status: "completed"
started: "2026-05-23 11:56"
completed: "2026-05-23 11:58"
time_spent: "~2m"
---

# Task Record: 4 迁移 lesson 命令使用 pkg/infocmd/

## Summary
Refactored pkg/lesson/lesson.go to use infocmd.Discover and infocmd.FindByID instead of manual directory scanning, frontmatter parsing, and sorting. Removed the standalone parseFrontmatter function (now uses infocmd.ParseFrontmatter). Preserved lesson-specific logic: Name-based ID, Category inference from filename prefix, Date->Created fallback, and original error message format.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/lesson/lesson.go

### Key Decisions
- Used infocmd.ScanConfig with flat mode (IsSubdir=false) for lesson's docs/lessons/*.md scanning pattern
- Wrapped FindByID error to preserve 'lesson not found' message (infocmd uses generic 'item not found')
- Kept inferCategory and categoryPrefixes as package-internal since they are lesson-specific logic
- ParseEntry closure handles the Date->Created fallback and Category inference, keeping the infocmd package generic

## Test Results
- **Tests Executed**: Yes
- **Passed**: 17
- **Failed**: 0
- **Coverage**: 94.4%

## Acceptance Criteria
- [x] pkg/lesson/lesson.go uses infocmd.Discover and infocmd.FindByID
- [x] lesson Name identifier correctly mapped (not Slug)
- [x] lesson Category inference preserved (from filename prefix)
- [x] lesson Date->Created fallback preserved
- [x] forge lesson command output byte-identical to pre-refactor
- [x] All existing lesson tests pass
- [x] No standalone parseFrontmatter copy in pkg/lesson/

## Notes
Coverage improved from 84.6% to 94.4% because the manual scan/sort code was replaced by infocmd delegation, reducing executable statements in lesson.go. The lesson.go file went from 180 lines to 97 lines.
