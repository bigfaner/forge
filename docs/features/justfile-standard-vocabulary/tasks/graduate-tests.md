---
id: "T-test-4"
title: "Graduate Test Scripts"
priority: "P1"
estimated_time: "30min"
dependencies: ["T-test-3"]
status: pending
---

# Graduate Test Scripts

## Description

Call `/graduate-tests` skill to migrate feature test scripts from `tests/e2e/justfile-standard-vocabulary/` to the project-wide regression suite at `tests/e2e/<target>/`.

This task is a gate: it only proceeds if e2e tests are passing.

## Reference Files

- `testing/results/latest.md` — Must show status = PASS before graduating
- `tests/e2e/justfile-standard-vocabulary/` — Source scripts to migrate
- `tests/e2e/` — Destination regression suite

## Acceptance Criteria

- [ ] `testing/results/latest.md` shows status = PASS
- [ ] `tests/e2e/.graduated/justfile-standard-vocabulary` marker exists
- [ ] Spec files present in `tests/e2e/<category>/`

## User Stories

No direct user story mapping. This is a standard test graduation task.

## Implementation Notes

**Step 1: Verify e2e passed**

Read `testing/results/latest.md`. Check status field.

- Status = PASS → proceed to Step 2
- Status = FAIL → mark task `blocked` and stop:
  ```
  e2e tests are still failing (see testing/results/latest.md).
  Wait for fix tasks to complete, then unblock:
    task status T-test-4 pending
  ```

**Step 2: Graduate**

Run `/graduate-tests` skill. The skill will:
- Read each spec file and understand its content
- Decide classification (split / merge / keep as-is)
- Migrate to `tests/e2e/<category>/`
- Rewrite import paths
- Create graduation marker `tests/e2e/.graduated/justfile-standard-vocabulary`

**Step 3: Record**

Mark task completed. The `all-completed` hook will run regression (`just test-e2e` without `--feature`) on the next session end to verify the graduated scripts integrate cleanly with the existing suite.
