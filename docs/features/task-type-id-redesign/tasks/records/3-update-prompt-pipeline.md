---
status: "completed"
started: "2026-05-19 14:43"
completed: "2026-05-19 14:48"
time_spent: "~5m"
---

# Task Record: 3 Update prompt pipeline

## Summary
Updated prompt pipeline: renamed 17 template files to match new prefix-based type names, updated typeToTemplate map with all new type-to-file mappings, created validation-code.md and validation-ux.md prompt templates for the two new validation types, and added them to the typeToTemplate map. genScriptBases already uses readable IDs from task 2.

## Changes

### Files Created
- forge-cli/pkg/prompt/data/validation-code.md
- forge-cli/pkg/prompt/data/validation-ux.md

### Files Modified
- forge-cli/pkg/prompt/prompt.go

### Key Decisions
- Template filename follows hard rule: type name with '.' replaced by '-' (e.g., coding.feature -> coding-feature.md)
- validation-code.md follows gate.md structure with code quality checks and quality gate
- validation-ux.md focuses on UX-specific validation without compile/test gate since it validates user-facing behavior
- genScriptBases unchanged since task 2 already migrated to readable IDs (T-test-gen-scripts, T-quick-gen-and-run)

## Test Results
- **Tests Executed**: Yes
- **Passed**: 23
- **Failed**: 0
- **Coverage**: 90.6%

## Acceptance Criteria
- [x] prompt.Synthesize() resolves new type names to correct template files
- [x] All 17 template files renamed per mapping
- [x] validation-code.md and validation-ux.md prompt templates created
- [x] genScriptBases updated to new ID prefixes

## Notes
无
