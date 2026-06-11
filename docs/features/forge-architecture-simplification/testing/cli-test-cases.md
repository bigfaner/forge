---
feature: "forge-architecture-simplification"
type: CLI
generated: "2026-05-22"
---

# CLI Test Cases: Forge Architecture Simplification

## TC-001: Concurrent Claim Lock — Exactly One Succeeds
- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/concurrent-claim-lock
- **Pre-conditions**: A task is available and unclaimed; two concurrent `forge task claim` invocations target the same task ID
- **Steps**:
  1. Run `forge task claim <id>` concurrently from two terminals/goroutines
  2. Wait for both to complete
- **Expected**: Exactly one invocation succeeds (exit code 0, task assigned). The other returns exit code 1 with AIError containing lock conflict message (Hint: not a corrupted index.json).
- **Priority**: P0

## TC-002: Concurrent Index Write — Both Changes Applied
- **Source**: Story 1 / AC-2
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/concurrent-index-write
- **Pre-conditions**: One process running `forge task index` while another runs `forge task submit` on a different task in the same index
- **Steps**:
  1. Start `forge task index` (long-running read/reindex)
  2. Concurrently run `forge task submit --result success` on a task
  3. Wait for both to complete
- **Expected**: index.json is valid JSON. Both changes (from reindex and submit) are reflected. No truncation or data loss.
- **Priority**: P0

## TC-003: Malformed Index — AIError with Corruption Hint
- **Source**: Story 1 / AC-3
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/malformed-index-returns-hint
- **Pre-conditions**: index.json is corrupted (e.g., truncated JSON, invalid syntax)
- **Steps**:
  1. Run `forge task claim`
- **Expected**: Exit code 1. Output is an AIError with Hint containing "index.json may be corrupted, run forge task index to rebuild". No panic.
- **Priority**: P0

## TC-004: Lock Timeout — Retry Action
- **Source**: Story 1 / AC-4
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/lock-timeout-retry
- **Pre-conditions**: Lock on index.json is held by another process and not released within 5s
- **Steps**:
  1. Hold the advisory lock on index.json
  2. Run `forge task submit --result success` on a task
- **Expected**: Exit code 1. AIError with Action "retry in a few seconds".
- **Priority**: P0

## TC-005: Submit on Completed Task — Error
- **Source**: Story 2 / AC-5 (Story 2 / AC-1)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/completed-task-rejected
- **Pre-conditions**: A task exists with status `completed`
- **Steps**:
  1. Run `forge task submit --result success` on the completed task
- **Expected**: Exit code 1. Error message: "task already completed, create a subtask if re-work needed". The task state remains `completed`.
- **Priority**: P0

## TC-006: Reopen Rejected Task — Transitions to Pending
- **Source**: Story 2 / AC-6 (Story 7 / AC-1)
- **Type**: CLI
- **Target**: cli/task-reopen
- **Test ID**: cli/task-reopen/rejected-task-becomes-pending
- **Pre-conditions**: A task exists with status `rejected`
- **Steps**:
  1. Run `forge task reopen <id>`
- **Expected**: Exit code 0. Task status changes to `pending`.
- **Priority**: P0

## TC-007: Reopen Completed Task — Error
- **Source**: Story 2 / AC-7 (Story 7 / AC-3)
- **Type**: CLI
- **Target**: cli/task-reopen
- **Test ID**: cli/task-reopen/completed-task-error
- **Pre-conditions**: A task exists with status `completed`
- **Steps**:
  1. Run `forge task reopen <id>`
- **Expected**: Exit code 1. Error message: "task already completed, create a subtask if re-work needed".
- **Priority**: P0

## TC-008: Status Command Read-Only — Rejects Arguments
- **Source**: Story 2 / AC-8
- **Type**: CLI
- **Target**: cli/task-status
- **Test ID**: cli/task-status/rejects-status-argument
- **Pre-conditions**: A task exists
- **Steps**:
  1. Run `forge task status <id>` (with only the task ID, no status argument), or with a status argument
- **Expected**: If status argument is provided, exit code 1 with error "task status is read-only. Use forge task submit to complete a task."
- **Priority**: P0

## TC-009: Worktree Remove Error — Structured AIError Output
- **Source**: Story 3 / AC-9
- **Type**: CLI
- **Target**: cli/worktree-remove
- **Test ID**: cli/worktree-remove/structured-error-output
- **Pre-conditions**: Worktree has uncommitted changes that prevent removal
- **Steps**:
  1. Run `forge worktree remove <name>`
- **Expected**: Exit code 1. Error output includes Code, Message, Cause, Hint, and Action fields (AIError format, not raw `fmt.Errorf`).
- **Priority**: P1

## TC-010: Submit Lock Conflict — AIError Format
- **Source**: Story 3 / AC-10
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/lock-conflict-aierror
- **Pre-conditions**: Lock conflict occurs during submit
- **Steps**:
  1. Induce a lock conflict on `forge task submit`
  2. Run `forge task submit --result success`
- **Expected**: Exit code 1. Output is AIError with Hint "retry in a few seconds". Not raw stderr + exit 1.
- **Priority**: P1

## TC-011: Consistent Error Format Across All Commands
- **Source**: Story 3 / AC-11
- **Type**: CLI
- **Target**: cli/any-command
- **Test ID**: cli/all-commands/consistent-error-format
- **Pre-conditions**: Any forge command is invoked in a way that produces an error
- **Steps**:
  1. Invoke multiple forge commands that produce errors (e.g., `forge task submit` on completed, `forge worktree remove` with uncommitted changes, `forge config get` with no config)
  2. Inspect error output across all commands
- **Expected**: All error outputs follow the same AIError structure (Code, Message, Cause, Hint, Action). No mixed `fmt.Errorf`/AIError formats.
- **Priority**: P1

## TC-012: Eval Max Iterations — Original Documents Restored
- **Source**: Story 4 / AC-12
- **Type**: CLI
- **Target**: cli/eval
- **Test ID**: cli/eval/max-iterations-restore-backup
- **Pre-conditions**: Eval pipeline configured with a max iteration count; documents exist at Step 1
- **Steps**:
  1. Run an eval that will reach max iterations without meeting target score
  2. Check final document state after eval completes
- **Expected**: Original documents are restored from the Step 1 backup. Reviser changes are rolled back.
- **Priority**: P1

## TC-013: Eval Scorer Malformed Output — Pipeline Halts
- **Source**: Story 4 / AC-13
- **Type**: CLI
- **Target**: cli/eval
- **Test ID**: cli/eval/scorer-malformed-halts
- **Pre-conditions**: The eval scorer produces malformed/unparseable output
- **Steps**:
  1. Run an eval iteration where the scorer returns malformed data
- **Expected**: Pipeline halts with an error message. No crash or silent ignore.
- **Priority**: P1

## TC-014: Eval Reviser — Same Project Context as Scorer
- **Source**: Story 4 / AC-14
- **Type**: CLI
- **Target**: cli/eval
- **Test ID**: cli/eval/reviser-gets-project-context
- **Pre-conditions**: Eval pipeline is configured with project conventions and business rules
- **Steps**:
  1. Run an eval iteration
  2. Inspect the prompt provided to the reviser
- **Expected**: The reviser receives the same project context (conventions, business rules) as the scorer.
- **Priority**: P1

## TC-015: Config Set and Get Round-Trip
- **Source**: Story 5 / AC-15
- **Type**: CLI
- **Target**: cli/config-set
- **Test ID**: cli/config-set/set-and-get-roundtrip
- **Pre-conditions**: Valid forge project config exists
- **Steps**:
  1. Run `forge config set auto.cleanCode true`
  2. Run `forge config get auto.cleanCode`
- **Expected**: The second command returns `true`.
- **Priority**: P1

## TC-016: Config Get for Non-GitPush Field
- **Source**: Story 5 / AC-16
- **Type**: CLI
- **Target**: cli/config-get
- **Test ID**: cli/config-get/e2etest-field-queryable
- **Pre-conditions**: Config file has `auto.e2eTest` field set to a value
- **Steps**:
  1. Run `forge config get auto.e2eTest`
- **Expected**: The current value is displayed.
- **Priority**: P1

## TC-017: Config Init — All Auto Fields Configured
- **Source**: Story 5 / AC-17
- **Type**: CLI
- **Target**: cli/config-init
- **Test ID**: cli/config-init/all-auto-fields-configured
- **Pre-conditions**: No existing config, or fresh project
- **Steps**:
  1. Run `forge config init`
  2. Complete the wizard flow
  3. Run `forge config get auto.e2eTest`
  4. Run `forge config get auto.consolidateSpecs`
  5. Run `forge config get auto.cleanCode`
  6. Run `forge config get auto.gitPush`
- **Expected**: All 4 auto fields (e2eTest, consolidateSpecs, cleanCode, gitPush) are configured and queryable.
- **Priority**: P1

## TC-018: Config Get on Empty Config — Meaningful Error
- **Source**: Story 5 / AC-18
- **Type**: CLI
- **Target**: cli/config-get
- **Test ID**: cli/config-get/empty-config-returns-error
- **Pre-conditions**: config.yaml is empty or missing
- **Steps**:
  1. Run `forge config get auto.gitPush`
- **Expected**: Exit code 1. A meaningful error is returned (not a panic or empty output).
- **Priority**: P1

## TC-019: Quality Gate Fix-Task — Real SourceTaskID
- **Source**: Story 6 / AC-19
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/fix-task-real-source-id
- **Pre-conditions**: Quality gate creates a fix-task for step "2.1"
- **Steps**:
  1. Trigger a quality gate that creates a fix-task
  2. Inspect the created fix-task's SourceTaskID field
- **Expected**: SourceTaskID is the actual blocked task ID (not `"quality-gate:2.1"` sentinel).
- **Priority**: P1

## TC-020: Quality Gate Fix-Task Cap — Active Only
- **Source**: Story 6 / AC-20
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/fix-task-cap-active-only
- **Pre-conditions**: 3 fix-tasks exist for a step but all are completed
- **Steps**:
  1. Quality gate evaluates whether to create another fix-task
- **Expected**: Creation is allowed (cap counts active fix-tasks only, not lifetime).
- **Priority**: P1

## TC-021: Quality Gate No Feature — Exit Code 1
- **Source**: Story 6 / AC-21
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/no-feature-exit-1
- **Pre-conditions**: Quality gate is invoked with no feature configured
- **Steps**:
  1. Run the quality gate command without a feature
- **Expected**: Exit code 1 with an error message (not silent exit 0).
- **Priority**: P1

## TC-022: Reopen Skipped Task — Transitions to Pending
- **Source**: Story 7 / AC-2
- **Type**: CLI
- **Target**: cli/task-reopen
- **Test ID**: cli/task-reopen/skipped-task-becomes-pending
- **Pre-conditions**: A task exists with status `skipped`
- **Steps**:
  1. Run `forge task reopen <id>`
- **Expected**: Exit code 0. Task status changes to `pending`.
- **Priority**: P0

## TC-023: Reopen In-Progress Task — Error
- **Source**: Story 7 / AC-4
- **Type**: CLI
- **Target**: cli/task-reopen
- **Test ID**: cli/task-reopen/in-progress-task-error
- **Pre-conditions**: A task exists with status `in_progress`
- **Steps**:
  1. Run `forge task reopen <id>`
- **Expected**: Exit code 1. Error message: "task is not rejected or skipped".
- **Priority**: P1

## TC-024: Index Write Atomicity — Temp+Rename Pattern
- **Source**: Spec DR-1
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/atomic-write-temp-rename
- **Pre-conditions**: Index write operation is triggered (submit, claim, status, add, build, migrate)
- **Steps**:
  1. Trigger an index write from any forge command
  2. Inspect the write mechanism at the code level
- **Expected**: All index writers use temp file + rename (atomic write), not `os.WriteFile`.
- **Priority**: P1

## TC-025: Index Write Concurrency — Advisory Lock
- **Source**: Spec DR-2
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/advisory-lock-on-write
- **Pre-conditions**: Any index write operation is triggered
- **Steps**:
  1. Trigger index writes from claim, submit, status, add, build, migrate
  2. Verify lock acquisition at the code level
- **Expected**: All index writers acquire an advisory file lock before writing.
- **Priority**: P1

## TC-026: SaveState Atomic — Temp+Rename
- **Source**: Spec DR-3
- **Type**: CLI
- **Target**: cli/state
- **Test ID**: cli/state/savestate-atomic
- **Pre-conditions**: A state save operation is triggered
- **Steps**:
  1. Trigger `SaveState` call
  2. Inspect write mechanism
- **Expected**: `SaveState` uses temp+rename, not `os.WriteFile`.
- **Priority**: P1

## TC-027: ClearForgeState Writes False — No Deletion
- **Source**: Spec DR-4
- **Type**: CLI
- **Target**: cli/state
- **Test ID**: cli/state/clearforgestate-writes-false
- **Pre-conditions**: A `ClearForgeState` call is triggered
- **Steps**:
  1. Trigger ClearForgeState
  2. Check state file contents after operation
- **Expected**: State file contains `false` instead of being deleted.
- **Priority**: P1

## TC-028: Auto-Downgrade Sets BlockedReason
- **Source**: Spec BC-4
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/auto-downgrade-blocked-reason
- **Pre-conditions**: Submit is called with testsFailed >= threshold, triggering auto-downgrade
- **Steps**:
  1. Run `forge task submit --result blocked` or trigger auto-downgrade
- **Expected**: Task status is `blocked`. BlockedReason is set to a descriptive value (e.g., "auto-downgrade: testsFailed=2").
- **Priority**: P1

## TC-029: Orphan Cleanup — Default Clean with Warning
- **Source**: Spec BC-8
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/orphan-cleanup-with-warning
- **Pre-conditions**: Orphan entries exist in the index
- **Steps**:
  1. Run `forge task index` which triggers orphan cleanup
- **Expected**: Orphans are cleaned up by default. A warning is printed for each orphan removed.
- **Priority**: P1

## TC-030: Test Promote Rejects Path Traversal
- **Source**: Spec EC-6
- **Type**: CLI
- **Target**: cli/test-promote
- **Test ID**: cli/test-promote/rejects-path-traversal
- **Pre-conditions**: Test promote is invoked with a path containing `../`
- **Steps**:
  1. Run `forge test promote <path-with-../>`
- **Expected**: Exit code 1. Input validation rejects the path traversal. AIError with appropriate message.
- **Priority**: P1

## TC-031: Test Verify Parse Failure — Returns Error
- **Source**: Spec EC-7
- **Type**: CLI
- **Target**: cli/test-verify
- **Test ID**: cli/test-verify/parse-failure-returns-error
- **Pre-conditions**: Contract parsing fails during test verify
- **Steps**:
  1. Run `forge test verify` on contracts that fail to parse
- **Expected**: Error is returned (not silent return of zero value).
- **Priority**: P1

## TC-032: Test Verify No Fact Table — Unverifiable
- **Source**: Spec EC-8
- **Type**: CLI
- **Target**: cli/test-verify
- **Test ID**: cli/test-verify/no-fact-table-unverifiable
- **Pre-conditions**: A test result has no Fact Table
- **Steps**:
  1. Run `forge test verify`
- **Expected**: The result is marked as "unverifiable", not silently marked as OK.
- **Priority**: P1

## TC-033: Eval Parse Failure — Abort with Error
- **Source**: Spec ES-1
- **Type**: CLI
- **Target**: cli/eval
- **Test ID**: cli/eval/parse-failure-abort
- **Pre-conditions**: Eval pipeline encounters an unparseable result
- **Steps**:
  1. Run eval pipeline with malformed input
- **Expected**: Pipeline aborts and outputs an error (does not crash or ignore).
- **Priority**: P1

## TC-034: Quality Gate Returns Error on Infra Failure
- **Source**: Spec EC-3
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/infra-failure-returns-error
- **Pre-conditions**: Quality gate encounters an infrastructure error (e.g., file read failure)
- **Steps**:
  1. Trigger quality gate in an environment with infra failure
- **Expected**: Exit code 1. Error is returned (not silently swallowed).
- **Priority**: P1

## TC-035: Test Command Uses AIError
- **Source**: Spec EC-5
- **Type**: CLI
- **Target**: cli/test
- **Test ID**: cli/test/error-uses-aierror
- **Pre-conditions**: A forge test command encounters an error
- **Steps**:
  1. Run `forge test` in a way that produces an error
- **Expected**: Error output is AIError format (not raw `os.Exit(1)`).
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/task-claim | P0 |
| TC-002 | Story 1 / AC-2 | CLI | cli/task-index | P0 |
| TC-003 | Story 1 / AC-3 | CLI | cli/task-claim | P0 |
| TC-004 | Story 1 / AC-4 | CLI | cli/task-submit | P0 |
| TC-005 | Story 2 / AC-5 | CLI | cli/task-submit | P0 |
| TC-006 | Story 2 / AC-6 | CLI | cli/task-reopen | P0 |
| TC-007 | Story 2 / AC-7 | CLI | cli/task-reopen | P0 |
| TC-008 | Story 2 / AC-8 | CLI | cli/task-status | P0 |
| TC-009 | Story 3 / AC-9 | CLI | cli/worktree-remove | P1 |
| TC-010 | Story 3 / AC-10 | CLI | cli/task-submit | P1 |
| TC-011 | Story 3 / AC-11 | CLI | cli/any-command | P1 |
| TC-012 | Story 4 / AC-12 | CLI | cli/eval | P1 |
| TC-013 | Story 4 / AC-13 | CLI | cli/eval | P1 |
| TC-014 | Story 4 / AC-14 | CLI | cli/eval | P1 |
| TC-015 | Story 5 / AC-15 | CLI | cli/config-set | P1 |
| TC-016 | Story 5 / AC-16 | CLI | cli/config-get | P1 |
| TC-017 | Story 5 / AC-17 | CLI | cli/config-init | P1 |
| TC-018 | Story 5 / AC-18 | CLI | cli/config-get | P1 |
| TC-019 | Story 6 / AC-19 | CLI | cli/quality-gate | P1 |
| TC-020 | Story 6 / AC-20 | CLI | cli/quality-gate | P1 |
| TC-021 | Story 6 / AC-21 | CLI | cli/quality-gate | P1 |
| TC-022 | Story 7 / AC-2 | CLI | cli/task-reopen | P0 |
| TC-023 | Story 7 / AC-4 | CLI | cli/task-reopen | P1 |
| TC-024 | Spec DR-1 | CLI | cli/task-index | P1 |
| TC-025 | Spec DR-2 | CLI | cli/task-submit | P1 |
| TC-026 | Spec DR-3 | CLI | cli/state | P1 |
| TC-027 | Spec DR-4 | CLI | cli/state | P1 |
| TC-028 | Spec BC-4 | CLI | cli/task-submit | P1 |
| TC-029 | Spec BC-8 | CLI | cli/task-index | P1 |
| TC-030 | Spec EC-6 | CLI | cli/test-promote | P1 |
| TC-031 | Spec EC-7 | CLI | cli/test-verify | P1 |
| TC-032 | Spec EC-8 | CLI | cli/test-verify | P1 |
| TC-033 | Spec ES-1 | CLI | cli/eval | P1 |
| TC-034 | Spec EC-3 | CLI | cli/quality-gate | P1 |
| TC-035 | Spec EC-5 | CLI | cli/test | P1 |