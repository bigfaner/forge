---
id: "5"
title: "Migrate qualitygate stderr calls"
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

# 5: Migrate qualitygate stderr calls

## Description

Migrate all `fmt.Fprintf(os.Stderr, ...)` and `fmt.Fprintln(os.Stderr, ...)` call sites in qualitygate command files to `forgelog` API. This is the heaviest area (~27 call sites across 3 files).

## Reference Files
- `docs/proposals/forge-cli-logging/proposal.md` — Call-Site Categorization Table (One-Time Migration Reference), Fprintln handling
- `forge-cli/internal/cmd/qualitygate/quality_gate.go` — 6 Fprintf + 5 Fprintln calls
- `forge-cli/internal/cmd/qualitygate/quality_gate_lifecycle.go` — 16 Fprintf + 11 Fprintln calls
- `forge-cli/internal/cmd/qualitygate/quality_gate_report.go` — Fprintf calls (verify exact count)

## Acceptance Criteria
- [ ] AC-1: All `fmt.Fprintf(os.Stderr, ...)` and `fmt.Fprintln(os.Stderr, ...)` calls in qualitygate/*.go are replaced with forgelog calls; no stderr writes remain
- [ ] AC-2: Console output from qualitygate files is byte-identical to pre-migration behavior; SC-2

## Hard Rules
- 仅修改以下文件：`internal/cmd/qualitygate/quality_gate.go`, `internal/cmd/qualitygate/quality_gate_lifecycle.go`, `internal/cmd/qualitygate/quality_gate_report.go`

## Implementation Notes
- Heaviest migration area: quality_gate_lifecycle.go has ~27 stderr writes
- Prefix classification: `ERROR:` (incl. 2-space indent) → `forgelog.Error()`, `WARNING:` (incl. 2-space indent) → `forgelog.Warn()`, prefixless progress/status → `forgelog.Info()`
- Fprintln callers must add `\n` explicitly
- Quality gate messages include indented variants (`  ERROR:`, `  WARNING:`) — prefix stripping handles these
