# Fix Tasks from Audit Findings

**Source**: `consolidated-report.md` — Global Doc-Code Consistency Audit
**Generated**: 2026-06-03
**Audit baseline**: commits `85421b10` through `062992a4`

---

## Task Templates Legend

| Template | Executor | Key Trait |
|----------|----------|-----------|
| **fix-type** | task-executor independent | Self-contained, no human checkpoint |
| **review-type** | task-executor pauses for human confirmation | Knowledge base deletion/merge requires human sign-off |
| **cross-layer-verification-type** | task-executor reads dependency report first | Depends on prior audit reports (cross-layer influence lists) |

---

## Immediate — P0 (Release-Blocking)

### FT-001 [fix-type] Fix CLAUDE.md Plugin Subdirectory List

**Priority**: P0
**Severity**: Critical — AI agents will create files at wrong paths

**Context**:
CLAUDE.md line 19 lists `plugins/forge/` subdirectories as "skills, commands, agents, hooks, references, scripts". Only 4 top-level directories actually exist: `agents/`, `commands/`, `hooks/`, `skills/`. The `references/` and `scripts/` directories do NOT exist as top-level directories under `plugins/forge/`. AI agents reading CLAUDE.md will create files at `plugins/forge/references/` or `plugins/forge/scripts/`, which are incorrect locations.

**Target file**: `CLAUDE.md` line 19

**Action**:
1. Open `CLAUDE.md`
2. Locate the line listing plugin subdirectories (around line 19, inside the `<MANDATORY>` block)
3. Change the subdirectory list from "skills, commands, agents, hooks, references, scripts" to "agents, commands, hooks, skills"
4. Add a clarification note that `references` and `scripts` are skill-internal subdirectories, not top-level plugin directories

**Verification**: Confirm only actual top-level directories under `plugins/forge/` are listed. Run `ls -d plugins/forge/*/` to verify.

**Self-contained**: Yes — single file, single line change, no dependencies on other fix tasks.

---

## First Fix Batch — P1 (High)

### FT-002 [fix-type] Update README.md Version Badge and Counts

**Priority**: P1
**Severity**: High — users and agents get incorrect version/feature model

**Context**:
README.md has three inaccuracies: (1) line 5 version badge shows 5.6.0 but plugin version is 3.0.0-rc.41; (2) line 418 states command count is 18 but actual count is 16; (3) line 393 lists `test.verify-regression` task type which does not exist.

**Target file**: `README.md`

**Actions**:
1. Line 5: Update version badge to match current plugin version in `plugins/forge/.claude-plugin/plugin.json`
2. Line 418: Change command count from 18 to 16
3. Line 393: Remove `test.verify-regression` from the type table

**Verification**:
- Run `grep '"version"' plugins/forge/.claude-plugin/plugin.json` to confirm correct version
- Run `ls plugins/forge/commands/*.md | wc -l` to confirm command count
- Grep codebase for `verify-regression` to confirm it does not exist

**Self-contained**: Yes.

---

### FT-003 [fix-type] Fix ARCHITECTURE.md Command Count and Broken Link

**Priority**: P1
**Severity**: High — incorrect architecture model and dead link

**Context**:
ARCHITECTURE.md has two issues: (1) line 33 states commands count is 18 but actual is 16; (2) lines 222 and 485 contain broken links to `docs/reference/test-type-model.md` which does not exist — the actual path is `plugins/forge/skills/test-guide/references/test-type-model.md`.

**Target file**: `docs/ARCHITECTURE.md`

**Actions**:
1. Line 33: Update command count from 18 to 16
2. Lines 222, 485: Replace `docs/reference/test-type-model.md` with the correct path `plugins/forge/skills/test-guide/references/test-type-model.md`

**Verification**:
- Run `ls plugins/forge/commands/*.md | wc -l` for command count
- Run `ls plugins/forge/skills/test-guide/references/test-type-model.md` to verify path exists

**Self-contained**: Yes.

---

### FT-004 [fix-type] Fix initialization.md Surface Priority Order

**Priority**: P1
**Severity**: High — incorrect documentation of surface detection priority

**Context**:
`docs/user-guide/initialization.md` lines 300-306 reverse the surface priority for CLI and TUI. The documentation lists them in wrong order relative to what the code actually implements.

**Target file**: `docs/user-guide/initialization.md`

**Actions**:
1. Locate the surface priority table around lines 300-306
2. Swap the CLI and TUI rows to match the actual code priority order
3. Verify the full priority list matches the code in `internal/surface/` or equivalent package

**Verification**: Compare the documented surface priority order against the code's surface detection logic.

**Self-contained**: Yes.

---

### FT-005 [fix-type] Update plugin.md Hook Types and Events

**Priority**: P1
**Severity**: High — incomplete hook documentation causes agents to miss hook capabilities

**Context**:
`docs/official-references/plugin.md` has two issues: (1) lines 141-145 list only 4 hook types but actual count is 5 — missing `mcp_tool`; (2) lines 111-138 hook events table is missing `UserPromptExpansion` and `PostToolBatch` events.

**Target file**: `docs/official-references/plugin.md`

**Actions**:
1. Lines 141-145: Add `mcp_tool` to the hook types list
2. Lines 111-138: Add `UserPromptExpansion` and `PostToolBatch` to the hook events table with correct descriptions

**Verification**:
- Grep codebase for `mcp_tool` hook type registration to confirm its existence
- Grep for `UserPromptExpansion` and `PostToolBatch` event names in hook-related code

**Self-contained**: Yes.

---

### FT-006 [fix-type] Fix surface-orchestration.md Regex and Paths

**Priority**: P1
**Severity**: High — incorrect regex prevents valid surface keys

**Context**:
`docs/conventions/surface-orchestration.md` line 62 documents surface-key regex as `[a-zA-Z0-9_-]` but the actual code regex is `^[a-z][a-z0-9-]*$` (lowercase only, must start with letter). Also line 42 references non-existent `docs/reference/test-type-model.md`.

**Target file**: `docs/conventions/surface-orchestration.md`

**Actions**:
1. Line 62: Update regex from `[a-zA-Z0-9_-]` to `^[a-z][a-z0-9-]*$`
2. Line 42: Update path reference from `docs/reference/test-type-model.md` to `plugins/forge/skills/test-guide/references/test-type-model.md`

**Verification**: Grep codebase for the actual surface-key regex pattern to confirm.

**Self-contained**: Yes.

---

### FT-007 [fix-type] Fix task-lifecycle.md Stale Path Reference

**Priority**: P1
**Severity**: High — dead reference prevents agents from finding test type model

**Context**:
`docs/conventions/task-lifecycle.md` line 42 references non-existent `docs/reference/test-type-model.md`. The actual file is at `plugins/forge/skills/test-guide/references/test-type-model.md`.

**Target file**: `docs/conventions/task-lifecycle.md`

**Actions**:
1. Line 42: Replace `docs/reference/test-type-model.md` with `plugins/forge/skills/test-guide/references/test-type-model.md`

**Verification**: Confirm the target path exists.

**Self-contained**: Yes.

---

### FT-008 [fix-type] Update quality-gate.md Surface-Aware Mode

**Priority**: P1
**Severity**: High — missing surface-aware phase documentation

**Context**:
`docs/business-rules/quality-gate.md` lines 17-19 omit the surface-aware mode probe and lifecycle details. The quality gate now has a surface-aware Phase 3 that probes surfaces before testing, but this is not documented.

**Target file**: `docs/business-rules/quality-gate.md`

**Actions**:
1. Add a description of the surface-aware Phase 3 after the existing phase descriptions
2. Include: how the surface probe works, which surfaces are probed, what happens when a probe fails
3. Reference the surface-orchestration convention for detailed probe behavior

**Verification**: Compare against actual quality gate implementation in `internal/cmd/qualitygate/`.

**Self-contained**: Yes.

---

### FT-009 [fix-type] Fix naming.md Non-Existent Package References

**Priority**: P1
**Severity**: High — references to non-existent packages mislead agents

**Context**:
`docs/conventions/naming.md` has three stale package references: (1) line 191 lists non-existent `pkg/version` package (version is in `pkg/types`); (2) lines 196-197 list non-existent `pkg/lesson` and `pkg/research` packages.

**Target file**: `docs/conventions/naming.md`

**Actions**:
1. Line 191: Remove the `pkg/version` row from the package table; update `pkg/types` description to note it includes version types
2. Lines 196-197: Remove the `pkg/lesson` and `pkg/research` rows from the package table

**Verification**:
- Run `ls internal/pkg/version/` to confirm it does not exist
- Run `ls internal/pkg/lesson/` and `ls internal/pkg/research/` to confirm they do not exist
- Run `grep -r "version" internal/pkg/types/` to confirm version types location

**Self-contained**: Yes.

---

### FT-010 [fix-type] Update package-organization.md Counts and Missing Subpackage

**Priority**: P1
**Severity**: High — stale counts and missing documentation for subpackage

**Context**:
`docs/conventions/package-organization.md` has two issues: (1) line 40 stale subpackage count (says 15 but actual is 18 top-level files); (2) lines 37-39 missing `qualitygate/` subpackage documentation. Also `docs/` subpackage was removed but count not updated.

**Target file**: `docs/conventions/package-organization.md`

**Actions**:
1. Line 40: Update the subpackage count to match current actual count
2. Lines 37-39: Add `qualitygate/` subpackage to the subpackage table with description of its contents
3. Verify and update the overall composition description if `docs/` subpackage was removed

**Verification**: Run `ls internal/cmd/` and count subdirectories to get actual count.

**Self-contained**: Yes.

---

### FT-011 [fix-type] Fix skill-structure.md Constraint Violations

**Priority**: P1
**Severity**: High — documented constraint does not match reality

**Context**:
`docs/conventions/skill-structure.md` has two issues: (1) line 38 states a constraint violated — 6 skills use auxiliary directories beyond `rules/` and `templates/` (actual auxiliary directories include `experts/`, `rubrics/`, `references/`, `rules/`, `data/`, `types/`); (2) line 37 notes 5 SKILL.md files exceed the 350-line limit (max: 504 lines).

**Target file**: `docs/conventions/skill-structure.md`

**Actions**:
1. Line 38: Update the constraint to document the allowed auxiliary directory types: `rules/`, `templates/`, `experts/`, `rubrics/`, `references/`, `data/`, `types/`
2. Line 37: Either (a) update the line limit to accommodate current maximum (e.g., 500 lines) with rationale, or (b) flag the 5 oversized SKILL.md files for content extraction into rules/ and templates/ files as a separate follow-up action
3. Document the extraction recommendation if approach (b) is chosen

**Verification**: Run `find plugins/forge/skills/ -maxdepth 2 -type d | sort` to list all actual auxiliary directories. Run `wc -l plugins/forge/skills/*/SKILL.md | sort -rn | head -6` to confirm line counts.

**Self-contained**: Yes.

---

### FT-012 [fix-type] Fix testing/index.md Missing Surface Types

**Priority**: P1
**Severity**: High — incomplete testing documentation

**Context**:
`docs/conventions/testing/index.md` lines 10-12 only lists `cli` surface but omits `web`, `api`, `tui`, and `mobile`. Also states `tests/cli/` as test directory but it does not exist.

**Target file**: `docs/conventions/testing/index.md`

**Actions**:
1. Lines 10-12: Add all 5 surface types (`cli`, `web`, `api`, `tui`, `mobile`) to the surface list
2. Update the test directory pattern to reflect the actual test structure (`tests/<journey>/` instead of `tests/cli/`)

**Verification**: Run `ls tests/` to confirm actual test directory structure. Run `grep -r "surface" internal/` for surface type definitions.

**Self-contained**: Yes.

---

## Second Fix Batch — P2 (Medium)

### FT-013 [fix-type] Fix ARCHITECTURE.md Quality Gate and Eval References

**Priority**: P2
**Severity**: Medium — imprecise descriptions increase debugging cost

**Context**:
ARCHITECTURE.md has two issues: (1) lines 370-371 describe FullGateSequence imprecisely; the code uses a three-step process but documentation is vague; (2) line 101 references "T-eval-doc" task type which does not exist.

**Target file**: `docs/ARCHITECTURE.md`

**Actions**:
1. Lines 370-371: Update FullGateSequence description to accurately reflect the three-step quality gate process
2. Line 101: Remove or correct the "T-eval-doc" reference

**Verification**: Compare against quality gate implementation in `internal/cmd/qualitygate/`.

**Self-contained**: Yes.

---

### FT-014 [fix-type] Fix architecture-overview.md Quality Gate Description

**Priority**: P2
**Severity**: Medium — uses generic "test" instead of specific "unit-test" step

**Context**:
`docs/official-references/architecture-overview.md` line 134 describes quality gate using generic "test" instead of the actual "unit-test" step. Also line 237 lists non-existent `docs/reference/` directory.

**Target file**: `docs/official-references/architecture-overview.md`

**Actions**:
1. Line 134: Change "test" to "unit-test" to match actual quality gate step name
2. Line 237: Remove reference to `docs/reference/` directory or update to correct location

**Verification**: Grep quality gate code for actual step names.

**Self-contained**: Yes.

---

### FT-015 [fix-type] Fix initialization.md YAML Example Contradiction

**Priority**: P2
**Severity**: Medium — example contradicts defaults table

**Context**:
`docs/user-guide/initialization.md` lines 152-153 YAML example shows `full: true` but the defaults table says `full: false`.

**Target file**: `docs/user-guide/initialization.md`

**Actions**:
1. Lines 152-153: Align the YAML example with the defaults table (either update example to show `full: false` or update the defaults table with clarification that `full` defaults to `false` but the example shows an override)

**Verification**: Check code default value for the `full` configuration option.

**Self-contained**: Yes.

---

### FT-016 [fix-type] Fix skills-ref.md Supporting File Types

**Priority**: P2
**Severity**: Medium — incomplete file type list

**Context**:
`docs/official-references/skills-ref.md` lines 106-117 list supporting file types for Forge plugin but miss `rules/`, `rubrics/`, `experts/`, `types/`, and `data/` directory types.

**Target file**: `docs/official-references/skills-ref.md`

**Actions**:
1. Lines 106-117: Add missing file types to the supporting file types table: `rules/`, `rubrics/`, `experts/`, `types/`, `data/`

**Verification**: Run `find plugins/forge/skills/ -maxdepth 2 -type d | sort` to list all actual directory types.

**Self-contained**: Yes.

---

### FT-017 [fix-type] Fix plugin.md Standard Layout and Agent Fields

**Priority**: P2
**Severity**: Medium — incomplete standard layout documentation

**Context**:
`docs/official-references/plugin.md` has two issues: (1) lines 522-553 standard layout omits `userConfig` and `channels` fields; (2) line 70 Forge agent uses non-standard frontmatter fields (`color`, `memory`, `inputs`) that are not documented.

**Target file**: `docs/official-references/plugin.md`

**Actions**:
1. Lines 522-553: Add `userConfig` and `channels` to the standard layout description
2. Line 70: Document the non-standard frontmatter fields used by Forge agent (`color`, `memory`, `inputs`) or note them as Forge-specific extensions

**Verification**: Check `plugins/forge/.claude-plugin/plugin.json` and `plugins/forge/agents/task-executor.md` for actual field usage.

**Self-contained**: Yes.

---

### FT-018 [fix-type] Fix surface-orchestration.md Probe and Teardown Descriptions

**Priority**: P2
**Severity**: Medium — behavior description does not match code

**Context**:
`docs/conventions/surface-orchestration.md` has multiple inaccuracies: (1) line 25 describes probe as HTTP endpoint checking but code is recipe-based; (2) line 54 describes teardown idempotency as PID-based but code is recipe-based; (3) line 40 states probe failure exit code is 1 but hook exits with code 0; (4) line 47 omits `mobile` from probe hard-gate scope (only lists web, api); (5) line 30 mobile sequence omits optional `test-setup` step.

**Target file**: `docs/conventions/surface-orchestration.md`

**Actions**:
1. Line 25: Update probe behavior description from "HTTP endpoint checking" to "recipe-based probing"
2. Line 54: Update teardown idempotency description from "PID-based" to "recipe-based"
3. Line 40: Correct probe failure exit code from 1 to 0
4. Line 47: Add `mobile` to the probe hard-gate scope list
5. Line 30: Add optional `test-setup` step to mobile sequence

**Verification**: Read relevant probe and teardown code in `internal/surface/` or `internal/cmd/` to confirm actual behavior.

**Self-contained**: Yes.

---

### FT-019 [fix-type] Fix constants.md Stale References

**Priority**: P2
**Severity**: Medium — stale line references and non-existent constant

**Context**:
`docs/conventions/constants.md` has three issues: (1) line 30 stale line reference for gitignore entry (init.go:46 should be init.go:43); (2) line 31 `defaultHealthPath` constant does not exist — path is inline default; (3) lines 24-33 "All extracted" claim contradicted by testrunner literals.

**Target file**: `docs/conventions/constants.md`

**Actions**:
1. Line 30: Update line reference from `init.go:46` to `init.go:43` (or current actual line)
2. Line 31: Remove `defaultHealthPath` entry or update to reflect it is an inline default, not a named constant
3. Lines 24-33: Add a note that some constants may exist as inline defaults in testrunner code, not as extracted constants

**Verification**: Run `grep -n "gitignore" internal/cmd/init.go` and `grep -rn "defaultHealthPath" internal/` to confirm.

**Self-contained**: Yes.

---

### FT-020 [fix-type] Fix forge-cli-reference.md Quality Gate Path

**Priority**: P2
**Severity**: Medium — incorrect source file path

**Context**:
`docs/conventions/forge-cli-reference.md` line 22 lists incorrect quality-gate source file path, missing the `qualitygate/` subpackage.

**Target file**: `docs/conventions/forge-cli-reference.md`

**Actions**:
1. Line 22: Update quality-gate source file path to include the `qualitygate/` subpackage (e.g., `internal/cmd/qualitygate/` instead of `internal/cmd/`)

**Verification**: Run `ls internal/cmd/qualitygate/` to confirm subpackage exists.

**Self-contained**: Yes.

---

### FT-021 [fix-type] Fix forge-distribution.md Hooks Directory Tree

**Priority**: P2
**Severity**: Medium — incomplete directory tree

**Context**:
`docs/conventions/forge-distribution.md` lines 48-56 hooks directory tree omits `run-hook.cmd` and `debug` entries.

**Target file**: `docs/conventions/forge-distribution.md`

**Actions**:
1. Lines 48-56: Add `run-hook.cmd` and `debug` to the hooks directory tree

**Verification**: Run `ls plugins/forge/hooks/` to confirm actual file listing.

**Self-contained**: Yes.

---

### FT-022 [fix-type] Fix package-organization.md Stale Deviation

**Priority**: P2
**Severity**: Medium — deviation describes non-existent subpackage

**Context**:
`docs/conventions/package-organization.md` line 121 deviation D7 describes non-existent `cmd/docs` subpackage which has already been removed.

**Target file**: `docs/conventions/package-organization.md`

**Actions**:
1. Line 121: Remove or update deviation D7 since `cmd/docs` subpackage no longer exists

**Verification**: Run `ls internal/cmd/docs/` to confirm it does not exist.

**Self-contained**: Yes.

---

### FT-023 [fix-type] Fix prompt-template-hierarchy.md Unused Tag

**Priority**: P2
**Severity**: Medium — documents a tag that is never used

**Context**:
`docs/conventions/prompt-template-hierarchy.md` lines 30-33 define `TASK-CONSTRAINTS` tag but it is not used in any template.

**Target file**: `docs/conventions/prompt-template-hierarchy.md`

**Actions**:
1. Lines 30-33: Either remove the `TASK-CONSTRAINTS` tag definition or add a note that it is reserved for future use
2. If removing, verify no template references it

**Verification**: Run `grep -r "TASK-CONSTRAINTS" plugins/forge/` to confirm no usage.

**Self-contained**: Yes.

---

### FT-024 [fix-type] Fix testing/cli/core.md Stale References

**Priority**: P2
**Severity**: Medium — references non-existent test directory and incorrect build tag scope

**Context**:
`docs/conventions/testing/cli/core.md` has two issues: (1) line 12 `tests/cli/` does not exist as a test directory; (2) line 14 `cli_functional` build tag is used across all journey tests, not CLI-specific.

**Target file**: `docs/conventions/testing/cli/core.md`

**Actions**:
1. Line 12: Update test directory reference from `tests/cli/` to actual test directory structure
2. Line 14: Clarify that `cli_functional` build tag applies to all journey tests, not just CLI-specific tests

**Verification**: Run `ls tests/` for actual directory structure. Run `grep -r "cli_functional" tests/` for build tag usage.

**Self-contained**: Yes.

---

### FT-025 [fix-type] Fix naming.md Deviation N3 and surface-cli.md Quote Style

**Priority**: P2
**Severity**: Medium — deviation description contradicts code reality

**Context**:
Two issues: (1) `docs/conventions/naming.md` lines 244-248 deviation N3 describes a redundant prefix but code still uses it; (2) `docs/official-references/surface-cli.md` line 19 error message quote style differs (single quotes vs backticks).

**Target files**: `docs/conventions/naming.md`, `docs/official-references/surface-cli.md`

**Actions**:
1. `naming.md` lines 244-248: Update deviation N3 — if code still uses the prefix, remove the "redundant" characterization or update to reflect current status
2. `surface-cli.md` line 19: Align error message quote style to match code convention (update single quotes to backticks if code uses backticks, or vice versa)

**Verification**: Grep for the prefix in code. Check error message format in CLI code.

**Self-contained**: Yes.

---

## Third Fix Batch — P3 (Low / Deferred)

### FT-026 [fix-type] Fix README.md Minor Inaccuracies

**Priority**: P3
**Severity**: Low — minor inaccuracies, no functional impact

**Context**:
README.md has three low-severity issues: (1) line 159 `task validate-index` should be `task validate`; (2) lines 261-265 missing `--json` flag in surfaces flags reference; (3) line 379 test category header says "5 types" but should say "4 types".

**Target file**: `README.md`

**Actions**:
1. Line 159: Change `task validate-index` to `task validate`
2. Lines 261-265: Add `--json` flag to surfaces flags reference
3. Line 379: Change "5 types" to "4 types" (or correct count based on current code)

**Verification**: Run `forge task --help` for correct command name. Run `forge surfaces --help` for flags.

**Self-contained**: Yes.

---

### FT-027 [fix-type] Fix ARCHITECTURE.md Missing init-justfile Subsystem

**Priority**: P3
**Severity**: Low — undocumented subsystem

**Context**:
ARCHITECTURE.md does not document the `init-justfile` skill in v3.0.0 subsystems.

**Target file**: `docs/ARCHITECTURE.md`

**Actions**:
1. Add `init-justfile` skill to the subsystems section with brief description of its purpose (project initialization with Justfile templates)

**Verification**: Confirm `init-justfile` skill exists at `plugins/forge/skills/init-justfile/SKILL.md`.

**Self-contained**: Yes.

---

### FT-028 [fix-type] Fix hooks guide.md Stale Reference

**Priority**: P3
**Severity**: Low — references non-existent directory

**Context**:
`plugins/forge/hooks/guide.md` line 10 references non-existent `docs/reference/` directory.

**Target file**: `plugins/forge/hooks/guide.md`

**Actions**:
1. Line 10: Remove or update reference to `docs/reference/` directory

**Verification**: Confirm `docs/reference/` does not exist.

**Self-contained**: Yes.

---

### FT-029 [fix-type] Fix L2 P3 Conventions Issues

**Priority**: P3
**Severity**: Low — minor documentation gaps

**Context**:
Multiple L2 P3 issues: (1) `quality-gate.md` line 17 ambiguous "just unit-test" description; (2) `quality-gate.md` line 19 missing surface-specific test recipes; (3) `forge-distribution.md` lines 79-84 hooks table omits SessionEnd and SubagentStop; (4) `forge-distribution.md` lines 188-192 run-tasks misclassified as skill in pipeline diagrams; (5) `constants.md` line 169 ANSI code deviation status inconsistency; (6) `code-structure.md` missing validate_index.go rename in CS-2 deviation.

**Target files**: `docs/business-rules/quality-gate.md`, `docs/conventions/forge-distribution.md`, `docs/conventions/constants.md`, `docs/conventions/code-structure.md`

**Actions**:
1. `quality-gate.md`: Clarify "just unit-test" has probe chain fallback; add surface-specific test recipes description
2. `forge-distribution.md`: Add SessionEnd and SubagentStop to hooks table; fix run-tasks classification (command, not skill) in pipeline diagrams
3. `constants.md`: Align ANSI code deviation status with enum-constants.md
4. `code-structure.md`: Add validate_index.go rename to CS-2 deviation

**Verification**: Cross-reference each fix against the relevant code.

**Self-contained**: Yes.

---

### FT-030 [fix-type] Fix L2 P3 Minor Issues

**Priority**: P3
**Severity**: Low — content quality improvements

**Context**:
Remaining L2 P3 issues: (1) `testing/cli/index.md` minimal content — only contains link to core.md; (2) `skill-self-containment.md` very brief — lacks concrete examples; (3) `package-organization.md` lines 28-30 non-command files (`output.go`, `styles.go`) in cmd/ top level.

**Target files**: `docs/conventions/testing/cli/index.md`, `docs/conventions/skill-self-containment.md`, `docs/conventions/package-organization.md`

**Actions**:
1. `testing/cli/index.md`: Expand with overview content for CLI testing conventions
2. `skill-self-containment.md`: Add concrete examples of self-contained vs non-self-contained skill definitions
3. `package-organization.md`: Document the non-command files in cmd/ top level (output.go, styles.go) with rationale or flag as deviation

**Verification**: Read each file after modification to confirm content is complete.

**Self-contained**: Yes.

---

## Knowledge Base — Review-Type Tasks (Human Confirmation Required)

> **HUMAN CONFIRMATION REQUIRED**: The following tasks involve deletion, archival, or merging of knowledge base entries. Per project policy, these actions cannot be executed automatically. The task-executor will pause at the confirmation checkpoint and wait for human sign-off before proceeding.
>
> **SLA**: Proposal author checks daily. If no response within 3 working days, escalate to P1 and send reminder via project collaboration channel.

### RT-001 [review-type] Archive Outdated Resolved-Bug Lessons (12 items)

**Priority**: P2 (escalates to P1 after 3 working days without confirmation)
**Human Confirmation Required**: YES — deletion of knowledge base entries

**Context**:
The following 12 lessons describe bugs or issues that have been fully resolved. The code references in these entries no longer exist or have been restructured. Keeping them creates noise in the knowledge base.

**Proposed action**: Archive (move to `docs/lessons/_archived/`) the following items:

| # | File | Reason |
|---|------|--------|
| 1 | `gotcha-drift-detection-task-runtime.md` | Empty template problem fixed |
| 2 | `gotcha-duplicate-test-runs.md` | All code references restructured |
| 3 | `gotcha-embedded-template-name-mismatch.md` | Dot-to-hyphen conversion implemented |
| 4 | `gotcha-e2e-skill-monorepo-path-mismatch.md` | `run-e2e-tests` skill removed |
| 5 | `gotcha-e2e-test-quality-antipatterns.md` | `tests/e2e/` structure reorganized |
| 6 | `gotcha-eval-prd-use-zcode-agents.md` | Subagent types not used in current implementation |
| 7 | `gotcha-eval-subagent-type.md` | Duplicate of above, also outdated |
| 8 | `gotcha-fix-task-dependency-chain.md` | SourceTaskID map-key bug fixed |
| 9 | `gotcha-forge-task-index-per-type-duplicate.md` | Example feature directory removed |
| 10 | `gotcha-graduation-dual-module-drift.md` | Graduation/staging system removed |
| 11 | `gotcha-go-test-staging-graduation-friction.md` | Graduation system removed |
| 12 | `gotcha-hook-unbounded-test-timeout.md` | Node.js/Playwright e2e infrastructure removed |

**Execution Steps**:
1. For each file, read its content and verify it is indeed about a resolved issue
2. Create `docs/lessons/_archived/` directory if it does not exist
3. Move each file to `docs/lessons/_archived/`
4. **CHECKPOINT**: Pause and present the list to human for confirmation before executing moves

**Self-contained**: Yes — each item can be evaluated independently.

---

### RT-002 [review-type] Archive Outdated Stale Reference Lessons (9 items)

**Priority**: P2 (escalates to P1 after 3 working days without confirmation)
**Human Confirmation Required**: YES — deletion of knowledge base entries

**Context**:
The following 9 lessons contain stale test/tool references. The tools, commands, or test structures they describe no longer exist in the codebase.

**Proposed action**: Archive (move to `docs/lessons/_archived/`) the following items:

| # | File | Reason |
|---|------|--------|
| 1 | `gotcha-quick-tasks-no-commit.md` | Commit step added to quick-tasks SKILL.md |
| 2 | `gotcha-quick-tasks-stale-detect-command.md` | `forge test detect` reference removed |
| 3 | `gotcha-redundant-manual-e2e-verification.md` | References old `task-cli/` directory |
| 4 | `gotcha-task-executor-never-returns.md` | Termination constraint fully implemented |
| 5 | `gotcha-task-executor-auto-claim.md` | References outdated `zcode:` namespace |
| 6 | `gotcha-task-type-documentation-vs-doc.md` | Template bug fixed |
| 7 | `gotcha-test-pipeline-no-languages.md` | `interfaces` config replaced by `surfaces` |
| 8 | `gotcha-test-script-staging-vs-graduation.md` | Staging/graduation system removed |
| 9 | `gotcha-journey-hallucination-revision-death-spiral.md` | References `tests/e2e/` and `docs/reference/` |

**Execution Steps**:
1. For each file, read its content and verify it references non-existent code paths
2. Move each file to `docs/lessons/_archived/`
3. **CHECKPOINT**: Pause and present the list to human for confirmation before executing moves

**Self-contained**: Yes.

---

### RT-003 [review-type] Archive Outdated Wrong-Project Lesson (1 item)

**Priority**: P2 (escalates to P1 after 3 working days without confirmation)
**Human Confirmation Required**: YES — deletion of knowledge base entry

**Context**:
`gotcha-api-no-api-prefix.md` references `backend/` and `frontend/` directories that do not exist in this repository. This lesson belongs to a different project entirely.

**Proposed action**: Delete (not archive) this file as it does not belong in this repository.

**Execution Steps**:
1. Read the file to confirm it references directories not in this repo
2. **CHECKPOINT**: Present to human for confirmation
3. Delete the file

**Self-contained**: Yes.

---

### RT-004 [review-type] Merge Duplicate Knowledge Base Entries (3 pairs)

**Priority**: P2 (escalates to P1 after 3 working days without confirmation)
**Human Confirmation Required**: YES — merging knowledge base entries

**Context**:
The audit identified 16 items as potential duplicates. After cross-batch review, only 3 pairs are formally recommended for merge. The recommended action is to keep the more complete version and archive the duplicate.

**Proposed merges**:

| Duplicate (Archive) | Primary (Keep) | Reason |
|---------------------|---------------|--------|
| `gotcha-eval-subagent-type.md` | `gotcha-eval-prd-use-zcode-agents.md` | PRD version has full architectural analysis |
| `gotcha-graduation-dual-module-drift.md` | `gotcha-go-test-staging-graduation-friction.md` | Staging friction is more comprehensive |
| `arch-prototype-navigation-contract.md` | `arch-forge-skill-gap-analysis.md` | Gap analysis contains all navigation contract info plus additional proposals |

**Execution Steps**:
1. For each pair, read both files side-by-side
2. Verify the primary is indeed more complete
3. Merge any unique content from duplicate into primary (if any)
4. **CHECKPOINT**: Present merge plan to human for confirmation
5. Move duplicate files to `docs/lessons/_archived/`

**Self-contained**: Yes — each merge pair can be evaluated independently.

---

### RT-005 [review-type] Decide Fate of Empty Decision Stubs (5 items)

**Priority**: P3 (escalates to P1 after 3 working days without confirmation)
**Human Confirmation Required**: YES — deletion of decision files

**Context**:
The following decision files contain only header rows with no actual decisions:

1. `docs/decisions/data-model.md`
2. `docs/decisions/dependencies.md`
3. `docs/decisions/error-handling.md`
4. `docs/decisions/local-dev-deployment.md`
5. `docs/decisions/security.md`

**Options**:
- **A**: Keep as placeholders for future decisions
- **B**: Remove to reduce directory clutter

**Execution Steps**:
1. Read each file to confirm it is empty (header only)
2. **CHECKPOINT**: Present options A or B to human for decision
3. Execute chosen option

**Self-contained**: Yes.

---

### RT-006 [review-type] Archive Additional Outdated Items from Batches 5-6 (9 items)

**Priority**: P2 (escalates to P1 after 3 working days without confirmation)
**Human Confirmation Required**: YES — archival of knowledge base entries

**Context**:
9 additional lessons identified in L3 batches 5-6 have all code references stale or describe removed subsystems. These are items 29-37 from the consolidated report's outdated items list.

**Proposed action**: Archive these items after individual verification.

**Execution Steps**:
1. For each item, read its content and verify all code references are stale
2. Move to `docs/lessons/_archived/`
3. **CHECKPOINT**: Present the full list with verification evidence to human for confirmation

**Self-contained**: Yes.

---

### RT-007 [review-type] Archive Outdated Decisions (6 items)

**Priority**: P2 (escalates to P1 after 3 working days without confirmation)
**Human Confirmation Required**: YES — archival of decision entries

**Context**:
6 individual decisions across `architecture.md`, `manifest.md`, `testing.md`, and other decision files contain stale references or counts that no longer reflect reality. These are items 23-28 from the consolidated report's outdated items list.

**Proposed action**: Mark these decisions as superseded with a note explaining what changed, or archive them if the entire decision is no longer relevant.

**Execution Steps**:
1. Read each decision entry
2. Determine if the decision is partially valid (mark as superseded) or fully obsolete (archive)
3. **CHECKPOINT**: Present recommendations to human for confirmation
4. Apply approved changes

**Self-contained**: Yes.

---

## Cross-Layer Verification Tasks

> These tasks depend on the cross-layer influence lists from the consolidated audit report (Section "Cross-Layer Verification"). The task-executor must read the relevant source report before executing.

### CV-001 [cross-layer-verification-type] Verify docs/reference/ Path References Across All Layers

**Priority**: P2
**Dependency**: Consolidated report Section "Cross-Layer Verification" — `docs/reference/` path stale cross-reference

**Context**:
Multiple documents across L1, L2, and L3 reference `docs/reference/` as a directory, but it has never existed. The actual reference files are at `plugins/forge/skills/test-guide/references/`. The consolidated report identified 3 L3 lessons affected and 3 L2 items with stale references.

**Cross-layer influence list** (from consolidated report):
- L2-business-rules: All lessons referencing `docs/reference/` — 3 lessons affected (marked needs-update/outdated)
- L1-core-docs: `test-type-model.md` path broken — 2 lessons affected
- L2-conventions-batch1: Lessons referencing `docs/reference/` — additional items may exist

**Actions**:
1. Read the consolidated report cross-layer verification section for `docs/reference/` entries
2. For each affected L3 lesson: verify the stale reference and update the path to the correct location (or mark for archival if the entire lesson is about the non-existent path)
3. Verify all L1/L2 fixes (from FT-003, FT-006, FT-007) have updated `docs/reference/` references to the correct path
4. Output a verification summary listing all corrected paths

**Verification**: After all referenced fixes are applied, run `grep -r "docs/reference/" docs/` to confirm no stale references remain.

**Self-contained**: Yes — includes all context needed. Does not depend on other fix tasks completing first, but should be executed after or in parallel with FT-003, FT-006, FT-007.

---

### CV-002 [cross-layer-verification-type] Verify Plugin Subdirectory References Across Layers

**Priority**: P2
**Dependency**: Consolidated report Section "Cross-Layer Verification" — `plugins/forge/` subdirectory count cross-reference

**Context**:
The P0 CLAUDE.md fix (FT-001) corrects the top-level plugin subdirectory list. The consolidated report identified 5 L3 lessons that reference `skills/run-tasks/SKILL.md` which should be `commands/run-tasks.md`.

**Cross-layer influence list** (from consolidated report):
- L2-business-rules: Lessons referencing plugin directory structure — 5 lessons reference `skills/run-tasks/SKILL.md` (should be `commands/run-tasks.md`)

**Actions**:
1. Read the consolidated report cross-layer verification section for plugin subdirectory entries
2. For each affected L3 lesson: update `skills/run-tasks/SKILL.md` to `commands/run-tasks.md` (or mark for archival if the lesson is about the old skill structure)
3. Verify FT-001 (CLAUDE.md fix) has been applied correctly
4. Output a verification summary

**Verification**: Run `grep -r "skills/run-tasks" docs/lessons/` to check for remaining stale references.

**Self-contained**: Yes.

---

### CV-003 [cross-layer-verification-type] Verify Test Path Restructuring Impact Across Layers

**Priority**: P2
**Dependency**: Consolidated report Section "Cross-Layer Verification" — `tests/e2e/` reorganized cross-reference

**Context**:
The test infrastructure was reorganized from `tests/e2e/features/<slug>/` to `tests/<journey>/`. The consolidated report identified 15+ lessons referencing non-existent `tests/e2e/` paths, plus 2 L2 items with stale test path references.

**Cross-layer influence list** (from consolidated report):
- L2-conventions-batch2: Lessons referencing old test paths — 15+ lessons reference non-existent `tests/e2e/` paths
- Additional: Lessons about CLI test organization — 2 lessons reference non-existent `tests/cli/` path

**Actions**:
1. Read the consolidated report cross-layer verification section for test path restructuring entries
2. For each affected L3 lesson: determine if the lesson is about the old test structure (archive) or still relevant but needs path update (update paths from `tests/e2e/` to `tests/<journey>/`)
3. Verify L2 convention fixes (FT-012, FT-024) have updated test path references
4. Output a verification summary categorizing lessons as: archived, updated, or confirmed-valid

**Verification**: Run `grep -r "tests/e2e" docs/lessons/` and `grep -r "tests/cli/" docs/` to check for remaining stale references.

**Self-contained**: Yes.

---

### CV-004 [cross-layer-verification-type] Verify Quality Gate Flow Consistency Across Layers

**Priority**: P2
**Dependency**: Consolidated report Section "Cross-Layer Verification" — quality gate flow differs from docs

**Context**:
L1 findings show quality gate documentation is imprecise. L3 has 4 lessons describing evolved quality gate behavior. Cross-layer verification needed to ensure all quality gate documentation and knowledge base entries are consistent.

**Cross-layer influence list** (from consolidated report):
- L1-core-docs: Lessons about quality gate behavior — 4 lessons describe evolved behavior

**Actions**:
1. Read the consolidated report cross-layer verification section for quality gate entries
2. For each affected L3 lesson: verify whether the lesson describes the current (v3.0.0) quality gate flow or an outdated version
3. Verify L1 fixes (FT-008, FT-013, FT-014) have updated quality gate descriptions
4. Output a verification summary

**Verification**: Cross-reference quality gate code in `internal/cmd/qualitygate/` with all updated documentation and knowledge base entries.

**Self-contained**: Yes.

---

### CV-005 [cross-layer-verification-type] Verify Prompt Template Path Changes Across Layers

**Priority**: P2
**Dependency**: Consolidated report Section "Cross-Layer Verification" — `prompt/data/` renamed to `prompt/templates/`

**Context**:
`prompt/data/` was renamed to `prompt/templates/`. The consolidated report identified 3 L3 lessons affected by this path change.

**Cross-layer influence list** (from consolidated report):
- L2-conventions-batch1: Lessons referencing prompt template paths — 3 lessons affected

**Actions**:
1. Read the consolidated report cross-layer verification section for prompt template path entries
2. For each affected L3 lesson: update `prompt/data/` references to `prompt/templates/`
3. Verify no L2 convention documents reference the old `prompt/data/` path
4. Output a verification summary

**Verification**: Run `grep -r "prompt/data/" docs/` to confirm no stale references remain.

**Self-contained**: Yes.

---

### CV-006 [cross-layer-verification-type] Verify Skill Rename and Agent Removal Across Layers

**Priority**: P2
**Dependency**: Consolidated report Section "L3 Findings Cross-Referenced Against L2 (Reverse Feedback)"

**Context**:
Reverse feedback from L3 identified: (1) `record-task` skill renamed to `submit-task` — 4 lessons reference old name; (2) `error-fixer` agent removed — 2 lessons reference non-existent agent; (3) `pkg/template/` package removed — 3 lessons reference non-existent package. These L3 findings should trigger verification of corresponding L2 documents.

**Cross-layer influence list** (from consolidated report):
- `record-task` -> `submit-task`: 4 lessons reference old name; skill-structure.md should verify rename is reflected
- `error-fixer` agent removed: 2 lessons reference non-existent agent; forge-distribution.md should verify agent list completeness
- `pkg/template/` removed: 3 lessons reference non-existent package; package-organization.md should verify no template package deviation entry needed

**Actions**:
1. Read the consolidated report reverse feedback section
2. Verify `skill-structure.md` reflects the `record-task` -> `submit-task` rename
3. Verify `forge-distribution.md` agent list does not include `error-fixer`
4. Verify `package-organization.md` does not reference `pkg/template/`
5. Update affected L3 lessons: change `record-task` to `submit-task`, mark error-fixer and template package lessons as outdated
6. Output a verification summary

**Verification**: Run `grep -r "record-task" docs/`, `grep -r "error-fixer" docs/`, `grep -r "pkg/template/" docs/`.

**Self-contained**: Yes.

---

### CV-007 [cross-layer-verification-type] Verify qualitygate Subpackage and Surface-Specific Test Recipes

**Priority**: P2
**Dependency**: Consolidated report Section "Cross-Layer Verification" — `qualitygate/` subpackage undocumented, surface-specific test recipes missing

**Context**:
L2 identified `qualitygate/` subpackage as undocumented and surface-specific test recipes as missing from quality-gate.md. L3 has 3 lessons referencing old flat `cmd/quality_gate.go` path.

**Cross-layer influence list** (from consolidated report):
- L2-conventions-batch2: Lessons about quality gate architecture — 3 lessons reference old flat `cmd/quality_gate.go` path

**Actions**:
1. Read the consolidated report cross-layer verification section for qualitygate entries
2. For each affected L3 lesson: update `cmd/quality_gate.go` references to `cmd/qualitygate/quality_gate.go` (or mark as outdated if the lesson is about the old structure)
3. Verify FT-010 (package-organization.md fix) has documented the `qualitygate/` subpackage
4. Verify FT-008 (quality-gate.md fix) has documented surface-aware mode
5. Output a verification summary

**Verification**: Run `grep -r "cmd/quality_gate" docs/` and verify `ls internal/cmd/qualitygate/` matches documentation.

**Self-contained**: Yes.

---

## Needs-Update Items — Batch Fix Task

### FT-031 [fix-type] Batch Update 61 Knowledge Base Items with Stale Paths

**Priority**: P3
**Severity**: Low — core insights valid, only paths outdated

**Context**:
61 knowledge base items contain valid core insights but have outdated file paths, moved code references, or stale examples. Common path update patterns:
- **File path moves to subdirectories** (25+ items): `cmd/submit.go` -> `cmd/task/submit.go`, `cmd/claim.go` -> `cmd/task/claim.go`, `cmd/quality_gate.go` -> `cmd/qualitygate/quality_gate.go`
- **`tests/e2e/` -> `tests/<journey>/`** (15+ items): Test directory restructuring
- **`record-task` -> `submit-task`** (4 items): Skill rename
- **`skills/run-tasks/SKILL.md` -> `commands/run-tasks.md`** (5 items): Command/skill classification fix

**Actions**:
1. For each item in the needs-update list (see consolidated report Section "L3 Knowledge Base Detailed Findings — Needs-Update Items"):
   a. Read the item content
   b. Identify stale path references
   c. Apply the appropriate path transformation from the patterns above
   d. Verify the new path exists in the codebase
2. Output a summary of all updates applied

**Verification**: For each updated item, confirm the new code references exist via `find` or `grep`.

**Self-contained**: Yes — includes all context and path mapping rules.

---

## Summary Statistics

| Category | Count | Template Types |
|----------|-------|---------------|
| P0 (Immediate) | 1 | 1 fix-type |
| P1 (First Batch) | 11 | 11 fix-type |
| P2 (Second Batch) | 13 | 11 fix-type, 6 review-type, 7 cross-layer-verification-type |
| P3 (Third Batch) | 5 | 3 fix-type, 1 review-type, 1 fix-type (batch) |
| **Total** | **37** | **27 fix-type**, **7 review-type**, **7 cross-layer-verification-type** |

### Review-Type Tasks Requiring Human Confirmation

| Task ID | Description | Items Affected |
|---------|-------------|----------------|
| RT-001 | Archive resolved-bug lessons | 12 |
| RT-002 | Archive stale-reference lessons | 9 |
| RT-003 | Delete wrong-project lesson | 1 |
| RT-004 | Merge duplicate entries | 3 pairs |
| RT-005 | Decide empty decision stubs | 5 |
| RT-006 | Archive additional outdated items | 9 |
| RT-007 | Archive outdated decisions | 6 |
| **Total** | | **~45 items** |
