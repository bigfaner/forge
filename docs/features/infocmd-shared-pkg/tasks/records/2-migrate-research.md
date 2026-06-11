---
status: "completed"
started: "2026-05-23 11:41"
completed: "2026-05-23 11:47"
time_spent: "~6m"
---

# Task Record: 2 迁移 research 命令使用 pkg/infocmd/

## Summary
Migrated research command to use pkg/infocmd/ shared package. Replaced manual Discover/FindBySlug/parseFrontmatter implementations with infocmd.Discover and infocmd.FindByID delegations. Preserved all external behavior including error message format, sorting logic, and frontmatter filtering (skip files without Topic or Mode).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/research/research.go

### Key Decisions
- Used package-level scanConfig variable to define infocmd.ScanConfig, keeping Discover/FindBySlug as thin wrappers
- Preserved Topic+Mode empty filter in ParseEntry callback to match original behavior of skipping files without meaningful frontmatter
- FindBySlug wraps infocmd.FindByID error to maintain 'research report not found: X' error message format for backward compatibility
- Unused parameters in ParseEntry lambda renamed to _ to pass lint (path, modTime not needed for research's flat file mode)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 19
- **Failed**: 0
- **Coverage**: 100.0%

## Acceptance Criteria
- [x] pkg/research/research.go uses infocmd.Discover and infocmd.FindByID
- [x] internal/cmd/research.go list/detail rendering preserves original output format
- [x] forge research command output byte-identical to pre-refactor
- [x] All existing research tests pass
- [x] No parseFrontmatter copy in pkg/research/

## Notes
research.go reduced from 163 lines to 82 lines. All 19 tests pass with 100% coverage. No changes needed in internal/cmd/research.go since Discover/FindBySlug signatures were preserved.
