---
title: "Forge CLI Structured Logging"
slug: forge-cli-logging
status: draft
created: 2026-06-04
intent: enhancement
domains: [cli, logging, diagnostics, developer-experience]
---

# Problem

Forge CLI outputs diagnostic messages exclusively to stderr. These messages are ephemeral — lost once the process exits. When the run-tasks dispatcher runs autonomous loops, diagnostic output scrolls away in subagent sessions and becomes impossible to trace after the fact.

Recent example: `autoRestoreSourceTask` silently returned without restoring a blocked source task. No log existed to diagnose why. We had to speculate about the root cause from a lesson document instead of reading actual logs.

### Evidence

~100 stderr write call sites across the CLI (`fmt.Fprintf` + `fmt.Fprintln` patterns in `internal/` and `pkg/`), plus 1 `slog.Warn` and 1 `log.Printf`. Zero persisted diagnostics. Approximate counts (exact verification at implementation time):

| Pattern | internal/ | pkg/ | forensic/ (excluded) | Migratable |
|---------|-----------|------|---------------------|------------|
| `fmt.Fprintf(os.Stderr, ...)` | ~55 | 8 | 5 | ~63 |
| `fmt.Fprintln(os.Stderr, ...)` | ~32 | 5 | 5 | ~37 |
| `slog.Warn(...)` | 0 | 1 | 0 | 1 |
| `log.Printf(...)` | 0 | 1 | 0 | 1 |
| **Total** | ~87 | ~15 | 10 | **~102** |

Counts are approximate; exact verification at implementation time via:
```bash
grep -r 'fmt.Fprintf(os.Stderr\|fmt.Fprintln(os.Stderr' forge-cli/internal/ forge-cli/pkg/ --include='*.go' | grep -v testdata | grep -v 'forensic/'
```

### Urgency

Every future incident requires code archaeography instead of log-based diagnosis. The AUTO-RESTORE debugging session that prompted this proposal took hours of code archaeography that a single log line would have resolved.

# Proposed Solution

Add a lightweight file-based logging layer to forge-cli with a three-layer elegant architecture:

## 1. Architecture Layer — Backend Abstraction

`forgelog` acts as a unified output gateway. All diagnostic writes flow through it, with console and file as independent backends:

```go
// Backend is a log output target.
type Backend interface {
    Write(level LogLevel, timestamp time.Time, msg string)
    Close() error
}

// ConsoleBackend writes to stderr with original format.
// Output: just the message as-is (preserves current behavior exactly).
// No level filtering — always outputs all messages sent to it.
// Note: the existing verbose gate (base.Debugf) remains in the CALLER code,
// not in forgelog. Callers that currently gate debug output behind a flag
// continue to do so — they simply call forgelog.Debug() only when the flag
// is set. ConsoleBackend outputs whatever the caller sends, preserving
// byte-identical console behavior.
type ConsoleBackend struct{}

// FileBackend writes to log file with structured format.
// Output: 2026-06-04T17:30:00.123 [WARN] message
type FileBackend struct {
    mu   sync.Mutex
    file *os.File
}
```

Each `forgelog` call dispatches to all registered backends sequentially in registration order: ConsoleBackend first, then FileBackend. ConsoleBackend outputs the raw message; FileBackend adds timestamp+level prefix. The two are decoupled — changing file format never affects console output, and vice versa. FileBackend write errors are silently ignored so that a stalled or failing filesystem does not propagate errors to the caller. Note: because dispatch is synchronous, a severely stalled filesystem write could block subsequent console writes; `FORGE_NO_LOG=1` provides an escape hatch to disable file logging entirely if this becomes an issue in practice.

## 2. Format Layer — Dual Format Design

| Output Target | Format | Level Filtering | Rationale |
|--------------|--------|----------------|-----------|
| **Console (stderr)** | Original message, unchanged | None — all messages always output | Zero behavioral change. Existing scripts, pipes, and user expectations are preserved exactly |
| **Log file** | `2006-01-02T15:04:05.000 [LEVEL] message` | Yes — only messages at or above configured level | Structured prefix enables grep/filter/trace. Local time with millisecond precision |

Console output is byte-identical to pre-migration behavior. A test can diff stderr output before and after migration and expect zero changes.

## 3. API Layer — Printf-Style One-Liners

```go
package forgelog

// Init initializes the logging layer.
// - Creates .forge/logs/ on demand via os.MkdirAll(logDir, 0700).
// - Falls back to console-only if directory creation fails.
// - Checks FORGE_NO_LOG=1 — if set, skips FileBackend.
// logsDir is derived as: filepath.Join(projectRoot, ForgeLogsDir)
// where projectRoot is the resolved project directory and ForgeLogsDir
// is the constant ".forge/logs" defined in pkg/feature/constants.go.
func Init(config *forgeconfig.LogsConfig, logsDir string) error

// Printf-style: Warn("WARNING: task %s not found", id)
// One call dispatches to all backends.
func Debug(format string, args ...interface{})
func Info(format string, args ...interface{})
func Warn(format string, args ...interface{})
func Error(format string, args ...interface{})

// Close releases all backend resources (file handles). Call via defer in each command's runE.
// With O_APPEND per-write and no buffering, Close() is not strictly needed for data safety —
// it only releases the file handle cleanly. os.Exit/log.Fatal bypass defer but cannot lose
// data since all writes were already issued to the OS. The OS reclaims all file handles on
// process exit, so a leaked handle from a bypassed defer does not cause resource exhaustion.
// The existing log.Printf call site will be migrated to forgelog, eliminating log.Fatal usage.
func Close()
```

Migration at each call site is a one-line mechanical change:
```go
// Before:
fmt.Fprintf(os.Stderr, "WARNING: task %s not found\n", id)
// After:
forgelog.Warn("WARNING: task %s not found\n", id)
```

## Core Behaviors

1. **Log file per command execution**: Each `forge` invocation writes to `.forge/logs/<ISO-8601-datetime>-<pid>.log`. PID suffix prevents collisions when multiple commands run within the same second (e.g., concurrent subagent sessions).

2. **Level-filtered (file only)**: Four levels (DEBUG, INFO, WARN, ERROR). Only messages at or above the configured level are written to the **file backend**. The **console backend always outputs all messages** — this preserves current behavior exactly. Configured via `.forge/config.yaml` `logs.level` (default: `info`). Console has no level filter; the level config controls only what persists to disk.

3. **Auto-cleanup**: On each command startup, **after** the new log file is opened, delete log files older than `retentionDays` (default 7). Cleanup errors are silently ignored.

4. **Emergency disable**: `FORGE_NO_LOG=1` env var or `logs.enabled: false` in config skips FileBackend initialization. ConsoleBackend continues unchanged — identical to pre-migration behavior. Env var takes precedence over config.

5. **Directory auto-creation**: `forgelog.Init()` creates `.forge/logs/` on demand via `os.MkdirAll(logDir, 0700)`. `forge init` does NOT create this directory — it only exists when logging actually happens.

### Innovation Highlights

This is a straightforward adoption of a standard backend-pattern logging architecture. The creative insight is the **zero-change console contract**: by treating console as a first-class backend that preserves original format byte-for-byte, the migration is safe by construction. No integration test suite is needed to verify console output — if the console backend's Write method outputs the raw message unchanged, behavioral equivalence is guaranteed.

The per-invocation log file (timestamp+PID naming) borrows from web server access log patterns (nginx, apache) where each process writes to its own file, avoiding contention without lock files.

### Call-Site Categorization Table (One-Time Migration Reference)

This table classifies all existing stderr call sites for the one-time migration. It is **not a runtime feature** — prefix parsing exists only during migration to determine the correct `forgelog` level function for each call site. New code calls `forgelog.Warn()` directly and needs no prefix convention. After migration, this table has no ongoing maintenance value.

All ~102 stderr write call sites are classified below. Prefix parsing is **only for migrating existing call sites** — new code calls `forgelog.Warn()` directly and needs no prefix convention.

| Prefix Pattern | Level | Source Areas | Notes |
|---|---|---|---|
| `ERROR: ...` / `error: ...` (case-insensitive) | ERROR | init.go, init_config.go, quality_gate.go, errors.go, upgrade.go, etc. | Case-insensitive match; includes 2-space indented variants |
| `  ERROR: ...` (2-space indent) | ERROR | quality_gate_lifecycle.go | Indented context messages |
| `ERROR_CODE: ...` / `CAUSE: ...` / `ACTION: ...` | ERROR | errors.go | Structured error block parts |
| `HINT: ...` | INFO | errors.go | Remediation suggestions within structured error blocks — not errors themselves |
| `WARNING: ...` (no indent) | WARN | task/*.go, init.go, qualitygate/*.go, errors.go, submit.go | Standard warning prefix |
| `  WARNING: ...` (2-space indent) | WARN | quality_gate_lifecycle.go, serverprobe.go | Indented context messages |
| `AUTO-RESTORE-SKIP: ...` | WARN | submit.go | Skip conditions for auto-restore |
| `AUTO-RESTORE: ...` | INFO | submit.go | Successful restore notification |
| `SOURCE-RESOLVE: ...` | INFO | pkg/task/add.go | Source resolution trace |
| `NOTE: ...` | INFO | index.go | Informational notes |
| `[debug] ...` | DEBUG | base/output.go | Debug trace messages |
| `[feature:complete] Error: ...` | ERROR | feature_complete.go | Compound prefix; longest-prefix matching |
| `[feature:complete] Status ...` | INFO | feature_complete.go | Status update |
| `[feature:complete] Warning: ...` | WARN | feature_complete.go | Warning within feature completion |
| `[feature:complete] Push failed: ...` | ERROR | feature_complete.go | Push failure |
| `FAIL: ...` / `OK: ...` | WARN / INFO | pkg/serverprobe/serverprobe.go | Probe result messages |
| `Warning: ...` (mixed case) | WARN | pkg/task/state.go | Case-insensitive matching |
| Prefixless (progress/status) | INFO | qualitygate/*.go, output.go, base/output.go | Progress bars, orchestration status, probe messages |
| Forensic (all patterns) | EXCLUDED | forensic/extract.go | stderr is primary output channel; excluded |

**Matching priority**: Longest-prefix-first. Compound prefixes (`[feature:complete] Error:`) match before simple prefixes (`ERROR:`).

**Fallback**: Any message without a matching prefix defaults to INFO level. Note: prefixless sites (~30+) require individual review at implementation time — the "one-line mechanical" migration claim applies primarily to prefixed sites. Prefixless sites will be audited and assigned correct levels case-by-case.

**Case-insensitive**: `strings.HasPrefix(strings.ToUpper(strings.TrimSpace(msg)), prefix)`. Leading whitespace stripped before matching.

**Fprintln handling**: `fmt.Fprintln(os.Stderr, msg)` appends `\n` automatically. **forgelog functions do NOT append `\n`** — the formatted message is output exactly as-is. Migration rule for Fprintln callers: the caller must include `\n` explicitly (e.g., `fmt.Fprintln(os.Stderr, "msg")` becomes `forgelog.Info("msg\n")`).

### .forge/config.yaml Extension

```yaml
logs:
  enabled: true         # set false to disable file logging (rollback mechanism)
  level: info           # debug | info | warn | error (default: info)
  retentionDays: 7      # minimum 1 (default: 7)
```

Config struct:
```go
type LogsConfig struct {
    Enabled       bool   `yaml:"enabled"`        // default: true; set false to disable file logging
    Level         string `yaml:"level"`          // default: "info"
    RetentionDays int    `yaml:"retentionDays"`  // default: 7
}
// Added to existing Config struct:
// Logs *LogsConfig `yaml:"logs,omitempty"`
```

The `omitempty` tag ensures existing configs without `logs` section deserialize cleanly — `Logs` is `nil`, defaults applied in `forgelog.Init()`.

# Requirements Analysis

### Key Scenarios

- **Happy path**: User runs `forge task submit`; diagnostic messages appear on console (unchanged) and are persisted to `.forge/logs/` for later review
- **Concurrent invocations**: Multiple subagent sessions under `run-tasks` each write to separate log files; no contention
- **Config missing**: User has no `logs` section in config; logging works with defaults (info level, 7-day retention)
- **Disk failure**: `.forge/logs/` creation fails (read-only filesystem, permission denied); falls back to console-only mode silently
- **Emergency disable**: User sets `FORGE_NO_LOG=1`; file logging disabled, console output unchanged
- **Log cleanup**: On startup, files older than `retentionDays` are deleted; current invocation's log is never deleted

### Non-Functional Requirements

- **Performance**: Per-write `O_APPEND` direct to file. No bufio layer. Acceptable for typical volumes (~50-200 lines per invocation). `O_APPEND` provides atomic appends under `PIPE_BUF` on most OS.
- **Security**: Log files created with mode `0600` (owner-only). Directories with `0700`. `.forge/logs/` added to `.gitignore` by `forge init`.
- **Concurrency**: FileBackend uses `sync.Mutex` to serialize writes. ConsoleBackend writes to stderr (already thread-safe in Go's `fmt.Fprintf`).
- **Data safety**: No user-space buffering — each write is issued to the OS (via `O_APPEND`) before the function returns. Persistence to stable storage depends on OS behavior. Normal exits, panics (with recover), and `os.Exit` will not lose data since writes reach the kernel before the call returns. `SIGKILL` may lose data if the kernel has not flushed its buffer cache.

### Constraints & Dependencies

- No external dependencies — `forgelog` uses only stdlib (`fmt`, `os`, `sync`, `time`, `strings`)
- No changes to plugin layer (agents/commands/skills) — this is CLI-only
- Go 1.21+ compatibility required
- **Governance**: all new CLI diagnostic output must use forgelog. Direct `fmt.Fprintf(os.Stderr, ...)` calls are banned in new code outside the forensic command (which is excluded from forgelog). This prevents gradual migration erosion.
- Windows: file permission modes (`0600`/`0700`) are advisory-only on Windows; `O_APPEND` atomicity guarantees are Unix-specific. Log files will be created on Windows but without owner-only enforcement. PID-based naming works cross-platform

# Alternatives & Industry Benchmarking

### Industry Solutions

Most CLI tools use one of: (1) stderr-only output (current forge approach), (2) `log/slog` with structured handlers, (3) per-invocation log files. The proposed solution uses pattern (3). Notable CLI tools using per-invocation log files: **cargo** writes build diagnostics to `target/debug/` with per-build output; **kubectl** persists event logs with `--v` verbosity control per invocation; **docker** uses `json-file` log driver with per-container log files. Unlike these tools, forgelog targets a single-user CLI (not a daemon) and prioritizes zero-change console output over structured schemas. The backend abstraction is inspired by slog's Handler interface but simplified for human-readable diagnostics — no key-value pairs, no group support, no structured schema needed.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Zero improvement. Every incident = code archaeography | Rejected: cost of inaction too high |
| `log/slog` structured logging | Go stdlib (1.21+) | Stdlib, structured key-value, maintained by Go team | Two custom slog.Handlers needed: one stripping structure for raw console (~40 lines), one formatting for file (~30 lines), plus slog.Logger wiring (~20 lines). Total: ~90 lines vs forgelog's ~150 lines — slog saves ~60 lines but adds a dependency on slog's Handler contract, key-value plumbing, and group semantics that forgelog never uses. The "complexity mismatch" is cognitive, not line-count: future maintainers must understand slog's Handler lifecycle (Enabled, Handle, WithAttrs, WithGroup) for a use case that needs none of them | Rejected: cognitive overhead mismatch |
| Env var toggle (`FORGE_LOG_FILE=1`) | Common CLI pattern | Simpler opt-in | Doesn't integrate with existing config. No auto-cleanup | Rejected: poor UX |
| **Backend-pattern forgelog** | Web server access log pattern | Minimal, zero-change console, pluggable, auto-cleanup | Team-maintained (~150 lines) | **Selected: matches problem scope precisely** |

# Feasibility Assessment

### Technical Feasibility

Straightforward Go implementation. No external dependencies. The Backend interface is 2 methods. Each call site migration is a one-line mechanical change.

### Resource & Timeline

Single PR, mechanically verifiable. Implementation order: (1) `pkg/forgelog` package + tests, (2) config extension, (3) gitignore entry, (4) migrate all ~102 call sites.

### Dependency Readiness

No upstream dependencies. Config struct extension is additive (omitempty).

# Assumptions Challenged

| Assumption | Challenge Tool | Finding |
|------------|---------------|---------|
| "Need `forge init` to create logs dir" | Occam's Razor | Overturned: forgelog.Init() auto-creates on demand. No need to involve forge init in directory creation. forge init only adds `.forge/logs/` to .gitignore. |
| "Need 80 call sites migrated" | Stress Test | Refined: actual scope is ~102 sites (Fprintf + Fprintln + slog + log). Proposal's 80 count was Fprintf-only, missing Fprintln patterns. |
| "stderr-first-then-file ordering matters" | Assumption Flip | Overturned: with Backend abstraction, ordering is backend registration order. ConsoleBackend first, FileBackend second. If file write fails, console has already written. Same safety, cleaner abstraction. |
| "forgelog needs prefix parsing at runtime" | Occam's Razor | Confirmed but scope-limited: prefix parsing exists only in the migration layer to classify existing call sites. New code uses explicit level calls. No ongoing maintenance burden. |

# Scope

### In Scope

- `pkg/forgelog` package: LogLevel enum, Backend interface, ConsoleBackend, FileBackend, Init/Close, printf-style API
- `.forge/config.yaml` `logs` section: `level` (debug/info/warn/error), `retentionDays` (default 7)
- `LogsConfig` struct added to `forgeconfig.Config` with `omitempty`
- Config validation in `forgelog.Init()`: unrecognized level → default `info`; `retentionDays < 1` → default 7
- `ForgeLogsDir` constant in `pkg/feature/constants.go`
- `.forge/logs/` entry in `gitignoreEntries` (init.go)
- `forgelog.Init()` called early in each command's `runE` function; `defer forgelog.Close()` follows. **Known gap**: messages emitted before `Init()` (e.g., in `init()` functions or package-level vars) are not captured in the log file — they only appear on console. This is acceptable since pre-Init messages are rare and non-diagnostic
- Migrate all ~102 stderr write call sites in `internal/` and `pkg/` to `forgelog` API (single PR)
- Exclude all forensic command call sites from migration
- **Emergency disable**: `FORGE_NO_LOG=1` env var or `logs.enabled: false` config skips FileBackend
- **Directory auto-creation**: `forgelog.Init()` creates `.forge/logs/` via `os.MkdirAll(logDir, 0700)`
- **Log file permissions**: `0600` for files, `0700` for directory

### Out of Scope

- Structured/JSON log format
- Log rotation (one file per invocation is sufficient)
- Log aggregation or remote shipping
- CLI log viewer command (`forge log` / `forge logs`)
- Migrating test code (tests continue using stderr directly)
- Changes to the plugin (agents/commands/skills)
- `forge init` creating `.forge/logs/` directory
- Migrating the `forensic` command (stderr is its primary output channel)

# Key Risks

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Log file contention under concurrent commands | Low | Low | Per-invocation filename (timestamp+PID). No shared file |
| Disk accumulation | Low | Low | Auto-cleanup on each startup. Typical budget: ~5-10KB/invocation. Worst case for run-tasks loops: a single loop dispatching 200 subagent invocations at ~50KB each = ~10MB/day. Over 7-day retention = ~70MB. Sustained heavy use (1000 invocations/day) = ~350MB over 7 days — still negligible on modern disks |
| Config parsing failure blocks logging | Low | Low | Hardcoded defaults applied when config missing or invalid |
| Sensitive info in log files | Medium | Medium | File mode `0600`, dir `0700`, `.gitignore` entry. No redaction at log time — diagnostic value > risk for local files. In CI/container environments, `FORGE_NO_LOG=1` or `logs.enabled: false` recommended to prevent sensitive data persistence in shared or ephemeral filesystems |
| Logging layer regression | Low | High | `FORGE_NO_LOG=1` env var or `logs.enabled: false` config for rollback. Migration is one-line mechanical changes |
| Migration count inaccuracy | Medium | Low | Exact verification at implementation time. Grep-based CI check in PR |
| `fmt.Fprintln` edge cases | Low | Low | Fprintln adds `\n` automatically; forgelog does not append `\n`. Migration must add `\n` to each former Fprintln call site |

# Success Criteria

- [ ] SC-1: `forge task submit` writes AUTO-RESTORE diagnostic to `.forge/logs/<datetime>-<pid>.log` with structured format — each line matches regex `^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3} \[(DEBUG|INFO|WARN|ERROR)\] .+`
- [ ] SC-2: Console output is byte-identical before and after migration — diff of stderr shows zero changes (console backend has no level filter, always outputs all messages)
- [ ] SC-3: Level filtering works on file backend: `level: warn` suppresses DEBUG and INFO messages in log file only; console still shows all messages
- [ ] SC-4: Auto-cleanup deletes files older than `retentionDays`; active log file is never deleted by its own cleanup pass
- [ ] SC-5: `forge init` adds `.forge/logs/` to `.gitignore` but does NOT create the directory
- [ ] SC-6: `.forge/logs/` is auto-created by `forgelog.Init()` on first log write; falls back to console-only if creation fails
- [ ] SC-7: Config missing or malformed falls back to defaults (info level, 7 days retention) without error
- [ ] SC-8: Migration completeness: `grep -rE 'fmt\.Fprintf\(os\.Stderr|fmt\.Fprintln\(os\.Stderr|slog\.Warn|log\.Printf' forge-cli/internal/ forge-cli/pkg/ --include='*.go' | grep -v testdata | grep -v 'forensic/'` returns 0 results
- [ ] SC-9: Concurrent commands produce separate log files with distinct PIDs
- [ ] SC-10: `FORGE_NO_LOG=1` or `logs.enabled: false` suppresses file logging; stderr output unchanged; no `.forge/logs/` directory created
- [ ] SC-11: Log file created with mode `0600`; `.forge/logs/` directory created with mode `0700` (Unix only; advisory on Windows)
- [ ] SC-12: Config validation: `level: "bogus"` falls back to `info`; `retentionDays: -1` falls back to `7` — verified by checking default values in written log file content and cleanup behavior
