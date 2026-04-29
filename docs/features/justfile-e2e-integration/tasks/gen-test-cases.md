---
id: "T-test-1"
title: "Generate e2e Test Cases"
priority: "P1"
estimated_time: "1-2h"
dependencies: [{{T_TEST_1_DEP}}]
status: pending
---

# Generate e2e Test Cases

## Description

Call `/gen-test-cases` skill to generate structured test case documentation from PRD acceptance criteria.

Each test case includes:
- Source: Specific acceptance criterion from PRD
- Type: UI / API / CLI
- Target: Test target path (e.g., ui/login, api/auth)
- Test ID: Unique identifier (e.g., ui/login/login-with-valid-credentials)
- Pre-conditions, Steps, Expected, Priority

## Reference Files

- `prd/prd-spec.md` — PRD specification
- `prd/prd-user-stories.md` — User stories (with Given/When/Then acceptance criteria)
- `prd/prd-ui-functions.md` — UI function requirements (optional)

## Acceptance Criteria

- [ ] `testing/test-cases.md` file created
- [ ] Each test case includes Target and Test ID fields
- [ ] All test cases traceable to PRD acceptance criteria
- [ ] Test cases grouped by type (UI → API → CLI)

## User Stories

No direct user story mapping. This is a standard test generation task.

## Implementation Notes

1. Run `/gen-test-cases` skill
2. Verify generated `testing/test-cases.md` contains Target and Test ID fields
3. If PRD has no UI/API/CLI requirements, mark task as skipped with explanation
