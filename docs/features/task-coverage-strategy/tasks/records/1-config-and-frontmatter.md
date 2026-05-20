---
status: "completed"
started: "2026-05-20 13:05"
completed: "2026-05-20 13:13"
time_spent: "~8m"
---

# Task Record: 1 Add coverage config schema, parsing, and frontmatter field

## Summary
Added coverage config schema (CoverageConfig, CoverageStrategy), parsing with built-in defaults, GetConfigValue support for coverage.* dot-notation, FrontmatterData.Coverage field (*int), Task.Coverage field propagation from frontmatter through BuildIndex, and updated forge-config.example.yaml with coverage documentation.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/forgeconfig/config_test.go
- forge-cli/pkg/task/frontmatter.go
- forge-cli/pkg/task/frontmatter_test.go
- forge-cli/pkg/task/types.go
- forge-cli/pkg/task/build.go
- forge-cli/pkg/task/build_test.go
- forge-cli/internal/cmd/testdata/forge-config.example.yaml

### Key Decisions
- Used map[string]CoverageStrategy with yaml:",inline" for extensible per-type coverage config (Hard Rule)
- CoverageStrategy.Percentage is *int (nil for maintain strategy, per Hard Rule)
- FrontmatterData.Coverage is *int to distinguish unset from zero value
- Task.Coverage is *int with json:",omitempty" so nil values are omitted from index.json
- CoverageConfigDefaults() returns fresh map each call to prevent mutation issues
- ReadCoverageConfig merges user config on top of defaults (partial override supported)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 236
- **Failed**: 0
- **Coverage**: 88.9%

## Acceptance Criteria
- [x] Config struct has Coverage *CoverageConfig field
- [x] CoverageConfig uses map structure with percentage and maintain strategies
- [x] Built-in defaults: coding.feature 80%, coding.enhancement/fix 60%, refactor/cleanup/clean maintain
- [x] Missing coverage config returns defaults without error
- [x] GetConfigValue supports coverage.* dot-notation queries
- [x] FrontmatterData has optional Coverage *int field
- [x] forge-config.example.yaml updated with coverage config example
- [x] Existing tests pass

## Notes
Coverage percentage is weighted average across forgeconfig (88.0%) and task (89.8%) packages.
