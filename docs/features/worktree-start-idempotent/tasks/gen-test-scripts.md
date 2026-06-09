---
id: "T-test-gen-scripts"
title: "Generate CLI Functional Test Scripts"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["T-test-gen-contracts"]
type: "test.gen-scripts"
surface-key: "."
surface-type: "cli"
---

Generate executable test scripts for the worktree-start-idempotent feature.
Test type: cli.

## Feature Paths

Discover the feature's testing directory layout before starting:
```bash
ls docs/features/worktree-start-idempotent/testing/                                 # journeys
ls docs/features/worktree-start-idempotent/testing/<journey>/contracts/              # contracts
```

Read the approved test cases and generate scripts using the framework from the surface.

## Acceptance Criteria

- [ ] All acceptance criteria met

Type: **cli**
