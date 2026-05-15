---
status: "completed"
started: "2026-05-15 21:54"
completed: "2026-05-15 22:03"
time_spent: "~9m"
---

# Task Record: 2 Update task ID inference and prompt template for per-type tasks

## Summary
Update InferType() and prompt template to recognize per-type task IDs (e.g. T-test-2-api, T-test-2a-tui). Added typeSuffixedID() and ExtractTypeSuffix() helpers to infer.go. Updated InferType() switch cases for T-test-2/3/4 and T-quick-2/3/4 to match type-suffixed variants. Updated renderTemplate() in prompt.go to extract type suffix and inject --type <capability> into the gen-scripts template. Template omits --type when no type suffix present (backward compatible).

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/infer.go
- forge-cli/pkg/task/infer_test.go
- forge-cli/pkg/prompt/prompt.go
- forge-cli/pkg/prompt/prompt_test.go
- forge-cli/pkg/prompt/data/test-pipeline-gen-scripts.md
- forge-cli/scripts/version.txt

### Key Decisions
- typeSuffixedID() is a pure pattern matcher -- InferType() controls which bases accept type suffixes (T-test-2/3/4 and T-quick-2/3/4 only, not T-test-1 or T-test-4.5)
- ExtractTypeSuffix() is exported for use by prompt package to render --type arg
- Template uses {{TEST_TYPE_ARG}} placeholder that expands to ' --type <cap>' or empty string, avoiding conditional template logic

## Test Results
- **Tests Executed**: Yes
- **Passed**: 62
- **Failed**: 0
- **Coverage**: 90.2%

## Acceptance Criteria
- [x] InferType("T-test-2-api") returns TypeTestPipelineGenScripts
- [x] InferType("T-test-2-tui") returns TypeTestPipelineGenScripts
- [x] InferType("T-test-2-cli") returns TypeTestPipelineGenScripts
- [x] InferType("T-quick-2-api") returns TypeTestPipelineGenScripts
- [x] Profile-suffixed + type-suffixed: T-test-2a-api, T-quick-2b-tui
- [x] Existing patterns still work: T-test-2, T-test-2a, T-quick-2, T-quick-2a
- [x] Prompt template passes --type <capability> when task ID contains type suffix
- [x] Prompt template omits --type when task ID has no type suffix (backward compatible)

## Notes
无
