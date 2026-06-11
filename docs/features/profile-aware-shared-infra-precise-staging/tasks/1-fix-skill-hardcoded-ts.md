---
id: "1"
title: "Fix SKILL.md hardcoded .ts references in Step 3.5 and post-generation"
priority: "P0"
estimated_time: "30m"
dependencies: []
type: "documentation"
mainSession: false
---

# 1: Fix SKILL.md hardcoded .ts references in Step 3.5 and post-generation

## Description

gen-test-scripts SKILL.md contains 3 hardcoded `.ts` filename references that mislead AI agents into creating/referencing TypeScript files even when the active profile is `go-test`. The profile manifest already provides the correct filenames via `templates.helpers` — these hardcoded values bypass that mechanism.

## Reference Files
- `docs/proposals/profile-aware-shared-infra-precise-staging/proposal.md` — Source proposal
- `docs/lessons/gotcha-gen-test-scripts-ts-residue.md` — Incident analysis

## Affected Files

### Create
| File | Description |
|------|-------------|

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/gen-test-scripts/SKILL.md` | Replace 3 hardcoded `.ts` references with profile-manifest-derived references |

### Delete
| File | Reason |
|------|--------|

## Acceptance Criteria

- [ ] Line ~180 (Task Splitting Guard HARD-RULE): Replace "auth-setup.ts, playwright.config.ts, helpers.ts" with generic description referencing profile manifest templates
- [ ] Line ~363 (Auth Infrastructure point 4): Replace "Verify `helpers.ts`" with "Verify the profile's helpers file (from `manifest.templates.helpers`)"
- [ ] Line ~502 (Post-generation helper merge): Replace "verify `helpers.ts`" with "verify the profile's helpers file (from `manifest.templates.helpers`)"
- [ ] grep for remaining hardcoded `.ts` references in Step 3.5 auth section and post-generation sections returns zero results (excluding code examples and the playwright-specific branch in points 2-3 which are correctly scoped)
- [ ] web-playwright profile behavior is unchanged — the references still resolve to `helpers.ts` via manifest

## Hard Rules

- Do NOT modify Steps outside of Step 3.5 auth section (lines ~352-379) and post-generation helper merge (line ~502) and Task Splitting Guard (line ~180)
- Do NOT change the playwright-specific branch (lines 359-361) — those are correctly scoped to "For playwright profile"

## Implementation Notes

- The go-test profile manifest defines `templates.helpers: helpers.go` — after this fix, the skill should reference that file instead of hardcoded `helpers.ts`
- Keep the "For playwright profile:" branch at lines 359-361 as-is — those `.ts` references are correctly scoped to a specific profile
- The Task Splitting Guard (line 180) is descriptive text, not a direct instruction — make it generic enough to cover all profiles without listing specific filenames
