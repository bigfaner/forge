---
status: "completed"
started: "2026-05-19 01:27"
completed: "2026-05-19 01:31"
time_spent: "~4m"
---

# Task Record: 5 Move forge-config schema and example YAML to CLI, update test paths

## Summary
Moved forge-config.schema.json and forge-config.example.yaml from plugins/forge/references/shared/ to forge-cli/internal/cmd/testdata/ (conventional Go test data location). Updated filepath.Join paths in config_schema_test.go to point to testdata/ instead of the old cross-component path. Verified all 5 config schema tests pass. Removed old files. Bumped CLI version from 4.1.0 to 4.1.1 (patch for dead code path cleanup).

## Changes

### Files Created
- forge-cli/internal/cmd/testdata/forge-config.schema.json
- forge-cli/internal/cmd/testdata/forge-config.example.yaml

### Files Modified
- forge-cli/internal/cmd/config_schema_test.go
- forge-cli/scripts/version.txt

### Key Decisions
- Used forge-cli/internal/cmd/testdata/ as the new location since it already exists and is the conventional Go test data directory
- Bumped patch version (4.1.0 -> 4.1.1) per CLI CLAUDE.md rules for dead code path cleanup

## Test Results
- **Tests Executed**: Yes
- **Passed**: 5
- **Failed**: 0
- **Coverage**: 3.3%

## Acceptance Criteria
- [x] forge-config.schema.json and forge-config.example.yaml exist in their new location under forge-cli/
- [x] config_schema_test.go paths updated to read from new location
- [x] go test ./forge-cli/internal/cmd/ -run TestConfigSchema -v passes
- [x] Old files removed from plugins/forge/references/shared/

## Notes
All 5 config schema tests (TestConfigSchemaAutoBlock, TestConfigSchemaAutoDefaults, TestConfigSchemaBackwardCompatible, TestConfigExampleDocumentsAllAutoFields, TestConfigSchemaTestFrameworkFields) pass with the new testdata/ paths.
