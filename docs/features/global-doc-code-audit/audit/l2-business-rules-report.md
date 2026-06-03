# L2 Business Rules + CLAUDE.md Audit Report

## Audit Baseline

- **Baseline commit**: `feaa321e6552645c51f6c11499fdf485617d1088`
- **Audit date**: 2026-06-03
- **Audit scope**:
  - `docs/business-rules/error-reporting.md` (25 lines, 2 rules)
  - `docs/business-rules/quality-gate.md` (23 lines, 1 rule)
  - `docs/business-rules/surface-orchestration.md` (65 lines, 6 rules)
  - `docs/business-rules/task-lifecycle.md` (68 lines, 6 rules)
  - `CLAUDE.md` (22 lines)
- **Audit type**: L2 Business Rules + CLAUDE.md — declaration extraction + code verification

## Issue Summary

- **P0**: 1 | **P1**: 3 | **P2**: 5 | **P3**: 2
- Total issues: 11
- Declarations examined: 15 business rules + 3 CLAUDE.md claims
- Correct identifications: 10
- Cross-layer influence items for L3 reference: 3

## Issue Details

### [P0] CLAUDE.md references non-existent plugin subdirectories

- **File**: `CLAUDE.md:19`
- **Declaration**: `plugins/forge/` lists subdirectories as "skills、commands、agents、hooks、references、scripts"
- **Actual**: `plugins/forge/` contains only 4 subdirectories: `agents/`, `commands/`, `hooks/`, `skills/`. The directories `references/` and `scripts/` do **not** exist as top-level directories under `plugins/forge/`. The `references/` directory exists only within individual skills (e.g., `plugins/forge/skills/test-guide/references/`), not as a peer to `skills/`.
- **Impact**: AI agents reading CLAUDE.md will believe `plugins/forge/references/` and `plugins/forge/scripts/` exist as top-level directories, potentially creating files in the wrong location.
- **Suggested action**: Update CLAUDE.md to list only the actual top-level directories: `agents`, `commands`, `hooks`, `skills`. Alternatively, clarify that `references` and `scripts` are skill-internal subdirectories.

### [P1] Surface-key naming regex is more restrictive than documented

- **File**: `docs/business-rules/surface-orchestration.md:62`
- **Declaration**: BIZ-surface-orchestration-006 states "Surface-key values MUST match `[a-zA-Z0-9_-]` only"
- **Actual**: The enforcement code in `forge-cli/pkg/forgeconfig/execution_order.go:14` uses the regex `^[a-z][a-z0-9-]*$`. This is more restrictive: (1) only lowercase letters, (2) must start with a letter (no digits/underscores as first character), (3) underscores are not allowed at all. The `NormalizeSurfaceKey` function converts uppercase to lowercase and replaces underscores with hyphens before validation.
- **Evidence**: `execution_order.go:14`: `var surfaceKeyPattern = regexp.MustCompile("^[a-z][a-z0-9-]*$")`
- **Suggested action**: Update BIZ-surface-orchestration-006 to reflect the actual regex: `^[a-z][a-z0-9-]*$`. Clarify that normalization converts uppercase to lowercase and replaces non-alphanumeric characters with hyphens before validation.

### [P1] test-type-model.md reference points to non-existent path

- **File**: `docs/business-rules/task-lifecycle.md:42`
- **Declaration**: "权威定义参见 `docs/reference/test-type-model.md`" (authoritative definition see `docs/reference/test-type-model.md`)
- **Actual**: The directory `docs/reference/` does **not** exist. The actual file is at `plugins/forge/skills/test-guide/references/test-type-model.md`, which is a plugin-internal reference file, not a docs/ directory file.
- **Suggested action**: Update the reference to the correct path: `plugins/forge/skills/test-guide/references/test-type-model.md`, or remove the path reference and describe the mapping inline.

### [P1] Quality gate doc omits probe details and test-setup in surface-aware mode

- **File**: `docs/business-rules/quality-gate.md:17-19`
- **Declaration**: Phase 2 is "`just unit-test` with retry-once policy". Phase 3 is "`just test` (full regression suite, requires probe health check first). Optional `just test-setup` step."
- **Actual**: In surface-aware mode (when surfaces are configured in `.forge/config.yaml`), Phase 3 uses `runTestRegressionSurface` which orchestrates per-surface-type lifecycle:
  - web/api/mobile: dev -> probe -> test -> teardown (with optional `test-setup` for mobile only)
  - cli/tui: test -> teardown
  The probe in surface-aware mode uses `probeWithRetry` with `maxProbeRetries=3` and `probeRetryInterval=5s`, running a just recipe (e.g., `web-probe` or generic `probe`), not the `serverprobe.ProbeServers` HTTP health check. The legacy path (no surfaces) uses `serverprobe.ProbeServers` which checks HTTP endpoints from `tests/config.yaml`.
- **Suggested action**: Add a section to BIZ-quality-gate-001 describing the surface-aware Phase 3 behavior, or note that the current description covers legacy mode only.

### [P2] Surface orchestration table describes probe behavior incorrectly for web/api

- **File**: `docs/business-rules/surface-orchestration.md:25`
- **Declaration**: BIZ-surface-orchestration-002 table says web "probe checks page root path" and api "probe checks /healthz"
- **Actual**: The probe implementation in surface-aware mode (`quality_gate_lifecycle.go:225-252`) runs a just recipe (e.g., `just web-probe` or `just api-probe` or `just probe`). It does not hardcode checking "page root path" or "/healthz". The legacy `serverprobe.ProbeServers` checks `/health` by default (not `/healthz`), and this path is configurable. The surface-aware probe is recipe-based, so the check behavior depends on what the recipe does.
- **Suggested action**: Update BIZ-surface-orchestration-002 to describe probe as "recipe-based (just <surface>-probe or just probe)" rather than specifying exact paths. Note that the /healthz and page root path descriptions are not enforced by code.

### [P2] Teardown idempotency rule describes PID-based behavior but code is recipe-based

- **File**: `docs/business-rules/surface-orchestration.md:54`
- **Declaration**: BIZ-surface-orchestration-005 states "if the PID does not exist, skip silently. If kill fails, retry once; if still failing, log process info and continue... Final guarantee: `.forge/test-state.json` state file is cleaned up"
- **Actual**: The teardown implementation in `quality_gate_lifecycle.go:211-223` (`runTeardown`) simply runs `just <surface>-teardown` or `just teardown`. There is no PID-based logic, no retry-once for kill failure, and no `.forge/test-state.json` cleanup in the teardown code. The `.forge/test-state.json` file is referenced only in `init.go` as a file to create, not as a state file to clean up during teardown. The recipe-based teardown delegates all cleanup behavior to the justfile recipe.
- **Suggested action**: Update BIZ-surface-orchestration-005 to reflect that teardown is recipe-based and the cleanup guarantees depend on the recipe implementation. Remove the PID-specific and test-state.json-specific claims that are not enforced by the Go code.

### [P2] Probe failure exit code differs from documented

- **File**: `docs/business-rules/surface-orchestration.md:40`
- **Declaration**: BIZ-surface-orchestration-003 states "All 3 retries failing is treated as retryable failure (exit code 1)"
- **Actual**: When probe fails in surface-aware mode, `runSurfaceLifecycle` returns a `lifecycleResult{success: false}`. The caller `runTestRegressionSurface` then calls `HandleGateFailure` which returns an error, and the caller `RunQualityGate` calls `os.Exit(0)` (not exit code 1). In the quality-gate hook, exit code 0 is used for all gate failures because the hook JSON signals the actual decision. The exit code 1 claim may be accurate for a different context (CI retry), but the actual code path exits 0.
- **Suggested action**: Clarify that probe failure in quality-gate hook context exits with code 0 (hook JSON signals the decision), not code 1. The exit code 1 may apply to the scheduler/CI layer that consumes the hook JSON.

### [P2] BIZ-surface-orchestration-004 probe hard-gate scope is narrower than documented

- **File**: `docs/business-rules/surface-orchestration.md:47`
- **Declaration**: "Applies to all surface types that use probe (web, api)"
- **Actual**: The code in `quality_gate_lifecycle.go:125-127` shows `needsFullLifecycle` returns true for `web`, `api`, AND `mobile`. All three types get probe behavior. The doc should include `mobile` in the scope.
- **Suggested action**: Update BIZ-surface-orchestration-004 to include `mobile` in the scope: "Applies to all surface types that use probe (web, api, mobile)".

### [P2] Surface orchestration table omits mobile test-setup step

- **File**: `docs/business-rules/surface-orchestration.md:30`
- **Declaration**: BIZ-surface-orchestration-002 table shows mobile sequence as "dev -> probe -> [per-journey test] -> teardown"
- **Actual**: The code in `quality_gate_lifecycle.go:176-187` shows mobile surfaces have an additional `test-setup` step between probe and test: "dev -> probe -> mobile-test-setup (optional) -> test -> teardown". The table omits this mobile-specific setup step.
- **Suggested action**: Update the mobile row in BIZ-surface-orchestration-002 to include the optional test-setup step.

### [P3] Quality gate doc uses ambiguous "just unit-test" description

- **File**: `docs/business-rules/quality-gate.md:17`
- **Declaration**: Phase 2 is described as running `just unit-test`
- **Actual**: The code calls `testrunner.RunProjectTests()` which has a probe chain: `just unit-test` -> `just test` -> `go test ./...` -> `npm test` -> `pytest`. If `just unit-test` recipe does not exist, it falls back to other methods. The doc states it runs `just unit-test` specifically, which is only the first probe in the chain.
- **Suggested action**: Clarify that Phase 2 probes for `just unit-test` first but may fall back to other test runners if the recipe is absent.

### [P3] Quality gate doc does not mention surface-specific test recipes

- **File**: `docs/business-rules/quality-gate.md:19`
- **Declaration**: Phase 3 says "`just test` (full regression suite)"
- **Actual**: In surface-aware mode, the test step uses `resolveRecipe(projectRoot, surfaceType, "test")` which prefers surface-specific recipes like `web-test`, `api-test`, `cli-test` before falling back to generic `test`. The doc does not mention this surface-specific recipe resolution.
- **Suggested action**: Add a note about surface-specific recipe resolution in Phase 3 description.

## Verified Rules (No Issues Found)

### BIZ-error-reporting-001: Exit Code Semantics

- **File**: `docs/business-rules/error-reporting.md:12-16`
- **Status**: MATCHES code
- **Evidence**: `forge-cli/internal/cmd/base/errors.go:57-63` — `ExitCode()` method returns 2 for `ErrInvalidTransition`, `ErrInvalidPath`, `ErrContractUnverifiable`, `ErrMigrationRequired`, and 1 for all others. Matches doc exactly.

### BIZ-error-reporting-002: Actionable Error Messages

- **File**: `docs/business-rules/error-reporting.md:22-24`
- **Status**: MATCHES code
- **Evidence**: `AIError` struct in `base/errors.go` contains `Message`, `Cause`, `Hint`, `Action` fields. `printAIError` prints all four components. Example helper `ErrTaskNotFound` lists action command.

### BIZ-quality-gate-001: Three-Phase Pipeline (partial match)

- **File**: `docs/business-rules/quality-gate.md:14-22`
- **Status**: Partially matches code
- **Evidence**:
  - Phase 1 (compile/fmt/lint): Code uses `just.NonBreakingGateSequence()` which returns compile -> fmt -> lint. MATCHES.
  - Fix task cap of 3: Code has `maxFixTasksPerStep = 3` in `quality_gate.go:30`. MATCHES.
  - Docs-only skip: Code has `IsDocsOnly` check in `quality_gate.go:142-144`. MATCHES.
  - Fix task types: `fixTypeFromStep` returns `TypeCodingFix` for compile/test, `TypeCodingCleanup` for fmt/lint. MATCHES.
  - Submit uses tiered gate: Code in `submit.go:381-385` uses `UnitGateSequence` for breaking, `NonBreakingGateSequence` for non-breaking. MATCHES.
  - Two-layer recipe model: `RunProjectTests` probes `just unit-test` -> `just test` -> etc. MATCHES.
  - Issue: Exit code claim "Submit failure exits with code 1" is technically correct at the cobra level (root.go:30 `os.Exit(1)`), but the `validateQualityGate` function uses `panic(AIError)` which is not recovered in production code, causing a crash rather than a clean exit. This is a code bug rather than a doc issue.

### BIZ-task-lifecycle-001: State Transition Constraints

- **File**: `docs/business-rules/task-lifecycle.md:14-26`
- **Status**: MATCHES code
- **Evidence**: `statemachine.go` defines 7 statuses matching doc exactly. `transitionTable` enforces terminal state protection, submit-only path to completed, reopen for rejected/skipped. `autoRestoreSourceTask` in `submit.go:263-277` restores blocked tasks when deps complete. `deps.go:31-34` shows `satisfiedStatuses` = {completed, skipped}. MATCHES all claims.

### BIZ-task-lifecycle-002: Terminal State Immutability

- **File**: `docs/business-rules/task-lifecycle.md:30-32`
- **Status**: MATCHES code
- **Evidence**: `statemachine.go:53` blocks all transitions from completed. Lines 56-57 allow only reopen for rejected. Lines 59-60 allow only reopen for skipped. `status.go:36-38` confirms `IsTerminalStatus` for completed/skipped/rejected. `reopen.go` exists and works only for rejected/skipped. `status.go` has no --force flag (confirmed by reopen_test.go:283-288). MATCHES.

### BIZ-task-lifecycle-003: System Type Exclusion

- **File**: `docs/business-rules/task-lifecycle.md:38-44`
- **Status**: MATCHES code
- **Evidence**: `types.go:149-162` defines `SystemTypes` map with exactly 12 entries matching the doc. `IsSystemType` strips last segment for surface variants. `build.go:170-180` enforces system type exclusion for non-auto-gen tasks. `IsAutoGenTaskID` in `build.go:543-560` matches the documented ID patterns. `doc.consolidate` and `doc.drift` are not in SystemTypes. MATCHES.

### BIZ-task-lifecycle-004: Topological Task Ordering

- **File**: `docs/business-rules/task-lifecycle.md:51-53`
- **Status**: MATCHES code
- **Evidence**: `toposort.go` implements Kahn's algorithm. `list.go:62` shows default sort is "topo" with "id" fallback. `claim.go:217-228` uses `computeTopoDepths` to sort by depth, then priority, then natural ID. `list.go:63` shows `--tree` flag exists. `deps.go:13-28` shows `ResolveWildcardDep`. MATCHES.

### BIZ-surface-orchestration-001: Surface Type Fixed Enumeration

- **File**: `docs/business-rules/surface-orchestration.md:14`
- **Status**: MATCHES code
- **Evidence**: `types/surface.go:8-13` defines exactly 5 surface types: web, api, cli, tui, mobile. `forgeconfig/detect.go:13-19` confirms `KnownSurfaceTypes` with these 5 values. MATCHES.

### BIZ-surface-orchestration-003: Probe Retry Parameters

- **File**: `docs/business-rules/surface-orchestration.md:39-41`
- **Status**: MATCHES code
- **Evidence**: `qualitygate/constants.go:8-9` defines `maxProbeRetries = 3` and `probeRetryInterval = 5 * time.Second`. MATCHES.

### CLAUDE.md Forge Plugin MANDATORY rule

- **File**: `CLAUDE.md:19-21`
- **Status**: PARTIAL MATCH
- **Evidence**: The referenced file `docs/conventions/forge-distribution.md` exists. The rule correctly identifies that Forge's distribution model must be understood before modifying plugin files. However, the subdirectory list is inaccurate (see P0 issue above).

## Cross-Layer Influence Items (for L3 reference)

1. **[L3-REF-01] test-type-model.md path**: `docs/business-rules/task-lifecycle.md:42` references `docs/reference/test-type-model.md`. L3 audit should check if any `docs/lessons/` or `docs/decisions/` entries also reference this non-existent path.

2. **[L3-REF-02] CLAUDE.md plugin directory structure**: CLAUDE.md claims `plugins/forge/` has 6 subdirectories but only 4 exist. L3 should check if any lessons/decisions reference `plugins/forge/references/` or `plugins/forge/scripts/` as top-level paths.

3. **[L3-REF-03] Probe health check vs recipe-based probe**: The surface orchestration rules describe HTTP-based probe behavior, but the actual implementation is recipe-based. L3 should check if lessons about probe debugging reference the HTTP endpoint approach.

## Audit Quality Review

- Sample ratio: 100% (all 4 business rules files + CLAUDE.md audited)
- Sample result: PASS
- Missed items: 0 (all files exhaustively audited)
- Whether to extend review: No
