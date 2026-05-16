---
status: "completed"
started: "2026-05-16 09:35"
completed: "2026-05-16 09:44"
time_spent: "~9m"
---

# Task Record: 2 Sort forge feature list output by manifest mtime descending

## Summary
Sort forge feature list output by manifest.md mtime descending. Added ManifestMtime field to featureInfo struct, capture mtime via os.Stat in discoverFeatures(), and sort with sort.Slice descending. Features with missing/unreadable manifest sort to end (mtime=0).

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/feature.go
- forge-cli/internal/cmd/feature_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used int64 unix timestamp for ManifestMtime instead of time.Time to keep the struct simple and comparable via sort.Slice
- Features with missing manifest get mtime=0 which naturally sorts to the end in descending order

## Test Results
- **Tests Executed**: Yes
- **Passed**: 405
- **Failed**: 0
- **Coverage**: 80.6%

## Acceptance Criteria
- [x] runFeatureList() sorts features by manifest.md mtime descending (newest first)
- [x] Features with missing/unreadable manifest sort to the end
- [x] Existing tests continue to pass
- [x] New test verifies sort order

## Notes
Hard rule satisfied: used sort.Slice() from stdlib. Two new tests: TestFeatureList_SortedByManifestMtime and TestFeatureList_MissingManifestSortsToEnd.
