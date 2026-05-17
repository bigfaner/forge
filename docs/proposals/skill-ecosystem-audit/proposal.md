---
created: 2026-05-17
author: "faner + Claude"
status: Approved
revised: 2026-05-17
---

# Proposal: Skill Ecosystem Audit — Distributability, Cleanup, Dedup (v2)

## Problem

Forge is a distributed Claude Code plugin (v3.0.0-beta-7). Three prior issues have been resolved by recent PRs (#104 gen-test-scripts dispatcher, node_modules removal, gen-test-cases/types/ creation), but **7 issues remain** that affect correctness when deployed to end users.

The core issue: **42 hardcoded `plugins/forge/` path references** across 18 files. When a skill runs from the installed plugin cache (`~/.claude/plugins/cache/forge/forge/<version>/`), these paths resolve to the user's project root — not the plugin directory — causing silent "file not found" failures.

## Resolved Issues (since v1)

| Issue | Resolution | PR |
|-------|-----------|-----|
| gen-test-cases `types/` directory missing | 5 type files created | Earlier PR |
| 27MB `node_modules/` in gen-test-scripts | Directory removed | Earlier PR |
| gen-test-scripts monolithic SKILL.md | Refactored to dispatcher pattern (231 lines, `types/` dispatch) | #104 |

## Remaining Issues

### P0 — Breaking / Functionality Failure

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 1 | **42 `plugins/forge/` hardcoded paths** across 18 files | 12 skill files + 4 task templates + 3 command files | Claude gets "file not found" when reading rubrics, references, templates from installed plugin |
| 2 | **Forensic hardcoded developer paths** — `~/.claude/projects/-Users-fanhuifeng-...` (4 occurrences) | `forensic/SKILL.md:89,93,100` | Breaks for every user except the original developer |

### P1 — Architecture / Maintainability

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 3 | **8 Playwright template files in gen-test-scripts** — `api.spec.ts`, `auth-setup.ts`, `cli.spec.ts`, `helpers.ts`, `package.json`, `playwright.config.ts`, `playwright-ui.spec.ts`, `tsconfig.json` | `gen-test-scripts/templates/` | Dead code; profile system already provides these; SKILL.md dispatches to `types/` not local templates |
| 4 | **breakdown-tasks ↔ quick-tasks duplication** — identical Type Assignment table, Intent Propagation, Step 0 profile resolution, plus 6 near-duplicate template pairs | `breakdown-tasks/SKILL.md:310-330` ≡ `quick-tasks/SKILL.md:96-116`; shared templates in both | Changes must be manually synced; 11 files maintained in parallel |

### P2 — Cleanup / Polish

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 5 | **2 stale `record-task` references** | `agents/task-executor.md:42`, `breakdown-tasks/templates/consolidate-specs.md:88` | Misleading documentation |
| 6 | **eval validate-ux sub-pipeline inline** — 90+ lines of project-type detection, PRD-to-operation translation, snapshot format, 7 impact types embedded in SKILL.md | `eval/SKILL.md:134-225` | Cognitive overload; inconsistent with gen-test-scripts extraction pattern |
| 7 | **breakdown-tasks scope heuristic** — classifies `src/` as frontend, wrong for Go/Rust backend projects | `breakdown-tasks/SKILL.md:290-304` | Incorrect task-to-scope mapping for backend projects |

### Full `plugins/forge/` Path Inventory

```
# Skills (12 files, 30 occurrences)
eval/SKILL.md:85,347                         → plugins/forge/skills/eval/rubrics/<type>.md
tech-design/SKILL.md:191                     → plugins/forge/references/shared/decision-logging.md
consolidate-specs/SKILL.md:177               → plugins/forge/references/shared/decision-logging.md
gen-test-scripts/SKILL.md                    → 6 occurrences (references, sitemap)
run-e2e-tests/SKILL.md:132                   → plugins/forge/skills/run-e2e-tests/templates/e2e-report.md
graduate-tests/SKILL.md:120,122,154           → 3 template references
improve-harness/SKILL.md:31                  → plugins/forge/skills/eval/rubrics/harness.md
gen-test-cases/SKILL.md:101                  → plugins/forge/skills/gen-test-cases/types/{type}.md
forensic/SKILL.md                            → 3 occurrences (hardcoded user paths)
init-justfile/SKILL.md                       → 11 occurrences (WORK CORRECTLY — exempt)

# Commands (3 files, 7 occurrences)
record-decision.md                           → 3 occurrences
gen-sitemap.md                               → 2 occurrences
extract-design-md.md                         → 2 occurrences

# Task templates (4 files, 4 occurrences)
breakdown-tasks/templates/validate-ux-task.md:48
breakdown-tasks/templates/validate-code-task.md:45
quick-tasks/templates/validate-ux-task.md:48
quick-tasks/templates/validate-code-task.md:45

# Already resolved: gen-test-cases types/ path now exists ✓
```

**Note**: `init-justfile` has 11 `plugins/forge/` references that work correctly in the installed plugin (templates are loaded relative to the skill directory during init). These are **exempt** from path fixes.

## Why Now

v3.0.0-beta-7 is in pre-release. The 42 path references are a known correctness gap — every skill that cross-references another skill's resources silently fails in the installed plugin. The recent gen-test-scripts dispatcher refactor (#104) proved the `types/` dispatch pattern works; this proposal applies the same "fix paths + clean dead code" pattern to the remaining skills.

## Proposed Solution

**Three workstreams**, prioritized by impact:

### W1: Dead Code Cleanup (~1.5h, High clarity)

- **Fix 2 stale `record-task` references**:
  - `agents/task-executor.md:42` — correct comment to reference `submit-task`
  - `breakdown-tasks/templates/consolidate-specs.md:88` — change `/record-task` to `/submit-task`
- **Fix forensic hardcoded paths**: Replace 4 `~/.claude/projects/-Users-fanhuifeng-...` references with generic instructions or `forge forensic` CLI commands. No user-specific paths should appear in distributed skills.

### W2: Distributability Fixes (~5h, Critical correctness)

- **Fix `plugins/forge/` path references**: Replace all 31 non-exempt occurrences with **relative-from-self paths**.

  **Approach**: Skills instruct Claude to resolve paths relative to the skill file's own location. For example, `plugins/forge/skills/eval/rubrics/harness.md` becomes `../eval/rubrics/harness.md` (relative to the referencing skill file).

  **Why relative-from-self**: No build step, no runtime variable, matches npm convention. `${CLAUDE_PLUGIN_ROOT}` is only available to hooks/scripts, not skills.

  **Grouping for implementation**:
  - Skills: 30 occurrences across 11 files (excluding init-justfile)
  - Commands: 7 occurrences across 3 files
  - Task templates: 4 occurrences across 4 files
- **Remove 8 Playwright template files** from `gen-test-scripts/templates/`: `.ts`, `.json` files that are dead code since the dispatcher refactor. The profile system and `types/` dispatch handle all test generation. Keep `validate-specs.mjs`, `validate-specs.test.mjs`, and `__test_fixtures__/` (used by the validation pipeline).

### W3: Dedup & Simplification (~4h, Long-term health)

- **Extract shared logic from breakdown-tasks ↔ quick-tasks**:
  - Shared Type Assignment table → `references/shared/type-assignment.md`
  - Shared Intent Propagation → `references/shared/intent-propagation.md`
  - Shared Step 0 profile resolution → `references/shared/step0-profile-resolution.md`
  - Both skills reference these shared files instead of duplicating inline
- **Simplify eval**: Extract 90+ line validate-ux inline sub-pipeline to `skills/eval/rubrics/validate-ux-pipeline.md` (not a new skill — a rubric file consumed by eval)
- **Fix breakdown-tasks scope heuristic**: Correct `src/` classification for backend projects by checking for Go/Rust indicators before classifying as frontend

## Alternatives

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Three workstreams, prioritized** | Addresses all 7 remaining issues; sequenced by risk | ~10.5h total effort | **Selected** |
| Only W1+W2 (skip dedup) | Fixes all breaking issues | Leaves duplication and cognitive load in eval | Partial: acceptable if time-constrained |
| Inline all referenced content | Zero cross-file dependencies | Increases skill file sizes significantly; duplicates content already in rubrics | Rejected: trades one problem for another |
| Build-time path rewriting | Zero runtime changes | Requires build step; fragile; doesn't fix dead code | Rejected: adds complexity without addressing root cause |

## Scope

### In Scope

1. Fix 2 stale `record-task` references (W1)
2. Fix forensic hardcoded developer paths — 4 occurrences (W1)
3. Fix 31 `plugins/forge/` path references with relative-from-self (W2)
4. Remove 8 Playwright template files from gen-test-scripts (W2)
5. Extract shared breakdown-tasks ↔ quick-tasks logic to reference files (W3)
6. Extract eval validate-ux sub-pipeline to rubric file (W3)
7. Fix breakdown-tasks scope heuristic for backend projects (W3)

### Out of Scope

- CI path-guard linter (deferred to separate follow-up)
- Core pipeline skills logic (brainstorm, write-prd, ui-design, tech-design)
- CLI binary changes
- New skill creation
- User-facing command name changes
- `init-justfile` template references (work correctly in installed plugin)
- Naming changes (`gen-` prefix)
- Already-resolved issues (gen-test-cases/types/, node_modules, gen-test-scripts dispatcher)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Relative-from-self resolution does not work at plugin cache path | M | H | Each path fix includes a fallback instruction in the skill: "If this file cannot be found, proceed using criteria described inline." Critical content has a compressed inline summary as degradation fallback. |
| Removing Playwright templates breaks test generation | L | H | Profile system + types/ dispatch already generate templates. The local copies are confirmed dead code by the dispatcher refactor in PR #104. |
| Deduplicating breakdown-tasks/quick-tasks breaks task generation | M | H | Generate tasks for the same test PRD via both paths before and after dedup. Diff the outputs. Identical = safe. |
| Extracting validate-ux from eval changes scoring behavior | L | M | Run eval on existing test cases before/after. Compare scores. |
| Cross-platform path separator issues (Windows backslash) | L | M | All relative paths use forward slashes. Claude's Read tool normalizes separators. |

## Success Criteria

- [ ] Zero stale `record-task` references in active source files
  `grep -r 'record-task' agents/ skills/ templates/` returns 0 hits
- [ ] Zero hardcoded user-specific paths in forensic
  `grep -rn '~/.claude/projects/' skills/forensic/SKILL.md` returns 0 hits
- [ ] Zero `plugins/forge/` hardcoded paths that fail at runtime
  `grep -rn 'plugins/forge/' skills/ agents/ references/ commands/ | grep -v 'init-justfile'` returns 0 hits
- [ ] `gen-test-scripts/templates/` contains no Playwright-specific `.ts`/`.json` files
  `ls skills/gen-test-scripts/templates/*.ts skills/gen-test-scripts/templates/*.json 2>/dev/null` returns empty
- [ ] Shared reference files exist and are referenced by both breakdown-tasks and quick-tasks
  `ls references/shared/type-assignment.md references/shared/intent-propagation.md references/shared/step0-profile-resolution.md` succeeds
- [ ] eval/SKILL.md under 280 lines (was 366)
- [ ] breakdown-tasks scope heuristic checks for Go/Rust indicators before classifying `src/` as frontend
- [ ] All 12 slash commands produce identical behavior to pre-migration
