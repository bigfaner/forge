---
type: test.gen-scripts
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
Generate executable test scripts for the {{.FeatureSlug}} feature.{{if .SurfaceType}}
Test type: {{.SurfaceType}}.{{end}}

Read the approved test cases and generate scripts using the framework from the surface.
