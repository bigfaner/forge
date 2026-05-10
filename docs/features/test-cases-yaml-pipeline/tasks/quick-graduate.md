---
id: "T-quick-4"
title: "Graduate Quick Test Scripts"
priority: "P1"
estimated_time: "15min"
dependencies: ["T-quick-3"]
status: pending
noTest: false
mainSession: false
---

# Graduate Quick Test Scripts

## Description

Call `/graduate-tests` skill to migrate feature test scripts from `tests/e2e/features/test-cases-yaml-pipeline/` to the
project-wide regression suite at `tests/e2e/<target>/`.

This task is a gate: it only proceeds if e2e tests are passing.

## Reference Files

- `tests/e2e/features/test-cases-yaml-pipeline/results/latest.md` — Must show status = PASS before graduating
- `tests/e2e/features/test-cases-yaml-pipeline/` — Source scripts to migrate
- `tests/e2e/` — Destination regression suite

## Acceptance Criteria

- [ ] `tests/e2e/features/test-cases-yaml-pipeline/results/latest.md` shows status = PASS
- [ ] `tests/e2e/.graduated/test-cases-yaml-pipeline` marker exists
- [ ] Spec files present in `tests/e2e/<module>/`

## Implementation Notes

**Step 1: Verify e2e passed**

Read `tests/e2e/features/test-cases-yaml-pipeline/results/latest.md`. Check status field.

- Status = PASS → proceed to Step 2
- Status = FAIL → mark task `blocked` and stop

**Step 2: Graduate**

Run `/graduate-tests` skill. The skill will:
- Read each spec file and understand its content
- Decide classification by functional module (split / merge / keep as-is)
- Migrate to `tests/e2e/<module>/`
- Validate TypeScript compilation post-migration
- Rewrite import paths
- Create graduation marker `tests/e2e/.graduated/test-cases-yaml-pipeline`

**Step 3: Record**

Mark task completed. T-quick-5 will run full regression to verify the
graduated scripts integrate cleanly with the existing suite.
