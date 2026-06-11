---
created: 2026-05-17
author: "faner + Claude"
status: Approved
revised: 2026-05-17
---

# Proposal: Skill Ecosystem Audit — Distributability, Cleanup, Frontmatter, Dedup (v3)

## Problem

Forge is a distributed Claude Code plugin (v3.0.0-beta-7). Three categories of issues remain after recent path migration commits:

1. **Correctness gaps**: 5 hardcoded paths in template files, forensic developer-specific paths
2. **Frontmatter non-compliance**: 18 skills/commands use wrong field names (`allowed_tools` vs `allowed-tools`) or non-standard formats, missing standard fields (`argument-hint`, `arguments`, `effort`)
3. **Maintainability debt**: eval inline pipeline, breakdown-tasks ↔ quick-tasks duplication, scope heuristic

## Resolved Issues

| Issue | Resolution | Commit |
|-------|-----------|--------|
| 42 `plugins/forge/` paths in SKILL.md/command files (37 of 42) | Migrated to `${CLAUDE_SKILL_DIR}` (14 uses) and relative paths | `f2265e1`, `90f91b2`, `8259b18`, `2e9739e` |
| gen-test-cases `types/` directory missing | 5 type files created | Earlier PR |
| 27MB `node_modules/` in gen-test-scripts | Directory removed | Earlier PR |
| gen-test-scripts monolithic SKILL.md | Refactored to dispatcher pattern | #104 |
| `record-decision` command (superseded by `/learn`) | Command removed entirely | `9a42c39` |
| Skills reference docs outdated | Synced from upstream Claude Code docs | `54d5107` |

**Key finding from reference docs**: `${CLAUDE_SKILL_DIR}` is the official mechanism for referencing plugin-internal files from SKILL.md content. It resolves to the directory containing the SKILL.md file.

## Remaining Issues

### P0 — Correctness

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 1 | **`allowed_tools` (underscore) instead of `allowed-tools` (hyphen)** | 15 files (13 commands + 2 skills) | Permission grants may silently not work — official field name is `allowed-tools` |
| 2 | **5 `plugins/forge/` paths in template files** | 4 task templates + 1 report template | Task validators reference rubrics via broken paths |
| 3 | **Forensic hardcoded developer paths** (4 occurrences) | `forensic/SKILL.md:89,93,100` | Breaks for every user except the original developer |

### P1 — Frontmatter Compliance

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 4 | **`allowed-tools` format uses JSON arrays** (`["Bash", "Read"]`) instead of space-separated string or YAML list | 17 files | Non-standard format; may not parse correctly |
| 5 | **`argument-hints` (plural, complex object) instead of `argument-hint` (singular, string)** | 9 command files | Non-standard field; autocomplete hint won't display |
| 6 | **Missing `argument-hint`** on skills/commands that accept arguments | ~15 files | No autocomplete guidance for users |
| 7 | **Missing `effort`** on complex skills (eval, forensic, tech-design, ui-design, write-prd) | ~5 files | These skills need deep reasoning; default effort may be insufficient |

### P2 — Maintainability

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 8 | **1 stale `record-task` reference** | `breakdown-tasks/templates/consolidate-specs.md:88` | Misleading documentation |
| 9 | **8 Playwright template files** (dead code) | `gen-test-scripts/templates/` | Profile system already provides these |
| 10 | **breakdown-tasks ↔ quick-tasks duplication** | Both skill directories | Changes must be manually synced |
| 11 | **eval validate-ux sub-pipeline inline** (90+ lines) | `eval/SKILL.md:134-225` | Cognitive overload |
| 12 | **breakdown-tasks scope heuristic** classifies `src/` as frontend | `breakdown-tasks/SKILL.md:290` | Wrong for Go/Rust backend projects |

### Frontmatter Audit Detail

#### Field Name Corrections

```
allowed_tools → allowed-tools  (15 files)

Commands: clean-code, execute-task, extract-design-md, fix-bug, gen-sitemap,
          git-checkout, git-commit, init-forge, quick, run-tasks, simplify-skill
Skills:   clean-code/SKILL.md, init-justfile/SKILL.md
```

#### Format Corrections

```
["Bash", "Read"] → Bash Read  (17 files)

All files with allowed-tools use JSON array syntax.
init-justfile uses single-quote array ['Bash', 'Read'].
Correct format: space-separated string or YAML list.
```

#### `argument-hints` → `argument-hint` + `arguments`

```
Current (9 command files use complex object):
  argument-hints:
    - name: target
      description: Target score threshold
      required: false

Should be:
  argument-hint: "[--target 900] [--iterations 3] [--scope docs|full]"
  arguments: [target, iterations, scope]
```

Files: eval-consistency, eval-design, eval-prd, eval-proposal, eval-test-cases, eval-ui, extract-design-md, fix-bug, gen-sitemap, git-checkout, git-commit, simplify-skill

#### Missing `argument-hint` (skills)

| Skill | Suggested `argument-hint` |
|-------|--------------------------|
| brainstorm | `[idea or feature description]` |
| breakdown-tasks | (no args — reads from tech-design) |
| consolidate-specs | `[--slug <feature-slug>]` |
| eval | `[--type <type>] [--target 900] [--iterations 3]` |
| forensic | `[session-id or keywords]` |
| gen-test-cases | (no args — reads from PRD) |
| gen-test-scripts | (no args — reads from test-cases) |
| graduate-tests | `[--slug <feature-slug>]` |
| improve-harness | (no args — reads from eval report) |
| learn | `[decision\|lesson\|convention topic description]` |
| quick-tasks | (no args — reads from proposal) |
| run-e2e-tests | (no args — reads from test scripts) |
| submit-task | `[task-id]` |
| tech-design | (no args — reads from PRD) |
| ui-design | (no args — reads from PRD) |
| write-prd | `[feature description or requirements]` |

#### Missing `effort`

| Skill | Suggested `effort` | Rationale |
|-------|-------------------|-----------|
| eval | `high` | Multi-round scoring, adversarial revision |
| forensic | `max` | Deep analysis of session transcripts |
| tech-design | `high` | Complex architectural decisions |
| ui-design | `high` | Multi-platform design with prototyping |
| write-prd | `high` | Long collaborative dialogue, stakeholder analysis |

#### Undocumented Fields

| Field | Files | Action |
|-------|-------|--------|
| `conventions` | gen-test-cases/SKILL.md, gen-test-scripts/SKILL.md | Keep — used by consolidate-specs to detect relevant skills; not in reference docs but functionally valid |

## Proposed Solution

### Dual Strategy (proven, already applied)

- **SKILL.md and command files**: Use `${CLAUDE_SKILL_DIR}` for cross-skill references.
- **Template and non-SKILL files**: Use relative-from-self paths.

### Four Workstreams

#### W1: Frontmatter Compliance (~2h, P0/P1 correctness)

**Fix all frontmatter fields across 18 skills and 18 commands.**

For each file:

1. **`allowed_tools` → `allowed-tools`**: Rename field (15 files)
2. **Format**: `["Bash", "Read"]` → `Bash Read` (17 files)
3. **`argument-hints` → `argument-hint`**: Replace complex object with simple string (9 command files). Add `arguments` list only where skill content uses named `$name` substitution.
4. **Add `argument-hint`**: To skills/commands that accept arguments but lack it (~15 files)
5. **Add `effort`**: To complex skills that need deep reasoning (5 files)

#### W2: Cleanup (~1h, High clarity)

- **Fix forensic hardcoded paths**: Replace 4 developer-specific paths with generic instructions using `${CLAUDE_SESSION_ID}` or `forge forensic` CLI commands.
- **Fix 1 stale `record-task` reference**: `consolidate-specs.md:88` → `/submit-task`.

#### W3: Distributability Completion (~1h, Critical correctness)

- **Fix 5 `plugins/forge/` paths in templates**: Replace with relative-from-self paths.
- **Remove 8 Playwright template files** from `gen-test-scripts/templates/`. Keep `validate-specs.mjs`, `validate-specs.test.mjs`, and `__test_fixtures__/`.

#### W4: Dedup & Simplification (~4h, Long-term health)

- **Extract shared logic from breakdown-tasks ↔ quick-tasks** to `references/shared/`.
- **Extract eval validate-ux sub-pipeline** to rubric file.
- **Fix breakdown-tasks scope heuristic** for backend projects.

## Scope

### In Scope

1. Fix `allowed_tools` → `allowed-tools` field name (15 files) (W1)
2. Fix `allowed-tools` format to space-separated string (17 files) (W1)
3. Fix `argument-hints` → `argument-hint` + `arguments` (9 files) (W1)
4. Add `argument-hint` to ~15 skills/commands (W1)
5. Add `effort` to 5 complex skills (W1)
6. Fix forensic hardcoded developer paths (W2)
7. Fix 1 stale `record-task` reference (W2)
8. Fix 5 `plugins/forge/` paths in templates (W3)
9. Remove 8 Playwright template files (W3)
10. Extract shared breakdown-tasks ↔ quick-tasks logic (W4)
11. Extract eval validate-ux sub-pipeline (W4)
12. Fix breakdown-tasks scope heuristic (W4)

### Out of Scope

- CI path-guard linter
- Core pipeline skills logic changes
- CLI binary changes
- New skill creation
- `init-justfile` template path references (work correctly)
- `conventions` field in gen-test-cases/gen-test-scripts (undocumented but functional — keep)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `allowed_tools` (underscore) is actually accepted by Claude Code | M | M | Test: create a skill with `allowed_tools`, invoke it, verify tool access. If both forms work, rename still improves spec compliance. |
| `argument-hint` + `arguments` changes break existing command invocation | L | M | The `$ARGUMENTS` and `$N` substitution works with or without `arguments` field. Adding `arguments` only enables `$name` aliases. |
| `effort: high` causes slower responses on simple invocations | L | L | Effort only applies while skill is active. User can override per-invocation. |
| Relative-from-self paths in templates don't resolve | L | M | Claude's Read tool normalizes paths. 3-level traversal is straightforward. |
| Removing Playwright templates breaks test generation | L | H | Profile system + types/ dispatch already generate templates. Dead code confirmed. |
| Dedup breaks task generation | M | H | Generate tasks via both paths before/after. Diff outputs. |

## Success Criteria

- [ ] Zero `allowed_tools` (underscore) in any skill/command frontmatter
  `grep -rn 'allowed_tools' plugins/forge/skills/ plugins/forge/commands/` returns 0 hits
- [ ] All `allowed-tools` values use space-separated format (no JSON arrays)
  `grep -rn 'allowed-tools: \[' plugins/forge/skills/ plugins/forge/commands/` returns 0 hits
- [ ] Zero `argument-hints` (plural) in any frontmatter
  `grep -rn 'argument-hints' plugins/forge/` returns 0 hits
- [ ] All argument-accepting skills/commands have `argument-hint`
- [ ] eval, forensic, tech-design, ui-design, write-prd have `effort` set
- [ ] Zero `plugins/forge/` hardcoded paths in any source file
- [ ] Zero hardcoded user-specific paths in forensic
- [ ] Zero stale `record-task` references
- [ ] `gen-test-scripts/templates/` contains no `.ts`/`.json` files
- [ ] Shared reference files exist for breakdown-tasks ↔ quick-tasks
- [ ] eval/SKILL.md under 280 lines
- [ ] breakdown-tasks scope heuristic handles Go/Rust `src/` correctly
- [ ] All slash commands produce identical behavior to pre-migration
