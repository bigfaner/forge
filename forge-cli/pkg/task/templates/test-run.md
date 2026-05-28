---
type: test.run
category: test
identity:
  - TaskID
  - TaskType
  - FeatureSlug
context:
  - Mode
  - SurfaceKey
  - SurfaceType
  - SurfaceTypes
  - AcceptanceCriteria
---
Execute staged test scripts for the {{.FeatureSlug}} feature.

## Feature Context
{{if .SurfaceKey}}- Scope: {{.SurfaceKey}}{{end}}

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.
