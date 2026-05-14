---
status: "completed"
started: "2026-05-14 17:13"
completed: "2026-05-14 17:27"
time_spent: "~14m"
---

# Task Record: fix-2 Fix: 14 e2e test failures — subcommands not wired up (config, proposal, lesson, feature list/status)

## Summary
Fix 14 e2e test failures: config get now auto-detects project root and silences error output for missing keys; config.yaml populated with project-type and capabilities; forge binary rebuilt with all new commands (config, proposal, lesson, feature list/status)

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/config.go
- forge-cli/scripts/version.txt
- .forge/config.yaml

### Key Decisions
- Added SilenceErrors+SilenceUsage on configGetCmd to produce empty CombinedOutput for TC-006 (missing key should produce no output)
- Changed resolveProjectRoot to auto-detect via project.FindProjectRoot() instead of defaulting to '.' which failed when forge binary was invoked from forge-cli/ subdirectory
- Populated .forge/config.yaml with project-type: backend and capabilities: [compile, test, lint] to satisfy TC-004 and TC-005

## Test Results
- **Tests Executed**: Yes
- **Passed**: 723
- **Failed**: 0
- **Coverage**: 80.8%

## Acceptance Criteria
- [x] forge config get project-type returns 'backend' with exit 0
- [x] forge config get capabilities returns array items with exit 0
- [x] forge config get nonexistent-key exits 1 with no output
- [x] forge proposal lists proposals with correct table columns
- [x] forge proposal <slug> shows detail view
- [x] forge lesson lists lessons with correct table columns
- [x] forge lesson <name> shows detail view
- [x] forge feature list shows all features with table columns
- [x] forge feature status <slug> shows detailed status
- [x] forge feature (no args) shows current feature

## Notes
Root cause was two-fold: (1) installed forge binary at ~/.forge/bin/forge was stale and did not include config/proposal/lesson commands, (2) config get used resolveProjectRoot defaulting to '.' instead of auto-detecting project root, and missing key errors were printed via cobra's default error handler to CombinedOutput. Fixed by rebuilding binary, auto-detecting project root, and silencing cobra error output on configGetCmd.
