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

~100 stderr write call sites across the CLI (`fmt.Fprintf` + `fmt.Fprintln` patterns in `internal/` and `pkg/`), plus 1 `slog.Warn` and 1 `log.Printf`. Zero persisted diagnostics. Exact counts:

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
// No level filtering — always outputs all messages.
type ConsoleBackend struct{}

// FileBackend writes to log file with structured format.
// Output: 2026-06-04T17:30:00.123 [WARN] message
type FileBackend struct {
    mu   sync.Mutex
    file *os.File
}
```

Each `forgelog` call dispatches to all registered backends. ConsoleBackend outputs the raw message; FileBackend adds timestamp+level prefix. The two are decoupled — changing file format never affects console output, and vice versa.

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
func Init(config *forgeconfig.LogsConfig, logsDir string) error

// Printf-style: Warn("WARNING: task %s not found", id)
// One call dispatches to all backends.
func Debug(format string, args ...interface{})
func Info(format string, args ...interface{})
func Warn(format string, args ...interface{})
func Error(format string, args ...interface{})

// Close closes all backends. Call via defer in each command's runE.
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

4. **Emergency disable**: `FORGE_NO_LOG=1` env var skips FileBackend initialization. ConsoleBackend continues unchanged — identical to pre-migration behavior.

5. **Directory auto-creation**: `forgelog.Init()` creates `.forge/logs/` on demand via `os.MkdirAll(logDir, 0700)`. `forge init` does NOT create this directory — it only exists when logging actually happens.

### Innovation Highlights

This is a straightforward adoption of a standard backend-pattern logging architecture. The creative insight is the **zero-change console contract**: by treating console as a first-class backend that preserves original format byte-for-byte, the migration is safe by construction. No integration test suite is needed to verify console output — if the console backend's Write method outputs the raw message unchanged, behavioral equivalence is guaranteed.

The per-invocation log file (timestamp+PID naming) borrows from web server access log patterns (nginx, apache) where each process writes to its own file, avoiding contention without lock files.

### Call-Site Categorization Table

All ~102 stderr write call sites are classified below. Prefix parsing is **only for migrating existing call sites** — new code calls `forgelog.Warn()` directly and needs no prefix convention.

| Prefix Pattern | Level | Source Areas | Notes |
|---|---|---|---|
| `ERROR: ...` (no indent) | ERROR | init.go, init_config.go, quality_gate.go, errors.go, etc. | Uppercase `ERROR:` prefix |
| `  ERROR: ...` (2-space indent) | ERROR | quality_gate_lifecycle.go | Indented context messages |
| `error: ...` (lowercase) | ERROR | upgrade.go | Case-insensitive match |
| `ERROR_CODE: ...` / `CAUSE: ...` / `ACTION: ...` | ERROR | errors.go | Structured error block parts |
| `HINT: ...` | ERROR | errors.go | Part of structured error/warning blocks |
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

**Fallback**: Any message without a matching prefix defaults to INFO level.

**Case-insensitive**: `strings.HasPrefix(strings.ToUpper(strings.TrimSpace(msg)), prefix)`. Leading whitespace stripped before matching.

**Fprintln handling**: `fmt.Fprintln(os.Stderr, msg)` is equivalent to `fmt.Fprintf(os.Stderr, msg + "\n")`. Migration treats them identically — `forgelog.Warn(msg)` handles the trailing newline.

### .forge/config.yaml Extension

```yaml
logs:
  level: info           # debug | info | warn | error (default: info)
  retentionDays: 7      # minimum 1 (default: 7)
```

Config struct:
```go
type LogsConfig struct {
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
- **Data safety**: No buffering — each write is persisted before function returns. `os.Exit`, panic, and signal kills cannot lose data.

### Constraints & Dependencies

- No external dependencies — `forgelog` uses only stdlib (`fmt`, `os`, `sync`, `time`, `strings`)
- No changes to plugin layer (agents/commands/skills) — this is CLI-only
- Go 1.21+ compatibility required

# Alternatives & Industry Benchmarking

### Industry Solutions

Most CLI tools use one of: (1) stderr-only output (current forge approach), (2) `log/slog` with structured handlers, (3) per-invocation log files (web server pattern). The proposed solution uses pattern (3) with a backend abstraction inspired by slog's Handler interface but simplified for human-readable diagnostic logs.

### Comparison Table

| Approach | Source | Pros | Cons | Verdict |
|----------|--------|------|------|---------|
| Do nothing | — | Zero cost | Zero improvement. Every incident = code archaeography | Rejected: cost of inaction too high |
| `log/slog` structured logging | Go stdlib (1.21+) | Stdlib, structured key-value, maintained by Go team | Over-engineered for human-readable diagnostics. Dual-output (stderr original + file structured) doesn't map to slog's single-format Handler | Rejected: complexity mismatch |
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
- `forgelog.Init()` called early in each command's `runE` function; `defer forgelog.Close()` follows
- Migrate all ~102 stderr write call sites in `internal/` and `pkg/` to `forgelog` API (single PR)
- Exclude all forensic command call sites from migration
- **Emergency disable**: `FORGE_NO_LOG=1` env var skips FileBackend
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
| Disk accumulation | Low | Low | Auto-cleanup on each startup. Budget: ~5-10KB/invocation, 7-day default ≈ 7-70MB worst case |
| Config parsing failure blocks logging | Low | Low | Hardcoded defaults applied when config missing or invalid |
| Sensitive info in log files | Medium | Medium | File mode `0600`, dir `0700`, `.gitignore` entry. No redaction at log time — diagnostic value > risk for local files |
| Logging layer regression | Low | High | `FORGE_NO_LOG=1` emergency disable. Migration is one-line mechanical changes |
| Migration count inaccuracy | Medium | Low | Exact verification at implementation time. Grep-based CI check in PR |
| `fmt.Fprintln` edge cases | Low | Low | Fprintln adds `\n` automatically; forgelog API expects caller to include `\n`. Migration strips trailing newline from Fprintln calls |

# Success Criteria

- [ ] SC-1: `forge task submit` writes AUTO-RESTORE diagnostic to `.forge/logs/<datetime>-<pid>.log` with structured format (`timestamp [LEVEL] message`)
- [ ] SC-2: Console output is byte-identical before and after migration — diff of stderr shows zero changes (console backend has no level filter, always outputs all messages)
- [ ] SC-3: Level filtering works on file backend: `level: warn` suppresses DEBUG and INFO messages in log file only; console still shows all messages
- [ ] SC-4: Auto-cleanup deletes files older than `retentionDays`; active log file is never deleted by its own cleanup pass
- [ ] SC-5: `forge init` adds `.forge/logs/` to `.gitignore` but does NOT create the directory
- [ ] SC-6: `.forge/logs/` is auto-created by `forgelog.Init()` on first log write; falls back to console-only if creation fails
- [ ] SC-7: Config missing or malformed falls back to defaults (info level, 7 days retention) without error
- [ ] SC-8: Migration completeness: `grep -r 'fmt.Fprintf(os.Stderr\|fmt.Fprintln(os.Stderr' forge-cli/internal/ forge-cli/pkg/ --include='*.go' | grep -v testdata | grep -v 'forensic/'` returns 0 results
- [ ] SC-9: Concurrent commands produce separate log files with distinct PIDs
- [ ] SC-10: `FORGE_NO_LOG=1` suppresses file logging; stderr output unchanged; no `.forge/logs/` directory created
