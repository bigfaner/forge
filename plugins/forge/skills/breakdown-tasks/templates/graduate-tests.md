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

Call `/graduate-tests` skill to migrate feature test scripts from `tests/e2e/features/<slug>/` to the
project-wide regression suite at `tests/e2e/<target>/`.

This task is a gate: it only proceeds if e2e tests are passing.

## Reference Files

- `tests/e2e/features/<slug>/results/latest.md` — Must show status = PASS before graduating
- `tests/e2e/features/<slug>/` — Source scripts to migrate
- `tests/e2e/` — Destination regression suite

## Acceptance Criteria

- [ ] `tests/e2e/features/<slug>/results/latest.md` shows status = PASS
- [ ] `tests/e2e/.graduated/<slug>` marker exists
- [ ] Spec files present in `tests/e2e/<module>/`

## User Stories

No direct user story mapping. This is a standard test graduation task.

## Implementation Notes

**Step 1: Verify e2e passed**

Read `tests/e2e/features/<slug>/results/latest.md`. Check status field.

- Status = PASS → proceed to Step 2
- Status = FAIL → mark task `blocked` and stop:
  ```
  e2e tests are still failing (see tests/e2e/features/<slug>/results/latest.md).
  Wait for fix tasks to complete, then unblock:
    task status T-test-4 pending
  ```

**Step 2: Graduate**

Run `/graduate-tests` skill. The skill will:
- Read each spec file and understand its content
- Decide classification by functional module (split / merge / keep as-is)
- Migrate to `tests/e2e/<module>/` (functional domain, not feature slug)
- Validate TypeScript compilation post-migration
- Rewrite import paths
- Create graduation marker `tests/e2e/.graduated/<slug>` with module metadata

**Step 3: Record**

Mark task completed. T-test-4.5 will run full regression to verify the
graduated scripts integrate cleanly with the existing suite.
