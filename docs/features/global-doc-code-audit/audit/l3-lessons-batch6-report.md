# L3 Lessons Audit Report — Batch 6 (gotcha-task-executor-invisible-thinking-time to lesson-tui-visual-verify)

## Audit Baseline

- **Baseline commit**: d1e0656a (global-doc-code-audit branch)
- **Audit date**: 2026-06-03
- **Audit scope**: 20 lesson files from `docs/lessons/` (alphabetically: gotcha-task-executor-invisible-thinking-time.md through lesson-tui-visual-verify.md)

## Classification Distribution

| Classification | Count |
|----------------|-------|
| code-reference | 11 |
| process-standard | 5 |
| experience-summary | 4 |

## Status Summary

| Status | Count |
|--------|-------|
| valid | 4 |
| needs-update | 11 |
| outdated | 5 |
| duplicate | 0 |

Note: Cluster 4 identifies gotcha-task-type-documentation-vs-doc.md (marked outdated) as functionally overlapping with gotcha-task-type-for-md-files.md (marked needs-update). The overlap is documented in the Duplicate Detection section but neither is formally marked duplicate since item 10 describes a fixed template bug while item 11 describes an ongoing classification issue.

## Cross-Layer Influence Check

The following L1/L2 audit findings affect items in this batch:

| L1/L2 Finding | Affected Lesson | Impact |
|---------------|-----------------|--------|
| L2 conventions-batch1: `quality_gate.go` moved to `qualitygate/` subpackage | gotcha-task-executor-thinking-overhead | Referenced path `forge-cli/internal/cmd/quality_gate.go` is outdated; actual path is `forge-cli/internal/cmd/qualitygate/quality_gate.go` |
| L2 conventions-batch2: `tests/e2e/` directory does not exist | gotcha-test-script-staging-vs-graduation | All references to `tests/e2e/` paths are invalid; test infrastructure reorganized to journey-based `tests/<journey>/` |
| L2 conventions-batch1: `prompt/data/` renamed to `prompt/templates/` | gotcha-task-reference-files-scope-creep | Referenced path `forge-cli/pkg/prompt/data/coding-enhancement.md` should be `forge-cli/pkg/prompt/templates/coding-enhancement.md` |
| L2 conventions-batch2: surfaces config replaced interfaces | gotcha-test-pipeline-no-languages | Lesson references `interfaces` config field; current config uses `surfaces` field with `ReadSurfaces()` instead of `ReadInterfaces()` |
| L2 conventions-batch1: `run-tasks` is a command, not a skill | gotcha-task-executor-never-returns | Reference to `plugins/zcode/agents/task-executor.md` is incorrect; actual location is `plugins/forge/agents/task-executor.md` |
| L2 conventions-batch2: TUI conventions not in docs/conventions/ | lesson-tui-tech-design-mockup, lesson-tui-visual-verify | Referenced `docs/conventions/tui-layout-ui.md` and `docs/conventions/tui-dynamic-content.md` do not exist |

## Item-by-Item Analysis

### 1. gotcha-task-executor-invisible-thinking-time.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents a forensic tool blind spot where extended thinking time is invisible in JSONL entries. The core insight (forensic duration calculation misses API-level thinking time) remains valid. The referenced file `forge-cli/internal/cmd/forensic.go` does NOT exist as a single file — the forensic implementation is split across `forge-cli/internal/cmd/forensic/` subpackage (extract.go, helpers.go, search.go, subagents.go, types.go, commands.go, register.go). The `forensic extract` duration calculation may have been updated since the lesson was written. The cross-referenced lesson `gotcha-task-executor-thinking-overhead` EXISTS and is in this batch. The `gotcha-prompt-template-complexity-agnostic` lesson was NOT found in this batch's range.
- **Code path verification**: `forge-cli/internal/cmd/forensic.go` MISSING (restructured to `forge-cli/internal/cmd/forensic/` directory with multiple files)
- **Required update**: (1) Update file reference from `forge-cli/internal/cmd/forensic.go` to `forge-cli/internal/cmd/forensic/` package; (2) Verify if the duration calculation blind spot has been addressed in the restructured code

### 2. gotcha-task-executor-never-returns.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes task-executor continuing to claim tasks after completing one, due to missing termination instruction. The fix has been FULLY IMPLEMENTED: `plugins/forge/agents/task-executor.md` now contains `<EXTREMELY-IMPORTANT>` block with "ONE TASK PER INVOCATION — after completing, STOP immediately" and "FORBIDDEN: run forge task claim, read index.json, or start any subsequent task." The root cause (missing termination constraint) was fixed by adding explicit stop directives. The referenced path `plugins/zcode/agents/task-executor.md` is incorrect — the actual path is `plugins/forge/agents/task-executor.md`.
- **Code path verification**: `plugins/forge/agents/task-executor.md` EXISTS with full termination directives; `plugins/zcode/agents/task-executor.md` MISSING
- **Recommendation**: Mark as outdated — the described bug has been fully fixed. The general insight (agent definitions need explicit termination conditions) could be preserved as a generalized experience-summary, but the specific problem and solution are resolved.

### 3. gotcha-task-executor-proposal-inaccessible.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents task-executor failing to read a proposal due to invalid Reference Files paths, causing implementation to deviate from design. The core architectural insight ("AC is a verification checklist, not a spec; executor needs access to the full proposal") remains universally valid. The referenced `docs/proposals/pipeline-topology-registry/proposal.md` EXISTS. The referenced `forge-cli/pkg/task/pipeline.go` EXISTS and now contains the `PipelineRegistry` implementation with `PipelineNode` structs, confirming the feature was implemented. The specific task file `docs/features/pipeline-topology-registry/tasks/1-define-pipeline-registry.md` EXISTS. The execution record at `docs/features/pipeline-topology-registry/tasks/records/1-define-pipeline-registry.md` EXISTS.
- **Code path verification**: All referenced files EXIST: `docs/proposals/pipeline-topology-registry/proposal.md`, `forge-cli/pkg/task/pipeline.go`, task file, and record
- **Required update**: (1) The lesson's insight remains valid — note that the quick-tasks SKILL.md now has a `<HARD-RULE>` requiring the first Reference Files entry to be the full proposal path; (2) The lesson could reference this fix as evidence that the root cause was partially addressed at the skill level

### 4. gotcha-task-executor-redundant-search.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes task-executor performing redundant code searches even when task files contain explicit target file paths, and recommends adding "skip search" constraints to Hard Rules. This is a general process pattern that does not depend on specific code paths. The solution (adding Hard Rules to skip search when targets are explicit) is a reusable pattern applicable to any task-executor invocation. No code references to verify.
- **Code path verification**: N/A — process-standard with no specific code references
- **Recommendation**: Keep as valid. The pattern (redundant exploration when information is already complete) is universally applicable to task-executor usage.

### 5. gotcha-task-executor-stops-at-step1.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson documents task-executor stopping after Step 1 output without doing actual work, and provides escalation protocols (direct implementation or concrete-spec agent). This is a behavioral pattern that remains applicable — LLM agents can still truncate execution. The escalation path (Option A: direct implementation, Option B: concrete spec) is sound advice. No specific code paths are referenced that could become outdated. The observation about tool count dropping across retries (19->14->4) is a general LLM behavior pattern.
- **Code path verification**: N/A — process-standard describing general agent behavior
- **Recommendation**: Keep as valid. The escalation protocol remains applicable.

### 6. gotcha-task-executor-thinking-overhead.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents 87% thinking overhead in task-executor due to redundant searches, coarse task granularity, and quality gate false positives. The core insights remain valid: task decomposition, batch-search strategy, quality gate baseline comparison. However, specific code references need updating: (1) `docs/forensics/task1-duration/evidence/evidence.json` EXISTS; (2) `forge-cli/internal/cmd/quality_gate.go` is OUTDATED — the actual path is `forge-cli/internal/cmd/qualitygate/quality_gate.go`; (3) The referenced `resolveBreakdownDeps` function in `autogen.go:253-289` no longer exists — the test pipeline was completely restructured to use `PipelineRegistry` in `pipeline.go` with `PipelineNode` structs.
- **Code path verification**: `docs/forensics/task1-duration/evidence/evidence.json` EXISTS; `forge-cli/internal/cmd/quality_gate.go` MISSING (moved to `qualitygate/`); `resolveBreakdownDeps` function NOT FOUND in autogen.go
- **Required update**: (1) Update quality gate path to `forge-cli/internal/cmd/qualitygate/quality_gate.go`; (2) Remove or update the `resolveBreakdownDeps` reference since the test pipeline architecture was fundamentally restructured

### 7. gotcha-task-index-preserve-deps.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents `forge task index` preserving old dependencies in index.json when .md files are updated, causing claim to skip business tasks. The core problem (index.json as runtime state vs .md as edit source) is architecturally valid. The referenced `forge-cli/internal/cmd/task.go` is OUTDATED — task commands are now in `forge-cli/internal/cmd/task/` subpackage (index.go, claim.go, etc.). The specific file `docs/features/slim-task-prompt-templates/tasks/index.json` EXISTS. The proposed fix (distinguish runtime vs declarative fields during index) is a design recommendation.
- **Code path verification**: `forge-cli/internal/cmd/task.go` MISSING (restructured to `forge-cli/internal/cmd/task/` directory); `docs/features/slim-task-prompt-templates/tasks/index.json` EXISTS
- **Required update**: (1) Update task.go reference to `forge-cli/internal/cmd/task/` subpackage; (2) Verify if the index preserve behavior has been fixed in the restructured command

### 8. gotcha-task-reference-files-scope-creep.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents Reference Files pointing to proposal full text causing task-executor to fix issues outside its scope. The core insight (Reference Files as scope boundary leak point) remains valid and well-documented. The referenced `forge-cli/pkg/prompt/data/coding-enhancement.md` path is OUTDATED — should be `forge-cli/pkg/prompt/templates/coding-enhancement.md`. The cross-referenced lessons `gotcha-prompt-template-complexity-agnostic`, `gotcha-task-executor-invisible-thinking-time`, and `gotcha-quick-tasks-merge-threshold` are from different batches.
- **Code path verification**: `forge-cli/pkg/prompt/data/coding-enhancement.md` MISSING (correct path: `forge-cli/pkg/prompt/templates/coding-enhancement.md`)
- **Required update**: Update the prompt data path from `data/` to `templates/`

### 9. gotcha-task-reference-source-anchor-misread.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson documents quick-tasks generating fabricated source annotations (wrong paths, invented section names, ambiguous format). The core insight (LLMs hallucinate reference metadata) remains critically important. The fix has been PARTIALLY IMPLEMENTED: `plugins/forge/skills/quick-tasks/SKILL.md` now contains a `<HARD-RULE>` block that (1) requires grepping actual headers before writing Reference Files, (2) mandates the first entry be the full proposal path, (3) prohibits fabricating headers. However, the old `(source: ...)` format may still appear in historical task files. The referenced `plugins/forge/skills/quick-tasks/` EXISTS, `plugins/forge/agents/task-executor/` EXISTS (as `task-executor.md` file, not directory), and the pipeline-topology-registry task files all EXIST.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS with HARD-RULE fix; `plugins/forge/agents/task-executor/` MISSING (actual: `plugins/forge/agents/task-executor.md` — a file, not directory)
- **Required update**: (1) Note that the quick-tasks SKILL.md now has HARD-RULE preventing header fabrication; (2) Update agent path reference from `task-executor/` directory to `task-executor.md` file

### 10. gotcha-task-type-documentation-vs-doc.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes `type: "documentation"` being rejected by CLI validation (which expects `"doc"`). The template fix has been FULLY IMPLEMENTED: `plugins/forge/skills/quick-tasks/templates/task-doc.md` now uses `type: "doc"` as its default. The SKILL.md template selection logic guides agents to the correct template. The CLI's `ValidTypes` in `types.go` confirms `TypeDoc = "doc"` is the correct constant. The Windows cache path `C:\Users\panda\.claude\plugins\cache\forge\forge\3.0.0-rc.5\skills\quick-tasks\templates\task-doc.md` is an invalid local reference.
- **Code path verification**: `task-doc.md` template now has `type: "doc"` (FIXED); `types.go` has `TypeDoc = "doc"` (unchanged); Windows cache path is invalid
- **Recommendation**: Mark as outdated — the template was fixed and the type mismatch can no longer occur when using the correct template. The general principle (cross-reference template defaults against CLI validation) is covered by other lessons.

### 11. gotcha-task-type-for-md-files.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes pure .md file tasks being incorrectly typed as `coding.enhancement`, causing unnecessary quality gate execution. The core insight (use operation semantics, not file location, to determine task type) remains valid. The referenced code files all EXIST: `forge-cli/pkg/task/types.go` (Type constants), `forge-cli/pkg/task/category.go` (CategoryForType), `forge-cli/pkg/task/build.go` (IsTestableType at line 497). The specific `slim-task-prompt-templates` feature index EXISTS. However, the template system now has separate `task.md` (with `type: "{{TYPE}}"` and guidance for all types) and `task-doc.md` (hardcoded `type: "doc"`) templates, which partially addresses the root cause by providing a template that defaults to doc type.
- **Code path verification**: All referenced files EXIST: types.go, category.go, build.go, index.json
- **Required update**: (1) Note that the dual-template approach (task.md vs task-doc.md) partially addresses the issue; (2) The CLI `BuildIndex` warning for `coding.*` tasks with only .md files has NOT been implemented — this remains an open improvement

### 12. gotcha-test-chain-not-linked-to-last-business-gate.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes auto-generated test tasks not being linked to the last business gate, causing execution order errors. The entire test pipeline architecture has been FUNDAMENTALLY RESTRUCTURED. The `resolveBreakdownDeps` function referenced at `autogen.go:253-289` NO LONGER EXISTS. The test pipeline now uses a `PipelineRegistry` approach in `forge-cli/pkg/task/pipeline.go` where each `PipelineNode` has explicit `DependsOn` entries using `DepRef` with resolver functions like `ResolveHighestGateOrLastBiz` and `ResolveLastBusinessTask`. The `T-clean-code` node correctly depends on `ResolveHighestGateOrLastBiz`. The `T-test-gen-journeys` node depends on both `T-review-doc` and `T-clean-code` if generated. The specific recurrence with `T-clean-code` described in the 2026-05-28 update is addressed by the new architecture. The referenced `validate_index.go` path is OUTDATED — validation is now at `forge-cli/internal/cmd/task/validate.go`.
- **Code path verification**: `resolveBreakdownDeps` MISSING; `autogen.go:253-289` range no longer contains this function; `forge-cli/pkg/task/pipeline.go` EXISTS with `PipelineRegistry`; `forge-cli/internal/cmd/validate_index.go` MISSING (actual: `forge-cli/internal/cmd/task/validate.go`); `docs/features/unify-enum-constants/tasks/clean-code.md` EXISTS
- **Recommendation**: Mark as outdated — the entire problem domain was resolved by the PipelineRegistry architectural redesign. The general principle (auto-generated tasks must link to business chain) is now enforced structurally by the `ResolveHighestGateOrLastBiz` resolver pattern.

### 13. gotcha-test-pipeline-no-languages.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes `forge task index` not generating test pipeline tasks because `interfaces` was not configured in `.forge/config.yaml`. The system has been FUNDAMENTALLY RESTRUCTURED: (1) The `interfaces` config field has been replaced by `surfaces` — the current `.forge/config.yaml` has `surfaces: cli` instead of `interfaces`; (2) `ReadInterfaces()` NO LONGER EXISTS in `forge-cli/pkg/forgeconfig/detect.go` — replaced by `ReadSurfaces()` (line 23); (3) `DetectLanguages`, `ReadLanguages`, `UnionLanguageInterfaces` have ALL been removed from the codebase; (4) The `GetBreakdownTestTasks` and `GetQuickTestTasks` bridge functions still exist but are marked as "deprecated" and now call `GenerateTestTasks` with surfaces parameter; (5) The guard `if len(interfaces) == 0 { return nil }` has been replaced by `if len(surfaces) == 0` at line 152 of autogen.go.
- **Code path verification**: `ReadInterfaces` MISSING from detect.go; `DetectLanguages`, `ReadLanguages`, `UnionLanguageInterfaces` all MISSING; `ReadSurfaces()` EXISTS at detect.go:23; config.yaml uses `surfaces: cli`; deprecated bridge functions still call GenerateTestTasks
- **Recommendation**: Mark as outdated — the entire terminology and mechanism was replaced (interfaces -> surfaces, language-based -> surface-type-based). The lesson's historical context (section "Historical Context (已解决)") already acknowledges the restructuring, but the main content still references the old system. The lesson could be kept as a historical record of the migration, but all actionable paths are invalid.

### 14. gotcha-test-script-staging-vs-graduation.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes confusion between staging location (`docs/features/<slug>/testing/scripts/`) and regression suite location (`tests/e2e/<slug>/`). The `tests/e2e/` directory NO LONGER EXISTS — test infrastructure has been reorganized to journey-based `tests/<journey>/` directories with Go modules. The "staging vs graduation" two-stage design has been replaced by a new pipeline where: (1) test journeys are generated via `T-test-gen-journeys` node; (2) test scripts are generated per surface type via `T-test-gen-scripts-{surface-type}`; (3) tests run per surface key via `T-test-run-{surface-key}`. The concept of "graduating" scripts from staging to regression suite no longer applies in the current architecture. The `gen-test-scripts` skill EXISTS but operates under the new pipeline model.
- **Code path verification**: `tests/e2e/` MISSING; test directories use journey-based naming; `plugins/forge/skills/gen-test-scripts/` EXISTS
- **Recommendation**: Mark as outdated — the staging/graduation model and `tests/e2e/` paths have been completely replaced by the journey-based test infrastructure.

### 15. hook-stop-e2e-blocking.md

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The lesson describes Stop hook blocking conversations when e2e tests fail after all tasks complete. The Stop hook mechanism EXISTS in `plugins/forge/hooks/hooks.json` — it runs `forge quality-gate`. The `quality-gate` command checks `.forge/state.json` for `allCompleted` flag (confirmed at `forge-cli/internal/cmd/qualitygate/quality_gate.go` lines 82-88). The core issue (Stop hook not distinguishing task execution from other work) remains architecturally valid. However, the specific command flow described (`task all-completed` → `just test-e2e`) is outdated — the current implementation uses `forge quality-gate` which runs a more sophisticated lifecycle including fix-task creation. The lesson references `task all-completed` and `task feature` commands which may have been restructured.
- **Code path verification**: `plugins/forge/hooks/hooks.json` EXISTS (Stop hook: `forge quality-gate`); `forge-cli/internal/cmd/qualitygate/quality_gate.go` EXISTS with allCompleted check; `forge-cli/internal/cmd/task.go` MISSING (restructured to task/ subpackage)
- **Required update**: (1) Update the hook command reference from `task all-completed` to `forge quality-gate`; (2) Note that the quality-gate now has a fix-task lifecycle, not just retry-and-block; (3) The troubleshooting path remains valid (check hooks.json, run the hook command manually, fix failing tests)

### 16. lesson-forge-tui-pipeline-gap.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: This is a comprehensive meta-lesson analyzing why the Forge pipeline fails to prevent vibe coding for TUI features. The systematic diagnosis (9 proposed skill modifications, 0 implemented) remains accurate. The proposed improvement plan (P0/P1/P2 priorities) is well-structured. However: (1) Referenced skill paths use old structure — `skills/tech-design/SKILL.md`, `skills/eval-design/templates/rubric.md`, `skills/breakdown-tasks/SKILL.md` — these should be under `plugins/forge/skills/`; (2) The lesson references `lesson-vibe-coding-scope-control.md` which EXISTS; (3) The Pipeline Stage table references pipeline stages (`/brainstorm`, `/write-prd`, `/tech-design`, `/eval-design`, `/breakdown-tasks`, `/execute-task`) — these skill names remain current. The core conclusion (conventions exist but pipeline doesn't enforce them) is a meta-observation that remains valid regardless of code changes.
- **Code path verification**: `lesson-vibe-coding-scope-control.md` EXISTS; skill paths should be prefixed with `plugins/forge/`; referenced TUI convention files NOT in `docs/conventions/`
- **Required update**: (1) Update skill path references to use `plugins/forge/skills/` prefix; (2) Verify which of the 9 proposed modifications have been implemented since the lesson was written; (3) Note that TUI conventions may have been relocated from `docs/conventions/` to plugin-level rules files

### 17. lesson-gate-force-over-fix.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson documents using `--force` to bypass a gate task that caught a real gap, instead of fixing the gap inline. The principle (gate rejection is signal, not obstacle; fix or block, never force) is a universal process standard. No specific code paths to verify — the lesson describes a behavioral pattern and decision framework. The `--force` flag concept and gate task structure (`.gate` suffix) still exist in the codebase (confirmed: `IDSuffixGate = ".gate"` in types.go).
- **Code path verification**: N/A — process-standard with no specific code references
- **Recommendation**: Keep as valid. The decision framework is timeless and widely applicable.

### 18. lesson-guide-missing-proposals-dir.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson documents guide.md missing the `docs/proposals/` directory entry, causing agents to write proposals to wrong locations. The fix has been APPLIED — the lesson states "已在 guide.md 中补充 proposals/ 条目" and `plugins/forge/hooks/guide.md` EXISTS at the correct location. The principle (guide.md must catalog all convention-governed directories) remains valid as a process standard. The lesson also notes the broader pattern: when creating new directories with conventions, guide.md must be updated.
- **Code path verification**: `plugins/forge/hooks/guide.md` EXISTS
- **Recommendation**: Keep as valid. The process pattern (synchronizing guide.md with directory conventions) remains applicable for any future directory additions.

### 19. lesson-tui-tech-design-mockup.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: The lesson specifies ASCII layout mockup requirements for TUI features in tech-design. The detailed specification (panel layout, dimensions, boundary scenarios, character palette, color mapping) is comprehensive and well-structured. However: (1) It references `tui-layout-ui.md` color table at `docs/conventions/tui-layout-ui.md` which DOES NOT EXIST — TUI conventions are now in plugin-level files (e.g., `plugins/forge/skills/ui-design/rules/tui-panel-requirements.md`, `plugins/forge/skills/run-tests/rules/surfaces/tui.md`); (2) The lesson proposes adding mockup requirements to `/tech-design` SKILL.md — this may or may not have been implemented. The boundary scenario table (5 required scenarios) and character palette specification are independently valuable regardless of code state.
- **Code path verification**: `docs/conventions/tui-layout-ui.md` MISSING; TUI convention files exist at `plugins/forge/skills/*/rules/surfaces/tui.md` and `plugins/forge/skills/ui-design/rules/tui-panel-requirements.md`
- **Required update**: (1) Update the color palette reference from `docs/conventions/tui-layout-ui.md` to actual TUI convention locations under `plugins/forge/skills/`; (2) Verify if tech-design SKILL.md has been updated to include TUI mockup gate

### 20. lesson-tui-visual-verify.md

- **Classification**: experience-summary
- **Status**: needs-update
- **Justification**: The lesson specifies visual verification criteria for TUI tasks (golden file comparison, boundary conditions, test data realism). The verification template and the "3 types of visual checks" framework are well-designed. However: (1) It references `docs/conventions/tui-layout-ui.md` which DOES NOT EXIST — same issue as lesson-tui-tech-design-mockup; (2) The golden test code examples reference `lipgloss.Width()` and `internal/model/*.go` View() functions — these patterns remain valid for the TUI codebase; (3) The breakdown-tasks integration point (auto-appending TUI verify template) may or may not have been implemented. The 4 mandatory boundary values table (CJK, long path, multi-digit, empty field) is a universally applicable testing pattern.
- **Code path verification**: `docs/conventions/tui-layout-ui.md` MISSING; TUI-related skill files EXIST under `plugins/forge/skills/`
- **Required update**: (1) Update color palette reference from `docs/conventions/tui-layout-ui.md` to actual plugin-level TUI convention files; (2) Verify if breakdown-tasks skill has been updated with TUI verify template injection

## Duplicate Detection

### Cluster 1: Task-executor Scope Control and Search Efficiency

| Item | Status | Classification |
|------|--------|----------------|
| gotcha-task-executor-redundant-search.md | valid | process-standard |
| gotcha-task-reference-files-scope-creep.md | needs-update | code-reference |
| gotcha-task-executor-invisible-thinking-time.md | needs-update | code-reference |
| gotcha-task-executor-thinking-overhead.md | needs-update | code-reference |

**Verdict**: NOT duplicates. Item 1 addresses redundant search behavior. Item 2 addresses Reference Files causing scope leakage. Items 3-4 both address thinking overhead but from different angles (item 3: forensic tool blind spot for wall-clock time; item 4: 87% thinking ratio with quality gate false positives). Items 3 and 4 share the theme of thinking overhead but document different incidents with different root causes and solutions — item 3 focuses on forensic measurement, item 4 on task decomposition and quality gate.

### Cluster 2: Task-executor Termination and Execution Behavior

| Item | Status | Classification |
|------|--------|----------------|
| gotcha-task-executor-never-returns.md | outdated | code-reference |
| gotcha-task-executor-stops-at-step1.md | valid | process-standard |

**Verdict**: NOT duplicate. Item 1 describes executor continuing past one task (over-execution). Item 2 describes executor stopping before completing one task (under-execution). They are opposite failure modes of the same agent.

### Cluster 3: Proposal/Reference File Accessibility

| Item | Status | Classification |
|------|--------|----------------|
| gotcha-task-executor-proposal-inaccessible.md | needs-update | code-reference |
| gotcha-task-reference-source-anchor-misread.md | needs-update | code-reference |
| gotcha-task-reference-files-scope-creep.md | needs-update | code-reference |

**Verdict**: Partial overlap but NOT duplicates. Item 1 documents executor unable to read proposal (missing path). Item 2 documents fabricated source annotations (hallucinated headers). Item 3 documents scope leakage from reading too much of the proposal. All three relate to Reference Files but describe distinct problems: can't-read, wrong-metadata, and over-read. All three are independently actionable.

### Cluster 4: Task Type Classification

| Item | Status | Classification |
|------|--------|----------------|
| gotcha-task-type-documentation-vs-doc.md | outdated | code-reference |
| gotcha-task-type-for-md-files.md | needs-update | code-reference |

**Verdict**: DUPLICATE. Both lessons describe the same root problem: tasks involving only .md files being incorrectly typed as `coding.*` instead of `doc`. Item 10 (`documentation-vs-doc`) focuses on the template defaulting to `type: "documentation"` (now fixed). Item 11 (`for-md-files`) focuses on the broader classification rule (operation semantics over file location). Item 11 is more comprehensive and subsumes item 10. **Recommendation**: Keep item 11 (gotcha-task-type-for-md-files.md) as the primary reference. Mark item 10 as duplicate of item 11.

### Cluster 5: TUI Feature Pipeline Gaps

| Item | Status | Classification |
|------|--------|----------------|
| lesson-forge-tui-pipeline-gap.md | needs-update | experience-summary |
| lesson-tui-tech-design-mockup.md | needs-update | experience-summary |
| lesson-tui-visual-verify.md | needs-update | experience-summary |

**Verdict**: NOT duplicates. Item 16 is the meta-lesson (systematic diagnosis of pipeline gaps). Item 19 is the specific tech-design mockup specification. Item 20 is the specific visual verification specification. They form a parent-child relationship: item 16 references items 19 and 20 as sub-lessons. All three are independently valuable and serve different purposes.

### Cluster 6: Test Pipeline Terminology and Infrastructure

| Item | Status | Classification |
|------|--------|----------------|
| gotcha-test-chain-not-linked-to-last-business-gate.md | outdated | code-reference |
| gotcha-test-pipeline-no-languages.md | outdated | code-reference |
| gotcha-test-script-staging-vs-graduation.md | outdated | code-reference |

**Verdict**: NOT duplicates. Item 12 addresses dependency wiring (test chain not linked to business gate). Item 13 addresses task generation (no test tasks generated due to missing config). Item 14 addresses location confusion (staging vs graduation paths). All three are outdated because the entire test pipeline was restructured to use PipelineRegistry, but they describe different aspects of the old system.

### Cluster 7: Quality Gate Interactions

| Item | Status | Classification |
|------|--------|----------------|
| hook-stop-e2e-blocking.md | needs-update | process-standard |
| lesson-gate-force-over-fix.md | valid | process-standard |
| gotcha-task-executor-thinking-overhead.md | needs-update | code-reference |

**Verdict**: NOT duplicates. Item 15 addresses Stop hook blocking on e2e failures. Item 17 addresses bypassing gate validation with --force. Item 6 addresses quality gate false positives causing submit retries. All three relate to quality gates but describe completely different interaction patterns.

## Audit Quality Review

- **Sampling ratio**: 10% | **Sample**: items 2 (gotcha-task-executor-never-returns), 17 (lesson-gate-force-over-fix) | **Result**: PASS | **Missed items**: 0 | **Extended review**: No
- **Cross-batch duplicate check**: No duplicates found with items in batches 1-5. The task-type classification issue (Cluster 4) is internal to this batch.
- **Path verification coverage**: 100% of code-referenced paths verified via find/grep
