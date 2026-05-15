---
id: "2"
title: "Trim submit, query, and status output to essential fields"
priority: "P0"
estimated_time: "30m"
dependencies: []
scope: "backend"
breaking: false
type: "implementation"
mainSession: false
---

# 2: Trim submit, query, and status output to essential fields

## Description

Apply the same lean-output principle to `submit.go`, `query.go`, and `status.go`. These commands output fields that downstream consumers never act on. Trim each to only essential fields.

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
| `forge-cli/internal/cmd/submit.go` | Remove TASK_ID and RECORD_FILE from non-JSON output block; keep STATUS only |
| `forge-cli/internal/cmd/query.go` | Remove KEY, TITLE, PRIORITY, ESTIMATED_TIME, DEPENDENCIES, FILE, RECORD; keep TASK_ID, STATUS; add SCOPE (conditional); BREAKING already conditional |
| `forge-cli/internal/cmd/status.go` | Remove KEY, TITLE, DEPENDENCIES from both query-mode and update-mode output blocks; keep TASK_ID, STATUS |

### Delete
| File | Reason |
|------|--------|
| (none) |

## Acceptance Criteria

- [ ] `forge task submit` (non-JSON, non-quiet) outputs exactly 1 field: STATUS
- [ ] `forge task query` outputs exactly TASK_ID + STATUS + SCOPE (when non-empty) + BREAKING (when true)
- [ ] `forge task status` (query mode) outputs exactly TASK_ID + STATUS
- [ ] `forge task status` (update mode) outputs exactly TASK_ID + STATUS
- [ ] `forge task status` (unmet deps warning) outputs TASK_ID + STATUS + WARNING line
- [ ] JSON mode (`--json`) in submit is NOT changed
- [ ] All existing unit tests pass after updates

## Hard Rules

- Do NOT change JSON output mode (`--json` flag) in submit.go — that is a separate interface.
- Do NOT change the quiet mode (`--quiet` flag) in submit.go.

## Implementation Notes

**submit.go target output (non-JSON, non-quiet):**
```
---
STATUS: completed
---
```

**query.go target output:**
```
---
TASK_ID: 1
STATUS: in_progress
SCOPE: backend
BREAKING: true
---
```
(SCOPE omitted when empty; BREAKING omitted when false. Note: query currently doesn't output SCOPE — need to add it from the task struct.)

**status.go target output (both modes):**
```
---
TASK_ID: 1
STATUS: in_progress
---
```
Plus WARNING line before closing `---` for unmet-deps case.

**Tests requiring updates in `output_contract_test.go`:**
- `TestContract_Record_Completed` — update to check only STATUS field (remove TASK_ID and RECORD_FILE assertions)
- `TestContract_Record_Blocked` — update to check only STATUS field

**Version bump**: Patch bump in `scripts/version.txt` (dead code removal / output simplification). Coordinate with Task 1 to avoid merge conflict — both tasks bump the same file. Whichever task runs second should handle the version bump; if Task 1 already bumped, just ensure the final version is correct.
