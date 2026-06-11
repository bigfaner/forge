---
status: "completed"
started: "2026-05-23 11:48"
completed: "2026-05-23 11:55"
time_spent: "~7m"
---

# Task Record: 3 迁移 proposal 命令使用 pkg/infocmd/

## Summary
Migrated pkg/proposal to use infocmd.Discover and infocmd.FindByID, removing the duplicate parseFrontmatter and manual directory scanning. Sorting moved from cmd layer into infocmd.Discover (already built-in). Cross-reference logic (PRD check, Feature status) preserved via ParseEntry callback. All tests pass, coverage 96.0%.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/proposal/proposal.go
- forge-cli/internal/cmd/proposal.go

### Key Decisions
- Wrapped infocmd.FindByID error in proposal.FindBySlug to preserve 'proposal not found' error message for backward compatibility
- Moved sort into infocmd.Discover (already built-in Created desc + mtime fallback), removed sort.Slice from cmd layer
- Kept fileExists as private helper in proposal package since it's a minimal one-liner used for PRD cross-reference
- scanConfig takes projectRoot as parameter to make ParseEntry closure capture the project root for cross-reference paths

## Test Results
- **Tests Executed**: Yes
- **Passed**: 699
- **Failed**: 0
- **Coverage**: 96.0%

## Acceptance Criteria
- [x] pkg/proposal/proposal.go uses infocmd.Discover and infocmd.FindByID instead of manual implementation
- [x] Proposal detail PRD (yes/no) and Feature status fields preserved
- [x] Proposal list sorting (Created descending) unchanged
- [x] forge proposal command output byte-identical to pre-refactor
- [x] All existing proposal-related tests pass
- [x] No standalone parseFrontmatter copy in pkg/proposal/

## Notes
worktree/worktree.go uses proposal.Discover via the public API which is preserved, so no changes were needed there.
