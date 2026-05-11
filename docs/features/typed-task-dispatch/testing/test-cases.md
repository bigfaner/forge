---
feature: "typed-task-dispatch"
sources:
  - docs/features/typed-task-dispatch/prd/prd-user-stories.md
  - docs/features/typed-task-dispatch/prd/prd-spec.md
generated: "2026-05-11"
---

# Test Cases: typed-task-dispatch

> **WARNING**: sitemap.json not found — Element set to `sitemap-missing`. This is a pure CLI feature; no UI tests apply.

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| **Integration** | **0** |
| API  | 0     |
| CLI  | 16    |
| **Total** | **16** |

> **Note**: This feature is a pure CLI/agent tool. No UI or API interfaces are exposed. All test cases are CLI type.

---

## CLI Test Cases

## TC-001: doc-generation.summary task executes without TDD steps

- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/run-tasks
- **Test ID**: cli/run-tasks/doc-generation-summary-task-executes-without-tdd-steps
- **Pre-conditions**: A task with `type: doc-generation.summary` exists in index.json with status `pending`; forge is installed and `run-tasks` is available
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a feature with an index.json containing a task with `type: doc-generation.summary` and status `pending`
  2. Run `run-tasks` (or invoke the dispatch path for that task)
  3. Observe the prompt delivered to task-executor
- **Expected**: The synthesized prompt delivered to task-executor contains no TDD-related steps (no RED/GREEN/REFACTOR, no `just test` invocation); the prompt reflects the doc-generation execution flow
- **Priority**: P0

---

## TC-002: fix task executes five-step diagnostic flow

- **Source**: Story 1 / AC-2
- **Type**: CLI
- **Target**: cli/run-tasks
- **Test ID**: cli/run-tasks/fix-task-executes-five-step-diagnostic-flow
- **Pre-conditions**: A task with `type: fix` exists in index.json with status `pending`; forge is installed
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a feature with an index.json containing a task with `type: fix` and status `pending`
  2. Run `run-tasks` to dispatch the fix task
  3. Inspect the prompt delivered to task-executor
- **Expected**: The prompt contains the five-step flow: diagnose → locate → fix → verify → commit; after task-executor executes, `go build ./...` or `go test ./...` exits with code 0
- **Priority**: P0

---

## TC-003: New type template generates correct prompt output

- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/new-type-template-generates-correct-prompt-output
- **Pre-conditions**: A new task type has been added by placing a markdown template file in the task-cli prompt templates directory and registering it in the type enum; a task of that new type exists in index.json
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Add a new markdown template file for the new type in the task-cli embed directory
  2. Register the new type in the type enum
  3. Create a task in index.json with the new type
  4. Run `task prompt <id>` for that task
  5. Inspect stdout
- **Expected**: stdout contains the task ID, scope, and execution steps from the new template file; stdout does not contain execution steps from any other type's template; Go unit tests can be written to cover the new type without modifying task-executor.md or any task template file
- **Priority**: P1

---

## TC-004: Template syntax error or unregistered type causes non-zero exit

- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/template-syntax-error-or-unregistered-type-causes-non-zero-exit
- **Pre-conditions**: A task exists in index.json referencing a type whose template file has a syntax error (e.g., malformed placeholder), or whose type is not registered in the enum
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a task in index.json with a type that is either unregistered or has a broken template
  2. Run `task prompt <id>` for that task
  3. Check exit code, stdout, and stderr
- **Expected**: Exit code is non-zero; stderr contains the template file path and a specific error reason; stdout is empty
- **Priority**: P1

---

## TC-005: task prompt outputs complete synthesized prompt within 500ms

- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/task-prompt-outputs-complete-synthesized-prompt-within-500ms
- **Pre-conditions**: Current feature has an in_progress task of any valid type in index.json; `.forge/state.json` is present with the correct feature slug
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Ensure `.forge/state.json` references the current feature
  2. Run `task prompt <id>` for an in_progress task
  3. Measure elapsed time and inspect stdout
- **Expected**: stdout contains the complete synthesized prompt including: task ID, scope, task type, execution steps from the corresponding type template, and phase summary path (if the task is the first task of a new phase); command completes in under 500ms; exit code is 0
- **Priority**: P0

---

## TC-006: task prompt with missing type or missing template outputs error to stderr

- **Source**: Story 3 / AC-2
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/task-prompt-with-missing-type-or-missing-template-outputs-error-to-stderr
- **Pre-conditions**: A task exists in index.json with either no `type` field, or a `type` value for which no template file exists
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create a task in index.json with a missing or invalid type
  2. Run `task prompt <id>` for that task
  3. Check exit code, stdout, and stderr
- **Expected**: Error description is written to stderr; exit code is non-zero; stdout is empty
- **Priority**: P0

---

## TC-007: task migrate fills type fields correctly for old index.json

- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/task-migrate
- **Test ID**: cli/task-migrate/task-migrate-fills-type-fields-correctly-for-old-indexjson
- **Pre-conditions**: An index.json exists without any `type` fields; all tasks have status other than `in_progress` (e.g., pending, completed, blocked)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Prepare an index.json with tasks of various IDs (e.g., `1.summary`, `1.gate`, `T-test-1`, `T-test-2`, `fix-1`, `disc-2`, `1.1`) and no `type` fields
  2. Run `task migrate`
  3. Inspect the updated index.json
  4. Run `task validate`
- **Expected**: All tasks in index.json have a `type` field populated according to the inference rules (e.g., `.summary` → `doc-generation.summary`, `.gate` → `gate`, `T-test-1` → `test-pipeline.gen-cases`, `fix-` prefix → `fix`, others → `implementation`); task statuses are unchanged; `task validate` reports no errors
- **Priority**: P0

---

## TC-008: task migrate rejects in_progress tasks

- **Source**: Story 4 / AC-2
- **Type**: CLI
- **Target**: cli/task-migrate
- **Test ID**: cli/task-migrate/task-migrate-rejects-in-progress-tasks
- **Pre-conditions**: An index.json exists with at least one task in `in_progress` status
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Prepare an index.json with one or more tasks having `status: in_progress`
  2. Run `task migrate`
  3. Check exit code, stderr output, and index.json content
- **Expected**: Command exits with a non-zero code and outputs an error message prompting the user to complete or manually resolve in-progress tasks before migrating; index.json is not modified
- **Priority**: P0

---

## TC-009: breakdown-tasks generates index.json with type fields for all tasks

- **Source**: Story 5 / AC-1
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/breakdown-tasks-generates-indexjson-with-type-fields-for-all-tasks
- **Pre-conditions**: A complete tech-design document exists for the feature; breakdown-tasks skill is available
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Provide a complete tech-design document
  2. Run `breakdown-tasks`
  3. Inspect the generated index.json
  4. Run `task validate`
- **Expected**: Every task in the generated index.json has a `type` field; `task validate` reports no errors; type values are consistent with the actual task nature (e.g., implementation tasks have `type: implementation`, summary tasks have `type: doc-generation.summary`, gate tasks have `type: gate`)
- **Priority**: P1

---

## TC-010: breakdown-tasks falls back to implementation for unrecognized task descriptions

- **Source**: Story 5 / AC-2
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/breakdown-tasks-falls-back-to-implementation-for-unrecognized-task-descriptions
- **Pre-conditions**: A tech-design document contains at least one task whose description does not match any known type inference rule
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Provide a tech-design document with an ambiguous or novel task description
  2. Run `breakdown-tasks`
  3. Inspect the generated index.json and stderr output
- **Expected**: The ambiguous task has `type: implementation` in the generated index.json; stderr contains a warning with the task ID and the reason it could not be inferred; the overall generation flow completes without interruption and all other tasks are generated correctly
- **Priority**: P1

---

## TC-011: execute-task uses task prompt routing instead of TASK_FILE + NO_TEST

- **Source**: Story 6 / AC-1
- **Type**: CLI
- **Target**: cli/execute-task
- **Test ID**: cli/execute-task/execute-task-uses-task-prompt-routing-instead-of-task-file-no-test
- **Pre-conditions**: index.json contains a task with `type: implementation` and status `pending`; execute-task skill is available
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a feature with an implementation task in index.json
  2. Invoke execute-task for that task
  3. Inspect the parameters passed to forge:task-executor
- **Expected**: execute-task calls `task prompt <id>` to synthesize the prompt and passes it as the `prompt` parameter to `Agent(forge:task-executor)`; the old `TASK_FILE` and `NO_TEST` parameter combination is not used; the prompt content received by task-executor is identical to what run-tasks would produce for the same task
- **Priority**: P0

---

## TC-012: execute-task marks task blocked when task prompt fails

- **Source**: Story 6 / AC-2
- **Type**: CLI
- **Target**: cli/execute-task
- **Test ID**: cli/execute-task/execute-task-marks-task-blocked-when-task-prompt-fails
- **Pre-conditions**: index.json contains a task whose `task prompt <id>` call will fail (e.g., missing type field or missing template)
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a task in index.json with a missing or invalid type
  2. Invoke execute-task for that task
  3. Check the task status in index.json and stderr output
- **Expected**: execute-task marks the task as `blocked` in index.json; stderr contains the error message from `task prompt`; the failure is not silent
- **Priority**: P0

---

## TC-013: run-tasks dispatches fix task via task prompt with five-step prompt

- **Source**: Story 7 / AC-1
- **Type**: CLI
- **Target**: cli/run-tasks
- **Test ID**: cli/run-tasks/run-tasks-dispatches-fix-task-via-task-prompt-with-five-step-prompt
- **Pre-conditions**: index.json contains a task with `type: fix` and status `pending`; run-tasks.md does not contain any dispatch call to `forge:error-fixer`
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up a feature with a fix task in index.json
  2. Run `run-tasks`
  3. Inspect the prompt delivered to task-executor and the content of run-tasks.md
- **Expected**: task-executor receives a prompt generated by `task prompt <id>` that contains the five-step flow (diagnose → locate → fix → verify → commit); run-tasks.md contains no dispatch call to `forge:error-fixer`
- **Priority**: P0

---

## TC-014: run-tasks uses --fix-record-missed when record file is absent

- **Source**: Story 7 / AC-2
- **Type**: CLI
- **Target**: cli/run-tasks
- **Test ID**: cli/run-tasks/run-tasks-uses-fix-record-missed-when-record-file-is-absent
- **Pre-conditions**: A task has completed execution but its record file is missing from `tasks/records/`; run-tasks is configured to detect this condition
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Simulate a task completion where the record file was not written
  2. Run `run-tasks` (or trigger the record-missing detection path)
  3. Inspect which command is called and which agent is dispatched
- **Expected**: run-tasks calls `task prompt <id> --fix-record-missed` and dispatches the resulting prompt to `forge:task-executor`; `forge:error-fixer` is not called; error-fixer.md dispatch entry is absent from run-tasks.md
- **Priority**: P0

---

## TC-015: task validate accepts valid type enum values and rejects invalid ones

- **Source**: PRD Spec / Functional Specs — task validate command extension
- **Type**: CLI
- **Target**: cli/task-validate
- **Test ID**: cli/task-validate/task-validate-accepts-valid-type-enum-values-and-rejects-invalid-ones
- **Pre-conditions**: index.json exists with tasks; `task validate` command is available
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Create an index.json where all tasks have valid `type` values from the 11-type enum
  2. Run `task validate` — verify it reports no errors
  3. Modify one task to have an invalid `type` value (e.g., `"type": "unknown-type"`)
  4. Run `task validate` again — verify it reports an error
  5. Remove the `type` field from a task entirely
  6. Run `task validate` — verify it reports a missing required field error
- **Expected**: Step 2 exits 0 with no errors; step 4 exits non-zero with an error identifying the invalid type value and task ID; step 6 exits non-zero with an error identifying the missing required `type` field and task ID
- **Priority**: P1

---

## TC-016: task prompt phase boundary detection injects phase summary path

- **Source**: PRD Spec / task prompt internal flow — phase boundary detection
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/task-prompt-phase-boundary-detection-injects-phase-summary-path
- **Pre-conditions**: index.json contains completed tasks from phase N and a pending task that is the first task of phase N+1; a phase summary file exists for phase N
- **Route**: N/A
- **Element**: sitemap-missing
- **Steps**:
  1. Set up index.json with completed tasks in phase 1 and a pending task as the first task of phase 2
  2. Run `task prompt <id>` for the first phase-2 task
  3. Inspect stdout
- **Expected**: stdout contains the phase summary path for phase 1 injected into the prompt; for a task that is not the first task of a new phase, the phase summary path is absent from the prompt
- **Priority**: P1

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/run-tasks | P0 |
| TC-002 | Story 1 / AC-2 | CLI | cli/run-tasks | P0 |
| TC-003 | Story 2 / AC-1 | CLI | cli/task-prompt | P1 |
| TC-004 | Story 2 / AC-2 | CLI | cli/task-prompt | P1 |
| TC-005 | Story 3 / AC-1 | CLI | cli/task-prompt | P0 |
| TC-006 | Story 3 / AC-2 | CLI | cli/task-prompt | P0 |
| TC-007 | Story 4 / AC-1 | CLI | cli/task-migrate | P0 |
| TC-008 | Story 4 / AC-2 | CLI | cli/task-migrate | P0 |
| TC-009 | Story 5 / AC-1 | CLI | cli/breakdown-tasks | P1 |
| TC-010 | Story 5 / AC-2 | CLI | cli/breakdown-tasks | P1 |
| TC-011 | Story 6 / AC-1 | CLI | cli/execute-task | P0 |
| TC-012 | Story 6 / AC-2 | CLI | cli/execute-task | P0 |
| TC-013 | Story 7 / AC-1 | CLI | cli/run-tasks | P0 |
| TC-014 | Story 7 / AC-2 | CLI | cli/run-tasks | P0 |
| TC-015 | PRD Spec / task validate extension | CLI | cli/task-validate | P1 |
| TC-016 | PRD Spec / task prompt phase boundary detection | CLI | cli/task-prompt | P1 |
