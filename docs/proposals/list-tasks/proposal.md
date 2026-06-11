---
created: "2026-05-23"
status: Draft
intent: "Add forge task list command to show all tasks for current feature"
---

# Proposal: forge task list Command

## Problem

There is no CLI command to list individual tasks for a feature. Current commands only show:
- `forge feature status <slug>` — aggregate counts (e.g. `pending: 3, completed: 5`)
- `forge task status <id>` — single task status
- `forge task query <id>` — single task details

Users must open `index.json` or `manifest.md` manually to see all task IDs, titles, and statuses at a glance. This is especially inconvenient when using `forge task claim` or `forge task transition` which require knowing task IDs.

## Solution

Add `forge task list` subcommand that displays all tasks for the current feature in a table:

```
── TASKS ────────────────────────────────────────
  8 found  (feature: feature-set-command)

  ID     TYPE              TITLE                         STATUS
  -----  ----------------  ----------------------------  --------
  1      coding.feature    Add feature set subcommand     completed
  2      coding.feature    Add feature list subcommand    completed
  3      coding.test       Unit tests for feature         in_progress
  T-1    gate              Compile check                  pending
──────────────────────────────────────────────────────────
```

### Behavior

- Resolves current feature via `feature.RequireFeature()` (same as other task commands)
- Loads `index.json`, iterates all tasks in the map
- Sorts by ID (natural order: `1, 2, 3, T-1, T-2`)
- Displays: ID, Type, Title (truncated if too long), Status
- Header shows total count and feature slug

### Options

- No extra flags needed for initial version. Can add `--status <status>` filter later.

## Alternatives

| Alternative | Trade-off |
|-------------|-----------|
| **Do nothing** — users read index.json | No structured CLI view, poor DX |
| **Extend forge feature status** — add task table there | Feature status is aggregate-focused; mixing concerns |
| **forge feature tasks** — new feature subcommand | Breaks the convention that task operations live under `forge task` |

## Scope

### In Scope
- `internal/cmd/task/list.go` — new subcommand
- Registration in `task/register.go`
- Natural-sort task IDs

### Out of Scope
- Status filtering (`--status` flag)
- Verbose mode (task file paths, dependencies)
- Cross-feature task listing

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| Large task lists (>20) produce verbose output | Low | Title column truncation, sorted by ID keeps it scannable |
| Feature with no index.json | Low | Print "no tasks found" matching existing patterns |

## Success Criteria

- [ ] `forge task list` shows all tasks for current feature in a table
- [ ] Tasks are sorted by ID in natural order
- [ ] Output includes ID, Type, Title (truncated), Status columns
- [ ] Graceful handling of features with no tasks
- [ ] Unit tests cover: normal list, empty feature, sorted output
- [ ] Version bump in `scripts/version.txt` (minor: new command)
