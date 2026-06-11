---
id: "2"
title: "Query verbose mode with RELATED_FIXES"
priority: "P1"
estimated_time: "1.5h"
dependencies: []
scope: "backend"
breaking: false
type: "feature"
mainSession: false
---

# 2: Query verbose mode with RELATED_FIXES

## Description
Extend the existing `forge task query` stub with a `--verbose` / `-v` flag. Default output stays lean (TASK_ID, STATUS, SCOPE, BREAKING). Verbose mode adds all remaining fields (KEY, TITLE, PRIORITY, TYPE, DEPENDENCIES, TASK_FILE, RECORD_FILE) plus a RELATED_FIXES reverse lookup that finds all fix tasks spawned from the queried task.

The `sourceTaskID` field already exists in the Task model and is populated when fix tasks are created. The reverse lookup simply iterates all tasks in the index collecting matches.

## Reference Files
- `docs/proposals/task-record-immutability/proposal.md` — Source proposal
- `forge-cli/internal/cmd/query.go` — Primary implementation target
- `forge-cli/pkg/task/types.go` — Task struct with SourceTaskID field (line 99)
- `forge-cli/pkg/feature/` — Feature slug utilities for constructing file paths

## Acceptance Criteria
- [ ] `forge task query <id>` output unchanged (TASK_ID, STATUS, SCOPE if set, BREAKING if true)
- [ ] `forge task query <id> --verbose` displays: KEY, TASK_ID, TITLE, STATUS, PRIORITY, TYPE, SCOPE (if set), DEPENDENCIES (multi-line if multiple), TASK_FILE, RECORD_FILE
- [ ] `forge task query <id> --verbose` shows RELATED_FIXES when fix tasks exist: `<id> [<status>] <title>` per line
- [ ] `forge task query <id> --verbose` omits RELATED_FIXES when no fixes exist
- [ ] `forge task query <id> -v` works as shorthand for `--verbose`
- [ ] TASK_FILE and RECORD_FILE paths are constructed from feature slug + task File/Record fields

## Hard Rules
- Default mode output MUST NOT change — zero regression risk for existing consumers
- Use existing `PrintField` / `PrintFieldIfNotEmpty` / `PrintBlockStart` / `PrintBlockEnd` helpers
- RELATED_FIXES uses `task.SourceTaskID` field (already exists in data model)

## Implementation Notes
- Add `queryVerbose` bool flag to `queryCmd` via `Flags().BoolVarP`
- Extract the verbose block into a separate function to keep `runQuery` clean
- For KEY: use the map key from `task.FindTask` (first return value)
- For TASK_FILE/RECORD_FILE: use `feature.GetTaskFile(featureSlug, t.File)` pattern already used in submit.go
- For RELATED_FIXES: iterate `index.TasksMap()` or equivalent, collect tasks where `SourceTaskID == t.ID`
