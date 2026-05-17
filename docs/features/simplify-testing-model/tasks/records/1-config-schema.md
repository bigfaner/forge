---
status: "completed"
started: "2026-05-17 13:46"
completed: "2026-05-17 13:57"
time_spent: "~11m"
---

# Task Record: 1 Config schema: replace TestProfiles/Capabilities with Interfaces/Languages

## Summary
Fix config_test.go in internal/cmd to use new schema field names (interfaces, languages) instead of old names (capabilities, test-profiles). The ForgeConfig struct, ReadLanguages(), ReadInterfaces(), and all embed.go exported symbols were already migrated in prior work. This task fixed the remaining test file references that used old YAML keys and CLI key names.

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/config_test.go

### Key Decisions
- The ForgeConfig struct, config.go, and embed.go were already migrated before this task started -- only test files needed updating
- Did not rename profiles/ directory or change embed paths per Hard Rules (those belong to Task 2)
- Did not change detection logic per Hard Rules (that belongs to Task 3)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 20
- **Failed**: 0
- **Coverage**: 79.0%

## Acceptance Criteria
- [x] ForgeConfig struct has Interfaces []string and Languages []string fields; TestProfiles and Capabilities removed
- [x] ReadLanguages() function exists
- [x] ReadInterfaces() function exists
- [x] Zero Go exported symbols contain capability or Capability
- [x] Config YAML field names: interfaces, languages
- [x] go build ./... passes

## Notes
The core schema migration (ForgeConfig struct, config.go functions, embed.go symbols) was already completed before this task. The remaining work was fixing internal/cmd/config_test.go which used old YAML keys (test-profiles:, capabilities:) and old CLI key names in test setup and assertions.
