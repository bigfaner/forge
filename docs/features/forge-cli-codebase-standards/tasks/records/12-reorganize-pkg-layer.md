---
status: "completed"
started: "2026-05-31 10:16"
completed: "2026-05-31 10:27"
time_spent: "~11m"
---

# Task Record: 12 Reorganize pkg/ layer

## Summary
Reorganized pkg/ layer from 17 packages to 14 by merging: pkg/version/ into pkg/types/ (version info now in types), pkg/lesson/ into pkg/infocmd/ (lesson discovery), pkg/research/ into pkg/infocmd/ (research discovery). Added doc.go to merged packages listing sub-domains and responsibility boundaries. Updated Makefile ldflags, install script, .golangci.yml exclusions, and package-organization.md to reflect new structure.

## Changes

### Files Created
- forge-cli/pkg/types/doc.go
- forge-cli/pkg/types/version.go
- forge-cli/pkg/infocmd/doc.go
- forge-cli/pkg/infocmd/lesson.go
- forge-cli/pkg/infocmd/lesson_test.go
- forge-cli/pkg/infocmd/research.go
- forge-cli/pkg/infocmd/research_test.go

### Files Modified
- forge-cli/Makefile
- forge-cli/scripts/install-local.sh
- forge-cli/.golangci.yml
- forge-cli/internal/cmd/version.go
- forge-cli/internal/cmd/lesson.go
- forge-cli/internal/cmd/research.go
- forge-cli/pkg/types/types_test.go
- docs/conventions/package-organization.md

### Key Decisions
- Merged pkg/version/ into pkg/types/ because both are leaf packages with zero internal dependencies and version is just constants/getters
- Merged pkg/lesson/ and pkg/research/ into pkg/infocmd/ because both only depend on infocmd's ScanConfig framework and are small (2 files each)
- Renamed Discover->DiscoverLessons/DiscoverReports and FindByName/FindBySlug->FindLessonByName/FindReportBySlug to avoid name collisions in merged package
- Kept pkg/project/ separate due to 31 cmd-layer consumers making merge cost disproportionate to benefit
- Execution order followed Hard Rules: leaf packages first (version), then domain packages (lesson, research)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 25
- **Failed**: 0
- **Coverage**: 85.2%

## Acceptance Criteria
- [x] pkg/ layer has at most 14 packages (SC-9)
- [x] Each merged package contains doc.go listing sub-domains and responsibility boundaries
- [x] No circular dependencies: go build ./... passes after merge
- [x] go build ./... and go test ./... all pass (SC-11)
- [x] Large file split feasibility assessment produced and recorded (SC-13)

## Notes
Task 7 audit confirmed no cross-module dependencies, so full execution was possible. Large file assessment: pkg/task/build_test.go (2784 lines), pkg/task/autogen_test.go (2675 lines), pkg/forgeconfig/config.go (1272 lines), and pkg/task/pipeline.go (1102 lines) are the largest files that may benefit from future splitting.
