---
type: test.gen-scripts
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
Generate executable test scripts for the {{.FeatureSlug}} feature.{{if .SurfaceType}}
Test type: {{.SurfaceType}}.{{end}}

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/{{.FeatureSlug}}/testing/                                 # journeys
ls docs/features/{{.FeatureSlug}}/testing/<journey>/contracts/              # contracts
```

Read the approved test cases and generate scripts using the framework from the surface.
