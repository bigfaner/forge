---
created: 2026-05-14
author: faner
status: Draft
---

# Proposal: Task Record Write-Once Protection & Query Verbose Mode

## Problem

Task execution records (`records/<task>.md`) serve as audit logs documenting what was done, when, and with what results. Currently, `forge task submit` uses `os.WriteFile` with no existence check — it silently overwrites any existing record.

This was exploited unintentionally: after fix-1 completed for task 2, the agent re-ran `forge task submit` for task 2, replacing the original implementation record ("8 files created, 5 modified, 89.1% coverage") with a verification-only record ("no changes, 80.6% coverage"). The document-level truth was lost.

Additionally, the existing `forge task query` command (added as a stub during `cli-lean-output`) only displays 4 fields. Agents often need full task context — related fix tasks, file paths, priority, dependencies — to make informed decisions. Currently they must manually parse `index.json`.

## Proposed Solution

Two changes, delivered in a single commit:

### 1. Write-Once Protection in `forge task submit`

**Location**: `forge-cli/internal/cmd/submit.go`, before the `os.WriteFile` call (line 168).

**Behavior**:
- If record file exists and `--force` is not set → exit with error: "Record for task X already exists at <path>. Use --force to overwrite, or create a fix task instead."
- If record file exists and `--force` is set → overwrite with stderr warning: "WARNING: Overwriting existing record at <path>"
- If record file does not exist → write normally (current behavior)

**Rationale**: Reuses the existing `--force` flag (currently used to bypass quality gate). The "broken record needs rewrite" scenario is covered by `--force`. Normal workflow never needs to overwrite — fix tasks have independent records. Dual semantics is acceptable since `--force` universally means "skip safety checks".

### 2. `--verbose` Flag for `forge task query`

**Location**: Extend existing `forge-cli/internal/cmd/query.go` (currently a 4-field stub).

**Behavior**:
- **Default mode** (current behavior unchanged): display TASK_ID, STATUS, SCOPE (if set), BREAKING (if true)
- **Verbose mode** (`--verbose` / `-v`): display all task fields + RELATED_FIXES reverse lookup

**Verbose output format**:
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

**All verbose fields**:
| Field | Source |
|-------|--------|
| KEY | index.json map key |
| TASK_ID | task.ID |
| TITLE | task.Title |
| STATUS | task.Status |
| PRIORITY | task.Priority |
| TYPE | task.Type |
| SCOPE | task.Scope (omit if empty) |
| DEPENDENCIES | task.Dependencies (multi-line if multiple) |
| TASK_FILE | constructed from feature slug + task.File |
| RECORD_FILE | constructed from feature slug + task.Record |
| RELATED_FIXES | reverse lookup by sourceTaskID (omit if none) |

**Design decision**: Default mode respects the `cli-lean-output` principle (minimal noise for routine checks). Verbose mode provides full context for agents who need to understand task relationships and locate files.

## Requirements Analysis

### Key Scenarios

- Agent completes task, submits record → write succeeds (record does not exist)
- Agent accidentally re-submits same task → blocked by write-once protection with clear error message
- Developer needs to rewrite a broken record → `--force` allows overwrite with warning
- Agent checks task status quickly → `forge task query 2` shows essential fields
- Agent needs full context for decision-making → `forge task query 2 --verbose` shows all fields + related fixes
- Agent queries a task with no fixes → RELATED_FIXES field omitted in verbose output

### Non-Functional Requirements

- Write-once check: single `os.Stat` call — negligible performance impact
- RELATED_FIXES lookup: linear scan of index.json — fine for typical size (<200 tasks)

### Constraints & Dependencies

- `sourceTaskID` field already exists in `task.Task` struct (`pkg/task/types.go:99`)
- Fix task creation logic already populates `sourceTaskID` (`pkg/task/add.go`)
- `query.go` already registered in CLI (`root.go:53`)

## Alternatives & Industry Benchmarking

### Industry Solutions

Immutable audit logs are a standard pattern in CI/CD systems (GitHub Actions, Jenkins) where run records are never overwritten — new runs create new records.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | No effort | Record corruption risk remains; agents lack query depth | Rejected: known data loss vector |
| Full verbose by default | Original proposal | Simpler implementation | Contradicts cli-lean-output decision; noise for routine checks | Rejected: violates lean principle |
| **Lean default + --verbose** | This proposal | Respects lean output; on-demand detail | Two code paths in query | **Selected: best of both worlds** |

## Scope

### In Scope
- Write-once protection in `forge task submit` (existence check + error)
- `--force` flag reuse for intentional record overwrite with stderr warning
- `--verbose` / `-v` flag on `forge task query` for full field display
- RELATED_FIXES reverse lookup via `sourceTaskID` in verbose mode
- Tests for both changes (coverage ≥ 80%)

### Out of Scope
- Modifying existing `forge task status` command
- Modifying default query output (4-field lean stays as-is)
- Record-level append functionality (not needed — fix records are independent files)
- Version bump (minor: new query flag, new submit behavior)

## Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Agent scripts that re-submit existing records break | Low | Medium | Clear error message guides to --force or fix task |
| --force dual semantics confuses users | Low | Low | Help text documents both behaviors; error message mentions --force |
| Verbose output format becomes stale as Task struct evolves | Medium | Low | Use reflection or iterate struct tags to auto-generate fields |

## Success Criteria

- [ ] `forge task submit <id>` fails with descriptive error when record file already exists
- [ ] `forge task submit <id> --force` overwrites existing record with stderr warning
- [ ] `forge task query <id>` output unchanged (4 fields: TASK_ID, STATUS, SCOPE, BREAKING)
- [ ] `forge task query <id> --verbose` displays all task fields (KEY, TITLE, PRIORITY, TYPE, SCOPE, DEPENDENCIES, TASK_FILE, RECORD_FILE)
- [ ] `forge task query <id> --verbose` shows RELATED_FIXES for tasks with associated fix records
- [ ] `forge task query <id> --verbose` omits RELATED_FIXES when no fixes exist
- [ ] `forge task status <id>` behavior unchanged
- [ ] Test coverage ≥ 80% for new and modified code
