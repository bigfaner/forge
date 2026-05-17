---
id: "6"
title: "Migrate all 10 skills to forge testing CLI"
priority: "P1"
estimated_time: "2h"
dependencies: ["4"]
scope: "all"
breaking: false
type: "refactor"
mainSession: false
---

# 6: Migrate all 10 skills to forge testing CLI

## Description
Update all 10 skills listed in proposal D6 table to use `forge testing` CLI commands instead of `forge profile`. Replace `capabilities` references with `interfaces`. Migrate each skill according to its specific change from the D6 table. This is the primary skill-facing migration work.

## Reference Files
- `docs/proposals/simplify-testing-model/proposal.md` — Source proposal (D6: Skills impact table)
- `plugins/forge/skills/` — All skill directories

## Acceptance Criteria
- `grep -r "forge profile" plugins/` returns zero matches (excluding test files)
- `grep -rE "\bcapabilities\b" plugins/` returns zero matches in skill instruction text (excluding test files and comments)
- All 10 skills call `forge testing` commands:
  - gen-test-scripts: `forge testing get generate` (no profile name argument)
  - run-e2e-tests: `forge testing get run`
  - graduate-tests: `forge testing get graduate`
  - gen-test-cases: uses `interfaces` terminology
  - breakdown-tasks: per-language task expansion (not per-profile)
  - quick-tasks: per-language task expansion (not per-profile)
  - init-justfile: `forge testing get justfile`
  - tech-design: no profile selection step, auto-detect
- No skill references v2 profile names (go-test, web-playwright, rust-test, pytest, java-junit, maestro)

## Hard Rules
- Each skill must be tested individually after migration — at minimum, verify the skill loads without error
- Do not add backward-compat fallbacks for `forge profile` in skills — clean migration only
- `eval-test-cases` dynamic scoring dimension changes from capabilities → interfaces

## Implementation Notes
- D6 table maps each skill to its specific change. Start with the most-used skills (gen-test-scripts, run-e2e-tests) and work through the list.
- `grep -r "forge profile" plugins/` before and after serves as the definitive completeness check
