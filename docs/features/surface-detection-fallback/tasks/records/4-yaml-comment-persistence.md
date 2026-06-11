---
status: "completed"
started: "2026-05-24 21:49"
completed: "2026-05-24 21:54"
time_spent: "~5m"
---

# Task Record: 4 Persist source annotation as YAML comment in config

## Summary
Implemented YAML comment persistence for surface source annotations using yaml.Node round-trip API

## Changes

### Files Created
- forge-cli/pkg/forgeconfig/yaml_comment.go
- forge-cli/pkg/forgeconfig/yaml_comment_test.go

### Files Modified
无

### Key Decisions
- Used yaml.Node round-trip API (Encode + Marshal) per Hard Rule, not string concatenation or regex injection
- Added WriteConfigWithSources as a separate function rather than modifying writeConfig, to avoid breaking existing callers
- Comment format uses LineComment on the surfaces value node: '# source: inference:cmd-dir'
- ReadSurfaceComment extracts comment via yaml.Node unmarshal for TUI re-detection context display

## Test Results
- **Tests Executed**: Yes
- **Passed**: 11
- **Failed**: 0
- **Coverage**: 87.9%

## Acceptance Criteria
- [x] After forge init infers cli from cmd/ directory, config file contains surfaces: cli # source: inference:cmd-dir
- [x] Comment is present in file but ignored by YAML unmarshaler (no schema change to SurfacesMap)
- [x] On re-detection, existing comment is read and displayed for context in the TUI prompt
- [x] If comment is stripped during round-trip, detection still works correctly
- [x] Integration test reads config file after write and asserts comment is present

## Notes
Comment appended as LineComment on value node. Both scalar and map forms supported. Full package coverage 87.9% exceeds 80% target.
