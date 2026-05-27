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
  - CasesGeneratedFormatted
  - CasesEvaluatedFormatted
  - ScriptsCreatedFormatted
  - TestResultsFormatted
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

## Cases Generated
{{.CasesGeneratedFormatted}}

## Cases Evaluated
{{.CasesEvaluatedFormatted}}

## Scripts Created
{{.ScriptsCreatedFormatted}}

## Test Results
{{.TestResultsFormatted}}

## Acceptance Criteria
{{.AcceptanceCriteriaFormatted}}

## Notes
{{.Notes}}
