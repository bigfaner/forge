---
id: "4"
title: "Migrate base, init, and upgrade stderr calls"
priority: "P1"
estimated_time: "1.5h"
complexity: "medium"
dependencies: [3]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 4: Migrate base, init, and upgrade stderr calls

## Description

Migrate all `fmt.Fprintf(os.Stderr, ...)` and `fmt.Fprintln(os.Stderr, ...)` call sites in base, init, and upgrade command files to `forgelog` API. Each call site is a one-line mechanical change classified by prefix pattern.

## Reference Files
- `docs/proposals/forge-cli-logging/proposal.md` — Call-Site Categorization Table (One-Time Migration Reference), Fprintln handling
- `forge-cli/internal/cmd/init.go` — 6 Fprintf calls
- `forge-cli/internal/cmd/base/errors.go` — 8 Fprintf + 4 Fprintln calls
- `forge-cli/internal/cmd/base/output.go` — 1 Fprintf call (Debugf)
- `forge-cli/internal/cmd/upgrade.go` — 5 Fprintf calls (lowercase `error:`)

## Acceptance Criteria
- [ ] AC-1: All `fmt.Fprintf(os.Stderr, ...)` and `fmt.Fprintln(os.Stderr, ...)` calls in init.go, base/errors.go, base/output.go, upgrade.go are replaced with forgelog calls; no stderr writes remain in these files
- [ ] AC-2: Console output from these files is byte-identical to pre-migration behavior; SC-2

## Hard Rules
- 仅修改以下文件：`internal/cmd/init.go`, `internal/cmd/base/errors.go`, `internal/cmd/base/output.go`, `internal/cmd/upgrade.go`

## Implementation Notes
- Prefix classification: `ERROR:` → `forgelog.Error()`, `error:` (lowercase) → `forgelog.Error()`, `WARNING:` → `forgelog.Warn()`, `[debug]` → `forgelog.Debug()`, prefixless → `forgelog.Info()`
- Fprintln callers must add `\n` explicitly (forgelog does not append `\n`)
- `base.Debugf(verbose, format, args...)` keeps verbose gate in caller: `if verbose { forgelog.Debug(...) }`
- init_config.go should be checked for stderr calls and included if found
