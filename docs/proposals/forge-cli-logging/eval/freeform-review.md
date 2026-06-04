# Freeform Expert Review: Forge CLI Structured Logging

**Reviewer**: CLI Diagnostics & Logging Architect
**Document**: `docs/proposals/forge-cli-logging/proposal.md`
**Date**: 2026-06-04

---

## Section 1: Background Assessment

This proposal addresses a real and costly gap in the Forge CLI's observability story. The problem statement is grounded in a concrete incident — the `autoRestoreSourceTask` silent-failure episode — where hours of code archaeography were needed to speculate about root cause, when a single persisted log line would have resolved it immediately. That's the kind of war story that justifies infrastructure work, and the proposal deserves credit for anchoring the motivation in an actual event rather than an abstract "wouldn't it be nice."

The core technical approach is sensible and appropriately scoped: a per-invocation log file under `.forge/logs/`, level-filtered writes via a new `pkg/forgelog` package, simultaneous output to both file and stderr (the "dual output" guarantee), and auto-cleanup on startup. This is the right shape for a first-pass logging system in a CLI tool — file-based, zero-dependency, human-readable, with a clear path to more structured formats later.

The proposal rests on several assumptions that are mostly sound but worth surfacing explicitly: (1) that the `.forge/` directory always exists and is writable at the point where logging needs to start, (2) that the 64 `fmt.Fprintf(os.Stderr, ...)` call sites can be cleanly categorized into four levels based on their message prefix, (3) that one file per invocation is sufficient (no need for rotation within a long-running process), and (4) that the existing config infrastructure (`forgeconfig.Config` struct, `.forge/config.yaml` parsing) can absorb a new `logs` section without breaking backward compatibility.

On the surface, the proposal reads as a well-bounded, low-risk enhancement. The alternatives analysis is honest — it correctly dismisses structured JSON logging as over-engineered for the current need and correctly notes that an environment-variable toggle lacks integration with the existing config system. The risk table identifies the obvious concerns (contention, disk accumulation, config failure, I/O performance).

Where the proposal begins to fray is in the details of its categorization scheme, the precision of its call-site accounting, and its silence on several failure modes that matter in practice. The following sections dig into these concerns.

---

## Section 2: Key Risk Identification

风险：

The proposal states "64 `fmt.Fprintf(os.Stderr, ...)` call sites across the CLI" but the actual codebase count is 72 call sites in `forge-cli/internal/` alone (96 total including test code, which the proposal correctly scopes out). This 8-call-site discrepancy suggests the proposal was written against an older snapshot of the code. More critically, the discrepancy means the migration plan is under-counting its workload from the start — any implementation plan derived from the "64" figure will discover unmigrated call sites late in the process.

> "Evidence: 64 `fmt.Fprintf(os.Stderr, ...)` call sites across the CLI, zero persisted diagnostics."

This matters because accurate counting is the foundation of migration completeness. If the count is wrong, the confidence that "all call sites are migrated" is undermined.

---

问题：

The categorization scheme in the "Categorized output" section does not account for a significant number of stderr calls that lack any recognizable prefix. Analysis of the actual codebase reveals approximately 15 call sites in `forge-cli/internal/` that use neither `ERROR:`, `WARNING:`, `error:`, `NOTE:`, `AUTO-RESTORE`, `HINT:`, `CAUSE:`, `ACTION:`, `ERROR_CODE:`, `[debug]`, nor `[feature:` prefixes. These are primarily in `qualitygate/quality_gate_lifecycle.go` (progress messages like `"  Starting dev server (just %s)..."`, `"  Running tests (just %s)..."`) and `forensic/extract.go` (formatted diagnostic output). The proposal's categorization table lists only four prefix patterns:

> "- `ERROR:`, `AUTO-RESTORE-SKIP:` -> ERROR"
> "- `WARNING:`, `AUTO-RESTORE-SKIP:` (degraded) -> WARN"
> "- `AUTO-RESTORE:`, `SOURCE-RESOLVE:`, `NOTE:` -> INFO"
> "- `[debug]` -> DEBUG"

This leaves at least three categories of calls unmapped: (a) progress/status messages in qualitygate lifecycle (e.g., `"=== All tasks completed for feature: %s ==="`), (b) lowercase `"error:"` prefix calls in `upgrade.go` (5 call sites), and (c) the entire forensic output subsystem (5 call sites). Without a mapping for these, the migration will either skip them (leaving some stderr-only, breaking the dual-output promise) or force ad-hoc decisions during implementation.

This matters because the categorization scheme is the linchpin of the entire design — it bridges the existing stderr-only world to the new leveled world. If the bridge has gaps, messages fall through.

---

问题：

The proposal lists `AUTO-RESTORE-SKIP:` under two different levels:

> "- `ERROR:`, `AUTO-RESTORE-SKIP:` -> ERROR"
> "- `WARNING:`, `AUTO-RESTORE-SKIP:` (degraded) -> WARN"

But the actual codebase contains three `AUTO-RESTORE-SKIP` call sites in `task/submit.go`, and none of them carry a "(degraded)" qualifier in the message text. They are:

- `"AUTO-RESTORE-SKIP: source task %s not found in index"`
- `"AUTO-RESTORE-SKIP: source task %s is %s (not blocked)"`
- `"AUTO-RESTORE-SKIP: source task %s has unmet deps: %v"`

These are all condition-based skips in the same function — there is no runtime mechanism to distinguish "ERROR-level skip" from "WARN-level skip" based on the message prefix alone, since the prefix is identical across all three. The proposal implies a categorization that cannot be mechanically applied without additional logic beyond prefix matching.

This matters because if the categorization rule cannot be implemented as described, the implementer will either need to introduce sub-prefixes (changing stderr output, risking backward compatibility) or pick one level for all `AUTO-RESTORE-SKIP` messages (losing the nuance the proposal intends).

---

风险：

The proposal does not address the `forensic` command's stderr output pattern. The `forensic/extract.go` file contains 5 `fmt.Fprintf(os.Stderr, ...)` calls that produce formatted diagnostic output (session summaries, tool timing breakdowns, top-N lists). These calls serve a fundamentally different purpose from the diagnostic/warning/error messages elsewhere — they are the primary output of the `forensic` command, not side-channel diagnostics. Including them in the logging system would create duplicate output and confuse the forensic command's designed output contract.

> "Migrate existing stderr calls to use `forgelog.Warn()`, `forgelog.Error()`, etc."

The blanket "migrate existing stderr calls" directive, without an explicit carve-out for commands where stderr IS the output, risks over-migration. The forensic command should be excluded from the migration scope, or the proposal needs a principle for distinguishing "diagnostic stderr" from "intentional output stderr."

---

风险：

The `upgrade.go` file uses lowercase `"error:"` prefix (5 call sites), which does not match any pattern in the categorization table:

> `fmt.Fprintf(os.Stderr, "error: failed to fetch latest release: %v\n", err)`
> `fmt.Fprintf(os.Stderr, "error: failed to parse release info: %v\n", err)`

The categorization table only recognizes uppercase `"ERROR:"`. If the categorization is implemented as case-sensitive prefix matching (the natural implementation), these 5 calls will fall into an unmapped category. If implemented as case-insensitive, the proposal should say so explicitly. This matters because `upgrade.go` is a user-facing command path — errors in version upgrade are high-stakes, and losing them from the log would be exactly the kind of gap the proposal aims to prevent.

---

风险：

The filename scheme `.forge/logs/<ISO-8601-datetime>.log` (e.g., `2026-06-04T17-30-00.log`) uses second-level precision. If two `forge` commands are invoked within the same second — which is realistic in scripted or parallel CI scenarios — they will collide on the same filename. The proposal's risk table acknowledges "Log file contention under concurrent commands" but claims it has "Low" likelihood:

> "Low — each invocation creates a unique timestamped file"

Second-precision timestamps are not unique per invocation under concurrent use. This is not a theoretical concern: a CI pipeline running `forge task claim && forge task submit` in rapid succession, or a `run-tasks` dispatcher spawning parallel subagents, could easily produce two invocations in the same second.

> "Use per-invocation filename; no shared file"

The mitigation describes the intent (unique filenames) but the mechanism (ISO-8601 timestamp) does not deliver it under concurrency. The filename needs a disambiguator — either a random suffix, a PID component, or a monotonic counter — to truly guarantee uniqueness.

---

风险：

The proposal specifies that `forgelog.Init()` is "called early in each command's `runE` function," but there are 46 `RunE` entry points across the command tree. The proposal does not address how `Init()` interacts with the `forge init` command itself, which is responsible for creating the `.forge/` directory and is listed in-scope for creating `.forge/logs/`. If `forgelog.Init()` tries to write to `.forge/logs/` before `forge init` has run (e.g., on a fresh project), the logs directory will not exist and the file creation will fail.

> "`forge init` ensures `.forge/logs/` directory exists"
> "forgelog.Init() called early in each command's runE function"

This creates a bootstrapping paradox: `forge init` is itself a command with a `runE` that would call `forgelog.Init()`, which needs the `.forge/logs/` directory to exist, but that directory is created by the same `forge init` command. The proposal does not specify whether `forgelog.Init()` should create the directory on demand (in which case `forge init`'s directory creation is redundant) or fail gracefully (in which case the first invocation of any command on a fresh project silently disables logging).

---

问题：

The proposal states that dual output preserves the current stderr behavior, but does not define what "identical" means when the log file write fails:

> "Dual output: Messages write to both log file and stderr (current behavior preserved)"

If the file write fails (disk full, permissions, path too long), the proposal must guarantee that stderr output continues uninterrupted. But it does not specify the failure mode. The obvious implementation (write to file first, then stderr) would introduce a latency spike on file-write failure that could disrupt interactive use. The safe implementation (write to stderr first, then file) means stderr output is always delivered, but the proposal does not mandate this ordering. This is the kind of subtle contract that, if left unspecified, gets implemented differently by different developers and leads to inconsistent behavior across call sites.

---

风险：

The proposal claims performance impact is negligible:

> "Low — log writes are append-only, <1KB per invocation"
> "No buffering needed; `os.OpenFile` with `O_APPEND` is efficient"

This estimate of "<1KB per invocation" appears to be based on a typical interactive command run. But the `run-tasks` dispatcher runs autonomous loops, and the quality gate lifecycle runs multi-step test suites with per-step progress messages. A single `forge qualitygate run` invocation could easily produce 20+ progress messages, each with formatted output. The `forensic` command can produce extensive formatted output. The "<1KB" assumption may not hold for these high-output commands. More importantly, the real performance concern is not throughput but latency: each `forgelog.Warn()` call now involves a synchronous file write (no buffering is specified), which means every diagnostic message pays the cost of a disk flush. On network filesystems (NFS home directories in corporate environments) or encrypted disk setups, this can add measurable latency per message.

The proposal's statement "No buffering needed" is an optimization decision that should be backed by measurement, not asserted. A simple `bufio.Writer` wrapping the file handle would batch writes with negligible complexity increase.

---

问题：

The proposal does not specify the timestamp format within log file lines. The categorization section describes level mapping but not the actual log line format:

> "Each `forge` invocation writes to `.forge/logs/<ISO-8601-datetime>.log`"

This describes the filename format but not the per-line format inside the file. When diagnosing the `autoRestoreSourceTask` issue that motivated this proposal, a developer needs to know when each message was emitted relative to other events. Without timestamps in the log lines, the log file is just a sequential transcript — useful but significantly less valuable than a timestamped one. This is a surprising omission for a logging proposal.

---

风险：

The auto-cleanup mechanism runs "on each command startup" and deletes files older than `retentionDays`:

> "On each command startup, delete log files older than `retentionDays` (default 7)"

The proposal does not specify whether this cleanup happens before or after the new log file is created. If cleanup runs first (the natural order), and the current invocation's log file somehow gets a timestamp that makes it appear older than the retention threshold (clock skew, manual testing with `touch -t`), the command would delete its own active log file. This is unlikely but not impossible, especially on CI systems with clock synchronization issues. More practically, the proposal does not address what happens when cleanup fails (e.g., permission error on an old log file). Does the command proceed normally, or does the cleanup error propagate?

---

问题：

The proposal's `Config` section mentions `.forge/config.yaml` `logs` section but the existing `forgeconfig.Config` struct has no `logs` field:

> ".forge/config.yaml logs section: level (debug/info/warn/error), retentionDays (default 7)"

The current `Config` struct contains: `Version`, `ProjectType`, `Auto`, `Worktree`, `Coverage`, `TestFramework`, `Languages`, `Surfaces`, `ExecutionOrder`. Adding a `Logs` field is straightforward but the proposal does not discuss the migration path for existing config files — will `forge config set logs.level warn` work immediately, or does the config infrastructure need changes to support nested key setting? The `config set` command supports dot-notation (`auto.gitPush true`), but its implementation is tightly coupled to the known field set.

---

## Section 3: Improvement Suggestions

建议：

**Re-audit all call sites and produce a complete categorization table.** Before implementation begins, generate an exhaustive list of all 72 `fmt.Fprintf(os.Stderr, ...)` call sites in `forge-cli/internal/` and assign each to a level. This table should be included in the proposal as an appendix. For the approximately 15 call sites that lack recognizable prefixes (qualitygate progress messages, forensic output), define explicit rules: qualitygate progress messages map to INFO, forensic output is excluded from migration entirely. This addresses the first and second risks identified above — the call-site count discrepancy and the categorization gap.

---

建议：

**Resolve the `AUTO-RESTORE-SKIP` dual-level ambiguity.** Since all three `AUTO-RESTORE-SKIP` call sites share the same prefix and differ only in the reason text, the proposal should choose one level for all of them. WARN is the natural choice — these are expected condition-based skips, not errors. If finer distinction is needed, introduce explicit level hints in the message structure (e.g., `"AUTO-RESTORE-SKIP [warn]:"` vs `"AUTO-RESTORE-SKIP [error]:"`) but this changes the stderr contract and should be weighed carefully. This addresses the `AUTO-RESTORE-SKIP` categorization inconsistency.

---

建议：

**Add a disambiguator to the log filename.** Replace the pure-ISO-8601 filename with a format that includes either a PID or a short random suffix: `2026-06-04T17-30-00.12345.log` (PID) or `2026-06-04T17-30-00-a3f2.log` (4-hex-char random). This eliminates the second-precision collision risk entirely with minimal complexity. The PID approach has the advantage of being deterministic and correlatable with process listings. This addresses the filename collision risk under concurrent invocation.

---

建议：

**Specify the per-line log format with timestamps.** Define the log line format as `2006-01-02T15:04:05.000 [LEVEL] message`. The millisecond precision in the timestamp enables correlating events across parallel subagent sessions. Include this in the proposal's scope section as an explicit format contract. This addresses the missing timestamp format concern.

---

建议：

**Define `forgelog.Init()` failure semantics explicitly.** Specify that `Init()` must: (1) attempt to create `.forge/logs/` if it does not exist (using `os.MkdirAll`), (2) fall back to stderr-only mode if directory creation or file opening fails, and (3) expose a `func IsReady() bool` check that callers can use to skip unnecessary work. This resolves the bootstrapping paradox (Init creates the directory itself, no dependency on `forge init` having run first) and the dual-output failure mode (stderr always works). The `forge init` command's directory creation becomes a belt-and-suspenders guarantee rather than a prerequisite.

---

建议：

**Adopt buffered writes with flush-on-close.** Wrap the log file handle in a `bufio.Writer` and flush explicitly in a `Close()` or `Sync()` function called via `defer` in each command's `runE`. This batches multiple small writes into larger I/O operations without changing the observable behavior (log lines appear in the file after the command exits). For crash scenarios, the worst case is losing the last few lines — acceptable for a diagnostic log. This addresses the synchronous I/O performance concern without introducing the complexity of a background writer goroutine.

---

建议：

**Explicitly exclude forensic and intentional-stderr-output commands from migration.** Add a "Non-migrated commands" subsection under Scope that lists `forensic` (where stderr IS the output) and any other commands that produce structured output via stderr. These commands should continue using `fmt.Fprintf(os.Stderr, ...)` directly. The `forgelog` package's documentation should include a guideline: "Use forgelog for diagnostic messages. Use fmt.Fprintf(os.Stderr, ...) for user-facing output." This addresses the over-migration risk.

---

建议：

**Normalize the `error:` vs `ERROR:` prefix inconsistency.** The 5 call sites in `upgrade.go` use lowercase `"error:"` while the rest of the codebase uses uppercase `"ERROR:"`. The proposal should include a one-time normalization pass that aligns these to uppercase before (or as part of) the migration. This ensures the prefix-based categorization works with a single case-sensitive matching rule, and also improves consistency in the existing stderr output as a side benefit. This addresses the case-sensitivity gap in the categorization scheme.

---

建议：

**Specify cleanup ordering and error handling.** State explicitly that cleanup runs after the new log file is opened (never before), and that cleanup errors are logged to stderr but do not prevent the command from proceeding. This prevents the self-deletion edge case and ensures cleanup is best-effort without becoming a failure vector.

---

The proposal after adopting these suggestions would look substantially the same in architecture — per-invocation files, level filtering, dual output, auto-cleanup — but with sharper edges: a complete categorization table that accounts for every call site, collision-resistant filenames, explicit bootstrapping and failure semantics, buffered writes, and clear boundaries on what gets migrated versus what stays as-is. The core idea is sound; these improvements tighten the specification enough that implementation can proceed without ambiguity-driven rework.
