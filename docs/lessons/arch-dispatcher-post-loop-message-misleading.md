---
created: "2026-05-26"
tags: [architecture, testing]
---

# Dispatcher post-loop message misleads when test tasks already ran

## Problem

After `/run-tasks` completes, the dispatcher prints a fixed template message: "T-test-run and T-test-verify-regression handle e2e verification and regression automatically." This message implies these tasks still need to run, but they may have already been executed during the dispatcher loop itself.

## Root Cause

1. **Surface**: The post-completion section of `/run-tasks` SKILL.md contains a hardcoded message intended for cases where T-test-run and T-test-verify-regression are NOT in the task list (e.g., feature missing these auto-generated tasks).
2. **Mechanism**: The dispatcher loop executes ALL pending tasks including test pipeline tasks. When the loop exits with "no pending tasks available", the agent blindly appends the template message regardless of whether those tasks already ran.
3. **Structural**: The skill template mixes two concerns: (a) informing the user about test task purpose, and (b) suggesting next steps. The message should be conditional on whether those tasks actually ran.

## Solution

Before printing the post-loop summary, check which tasks were actually completed in the loop. If T-test-run and T-test-verify-regression already have status "completed" in index.json, skip or adapt the message:

- If both completed: "All tasks completed including e2e tests (T-test-run) and regression verification (T-test-verify-regression)."
- If neither exists in index.json: print the original suggestion message.
- If only one completed: mention the remaining one.

## Reusable Pattern

Post-loop summary messages must reflect actual execution state, not template defaults. After a loop that dynamically executes all available tasks, the completion message should be derived from index.json status, not hardcoded text.

## Related Files

- plugins/forge/skills/run-tasks/SKILL.md (Post-Completion section)
