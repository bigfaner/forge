# L3 Lessons Audit Report — Batch 2 (gotcha-breaking-change-quality-gate-deadlock to gotcha-fix-task-broad-scope)

## Audit Baseline

- **Baseline commit**: 2eb35df2 (docs(audit): add L3 lessons batch 1 audit report)
- **Audit date**: 2026-06-03
- **Audit scope**: 20 lesson files from `docs/lessons/` (alphabetically: gotcha-breaking-change-quality-gate-deadlock.md through gotcha-fix-task-broad-scope.md)

## Classification Distribution

| Classification | Count |
|----------------|-------|
| code-reference | 11 |
| process-standard | 5 |
| experience-summary | 4 |

## Validity Summary

| Status | Count |
|--------|-------|
| valid | 7 |
| outdated | 7 |
| needs-update | 6 |

Note: Duplicate detection identified 1 formal duplicate (gotcha-eval-subagent-type.md as duplicate of gotcha-eval-prd-use-zcode-agents.md). The duplicate item is counted in its item-by-item status (outdated).

## Cross-Layer Influence Check

The following L1/L2 audit findings affect items in this batch:

| L1/L2 Finding | Affected Lesson | Impact |
|---------------|-----------------|--------|
| L2 business-rules: `docs/reference/` does not exist; test-type-model.md is at `plugins/forge/skills/test-guide/references/` | gotcha-docs-only-needs-code-audit | Lesson references `plugins/forge/references/shared/type-assignment.md` which DOES NOT EXIST |
| L2 conventions-batch1: `run-tasks` is a command (not a skill) | gotcha-duplicate-test-runs | Lesson references `plugins/forge/commands/execute-task.md` (correct) — no impact |
| L2 conventions-batch2: `tests/cli/` does not exist as a test directory | gotcha-e2e-test-quality-antipatterns | Related: `tests/e2e/` structure has been completely reorganized — all referenced test files are missing |
| L1 core-docs: quality gate flow differs from documented flow | gotcha-breaking-task-quality-gate-test-scope, gotcha-breaking-change-quality-gate-deadlock | Lessons describe quality gate behavior that may have evolved |
| L2 business-rules: CLAUDE.md plugin subdirectory list incomplete | gotcha-embedded-template-name-mismatch | Lesson references `forge-cli/pkg/template/` which DOES NOT EXIST; templates are now in `forge-cli/pkg/task/templates/` |

## Item Details

### 1. gotcha-breaking-change-quality-gate-deadlock.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes a reusable pattern for breaking change task scoping: include minimal caller updates in the breaking change task itself, scoped to "compiles + existing tests pass". The referenced files exist: `forge-cli/pkg/forgeconfig/config.go` EXISTS, `forge-cli/pkg/forgeconfig/detect.go` EXISTS, `forge-cli/pkg/task/build.go` EXISTS. The pattern (grep for references, classify into compilation-blocking vs non-blocking, include blocking references) is toolchain-independent and broadly applicable regardless of implementation changes. The specific code example (`ReadInterfaces` -> `ReadSurfaces`) is historical but the pattern remains sound.
- **Code path verification**: All 3 referenced files EXIST

### 2. gotcha-breaking-task-quality-gate-test-scope.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson references `plugins/forge/agents/task-executor/` which DOES NOT EXIST as a directory (the actual agent file is `plugins/forge/agents/task-executor.md`, a single file). It references `docs/features/pipeline-topology-registry/tasks/2-refactor-task-generation.md` which EXISTS. It references `forge-cli/pkg/task/autogen_test.go` which EXISTS. The core insight about the structural conflict between planning model (breaking tasks defer test fixes) and execution model (quality gate requires all tests pass) remains valid. The proposed solutions (two quality gate modes for breaking tasks, `exclude-tests` frontmatter field) describe improvements that may or may not have been implemented — the lesson should note the current state.
- **Code path verification**: `plugins/forge/agents/task-executor/` MISSING (correct: `plugins/forge/agents/task-executor.md`), task file EXISTS, autogen_test.go EXISTS
- **Required update**: (1) Fix agent path from directory to file; (2) Note whether the short-term/medium-term solutions have been implemented

### 3. gotcha-characterization-test-vs-refactoring.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes the tension between characterization tests and intentional refactoring, proposing the PRESERVE/EVOLVE Impact Declaration pattern. This is a generalized design pattern inspired by Michael Feathers' "Working Effectively with Legacy Code". The referenced task record `docs/features/forge-architecture-simplification/tasks/records/2.8-quality-gate-fixes.md` EXISTS. The core insight (characterization tests become obstacles during intentional behavior changes; breaking tasks must own the characterization test updates) is a universally applicable testing pattern independent of specific forge implementation.
- **Code path verification**: Task record EXISTS; pattern-level lesson with minimal code path dependencies

### 4. gotcha-dispatcher-ignores-compilation-diagnostics.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson references task-executor and dispatcher components. It describes a critical issue where the executor falsified quality gate results ("49 passed, 0 failed" while compilation errors existed). The core insight — "Verify, don't trust" for self-reported test results — remains universally valid. However, the specific code references need updating: (1) the quality gate implementation has been refactored into `forge-cli/internal/cmd/qualitygate/` (multiple files); (2) `submit.go` is now at `forge-cli/internal/cmd/task/submit.go` (moved to subdirectory); (3) the proposed two-layer defense (executor-level + dispatcher-level verification) should be assessed against current implementation.
- **Code path verification**: `forge-cli/internal/cmd/submit.go` MISSING (moved to `forge-cli/internal/cmd/task/submit.go`); qualitygate/ directory EXISTS
- **Required update**: (1) Update submit.go path to `forge-cli/internal/cmd/task/submit.go`; (2) Verify whether the proposed dispatcher-level diagnostic signal check has been implemented

### 5. gotcha-docs-only-needs-code-audit.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson references `forge-cli/pkg/task/build.go` (EXISTS) for the `testableTypes` map. However, the `testableTypes` map variable NO LONGER EXISTS — it has been replaced by the `IsTestableType()` function which uses a prefix check (`strings.HasPrefix(typ, "coding.") || typ == TypeCleanCode`). The lesson's solution section ("Added TypeCleanup and TypeRefactor to testableTypes map") describes a historical fix that has been superseded by the prefix-based approach. The lesson also references `plugins/forge/references/shared/type-assignment.md` which DOES NOT EXIST. The proposal reference `docs/proposals/task-type-code-docs-boundary/proposal.md` EXISTS. The core insight (audit every code path, not just the most obvious one) remains valid.
- **Code path verification**: `build.go` EXISTS, `testableTypes` map MISSING (replaced by `IsTestableType` function), `plugins/forge/references/shared/type-assignment.md` MISSING, `submit.go` moved to `cmd/task/` subdirectory
- **Required update**: (1) Note that `testableTypes` map has been replaced by prefix-based `IsTestableType()`; (2) Update `type-assignment.md` path or note it no longer exists; (3) Update `submit.go` path

### 6. gotcha-drift-detection-task-runtime.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes `doc.drift` tasks having an empty template ("Execute this test pipeline task") that causes long execution times. The current `doc.drift` template at `forge-cli/pkg/task/templates/doc-drift.md` has been significantly improved — it now includes a focused discovery strategy (git diff to identify changed files, scope-limited spec verification, explicit skip conditions). The root cause described in the lesson (empty template causing full codebase scan) has been resolved. The lesson's solution recommendations (skip for pure code refactors, inject context in templates, manual execution) were partially implemented via the improved template.
- **Code path verification**: `forge-cli/pkg/task/templates/doc-drift.md` EXISTS with substantive content (discovery strategy, git diff scoping, skip conditions)
- **Recommendation**: Mark as resolved/outdated — the empty template problem has been fixed. The general pattern about auto-generated tasks with empty templates remains valid, but the specific doc.drift case is no longer applicable.

### 7. gotcha-duplicate-test-runs.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes three layers of redundant test execution: (1) task type templates in `forge-cli/pkg/prompt/data/feature.md`, (2) submit-task SKILL.md metrics collection, (3) `forge task submit` CLI `validateQualityGate()`, (4) dispatcher breaking gate. Multiple structural changes render these references obsolete: (a) `forge-cli/pkg/prompt/data/` DOES NOT EXIST — prompt templates are now at `forge-cli/pkg/prompt/templates/` with different filenames (e.g., `coding-feature.md`); (b) `forge-cli/internal/cmd/submit.go` has moved to `forge-cli/internal/cmd/task/submit.go`; (c) `forge-cli/pkg/profile/framework.go` DOES NOT EXIST; (d) `plugins/forge/skills/run-e2e-tests/SKILL.md` DOES NOT EXIST (run-e2e-tests skill was removed). The core insight (designate one authoritative checkpoint) remains valid, but all specific code references are outdated.
- **Code path verification**: `prompt/data/` MISSING, `submit.go` moved, `profile/framework.go` MISSING, `run-e2e-tests/SKILL.md` MISSING
- **Recommendation**: Mark as outdated — all code references have been restructured. The general pattern is valid but the specific analysis no longer maps to the codebase.

### 8. gotcha-e2e-env-override-isolation.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: The lesson describes E2E tests inheriting project-root env vars (`CLAUDE_PROJECT_DIR`, `PROJECT_ROOT`) from the parent process. The referenced file `forge-cli/pkg/project/root.go` EXISTS and still contains `FindProjectRoot()`. The core insight and solution pattern (strip env vars in test helpers) remain valid. However, the referenced test file `tests/e2e/task_record_immutability_cli_test.go` DOES NOT EXIST — the `tests/e2e/` directory structure has been reorganized (tests now live in `forge-cli/tests/` with subdirectories like `forge-commands/`, `task-lifecycle/`, etc.). The cross-reference to `gotcha-fix-task-index-test-isolation.md` EXISTS.
- **Code path verification**: `forge-cli/pkg/project/root.go` EXISTS, `tests/e2e/task_record_immutability_cli_test.go` MISSING
- **Required update**: (1) Update test file path to current location under `forge-cli/tests/`; (2) Note that the test directory structure has been reorganized

### 9. gotcha-e2e-script-generation.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: This is an extensive structural analysis of the gen-test-scripts SKILL.md's blind spots (missing Code Reconnaissance step, missing Frontend UI coverage, fallback referencing nonexistent data, strategy/syntax mixing, UNKNOWN without completeness gate, hardcoded UI probe paths). The lesson documents 6 layers of root cause with corresponding fixes applied between 2026-04-28 and 2026-05-15. The `gen-test-scripts` skill EXISTS at `plugins/forge/skills/gen-test-scripts/SKILL.md`. The architecture principle ("SKILL.md owns strategy decisions, generate.md owns framework syntax") is a reusable design pattern. The timeline of fixes provides historical context. While the specific fixes were already applied, the lesson serves as a record of the blind spots found and the architectural principle established.
- **Code path verification**: `plugins/forge/skills/gen-test-scripts/SKILL.md` EXISTS

### 10. gotcha-e2e-skill-monorepo-path-mismatch.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes `run-e2e-tests` skill path mismatches in a monorepo. The referenced skill `plugins/forge/skills/run-e2e-tests/SKILL.md` DOES NOT EXIST — the run-e2e-tests skill has been removed from the codebase. The referenced test file `forge-cli/tests/e2e/features/spec-drift-detection/spec_drift_detection_cli_test.go` DOES NOT EXIST — the entire `tests/e2e/features/` directory structure has been removed. The task file `docs/features/spec-drift-detection/tasks/quick-run-tests-go-test.md` EXISTS. With the skill removed and the test structure reorganized, this lesson describes a configuration that no longer exists in the codebase.
- **Code path verification**: `run-e2e-tests/SKILL.md` MISSING, `e2e-report.md` template MISSING, spec-drift test file MISSING, task file EXISTS
- **Recommendation**: Mark as outdated — the skill and test structure referenced have been removed. The general pattern about path baseline ambiguity in monorepos remains valid, but all specific references are stale.

### 11. gotcha-e2e-test-binary-isolation.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: The lesson describes E2E tests depending on the system-installed forge binary instead of building from source. The referenced convention file `docs/conventions/testing-isolation.md` DOES NOT EXIST (but `docs/conventions/testing/` directory EXISTS with `cli`, `index.md`). The cross-reference to `gotcha-e2e-env-override-isolation.md` EXISTS. The core pattern (build binary from source in TestMain, use full path in exec.Command) remains valid. However, the test directory structure has changed significantly — tests are now in `forge-cli/tests/` rather than the old `tests/e2e/` structure.
- **Code path verification**: `docs/conventions/testing-isolation.md` MISSING (testing conventions now in `docs/conventions/testing/` directory)
- **Required update**: (1) Update convention path to `docs/conventions/testing/`; (2) Note test directory restructuring

### 12. gotcha-e2e-test-quality-antipatterns.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson audits 67 e2e test cases finding 52% with quality issues. ALL five referenced test files DO NOT EXIST: `tests/e2e/simplify_e2e_tests_cli_test.go`, `tests/e2e/cli_lean_output_cli_test.go`, `tests/e2e/cli_list_reverse_chronological_cli_test.go`, `tests/e2e/fix_task_claim_priority_cli_test.go`, `tests/e2e/feature_set_command_cli_test.go`. The entire `tests/e2e/` directory structure has been reorganized — tests now live under `forge-cli/tests/` with a flat subdirectory structure (forge-commands, task-lifecycle, etc.) rather than the old `tests/e2e/features/<slug>/` layout. The "graduation workflow" that caused duplicate tests (`/graduate-tests` copying from features/ to root) is based on the old directory structure that no longer exists. The cross-reference to `gotcha-recursive-go-test-process-explosion.md` EXISTS.
- **Code path verification**: ALL 5 referenced test files MISSING; `tests/e2e/` directory structure MISSING
- **Recommendation**: Mark as outdated — the entire test structure described has been reorganized. The quality antipatterns (recursive tests, conditional skips without fixtures, vacuous assertions) remain valid as general testing guidance, but all specific file references are stale.

### 13. gotcha-embedded-template-name-mismatch.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes a mismatch between embedded template filenames (`coding-fix.md` with hyphens) and the lookup path (`coding.fix` with dots). This problem has been RESOLVED: the current `autogenTemplatePath()` function in `forge-cli/pkg/task/autogen.go` automatically converts dots to hyphens (`strings.ReplaceAll(typeName, ".", "-")`), and the actual template files use hyphenated names (`coding.fix.md` in the `templates/` directory). The referenced file paths are also outdated: `forge-cli/pkg/template/template.go` DOES NOT EXIST (templates are now handled in `forge-cli/pkg/task/autogen.go` via `//go:embed templates/*.md`). `forge-cli/pkg/template/data/coding.fix.md` DOES NOT EXIST (correct path is `forge-cli/pkg/task/templates/coding.fix.md`). `forge-cli/pkg/task/add.go` EXISTS.
- **Code path verification**: `forge-cli/pkg/template/template.go` MISSING, `forge-cli/pkg/template/data/coding.fix.md` MISSING, `forge-cli/pkg/task/add.go` EXISTS, `forge-cli/pkg/task/autogen.go` EXISTS with dot-to-hyphen conversion
- **Recommendation**: Mark as outdated/resolved — the name mismatch has been fixed by automatic dot-to-hyphen conversion in `autogenTemplatePath()`.

### 14. gotcha-eval-loop-decision-gate.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson establishes the rule that the eval decision gate must fire after every scorer run, regardless of user instructions. The current eval SKILL.md (`plugins/forge/skills/eval/SKILL.md`) implements this correctly: the flowchart shows "score >= target" check immediately after scorer output, with separate paths for "Go to Step 5" (stop) vs "Go to Step 4" (revise). The step 3b decision table explicitly states "Score >= target → Go to Step 5" and "Score < target, ITERATION < MAX_ITERATIONS → Go to Step 4". The principle (decision gate is a hard rule, not affected by user "continue" instructions) is correctly implemented and the lesson documents this invariant.
- **Code path verification**: `plugins/forge/skills/eval/SKILL.md` EXISTS, decision gate logic confirmed correct

### 15. gotcha-eval-prd-use-zcode-agents.md

- **Classification**: process-standard
- **Status**: outdated
- **Justification**: The lesson recommends using `zcode:proposal-scorer` / `zcode:proposal-reviser` subagent types for eval-prd and eval-design. The current eval SKILL.md does NOT reference any `zcode:` subagent types. The current architecture uses `general-purpose` agents for scorer and reviser (spawned with `model: "sonnet"`), not specialized zcode subagents. The principle (main session owns the loop, scorer/reviser are separate single-responsibility agents) is correctly implemented in the current eval SKILL.md — but the specific subagent type names (`doc-scorer`, `doc-reviser`, `zcode:proposal-scorer`, `zcode:proposal-reviser`) referenced in this lesson do not match the current implementation.
- **Code path verification**: Current eval SKILL.md uses `general-purpose` agents with `model: "sonnet"`, not zcode: or doc- prefixed subagents
- **Recommendation**: Mark as outdated — the specific subagent type recommendations (`zcode:*`, `doc-scorer`, `doc-reviser`) do not match the current implementation. The general principle (main session as dispatcher, scorer/reviser as single-responsibility agents) remains valid and is correctly implemented.

### 16. gotcha-eval-reviser-too-many-attacks.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The lesson describes the performance issue when eval reviser processes too many attack points (>5 findings, >500 line document). The referenced files both exist: `plugins/forge/skills/eval/rules/reviser-composition.md` EXISTS, `plugins/forge/skills/eval/rubrics/proposal.md` EXISTS. The practical guidance (predict processing time, batch by severity, limit MAX_ITERATIONS to 2, monitor for 20-minute timeout) remains applicable to the current eval pipeline. The structural improvement direction (batch-mode reviser, 2-3 attack points per subagent call) remains a valid optimization proposal.
- **Code path verification**: Both referenced files EXIST

### 17. gotcha-eval-rollback-destroys-improvements.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The lesson describes the eval pipeline's rollback mechanism destroying valuable improvements. The current eval SKILL.md has evolved to implement a more nuanced two-level rollback strategy: (1) Inner level restores to post-P0.5 state (pre-revised checkpoint), max 1 per run; (2) The rollback decision is now conditional — it checks whether score improved or degraded. The SKILL.md explicitly states "No rollback for non-proposal types or when pre-revision was skipped." This indicates the lesson's core concern has been partially addressed in the current implementation. The referenced file `plugins/forge/skills/eval/SKILL.md` EXISTS. The proposal directory `docs/proposals/auto-gen-journeys-contracts/` EXISTS.
- **Code path verification**: Both referenced paths EXIST
- **Note**: The lesson's recommendation (user decides, don't auto-rollback) has been partially implemented in the current two-level rollback design

### 18. gotcha-eval-rubric-misses-disguised-patches.md

- **Classification**: experience-summary
- **Status**: valid
- **Justification**: The lesson identifies a fundamental blindspot in the design eval rubric: it measures document-quality metrics (completeness, correctness) but not design-quality metrics (structural quality, intent-output consistency). The referenced rubric `plugins/forge/skills/eval/rubrics/design.md` EXISTS. The referenced design document `docs/features/forge-architecture-simplification/design/tech-design.md` EXISTS. The eval iteration `docs/features/forge-architecture-simplification/eval/iteration-2.md` EXISTS. The core insight — eval rubrics for refactoring designs need an "intent-output consistency" dimension to catch disguised patches — is a generalized principle that remains applicable regardless of rubric evolution.
- **Code path verification**: All 3 referenced files EXIST

### 19. gotcha-eval-subagent-type.md

- **Classification**: experience-summary
- **Status**: outdated
- **Justification**: The lesson recommends using `doc-scorer` / `doc-reviser` subagent types instead of `error-fixer`. The current eval SKILL.md does NOT reference any of these subagent type names. The current implementation uses `general-purpose` agents spawned with `model: "sonnet"` for both scorer and reviser roles. The `zcode:error-fixer` type mentioned in the problem description is also not found in the current codebase. While the principle (match subagent type to task semantics) is valid, the specific type names in this lesson (`doc-scorer`, `doc-reviser`, `error-fixer`) do not correspond to any current subagent configuration.
- **Code path verification**: No `doc-scorer`, `doc-reviser`, or `error-fixer` subagent types found in current eval implementation
- **Recommendation**: Mark as outdated — the subagent type names are not used in the current implementation. This is a duplicate concern of item 15 (gotcha-eval-prd-use-zcode-agents) which recommends different but equally outdated type names.

### 20. gotcha-fix-task-broad-scope.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes quality-gate fix tasks with overly broad scope (20+ failures across 4 test suites in a single fix task). The referenced hook directory `plugins/forge/hooks/` EXISTS. The `tests/results/raw-output.txt` does NOT exist at the current path — tests are now under `forge-cli/tests/results/` (a `.gitkeep` exists there). The cross-reference to `gotcha-split-rules-operational-blindness.md` EXISTS. The two-layer solution (per-suite task splitting via grep, baseline comparison filtering) is a sound operational pattern. The specific shell commands (`grep "^FAIL\s" tests/results/raw-output.txt`) reference the old path.
- **Code path verification**: `plugins/forge/hooks/` EXISTS, `tests/results/raw-output.txt` MISSING (path is now `forge-cli/tests/results/`), `gotcha-split-rules-operational-blindness.md` EXISTS
- **Required update**: (1) Update `tests/results/` path references to `forge-cli/tests/results/`; (2) Verify if the proposed per-suite splitting has been implemented in the hook

## Duplicate Detection (Topic Clustering)

### Cluster 1: Eval Subagent Type Selection

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-eval-prd-use-zcode-agents.md | outdated | **KEEP** — more detailed, includes architectural comparison with eval-proposal's correct dispatch pattern |
| gotcha-eval-subagent-type.md | outdated | **DUPLICATE of gotcha-eval-prd-use-zcode-agents** — same core issue (wrong subagent type for eval tasks) but less detailed. Both are outdated with respect to current `general-purpose` agent implementation, but the prd-use-zcode-agents version provides the full architectural analysis and correct pattern. |

**Verdict**: `gotcha-eval-subagent-type.md` is a duplicate of `gotcha-eval-prd-use-zcode-agents.md`. Keep the prd version (more complete architectural analysis). Both are outdated but the prd version has higher salvage value for the principle-level guidance.

### Cluster 2: Quality Gate + Breaking Task Interactions

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-breaking-change-quality-gate-deadlock.md | valid | Keep — describes compilation boundary scoping for breaking tasks |
| gotcha-breaking-task-quality-gate-test-scope.md | needs-update | Keep — describes structural conflict between breaking task deferral and quality gate requirements |
| gotcha-dispatcher-ignores-compilation-diagnostics.md | needs-update | Keep — describes self-reported quality gate results being falsified |

**Verdict**: NOT duplicate. These describe different aspects of the quality gate system: scoping (#1), test deferral conflict (#2), and result falsification (#3). All are independently valuable.

### Cluster 3: E2E Test Infrastructure

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-e2e-env-override-isolation.md | needs-update | Keep — env var isolation pattern |
| gotcha-e2e-test-binary-isolation.md | needs-update | Keep — binary build isolation pattern |
| gotcha-e2e-skill-monorepo-path-mismatch.md | outdated | Keep — monorepo path resolution pattern (though skill removed) |
| gotcha-e2e-test-quality-antipatterns.md | outdated | Keep — test quality antipattern catalog (though files reorganized) |
| gotcha-e2e-script-generation.md | valid | Keep — SKILL.md structural blind spots analysis |

**Verdict**: NOT duplicate. These cover different E2E infrastructure issues: env isolation (#8), binary isolation (#11), path resolution (#10), test quality (#12), and skill design (#9). All address distinct failure modes.

### Cluster 4: Eval Pipeline Issues

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-eval-loop-decision-gate.md | valid | Keep — decision gate rule |
| gotcha-eval-reviser-too-many-attacks.md | valid | Keep — reviser performance |
| gotcha-eval-rollback-destroys-improvements.md | valid | Keep — rollback strategy |
| gotcha-eval-rubric-misses-disguised-patches.md | valid | Keep — rubric blindspot |

**Verdict**: NOT duplicate. These cover distinct eval pipeline issues: decision gate enforcement (#14), reviser performance (#16), rollback strategy (#17), and rubric completeness (#18). All are independently valuable.

### Cluster 5: Drift Detection and Task Template Quality

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-drift-detection-task-runtime.md | outdated | **KEEP** — describes the empty-template problem (now resolved) |
| gotcha-fix-task-broad-scope.md | needs-update | Keep — fix task scope granularity (different root cause) |

**Verdict**: NOT duplicate. Different subsystems and root causes. The drift detection item describes empty template content; the fix task item describes overly broad automatic task generation.

### Cross-Batch Duplicate Check (vs Batch 1)

The following items in this batch share topic areas with Batch 1 items but are NOT duplicates:

| This Batch Item | Batch 1 Item | Relationship |
|-----------------|--------------|-------------|
| gotcha-dispatcher-ignores-compilation-diagnostics.md | gotcha-ac-self-report-without-verification.md | Same theme (self-reported results trust), but this item focuses on quality gate falsification while Batch 1 focuses on AC self-reporting without verification |
| gotcha-breaking-change-quality-gate-deadlock.md | gotcha-breaking-change-integration-test-blast-radius.md | Same theme (breaking change scoping), but this item addresses compilation boundaries while Batch 1 addresses integration test fixture scope |

## Audit Quality Review

- **Sampling ratio**: 100% (all 20 items audited)
- **Cross-layer check**: Performed against L1 core-docs report, L2 conventions batch 1/2 reports, and L2 business-rules report
- **Code path verification**: All `find`/`grep` checks executed against current codebase at baseline commit
- **Key finding**: 3 items (#6 drift-detection, #7 duplicate-test-runs, #13 embedded-template) describe problems that have been fully resolved in the current codebase
- **Key finding**: 2 items (#10 e2e-skill-monorepo, #12 e2e-test-quality) reference skills and test structures that have been removed or completely reorganized
- **Key finding**: 2 items (#15 eval-prd-use-zcode, #19 eval-subagent-type) recommend subagent types that do not exist in the current implementation — identified as duplicates of each other
- **Key finding**: Multiple items reference `tests/e2e/` directory structure which no longer exists — tests have been reorganized under `forge-cli/tests/`
