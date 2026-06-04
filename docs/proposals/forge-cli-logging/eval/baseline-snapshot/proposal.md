---
title: "Forge CLI Structured Logging"
slug: forge-cli-logging
status: draft
created: 2026-06-04
intent: enhancement
domains: [cli, logging, diagnostics, developer-experience]
---

# Problem

Forge CLI outputs diagnostic messages (AUTO-RESTORE-SKIP, WARNING, ERROR) exclusively to stderr. These messages are ephemeral — lost once the process exits. When the run-tasks dispatcher runs autonomous loops, diagnostic output scrolls away in subagent sessions and becomes impossible to trace after the fact.

Recent example: `autoRestoreSourceTask` silently returned without restoring a blocked source task. No log existed to diagnose why (not found? not blocked? unmet deps?). We had to speculate about the root cause from a lesson document instead of reading actual logs.

**Evidence**: 64 `fmt.Fprintf(os.Stderr, ...)` call sites across the CLI, zero persisted diagnostics.

**Cost of inaction**: Every future incident requires code-level speculation instead of log-based diagnosis. The AUTO-RESTORE debugging session that prompted this proposal took hours of code archaeography that a single log line would have resolved.

# Proposed Solution

Add a lightweight file-based logging layer to forge-cli:

1. **Log file per command execution**: Each `forge` invocation writes to `.forge/logs/<ISO-8601-datetime>.log` (e.g., `2026-06-04T17-30-00.log`)
2. **Level-filtered**: Four levels (DEBUG, INFO, WARN, ERROR). Only messages at or above the configured level are written. Configured via `.forge/config.yaml`
3. **Dual output**: Messages write to both log file and stderr (current behavior preserved)
4. **Auto-cleanup**: On each command startup, delete log files older than `retentionDays` (default 7)
5. **Categorized output**: Existing stderr messages are classified into levels based on their prefix:
   - `ERROR:`, `AUTO-RESTORE-SKIP:` → ERROR
   - `WARNING:`, `AUTO-RESTORE-SKIP:` (degraded) → WARN
   - `AUTO-RESTORE:`, `SOURCE-RESOLVE:`, `NOTE:` → INFO
   - `[debug]` → DEBUG

# Alternatives

## A. Do nothing

Keep stderr-only output. Zero implementation cost but zero improvement in debuggability. Every future incident remains a code archaeography exercise.

## B. Environment variable toggle (FORGE_LOG_FILE)

Use `FORGE_LOG_FILE=1 forge task claim` to opt into file logging. Simpler config but requires users to remember the env var. Doesn't integrate with the existing `.forge/config.yaml` system. No auto-cleanup.

## C. Structured JSON logging

Write JSON-formatted logs for machine parseability. Over-engineered for the current need (human diagnosis). Can be added later as a format option if needed.

**Recommendation**: Proceed with the proposed solution. It balances debuggability with minimal complexity, integrates with existing config, and provides auto-cleanup.

# Scope

## In Scope

- `pkg/forgelog` package: log level enum, file writer with auto-cleanup, dual-output (file + stderr)
- `.forge/config.yaml` `logs` section: `level` (debug/info/warn/error), `retentionDays` (default 7)
- Add `ForgeLogsDir` constant to `pkg/feature/constants.go`
- Migrate existing stderr calls to use `forgelog.Warn()`, `forgelog.Error()`, etc.
- `.forge/logs/` entry in `gitignoreEntries` (init.go)
- `forge init` ensures `.forge/logs/` directory exists
- `forgelog.Init()` called early in each command's `runE` function

## Out of Scope

- Structured/JSON log format
- Log rotation (one file per invocation is sufficient)
- Log aggregation or remote shipping
- Migrating test code (tests continue using stderr directly)
- Changes to the plugin (agents/commands/skills) — this is CLI-only

# Risks

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| Log file contention under concurrent commands | Low — each invocation creates a unique timestamped file | Use per-invocation filename; no shared file |
| Disk accumulation if retention too long | Low — default 7 days, configurable | Auto-cleanup on each command startup |
| Config parsing failure blocks logging | Low — defaults applied when config missing | Hardcoded defaults: level=info, retention=7 days |
| Performance impact of file I/O on every log call | Low — log writes are append-only, <1KB per invocation | No buffering needed; `os.OpenFile` with `O_APPEND` is efficient |

# Success Criteria

| ID | Criterion | Verification |
|----|-----------|-------------|
| SC-1 | `forge task submit` writes AUTO-RESTORE diagnostic to `.forge/logs/<datetime>.log` | Run submit with fix-task scenario; verify log file exists with expected content |
| SC-2 | Log level filtering works: setting `level: warn` suppresses INFO messages | Set config to warn; run command; verify only WARN/ERROR in log |
| SC-3 | Auto-cleanup deletes files older than `retentionDays` | Create log file with old timestamp; run any forge command; verify deletion |
| SC-4 | `forge init` creates `.forge/logs/` directory and adds `.forge/logs/` to `.gitignore` | Fresh project; run `forge init`; verify directory and gitignore entry |
| SC-5 | Dual output: same message appears in both stderr and log file | Run command; compare stderr output with log file content |
| SC-6 | Config missing or malformed falls back to defaults (info level, 7 days) | Delete config section; run command; verify logging still works with defaults |
