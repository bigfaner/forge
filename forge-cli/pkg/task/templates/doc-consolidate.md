---
type: doc.consolidate
category: doc
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
Extract and consolidate business rules and tech specs from the {{.FeatureSlug}} feature.

## Feature Context
{{if .SurfaceKey}}- Scope: {{.SurfaceKey}}{{end}}

## Discovery Strategy
1. Scan docs/features/{{.FeatureSlug}}/ for all feature documents (PRD, design, task records)
2. Scan docs/proposals/{{.FeatureSlug}}/ for proposal
3. Extract rules and specs from discovered documents
4. Compare against existing specs in docs/business-rules/ and docs/conventions/

Run in non-interactive mode: auto-integrate all CROSS items. Commit with [auto-specs] tag.

## Acceptance Criteria

{{.AcceptanceCriteria}}
