---
id: "T-test-1"
title: "Generate e2e Test Cases"
priority: "P1"
estimated_time: "1-2h"
dependencies: ["4.gate"]
status: pending
type: "test-pipeline.gen-cases"
---

# T-test-1: Generate e2e Test Cases

## Description

Call `/gen-test-cases` skill to generate structured test case documentation from PRD acceptance criteria.

Each test case includes:
- Source: Specific acceptance criterion from PRD
- Type: CLI (this feature is pure CLI/agent, no UI/API)
- Target: Test target (e.g., cli/task-prompt, cli/task-migrate)
- Test ID: Unique identifier
- Pre-conditions, Steps, Expected, Priority

## Reference Files

- `prd/prd-spec.md` — PRD specification
- `prd/prd-user-stories.md` — User stories (with Given/When/Then acceptance criteria)

## Acceptance Criteria

- [ ] `testing/test-cases.md` file created
- [ ] Each test case includes Target and Test ID fields
- [ ] All test cases traceable to PRD acceptance criteria (7 user stories)
- [ ] Test cases grouped by type (CLI)

## User Stories

No direct user story mapping. This is a standard test generation task.

## Implementation Notes

1. If `docs/sitemap/sitemap.json` does not exist, skip `/gen-sitemap` (pure CLI feature, no UI)
2. Run `/gen-test-cases` skill
3. Verify generated `testing/test-cases.md` contains Target and Test ID fields
