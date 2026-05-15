---
id: "1"
title: "Trim claim output to essential fields only"
priority: "P0"
estimated_time: "45m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 1: Trim claim output to essential fields only

## Description

Remove 10 non-essential fields from `printTaskDetails()` in `claim.go`. The claim command currently outputs 17 fields; only 7 are consumed by downstream callers (run-tasks dispatcher, execute-task skill). The remaining 10 (KEY, TITLE, PRIORITY, STATUS, ESTIMATED_TIME, DEPENDENCIES, TYPE, PROFILE, NO_TEST, RECORD) are dead weight that wastes LLM context.

Additionally, make boolean fields BREAKING and MAIN_SESSION conditional — only print when `true` (absence = false).

## Reference Files

- `docs/proposals/cli-lean-output/proposal.md` — Source proposal

## Affected Files

### Create
| File | Description |
|------|-------------|
| (none) |

### Modify
| File | Changes |
|------|---------|
| `forge-cli/internal/cmd/claim.go` | Rewrite `printTaskDetails()` to output only TASK_ID, FEATURE, FILE, SCOPE (conditional), BREAKING (conditional, only when true), MAIN_SESSION (conditional, only when true) |
| `forge-cli/internal/cmd/output_contract_test.go` | Update contract tests: remove assertions for KEY, TITLE, PRIORITY, STATUS, TYPE, NO_TEST, RECORD; update BREAKING/MAIN_SESSION to conditional expectations; update field-order assertions |

### Delete
| File | Reason |
|------|--------|
| (none) |

## Acceptance Criteria

- [ ] `printTaskDetails()` outputs exactly: TASK_ID, FEATURE, FILE, SCOPE (only when non-empty), BREAKING (only when true), MAIN_SESSION (only when true)
- [ ] `printNewTask()` still wraps with ACTION: CLAIMED + the trimmed fields
- [ ] `printContinueTask()` still wraps with ACTION: CONTINUE + trimmed fields + STARTED_AT
- [ ] Removed fields no longer appear: KEY, TITLE, PRIORITY, STATUS, ESTIMATED_TIME, DEPENDENCIES, TYPE, PROFILE, NO_TEST, RECORD
- [ ] Boolean fields (BREAKING, MAIN_SESSION) absent when false, present with "true" when true
- [ ] All existing unit tests pass after updates

## Hard Rules

- Do NOT change the `key` parameter signature of `printTaskDetails` — it is still used internally for routing. Only remove it from **output**.
- Do NOT remove or rename `PrintFieldIfNotEmpty` / `PrintFieldIfNotEmptySlice` helpers in `output.go` — they may be used elsewhere.

## Implementation Notes

**Target output for `printNewTask` (CLAIMED):**
```
---
ACTION: CLAIMED
TASK_ID: 1
FEATURE: cli-lean-output
FILE: Z:/project/.../tasks/1.md
SCOPE: backend
BREAKING: true
MAIN_SESSION: true
---
```
(Boolean fields omitted when false; SCOPE omitted when empty.)

**Tests requiring updates in `output_contract_test.go`:**
- `TestContract_Claim_NewTask` — remove mandatory checks for KEY, TITLE, PRIORITY, STATUS, BREAKING, MAIN_SESSION, TYPE, NO_TEST, RECORD; update to check only TASK_ID, FEATURE, FILE, SCOPE, and conditional BREAKING/MAIN_SESSION
- `TestContract_Claim_NewTask_ConditionalAbsent` — remove ESTIMATED_TIME/DEPENDENCIES/SCOPE/PROFILE absent checks; add BREAKING/MAIN_SESSION absent-when-false checks
- `TestContract_Claim_ProfilePresent` / `TestContract_Claim_ProfileAbsent` — DELETE these tests (PROFILE removed from output)
- `TestContract_Claim_FieldOrder` — update: ACTION first, then TASK_ID, FEATURE, FILE, SCOPE, BREAKING, MAIN_SESSION (no KEY or TYPE)
- `TestContract_Claim_Continue` — verify trimmed fields still present with ACTION: CONTINUE + STARTED_AT

**Tests requiring updates in `claim_test.go`:**
- `TestPrintTaskDetails_BreakingInOutput` — update "breaking false" subtest: BREAKING should be absent (not "BREAKING: false")
- `TestPrintTaskDetails_ScopeInOutput` — keep as-is (SCOPE conditional logic unchanged)
- `TestPrintTaskDetails_TypeInOutput` — DELETE (TYPE removed from output)
- `TestPrintTaskDetails_ProfileInOutput` — DELETE (PROFILE removed from output)

**Version bump**: Patch bump in `scripts/version.txt` (dead code removal / output simplification).
