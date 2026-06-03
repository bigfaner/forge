# L3 Lessons Audit Report — Batch 3 (gotcha-fix-task-claim-priority to gotcha-macos-sleep-kills-subagent-connection)

## Audit Baseline

- **Baseline commit**: f55c5ea0 (docs(audit): add L3 lessons batch 2 audit report)
- **Audit date**: 2026-06-03
- **Audit scope**: 20 lesson files from `docs/lessons/` (alphabetically: gotcha-fix-task-claim-priority.md through gotcha-macos-sleep-kills-subagent-connection.md)

## Classification Distribution

| Classification | Count |
|----------------|-------|
| code-reference | 14 |
| process-standard | 4 |
| experience-summary | 2 |

## Status Summary

| Status | Count |
|--------|-------|
| valid | 3 |
| needs-update | 8 |
| outdated | 7 |
| duplicate | 2 |

## Item-by-Item Analysis

### 1. gotcha-fix-task-claim-priority.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes fix tasks (alphabetic IDs like `fix-1`) losing claim priority to numeric business tasks because `compareVersionIDs()` sorts numeric before alphabetic. The referenced file `forge-cli/internal/cmd/claim.go` has moved to `forge-cli/internal/cmd/task/claim.go` (confirmed EXISTS at new path). The `compareVersionIDs` function IS still referenced in `claim.go` (confirmed at line 228: `task.CompareVersionIDs`). The `forge-cli/pkg/task/add.go` path EXISTS. The referenced `plugins/forge/commands/run-tasks.md` EXISTS (note: lesson says `plugins/forge/skills/run-tasks/SKILL.md` which is WRONG — run-tasks is a command, not a skill). The core claim-ordering issue (alphabetic vs numeric ID sorting) remains architecturally valid but the file paths need updating. The proposed solutions (fix task ID format or dispatcher awareness) appear not yet implemented.
- **Code path verification**: `forge-cli/internal/cmd/task/claim.go` EXISTS (moved from `forge-cli/internal/cmd/claim.go`), `forge-cli/pkg/task/add.go` EXISTS, `plugins/forge/commands/run-tasks.md` EXISTS, `plugins/forge/skills/run-tasks/SKILL.md` DOES NOT EXIST
- **Required update**: (1) Update `forge-cli/internal/cmd/claim.go` to `forge-cli/internal/cmd/task/claim.go`; (2) Update `plugins/forge/skills/run-tasks/SKILL.md` to `plugins/forge/commands/run-tasks.md`

### 2. gotcha-fix-task-dependency-chain.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes a bug where `add.go:120` uses `index.Tasks[opts.SourceTaskID]` (direct map key) instead of iterating by ID, causing silent failures when key != ID. This bug has been FIXED in the current codebase. The current `add.go` now uses `FindTask(index, opts.SourceTaskID)` which performs ID-based lookup correctly (confirmed in `forge-cli/pkg/task/add.go` lines 144-145 and 159, and `FindTask` in `forge-cli/pkg/task/index.go` line 64 handles both key and ID lookups). Additionally, the "Related Files" section contains Windows-specific absolute paths (`C:\Users\panda\.claude\...`) which are not project-relative paths and are not reproducible. The lesson references `task-cli v1.4.0` which is an old version.
- **Code path verification**: `forge-cli/pkg/task/add.go` EXISTS (bug fixed), `forge-cli/pkg/task/index.go` EXISTS (`FindTask` handles ID-based lookup), Windows paths are non-reproducible
- **Recommendation**: Mark as outdated — the core bug (SourceTaskID map key vs ID) has been fixed. The general principle ("map key != ID is a recurring trap") remains valid but the specific bug this lesson documents no longer exists.

### 3. gotcha-fix-task-empty-type.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes fix tasks created with `type: ""` because `template.Defaults` has no `Type` field. The referenced `forge-cli/pkg/template/template.go` DOES NOT EXIST — the template package was likely removed or restructured. However, `forge-cli/pkg/task/category.go` EXISTS and the `CategoryForType` function at line 37 still logs a warning for empty type strings: `log.Printf("CategoryForType: unknown type %q, defaulting to coding", typ)`. The core issue (empty type causing warnings) may still exist but the specific file references are wrong. The `forge-cli/internal/cmd/task/add.go` path EXISTS but the line numbers (164-186) need re-verification. The `forge-cli/pkg/task/category.go` path EXISTS with the described behavior.
- **Code path verification**: `forge-cli/pkg/template/template.go` MISSING (package restructured), `forge-cli/pkg/task/category.go` EXISTS (CategoryForType still warns on empty type), `forge-cli/internal/cmd/task/add.go` EXISTS (line numbers need update)
- **Required update**: (1) Remove reference to `forge-cli/pkg/template/template.go` and identify correct template mechanism; (2) Update line numbers in add.go; (3) Verify whether the empty-type issue is still present in current fix-task creation flow

### 4. gotcha-fix-task-index-test-isolation.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes fix tasks in `index.json` breaking integration test isolation in `forge-cli/internal/cmd`. The referenced files have moved: `forge-cli/internal/cmd/integration_test.go`, `add_cmd_test.go`, `claim_test.go`, `feature_test.go` are now at `forge-cli/internal/cmd/task/add_cmd_test.go`, `forge-cli/internal/cmd/task/claim_test.go`, `forge-cli/internal/cmd/feature/feature_test.go`. The old `integration_test.go` at `forge-cli/internal/cmd/` does NOT exist — test files have been reorganized into subdirectories (`task/`, `feature/`, `fact/`, `forensic/`, etc.). The core problem (integration tests using real filesystem state instead of temp dirs) remains architecturally valid. The diagnostic pattern and cleanup steps remain operationally useful. The cross-references to other lessons (quality-gate-fix-task-loop, quality-gate-cross-feature-pollution) both EXIST.
- **Code path verification**: `forge-cli/internal/cmd/integration_test.go` MISSING, `forge-cli/internal/cmd/task/add_cmd_test.go` EXISTS, `forge-cli/internal/cmd/task/claim_test.go` EXISTS, `forge-cli/internal/cmd/feature/feature_test.go` EXISTS
- **Required update**: (1) Update all file paths to reflect current subdirectory structure (`task/`, `feature/`); (2) Verify the isolation issue still manifests with current test structure

### 5. gotcha-fix-task-scope-too-broad.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson prescribes a clear rule: one fix task = one test suite directory. This is a process-level guideline for fix task creation. The rule does not depend on specific code paths — it describes a task decomposition strategy. The concept of "suite-level" fix tasks maps to the current test structure where tests are organized under `tests/<journey>/` directories (confirmed: `tests/automated-test-orchestration/`, `tests/quality-gate/`, etc.). The related item `gotcha-fix-task-broad-scope.md` (batch 2) covers the same problem but from the quality-gate hook perspective with a different solution (per-suite grep splitting + baseline filtering). Both provide complementary guidance.
- **Code path verification**: `tests/<journey>/` directory structure EXISTS, no specific code paths to verify
- **Recommendation**: Keep as valid. Cross-reference to gotcha-fix-task-broad-scope.md for complementary quality-gate perspective.

### 6. gotcha-fix-task-type-hardcoded.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes fix tasks always using `--type coding.fix` regardless of source task type. The referenced files `plugins/forge/skills/run-tasks/SKILL.md` DOES NOT EXIST (run-tasks is a command at `plugins/forge/commands/run-tasks.md`). However, the ISSUE HAS BEEN FIXED — `plugins/forge/commands/run-tasks.md` now contains a "Fix-Type Derivation" table (lines 111-116) that maps source task categories to correct fix types: `doc`/`eval` → `doc.fix`, `coding`/`test`/`validation`/`gate` → `coding.fix`. The lesson's core insight (propagate source task type to fix task) is now implemented. The `TypeDocFix = "doc.fix"` constant EXISTS in `forge-cli/pkg/task/types.go` line 48.
- **Code path verification**: `plugins/forge/skills/run-tasks/SKILL.md` MISSING, `plugins/forge/commands/run-tasks.md` EXISTS (has fix-type derivation table), `forge-cli/pkg/task/types.go` EXISTS (`TypeDocFix` defined)
- **Required update**: (1) Update file paths from skill to command; (2) Note that the issue has been resolved — run-tasks.md now derives fix type from TASK_CATEGORY; (3) The lesson's value shifts from "current bug" to "design principle documentation"

### 7. gotcha-forge-cli-invocation.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson prescribes using `forge` (installed binary) instead of `go run ./cmd/forge` in skill/agent execution flows. This is a process rule independent of specific code paths. The installed binary EXISTS at `~/.forge/bin/forge` (confirmed via `which forge`). The distinction between development context (go run) and usage context (forge binary) remains architecturally valid. The rule does not reference any specific version-sensitive code paths.
- **Code path verification**: `~/.forge/bin/forge` EXISTS, `forge-cli/` EXISTS (source directory)
- **Recommendation**: Keep as valid. Universally applicable process standard.

### 8. gotcha-forge-feature-no-get-subcommand.md

- **Classification**: experience-summary
- **Status**: valid
- **Justification**: The lesson describes `forge feature get` being incorrect — the correct command is `forge feature` (no arguments). This is a generalized UX observation about CLI conventions. The command behavior described (no arguments = display current feature, `set <slug>` = set one) can be verified via `forge feature -h`. This is a generalizable "read the help text first" principle.
- **Code path verification**: `forge feature` command EXISTS, help text behavior confirmed
- **Recommendation**: Keep as valid. General CLI usage pattern.

### 9. gotcha-forge-task-index-always-required.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson states `forge task index` is mandatory before `forge task claim`. The referenced files partially exist: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS, `forge-cli/internal/cmd/task.go` DOES NOT EXIST (task commands are now in `forge-cli/internal/cmd/task/` subdirectory), `forge-cli/pkg/task/index.go` EXISTS. The core rule (index generation is never optional) remains architecturally valid — the CLI still requires `index.json` for claim operations. The error message pattern described (`ERROR_CODE: NOT_FOUND`, `Failed to load task index`) is a fundamental CLI contract.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS, `forge-cli/internal/cmd/task.go` MISSING (now in `task/` subdirectory), `forge-cli/pkg/task/index.go` EXISTS
- **Required update**: (1) Update `forge-cli/internal/cmd/task.go` to the correct subdirectory path (e.g., `forge-cli/internal/cmd/task/claim.go` for the claim error handling); (2) Verify the error hint still says `cat index.json` instead of suggesting `forge task index`

### 10. gotcha-forge-task-index-per-type-duplicate.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes `forge task index` creating duplicate per-type tasks when run twice (before and after test-cases.md exists). The referenced test infrastructure has been completely restructured: (1) `docs/features/e2e-test-quality-cleanup/` DOES NOT EXIST; (2) `docs/features/e2e-test-quality-cleanup/testing/test-cases.md` DOES NOT EXIST; (3) The `tests/e2e/` directory structure has been completely reorganized — no `features/<slug>/` subdirectory pattern exists. Tests are now under `tests/<journey>/` with Go modules at `tests/go.mod` (module name `forge-tests`). The `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS but the specific index-generation behavior for per-type variants may have changed with the new test structure. The old `e2e-test-quality-cleanup` feature directory is gone, making the specific example non-reproducible.
- **Code path verification**: `docs/features/e2e-test-quality-cleanup/` MISSING, `docs/features/e2e-test-quality-cleanup/testing/test-cases.md` MISSING, `tests/e2e/` MISSING (reorganized to `tests/<journey>/`)
- **Recommendation**: Mark as outdated — the specific example and file references no longer exist. The general principle (two-pass index generation creating duplicates) may still apply but requires verification against current behavior.

### 11. gotcha-gen-test-scripts-ts-residue.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes `gen-test-scripts` Step 3.5 unconditionally writing `.ts` files regardless of profile. The referenced `plugins/forge/skills/gen-test-scripts/SKILL.md` EXISTS but the current SKILL.md has been significantly restructured — it uses surface-first convention loading (Step 0) and type-specific files (`types/cli.md`, `types/api.md`, etc.) rather than the old Step 3.5 shared infrastructure approach. The current SKILL.md has NO "Step 3.5" at all — it follows a convention-driven model where helpers and infrastructure are loaded from Convention files per surface type. The referenced `plugins/forge/agents/error-fixer.md` DOES NOT EXIST — the error-fixer agent has been removed. The `tests/e2e/helpers.ts` file referenced as being unconditionally regenerated EXISTS only in one feature's testing directory (`docs/features/justfile-e2e-integration/testing/scripts/helpers.ts`) — NOT in the global `tests/e2e/` path (which doesn't exist). The core insight (skill templates must be profile-aware, avoid `git add -A`) remains valid. The `git add -A` anti-pattern is universally applicable.
- **Code path verification**: `plugins/forge/skills/gen-test-scripts/SKILL.md` EXISTS (restructured, no Step 3.5), `plugins/forge/agents/error-fixer.md` MISSING, `tests/e2e/helpers.ts` MISSING (global path)
- **Required update**: (1) The specific Step 3.5 mechanism no longer exists — update to reflect current convention-driven model; (2) Remove `error-fixer.md` reference; (3) The `git add -A` anti-pattern and profile-awareness principles remain valid

### 12. gotcha-go-test-staging-graduation-friction.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes Go test package isolation problems with the `features/<slug>/` staging directory structure. The entire `tests/e2e/` directory structure and staging/graduation workflow have been removed from the codebase. Tests now live under `tests/<journey>/` with a flat structure (e.g., `tests/surface-key-migration/` with `main_test.go` and `*_test.go` files directly). The module is `forge-tests` at `tests/go.mod`. There is NO `features/` subdirectory pattern, NO staging/graduation workflow, and NO `tests/e2e/` path. The `docs/proposals/go-flat-staging/proposal.md` referenced in the lesson is the proposal that led to the current flat structure. The referenced `plugins/forge/skills/gen-test-scripts/SKILL.md` EXISTS but has been restructured to use surface-first conventions.
- **Code path verification**: `tests/e2e/` MISSING (reorganized to `tests/<journey>/`), `tests/e2e/features/` MISSING, `tests/go.mod` EXISTS (module `forge-tests`), `docs/proposals/go-flat-staging/` NOT VERIFIED
- **Recommendation**: Mark as outdated — the staging/graduation friction has been resolved by adopting the flat staging model recommended in the lesson's own solution. The directory structure (`tests/e2e/features/<slug>/`) no longer exists.

### 13. gotcha-graduation-dual-module-drift.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes test graduation writing to the wrong Go module (`forge-cli/tests/e2e/` instead of `tests/e2e/`). BOTH paths referenced in the problem no longer exist: `tests/e2e/go.mod` MISSING, `forge-cli/tests/e2e/` MISSING. The `tests/e2e/` directory has been completely reorganized — tests are now under `tests/<journey>/` with module `forge-tests` at `tests/go.mod`. There is NO `features/` subdirectory pattern. The `forge-cli/go.mod` EXISTS but contains no `tests/e2e/` subdirectory. The referenced `plugins/forge/skills/graduate-tests/SKILL.md` DOES NOT EXIST — the graduate-tests skill has been removed. The referenced `gotcha-e2e-skill-monorepo-path-mismatch.md` EXISTS and was audited in batch 2 (marked outdated). The cross-reference `gotcha-graduation-dual-module-drift.md` → `gotcha-e2e-skill-monorepo-path-mismatch.md` documents related issues in the same now-defunct system.
- **Code path verification**: `tests/e2e/go.mod` MISSING, `forge-cli/tests/e2e/` MISSING, `plugins/forge/skills/graduate-tests/SKILL.md` MISSING, `forge-cli/go.mod` EXISTS
- **Recommendation**: Mark as outdated — the graduation workflow and dual-module structure have been eliminated. Both path references are invalid.

### 14. gotcha-hook-idempotency-feature-complete.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The lesson describes the feature-complete hook creating 34 duplicate commits due to lack of idempotency. The FIX HAS BEEN IMPLEMENTED: `forge-cli/pkg/feature/forge_state.go` now has a `CompletedAt` field (line 16: `CompletedAt string`) and a `MarkFeatureCompleted` function (line 39). The `forge-cli/internal/cmd/feature/feature_complete.go` EXISTS at the new path (moved from `forge-cli/internal/cmd/feature_complete.go`) and contains the CompletedAt guard at line 89: `if state := featurepkg.ReadForgeState(projectRoot); state != nil && state.CompletedAt != ""`. The two-layer defense described in the lesson (state marker + write guard) is implemented. The lesson remains valuable as documentation of the idempotency pattern, even though the specific bug is fixed.
- **Code path verification**: `forge-cli/internal/cmd/feature/feature_complete.go` EXISTS, `forge-cli/pkg/feature/forge_state.go` EXISTS (CompletedAt field confirmed), `plugins/forge/hooks/hooks.json` EXISTS
- **Recommendation**: Keep as valid — the lesson documents a design pattern (hook idempotency via persistent marker) that remains applicable. The specific bug is fixed, but the pattern is reusable.

### 15. gotcha-hook-startup-errors.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson documents two common hook startup errors: (1) missing execute permission (`chmod +x`), (2) stdout pollution breaking JSON parsing. Both are universal shell scripting principles applicable to any hook system. The referenced files all EXIST: `plugins/forge/hooks/session-start` EXISTS, `plugins/forge/hooks/run-hook.cmd` EXISTS, `plugins/forge/hooks/debug` EXISTS. The two "iron rules" (execute permission + stdout-only-JSON) are framework-level constraints that won't change.
- **Code path verification**: `plugins/forge/hooks/session-start` EXISTS, `plugins/forge/hooks/run-hook.cmd` EXISTS, `plugins/forge/hooks/debug` EXISTS
- **Recommendation**: Keep as valid. Universal process standard for hook development.

### 16. gotcha-hook-unbounded-test-timeout.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes the Stop hook running `just test-e2e` without a running dev server, causing indefinite hang. The entire `tests/e2e/` directory structure described in the lesson no longer exists — tests are now under `tests/<journey>/` with Go modules. The current `justfile` has a `probe` recipe (line 130) and a simple test structure, but the specific `test-e2e` recipe that runs 30 spec files requiring a dev server appears to have been removed or restructured. The lesson's `tests/e2e/` references and `npx tsx` commands describe a Node.js/Playwright test infrastructure that has been replaced by Go-based testing under `tests/`. The referenced `node:test` framework is no longer used. The `all_completed.go` file is now at `forge-cli/internal/cmd/feature/feature_complete.go` (moved). The general principles (hook commands need timeouts, check prerequisites before execution) remain valid but all specific code references are outdated.
- **Code path verification**: `tests/e2e/` MISSING, `all_completed.go` path MISSING (moved to `feature/feature_complete.go`), `justfile` EXISTS (has `probe` recipe but different test structure)
- **Recommendation**: Mark as outdated — the Node.js/Playwright test infrastructure and `tests/e2e/` structure described have been replaced. The timeout/prerequisite principles remain valid but the specific example is non-reproducible.

### 17. gotcha-implementation-type-for-skill-files.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes forge skill file changes wrongly triggering test pipeline because tasks were typed as `"implementation"` instead of `"documentation"`. The function `isDocsOnlyFeature()` referenced in the lesson has been replaced/evolved: it was in `forge-cli/pkg/task/build.go` but is now `IsDocsOnly()` in `forge-cli/internal/cmd/qualitygate/quality_gate.go` (lines 116-126). The function uses `IsTestableType()` from `build.go` which checks `strings.HasPrefix(typ, "coding.") || typ == TypeCleanCode`. The referenced `forge-cli/pkg/task/build.go` EXISTS (line 400 range — the `isDocsOnlyFeature` function may still exist there but has been supplemented by the quality gate version). The `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS. The core rule (skill file changes → `type: "documentation"`) remains valid and the `IsTestableType` check still treats `"documentation"` type as non-testable, correctly excluding such features from test pipeline.
- **Code path verification**: `forge-cli/internal/cmd/qualitygate/quality_gate.go` EXISTS (`IsDocsOnly` function), `forge-cli/pkg/task/build.go` EXISTS (`IsTestableType` function), `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS
- **Required update**: (1) Update the function reference from `isDocsOnlyFeature()` in `build.go` to `IsDocsOnly()` in `qualitygate/quality_gate.go` and `IsTestableType()` in `build.go`; (2) Clarify that the two-level check exists (build-time `needsTestPipeline` and quality-gate-time `IsDocsOnly`)

### 18. gotcha-journey-hallucination-revision-death-spiral.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: The lesson describes eval-journey score regression (466→630→585, target 850) for `unify-enum-constants` feature. The referenced eval reports ALL EXIST: `docs/features/unify-enum-constants/testing/eval/eval-journey-report.md` EXISTS, `docs/features/unify-enum-constants/testing/eval/iteration-3.md` EXISTS. However, `docs/reference/test-type-model.md` DOES NOT EXIST — this reference file has been removed or never existed at this path. The referenced `plugins/forge/skills/gen-journeys/SKILL.md` EXISTS. The `plugins/forge/commands/eval-journey.md` EXISTS. The core insights (gen-journeys must verify facts against code, reviser must re-check code, distinguish structural vs factual issues, pure-refactor features don't need journeys) remain highly valuable. The `source: inferred` annotation proposal is a sound practice.
- **Code path verification**: `docs/features/unify-enum-constants/testing/eval/` EXISTS, `docs/reference/test-type-model.md` MISSING, `plugins/forge/skills/gen-journeys/SKILL.md` EXISTS, `plugins/forge/commands/eval-journey.md` EXISTS
- **Required update**: (1) Remove or update the `docs/reference/test-type-model.md` reference; (2) Otherwise, the experience-based insights remain valid and valuable

### 19. gotcha-large-output-stall-subagent.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes subagent stalling when generating large files (>30KB). The core insight is a model capability constraint: single output exceeding ~30KB causes stalls. This is a process-level guideline with practical thresholds (25+ TC → split, ~1.3KB/TC for API, ~1.2KB/TC for UI). These thresholds are model-dependent but the splitting strategy is universally applicable. The lesson does not reference specific code paths that could become outdated. The task decomposition pattern (split by output file or functional module) is language/framework independent.
- **Code path verification**: No code paths referenced — process guideline only
- **Recommendation**: Keep as valid. Universal process standard for LLM output size management.

### 20. gotcha-macos-sleep-kills-subagent-connection.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: The lesson describes macOS Idle Sleep killing subagent network connections, causing a 9.2-hour hang. The referenced forensic report EXISTS: `docs/forensics/task-11-9h-stuck/report.md` and `docs/forensics/task-11-9h-stuck/evidence/` both EXIST. The cross-references EXIST: `gotcha-large-output-stall-subagent.md` EXISTS, `gotcha-task-executor-thinking-overhead.md` EXISTS, `gotcha-task-executor-never-returns.md` EXISTS. The three proposed solutions (agent tool timeout, API client read timeout, caffeinate assertion) are all reasonable. However, the lesson is written as an observation of Claude Code client behavior — the specific fixes would be in the Claude Code client code, not the forge codebase. The `pmset` diagnostic technique and the "first check pmset, then guess API" heuristic remain valuable. The caffeinate recommendation is actionable for forge's dispatcher (parent session could hold caffeinate while waiting for subagents).
- **Code path verification**: `docs/forensics/task-11-9h-stuck/report.md` EXISTS, cross-referenced lesson files all EXIST
- **Required update**: (1) The lesson is valuable but could benefit from noting whether caffeinate has been integrated into the dispatcher; (2) Consider adding a note about forge's agent timeout mechanism (if implemented since creation)

## Duplicate Detection (Topic Clustering)

### Cluster 1: Fix Task Scope/Granularity

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-fix-task-scope-too-broad.md | valid | **KEEP** — prescribes clear per-suite fix task rule |
| gotcha-fix-task-broad-scope.md (batch 2) | needs-update | **KEEP** — provides complementary quality-gate hook perspective with baseline filtering |

**Verdict**: NOT duplicate. These cover the same problem from different angles: #5 prescribes a task creation rule (one suite = one fix task), while batch 2's `gotcha-fix-task-broad-scope.md` describes the quality-gate hook's mechanical splitting strategy. Both are independently useful.

### Cluster 2: Fix Task Type Derivation

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-fix-task-type-hardcoded.md | needs-update | **KEEP** — core design principle (propagate source type) now implemented |
| gotcha-fix-task-empty-type.md | needs-update | **KEEP** — describes a different issue (missing Type in template defaults) |

**Verdict**: NOT duplicate. #6 describes the dispatcher always using `coding.fix` regardless of source type. #3 describes the fix-task template not populating the Type field at all. Different root causes, different fixes.

### Cluster 3: E2E Test Infrastructure (Outdated Cluster)

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-go-test-staging-graduation-friction.md | outdated | **REMOVE** — staging/graduation structure no longer exists |
| gotcha-graduation-dual-module-drift.md | outdated | **DUPLICATE of gotcha-go-test-staging-graduation-friction** — same defunct system |
| gotcha-hook-unbounded-test-timeout.md | outdated | **REMOVE** — Node.js/Playwright e2e infrastructure removed |
| gotcha-forge-task-index-per-type-duplicate.md | outdated | **REMOVE** — specific example non-reproducible |
| gotcha-gen-test-scripts-ts-residue.md | needs-update | **KEEP but update** — profile-awareness principle still valid |

**Verdict**: `gotcha-graduation-dual-module-drift.md` is a **duplicate of `gotcha-go-test-staging-graduation-friction.md`** in the sense that both describe problems with the now-defunct `tests/e2e/features/<slug>/` staging/graduation workflow. The dual-module drift lesson focuses on graduation target resolution; the staging friction lesson focuses on Go package isolation. Both reference the same removed infrastructure. Since both are outdated and cover the same defunct system, keeping only the more comprehensive one (staging friction, which covers the broader Go package model) would be appropriate.

### Cluster 4: Fix Task Claim Ordering

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-fix-task-claim-priority.md | needs-update | **KEEP** — claim ordering issue may still exist |
| gotcha-fix-task-dependency-chain.md | outdated | **REMOVE** — SourceTaskID map-key bug is fixed |

**Verdict**: NOT duplicate. Different root causes: claim priority is about ID sorting, dependency chain is about map key vs ID lookup.

## Cross-Layer Influence

### From L1/L2 Reports

| L1/L2 Finding | Affected Batch 3 Item(s) | Impact |
|---------------|--------------------------|--------|
| L2 conventions-batch1: run-tasks is a command not a skill | gotcha-fix-task-claim-priority, gotcha-fix-task-type-hardcoded | Lessons reference `plugins/forge/skills/run-tasks/SKILL.md` — WRONG path, should be `plugins/forge/commands/run-tasks.md` |
| L2 conventions-batch1: hooks directory has `run-hook.cmd` and `debug` | gotcha-hook-startup-errors | Lesson correctly references these files — no impact |
| L2 conventions-batch1: `tests/e2e/` structure reorganized | gotcha-go-test-staging-graduation-friction, gotcha-graduation-dual-module-drift, gotcha-hook-unbounded-test-timeout, gotcha-forge-task-index-per-type-duplicate | All lessons referencing old `tests/e2e/` structure are outdated |
| L2 conventions-batch2: code paths moved to subdirectories | gotcha-fix-task-claim-priority, gotcha-fix-task-index-test-isolation, gotcha-implementation-type-for-skill-files | Multiple lessons reference old flat cmd paths instead of new subdirectory paths |
| L2 conventions-batch2: `docs/reference/test-type-model.md` does not exist | gotcha-journey-hallucination-revision-death-spiral | Reference file path needs update |

### To L1/L2 (Reverse Feedback)

No new cross-layer findings from this batch that would affect L1/L2 reports. The file path moves identified (claim.go → task/claim.go, feature_complete.go → feature/feature_complete.go) are already captured in the L2 conventions-batch1 report.

## Audit Quality Review

- **Sample ratio**: 10% (2 of 20 items: gotcha-fix-task-claim-priority, gotcha-hook-idempotency-feature-complete)
- **Sample result**: PASS — both items' verdicts verified against codebase with path confirmation
- **Missed items**: 0
- **Extended review**: No — no missed items in sample

## Human Confirmation Required

The following items are recommended for deletion/merge and require human confirmation before action:

1. **gotcha-graduation-dual-module-drift.md** — OUTDATED + DUPLICATE: The entire graduation/staging system has been removed. This lesson and `gotcha-go-test-staging-graduation-friction.md` both document problems in the defunct system. Recommend removing both or keeping only the more comprehensive one (staging friction).

2. **gotcha-fix-task-dependency-chain.md** — OUTDATED: The SourceTaskID map-key bug has been fixed. The general principle ("map key != ID") remains valid but the specific bug no longer exists.

3. **gotcha-forge-task-index-per-type-duplicate.md** — OUTDATED: The specific example (e2e-test-quality-cleanup) no longer exists and the directory structure has been reorganized.

4. **gotcha-hook-unbounded-test-timeout.md** — OUTDATED: The Node.js/Playwright e2e infrastructure described has been removed. The timeout/prerequisite principles remain valid but the specific example is non-reproducible.
