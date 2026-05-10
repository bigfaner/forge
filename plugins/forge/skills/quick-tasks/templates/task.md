---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "{{PRIORITY}}"
estimated_time: "{{ESTIMATED_TIME}}"
dependencies: [{{DEPENDENCIES}}]
status: pending
breaking: false
mainSession: false
---

# {{ID}}: {{TITLE}}

## Description
{{DESCRIPTION}}

## Reference Files
- `docs/proposals/<slug>/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| {{NEW_FILES}} |

### Modify
| File | Changes |
|------|---------|
| {{MODIFIED_FILES}} |

### Delete
| File | Reason |
|------|--------|
| {{DELETED_FILES}} |

## Acceptance Criteria
{{ACCEPTANCE_CRITERIA}}

## Implementation Notes
{{NOTES}}

## Execution Workflow

1. Write failing tests for each acceptance criterion (RED phase).
   - Command: `just test [scope]`
   - Success: at least one test fails (exit non-zero) confirming test infrastructure works.
   - Failure: no tests discovered or test runner crashes → fix test setup, retry this step.
2. Implement minimum code to make all tests pass (GREEN phase).
   - Command: `just test [scope]`
   - Success: all tests pass (exit 0).
   - Failure: tests still fail → extend implementation, retry this step.
3. Refactor while keeping tests green (REFACTOR phase).
   - Command: `just test [scope]`
   - Success: all tests pass (exit 0) after refactoring.
   - Failure: refactoring broke tests → revert change, retry refactor.
4. Run quality gate in strict sequence, stopping at first failure:
   - `just compile [scope]` — Success: exit 0. Failure: fix compilation errors, retry from this step.
   - `just fmt [scope]` — Success: exit 0. Failure: task is blocked (formatting requires manual review).
   - `just lint [scope]` — Success: exit 0. Failure: self-fix once; if still failing, task is blocked.
   - `just test [scope]` — Success: exit 0. Failure: fix failing tests, retry from compile step above.
5. Stop. Proceed to Step 3 (Record).
