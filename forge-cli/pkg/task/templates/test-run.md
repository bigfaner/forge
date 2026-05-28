---
type: test.run
category: test
variables:
  - TaskID
  - TaskType
  - FeatureSlug
  - Mode
  - SurfaceKey
  - SurfaceType
  - SurfaceTypes
  - AcceptanceCriteria
  - DocTaskCriteria
---
Execute staged test scripts for the {{.FeatureSlug}} feature.

## Feature Context
{{if .SurfaceKey}}- Scope: {{.SurfaceKey}}{{end}}

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.
