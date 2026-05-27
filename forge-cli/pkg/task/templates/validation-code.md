---
type: validation.code
category: validation
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
Validate code quality for the {{.FeatureSlug}} feature.

## Validation Criteria
{{.AcceptanceCriteria}}

## Additional Checks
- Check docs/conventions/ for project-specific quality standards (read each file's `domains` frontmatter to determine relevance)
- Run the quality gate: just compile → just fmt → just lint → just test
