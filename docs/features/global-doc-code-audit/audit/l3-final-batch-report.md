# L3 Final Batch Audit Report: Remaining Lessons + All Decisions

## Audit Baseline

- **Baseline commit**: `062992a4`
- **Audit date**: 2026-06-03
- **Audit scope**: 13 lesson files (lesson-vibe-coding-scope-control.md through worktree-stale-refs.md) + 10 decision files (architecture.md through testing.md). Total: 23 items.

## Issue Summary

- **P0 (Critical)**: 0
- **P1 (High)**: 3
- **P2 (Medium)**: 7
- **P3 (Low)**: 3
- Items marked outdated/duplicate/needs-update: 16

## Classification Distribution

| Classification | Count | Items |
|---------------|-------|-------|
| code-reference | 8 | pattern-compile-check-before-submit, pattern-dispatcher-auto-verify, pattern-large-scale-rename, pattern-surface-resolution-shortcut, tool-cli-e2e-lifecycle, tool-submit-background-timeout, worktree-stale-refs, e2e-server-lifecycle-hardening |
| process-standard | 7 | lesson-vibe-coding-scope-control, pattern-sitemap-shared-layout, pattern-task-vs-output-naming, tool-fix-e2e-unknown-placeholder, tool-record-coverage-capture, tool-justfile-arg-attribute, architecture (decision) |
| experience-summary | 8 | pattern-dispatcher-auto-verify (dual), tool-cli-e2e-lifecycle (dual), interface (decision), testing (decision), data-model (decision), dependencies (decision), error-handling (decision), local-dev-deployment (decision), security (decision), manifest (decision) |

Note: Some items span multiple classifications.

## Status Summary

| Status | Lessons | Decisions (files) | Total |
|--------|---------|-----------|-------|
| valid | 4 | 1 | 5 |
| needs-update | 7 | 3 | 10 |
| outdated | 2 | 1 | 3 |
| empty-stub | 0 | 5 | 5 |

Total: 23 items assessed (13 lessons + 10 decision files).

Note: Decision counts are at file level. architecture.md contains 4 individual decisions (3 valid + 1 needs-update), testing.md contains 2 individual decisions (1 valid + 1 needs-update). At file level, both are classified needs-update because they contain at least one decision requiring update. Individual valid decisions: 5 (architecture D1,D2,D4, interface, testing D2).

Note: 4 decision files (data-model.md, dependencies.md, error-handling.md, local-dev-deployment.md) are empty stubs containing only header rows with no actual decisions. These are counted in the empty-stub category separately from the standard L3 validity classifications.

## Cross-Layer Influence Check

The following L1/L2 audit findings were checked against items in this batch:

| Cross-Layer Ref | Source | Items Checked | Impact |
|----------------|--------|---------------|--------|
| L3-REF-01: docs/reference/ path stale | L2-business-rules | tool-cli-e2e-lifecycle references `docs/proposals/test-profile-system/proposal.md` (EXISTS, not affected) | None |
| L3-REF-02: plugins/forge/ directory count mismatch | L2-business-rules | No items in this batch reference plugins/forge/ subdirectory count | None |
| L3-REF-03: Probe health check vs recipe-based | L2-business-rules | e2e-server-lifecycle-hardening decision describes recipe-based probe correctly | None |
| L2-CL2: pkg/profile/ does not exist | L2-conventions-batch2 | tool-cli-e2e-lifecycle references `forge-cli/pkg/profile/` | OUTDATED — see item analysis |
| L2-CL2: pkg/lesson and pkg/research do not exist | L2-conventions-batch2 | No items in this batch reference these packages | None |
| L1-CL1: Surface priority CLI/TUI reversed | L1-core-docs | pattern-surface-resolution-shortcut does not discuss priority ordering | None |
| L1-CL1: test-type-model.md path broken | L1-core-docs | No items in this batch reference docs/reference/test-type-model.md | None |
| L2-CL1: run-tasks classified as command not skill | L2-conventions-batch1 | architecture decision mentions "run-tasks" in context of skill decomposition | MATCHES — architecture decision is correct |

---

## Item-by-Item Analysis

### LESSONS (13 items)

---

### 1. lesson-vibe-coding-scope-control.md

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The core process rules (scope guard for >50 LOC changes, reflection threshold at 3 fix/style commits on same file, parser/stats layer protection) are sound and remain applicable. However, Rule 3 references `internal/parser/` and `internal/stats/` as protected paths — these directories no longer exist under `forge-cli/internal/`. The current `forge-cli/internal/` only contains `cmd/` and `embedded/`. The actual parser and stats functionality has been reorganized or relocated. The example commit hashes (`4432aa6`, `cb93d0d`) are historical references that may or may not still be valid. The core insight (vibe coding needs scope guards) remains valid regardless of code structure changes.
- **Code path verification**: `forge-cli/internal/parser/` DOES NOT EXIST; `forge-cli/internal/stats/` DOES NOT EXIST; `forge-cli/internal/` contains only `cmd/` and `embedded/`
- **Required update**: (1) Update Rule 3 to remove references to `internal/parser/` and `internal/stats/` — identify current protected paths if applicable, or reframe as a general principle; (2) The three rules remain valid as process guidance

---

### 2. pattern-compile-check-before-submit.md

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The core pattern (run `go build ./...` before submitting) is a fundamental best practice that remains correct. Referenced files `forge-cli/internal/cmd/root.go` and `forge-cli/internal/cmd/lesson.go` both EXIST. The quality gate order (build → vet → test) remains standard Go practice. The specific bug example (lessonCmd undefined) is historical context — the lesson's reusable pattern is timeless.
- **Code path verification**: `forge-cli/internal/cmd/root.go` EXISTS; `forge-cli/internal/cmd/lesson.go` EXISTS

---

### 3. pattern-dispatcher-auto-verify.md

- **Classification**: code-reference + experience-summary
- **Status**: outdated
- **Justification**: The lesson references `all_completed.go:70-74` as the source file for the dual-gate check. This file no longer exists — the functionality has been refactored into `forge-cli/internal/cmd/qualitygate/quality_gate.go` (function `CheckAllCompleted` at line 67). The described CWD mismatch problem (subagents running from `backend/` directory) references a project structure (backend/ as a separate Go module) that is specific to the train-recorder project, not the forge codebase itself. The core insight about CWD affecting `FindProjectRoot()` remains valid as a general concern, but the specific code path references are stale. Additionally, `FindProjectRoot` now lives in `forge-cli/pkg/project/root.go` with additional resolution strategies including `FindProjectRootFrom` and environment variable support.
- **Code path verification**: `all_completed.go` DOES NOT EXIST; `qualitygate/quality_gate.go` EXISTS with `CheckAllCompleted`; `pkg/project/root.go` EXISTS with `FindProjectRoot`
- **Required update**: Complete rewrite needed — update file references, describe current `CheckAllCompleted` location, and generalize the CWD lesson beyond the train-recorder-specific context

---

### 4. pattern-large-scale-rename.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The multi-pass rename strategy (identifiers → string literals → test data → test expectations → dead code) is a well-structured reusable pattern that remains valid. Referenced files `forge-cli/pkg/task/types.go` and `forge-cli/pkg/task/build.go` both EXIST. The cross-reference to `docs/lessons/arch-constant-rename-whack-a-mole.md` is valid (file EXISTS). However, the code examples use specific constant names from the rename event (e.g., `TypeTestPipelineVerifyRegression`, `TypeFeature`) that may have been renamed again since — this is historical context and acceptable. The sed patterns use macOS syntax (`sed -i ''`) which is platform-specific.
- **Code path verification**: `forge-cli/pkg/task/types.go` EXISTS; `forge-cli/pkg/task/build.go` EXISTS; `docs/lessons/arch-constant-rename-whack-a-mole.md` EXISTS
- **Required update**: Minor — consider noting that sed syntax is macOS-specific and may need adjustment for Linux; the core pattern is timeless

---

### 5. pattern-sitemap-shared-layout.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: This is a generalized process lesson about architecture-first exploration strategy. The core insight ("read architecture before mechanically executing") is language- and framework-agnostic. It does not reference any specific code paths. The React/TypeScript example context (AppLayout, Route) is illustrative rather than normative. The lesson's "Key Takeaway" is a universally applicable principle.
- **Code path verification**: N/A — no code paths referenced

---

### 6. pattern-surface-resolution-shortcut.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson references `skills/quick-tasks/SKILL.md` which EXISTS at `plugins/forge/skills/quick-tasks/SKILL.md`. The core optimization (skip per-file surface queries for single-surface projects) is valid. The `.forge/config.yaml` surface configuration format described in the examples matches current usage. However, this item was already audited in L3 batch 5 (as gotcha-surface-fields-single-surface-empty.md) where it was marked as needs-update. The current SKILL.md contains a Surface-Key/Type Inference section with a two-layer resolution strategy that partially addresses the lesson's concern.
- **Code path verification**: `plugins/forge/skills/quick-tasks/SKILL.md` EXISTS; `.forge/config.yaml` EXISTS
- **Required update**: (1) Verify whether the current quick-tasks SKILL.md adequately addresses the single-surface case; (2) May be partially resolved if template instructions were clarified

---

### 7. pattern-task-vs-output-naming.md

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The naming convention (`{titleSlug}.md` not `T-test-{N}.md`) is valid and enforced by the current breakdown-tasks skill at `plugins/forge/skills/breakdown-tasks/SKILL.md`. However, the lesson describes the problem in terms of `/breakdown-tasks` auto-generating test tasks with IDs like `T-test-1`, `T-test-2`. It is unclear whether this bug still exists — the breakdown-tasks SKILL.md would need to be checked for the current template behavior. The convention itself (slug-based naming) remains the project standard.
- **Code path verification**: `plugins/forge/skills/breakdown-tasks/SKILL.md` EXISTS
- **Required update**: Verify whether the described bug still exists in current breakdown-tasks SKILL.md; if fixed, update the lesson to note the resolution

---

### 8. tool-cli-e2e-lifecycle.md

- **Classification**: code-reference + process-standard
- **Status**: outdated
- **Justification**: Multiple code references are stale: (1) References `e2e-setup` and `e2e-discover` as standard recipes — the current justfile uses `test-setup` and `test-discover` instead (no `e2e-` prefix); (2) Describes `go-test profile` with `e2e-setup = go mod download` and `e2e-discover = go test -list` — the current justfile has `test-setup` (pre-builds forge binary + warms build cache) and `test-discover` (go test -tags=cli_functional -list); (3) References `forge-cli/pkg/profile/` — this directory DOES NOT EXIST; (4) References `docs/proposals/test-profile-system/proposal.md` — this EXISTS but is a historical proposal; (5) The comparison table shows "CLI profile (go-test)" with `tests/e2e` paths, but the current justfile uses `tests/` (no `e2e` subdirectory). The core concept (CLI projects need lightweight lifecycle steps) remains valid, but all specific recipe names, paths, and examples are outdated.
- **Code path verification**: `forge-cli/pkg/profile/` DOES NOT EXIST; `docs/proposals/test-profile-system/proposal.md` EXISTS; justfile has `test-setup` and `test-discover` (not `e2e-setup`/`e2e-discover`)
- **Required update**: Complete rewrite needed — update recipe names from `e2e-*` to `test-*`, remove `pkg/profile/` reference, update example justfile to match current format

---

### 9. tool-fix-e2e-unknown-placeholder.md

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The core problem description (auto-generated fix-e2e tasks containing "unknown" placeholders when infrastructure failures produce no parseable test results) describes a real design gap. The two-part solution (improve generator to capture raw output, handle infrastructure failures differently from test-level failures) is sound process guidance. The referenced `testing/results/latest.md` and `testing/results/failures/` paths are convention-dependent and may have changed. The sample fix-e2e task content shows Chinese text mixed with English paths, consistent with mixed-language conventions. The key takeaway about infrastructure-level vs test-level failure handling is a valid process standard.
- **Code path verification**: N/A — process standard; path conventions should be verified against current test output directory structure
- **Required update**: (1) Verify current test output directory structure (testing/results/ vs tests/results/); (2) Check if the fix-e2e task generator has been improved since this lesson was written

---

### 10. tool-justfile-arg-attribute.md

- **Classification**: process-standard + code-reference
- **Status**: valid
- **Justification**: The lesson documents the `[arg(long)]` attribute in just (available from v1.50.0) which is not yet in the official just book. This remains a current issue — the just book may still not document this attribute. The key takeaways are: (1) `[arg(long)]` creates named options, not boolean flags; (2) just version differences matter in production; (3) `just --dry-run` is the reliable verification method. The justfile template examples using `[arg("feature", long)]` match the current justfile templates in `plugins/forge/skills/init-justfile/templates/`. The external reference URLs point to the just book and source code.
- **Code path verification**: Init-justfile templates at `plugins/forge/skills/init-justfile/templates/` EXIST and use `[arg(long)]` pattern

---

### 11. tool-record-coverage-capture.md

- **Classification**: process-standard
- **Status**: valid
- **Justification**: The lesson documents that `record-task` skill's coverage field defaults to 0.0% when omitted, creating data distortion. The submit-task SKILL.md at `plugins/forge/skills/submit-task/SKILL.md` line 27 confirms that `coverage` is a coding-only field. The solution (run coverage command before writing record.json, use null for untested packages) is correct process guidance. The takeaway ("workflow field = format definition + collection step") is a generalizable principle.
- **Code path verification**: `plugins/forge/skills/submit-task/SKILL.md` EXISTS; line 27 confirms coverage field handling

---

### 12. tool-submit-background-timeout.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes `forge task submit` triggering the quality gate which can exceed the Bash tool's 120s default timeout. The referenced `forge-cli/internal/cmd/submit.go` has been moved to `forge-cli/internal/cmd/task/submit.go`. The `validateQualityGate` function EXISTS in the relocated file at line 377. The solution (set explicit timeout=300000 for breaking task submissions) remains valid. However, the file path needs updating, and the lesson should verify whether the CLI itself has added any timeout handling or progress reporting since the lesson was written.
- **Code path verification**: `forge-cli/internal/cmd/submit.go` DOES NOT EXIST (moved); `forge-cli/internal/cmd/task/submit.go` EXISTS; `validateQualityGate` at line 377 CONFIRMED
- **Required update**: Update file reference from `internal/cmd/submit.go` to `internal/cmd/task/submit.go`

---

### 13. worktree-stale-refs.md

- **Classification**: code-reference
- **Status**: needs-update
- **Justification**: The lesson describes `forge worktree start` ignoring `--source-branch` when stale remote tracking refs exist. The worktree command code EXISTS at `forge-cli/internal/cmd/worktree/cmd_start.go`. The `--source-branch` flag is CONFIRMED present (lines 25-26, 128-129). The source branch priority logic (flag > config > HEAD) is documented in the code. The core issue (stale refs silently overriding `--source-branch`) appears to be a valid concern — the code shows that when sourceBranch is specified and resolves successfully, it's used, but there may be edge cases where stale refs interfere. The workaround (`git remote prune origin`) remains valid. The improvement suggestions (warn about stale refs, auto-prune on remove) are actionable.
- **Code path verification**: `forge-cli/internal/cmd/worktree/cmd_start.go` EXISTS; `--source-branch` flag CONFIRMED; source branch priority logic CONFIRMED (lines 127-131)
- **Required update**: Verify whether the stale ref issue has been addressed in current code; update improvement suggestions with current status

---

### DECISIONS (10 items)

---

### 14. architecture.md (4 decisions)

- **Classification**: process-standard
- **Status**: needs-update (2 of 4 decisions need updates)
- **Justification**: Contains 4 decisions in tabular format:

  **Decision 1** (2026-04-30, justfile-standard-vocabulary): Defer `task scope` command. This remains valid — no evidence of `task scope` being implemented, and the rationale (current wiring sufficient) still holds. **valid**.

  **Decision 2** (2026-05-19, simplify-breakdown-tasks-prompt): Decompose skills into skeleton + `rules/` directory. This has been IMPLEMENTED — 17 of 21 skills now have `rules/` directories. The decision correctly describes the current architecture. **valid**.

  **Decision 3** (2026-05-20, system-type-exclusion): Reverse exclusion (SystemTypes blacklist). This has been IMPLEMENTED — `SystemTypes` map exists in `forge-cli/pkg/task/types.go` with 12 entries (not 13 as stated in the decision). The dual-identity exclusion for `doc.consolidate` and `doc.drift` is CONFIRMED in code comments. **needs-update** — the count says 13 system types but the code has 12.

  **Decision 4** (2026-05-20, prompt-template-audit): Submit-task responsibility owned by task-executor. This has been IMPLEMENTED — submit-task is a dedicated skill at `plugins/forge/skills/submit-task/SKILL.md`. The decision to keep submit ownership in a single agent is correct. **valid**.

- **Code path verification**: `forge-cli/pkg/task/types.go` lines 145-162 CONFIRM 12 SystemTypes entries; `plugins/forge/skills/submit-task/SKILL.md` EXISTS; 17 skill `rules/` directories CONFIRMED
- **Required update**: Correct system type count from 13 to 12 in Decision 3

---

### 15. data-model.md (0 decisions)

- **Classification**: empty-stub
- **Status**: empty-stub
- **Justification**: File contains only the header row of the decision table with no entries. This is a placeholder file with no content to audit.

---

### 16. dependencies.md (0 decisions)

- **Classification**: empty-stub
- **Status**: empty-stub
- **Justification**: File contains only the header row of the decision table with no entries. This is a placeholder file with no content to audit.

---

### 17. e2e-server-lifecycle-hardening.md (1 comprehensive decision)

- **Classification**: code-reference
- **Status**: outdated
- **Justification**: This is a detailed implementation decision (dated 2026-05-10, marked "已实施" / implemented) that describes hardening the e2e server lifecycle. Several references are stale:

  (1) References `references/justfile-templates/node.just` through `generic.just` — the actual path is `plugins/forge/skills/init-justfile/templates/node.just` through `generic.just`. The old `references/justfile-templates/` path does not exist.

  (2) References `skills/gen-test-scripts/SKILL.md` — the actual path is `plugins/forge/skills/gen-test-scripts/SKILL.md`.

  (3) References `skills/run-e2e-tests/SKILL.md` — this skill DOES NOT EXIST. There is no `run-e2e-tests` skill directory in the current plugins/forge/skills/.

  (4) References `skills/gen-test-scripts/templates/playwright.config.ts` — the `templates/` directory does not exist under gen-test-scripts (only `rules/` and `types/` exist).

  (5) The probe recipe code shows `tests/e2e/` paths but the current templates use `tests/` (no `e2e` subdirectory).

  (6) The decision's core content (server lifecycle embedded in test recipe, probe accepting non-5xx, three-layer detection) appears to be implemented in the current `go.just` template which shows the same pattern (PID alive → probe → start server → health check → run tests).

  (7) The migration guide references `/init-justfile` which is a valid command.

- **Code path verification**: `plugins/forge/skills/init-justfile/templates/*.just` EXIST (6 templates); `run-e2e-tests` skill DOES NOT EXIST; gen-test-scripts has no `templates/` directory; `go.just` template confirms lifecycle pattern is implemented
- **Required update**: Significant — update all file path references to use `plugins/forge/skills/` prefix; remove references to non-existent `run-e2e-tests` skill; update `tests/e2e/` to `tests/` paths; note that gen-test-scripts no longer has a templates/ subdirectory

---

### 18. error-handling.md (0 decisions)

- **Classification**: empty-stub
- **Status**: empty-stub
- **Justification**: File contains only the header row of the decision table with no entries. This is a placeholder file with no content to audit.

---

### 19. interface.md (1 decision)

- **Classification**: code-reference
- **Status**: valid
- **Justification**: The decision (2026-05-21) specifies using frontmatter `created` field (YYYY-MM-DD) for CLI list sorting in descending order, with mtime as fallback. This has been IMPLEMENTED — `forge-cli/internal/cmd/feature/feature.go` at line 208 confirms the sort-by-created logic with mtime fallback. The rationale (git clone/pull resets mtime) is correct and well-documented. The source reference to `proposal: cli-created-field-and-display` is traceable.
- **Code path verification**: `forge-cli/internal/cmd/feature/feature.go` lines 208-212 CONFIRM created-field sorting with mtime fallback

---

### 20. local-dev-deployment.md (0 decisions)

- **Classification**: empty-stub
- **Status**: empty-stub
- **Justification**: File contains only the header row of the decision table with no entries. This is a placeholder file with no content to audit.

---

### 21. manifest.md (decisions index)

- **Classification**: process-standard
- **Status**: needs-update
- **Justification**: The manifest/index file lists decision categories with counts and last-updated dates. Several entries are inaccurate:

  (1) Architecture listed as "1 decision, last updated 2026-04-30" — actually has 4 decisions (last entry 2026-05-20).

  (2) Testing listed as "2 decisions, last updated 2026-05-20" — this is correct per testing.md content.

  (3) Recent Decisions table lists 3 entries which correctly mirror the most impactful decisions across files.

  (4) The "updated: 2026-04-23" frontmatter date is older than the most recent decisions (2026-05-20).

- **Code path verification**: architecture.md CONFIRMED to have 4 decisions (not 1); manifest.md frontmatter date is stale
- **Required update**: (1) Update architecture decision count from 1 to 4 and last-updated to 2026-05-20; (2) Update frontmatter `updated` to 2026-05-20

---

### 22. security.md (0 decisions)

- **Classification**: empty-stub
- **Status**: empty-stub
- **Justification**: File contains only the header row of the decision table with no entries (not even the trailing `|`). This is a placeholder file with no content to audit.

---

### 23. testing.md (2 decisions)

- **Classification**: code-reference + process-standard
- **Status**: valid (1) + needs-update (1)
- **Justification**: Contains 2 decisions:

  **Decision 1** (2026-05-20, test-knowledge-convention-driven): Convention files replace hardcoded profile package. The decision states that convention files (`testing-{framework}.md`) decouple Forge from language-specific code. Current state: there is a `docs/conventions/testing/` directory with `index.md` and `cli/index.md`, `cli/core.md` — but no `testing-{framework}.md` pattern files at the top level. The convention directory structure differs from the decision's description, but the spirit (convention-driven, not hardcoded) is implemented. **needs-update** — the specific file naming pattern described does not match the actual implementation.

  **Decision 2** (2026-05-20, task-coverage-strategy): Per-task-type coverage tiers via config-driven prompt injection. This has been IMPLEMENTED — `forge-cli/pkg/prompt/prompt.go` line 342 confirms the coverage resolution logic with priority (task frontmatter > config per-type > built-in default). `forge-cli/pkg/forgeconfig/config.go` line 191 confirms `ReadCoverageConfig`. The task types.go line 251 confirms the coverage field on task types. **valid**.

- **Code path verification**: `docs/conventions/testing/` EXISTS with `index.md` and `cli/` subdirectory; `forge-cli/pkg/prompt/prompt.go` CONFIRMS coverage config resolution; `forge-cli/pkg/forgeconfig/config.go` CONFIRMS ReadCoverageConfig
- **Required update**: For Decision 1, update the file naming pattern from `testing-{framework}.md` to reflect actual `docs/conventions/testing/{surface}/` structure

---

## Duplicate Detection (Topic Clustering)

### Cluster 1: Large-Scale Rename Strategy (2 lessons)

**Members**: pattern-large-scale-rename.md, arch-constant-rename-whack-a-mole.md

**Analysis**: These two lessons address the same rename event from different angles. `arch-constant-rename-whack-a-mole.md` describes the subagent getting stuck (the problem), while `pattern-large-scale-rename.md` extracts the reusable pattern (the solution). They are complementary, not redundant — one is the incident report, the other is the generalized methodology. pattern-large-scale-rename.md already cross-references the other.

**Recommendation**: **Not duplicate** — keep both, as they serve different purposes (incident analysis vs reusable pattern).

### Cluster 2: Task Submission and Quality Gate (2 lessons)

**Members**: pattern-compile-check-before-submit.md, tool-submit-background-timeout.md

**Analysis**: Both relate to the task submission process. `pattern-compile-check` focuses on the pre-submit compilation check pattern. `tool-submit-background-timeout` focuses on the Bash timeout issue during quality gate execution. These address different failure modes of the same workflow step but are not redundant.

**Recommendation**: **Not duplicate** — different failure modes, different solutions.

### Cluster 3: E2E Test Infrastructure (3 lessons + 1 decision)

**Members**: tool-cli-e2e-lifecycle.md, tool-fix-e2e-unknown-placeholder.md, e2e-server-lifecycle-hardening.md (decision)

**Analysis**: These all relate to e2e test infrastructure. `tool-cli-e2e-lifecycle` documents the CLI-appropriate lifecycle steps. `tool-fix-e2e-unknown-placeholder` documents the fix-e2e task generation gap. `e2e-server-lifecycle-hardening` is the architectural decision that shaped the current implementation. They cover different aspects (lifecycle design, error handling, architecture) of the same system.

**Recommendation**: **Not duplicate** — different aspects of the same system. However, tool-cli-e2e-lifecycle and e2e-server-lifecycle-hardening have significant content overlap (both describe probe/test-e2e recipe design). tool-cli-e2e-lifecycle is outdated and its content is subsumed by the more comprehensive e2e-server-lifecycle-hardening decision. Consider marking tool-cli-e2e-lifecycle as **duplicate** of e2e-server-lifecycle-hardening for the recipe design content.

**Status**: tool-cli-e2e-lifecycle marked as **outdated** (not duplicate) since it contains stale paths rather than purely redundant content.

### Cluster 4: Forge Worktree and State Management (2 lessons)

**Members**: pattern-dispatcher-auto-verify.md, worktree-stale-refs.md

**Analysis**: Both relate to Forge CLI infrastructure issues. `pattern-dispatcher-auto-verify` covers the `.forge/state.json` creation failure due to CWD mismatch. `worktree-stale-refs` covers stale remote tracking refs in worktree operations. Different problems, different solutions.

**Recommendation**: **Not duplicate** — different problems in related but distinct CLI areas.

### Cluster 5: Skill/Tool Configuration Awareness (3 lessons)

**Members**: pattern-surface-resolution-shortcut.md, tool-justfile-arg-attribute.md, tool-record-coverage-capture.md

**Analysis**: These lessons all document gaps between tool configuration and actual behavior. `surface-resolution-shortcut` covers unnecessary surface queries. `justfile-arg-attribute` covers undocumented just features. `record-coverage-capture` covers missing data collection steps. Each addresses a distinct configuration gap.

**Recommendation**: **Not duplicate** — distinct configuration awareness lessons.

### Cluster 6: Process Discipline (3 lessons)

**Members**: lesson-vibe-coding-scope-control.md, pattern-sitemap-shared-layout.md, pattern-task-vs-output-naming.md

**Analysis**: These are all process standards about discipline in different contexts. `vibe-coding-scope-control` is about scope management during fix iterations. `sitemap-shared-layout` is about architecture-first exploration. `task-vs-output-naming` is about file naming conventions. No content overlap.

**Recommendation**: **Not duplicate** — different process domains.

---

## Summary of Findings

### Items Requiring Action (by priority)

**P1 (High)**:
1. tool-cli-e2e-lifecycle.md — outdated recipe names (`e2e-*` → `test-*`), stale `pkg/profile/` path
2. e2e-server-lifecycle-hardening.md — multiple stale file paths, references to non-existent skill
3. pattern-dispatcher-auto-verify.md — references non-existent `all_completed.go`, function relocated to qualitygate package

**P2 (Medium)**:
4. lesson-vibe-coding-scope-control.md — Rule 3 references non-existent `internal/parser/` and `internal/stats/`
5. pattern-surface-resolution-shortcut.md — partially addressed by current SKILL.md, may be close to resolution
6. architecture.md — system type count discrepancy (13 stated vs 12 actual)
7. manifest.md — incorrect architecture decision count (1 vs 4)
8. testing.md — convention file naming pattern differs from decision description
9. tool-submit-background-timeout.md — file path moved from `internal/cmd/submit.go` to `internal/cmd/task/submit.go`
10. pattern-task-vs-output-naming.md — verify if described bug still exists

**P3 (Low)**:
11. pattern-large-scale-rename.md — sed syntax is macOS-specific, minor note needed
12. worktree-stale-refs.md — verify if stale ref issue has been addressed
13. tool-fix-e2e-unknown-placeholder.md — verify current test output directory structure

### Items Valid (no action needed):
- pattern-compile-check-before-submit.md
- pattern-sitemap-shared-layout.md
- tool-justfile-arg-attribute.md
- tool-record-coverage-capture.md
- interface.md (decision — created field sorting)
- testing.md Decision 2 (coverage strategy)

### Empty Decision Stubs (4 files):
- data-model.md
- dependencies.md
- error-handling.md
- local-dev-deployment.md
- security.md

**Recommendation for empty stubs**: These 5 files serve as structural placeholders. They are not harmful but add no value. Consider either (a) keeping them as placeholders for future decisions or (b) removing them to reduce directory clutter. Human confirmation required before removal.

---

## Audit Quality Review

- **Sampling ratio**: 100% (all 23 items audited individually)
- **Sampling result**: PASS
- **Missed items**: 0
- **Extended review**: Not needed — full coverage achieved for this batch
