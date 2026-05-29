---
status: "blocked"
started: "2026-05-29 11:21"
completed: "N/A"
time_spent: ""
---

# Task Record: T-review-doc Review Documentation Quality

## Summary
Reviewed documentation quality for derive-fix-task-type feature. 6 of 7 AC items pass. One AC item (TYPE listed as extractable field) fails: forge task claim outputs TYPE but execute-task.md and run-tasks.md do not list TYPE in their extract fields. This fix requires modifying plugins/ files, which is out of scope for doc.review task (docs/ only).

## Changes

### Files Created
无

### Files Modified
无

### Key Decisions
无

## Document Metrics
AC pass rate: 6/7 (85%), files scanned: 4 skill files + 1 template + 1 proposal

## Referenced Documents
- docs/proposals/derive-fix-task-type/proposal.md

## Review Status
reviewed

## Acceptance Criteria
- [x] doc-fix.md template exists at forge-cli/pkg/task/templates/doc-fix.md
- [x] Template contains fix instructions scoped to doc-type failures: no code quality gates, no test execution, only markdown/content fixes
- [x] GetTaskTemplate("doc.fix") returns the template content without error
- [x] Error-handling instructions in task-executor.md, execute-task.md, run-tasks.md, submit-task/SKILL.md use derivation rule
- [x] Derivation rule table documented in at least one canonical location
- [ ] TYPE and TASK_CATEGORY documented as extractable fields from forge task claim output in skill files
- [x] grep -rn 'type coding\.fix' plugins/forge/ --include='*.md' returns zero matches in error-handling contexts

## Notes
AC-3 (TYPE as extractable field) fails: TASK_CATEGORY is listed in execute-task.md and run-tasks.md extract fields, but TYPE (output by claim.go line 288) is not. Fix requires adding TYPE to the extract field lists in execute-task.md line 24-31 and run-tasks.md line 55. This is a plugins/ change, outside doc.review scope.
