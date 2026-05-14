---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: Task Record Write-Once Protection & Task Query Command

## Problem

Task execution records (`records/<task>.md`) serve as audit logs documenting what was done, when, and with what results. Currently, `forge task submit` uses `os.WriteFile` with no existence check — it silently overwrites any existing record.

This was exploited unintentionally: after fix-1 completed for task 2, the agent re-ran `forge task submit` for task 2, replacing the original implementation record ("8 files created, 5 modified, 89.1% coverage") with a verification-only record ("no changes, 80.6% coverage"). The document-level truth was lost.

Additionally, there is no CLI command to inspect task details including related fix tasks. The existing `forge task status <id>` shows only KEY/TASK_ID/STATUS/TITLE/DEPENDENCIES — insufficient for agents who need to understand task context, locate record files, or discover related fix tasks.

## Proposed Solution

Two changes, delivered in a single commit:

### 1. Write-Once Protection in `forge task submit`

**Location**: `forge-cli/internal/cmd/submit.go`, before the `os.WriteFile` call (line 168).

**Behavior**:
- If record file exists and `--force` is not set → exit with error: "Record for task X already exists at <path>. Use --force to overwrite, or create a fix task instead."
- If record file exists and `--force` is set → overwrite with stderr warning: "WARNING: Overwriting existing record at <path>"
- If record file does not exist → write normally (current behavior)

**Rationale**: Reuses the existing `--force` flag (currently used to bypass quality gate). The "broken record needs rewrite" scenario is covered by `--force`. Normal workflow never needs to overwrite — fix tasks have independent records.

### 2. `forge task query <id>` Command

**New file**: `forge-cli/internal/cmd/query.go` + `query_test.go`

**Behavior**: Display all task metadata from index.json, plus a reverse-lookup of related fix tasks.

**Output format** (using existing PrintField/PrintBlock style):

```
>>>
KEY             2-info-commands
TASK_ID         2
TITLE           Info commands (proposal, feature, lesson)
STATUS          completed
PRIORITY        P1
TYPE            implementation
SCOPE           backend
DEPENDENCIES    1
TASK_FILE       docs/features/forge-info-commands/tasks/2-info-commands.md
RECORD_FILE     docs/features/forge-info-commands/tasks/records/2-info-commands.md
RELATED_FIXES   fix-1 [completed] Fix: compilation errors in task 2
<<<
```

**RELATED_FIXES logic**: Iterate all tasks in index.json, collect entries where `sourceTaskID` matches the queried task's ID. Display each as `<id> [<status>] <title>`, one per line. If no related fixes exist, omit the field entirely.

**All fields displayed**:
| Field | Source |
|-------|--------|
| KEY | index.json map key |
| TASK_ID | task.ID |
| TITLE | task.Title |
| STATUS | task.Status |
| PRIORITY | task.Priority |
| TYPE | task.Type |
| SCOPE | task.Scope |
| DEPENDENCIES | task.Dependencies (multi-line if multiple) |
| TASK_FILE | constructed from feature slug + task.File |
| RECORD_FILE | constructed from feature slug + task.Record |
| RELATED_FIXES | reverse lookup by sourceTaskID (omit if none) |

**Relationship with `forge task status`**: `status` retains its dual-purpose behavior (query + update). `query` is a read-only, information-rich alternative. No breaking changes.

## Scope

### In Scope
- Write-once protection in `forge task submit` (overwrite detection + error)
- `--force` flag reuse for intentional overwrite
- New `forge task query <id>` subcommand
- Reverse fix-lookup via `sourceTaskID`
- Tests for both changes (coverage ≥ 80%)

### Out of Scope
- Modifying existing `forge task status` command
- Record-level append functionality (not needed — fix records are independent files)
- Lesson document updates

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent scripts that re-submit existing records break | Low | Medium | Clear error message guides to --force or fix task |
| RELATED_FIXES slow on large index | Low | Low | index.json is small (<200 tasks typically), linear scan is fine |

## Success Criteria

- [ ] `forge task submit <id>` fails with descriptive error when record file already exists
- [ ] `forge task submit <id> --force` overwrites existing record with stderr warning
- [ ] `forge task query <id>` displays all task fields
- [ ] `forge task query <id>` shows RELATED_FIXES for tasks with associated fix records
- [ ] `forge task query <id>` omits RELATED_FIXES when no fixes exist
- [ ] `forge task status <id>` query behavior unchanged
- [ ] Test coverage ≥ 80% for new and modified code

## Next Steps

1. Run `/eval-proposal` to evaluate
2. Implement directly (small scope — no PRD/design needed)
