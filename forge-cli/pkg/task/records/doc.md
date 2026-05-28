---
type: record
category: record
identity:
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
  - DocMetricsFormatted
  - ReferencedDocsFormatted
  - ReviewStatusFormatted
  - AcceptanceCriteriaFormatted
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

## Document Metrics
{{.DocMetricsFormatted}}

## Referenced Documents
{{.ReferencedDocsFormatted}}

## Review Status
{{.ReviewStatusFormatted}}

## Acceptance Criteria
{{.AcceptanceCriteriaFormatted}}

## Notes
{{.Notes}}
