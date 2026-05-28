---
type: eval.journey
category: eval
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
Evaluate Journey quality for the {{.FeatureSlug}} feature using the 6-dimension rubric (1000-point scale).

## Discovery Strategy
Scan `docs/features/{{.FeatureSlug}}/testing/journeys/` for all Journey files listed in `manifest.md`.

For each Journey:
1. Run `/eval --type journey` using the journey rubric (`eval/rubrics/journey.md`)
2. Scoring dimensions: Completeness, Semantic Purity, Precondition Exclusivity, Fact Alignment, Surface Fitness, Internal Consistency
3. Target score: 850/1000 with all dimensions above min thresholds

If any Journey fails evaluation after max iterations, report the failure and abort. Do not proceed to gen-contracts with low-quality Journeys.

## Acceptance Criteria
- [ ] All Journeys scored >= 850/1000
- [ ] All dimensions above min threshold per rubric
- [ ] Eval report written to `testing/journeys/.eval-report.md`
