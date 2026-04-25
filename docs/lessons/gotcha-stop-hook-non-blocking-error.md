# Stop Hook: Failed with Non-Blocking Status Code

## Problem

Claude Code reports:
```
Stop hook error: Failed with non-blocking status code: No stderr output.
```

The hook failure is silent — no diagnostic info, no indication of what went wrong.

## Root Cause

A Stop hook (configured in `.claude/settings.json` under `hooks.Stop`) exited with a non-zero status code but wrote nothing to stderr. Claude Code surfaces the failure but has no message to show.

Non-blocking hooks don't halt Claude's execution, but they do emit this warning when they fail silently.

## Solution

Two options:

**Option A — Fix the hook script** so it exits 0 on success and writes to stderr on failure:
```bash
#!/bin/bash
some-command || { echo "some-command failed: $?" >&2; exit 1; }
```

**Option B — Remove or disable the hook** if it's no longer needed:
```json
// .claude/settings.json
{
  "hooks": {
    "Stop": []
  }
}
```

## Key Takeaway

Always write to stderr in hook scripts when failing. Claude Code only shows stderr output in hook error messages — a silent non-zero exit produces the unhelpful "No stderr output" warning. Add `>&2` redirects to any diagnostic `echo` in hook scripts.
