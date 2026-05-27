---
status: "completed"
started: "2026-05-28 01:03"
completed: "2026-05-28 01:10"
time_spent: "~7m"
---

# Task Record: 3 简化 cleanTemplateOutput 并添加模板校验

## Summary
Simplified cleanTemplateOutput to collapseBlankLines (blank-line collapsing only), removed 4 conditional deletion modes and isLabelWithEmptyValue helper. Enhanced ValidatePromptTemplates with missingkey=error zero-value struct Execute to catch field misspellings at startup.

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go

### Key Decisions
- Replaced cleanTemplateOutput (4 conditional modes + blank-line collapse) with collapseBlankLines (blank-line collapse only) since templates now handle all conditional logic via {{if}} blocks
- Added io.Discard Execute validation in ValidatePromptTemplates to catch field misspellings using missingkey=error with zero-value promptTemplateData
- Evolved 5 direct cleanTemplateOutput tests to test collapseBlankLines; preserved all Synthesize integration tests

## Test Results
- **Tests Executed**: Yes
- **Passed**: 58
- **Failed**: 0
- **Coverage**: 73.7%

## Acceptance Criteria
- [x] cleanTemplateOutput() only keeps blank-line collapse logic, 4 conditional deletion modes removed
- [x] ValidatePromptTemplates() uses missingkey=error + zero-value promptTemplateData Execute to io.Discard
- [x] go build ./... passes, forge prompt get-by-task-id output functionally equivalent (blank-line diffs allowed)

## Notes
Coverage 73.7% is unchanged from pre-refactor baseline. The removed dead code (cleanTemplateOutput conditional logic, isLabelWithEmptyValue) is no longer needed since Task 2 migrated all templates to use {{if}} blocks for conditional content.
