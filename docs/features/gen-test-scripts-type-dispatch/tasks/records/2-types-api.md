---
status: "completed"
started: "2026-05-17 13:15"
completed: "2026-05-17 13:18"
time_spent: "~3m"
---

# Task Record: 2 Create types/api.md instruction file

## Summary
Created types/api.md instruction file for API test script generation with reconnaissance strategy, Fact Table required keys, verification method, generation patterns, and API antipattern guards.

## Changes

### Files Created
- plugins/forge/skills/gen-test-scripts/types/api.md

### Files Modified
无

### Key Decisions
- Included a Required Reads subsection in Reconnaissance Strategy to preserve the SKILL.md source categories (Router files, Config files, API handlers, Auth implementation) with discovery guidance
- Added ROUTE_* as a Fact Table key pattern alongside API_PORT and AUTH_ENDPOINT to support projects where route discovery is the primary reconnaissance outcome
- Added Error Contract Testing subsection under Generation Patterns to ensure API tests verify error response structure, not just happy-path status codes
- Authentication Integration table mirrors the SKILL.md Step 1 auth classification but focuses on API-specific generation strategy (header injection, token caching)

## Test Results
- **Tests Executed**: No
- **Passed**: 0
- **Failed**: 0
- **Coverage**: N/A (task has no tests)

## Acceptance Criteria
- [x] plugins/forge/skills/gen-test-scripts/types/api.md exists
- [x] Frontmatter declares type: api and conventions: [testing-api.md]
- [x] Contains a Reconnaissance Strategy section with API-specific search patterns
- [x] Contains a Fact Table Required Keys section listing minimum keys for API type
- [x] Contains a Verification Method section describing how to confirm the project exposes an API
- [x] Contains a Generation Patterns section describing how API test cases translate to executable scripts
- [x] Contains an API Antipattern Guards section
- [x] At least 3 section headings are unique to this file

## Notes
Documentation task (noTest). Hard rules verified: generation patterns reference profile generate.md for syntax, reconnaissance patterns cite actual grep commands.
