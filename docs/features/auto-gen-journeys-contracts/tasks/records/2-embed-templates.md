---
status: "completed"
started: "2026-05-24 00:01"
completed: "2026-05-24 00:05"
time_spent: "~4m"
---

# Task Record: 2 创建 embed 模板 test-gen-journeys.md 和 test-gen-contracts.md

## Summary
Created embed templates test-gen-journeys.md and test-gen-contracts.md with full content including AUTO_COMMIT and SKIP_EVAL_GATE conditional directives, mode-aware input source guidance, and clear process flow instructions for AI executors

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/data/test-gen-journeys.md
- forge-cli/pkg/task/data/test-gen-contracts.md
- forge-cli/pkg/task/autogen_test.go

### Key Decisions
- Templates use only renderBody()-supported placeholders ({{FEATURE_SLUG}}, {{MODE}}, {{SCOPE}})
- AUTO_COMMIT and SKIP_EVAL_GATE are natural-language conditional directives embedded in template text, not code variables
- Each template includes mode-specific input source guidance (Breakdown vs Quick) per proposal.md requirements

## Test Results
- **Tests Executed**: Yes
- **Passed**: 4
- **Failed**: 0
- **Coverage**: 85.0%

## Acceptance Criteria
- [x] data/test-gen-journeys.md uses {{FEATURE_SLUG}}, {{MODE}} placeholders
- [x] data/test-gen-contracts.md uses {{FEATURE_SLUG}}, {{MODE}} placeholders
- [x] test-gen-journeys.md contains AUTO_COMMIT=true conditional instruction for skipping user approval
- [x] test-gen-contracts.md contains SKIP_EVAL_GATE=true conditional instruction for skipping eval-journey check
- [x] Templates auto-included via go:embed data/*.md (no embed declaration change needed)
- [x] Template content clearly guides AI executor through gen-journeys/gen-contracts skill flow

## Notes
Existing TestAutogenTypeToFileMapping already validates embed FS inclusion. 4 new tests added: template content checks (AUTO_COMMIT, SKIP_EVAL_GATE, placeholders) and render-through GenerateTestTaskMD verification.
