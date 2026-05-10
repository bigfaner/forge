---
feature: "task-executor-skeleton"
sources:
  - prd/prd-user-stories.md
  - prd/prd-spec.md
generated: "2026-05-10"
---

# Test Cases: task-executor-skeleton

## Summary

| Type | Count |
|------|-------|
| UI   | 0   |
| **Integration** | **0** |
| API  | 0  |
| CLI  | 17  |
| **Total** | **17** |

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
- **Element**: N/A
- **Steps**:
  1. Prepare a task template file containing a `## Execution Workflow` section with non-empty body content (e.g., "Run the build and report results")
  2. Dispatch the task: `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  3. After execution, run `cat .forge/tasks/{task-id}/record.json | jq '.steps[1].output'` to read the Step 2 output
- **Expected**: Step 2 instructions contain the body of `## Execution Workflow` from the template, not the hardcoded TDD steps (RED/GREEN/REFACTOR)
- **Priority**: P0

## TC-002: Missing Execution Workflow falls back to TDD and Quality Gate
- **Source**: Story 1 / AC-2
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/missing-execution-workflow-fallback-tdd
- **Pre-conditions**: Task template file exists with valid frontmatter but does NOT contain a `## Execution Workflow` section
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Prepare a task template file without a `## Execution Workflow` section
  2. Dispatch the task: `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  3. After execution, run `cat .forge/tasks/{task-id}/record.json | jq '.steps[1].output'` to inspect Step 2 output
  4. Run `cat .forge/tasks/{task-id}/record.json | jq '.steps[2].output'` to inspect Step 3 output
- **Expected**: Task-executor falls back to the TDD + Quality Gate flow (Step 2 = TDD implementation, Step 3 = quality gate verification). Behavior is identical to the pre-feature baseline.
- **Priority**: P0

## TC-003: Empty Execution Workflow body triggers warning and TDD fallback
- **Source**: Story 1 / AC-3
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/empty-execution-workflow-warning-tdd-fallback
- **Pre-conditions**: Task template file exists with valid frontmatter and contains a `## Execution Workflow` header but the body is empty (no content after the header)
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Prepare a task template file with `## Execution Workflow` header followed by no body content
  2. Dispatch the task: `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  3. After execution, run `grep -i "configuration error\|config error\|empty workflow" .forge/tasks/{task-id}/record.json` to check for the warning
  4. Run `cat .forge/tasks/{task-id}/record.json | jq '.steps[1].output'` to verify Step 2 contains TDD content
- **Expected**: A configuration error warning is recorded in the execution log. Task-executor falls back to TDD + Quality Gate flow. The task is not blocked by the empty workflow.
- **Priority**: P1

### Execution Workflow Behavior

## TC-004: Execution-type task creates fix task on failure without TDD retry
- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/execution-task-creates-fix-task-on-failure
- **Pre-conditions**: Task uses template T-test-3 whose Execution Workflow contains "on failure, create a fix task, do not retry". A test file at `tests/e2e/smoke.test.ts` contains `assert(false)` to force a deterministic failure.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Write `assert(false)` into `tests/e2e/smoke.test.ts` to create a guaranteed failing test
  2. Dispatch a task using template T-test-3: `task-cli execute-task --template T-test-3 --task-file .forge/tasks/{task-id}.md`
  3. After execution completes, read `.forge/tasks/{task-id}/record.json` and inspect the `steps` array
  4. Run `ls .forge/tasks/` and check for a newly created fix-task file
- **Expected**: `record.json` shows exactly one attempt at the workflow step -- no retry entries. A fix-task file exists in `.forge/tasks/`. The `steps` array does NOT contain any entries with "RED", "GREEN", "REFACTOR", or "TDD cycle" in the output text. Task status is `failed`, not `in-progress`.
- **Priority**: P0

## TC-005: Step 2 output uses Execution Workflow terminology not TDD terminology
- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/step2-output-uses-workflow-terminology
- **Pre-conditions**: Task with a valid Execution Workflow has been executed and completed (success or failure), producing a record file at `.forge/tasks/{task-id}/record.json`
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Execute a task that uses an Execution Workflow (non-TDD) via `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  2. Read the execution record: `cat .forge/tasks/{task-id}/record.json | jq '.steps[1].output'` (Step 2 output)
  3. Run `grep -c "TDD implementation\|RED/GREEN/REFACTOR\|TDD cycle" .forge/tasks/{task-id}/record.json` to count forbidden keyword matches
- **Expected**: `jq` output for `steps[1].output` contains the phrase "Execution Workflow" (or "Implementation" for workflow-based tasks). `grep -c` returns `0` -- the record file contains zero matches for "TDD implementation", "RED/GREEN/REFACTOR", or "TDD cycle".
- **Priority**: P1

## TC-006: Execution-type task skips Quality Gate and proceeds to record and commit
- **Source**: Story 2 / AC-3
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/execution-task-skips-quality-gate
- **Pre-conditions**: Task uses a template with an Execution Workflow. The workflow completes successfully.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Dispatch a task with an Execution Workflow that completes successfully: `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  2. After execution, run `cat .forge/tasks/{task-id}/record.json | jq '.steps[2].output'` to inspect Step 3 (record + commit) output
  3. Run `grep -c "compile\|fmt.*check\|lint\|test" .forge/tasks/{task-id}/record.json` to verify Quality Gate was skipped
- **Expected**: `steps[2].output` shows the task proceeding directly to record + commit after Step 2. The `grep -c` for "compile", "fmt", "lint", "test" in the record returns `0` or the matches are limited to the Execution Workflow body content, not Quality Gate execution logs. The Quality Gate sequence (compile -> fmt -> lint -> test) is NOT present in any step output.
- **Priority**: P0

### noTest Removal Verification

## TC-007: Grep noTest and NO_TEST across all harness files yields zero matches
- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/task-cli
- **Test ID**: cli/task-cli/grep-notest-zero-matches
- **Pre-conditions**: All code changes for the feature have been applied. The harness codebase is in its final state.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `grep -ri "notest" --include="*.md" --include="*.go" --include="*.json" .forge/ cmd/ internal/ pkg/` across the forge harness and task-cli directories
  2. Run `grep -ri "no_test" --include="*.md" --include="*.go" --include="*.json" .forge/ cmd/ internal/ pkg/` across the same directories
  3. Run `grep -ri "NO_TEST" --include="*.md" --include="*.go" --include="*.json" .forge/ cmd/ internal/ pkg/` across the same directories
- **Expected**: All three grep commands return zero matches. No file in task-cli/agent/command/skill directories contains any variant of noTest/NO_TEST/no_test.
- **Priority**: P0

## TC-008: task-cli Go code has no noTest conditional branches
- **Source**: Story 3 / AC-2
- **Type**: CLI
- **Target**: cli/task-cli
- **Test ID**: cli/task-cli/golang-no-notest-branches
- **Pre-conditions**: All code changes for the feature have been applied.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `grep -n "noTest\|NoTest\|NO_TEST" pkg/task/types.go` to check for noTest fields in the Task struct
  2. Run `grep -n "noTest\|NoTest\|NO_TEST" internal/record/record.go` to check for noTest conditional branches in record logic
  3. Run `grep -rn "noTest\|NoTest\|NO_TEST" --include="*.go" .` to scan all Go files for any noTest-based logic
- **Expected**: No noTest-related fields exist in types.go. No conditional branches in record.go (or any Go file) reference noTest. The `Task` struct and related types are clean of noTest.
- **Priority**: P0

## TC-009: All 16 task templates have no noTest in frontmatter
- **Source**: Story 3 / AC-3
- **Type**: CLI
- **Target**: cli/task-templates
- **Test ID**: cli/task-templates/all-templates-no-notest-frontmatter
- **Pre-conditions**: All task templates have been updated for the feature.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `find .forge/templates/breakdown-tasks/ -name "*.md" ! -name "manifest-*.md" ! -name "eval-test-cases.md" | wc -l` to confirm template count is 10
  2. Run `find .forge/templates/quick-tasks/ -name "*.md" ! -name "manifest-quick.md" | wc -l` to confirm template count is 6
  3. Run `grep -rl "noTest" .forge/templates/breakdown-tasks/ .forge/templates/quick-tasks/ --include="*.md"` to find any template containing `noTest` in its YAML frontmatter
  4. For each template listed (if any), run `head -20 {file}` to inspect the frontmatter and confirm the `noTest` key exists
- **Expected**: All 16 templates have valid frontmatter with zero occurrences of the `noTest` key. Template count matches expected (10 breakdown + 6 quick = 16).
- **Priority**: P0

## TC-010: index.schema.json files have no noTest field definition and validate all templates
- **Source**: Story 3 / AC-4
- **Type**: CLI
- **Target**: cli/task-cli
- **Test ID**: cli/task-cli/schema-no-notest-validates-templates
- **Pre-conditions**: Schema files and templates are in their final state. `ajv` validator is available.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `cat .forge/templates/breakdown-tasks/index.schema.json | jq '.properties | has("noTest")'` to check for noTest property definition
  2. Run `cat .forge/templates/quick-tasks/index.schema.json | jq '.properties | has("noTest")'` to check for noTest property definition
  3. Run `ajv validate -s .forge/templates/breakdown-tasks/index.schema.json -d ".forge/templates/breakdown-tasks/*.md"` to validate all breakdown templates
  4. Run `ajv validate -s .forge/templates/quick-tasks/index.schema.json -d ".forge/templates/quick-tasks/*.md"` to validate all quick templates
- **Expected**: Neither schema file contains a `noTest` property definition. `ajv validate` passes for all templates against their respective schemas (zero validation errors).
- **Priority**: P1

## TC-011: Command docs run-tasks.md and execute-task.md have no NO_TEST references
- **Source**: Story 3 / AC-5
- **Type**: CLI
- **Target**: cli/commands
- **Test ID**: cli/commands/run-tasks-execute-task-no-notest
- **Pre-conditions**: Command docs are in their final state.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `grep -i "NO_TEST" agents/commands/run-tasks.md` to search for NO_TEST references in run-tasks
  2. Run `grep -i "NO_TEST" agents/commands/execute-task.md` to search for NO_TEST references in execute-task
  3. Run `grep -A5 "claim" agents/commands/run-tasks.md | grep -i "no_test"` to verify claim parsing section has no NO_TEST extraction logic
  4. Run `grep -A5 "dispatch" agents/commands/run-tasks.md | grep -i "no_test"` to verify dispatch prompt section has no NO_TEST parameter passing
- **Expected**: Zero matches for NO_TEST in either file. Claim parsing logic does not extract a NO_TEST variable. Dispatch prompt does not include NO_TEST as a parameter to subagents.
- **Priority**: P1

## TC-012: task-executor.md Step 2-3 has no NO_TEST references and uses workflow injection
- **Source**: Story 3 / AC-6
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/step2-3-no-notest-workflow-injection
- **Pre-conditions**: task-executor.md is in its final state.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `cat agents/task-executor.md | grep -i "NO_TEST"` to search for NO_TEST references
  2. Run `grep -A10 "Step 2" agents/task-executor.md | grep -i "execution workflow"` to verify Step 2 describes reading `## Execution Workflow`
  3. Run `grep -A10 "Step 3" agents/task-executor.md | grep -i "noTest\|NO_TEST"` to verify Step 3 has no noTest-based conditional logic
  4. Run `grep -c "## Execution Workflow" agents/task-executor.md` to confirm workflow reading logic exists in the file
- **Expected**: Zero matches for NO_TEST in task-executor.md. Step 2 describes reading `## Execution Workflow` from the task file and injecting it into the agent prompt. Step 3 does not reference noTest for quality gate bypass decisions.
- **Priority**: P0

### Failure Handling

## TC-013: Missing or unparseable task file sets status to failed with error log
- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/missing-task-file-status-failed
- **Pre-conditions**: A task is dispatched with reference to a non-existent file path (e.g., `.forge/tasks/nonexistent-task.md`).
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `task-cli execute-task --task-file .forge/tasks/nonexistent-task.md` to dispatch a task referencing a missing file
  2. After execution, run `cat .forge/tasks/nonexistent-task/record.json | jq '.status'` to check the task status
  3. Run `cat .forge/tasks/nonexistent-task/record.json | jq '.steps | length'` to verify no Step 2 execution occurred
  4. Run `cat .forge/tasks/nonexistent-task/record.json | jq '.error'` to inspect the error log entry
- **Expected**: Task status is set to `failed`. An error log entry describes the file access or parsing failure. Step 2 (execution) is NOT entered. The agent records the failure and proceeds to commit the error record.
- **Priority**: P0

## TC-014: Workflow failure with explicit failure instruction followed correctly
- **Source**: Story 4 / AC-2
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/workflow-failure-explicit-instruction-followed
- **Pre-conditions**: Task template has an Execution Workflow containing "on failure, create a fix task". A test file at `tests/e2e/smoke.test.ts` contains `assert(false)` to guarantee failure during workflow execution.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Write `assert(false)` into `tests/e2e/smoke.test.ts` to create a deterministic failure
  2. Dispatch the task: `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  3. After execution, run `ls .forge/tasks/ | grep fix` to verify a fix task was created
  4. Run `cat .forge/tasks/{task-id}/record.json | jq '.status'` to check task status
  5. Run `cat .forge/tasks/{task-id}/record.json | jq '.steps[] | select(.output | test("RED|GREEN|REFACTOR"))'` to verify no TDD keywords appear
- **Expected**: Agent executes the explicit failure instruction (e.g., creates a fix task). Agent then stops execution. Task status is set to `failed`. Agent does NOT retry the failed step.
- **Priority**: P0

## TC-015: Workflow failure without explicit instruction records and stops
- **Source**: Story 4 / AC-3
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/workflow-failure-no-instruction-records-and-stops
- **Pre-conditions**: Task template has an Execution Workflow without any explicit failure handling instructions (no "on failure" clause). The workflow references a command that will fail, e.g., `npx playwright test --config=nonexistent.config.ts` which guarantees a command-not-found error.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Dispatch the task: `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  2. After execution, run `cat .forge/tasks/{task-id}/record.json | jq '.status'` to verify status is `failed`
  3. Run `cat .forge/tasks/{task-id}/record.json | jq '.steps[-1].output'` to read the failure reason logged by the agent
  4. Run `grep -c "retry\|RED\|GREEN\|REFACTOR\|TDD cycle" .forge/tasks/{task-id}/record.json` to confirm no TDD retry loop occurred
- **Expected**: Agent records the failure reason in the execution log. Agent stops execution. Task status is set to `failed`. Agent does NOT enter a TDD loop or retry the failed step.
- **Priority**: P1

## TC-016: Multi-step workflow mid-failure records completed steps and failure point
- **Source**: Story 4 / AC-4
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/multi-step-mid-failure-records-progress
- **Pre-conditions**: Task template has a multi-step Execution Workflow with 3 steps: step 1 runs `echo "step 1 done"`, step 2 runs `exit 1` (guaranteed failure), step 3 runs `echo "step 3 done"`. The workflow body contains these three steps in sequence.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Dispatch the task: `task-cli execute-task --task-file .forge/tasks/{task-id}.md`
  2. After execution, run `cat .forge/tasks/{task-id}/record.json | jq '.status'` to verify status is `failed`
  3. Run `cat .forge/tasks/{task-id}/record.json | jq '.completed_steps'` to inspect the completed steps summary
  4. Run `cat .forge/tasks/{task-id}/record.json | jq '.failure_point'` to inspect the failure point description
- **Expected**: `completed_steps` contains at least one entry with output matching "step 1 done". `failure_point` contains a string referencing step 2 with error details including the `exit 1` code. `steps` array length is exactly 2 (step 1 completed + step 2 failed) -- step 3 was not executed. Task status is `failed`.
- **Priority**: P1

### End-to-End Integration

## TC-017: Full dispatch-to-commit pipeline with Execution Workflow template
- **Source**: Story 1 / AC-1, Story 2 / AC-3
- **Type**: CLI
- **Target**: cli/task-executor
- **Test ID**: cli/task-executor/e2e-dispatch-to-commit-workflow
- **Pre-conditions**: A task template exists at `.forge/templates/breakdown-tasks/T-test-3.md` with a valid `## Execution Workflow` section. A task file exists at `.forge/tasks/{task-id}.md` referencing that template. The working tree is clean.
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `task-cli run-tasks --task-file .forge/tasks/{task-id}.md` to dispatch the task through run-tasks, which delegates to task-executor
  2. After dispatch completes, verify task-executor received the correct template: `cat .forge/tasks/{task-id}/record.json | jq '.template'` should return the template file name
  3. Verify Step 2 used the Execution Workflow body: `cat .forge/tasks/{task-id}/record.json | jq '.steps[1].output'` should contain the workflow body text from the template
  4. Verify no TDD keywords: `grep -c "RED/GREEN/REFACTOR\|TDD implementation\|TDD cycle" .forge/tasks/{task-id}/record.json` returns `0`
  5. Verify Quality Gate was skipped: `cat .forge/tasks/{task-id}/record.json | jq '.steps[2].output'` shows record + commit, no compile/fmt/lint/test sequence
  6. Verify record committed: `cat .forge/tasks/{task-id}/record.json | jq '.status'` returns `"completed"` and a git commit exists with the record
- **Expected**: The task flows through the complete pipeline: run-tasks dispatches to task-executor, task-executor reads the Execution Workflow from the template file and injects it into Step 2, Step 2 executes the workflow without TDD, Step 3 records results and commits. Template name in record matches the dispatched template. Step 2 output contains workflow body text. Zero TDD keywords in the record. Status is `completed`. A git commit containing the record file exists.
- **Priority**: P0

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
| TC-017 | Story 1 / AC-1, Story 2 / AC-3 | CLI | cli/task-executor | P0 |

---

## Route Validation

_No route validation — this feature has no UI/API routes. All test cases verify CLI/agent behavior through file inspection and command execution._
