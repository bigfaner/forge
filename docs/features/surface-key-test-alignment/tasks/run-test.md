---
id: "T-test-run"
title: "Run CLI Functional Test"
priority: "P1"
estimated_time: "30min-1h"
dependencies: ["T-test-gen-journeys"]
type: "test.run"
surface-key: "."
surface-type: "cli"
---

Execute staged test scripts for the surface-key-test-alignment feature.

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/surface-key-test-alignment/testing/                                 # journeys
ls docs/features/surface-key-test-alignment/testing/<journey>/contracts/              # contracts
```

## Feature Context
- Scope: .

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.

Type: **cli**

## Acceptance Criteria
- [ ] All staged test scripts executed
- [ ] Test results recorded and failures diagnosed
