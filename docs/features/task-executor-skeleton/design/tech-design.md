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

**Rules for workflow content** (enforced in template review):
- Must include explicit stop/termination condition
- Must not use open-ended instructions ("continue until...", "keep trying...")
- Must specify commands to run and success/failure criteria
- Must specify failure handling (create fix task? halt? record partial?)

### Interface 3: task-cli Claim Output

```
BEFORE: KEY, TASK_ID, FILE, BREAKING, MAIN_SESSION, NO_TEST, SCOPE, FEATURE
AFTER:  KEY, TASK_ID, FILE, BREAKING, MAIN_SESSION, SCOPE, FEATURE
```

`NO_TEST` field removed. Downstream consumers (run-tasks.md, execute-task.md) stop parsing it.

### Interface 4: task-cli Record Behavior

```
BEFORE:
  noTest → coverage = -1.0 (auto), skip quality gate pre-check
  else   → require test evidence, run quality gate

AFTER:
  All tasks: require test evidence OR explicit --force
  Quality gate pre-check runs for all completed tasks
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

| Error Case | Detection | Behavior |
|---|---|---|
| Task file unreadable | Step 1 (read failure) | Status = `failed`, log error, skip Step 2 |
| Workflow heading empty | Step 2 Case C | Log warning, read default template |
| No workflow heading | Step 2 Case B | Read default template (`task.md`) |
| Default template missing | Step 2 Case B (file not found) | Status = `failed`, log "default template not found" |
| Workflow execution fails | Step 2 (agent reports) | Status = `failed` in Step 3, record failure reason |
| Agent timeout | External (session limit) | Status stays `in_progress`, re-claimable |
| No test evidence at record | task-cli pre-check | Error, suggest `--force` |

## Integration Specs

No existing-page integrations — not applicable.

## Testing Strategy

### Per-Layer Test Plan

| Layer | Test Type | Tool | What to Test |
|---|---|---|---|
| task-cli (Go) | Unit | go test | NoTest removal: update/remove noTest test cases, remaining tests pass |
| task-cli (Go) | Integration | go test | `task record` quality gate runs for all tasks (no noTest bypass) |
| task-cli (Go) | Build | go build | Compilation passes after field removal |
| Agent prompt | Manual | Execute T-test-3 | Step 2 output = "Execution Workflow", not "TDD implementation" |
| Grep | Manual | grep -ri noTest | Zero matches across .go, .md, .json files |

### Key Test Scenarios

1. `task record` with no test evidence → error (previously auto-passed for noTest tasks)
2. `task record --force` with no test evidence → success
3. `task record` with test evidence → success, quality gate runs
4. Task file without `## Execution Workflow` → agent reads default template → follows TDD
5. Task file with `## Execution Workflow` → agent follows workflow steps
6. T-test-3 execution → <5 min, creates fix tasks on failure (no TDD retry)

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
