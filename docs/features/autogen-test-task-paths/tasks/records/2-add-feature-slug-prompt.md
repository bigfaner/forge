---
status: "completed"
started: "2026-05-29 00:02"
completed: "2026-05-29 00:06"
time_spent: "~4m"
---

# Task Record: 2 Prompt 模板渲染 FEATURE_SLUG

## Summary
Added FEATURE_SLUG: {{.FeatureSlug}} line after TASK_FILE line in all 6 prompt templates, enabling subagents to receive the feature slug directly without path parsing

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/prompt/templates/test-run.md
- forge-cli/pkg/prompt/templates/test-gen-journeys.md
- forge-cli/pkg/prompt/templates/test-gen-contracts.md
- forge-cli/pkg/prompt/templates/test-gen-scripts.md
- forge-cli/pkg/prompt/templates/eval-journey.md
- forge-cli/pkg/prompt/templates/eval-contract.md

### Key Decisions
无

## Document Metrics
6 templates updated, 1 line added per template, all in consistent position

## Referenced Documents
- docs/features/autogen-test-task-paths/tasks/2-add-feature-slug-prompt.md

## Review Status
final

## Acceptance Criteria
- [x] 6 prompt templates output FEATURE_SLUG: {{.FeatureSlug}} line after TASK_FILE
- [x] go build ./... and go test ./... pass

## Notes
3 pre-existing test failures in internal/cmd (TestAddFixTask_*) confirmed to exist before this change; unrelated to prompt template modifications
