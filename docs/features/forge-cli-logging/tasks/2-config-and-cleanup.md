---
id: "2"
title: "Add config extension, validation, and auto-cleanup"
priority: "P1"
estimated_time: "1.5h"
complexity: "medium"
dependencies: [1]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 2: Add config extension, validation, and auto-cleanup

## Description

Extend `.forge/config.yaml` with a `logs` section (level, retentionDays, enabled). Add `LogsConfig` struct to `forgeconfig.Config` with `omitempty`. Implement config validation with safe defaults. Add auto-cleanup of old log files on startup. Add `ForgeLogsDir` constant.

## Reference Files
- `docs/proposals/forge-cli-logging/proposal.md` — Core Behaviors (auto-cleanup, emergency disable), .forge/config.yaml Extension, Constraints & Dependencies
- `forge-cli/pkg/forgeconfig/config.go` — Config struct and ReadConfig function

## Acceptance Criteria
- [ ] AC-1: On startup, after new log file is opened, delete log files older than `retentionDays` (default 7); active log never deleted by its own cleanup; SC-4
- [ ] AC-2: Config missing or malformed falls back to defaults (level=info, retentionDays=7, enabled=true); invalid level falls back to info; retentionDays<1 falls back to 7; SC-7
- [ ] AC-3: `FORGE_NO_LOG=1` env var or `logs.enabled: false` config skips FileBackend initialization; env var takes precedence; SC-10

## Hard Rules
- `omitempty` on Logs field — existing configs deserialize cleanly
- Cleanup errors silently ignored (e.g., locked files)

## Implementation Notes
- `ForgeLogsDir` constant = `".forge/logs"` in `pkg/feature/constants.go`
- Cleanup runs after new log file is opened (never delete active log)
- LogsConfig struct: `Enabled bool (default true)`, `Level string (default "info")`, `RetentionDays int (default 7)`
- Validation in `forgelog.Init()`: normalize level to lowercase, check retentionDays >= 1
