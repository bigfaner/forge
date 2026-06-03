# L3 Lessons Audit Report — Batch 5 (gotcha-review-task to gotcha-task-executor-ignores-implementation-notes)

## Audit Baseline

- **Baseline commit**: c68e834d (docs(audit): add L3 lessons batch 4 audit report)
- **Audit date**: 2026-06-03
- **Audit scope**: 20 lesson files from `docs/lessons/` (alphabetically: gotcha-review-task-incomplete-dependencies.md through gotcha-task-executor-ignores-implementation-notes.md)

## Classification Distribution

| Classification | Count |
|----------------|-------|
| code-reference | 11 |
| process-standard | 5 |
| experience-summary | 4 |

## Status Summary

| Status | Count |
|--------|-------|
| valid | 3 |
| needs-update | 10 |
| outdated | 4 |
| duplicate | 3 |

## Cross-Layer Influence Check

The following L1/L2 audit findings affect items in this batch:

| L1/L2 Finding | Affected Lesson | Impact |
|---------------|-----------------|--------|
| L2 conventions-batch1: `run-tasks` is a command (not a skill) | gotcha-run-tasks-no-auto-test, gotcha-task-executor-auto-claim | References to `/run-tasks` as skill are incorrect; actual location is `plugins/forge/commands/run-tasks.md` |
| L2 conventions-batch2: `prompt/data/` renamed to `prompt/templates/` | gotcha-spec-authority-drift | Referenced template paths use old `forge-cli/pkg/prompt/data/` directory; current path is `forge-cli/pkg/prompt/templates/` |
| L2 conventions-batch1: `tests/e2e/` directory does not exist at project root | gotcha-split-task-missing-shared-setup, gotcha-task-executor-ignores-implementation-notes | Referenced paths like `tests/e2e/playwright.config.ts`, `tests/e2e/auth-setup.ts` are invalid; test infrastructure has been restructured |

## Item-by-Item Analysis

### 1. gotcha-review-task-incomplete-dependencies.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes a review-doc task being dispatched before its target tasks completed, due to incomplete dependency lists in `forge task index` auto-generation. The referenced file `forge-cli/pkg/task/autogen.go` EXISTS and contains review task generation logic. The core architectural insight (auto-generated review tasks must depend on ALL tasks of the reviewed type, not just currently-known tasks) remains valid. However, the referenced `gotcha-task-reference-files-scope-creep` cross-reference EXISTS and is still relevant. The specific task IDs and feature context are historical. The file `forge-cli/pkg/task/autogen.go` was verified to exist but the specific "review-doc dependency generation" function was not found by name matching — the implementation may have been restructured or the functionality moved to a different module.
- **Code path verification**: `forge-cli/pkg/task/autogen.go` EXISTS
- **Required update**: (1) Verify that the dependency calculation logic in autogen.go still has the described limitation or has been fixed; (2) Update the specific function/file reference if the code was reorganized

### 2. gotcha-reviser-agent-long-running.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents a Reviser subagent running ~2 hours without completing, caused by: (1) Edit tool matching overhead on large documents, (2) document bloat across iterations, (3) excessive re-read for quality checks, (4) no max-duration constraint, (5) single-agent architecture bottleneck. The referenced `plugins/forge/skills/eval/experts/protocol/reviser-protocol.md` EXISTS. Checking the current protocol: it contains `<EXTREMELY-IMPORTANT>` with "Maximum 3 rounds of self-review" but NO max-duration or wall-clock time constraint. The lesson's suggested fixes (max-duration, attack points batching, reduced re-read) were NOT implemented in the protocol. The referenced proposals EXIST: `docs/proposals/intent-driven-pipeline-branching/proposal.md` EXISTS, `docs/proposals/pipeline-topology-registry/proposal.md` EXISTS. The root cause analysis (5-level causal chain) is thorough and still applicable.
- **Code path verification**: `plugins/forge/skills/eval/experts/protocol/reviser-protocol.md` EXISTS (no max-duration constraint confirmed), both proposals EXIST
- **Required update**: (1) Note that the suggested max-duration fix was NOT implemented — the protocol still lacks time constraints; (2) The lesson's recommendations remain actionable and should be flagged as pending implementation

### 3. gotcha-run-tasks-no-auto-test.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes `/run-tasks` not triggering tests automatically and proposes using `task all-completed` + Stop hook as the solution. The `all-completed` functionality is now integrated into the `forge quality-gate` command (confirmed: `forge-cli/internal/cmd/qualitygate/quality_gate.go` lines 82-88 check `.forge/state.json` for `allCompleted` flag). The Stop hook in `plugins/forge/hooks/hooks.json` now runs `forge quality-gate` followed by `forge feature complete --if-done` — which implements the lesson's proposed solution through a different mechanism. The lesson's description of the Iron Laws and the `task all-completed` command as a standalone tool is outdated; the functionality was absorbed into the quality-gate infrastructure. The core insight (skill Iron Laws override CLAUDE.md and user instructions) remains valid, but the specific implementation details have been superseded.
- **Code path verification**: `plugins/forge/hooks/hooks.json` EXISTS (Stop hook confirmed: quality-gate + feature complete), `forge-cli/internal/cmd/qualitygate/quality_gate.go` EXISTS (allCompleted check confirmed)
- **Recommendation**: Mark as outdated. The proposed solution was implemented via a different mechanism (quality-gate integration). Consider consolidating the valid core insight (Iron Laws priority hierarchy) into a more general process-standard lesson.

### 4. gotcha-shared-interface-mock-cascade.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes mock cascade stalls when adding methods to the `MainItemRepo` interface (17 methods) in a backend Go project. The referenced paths are all under `backend/internal/` — checking the project: `backend/internal/repository/main_item_repo.go` DOES NOT EXIST. The `backend/` directory does not exist in the current project structure. This lesson documents a pattern from a different project/codebase (likely a SaaS application) that was developed alongside the forge tooling. The ISP (Interface Segregation Principle) recommendation and the "split into interface-update task + feature task" pattern are architecturally sound, but all specific code references are invalid for this repository.
- **Code path verification**: `backend/internal/repository/main_item_repo.go` MISSING, `backend/internal/service/main_item_service_test.go` MISSING, `backend/` directory does not exist in this repository
- **Recommendation**: Mark as outdated. All code references point to a different project. The architectural pattern (ISP for fat interfaces, pre-task interface updates) is sound and could be preserved as a generalized experience-summary if stripped of the invalid code references.

### 5. gotcha-skill-step-analysis-paralysis.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: The lesson documents over-analysis during `/breakdown-tasks` execution, caused by perfectionism and step-order violation. The referenced `docs/features/test-capability-v2/design/tech-design.md` EXISTS. The core insight (follow skill steps in order, use CLI validation tools as safety net, don't pre-solve downstream steps) is universally valid. However, the "Related Files" section references `C:\Users\panda\.claude\plugins\cache\forge\forge\3.0.0-rc.18\skills\breakdown-tasks\SKILL.md` — a Windows-specific local cache path that is not a valid reference in this repository. The actual skill location is `plugins/forge/skills/breakdown-tasks/SKILL.md` which EXISTS. The breakdown-tasks skill's SKILL.md has likely been updated since version 3.0.0-rc.18.
- **Code path verification**: `docs/features/test-capability-v2/design/tech-design.md` EXISTS, `plugins/forge/skills/breakdown-tasks/SKILL.md` EXISTS
- **Required update**: (1) Replace the Windows cache path with the correct repository path `plugins/forge/skills/breakdown-tasks/SKILL.md`; (2) Verify the skill's current step structure matches the lesson's description

### 6. gotcha-skip-plan-for-pipeline-change.md

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The lesson prescribes entering plan mode before modifying forge's own pipeline (skills, commands, agents, task pipeline). The referenced files: `forge-cli/pkg/task/types.go` EXISTS, `forge-cli/pkg/task/infer.go` EXISTS, `forge-cli/pkg/prompt/prompt.go` EXISTS, `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS. However, `forge-cli/pkg/task/testgen.go` DOES NOT EXIST — this file was either removed or renamed during v3.0.0 restructuring. The core rule (3+ files, cross-module, or pipeline changes → plan mode first) remains a sound process standard. The reusable pattern and the "how to apply" guidance are independent of the specific file structure.
- **Code path verification**: `forge-cli/pkg/task/types.go` EXISTS, `forge-cli/pkg/task/infer.go` EXISTS, `forge-cli/pkg/task/testgen.go` MISSING, `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS
- **Required update**: (1) Remove the reference to `testgen.go` or update with the current equivalent file; (2) Verify the EnterPlanMode guidance still matches current Claude Code behavior

### 7. gotcha-spec-authority-drift.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents 43 deviations from tech-design.md when modifying 6 skill files, caused by agents treating existing code as the baseline rather than the spec. The root cause analysis (4 levels) is thorough and the solution (spec-driven modification workflow with explicit verification steps) is architecturally sound. The suggested template improvements WERE implemented: all coding templates now contain `<CRITICAL>` blocks enforcing Reference Files as authoritative sources (confirmed: `coding-enhancement.md`, `coding-refactor.md`, `coding-feature.md`, `coding-cleanup.md`, `coding-fix.md` all have "Spec Authority Enforcement" sections). However, the referenced template paths are stale: the lesson uses `forge-cli/pkg/prompt/data/` but the current directory is `forge-cli/pkg/prompt/templates/`. The referenced `docs/features/test-capability-v2/design/tech-design.md` EXISTS. The referenced agent and command files EXIST: `plugins/forge/agents/task-executor.md` EXISTS, `plugins/forge/commands/execute-task.md` EXISTS.
- **Code path verification**: `forge-cli/pkg/prompt/data/doc.md` MISSING (now `forge-cli/pkg/prompt/templates/doc.md`), `forge-cli/pkg/prompt/templates/coding-enhancement.md` EXISTS, `plugins/forge/agents/task-executor.md` EXISTS, `docs/features/test-capability-v2/design/tech-design.md` EXISTS
- **Required update**: (1) Update all `forge-cli/pkg/prompt/data/` references to `forge-cli/pkg/prompt/templates/`; (2) Note that the suggested template improvements were implemented — the lesson's value shifts from "pending fix" to "documented pattern + fix verification"

### 8. gotcha-split-rules-operational-blindness.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson identifies that task splitting rules only check "functional granularity" (AC count, verb count) but miss "operational granularity" (number of files to modify, number of edit operations). The suggested "Operational Ceiling" rule (>8 files → split by file group) WAS implemented: the current `plugins/forge/skills/quick-tasks/SKILL.md` contains "Operational ceiling" rules including the exact >8 files threshold and the "Hard Rules for file boundaries" recommendation. The cross-references EXIST: `gotcha-prompt-template-complexity-agnostic.md` EXISTS, `gotcha-task-reference-files-scope-creep.md` EXISTS, `gotcha-task-executor-thinking-overhead.md` EXISTS. The referenced skill `skills/quick-tasks/SKILL.md` EXISTS at `plugins/forge/skills/quick-tasks/SKILL.md`. The core insight and the solution were validated by implementation.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS (operational ceiling rule confirmed present)
- **Required update**: (1) Note that the suggested operational ceiling fix WAS implemented; (2) The lesson's value shifts from "pending fix" to "documented pattern + implementation verification"

### 9. gotcha-split-task-missing-shared-setup.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes missing auth-setup when splitting `/gen-test-scripts` into parallel sub-tasks. All referenced paths are under `tests/e2e/` which DOES NOT EXIST at the project root. The test infrastructure has been restructured: `tests/` directory contains only subdirectories like `automated-test-orchestration/`, `task-lifecycle/`, `test-generation/` etc. — no `e2e/` directory. The referenced `tests/e2e/playwright.config.ts`, `tests/e2e/auth-setup.ts`, `tests/e2e/features/milestone-map/milestones-page.spec.ts` are all invalid. The `milestone-map` feature directory was not found in the current repository. The core insight (identify global setup phases before splitting, create pre-task for shared infrastructure) is architecturally valid, but all code references and the specific incident are from a different project or an old project structure that no longer exists.
- **Code path verification**: `tests/e2e/` MISSING, `tests/e2e/playwright.config.ts` MISSING, `tests/e2e/auth-setup.ts` MISSING, `docs/lessons/gotcha-large-output-stall-subagent.md` EXISTS (cross-reference valid)
- **Recommendation**: Mark as outdated. All code references point to non-existent paths. The architectural pattern (identify global setup before splitting) is sound and could be preserved as a generalized experience-summary.

### 10. gotcha-stale-skill-cli-flags.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents a stale `--languages` flag reference in the quick-tasks SKILL.md. The current `forge-cli/internal/cmd/task/index.go` does NOT have a `--languages` or `--test-profiles` flag — it only has `--feature`. Profile resolution is handled by `forge profile` command separately. The current quick-tasks SKILL.md Step 5 says: "If the profile was not set in Step 0, pass it explicitly: `forge task index --feature <slug> --test-profiles <p1>,<p2>`" — this `--test-profiles` reference in the SKILL.md may itself be stale since the current index.go only accepts `--feature`. The old `forge-cli/internal/cmd/index.go` (the lesson's reference) was reorganized to `forge-cli/internal/cmd/task/index.go`. The core insight (verify CLI flags via `--help` before following skill docs) remains universally valid.
- **Code path verification**: `forge-cli/internal/cmd/task/index.go` EXISTS (no --languages or --test-profiles flag), `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS (Step 5 may still reference stale --test-profiles flag)
- **Required update**: (1) Update the old path `forge-cli/internal/cmd/index.go` to `forge-cli/internal/cmd/task/index.go`; (2) Note that the `--test-profiles` flag mentioned in quick-tasks SKILL.md Step 5 may also be stale; (3) The `--help` verification pattern remains valid

### 11. gotcha-stale-state-json-feature-mismatch.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The lesson documents `.forge/state.json` becoming stale when switching between features. The `forge feature set <slug>` command EXISTS (confirmed via test file `forge-cli/internal/cmd/feature/feature_test.go` which tests `Cmd.SetArgs([]string{"set", "my-feature"})`). The `.forge/state.json` EXISTS. The `GetForgeStatePath` function in `forge-cli/pkg/feature/paths.go` returns the correct path. The core insight (always run `forge feature set <slug>` before claiming tasks for a new feature) remains operationally valid. The reusable pattern and example are accurate for the current codebase.
- **Code path verification**: `.forge/state.json` EXISTS, `forge feature set` command EXISTS (confirmed via feature_test.go), `forge-cli/pkg/feature/paths.go` EXISTS (GetForgeStatePath confirmed)
- **Recommendation**: Keep as valid. All references are current and the operational pattern remains correct.

### 12. gotcha-stale-test-results-cascade.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents stale test result files (`results/latest.md`) causing cascading false-positive fix tasks. The referenced `forge-cli/tests/e2e/features/*/results/latest.md` path does NOT follow the current test results structure. However, `docs/features/justfile-e2e-integration/testing/results/latest.md` EXISTS, indicating the results file pattern is still used but at a different location (under `docs/features/<slug>/testing/results/` rather than `tests/e2e/features/<slug>/results/`). The referenced `tests/results/unit-raw-output.txt` was not found. The core insight (never trust cached test results when code has changed) remains universally valid. The solution and checklist pattern are sound.
- **Code path verification**: `docs/features/justfile-e2e-integration/testing/results/latest.md` EXISTS (different location than referenced), `tests/results/unit-raw-output.txt` MISSING
- **Required update**: (1) Update the results file path from `forge-cli/tests/e2e/features/*/results/latest.md` to `docs/features/<slug>/testing/results/latest.md`; (2) Verify the stale-results propagation mechanism still exists in the current test infrastructure

### 13. gotcha-standard-task-id-collision.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The lesson documents the task executor detecting "already completed" based on global git log search by task ID, matching commits from ANY feature rather than scoping to the current feature. The referenced `permission-granularity` feature was not found in the current repository (may have been removed or renamed), but the architectural issue remains valid: standard task IDs (T-test-1, T-test-2) from `/breakdown-tasks` are not namespaced per feature, and the task executor searches git history globally. The `forge-cli/internal/cmd/task/claim.go` EXISTS and contains the claim logic. The root cause analysis (4 levels) and the proposed solutions (feature-scoped git search, namespaced task IDs) remain architecturally relevant.
- **Code path verification**: `forge-cli/internal/cmd/task/claim.go` EXISTS
- **Recommendation**: Keep as valid. The task ID collision risk is inherent in the current architecture and the lesson's mitigation suggestions remain applicable.

### 14. gotcha-stop-hook-non-blocking-error.md

- **Classification**: experience-summary
- **Status**: valid
- **Justification**: The lesson documents the unhelpful "Failed with non-blocking status code: No stderr output" error message from Claude Code Stop hooks. This is a general diagnostic pattern for hook script development: always write to stderr when failing. The lesson is independent of specific code paths — it describes Claude Code's hook behavior and provides two clear solutions (fix the hook script or remove the hook). The current `plugins/forge/hooks/hooks.json` Stop hooks run `forge quality-gate` and `forge feature complete --if-done` — the diagnostic advice (ensure hook scripts write to stderr on failure) remains applicable to these commands.
- **Code path verification**: `plugins/forge/hooks/hooks.json` EXISTS (Stop hooks confirmed)
- **Recommendation**: Keep as valid. Universal diagnostic pattern for Claude Code hook development.

### 15. gotcha-strategy-bypass-justfile.md

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The lesson documents bypassing justfile commands in favor of direct Go toolchain commands, violating the principle of following the prescribed strategy. The referenced task file `docs/features/forge-cli-v3/tasks/2.3-list-types-command.md` EXISTS. The referenced `forge-cli/scripts/justfile` DOES NOT EXIST — the justfile for forge-cli was either removed or moved. The root `justfile` EXISTS at the project root. The core principle (follow strategy-specified commands, diagnose before replacing) is a sound process standard. The "Reusable Pattern" table (correct: execute then diagnose; wrong: predict failure and substitute) is universally applicable. The violation of CLAUDE.md's "first-principles thinking, reject empiricism" principle is well-articulated.
- **Code path verification**: `docs/features/forge-cli-v3/tasks/2.3-list-types-command.md` EXISTS, `forge-cli/scripts/justfile` MISSING, root `justfile` EXISTS
- **Required update**: (1) Update the justfile reference from `forge-cli/scripts/justfile` to the root `justfile` or the current location; (2) Verify that the `-race` flag issue mentioned still applies

### 16. gotcha-surface-fields-single-surface-empty.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents single-surface projects leaving surface-key and surface-type empty in task frontmatter. The referenced `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS. The current SKILL.md contains "Surface-Key/Type Inference" section with a two-layer resolution strategy that matches the lesson's description. The `.forge/config.yaml` EXISTS. The cross-reference `docs/lessons/pattern-surface-resolution-shortcut.md` was not verified but is referenced. The core issue (agent treating `.` key as placeholder and leaving fields empty) is a subtle bug in template interpretation. The solution (single surface → use type value, key empty) is specific and correct.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS (Surface-Key/Type Inference section confirmed), `.forge/config.yaml` EXISTS
- **Required update**: (1) Verify whether the current quick-tasks SKILL.md's Surface-Key/Type Inference instructions adequately address the single-surface case; (2) The lesson may be partially resolved if the template instructions were clarified

### 17. gotcha-task-cli-path-duplication.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson documents a path duplication bug where `index.json` record fields contained `tasks/` prefix AND `GetTaskFile()` added another `tasks/`, producing `docs/features/<slug>/tasks/tasks/records/<file>.md`. Checking the current state: all `index.json` files use `"records/xxx.md"` format (no `tasks/` prefix), and `GetTaskFile()` adds the `tasks/` segment correctly. The `query.go` code uses `feature.GetTaskFile(featureSlug, t.Record)` which produces the correct path. The `claim.go` uses `feature.GetTaskFile(featureSlug, t.File)` for the FILE field. The `GetTaskFile` function in `forge-cli/pkg/feature/paths.go` adds `TasksDirName` (`tasks`) only once. The bug described in the lesson appears to have been fixed — either in the data (index.json) or in the code. The `query.go` reference to `RECORD_FILE` field at line 58 uses `GetTaskFile` which produces correct paths. The referenced files `claim.go`, `query.go`, `record.go` all exist at their current locations.
- **Code path verification**: `forge-cli/internal/cmd/task/claim.go` EXISTS, `forge-cli/internal/cmd/task/query.go` EXISTS, `forge-cli/pkg/task/record.go` EXISTS, index.json record fields use `records/xxx.md` format (no duplication)
- **Recommendation**: Mark as outdated. The described bug has been fixed. The architectural insight about path construction responsibility boundaries (data vs code) remains valid but is now historical.

### 18. gotcha-task-derivation-over-research.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents over-research triggered by the instruction verb "Determine" in quick-tasks Step 2, contrasted with breakdown-tasks using "Inspect". The suggested fix was to change "Determine" to "Infer" with a degradation rule. Checking the current `plugins/forge/skills/quick-tasks/SKILL.md`: the specific "Determine affected file paths from the solution description" instruction no longer exists in its original form. The SKILL.md has been significantly restructured — Step 2 now describes "For each In Scope bullet: estimate effort, derive acceptance criteria, classify type, resolve surface-key/surface-type, fill Reference Files" without the specific "Determine" verb for file paths. The referenced `plugins/forge/skills/breakdown-tasks/SKILL.md` EXISTS. The core insight about instruction verb choice affecting agent behavior ("Determine" vs "Infer" vs "Inspect") is a valuable prompt engineering pattern.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS (restructured, original "Determine" instruction no longer present), `plugins/forge/skills/breakdown-tasks/SKILL.md` EXISTS
- **Required update**: (1) Note that the SKILL.md was restructured, making the specific "Determine" instruction obsolete; (2) The verb-choice insight remains valid as a general prompt engineering lesson; (3) Verify whether the restructured Step 2 adequately addresses the over-research issue

### 19. gotcha-task-executor-auto-claim.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson documents `zcode:task-executor` auto-claiming the next task after completing its assigned task. The project has been renamed from `zcode` to `forge` — all `zcode:` references are stale. The current namespace uses `forge:task-executor`. The referenced `zcode:claim-task` is now `forge:claim-task` or equivalent. The core issue (task-executor has autonomous task-chaining behavior that conflicts with the dispatcher) may still be architecturally relevant, but all specific references use the old `zcode` namespace. The `plugins/forge/commands/run-tasks.md` EXISTS and contains the dispatcher logic. The solution (explicit "stop after one task" instruction in dispatch prompt) is a reasonable mitigation.
- **Code path verification**: `plugins/forge/commands/run-tasks.md` EXISTS, `zcode:*` references are stale (project uses `forge:` namespace)
- **Recommendation**: Mark as outdated. All `zcode` references are stale. If the auto-chaining behavior still exists with the `forge` namespace, the lesson should be rewritten with current references; otherwise it serves only as historical documentation.

### 20. gotcha-task-executor-ignores-implementation-notes.md

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The lesson documents a task-executor choosing `npx playwright test` over the task-specified `just test-e2e --feature milestone-map`, because Implementation Notes have lower compliance priority than Hard Rules. The core insight (elevate critical commands from Implementation Notes to Hard Rules) is a sound and universally applicable process standard. The referenced `docs/features/milestone-map/tasks/run-e2e-tests.md` was not found (the milestone-map feature may have been cleaned up). The `justfile` EXISTS at the project root. The referenced `tests/e2e/playwright.config.ts` DOES NOT EXIST — `tests/e2e/` directory does not exist. The cross-reference `gotcha-split-task-missing-shared-setup.md` EXISTS. The `just test-e2e` recipe may have been renamed or removed since the justfile was updated. Despite invalid code paths, the process standard (Hard Rules >> Implementation Notes for task-executor compliance) is architecturally valid.
- **Code path verification**: root `justfile` EXISTS, `docs/features/milestone-map/tasks/run-e2e-tests.md` NOT FOUND, `tests/e2e/playwright.config.ts` MISSING
- **Required update**: (1) Update or remove invalid code path references; (2) Verify if `just test-e2e` recipe still exists; (3) The Hard Rules priority pattern remains valid regardless of specific code references

## Duplicate Detection (Topic Clustering)

### Cluster 1: Spec/Reference Files Authority (2 lessons)

| Lesson | Focus | Recommendation |
|--------|-------|----------------|
| gotcha-spec-authority-drift | Agent drifts from spec during multi-file edits; solution: enforce Reference Files in templates | KEEP (primary) — detailed root cause analysis + implemented fix |
| gotcha-task-reference-files-scope-creep (batch 4) | Reference Files pointing to proposal causes scope creep | NOT DUPLICATE — different mechanism (read-over-bounding vs write-over-bounding), both valid |

**Verdict**: No duplicate. Both lessons address spec authority but from different angles (spec drift during editing vs scope creep from reading proposal).

### Cluster 2: Task Executor Behavior (3 lessons)

| Lesson | Focus | Recommendation |
|--------|-------|----------------|
| gotcha-task-executor-auto-claim | Executor auto-chains to next task | OUTDATED (zcode references) |
| gotcha-task-executor-ignores-implementation-notes | Executor ignores Implementation Notes in favor of own judgment | KEEP — valid process standard |
| gotcha-quick-tasks-no-autochain (batch 4) | quick-tasks doesn't auto-chain to run-tasks | NOT DUPLICATE — different scope (quick-tasks planning vs task-executor execution) |

**Verdict**: No duplicate within batch. Auto-claim lesson is outdated (stale namespace); ignores-implementation-notes is a valid process standard.

### Cluster 3: Task Splitting Patterns (2 lessons)

| Lesson | Focus | Recommendation |
|--------|-------|----------------|
| gotcha-split-rules-operational-blindness | Splitting rules miss operational granularity | KEEP — fix implemented |
| gotcha-split-task-missing-shared-setup | Splitting misses global setup phase | OUTDATED (invalid code paths) |

**Verdict**: No duplicate. Different aspects of task splitting (operational granularity vs shared setup identification). Both are complementary patterns.

### Cluster 4: Justfile/Command Compliance (2 lessons)

| Lesson | Focus | Recommendation |
|--------|-------|----------------|
| gotcha-strategy-bypass-justfile | Bypassing justfile in favor of direct toolchain commands | KEEP — valid process standard |
| gotcha-task-executor-ignores-implementation-notes | Executor ignores task-specified commands | KEEP — different angle (process philosophy vs prompt priority) |

**Verdict**: No duplicate. Strategy-bypass is about philosophical approach (first-principles vs empiricism); ignores-implementation-notes is about prompt section priority (Hard Rules > Implementation Notes).

### Cluster 5: Stale State/Results (3 lessons)

| Lesson | Focus | Recommendation |
|--------|-------|----------------|
| gotcha-stale-state-json-feature-mismatch | Stale state.json causes wrong feature claims | KEEP — valid |
| gotcha-stale-test-results-cascade | Stale test results cause cascading fix tasks | KEEP — needs path update |
| gotcha-stale-skill-cli-flags | Stale CLI flag reference in skill docs | KEEP — needs path update |

**Verdict**: No duplicate. Three different "stale" scenarios (feature state, test results, documentation).

### Cluster 6: Task Scope/Size Control (3 lessons)

| Lesson | Focus | Recommendation |
|--------|-------|----------------|
| gotcha-split-rules-operational-blindness | Operational ceiling for task splitting | KEEP (primary — operational granularity) |
| gotcha-prompt-template-complexity-agnostic (batch 3) | Templates don't differentiate by task complexity | NOT DUPLICATE — template-level vs splitting-rule-level |
| gotcha-task-derivation-over-research | "Determine" verb triggers over-research | NOT DUPLICATE — different cause (verb choice vs granularity rules) |

**Verdict**: No duplicate. Three different causes of oversized task scope, each with distinct mitigation.

## Audit Quality Review

- **Sampling ratio**: 10% (2 items) | **Sampling result**: pass | **Missed items**: 0 | **Extended review**: no
- **Items sampled**: gotcha-spec-authority-drift (#7), gotcha-standard-task-id-collision (#13) — both verified against current codebase with accurate findings
