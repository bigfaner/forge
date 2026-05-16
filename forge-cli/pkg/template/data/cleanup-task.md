---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "P0"
estimated_time: "15min"
dependencies: []
status: pending
breaking: true
type: "cleanup"
---

# {{TITLE}}

## Root Cause

{{DESCRIPTION}}

## Reference Files

- Source: {{SOURCE_FILES}}
- Tool output: {{TEST_RESULTS}}

## Cleanup Guidelines

Fix only the reported style/lint issues. Do not refactor adjacent code.

1. Read the tool output and identify each violation
2. Fix each violation with minimal changes
3. Re-run the failing tool to confirm the fix

## Verification

After fixing, verify the cleanup works:
1. `just test [scope]` — must pass

E2e regression is verified by the dispatcher, not by this cleanup task.

When this task is recorded as completed via `task record`, the source task {{SOURCE_TASK_ID}} is automatically restored to pending if all its dependencies are completed.
