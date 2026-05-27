---
id: "{{ID}}"
title: "{{TITLE}}"
priority: "P0"
estimated_time: "15min"
dependencies: []
status: pending
breaking: true
type: "coding.cleanup"
surface-key: ""
surface-type: ""
---

# {{TITLE}}

## Root Cause

{{DESCRIPTION}}

## Reference Files

- Source: {{SOURCE_FILES}}
- Tool output: {{TEST_RESULTS}}

## Surface Inference

This cleanup-task was created by the quality-gate hook. If `surface-key` and `surface-type` above are empty, infer them at execution time:

1. Parse `{{SOURCE_FILES}}` to extract the first file path (comma-separated).
2. Run `forge surfaces --json <file-path>` to resolve surface-key/type.
3. Use the resolved surface-type to load the appropriate `rules/surfaces/<type>.md` for test orchestration guidance.

If `forge surfaces --json` fails (no surfaces configured, command not found), proceed without surface information — this does not block the cleanup.

## Cleanup Guidelines

Fix only the reported style/lint issues. Do not refactor adjacent code.

1. Read the tool output and identify each violation
2. Fix each violation with minimal changes
3. Re-run the failing tool to confirm the fix

## Verification

After fixing, verify the cleanup works:
1. Run targeted tests on changed packages: `go test -race ./affected/package/...`
2. Replace the path with the actual packages you modified

> **Note:** Full project-wide tests run at CLI submit (`forge task submit`) — agent runs targeted tests only.

E2e regression is verified by the dispatcher, not by this cleanup task.

When this task is recorded as completed via `task record`, the source task {{SOURCE_TASK_ID}} is automatically restored to pending if all its dependencies are completed.
