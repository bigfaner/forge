---
id: "6"
title: "Migrate task, feature_complete, and pkg stderr calls"
priority: "P1"
estimated_time: "1h"
complexity: "medium"
dependencies: [4, 5]
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 6: Migrate task, feature_complete, and pkg stderr calls

## Description

Migrate all remaining `fmt.Fprintf(os.Stderr, ...)`, `fmt.Fprintln(os.Stderr, ...)`, `slog.Warn(...)`, and `log.Printf(...)` call sites in task/submit, feature_complete, and pkg/ files to `forgelog` API. This is the final migration task — after completion, the grep verification in SC-8 must return 0 results.

## Reference Files
- `docs/proposals/forge-cli-logging/proposal.md` — Call-Site Categorization Table (One-Time Migration Reference), Fprintln handling
- `forge-cli/internal/cmd/task/submit.go` — AUTO-RESTORE / AUTO-RESTORE-SKIP / WARNING calls
- `forge-cli/internal/cmd/feature_complete.go` — Compound `[feature:complete]` prefix calls
- `forge-cli/pkg/serverprobe/serverprobe.go` — FAIL/OK/WARNING probe messages
- `forge-cli/pkg/just/just.go` — ERROR/WARNING just execution messages
- `forge-cli/pkg/task/state.go` — Mixed-case `Warning:` calls
- `forge-cli/pkg/task/category.go` — 1 `log.Printf` call
- `forge-cli/pkg/testrunner/testrunner.go` — WARNING call
- `forge-cli/pkg/forgeconfig/detect.go` — 1 `slog.Warn` call

## Acceptance Criteria
- [ ] AC-1: `grep -rE 'fmt\.Fprintf\(os\.Stderr|fmt\.Fprintln\(os\.Stderr|slog\.Warn|log\.Printf' forge-cli/internal/ forge-cli/pkg/ --include='*.go' | grep -v testdata | grep -v 'forensic/'` returns 0 results; SC-8
- [ ] AC-2: Console output from all migrated files is byte-identical to pre-migration behavior; SC-2

## Hard Rules
- 仅修改以下文件：`internal/cmd/task/submit.go`, `internal/cmd/feature_complete.go`, `pkg/serverprobe/serverprobe.go`, `pkg/just/just.go`, `pkg/task/state.go`, `pkg/task/category.go`, `pkg/testrunner/testrunner.go`, `pkg/forgeconfig/detect.go`

## Implementation Notes
- Compound prefix matching: `[feature:complete] Error:` → `forgelog.Error()`, `[feature:complete] Warning:` → `forgelog.Warn()`, `[feature:complete] Status` → `forgelog.Info()`, `[feature:complete] Push failed:` → `forgelog.Error()`
- `AUTO-RESTORE-SKIP:` → `forgelog.Warn()`, `AUTO-RESTORE:` → `forgelog.Info()`
- `FAIL:` → `forgelog.Warn()`, `OK:` → `forgelog.Info()`
- `Warning:` (mixed case) → `forgelog.Warn()` (case-insensitive)
- `log.Printf` in category.go → `forgelog.Info()` (with `\n` appended)
- `slog.Warn` in detect.go → `forgelog.Warn()` (with `\n` appended)
- Forensic/extract.go call sites are EXCLUDED from migration
