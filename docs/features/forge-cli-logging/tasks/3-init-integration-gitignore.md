---
id: "3"
title: "Integrate forgelog into commands and update gitignore"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [2]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 3: Integrate forgelog into commands and update gitignore

## Description

Wire `forgelog.Init()` into each command's `runE` function with `defer forgelog.Close()`. Add `.forge/logs/` to `gitignoreEntries` in `init.go`. Verify that `forge task submit` writes AUTO-RESTORE diagnostics to the log file and that `forge init` adds gitignore but does NOT create the logs directory.

## Reference Files
- `docs/proposals/forge-cli-logging/proposal.md` — Core Behaviors (directory auto-creation, emergency disable), Scope (In Scope)
- `forge-cli/internal/cmd/root.go` — Command registration and root runE
- `forge-cli/internal/cmd/init.go` — gitignoreEntries and directory creation

## Acceptance Criteria
- [ ] AC-1: `forge init` adds `.forge/logs/` to `.gitignore` but does NOT create `.forge/logs/` directory; SC-5
- [ ] AC-2: `forge task submit` with a fix-task scenario writes AUTO-RESTORE diagnostic to `.forge/logs/<datetime>-<pid>.log` with structured format `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3} \[(DEBUG|INFO|WARN|ERROR)\] .+`; SC-1

## Hard Rules
- `forgelog.Init()` called after config loading, before command logic
- `defer forgelog.Close()` immediately after Init

## Implementation Notes
- Init placement: after config is loaded (forgelog needs LogsConfig), before command-specific logic
- Pre-Init messages (config loading, flag parsing) are not captured — this is accepted
- Consider adding Init in a shared helper called from each command's runE, or in PersistentPreRunE if config is available there
- gitignoreEntries: add `".forge/logs/"` entry with `# Forge` comment group
