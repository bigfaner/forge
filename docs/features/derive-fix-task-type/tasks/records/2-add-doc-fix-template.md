---
status: "completed"
started: "2026-05-29 11:11"
completed: "2026-05-29 11:14"
time_spent: "~3m"
---

# Task Record: 2 Add doc-fix task template

## Summary
Created doc-fix.md task template for doc-category fix tasks. Template adapts coding.fix.md structure but removes Surface Inference, Fix Boundaries (dev server/npm/test), and Verification sections. Retains dual-frontmatter pattern, Go template placeholders, Root Cause, Reference Files, and auto-restore note.

## Changes

### Files Created
- forge-cli/pkg/task/templates/doc-fix.md

### Files Modified
无

### Key Decisions
无

## Document Metrics
52 lines, 4 sections (Root Cause, Reference Files, Content Fix Guidance, auto-restore note)

## Referenced Documents
- forge-cli/pkg/task/templates/coding.fix.md
- forge-cli/pkg/task/tasktemplate.go

## Review Status
final

## Acceptance Criteria
- [x] doc-fix.md template exists at forge-cli/pkg/task/templates/doc-fix.md
- [x] Template contains fix instructions scoped to doc-type failures: no code quality gates, no test execution, only markdown/content fixes
- [x] GetTaskTemplate("doc.fix") returns the template content without error

## Notes
tasktemplate.go already had doc.fix defaults registered (Priority: P0, Breaking: false, EstimatedTime: 30min, Type: doc.fix, IDPrefix: doc-fix). Go embed picks up the new file automatically.
