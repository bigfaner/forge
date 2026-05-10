---
feature: "task-executor-skeleton"
sources:
  - docs/features/task-executor-skeleton/prd/prd-user-stories.md
  - docs/features/task-executor-skeleton/prd/prd-spec.md
generated: "2026-05-10"
---

# Test Cases: task-executor-skeleton

> **WARNING**: sitemap.json not found — Element set to `sitemap-missing`. Run `/gen-sitemap` for precise element references.

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 16  |
| **Total** | **16** |

> **Note**: This feature has no UI surface. All changes are internal to the forge harness (agent prompts, task templates, CLI code, skill docs). Only CLI test cases are generated.

---

## UI Test Cases

_No UI test cases — feature has no UI surface._

---

## API Test Cases

_No API test cases — feature exposes no HTTP endpoints._

---

## CLI Test Cases

### Workflow Detection & Injection

## TC-001: Execution Workflow detected in task template replaces TDD
- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/execution-workflow-detected-replaces-tdd
- **Pre-conditions**: Task template file exists with valid frontmatter and contains `## Execution Workflow` section with non-empty body
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Prepare a task template file containing a `## Execution Workflow` section with non-empty body content (e.g., "Run the build and report results")
  2. Dispatch the task to task-executor
  3. Inspect the agent prompt injected for Step 2
- **Expected**: Step 2 instructions contain the body of `## Execution Workflow` from the template, not the hardcoded TDD steps (RED/GREEN/REFACTOR)
- **Priority**: P0

## TC-002: Missing Execution Workflow falls back to TDD and Quality Gate
- **Source**: Story 1 / AC-2
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/missing-execution-workflow-fallback-tdd
- **Pre-conditions**: Task template file exists with valid frontmatter but does NOT contain a `## Execution Workflow` section
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Prepare a task template file without a `## Execution Workflow` section
  2. Dispatch the task to task-executor
  3. Inspect the agent prompt and execution behavior
- **Expected**: Task-executor falls back to the TDD + Quality Gate flow (Step 2 = TDD implementation, Step 3 = quality gate verification). Behavior is identical to the pre-feature baseline.
- **Priority**: P0

## TC-003: Empty Execution Workflow body triggers warning and TDD fallback
- **Source**: Story 1 / AC-3
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/empty-execution-workflow-warning-tdd-fallback
- **Pre-conditions**: Task template file exists with valid frontmatter and contains a `## Execution Workflow` header but the body is empty (no content after the header)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Prepare a task template file with `## Execution Workflow` header followed by no body content
  2. Dispatch the task to task-executor
  3. Inspect execution logs and agent prompt
- **Expected**: A configuration error warning is recorded in the execution log. Task-executor falls back to TDD + Quality Gate flow. The task is not blocked by the empty workflow.
- **Priority**: P1

### Execution Workflow Behavior

## TC-004: Execution-type task creates fix task on failure without TDD retry
- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/execution-task-creates-fix-task-on-failure
- **Pre-conditions**: Task uses a template (e.g., T-test-3) with an Execution Workflow containing explicit failure instructions (e.g., "create fix task, do not retry"). A controlled failure condition is set up (e.g., failing e2e test).
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Dispatch a task with an Execution Workflow that includes failure handling instructions
  2. Trigger a failure condition during execution
  3. Observe agent behavior after the failure
- **Expected**: Agent creates a fix task as specified by the workflow and stops. Agent does NOT enter a TDD loop or retry the failed step. Task execution time is significantly reduced compared to the pre-feature TDD-retry behavior.
- **Priority**: P0

## TC-005: Step 2 output uses Execution Workflow terminology not TDD terminology
- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/step2-output-uses-workflow-terminology
- **Pre-conditions**: Task with a valid Execution Workflow has been executed and completed (success or failure)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Execute a task that uses an Execution Workflow (non-TDD)
  2. Read the execution record for Step 2
  3. Check for keyword presence
- **Expected**: Step 2 output contains the phrase "Execution Workflow" (or "Implementation" for workflow-based tasks). Step 2 output does NOT contain "TDD implementation", "RED/GREEN/REFACTOR", or "TDD cycle" keywords.
- **Priority**: P1

## TC-006: Execution-type task skips Quality Gate and proceeds to record and commit
- **Source**: Story 2 / AC-3
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/execution-task-skips-quality-gate
- **Pre-conditions**: Task uses a template with an Execution Workflow. The workflow completes successfully.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Dispatch a task with an Execution Workflow that completes successfully
  2. Observe the execution flow after Step 2 completion
- **Expected**: After workflow execution completes, task-executor proceeds directly to Step 3 (record + commit). The Quality Gate sequence (compile -> fmt -> lint -> test) is NOT executed. Step 3 output reflects "SKIPPED" for the quality gate.
- **Priority**: P0

### noTest Removal Verification

## TC-007: Grep noTest and NO_TEST across all harness files yields zero matches
- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/task-cli
- **Test ID**: cli/task-cli/grep-notest-zero-matches
- **Pre-conditions**: All code changes for the feature have been applied. The harness codebase is in its final state.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Run `grep -ri "notest" --include="*.md" --include="*.go" --include="*.json"` across the entire task-cli and forge harness directory
  2. Run `grep -ri "no_test" --include="*.md" --include="*.go" --include="*.json"` across the same directories
  3. Run `grep -ri "NO_TEST" --include="*.md" --include="*.go" --include="*.json"` across the same directories
- **Expected**: All three grep commands return zero matches. No file in task-cli/agent/command/skill directories contains any variant of noTest/NO_TEST/no_test.
- **Priority**: P0

## TC-008: task-cli Go code has no noTest conditional branches
- **Source**: Story 3 / AC-2
- **Type**: CLI
- **Target**: cli/task-cli
- **Test ID**: cli/task-cli/golang-no-notest-branches
- **Pre-conditions**: All code changes for the feature have been applied.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Inspect `types.go` for any field or struct member named `noTest`, `NoTest`, or `NO_TEST`
  2. Inspect `record.go` for any conditional branch based on a noTest field
  3. Inspect all Go files for any logic that branches on the presence or absence of a noTest flag
- **Expected**: No noTest-related fields exist in types.go. No conditional branches in record.go (or any Go file) reference noTest. The `Task` struct and related types are clean of noTest.
- **Priority**: P0

## TC-009: All 16 task templates have no noTest in frontmatter
- **Source**: Story 3 / AC-3
- **Type**: CLI
- **Target**: cli/task-templates
- **Test ID**: cli/task-templates/all-templates-no-notest-frontmatter
- **Pre-conditions**: All task templates have been updated for the feature.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Enumerate all breakdown task templates (10 files, excluding manifest-update-tasks.md and eval-test-cases.md)
  2. Enumerate all quick-task templates (6 files, excluding manifest-quick.md)
  3. Parse each template's YAML frontmatter
  4. Check for the presence of a `noTest` key
- **Expected**: All 16 templates have valid frontmatter with zero occurrences of the `noTest` key. Template count matches expected (10 breakdown + 6 quick = 16).
- **Priority**: P0

## TC-010: index.schema.json files have no noTest field definition and validate all templates
- **Source**: Story 3 / AC-4
- **Type**: CLI
- **Target**: cli/task-cli
- **Test ID**: cli/task-cli/schema-no-notest-validates-templates
- **Pre-conditions**: Schema files and templates are in their final state. `ajv` validator is available.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `index.schema.json` for breakdown templates
  2. Read `index.schema.json` for quick templates
  3. Verify no `noTest` property definition exists in either schema
  4. Run `ajv validate` against all templates using the corresponding schema
- **Expected**: Neither schema file contains a `noTest` property definition. `ajv validate` passes for all templates against their respective schemas (zero validation errors).
- **Priority**: P1

## TC-011: Command docs run-tasks.md and execute-task.md have no NO_TEST references
- **Source**: Story 3 / AC-5
- **Type**: CLI
- **Target**: cli/commands
- **Test ID**: cli/commands/run-tasks-execute-task-no-notest
- **Pre-conditions**: Command docs are in their final state.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `commands/run-tasks.md`
  2. Read `commands/execute-task.md`
  3. Grep both files for `NO_TEST` (case-insensitive)
  4. Verify claim parsing section has no NO_TEST extraction logic
  5. Verify dispatch prompt section has no NO_TEST parameter passing
- **Expected**: Zero matches for NO_TEST in either file. Claim parsing logic does not extract a NO_TEST variable. Dispatch prompt does not include NO_TEST as a parameter to subagents.
- **Priority**: P1

## TC-012: task-executor.md Step 2-3 has no NO_TEST references and uses workflow injection
- **Source**: Story 3 / AC-6
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/step2-3-no-notest-workflow-injection
- **Pre-conditions**: task-executor.md is in its final state.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Read `agents/task-executor.md`
  2. Grep for `NO_TEST` in the file (case-insensitive)
  3. Inspect Step 2 definition for workflow reading and injection logic
  4. Inspect Step 3 definition for absence of noTest-based conditional logic
- **Expected**: Zero matches for NO_TEST in task-executor.md. Step 2 describes reading `## Execution Workflow` from the task file and injecting it into the agent prompt. Step 3 does not reference noTest for quality gate bypass decisions.
- **Priority**: P0

### Failure Handling

## TC-013: Missing or unparseable task file sets status to failed with error log
- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/missing-task-file-status-failed
- **Pre-conditions**: A task is dispatched but its corresponding task file is missing or has corrupted/unparseable frontmatter.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Dispatch a task referencing a non-existent task file (or a file with invalid YAML frontmatter)
  2. Observe task-executor behavior
  3. Check the task status and execution log
- **Expected**: Task status is set to `failed`. An error log entry describes the file access or parsing failure. Step 2 (execution) is NOT entered. The agent records the failure and proceeds to commit the error record.
- **Priority**: P0

## TC-014: Workflow failure with explicit failure instruction followed correctly
- **Source**: Story 4 / AC-2
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/workflow-failure-explicit-instruction-followed
- **Pre-conditions**: Task template has an Execution Workflow that includes an explicit failure instruction (e.g., "on failure, create a fix task"). A failure condition occurs during workflow execution.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Dispatch a task with a workflow containing explicit failure handling instructions
  2. Trigger a failure during workflow execution
  3. Observe agent response to the failure
- **Expected**: Agent executes the explicit failure instruction (e.g., creates a fix task). Agent then stops execution. Task status is set to `failed`. Agent does NOT retry the failed step.
- **Priority**: P0

## TC-015: Workflow failure without explicit instruction records and stops
- **Source**: Story 4 / AC-3
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/workflow-failure-no-instruction-records-and-stops
- **Pre-conditions**: Task template has an Execution Workflow without any explicit failure handling instructions. A failure condition occurs during workflow execution.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Dispatch a task with a workflow that has no failure handling instructions
  2. Trigger a failure during workflow execution
  3. Observe agent response to the failure
- **Expected**: Agent records the failure reason in the execution log. Agent stops execution. Task status is set to `failed`. Agent does NOT enter a TDD loop or retry the failed step.
- **Priority**: P1

## TC-016: Multi-step workflow mid-failure records completed steps and failure point
- **Source**: Story 4 / AC-4
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/multi-step-mid-failure-records-progress
- **Pre-conditions**: Task template has a multi-step Execution Workflow (3+ steps). A controlled failure occurs at step 2.
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Dispatch a task with a multi-step workflow
  2. Allow step 1 to complete successfully
  3. Trigger a failure at step 2
  4. Inspect the failure record
- **Expected**: The failure record includes a summary of completed steps (step 1 result). The failure record describes the failure point (step 2 error details). Task status is `failed` with the partial progress documented.
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/task-executor | P0 |
| TC-002 | Story 1 / AC-2 | CLI | cli/task-executor | P0 |
| TC-003 | Story 1 / AC-3 | CLI | cli/task-executor | P1 |
| TC-004 | Story 2 / AC-1 | CLI | cli/task-executor | P0 |
| TC-005 | Story 2 / AC-2 | CLI | cli/task-executor | P1 |
| TC-006 | Story 2 / AC-3 | CLI | cli/task-executor | P0 |
| TC-007 | Story 3 / AC-1 | CLI | cli/task-cli | P0 |
| TC-008 | Story 3 / AC-2 | CLI | cli/task-cli | P0 |
| TC-009 | Story 3 / AC-3 | CLI | cli/task-templates | P0 |
| TC-010 | Story 3 / AC-4 | CLI | cli/task-cli | P1 |
| TC-011 | Story 3 / AC-5 | CLI | cli/commands | P1 |
| TC-012 | Story 3 / AC-6 | CLI | cli/task-executor | P0 |
| TC-013 | Story 4 / AC-1 | CLI | cli/task-executor | P0 |
| TC-014 | Story 4 / AC-2 | CLI | cli/task-executor | P0 |
| TC-015 | Story 4 / AC-3 | CLI | cli/task-executor | P1 |
| TC-016 | Story 4 / AC-4 | CLI | cli/task-executor | P1 |

---

## Route Validation

_No route validation — this feature has no UI/API routes. All test cases verify CLI/agent behavior through file inspection and command execution._
