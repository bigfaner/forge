---
id: "fix-5"
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
...
--- FAIL: TestFindRootInfoFrom/finds_root_from_subdirectory (0.00s)
root_test.go:205: Path = "Z:\\project\\ai\\forge", want "C:\\Users\\panda\\AppData\\Local\\Temp\\TestFindRootInfoFromfinds_root_from_subdirectory3346198918\\001"
FAIL
FAIL	task-cli/pkg/project	0.494s
ok  	task-cli/pkg/task	(cached)
ok  	task-cli/pkg/template	(cached)
ok  	task-cli/pkg/testrunner	(cached)
ok  	task-cli/pkg/version	(cached)
FAIL
error: Recipe `test` failed with exit code 1
```

## Reference Files

- Source: claim_test.go, integration_test.go, root_test.go
- Test script: just test
- Test results: tests/results/unit-raw-output.txt

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
