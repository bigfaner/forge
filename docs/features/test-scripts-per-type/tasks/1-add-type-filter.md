---
id: "1"
title: "Add --type filter to gen-test-scripts skill"
priority: "P0"
estimated_time: "1-2h"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Add --type filter to gen-test-scripts skill

## Description

Add a `--type <capability>` argument to the gen-test-scripts skill that filters script generation to a single test type. When specified, the skill skips all other type groups entirely — no Fact Table verification, no locator mapping, no spec generation for non-matching types. Shared infrastructure (helpers, config, auth-setup) is always generated/verified regardless of type filter, using existing idempotent logic.

This is the foundation for per-type task splitting: downstream tasks will create separate T-test-2-ui/T-test-2-api/T-test-2-cli tasks that each invoke gen-test-scripts with `--type`.

## Reference Files
- `docs/proposals/test-scripts-per-type/proposal.md` — Source proposal
- `plugins/forge/skills/gen-test-scripts/SKILL.md` — Skill to modify

## Acceptance Criteria
- [ ] gen-test-scripts accepts `--type <capability>` argument (e.g., `--type tui`, `--type api`, `--type cli`)
- [ ] When `--type` is specified, only test cases of that type are processed (Step 1 grouping filters to the specified type)
- [ ] Fact Table verification (Step 1.5) is skipped for non-matching types
- [ ] Locator mapping (Step 3) is skipped for non-UI types when `--type api` or `--type cli`
- [ ] Spec generation (Step 4) only produces files for the specified type
- [ ] Shared infrastructure (Step 3.5) always runs regardless of `--type` value
- [ ] Without `--type`, behavior is unchanged (generates all types as before)
- [ ] Type value matches profile capability names (e.g., `tui` for go-test, `web-ui` for web-playwright)

## Hard Rules
- MUST NOT break existing behavior when `--type` is not specified
- Shared infrastructure generation MUST remain idempotent (create-only-if-missing, merge-for-helpers)

## Implementation Notes
- The type filter should be applied early in the pipeline — at Step 1 (test case grouping) after capabilities are resolved
- Profile capabilities already determine supported types; the `--type` value must be one of the profile's declared capabilities
- The `--type` argument should be validated against profile capabilities; invalid types should error with a clear message
- For go-test profile: capabilities are `[tui, api, cli]`, so valid types are `tui`, `api`, `cli`
- Key risk from proposal: shared infra race condition — mitigated by existing idempotent writes (create-only-if-missing)
