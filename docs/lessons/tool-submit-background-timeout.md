---
created: "2026-05-20"
tags: [testing, local-dev-deployment]
---

# forge task submit with quality gate triggers background execution

## Problem

When running `forge task submit` for `breaking: true` tasks, the CLI runs the full quality gate (compile→fmt→lint→test). On this codebase, `just test` alone takes ~90 seconds, and the total gate can approach or exceed the Bash tool's default 120-second timeout. The system automatically moves the command to background execution, requiring `TaskOutput` to check the result — adding an extra round-trip and delaying feedback.

## Root Cause

1. **Immediate cause**: `forge task submit` ran as a background task, requiring `TaskOutput` to retrieve the result.
2. **Proximate cause**: The full quality gate (compile→fmt→lint→test) for a project with ~90s test suite approaches the Bash tool's 120s default timeout, triggering automatic background execution.
3. **Systemic cause**: The Bash tool has a fixed default timeout (120s) and no automatic retry with longer timeout. Long-running quality gates always cross this boundary.

## Solution

When calling `forge task submit` for tasks that will trigger the quality gate (breaking tasks, coding.* types), set an explicit timeout:

```
Bash(command="forge task submit <ID> --data record.json", timeout=300000)
```

For doc-type tasks or when using `--force`, the default timeout is sufficient.

## Reusable Pattern

Any CLI command that runs `just test` as part of its workflow needs `timeout: 300000` (5 minutes) in the Bash tool call. This includes:
- `forge task submit` (breaking tasks)
- `just test` (direct invocation)
- `forge quality-gate`

## Example

```bash
# Wrong — may auto-background
forge task submit 3 --data record.json

# Right — explicit timeout for quality gate
forge task submit 3 --data record.json  # with timeout=300000 in Bash tool
```

## Related Files

- `forge-cli/internal/cmd/submit.go` — `validateQualityGate()` runs the gate
