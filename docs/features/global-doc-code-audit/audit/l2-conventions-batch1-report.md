# L2 Conventions Batch 1 Audit Report

## Audit Baseline

- **Base commit**: `585d97075d4b3f88309e0f2a76aaf36d150b37da`
- **Audit date**: 2026-06-03
- **Audit scope**: 8 files in `docs/conventions/`:
  1. code-structure.md
  2. constants.md
  3. dead-code.md
  4. dispatcher-quality.md
  5. enum-constants.md
  6. error-handling.md
  7. forge-cli-reference.md
  8. forge-distribution.md

## Issue Summary

- P0: 0 | P1: 0 | P2: 5 | P3: 4

## Issue Details

### [P2] constants.md: Stale Line Reference for Gitignore Entry

- **File**: `docs/conventions/constants.md:30`
- **Declaration**: `"tests/results/"` at `init.go:46` is listed as a gitignore path entry.
- **Actual**: The `init.go` file at line 43 (not 46) uses `feature.TestResultsDir + "/"` -- a constant reference, not a literal string. The literal `"tests/results/"` no longer exists in init.go; the code references `TestResultsDir` constant from `pkg/feature/constants.go`.
- **Suggested action**: Update line reference from `init.go:46` to `init.go:43` and update the value column from `"tests/results/"` to `feature.TestResultsDir + "/"`. Consider adding a note that this is now a constant reference, making it an even stronger example of constant extraction.

### [P2] constants.md: False "Extracted" Claim for defaultHealthPath

- **File**: `docs/conventions/constants.md:31`
- **Declaration**: Table lists `defaultHealthPath` at `pkg/serverprobe/constants.go` as a path constant example with current state "all extracted".
- **Actual**: `pkg/serverprobe/constants.go` only contains `defaultProbeTimeout`. The `/health` path is an inline default in `pkg/serverprobe/serverprobe.go:31` (`path = "/health"`), NOT a named constant. There is no `defaultHealthPath` constant anywhere in the codebase.
- **Suggested action**: Either (a) create the `defaultHealthPath` constant in `pkg/serverprobe/constants.go` and use it in `serverprobe.go`, or (b) remove this row from the constants table and mark it as an acceptable inline value (similar to the gitignore entry justification).

### [P2] constants.md: "All Extracted" Claim Contradicted by testrunner Literals

- **File**: `docs/conventions/constants.md:24-33,152-159`
- **Declaration**: Path constants section states "current state -- all extracted" and the deviation analysis marks P1/P2 as "Fixed".
- **Actual**: `pkg/testrunner/test_results.go` contains four instances of literal path strings that should use existing constants:
  - Line 19: `"tests"` and `"results"` (should use `feature.TestResultsDir` or `feature.GetTestResultsDir()`)
  - Line 24: `"raw-output.txt"` (should use `feature.TestOutputFileName`)
  - Line 28: `"unit-raw-output.txt"` (should use `feature.UnitTestOutputFileName`)
- **Suggested action**: (1) Update `test_results.go` to import and use the constants from `pkg/feature/constants.go`. (2) Add a deviation entry in constants.md for the testrunner literals until fixed. (3) Change "current state -- all extracted" to acknowledge the remaining deviation.

### [P2] forge-cli-reference.md: Incorrect Source File for quality-gate Command

- **File**: `docs/conventions/forge-cli-reference.md:22`
- **Declaration**: `forge quality-gate` source file listed as `quality_gate.go` (top-level in `internal/cmd/`).
- **Actual**: The source file is at `internal/cmd/qualitygate/quality_gate.go` (inside the `qualitygate/` subpackage). The top-level `quality_gate.go` file does not exist at `internal/cmd/quality_gate.go`.
- **Suggested action**: Update the source file column from `quality_gate.go` to `qualitygate/quality_gate.go`.

### [P2] forge-distribution.md: Incomplete hooks/ Directory Tree

- **File**: `docs/conventions/forge-distribution.md:48-56`
- **Declaration**: The directory tree for hooks/ shows three entries: `hooks.json`, `session-start`, `guide.md`.
- **Actual**: The hooks/ directory contains five entries: `hooks.json`, `session-start`, `guide.md`, `run-hook.cmd`, `debug`. The tree omits `run-hook.cmd` (the cross-platform polyglot hook wrapper used in hooks.json) and `debug` (a debug utility script).
- **Suggested action**: Update the directory tree to include `run-hook.cmd` and `debug`. Add `run-hook.cmd` to the component table in Section 2 as the hook execution wrapper. Optionally mention `debug` as a development utility (not user-facing).

### [P3] forge-distribution.md: Hooks Table Omits SessionEnd and SubagentStop

- **File**: `docs/conventions/forge-distribution.md:79-84`
- **Declaration**: Component table for hooks/ describes three hook events: SessionStart, SubagentStart (inject guide.md), and Stop (quality-gate + feature complete).
- **Actual**: `hooks.json` also defines `SessionEnd` and `SubagentStop` hooks, both running `forge cleanup`. These are active hooks that execute cleanup logic but are not documented in the component table.
- **Suggested action**: Add SessionEnd and SubagentStop to the hooks component table, noting they run `forge cleanup` for state cleanup.

### [P3] forge-distribution.md: run-tasks Misclassified as Skill in Pipeline Diagrams

- **File**: `docs/conventions/forge-distribution.md:188-192`
- **Declaration**: Both pipeline diagrams show `/run-tasks` alongside skills (brainstorm, write-prd, etc.). Section 6 "辅助 Skill" lists `/run-tasks` as a skill.
- **Actual**: In the plugin structure, `run-tasks` is a command (`plugins/forge/commands/run-tasks.md`), not a skill. It is registered as a slash command entry point, not a skill with SKILL.md. The distinction matters because commands and skills have different path resolution rules (Section 5).
- **Suggested action**: In the pipeline diagrams and "辅助 Skill" section, clarify that `/run-tasks` is a command, or move it to a separate "辅助 Commands" subsection. The pipeline flow is correct; only the classification label is misleading.

### [P3] constants.md: Inconsistent Deviation Status for ANSI Codes

- **File**: `docs/conventions/constants.md:169`
- **Declaration**: Color deviation C4 lists `"\033[33m"` / `"\033[0m"` at `list.go` as "Fixed: ANSI codes cleaned up".
- **Actual**: The enum-constants.md TECH-const-003 section claims `colorCycleMarker = "\033[33m"` and `colorReset = "\033[0m"` are defined in `internal/cmd/styles.go` as the target pattern, but these constants do not exist in the current codebase (`styles.go` only has hex color constants). The ANSI codes in `list.go` have been cleaned up (confirmed), but the pattern shown as the "correct" example in enum-constants.md is aspirational, not current.
- **Suggested action**: Verify whether the ANSI constants `colorCycleMarker`/`colorReset` were ever created. If not, either (a) create them in `styles.go` and reference them, or (b) update enum-constants.md TECH-const-003 to remove the aspirational ANSI constant examples and note that ANSI cleanup was done by removing the inline usage rather than extracting to named constants.

### [P3] code-structure.md: validate_index.go Rename Reference

- **File**: `docs/conventions/code-structure.md` (implied by dead-code.md:83)
- **Declaration**: dead-code.md DC-2 deviation detail states "file has been renamed from validate_index.go to validate.go".
- **Actual**: `internal/cmd/task/validate.go` exists. No `validate_index.go` exists. The rename is confirmed. However, the code-structure.md deviation table does not mention this rename; only dead-code.md does. This is a minor cross-reference inconsistency -- the code-structure.md CS-2 entry should ideally mention the file rename for completeness.
- **Suggested action**: Add a note to code-structure.md CS-2 deviation entry mentioning the file rename from `validate_index.go` to `validate.go` for traceability.

## Verified Claims (No Issues Found)

The following declarations were verified as accurate against the current codebase:

### code-structure.md
- CS-1: `cmd.Debugf` duplicate removed -- only `base.Debugf` exists in `base/output.go` (line 85). MATCHES.
- CS-2: `getTaskPhase` alias removed, direct calls to `task.GetTaskPhase()`. MATCHES.
- CS-3: `checkExistingTaskState` alias removed. MATCHES.
- CS-4: `compareVersionIDs` alias removed. MATCHES.
- CS-5: `FrontmatterData.Scope` removed from `frontmatter.go`. `Task.Scope` retained in `types.go` for migration. MATCHES.
- TECH-code-structure-004: No violations of dependency direction in `internal/cmd/task/` -- imports only `base` and `pkg/*`. MATCHES.

### constants.md
- `TestOutputFileName = "raw-output.txt"` in `pkg/feature/constants.go:61`. MATCHES.
- `UnitTestOutputFileName = "unit-raw-output.txt"` in `pkg/feature/constants.go:62`. MATCHES.
- `colorModeHighlight`, `colorConflict`, `colorSource` in `internal/cmd/styles.go`. MATCHES.
- `probeRetryInterval`, `maxProbeRetries` in `qualitygate/constants.go`. MATCHES.
- `defaultProbeTimeout` in `pkg/serverprobe/constants.go`. MATCHES.
- `defaultLockTimeout`, `lockRetryBackoff` in `pkg/index/lock.go`. MATCHES.
- `fallbackSortPriority` in `task/list.go:24`, `unreachableDepth` in `task/claim.go:23`. MATCHES.
- `conciseErrorMaxLines = 5`, `maxSourceFiles = 10` in `qualitygate/constants.go`. MATCHES.
- `statusColor` function in `task/tree.go` uses named color strings. MATCHES.
- `SlugColMinWidth = 30`, `SlugColMaxWidth = 60` in `base/output.go`. MATCHES.

### dead-code.md
- All 5 deviations (DC-1 through DC-5) confirmed as resolved. MATCHES.
- `CheckLegacyScope` retained in `pkg/task/migrate.go` for migration detection. MATCHES.
- Cleanup procedures accurately describe the process used. MATCHES.

### enum-constants.md
- `pkg/types/` contains `status.go`, `surface.go`, `priority.go`. MATCHES.
- `AllStatuses()`, `AllPriorities()`, `AllSurfaceTypes()` exist. MATCHES.
- `IsTerminalStatus()` exists in `pkg/types/status.go`. MATCHES.
- Type aliases in `pkg/feature/constants.go` for backward compatibility. MATCHES.
- Task type constants remain in `pkg/task/types.go` (not migrated to `pkg/types/`). MATCHES.
- `pkg/types/` has zero forge-cli internal imports (leaf package). MATCHES.

### error-handling.md
- Error messages use `fmt.Errorf("context: detail")` pattern. MATCHES.
- Warning/diagnostic output uses `fmt.Fprintf(os.Stderr, ...)`. MATCHES.
- No error messages found going to stdout. MATCHES.

### forge-cli-reference.md
- All 13 top-level commands verified against `root.go` init() function. MATCHES.
- All 13 task subcommands verified against `task/register.go`. MATCHES.
- All worktree subcommands (5) verified. MATCHES.
- All forensic subcommands (3) verified. MATCHES.
- All fact subcommands (3) verified. MATCHES.
- Feature subcommands (4) verified. MATCHES.
- Config subcommands (3) verified. MATCHES.
- `forge config get mode` returns `quick`/`full`/`none` based on path detection. MATCHES.
- Removed commands section: `forge probe`, `forge e2e`, `forge test` subcommands confirmed absent from codebase. MATCHES.
- List sorting by `created` frontmatter descending. MATCHES.
- Source file references (except quality-gate noted above) all exist. MATCHES.

### forge-distribution.md
- Plugin directory structure (`plugins/forge/`) verified: `.claude-plugin/`, `agents/`, `commands/`, `hooks/`, `skills/`. MATCHES.
- `plugin.json` exists at `.claude-plugin/plugin.json`. MATCHES.
- `agents/task-executor.md` exists. MATCHES.
- `hooks.json` uses `${CLAUDE_PLUGIN_ROOT}` for hook paths. MATCHES.
- `skills/eval/experts/` structure matches: `protocol/`, `freeform/`, `scorer/` with exact file names. MATCHES.
- Expert files: `scorer-protocol.md`, `reviser-protocol.md`, 8 scorer roles, 5 freeform files. All MATCHES.
- Path resolution rules: no `${CLAUDE_SKILL_DIR}` usage found in skills. MATCHES.
- No absolute `plugins/forge/...` paths in skills. MATCHES.
- `BuildIndexOpts.Intent` field exists in `pkg/task/build.go`. MATCHES.
- Intent types (`new-feature`, `enhancement`, `refactor`, `cleanup`, `fix`, `doc`) match proposal frontmatter values. MATCHES.
- Pipeline skills/commands verified: brainstorm, write-prd, ui-design, tech-design, breakdown-tasks, quick-tasks, submit-task, gen-journeys, gen-contracts, gen-test-scripts, run-tests all exist. MATCHES.
- Auxiliary skills verified: consolidate-specs, learn, forensic, init-justfile, eval, clean-code, deep-research, extract-design-md, gen-web-sitemap, test-guide. MATCHES.

## Cross-Layer Influence Items (for L3 Reference)

The following findings from this audit may affect knowledge base entries in `docs/lessons/` or `docs/decisions/`:

1. **Path constant extraction lessons**: If any lesson entries reference `defaultHealthPath` or claim all path constants are extracted, those entries need updating since `defaultHealthPath` does not exist and `testrunner` still uses literals.

2. **Hook lifecycle lessons**: If any lesson entries describe hook events, they may be incomplete if they only cover SessionStart/Stop and omit SessionEnd/SubagentStop.

3. **run-tasks classification**: If any decisions reference run-tasks as a "skill", the classification should be corrected to "command".

## Audit Quality Review

- Sample ratio: 100% (all 8 target files fully audited, every declaration verified)
- Sample result: PASS
- Missed items: 0
- Expand review: No
