---
type: validation.code
category: validation
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
Validate code quality for the {{.FeatureSlug}} feature.

## Acceptance Criteria

{{.AcceptanceCriteria}}

## Validation Criteria

## Additional Checks
- Check docs/conventions/ for project-specific quality standards (read each file's `domains` frontmatter to determine relevance)
- Run the quality gate: just compile → just fmt → just lint → just test
