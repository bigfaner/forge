---
id: "T-test-4"
title: "Graduate Test Scripts"
priority: "P1"
estimated_time: "30min"
dependencies: ["T-test-3"]
status: pending
noTest: false
mainSession: false
---

# Graduate Test Scripts

## Description

Call `/graduate-tests` skill to migrate feature test scripts to the project-wide regression suite.

## Reference Files

- `tests/e2e/features/task-executor-skeleton/results/latest.md` — Must show PASS
- `tests/e2e/features/task-executor-skeleton/` — Source scripts

## Acceptance Criteria

- [ ] `tests/e2e/features/task-executor-skeleton/results/latest.md` shows status = PASS
- [ ] `tests/e2e/.graduated/task-executor-skeleton` marker exists
- [ ] Spec files present in `tests/e2e/<module>/`

## User Stories

No direct user story mapping. This is a standard test graduation task.

## Implementation Notes

1. Verify e2e passed (read latest.md)
2. Run `/graduate-tests` skill
3. Record task completed
