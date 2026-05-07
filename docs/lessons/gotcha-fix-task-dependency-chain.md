---
created: "2026-05-07"
tags: [dependencies, testing]
---

# task add --source-task-id Uses ID to Lookup, But Index Stores Tasks by Key

## Problem

When the dispatcher added a fix task via `task add --template fix-task --source-task-id T-test-3`, the fix task was created with `sourceTaskID: "T-test-3"` but the reverse dependency (disc-1 appended to T-test-3's dependencies) was **silently skipped**. The source task `run-e2e-tests` had `id: "T-test-3"` but was stored under key `"run-e2e-tests"` in `index.json`. The lookup `index.Tasks["T-test-3"]` returned not found.

This caused:
1. disc-1 was added with `sourceTaskID` annotation but no dependency edge back to the source
2. T-test-3 was forced to `completed` to unblock the pipeline (bypassing quality gate)
3. Downstream tasks (T-test-4, T-test-4.5) could be claimed before disc-1 completed
4. Only disc-1's P0 priority prevented a wrong execution order

## Root Cause

**Causal chain (4 levels deep):**

1. **Symptom**: disc-1 was not added as a dependency of T-test-3 in `index.json`
2. **Direct cause**: The source task lookup at `add.go:120` uses `opts.SourceTaskID` (the task ID, e.g. `"T-test-3"`) as the map key, but the index stores tasks under their **slug key** (e.g. `"run-e2e-tests"`)
3. **Root cause**: Two different lookup patterns exist in `add.go`:
   - **Dependency validation** (line 83): iterates all tasks, compares `t.ID == dep` — works correctly with IDs
   - **SourceTaskID lookup** (line 120): uses `index.Tasks[opts.SourceTaskID]` — direct map access by ID, fails when key != ID
4. **Trigger condition**: Any task whose key differs from its ID (all tasks defined in `index.json` with a slug key like `"run-e2e-tests"` that maps to id `"T-test-3"`)

**The code at `add.go:118-127`:**
```go
if opts.SourceTaskID != "" {
    srcTask, found := index.Tasks[opts.SourceTaskID]  // BUG: uses ID as map key
    if found {
        if !slices.Contains(srcTask.Dependencies, opts.ID) {
            srcTask.Dependencies = append(srcTask.Dependencies, opts.ID)
            index.Tasks[opts.SourceTaskID] = srcTask
        }
    }
    // Silent no-op when found == false
}
```

**The fix**: lookup should iterate by ID (like the dependency validation does), not by direct map key:
```go
if opts.SourceTaskID != "" {
    for key, t := range index.Tasks {
        if t.ID == opts.SourceTaskID {
            if !slices.Contains(t.Dependencies, opts.ID) {
                t.Dependencies = append(t.Dependencies, opts.ID)
                index.Tasks[key] = t
            }
            break
        }
    }
}
```

## Solution

**Task-cli bug**: `pkg/task/add.go` line 120 should iterate over `index.Tasks` matching by `t.ID == opts.SourceTaskID` instead of using the ID as a direct map key. The same pattern is already used correctly for dependency validation at lines 82-87.

**Workaround for dispatcher**: After `task add --source-task-id`, verify the dependency was actually added by reading `index.json`. If not, manually edit the source task's dependencies.

## Reusable Pattern

**Map key != ID is a recurring trap in task-cli.** When:
- Tasks are defined in `index.json` with slug keys (e.g. `"run-e2e-tests"`) and separate `id` fields (e.g. `"T-test-3"`)
- Dynamically added tasks (disc-1) use their ID as both key and ID
- Any code that does `index.Tasks[someID]` will fail silently for static tasks

**Rule**: Always iterate `for key, t := range index.Tasks { if t.ID == targetID ... }` when looking up by task ID. Only use direct map access when you have the actual key.

## Example

```
# index.json structure — key ≠ id for static tasks
"run-e2e-tests": { "id": "T-test-3", ... }   # key=run-e2e-tests, id=T-test-3
"disc-1":         { "id": "disc-1", ... }      # key=disc-1, id=disc-1 (same)

# task add --source-task-id T-test-3 tries:
index.Tasks["T-test-3"]  → not found (key is "run-e2e-tests")

# Correct lookup:
for key, t := range index.Tasks {
    if t.ID == "T-test-3" {  // found at key "run-e2e-tests"
```

## Related Files

- `C:\Users\panda\.claude\plugins\marketplaces\forge\task-cli\pkg\task\add.go` — bug at line 120
- `C:\Users\panda\.claude\plugins\marketplaces\forge\task-cli\pkg\task\types.go` — Task struct with SourceTaskID field
- `docs/features/agent-task-dashboard/tasks/index.json` — task index showing key vs ID mismatch

## References

- task-cli v1.4.0
