---
status: "completed"
started: "2026-05-16 09:27"
completed: "2026-05-16 09:34"
time_spent: "~7m"
---

# Task Record: 1 Sort forge proposal output by created date descending

## Summary
Added descending sort by Created date to runProposalList() using sort.Slice(). The Created field is stored as YYYY-MM-DD which sorts correctly lexicographically, so no time.Time parsing was needed. Also bumped version to 3.12.1 (patch).

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/proposal.go
- forge-cli/internal/cmd/proposal_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used lexicographic string comparison on YYYY-MM-DD format instead of parsing to time.Time, since the format sorts correctly and avoids the overhead of time parsing
- Placed sort in runProposalList() rather than Discover() to keep Discover() as a pure data retrieval function and sort only when displaying

## Test Results
- **Tests Executed**: Yes
- **Passed**: 6
- **Failed**: 0
- **Coverage**: 80.5%

## Acceptance Criteria
- [x] runProposalList() sorts proposals by Created date descending (newest first)
- [x] Proposals without created frontmatter still sort correctly (fallback mtime)
- [x] Existing tests continue to pass
- [x] New test verifies sort order

## Notes
The mtime fallback case is implicitly covered: Discover() already sets Created to ModTime formatted as YYYY-MM-DD, so the same lexicographic sort works for both frontmatter-sourced and mtime-fallback dates.
