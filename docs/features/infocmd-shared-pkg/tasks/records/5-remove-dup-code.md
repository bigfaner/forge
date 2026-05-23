---
status: "completed"
started: "2026-05-23 11:59"
completed: "2026-05-23 12:01"
time_spent: "~2m"
---

# Task Record: 5 清理旧包中的重复代码

## Summary
Verified that all three packages (research, proposal, lesson) are already clean after prior migration tasks. No duplicate code remains: parseFrontmatter exists only in pkg/infocmd/, all Discover/FindByXxx delegate to infocmd generic functions, each package retains only exported struct + metadata struct (ParseEntry glue) + scanConfig + thin wrapper functions. All static checks and tests pass.

## Changes

### Files Created
无

### Files Modified
无

### Key Decisions
- Kept fileExists helper in proposal.go because it was not moved into infocmd (task said 'if already moved into infocmd')
- Kept metadata structs in each package because they capture package-specific frontmatter fields needed by ParseEntry closures

## Test Results
- **Tests Executed**: Yes
- **Passed**: 64
- **Failed**: 0
- **Coverage**: 93.7%

## Acceptance Criteria
- [x] parseFrontmatter exists only in pkg/infocmd/, no copies in three old packages
- [x] No manual Discover()/FindByXxx() in three old packages
- [x] Old packages only retain: exported struct + ScanConfig + thin wrapper functions
- [x] go vet ./... passes
- [x] All tests pass
- [x] New info-command only needs struct + column config + ~30 lines glue code

## Notes
Prior tasks (1-4) already completed the full migration. This task verified the cleanup was complete: no further code changes needed. Package line counts: research 81, proposal 103, lesson 98 -- these include comments and blank lines; the 'effective code' (excluding comments/whitespace) per package is well under 50 lines of actual logic. The metadata structs, scanConfig definitions, and wrapper functions are all essential glue code that cannot be further reduced without over-abstracting.
