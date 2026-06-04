---
title: "Forge CLI Structured Logging"
slug: forge-cli-logging
status: draft
created: 2026-06-04
intent: enhancement
domains: [cli, logging, diagnostics, developer-experience]
---

# Problem

Forge CLI outputs diagnostic messages (AUTO-RESTORE-SKIP, WARNING, ERROR) exclusively to stderr. These messages are ephemeral — lost once the process exits. When the run-tasks dispatcher runs autonomous loops, diagnostic output scrolls away in subagent sessions and becomes impossible to trace after the fact. While the motivating incident occurred in run-tasks, all forge commands produce diagnostic stderr — logging all commands provides consistent debuggability and avoids the complexity of selective per-command logging. A targeted dispatcher-only approach would require a separate logging configuration surface for each command category, adding more complexity than logging uniformly.

Recent example: `autoRestoreSourceTask` silently returned without restoring a blocked source task. No log existed to diagnose why (not found? not blocked? unmet deps?). We had to speculate about the root cause from a lesson document instead of reading actual logs.

**Evidence**: 72 `fmt.Fprintf(os.Stderr, ...)` call sites in `internal/` + 8 call sites in `pkg/` = 80 total across the CLI, zero persisted diagnostics.

**Cost of inaction**: Every future incident requires code-level speculation instead of log-based diagnosis. The AUTO-RESTORE debugging session that prompted this proposal took hours of code archaeography that a single log line would have resolved.

# Proposed Solution

Add a lightweight file-based logging layer to forge-cli:

1. **Log file per command execution**: Each `forge` invocation writes to `.forge/logs/<ISO-8601-datetime>-<pid>.log` (e.g., `2026-06-04T17-30-00-45231.log`). The PID suffix prevents filename collisions when multiple forge commands are invoked within the same second (e.g., concurrent subagent sessions under `run-tasks`).
2. **Level-filtered**: Four levels (DEBUG, INFO, WARN, ERROR). Only messages at or above the configured level are written. Configured via `.forge/config.yaml`
3. **Dual output**: Messages write to both log file and stderr (current behavior preserved). Write ordering is **stderr-first-then-file**: the message is always written to stderr first, then to the log file. If file write fails (disk full, permission error, etc.), the stderr output has already occurred — no diagnostic information is lost. File write errors are silently ignored to avoid recursive logging. Each forgelog call performs a direct write to the file (opened with `O_APPEND`), ensuring messages are persisted immediately — no data is lost on unclean exit (`os.Exit`, panic, signal kill). `O_APPEND` provides atomic appends on most operating systems for writes under `PIPE_BUF` (typically 4KB or more), making per-write calls efficient without a bufio layer.
4. **Auto-cleanup**: On each command startup, **after** the new log file has been successfully opened and is being written to, delete log files older than `retentionDays` (default 7). Cleanup errors are best-effort and silently ignored (e.g., if a file is locked by another process). Running cleanup after opening the new log file ensures the active log is never deleted by its own cleanup pass.
5. **Per-line log format**: Each log line follows the format `2006-01-02T15:04:05.000 [LEVEL] message\n`. Timestamp is local time with millisecond precision. Level is one of `DEBUG`, `INFO`, `WARN`, `ERROR`. Message is the original stderr text with its prefix preserved (e.g., `WARNING: task not found` remains as-is in the message field — the `[LEVEL]` tag provides the structured level).
6. **Categorized output**: Existing stderr messages are classified into levels based on their prefix. Prefix parsing is **only for migrating existing call sites** — new code calls `forgelog.Warn()` etc. directly and needs no prefix convention.
7. **Emergency disable**: Set `FORGE_NO_LOG=1` environment variable to disable all file logging. When set, `forgelog` functions write to stderr only (identical to pre-migration behavior). This provides an escape hatch if the logging layer causes regressions.

### Call-Site Categorization Table

All 80 `fmt.Fprintf(os.Stderr, ...)` call sites are classified below (72 in `internal/`, 8 in `pkg/`). This table covers all known prefixes found by `grep -r 'fmt.Fprintf(os.Stderr' forge-cli/internal/ forge-cli/pkg/ --include='*.go' | grep -v testdata`. Any future message that does not match a listed prefix defaults to INFO level.

| Prefix Pattern | Level | Count | Source Files | Notes |
|---|---|---|---|---|
| `ERROR: ...` (no indent) | ERROR | 10 | init.go, init_config.go, quality_gate.go, errors.go, etc. | Uppercase `ERROR:` prefix |
| `  ERROR: ...` (2-space indent) | ERROR | 2 | quality_gate_lifecycle.go | Indented context messages; prefix stripping removes leading whitespace |
| `error: ...` (lowercase) | ERROR | 5 | upgrade.go | Matched case-insensitively; normalized to ERROR level |
| `ERROR_CODE: ...` | ERROR | 1 | errors.go | Part of `printAIError` structured error block |
| `CAUSE: ...` | ERROR | 1 | errors.go | Part of `printAIError` structured error block |
| `HINT: ...` | ERROR | 2 | errors.go, errors.go | Part of structured error/warning output blocks |
| `ACTION: ...` | ERROR | 1 | errors.go | Part of `printAIError` structured error block |
| `WARNING: ...` (no indent) | WARN | 22 | task/*.go, init.go, qualitygate/*.go, errors.go, submit.go | Standard warning prefix |
| `  WARNING: ...` (2-space indent) | WARN | 1 | quality_gate_lifecycle.go | Indented context message |
| `WARNING: ...` (compound, multi-line) | WARN | 1 | submit.go | Multi-line message starting with `---\nWARNING:` and containing embedded `HINT:`; classified by leading prefix |
| `AUTO-RESTORE-SKIP: ...` | WARN | 3 | submit.go | Skip conditions for auto-restore |
| `AUTO-RESTORE: ...` | INFO | 1 | submit.go | Successful restore notification |
| `SOURCE-RESOLVE: ...` | INFO | 1 | task/add.go (pkg/) | Source resolution trace |
| `NOTE: ...` | INFO | 1 | index.go | Informational notes |
| `[debug] ...` | DEBUG | 1 | base/output.go | Debug trace messages |
| `[feature:complete] Error: ...` | ERROR | 1 | feature_complete.go | Compound prefix; longest-prefix matching |
| `[feature:complete] Status ...` | INFO | 1 | feature_complete.go | Status update |
| `[feature:complete] Warning: ...` | WARN | 1 | feature_complete.go | Warning within feature completion |
| `[feature:complete] Push failed: ...` | ERROR | 1 | feature_complete.go | Push failure |
| Prefixless (progress/status) | INFO | 11 | qualitygate/*.go: quality_gate.go (3 — orchestration status messages "Running quality gate...", "Quality gate passed"), quality_gate_lifecycle.go (4 — lifecycle step progress "Checking X...", step completion messages), quality_gate_report.go (2 — report formatting status), output.go (1 — probe status message), base/output.go (1 — general progress indicator) | Progress bars, orchestration status, probe messages |
| Forensic (all prefixes) | EXCLUDED | 5 | forensic/extract.go | Forensic command's stderr is its primary output channel; excluded from migration |
| `Warning: ...` (mixed case) | WARN | 1 | pkg/task/state.go | Mixed-case prefix; caught by case-insensitive matching |
| `FAIL: ...` / `OK: ...` | WARN / INFO | 2 | pkg/serverprobe/serverprobe.go | Probe result messages |
| `ERROR: ...` (pkg/) | ERROR | 1 | pkg/just/just.go | Just execution error |
| `WARNING: ...` (pkg/) | WARN | 3 | pkg/just/just.go, pkg/testrunner/testrunner.go, pkg/serverprobe/serverprobe.go | Package-layer warnings (includes 1 with 2-space indent in serverprobe) |

**Matching priority**: Prefix matching uses **longest-prefix-first** ordering. When a message has a compound prefix (e.g., `[feature:complete] Error:`), the entire compound prefix is matched first. The matching order is: (1) `[feature:complete] Error:` / `[feature:complete] Warning:` / `[feature:complete] Push failed:` → (2) `AUTO-RESTORE-SKIP:` → (3) `AUTO-RESTORE:` → (4) `SOURCE-RESOLVE:` → (5) `ERROR:` / `ERROR_CODE:` / `CAUSE:` / `ACTION:` / `HINT:` → (6) `WARNING:` → (7) `NOTE:` → (8) `[debug]` → (9) `FAIL:` / `OK:` → (10) default INFO.

**Fallback rule for uncategorized messages**: Any stderr call site that (a) does not belong to the `forensic` command, and (b) lacks a matching prefix from the table above, is assigned INFO level. This ensures forward compatibility — new messages that do not match any prefix are still logged rather than silently dropped.

**Forensic exclusion**: The `forensic` command writes structured diagnostic output to stderr by design. All 5 call sites in `forensic/extract.go` are explicitly excluded from the migration scope. Migrating them would produce duplicate output and break the command's output contract with consuming tools.

**Case-insensitive prefix matching**: Prefix matching uses `strings.HasPrefix(strings.ToUpper(strings.TrimSpace(msg)), prefix)` to handle both `ERROR:` and `error:`, as well as mixed-case `Warning:`, uniformly. Leading whitespace is stripped before matching to handle indented messages (`  ERROR:`, `  WARNING:`). No source-code normalization is required.

**Note on prefix parsing scope**: The prefix-based classification scheme is a migration tool only. After migration, all existing call sites use explicit level-based `forgelog` API calls. New code calls `forgelog.Warn()`, `forgelog.Error()`, etc. directly — no prefix convention needed. The prefix matching logic exists solely in the migration layer and does not impose an ongoing maintenance burden.

# Alternatives

## A. Do nothing

Keep stderr-only output. Zero implementation cost but zero improvement in debuggability. Every future incident remains a code archaeography exercise.

## B. Environment variable toggle (FORGE_LOG_FILE)

Use `FORGE_LOG_FILE=1 forge task claim` to opt into file logging. Simpler config but requires users to remember the env var. Doesn't integrate with the existing `.forge/config.yaml` system. No auto-cleanup.

## C. Structured JSON logging

Write JSON-formatted logs for machine parseability. Over-engineered for the current need (human diagnosis). Can be added later as a format option if needed.

**Recommendation**: Proceed with the proposed solution. It balances debuggability with minimal complexity, integrates with existing config, and provides auto-cleanup.

## D. Use Go standard library `log/slog`

Go 1.21+ provides `log/slog` with structured logging, leveled output, and handler customization. However, slog is designed for structured key-value logging (e.g., `slog.Info("task restored", "id", taskID, "status", "pending")`). This proposal needs simpler human-readable diagnostic logs with a fixed timestamp+level+message format, dual output to both stderr and a per-invocation file, and auto-cleanup of old log files. slog's `Handler` interface could theoretically be implemented for this, but the overhead is unjustified: the `forgelog` package is ~150 lines of straightforward Go that handles file creation, per-write `O_APPEND` output, and cleanup in one place. Using slog would require a custom `slog.Handler` implementation for the dual-output pattern (stderr + file with different formatting for each), plus the same file-management code, adding complexity without benefit. The specific dual-output requirement — stderr-first-then-file with prefix-preserved messages on stderr and timestamp+level+message in the file — does not map cleanly to slog's `Handler.Handle()` method, which produces a single formatted record. Implementing this would require either two handler instances (one for stderr, one for file) with careful coordination, or a single handler that internally manages two writers — either way, more code and more complexity than the direct `forgelog` approach. **TCO acknowledgment**: slog is maintained by the Go team indefinitely, while forgelog is maintained by the forge team. This is a deliberate trade-off: forgelog's narrow scope (~150 lines, no dependencies, printf-style API) minimizes the maintenance surface. If structured/machine-parseable logging is needed in the future, slog can be adopted then — the `forgelog` API is intentionally narrow and can be replaced without affecting call sites.

# Scope

## In Scope

- `pkg/forgelog` package: log level enum, file writer with auto-cleanup, dual-output (file + stderr)
- **`forgelog` public API**:
  ```go
  package forgelog

  // Init initializes the logging layer. Creates .forge/logs/ on demand.
  // Falls back to stderr-only mode if directory creation fails.
  // Checks FORGE_NO_LOG=1 env var — if set, skips file logging entirely.
  func Init(config *forgeconfig.LogsConfig, logsDir string) error

  // Debug/Info/Warn/Error write to both stderr and the log file (if active).
  // Format: 2006-01-02T15:04:05.000 [LEVEL] message\n
  // These are printf-style: Warn("task %s not found", id)
  // All functions are safe for concurrent use. The underlying file handle
  // uses a sync.Mutex to serialize writes.
  func Debug(msg string, args ...interface{})
  func Info(msg string, args ...interface{})
  func Warn(msg string, args ...interface{})
  func Error(msg string, args ...interface{})

  // Close closes the log file handle.
  // Call via defer in each command's runE function.
  // Not strictly necessary with O_APPEND per-write writes (no buffer to flush),
  // but ensures clean file handle release.
  func Close()
  ```
- `.forge/config.yaml` `logs` section: `level` (debug/info/warn/error), `retentionDays` (default 7)
- **Config struct extension**: Add a `Logs` field to `forgeconfig.Config`:
  ```go
  type LogsConfig struct {
      Level          string `yaml:"level"`           // default: "info"
      RetentionDays  int    `yaml:"retentionDays"`   // default: 7
  }
  // Added to existing Config struct:
  // Logs *LogsConfig `yaml:"logs,omitempty"`
  ```
  The `omitempty` YAML tag ensures existing config files without a `logs` section deserialize cleanly — `Logs` will be `nil`, and defaults are applied in `forgelog.Init()`. No config migration script is needed.
- **Config validation**: `forgelog.Init()` validates the config and applies safe defaults for invalid values. Level must be one of `debug`, `info`, `warn`, `error` (case-insensitive); any unrecognized value falls back to `info`. `retentionDays` must be >= 1; values of 0 or negative fall back to the default of 7. This prevents edge cases like `retentionDays: 0` from deleting the active log file.
- Add `ForgeLogsDir` constant to `pkg/feature/constants.go`
- Migrate existing stderr calls in `internal/` to use `forgelog.Warn()`, `forgelog.Error()`, etc.
- **Migrate existing stderr calls in `pkg/`** (8 call sites in serverprobe, just, task, testrunner) to use `forgelog` API. The `pkg/` layer is included in scope because these packages are internal to the CLI (not a public library) and their stderr output is equally ephemeral.
- `.forge/logs/` entry in `gitignoreEntries` (init.go)
- `forge init` ensures `.forge/logs/` directory exists
- `forgelog.Init()` called early in each command's `runE` function
- **Bootstrap safety**: `forgelog.Init()` creates `.forge/logs/` on demand via `os.MkdirAll(logDir, 0700)` (owner-only directory). Log files are created with mode `0600` (owner-only read/write) to prevent exposing potentially sensitive diagnostic content (file paths, task content, config values) on shared systems. If directory creation fails (e.g., read-only filesystem, permission denied), `Init()` falls back to stderr-only mode — logging functions become no-ops for the file writer but stderr output continues unchanged. This resolves the paradox where `forgelog.Init()` is called before `forge init` has run.
- **Migration strategy**: Migrate all 80 call sites in a single PR to avoid partial-migration state. The PR is mechanically verifiable: `grep -r 'fmt.Fprintf(os.Stderr' forge-cli/internal/ forge-cli/pkg/ --include='*.go' | grep -v testdata | grep -v 'forensic/'` must return 5 results (the excluded forensic call sites) after migration. Each migrated call site is a one-line change (`fmt.Fprintf(os.Stderr, "WARNING: ...")` → `forgelog.Warn("WARNING: ...")`) with no behavioral change to stderr output. **Note on alternative stderr patterns**: The 80-count is based on `fmt.Fprintf(os.Stderr, ...)` which is the dominant pattern in the codebase. A secondary sweep will check for `fmt.Fprintln(os.Stderr, ...)`, `os.Stderr.WriteString(...)`, and `log` package usage (`log.Printf`, `log.Println` with default stderr output). If any such call sites are found, they are included in the same migration PR using the same classification rules.

## Out of Scope

- Structured/JSON log format
- Log rotation (one file per invocation is sufficient)
- Log aggregation or remote shipping
- Migrating test code (tests continue using stderr directly)
- Changes to the plugin (agents/commands/skills) — this is CLI-only
- **Migrating the `forensic` command**: `forensic` uses stderr as its primary structured output channel for session analysis. Its 5 call sites in `forensic/extract.go` are excluded from migration to preserve output contract with consuming tools. Migrating them would produce duplicate output and break existing pipelines.
- **CLI log viewer command** (`forge log` / `forge logs`): A dedicated command to list, search, and filter log files from `.forge/logs/` is not included. Users can read logs directly with standard tools (`cat .forge/logs/2026-06-04T*`, `grep`, etc.). A log viewer command is a natural follow-up enhancement once the logging infrastructure is in place and usage patterns are better understood.

# Risks

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| Log file contention under concurrent commands | Low — each invocation creates a unique timestamped file | Use per-invocation filename; no shared file |
| Disk accumulation if retention too long | Low — default 7 days, configurable | Auto-cleanup on each command startup. **Disk budget**: typical invocation ~5-10KB, `run-tasks` ~50-100KB. 7-day retention at 100 invocations/day ≈ 7-70MB worst case. |
| Config parsing failure blocks logging | Low — defaults applied when config missing | Hardcoded defaults: level=info, retention=7 days. Invalid values (unrecognized level, retentionDays < 1) fall back to defaults. |
| Performance impact of file I/O on every log call | Low — per-write `O_APPEND` is efficient | Each `forgelog` call writes directly to the file via `O_APPEND` (no bufio layer). `O_APPEND` provides atomic appends under `PIPE_BUF` on most OS, so no data is lost on unclean exit. Per-write syscall overhead is acceptable for typical volumes (~50-200 log lines per invocation). `run-tasks` auto-loops produce ~10-50KB per invocation — well within performance budget. |
| Sensitive information in log files | Medium — ERROR/WARNING messages may include file paths, task content, config values | Log files are created with mode `0600` (owner-only), directories with `0700`. `.forge/logs/` is added to `.gitignore` by `forge init`. Document that `.forge/` directory must not be committed. Do not redact at log time — the diagnostic value of full messages outweighs the risk for local-only files. |
| Logging layer regression causes command failures | Low — one-line mechanical changes per call site | **Emergency disable**: `FORGE_NO_LOG=1` env var skips all file logging, reverting to stderr-only behavior. **Rollback**: revert the migration PR; no schema or config changes require cleanup. |
| Data loss on unclean exit | None — resolved by design | Each `forgelog` call writes directly via `O_APPEND` with no buffering. Messages are persisted to disk before the function returns. `os.Exit`, panic, and signal kills cannot lose log data because there is no buffer to flush. |

# Success Criteria

| ID | Criterion | Verification |
|----|-----------|-------------|
| SC-1 | `forge task submit` writes AUTO-RESTORE diagnostic to `.forge/logs/<datetime>-<pid>.log` | Run submit with fix-task scenario; verify log file exists with expected content |
| SC-2 | Log level filtering works: setting `level: warn` suppresses INFO messages | Set config to warn; run command; verify only WARN/ERROR in log |
| SC-3 | Auto-cleanup deletes files older than `retentionDays` | Create log file with old timestamp; run any forge command; verify deletion |
| SC-4 | `forge init` creates `.forge/logs/` directory and adds `.forge/logs/` to `.gitignore` | Fresh project; run `forge init`; verify directory and gitignore entry |
| SC-5 | Dual output: same message appears in both stderr and log file | Run command; compare stderr output with log file content |
| SC-6 | Config missing or malformed falls back to defaults (info level, 7 days) | Delete config section; run command; verify logging still works with defaults |
| SC-7 | Migration completeness: all non-excluded call sites use forgelog API | `grep -r 'fmt.Fprintf(os.Stderr' forge-cli/internal/ forge-cli/pkg/ --include='*.go' \| grep -v testdata \| grep -v 'forensic/'` returns 0 results |
| SC-8 | Concurrent commands produce separate log files | Run two `forge task claim` commands simultaneously; verify two distinct `.log` files with different PIDs |
| SC-9 | Emergency disable: `FORGE_NO_LOG=1` suppresses file logging without affecting stderr | Set env var; run command; verify no `.forge/logs/` file created; stderr output unchanged |
