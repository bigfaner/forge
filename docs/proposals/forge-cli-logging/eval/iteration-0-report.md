# Iteration 0 Report: Pre-Revision (Freeform Findings)

**Iteration**: 0
**Title**: Pre-Revision (Freeform Findings)

## ATTACK_POINTS

1. **[high]** Data-safety claim is factually inaccurate — `O_APPEND` writes to kernel buffer cache, not stable storage. Claim "signal kills cannot lose data" is false for SIGKILL. | quote: "No buffering — each write is persisted before function returns. os.Exit, panic, and signal kills cannot lose data." | improvement: Soften claim to "each write is issued to the OS before the function returns; persistence to stable storage depends on OS behavior." Remove absolute signal-kill guarantee.

2. **[high]** Fprintln newline handling is ambiguous — API contract does not specify whether forgelog functions append \n internally. This directly threatens SC-2 (byte-identical console). | quote: "forgelog.Warn(msg) handles the trailing newline" — "handles" is ambiguous | improvement: State explicitly: "forgelog functions do NOT append \n. The formatted message is output exactly as-is." Add migration rules for Fprintln callers.

3. **[high]** Synchronous sequential backend dispatch means a stalled file write blocks console output — the exact opposite of the design intent. | quote: "ConsoleBackend first, FileBackend second. If file write fails, console has already written." | improvement: Note this risk explicitly. Consider acknowledging that ConsoleBackend should complete before FileBackend starts, or document that file write errors are silently ignored (not blocked on).

4. **[high]** os.Exit/log.Fatal bypasses defer, so Close() never runs on fatal exits. | quote: "Close closes all backends. Call via defer in each command's runE." | improvement: Document that with O_APPEND per-write, Close is not strictly needed for data safety. Add note that log.Fatal paths should call forgelog.Flush() if added, or migrate log.Fatal usage to forgelog + manual exit.

5. **[medium]** HINT: prefix mapped to ERROR is semantically wrong — hints are remediation suggestions, not errors. | quote: "HINT: ... | ERROR | errors.go | Part of structured error/warning blocks" | improvement: Reclassify HINT: from ERROR to INFO.

6. **[medium]** SC-8 verification grep only covers Fprintf/Fprintln patterns, missing slog.Warn and log.Printf. | quote: "grep -r 'fmt.Fprintf(os.Stderr\|fmt.Fprintln(os.Stderr' ... returns 0 results" | improvement: Add slog and log patterns to SC-8 grep, or add a separate SC.

7. **[medium]** Call-site counts are approximate and likely inaccurate. Baseline scorer found discrepancies (proposal says ~102, actual may be ~107). | quote: "Total ~87 ~15 10 ~102" | improvement: Add note that exact counts will be verified at implementation time, or replace with "approximately 100+" to avoid precision claims.

8. **[medium]** No rollback plan beyond FORGE_NO_LOG=1 env var — no config-based disable, no versioned rollout. | quote: "Emergency disable: FORGE_NO_LOG=1 env var skips FileBackend" | improvement: Add `logs.enabled: false` config option as a persistent alternative to the env var, or add explicit rollback documentation.

9. **[medium]** Prefixless messages defaulting to INFO may misclassify semantically important messages. 30+ prefixless sites require per-site judgment, contradicting "one-line mechanical" migration claim. | quote: "Fallback: Any message without a matching prefix defaults to INFO level." | improvement: Add explicit note that prefixless call sites require individual review and the "mechanical" claim applies only to prefixed sites. List known prefixless sites with their correct levels.

10. **[medium]** slog dismissal is factually incorrect — slog supports multiple handlers. | quote: "slog's single-format Handler" | improvement: State actual reason for rejection honestly: "slog's Handler interface could support dual-output, but the overhead of custom handler implementations exceeds the value for a ~150-line human-readable logging package."

11. **[low]** Categorization table lists ERROR: and error: as separate rows despite case-insensitive matching — confusing for implementers. | quote: "ERROR: ... (no indent) | ERROR" and "error: ... (lowercase) | ERROR" as separate rows | improvement: Merge into one row with footnote about case-insensitive matching.

12. **[low]** Windows compatibility not addressed — Unix file modes (0600/0700), O_APPEND atomicity, PID naming are all Unix-centric. | quote: "Log files created with mode 0600 (owner-only)" | improvement: Add note in Constraints about Windows behavior, or state Windows support is out of scope.

13. **[low]** Sensitive data risk assessment only covers local single-user scenarios, not CI/container. | quote: "No redaction at log time — diagnostic value > risk for local files" | improvement: Expand assessment to note CI/container risks and recommend FORGE_NO_LOG=1 in those environments.

14. **[low]** Pre-Init message capture gap — messages emitted before forgelog.Init() are lost from log file. | quote: "forgelog.Init() called early in each command's runE function" | improvement: Acknowledge this gap explicitly in the design.

15. **[low]** No success criterion verifies file permissions (0600). | quote: "Log files created with mode 0600 (owner-only)" in NFR | improvement: Add SC for permission verification.

## BORDERLINE_FINDINGS

None.

## SKIPPED_FINDINGS

None — all findings classified as factual corrections or structural suggestions.

## Triage Summary

| Triage Layer | Count | Action |
|-------------|-------|--------|
| Factual correction | 7 | Direct edit |
| Structural suggestion | 8 | Edit where verifiable inconsistency exists |
| Subjective preference | 0 | None skipped |

**Hit Rate**: 18/16 = 1.125 (over-counted due to multiple findings per keyword paragraph)

## Rubric

All dimensions: N/A (freeform-driven pre-revision)
