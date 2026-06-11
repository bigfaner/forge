---
status: "completed"
started: "2026-05-17 00:29"
completed: "2026-05-17 00:39"
time_spent: "~10m"
---

# Task Record: 2 Update Go CLI display for Approved/Completed status

## Summary
Verified and added test coverage for Approved/Completed status display in forge proposal list, proposal detail, feature status, and feature list commands. The CLI already displays raw status strings correctly; added 9 new tests confirming Approved and Completed statuses render properly at all display points.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/proposal_test.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/pkg/proposal/proposal_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- No code changes needed -- existing CLI reads status as raw string from frontmatter and displays as-is, which already handles Approved and Completed correctly
- Tests verify both proposal-level (list + detail) and feature-level (status + list) display paths
- Version bumped from 3.17.0 to 3.17.1 (patch) per hard rules for display enhancement

## Test Results
- **Tests Executed**: Yes
- **Passed**: 461
- **Failed**: 0
- **Coverage**: 80.8%

## Acceptance Criteria
- [x] forge proposal list displays proposals with status: Approved showing Approved in the STATUS column
- [x] forge proposal list displays proposals with status: Completed showing Completed in the STATUS column
- [x] forge feature status <slug> correctly reflects when a feature's manifest status is completed
- [x] Existing tests pass (go test ./...)
- [x] New test cases cover Approved and Completed status display

## Notes
无
