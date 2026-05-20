---
status: "completed"
started: "2026-05-20 21:13"
completed: "2026-05-20 21:20"
time_spent: "~7m"
---

# Task Record: 1 Add RunTasks ModeToggle to AutoConfig

## Summary
Add RunTasks ModeToggle field to AutoConfig struct with default Quick:true Full:false, following existing pattern (CleanCode/Validation). Updated all related functions: AutoConfigDefaults, IsZero, WithDefaults, parseAutoRaw, applyDefaults, getAutoKeyValue.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/forgeconfig/config_test.go

### Key Decisions
- Reused existing ModeToggle pattern, no new config structures introduced
- RunTasks placed before GitPush in struct to group ModeToggle fields together
- getAutoKeyValue returns runTasks as 'quick:<bool> full:<bool>' format for CLI display

## Test Results
- **Tests Executed**: Yes
- **Passed**: 73
- **Failed**: 0
- **Coverage**: 88.2%

## Acceptance Criteria
- [x] AutoConfig struct has RunTasks ModeToggle field with yaml:"runTasks" tag
- [x] AutoConfigDefaults() sets RunTasks: ModeToggle{Quick: true, Full: false}
- [x] IsZero() checks include RunTasks field
- [x] WithDefaults() handles RunTasks zero-value filling
- [x] applyDefaults() handles RunTasks raw/default value logic
- [x] forge config get auto.runTasks returns correct default quick:true full:false
- [x] Backward compatible: unconfigured runTasks uses default values
- [x] Existing tests pass, new field has test coverage

## Notes
无
