---
iteration: 1
title: "Adversarial Rubric Evaluation — Iteration 1"
scorer: CTO Adversary
date: 2026-06-04
---

# Iteration 1: Adversarial Rubric Evaluation

## Phase 1: Reasoning Audit

### Problem -> Solution Trace

**Problem**: stderr output is ephemeral, lost after process exit, making post-incident diagnosis impossible. A concrete incident (autoRestoreSourceTask) required hours of code archaeography.

**Solution**: Add a file-based logging layer (forgelog) with dual backend — console preserves original output byte-for-byte, file persists structured messages.

**Verdict**: The solution directly addresses the stated problem. The zero-change console contract ensures the solution does not reintroduce any behavioral regression. Sound trace.

### Solution -> Evidence Trace

**Evidence**: ~102 stderr call sites quantified by pattern type with a grep-reproducible command. One concrete incident documented.

**Verdict**: The evidence table is significantly stronger after pre-revision — it now includes Fprintln, slog, and log.Printf patterns alongside Fprintf. The counts are marked "approximate; exact verification at implementation time." This is honest. The grep command for verification is provided, which is reproducible. The categorization table for prefix-based migration is thorough, covering compound prefixes like `[feature:complete] Error:` and edge cases like 2-space indented variants.

### Evidence -> Success Criteria Trace

SCs test: persistence (SC-1), dual output fidelity (SC-2), level filtering (SC-3), auto-cleanup (SC-4), gitignore (SC-5), auto-creation (SC-6), defaults (SC-7), migration completeness via grep (SC-8), concurrent files (SC-9), emergency disable (SC-10), permissions (SC-11).

**Verdict**: The SC set covers the full lifecycle. SC-8 (grep-based migration completeness) directly addresses the "all call sites migrated" claim. SC-9 addresses concurrency. SC-10 addresses emergency disable. The pre-revision added the migration completeness SC that the prior iteration flagged as missing.

### Self-Contradiction Check — SC Consistency Deep-Dive

Cluster by affected area:

**Config cluster**: SC-3 (level filtering), SC-7 (defaults fallback), SC-10 (emergency disable).
- SC-3 requires level filtering on file. SC-7 requires defaults when config missing. SC-10 requires FORGE_NO_LOG=1 to disable. These are mutually compatible.

**File lifecycle cluster**: SC-1 (write), SC-4 (cleanup), SC-6 (auto-create), SC-11 (permissions).
- SC-4 deletes old files. SC-6 creates directory. SC-1 writes. SC-11 sets permissions. Ordering: SC-6 -> SC-1 -> SC-4 (cleanup runs after new file open per "Core Behaviors" item 3). No contradiction.

**Dual output cluster**: SC-1 (file persistence), SC-2 (console unchanged).
- SC-1 requires structured format in file. SC-2 requires byte-identical console. These are independent backends — no contradiction.

**Migration cluster**: SC-2 (console unchanged), SC-8 (no remaining stderr writes).
- SC-2 requires console unchanged. SC-8 requires all call sites migrated. If all call sites go through forgelog, and ConsoleBackend outputs raw message, both are satisfied. Compatible.

**Potential tension**: SC-10 ("no .forge/logs/ directory created" when disabled) vs SC-6 ("auto-created by Init"). The proposal resolves this: when FORGE_NO_LOG=1 is set, FileBackend is never initialized, so the directory is never created. Compatible.

**No contradictions found within the SC set.**

---

## Phase 2: Rubric Scoring with Verification Stance

### 1. Problem Definition: 98/110

**Problem stated clearly (38/40)**: The problem is unambiguous — stderr is ephemeral, diagnostics are lost post-exit. The autoRestoreSourceTask incident provides a concrete, relatable example. Minor deduction: the problem statement says "diagnostic output scrolls away in subagent sessions" which conflates two distinct issues (no persistence vs no visibility in subagent contexts). These are related but different problems.

**Evidence provided (35/40)**: The evidence table quantifies ~102 call sites by pattern and location, with a reproducible grep command. This is strong. Minor deduction: the counts are approximate. The evidence section says "~55" and "~32" for internal/ Fprintf and Fprintln respectively — the tilde prefix signals uncertainty but the table is presented as authoritative. A two-sentence note confirming "counts verified on date X via grep command Y" would strengthen this from "directional" to "verified."

**Urgency justified (25/30)**: The cost of inaction is clearly stated — "hours of code archaeography that a single log line would have resolved." This is specific and compelling. Deduction: "hours" is unquantified (2 hours? 8 hours?). A concrete time estimate would make urgency more persuasive.

### 2. Solution Clarity: 108/120

**Approach is concrete (37/40)**: The three-layer architecture (Backend, Format, API) is well-specified. The Backend interface is defined with Go code. The dual format design is captured in a table. The API layer shows exact function signatures. The migration example shows before/after code. Deduction: the `Init` function signature references `*forgeconfig.LogsConfig` and `logsDir string` but does not explain where `logsDir` is derived from — is it always `.forge/logs/` relative to CWD? What if CWD is not the project root?

**User-facing behavior described (40/45)**: The zero-change console contract is precisely stated: "byte-identical to pre-migration behavior." The file format is specified: `2006-01-02T15:04:05.000 [LEVEL] message`. Level filtering behavior is specified per-backend. The per-invocation file naming (`<ISO-8601-datetime>-<pid>.log`) is clear. Deduction: the proposal does not specify what happens when a user views the log file — is there a recommended tool, or is `cat`/`grep` expected? This is minor since the format is human-readable, but "developer experience" is listed as a domain and no DX for log consumption is described.

**Technical direction clear (31/35)**: The forgelog package structure, Backend interface, config extension with YAML tags, and migration strategy are all clear. The Fprintln handling rule is explicitly addressed. Deduction: the dispatch to backends is "sequential in registration order" — this means a slow file write blocks the next operation. The proposal acknowledges this ("a severely stalled filesystem write could block subsequent console writes") and provides FORGE_NO_LOG=1 as escape hatch, but does not consider whether async file writing is a better default.

### 3. Industry Benchmarking: 85/120

**Industry solutions referenced (28/40)**: The proposal references `log/slog` (Go stdlib 1.21+) as a concrete alternative with specific pros/cons. It also references web server access log patterns (nginx, apache). The comparison table includes four alternatives. Deduction: no reference to how specific comparable CLI tools handle this — e.g., `cargo` writes to `target/debug/` logs, `kubectl` uses `--v` verbosity flags, `docker` uses `json-file` log driver. These are directly relevant benchmarks from the CLI domain.

**At least 3 meaningful alternatives (22/30)**: Four alternatives are listed: do nothing, slog, env var toggle, backend-pattern forgelog. "Do nothing" is required. slog is a genuine alternative. Env var toggle is a lightweight variant. The set is adequate but "env var toggle" is somewhat of a straw man — it is so obviously inferior that its inclusion feels performative.

**Honest trade-off comparison (18/25)**: The slog rejection rationale is honest: "the overhead of implementing two custom handlers ... exceeds the value for a ~150-line logging package." This is a reasonable engineering judgment. Deduction: the proposal does not quantify the slog overhead. How many lines would two custom handlers require? If it is 200 lines vs 150, the "complexity mismatch" argument weakens.

**Chosen approach justified against benchmarks (17/25)**: The backend-pattern is justified as matching "problem scope precisely." The zero-change console contract is the key differentiator over slog. This is sound. Deduction: the proposal does not address that slog's Handler interface IS a backend pattern — the forgelog Backend interface is essentially a simplified slog.Handler. The proposal could acknowledge this lineage more explicitly and explain why the simplification is warranted (e.g., no structured key-value needed, no group support needed).

### 4. Requirements Completeness: 95/110

**Scenario coverage (35/40)**: The Key Scenarios section covers happy path, concurrent invocations, config missing, disk failure, emergency disable, and log cleanup. This is thorough. Deduction: no scenario for what happens when the log file grows very large within a single invocation (e.g., a long-running `run-tasks` loop producing thousands of lines). The per-invocation model means a single run could produce a multi-MB file. This is not a problem per se, but should be acknowledged.

**Non-functional requirements (35/40)**: Performance (O_APPEND, no bufio), security (0600/0700, .gitignore), concurrency (Mutex, per-file), data safety (no user-space buffering, O_APPEND atomicity) — all specified. The Windows constraint is honestly acknowledged. The data safety analysis (pre-revised) is thorough, covering normal exits, panics, os.Exit, and SIGKILL. Deduction: no explicit statement about memory overhead — each forgelog call allocates a string. For typical volumes (50-200 lines) this is negligible, but the proposal does not state this assumption.

**Constraints & dependencies (25/30)**: No external dependencies, Go 1.21+ compatibility, CLI-only (no plugin changes), Windows permission advisory. Well-specified. Deduction: the `forgeconfig.LogsConfig` dependency means the config package must be modified before forgelog can use it. The proposal does not mention this ordering constraint explicitly.

### 5. Solution Creativity: 72/100

**Novelty over industry baseline (30/40)**: The zero-change console contract is the key innovation — treating console as a first-class backend that preserves original format byte-for-byte. This is not novel in the absolute sense (it is the null hypothesis), but recognizing that the simplest migration is one that changes nothing observable is a valuable insight. The per-invocation file with PID suffix is standard practice. Deduction: the proposal acknowledges "This is a straightforward adoption of a standard backend-pattern logging architecture" — honesty appreciated, but the novelty is modest.

**Cross-domain inspiration (22/35)**: The proposal references web server access log patterns (nginx, apache) for per-invocation files. The backend abstraction is inspired by slog's Handler interface. These are relevant cross-references. Deduction: no cross-domain inspiration beyond standard logging patterns. For example, database WAL (write-ahead logging) patterns, audit log patterns, or observability tracing (OpenTelemetry) concepts are not referenced. These are not necessary but would strengthen the creativity dimension.

**Simplicity of insight (20/25)**: The insight that "treating console as a backend that outputs the raw message guarantees behavioral equivalence by construction" is elegant. The Fprintln handling rule is simple and correct. The prefix-based migration categorization with longest-prefix-first matching is a clean algorithm. Deduction: the prefix parsing at migration time adds unnecessary complexity — the migration could simply assign levels explicitly at each call site. The proposal acknowledges this ("prefix parsing exists only in the migration layer") but the categorization table still adds conceptual weight.

### 6. Feasibility: 90/100

**Technical feasibility (37/40)**: The Backend interface is 2 methods. The implementation is straightforward Go using only stdlib. The migration is mechanical one-line changes. No showstoppers. Deduction: migrating ~102 call sites in a single PR is a large diff. While each change is mechanical, the aggregate diff size makes review harder. The proposal does not discuss whether this could be split.

**Resource & timeline feasibility (26/30)**: "Single PR, mechanically verifiable." Implementation order is specified: (1) forgelog package + tests, (2) config extension, (3) gitignore entry, (4) migrate all ~102 call sites. This is a clear plan. Deduction: no time estimate (hours? days?). The proposal says "mechanically verifiable" but does not estimate how long the mechanical work takes.

**Dependency readiness (27/30)**: No external dependencies. Config struct extension is additive (omitempty). Greenfield package. Deduction: the proposal does not address test infrastructure — how do existing tests interact with forgelog? Do tests need to initialize forgelog? Are there test-specific backends?

### 7. Scope Definition: 72/80

**In-scope items are concrete (27/30)**: Each in-scope item is a deliverable: pkg/forgelog package, config section, LogsConfig struct, validation, constants, gitignore entry, Init/Close calls, migration of ~102 sites, emergency disable, directory auto-creation, file permissions. The "Known gap" about pre-Init messages is honestly acknowledged. Deduction: the migration of ~102 call sites is stated as "single PR" but does not specify whether this includes forensic exclusion as a gate (verify forensic sites are NOT migrated).

**Out-of-scope explicitly listed (22/25)**: Seven items listed: JSON format, log rotation, remote shipping, CLI log viewer, test code migration, plugin changes, forensic command. These are specific and well-scoped. The forensic exclusion is well-justified ("stderr is its primary output channel"). Deduction: "Changes to the plugin (agents/commands/skills)" — does this mean the plugin's internal logging is out of scope, or that no new logging skill/command is in scope? The boundary could be clearer.

**Scope is bounded (23/25)**: CLI-only boundary is clear. "No plugin changes" is clear. Single PR scope is bounded. Deduction: the proposal does not specify when future code should use forgelog vs direct stderr. Is there a "forgelog for all diagnostics, fmt.Fprintf for user prompts" rule? This governance gap could lead to drift.

### 8. Risk Assessment: 78/90

**Risks identified (26/30)**: Seven risks identified, covering contention, disk accumulation, config parsing failure, sensitive info, logging regression, migration count inaccuracy, and Fprintln edge cases. This is a thorough set. The sensitive info risk is notable for its nuanced treatment ("diagnostic value > risk for local files") with CI/container recommendations. Deduction: no risk for "log file fills disk faster than cleanup" — the auto-cleanup only runs on startup, so a long-running loop producing large logs is not mitigated.

**Likelihood + impact rated (26/30)**: Ratings are provided and honest. "Sensitive info in log files" is rated Medium/Medium — appropriate. "Logging layer regression" is Low/High — appropriate (rare but severe). "Migration count inaccuracy" is Medium/Low — honest. Deduction: "Disk accumulation" is rated Low/Low with budget "7-70MB worst case" — this assumes ~10KB per invocation and ~1000 invocations in 7 days, but does not account for `run-tasks` loops that could produce much larger files.

**Mitigations are actionable (26/30)**: FORGE_NO_LOG=1 is actionable. Per-invocation files mitigate contention. Grep-based CI check for migration verification. File permissions for security. All mitigations are concrete. Deduction: the mitigation for sensitive info ("FORGE_NO_LOG=1 recommended in CI") shifts responsibility to the user without tooling support. An alternative mitigation (log redaction, env-var-based path filtering) is not discussed.

### 9. Success Criteria: 72/80

**Criteria are measurable and testable (27/30)**: SC-1 through SC-11 are all objectively verifiable. SC-2 ("diff of stderr shows zero changes") is particularly strong — it is a concrete, automatable test. SC-8 (grep-based migration completeness) is precise and CI-friendly. SC-10 is measurable (set env var, verify no directory created). Deduction: SC-1 ("writes AUTO-RESTORE diagnostic to .forge/logs/... with structured format") does not specify what "structured format" means — is it the timestamp+level+message format from section 2? If so, SC-1 should reference it explicitly.

**Coverage is complete (22/25)**: Persistence (SC-1), dual output (SC-2), filtering (SC-3), cleanup (SC-4), gitignore (SC-5), auto-creation (SC-6), defaults (SC-7), migration completeness (SC-8), concurrency (SC-9), emergency disable (SC-10), permissions (SC-11). This covers all in-scope items. Deduction: no SC verifies the config validation behavior (unrecognized level -> default info, retentionDays < 1 -> default 7) mentioned in the in-scope list.

**SC internal consistency (23/25)**: All SCs are mutually satisfiable. SC-10 (no .forge/logs/ when disabled) is compatible with SC-6 (auto-created when enabled) because they operate under different conditions. SC-2 (console unchanged) is compatible with SC-8 (all sites migrated) because ConsoleBackend outputs raw messages. Deduction: SC-1 references "structured format (timestamp [LEVEL] message)" but this format is defined in the Format Layer section, not in the SC itself. If the format specification changes, SC-1's verification criteria become ambiguous.

### 10. Logical Consistency: 82/90

**Solution addresses the stated problem (33/35)**: The file-based logging layer directly solves the "ephemeral diagnostics" problem. The dual-backend design ensures no behavioral regression. The emergency disable provides a safety net. The auto-cleanup prevents indefinite accumulation. Direct alignment between problem and solution.

**Scope <-> Solution <-> SC aligned (26/30)**: The in-scope items map to solution components. The SCs verify in-scope deliverables. The alignment is strong. Deduction: the in-scope list mentions "Config validation in forgelog.Init(): unrecognized level -> default info; retentionDays < 1 -> default 7" but no SC verifies this validation behavior. There is a gap between scope and SC for config validation.

**Requirements <-> Solution coherent (23/25)**: The requirements (persistence, dual output, level filtering, auto-cleanup, config integration) map cleanly to solution components. No orphan requirements or solution features without requirements. Deduction: the "prefixless sites require individual review" caveat in the Call-Site Categorization Table introduces ambiguity — the solution claims mechanical migration but acknowledges non-mechanical cases. This does not undermine coherence but reduces precision.

---

## Phase 3: Blindspot Hunt

[blindspot-1] **Duplicate YAML code fence.** The config section at line 170 has ```` ```yaml ```` followed immediately by another ```` ```yaml ```` at line 172. This renders as a broken code block in most Markdown processors. The first fence is empty, and the second contains the actual config. This is a formatting artifact but signals insufficient proofreading of a document that has been through pre-revision.

[blindspot-2] **No governance rule for future code.** The proposal migrates ~102 existing call sites but does not specify a rule for new code. Should new diagnostic messages use `forgelog.Info()` or `fmt.Fprintf(os.Stderr, ...)`? Without an explicit rule, the migration can be gradually undone by new code that bypasses forgelog. The proposal mentions "new code calls forgelog.Warn() directly and needs no prefix convention" but this is a migration note, not a governance rule. A one-line rule in the Constraints section would suffice.

[blindspot-3] **ConsoleBackend has no level filtering, but the API accepts levels.** The `ConsoleBackend.Write(level, timestamp, msg)` receives the level parameter but ignores it ("No level filtering — always outputs all messages"). This means `forgelog.Debug("verbose trace\n")` will always appear on console even when the user would prefer less verbosity. The proposal explicitly preserves this behavior ("preserves current behavior exactly") but does not discuss whether this is desirable. Currently, debug messages (prefixed `[debug]`) are gated behind some condition in the existing code — the migration would make them unconditional on console. This could change observed behavior, contradicting SC-2.

[blindspot-4] **The `Close()` function and `defer` placement.** The proposal says "Call via defer in each command's runE." But Go's `defer` does not execute on `os.Exit()`. The proposal acknowledges this in the data safety section ("os.Exit/log.Fatal bypass defer but cannot lose data since all writes were already issued to the OS"). However, if Close() is "not strictly needed for data safety" and os.Exit bypasses it, then the file handle is leaked on os.Exit paths. On long-running systems this is irrelevant, but it means the `defer forgelog.Close()` pattern provides less cleanup than claimed.

[blindspot-5] **The migration changes error handling semantics.** Currently, `fmt.Fprintf(os.Stderr, ...)` returns `(int, error)` which is universally ignored. After migration, `forgelog.Warn(...)` has no return value. This is fine — but the proposal does not address call sites that currently DO check the error return of fmt.Fprintf (if any exist). If any call site handles the write error (retries, fallback), migration to forgelog silently drops that error handling.

[blindspot-6] **retentionDays cleanup runs on startup, but what about commands that never complete?** A `run-tasks` loop that runs for days continuously writes to the same log file. The cleanup only runs on command startup, so old files from completed invocations are cleaned up — but the active file grows unbounded. For typical volumes this is fine, but the proposal does not state this assumption.

---

## Bias Detection Report

Annotated regions (marked with `<!-- pre-revised: ... -->`):
- 10 pre-revised markers covering ~14 paragraphs
- Attack points in annotated regions: 3
  - [Solution Clarity] Duplicate YAML fence (blindspot-1, trivial)
  - [Logical Consistency] Close() vs defer vs os.Exit (blindspot-4, pre-revised section)
  - [Solution Clarity] ConsoleBackend blocking concern (acknowledged in pre-revised text, no deduction)
- Density: 3/14 = 0.21

Unannotated regions:
- ~25 paragraphs without pre-revised markers
- Attack points in unannotated regions: 12
- Density: 12/25 = 0.48

Ratio (annotated/unannotated): 0.44

The lower attack density in annotated regions indicates the pre-revision addressed several structural gaps effectively. The annotated regions received focused attention during revision and consequently have fewer issues. The unannotated regions, particularly the Call-Site Categorization Table and the config section, contain the remaining weaknesses.

---

## Summary Table

| Dimension | Score | Max |
|-----------|-------|-----|
| Problem Definition | 98 | 110 |
| Solution Clarity | 108 | 120 |
| Industry Benchmarking | 85 | 120 |
| Requirements Completeness | 95 | 110 |
| Solution Creativity | 72 | 100 |
| Feasibility | 90 | 100 |
| Scope Definition | 72 | 80 |
| Risk Assessment | 78 | 90 |
| Success Criteria | 72 | 80 |
| Logical Consistency | 82 | 90 |
| **Total** | **852** | **1000** |

## ATTACK_POINTS

1. [Industry Benchmarking]: No reference to comparable CLI tools' logging practices — the comparison table includes slog and generic patterns but not cargo, kubectl, docker, or similar CLI tools that face the same diagnostic persistence problem. Quote: "Most CLI tools use one of: (1) stderr-only output (current forge approach), (2) log/slog with structured handlers, (3) per-invocation log files" — this categorizes the industry in three buckets without citing specific tools for bucket (3). Must improve: reference at least 2 specific CLI tools with per-invocation logging and how their approach compares.

2. [Industry Benchmarking]: slog rejection rationale is qualitative, not quantitative — quote: "the overhead of implementing two custom handlers (one for raw console, one for structured file) exceeds the value for a ~150-line logging package." The word "overhead" is vague. If the custom handlers would be 30 lines each (60 total), the difference is 60 vs 150 lines — forgelog is larger. Must improve: quantify the slog implementation cost (estimated lines) to make the comparison honest.

3. [Solution Clarity]: Duplicate YAML code fence in config section — line 170 has ```` ```yaml ```` followed immediately by ```` ```yaml ```` at line 172. The first fence is empty, creating a rendering artifact. Must improve: remove the duplicate fence.

4. [Solution Clarity]: Init function's logsDir parameter origin is unspecified — quote: "Init(config *forgeconfig.LogsConfig, logsDir string) error". Where does logsDir come from? Is it always `.forge/logs/` relative to CWD? Relative to project root? Must improve: specify the source of logsDir.

5. [Requirements Completeness]: [blindspot-3] ConsoleBackend always outputs all messages including DEBUG — quote: "No level filtering — always outputs all messages." Currently, `[debug]`-prefixed messages appear to be gated behind some condition in the existing code (output.go). After migration, `forgelog.Debug()` would unconditionally write to console, potentially changing observed behavior and contradicting SC-2 ("Console output is byte-identical before and after migration"). Must improve: verify that existing debug messages are unconditional, or specify that ConsoleBackend must replicate the existing gating logic.

6. [Solution Creativity]: The proposal self-describes as "a straightforward adoption of a standard backend-pattern logging architecture" — the creativity is in the zero-change console contract, which is essentially a well-executed null hypothesis. The per-invocation file pattern is standard web server practice. The Backend interface is a simplified slog.Handler. Must improve: while honesty is appreciated, the novelty score reflects that this is competent engineering rather than creative problem-solving.

7. [Scope Definition]: No governance rule for future diagnostic code — the proposal migrates existing call sites but does not state "all new diagnostic output must use forgelog." Quote from Constraints: "No changes to plugin layer — this is CLI-only" specifies what is excluded but not what is required going forward. Must improve: add a one-line rule in Constraints or Scope stating that new CLI diagnostic code must use forgelog.

8. [Success Criteria]: SC-1 references "structured format (timestamp [LEVEL] message)" but does not define it — the format is specified in the Format Layer section, not in the SC itself. If the format changes during implementation, SC-1's verification criteria become ambiguous. Must improve: either define the format inline in SC-1 or add a dedicated SC for log line format compliance.

9. [Success Criteria]: No SC for config validation behavior — the in-scope list includes "Config validation in forgelog.Init(): unrecognized level -> default info; retentionDays < 1 -> default 7" but no SC verifies this. Must improve: add SC-12 verifying that invalid config values fall back to defaults.

10. [Risk Assessment]: Disk accumulation risk budget assumes ~10KB/invocation but does not account for run-tasks loops — quote: "~5-10KB/invocation, 7-day default ~ 7-70MB worst case." A single run-tasks loop with hundreds of subagent invocations could produce 500KB+ per invocation. The "worst case" is not worst case. Must improve: calculate a genuine worst case (e.g., 1000 invocations at 50KB each over 7 days) and confirm it is acceptable.

11. [Logical Consistency]: conflict-with-pre-revision — The pre-revised Close() documentation (lines 96-101) states "os.Exit/log.Fatal bypass defer but cannot lose data since all writes were already issued to the OS." But this means the file handle is leaked on os.Exit paths, which contradicts the in-scope item "defer forgelog.Close()" being a sufficient cleanup mechanism. Must improve: either acknowledge the handle leak as acceptable (OS reclaims on process exit) or specify an alternative cleanup path.

12. [Solution Creativity]: Prefix-based migration categorization adds complexity that the migration itself makes unnecessary — quote: "prefix parsing exists only in the migration layer to classify existing call sites." If prefix parsing is only for migration, the categorization table's 16-row complexity is a documentation burden that provides no ongoing value. Must improve: note that the categorization table is a one-time migration guide, not a runtime feature.
