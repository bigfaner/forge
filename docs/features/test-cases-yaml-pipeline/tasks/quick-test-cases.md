---
id: "T-quick-1"
title: "Generate Quick Test Cases"
priority: "P1"
estimated_time: "30min-1h"
dependencies: ["3"]
status: pending
noTest: true
mainSession: false
---

# Generate Quick Test Cases

## Description

Generate structured test cases from the proposal's Success Criteria.

Each test case includes:
- Source: Specific success criterion from proposal
- Type: UI / API / CLI
- Target: Test target path (e.g., ui/login, api/auth)
- Test ID: Unique identifier (e.g., ui/login/login-with-valid-credentials)
- Pre-conditions, Steps, Expected, Priority

## Reference Files

- `docs/proposals/test-cases-yaml-pipeline/proposal.md` — Source proposal with Success Criteria

## Acceptance Criteria

- [ ] `testing/test-cases.md` file created in `docs/features/test-cases-yaml-pipeline/testing/`
- [ ] Each test case includes Target and Test ID fields
- [ ] All test cases traceable to proposal Success Criteria
- [ ] Test cases grouped by type (UI → API → CLI)

## Implementation Notes

1. Read `docs/proposals/test-cases-yaml-pipeline/proposal.md`, extract Success Criteria section
2. Call `/gen-test-cases` skill — provide proposal Success Criteria as the source
3. No sitemap prerequisite (quick mode skips `/gen-sitemap`)
4. Each Success Criterion checkbox becomes one or more test cases
5. If proposal has no testable criteria, mark task as skipped with explanation
