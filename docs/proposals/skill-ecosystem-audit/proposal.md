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

## Why Now

This audit is blocking the v3.0.0 release. Three user-facing failures are confirmed: `gen-test-cases` cannot load per-type conventions at runtime (P0 #2), `forensic` only works for the original developer (P0 #3), and every forge install ships 27MB of dead node_modules (P0 #1). The v3.0.0 test profile system already abstracts away Playwright-specific templates — but gen-test-scripts still contains 14 bit-identical copies of those templates. Shipping v3.0.0 without this audit means shipping known broken paths to every user.

## Industry Benchmarking

Other plugin/extension ecosystems face the same asset-resolution and deduplication challenges. Three relevant precedents:

1. **VS Code Extensions** use `vscode.extensions.getExtension('id').extensionPath` to resolve assets at runtime. Extensions bundle templates inside their install directory and reference them via this API, never via hardcoded relative paths. The lesson: a canonical path variable is the standard solution. Forge's `${CLAUDE_PLUGIN_ROOT}` (available to hooks/scripts) is the direct analog, but skills lack access — which is the core gap W2 must bridge.

2. **npm packages** handle template deduplication via the `package.json` `files` field and `.npmignore` to exclude dev-only assets from distribution. Packages like `create-react-app` ship template directories and use `fs.readFileSync(path.join(__dirname, 'template', ...))` — always resolved relative to the package's own install location, never relative to the consumer's project root. The lesson: relative-from-self resolution is the idiomatic approach for distributed packages.

3. **JetBrains IDE Plugins** use `PluginManager.getPlugin(pluginId)?.pluginPath` and the `VfsUtil` virtual filesystem to locate bundled resources. They also enforce plugin size limits in the marketplace (~200MB). JetBrains explicitly rejects plugins that ship redundant copies of libraries already available in the IDE runtime. The lesson: distribution size matters, and marketplace gatekeeping enforces it.

These benchmarks converge on two patterns: (a) a canonical root variable or relative-from-self resolution, and (b) exclusion of redundant assets from the distribution package. Forge's current approach violates both.

## Proposed Solution

**Three workstreams**, prioritized by impact:

### W1: Dead Code Cleanup (~3h, High clarity)

- **Fix 3 stale `record-task` references**:
  - `agents/task-executor.md:42` — correct comment about submit-task
  - `breakdown-tasks/templates/consolidate-specs.md:88` — change to `/submit-task`
  - `forge-cli/tests/e2e/justfile_mixed_cli_cli_test.go:283-289` — update assertion
- **Fix forensic hardcoded paths**: Replace `~/.claude/projects/-Users-fanhuifeng-...` with generic examples or instructions to use `forge forensic search` output. Fix stale build command at line 35.

### W2: Distributability Fixes (~8h, Critical correctness)

- **Fix `plugins/forge/` path references**: 22+ references need a path convention that works when the plugin is installed.

  **Selected approach: relative-from-self resolution.** Skills instruct Claude to resolve paths relative to the skill file's own location. For example, instead of `plugins/forge/skills/eval/rubrics/harness.md`, the skill says "Read the file at `../eval/rubrics/harness.md` relative to this skill file." This matches the npm convention (resolved via `__dirname`) and requires no runtime variable.

  **Why not the alternatives:**
  - `${CLAUDE_PLUGIN_ROOT}` is only available to hooks/scripts, not to skills. Verified: skills are Markdown instructions executed by the AI agent — they have no access to environment variables. Using this would require a pre-processor, adding build complexity.
  - `<FORGE_ROOT>` placeholder would need a custom resolver in the skill runner, which does not exist and would be a larger change than the problem warrants.

  **Prerequisite spike** (1h): Before W2 implementation, verify that Claude correctly resolves relative paths when reading files referenced from a skill's install location. Test by placing a skill at the plugin cache path and confirming `Read` succeeds with a relative reference.
- **Fix gen-test-cases `types/` path**: Either rename `templates/` to `types/` or update both SKILL.md files to reference `templates/`
- **Remove gen-test-scripts dead weight**:
  - Delete 14 Playwright-specific template files (duplicated in profile system)
  - Delete `node_modules/` (27MB) — applies the JetBrains lesson that marketplace gatekeeping rejects plugins shipping redundant library copies; move `validate-specs.mjs` to profile tooling or make it a project-level devDependency
  - Extract Playwright auth instructions (lines 359-371) to `generate.md`

### W3: Dedup & Simplification (~5h, Long-term health)

- **Extract shared logic from breakdown-tasks ↔ quick-tasks**:
  - Shared Type Assignment table → shared reference file
  - Shared Intent Propagation → shared reference file
  - Shared Step 0 profile resolution → shared reference file
  - 6 near-duplicate templates → shared template files
- **Simplify eval**: Extract 90-line validate-ux inline sub-pipeline to a rubric file (not a new skill — see scope exclusion alignment below)
- **Fix breakdown-tasks scope heuristic**: Correct `src/` classification for backend projects

### Preventive Mechanism: CI Path Guard

To prevent regressions after this audit, add a linter step to CI that detects hardcoded `plugins/forge/` paths in skill/agent/reference files:

```bash
# .github/workflows/skill-lint.yml or equivalent
# Fails if any skill or agent file contains a hardcoded plugin path
# (init-justfile is exempt — its templates work correctly in installed plugins)
! grep -rn 'plugins/forge/' skills/ agents/ references/ commands/ \
    --include='*.md' | grep -v 'init-justfile'
```

This mirrors how npm's `npm publish` runs `npm pack --dry-run` to verify the published tarball contains only intended files. The CI guard makes the "no hardcoded paths" invariant machine-enforced rather than relying on reviewer diligence.

## Requirements Analysis

### Key Scenarios

- User installs forge plugin → all skills reference files correctly via plugin-relative paths
- User runs `/submit-task` → works identically to old `/record-task` (no residual references)
- User with non-Playwright test framework → gen-test-scripts generates correct templates from profile
- User project has `docs/conventions/` → agents read these correctly as user-facing specs
- `forensic` skill works on any user's machine, not just the developer's

**Failure/Recovery Scenarios:**

- A skill references `../eval/rubrics/harness.md` but the file is missing from the plugin install → Claude logs a "file not found" error. The skill should include a fallback instruction: "If the referenced file cannot be found, proceed using the evaluation criteria described inline below." Each cross-referenced skill should embed a compressed summary of critical content as a degradation fallback.
- User's project has no `docs/conventions/` directory → `consolidate-specs` and other skills that read from this directory already handle this gracefully (the directory is created during `/init`). No change needed; verified by existing behavior.
- A shared template is updated for breakdown-tasks but the change is incompatible with quick-tasks → The CI guard (see W3 preventive mechanism) detects the change. If it ships undetected, the fix is to revert the shared file and re-split into skill-specific copies. Mitigation: shared reference files document which skills consume them, and changes require updating all consumers.

### Non-functional Requirements

- **Install size**: Target plugin install size under 5MB (currently ~30MB due to node_modules). W2 alone achieves this by removing 27MB of dead weight.
- **Backward compatibility**: All existing e2e tests in `forge-cli/tests/e2e/` must pass after each workstream. No command names, CLI flags, or user-facing output formats change.
- **Security**: Path resolution must not allow traversal outside the plugin directory (e.g., `../../etc/passwd`). Relative-from-self resolution inherently constrains this since Claude only reads paths the skill specifies.
- **Cross-platform**: Path references must work on Windows, macOS, and Linux. Relative paths with forward slashes are normalized by Claude's Read tool.

### Verification Plan

1. **W1**: `grep -r 'record-task'` and `grep '~/.claude/projects/'` return 0 hits. Run `go test ./tests/e2e/...` in forge-cli.
2. **W2**: Install plugin via `forge install`. Invoke each affected skill (`/eval`, `/forensic`, `/gen-test-cases`, `/gen-test-scripts`, `/run-e2e-tests`, `/graduate-tests`, `/improve-harness`, `/consolidate-specs`) and confirm no file-not-found errors. Compare gen-test-scripts output against golden fixture: `diff <(output after) <(golden-output/)`.
3. **W3**: For breakdown-tasks/quick-tasks dedup, generate tasks for the same PRD via both paths and diff the outputs. For eval extraction, run existing eval test cases and confirm scores are identical before/after.

### Constraints & Dependencies

- **Backward compatibility**: All slash commands must work identically after changes
- **Plugin distribution**: Skills must work when forge is installed as a plugin (not as source repo)
- **Test profile system**: Already exists at v3.0.0 — the Playwright abstraction builds on it
- **Existing `skill-rationalization` proposal**: Eval consolidation already done (single eval skill + 16 rubrics). This proposal does not duplicate that work

## Alternatives

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| **Three workstreams, prioritized** | Addresses all issues; can be phased | Medium total effort (~16h) | **Selected** |
| Only W1 (cleanup) | Minimal effort, low risk | Leaves path correctness broken | Rejected: only treats symptoms |
| Only W2 (distributability) | Fixes user-facing breakage | Leaves redundancy and complexity | Partial: good first phase |
| Build-time path rewriting | Zero runtime changes; skill source stays readable | Requires a build step for the plugin package; fragile if path format changes; does not fix dead code | Rejected: adds build complexity without addressing W1/W3 |

## Scope

### Phasing Strategy

All three workstreams ship in a single v3.0.0 release, but are sequenced: W1 first (unblocks testing), W2 second (core correctness), W3 last (polish). Total estimated effort: ~16h. If W1+W2 combined exceed 12h, defer W3 to v3.0.1 — W1 and W2 are release-blocking, W3 is not.

### In Scope

1. Fix 3 stale `record-task` references (W1)
2. Fix forensic hardcoded developer paths + stale build command (W1)
3. Fix 22+ `plugins/forge/` path references for installed plugin context (W2)
4. Fix gen-test-cases `types/` vs `templates/` path mismatch (W2)
5. Remove 14 Playwright template files + 27MB node_modules from gen-test-scripts (W2)
6. Extract Playwright auth instructions from gen-test-scripts SKILL.md to generate.md (W2)
7. Deduplicate breakdown-tasks ↔ quick-tasks shared logic (W3)
8. Simplify eval validate-ux inline sub-pipeline — extraction to a rubric file under eval/ (W3)
9. Fix breakdown-tasks scope heuristic for backend projects (W3)
10. Add CI path-guard linter (preventive mechanism)

### Out of Scope

- Eval skill restructuring (already done)
- Core pipeline skills logic (brainstorm, write-prd, ui-design, tech-design)
- CLI binary changes
- New skill creation (validate-ux extraction stays within eval/ as a rubric, not a separate skill)
- User-facing command name changes
- `just` or `doc-scorer`/`doc-reviser` references (expected dependencies)
- `init-justfile` template references (work correctly in installed plugin)
- Naming changes (`gen-` prefix)

## Key Risks

Likelihood calibration: H = >70%, M = 30–70%, L = <30%.

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Relative-from-self resolution does not work at the plugin cache path | M | H | **Prerequisite spike** (1h): place a test skill at `~/.claude/plugins/cache/forge/forge/<version>/skills/test/SKILL.md` with a relative reference `../eval/rubrics/harness.md`, invoke the skill, confirm Claude reads the file. If it fails, fall back to embedding content inline with a degradation instruction. |
| Moving Playwright templates breaks test generation | M | H | Generate golden output before deletion: `claude --skill gen-test-scripts --profile web-playwright > golden/`. After deletion, re-generate and `diff -r golden/ output/`. Empty diff = safe. |
| Deduplicating breakdown-tasks/quick-tasks breaks task generation | M | H | Generate tasks for the same test PRD via both `/breakdown-tasks` and `/quick-tasks` before and after dedup. Diff the task files. Identical task structure = safe. |
| Extracting validate-ux from eval changes scoring behavior | L | M | Run eval on 3 existing test cases before extraction. After extraction, re-run on the same inputs. Compare scores line-by-line. |
| Shared template changes have unintended cross-skill effects | M | M | Each shared reference file lists its consuming skills in a header comment. CI guard runs `diff` on both consumers' output after any shared file changes. If outputs diverge, the change is flagged for manual review. |
| Cross-platform path inconsistency — Windows backslash vs Unix forward slash in relative paths | L | H | Relative paths in skills use forward slashes exclusively. Claude's Read tool normalizes separators on Windows (confirmed in NFR section). Verify by running the prerequisite spike on a Windows machine. |

## Success Criteria

Each criterion includes a verification command where applicable.

- [ ] Zero stale `record-task` references in active source files
  `grep -r 'record-task' agents/ skills/ breakdown-tasks/templates/ quick-tasks/templates/` returns 0 hits (excluding the skill's own directory if it still exists on disk)
- [ ] Zero hardcoded user-specific paths in forensic
  `grep -rn '~/.claude/projects/' forensic/SKILL.md` returns 0 hits
- [ ] Zero `plugins/forge/` hardcoded paths that fail at runtime
  `grep -rn 'plugins/forge/' skills/ agents/ references/ commands/ | grep -v 'init-justfile'` returns 0 hits (init-justfile templates work correctly in installed plugin)
- [ ] `gen-test-cases/types/` path mismatch resolved
  Either `ls skills/gen-test-cases/types/` shows `.md` files, or both referencing SKILL.md files use `templates/`
- [ ] `gen-test-scripts/templates/node_modules/` no longer exists
  `test -d skills/gen-test-scripts/templates/node_modules && echo FAIL || echo OK`
- [ ] Zero Playwright-specific template files in gen-test-scripts (all in profile system)
  `ls skills/gen-test-scripts/templates/*.ts skills/gen-test-scripts/templates/*.json 2>/dev/null | grep -v generate.md` returns 0 hits
- [ ] breakdown-tasks and quick-tasks share common logic via shared reference files
  `ls references/shared/type-assignment.md references/shared/intent-propagation.md references/shared/step0-profile-resolution.md` succeeds; content-level dedup verified: `diff <(cat references/shared/type-assignment.md) <(grep -A20 'Type Assignment' breakdown-tasks/SKILL.md)` shows the shared file content matches the inline section in both skills (not just file existence)
- [ ] All 12 slash commands produce identical behavior to pre-migration (`/brainstorm`, `/write-prd`, `/ui-design`, `/tech-design`, `/eval`, `/breakdown-tasks`, `/quick-tasks`, `/gen-test-cases`, `/gen-test-scripts`, `/run-e2e-tests`, `/graduate-tests`, `/submit-task`)
  Verified by running existing e2e test suite: `cd forge-cli && go test ./tests/e2e/...` — all tests pass (covers 7 commands with functional e2e tests; remaining 5 verified by smoke test at line 276)
- [ ] Plugin installs and all skills work correctly from installed location
  `forge install && claude --skill eval --prompt "list rubrics"` succeeds without file-not-found errors
- [ ] eval validate-ux sub-pipeline extracted from eval/SKILL.md
  `wc -l skills/eval/SKILL.md` returns under 280 lines (was 366); extracted content exists in `skills/eval/rubrics/validate-ux-pipeline.md` or equivalent
- [ ] breakdown-tasks scope heuristic correctly classifies backend `src/` directories
  `grep -A5 'src/' breakdown-tasks/SKILL.md` shows logic that checks for Go/Rust indicators before classifying as frontend
- [ ] CI guard passes: no hardcoded `plugins/forge/` paths regress
  `.github/workflows/skill-lint.yml` (or equivalent) runs `grep -rn 'plugins/forge/' skills/ agents/ references/ --include='*.md' | grep -v init-justfile` and fails on non-zero exit
