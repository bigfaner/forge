---
id: "fix-3"
title: "Fix: lint failure in all-completed quality gate"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: lint failure in all-completed quality gate

## Root Cause

Quality gate step `just lint` failed during all-completed hook.

Error output saved to: `tests/results/unit-raw-output.txt`

Concise error:
```
...
^
99 issues:
* errcheck: 50
* gocritic: 11
* ineffassign: 1
* revive: 29
* staticcheck: 5
* unparam: 2
* whitespace: 1
error: Recipe `lint` failed with exit code 1
```

## Reference Files

- Source: add_cmd_test.go, all_completed.go, all_completed_test.go, root_test.go, add_test.go, state_test.go, check_test.go, forensic.go, integration_test.go, validate.go
- Test script: just lint
- Test results: tests/results/unit-raw-output.txt

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
