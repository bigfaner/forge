---
id: "2"
title: "Fix argument-hints and add missing argument-hint + arguments + effort"
priority: "P1"
estimated_time: "1h"
dependencies: ["1"]
type: "cleanup"
scope: "all"
breaking: false
mainSession: false
---

# 2: Fix argument-hints and add missing argument-hint + arguments + effort

## Description

Nine command files use a non-standard `argument-hints` field (plural) with complex object structure (name/description/required). The official Claude Code reference documents two separate fields: `argument-hint` (singular, string for autocomplete) and `arguments` (list for `$name` substitution).

Additionally, ~15 skills/commands lack `argument-hint` entirely, and 5 complex skills lack `effort` settings.

## Reference Files
- `docs/proposals/skill-ecosystem-audit/proposal.md` — Source proposal (W1, items 3-5)
- `docs/official-references/skills-ref.md` — Official frontmatter reference

## Affected Files

### Modify

**Fix `argument-hints` → `argument-hint` + `arguments` (9 command files):**

| File | Current | Target |
|------|---------|--------|
| `commands/eval-consistency.md` | Complex `argument-hints` object | `argument-hint: "[--target 900] [--iterations 3] [--scope docs\|full]"` + `arguments: [target, iterations, scope]` |
| `commands/eval-design.md` | Complex `argument-hints` object | `argument-hint: "[--target 900] [--iterations 3]"` + `arguments: [target, iterations]` |
| `commands/eval-prd.md` | Complex `argument-hints` object | `argument-hint: "[--target 900] [--iterations 3]"` + `arguments: [target, iterations]` |
| `commands/eval-proposal.md` | Complex `argument-hints` object | `argument-hint: "[--target 900] [--iterations 3]"` + `arguments: [target, iterations]` |
| `commands/eval-test-cases.md` | Complex `argument-hints` object | `argument-hint: "[--target 900] [--iterations 6]"` + `arguments: [target, iterations]` |
| `commands/eval-ui.md` | Complex `argument-hints` object | `argument-hint: "[--target 950] [--iterations 3]"` + `arguments: [target, iterations]` |
| `commands/extract-design-md.md` | Complex `argument-hints` object | `argument-hint: "[url] [--platform web\|mobile\|tui]"` + `arguments: [url, platform]` |
| `commands/fix-bug.md` | Complex `argument-hints` object | `argument-hint: "[error-msg] [scope]"` + `arguments: [error-msg, scope]` |
| `commands/gen-sitemap.md` | Complex `argument-hints` object | `argument-hint: "[base-url] [api-base-url]"` + `arguments: [base-url, api-base-url]` |
| `commands/git-checkout.md` | Complex `argument-hints` object | `argument-hint: "[source-branch]"` + `arguments: [source-branch]` |
| `commands/git-commit.md` | Complex `argument-hints` object | `argument-hint: "[scope]"` + `arguments: [scope]` |
| `commands/simplify-skill.md` | `argument-hints: skill name` (simple string) | `argument-hint: "[skill-name]"` + `arguments: [skill-name]` |

**Add missing `argument-hint` (skills with no hint):**

| File | `argument-hint` to add |
|------|----------------------|
| `skills/brainstorm/SKILL.md` | `[idea or feature description]` |
| `skills/consolidate-specs/SKILL.md` | `[--slug <feature-slug>]` |
| `skills/eval/SKILL.md` | `[--type <type>] [--target 900] [--iterations 3]` |
| `skills/forensic/SKILL.md` | `[session-id or keywords]` |
| `skills/graduate-tests/SKILL.md` | `[--slug <feature-slug>]` |
| `skills/learn/SKILL.md` | `[decision\|lesson\|convention topic description]` |
| `skills/submit-task/SKILL.md` | `[task-id]` |
| `skills/write-prd/SKILL.md` | `[feature description or requirements]` |

**Add `effort` (complex skills needing deep reasoning):**

| File | `effort` value | Rationale |
|------|---------------|-----------|
| `skills/eval/SKILL.md` | `high` | Multi-round scoring, adversarial revision |
| `skills/forensic/SKILL.md` | `max` | Deep analysis of session transcripts |
| `skills/tech-design/SKILL.md` | `high` | Complex architectural decisions |
| `skills/ui-design/SKILL.md` | `high` | Multi-platform design with prototyping |
| `skills/write-prd/SKILL.md` | `high` | Long collaborative dialogue |

## Acceptance Criteria

- [ ] Zero files contain `argument-hints` (plural)
  `grep -rn 'argument-hints' plugins/forge/` returns 0 hits
- [ ] All eval commands have `argument-hint` with `[--target]` and `[--iterations]`
- [ ] eval, forensic, tech-design, ui-design, write-prd have `effort` set in frontmatter
- [ ] All argument-accepting skills have `argument-hint` field

## Hard Rules

- Only modify frontmatter sections. Do not change skill/command content.
- `arguments` field: only add where the skill content uses positional `$N` or named `$name` substitution. Skills using only `$ARGUMENTS` (full string) do not need `arguments`.
- `effort` values: only `high` or `max`. Do not use `low` or `medium`.

## Implementation Notes

- The `argument-hints` complex object format is entirely non-standard. Replace with simple string `argument-hint`.
- `arguments` list maps to positional parameters for `$name` substitution in content. Only add when the skill body actually uses named references.
- Skills that read from files (breakdown-tasks, quick-tasks, gen-test-cases, gen-test-scripts, run-e2e-tests, improve-harness, tech-design, ui-design) don't need `argument-hint` since they take no user arguments — they derive context from the project file system.
