---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "P0"
estimated_time: "30min"
dependencies: []
status: pending
breaking: true
---

# {{TITLE}}

## Root Cause

{{DESCRIPTION}}

## Reference Files

- Source: {{SOURCE_FILES}}
- Test script: {{TEST_SCRIPT}}
- Test results: {{TEST_RESULTS}}

## Verification

After fixing, verify the fix works:
1. `just test [scope]` — must pass
2. If UI/page related: `just test-e2e --feature <slug>` — must also pass

When this task is recorded as completed via `task record`, the source task {{SOURCE_TASK_ID}} is automatically restored to pending if all its dependencies are completed.
