---
id: "T-test-1"
title: "Generate e2e Test Cases"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["3.gate"]
status: pending
---

# Generate e2e Test Cases

## Description

Call `/gen-test-cases` skill to generate structured test case documentation from PRD acceptance criteria.

Each test case includes:
- Source: Specific acceptance criterion from PRD
- Type: UI / API / CLI
- Target: Test target path (e.g., cli/justfile, cli/init-justfile)
- Test ID: Unique identifier
- Pre-conditions, Steps, Expected, Priority

## Reference Files

- `prd/prd-spec.md` — PRD specification
- `prd/prd-user-stories.md` — User stories (with Given/When/Then acceptance criteria)

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
3. Key test scenarios from tech-design: init-justfile adaptive generation, scope parameter validation, project-type output, scope mismatch, skill migration completeness
