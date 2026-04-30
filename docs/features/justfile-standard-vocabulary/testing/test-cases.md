---
feature: "justfile-standard-vocabulary"
sources:
  - prd/prd-user-stories.md
  - prd/prd-spec.md
generated: "2026-04-30"
---

# Test Cases: justfile-standard-vocabulary

## Summary

| Type | Count |
|------|-------|
| UI   | 0     |
| API  | 0     |
| CLI  | 25    |
| **Total** | **25** |

---

## UI Test Cases

_No UI test cases — this feature has no user interface._

---

## API Test Cases

_No API test cases — this feature has no HTTP endpoints._

---

## CLI Test Cases

### TC-001: Skill commands use standard just verbs

- **Source**: Story 1 / AC-1
- **Type**: CLI
- **Target**: cli/skill-execution
- **Test ID**: cli/skill-execution/skill-commands-use-standard-just-verbs
- **Pre-conditions**: A forge project with a generated justfile exists
- **Steps**:
  1. Scan all skill/agent/command files for raw toolchain commands (`go test`, `npm run build`, `npx serve`, `go build`, `cargo build`)
  2. Verify each build/test/compile operation references `just <verb>` instead
- **Expected**: Zero occurrences of direct language toolchain commands in skill/agent/command files; all operations invoke `just <verb>`
- **Priority**: P0

### TC-002: Pure backend project executes correct toolchain via just test

- **Source**: Story 1 / AC-2
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/pure-backend-test-calls-go-test
- **Pre-conditions**: A pure backend Go project with `go.mod` exists; justfile generated via `/init-justfile` for backend type
- **Steps**:
  1. Run `just test` in the project directory
  2. Verify the underlying command is `go test -race ./...`
- **Expected**: `just test` invokes `go test -race ./...` and exits with code 0 when tests pass
- **Priority**: P0

### TC-003: Mixed project scope parameter targets frontend only

- **Source**: Story 1 / AC-3
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/mixed-project-scope-frontend-build
- **Pre-conditions**: A mixed project with both `package.json` (frontend) and `go.mod` (backend) exists; justfile generated via `/init-justfile` for mixed type
- **Steps**:
  1. Run `just build frontend`
  2. Verify only the frontend build executes (e.g., `npm run build` or equivalent)
  3. Verify backend build is not triggered
- **Expected**: Only frontend build commands run; backend build is skipped
- **Priority**: P0

### TC-004: Frontend project detection and scope-free justfile generation

- **Source**: Story 2 / AC-1
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/frontend-project-no-scope-params
- **Pre-conditions**: A project directory containing only `package.json` (no `go.mod`, `Cargo.toml`, or `pyproject.toml`)
- **Steps**:
  1. Run `/init-justfile`
  2. Inspect generated justfile for scope parameters in recipe signatures
  3. Run `just project-type`
- **Expected**: Generated justfile recipes have no scope parameter; `just project-type` outputs `frontend`
- **Priority**: P0

### TC-005: Backend project detection and scope-free justfile generation

- **Source**: Story 2 / AC-2
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/backend-project-no-scope-params
- **Pre-conditions**: A project directory containing only `go.mod` (no `package.json`)
- **Steps**:
  1. Run `/init-justfile`
  2. Inspect generated justfile for scope parameters in recipe signatures
  3. Run `just project-type`
- **Expected**: Generated justfile recipes have no scope parameter; `just project-type` outputs `backend`
- **Priority**: P0

### TC-006: Mixed project detection and scope-aware justfile generation

- **Source**: Story 2 / AC-3
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/mixed-project-with-scope-params
- **Pre-conditions**: A project directory containing both `package.json` and `go.mod`
- **Steps**:
  1. Run `/init-justfile`
  2. Inspect generated justfile for scope parameters in scoped recipe signatures
  3. Run `just project-type`
- **Expected**: Scoped recipes accept `frontend`/`backend` parameter; `just project-type` outputs `mixed`
- **Priority**: P0

### TC-007: Mixed project tasks receive scope field in index.json

- **Source**: Story 3 / AC-1
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/mixed-project-tasks-have-scope
- **Pre-conditions**: A mixed project with a finalized tech design document
- **Steps**:
  1. Run `/breakdown-tasks` on the mixed project tech design
  2. Open generated `index.json`
  3. Verify each task object contains a `scope` field
- **Expected**: Every task in `index.json` has a `scope` field with value `frontend`, `backend`, or `all`
- **Priority**: P0

### TC-008: Frontend-only task scope marked as frontend

- **Source**: Story 3 / AC-2
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/frontend-only-task-scope-frontend
- **Pre-conditions**: A mixed project with a task that only affects files in frontend directories (e.g., `web/src/components/Button.tsx`)
- **Steps**:
  1. Run `/breakdown-tasks`
  2. Find the task whose affected files are all in frontend directories
  3. Check the task's `scope` field
- **Expected**: Task scope is `frontend`
- **Priority**: P1

### TC-009: Cross-scope task marked as all

- **Source**: Story 3 / AC-3
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/cross-scope-task-scope-all
- **Pre-conditions**: A mixed project with a task that affects files in both frontend and backend directories
- **Steps**:
  1. Run `/breakdown-tasks`
  2. Find the task whose affected files span both frontend and backend directories
  3. Check the task's `scope` field
- **Expected**: Task scope is `all`
- **Priority**: P1

### TC-010: Non-mixed project tasks all receive scope all

- **Source**: Story 3 / AC-4
- **Type**: CLI
- **Target**: cli/breakdown-tasks
- **Test ID**: cli/breakdown-tasks/non-mixed-project-all-tasks-scope-all
- **Pre-conditions**: A non-mixed project (pure frontend or pure backend) with a finalized tech design
- **Steps**:
  1. Run `/breakdown-tasks`
  2. Open generated `index.json`
  3. Check every task's `scope` field
- **Expected**: All tasks have scope `all`
- **Priority**: P1

### TC-011: Successful command exits with code 0

- **Source**: Story 4 / AC-1
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/successful-command-exit-code-0
- **Pre-conditions**: A forge project with a valid generated justfile; code is in a passing state
- **Steps**:
  1. Run `just compile`
  2. Check exit code
  3. Verify stdout contains normal output (no errors)
- **Expected**: Exit code is 0; stdout contains compilation output without error messages
- **Priority**: P0

### TC-012: Failed command exits with non-zero code and stderr output

- **Source**: Story 4 / AC-2
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/failed-command-exit-nonzero-stderr
- **Pre-conditions**: A forge project with a generated justfile; code is in a failing state (e.g., type error introduced)
- **Steps**:
  1. Run `just compile`
  2. Check exit code
  3. Check stderr for error information
- **Expected**: Exit code is non-0; stderr contains error details
- **Priority**: P0

### TC-013: Compile with type errors outputs details to stderr

- **Source**: Story 4 / AC-3
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/compile-type-errors-stderr-details
- **Pre-conditions**: A forge project with code containing type errors
- **Steps**:
  1. Run `just compile`
  2. Check exit code is non-0
  3. Verify error details (file name, line number, error description) are present in stderr
- **Expected**: Exit code non-0; stderr contains parseable error information with file, line, and description
- **Priority**: P1

### TC-014: Consecutive commands all succeed with exit code 0

- **Source**: Story 4 / AC-4
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/consecutive-commands-all-succeed
- **Pre-conditions**: A forge project with valid generated justfile; code is in a passing state
- **Steps**:
  1. Run `just install`
  2. Verify exit code 0
  3. Run `just compile`
  4. Verify exit code 0
  5. Run `just test`
  6. Verify exit code 0
- **Expected**: All three commands exit with code 0; no human intervention required
- **Priority**: P1

### TC-015: Scope mismatch shows warning and falls back

- **Source**: Story 5 / AC-1
- **Type**: CLI
- **Target**: cli/skill-execution
- **Test ID**: cli/skill-execution/scope-mismatch-warning-fallback
- **Pre-conditions**: A pure backend project; a task with scope=frontend
- **Steps**:
  1. Skill reads task scope = `frontend`
  2. Skill executes `just project-type` and receives `backend`
  3. Skill detects mismatch
- **Expected**: Warning message `[forge] scope=frontend but project-type=backend; falling back to just build` is displayed; skill executes `just build` without scope parameter
- **Priority**: P0

### TC-016: Mixed project with matching scope executes normally

- **Source**: Story 5 / AC-2
- **Type**: CLI
- **Target**: cli/skill-execution
- **Test ID**: cli/skill-execution/mixed-project-matching-scope-normal
- **Pre-conditions**: A mixed project; a task with scope=frontend
- **Steps**:
  1. Skill reads task scope = `frontend`
  2. Skill executes `just project-type` and receives `mixed`
  3. Skill executes `just build frontend`
- **Expected**: Command `just build frontend` executes normally; no warning message
- **Priority**: P0

### TC-017: Mixed project with invalid scope exits with error

- **Source**: Spec 5.3 (scope validation rules, row 2)
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/mixed-invalid-scope-exits-error
- **Pre-conditions**: A mixed project with a generated justfile supporting scope parameters
- **Steps**:
  1. Run `just build foo` (invalid scope value)
  2. Check exit code
  3. Check stderr for error message
- **Expected**: Exit code is 1; stderr contains `[forge] invalid scope 'foo'; expected frontend/backend`
- **Priority**: P0

### TC-018: No marker files detected causes init-justfile to error

- **Source**: Spec 5.2 (project type detection rules)
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/no-marker-files-error
- **Pre-conditions**: An empty project directory with no `package.json`, `go.mod`, `Cargo.toml`, or `pyproject.toml`
- **Steps**:
  1. Run `/init-justfile` in the empty directory
- **Expected**: Command aborts with error message indicating no known project marker files detected
- **Priority**: P1

### TC-019: Existing justfile triggers user confirmation

- **Source**: Spec 5.2 (adaptive generation flow) and Spec (compatibility requirements)
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/existing-justfile-prompts-confirmation
- **Pre-conditions**: A project directory with an existing justfile that has no forge boundary markers
- **Steps**:
  1. Run `/init-justfile`
  2. Observe whether a confirmation prompt appears
- **Expected**: A confirmation prompt is shown before overwriting the existing justfile
- **Priority**: P1

### TC-020: Boundary markers present triggers idempotent merge

- **Source**: Spec (maintainability requirements)
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/boundary-markers-idempotent-merge
- **Pre-conditions**: A project with an existing justfile containing forge boundary markers (`# --- forge standard recipes ---` / `# --- end forge standard recipes ---`) and custom recipes outside markers
- **Steps**:
  1. Run `/init-justfile`
  2. Verify only the marked section is replaced
  3. Verify custom recipes outside markers are preserved
- **Expected**: Only the forge standard recipe section is updated; custom recipes remain unchanged
- **Priority**: P1

### TC-021: Project-type recipe outputs deterministic single word

- **Source**: Spec 5.1 (command vocabulary table) and Spec (agent-friendly: deterministic output)
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/project-type-deterministic-single-word
- **Pre-conditions**: A forge project with a generated justfile
- **Steps**:
  1. Run `just project-type`
  2. Capture stdout output
  3. Run `just project-type` again
  4. Compare outputs
- **Expected**: Both runs produce the same single-word output (`frontend`, `backend`, or `mixed`); no side effects from either run
- **Priority**: P1

### TC-022: All 15 standard commands are present in generated justfile

- **Source**: Spec 5.1 (command vocabulary table)
- **Type**: CLI
- **Target**: cli/init-justfile
- **Test ID**: cli/init-justfile/all-15-standard-commands-present
- **Pre-conditions**: A project with valid marker files for any project type
- **Steps**:
  1. Run `/init-justfile`
  2. List all recipe names in the generated justfile
  3. Verify all 15 standard commands are present: `compile`, `build`, `run`, `dev`, `test`, `test-e2e`, `lint`, `fmt`, `check`, `clean`, `install`, `ci`, `e2e-setup`, `e2e-verify`, `project-type`
- **Expected**: All 15 recipe names exist in the generated justfile
- **Priority**: P0

### TC-023: Just project-type failure triggers fallback in skill

- **Source**: Spec 5.3 (scope validation rules, row 4)
- **Type**: CLI
- **Target**: cli/skill-execution
- **Test ID**: cli/skill-execution/project-type-failure-fallback
- **Pre-conditions**: A project with an old justfile that has no `project-type` recipe; a task with scope=frontend
- **Steps**:
  1. Skill reads task scope = `frontend`
  2. Skill executes `just project-type` which fails (non-zero exit)
  3. Skill detects failure
- **Expected**: Skill logs `[forge] just project-type failed (exit N); falling back to just verb` and executes `just <verb>` without scope
- **Priority**: P1

### TC-024: Unexpected project-type output triggers fallback

- **Source**: Spec 5.3 (scope validation rules, row 5)
- **Type**: CLI
- **Target**: cli/skill-execution
- **Test ID**: cli/skill-execution/unexpected-project-type-fallback
- **Pre-conditions**: A project where `just project-type` returns an unexpected string (e.g., `unknown`)
- **Steps**:
  1. Skill reads task scope = `frontend`
  2. Skill executes `just project-type` and receives unexpected output (not `frontend`/`backend`/`mixed`)
  3. Skill detects unexpected output
- **Expected**: Skill logs `[forge] just project-type returned unexpected output 'XYZ'; falling back to just verb` and executes `just <verb>` without scope
- **Priority**: P1

### TC-025: Idempotent recipes produce no side effects on repeat runs

- **Source**: Spec (agent-friendly: idempotency)
- **Type**: CLI
- **Target**: cli/justfile
- **Test ID**: cli/justfile/idempotent-recipes-no-side-effects
- **Pre-conditions**: A forge project with a generated justfile; dependencies already installed
- **Steps**:
  1. Run `just install`
  2. Run `just install` again
  3. Run `just e2e-setup`
  4. Run `just e2e-setup` again
- **Expected**: All four commands exit with code 0; no errors or unexpected side effects from repeated execution
- **Priority**: P2

---

## Traceability

| TC ID | Source | Type | Target | Priority |
|-------|--------|------|--------|----------|
| TC-001 | Story 1 / AC-1 | CLI | cli/skill-execution | P0 |
| TC-002 | Story 1 / AC-2 | CLI | cli/justfile | P0 |
| TC-003 | Story 1 / AC-3 | CLI | cli/justfile | P0 |
| TC-004 | Story 2 / AC-1 | CLI | cli/init-justfile | P0 |
| TC-005 | Story 2 / AC-2 | CLI | cli/init-justfile | P0 |
| TC-006 | Story 2 / AC-3 | CLI | cli/init-justfile | P0 |
| TC-007 | Story 3 / AC-1 | CLI | cli/breakdown-tasks | P0 |
| TC-008 | Story 3 / AC-2 | CLI | cli/breakdown-tasks | P1 |
| TC-009 | Story 3 / AC-3 | CLI | cli/breakdown-tasks | P1 |
| TC-010 | Story 3 / AC-4 | CLI | cli/breakdown-tasks | P1 |
| TC-011 | Story 4 / AC-1 | CLI | cli/justfile | P0 |
| TC-012 | Story 4 / AC-2 | CLI | cli/justfile | P0 |
| TC-013 | Story 4 / AC-3 | CLI | cli/justfile | P1 |
| TC-014 | Story 4 / AC-4 | CLI | cli/justfile | P1 |
| TC-015 | Story 5 / AC-1 | CLI | cli/skill-execution | P0 |
| TC-016 | Story 5 / AC-2 | CLI | cli/skill-execution | P0 |
| TC-017 | Spec 5.3 / row 2 | CLI | cli/justfile | P0 |
| TC-018 | Spec 5.2 / detection | CLI | cli/init-justfile | P1 |
| TC-019 | Spec 5.2 / flow | CLI | cli/init-justfile | P1 |
| TC-020 | Spec / maintainability | CLI | cli/init-justfile | P1 |
| TC-021 | Spec 5.1 + agent-friendly | CLI | cli/justfile | P1 |
| TC-022 | Spec 5.1 / vocabulary | CLI | cli/init-justfile | P0 |
| TC-023 | Spec 5.3 / row 4 | CLI | cli/skill-execution | P1 |
| TC-024 | Spec 5.3 / row 5 | CLI | cli/skill-execution | P1 |
| TC-025 | Spec / idempotency | CLI | cli/justfile | P2 |

---

## Route Validation

_Omitted — this feature has no web routes. All test cases are CLI type targeting justfile commands, init-justfile skill, breakdown-tasks skill, and skill execution behavior._
