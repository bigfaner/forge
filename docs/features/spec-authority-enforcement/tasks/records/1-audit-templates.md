---
status: "completed"
started: "2026-05-24 09:42"
completed: "2026-05-24 09:44"
time_spent: "~2m"
---

# Task Record: 1 Audit all 19 prompt templates for Reference Files strengthening

## Summary
Audited all 19 prompt templates under forge-cli/pkg/prompt/data/*.md for Reference Files strengthening needs. 9 templates need strengthening (5 coding.*, 2 doc.*, 3 validation/gate), 10 can be skipped (delegate-only, informational, or recovery templates).

## Changes

### Files Created
无

### Files Modified
无

### Key Decisions
无

## Document Metrics
19 templates audited, 9 need strengthening, 10 skip

## Referenced Documents
- docs/proposals/spec-authority-enforcement/proposal.md
- docs/features/spec-authority-enforcement/tasks/1-audit-templates.md

## Review Status
final

## Acceptance Criteria
- [x] All 19 template files in forge-cli/pkg/prompt/data/ have been read and analyzed
- [x] Each template classified as needs-strengthening or skip with explicit reasoning
- [x] Audit criteria applied: (a) coding/doc task type, (b) Step 1 contains read-task-file, (c) spec-driven implementation tasks
- [x] Output structured audit table: Template Name | Needs Strengthening | Reason | Current Step 1 Location | Current Verify Step Location

## Notes
Needs-strengthening (9): coding-feature, coding-enhancement, coding-refactor, coding-fix, coding-cleanup, doc, doc-review, gate, validation-code, validation-ux. Skip (10): clean-code, doc-consolidate, doc-drift, doc-summary, test-gen-and-run, test-gen-scripts, test-run, test-verify-regression, fix-record-missed. Task 2 should target the 9 templates that need strengthening.
