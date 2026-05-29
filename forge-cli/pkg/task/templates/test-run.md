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

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/{{.FeatureSlug}}/testing/                                 # journeys
ls docs/features/{{.FeatureSlug}}/testing/<journey>/contracts/              # contracts
```

## Feature Context
{{if .SurfaceKey}}- Scope: {{.SurfaceKey}}{{end}}

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.
