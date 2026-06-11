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
  - GateChecksFormatted
  - GatePassedFormatted
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

{{end}}## Gate Checks
{{.GateChecksFormatted}}

## Gate Status
- **Passed**: {{.GatePassedFormatted}}

## Notes
{{.Notes}}
