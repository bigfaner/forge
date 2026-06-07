---
status: "completed"
started: "2026-06-07 21:48"
completed: "2026-06-07 21:56"
time_spent: "~8m"
---

# Task Record: 3 Fix CLI help text for task/feature/quality-gate commands

## Summary
Fixed CLI help text for 6 commands (C2-C6, C11): updated cobra Long/Short descriptions in cleanup, task claim, task validate, task add, feature, and quality-gate to accurately reflect actual code behavior

## Changes

### Files Created
无

### Files Modified
- forge-cli/internal/cmd/cleanup.go
- forge-cli/internal/cmd/task/claim.go
- forge-cli/internal/cmd/task/validate.go
- forge-cli/internal/cmd/task/add.go
- forge-cli/internal/cmd/feature/feature.go
- forge-cli/internal/cmd/qualitygate/quality_gate.go

### Key Decisions
- Listed all 15 validations in task validate Long (not just the original 5) to match actual validator.run() calls
- Included docs-only skip, retry-once policy, and fix task auto-creation in quality-gate Long
- Added fix-task usage pattern with --source-task-id/--block-source in task add Long

## Test Results
- **Tests Executed**: Yes
- **Passed**: 9
- **Failed**: 0
- **Coverage**: 74.1%

## Acceptance Criteria
- [x] forge cleanup --help Long includes blocked/suspended/rejected status cleanup
- [x] forge task claim --help Long includes auto-unblock behavior
- [x] forge task validate --help Long lists all 12+ validation steps
- [x] forge task add --help Long includes usage overview with fix-task vs regular task differences
- [x] forge feature --help Long includes set subcommand and behavior description
- [x] forge quality-gate --help Long includes fix task auto-creation, retry-once, docs-only skip

## Notes
All changes are string-only edits to cobra Command Long/Short fields. No logic changes. go build, go fmt, go vet, golangci-lint all pass. Targeted tests for all affected packages pass (9 suites).
