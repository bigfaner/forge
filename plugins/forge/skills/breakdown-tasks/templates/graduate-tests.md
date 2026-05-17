---
id: "T-test-4"
title: "Promote Test Scripts"
priority: "P1"
estimated_time: "15min"
dependencies: ["T-test-3"]
status: pending
noTest: false
mainSession: false
---

# Promote Test Scripts

## Description

Run `forge test promote <journey>` to promote the journey's test scripts from @feature tags to @regression tags.

This task is a gate: it only proceeds if e2e tests are passing.

## Reference Files

- `tests/<journey>/results/latest.md` — Must show status = PASS before promoting
- `tests/<journey>/` — Journey test scripts with @feature tags

## Acceptance Criteria

- [ ] Journey tests pass (verified by promote command)
- [ ] All @feature tags replaced with @regression tags
- [ ] No code changes other than tag replacements (verified via `git diff`)

## User Stories

No direct user story mapping. This is a standard test promotion task.

## Implementation Notes

**Step 1: Verify e2e passed**

Read `tests/<journey>/results/latest.md`. Check status field.

- Status = PASS -> proceed to Step 2
- Status = FAIL -> mark task `blocked` and stop:
  ```
  e2e tests are still failing (see tests/<journey>/results/latest.md).
  Wait for fix tasks to complete, then unblock:
    task status T-test-4 pending
  ```

**Step 2: Promote**

Run `forge test promote <journey>`. The command will:
- Run the journey's tests automatically
- On pass, replace @feature with @regression in all test files under the journey
- Refuse promotion if any test fails

**Step 3: Verify**

Run `git diff` to confirm only tag changes were made, no other code modifications.

**Step 4: Record**

Mark task completed. T-test-4.5 will run full regression to verify the
promoted scripts integrate cleanly with the existing suite.
