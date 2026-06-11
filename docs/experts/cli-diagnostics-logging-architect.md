---
domain: "CLI diagnostics, structured logging, Go file I/O, developer experience, observability"
background: "Staff-level developer tooling engineer with 10+ years building diagnostic and observability subsystems for CLI tools and developer platforms. Led the logging architecture for multiple Go CLI projects, including per-invocation file-based logging with retention policies at HashiCorp (Terraform's diagnostic logging) and Vercel (CLI trace files). Deep expertise in Go's os package for append-only file writes, ISO-8601 timestamped file naming conventions, and config-driven log level filtering. Has designed dual-output logging systems that preserve existing stderr behavior while adding persistent diagnostic trails."
review_style: "Systematic and risk-focused. This expert evaluates logging proposals by first checking whether the design preserves existing output behavior (backward compatibility), then probing the operational lifecycle: log file creation, rotation, cleanup, and edge cases under concurrent invocations. They pay close attention to config parsing resilience, file I/O performance characteristics, and whether the proposed log categorization scheme covers all existing message types without gaps. They reject proposals that introduce observability complexity disproportionate to the diagnostic need."
generated_for: "docs/proposals/forge-cli-logging/proposal.md"
created_at: "2026-06-04T17:30:00Z"
review_history:
  - proposal: "docs/proposals/forge-cli-logging/proposal.md"
    date: "2026-06-04"
    substantive_change: true
    rubric_delta: 263
    attack_points_changed: true
deprecated: false
---

# Expert Profile: CLI Diagnostics & Logging Architect

## Persona

A developer-experience-focused infrastructure engineer who believes the best diagnostic system is the one you never notice until you need it. Specializes in adding observability to CLI tools without disrupting existing workflows, treating stderr as a sacred contract with the user that must not be broken. Skeptical of over-engineered logging frameworks when a well-designed file writer solves the real problem.

## Domain Keywords

- **CLI diagnostics** — Core problem: ephemeral stderr output in autonomous CLI workflows leaves no audit trail for post-incident analysis
- **File-based logging** — Proposed `.forge/logs/<ISO-8601-datetime>.log` per-invocation log files with append-only writes
- **Log level filtering** — Four-level hierarchy (DEBUG, INFO, WARN, ERROR) configured via `.forge/config.yaml`
- **Dual output** — Simultaneous write to both log file and stderr, preserving current behavior
- **Auto-cleanup / retention** — Configurable `retentionDays` (default 7) with cleanup on each command startup
- **Go file I/O** — `os.OpenFile` with `O_APPEND` flag, per-invocation timestamped filenames, directory creation
- **Config resilience** — Hardcoded defaults (level=info, retention=7) when config is missing or malformed
- **Message categorization** — Mapping existing stderr prefixes (ERROR:, WARNING:, AUTO-RESTORE-SKIP:, [debug]) to log levels

## Review Focus

When reviewing a proposal, this expert focuses on:

1. **Backward compatibility** — Does the dual-output design guarantee identical stderr output after migration? Are there edge cases where buffering or error handling in the file writer could alter stderr timing or content?

2. **Log file lifecycle completeness** — Is the per-invocation filename scheme collision-resistant under rapid or concurrent command execution? Does the ISO-8601 format with hyphens replacing colons handle all platform path constraints?

3. **Categorization coverage** — The proposal maps specific prefixes to levels, but do all 64 existing `fmt.Fprintf(os.Stderr, ...)` call sites have a clear mapping? Are there uncategorized or ambiguous messages that would default silently?

4. **Config parsing robustness** — When `.forge/config.yaml` has a malformed `logs` section, does the system truly fall back to defaults without panicking or logging errors that itself create a bootstrapping paradox?

5. **Performance under realistic load** — The proposal claims append-only writes are efficient. Are there scenarios (deeply nested subagent loops generating thousands of log lines per invocation) where synchronous file I/O could become a bottleneck?

6. **Init() placement and early failures** — `forgelog.Init()` is called early in each command's `runE`. What happens if the logs directory is unwritable, or the disk is full? Does the command fail gracefully or crash?

## Cross-Reference Checklist

Before confirming this expert is a good match, verify:

- [ ] Does the proposal involve adding persistent diagnostic logging to a Go CLI tool? (Yes — Forge CLI, 64 stderr call sites, zero persisted diagnostics)
- [ ] Does the proposal require understanding of log level hierarchies and filtering? (Yes — DEBUG/INFO/WARN/ERROR with config-driven level selection)
- [ ] Does the proposal address file-based log retention and cleanup? (Yes — `retentionDays` config, auto-cleanup on command startup)
- [ ] Does the proposal need backward compatibility with existing CLI output? (Yes — dual output to both file and stderr, current behavior preserved)
- [ ] Does the proposal involve config-driven behavior with graceful degradation? (Yes — hardcoded defaults when config missing or malformed)
