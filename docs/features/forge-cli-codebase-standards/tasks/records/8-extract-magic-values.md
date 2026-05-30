---
status: "completed"
started: "2026-05-30 22:52"
completed: "2026-05-30 23:03"
time_spent: "~11m"
---

# Task Record: 8 Extract magic values to named constants

## Summary
Extracted all magic values to named constants: path constants (tests/results/raw-output.txt, unit-raw-output.txt), color constants (#7DCFFF, #FF8700, #9ECE6A), sentinel constants (99999), retry parameters, timeout values, and converted octal permissions to 0o prefix format.

## Changes

### Files Created
- forge-cli/internal/cmd/styles.go
- forge-cli/internal/cmd/constants.go
- forge-cli/pkg/serverprobe/constants.go

### Files Modified
- forge-cli/pkg/feature/constants.go
- forge-cli/internal/cmd/quality_gate.go
- forge-cli/internal/cmd/init.go
- forge-cli/internal/cmd/init_surfaces.go
- forge-cli/internal/cmd/task/list.go
- forge-cli/internal/cmd/task/claim.go
- forge-cli/internal/cmd/task/index.go
- forge-cli/internal/cmd/task/submit.go
- forge-cli/internal/cmd/feature/feature_complete.go
- forge-cli/internal/cmd/forensic/extract.go
- forge-cli/pkg/testrunner/test_results.go
- forge-cli/pkg/task/index.go
- forge-cli/pkg/task/add.go
- forge-cli/pkg/task/frontmatter.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/stage_gates.go
- forge-cli/pkg/task/state.go
- forge-cli/pkg/index/lock.go
- forge-cli/pkg/index/atomic.go
- forge-cli/pkg/feature/forge_state.go
- forge-cli/pkg/feature/feature.go
- forge-cli/pkg/serverprobe/serverprobe.go

### Key Decisions
- Path constants (TestOutputFileName, UnitTestOutputFileName) added to pkg/feature/constants.go as shared constants
- Color constants centralized in internal/cmd/styles.go (colorModeHighlight, colorConflict, colorSource)
- ANSI escape codes placed in internal/cmd/task/list.go since task package cannot import cmd package
- Retry/tuning parameters in internal/cmd/constants.go (maxProbeRetries, probeRetryInterval, conciseErrorMaxLines, maxSourceFiles)
- Probe timeout in pkg/serverprobe/constants.go (defaultProbeTimeout)
- Lock retry backoff in pkg/index/lock.go (lockRetryBackoff)
- Sentinel constants (fallbackSortPriority, unreachableDepth) defined locally in their respective files with doc comments
- Octal permissions standardized to 0o prefix across all production files

## Test Results
- **Tests Executed**: Yes
- **Passed**: 3901
- **Failed**: 0
- **Coverage**: 74.6%

## Acceptance Criteria
- [x] SC-1: grep tests/results/ path literals returns zero in production code
- [x] SC-2: grep lipgloss.Color('# hex literals returns zero
- [x] SC-3: grep 99999 sentinel values returns only constant definitions
- [x] SC-4: grep 0644|0755 old octal returns zero
- [x] SC-11: go build and go test all pass

## Notes
SC-3 note: 99999 still appears in constant definitions (fallbackSortPriority, unreachableDepth) which is the spec-prescribed target state. The values are no longer used inline. All tests pass with zero failures across all packages.
