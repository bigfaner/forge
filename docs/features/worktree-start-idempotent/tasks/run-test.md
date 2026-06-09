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

Execute staged test scripts for the worktree-start-idempotent feature.

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/worktree-start-idempotent/testing/                                 # journeys
ls docs/features/worktree-start-idempotent/testing/<journey>/contracts/              # contracts
```

## Feature Context
- Scope: .

Run all staged test scripts. If tests fail, identify root cause, apply minimal fix, and re-run.

## Acceptance Criteria

- [ ] All acceptance criteria met

### Hard Acceptance Criteria (non-negotiable)

- [ ] All test cases MUST pass — no skipped tests, no expected failures, no TODO placeholders
- [ ] Tests MUST verify actual functional behavior — no placeholder tests, no always-pass mocks, no stub assertions that validate nothing

Type: **cli**
