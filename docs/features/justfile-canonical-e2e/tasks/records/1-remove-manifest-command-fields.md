---
status: "completed"
started: "2026-05-15 00:31"
completed: "2026-05-15 00:38"
time_spent: "~7m"
---

# Task Record: 1 Remove run.* and graduate.* fields from all 6 manifest.yaml files

## Summary
Removed run.* and graduate.* command fields from all 6 profile manifest.yaml files. These fields were never parsed by any Go code and duplicated commands already canonically defined in justfile-recipes, establishing justfile-recipes as the single source of truth.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/profile/profiles/go-test/manifest.yaml
- forge-cli/pkg/profile/profiles/java-junit/manifest.yaml
- forge-cli/pkg/profile/profiles/maestro/manifest.yaml
- forge-cli/pkg/profile/profiles/pytest/manifest.yaml
- forge-cli/pkg/profile/profiles/rust-test/manifest.yaml
- forge-cli/pkg/profile/profiles/web-playwright/manifest.yaml

### Key Decisions
- Verified via grep that no Go code reads run.* or graduate.* fields from manifest YAML before removal
- The profileManifest Go struct only has Capabilities field, confirming run/graduate fields were dead data

## Test Results
- **Tests Executed**: Yes
- **Passed**: 14
- **Failed**: 0
- **Coverage**: 91.2%

## Acceptance Criteria
- [x] All 6 manifest.yaml files contain zero run.* fields
- [x] All 6 manifest.yaml files contain zero graduate.* fields
- [x] Remaining manifest fields (name, display, language, file-extension, test-directory, capabilities, templates) are untouched
- [x] go test ./... passes after removal
- [x] forge profile get <profile> --manifest still works for all 6 profiles

## Notes
Pre-existing test failure in internal/cmd (TestSaveIndexAndSignalCompletion_SaveIndexError) is unrelated to this change. Profile package tests all pass with 91.2% coverage.
