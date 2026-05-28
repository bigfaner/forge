---
type: code-quality.simplify
category: coding
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
Simplify and clean up code for the {{.FeatureSlug}} feature.

## Discovery Strategy
1. Run `git diff --name-only main...HEAD` to identify files changed by this feature
2. Focus cleanup on changed files only
3. The skill resolves scope: git diff > feature context > user-specified paths

Do NOT clean up files outside this feature's scope.
