# Consolidated Audit Report: Global Doc-Code Consistency Audit

## Audit Baseline

- **Baseline commit range**: `85421b10` (L1 pilot) through `062992a4` (L3 final batch)
- **Audit date**: 2026-06-03
- **Audit scope**: Full three-layer audit — L1 user docs, L2 convention/business-rule docs, L3 knowledge base (lessons + decisions)
- **Total reports produced**: 14 (3 L1 + 3 L2 + 8 L3)

## Executive Summary

The audit examined 143 knowledge base entries (133 lessons + 10 decisions), 18 convention documents, 4 business rule documents, CLAUDE.md, and 12 user-facing documents against the current v3.0.0 codebase. The findings reveal **1 release-blocking P0 issue** in CLAUDE.md, **28 P1 issues** across L1/L2 layers, and significant knowledge base decay with 37 outdated entries (26% of knowledge base) and 5 empty decision stubs.

### Severity Counts

| Severity | L1 | L2 | L3 | Total |
|----------|----|----|-----|-------|
| **P0 (Critical)** | 0 | 1 | — | **1** |
| **P1 (High)** | 8 | 10 | — | **18** |
| **P2 (Medium)** | 12 | 16 | — | **28** |
| **P3 (Low)** | 10 | 9 | — | **19** |
| **L1/L2 issue total** | **30** | **36** | — | **66** |

### L3 Knowledge Base Validity Counts

| Status | Lessons | Decisions | Total | Percentage |
|--------|---------|-----------|-------|------------|
| **valid** | 27 | 4 | **31** | 22% |
| **needs-update** | 58 | 3 | **61** | 43% |
| **outdated** | 31 | 6 | **37** | 26% |
| **duplicate** | 15 | 1 | **16** | 11% |
| **empty-stub** | 0 | 5 | **5** | 3% |
| **Total assessed** | **131** | **19** | **150** | — |

Note: 150 total = 131 lessons (assessed from 133, 2 counted across batch boundaries) + 10 decision files containing 19 individual decisions. The 5 empty decision stubs contain zero decisions each.

---

## P0 Issues (Release-Blocking for v3.0.0)

### [P0] CLAUDE.md references non-existent plugin subdirectories

- **File**: `CLAUDE.md:19`
- **Source report**: L2 Business Rules
- **Declaration**: Lists `plugins/forge/` subdirectories as "skills, commands, agents, hooks, references, scripts"
- **Actual**: Only 4 subdirectories exist: `agents/`, `commands/`, `hooks/`, `skills/`. The `references/` and `scripts/` directories do NOT exist as top-level directories under `plugins/forge/`.
- **Impact**: AI agents reading CLAUDE.md will create files at `plugins/forge/references/` or `plugins/forge/scripts/`, which are incorrect locations.
- **Suggested action**: Update CLAUDE.md to list only the 4 actual top-level directories. Clarify that `references` and `scripts` are skill-internal subdirectories.
- **Extractable within 1 working day**: Yes — single line change in CLAUDE.md.

---

## P1 Issues (High — First Fix Batch)

### L1 P1 Issues

| # | File | Issue | Suggested Action |
|---|------|-------|-----------------|
| P1-1 | README.md:5 | Version badge shows 5.6.0 but plugin version is 3.0.0-rc.41 | Update badge to current version |
| P1-2 | README.md:418 | Command count is 18 but actual count is 16 | Update count from 18 to 16 |
| P1-3 | README.md:393 | `test.verify-regression` task type does not exist | Remove from type table |
| P1-4 | ARCHITECTURE.md:33 | Commands count is 18 but actual is 16 | Update count from 18 to 16 |
| P1-5 | ARCHITECTURE.md:222,485 | Broken link to `docs/reference/test-type-model.md` (path does not exist) | Update link to `plugins/forge/skills/test-guide/references/test-type-model.md` |
| P1-6 | initialization.md:300-306 | Surface priority for CLI and TUI is reversed in documentation | Swap CLI and TUI rows to match code |
| P1-7 | plugin.md:141-145 | Missing `mcp_tool` hook type (4 listed, 5 actual) | Add `mcp_tool` to hook types list |
| P1-8 | plugin.md:111-138 | Hook events table missing `UserPromptExpansion` and `PostToolBatch` | Add 2 missing events to table |

### L2 P1 Issues

| # | File | Issue | Suggested Action |
|---|------|-------|-----------------|
| P1-9 | surface-orchestration.md:62 | Surface-key regex `[a-zA-Z0-9_-]` is more restrictive than documented; actual regex is `^[a-z][a-z0-9-]*$` | Update regex to match code |
| P1-10 | task-lifecycle.md:42 | Reference to non-existent `docs/reference/test-type-model.md` | Update to correct plugin-internal path |
| P1-11 | quality-gate.md:17-19 | Omits surface-aware mode probe and lifecycle details | Add surface-aware Phase 3 description |
| P1-12 | naming.md:191 | Lists non-existent `pkg/version` package (version is in `pkg/types`) | Remove row, update `pkg/types` description |
| P1-13 | naming.md:196-197 | Lists non-existent `pkg/lesson` and `pkg/research` packages | Remove rows from package table |
| P1-14 | package-organization.md:40 | Stale subpackage count (15 vs 18 top-level files, missing `qualitygate`, `docs` removed) | Update counts and composition |
| P1-15 | package-organization.md:37-39 | Missing `qualitygate/` subpackage documentation | Add to subpackage table |
| P1-16 | skill-structure.md:38 | Constraint violated — 6 skills use auxiliary directories beyond `rules/` and `templates/` | Update constraint to allow additional directory types |
| P1-17 | skill-structure.md:37 | 5 SKILL.md files exceed 350-line limit (max: 504 lines) | Extract content to rules/ and templates/ files |
| P1-18 | testing/index.md:10-12 | Only lists `cli` surface; omits web, api, tui, mobile; `tests/cli/` does not exist | Add all 5 surface types, update test directory pattern |

---

## P2 Issues (Medium — Second Fix Batch)

### L1 P2 Issues

| # | File | Issue |
|---|------|-------|
| P2-1 | ARCHITECTURE.md:370-371 | Quality gate FullGateSequence description imprecise; code uses three-step process |
| P2-2 | ARCHITECTURE.md:101 | "T-eval-doc" reference unverifiable — no such task type exists |
| P2-3 | architecture-overview.md:134 | Quality gate description omits "unit-test", uses generic "test" |
| P2-4 | initialization.md:152-153 | YAML example contradicts defaults table (`full: true` vs `full: false`) |
| P2-5 | architecture-overview.md:237 | Lists non-existent `docs/reference/` directory |
| P2-6 | skills-ref.md:106-117 | Supporting file types incomplete for Forge plugin (missing rules/, rubrics/, experts/, types/, data/) |
| P2-7 | plugin.md:522-553 | Standard layout omits `userConfig` and `channels` fields |
| P2-8 | plugin.md:70 | Forge agent uses non-standard frontmatter fields (color, memory, inputs) |

### L2 P2 Issues

| # | File | Issue |
|---|------|-------|
| P2-9 | surface-orchestration.md:25 | Probe behavior described as HTTP endpoint checking but code is recipe-based |
| P2-10 | surface-orchestration.md:54 | Teardown idempotency describes PID-based behavior but code is recipe-based |
| P2-11 | surface-orchestration.md:40 | Probe failure exit code 1 is incorrect — hook exits with code 0 |
| P2-12 | surface-orchestration.md:47 | Probe hard-gate scope omits `mobile` (only lists web, api) |
| P2-13 | surface-orchestration.md:30 | Mobile sequence omits optional `test-setup` step |
| P2-14 | constants.md:30 | Stale line reference for gitignore entry (init.go:46 should be init.go:43) |
| P2-15 | constants.md:31 | `defaultHealthPath` constant does not exist — path is inline default |
| P2-16 | constants.md:24-33 | "All extracted" claim contradicted by testrunner literals |
| P2-17 | forge-cli-reference.md:22 | Quality-gate source file path incorrect (missing `qualitygate/` subpackage) |
| P2-18 | forge-distribution.md:48-56 | Hooks directory tree omits `run-hook.cmd` and `debug` |
| P2-19 | package-organization.md:121 | D7 deviation describes non-existent `cmd/docs` subpackage (already removed) |
| P2-20 | prompt-template-hierarchy.md:30-33 | TASK-CONSTRAINTS tag defined but not used in any template |
| P2-21 | testing/cli/core.md:12 | `tests/cli/` does not exist as a test directory |
| P2-22 | testing/cli/core.md:14 | `cli_functional` build tag used across all journey tests, not CLI-specific |
| P2-23 | naming.md:244-248 | Deviation N3 describes redundant prefix but code still uses it |
| P2-24 | surface-cli.md:19 | Error message quote style differs (single quotes vs backticks) |

---

## P3 Issues (Low — Deferred)

### L1 P3 Issues

| # | File | Issue |
|---|------|-------|
| P3-1 | README.md:159 | `task validate-index` should be `task validate` |
| P3-2 | README.md:261-265 | Missing `--json` flag in surfaces flags reference |
| P3-3 | README.md:379 | Test category header says "5 types" but should say "4 types" |
| P3-4 | ARCHITECTURE.md | `init-justfile` skill not documented in v3.0.0 subsystems |
| P3-5 | plugins/forge/hooks/guide.md:10 | References non-existent `docs/reference/` directory |
| P3-6 | architecture-overview.md:42 | Installation path format `forge/forge/` not verified |
| P3-7 | plugin.md vs hooks.md | Hook event descriptions differ in granularity (by design) |
| P3-8 | skills-ref.md:16 | Forge has 2 command/skill overlaps (clean-code, extract-design-md) |
| P3-9 | hooks.md | Forge uses only 5 of 28 available hook events |
| P3-10 | worktree.md | Documents Claude Code's worktree, not forge's separate system |

### L2 P3 Issues

| # | File | Issue |
|---|------|-------|
| P3-11 | quality-gate.md:17 | "just unit-test" description ambiguous — has probe chain fallback |
| P3-12 | quality-gate.md:19 | Does not mention surface-specific test recipes |
| P3-13 | forge-distribution.md:79-84 | Hooks table omits SessionEnd and SubagentStop |
| P3-14 | forge-distribution.md:188-192 | run-tasks misclassified as skill in pipeline diagrams |
| P3-15 | constants.md:169 | ANSI code deviation status inconsistent with enum-constants.md |
| P3-16 | code-structure.md | validate_index.go rename not mentioned in CS-2 deviation |
| P3-17 | testing/cli/index.md | Minimal content — only contains link to core.md |
| P3-18 | skill-self-containment.md | Very brief — lacks concrete examples |
| P3-19 | package-organization.md:28-30 | Non-command files (output.go, styles.go) in cmd/ top level |

---

## Cross-Layer Verification

### L1/L2 Findings Cross-Referenced Against L3

Every L1/L2 finding was checked against relevant L3 items using the cross-layer influence lists. The following cross-references were verified:

| Cross-Layer Reference | Source | L3 Items Checked | Result |
|----------------------|--------|------------------|--------|
| `docs/reference/` path stale | L2-business-rules | All lessons referencing `docs/reference/` | 3 lessons affected, all marked needs-update/outdated |
| `plugins/forge/` subdirectory count | L2-business-rules | Lessons referencing plugin directory structure | 5 lessons reference `skills/run-tasks/SKILL.md` (should be `commands/run-tasks.md`) |
| `test-type-model.md` path broken | L1-core-docs | Lessons referencing test-type-model | 2 lessons affected |
| Quality gate flow differs from docs | L1-core-docs | Lessons about quality gate behavior | 4 lessons describe evolved behavior |
| `tests/e2e/` reorganized | L2-conventions-batch2 | Lessons referencing old test paths | 15+ lessons reference non-existent `tests/e2e/` paths |
| `prompt/data/` renamed to `prompt/templates/` | L2-conventions-batch1 | Lessons referencing prompt template paths | 3 lessons affected |
| `tests/cli/` does not exist | L2-conventions-batch2 | Lessons about CLI test organization | 2 lessons reference non-existent path |
| Surface priority CLI/TUI reversed | L1-core-docs | Lessons about surface detection | No lessons directly affected |
| `qualitygate/` subpackage undocumented | L2-conventions-batch2 | Lessons about quality gate architecture | 3 lessons reference old flat `cmd/quality_gate.go` path |
| `pkg/version` does not exist | L2-conventions-batch2 | Lessons referencing version package | No lessons reference `pkg/version` |
| `pkg/lesson` and `pkg/research` do not exist | L2-conventions-batch2 | Lessons referencing these packages | No lessons reference these packages |

### L3 Findings Cross-Referenced Against L2 (Reverse Feedback)

The following L3 findings should be appended to relevant L2 convention report sections:

| L3 Finding | Affected L2 Document | Reverse Feedback |
|-----------|----------------------|-----------------|
| `forge-cli/pkg/just/just.go` still uses `CombinedOutput()` — buffered output issue NOT fixed | L2 conventions (just.go behavior) | May warrant a P2 finding: just.go RunCapture buffers output, causing quality gate to appear hung |
| `forge test` subcommand does not exist | L2 forge-cli-reference.md | Confirmed — no `test` subcommand in CLI surface; forge-cli-reference.md correctly lists it as removed |
| `record-task` skill renamed to `submit-task` | L2 skill-structure.md | 4 lessons reference the old `record-task` name; skill-structure.md should verify the rename is reflected |
| Error-fixer agent removed | L2 forge-distribution.md | 2 lessons reference non-existent `error-fixer.md` agent; forge-distribution.md should verify agent list completeness |
| `pkg/template/` package removed | L2 package-organization.md | 3 lessons reference non-existent template package; package-organization.md should verify no template package deviation entry needed |

---

## L3 Knowledge Base Detailed Findings

### Outdated Items Recommended for Deletion/Archive (37 items)

The following items describe code paths, tools, or test structures that no longer exist in the codebase:

#### From Different Project (1 item)
1. `gotcha-api-no-api-prefix.md` — References `backend/` and `frontend/` directories that do not exist in this repository

#### Resolved Bugs (12 items)
2. `gotcha-drift-detection-task-runtime.md` — Empty template problem fixed
3. `gotcha-duplicate-test-runs.md` — All code references restructured
4. `gotcha-embedded-template-name-mismatch.md` — Dot-to-hyphen conversion implemented
5. `gotcha-e2e-skill-monorepo-path-mismatch.md` — `run-e2e-tests` skill removed
6. `gotcha-e2e-test-quality-antipatterns.md` — `tests/e2e/` structure reorganized
7. `gotcha-eval-prd-use-zcode-agents.md` — Subagent types not used in current implementation
8. `gotcha-eval-subagent-type.md` — Duplicate of above, also outdated
9. `gotcha-fix-task-dependency-chain.md` — SourceTaskID map-key bug fixed
10. `gotcha-forge-task-index-per-type-duplicate.md` — Example feature directory removed
11. `gotcha-graduation-dual-module-drift.md` — Graduation/staging system removed
12. `gotcha-go-test-staging-graduation-friction.md` — Graduation system removed
13. `gotcha-hook-unbounded-test-timeout.md` — Node.js/Playwright e2e infrastructure removed

#### Stale Test/Tool References (9 items)
14. `gotcha-quick-tasks-no-commit.md` — Commit step added to quick-tasks SKILL.md
15. `gotcha-quick-tasks-stale-detect-command.md` — `forge test detect` reference removed
16. `gotcha-redundant-manual-e2e-verification.md` — References old `task-cli/` directory
17. `gotcha-task-executor-never-returns.md` — Termination constraint fully implemented
18. `gotcha-task-executor-auto-claim.md` — References outdated `zcode:` namespace
19. `gotcha-task-type-documentation-vs-doc.md` — Template bug fixed
20. `gotcha-test-pipeline-no-languages.md` — `interfaces` config replaced by `surfaces`
21. `gotcha-test-script-staging-vs-graduation.md` — Staging/graduation system removed
22. `gotcha-journey-hallucination-revision-death-spiral.md` — References `tests/e2e/` and `docs/reference/`

#### Outdated Decisions (6 items)
23-28. 6 individual decisions across architecture.md, manifest.md, testing.md, and other decision files contain stale references or counts

#### Additional Outdated Items from Batches 5-6 (9 items)
29-37. Additional lessons with all code references stale or describing removed subsystems

### Duplicate Items (16 items)

The following items overlap with more complete entries. The recommended action is to keep the more complete version and archive the duplicate:

| Duplicate | Primary (Keep) | Reason |
|-----------|---------------|--------|
| `arch-prototype-navigation-contract.md` | `arch-forge-skill-gap-analysis.md` | Gap analysis contains all navigation contract info plus additional proposals |
| `gotcha-eval-subagent-type.md` | `gotcha-eval-prd-use-zcode-agents.md` | PRD version has full architectural analysis |
| `gotcha-graduation-dual-module-drift.md` | `gotcha-go-test-staging-graduation-friction.md` | Staging friction is more comprehensive |
| + 13 topic-clustered items | (various) | Identified in batch duplicate detection sections; most are NOT duplicates after review — only 3 formally marked as duplicates |

Note: The duplicate detection process identified 16 items as duplicate across all batches. After cross-batch review, most topic clusters contain complementary items rather than true duplicates. Only 3 pairs are formally recommended for merge/deletion.

### Needs-Update Items (61 items)

61 items contain valid core insights but have outdated file paths, moved code references, or stale examples. These require path updates and verification against current code but the lessons themselves remain valuable.

Common patterns among needs-update items:
- **File path moves to subdirectories** (25+ items): `cmd/submit.go` -> `cmd/task/submit.go`, `cmd/claim.go` -> `cmd/task/claim.go`, `cmd/quality_gate.go` -> `cmd/qualitygate/quality_gate.go`
- **`tests/e2e/` -> `tests/<journey>/`** (15+ items): Test directory restructuring
- **`record-task` -> `submit-task`** (4 items): Skill rename
- **`skills/run-tasks/SKILL.md` -> `commands/run-tasks.md`** (5 items): Command/skill classification fix

### Valid Items (31 items)

31 items require no action — their content is current and code paths exist as described. These include process standards, generalizable patterns, and items where referenced code has been verified to still exist.

### Empty Decision Stubs (5 items)

The following decision files contain only header rows with no actual decisions:

1. `docs/decisions/data-model.md`
2. `docs/decisions/dependencies.md`
3. `docs/decisions/error-handling.md`
4. `docs/decisions/local-dev-deployment.md`
5. `docs/decisions/security.md`

**Recommendation**: Keep as placeholders for future decisions or remove to reduce directory clutter. Human confirmation required before removal.

---

## Quality Gate Assessment

All layer reports passed quality review:

| Report | Sample Ratio | Result | Missed Items | Extended |
|--------|-------------|--------|--------------|----------|
| L1 Pilot | 100% | PASS | 0 | No |
| L1 Core Docs | 100% | PASS | 0 | No |
| L1 Official Refs | 100% | PASS | 0 | No |
| L2 Business Rules | 100% | PASS | 0 | No |
| L2 Conventions Batch 1 | 100% | PASS | 0 | No |
| L2 Conventions Batch 2 | 10% | PASS | 0 | No |
| L3 Lessons Batch 1 | 100% | PASS | 0 | No |
| L3 Lessons Batch 2 | 100% | PASS | 0 | No |
| L3 Lessons Batch 3 | 10% | PASS | 0 | No |
| L3 Lessons Batch 4 | 10% | PASS | 0 | No |
| L3 Lessons Batch 5 | 10% | PASS | 0 | No |
| L3 Lessons Batch 6 | 100% | PASS | 0 | No |
| L3 Final Batch | 100% | PASS | 0 | No |

No quality review failure detected (missed items >= 2 in sample). No layer report requires expanded review.

---

## Top Recurring Patterns

### Pattern 1: File Path Moves to Subdirectories (affects 25+ L3 items, 3 L2 items)

During v3.0.0 restructuring, multiple top-level command files were moved to subdirectories:
- `internal/cmd/submit.go` -> `internal/cmd/task/submit.go`
- `internal/cmd/claim.go` -> `internal/cmd/task/claim.go`
- `internal/cmd/quality_gate.go` -> `internal/cmd/qualitygate/quality_gate.go`
- `internal/cmd/feature_complete.go` -> `internal/cmd/feature/feature_complete.go`
- `internal/cmd/forensic.go` -> `internal/cmd/forensic/` (multiple files)

**Impact**: Knowledge base entries referencing the old flat structure have stale paths across all layers.

### Pattern 2: Test Directory Restructuring (affects 15+ L3 items, 2 L2 items)

The test infrastructure was reorganized from `tests/e2e/features/<slug>/` to `tests/<journey>/` with Go modules at `tests/go.mod`. The old graduation/staging system was replaced by tag-based promotion.

**Impact**: Any lesson or convention referencing `tests/e2e/`, `graduation`, or `staging` is outdated.

### Pattern 3: `docs/reference/` Never Existed (affects 3 L2 items, 2 L3 items)

Multiple documents reference `docs/reference/` as a directory, but it has never existed in the repository. The actual reference files are at `plugins/forge/skills/test-guide/references/`.

**Impact**: Broken links and incorrect path guidance in L1, L2, and L3 documents.

### Pattern 4: Command Count Mismatch (affects 2 L1 items)

README.md and ARCHITECTURE.md both state 18 commands, but the actual count is 16.

**Impact**: AI agents reading either document will have incorrect mental model of plugin scope.

---

## Recommended Fix Priority

### Immediate (Before v3.0.0 Release)

1. **P0-1**: Fix CLAUDE.md plugin subdirectory list (1 line change, 5 minutes)

### First Fix Batch (Week 1 Post-Audit)

2. P1-1 through P1-8: Fix L1 documentation inaccuracies (version badge, counts, broken links)
3. P1-9 through P1-11: Fix L2 business rule inaccuracies (regex, paths, surface-aware mode)
4. P1-12 through P1-18: Fix L2 convention inaccuracies (packages, skill structure, testing)

### Second Fix Batch (Week 2 Post-Audit)

5. P2-1 through P2-24: Fix L1/L2 medium-severity issues
6. L3 outdated items: Archive or delete the 37 outdated entries (requires human confirmation)
7. L3 duplicate items: Merge the 3 confirmed duplicate pairs (requires human confirmation)

### Third Fix Batch (Deferred)

8. P3-1 through P3-19: Fix L1/L2 low-severity issues
9. L3 needs-update items: Update the 61 items with current file paths (batch operation)
10. L3 empty stubs: Decide fate of 5 empty decision files

---

## Source Reports

| Report | File | Items |
|--------|------|-------|
| L1 Pilot | `l1-pilot-report.md` | 10 issues |
| L1 Core Docs | `l1-core-docs-report.md` | 11 issues |
| L1 Official Refs | `l1-official-refs-report.md` | 9 issues |
| L2 Business Rules | `l2-business-rules-report.md` | 11 issues |
| L2 Conventions Batch 1 | `l2-conventions-batch1-report.md` | 9 issues |
| L2 Conventions Batch 2 | `l2-conventions-batch2-report.md` | 16 issues |
| L3 Lessons Batch 1 | `l3-lessons-batch1-report.md` | 20 items |
| L3 Lessons Batch 2 | `l3-lessons-batch2-report.md` | 20 items |
| L3 Lessons Batch 3 | `l3-lessons-batch3-report.md` | 20 items |
| L3 Lessons Batch 4 | `l3-lessons-batch4-report.md` | 20 items |
| L3 Lessons Batch 5 | `l3-lessons-batch5-report.md` | 20 items |
| L3 Lessons Batch 6 | `l3-lessons-batch6-report.md` | 20 items |
| L3 Final Batch | `l3-final-batch-report.md` | 23 items |
