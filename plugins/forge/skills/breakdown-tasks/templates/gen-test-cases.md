---
id: "T-test-1"
title: "Generate e2e Test Cases"
priority: "P1"
estimated_time: "1-2h"
dependencies: [{{T_TEST_1_DEP}}]
status: pending
mainSession: false
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

1. If `docs/sitemap/sitemap.json` does not exist, run `/gen-sitemap` first
2. Run `/gen-test-cases` skill
3. Verify generated `testing/test-cases.md` contains Target and Test ID fields
4. If PRD has no UI/API/CLI requirements, mark task as skipped with explanation

## Execution Workflow

1. Verify sitemap exists.
   - Check: `docs/sitemap/sitemap.json` file exists.
   - Success: file found, proceed to step 2.
   - Failure: run `/gen-sitemap` skill, then re-check.
2. Run `/gen-test-cases` skill to generate test case documentation.
   - Action: invoke the skill; it writes `testing/test-cases.md`.
   - Success: `testing/test-cases.md` created with Target and Test ID fields present.
   - Failure: if PRD has no UI/API/CLI requirements, mark task as skipped with explanation.
3. Verify generated output.
   - Check: `testing/test-cases.md` exists and contains at least one test case with Target and Test ID fields.
   - Success: file exists, fields populated.
   - Failure: re-run the skill or fix the output manually.
4. Stop. Proceed to Step 3 (Record).
