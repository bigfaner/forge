---
created: "2026-05-17"
tags: [architecture]
---

# Post-Completion Status Transition and Commit via Stop Hook

## Problem

After all tasks complete in `/quick` mode, the dispatcher updates `manifest.md` and `proposal.md` status, but these changes are never committed. The user finds uncommitted status files after the pipeline finishes.

## Root Cause

Causal chain:
1. **Symptom**: manifest.md and proposal.md status updates are uncommitted after pipeline completion
2. **Direct cause**: The `/run-tasks` Post-Completion step only prints a summary and optionally pushes — no commit step
3. **Root cause**: Gap between `/quick` and `/run-tasks` responsibilities. Neither skill owns the post-completion status transition + commit
4. **Trigger condition**: Any `/quick` execution that completes all tasks successfully

## Solution

Use Claude Code's **Stop hook** for the post-completion status transition. This aligns with the existing quality gate mechanism and naturally handles the "fix tasks added by quality gate" case.

### Stop Hook Mechanism

Claude Code fires a `Stop` event when the agent finishes responding. Stop hooks can return a decision:

| Return value | Effect |
|-------------|--------|
| Normal exit (no JSON) | Agent stops normally |
| `{ decision: "block", reason: "..." }` | Agent continues working |

This creates a natural loop for quality gate → fix tasks:

```
Agent finishes → Stop hook #1 (quality gate)
  → Issues found: block → agent continues (fix tasks exist)
  → No issues: allow → agent stops → Stop hook #2 (status transition)
```

### `stop_hook_active` Prevents Infinite Loops

When a Stop hook returns `block`, the next Stop event includes `stop_hook_active: true`. The hook script can check this to enforce a max retry count.

### Proposed Hook Configuration

```json
{
  "hooks": {
    "Stop": [
      {
        "hooks": [
          { "type": "command", "command": "forge quality-gate" },
          { "type": "command", "command": "forge feature complete-if-done" }
        ]
      }
    ]
  }
}
```

- Hook #1 (quality gate): checks for issues, returns `block` if fix tasks added
- Hook #2 (status transition): runs only when hook #1 passes and agent truly stops. Updates proposal.md + manifest.md status and commits

### Open Question

Whether multiple Stop hooks execute sequentially (hook #2 always runs after #1) or whether hook #1 returning `block` prevents hook #2 from executing. **This needs testing** — the official docs describe single-hook behavior but don't explicitly state the multi-hook execution order when the first hook returns `block`.

## Reusable Pattern

**Rule**: Feature lifecycle status transitions should be handled by Stop hooks, not by /quick or /run-tasks skill code. This ensures:
- Status only changes after quality gate passes
- Fix tasks added by quality gate prevent premature completion
- The mechanism works for both /quick and full pipeline

**Prerequisite**: Implement `forge feature complete-if-done` CLI command that:
1. Checks all tasks are completed (no blocked/in_progress)
2. Checks `stop_hook_active` is false (not a retry loop)
3. Updates proposal.md + manifest.md status
4. Commits both files

## Related Files

- `docs/official-references/hooks.md` — Claude Code hooks reference
- `plugins/forge/skills/run-tasks/SKILL.md` — dispatcher protocol
- `plugins/forge/skills/quick/SKILL.md` — quick mode pipeline
