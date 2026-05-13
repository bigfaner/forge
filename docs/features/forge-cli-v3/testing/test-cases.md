---
feature: "forge-cli-v3"
sources:
  - docs/features/forge-cli-v3/prd/prd-user-stories.md
  - docs/features/forge-cli-v3/prd/prd-spec.md
generated: "2026-05-14"
---

# Test Cases: forge-cli-v3

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 41  |
| **Total** | **41** |

> **Note**: This is a CLI-only feature. No UI, API, or integration test cases apply.

---

## CLI Test Cases

### Command Discovery & Help

## TC-001: Help output shows correct command groups and top-level commands
- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/forge-help
- **Test ID**: cli/forge-help/help-output-shows-command-groups
- **Pre-conditions**: forge binary built and available in PATH
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge --help`
  2. Check exit code is 0
  3. Verify output contains 5 command groups: task, e2e, forensic, profile, prompt
  4. Verify output contains 5 top-level commands: feature, probe, cleanup, quality-gate, verify-task-done (plus version)
  5. Verify total entry count is <= 10
- **Expected**: Output displays exactly 5 command groups and 5 top-level commands (total 10 entries), exit code 0
- **Priority**: P0

## TC-002: Task subcommand help shows all 11 subcommands with proper descriptions
- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/task-help
- **Test ID**: cli/task-help/task-subcommand-help-shows-commands
- **Pre-conditions**: forge binary built and available in PATH
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task --help`
  2. Check exit code is 0
  3. Verify output lists 11 subcommands: claim, submit, status, query, check-deps, validate-index, verify-task-done, add, index, migrate, list-types
  4. For each subcommand, verify description follows "command-name + verb + object" pattern (e.g., "submit task execution result")
  5. Verify each description length <= 80 characters
- **Expected**: All 11 subcommands listed with self-describing names <= 80 chars each, exit code 0
- **Priority**: P0

## TC-003: Unknown top-level command returns error with suggestion
- **Source**: Story 1 / AC-2
- **Type**: CLI
- **Target**: cli/forge-help
- **Test ID**: cli/forge-help/unknown-command-returns-error-with-suggestion
- **Pre-conditions**: forge binary built and available in PATH
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge taks` (misspelled command)
  2. Check exit code is 1
  3. Verify stderr contains "unknown command" text
  4. Verify stderr contains a suggestion for the closest matching command
- **Expected**: Exit code 1, stderr contains "unknown command" and a command suggestion
- **Priority**: P1

## TC-004: Unknown task subcommand returns error with valid subcommand list
- **Source**: Story 1 / AC-3
- **Type**: CLI
- **Target**: cli/task-help
- **Test ID**: cli/task-help/unknown-subcommand-returns-error-with-list
- **Pre-conditions**: forge binary built and available in PATH
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task nonexistent-sub`
  2. Check exit code is 1
  3. Verify stderr contains "unknown subcommand" text
  4. Verify stderr lists valid subcommands
- **Expected**: Exit code 1, stderr contains "unknown subcommand" and lists all valid subcommands
- **Priority**: P1

### Prompt Commands

## TC-005: Get prompt by task ID returns correct prompt for implementation task
- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/prompt-get
- **Test ID**: cli/prompt-get/get-by-task-id-returns-correct-prompt
- **Pre-conditions**: Task T-impl-1 exists with type: implementation in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge prompt get-by-task-id T-impl-1`
  2. Check exit code is 0
  3. Verify output contains the implementation-type prompt
  4. Verify prompt contains substituted values for TASK_ID, TASK_FILE, SCOPE variables
- **Expected**: Exit code 0, output contains implementation prompt with variable substitutions applied
- **Priority**: P0

## TC-006: Get prompt for nonexistent task ID returns error
- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/prompt-get
- **Test ID**: cli/prompt-get/nonexistent-task-id-returns-error
- **Pre-conditions**: Task NONEXISTENT-999 does not exist in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge prompt get-by-task-id NONEXISTENT-999`
  2. Check exit code is 1
  3. Verify stderr contains "task not found" error message
  4. Verify stdout is empty
- **Expected**: Exit code 1, stderr contains "task not found", stdout empty
- **Priority**: P0

## TC-007: Get prompt for task with missing or invalid type returns error
- **Source**: Story 2 / AC-3
- **Type**: CLI
- **Target**: cli/prompt-get
- **Test ID**: cli/prompt-get/missing-or-invalid-type-returns-error
- **Pre-conditions**: A task exists in index.json with missing type field or type value not in known types list
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a task entry with missing or invalid type field
  2. Run `forge prompt get-by-task-id <id>` for that task
  3. Check exit code is 1
  4. Verify stderr contains "unknown task type" or "missing task type" error message
- **Expected**: Exit code 1, stderr contains appropriate error about task type
- **Priority**: P1

### Task Submit

## TC-008: Submit task success updates status and creates record
- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/success-updates-status-and-creates-record
- **Pre-conditions**: Task T-impl-1 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task submit T-impl-1 --result success --summary "completed implementation"`
  2. Check exit code is 0
  3. Verify index.json shows T-impl-1 status as completed
  4. Verify a record file exists in records/ directory for T-impl-1
- **Expected**: Exit code 0, index.json updated to completed, record file created
- **Priority**: P0

## TC-009: Submit task already in terminal state returns error
- **Source**: Story 3 / AC-2
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/already-terminal-state-returns-error
- **Pre-conditions**: Task T-impl-1 exists with status completed in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task submit T-impl-1 --result success --summary "retry"`
  2. Check exit code is 1
  3. Verify stderr contains "task already in terminal state" error message
  4. Verify index.json is unchanged
- **Expected**: Exit code 1, stderr contains "task already in terminal state", index.json unchanged
- **Priority**: P0

## TC-010: Submit task missing required --result flag returns error
- **Source**: Story 3 / AC-3
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/missing-result-flag-returns-error
- **Pre-conditions**: Task T-impl-1 exists in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task submit T-impl-1 --summary "missing result"`
  2. Check exit code is 1
  3. Verify stderr contains "required flag(s) not set: result" error message
- **Expected**: Exit code 1, stderr contains "required flag(s) not set: result"
- **Priority**: P0

## TC-011: Concurrent task submit handles lock contention correctly
- **Source**: Story 3 / AC-4
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/concurrent-submit-handles-lock-contention
- **Pre-conditions**: Task T-impl-1 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Simultaneously invoke `forge task submit T-impl-1 --result success --summary "agent-1"` and `forge task submit T-impl-1 --result success --summary "agent-2"` from two processes
  2. Verify exactly one process exits with code 0 (success)
  3. Verify the other process exits with code 1 and stderr contains "concurrent write conflict, retry"
  4. Verify index.json is valid JSON parseable by `jq .`
- **Expected**: One agent gets exit code 0, the other gets exit code 1 with "concurrent write conflict, retry"; index.json remains valid JSON
- **Priority**: P1

### Hook-Triggered Lifecycle

## TC-012: Cleanup removes terminal state task files
- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/cleanup
- **Test ID**: cli/cleanup/cleanup-removes-terminal-state-files
- **Pre-conditions**: At least one task in index.json is in completed, blocked, or rejected state; .forge/state.json exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge cleanup`
  2. Check exit code is 0
  3. Verify state files for completed/blocked/rejected tasks are removed
  4. Verify .forge/state.json is preserved unchanged
- **Expected**: Exit code 0, terminal task state files removed, .forge/state.json unchanged
- **Priority**: P0

## TC-013: Quality gate runs compile-fmt-lint-test sequence
- **Source**: Story 4 / AC-2
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/runs-compile-fmt-lint-test-sequence
- **Pre-conditions**: All tasks completed; project has compilable code
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge quality-gate`
  2. Verify the command executes compile, then fmt, then lint, then test in sequence
  3. If all pass: verify exit code is 0
  4. If any step fails: verify exit code is 1 and a P0 fix-task is created
- **Expected**: Sequential execution of compile->fmt->lint->test; exit code 0 if all pass, exit code 1 with fix-task creation if any fails
- **Priority**: P0

## TC-014: Cleanup with no terminal tasks outputs message
- **Source**: Story 4 / AC-3
- **Type**: CLI
- **Target**: cli/cleanup
- **Test ID**: cli/cleanup/cleanup-no-terminal-tasks-outputs-message
- **Pre-conditions**: No tasks in index.json are in completed, blocked, or rejected state
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge cleanup`
  2. Check exit code is 0
  3. Verify stdout contains "no tasks to clean up"
  4. Verify no files are deleted
- **Expected**: Exit code 0, stdout contains "no tasks to clean up", no files deleted
- **Priority**: P1

## TC-015: Quality gate creates new fix-task on repeated failure
- **Source**: Story 4 / AC-4
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/creates-new-fix-task-on-repeated-failure
- **Pre-conditions**: A fix-task exists but the corresponding fix also fails; quality-gate detects the same failing step
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge quality-gate` when a step fails that already has an existing fix-task
  2. Verify a new P0 fix-task is created (not overwriting the existing one)
  3. Verify the new fix-task title contains the failure step name and sequence number (e.g., "fix-compile-3")
  4. Verify grep for "fix-compile-" in index.json returns N results where N is the total count of fix-tasks for that step
- **Expected**: New fix-task created with incremented sequence number; existing fix-tasks preserved
- **Priority**: P1

## TC-016: Quality gate stops creating fix-tasks after max 3
- **Source**: Story 4 / AC-5
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/stops-creating-fix-tasks-after-max-3
- **Pre-conditions**: The same failing step already has 3 uncompleted fix-tasks in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge quality-gate` when the step fails again
  2. Check exit code is 1
  3. Verify stderr contains "max fix-tasks reached for <step>, manual intervention required"
  4. Verify no new fix-task is created
- **Expected**: Exit code 1, stderr contains max fix-tasks message, no new fix-task created
- **Priority**: P1

### E2E Test Commands

## TC-017: E2E run with configured profile executes correct suite
- **Source**: Story 5 / AC-1
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/with-configured-profile-executes-suite
- **Pre-conditions**: .forge/config.yaml has profile field set to a valid profile (e.g., web-playwright)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set profile in .forge/config.yaml
  2. Run `forge e2e run --feature my-feature`
  3. Verify stdout contains "profile: <profile-name>"
  4. Verify exit code matches the test suite result
- **Expected**: CLI reads config.yaml profile, outputs profile name, executes corresponding test suite, exit code reflects test result
- **Priority**: P0

## TC-018: E2E run with no profile configured returns error
- **Source**: Story 5 / AC-2
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/no-profile-configured-returns-error
- **Pre-conditions**: .forge/config.yaml has no profile field or profile field is empty
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Remove or empty the profile field in .forge/config.yaml
  2. Run `forge e2e run --feature my-feature`
  3. Check exit code is 1
  4. Verify stderr contains "no e2e profile configured" error message
- **Expected**: Exit code 1, stderr contains "no e2e profile configured"
- **Priority**: P0

## TC-019: E2E run with unknown profile returns error with valid profile list
- **Source**: Story 5 / AC-3
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/unknown-profile-returns-error-with-list
- **Pre-conditions**: .forge/config.yaml has profile field set to an unsupported value (e.g., "unknown-profile")
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set profile to "unknown-profile" in .forge/config.yaml
  2. Run `forge e2e run --feature my-feature`
  3. Check exit code is 1
  4. Verify stderr contains "unknown profile: unknown-profile"
  5. Verify stderr lists all supported valid profiles
- **Expected**: Exit code 1, stderr contains "unknown profile: unknown-profile" and lists supported profiles
- **Priority**: P0

## TC-020: E2E run with nonexistent feature returns error
- **Source**: Story 5 / AC-4
- **Type**: CLI
- **Target**: cli/e2e-run
- **Test ID**: cli/e2e-run/nonexistent-feature-returns-error
- **Pre-conditions**: No directory or test files exist for the specified feature name
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge e2e run --feature nonexistent-feature`
  2. Check exit code is 1
  3. Verify stderr contains "feature not found: nonexistent-feature"
- **Expected**: Exit code 1, stderr contains "feature not found: nonexistent-feature"
- **Priority**: P1

### Task Types

## TC-021: List types outputs all supported task types with descriptions
- **Source**: Story 6 / AC-1
- **Type**: CLI
- **Target**: cli/task-list-types
- **Test ID**: cli/task-list-types/list-types-outputs-all-with-descriptions
- **Pre-conditions**: Task types are registered in the system
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task list-types`
  2. Check exit code is 0
  3. Verify output contains all 11 task types (implementation, fix, gate, doc-generation.*, test-pipeline.*)
  4. Verify each type has a description in "verb + object" format
  5. Verify each description length <= 60 characters
- **Expected**: Exit code 0, all 11 types listed with verb+object descriptions <= 60 chars
- **Priority**: P0

## TC-022: List types with empty registry returns empty output
- **Source**: Story 6 / AC-2
- **Type**: CLI
- **Target**: cli/task-list-types
- **Test ID**: cli/task-list-types/empty-registry-returns-empty-output
- **Pre-conditions**: No task types are defined in the registry
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task list-types` with an empty type registry
  2. Check exit code is 0
  3. Verify stdout outputs empty list (0 lines of type records)
- **Expected**: Exit code 0, stdout contains no type records
- **Priority**: P2

### Forensic Commands

## TC-023: Forensic search scans history and returns matching sessions
- **Source**: Story 7 / AC-1
- **Type**: CLI
- **Target**: cli/forensic-search
- **Test ID**: cli/forensic-search/search-returns-matching-sessions
- **Pre-conditions**: history.jsonl exists with at least one recorded session
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge forensic search --project-path .`
  2. Check exit code is 0
  3. Verify output contains a list of matching sessions with session ID, timestamp, and skill name
- **Expected**: Exit code 0, output lists sessions with ID, timestamp, skill name
- **Priority**: P0

## TC-024: Forensic extract outputs evidence summary for valid session
- **Source**: Story 7 / AC-2
- **Type**: CLI
- **Target**: cli/forensic-extract
- **Test ID**: cli/forensic-extract/extract-outputs-evidence-summary
- **Pre-conditions**: A valid session JSONL file exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge forensic extract <valid-session-jsonl-path>`
  2. Check exit code is 0
  3. Verify output contains a compact evidence summary with key decision points, tool call sequences, and deviation points
- **Expected**: Exit code 0, output contains evidence summary with decision points, tool calls, and deviation nodes
- **Priority**: P0

## TC-025: Forensic subagents lists subagent transcripts
- **Source**: Story 7 / AC-3
- **Type**: CLI
- **Target**: cli/forensic-subagents
- **Test ID**: cli/forensic-subagents/subagents-lists-transcripts
- **Pre-conditions**: A valid session directory with subagent transcripts exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge forensic subagents <valid-session-dir-path>`
  2. Check exit code is 0
  3. Verify output lists all subagent transcript file paths and summary information
- **Expected**: Exit code 0, output lists subagent transcript paths and summaries
- **Priority**: P0

## TC-026: Forensic extract with nonexistent path returns error
- **Source**: Story 7 / AC-4
- **Type**: CLI
- **Target**: cli/forensic-extract
- **Test ID**: cli/forensic-extract/nonexistent-path-returns-error
- **Pre-conditions**: The specified file path does not exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge forensic extract /nonexistent/path.jsonl`
  2. Check exit code is 1
  3. Verify stderr contains "file not found: /nonexistent/path.jsonl"
- **Expected**: Exit code 1, stderr contains "file not found: /nonexistent/path.jsonl"
- **Priority**: P1

### Profile Commands

## TC-027: Profile detect scans project and outputs detected profiles
- **Source**: Story 8 / AC-1
- **Type**: CLI
- **Target**: cli/profile-detect
- **Test ID**: cli/profile-detect/detect-scans-and-outputs-profiles
- **Pre-conditions**: Project directory contains test framework config files (e.g., playwright.config.ts, *_test.go)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge profile detect`
  2. Check exit code is 0
  3. Verify output lists detected profiles with detection evidence (e.g., "web-playwright: found playwright.config.ts")
- **Expected**: Exit code 0, output shows detected profiles with evidence
- **Priority**: P0

## TC-028: Profile set updates config.yaml with valid profile
- **Source**: Story 8 / AC-2
- **Type**: CLI
- **Target**: cli/profile-set
- **Test ID**: cli/profile-set/set-updates-config-with-valid-profile
- **Pre-conditions**: .forge/config.yaml exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge profile set web-playwright`
  2. Check exit code is 0
  3. Verify .forge/config.yaml profile field is updated to "web-playwright"
- **Expected**: Exit code 0, config.yaml profile field set to "web-playwright"
- **Priority**: P0

## TC-029: Profile get outputs strategy file content for valid profile
- **Source**: Story 8 / AC-3
- **Type**: CLI
- **Target**: cli/profile-get
- **Test ID**: cli/profile-get/get-outputs-strategy-file-content
- **Pre-conditions**: The specified profile exists with strategy files
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge profile get web-playwright`
  2. Check exit code is 0
  3. Verify output contains the profile's strategy file content (generate.md, run.md, etc.)
- **Expected**: Exit code 0, output shows strategy file content for the profile
- **Priority**: P1

## TC-030: Profile set with invalid profile returns error with valid list
- **Source**: Story 8 / AC-4
- **Type**: CLI
- **Target**: cli/profile-set
- **Test ID**: cli/profile-set/invalid-profile-returns-error-with-list
- **Pre-conditions**: .forge/config.yaml exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge profile set nonexistent-profile`
  2. Check exit code is 1
  3. Verify stderr contains "unknown profile: nonexistent-profile"
  4. Verify stderr lists all supported valid profiles
- **Expected**: Exit code 1, stderr contains "unknown profile: nonexistent-profile" and lists supported profiles
- **Priority**: P0

### Error Handling (from Spec)

## TC-031: Task claim with no available tasks returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/no-available-tasks-returns-error
- **Pre-conditions**: All tasks are already claimed or in terminal state
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task claim` when no tasks are available
  2. Check exit code is 1
  3. Verify stderr contains "no available tasks to claim"
- **Expected**: Exit code 1, stderr contains "no available tasks to claim"
- **Priority**: P1

## TC-032: Task claim with corrupted index.json returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/corrupted-index-returns-error
- **Pre-conditions**: index.json does not exist or has invalid JSON format
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Remove or corrupt index.json
  2. Run `forge task claim`
  3. Check exit code is 1
  4. Verify stderr contains "failed to load index:" with reason
- **Expected**: Exit code 1, stderr contains "failed to load index:" with failure reason
- **Priority**: P1

## TC-033: Task check-deps with unmet dependency returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-check-deps
- **Test ID**: cli/task-check-deps/unmet-dependency-returns-error
- **Pre-conditions**: A task has a dependency on another task that is not completed
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task check-deps` for a task with unmet dependency
  2. Check exit code is 1
  3. Verify stderr contains "dependency not met:" with the dependency ID and its status
- **Expected**: Exit code 1, stderr contains "dependency not met: <dep-id> is <status>"
- **Priority**: P1

## TC-034: Task validate-index with invalid schema returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-validate-index
- **Test ID**: cli/task-validate-index/invalid-schema-returns-error
- **Pre-conditions**: index.json has schema validation errors
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task validate-index` with an index.json that fails schema validation
  2. Check exit code is 1
  3. Verify stderr contains "index validation failed:" with details
- **Expected**: Exit code 1, stderr contains "index validation failed:" with error details
- **Priority**: P1

## TC-035: Task status with nonexistent ID returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-status
- **Test ID**: cli/task-status/nonexistent-id-returns-error
- **Pre-conditions**: The specified task ID does not exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge task status NONEXISTENT-999`
  2. Check exit code is 1
  3. Verify stderr contains "task not found: NONEXISTENT-999"
- **Expected**: Exit code 1, stderr contains "task not found: NONEXISTENT-999"
- **Priority**: P1

## TC-036: Forensic search with no results returns empty output
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/forensic-search
- **Test ID**: cli/forensic-search/no-results-returns-empty-output
- **Pre-conditions**: history.jsonl exists but no sessions match the search criteria
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge forensic search --project-path .` with no matching sessions
  2. Check exit code is 0
  3. Verify stderr contains "no results found"
- **Expected**: Exit code 0, stderr contains "no results found"
- **Priority**: P2

## TC-037: Forensic search with missing records directory returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/forensic-search
- **Test ID**: cli/forensic-search/missing-records-dir-returns-error
- **Pre-conditions**: records/ directory does not exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge forensic search` when records/ directory is absent
  2. Check exit code is 1
  3. Verify stderr contains "records directory not found"
- **Expected**: Exit code 1, stderr contains "records directory not found"
- **Priority**: P2

## TC-038: Verify-task-done with incomplete tasks returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/verify-task-done
- **Test ID**: cli/verify-task-done/incomplete-tasks-returns-error
- **Pre-conditions**: At least one task in the current feature is not in terminal state
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge verify-task-done` with incomplete tasks present
  2. Check exit code is 1
  3. Verify stderr contains "incomplete tasks found:" with count
- **Expected**: Exit code 1, stderr contains "incomplete tasks found: <count>"
- **Priority**: P1

## TC-039: Task submit with concurrent write conflict returns retry error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/concurrent-write-conflict-returns-retry-error
- **Pre-conditions**: Task exists in in_progress state; another process holds the lock
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Simulate a lock contention scenario
  2. Run `forge task submit <id> --result success --summary "..."`
  3. Check exit code is 1
  4. Verify stderr contains "concurrent write conflict, retry"
- **Expected**: Exit code 1, stderr contains "concurrent write conflict, retry"
- **Priority**: P1

## TC-040: Task submit with missing index.json returns error
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/missing-index-returns-error
- **Pre-conditions**: index.json does not exist for the feature
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Remove index.json for the current feature
  2. Run `forge task submit <id> --result success --summary "..."`
  3. Check exit code is 1
  4. Verify stderr contains "index not found for feature:"
- **Expected**: Exit code 1, stderr contains "index not found for feature: <slug>"
- **Priority**: P1

## TC-041: Profile get with invalid profile returns error with valid list
- **Source**: Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/profile-get
- **Test ID**: cli/profile-get/invalid-profile-returns-error-with-list
- **Pre-conditions**: The specified profile does not exist
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `forge profile get nonexistent-profile`
  2. Check exit code is 1
  3. Verify stderr contains "unknown profile: nonexistent-profile"
  4. Verify stderr lists all supported valid profiles
- **Expected**: Exit code 1, stderr contains "unknown profile: nonexistent-profile" and lists supported profiles
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/forge-help | P0 |
| TC-002 | Story 1 / AC-1 | CLI | cli/task-help | P0 |
| TC-003 | Story 1 / AC-2 | CLI | cli/forge-help | P1 |
| TC-004 | Story 1 / AC-3 | CLI | cli/task-help | P1 |
| TC-005 | Story 2 / AC-1 | CLI | cli/prompt-get | P0 |
| TC-006 | Story 2 / AC-2 | CLI | cli/prompt-get | P0 |
| TC-007 | Story 2 / AC-3 | CLI | cli/prompt-get | P1 |
| TC-008 | Story 3 / AC-1 | CLI | cli/task-submit | P0 |
| TC-009 | Story 3 / AC-2 | CLI | cli/task-submit | P0 |
| TC-010 | Story 3 / AC-3 | CLI | cli/task-submit | P0 |
| TC-011 | Story 3 / AC-4 | CLI | cli/task-submit | P1 |
| TC-012 | Story 4 / AC-1 | CLI | cli/cleanup | P0 |
| TC-013 | Story 4 / AC-2 | CLI | cli/quality-gate | P0 |
| TC-014 | Story 4 / AC-3 | CLI | cli/cleanup | P1 |
| TC-015 | Story 4 / AC-4 | CLI | cli/quality-gate | P1 |
| TC-016 | Story 4 / AC-5 | CLI | cli/quality-gate | P1 |
| TC-017 | Story 5 / AC-1 | CLI | cli/e2e-run | P0 |
| TC-018 | Story 5 / AC-2 | CLI | cli/e2e-run | P0 |
| TC-019 | Story 5 / AC-3 | CLI | cli/e2e-run | P0 |
| TC-020 | Story 5 / AC-4 | CLI | cli/e2e-run | P1 |
| TC-021 | Story 6 / AC-1 | CLI | cli/task-list-types | P0 |
| TC-022 | Story 6 / AC-2 | CLI | cli/task-list-types | P2 |
| TC-023 | Story 7 / AC-1 | CLI | cli/forensic-search | P0 |
| TC-024 | Story 7 / AC-2 | CLI | cli/forensic-extract | P0 |
| TC-025 | Story 7 / AC-3 | CLI | cli/forensic-subagents | P0 |
| TC-026 | Story 7 / AC-4 | CLI | cli/forensic-extract | P1 |
| TC-027 | Story 8 / AC-1 | CLI | cli/profile-detect | P0 |
| TC-028 | Story 8 / AC-2 | CLI | cli/profile-set | P0 |
| TC-029 | Story 8 / AC-3 | CLI | cli/profile-get | P1 |
| TC-030 | Story 8 / AC-4 | CLI | cli/profile-set | P0 |
| TC-031 | Spec Error Handling | CLI | cli/task-claim | P1 |
| TC-032 | Spec Error Handling | CLI | cli/task-claim | P1 |
| TC-033 | Spec Error Handling | CLI | cli/task-check-deps | P1 |
| TC-034 | Spec Error Handling | CLI | cli/task-validate-index | P1 |
| TC-035 | Spec Error Handling | CLI | cli/task-status | P1 |
| TC-036 | Spec Error Handling | CLI | cli/forensic-search | P2 |
| TC-037 | Spec Error Handling | CLI | cli/forensic-search | P2 |
| TC-038 | Spec Error Handling | CLI | cli/verify-task-done | P1 |
| TC-039 | Spec Error Handling | CLI | cli/task-submit | P1 |
| TC-040 | Spec Error Handling | CLI | cli/task-submit | P1 |
| TC-041 | Spec Error Handling | CLI | cli/profile-get | P1 |
