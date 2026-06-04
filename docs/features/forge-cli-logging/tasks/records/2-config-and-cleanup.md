---
status: "completed"
started: "2026-06-05 07:10"
completed: "2026-06-05 07:22"
time_spent: "~12m"
---

# Task Record: 2 Add config extension, validation, and auto-cleanup

## Summary
Added LogsConfig struct to forgeconfig.Config with omitempty, ResolveLogsConfig with safe defaults, auto-cleanup of old log files via forgelog.Init(), and ForgeLogsDir constant. LogsConfig.Enabled uses *bool to distinguish absent (default true) from explicit false.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/forgeconfig/config.go
- forge-cli/pkg/forgeconfig/config_test.go
- forge-cli/pkg/forgelog/forgelog.go
- forge-cli/pkg/forgelog/forgelog_test.go

### Key Decisions
- Used *bool for Enabled field to distinguish YAML absent (nil -> true) from explicit false
- Moved LogsConfig from forgelog to forgeconfig as the authoritative config location
- ResolveLogsConfig handles validation (level defaults, retentionDays >= 1) so forgelog.Init stays simple
- ForgeLogsDir constant already existed in constants.go from task 1

## Test Results
- **Tests Executed**: Yes
- **Passed**: 119
- **Failed**: 0
- **Coverage**: 87.7%

## Acceptance Criteria
- [x] AC-1: Auto-cleanup deletes log files older than retentionDays after new log opened; active log never deleted
- [x] AC-2: Config missing/malformed falls back to defaults (level=info, retentionDays=7, enabled=true); invalid level->info; retentionDays<1->7
- [x] AC-3: FORGE_NO_LOG=1 or logs.enabled:false skips FileBackend; env var takes precedence

## Notes
Coverage: forgeconfig 84.7%, forgelog 87.7%. Hard rules verified: omitempty on Logs field, cleanup errors silently ignored.
