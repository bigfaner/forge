---
id: "fix-e2e-{{round}}-{{index}}"
title: "Fix e2e Test Failure: {{test_name}}"
priority: "P0"
estimated_time: "30min-2h"
dependencies: []
status: pending
---

# fix-e2e-{{round}}-{{index}}: Fix e2e Test Failure

## Description

This is fix attempt round {{round}}. Steps:

1. Read `testing/results/latest.md` for failure overview
2. Read `testing/results/failures/failure-{{test_case_id}}.md` for specific failure details
3. Locate root cause (code logic / test script / environment config)
4. Fix and verify

## Reference Files

- `testing/results/latest.md` — Test results overview
- `testing/results/failures/failure-{{test_case_id}}.md` — Failure details
- `testing/test-cases.md` — Test case document
- `testing/scripts/` — Test scripts directory

## Acceptance Criteria

- [ ] Root cause of failure identified
- [ ] Code or test script fixed
- [ ] All unit tests pass
