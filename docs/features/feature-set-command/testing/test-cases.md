---
feature: "feature-set-command"
sources:
  - docs/proposals/feature-set-command/proposal.md
  - docs/features/feature-set-command/tasks/1-feature-set-command.md
  - docs/features/feature-set-command/tasks/2-priority-chain-state.md
  - docs/features/feature-set-command/tasks/3-verbose-flag.md
generated: "2026-05-16"
---

# Test Cases: feature-set-command

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 20  |
| **Total** | **20** |

> **Note**: This is a pure CLI feature. No UI or API interfaces exist. All test cases exercise the `forge` binary via command-line invocations.

---

## CLI Test Cases

### Task 1: `forge feature set` subcommand

## TC-001: Set feature creates directory and state
- **Source**: Proposal Success Criteria #1, #2 / Task 1 AC #1, #2
- **Type**: CLI
- **Target**: cli/feature-set
- **Test ID**: cli/feature-set/set-feature-creates-directory-and-state
- **Pre-conditions**: Clean project directory with `go.mod`, no `.forge/state.json`, no `docs/features/my-feature/`
- **Steps**:
  1. Run `forge feature set my-feature`
  2. Check exit code is 0
  3. Verify stdout contains `FEATURE: my-feature`
  4. Verify `.forge/state.json` exists with `feature: "my-feature"` and `allCompleted: false`
  5. Verify `docs/features/my-feature/` directory structure exists (including tasks/process subdirs)
- **Expected**: Feature directory is created, state.json is written with correct slug and allCompleted=false, stdout confirms the feature
- **Priority**: P0

## TC-002: Set feature with empty slug returns error
- **Source**: Proposal Success Criteria #3 / Task 1 AC #3
- **Type**: CLI
- **Target**: cli/feature-set
- **Test ID**: cli/feature-set/set-feature-empty-slug-returns-error
- **Pre-conditions**: Clean project directory with `go.mod`
- **Steps**:
  1. Run `forge feature set ""`
  2. Check exit code is non-zero
  3. Verify `.forge/state.json` does NOT exist
- **Expected**: Command fails with error, no state.json or feature directory is created
- **Priority**: P0

## TC-003: Set feature prints slug to stdout
- **Source**: Task 1 AC #4
- **Type**: CLI
- **Target**: cli/feature-set
- **Test ID**: cli/feature-set/set-feature-prints-slug-to-stdout
- **Pre-conditions**: Clean project directory with `go.mod`
- **Steps**:
  1. Run `forge feature set test-slug`
  2. Capture stdout
- **Expected**: Stdout contains `FEATURE: test-slug`
- **Priority**: P0

## TC-004: Positional arg backward compatibility
- **Source**: Task 1 AC #5
- **Type**: CLI
- **Target**: cli/feature-positional
- **Test ID**: cli/feature-positional/positional-arg-backward-compatibility
- **Pre-conditions**: Clean project directory with `go.mod`
- **Steps**:
  1. Run `forge feature legacy-feature`
  2. Verify exit code is 0
  3. Verify `docs/features/legacy-feature/` directory exists
  4. Verify `.forge/state.json` does NOT exist (old behavior: positional arg does not write state)
- **Expected**: Existing `forge feature <slug>` behavior unchanged; directory created but state.json not written
- **Priority**: P0

## TC-005: Set feature with whitespace-only slug returns error
- **Source**: Task 1 Hard Rules (validate slug is non-empty before filesystem ops)
- **Type**: CLI
- **Target**: cli/feature-set
- **Test ID**: cli/feature-set/set-feature-whitespace-slug-returns-error
- **Pre-conditions**: Clean project directory with `go.mod`
- **Steps**:
  1. Run `forge feature set "   "`
  2. Check exit code is non-zero
  3. Verify `.forge/state.json` does NOT exist
- **Expected**: Command fails with error, no side effects on filesystem
- **Priority**: P1

## TC-006: Set feature idempotent on repeated calls
- **Source**: Proposal Key Scenarios (happy path) / inferred from EnsureFeatureDir + EnsureForgeState semantics
- **Type**: CLI
- **Target**: cli/feature-set
- **Test ID**: cli/feature-set/set-feature-idempotent-on-repeated-calls
- **Pre-conditions**: Clean project directory with `go.mod`
- **Steps**:
  1. Run `forge feature set my-feature`
  2. Verify success
  3. Run `forge feature set my-feature` again
  4. Verify success
  5. Verify `.forge/state.json` still has `feature: "my-feature"`
- **Expected**: Second call succeeds without error; state.json remains valid
- **Priority**: P1

## TC-007: Set feature overwrites previous feature in state
- **Source**: Proposal Key Scenarios (worktree mismatch) — user switches feature
- **Type**: CLI
- **Target**: cli/feature-set
- **Test ID**: cli/feature-set/set-feature-overwrites-previous-feature
- **Pre-conditions**: Clean project directory with `go.mod`, `forge feature set feature-a` already run
- **Steps**:
  1. Run `forge feature set feature-b`
  2. Verify `.forge/state.json` has `feature: "feature-b"`
  3. Verify `docs/features/feature-b/` directory exists
- **Expected**: State updated to new feature; old feature directory remains intact
- **Priority**: P1

### Task 2: `GetCurrentFeature()` priority chain

## TC-008: GetCurrentFeature returns state.json feature when present
- **Source**: Proposal Success Criteria #4 / Task 2 AC #1
- **Type**: CLI
- **Target**: cli/feature-query
- **Test ID**: cli/feature-query/get-current-feature-returns-state-json-feature
- **Pre-conditions**: Project with `.forge/state.json` containing `feature: "explicit-feature"`, feature directory `docs/features/explicit-feature/` exists
- **Steps**:
  1. Run `forge feature`
  2. Verify stdout contains `FEATURE: explicit-feature`
- **Expected**: Feature resolved from state.json, overriding any git context
- **Priority**: P0

## TC-009: GetCurrentFeature falls back when state.json absent
- **Source**: Proposal Success Criteria #6 / Task 2 AC #2
- **Type**: CLI
- **Target**: cli/feature-query
- **Test ID**: cli/feature-query/get-current-feature-falls-back-when-state-json-absent
- **Pre-conditions**: Project with git branch `feature/branch-feature`, `docs/features/branch-feature/` exists, no `.forge/state.json`
- **Steps**:
  1. Delete `.forge/state.json` if it exists
  2. Run `forge feature`
  3. Verify stdout contains `FEATURE: branch-feature`
- **Expected**: Resolution falls back to git context when state.json is deleted (quality-gate cleanup scenario)
- **Priority**: P0

## TC-010: GetCurrentFeatureWithSource returns correct source type
- **Source**: Task 2 AC #3
- **Type**: CLI
- **Target**: cli/feature-query
- **Test ID**: cli/feature-query/get-current-feature-with-source-returns-correct-source
- **Pre-conditions**: Multiple scenarios to test each source type
- **Steps**:
  1. With state.json set and feature dir existing: run `forge feature -v`, verify `(from: state.json)`
  2. On feature branch without state.json: run `forge feature -v`, verify `(from: branch)` or `(from: worktree)`
  3. With single feature dir, no state.json, no feature branch: run `forge feature -v`, verify `(from: features-dir)`
- **Expected**: Each resolution source is correctly identified and displayed
- **Priority**: P0

## TC-011: State.json with nonexistent feature directory falls through
- **Source**: Task 2 Hard Rules (when feature directory doesn't exist, skip to next priority)
- **Type**: CLI
- **Target**: cli/feature-query
- **Test ID**: cli/feature-query/state-json-with-nonexistent-dir-falls-through
- **Pre-conditions**: `.forge/state.json` with `feature: "ghost-feature"`, no `docs/features/ghost-feature/` directory, git branch `feature/fallback-feature` with directory existing
- **Steps**:
  1. Run `forge feature -v`
  2. Verify output does NOT show `ghost-feature`
  3. Verify output shows `fallback-feature (from: branch)` or similar
- **Expected**: state.json entry skipped when its feature directory is missing; falls back to git context
- **Priority**: P0

## TC-012: Corrupt state.json falls through silently
- **Source**: Task 2 Hard Rules (state.json read failure silently ignored)
- **Type**: CLI
- **Target**: cli/feature-query
- **Test ID**: cli/feature-query/corrupt-state-json-falls-through-silently
- **Pre-conditions**: `.forge/state.json` containing invalid JSON (e.g., `not json at all`), git branch with valid feature directory
- **Steps**:
  1. Run `forge feature -v`
  2. Verify no error about corrupt state.json
  3. Verify feature resolved from fallback source
- **Expected**: Corrupt state.json silently ignored; feature resolved from git or features-dir
- **Priority**: P1

## TC-013: State.json takes priority over git worktree
- **Source**: Proposal Key Scenarios (worktree mismatch)
- **Type**: CLI
- **Target**: cli/feature-query
- **Test ID**: cli/feature-query/state-json-takes-priority-over-git-worktree
- **Pre-conditions**: In a git worktree named `fix-auth-bug`, `.forge/state.json` contains `feature: "oauth-rewrite"`, `docs/features/oauth-rewrite/` exists
- **Steps**:
  1. Run `forge feature -v`
  2. Verify output shows `oauth-rewrite (from: state.json)`
- **Expected**: Explicit state.json selection overrides worktree-derived feature name
- **Priority**: P0

## TC-014: Existing callers unchanged after priority chain change
- **Source**: Task 2 AC #5, #6
- **Type**: CLI
- **Target**: cli/feature-query
- **Test ID**: cli/feature-query/existing-callers-unchanged-after-priority-chain
- **Pre-conditions**: Project with a single feature directory, no `.forge/state.json`, on main branch
- **Steps**:
  1. Run `forge feature`
  2. Verify feature resolves correctly via features-dir fallback
  3. Run existing unit tests for feature resolution
- **Expected**: All existing behavior preserved when state.json is absent
- **Priority**: P1

### Task 3: Verbose flag

## TC-015: Verbose shows state.json source
- **Source**: Proposal Success Criteria #5 / Task 3 AC #1
- **Type**: CLI
- **Target**: cli/feature-verbose
- **Test ID**: cli/feature-verbose/verbose-shows-state-json-source
- **Pre-conditions**: `.forge/state.json` with valid feature, feature directory exists
- **Steps**:
  1. Run `forge feature -v`
  2. Verify output contains `FEATURE: <slug> (from: state.json)`
- **Expected**: Verbose output includes resolution source as `state.json`
- **Priority**: P0

## TC-016: Verbose shows worktree source
- **Source**: Task 3 AC #2
- **Type**: CLI
- **Target**: cli/feature-verbose
- **Test ID**: cli/feature-verbose/verbose-shows-worktree-source
- **Pre-conditions**: In a git worktree with name matching a feature slug, no `.forge/state.json`
- **Steps**:
  1. Run `forge feature -v`
  2. Verify output contains `(from: worktree)`
- **Expected**: Verbose output shows `worktree` as resolution source
- **Priority**: P0

## TC-017: Verbose shows branch source
- **Source**: Task 3 AC #3
- **Type**: CLI
- **Target**: cli/feature-verbose
- **Test ID**: cli/feature-verbose/verbose-shows-branch-source
- **Pre-conditions**: On a git branch `feature/branch-feature`, feature directory exists, no `.forge/state.json`
- **Steps**:
  1. Run `forge feature -v`
  2. Verify output contains `(from: branch)`
- **Expected**: Verbose output shows `branch` as resolution source
- **Priority**: P0

## TC-018: Verbose shows features-dir source
- **Source**: Task 3 AC #4
- **Type**: CLI
- **Target**: cli/feature-verbose
- **Test ID**: cli/feature-verbose/verbose-shows-features-dir-source
- **Pre-conditions**: Single feature directory, no `.forge/state.json`, on main branch (no git feature context)
- **Steps**:
  1. Run `forge feature -v`
  2. Verify output contains `(from: features-dir)`
- **Expected**: Verbose output shows `features-dir` as resolution source
- **Priority**: P0

## TC-019: Verbose shows none when no feature set
- **Source**: Task 3 AC #5
- **Type**: CLI
- **Target**: cli/feature-verbose
- **Test ID**: cli/feature-verbose/verbose-shows-none-when-no-feature-set
- **Pre-conditions**: Project with no feature directories, no `.forge/state.json`
- **Steps**:
  1. Run `forge feature -v`
  2. Verify output contains `FEATURE: (none)`
- **Expected**: Verbose output shows `(none)` when no feature can be resolved
- **Priority**: P0

## TC-020: Verbose flag is local to feature command only
- **Source**: Task 3 AC #7 / Task 3 Hard Rules
- **Type**: CLI
- **Target**: cli/feature-verbose
- **Test ID**: cli/feature-verbose/verbose-flag-local-to-feature-command
- **Pre-conditions**: None
- **Steps**:
  1. Verify `forge feature -v` is recognized (no unknown flag error)
  2. Verify `forge feature set -v my-feature` is NOT recognized (unknown flag error or ignored)
  3. Verify `forge feature list -v` is NOT recognized
  4. Verify `forge feature status -v my-feature` is NOT recognized
- **Expected**: `-v` flag only applies to bare `forge feature`, not to subcommands
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Proposal SC #1,#2 / Task 1 AC #1,#2 | CLI | cli/feature-set | P0 |
| TC-002 | Proposal SC #3 / Task 1 AC #3 | CLI | cli/feature-set | P0 |
| TC-003 | Task 1 AC #4 | CLI | cli/feature-set | P0 |
| TC-004 | Task 1 AC #5 | CLI | cli/feature-positional | P0 |
| TC-005 | Task 1 Hard Rules | CLI | cli/feature-set | P1 |
| TC-006 | Proposal Key Scenarios | CLI | cli/feature-set | P1 |
| TC-007 | Proposal Key Scenarios | CLI | cli/feature-set | P1 |
| TC-008 | Proposal SC #4 / Task 2 AC #1 | CLI | cli/feature-query | P0 |
| TC-009 | Proposal SC #6 / Task 2 AC #2 | CLI | cli/feature-query | P0 |
| TC-010 | Task 2 AC #3 | CLI | cli/feature-query | P0 |
| TC-011 | Task 2 Hard Rules | CLI | cli/feature-query | P0 |
| TC-012 | Task 2 Hard Rules | CLI | cli/feature-query | P1 |
| TC-013 | Proposal Key Scenarios | CLI | cli/feature-query | P0 |
| TC-014 | Task 2 AC #5,#6 | CLI | cli/feature-query | P1 |
| TC-015 | Proposal SC #5 / Task 3 AC #1 | CLI | cli/feature-verbose | P0 |
| TC-016 | Task 3 AC #2 | CLI | cli/feature-verbose | P0 |
| TC-017 | Task 3 AC #3 | CLI | cli/feature-verbose | P0 |
| TC-018 | Task 3 AC #4 | CLI | cli/feature-verbose | P0 |
| TC-019 | Task 3 AC #5 | CLI | cli/feature-verbose | P0 |
| TC-020 | Task 3 AC #7 / Hard Rules | CLI | cli/feature-verbose | P1 |
