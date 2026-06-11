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
  - ScoreFormatted
  - FindingsFormatted
  - SeverityFormatted
  - PassedFormatted
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

{{end}}## Eval Score
- **Score**: {{.ScoreFormatted}}

## Findings
{{.FindingsFormatted}}

## Severity
- **Severity**: {{.SeverityFormatted}}

## Passed
- **Passed**: {{.PassedFormatted}}

## Acceptance Criteria
{{.AcceptanceCriteriaFormatted}}

## Notes
{{.Notes}}
