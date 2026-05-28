---
type: validation.ux
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
Validate user experience for the {{.FeatureSlug}} feature.

## Validation Criteria
{{.AcceptanceCriteria}}

## Additional Checks
- Read the UI design spec (if exists) at docs/features/{{.FeatureSlug}}/design/ui-design.md
- Check docs/conventions/ for UX-related standards
- Verify accessibility, usability, and consistency against design specs
