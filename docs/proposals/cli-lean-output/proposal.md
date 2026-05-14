---
created: 2026-05-14
author: "faner"
status: Draft
---

# Proposal: Remove Non-Essential Fields from Forge CLI Structured Output

## Problem

Forge CLI commands (claim, submit, query, status) output too many fields in their `---` blocks. Most are never consumed by downstream callers — they're informational noise that an LLM agent must parse and ignore, wasting context and increasing misparse risk.

### Evidence

`forge task claim` outputs 17 fields. Analysis shows:

| Field | Consumed by downstream? | Inferrable? |
|-------|------------------------|-------------|
| ACTION | Yes (run-tasks routing) | No |
| TASK_ID | Yes (everywhere) | No |
| BREAKING | Yes (quality gate) | From index |
| MAIN_SESSION | Yes (dispatch routing) | From index |
| SCOPE | Yes (just commands) | From index |
| FEATURE | Yes (E2E gate) | From state.json |
| FILE | Yes (agent reads task file) | Deterministic: `docs/features/<FEATURE>/tasks/<TASK_ID>.md` |
| **RECORD** | **No** — submit computes it internally | Deterministic path |
| **KEY** | **No** | From index |
| **TITLE** | **No** | From index |
| **PRIORITY** | **No** | From index |
| **STATUS** | **No** — always "in_progress" after claim | Trivial |
| **TYPE** | **No** — prompt.go computes independently | From index |
| **NO_TEST** | **No** — submit reads from index | From index |
| **ESTIMATED_TIME** | **No** | From index |
| **DEPENDENCIES** | **No** | From index |
| **PROFILE** | **No** — prompt.go reads from index | From index |

`forge task submit` outputs 3 fields:
- **TASK_ID** — caller already knows it
- **RECORD_FILE** — never consumed
- STATUS — essential (detects auto-downgrade to blocked)

11 of 17 claim fields and 2 of 3 submit fields are dead weight.

### Urgency

Every `forge task claim` invocation in run-tasks/execute-task forces the LLM to process 10 useless lines. In a 10-task feature, that's 100 lines of noise across the session. LLM context is expensive — trim it.

## Proposed Solution

**Remove all non-essential fields from structured output.** Only emit what downstream consumers actually use.

### claim output (before → after)

```
Before (17 fields):              After (7 fields):
---                              ---
ACTION: CLAIMED                  ACTION: CLAIMED
KEY: 1                           TASK_ID: 1
TASK_ID: 1                       FEATURE: cli-lean-output
TITLE: Some task title           FILE: Z:/project/.../tasks/1.md
PRIORITY: high                   SCOPE: backend
STATUS: in_progress              BREAKING: true
ESTIMATED_TIME: 30min            MAIN_SESSION: true
DEPENDENCIES: []                 ---
BREAKING: false
MAIN_SESSION: false
TYPE: default
SCOPE: backend
PROFILE: playwright
NO_TEST: false
FEATURE: cli-lean-output
FILE: Z:/project/.../tasks/1.md
RECORD: Z:/project/.../records/1.md
---
```

**Kept fields (consumed downstream):**
- `ACTION` — run-tasks routing logic
- `TASK_ID` — primary identifier
- `FEATURE` — E2E gate
- `FILE` — agent reads task markdown
- `SCOPE` — quality gate scope resolution (omit when empty)
- `BREAKING` — quality gate skip decision (omit when false)
- `MAIN_SESSION` — dispatch routing (omit when false)

### submit output (before → after)

```
Before (3 fields):     After (1 field):
---                    ---
TASK_ID: 1             STATUS: completed
RECORD_FILE: .../1.md  ---
STATUS: completed
---
```

### query output

Same trim: keep TASK_ID, STATUS, SCOPE, BREAKING. Drop KEY, TITLE, PRIORITY, ESTIMATED_TIME, DEPENDENCIES, FILE, RECORD.

## Requirements Analysis

### Key Scenarios

- `forge task claim` → agent receives minimal actionable fields
- `forge task submit` → agent only sees STATUS (the only thing it acts on)
- `forge task query` → agent only sees STATUS (the only thing it checks)

### Constraints & Dependencies

- Must keep all fields that run-tasks/execute-task parse from claim output
- Must keep STATUS from submit (record-task skill checks it)
- Boolean fields: omit when false (absence = false)

## Alternatives & Industry Benchmarking

| Approach | Pros | Cons | Verdict |
|----------|------|------|---------|
| Do nothing | No effort | 100+ lines of noise per feature run | Rejected |
| Add `--quiet` flag | Opt-in, backward compat | Downstream must remember flag; default still noisy | Rejected: wrong default |
| **Remove dead fields** | Clean default, less parsing risk, saves context | Slight breaking change for any consumer parsing TITLE etc. | **Selected: no known consumer parses these** |

## Scope

### In Scope

- `claim.go` — remove 10 non-essential fields from `printTaskDetails()`
- `submit.go` — remove TASK_ID and RECORD_FILE from output
- `query.go` — remove non-essential fields
- `status.go` — remove non-essential fields (if any)
- Boolean fields: BREAKING/MAIN_SESSION only print when true

### Out of Scope

- stderr cleanup (separate concern)
- Adding new fields
- JSON output mode
- Changing field names of kept fields

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Unknown consumer parses TITLE/KEY | L | M | grep confirms no consumer; easy to add back if needed |
| Agent needs RECORD path | L | L | forge task submit computes it internally; agent never reads record from claim |

## Success Criteria

- [ ] `forge task claim` outputs exactly 4-7 fields (ACTION + TASK_ID + FILE + FEATURE + conditional SCOPE/BREAKING/MAIN_SESSION)
- [ ] `forge task submit` outputs exactly 1 field (STATUS)
- [ ] `forge task query` outputs exactly 2-4 fields (TASK_ID + STATUS + conditional SCOPE/BREAKING)
- [ ] Boolean fields absent when false
- [ ] All existing unit tests pass

## Next Steps

- Proceed to `/quick-tasks` to generate implementation tasks
