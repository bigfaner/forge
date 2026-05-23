---
status: "{{.Status}}"
started: "{{.Started}}"
completed: "{{.Completed}}"
time_spent: "{{.TimeSpent}}"
---

# Task Record: {{.TaskID}} {{.TaskTitle}}

## Summary
{{.Summary}}

{{if .TypeReclassification}}## Type Reclassification
- Original: {{.TypeReclassification.OriginalType}}
- Actual: {{.TypeReclassification.ActualType}}
- Reason: {{.TypeReclassification.Reason}}

{{end}}## Changes

### Files Created
{{.FilesCreatedFormatted}}

### Files Modified
{{.FilesModifiedFormatted}}

### Key Decisions
{{.KeyDecisionsFormatted}}

## Test Results
- **Tests Executed**: {{.TestsExecuted}}
- **Passed**: {{.TestsPassed}}
- **Failed**: {{.TestsFailed}}
- **Coverage**: {{.CoverageFormatted}}

## Acceptance Criteria
{{.AcceptanceCriteriaFormatted}}

## Notes
{{.Notes}}
