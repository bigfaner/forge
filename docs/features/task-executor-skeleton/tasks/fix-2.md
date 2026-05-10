---
id: "fix-2"
title: "Fix: unit-test failure in all-completed quality gate"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: unit-test failure in all-completed quality gate

## Root Cause

Quality gate step `just test` failed during all-completed hook.

Error output saved to: `tests/results/unit-raw-output.txt`

Concise error:
```
go: -race requires cgo; enable cgo by setting CGO_ENABLED=1
error: Recipe `test` failed with exit code 2

```

## Reference Files

- Source: See error output for affected files
- Test script: just test
- Test results: tests/results/unit-raw-output.txt

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
