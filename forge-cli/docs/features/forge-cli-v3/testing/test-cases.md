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
| **Integration** | **6** |
| API  | 0  |
| CLI  | 75  |
| **Total** | **81** |

> **Note**: This is a CLI-only feature. No UI or API test cases apply. Integration TCs exercise multi-command workflows defined in the PRD flow descriptions.

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

### Task Claim (Happy Path)

## TC-042: Task claim assigns next available pending task
- **Source**: Spec Agent Task Execution Flow (claim step)
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/claim-assigns-next-available-pending-task
- **Pre-conditions**: index.json exists with at least one task in status pending; no tasks are in_progress
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task claim`
  2. Check exit code is 0
  3. Verify stdout contains the claimed task ID
  4. Run `cat index.json | jq '.tasks[] | select(.id=="<claimed-id>") | .status'` to confirm status is in_progress
- **Expected**: Exit code 0, stdout shows claimed task ID, index.json updated to in_progress for that task
- **Priority**: P0

### Task Add

## TC-043: Task add creates a new task entry in index.json
- **Source**: Spec Command Structure Table (task add)
- **Type**: CLI
- **Target**: cli/task-add
- **Test ID**: cli/task-add/add-creates-new-task-entry
- **Pre-conditions**: index.json exists; feature context is set
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task add --title "my-new-task" --type feature`
  2. Check exit code is 0
  3. Run `cat index.json | jq '.tasks | length'` to verify task count increased by 1
  4. Run `cat index.json | jq '.tasks[] | select(.title=="my-new-task")'` to verify new task exists with status pending and type feature
- **Expected**: Exit code 0, index.json contains new task with title "my-new-task", type "feature", status "pending"
- **Priority**: P1

### Task Index

## TC-044: Task index generates index.json from task markdown files
- **Source**: Spec Command Structure Table (task index)
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/index-generates-from-markdown
- **Pre-conditions**: Feature directory contains task markdown files but no index.json (or stale index.json)
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Remove existing index.json if present
  2. Run `forge task index`
  3. Check exit code is 0
  4. Verify index.json is created and is valid JSON parseable by `jq .`
  5. Verify each task markdown file has a corresponding entry in index.json with correct id, title, type, and status fields
- **Expected**: Exit code 0, index.json created with entries for all task markdown files, valid JSON structure
- **Priority**: P1

### Task Migrate

## TC-045: Task migrate updates index.json schema to current version
- **Source**: Spec Command Structure Table (task migrate)
- **Type**: CLI
- **Target**: cli/task-migrate
- **Test ID**: cli/task-migrate/migrate-updates-schema-version
- **Pre-conditions**: index.json exists with an older schema version format
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task migrate`
  2. Check exit code is 0
  3. Verify index.json schema version field is updated to the current version string
  4. Verify all existing task entries are preserved with correct field mappings
- **Expected**: Exit code 0, index.json schema version updated, all task data preserved
- **Priority**: P2

### Task Query

## TC-046: Task query filters tasks by status
- **Source**: Spec Command Structure Table (task query)
- **Type**: CLI
- **Target**: cli/task-query
- **Test ID**: cli/task-query/query-filters-by-status
- **Pre-conditions**: index.json contains tasks in multiple statuses (pending, in_progress, completed)
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task query --status pending`
  2. Check exit code is 0
  3. Verify every task in stdout output has status pending
  4. Verify no tasks with status in_progress or completed appear in output
  5. Run `forge task query --status completed` and verify only completed tasks appear
- **Expected**: Exit code 0, output contains only tasks matching the queried status; filtering is accurate for each status value
- **Priority**: P1

### Task Status (Happy Path)

## TC-047: Task status displays current status for valid task ID
- **Source**: Spec Command Structure Table (task status)
- **Type**: CLI
- **Target**: cli/task-status
- **Test ID**: cli/task-status/status-displays-for-valid-id
- **Pre-conditions**: Task T-impl-1 exists in index.json with status in_progress
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task status T-impl-1`
  2. Check exit code is 0
  3. Verify stdout contains the task ID "T-impl-1"
  4. Verify stdout contains status "in_progress"
- **Expected**: Exit code 0, stdout shows task ID and current status
- **Priority**: P0

### Feature Command

## TC-048: Feature get displays current feature context
- **Source**: Spec Top-Level Commands Table (feature)
- **Type**: CLI
- **Target**: cli/feature
- **Test ID**: cli/feature/get-displays-current-context
- **Pre-conditions**: Feature context is set (e.g., via previous `forge feature <name>` call)
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge feature`
  2. Check exit code is 0
  3. Verify stdout contains the current feature name/slug
- **Expected**: Exit code 0, stdout displays the active feature context name
- **Priority**: P1

## TC-049: Feature set updates current feature context
- **Source**: Spec Top-Level Commands Table (feature)
- **Type**: CLI
- **Target**: cli/feature
- **Test ID**: cli/feature/set-updates-context
- **Pre-conditions**: .forge/ directory exists
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge feature my-feature`
  2. Check exit code is 0
  3. Run `forge feature` to verify current context is "my-feature"
  4. Verify .forge/state.json or equivalent config reflects "my-feature"
- **Expected**: Exit code 0, subsequent `forge feature` call shows "my-feature"
- **Priority**: P1

### Probe Command

## TC-050: Probe performs HTTP health check and returns status
- **Source**: Spec Top-Level Commands Table (probe)
- **Type**: CLI
- **Target**: cli/probe
- **Test ID**: cli/probe/probe-performs-health-check
- **Pre-conditions**: A local service is running and reachable at the configured probe URL
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge probe`
  2. Check exit code is 0 when service is healthy
  3. Verify stdout contains HTTP status code (e.g., "200 OK")
  4. Stop the service and run `forge probe` again
  5. Check exit code is 1
  6. Verify stderr contains connection error or non-200 status message
- **Expected**: Exit code 0 with "200 OK" when service healthy; exit code 1 with error when service unreachable
- **Priority**: P1

### Version Command

## TC-051: Version outputs binary version information
- **Source**: Spec Top-Level Commands Table (version)
- **Type**: CLI
- **Target**: cli/version
- **Test ID**: cli/version/version-outputs-binary-info
- **Pre-conditions**: forge binary is built and available in PATH
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge version`
  2. Check exit code is 0
  3. Verify stdout contains a semver-compliant version string (e.g., "v3.0.0")
  4. Verify the version matches the value in scripts/version.txt
- **Expected**: Exit code 0, stdout shows version string matching scripts/version.txt
- **Priority**: P1

### E2E Subcommands

## TC-052: E2E setup initializes test environment for configured profile
- **Source**: Spec Command Structure Table (e2e setup)
- **Type**: CLI
- **Target**: cli/e2e-setup
- **Test ID**: cli/e2e-setup/setup-initializes-test-environment
- **Pre-conditions**: .forge/config.yaml has a valid profile configured (e.g., web-playwright)
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge e2e setup --feature my-feature`
  2. Check exit code is 0
  3. Verify test dependencies are installed (e.g., node_modules for Playwright profiles)
  4. Verify stdout contains "profile: <profile-name>"
- **Expected**: Exit code 0, test environment initialized according to profile, profile name displayed
- **Priority**: P1

## TC-053: E2E verify checks test artifacts are valid
- **Source**: Spec Command Structure Table (e2e verify)
- **Type**: CLI
- **Target**: cli/e2e-verify
- **Test ID**: cli/e2e-verify/verify-checks-test-artifacts
- **Pre-conditions**: .forge/config.yaml has a valid profile configured; test files exist for the feature
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge e2e verify --feature my-feature`
  2. Check exit code is 0 when test artifacts are valid
  3. Verify stdout contains verification pass message
  4. Corrupt a test file and run again
  5. Check exit code is 1
  6. Verify stderr contains specific validation error message
- **Expected**: Exit code 0 with valid artifacts; exit code 1 with descriptive error when artifacts invalid
- **Priority**: P1

## TC-054: E2E compile builds test code
- **Source**: Spec Command Structure Table (e2e compile)
- **Type**: CLI
- **Target**: cli/e2e-compile
- **Test ID**: cli/e2e-compile/compile-builds-test-code
- **Pre-conditions**: .forge/config.yaml has a valid profile configured; test source files exist
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge e2e compile --feature my-feature`
  2. Check exit code is 0
  3. Verify stdout contains compilation success message or compiled output path
  4. Verify compiled test artifacts exist in the expected output directory
- **Expected**: Exit code 0, test code compiled successfully, compiled artifacts exist
- **Priority**: P1

## TC-055: E2E discover lists available test files for feature
- **Source**: Spec Command Structure Table (e2e discover)
- **Type**: CLI
- **Target**: cli/e2e-discover
- **Test ID**: cli/e2e-discover/discover-lists-test-files
- **Pre-conditions**: Feature directory exists with multiple test files matching the configured profile
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge e2e discover --feature my-feature`
  2. Check exit code is 0
  3. Verify stdout lists all discoverable test files with file paths
  4. Verify each listed file path exists on disk
- **Expected**: Exit code 0, stdout lists all test files for the feature; listed files exist on disk
- **Priority**: P1

### Task Submit --result blocked

## TC-056: Task submit with --result blocked transitions to blocked state
- **Source**: Story 3 / AC-1 (extended: blocked result variant)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/result-blocked-transitions-to-blocked
- **Pre-conditions**: Task T-block-1 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-block-1 --result blocked --summary "blocked by external dependency"`
  2. Check exit code is 0
  3. Run `cat index.json | jq '.tasks[] | select(.id=="T-block-1") | .status'` and verify result is "blocked"
  4. Verify a record file exists in records/ directory for T-block-1
- **Expected**: Exit code 0, index.json updated to status "blocked", record file created
- **Priority**: P1

### State Machine Transition Validation

## TC-057: Valid transition pending to in_progress succeeds
- **Source**: Spec State Transition Constraints Table (pending row)
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-state/valid-pending-to-in-progress
- **Pre-conditions**: Task T-trans-1 exists with status pending in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task claim` (claims T-trans-1)
  2. Check exit code is 0
  3. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-1") | .status'` and verify result is "in_progress"
- **Expected**: Exit code 0, status transitioned from pending to in_progress
- **Priority**: P0

## TC-058: Valid transition pending to rejected succeeds
- **Source**: Spec State Transition Constraints Table (pending row)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/valid-pending-to-rejected
- **Pre-conditions**: Task T-trans-2 exists with status pending in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-trans-2 --result rejected --summary "duplicate of another task"`
  2. Check exit code is 0
  3. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-2") | .status'` and verify result is "rejected"
- **Expected**: Exit code 0, status transitioned from pending to rejected
- **Priority**: P0

## TC-059: Valid transition in_progress to completed succeeds
- **Source**: Spec State Transition Constraints Table (in_progress row)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/valid-in-progress-to-completed
- **Pre-conditions**: Task T-trans-3 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-trans-3 --result success --summary "implementation complete"`
  2. Check exit code is 0
  3. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-3") | .status'` and verify result is "completed"
- **Expected**: Exit code 0, status transitioned from in_progress to completed
- **Priority**: P0

## TC-060: Valid transition in_progress to blocked succeeds
- **Source**: Spec State Transition Constraints Table (in_progress row)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/valid-in-progress-to-blocked
- **Pre-conditions**: Task T-trans-4 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-trans-4 --result blocked --summary "blocked by upstream dependency"`
  2. Check exit code is 0
  3. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-4") | .status'` and verify result is "blocked"
- **Expected**: Exit code 0, status transitioned from in_progress to blocked
- **Priority**: P0

## TC-061: Valid transition blocked to in_progress succeeds
- **Source**: Spec State Transition Constraints Table (blocked row)
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-state/valid-blocked-to-in-progress
- **Pre-conditions**: Task T-trans-5 exists with status blocked in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task claim` (claims T-trans-5)
  2. Check exit code is 0
  3. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-5") | .status'` and verify result is "in_progress"
- **Expected**: Exit code 0, status transitioned from blocked to in_progress
- **Priority**: P0

## TC-062: Invalid transition pending to completed is rejected
- **Source**: Spec State Transition Constraints Table (pending row — allowed: in_progress, rejected only)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/invalid-pending-to-completed
- **Pre-conditions**: Task T-trans-6 exists with status pending in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-trans-6 --result success --summary "skip to completed"`
  2. Check exit code is 1
  3. Verify stderr contains "invalid state transition: pending → completed"
  4. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-6") | .status'` and verify status is still "pending" (unchanged)
- **Expected**: Exit code 1, stderr contains "invalid state transition: pending → completed", index.json unchanged
- **Priority**: P0

## TC-063: Invalid transition blocked to completed is rejected
- **Source**: Spec State Transition Constraints Table (blocked row — allowed: in_progress only)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/invalid-blocked-to-completed
- **Pre-conditions**: Task T-trans-7 exists with status blocked in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-trans-7 --result success --summary "skip blocked to completed"`
  2. Check exit code is 1
  3. Verify stderr contains "invalid state transition: blocked → completed"
  4. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-7") | .status'` and verify status is still "blocked" (unchanged)
- **Expected**: Exit code 1, stderr contains "invalid state transition: blocked → completed", index.json unchanged
- **Priority**: P0

## TC-064: Invalid transition in_progress to pending is rejected
- **Source**: Spec State Transition Constraints Table (in_progress row — allowed: completed, blocked, rejected only)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/invalid-in-progress-to-pending
- **Pre-conditions**: Task T-trans-8 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Attempt to transition T-trans-8 back to pending by calling `forge task submit T-trans-8 --result pending --summary "revert"`
  2. Check exit code is 1
  3. Verify stderr contains "invalid state transition: in_progress → pending" or "invalid result value: pending"
  4. Run `cat index.json | jq '.tasks[] | select(.id=="T-trans-8") | .status'` and verify status is still "in_progress" (unchanged)
- **Expected**: Exit code 1, transition rejected, index.json unchanged
- **Priority**: P0

## TC-065: Terminal state completed rejects all transitions
- **Source**: Spec State Transition Constraints Table (completed row — terminal)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/terminal-completed-rejects-all
- **Pre-conditions**: Task T-trans-9 exists with status completed in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-trans-9 --result success --summary "re-submit completed"`
  2. Check exit code is 1
  3. Verify stderr contains "task already in terminal state: completed"
  4. Verify index.json unchanged
- **Expected**: Exit code 1, stderr contains "task already in terminal state: completed", no status change
- **Priority**: P0

## TC-066: Terminal state rejected rejects all transitions
- **Source**: Spec State Transition Constraints Table (rejected row — terminal)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-state/terminal-rejected-rejects-all
- **Pre-conditions**: Task T-trans-10 exists with status rejected in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-trans-10 --result success --summary "re-submit rejected"`
  2. Check exit code is 1
  3. Verify stderr contains "task already in terminal state: rejected"
  4. Verify index.json unchanged
- **Expected**: Exit code 1, stderr contains "task already in terminal state: rejected", no status change
- **Priority**: P0

### Boundary and Edge Cases

## TC-072: Task submit with invalid --result value returns error
- **Source**: Spec Error Handling Table (task submit, invalid result value)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/invalid-result-value-returns-error
- **Pre-conditions**: Task T-edge-1 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-edge-1 --result invalid_state --summary "bad result value"`
  2. Check exit code is 1
  3. Verify stderr contains "invalid result value: invalid_state" or "unknown result: invalid_state"
  4. Run `cat index.json | jq '.tasks[] | select(.id=="T-edge-1") | .status'` and verify status is still "in_progress" (unchanged)
- **Expected**: Exit code 1, stderr contains error about invalid result value, index.json unchanged
- **Priority**: P0

## TC-073: Task submit with empty --summary string accepts submission
- **Source**: Spec Command Structure Table (task submit, --summary flag) + Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/empty-summary-accepts-or-rejects
- **Pre-conditions**: Task T-edge-2 exists with status in_progress in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit T-edge-2 --result success --summary ""`
  2. If accepted: check exit code is 0, verify index.json status is "completed", verify record file exists with empty summary field
  3. If rejected: check exit code is 1, verify stderr contains "summary must not be empty", verify index.json unchanged
  4. Verify the behavior matches the PRD specification for empty summary handling
- **Expected**: Exit code and behavior match PRD spec for empty --summary; index.json state is consistent with the acceptance/rejection
- **Priority**: P1

## TC-074: Task submit with special characters in task ID returns error
- **Source**: Spec Error Handling Table (task submit, task not found)
- **Type**: CLI
- **Target**: cli/task-submit
- **Test ID**: cli/task-submit/special-chars-in-task-id-returns-error
- **Pre-conditions**: No task exists with a special-character ID in index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task submit "T-<script>alert(1)</script>" --result success --summary "xss test"`
  2. Check exit code is 1
  3. Verify stderr contains "task not found" error
  4. Verify index.json is unchanged and no injection artifacts exist
  5. Run `forge task submit "T space id" --result success --summary "space test"`
  6. Check exit code is 1
  7. Verify stderr contains "task not found" error
  8. Run `forge task submit "T'; DROP TABLE tasks;--" --result success --summary "sql test"`
  9. Check exit code is 1
  10. Verify stderr contains "task not found" error, no data corruption in index.json
- **Expected**: Exit code 1 for all special-character IDs; stderr contains "task not found"; index.json remains valid JSON with no injection artifacts
- **Priority**: P1

## TC-075: Task index with no task markdown files creates empty index.json
- **Source**: Spec Command Structure Table (task index)
- **Type**: CLI
- **Target**: cli/task-index
- **Test ID**: cli/task-index/empty-directory-creates-empty-index
- **Pre-conditions**: Feature directory exists but contains no task markdown files (*.md); no existing index.json
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Create a feature directory with no markdown files
  2. Run `forge task index`
  3. Check exit code is 0
  4. Verify index.json is created and is valid JSON parseable by `jq .`
  5. Run `cat index.json | jq '.tasks | length'` and verify result is 0
- **Expected**: Exit code 0, index.json created with empty tasks array, valid JSON structure
- **Priority**: P1

## TC-076: Quality gate with partial failure reports specific failing step
- **Source**: Story 4 / AC-2 (quality gate sequence) + Spec Error Handling Table
- **Type**: CLI
- **Target**: cli/quality-gate
- **Test ID**: cli/quality-gate/partial-failure-reports-failing-step
- **Pre-conditions**: Project code compiles successfully but has lint errors (compile passes, lint fails); all tasks completed
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge quality-gate`
  2. Check exit code is 1
  3. Verify stdout or stderr indicates compile step passed
  4. Verify stderr indicates lint step failed with specific lint error details
  5. Verify a P0 fix-task is created in index.json with title containing "fix-lint-1" (not "fix-compile")
  6. Verify the fix-task type is "fix" and priority is P0
- **Expected**: Exit code 1, compile step passes, lint step fails, fix-task created targeting the lint step specifically (not compile)
- **Priority**: P0

## TC-077: Task add with missing required flags returns error
- **Source**: Spec Command Structure Table (task add)
- **Type**: CLI
- **Target**: cli/task-add
- **Test ID**: cli/task-add/missing-flags-returns-error
- **Pre-conditions**: index.json exists; feature context is set
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task add` with no flags
  2. Check exit code is 1
  3. Verify stderr contains "required" and references "title"
  4. Run `forge task add --title "no-type-task"` (missing --type, optional)
  5. Check exit code is 0 (--type is optional, InferType used as fallback)
  6. Run `forge task add --type feature` (missing --title)
  7. Check exit code is 1
  8. Verify stderr contains "required" and references "title"
- **Expected**: Exit code 1 when --title is missing; exit code 0 when only --type is missing (type is optional)
- **Priority**: P0

## TC-078: Concurrent task claim on same task resolves to single winner
- **Source**: Story 3 / AC-4 (concurrent conflict) + Spec Agent Task Execution Flow (claim step)
- **Type**: CLI
- **Target**: cli/task-claim
- **Test ID**: cli/task-claim/concurrent-claim-single-winner
- **Pre-conditions**: index.json has exactly one task T-concurrent-1 with status pending
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. From two shell processes, simultaneously run `forge task claim` (both targeting T-concurrent-1)
  2. Verify exactly one process exits with code 0 and stdout contains "T-concurrent-1"
  3. Verify the other process exits with code 1 and stderr contains "concurrent write conflict, retry" or "no available tasks to claim"
  4. Run `cat index.json | jq '.tasks[] | select(.id=="T-concurrent-1") | .status'` and verify result is "in_progress"
  5. Verify index.json is valid JSON parseable by `jq .`
- **Expected**: Exactly one claim succeeds (exit code 0), the other fails (exit code 1); index.json shows exactly one task claimed to in_progress; JSON integrity preserved
- **Priority**: P1

### Task check-deps and validate-index Happy Paths

## TC-080: Task check-deps with all dependencies met succeeds
- **Source**: Spec Agent Task Execution Flow (check-deps step) + Spec Command Structure Table (task check-deps)
- **Type**: CLI
- **Target**: cli/task-check-deps
- **Test ID**: cli/task-check-deps/all-deps-met-succeeds
- **Pre-conditions**: Task T-dep-1 exists in index.json with dependencies on tasks T-dep-A and T-dep-B, both of which have status completed
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task check-deps T-dep-1`
  2. Check exit code is 0
  3. Verify stdout contains "all dependencies met" or "dependencies satisfied"
  4. Verify each dependency task ID and its completed status is listed in output
- **Expected**: Exit code 0, stdout confirms all dependencies are satisfied, each dependency listed with completed status
- **Priority**: P0

## TC-081: Task validate-index with valid schema succeeds
- **Source**: Spec Agent Task Execution Flow (validate-index step) + Spec Command Structure Table (task validate-index)
- **Type**: CLI
- **Target**: cli/task-validate-index
- **Test ID**: cli/task-validate-index/valid-schema-succeeds
- **Pre-conditions**: index.json exists with valid schema (correct fields, valid types, all required fields present)
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge task validate-index`
  2. Check exit code is 0
  3. Verify stdout contains "index validation passed" or "index is valid"
  4. Verify the number of validated tasks is reported in output
- **Expected**: Exit code 0, stdout confirms validation passed, task count reported
- **Priority**: P0

### End-to-End Workflows (Integration)

## TC-067: Agent Flow — claim to submit lifecycle completes successfully
- **Source**: Spec Agent Task Execution Flow (full lifecycle)
- **Type**: Integration
- **Target**: cli/workflow-agent-flow
- **Test ID**: cli/workflow/agent-claim-execute-submit-lifecycle
- **Pre-conditions**: Feature "wf-test-1" exists with at least one pending implementation task T-wf-1; forge binary built
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge feature wf-test-1` to set feature context
  2. Run `forge prompt get-by-task-id T-wf-1` and verify exit code 0; capture prompt output
  3. Verify prompt contains substituted TASK_ID matching "T-wf-1"
  4. Run `forge task claim` and verify exit code 0; capture claimed task ID
  5. Run `cat index.json | jq '.tasks[] | select(.id=="T-wf-1") | .status'` and verify result is "in_progress"
  6. Run `forge task submit T-wf-1 --result success --summary "completed via agent workflow"`
  7. Check exit code is 0
  8. Run `cat index.json | jq '.tasks[] | select(.id=="T-wf-1") | .status'` and verify result is "completed"
  9. Verify records/ directory contains a record file for T-wf-1
- **Expected**: Full lifecycle succeeds: prompt retrieved, task claimed (pending->in_progress), task submitted (in_progress->completed), record created. All intermediate exit codes 0.
- **Priority**: P0

## TC-068: Agent Flow — quality gate failure creates fix-task and blocked escalation
- **Source**: Spec Agent Task Execution Flow (quality gate failure branch)
- **Type**: Integration
- **Target**: cli/workflow-agent-flow-failure
- **Test ID**: cli/workflow/agent-quality-gate-failure-escalation
- **Pre-conditions**: Feature "wf-test-2" exists with a pending task T-wf-2; project has a deliberate compile error
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge feature wf-test-2`
  2. Run `forge task claim` and verify exit code 0
  3. Run `forge task submit T-wf-2 --result success --summary "implementation done"`
  4. Verify exit code is 0
  5. Run `forge quality-gate` and verify exit code is 1 (compile fails)
  6. Verify index.json contains a new P0 fix-task with title containing the failed step name
  7. Run `forge task claim` to claim the fix-task
  8. Run `forge task submit <fix-task-id> --result blocked --summary "cannot fix, escalating"`
  9. Verify exit code is 0 and fix-task status is "blocked"
- **Expected**: Quality gate failure creates fix-task; fix-task can be claimed and submitted as blocked. Full failure-to-escalation path works end-to-end.
- **Priority**: P0

## TC-069: Hook Flow — cleanup then quality-gate then verify-task-done sequence
- **Source**: Spec CI/Hook Flow (SessionEnd -> Stop -> PreToolUse)
- **Type**: Integration
- **Target**: cli/workflow-hook-sequence
- **Test ID**: cli/workflow/hook-cleanup-qualitygate-verify-sequence
- **Pre-conditions**: Feature "wf-test-3" exists with some completed tasks and some in_progress tasks; project code compiles and passes tests
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge cleanup` and verify exit code is 0
  2. Verify state files for completed tasks are removed
  3. Run `forge quality-gate` and verify exit code is 0 (all pass)
  4. Run `forge verify-task-done` and verify exit code is 1 (in_progress tasks remain)
  5. Submit all remaining in_progress tasks: `forge task submit <id> --result success --summary "done"`
  6. Run `forge verify-task-done` again and verify exit code is 0 (all tasks terminal)
- **Expected**: Cleanup removes terminal state files. Quality gate passes when code is good. verify-task-done fails until all tasks are terminal, then passes. Full hook sequence behaves correctly.
- **Priority**: P0

## TC-070: Developer Flow — profile detect to set to e2e run
- **Source**: Spec Developer Flow (profile setup + e2e execution)
- **Type**: Integration
- **Target**: cli/workflow-developer-flow
- **Test ID**: cli/workflow/developer-profile-detect-set-e2e-run
- **Pre-conditions**: Project contains playwright.config.ts; no profile set in .forge/config.yaml
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge profile detect` and verify exit code is 0
  2. Verify stdout contains "web-playwright" with detection evidence
  3. Run `forge profile set web-playwright` and verify exit code is 0
  4. Verify .forge/config.yaml profile field is "web-playwright"
  5. Run `forge e2e discover --feature my-feature` and verify exit code is 0
  6. Run `forge e2e run --feature my-feature` and verify profile is used correctly
- **Expected**: Profile detected from project, set in config, and used by e2e run. Full developer setup-to-execution flow works.
- **Priority**: P0

## TC-071: Quality gate retry loop — fix-task creation up to max 3
- **Source**: Spec Agent Task Execution Flow (retry branch) + Story 4 / AC-4, AC-5
- **Type**: Integration
- **Target**: cli/workflow-quality-gate-retry
- **Test ID**: cli/workflow/quality-gate-retry-loop-max-3
- **Pre-conditions**: Feature "wf-test-4" exists; project has a persistent compile error
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge quality-gate` — exit code 1, verify first fix-task (e.g., "fix-compile-1") created in index.json
  2. Run `forge quality-gate` again — exit code 1, verify second fix-task "fix-compile-2" created
  3. Run `forge quality-gate` again — exit code 1, verify third fix-task "fix-compile-3" created
  4. Run `forge quality-gate` again — exit code 1, verify stderr contains "max fix-tasks reached for compile, manual intervention required"
  5. Verify index.json still has exactly 3 fix-tasks for compile step
- **Expected**: Quality gate creates exactly 3 fix-tasks for the same failing step, then stops and outputs max-fix-tasks message. No fourth fix-task created.
- **Priority**: P0

## TC-082: Agent Flow — quality gate failure, fix succeeds, retry passes, task submits
- **Source**: Spec Agent Task Execution Flow (recovery branch: quality gate fail -> fix -> retry -> succeed)
- **Type**: Integration
- **Target**: cli/workflow-agent-flow-recovery
- **Test ID**: cli/workflow/agent-quality-gate-recovery-success
- **Pre-conditions**: Feature "wf-test-5" exists with a pending implementation task T-wf-5; project has a fixable compile error (e.g., missing import that can be added); forge binary built
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `forge feature wf-test-5` to set feature context
  2. Run `forge task claim` and verify exit code 0; capture claimed task ID (T-wf-5)
  3. Run `forge task submit T-wf-5 --result success --summary "implementation done"`
  4. Verify exit code is 0 and T-wf-5 status is "completed" in index.json
  5. Run `forge quality-gate` — verify exit code is 1 (compile fails due to fixable error)
  6. Verify index.json contains a new P0 fix-task with title containing "fix-compile-1"
  7. Run `forge task claim` to claim the fix-task; verify exit code 0
  8. Fix the compile error in the project source code (e.g., add the missing import)
  9. Run `forge task submit <fix-task-id> --result success --summary "import added, compile fixed"`
  10. Verify exit code is 0 and fix-task status is "completed" in index.json
  11. Run `forge quality-gate` again — verify exit code is 0 (all steps pass now)
  12. Verify no new fix-tasks are created in index.json
  13. Run `forge verify-task-done` — verify exit code is 0 (all tasks terminal)
- **Expected**: Full recovery loop succeeds: task submitted -> quality gate fails -> fix-task created -> fix executed -> fix submitted -> quality gate retries and passes -> all tasks verified done. No residual fix-tasks remain in non-terminal state.
- **Priority**: P0

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
| TC-042 | Spec Agent Task Execution Flow | CLI | cli/task-claim | P0 |
| TC-043 | Spec Command Structure (task add) | CLI | cli/task-add | P1 |
| TC-044 | Spec Command Structure (task index) | CLI | cli/task-index | P1 |
| TC-045 | Spec Command Structure (task migrate) | CLI | cli/task-migrate | P2 |
| TC-046 | Spec Command Structure (task query) | CLI | cli/task-query | P1 |
| TC-047 | Spec Command Structure (task status) | CLI | cli/task-status | P0 |
| TC-048 | Spec Top-Level Commands (feature get) | CLI | cli/feature | P1 |
| TC-049 | Spec Top-Level Commands (feature set) | CLI | cli/feature | P1 |
| TC-050 | Spec Top-Level Commands (probe) | CLI | cli/probe | P1 |
| TC-051 | Spec Top-Level Commands (version) | CLI | cli/version | P1 |
| TC-052 | Spec Command Structure (e2e setup) | CLI | cli/e2e-setup | P1 |
| TC-053 | Spec Command Structure (e2e verify) | CLI | cli/e2e-verify | P1 |
| TC-054 | Spec Command Structure (e2e compile) | CLI | cli/e2e-compile | P1 |
| TC-055 | Spec Command Structure (e2e discover) | CLI | cli/e2e-discover | P1 |
| TC-056 | Story 3 / AC-1 (blocked variant) | CLI | cli/task-submit | P1 |
| TC-057 | Spec State Transition (pending row) | CLI | cli/task-state | P0 |
| TC-058 | Spec State Transition (pending row) | CLI | cli/task-state | P0 |
| TC-059 | Spec State Transition (in_progress row) | CLI | cli/task-state | P0 |
| TC-060 | Spec State Transition (in_progress row) | CLI | cli/task-state | P0 |
| TC-061 | Spec State Transition (blocked row) | CLI | cli/task-state | P0 |
| TC-062 | Spec State Transition (pending invalid) | CLI | cli/task-state | P0 |
| TC-063 | Spec State Transition (blocked invalid) | CLI | cli/task-state | P0 |
| TC-064 | Spec State Transition (in_progress invalid) | CLI | cli/task-state | P0 |
| TC-065 | Spec State Transition (completed terminal) | CLI | cli/task-state | P0 |
| TC-066 | Spec State Transition (rejected terminal) | CLI | cli/task-state | P0 |
| TC-067 | Spec Agent Task Execution Flow | Integration | cli/workflow-agent-flow | P0 |
| TC-068 | Spec Agent Flow (failure branch) | Integration | cli/workflow-agent-flow-failure | P0 |
| TC-069 | Spec CI/Hook Flow | Integration | cli/workflow-hook-sequence | P0 |
| TC-070 | Spec Developer Flow | Integration | cli/workflow-developer-flow | P0 |
| TC-071 | Spec Agent Flow (retry) + Story 4 / AC-4, AC-5 | Integration | cli/workflow-quality-gate-retry | P0 |
| TC-072 | Spec Error Handling (invalid result value) | CLI | cli/task-submit | P0 |
| TC-073 | Spec Command Structure (task submit, --summary) | CLI | cli/task-submit | P1 |
| TC-074 | Spec Error Handling (task not found) | CLI | cli/task-submit | P1 |
| TC-075 | Spec Command Structure (task index) | CLI | cli/task-index | P1 |
| TC-076 | Story 4 / AC-2 (partial failure) | CLI | cli/quality-gate | P0 |
| TC-077 | Spec Command Structure (task add, required flags) | CLI | cli/task-add | P0 |
| TC-078 | Story 3 / AC-4 + Spec Agent Flow (claim) | CLI | cli/task-claim | P1 |
| TC-080 | Spec Agent Flow (check-deps) + Spec Command Structure | CLI | cli/task-check-deps | P0 |
| TC-081 | Spec Agent Flow (validate-index) + Spec Command Structure | CLI | cli/task-validate-index | P0 |
| TC-082 | Spec Agent Task Execution Flow (recovery branch) | Integration | cli/workflow-agent-flow-recovery | P0 |
