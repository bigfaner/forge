---
type: doc.review
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
Review documentation quality for the {{.FeatureSlug}} feature ({{.Mode}} mode).

## Acceptance Criteria Summary

The following acceptance criteria are pre-extracted from doc tasks. Use these as the review baseline.

{{.DocTaskCriteria}}

## Discovery Strategy

Scan ONLY the following allowlist of directories for target documents:
- docs/features/{{.FeatureSlug}}/ (prd/, design/, testing/, and any subdirectories)
- docs/proposals/{{.FeatureSlug}}/

EXCLUDE the following from scanning — do NOT read or process these:
- tasks/ directory (task definitions are not deliverables)
- tasks/records/ directory (execution records are not deliverables)
- manifest.md (build artifact)
- index.json (build artifact)

Only .md files under the allowlist directories are target deliverables.

## Acceptance Criteria

{{.AcceptanceCriteria}}
