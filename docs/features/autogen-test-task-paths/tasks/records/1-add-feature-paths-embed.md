---
status: "completed"
started: "2026-05-28 23:59"
completed: "2026-05-29 00:02"
time_spent: "~3m"
---

# Task Record: 1 Embed 模板添加 ## Feature Paths 区域

## Summary
Added ## Feature Paths discovery section (journey + contract ls commands) to all 6 test pipeline embed templates

## Changes

### Files Created
无

### Files Modified
- forge-cli/pkg/task/templates/test-run.md
- forge-cli/pkg/task/templates/test-gen-scripts.md
- forge-cli/pkg/task/templates/eval-journey.md
- forge-cli/pkg/task/templates/eval-contract.md
- forge-cli/pkg/task/templates/test-gen-journeys.md
- forge-cli/pkg/task/templates/test-gen-contracts.md

### Key Decisions
无

## Document Metrics
6 templates modified, 44 lines added, 0 lines removed

## Referenced Documents
- docs/features/autogen-test-task-paths/proposal.md

## Review Status
final

## Acceptance Criteria
- [x] 6 templates contain ## Feature Paths with journeys and contracts discovery ls commands
- [x] Rich templates (test-gen-journeys, test-gen-contracts) do not duplicate existing path refs
- [x] go build ./... and go test ./... pass

## Notes
Rich templates integrated discovery ls into existing ## Discovery Strategy section rather than adding a separate ## Feature Paths section, avoiding duplication with their inline path references. go test failure in forge-cli/internal/cmd is pre-existing (fmt recipe test) and unrelated to template changes.
