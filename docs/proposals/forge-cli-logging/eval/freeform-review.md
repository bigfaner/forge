# Freeform Expert Review: Forge CLI Structured Logging (Revision)

**Reviewer**: CLI Diagnostics & Logging Architect
**Document**: `docs/proposals/forge-cli-logging/proposal.md`
**Date**: 2026-06-04

---

## Section 1: Background Assessment

This proposal adds a per-invocation file-based diagnostic logging layer to the Forge CLI. The motivating incident is concrete and persuasive: an `autoRestoreSourceTask` silent failure that required hours of code archaeography instead of a simple log lookup. This is exactly the kind of operational pain that justifies infrastructure investment.

The architecture is a backend-abstraction pattern: a `forgelog` package with `ConsoleBackend` (preserving stderr byte-for-byte) and `FileBackend` (adding timestamp+level prefix), dispatched to independently. Each `forge` invocation creates its own log file under `.forge/logs/` named `<ISO-8601-datetime>-<pid>.log`. Levels are file-only; console always outputs everything. Configuration lives in `.forge/config.yaml` under a new `logs` section with `omitempty` for backward compatibility. Auto-cleanup runs on startup after the new file is opened. An emergency disable via `FORGE_NO_LOG=1` is provided.

The revision has clearly benefited from prior review feedback. The filename now includes a PID suffix to prevent concurrent-invocation collisions. The call-site categorization table is exhaustive, covering compound prefixes, case-insensitive matching, indented variants, prefixless messages, and explicit forensic exclusion. The per-line format (`2006-01-02T15:04:05.000 [LEVEL] message`) is now specified. The `forgelog.Init()` failure semantics (auto-create directory, fallback to console-only) are explicit. The cleanup ordering (after new file open) is stated. These are all meaningful improvements.

The proposal's strength is its architectural restraint. It resists the temptation to adopt `log/slog` (correctly identified as a complexity mismatch for human-readable diagnostics) and instead builds a minimal, purpose-built logging layer (~150 lines). The zero-change console contract is the key design insight -- by treating console as a first-class backend that outputs raw messages, backward compatibility is guaranteed by construction.

Where the proposal still has gaps is in edge cases around the migration process itself, a subtle but significant tension in its data-safety claims, and several operational scenarios that the design doesn't fully address. The following sections detail these concerns.

---

## Section 2: Key Risk Identification

风险：

The proposal makes an explicit data-safety guarantee that creates a significant architectural tension. It states:

> "No buffering — each write is persisted before function returns. `os.Exit`, panic, and signal kills cannot lose data."

This guarantee is simultaneously the proposal's strongest safety claim and its most problematic constraint. Without buffering, every `forgelog.Warn()` call incurs a full syscall round-trip to the OS. Under the stated volumes (~50-200 lines per invocation), this is acceptable for interactive use. But the proposal explicitly targets the `run-tasks` autonomous dispatcher as a primary beneficiary, and a quality gate run can produce substantially more than 200 log lines across test orchestration, server probe output, and feature completion lifecycle messages.

More critically, this no-buffering claim is only partially true even on its own terms. `os.File.Write` with `O_APPEND` writes to the kernel buffer cache. It does NOT guarantee persistence to stable storage unless followed by `fsync`. A process killed by `SIGKILL` or a power loss can lose data that has been "written" via `Write` but not flushed to disk. The proposal's claim that "signal kills cannot lose data" is inaccurate for `SIGKILL` (which cannot be caught) unless every write is followed by `fsync`, which would be catastrophically slow.

This matters because the data-safety claim is used to justify the rejection of buffering. If the claim is weaker than stated, the cost-benefit analysis changes: a `bufio.Writer` with deferred flush would lose at most the last few lines on `SIGKILL`/power-loss (same practical risk as the current design), while substantially improving write throughput for high-volume invocations.

---

问题：

The `Fprintln` handling description contains a subtle behavioral trap:

> "`fmt.Fprintln(os.Stderr, msg)` is equivalent to `fmt.Fprintf(os.Stderr, msg + "\n")`. Migration treats them identically — `forgelog.Warn(msg)` handles the trailing newline."

This equivalence is not exact. `fmt.Fprintln` adds `\n` unconditionally, while `fmt.Fprintf(os.Stderr, msg + "\n")` adds it only if the caller remembered to include it. More importantly, the migration statement "forgelog.Warn(msg) handles the trailing newline" is ambiguous: does `forgelog.Warn` add a trailing newline to the message, or does it expect the caller to include one? The API signature `Warn(format string, args ...interface{})` suggests printf-style formatting, but the proposal never specifies whether the forgelog functions append `\n` internally or expect the caller to include it.

If `forgelog.Warn` appends `\n` (like `log.Println`), then migrating `fmt.Fprintf(os.Stderr, "WARNING: task %s not found\n", id)` to `forgelog.Warn("WARNING: task %s not found\n", id)` would produce a double newline. If it does NOT append `\n`, then migrating `fmt.Fprintln(os.Stderr, "WARNING: task")` to `forgelog.Warn("WARNING: task")` would lose the trailing newline in the file backend (since FileBackend adds `timestamp [LEVEL]` prefix but no trailing newline).

This matters because it directly impacts the zero-change console contract. ConsoleBackend must output the message exactly as it would have appeared via the original `fmt.Fprintf(os.Stderr, ...)`. If the newline handling is wrong, console output changes, and SC-2 (byte-identical stderr) fails.

---

问题：

The categorization table's "Prefixless (progress/status)" entry maps all prefixless messages to INFO:

> "Prefixless (progress/status) | INFO | qualitygate/*.go, output.go, base/output.go | Progress bars, orchestration status, probe messages"

The fallback rule reinforces this:

> "Fallback: Any message without a matching prefix defaults to INFO level."

But the proposal also states that prefix parsing is only for migrating existing call sites:

> "Prefix parsing is **only for migrating existing call sites** — new code calls `forgelog.Warn()` directly and needs no prefix convention."

This creates a contradiction: the migration is supposed to be a one-line mechanical change from `fmt.Fprintf(os.Stderr, msg)` to `forgelog.Info(msg)`, but for prefixless messages, the implementer must first determine whether each message is truly INFO or should be something else. This is not mechanical -- it requires judgment per call site. Progress bars that indicate normal operation are INFO, but what about a prefixless message like `"Failed to start dev server"` that currently exists in the codebase? Without an exhaustive per-call-site mapping, the implementer has no authoritative reference for these edge cases.

This matters because the proposal sells the migration as "one-line mechanical" multiple times. If a significant fraction of the ~102 call sites require judgment calls (not just prefix matching), the migration is more error-prone and time-consuming than advertised.

---

风险：

The `HINT:` prefix is mapped to ERROR level:

> "`HINT: ...` | ERROR | errors.go | Part of structured error/warning blocks"

But a hint is typically a suggestion for remediation, not an error itself. In the actual codebase, `HINT:` messages appear alongside `ERROR:` and `WARNING:` blocks as follow-up guidance. If `HINT:` is logged at ERROR level, then `forgelog.Error("HINT: run 'forge init' to fix this")` writes `[ERROR] HINT: run 'forge init' to fix this` to the log file, which is semantically misleading. A future developer grepping for `[ERROR]` in log files would find remediation suggestions mixed with actual errors, reducing the diagnostic value of the log.

This matters because log level semantics are a long-term contract. Once log files accumulate, search patterns and monitoring tools will be built around the assumption that `[ERROR]` means something went wrong. Including helpful hints at ERROR level pollutes this signal.

---

风险：

The proposal specifies that `forgelog.Init()` is "called early in each command's `runE` function" but does not address the interaction with Cobra's command initialization lifecycle. In a typical Cobra CLI, `PersistentPreRunE` is used for cross-cutting concerns like configuration loading. If `forgelog.Init()` requires a loaded config (it needs `forgeconfig.LogsConfig`), it must be called after config loading. But if config loading itself emits diagnostic messages (e.g., "WARNING: unknown field in config"), those messages would fire before `forgelog` is initialized and would go only to raw stderr, not to the log file.

> "`forgelog.Init()` called early in each command's `runE` function; `defer forgelog.Close()` follows"

This placement means that any diagnostic messages emitted during config loading, cobra command setup, or flag parsing are lost from the log file. For the motivating scenario (diagnosing mysterious behavior), these early messages can be critical -- a misconfigured flag or a config parsing issue could be exactly the root cause being investigated.

This matters because the proposal's value proposition is "never lose diagnostic output again." Messages emitted before `forgelog.Init()` are still lost, and there's no clear boundary for what constitutes "early enough."

---

问题：

The proposal defines a `Close()` function but does not specify its failure semantics:

> "`Close()` closes all backends. Call via `defer` in each command's `runE`."

When `Close()` is called via `defer`, it runs in the deferred-function chain. If the command exits via `log.Fatal()` (which calls `os.Exit(1)`), deferred functions do NOT run in Go -- `os.Exit` terminates immediately. The proposal currently has 1 existing `log.Printf` call site. If that call site or any future code path calls `log.Fatal()`, the `forgelog.Close()` defer is bypassed, and the log file is not properly closed.

While an unclosed `O_APPEND` file will eventually be flushed by the OS on process exit, this is an implementation detail of the Go runtime and the OS, not a guarantee the proposal should rely on for its data-safety claims. Additionally, the cleanup routine in `Close()` (if any final flush or metadata write is added later) would be skipped.

This matters because the proposal emphasizes data safety as a key property, and the most common way to lose data is exactly when something goes wrong enough to call `os.Exit` or `log.Fatal`.

---

风险：

The backend registration order is specified as "ConsoleBackend first, FileBackend second":

> "ConsoleBackend first, FileBackend second. If file write fails, console has already written."

But the Backend interface's `Write` method signature is `Write(level LogLevel, timestamp time.Time, msg string)`, and the proposal shows sequential dispatch to all registered backends. This means a slow or blocked file write (e.g., NFS stall, disk I/O wait on a spinning drive) would stall the console write for that message. The proposal's guarantee that "console output continues uninterrupted" is only true if file writes never block -- an assumption that does not hold on all filesystems.

For a CLI tool that runs interactively (quality gate progress bars, task status updates), a single stalled file write would freeze the user-visible output. This is the exact opposite of what the proposal intends -- the logging layer should never degrade the interactive experience.

This matters because the backend dispatch is synchronous and sequential. There is no timeout, no async write, and no fallback if the file backend becomes slow (as opposed to failing outright).

---

问题：

The proposal specifies case-insensitive prefix matching:

> "Case-insensitive: `strings.HasPrefix(strings.ToUpper(strings.TrimSpace(msg)), prefix)`. Leading whitespace stripped before matching."

But the categorization table lists `ERROR:` and `error:` as separate rows with different source areas:

> "`ERROR: ...` (no indent) | ERROR | init.go, init_config.go, quality_gate.go, errors.go, etc."
> "`error: ...` (lowercase) | ERROR | upgrade.go | Case-insensitive match"

If the matching is truly case-insensitive, these two rows are redundant -- they both match `ERROR:` after `strings.ToUpper`. Listing them separately suggests the implementer should treat them differently, but the matching algorithm makes no such distinction. This is a presentation inconsistency that could confuse during implementation. The `error:` row should either be removed (as it's subsumed by the case-insensitive `ERROR:` rule) or the table should note that these are all unified by the case-insensitive matching.

---

风险：

The proposal mentions log files containing potentially sensitive information:

> "Sensitive info in log files | Medium | Medium | File mode `0600`, dir `0700`, `.gitignore` entry. No redaction at log time -- diagnostic value > risk for local files"

The assessment that "diagnostic value > risk for local files" is reasonable for single-developer machines but does not account for shared development environments, containerized CI runners where log files might be collected as artifacts, or paired-programming scenarios. More importantly, the proposal's migration preserves the original message text verbatim. This means any error message that includes user-provided content (task names, file paths, configuration values) will be persisted to disk in a file that is never redacted and only deleted after 7 days.

The `0600` permission only protects against other users on the same machine. In container environments, files are often owned by root or a shared user. In CI, log files could be uploaded as build artifacts with broader visibility.

This matters because the proposal makes a deliberate decision to not redact, but the risk assessment considers only local single-user scenarios. The actual deployment surface is wider.

---

## Section 3: Improvement Suggestions

建议：

**Resolve the newline handling ambiguity explicitly in the API contract.** Specify that `forgelog.Warn(format, args...)` does NOT append a trailing newline -- the message is output exactly as `fmt.Fprintf(os.Stderr, format, args...)` would produce it. For the file backend, the format is `timestamp [LEVEL] formatted_msg` where `formatted_msg` already contains the caller's `\n` if any. Document this in a code comment on the public API functions. Migration rule: `fmt.Fprintf(os.Stderr, "WARNING: %s\n", x)` becomes `forgelog.Warn("WARNING: %s\n", x)` (newline preserved), and `fmt.Fprintln(os.Stderr, "WARNING: "+x)` becomes `forgelog.Warn("WARNING: %s\n", x)` (newline added). This addresses the Fprintln behavioral trap by making the newline handling rule unambiguous for implementers.

---

建议：

**Soften the data-safety claim and consider optional bufio for high-volume paths.** Replace the absolute claim "each write is persisted before function returns" with the more accurate "each write is issued to the OS before the function returns; persistence to stable storage depends on OS behavior." Add a `bufio.Writer` wrapper as an optional optimization, enabled by default, with an explicit `Flush()` in `Close()`. Document the trade-off: in normal operation, all lines are flushed on `Close()`; on `SIGKILL` or power loss, the last few lines may be lost regardless of buffering strategy. This addresses the architectural tension between the overstated data-safety claim and the practical performance characteristics.

---

建议：

**Reclassify `HINT:` from ERROR to INFO level.** Hints are remediation suggestions, not error conditions. The categorization should be:

- `ERROR_CODE:`, `CAUSE:`, `ACTION:` remain ERROR (they describe error context)
- `HINT:` maps to INFO (it describes a suggested action)

This preserves the semantic integrity of the ERROR level for log-file consumers and ensures that `grep '[ERROR]' *.log` returns only things that went wrong, not suggestions for how to fix them.

---

建议：

**Address pre-Init message capture.** Add a note in the design that acknowledges the gap: messages emitted before `forgelog.Init()` (during config loading, flag parsing) are not captured by the log file. If this gap is acceptable (config loading failures are rare and typically visible on console), state so explicitly. If not, consider a two-phase init: `forgelog.InitConsole()` called before config loading (captures to console only), then `forgelog.InitFile(config)` called after config is loaded (adds FileBackend retroactively). This addresses the early-message blind spot.

---

建议：

**Add a write timeout or async fallback for the FileBackend.** Even a simple 100ms deadline on the file write would prevent an NFS stall from freezing interactive output. The implementation could use a goroutine with a `select`/`time.After` pattern, or simply log a warning to console ("forgelog: file write stalled, degrading to console-only") and disable the FileBackend for the remainder of the invocation. This addresses the synchronous-stall risk without adding significant complexity.

---

建议：

**Consolidate the categorization table to remove redundancy from case-insensitive matching.** Since the matching algorithm applies `strings.ToUpper` uniformly, the `error:` row in upgrade.go should not appear as a separate entry. Instead, add a footnote to the `ERROR:` row: "Also covers lowercase `error:` variants in upgrade.go (matched via case-insensitive comparison)." This reduces confusion during implementation about whether these are truly separate categories or presentation artifacts.

---

建议：

**Expand the sensitive-data risk assessment to cover CI and container scenarios.** Add a note under the Key Risks table's "Sensitive info in log files" entry acknowledging that in CI environments, log files may be collected as artifacts with broader visibility. Recommend that CI configurations using Forge should set `FORGE_NO_LOG=1` or configure log output to a secured directory. This does not change the design but sets accurate expectations about the trust boundary of log file contents.

---

建议：

**Handle the `log.Fatal` / `os.Exit` path explicitly.** Add to the design that the existing `log.Printf` call site should be migrated to `forgelog` (replacing `log` package usage entirely), and that no code path in the CLI should use `log.Fatal()`. If `log.Fatal` must be preserved somewhere, add a `forgelog.Flush()` function that can be called before `os.Exit` in those specific paths, and document that `defer forgelog.Close()` does not run on `os.Exit`. This addresses the data-loss-on-fatal-exit risk.

---

The proposal after adopting these suggestions maintains its core architecture (backend abstraction, per-invocation files, zero-change console) but with sharper edge handling: correct data-safety claims, unambiguous newline semantics, improved level classification, protection against file-write stalls, and a more complete operational risk picture. The design is solid at its center; these changes address the perimeter where real-world friction lives.
