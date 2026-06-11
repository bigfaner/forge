---
id: "T-test-4"
title: "Graduate Test Scripts"
priority: "P1"
estimated_time: "30min"
dependencies: ["T-test-3"]
status: pending
type: "test-pipeline.graduate"
---

# T-test-4: Graduate Test Scripts

## Description

Call `/graduate-tests` skill to migrate feature test scripts from `tests/e2e/features/typed-task-dispatch/` to the project-wide regression suite.

## Acceptance Criteria

- [ ] Test scripts migrated to `tests/e2e/` regression suite
- [ ] `tests/e2e/.graduated/typed-task-dispatch` marker created
- [ ] TypeScript compilation passes after migration

## Implementation Notes

1. Verify e2e tests passed (check `latest.md`)
2. Run `/graduate-tests` skill
