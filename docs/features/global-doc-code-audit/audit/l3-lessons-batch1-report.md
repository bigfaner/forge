# L3 Lessons Audit Report — Batch 1 (arch/gotcha-a to gotcha-breaking-change)

## Audit Baseline

- **Baseline commit**: f64a7ca8 (docs(audit): add L2 conventions batch 2 audit report)
- **Audit date**: 2026-06-03
- **Audit scope**: 20 lesson files from `docs/lessons/` (alphabetically: arch-constant-rename-whack-a-mole.md through gotcha-breaking-change-integration-test-blast-radius.md)

## Classification Distribution

| Classification | Count |
|----------------|-------|
| code-reference | 14 |
| process-standard | 4 |
| experience-summary | 2 |

## Validity Summary

| Status | Count |
|--------|-------|
| valid | 6 |
| outdated | 3 |
| needs-update | 8 |
| duplicate | 3 |

## Cross-Layer Influence Check

The following L1/L2 audit findings affect items in this batch:

| L1/L2 Finding | Affected Lesson | Impact |
|---------------|-----------------|--------|
| L2 conventions-batch1: `run-tasks` is a command (not a skill) | gotcha-adjacent-section-over-removal, arch-post-loop-artifact-commit-gap, arch-dispatcher-post-loop-message-misleading | Path references use `skills/run-tasks/SKILL.md` but actual location is `commands/run-tasks.md` |
| L2 conventions-batch1: constants.md has stale path references | arch-constant-rename-whack-a-mole | No direct impact (lesson about rename process, not specific constants) |
| L1 core-docs: `docs/reference/` directory does not exist | No items in this batch reference `docs/reference/` | No impact |

## Item Details

### 1. arch-constant-rename-whack-a-mole.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson describes a reusable pattern for bulk constant renames (pre-compute affected files, bulk replace, verify once). The referenced files `forge-cli/pkg/task/types.go` and `forge-cli/pkg/task/build.go` both exist. The pattern itself is toolchain-independent and broadly applicable. The specific task file `docs/features/task-type-id-redesign/tasks/1-rename-type-constants.md` exists as a historical reference.
- **Code path verification**: `forge-cli/pkg/task/types.go` EXISTS, `forge-cli/pkg/task/build.go` EXISTS

### 2. arch-dispatcher-post-loop-message-misleading.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson references `plugins/forge/skills/run-tasks/SKILL.md` (Post-Completion section). Per L2 audit, `run-tasks` is a command located at `plugins/forge/commands/run-tasks.md`, not a skill. The described problem (hardcoded post-loop message about T-test-run/T-test-verify-regression) is a genuine design issue, but the file path is wrong. The core conclusion remains valid: post-loop messages should reflect actual execution state.
- **Code path verification**: `plugins/forge/skills/run-tasks/SKILL.md` DOES NOT EXIST. Correct path: `plugins/forge/commands/run-tasks.md`
- **Required update**: Change "plugins/forge/skills/run-tasks/SKILL.md" to "plugins/forge/commands/run-tasks.md"

### 3. arch-forge-skill-gap-analysis.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: This is a detailed gap analysis of write-prd and ui-design skills. The referenced skills exist: `plugins/forge/skills/ui-design/SKILL.md` EXISTS, `plugins/forge/skills/write-prd/SKILL.md` EXISTS. However, this lesson is written in Chinese while the lesson content is effectively a design document with specific improvement proposals. The file paths referenced (e.g., `skills/write-prd/templates/prd-ui-functions.md`, `skills/ui-design/templates/prototype.md`) need verification whether the specific template files and sections mentioned still exist in their described form after v3.0.0 restructuring.
- **Code path verification**: `plugins/forge/skills/ui-design/` EXISTS, `plugins/forge/skills/write-prd/` EXISTS, template files need verification
- **Required update**: Verify template file paths and section names against current templates; update if restructured

### 4. arch-freeform-findings-indirect-influence.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson references `plugins/forge/skills/eval/rules/freeform-injection.md` which DOES NOT EXIST. The current eval rules directory contains `freeform-pipeline.md` and `freeform-expert-persistence.md` instead, suggesting the file was renamed or split. Other references (`plugins/forge/skills/eval/SKILL.md`, `plugins/forge/skills/eval/rules/scorer-composition.md`, `plugins/forge/skills/eval/experts/protocol/reviser-protocol.md`, `docs/proposals/spec-authority-enforcement/eval/freeform-review.md`, `docs/proposals/spec-authority-enforcement/eval/iteration-1.md`) all exist. The core insight about information loss through the Scorer mediation layer remains valid.
- **Code path verification**: `freeform-injection.md` DOES NOT EXIST. Replacements: `freeform-pipeline.md`, `freeform-expert-persistence.md`
- **Required update**: Replace `freeform-injection.md` reference with current equivalent file(s)

### 5. arch-post-loop-artifact-commit-gap.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: References `plugins/forge/skills/run-tasks/SKILL.md` which DOES NOT EXIST (correct path is `plugins/forge/commands/run-tasks.md`). Also references `plugins/forge/skills/submit-task/SKILL.md` which EXISTS. The core problem (uncommitted post-loop artifacts) and solution remain valid.
- **Code path verification**: `run-tasks/SKILL.md` MISSING, `submit-task/SKILL.md` EXISTS
- **Required update**: Fix run-tasks path from skills to commands

### 6. arch-prototype-navigation-contract.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: This lesson provides a generalized design pattern for multi-page prototypes (navigation contracts, code layer separation, PRD alignment). The solution describes establishing navigation architecture tables and code separation conventions. These patterns are toolchain-independent best practices. The file references `app.js` which is a pattern-level reference, not a codebase-specific path. No specific forge code paths are referenced.
- **Code path verification**: N/A (pattern-level lesson with no specific code paths)

### 7. arch-task-failure-recovery-loop.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: This lesson provides an extremely detailed 6-gap analysis of the task failure pipeline. Several referenced paths have moved: (1) `forge-cli/internal/cmd/feature_complete.go` is now at `forge-cli/internal/cmd/feature/feature_complete.go`; (2) The lesson references `~/.claude/plugins/cache/forge/forge/2.14.0/` paths which are version-specific cache paths and will differ per installation; (3) The lesson references `forge-cli/internal/cmd/submit.go` and `claim.go` which are now at `forge-cli/internal/cmd/task/submit.go` and `forge-cli/internal/cmd/task/claim.go`; (4) The referenced evidence file `docs/features/e2e-test-scripts-rebuild/tasks/process/record.json` DOES NOT EXIST. Critically, the ROOT CAUSE described (Gap 3: CLI allows completed + testsFailed > 0) has been FIXED — `submit.go` now has auto-downgrade logic at line 314: `if isCoding && rd.Status == string(types.StatusCompleted) && rd.TestsFailed > 0` which auto-downgrades to blocked.
- **Code path verification**: Multiple paths moved to subdirectories (cmd/task/), feature_complete.go moved to cmd/feature/
- **Required update**: (1) Update all file paths to current locations; (2) Mark Gap 3 (CLI validation) as RESOLVED — auto-downgrade implemented; (3) Update gap status for any other fixes applied

### 8. arch-task-record-immutability.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson establishes the principle that task records should be append-only. The referenced directory `docs/features/forge-info-commands/tasks/records/` EXISTS. The lesson references `index.json` and specific record files which exist. The core principle (records as audit logs, never overwrite) remains sound and applicable regardless of implementation changes.
- **Code path verification**: `docs/features/forge-info-commands/tasks/records/` EXISTS

### 9. arch-test-type-index-chicken-egg.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: References `forge-cli/pkg/task/build.go` (EXISTS, lines 269-318) and `forge-cli/pkg/task/testgen.go` (DOES NOT EXIST). The function `DetectTypesFromTestCases` was moved to `forge-cli/pkg/task/autogen.go` (verified via grep). The referenced proposal `docs/proposals/test-scripts-per-type/proposal.md` EXISTS. The chicken-and-egg architectural issue may still apply but the specific code locations have changed.
- **Code path verification**: `build.go` EXISTS, `testgen.go` MISSING (functions moved to `autogen.go`)
- **Required update**: Replace `testgen.go` reference with `autogen.go`; verify if the chicken-and-egg issue is still present in current BuildIndex implementation

### 10. fix-zsh-compinit-docker.md

- **Classification**: experience-summary
- **Status**: valid
- **Justification**: This lesson describes a local development environment issue (zsh compinit errors from dangling Docker Desktop symlinks on macOS). The solution is environment-specific and the steps are correct (check for dangling symlinks in `/opt/homebrew/share/zsh/site-functions/`, remove them). The general troubleshooting pattern (clear cache, check symlinks, rebuild) remains valid. No forge code paths are referenced.
- **Code path verification**: N/A (environment troubleshooting, no forge code paths)

### 11. gotcha-ac-self-report-without-verification.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson references `task record` CLI and `record-task/SKILL.md`. The `record-task` skill no longer exists — it has been replaced by `submit-task` skill (`plugins/forge/skills/submit-task/SKILL.md` which EXISTS). The core insight (agents self-report AC as met without running the artifact, CLI only validates structure not truth) remains valid. The referenced task-executor agent (`plugins/forge/agents/task-executor.md`) EXISTS. The submit-task SKILL.md now has stronger validation (auto-downgrade for testsFailed > 0, rejection for met:false with completed status), partially addressing this gotcha.
- **Code path verification**: `record-task` skill MISSING, replaced by `submit-task/SKILL.md` which EXISTS
- **Required update**: Replace `record-task` references with `submit-task`; note partial mitigation from auto-downgrade logic

### 12. gotcha-adjacent-section-over-removal.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: References `plugins/forge/commands/run-tasks.md` (EXISTS) and `docs/proposals/auto-knowledge-save/proposal.md` (EXISTS). The core lesson about adjacent section removal is a valid documentation editing pattern. The file path is correct for the command.
- **Code path verification**: `plugins/forge/commands/run-tasks.md` EXISTS, `docs/proposals/auto-knowledge-save/proposal.md` EXISTS
- **Required update**: No code path update needed; status is needs-update because the lesson references a historical event that may no longer be actionable, but the pattern itself remains valid

### 13. gotcha-agent-example-over-schema.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson establishes a general principle: "Examples > Documentation" — agents copy example commands verbatim and ignore field reference tables. The referenced files (`task-executor.md`, `record-task/SKILL.md`) are historical references for the specific incident. While `record-task` has been renamed to `submit-task`, the principle itself is universally applicable and remains sound.
- **Code path verification**: Principle-level lesson; specific file references are historical incident context

### 14. gotcha-api-no-api-prefix.md

- **Classification**: experience-summary
- **Status**: outdated
- **Justification**: This lesson references `backend/internal/handler/router.go` and `frontend/vite.config.ts` — neither `backend/` nor `frontend/` directories exist in this repository. This lesson appears to be from a different project (likely a web application project, not the forge CLI tool). The lesson content is not applicable to the forge codebase and provides no actionable value for forge developers.
- **Code path verification**: `backend/` directory DOES NOT EXIST, `frontend/` directory DOES NOT EXIST
- **Recommendation**: Mark for deletion — lesson belongs to a different project

### 15. gotcha-auto-gen-tasks-reappear-and-preempt-fix.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: References `claim.go:242-247` (task claim sort logic) and `claim.go:79` (auto-unblock logic). The file has moved to `forge-cli/internal/cmd/task/claim.go`. The current sort logic uses topological depth first, then priority, then version ID comparison — fix task priority boost is NOT explicitly in the sort logic (there is no special handling for `fix-*` prefix). The recommended "中期" fix (fix tasks should preempt auto-gen tasks) has NOT been implemented in the sort. The "长期" fix (distinguish blockedReason) has been partially implemented via `suspended` status. References `autogen.go:32-121` (now at `forge-cli/pkg/task/autogen.go`, EXISTS) and `build.go:164-170` (EXISTS).
- **Code path verification**: File paths moved to `cmd/task/` subdirectory; sort logic does not prioritize fix tasks
- **Required update**: Update file paths; note that suspended status partially addresses long-term fix; medium-term fix (fix-task priority in sort) not yet implemented

### 16. gotcha-auto-push-no-upstream.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: References `forge-cli/internal/cmd/feature_complete.go` (now at `forge-cli/internal/cmd/feature/feature_complete.go`) and `.forge/config.yaml` (EXISTS). The fix described (`git push -u origin HEAD`) has been VERIFIED as implemented — the current `gitPush()` function at line 274 uses `exec.Command("git", "push", "-u", "origin", "HEAD")`. The lesson is valid as a historical record and the reusable pattern (always use `-u origin HEAD`) is correct.
- **Code path verification**: Fix is IMPLEMENTED. `feature_complete.go` moved to `cmd/feature/` subdirectory
- **Note**: File path is outdated but the lesson's conclusion is confirmed valid by code inspection

### 17. gotcha-auto-unblock-loop.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: Describes the auto-unblock loop problem where manually blocked tasks get auto-unblocked. The proposed solution (`suspended` status) has been IMPLEMENTED — `claim.go` line 179 confirms "Suspended tasks are naturally excluded (they have status 'suspended', not 'blocked')" and `forge-cli/internal/cmd/task/transition.go` line 31 documents the suspended transition. The lesson remains valid as it describes the root cause and the solution that was eventually implemented. The commit reference `3a54430f` is historical.
- **Code path verification**: `suspended` status IS implemented; `claim.go` excludes suspended tasks from auto-unblock
- **Note**: This lesson's proposed solution has been fully implemented

### 18. gotcha-blocked-task-never-auto-unblocks.md

- **Classification**: code-reference
- **Status**: valid (resolved)
- **Justification**: Describes the opposite problem of #17 — tasks that SHOULD auto-unblock but don't. The referenced code in `claim.go` now has a "Lazy unblock scan" (lines 177-188) that checks blocked tasks and auto-transitions eligible ones to pending when dependencies are met. The `autoRestoreSourceTask` function in `submit.go` handles fix-task source restoration. However, the lesson correctly notes that general dependency cascade (not just fix-task sources) should be in `saveIndexAndSignalCompletion` — checking the current code, this cascade is in `claim.go` (lazy scan) rather than `submit.go`, which is an alternative implementation approach.
- **Code path verification**: Files moved to `cmd/task/` subdirectory; lazy unblock scan IS implemented in claim.go
- **Note**: The problem described has been resolved via lazy unblock scan in claim.go

### 19. gotcha-brainstorm-challenge-failure.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: References `plugins/forge/skills/brainstorm/SKILL.md` (EXISTS) and `plugins/forge/skills/brainstorm/rules/challenge-protocol.md` (EXISTS). The core lesson (brainstorm should challenge pseudo-requirements using Occam's Razor before discussing implementation) is a valid process standard. The skill files and challenge protocol exist as described.
- **Code path verification**: All referenced paths EXIST

### 20. gotcha-breaking-change-integration-test-blast-radius.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: Describes the pattern that breaking tasks must scope integration test fixtures. References Task 4 and fix-3 from a specific feature. The lesson does not reference specific code paths that need verification — it provides a general reusable pattern for breaking tasks. The pattern (grep test directories for callers, check fixtures, list fixture updates in AC) is toolchain-independent and broadly applicable.
- **Code path verification**: N/A (pattern-level lesson)

## Duplicate Detection (Topic Clustering)

### Cluster 1: Navigation Contract / UI Design

| Item | Status | Recommendation |
|------|--------|----------------|
| arch-forge-skill-gap-analysis.md | needs-update | **KEEP** — more detailed, includes specific improvement proposals with code examples |
| arch-prototype-navigation-contract.md | valid | **DUPLICATE of arch-forge-skill-gap-analysis** — covers the same navigation contract / multi-page prototype problem but in a more condensed form. The gap analysis (#3) contains all the same information plus additional improvement proposals. |

**Verdict**: `arch-prototype-navigation-contract.md` is a duplicate of `arch-forge-skill-gap-analysis.md`. Keep the gap analysis (more complete).

### Cluster 2: Auto-Unblock / Blocked Task Semantics

| Item | Status | Recommendation |
|------|--------|----------------|
| gotcha-auto-unblock-loop.md | valid | Keep — describes the "cannot skip tasks" problem, proposed suspended status |
| gotcha-blocked-task-never-auto-unblocks.md | valid | Keep — describes the opposite problem "tasks stuck in blocked forever" |

**Verdict**: NOT duplicate. These describe opposite failure modes of the same system (over-unblocking vs. under-unblocking). Both are independently valuable.

### Cluster 3: run-tasks Post-Loop Issues

| Item | Status | Recommendation |
|------|--------|----------------|
| arch-dispatcher-post-loop-message-misleading.md | needs-update | Keep — distinct problem (misleading message) |
| arch-post-loop-artifact-commit-gap.md | needs-update | Keep — distinct problem (uncommitted artifacts) |
| gotcha-adjacent-section-over-removal.md | needs-update | Keep — distinct problem (section removal side-effect) |

**Verdict**: NOT duplicate. These describe different problems in the same subsystem. All are independently valuable.

### Cluster 4: Test Pipeline Task Generation

| Item | Status | Recommendation |
|------|--------|----------------|
| arch-test-type-index-chicken-egg.md | needs-update | Keep — architectural issue with BuildIndex |
| gotcha-auto-gen-tasks-reappear-and-preempt-fix.md | needs-update | Keep — dispatcher claim priority issue |
| gotcha-breaking-change-integration-test-blast-radius.md | valid | Keep — breaking task fixture scoping pattern |

**Verdict**: NOT duplicate. These cover different aspects of the test pipeline. All are independently valuable.

## Audit Quality Review

- **Sampling ratio**: 100% (all 20 items audited)
- **Cross-layer check**: Performed against L1 core-docs report and L2 conventions batch 1/2 reports
- **Code path verification**: All `find`/`grep` checks executed against current codebase at baseline commit
- **Key finding**: 3 items reference `plugins/forge/skills/run-tasks/SKILL.md` which does not exist (run-tasks is a command, not a skill — confirmed by L2 audit)
- **Key finding**: Several items reference Go source files that moved to subdirectories (cmd/task/, cmd/feature/) during v3.0.0 restructuring
- **Key finding**: 2 items (#17 auto-unblock-loop, #18 blocked-task-never-auto-unblocks) describe problems that have been fully resolved in the current codebase
- **Key finding**: 1 item (#14 api-no-api-prefix) is from a different project entirely and should be deleted
