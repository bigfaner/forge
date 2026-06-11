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

## Acceptance Criteria

{{.AcceptanceCriteria}}

### Hard Acceptance Criteria (non-negotiable)

- [ ] All test cases MUST pass — no skipped tests, no expected failures, no TODO placeholders
- [ ] Tests MUST verify actual functional behavior — no placeholder tests, no always-pass mocks, no stub assertions that validate nothing
