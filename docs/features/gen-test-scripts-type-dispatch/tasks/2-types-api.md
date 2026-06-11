---
id: "2"
title: "Create types/api.md instruction file"
priority: "P1"
estimated_time: "1-2h"
dependencies: []
type: "documentation"
mainSession: false
---

# 2: Create types/api.md instruction file

## Description

Extract API-specific generation logic from the monolithic `gen-test-scripts` SKILL.md into a dedicated `types/api.md` type instruction file. This file guides the agent when generating API test scripts — HTTP reconnaissance, endpoint Fact Table, request/response generation patterns, status code assertions.

Modeled after `gen-test-cases/types/api.md` structure.

## Reference Files
- `docs/proposals/gen-test-scripts-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-cases/types/api.md` — Reference structure (gen-test-cases API type file)
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Source of API-specific content to extract

## Affected Files

### Create
| File | Description |
|------|-------------|
| `plugins/forge/skills/gen-test-scripts/types/api.md` | API type instruction file with conventions frontmatter |

## Acceptance Criteria

- [ ] `plugins/forge/skills/gen-test-scripts/types/api.md` exists
- [ ] Frontmatter declares `type: api` and `conventions: [testing-api.md]`
- [ ] Contains a **Reconnaissance Strategy** section with API-specific search patterns (grep router handlers, endpoint definitions, HTTP method bindings, response schemas)
- [ ] Contains a **Fact Table Required Keys** section listing minimum keys for API type (API_PORT or route path entries like AUTH_ENDPOINT)
- [ ] Contains a **Verification Method** section describing how to confirm the project exposes an API (grep router/handler/endpoint patterns in Go/TS/JS)
- [ ] Contains a **Generation Patterns** section describing how API test cases translate to executable scripts (HTTP client usage, request construction, response assertion, status code checks, header validation)
- [ ] Contains an **API Antipattern Guards** section (hardcoded URLs, missing error contract tests, vacuous "returns success" assertions)
- [ ] At least 3 section headings are unique to this file (not shared with other type files)

## Hard Rules

- Generation patterns must reference the profile's `generate.md` for framework-specific HTTP client syntax — type files describe *what* to generate, not the exact import syntax
- Reconnaissance patterns must cite actual grep commands or search strategies

## Implementation Notes

- Current SKILL.md "Required reads" table has: Router files, Config files, API handlers, Auth implementation — extract these into the reconnaissance strategy
- Fact Table example (lines 263-268) shows API_PORT and AUTH_ENDPOINT — use as the basis for required keys
- Auth classification (lines 207-236) has API-relevant patterns (API key, Bearer token) — reference from generation patterns
- Step 4 API verification (line 391) uses grep for handler patterns — extract into verification method
