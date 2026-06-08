---
id: "T-test-run"
title: "Run CLI Functional Test"
priority: "P1"
estimated_time: "30min-1h"
dependencies: ["T-test-gen-scripts"]
type: "test.run"
surface-key: "."
surface-type: "cli"
---

Execute staged test scripts for the test-pipeline-interleaved feature.

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/test-pipeline-interleaved/testing/                                 # journeys
ls docs/features/test-pipeline-interleaved/testing/<journey>/contracts/              # contracts
```

## Feature Context
- Scope: .

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.

Type: **cli**

## Acceptance Criteria

- [ ] All acceptance criteria met
