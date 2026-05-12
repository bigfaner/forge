---
feature: "typed-task-dispatch"
sources:
  - docs/features/typed-task-dispatch/prd/prd-user-stories.md
  - docs/features/typed-task-dispatch/prd/prd-spec.md
generated: "2026-05-11"
---

# Test Cases: typed-task-dispatch

> **Note**: This is a pure CLI feature; no UI tests apply. Element is N/A for all test cases.

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| **Integration** | **0** |
| API  | 0     |
| CLI  | 20    |
| **Total** | **20** |

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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with content: `{"tasks": [{"id": "1.summary", "title": "Write summary doc", "type": "doc-generation.summary", "status": "pending"}]}`; set `.forge/state.json` to reference feature slug `test-feature`
  2. Run `task prompt 1.summary` and capture stdout to `prompt_out.txt`
  3. Inspect `prompt_out.txt`
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with content: `{"tasks": [{"id": "fix-1", "type": "fix", "status": "pending"}]}`; set `.forge/state.json` to `{"feature": "test-feature"}`
  2. Run `task prompt fix-1` and capture stdout to `prompt_out.txt`
  3. Inspect `prompt_out.txt` for the five-step flow
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
- **Element**: N/A
- **Steps**:
  1. Create `plugins/forge/task-cli/prompts/custom-audit.md` with content: `# custom-audit\n\nSteps:\n1. Run audit\n2. Report findings`; add `"custom-audit"` to the type enum in `pkg/task/type.go`
  2. Create `docs/features/test-feature/tasks/index.json` with content: `{"tasks": [{"id": "1.1", "type": "custom-audit", "status": "pending"}]}`; set `.forge/state.json` to `{"feature": "test-feature"}`
  3. Run `task prompt 1.1` and capture stdout to `prompt_out.txt`
  4. Inspect `prompt_out.txt`
- **Expected**: stdout contains the task ID, scope, and execution steps from the new template file; stdout does not contain execution steps from any other type's template; running `go test ./pkg/prompt/...` exits 0 and covers the new type
- **Priority**: P1

---

## TC-004: Template syntax error or unregistered type causes non-zero exit

- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/template-syntax-error-or-unregistered-type-causes-non-zero-exit
- **Pre-conditions**: A task exists in index.json referencing a type whose template file has a syntax error (e.g., malformed placeholder), or whose type is not registered in the enum
- **Route**: N/A
- **Element**: N/A
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with content: `{"tasks": [{"id": "1.1", "type": "implementation", "status": "in_progress"}]}`; set `.forge/state.json` to `{"feature": "test-feature"}`
  2. Run `task prompt 1.1` for an in_progress task
  3. Run `time task prompt 1.1 > /dev/null` and capture the real elapsed time reported by `time`
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
- **Element**: N/A
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with tasks of various IDs (`1.summary`, `1.gate`, `T-test-1`, `T-test-2`, `fix-1`, `disc-2`, `1.1`) and no `type` fields on any task
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with content: `{"tasks": [{"id": "1.1", "type": "implementation", "status": "in_progress"}, {"id": "1.2", "type": "fix", "status": "in_progress"}]}`; set `.forge/state.json` to `{"feature": "test-feature"}`
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/design/tech-design.md` with a complete tech-design document covering implementation tasks, a summary task, a gate task, and at least one test task
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/design/tech-design.md` with content that includes a task described as "Investigate and document the optimal caching strategy" — a description that does not match any known type inference rule
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with a task `{"id": "1.1", "type": "implementation", "status": "pending"}`; set `.forge/state.json` to reference feature slug `test-feature`
  2. Run `task prompt 1.1` and capture stdout to `execute_task_prompt.txt`
  3. Run `task prompt 1.1` again via the run-tasks dispatch path and capture stdout to `run_tasks_prompt.txt`
  4. Diff `execute_task_prompt.txt` against `run_tasks_prompt.txt`
- **Expected**: `diff execute_task_prompt.txt run_tasks_prompt.txt` produces no output (files are identical); neither file contains `TASK_FILE` or `NO_TEST` strings; exit code of both `task prompt` calls is 0
- **Priority**: P0

---

## TC-012: execute-task marks task blocked when task prompt fails

- **Source**: Story 6 / AC-2
- **Type**: CLI
- **Target**: cli/execute-task
- **Test ID**: cli/execute-task/execute-task-marks-task-blocked-when-task-prompt-fails
- **Pre-conditions**: index.json contains a task whose `task prompt <id>` call will fail (e.g., missing type field or missing template)
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with content: `{"tasks": [{"id": "1.1", "title": "Fix login bug", "status": "pending"}]}`; set `.forge/state.json` to `{"feature": "test-feature"}` (note: no `type` field on the task)
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with a task `{"id": "fix-1", "type": "fix", "status": "pending"}`; set `.forge/state.json` to reference feature slug `test-feature`
  2. Run `task prompt fix-1` and capture stdout to `prompt_out.txt`
  3. Inspect `prompt_out.txt` for the five-step flow and inspect `plugins/forge/skills/run-tasks.md` for any `forge:error-fixer` dispatch call
- **Expected**: task-executor receives a prompt generated by `task prompt <id>` that contains the five-step flow (diagnose → locate → fix → verify → commit); run-tasks.md contains no dispatch call to `forge:error-fixer`
- **Priority**: P0

---

## TC-014: run-tasks uses --fix-record-missed when record file is absent

- **Source**: Story 7 / AC-2
- **Type**: CLI
- **Target**: cli/run-tasks
- **Test ID**: cli/run-tasks/run-tasks-uses-fix-record-missed-when-record-file-is-absent
- **Pre-conditions**: A task with `id: fix-1` exists in index.json with `status: in_progress`; `tasks/records/fix-1.md` does not exist (was never written); `.forge/state.json` references the current feature slug
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with `{"id": "fix-1", "type": "fix", "status": "in_progress"}`; ensure `docs/features/test-feature/tasks/records/fix-1.md` does not exist
  2. Run `run-tasks` (or trigger the record-missing detection path by running `task check-records fix-1`)
  3. Capture stdout and inspect which command is invoked
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
- **Element**: N/A
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
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with content: `{"tasks": [{"id": "1.1", "type": "implementation", "status": "completed"}, {"id": "1.2", "type": "implementation", "status": "completed"}, {"id": "2.1", "type": "implementation", "status": "pending"}]}`; create `docs/features/test-feature/tasks/process/phase-1-summary.md` with any non-empty content; set `.forge/state.json` to `{"feature": "test-feature"}`
  2. Run `task prompt <id>` for the first phase-2 task
  3. Inspect stdout
- **Expected**: stdout contains the phase summary path for phase 1 injected into the prompt; for a task that is not the first task of a new phase, the phase summary path is absent from the prompt
- **Priority**: P1

---

## TC-017: eval-cases task executes in main session without subagent dispatch

- **Source**: PRD Spec §Scope — type == test-pipeline.eval-cases permanent exception
- **Type**: CLI
- **Target**: cli/run-tasks
- **Test ID**: cli/run-tasks/eval-cases-task-executes-in-main-session-without-subagent-dispatch
- **Pre-conditions**: index.json contains a task with `type: test-pipeline.eval-cases` and status `pending`; `.forge/state.json` references the current feature slug
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with `{"id": "T-eval-1", "type": "test-pipeline.eval-cases", "status": "pending"}`; set `.forge/state.json` to reference feature slug `test-feature`
  2. Run `run-tasks` and capture the full stdout/stderr log
  3. Search the log for any `Agent(forge:task-executor)` or subagent dispatch call targeting `T-eval-1`
- **Expected**: The log contains no `Agent(forge:task-executor)` call for task `T-eval-1`; the prompt for `T-eval-1` is executed directly in the main session (stdout shows the prompt content being acted on inline); exit code is 0
- **Priority**: P0

---

## TC-018: task prompt --fix-record-missed outputs record-recovery prompt to stdout

- **Source**: PRD Spec §Scope — task prompt `--fix-record-missed` mode
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/task-prompt-fix-record-missed-outputs-record-recovery-prompt
- **Pre-conditions**: A task with `id: 1.1` exists in index.json with `status: in_progress`; `tasks/records/1.1.md` does not exist; `.forge/state.json` references the current feature slug
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with `{"id": "1.1", "type": "implementation", "status": "in_progress"}`; ensure `docs/features/test-feature/tasks/records/1.1.md` does not exist
  2. Run `task prompt 1.1 --fix-record-missed` and capture stdout to `fix_record_prompt.txt`
  3. Inspect `fix_record_prompt.txt`
- **Expected**: stdout contains instructions to reconstruct the missing record file for task `1.1`, including the record file path `tasks/records/1.1.md` and the required record format; stdout does not contain the standard implementation execution steps; exit code is 0
- **Priority**: P0

---

## TC-019: quick-tasks generates index.json with type fields for all tasks

- **Source**: PRD Spec §Scope — quick-tasks skill generates tasks with type set automatically
- **Type**: CLI
- **Target**: cli/quick-tasks
- **Test ID**: cli/quick-tasks/quick-tasks-generates-indexjson-with-type-fields-for-all-tasks
- **Pre-conditions**: A feature proposal or brief description is available; quick-tasks skill is installed
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Run `quick-tasks` with a brief feature description (e.g., "add a config reload command")
  2. Inspect the generated `tasks/index.json`
  3. Run `task validate`
- **Expected**: Every task in the generated index.json has a non-empty `type` field; `task validate` exits 0 with no errors; type values match the task nature (e.g., implementation tasks have `type: implementation`, any summary task has `type: doc-generation.summary`)
- **Priority**: P1

---

## TC-020: task prompt marks task blocked when .forge/state.json is missing or unreadable

- **Source**: PRD Spec §Blocked State Lifecycle — state.json 读取失败 triggers blocked state
- **Type**: CLI
- **Target**: cli/task-prompt
- **Test ID**: cli/task-prompt/task-prompt-marks-task-blocked-when-state-json-missing
- **Pre-conditions**: index.json contains a task with `id: 1.1` and status `pending`; `.forge/state.json` does not exist or has permissions set to 000
- **Route**: N/A
- **Element**: N/A
- **Steps**:
  1. Create `docs/features/test-feature/tasks/index.json` with `{"id": "1.1", "type": "implementation", "status": "pending"}`
  2. Delete `.forge/state.json` (or run `chmod 000 .forge/state.json` to make it unreadable)
  3. Run `task prompt 1.1`
  4. Check exit code, stderr, and the `status` and `blocked_reason` fields of task `1.1` in index.json
- **Expected**: Exit code is non-zero; stderr contains a message referencing `.forge/state.json` and the read failure; task `1.1` in index.json has `status: blocked` and a non-empty `blocked_reason` field identifying the state.json failure; stdout is empty
- **Priority**: P0

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
| TC-017 | PRD Spec §Scope — eval-cases permanent exception | CLI | cli/run-tasks | P0 |
| TC-018 | PRD Spec §Scope — task prompt --fix-record-missed mode | CLI | cli/task-prompt | P0 |
| TC-019 | PRD Spec §Scope — quick-tasks type auto-generation | CLI | cli/quick-tasks | P1 |
| TC-020 | PRD Spec §Blocked State Lifecycle — state.json read failure | CLI | cli/task-prompt | P0 |
