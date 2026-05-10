---
created: 2026-05-10
prd: prd/prd-spec.md
status: Draft
---

# Technical Design: task-executor-skeleton

## Overview

Replace task-executor's hardcoded TDD (Steps 2-3) with a template-driven workflow mechanism. Every task template declares its own `## Execution Workflow`; task-executor reads and follows it. The default business task template (`breakdown-tasks/templates/task.md`) contains the TDD + Quality Gate workflow as its `## Execution Workflow`. Remove `noTest`/`NO_TEST` from the entire codebase.

The change affects 3 categories of files: (1) agent prompts, (2) task-cli Go code, (3) task templates and schemas.

## Architecture

### Step Renumbering

```
Current (6 steps):
  Step 0: Claim → Step 1: Read task → Step 2: TDD → Step 3: Quality Gate → Step 4: Record → Step 5: Commit
  NO_TEST=true → skip Steps 2-3

New (5 steps):
  Step 0: Claim → Step 1: Read task → Step 2: Execute Workflow → Step 3: Record → Step 4: Commit
  ## Execution Workflow present → follow it
  ## Execution Workflow absent  → fallback to default template (task.md) → follow its workflow
```

Old Steps 2+3 merge into new Step 2. Old Steps 4+5 become new Steps 3+4.

### Component Diagram

```
┌─────────────┐    claim output     ┌──────────────────┐
│  task-cli   │ ──────────────────> │  run-tasks.md     │
│  (Go)       │  (no NO_TEST)       │  (dispatcher)     │
└─────────────┘                     └────────┬─────────┘
                                             │ dispatch (no NO_TEST param)
                                             ▼
                                    ┌──────────────────┐
                                    │ task-executor.md  │
                                    │ (skeleton prompt) │
                                    │                   │
                                    │ Step 2:           │
                                    │ 1. Read task file │
                                    │ 2. Find ## Exec   │
                                    │    Workflow       │
                                    │ 3. Follow it OR   │
                                    │    fallback to    │
                                    │    default tmpl   │
                                    └────────┬─────────┘
                                             │ reads
                              ┌──────────────┼──────────────┐
                              ▼              ▼              ▼
                    ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
                    │ Task file    │ │ Task file    │ │ (no workflow)│
                    │ (TDD tasks)  │ │ (exec tasks) │ │              │
                    │              │ │              │ │              │
                    │ ## Execution │ │ ## Execution │ │ fallback:    │
                    │ Workflow:    │ │ Workflow:    │ │ read default │
                    │ TDD + QG     │ │ run+analyze  │ │ template     │
                    └──────────────┘ └──────────────┘ └──────┬───────┘
                                                            │
                                                            ▼
                                                  ┌──────────────┐
                                                  │ task.md      │
                                                  │ (default     │
                                                  │  template)   │
                                                  │              │
                                                  │ ## Execution │
                                                  │ Workflow:    │
                                                  │ TDD + QG     │
                                                  └──────────────┘
```

### Template Hierarchy

All task templates managed in `breakdown-tasks/templates/`:

| Template | Workflow Type | noTest (before) |
|----------|--------------|-----------------|
| `task.md` | **DEFAULT** — TDD + Quality Gate | `false` |
| `gate-task.md` | Verification + Quality Gate | `false` |
| `gen-test-cases.md` | Generate documentation | `true` |
| `eval-test-cases.md` | Evaluate + iterate | `true` |
| `gen-test-scripts.md` | Generate + verify compilation | `false` |
| `run-e2e-tests.md` | Execute + classify failures | `false` |
| `graduate-tests.md` | Migrate + verify | `false` |
| `verify-regression.md` | Run full suite | `false` |
| `consolidate-specs.md` | Extract + merge | `true` |
| `phase-summary-task.md` | Summarize phase | `true` |

Quick-tasks templates mirror a subset (6 templates, similar pattern).

### Dependencies

**Removed**: `noTest` field dependency chain:
```
types.go (NoTest) → claim.go (copy + print) → record.go (3 behaviors) → errors.go (suggestion)
                  → run-tasks.md (parse + dispatch) → task-executor.md (conditional)
```

**No new dependencies**. The workflow mechanism uses the existing task-file-reading flow. The fallback reads `plugins/forge/skills/breakdown-tasks/templates/task.md` — an existing file.

## Interfaces

### Interface 1: task-executor.md Step 2 (Skeleton)

The new Step 2 contains zero hardcoded TDD logic:

```markdown
## Step 2: Execute Workflow

<EXTREMELY-IMPORTANT>
You MUST determine the execution workflow for this task by following this exact procedure:

1. Read the task file specified in Step 1.
2. Search for a `## Execution Workflow` heading in the task file.
3. Based on what you find:

   **CASE A — `## Execution Workflow` heading exists with non-empty content:**
   The content under the heading (excluding the heading line itself, up to the next
   `##` heading or end of file) is your execution instructions. Follow these steps
   EXACTLY. Do not deviate, add, or skip steps.

   **CASE B — No `## Execution Workflow` heading found:**
   Read the default workflow template at:
   `plugins/forge/skills/breakdown-tasks/templates/task.md`
   Find its `## Execution Workflow` section and follow those steps.

   **CASE C — `## Execution Workflow` heading exists but content is empty:**
   Log: "WARNING: ## Execution Workflow heading present but empty. Falling back to default template."
   Then proceed as Case B.

4. Output after execution:
   - Success: `Step 2/4: [workflow description]... DONE`
   - Failure: `Step 2/4: [workflow description]... FAILED: [reason]`
</EXTREMELY-IMPORTANT>
```

### Interface 2: Task Template `## Execution Workflow` Section

Every task template adds this section after existing content. Example for the default template (`task.md`):

```markdown
## Execution Workflow

1. Write failing tests for each acceptance criterion (RED phase).
2. Implement minimum code to make all tests pass (GREEN phase).
3. Refactor while keeping tests green (REFACTOR phase).
4. Run quality gate in strict sequence, stopping at first failure:
   - `just compile [scope]`
   - `just fmt [scope]`
   - `just lint [scope]`
   - `just test [scope]`
5. If compile fails: fix, retry from compile.
   If fmt fails: task is blocked (formatting requires manual review).
   If lint fails: self-fix once, then blocked if still failing.
   If test fails: fix, retry from compile.
6. Stop. Proceed to Step 3 (Record).
```

**Workflow Content Model** — each `## Execution Workflow` section MUST satisfy this checklist. A template reviewer (or future linter) checks each item:

| # | Rule | Example (valid) | Counter-example (invalid) |
|---|------|-----------------|---------------------------|
| W1 | Steps are numbered and sequential | `1. Write failing tests...` | Bullet list without order |
| W2 | Each step names a concrete command or action | `Run just compile [scope]` | "Make sure compilation works" |
| W3 | Each step states success criteria | "All tests pass (exit 0)" | "Tests should be good" |
| W4 | Failure handling is explicit per step | "If lint fails: self-fix once, then blocked" | "Fix if broken" |
| W5 | Terminal step is a stop condition (no loops) | `6. Stop. Proceed to Step 3.` | "Repeat until all tests pass" |

**Automated validation**: `grep -c "^## Execution Workflow" templates/*.md` confirms every template has the heading. Content compliance (W1-W5) is enforced in template review with the checklist above — a future linter could parse numbered steps and detect missing `Stop` / `Proceed to Step` markers.

### Interface 3: task-cli Claim Output

`NoTest` field removed from both structs and the claim print output.

**types.go struct change:**

```go
// BEFORE (Task struct, line ~32):
NoTest bool `json:"noTest,omitempty"`

// BEFORE (TaskState struct, line ~120):
NoTest bool `json:"noTest,omitempty"`

// AFTER: both fields deleted entirely.
```

**claim.go function change:**

```go
// BEFORE (printTaskDetails, line ~310-325):
func printTaskDetails(key string, t *task.Task, projectRoot, featureSlug string) {
    PrintField("KEY", key)
    PrintField("TASK_ID", t.ID)
    // ... other fields ...
    PrintField("NO_TEST", strconv.FormatBool(t.NoTest))  // REMOVED
    PrintFieldIfNotEmpty("FEATURE", featureSlug)
    // ...
}

// AFTER:
func printTaskDetails(key string, t *task.Task, projectRoot, featureSlug string) {
    PrintField("KEY", key)
    PrintField("TASK_ID", t.ID)
    // ... other fields (unchanged) ...
    // NO_TEST line deleted
    PrintFieldIfNotEmpty("FEATURE", featureSlug)
    // ...
}
```

**claim.go state bootstrap change (executeClaim, line ~115-129):**

```go
// BEFORE:
state := &task.TaskState{
    // ...
    NoTest: t.NoTest,  // REMOVED
}

// AFTER:
state := &task.TaskState{
    // ... (NoTest line deleted)
}
```

Downstream consumers (run-tasks.md, execute-task.md) stop parsing the `NO_TEST` field from claim output.

### Interface 4: task-cli Record Behavior

The `NoTest` bypass is removed. Record validation applies uniformly to all tasks.

**record.go line ~113-116 deletion:**

```go
// BEFORE:
if t.NoTest && rd.Coverage >= 0 && rd.TestsPassed == 0 && rd.TestsFailed == 0 {
    rd.Coverage = -1.0
}

// AFTER: deleted entirely.
```

**record.go line ~121-124 quality gate condition change:**

```go
// BEFORE:
if rd.Status == "completed" && !recordForce && !t.NoTest {
    validateQualityGate(projectRoot, t.Scope)
}

// AFTER:
if rd.Status == "completed" && !recordForce {
    validateQualityGate(projectRoot, t.Scope)
}
```

**record.go template rendering change (fillRecordTemplate, line ~380):**

```go
// BEFORE:
formatTestsExecuted(rd.Coverage, t.NoTest)  // NoTest param
// AFTER:
formatTestsExecuted(rd.Coverage)             // NoTest param removed
```

**errors.go ErrNoTestEvidence hint update (line ~235):**

```go
// BEFORE:
"Either (1) run tests and report results, (2) add noTest:true to the task template, or (3) use --force to override",
// AFTER:
"Either (1) run tests and report results, or (2) use --force to override",
```

Rationale: The workflow itself determines what validation to run during Step 2. task-cli's record pre-check is a separate safety net. After removal, documentation-only tasks use `task record --force` or include a lightweight validation step (e.g., `just compile`) in their workflow.

## Data Models

Single-layer feature. Cross-Layer Data Map not applicable.

### Removed: NoTest Field

```
Task.NoTest      bool  // REMOVED from pkg/task/types.go:32
TaskState.NoTest bool  // REMOVED from pkg/task/types.go:120
```

No new data models. The `## Execution Workflow` content is embedded in task template markdown — no schema field needed for the workflow content. Only the `noTest` field is removed from schemas.

## Error Handling

### TaskStatus Enum (existing, unchanged)

Task statuses are already defined in `types.go` via `TaskIndex.StatusEnum`:

```go
// NewTaskIndex creates a new TaskIndex with default enum values.
func NewTaskIndex(feature string) *TaskIndex {
    return &TaskIndex{
        StatusEnum: []string{
            "pending",
            "in_progress",
            "completed",
            "blocked",
            "skipped",
            "rejected",
        },
        // ...
    }
}
```

These are the only valid values for `Task.Status` and `RecordData.Status`.

### Error Codes (existing, extended)

The existing `ErrorCode` type in `errors.go` defines codes as string constants. Two changes for this feature:

```go
// EXISTING (unchanged):
type ErrorCode string

const (
    ErrNoProject   ErrorCode = "NO_PROJECT"
    ErrNoFeature   ErrorCode = "NO_FEATURE"
    ErrInvalidInput ErrorCode = "INVALID_INPUT"
    ErrNotFound    ErrorCode = "NOT_FOUND"
    ErrConflict    ErrorCode = "CONFLICT"
    ErrValidation  ErrorCode = "VALIDATION_ERROR"
)

// CHANGE: update ErrNoTestEvidence hint to remove noTest option.
// No new error codes needed — all error cases map to existing codes.
```

### Error Propagation

Errors propagate through two channels:

1. **task-cli (Go)**: errors are `*AIError` structs with `Code`, `Message`, `Cause`, `Hint`, `Action` fields. Printed to stderr in a machine-readable format (`ERROR_CODE: ...`, `ERROR: ...`, `HINT: ...`). No error codes cross the Go/agent boundary.

2. **Agent prompts (markdown)**: errors are expressed as step output strings:
   - Success: `Step 2/4: [workflow description]... DONE`
   - Failure: `Step 2/4: [workflow description]... FAILED: [reason]`
   The failure reason is freeform text written into the task record by the agent.

### Agent-to-task-cli Error Translation

The agent does not rely on freeform text parsing. The error handoff works as follows:

1. **Agent sets status directly**: When a workflow step fails, the agent writes the task record file with `status: failed` (or `blocked`) in the YAML frontmatter and includes the failure reason in the `## Notes` section.
2. **Agent calls `task record`**: The agent invokes `task record --status failed` (or `blocked`), passing the status as a structured argument. The `notes` field carries the failure reason. task-cli reads these arguments — no text parsing of step output strings is required.
3. **task-cli persists**: `task record` writes the record file and updates `index.json` using the status provided as a CLI argument.

This means the "Step 2/4: ... FAILED: [reason]" output is a human-readable log line. The structured data flows through the `task record` CLI arguments, not through parsing the agent's freeform output.

### Error Persistence

Errors are recorded in two places:

1. **Task record file** (markdown): the `status` frontmatter field is set to `failed`, `blocked`, or `rejected`. The `## Notes` section contains the failure reason. This is the primary audit trail.
2. **index.json**: the task's `status` field is updated via `task record`. No separate error log is created.

### Error Cases

| Error Case | Detection | Status Set | Error Code (Go) | Behavior |
|---|---|---|---|---|
| Task file unreadable | Step 1 (read failure) | `failed` | `NOT_FOUND` | Log error, skip Step 2 |
| Workflow heading empty | Step 2 Case C | (continues) | (warning only) | Log warning, read default template |
| No workflow heading | Step 2 Case B | (continues) | (no error) | Read default template (`task.md`) |
| Default template missing | Step 2 Case B (file not found) | `failed` | `NOT_FOUND` | Agent logs "default template not found" |
| Workflow execution fails | Step 2 (agent reports) | `failed` | (agent text) | Record failure reason in task record notes |
| Workflow declares failure instructions | Step 2 (workflow-specific) | `blocked` | (agent text) | Workflow instructs agent to create fix task; agent sets status to `blocked` |
| Agent timeout | External (session limit) | `in_progress` | (none) | Status stays `in_progress`, task is re-claimable |
| No test evidence at record | task-cli pre-check | (exit) | `VALIDATION_ERROR` | Print ErrNoTestEvidence, suggest `--force` |

## Integration Specs

No existing-page integrations — not applicable.

## Testing Strategy

### Coverage Target

- **task-cli Go code**: maintain >= 80% coverage on changed files (`claim.go`, `record.go`, `errors.go`, `types.go`). This is the same threshold enforced by `task-cli/CLAUDE.md` for all Go code. Verify with `go test -race -cover ./...`.
- **Agent prompt tests**: manual execution with binary pass/fail criteria (see scenarios below).
- **Grep verification**: binary pass/fail (zero matches = pass).

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test |
|---|---|---|---|
| task-cli (Go) | Unit | `go test` | NoTest removal: update/remove noTest test cases, remaining tests pass |
| task-cli (Go) | Integration | `go test` | `task record` quality gate runs for all tasks (no noTest bypass) |
| task-cli (Go) | Build | `go build ./...` | Compilation passes after field removal |
| Agent prompt | Manual | Execute T-test-3; test scenarios 4-8 | Step 2 output = "Execution Workflow", not "TDD implementation"; Case C (scenario 8) logs warning and falls back |
| Grep | Manual | `grep -ri noTest` | Zero matches across .go, .md, .json files |

### Key Test Scenarios

| # | Scenario | Pass Criteria | Fail Criteria |
|---|----------|---------------|---------------|
| 1 | `task record` with no test evidence | Exit code 1, stderr contains `VALIDATION_ERROR` | Exit code 0 or no error message |
| 2 | `task record --force` with no test evidence | Exit code 0, record file created, status = "completed" | Exit code non-zero or missing record file |
| 3 | `task record` with test evidence | Exit code 0, quality gate runs, record file created | Quality gate skipped or exit code non-zero |
| 4 | Task file without `## Execution Workflow` | Agent reads `templates/task.md`, follows TDD + QG steps | Agent reports "no workflow found" or skips execution |
| 5 | Task file with `## Execution Workflow` | Agent follows workflow steps exactly as declared | Agent deviates from declared workflow steps |
| 6 | T-test-3 execution | Completes in <5 min; on failure, creates fix tasks (no TDD retry loop) | Exceeds 5 min or retries TDD instead of creating fix tasks |
| 7 | `grep -ri noTest` across repo | Zero matches in `.go`, `.md`, `.json` files | Any match found |
| 8 | Task file with `## Execution Workflow` heading but empty content | Agent logs "WARNING: ... empty. Falling back to default template.", reads default template, follows TDD + QG steps | Agent skips execution, errors out, or uses workflow-less behavior |

## Security Considerations

### Threat Model

Internal forge tooling only. No external attack surface. Primary risk: agent deviates from workflow instructions, causing incorrect task execution.

### Mitigations

- `<EXTREMELY-IMPORTANT>` tags enforce workflow adherence (same mechanism as "ONE TASK PER INVOCATION" — zero violations since deployment)
- Default template fallback prevents silent failures on malformed task files
- task-cli quality gate pre-check remains as safety net for recording

## PRD Coverage Map

| PRD AC | Design Component | Interface/Model |
|---|---|---|
| Step 2 reads `## Execution Workflow` | task-executor.md Step 2 (Case A) | Interface 1 |
| No workflow → fallback TDD + QG | task-executor.md Step 2 (Case B) → default template | Interface 1 + 2 |
| Empty workflow → warning + fallback | task-executor.md Step 2 (Case C) | Interface 1 |
| All 16 templates have `## Execution Workflow` | Template updates | Interface 2 |
| noTest zero match (grep) | All file removals | Removed data model |
| task-cli compiles + tests pass | Go code removal + test updates | Interface 3 + 4 |
| run-tasks.md no NO_TEST | commands/run-tasks.md | Interface 3 |
| execute-task.md no NO_TEST | commands/execute-task.md | Interface 3 |
| T-test-3 output = "Execution Workflow" | Manual test scenario 6 | Interface 1 |
| record-task skill no noTest | skills/record-task/SKILL.md | Skill doc update |
| quick-tasks skill no --no-test | skills/quick-tasks/SKILL.md | Skill doc update |
| consolidate-specs skill no noTest | skills/consolidate-specs/SKILL.md | Skill doc update |
| index.schema.json no noTest field | Schema updates | Removed data model |
| File unreadable → failed | Error handling row 1 | Error handling |
| Agent timeout → failed | Error handling row 6 | Error handling |
| Workflow failure → failed + reason | Error handling row 5 | Error handling |
| Partial completion → failed + steps summary | Error handling row 5 | Error handling |

## Open Questions

None. All technical decisions resolved by proposal + PRD + this design.

## Appendix

### File Change Inventory

#### Agent Prompts (3 files)
| File | Change |
|---|---|
| `plugins/forge/agents/task-executor.md` | Rewrite Step 2 (skeleton), merge old Steps 2+3, remove Steps 4→3 and 5→4 renumbering, remove NO_TEST input |
| `plugins/forge/commands/run-tasks.md` | Remove NO_TEST from claim parsing and dispatch prompt |
| `plugins/forge/commands/execute-task.md` | Remove NO_TEST references, adopt same skeleton Step 2 |

#### Task Templates — Breakdown (10 files)
| File | Change |
|---|---|
| `templates/task.md` | Add `## Execution Workflow` (TDD + QG), remove `noTest: false` |
| `templates/gate-task.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/gen-test-cases.md` | Add `## Execution Workflow`, remove `noTest: true` |
| `templates/eval-test-cases.md` | Add `## Execution Workflow`, remove `noTest: true` |
| `templates/gen-test-scripts.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/run-e2e-tests.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/graduate-tests.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/verify-regression.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/consolidate-specs.md` | Add `## Execution Workflow`, remove `noTest: true` |
| `templates/phase-summary-task.md` | Add `## Execution Workflow`, remove `noTest: true` |

#### Task Templates — Quick (6 files)
| File | Change |
|---|---|
| `templates/task.md` | Add `## Execution Workflow` (TDD + QG), remove `noTest: false` |
| `templates/quick-test-cases.md` | Add `## Execution Workflow`, remove `noTest: true` |
| `templates/quick-gen-scripts.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/quick-run-tests.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/quick-graduate.md` | Add `## Execution Workflow`, remove `noTest: false` |
| `templates/quick-verify-regression.md` | Add `## Execution Workflow`, remove `noTest: false` |

#### Schemas (2 files)
| File | Change |
|---|---|
| `breakdown-tasks/templates/index.schema.json` | Remove `noTest` field definition |
| `quick-tasks/templates/index.schema.json` | Remove `noTest` field definition |

#### Task-CLI Go Code (4 files + tests)
| File | Change |
|---|---|
| `task-cli/pkg/task/types.go` | Remove `NoTest` from Task and TaskState structs |
| `task-cli/internal/cmd/claim.go` | Remove NoTest copy + PrintField("NO_TEST") |
| `task-cli/internal/cmd/record.go` | Remove coverage auto-set, QG skip, noTest formatting |
| `task-cli/internal/cmd/errors.go` | Update ErrNoTestEvidence message |
| `task-cli/internal/cmd/record_test.go` | Update/remove noTest test cases |
| `task-cli/internal/cmd/integration_test.go` | Update/remove noTest integration tests |

#### Skill Docs (3 files)
| File | Change |
|---|---|
| `plugins/forge/skills/record-task/SKILL.md` | Remove noTest references |
| `plugins/forge/skills/quick-tasks/SKILL.md` | Remove `--no-test` flag and noTest references |
| `plugins/forge/skills/consolidate-specs/SKILL.md` | Remove noTest reference |

### Total: ~29 files modified
