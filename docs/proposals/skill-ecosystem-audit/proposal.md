---
created: 2026-05-17
author: "faner + Claude"
status: Draft
---

# Proposal: Skill Ecosystem Audit — Clarity, Redundancy, Distributability

## Problem

Forge is a distributed Claude Code plugin. Its 18 skills contain structural issues that affect correctness when deployed to end users. The audit is grounded in the actual distribution model: the installed plugin lives at `~/.claude/plugins/cache/forge/forge/<version>/` and contains `agents/`, `commands/`, `hooks/`, `references/`, `scripts/`, `skills/` — all at the top level.

### Key Distribution Insight

Skills are Markdown instructions executed by Claude's AI agent. When a skill says "Read `plugins/forge/skills/eval/rubrics/harness.md`", Claude tries to read from the **user's project root** — not from the plugin cache. The plugin system provides `${CLAUDE_PLUGIN_ROOT}` for hooks/scripts, but skills have no equivalent path resolution mechanism.

Additionally, some directories like `docs/conventions/` and `docs/business-rules/` are **forge-generated user-project specs** that agents must read during task execution. These are NOT forge-internal resources — they're expected output that users must follow.

## Per-Skill Scorecard

| Skill | Lines | Commits | Clarity | Hardcoded Paths | Bloat/Duplication | Verdict |
|-------|-------|---------|---------|-----------------|-------------------|---------|
| brainstorm | 97 | 6 | 5/5 | 0 | Clean | Healthy |
| write-prd | 260 | 13 | 5/5 | 0 | Clean | Healthy |
| ui-design | 313 | 13 | 4/5 | 0 | Clean | Healthy |
| tech-design | 231 | 15 | 5/5 | 1 | Clean | Minor issue |
| eval | 366 | 13 | 3/5 | 2 | 90-line validate-ux inline | Needs simplification |
| breakdown-tasks | 405 | 67 | 4/5 | 0* | 6+ dup templates with quick-tasks | Needs dedup |
| quick-tasks | 168 | 34 | 5/5 | 0* | 6+ dup templates with breakdown-tasks | Needs dedup |
| gen-test-cases | 136 | 20 | 5/5 | 1 | `types/` dir missing | **Bug** |
| gen-test-scripts | 529 | 33 | 4/5 | 2 | **27MB node_modules**, 14 dup Playwright templates | **Critical** |
| run-e2e-tests | 247 | 21 | 5/5 | 1 | Clean | Minor issue |
| graduate-tests | 209 | 11 | 5/5 | 3 | Clean | Minor issue |
| improve-harness | 158 | 10 | 5/5 | 2 | Clean | Minor issue |
| consolidate-specs | 418 | 10 | 4/5 | 1 | Complexity justified | Minor issue |
| submit-task | 199 | 2 | 5/5 | 0 | Clean | Healthy |
| record-task | 175 | 21 | — | — | **DEAD** (deleted from source, 3 stale refs) | Cleanup needed |
| forensic | 194 | 4 | 4/5 | 3+ | Hardcoded user-specific paths | **Bug** |
| learn-lesson | 123 | 4 | 5/5 | 0 | Clean | Healthy |
| init-justfile | 344 | 7 | 4/5 | 12 | Templates distributed with plugin — OK | Healthy |

\* validate-* task templates contain hardcoded paths (not in SKILL.md itself)

## Findings by Priority

### P0 — Breaking / Functionality Failure

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 1 | **27MB node_modules in gen-test-scripts** — 31 packages, exact copies of profile system templates | `gen-test-scripts/templates/node_modules/` + 14 `.ts/.json` files | Every install ships 27MB of dead weight. Templates are bit-identical to `web-playwright` profile. |
| 2 | **gen-test-cases `types/` dir does not exist** — 2 skills reference `plugins/forge/skills/gen-test-cases/types/{type}.md` but files live in `templates/` | `gen-test-cases/SKILL.md:101`, `gen-test-scripts/SKILL.md:50` | Runtime Read failure — agent cannot load per-type conventions |
| 3 | **forensic hardcoded developer paths** — `~/.claude/projects/-Users-fanhuifeng-...` | `forensic/SKILL.md:89,93,100` | Breaks for every user except the original developer. Also has stale build cmd at line 35 referencing `~/.zcode-forge-cli/task` |

### P1 — Architecture / Maintainability

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 4 | **22+ `plugins/forge/` hardcoded paths across skills** — resolve to user's project root, not plugin cache | 12 skill files + 4 task templates (see full inventory below) | Silent failure in distributed plugin — Claude gets "file not found" |
| 5 | **breakdown-tasks ↔ quick-tasks massive duplication** — identical Type Assignment table, Intent Propagation, Step 0 profile resolution, plus 6 near-duplicate template pairs + 5 test-pipeline parallel templates | `breakdown-tasks/SKILL.md:310-330` ≡ `quick-tasks/SKILL.md:96-116`; `templates/{task,task-doc,validate-code-task,validate-ux-task,index.json,index.schema.json}.md` in both | Changes must be manually synced. 11 files to maintain in parallel. |
| 6 | **14 Playwright-specific templates in gen-test-scripts** identical to `web-playwright` profile | `gen-test-scripts/templates/{api.spec.ts,auth-setup.ts,cli.spec.ts,helpers.ts,playwright-ui.spec.ts,playwright.config.ts,package.json,tsconfig.json}` | Dead code. SKILL.md correctly references `{profile-templates-dir}` but local copies remain. |

### P2 — Cleanup / Polish

| # | Issue | Files | Impact |
|---|-------|-------|--------|
| 7 | **record-task stale references** — already deleted from source, but 3 active files still reference it | `agents/task-executor.md:42` (wrong comment), `breakdown-tasks/templates/consolidate-specs.md:88` (invokes `/record-task`), `forge-cli/tests/e2e/justfile_mixed_cli_cli_test.go:283-289` (asserts "record-task" in task-executor.md) | Misleading documentation, failing test |
| 8 | **eval validate-ux sub-pipeline inline** — 90 lines of project-type detection, PRD-to-operation translation, snapshot format, 7 impact types embedded in SKILL.md | `eval/SKILL.md:134-225` | Cognitive overload; would benefit from extraction to separate skill or rubric |
| 9 | **breakdown-tasks Scope Assignment heuristic** — classifies `src/` as frontend, wrong for Go/Rust backend projects | `breakdown-tasks/SKILL.md:290-304` | Incorrect task-to-scope mapping for backend projects |
| 10 | **Naming inconsistency** — `gen` abbreviation in `gen-test-cases`/`gen-test-scripts` vs full verbs (`write`, `run`, `graduate`) | Skill directory names | Minor aesthetic issue |

### Full `plugins/forge/` Path Inventory

All files containing broken `plugins/forge/` references:

```
# References to files DISTRIBUTED with plugin (exist but wrong path)
eval/SKILL.md:85,347                    → plugins/forge/skills/eval/rubrics/<type>.md
tech-design/SKILL.md:191                → plugins/forge/references/shared/decision-logging.md
consolidate-specs/SKILL.md:177          → plugins/forge/references/shared/decision-logging.md
gen-test-scripts/SKILL.md:143           → plugins/forge/references/shared/sitemap.json
run-e2e-tests/SKILL.md:132              → plugins/forge/skills/run-e2e-tests/templates/e2e-report.md
graduate-tests/SKILL.md:120,122,154     → 3 template references
improve-harness/SKILL.md:31             → plugins/forge/skills/eval/rubrics/harness.md
improve-harness/templates/improvements.md:5 → same rubric path
init-justfile/SKILL.md:183-198          → 12 template references (work correctly)
forensic/SKILL.md:89,93,100             → hardcoded user-specific paths

# References to files that DO NOT EXIST
gen-test-cases/SKILL.md:101             → plugins/forge/skills/gen-test-cases/types/{type}.md (dir missing!)
gen-test-scripts/SKILL.md:50            → plugins/forge/skills/gen-test-cases/types/{type}.md (dir missing!)

# In task templates (written to user projects)
breakdown-tasks/templates/validate-ux-task.md:48    → plugins/forge/skills/eval/rubrics/validate-ux.md
breakdown-tasks/templates/validate-code-task.md:45  → plugins/forge/skills/eval/rubrics/validate-code.md
quick-tasks/templates/validate-ux-task.md:48        → plugins/forge/skills/eval/rubrics/validate-ux.md
quick-tasks/templates/validate-code-task.md:45      → plugins/forge/skills/eval/rubrics/validate-code.md
```

Plus 3 command files with hardcoded paths (gen-sitemap, record-decision, extract-design-md).

## Proposed Solution

**Three workstreams**, prioritized by impact:

### W1: Dead Code Cleanup (Low effort, High clarity)

- **Fix 3 stale `record-task` references**:
  - `agents/task-executor.md:42` — correct comment about submit-task
  - `breakdown-tasks/templates/consolidate-specs.md:88` — change to `/submit-task`
  - `forge-cli/tests/e2e/justfile_mixed_cli_cli_test.go:283-289` — update assertion
- **Fix forensic hardcoded paths**: Replace `~/.claude/projects/-Users-fanhuifeng-...` with generic examples or instructions to use `forge forensic search` output. Fix stale build command at line 35.

### W2: Distributability Fixes (Medium effort, Critical correctness)

- **Fix `plugins/forge/` path references**: 22+ references need a path convention that works when the plugin is installed. Options:
  - Use `${CLAUDE_PLUGIN_ROOT}` if available to skills (needs verification)
  - Use relative paths from the skill's own location
  - Define a `<FORGE_ROOT>` placeholder that the skill runner resolves
- **Fix gen-test-cases `types/` path**: Either rename `templates/` to `types/` or update both SKILL.md files to reference `templates/`
- **Remove gen-test-scripts dead weight**:
  - Delete 14 Playwright-specific template files (duplicated in profile system)
  - Delete `node_modules/` (27MB) — move `validate-specs.mjs` to profile tooling or make it a project-level devDependency
  - Extract Playwright auth instructions (lines 359-371) to `generate.md`

### W3: Dedup & Simplification (Medium effort, Long-term health)

- **Extract shared logic from breakdown-tasks ↔ quick-tasks**:
  - Shared Type Assignment table → shared reference file
  - Shared Intent Propagation → shared reference file
  - Shared Step 0 profile resolution → shared reference file
  - 6 near-duplicate templates → shared template files
- **Simplify eval**: Extract 90-line validate-ux inline sub-pipeline to a separate skill or rubric
- **Fix breakdown-tasks scope heuristic**: Correct `src/` classification for backend projects

## Requirements Analysis

### Key Scenarios

- User installs forge plugin → all skills reference files correctly via plugin-relative paths
- User runs `/submit-task` → works identically to old `/record-task` (no residual references)
- User with non-Playwright test framework → gen-test-scripts generates correct templates from profile
- User project has `docs/conventions/` → agents read these correctly as user-facing specs
- `forensic` skill works on any user's machine, not just the developer's

### Constraints & Dependencies

- **Backward compatibility**: All slash commands must work identically after changes
- **Plugin distribution**: Skills must work when forge is installed as a plugin (not as source repo)
- **Test profile system**: Already exists at v3.0.0 — the Playwright abstraction builds on it
- **Existing `skill-rationalization` proposal**: Eval consolidation already done (single eval skill + 16 rubrics). This proposal does not duplicate that work

## Alternatives

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Three workstreams, prioritized** | Addresses all issues; can be phased | Medium total effort | **Selected** |
| Only W1 (cleanup) | Minimal effort, low risk | Leaves path correctness broken | Rejected: only treats symptoms |
| Only W2 (distributability) | Fixes user-facing breakage | Leaves redundancy and complexity | Partial: good first phase |
| Rewrite from scratch | Clean slate | Massive effort, high risk | Rejected: overkill |

## Scope

### In Scope

1. Fix 3 stale `record-task` references
2. Fix forensic hardcoded developer paths + stale build command
3. Fix 22+ `plugins/forge/` path references for installed plugin context
4. Fix gen-test-cases `types/` vs `templates/` path mismatch
5. Remove 14 Playwright template files + 27MB node_modules from gen-test-scripts
6. Extract Playwright auth instructions from gen-test-scripts SKILL.md to generate.md
7. Deduplicate breakdown-tasks ↔ quick-tasks shared logic
8. Simplify eval validate-ux inline sub-pipeline
9. Fix breakdown-tasks scope heuristic for backend projects

### Out of Scope

- Eval skill restructuring (already done)
- Core pipeline skills logic (brainstorm, write-prd, ui-design, tech-design)
- CLI binary changes
- New skill creation
- User-facing command name changes
- `just` or `doc-scorer`/`doc-reviser` references (expected dependencies)
- `init-justfile` template references (work correctly in installed plugin)
- Naming changes (`gen-` prefix)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Path resolution convention doesn't work in all Claude Code versions | M | H | Test with multiple Claude Code versions; provide fallback |
| Moving Playwright templates breaks existing test generation | M | H | Verify gen-test-scripts output is identical before/after for Playwright profile |
| Deduplicating breakdown-tasks/quick-tasks breaks task generation | M | H | Verify task output is identical for both pipeline paths |
| Extracting validate-ux from eval changes eval behavior | L | M | Validate with existing eval-ui test cases |

## Success Criteria

- [ ] Zero stale `record-task` references in active source files
- [ ] Zero hardcoded user-specific paths in forensic
- [ ] Zero `plugins/forge/` hardcoded paths that fail at runtime (all use plugin-relative resolution or correct dir names)
- [ ] `gen-test-cases/types/` path mismatch resolved
- [ ] `gen-test-scripts/templates/node_modules/` no longer exists (27MB reduction)
- [ ] Zero Playwright-specific template files in gen-test-scripts (all in profile system)
- [ ] breakdown-tasks and quick-tasks share common logic via reference files (not copy-paste)
- [ ] All slash commands produce identical behavior to pre-migration
- [ ] Plugin installs and all skills work correctly from installed location
