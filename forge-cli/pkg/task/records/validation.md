---
type: record
category: record
variables:
  - Status
  - Started
  - Completed
  - TimeSpent
  - TaskID
  - TaskTitle
  - Summary
  - TypeReclassification
  - FilesCreatedFormatted
  - FilesModifiedFormatted
  - KeyDecisionsFormatted
  - ValidationPassedFormatted
  - IssuesFoundFormatted
  - AcceptanceCriteriaFormatted
  - Notes
---
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

## Pass/Fail Verdict
- **Status**: {{.ValidationPassedFormatted}}

## Issues Found
{{.IssuesFoundFormatted}}

## Acceptance Criteria
{{.AcceptanceCriteriaFormatted}}

## Notes
{{.Notes}}
