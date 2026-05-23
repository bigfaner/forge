---
title: "Error Reporting Rules"
domains: [exit-code, stderr, actionable, recovery, hint, AIError, blocking]
---

# Error Reporting Rules

_Source: feature/forge-cli-v3_

## Exit Codes

### BIZ-error-reporting-001: Consistent Exit Code Semantics

**Rule**: Exit code 0 = success (or intentional no-op, e.g., "no tasks to clean up"). Exit code 1 = retryable failure (not found, validation error, conflict). Exit code 2 = blocking errors (invalid transition, invalid path, contract unverifiable). AIError's `ExitCode()` method differentiates: `INVALID_TRANSITION`, `INVALID_PATH`, `CONTRACT_UNVERIFIABLE` return 2; all other error codes return 1. Cobra usage errors also produce exit code 2 by default.
**Context**: AI agents rely on exit codes to determine next action; consistent semantics prevent misinterpretation.
**Source**: feature/forge-cli-v3 BIZ-008

## Error Messages

### BIZ-error-reporting-002: Actionable Error Messages

**Rule**: Every error message on stderr MUST contain: (1) the specific failure reason, and (2) a hint for recovery when applicable. Example: "unknown profile: <value>" MUST be followed by listing all supported profiles.
**Context**: AI agents need self-correcting feedback loops without human intervention.
**Source**: feature/forge-cli-v3 BIZ-009
