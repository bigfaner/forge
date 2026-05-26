---
status: "completed"
started: "2026-05-26 23:01"
completed: "2026-05-26 23:04"
time_spent: "~3m"
---

# Task Record: 3 Update CLI documentation (WORKFLOW.md, OVERVIEW.md, README.md)

## Summary
Updated WORKFLOW.md, OVERVIEW.md, README.md to replace --template with --type syntax for forge task add

## Changes

### Files Created
无

### Files Modified
- forge-cli/docs/WORKFLOW.md
- forge-cli/docs/OVERVIEW.md
- README.md

### Key Decisions
无

## Document Metrics
3 files updated, 6 distinct changes, 0 --template references remaining

## Referenced Documents
- docs/proposals/unify-task-add-type-param/proposal.md
- forge-cli/CLAUDE.md
- forge-cli/internal/cmd/task/add.go
- forge-cli/pkg/template/template.go

## Review Status
final

## Acceptance Criteria
- [x] No documentation file contains --template in the context of forge task add
- [x] WORKFLOW.md flag table lists --type with description 'Task type, auto-discovers matching template'
- [x] All example commands use --type coding.fix instead of --template fix-task
- [x] README.md CLI parameter table reflects new --type flag

## Notes
make check-docs passed all 32 docsync tests. Did not touch task profile get --template references (different feature). WORKFLOW.md also updated: addFixTask() description, IDPrefix example, template defaults section expanded to include coding.cleanup.
