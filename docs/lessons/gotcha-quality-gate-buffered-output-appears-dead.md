---
created: "2026-05-24"
tags: [testing, local-dev-deployment]
---

# forge task submit Appears to Hang Due to Buffered Quality Gate Output

## Problem

`forge task submit` for a `breaking: true` task appears to hang/deadlock. No output is shown for minutes, leading the user to kill the process (exit code 137) and suspect a deadlock.

## Root Cause

1. **Surface**: `forge task submit` shows no output for extended time, appears frozen.
2. **Deeper**: The quality gate runs `just test` which executes `go test -race ./...` on the entire project (896+ tests). This takes several minutes with race detection enabled.
3. **Deepest**: `RunCapture()` in `pkg/just/just.go` uses `exec.Command.CombinedOutput()` which **buffers all output until the command completes** before printing. The user sees zero feedback during the entire test run — no progress indicator, no streaming output. The gate is running correctly but appears dead because of buffered I/O.

## Solution

This is a UX issue, not a deadlock. The gate completes correctly given enough time. Two possible improvements:
1. Stream output in real-time (use `cmd.StdoutPipe()`/`cmd.StderrPipe()` instead of `CombinedOutput()`)
2. Print a progress message before each gate step (e.g., "Running quality gate: test...")

Workaround: run `just test` manually before `forge task submit` to verify tests pass, then the gate's test step will be fast (cached).

## Reusable Pattern

When wrapping long-running CLI commands in tools:
- **Never buffer output silently** — stream it or print periodic progress.
- **Distinguish slow from stuck**: if a command produces no output for >30s, it's either buffering or hung. Either way, the user can't tell the difference.
- `CombinedOutput()` is fine for sub-second commands. For anything that could take >5s, use streaming I/O.

## Example

```go
// Current: buffers all output (appears to hang)
output, err := cmd.CombinedOutput()

// Better: stream output in real-time
cmd.Stdout = os.Stderr
cmd.Stderr = os.Stderr
err := cmd.Run()
```

## Related Files

- `forge-cli/pkg/just/just.go:55-61` — RunCapture function using CombinedOutput
- `forge-cli/internal/cmd/task/submit.go:366-379` — validateQualityGate calling RunGate
