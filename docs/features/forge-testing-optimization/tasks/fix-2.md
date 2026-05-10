---
id: "fix-2"
title: "Fix: compile failure in all-completed quality gate"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# Fix: compile failure in all-completed quality gate

## Root Cause

Quality gate step `just compile` failed during all-completed hook.

Error output saved to: `tests/results/unit-raw-output.txt`

Concise error:
```

[41m                                                                               [0m
[41m[37m                This is not the tsc command you are looking for                [0m
[41m                                                                               [0m

To get access to the TypeScript compiler, [34mtsc[0m, from the command line either:

- Use [1mnpm install typescript[0m to first add TypeScript to your project [1mbefore[0m using npx
- Use [1myarn[0m to avoid accidentally running code from un-installed packages
error: Recipe `compile` failed with exit code 1

```

## Reference Files

- Source: See error output for affected files
- Test script: just compile
- Test results: tests/results/unit-raw-output.txt

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task N/A (project-wide gate) is automatically restored to pending if all its dependencies are completed.
