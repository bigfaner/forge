---
id: "5"
title: "Update eval skill table + eval-test-cases command"
priority: "P1"
estimated_time: "1h"
dependencies: ["4"]
type: "documentation"
mainSession: false
---

# 5: Update eval skill table + eval-test-cases command

## Description

Add 5 new `test-cases-*` entries to the eval skill's prerequisite/location table and refactor the `eval-test-cases` command into a thin per-type dispatcher loop that invokes the eval skill once per active type with the matching per-type rubric.

## Reference Files
- `docs/proposals/gen-test-cases-per-type-dispatch/proposal.md` — Source proposal
- `plugins/forge/skills/eval/SKILL.md` — Eval skill (prerequisite table + location table to extend)
- `plugins/forge/commands/eval-test-cases.md` — Current thin wrapper to refactor

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/eval/SKILL.md` | Add 5 new entries to Prerequisites table and Default Doc Dir table; add to Parameters --type enum; add to Rubric Reference table; add type-specific pre-processing entry for `test-cases-*` types |
| `plugins/forge/commands/eval-test-cases.md` | Refactor from thin `Skill(forge:eval, args="--type test-cases")` wrapper into per-type dispatcher loop with legacy fallback |

### Delete
| File | Reason |
|------|--------|
| (none) |

## Acceptance Criteria
- [ ] Eval skill Prerequisites table has entries for `ui-test-cases`, `tui-test-cases`, `mobile-test-cases`, `api-test-cases`, `cli-test-cases` — each requiring `testing/{type}-test-cases.md`
- [ ] Eval skill Default Doc Dir table has entries mapping each `test-cases-*` type to `testing/` directory
- [ ] Eval skill Parameters `--type` enum includes all 5 new `test-cases-*` values
- [ ] Eval skill Rubric Reference table has entries for all 5 new rubrics (1000 scale, 900 target, 6 iterations)
- [ ] `eval-test-cases` command: if per-type files exist (`testing/*-test-cases.md`), loop through each active type invoking eval with `--type test-cases-{type}`
- [ ] `eval-test-cases` command: if no per-type files exist, fall back to `--type test-cases` for legacy monolithic mode
- [ ] Pre-processing for `test-cases-*` types: resolve test profile, pass profile capabilities to scorer (same as current `test-cases` pre-processing)

## Hard Rules
- Do NOT modify the eval skill's core scorer→gate→revise loop — only extend the prerequisite/location tables and parameter enum
- The eval-test-cases dispatcher must pass a single `{type}-test-cases.md` file path per invocation, NOT the entire `testing/` directory
- Legacy fallback must work without errors when only `testing/test-cases.md` exists

## Implementation Notes
- The 5 new prerequisite entries follow the same pattern as the existing `test-cases` entry
- The eval-test-cases dispatcher detects active types from the profile manifest (same as gen-test-cases Step 2.5), then loops: for each active type with a matching `testing/{type}-test-cases.md` file, invoke eval skill
- The dispatcher should aggregate per-type scores into a combined report showing pass/fail per type
