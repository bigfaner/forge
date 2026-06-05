---
id: "1"
title: "Implement core forgelog package"
priority: "P0"
estimated_time: "2h"
complexity: "high"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Implement core forgelog package

## Description

Create `pkg/forgelog` — the unified diagnostic output gateway for forge-cli. Implements the Backend abstraction pattern with ConsoleBackend (raw stderr, no level filter) and FileBackend (timestamp+level prefix, level-filtered). Each forge invocation writes to `.forge/logs/<ISO-8601-datetime>-<pid>.log`. Directory auto-created on demand via `os.MkdirAll`. File permissions 0600/0700.

## Reference Files
- `docs/proposals/forge-cli-logging/proposal.md` — Architecture Layer, Format Layer, API Layer, Core Behaviors
- `forge-cli/pkg/feature/constants.go` — ForgeLogsDir constant location
- `forge-cli/internal/cmd/base/output.go` — Existing Debugf pattern (verbose gate reference)

## Acceptance Criteria
- [ ] AC-1: `forgelog.Warn("WARNING: task %s not found\n", "x")` outputs `WARNING: task x not found\n` to stderr (byte-identical to `fmt.Fprintf(os.Stderr, ...)`); SC-2
- [ ] AC-2: FileBackend writes to `.forge/logs/<ISO-8601-datetime>-<pid>.log` with format `2006-01-02T15:04:05.000 [WARN] message`; level filtering suppresses below configured level in file only; SC-3
- [ ] AC-3: `forgelog.Init()` creates `.forge/logs/` via `os.MkdirAll(dir, 0700)` on demand; falls back to console-only if creation fails; SC-6
- [ ] AC-4: Two concurrent `forgelog.Init()` calls produce separate log files with distinct PIDs; SC-9
- [ ] AC-5: Log file created with mode `0600`, directory with `0700` (Unix; advisory on Windows); SC-11

## Hard Rules
- forgelog functions do NOT append `\n` — the formatted message is output exactly as-is
- FileBackend write errors are silently ignored (never propagate to caller)
- No external dependencies — stdlib only

## Implementation Notes
- Backend interface: `Write(level LogLevel, timestamp time.Time, msg string)` + `Close() error`
- ConsoleBackend dispatches first, FileBackend second (sequential)
- FileBackend uses `sync.Mutex` for concurrent write safety
- Each write uses `O_APPEND` for atomic appends; no bufio layer
- File naming: `time.Now().Format("2006-01-02T15-04-05") + fmt.Sprintf("-%d", os.Getpid()) + ".log"`
- `forgelog.Init()` accepts `*forgeconfig.LogsConfig` (nil = defaults) and `logsDir string` (filepath.Join(projectRoot, ForgeLogsDir))
- Verbose gate (base.Debugf) stays in caller code, not in forgelog — ConsoleBackend outputs whatever caller sends
