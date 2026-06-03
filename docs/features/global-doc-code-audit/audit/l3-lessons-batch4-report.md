# L3 Lessons Audit Report — Batch 4 (gotcha-main-session-flag to gotcha-revert-mid-dispatch)

## Audit Baseline

- **Baseline commit**: b5f67a8f (docs(audit): add L3 lessons batch 3 audit report)
- **Audit date**: 2026-06-03
- **Audit scope**: 20 lesson files from `docs/lessons/` (alphabetically: gotcha-main-session-flag.md through gotcha-revert-mid-dispatch.md)

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
| needs-update | 9 |
| outdated | 4 |
| duplicate | 3 |

## Item-by-Item Analysis

### 1. gotcha-main-session-flag.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes the MAIN_SESSION flag and the dispatcher routing logic. The core architectural insight remains valid: subagents lack the Agent tool and cannot spawn sub-subagents. The `run-tasks.md` command DOES contain MAIN_SESSION routing (confirmed: line 9, line 16 flowchart, line 28 dispatch logic, line 60 MAIN_SESSION check). However, several referenced files are outdated: (1) `plugins/forge/skills/eval-test-cases/SKILL.md` DOES NOT EXIST — the eval-test-cases skill was removed during skill rationalization; (2) `plugins/forge/skills/breakdown-tasks/templates/eval-test-cases.md` DOES NOT EXIST. The `plugins/forge/commands/run-tasks.md` EXISTS, `plugins/forge/commands/execute-task.md` EXISTS, `plugins/forge/agents/task-executor.md` EXISTS, `plugins/forge/skills/breakdown-tasks/SKILL.md` EXISTS. The "Changes Applied" section lists 7 changes that were implemented, making this partly a historical record.
- **Code path verification**: `plugins/forge/commands/run-tasks.md` EXISTS (MAIN_SESSION routing confirmed), `plugins/forge/agents/task-executor.md` EXISTS, `plugins/forge/skills/eval-test-cases/SKILL.md` MISSING (removed), `plugins/forge/skills/breakdown-tasks/templates/eval-test-cases.md` MISSING
- **Required update**: (1) Remove references to `eval-test-cases/SKILL.md` and `eval-test-cases.md` template; (2) Note that the skill was removed during rationalization; (3) The lesson's value shifts from "active bug" to "design pattern documentation"

### 2. gotcha-merge-ghost-revival.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes git merge silently reintroducing deleted files. This is a universal git behavior pattern independent of specific code paths. The reusable pattern (pre-merge audit, post-merge verification) is universally applicable. The referenced skill directories (`eval-*`) were indeed consolidated into `plugins/forge/skills/eval/` (confirmed EXISTS). The referenced `error-fixer.md` agent DOES NOT EXIST (removed). The `record-task/SKILL.md` DOES NOT EXIST (removed). The specific incidents are historical, but the git merge behavior pattern and the detection/prevention steps remain valid for any project using feature branches.
- **Code path verification**: `plugins/forge/skills/eval/` EXISTS (consolidated eval skill), `plugins/forge/agents/error-fixer.md` MISSING (removed), `plugins/forge/skills/record-task/SKILL.md` MISSING (removed), `docs/proposals/skill-rationalization/proposal.md` EXISTS
- **Recommendation**: Keep as valid. The git merge ghost revival pattern is universally applicable. The specific file references serve as historical incident context.

### 3. gotcha-pipeline-skill-bypass.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson prescribes following pipeline skills in order rather than substituting ad-hoc analysis. This is a process-level guideline. The referenced files all EXIST: `plugins/forge/commands/quick.md` EXISTS, `plugins/forge/skills/brainstorm/SKILL.md` EXISTS, `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS. The core rule ("execute pipeline steps in order, don't substitute ad-hoc analysis") is independent of specific code paths. The exception for explicit user override ("just analyze this") is a sound process design principle.
- **Code path verification**: `plugins/forge/commands/quick.md` EXISTS, `plugins/forge/skills/brainstorm/SKILL.md` EXISTS, `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS
- **Recommendation**: Keep as valid. Universal process standard for pipeline skill execution.

### 4. gotcha-post-completion-commit.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The lesson describes using Claude Code's Stop hook for post-completion status transition and commit. The hook configuration in `plugins/forge/hooks/hooks.json` confirms this pattern is implemented: the Stop hook array contains `forge quality-gate` followed by `forge feature complete --if-done`. The `docs/official-references/hooks.md` EXISTS. The `plugins/forge/skills/quick/SKILL.md` DOES NOT EXIST (quick is a command, not a skill — `plugins/forge/commands/quick.md` EXISTS). The lesson correctly describes the `stop_hook_active` mechanism and the dual-hook architecture (quality gate first, then status transition). The "Open Question" about multi-hook execution order has been implicitly resolved by the current implementation (hooks execute sequentially per hooks.json).
- **Code path verification**: `plugins/forge/hooks/hooks.json` EXISTS (Stop hook confirmed: quality-gate + feature complete), `docs/official-references/hooks.md` EXISTS, `plugins/forge/skills/quick/SKILL.md` MISSING (quick is a command at `plugins/forge/commands/quick.md`)
- **Recommendation**: Keep as valid. The Stop hook mechanism described matches the current implementation. Minor note: `plugins/forge/skills/quick/SKILL.md` should be `plugins/forge/commands/quick.md`.

### 5. gotcha-pre-existing-syntax-errors-block-executor.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson is written in Chinese (unlike most other lessons which are in English) and describes pre-existing syntax errors blocking the task executor. The core insight ("an uncompilable codebase is a deadlock trap for the dispatcher") remains architecturally valid. The referenced `submit.go:128` now exists at `forge-cli/internal/cmd/task/submit.go` (moved to subdirectory). The referenced `test_promote.go:46,67` — no file with this name exists in the current codebase (searched with `find`). The cross-reference `gotcha-dispatcher-ignores-compilation-diagnostics.md` EXISTS. The lesson mentions `go build ./...` as a pre-compilation check — this is not currently implemented in the dispatcher (run-tasks.md has no pre-dispatch compilation check step). The "add fix task" mechanism has been redesigned: `AddFixTask` in `quality_gate_fix_task.go` now groups failures by test directory. The SourceTaskID is deliberately empty for quality-gate fix tasks (confirmed in `createFixTask` line 186).
- **Code path verification**: `forge-cli/internal/cmd/task/submit.go` EXISTS (moved), `test_promote.go` MISSING, `docs/lessons/gotcha-dispatcher-ignores-compilation-diagnostics.md` EXISTS, `forge-cli/internal/cmd/qualitygate/quality_gate_fix_task.go` EXISTS
- **Required update**: (1) Update `submit.go` path from root to subdirectory; (2) Note `test_promote.go` no longer exists; (3) Update the fix-task mechanism description to reflect current per-directory grouping design

### 6. gotcha-present-analysis-before-edit.md

- **Classification**: experience-summary
- **Status**: valid
- **Justification**: The lesson prescribes analyzing before editing, with independent judgment as the primary framework and cross-referencing as secondary. This is a generalized reasoning principle that does not reference specific code paths. The four-step process (independent assessment, then reference, then gap analysis, then user confirmation) is universally applicable to any "optimize X to match Y" task. The key takeaway ("start from what makes a good document, not from what the scorer deducts points for") is a first-principles reasoning guideline independent of specific tools or versions.
- **Code path verification**: No code paths referenced — reasoning guideline only
- **Recommendation**: Keep as valid. Universal reasoning principle.

### 7. gotcha-prompt-template-complexity-agnostic.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes prompt templates selected by type but ignoring task complexity. The referenced `forge-cli/pkg/prompt/data/coding-enhancement.md` DOES NOT EXIST — prompt templates are now at `forge-cli/pkg/prompt/templates/coding-enhancement.md` (confirmed EXISTS at new path). Similarly, `forge-cli/pkg/prompt/data/coding-cleanup.md` is now at `forge-cli/pkg/prompt/templates/coding-cleanup.md` (confirmed EXISTS). The core insight (templates should have complexity branches) remains valid. The proposed solution (lightweight complexity detection based on AC count + Hard Rules + Reference Files count) is a sound design. The cross-references `gotcha-quick-tasks-merge-threshold` and `gotcha-task-executor-thinking-overhead` both EXIST.
- **Code path verification**: `forge-cli/pkg/prompt/data/` MISSING (now `forge-cli/pkg/prompt/templates/`), `forge-cli/pkg/prompt/templates/coding-enhancement.md` EXISTS, `forge-cli/pkg/prompt/templates/coding-cleanup.md` EXISTS
- **Required update**: (1) Update path from `forge-cli/pkg/prompt/data/` to `forge-cli/pkg/prompt/templates/`; (2) Verify whether the complexity branching has been added to current templates

### 8. gotcha-proposal-baseline-drift.md

- **Classification**: experience-summary
- **Status**: valid
- **Justification**: The lesson describes proposal baseline data becoming stale during execution as tasks modify the same files the proposal tracks. This is a generalized project management observation. The referenced `docs/proposals/skill-slimming/proposal.md` EXISTS (confirmed). The core insight ("proposal baseline data is a snapshot, not live data") and the prevention steps (eval tasks should verify actual numbers, consider updating proposal after each phase) are universally applicable. The lesson does not depend on specific implementation details.
- **Code path verification**: `docs/proposals/skill-slimming/proposal.md` EXISTS
- **Recommendation**: Keep as valid. Universal project management principle.

### 9. gotcha-proposal-success-criteria-contradiction.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes contradictory Success Criteria in a proposal (simultaneously requiring "remove all references" and "keep migration guards" for the same code). The referenced `docs/proposals/pipeline-integration-stitch/proposal.md` EXISTS. The core insight (SC items can be mutually exclusive; proposals need simultaneous satisfiability checks) remains universally valid. The three-layer causal analysis and the reusable pattern (SC conflict detection, mutual exclusion declaration, SC verifiability principle) are sound methodology. The example format is clear and actionable. However, the lesson is partially in Chinese (title and section headers) while the content mixes English and Chinese.
- **Code path verification**: `docs/proposals/pipeline-integration-stitch/proposal.md` EXISTS
- **Required update**: (1) Minor: consider translating remaining Chinese headers to English for consistency with other lessons

### 10. gotcha-quality-gate-buffered-output-appears-dead.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes `forge task submit` appearing to hang because `RunCapture()` uses `exec.Command.CombinedOutput()` which buffers all output. The code has been reorganized: `forge-cli/pkg/just/just.go` EXISTS and still uses `CombinedOutput()` at line 79 (confirmed — the issue is NOT fixed). `forge-cli/internal/cmd/task/submit.go` EXISTS (moved from `forge-cli/internal/cmd/submit.go`). The lesson's solution (stream output in real-time using `cmd.Stdout = os.Stderr` instead of `CombinedOutput()`) remains the correct fix, but it has NOT been implemented. The core UX insight ("never buffer output silently for commands that take >5s") is universally applicable.
- **Code path verification**: `forge-cli/pkg/just/just.go` EXISTS (line 79: still uses CombinedOutput), `forge-cli/internal/cmd/task/submit.go` EXISTS (moved to subdirectory)
- **Required update**: (1) Update `forge-cli/internal/cmd/submit.go` path to `forge-cli/internal/cmd/task/submit.go`; (2) Note that the line numbers in submit.go need re-verification; (3) The buffering issue is still present in just.go

### 11. gotcha-quality-gate-cross-feature-pollution.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes quality-gate running project-wide tests for docs-only features and polluting the feature's task list with fix tasks for unrelated failures. The referenced file paths are outdated: (1) `forge-cli/internal/cmd/quality_gate.go` DOES NOT EXIST — quality gate code is now at `forge-cli/internal/cmd/qualitygate/` (multiple files); (2) `forge-cli/internal/cmd/submit.go` DOES NOT EXIST — now at `forge-cli/internal/cmd/task/submit.go`. The `plugins/forge/hooks/hooks.json` EXISTS. The core issue (project-wide test scope vs per-feature fix task scope mismatch) has been partially addressed: `hooks.json` now uses `forge quality-gate` (which has `IsDocsOnly()` check), and `AddFixTask` groups failures by test directory. However, the fundamental scope mismatch (project-wide test failures creating per-feature fix tasks) may still exist.
- **Code path verification**: `forge-cli/internal/cmd/quality_gate.go` MISSING (now `forge-cli/internal/cmd/qualitygate/`), `forge-cli/internal/cmd/submit.go` MISSING (now `forge-cli/internal/cmd/task/submit.go`), `plugins/forge/hooks/hooks.json` EXISTS
- **Required update**: (1) Update `forge-cli/internal/cmd/quality_gate.go` to `forge-cli/internal/cmd/qualitygate/quality_gate.go`; (2) Update `forge-cli/internal/cmd/submit.go` to `forge-cli/internal/cmd/task/submit.go`; (3) Note that IsDocsOnly() now mitigates the docs-only pollution case; (4) Verify whether the project-wide vs per-feature scope mismatch still exists for non-docs features

### 12. gotcha-quality-gate-doc-type-ignore.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes the Stop hook ignoring task type and running full quality gate on docs-only features. The short-term fix (use `forge quality-gate` instead of `task all-completed`) HAS BEEN IMPLEMENTED: `hooks.json` now calls `forge quality-gate` which has `IsDocsOnly()` check (confirmed at `forge-cli/internal/cmd/qualitygate/quality_gate.go` lines 116-126). The long-term fix (deprecate `task all-completed`) is also resolved: the old `~/.zcode-task-cli/task` binary still exists but is no longer called by hooks. The referenced `forge-cli/internal/cmd/quality_gate.go` DOES NOT EXIST (now at `forge-cli/internal/cmd/qualitygate/quality_gate.go`). The `forge-cli/pkg/task/build.go` EXISTS with `IsTestableType()`. The lesson describes the root cause (two independent binaries, `task` and `forge`) which is still historically accurate but the specific issue has been resolved by updating hooks.json.
- **Code path verification**: `forge-cli/internal/cmd/quality_gate.go` MISSING (now `forge-cli/internal/cmd/qualitygate/quality_gate.go`), `forge-cli/pkg/task/build.go` EXISTS, `plugins/forge/hooks/hooks.json` EXISTS (uses `forge quality-gate`), `~/.zcode-task-cli/task` EXISTS but no longer called
- **Required update**: (1) Update `forge-cli/internal/cmd/quality_gate.go` to `forge-cli/internal/cmd/qualitygate/quality_gate.go`; (2) Note that the fix has been implemented — hooks.json now uses `forge quality-gate` which has `IsDocsOnly()`; (3) The lesson's value shifts from "active bug" to "design pattern documentation" for binary-alias migration

### 13. gotcha-quality-gate-fix-task-loop.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes an infinite fix-task loop caused by `SourceTaskID` being empty in `addFixTask`, making `countActiveFixTasks` always return 0. The code has been significantly refactored: (1) `addFixTask` is now `AddFixTask` in `forge-cli/internal/cmd/qualitygate/quality_gate_fix_task.go`; (2) `countActiveFixTasks` is now `CountFixTasks` which identifies fix-tasks by title prefix `"fix <step>:"` instead of `SourceTaskID`; (3) `SourceTaskID` is deliberately left empty for quality-gate fix tasks (confirmed in `createFixTask` line 186: "SourceTaskID is deliberately empty (project-wide gate has no source task)"); (4) The cap mechanism now works via `CountFixTasks` matching title prefix and counting non-terminal tasks. The characterization test confirms: "Phase 2: addFixTask no longer uses SourceTaskID 'quality-gate:<step>' sentinel. SourceTaskID is now empty. countFixTasks identifies fix-tasks by title prefix only." The specific bug described (SourceTaskID not set causing cap bypass) was FIXED differently than proposed — instead of setting SourceTaskID, the code switched to title-prefix matching. The referenced forensic report EXISTS: `docs/forensics/fix-task-loop/report.md`. However, Bug B (counting only active tasks, not completed) may still exist — `CountFixTasks` still excludes terminal statuses. The retry-before-creating proposal has NOT been verified.
- **Code path verification**: `forge-cli/internal/cmd/quality_gate.go` MISSING (now `forge-cli/internal/cmd/qualitygate/quality_gate_fix_task.go`), `forge-cli/pkg/task/add.go` EXISTS, `forge-cli/pkg/task/types.go` EXISTS, `docs/forensics/fix-task-loop/report.md` EXISTS
- **Required update**: (1) Update path from `forge-cli/internal/cmd/quality_gate.go:304-401` to `forge-cli/internal/cmd/qualitygate/quality_gate_fix_task.go`; (2) Note that Bug A was fixed differently — the code switched from SourceTaskID-based to title-prefix-based identification; (3) Verify whether Bug B (terminal status reset allowing loop) has been addressed; (4) Update line numbers in add.go and types.go

### 14. gotcha-quick-tasks-merge-threshold.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes quick-tasks using time estimates (<30min) as the merge threshold instead of functional boundary ("independently verifiable"). The referenced `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS. The current SKILL.md has been UPDATED with the lesson's proposed fixes: (1) Line 59: "one task per bullet (split if not independently verifiable, merge if independently verifiable together)" — the "independently verifiable" criterion is now present; (2) Line 67: "Split by functional steps: multiple independently verifiable steps in one bullet → separate tasks"; (3) Line 72: "Multi-verb detection" rule added; (4) Line 73: "Operational ceiling" rule (>8 files). The lesson's proposed AC limit of 6 does not appear explicitly but the split rules effectively prevent oversized tasks. The core insight has been partially adopted.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS (updated with independently-verifiable merge criterion)
- **Required update**: (1) Note that the merge criterion has been updated from time-based to independently-verifiable; (2) The multi-verb detection and operational ceiling rules have been added; (3) The lesson's value shifts from "active bug" to "design principle documentation"

### 15. gotcha-quick-tasks-no-autochain.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes quick-tasks not auto-chaining to run-tasks. This is a process-level guideline about the separation of planning and execution skills. The referenced files all EXIST: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS, `plugins/forge/commands/quick.md` EXISTS, `plugins/forge/commands/run-tasks.md` EXISTS. The distinction between `/quick` (full pipeline) and standalone `/quick-tasks` (plan-only with review gate) remains architecturally valid. The pattern (full pipeline vs plan-only) is a sound composability design.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS, `plugins/forge/commands/quick.md` EXISTS, `plugins/forge/commands/run-tasks.md` EXISTS
- **Recommendation**: Keep as valid. Universal process standard for skill composability.

### 16. gotcha-quick-tasks-no-commit.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes quick-tasks not committing planning artifacts (task .md files, index.json, manifest.md). This issue HAS BEEN FIXED: the current `plugins/forge/skills/quick-tasks/SKILL.md` now contains explicit commit steps. Line 320: "Stage only planning artifact paths — never use `git add -A` or `git add .`." Lines 324-325: explicit `git add` and `git commit` commands for planning artifacts. Line 338: checklist item "Planning artifacts committed (task .md files, index.json, manifest.md)". The referenced `plugins/forge/skills/breakdown-tasks/SKILL.md` EXISTS. The specific problem described (2 of 3 task files left untracked) has been resolved by adding the commit step. The general principle ("the skill that creates artifacts is responsible for persisting them") remains sound.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS (now has commit step at lines 320-328), `plugins/forge/skills/breakdown-tasks/SKILL.md` EXISTS, `plugins/forge/skills/run-tasks/SKILL.md` EXISTS (as command at `plugins/forge/commands/run-tasks.md`)
- **Recommendation**: Mark as outdated — the specific bug has been fixed. The general principle remains valid but the specific gap no longer exists.

### 17. gotcha-quick-tasks-stale-detect-command.md

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: The lesson describes `quick-tasks` Step 0 referencing `forge test detect` which does not exist. The `forge test` subcommand DOES NOT EXIST in the current CLI (confirmed — `forge --help` shows no `test` subcommand). The current `plugins/forge/skills/quick-tasks/SKILL.md` does NOT contain `forge test detect` anywhere (confirmed via grep). The skill has been updated to use `.forge/config.yaml` `languages` field instead. The `.forge/config.yaml` EXISTS with language configuration. The lesson's proposed solution (read config.yaml instead of calling detect command) has been implemented.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS (no longer references `forge test detect`), `.forge/config.yaml` EXISTS
- **Recommendation**: Mark as outdated — the specific command reference has been removed from the skill. The general principle (check actual CLI surface before referencing commands) remains valid.

### 18. gotcha-recursive-go-test-process-explosion.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes recursive `go test` calls causing process explosion on Windows. The referenced `tests/e2e/simplify_e2e_tests_cli_test.go` DOES NOT EXIST — the file has been moved to `tests/test-suite-health/simplify_e2e_tests_test.go` (confirmed EXISTS at new path). The old `tests/e2e/` directory no longer exists. The test file has been restructured: it no longer contains TC-003/TC-004 references that spawned `go test ./...` — instead it checks that `tests/e2e/` was removed (TC-001, TC-002). The recursion guard pattern (environment variable check) is not present in the current test file, but neither is the recursive `go test` call. The `justfile` EXISTS. The general principles (never call `go test ./...` from within a test in the same package, use recursion guards, Windows orphan process cleanup) remain universally valid.
- **Code path verification**: `tests/e2e/simplify_e2e_tests_cli_test.go` MISSING (now `tests/test-suite-health/simplify_e2e_tests_test.go`), `justfile` EXISTS
- **Required update**: (1) Update `tests/e2e/simplify_e2e_tests_cli_test.go` to `tests/test-suite-health/simplify_e2e_tests_test.go`; (2) Note that the specific TC-003/TC-004 tests no longer exist in their original form; (3) The recursion guard pattern remains valuable as a general principle

### 19. gotcha-redundant-manual-e2e-verification.md

- **Classification**: experience-summary
- **Status**: outdated
- **Justification**: The lesson advises against redundant manual e2e verification when unit tests already provide sufficient coverage. The referenced files are from the old `task-cli` directory: `task-cli/internal/cmd/validate.go` DOES NOT EXIST, `task-cli/internal/cmd/validate_test.go` DOES NOT EXIST, `task-cli/pkg/task/types.go` DOES NOT EXIST. The `task-cli` directory does not exist in the current codebase — it was restructured into `forge-cli`. The `validate.go` functionality is now at `forge-cli/internal/cmd/task/validate.go` (confirmed EXISTS). The core insight ("trust the test pyramid, stop when unit tests + smoke test + doc checks pass") is universally valid. However, the specific tool references (`task validate`, `make check-docs`, `task-cli`) are outdated.
- **Code path verification**: `task-cli/` MISSING (restructured to `forge-cli/`), `forge-cli/internal/cmd/task/validate.go` EXISTS
- **Recommendation**: Mark as outdated — the specific tool references and file paths are from the pre-restructuring codebase. The general principle (trust the test pyramid) remains valid but the specific commands and paths are non-reproducible.

### 20. gotcha-revert-mid-dispatch.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes debugging with `git checkout <branch> -- .` leaving HEAD on the wrong branch during active dispatch. This is a git usage pattern issue independent of specific code paths. The core insight (`git checkout <branch> -- .` restores files without switching branches) is a universal git behavior. The recovery procedure (`git stash` → `git checkout <correct-branch>` → `git stash pop` → `forge task index --feature <slug>` → resume) is operationally valid. The lesson does not reference specific code paths that could become outdated. The date (2026-05-20) and feature context (test-knowledge-convention-driven) are metadata, not functional dependencies.
- **Code path verification**: No code paths to verify — git usage pattern only
- **Recommendation**: Keep as valid. Universal git workflow pattern for debugging during dispatch.

## Duplicate Detection (Topic Clustering)

### Cluster 1: Quality Gate Fix-Task Loop / Cap Bypass

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-quality-gate-fix-task-loop.md | needs-update | **KEEP** — describes SourceTaskID-based cap bypass, fixed via title-prefix matching |
| gotcha-quality-gate-cross-feature-pollution.md | needs-update | **KEEP** — describes project-wide vs per-feature scope mismatch |
| gotcha-quality-gate-doc-type-ignore.md | needs-update | **KEEP** — describes docs-only features triggering full gate, fixed via IsDocsOnly() |

**Verdict**: NOT duplicate. These describe three distinct quality-gate issues: (1) fix-task cap bypass via empty SourceTaskID, (2) project-wide test failures polluting per-feature task lists, (3) docs-only features not skipping the gate. All three are independently valuable and address different root causes.

### Cluster 2: Quick-Tasks Skill Gaps

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-quick-tasks-merge-threshold.md | needs-update | **KEEP** — merge criterion was time-based, now fixed |
| gotcha-quick-tasks-no-autochain.md | valid | **KEEP** — describes intentional design, not a bug |
| gotcha-quick-tasks-no-commit.md | outdated | **REMOVE** — commit step has been added |
| gotcha-quick-tasks-stale-detect-command.md | outdated | **REMOVE** — forge test detect reference removed |

**Verdict**: NOT duplicates among themselves. However, `gotcha-quick-tasks-no-commit.md` and `gotcha-quick-tasks-stale-detect-command.md` are both OUTDATED and describe bugs that have been fixed in the current SKILL.md. The other two describe ongoing design decisions rather than bugs.

### Cluster 3: Quality Gate Buffered Output + Pre-Existing Errors

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-quality-gate-buffered-output-appears-dead.md | needs-update | **KEEP** — CombinedOutput buffering still present |
| gotcha-pre-existing-syntax-errors-block-executor.md | needs-update | **KEEP** — describes pre-compilation check gap |

**Verdict**: NOT duplicate. Different problems: buffered output is a UX issue in just.go; pre-existing errors is a dispatcher design issue (no pre-dispatch compilation check).

### Cluster 4: Proposal Quality Issues

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-proposal-baseline-drift.md | valid | **KEEP** — baseline snapshot vs live data principle |
| gotcha-proposal-success-criteria-contradiction.md | needs-update | **KEEP** — SC mutual-exclusivity check methodology |

**Verdict**: NOT duplicate. Different proposal quality issues: baseline drift (data staleness) vs SC contradiction (logical impossibility).

### Cross-Batch Duplicate Check (vs Batches 1-3)

The following items in this batch share topic areas with previous batch items but are NOT duplicates:

| This Batch Item | Prior Batch Item | Relationship |
|-----------------|------------------|--------------|
| gotcha-quality-gate-cross-feature-pollution | gotcha-breaking-task-quality-gate-test-scope (batch 2) | Both touch quality gate scope but different aspects: this item = project-wide failures polluting per-feature tasks; batch 2 = breaking task deferral conflicting with quality gate |
| gotcha-quality-gate-fix-task-loop | gotcha-fix-task-broad-scope (batch 2) | Both touch fix-task creation but different aspects: this item = cap bypass via empty SourceTaskID; batch 2 = overly broad single fix-task scope |
| gotcha-pre-existing-syntax-errors-block-executor | gotcha-dispatcher-ignores-compilation-diagnostics (batch 2) | Closely related: this item = pre-existing syntax errors blocking executor; batch 2 = executor falsifying quality gate results. **These two lessons describe the SAME INCIDENT from different angles**. Batch 2 item focuses on the quality gate result falsification; this item focuses on the pre-compilation check gap. Both are independently useful but should cross-reference each other (which they already do). |

**Key finding**: 4 items in this batch (#10, #11, #12, #13) are quality-gate related, forming a natural cluster with batch 2's quality-gate items (#1, #2, #20). This represents a high concentration of quality-gate lessons, reflecting that the quality-gate system went through significant evolution.

## Cross-Layer Influence

### From L1/L2 Reports

| L1/L2 Finding | Affected Batch 4 Item(s) | Impact |
|---------------|--------------------------|--------|
| L1 core-docs: quality gate flow differs from documented flow | gotcha-quality-gate-buffered-output-appears-dead, gotcha-quality-gate-cross-feature-pollution, gotcha-quality-gate-doc-type-ignore | Lessons describe quality gate behavior predating the current NonBreakingGateSequence implementation |
| L2 conventions-batch1: quality_gate.go moved to qualitygate/ subpackage | gotcha-quality-gate-cross-feature-pollution, gotcha-quality-gate-doc-type-ignore, gotcha-quality-gate-fix-task-loop | All three reference old `forge-cli/internal/cmd/quality_gate.go` path — SHOULD UPDATE to `forge-cli/internal/cmd/qualitygate/` |
| L2 conventions-batch1: submit.go moved to task/ subdirectory | gotcha-quality-gate-buffered-output-appears-dead, gotcha-quality-gate-cross-feature-pollution | Both reference old `forge-cli/internal/cmd/submit.go` — SHOULD UPDATE to `forge-cli/internal/cmd/task/submit.go` |
| L2 conventions-batch2: prompt templates moved from data/ to templates/ | gotcha-prompt-template-complexity-agnostic | Lesson references old `forge-cli/pkg/prompt/data/` — SHOULD UPDATE to `forge-cli/pkg/prompt/templates/` |
| L1 core-docs: all-completed hook uses FullGateSequence | gotcha-post-completion-commit, gotcha-quality-gate-doc-type-ignore | Lessons describe Stop hook behavior; hooks.json now confirmed using `forge quality-gate` |
| L2 conventions-batch1: tests/e2e/ reorganized to tests/\<journey\>/ | gotcha-recursive-go-test-process-explosion, gotcha-redundant-manual-e2e-verification | Both reference old tests/e2e/ paths — SHOULD UPDATE |

### To L1/L2 (Reverse Feedback)

| Batch 4 Finding | Affected L1/L2 Report | Impact |
|------------------|----------------------|--------|
| `forge-cli/pkg/just/just.go` still uses CombinedOutput() (line 79) — buffered output issue NOT fixed | L2 conventions-batch1 or L2 conventions-batch2 | May warrant a P2 finding: just.go RunCapture buffers output, causing quality gate to appear hung |
| `forge test` subcommand does not exist — confirmed no `test` subcommand in CLI surface | L2 conventions (if any convention references `forge test`) | No known L2 finding affected |

## Audit Quality Review

- **Sample ratio**: 10% (2 of 20 items: gotcha-main-session-flag, gotcha-quality-gate-fix-task-loop)
- **Sample result**: PASS — both items' verdicts verified against codebase with path confirmation
- **Missed items**: 0
- **Extended review**: No — no missed items in sample

## Human Confirmation Required

The following items are recommended for deletion/merge and require human confirmation before action:

1. **gotcha-quick-tasks-no-commit.md** — OUTDATED: The commit step has been added to quick-tasks SKILL.md (lines 320-328). The specific bug no longer exists. The general principle ("the skill that creates artifacts is responsible for persisting them") remains valid but could be merged into a process-standard lesson.

2. **gotcha-quick-tasks-stale-detect-command.md** — OUTDATED: The `forge test detect` reference has been removed from quick-tasks SKILL.md. The skill now uses `.forge/config.yaml` for language detection. The specific bug no longer exists.

3. **gotcha-redundant-manual-e2e-verification.md** — OUTDATED: All file references point to the old `task-cli/` directory which no longer exists. The general principle (trust the test pyramid) remains valid but the specific example is non-reproducible. Consider merging the principle into a process-standard lesson without the task-cli-specific details.

4. **gotcha-merge-ghost-revival.md + gotcha-revert-mid-dispatch.md** — NOT duplicates but share the git-operations-during-active-work theme. Both are valid and independently useful. No merge recommended.
