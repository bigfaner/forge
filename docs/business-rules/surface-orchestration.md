---
title: "Surface Orchestration Rules"
domains: [surface, orchestration, probe, teardown, recipe, surface-key, surface-type]
---

# Surface Orchestration Rules

_Source: feature/surface-aware-justfile_

## Surface Types

### BIZ-surface-orchestration-001: Surface Type Fixed Enumeration

**Rule**: Forge recognizes exactly 5 surface types: `web`, `api`, `cli`, `tui`, `mobile`. Surface-type determines the orchestration strategy (e.g., web/api/mobile require dev->probe->[per-journey test]->teardown; cli/tui use [per-journey test]->teardown). Surface-key names are user-defined in `.forge/config.yaml` surfaces field; surface-type values are constrained to this fixed set.
**Context**: Establishes the type-system that maps user-defined surface keys to orchestration strategies. Without fixed types, skills cannot determine the correct execution sequence.
**Source**: feature/surface-aware-justfile BIZ-001

## Orchestration Sequences

### BIZ-surface-orchestration-002: Surface Orchestration Sequence Table

**Rule**: Each surface type has a fixed orchestration sequence with per-journey test loops:

| Surface | Sequence | Key Differences |
|---------|----------|-----------------|
| web | dev -> probe -> [per-journey test] -> teardown | probe checks page root path; dev/probe once, test loops per journey |
| api | dev -> probe -> [per-journey test] -> teardown | probe checks /healthz; dev/probe once, test loops per journey |
| cli | [per-journey test] -> teardown | no dev, no probe, no build; test loops per journey |
| tui | [per-journey test] -> teardown | no dev, no probe, no build; test loops per journey |
| mobile | dev -> probe -> [per-journey test] -> teardown | dev/probe once, test loops per journey |

Dev server failures MUST NOT proceed to subsequent steps -- immediately teardown and exit. Test recipe format is `just <surface>-test <journey>` where `<journey>` is a directory name from `docs/features/<slug>/testing/`.
**Context**: Defines the complete execution sequence per surface type, consumed by run-tests skill via surface rule files. Journey isolation means dev/probe execute once, then test runs per-journey sequentially, then teardown once.
**Source**: feature/surface-aware-justfile BIZ-002

## Probe Behavior

### BIZ-surface-orchestration-003: Probe Retry Parameters

**Rule**: Probe checks use unified retry parameters: max 3 retries, 5-second interval between retries. Failure behavior: teardown + abort. All 3 retries failing is treated as retryable failure (exit code 1).
**Context**: Ensures consistent probe behavior across all surface types that require health checks (web, api, mobile).
**Source**: feature/surface-aware-justfile BIZ-003

### BIZ-surface-orchestration-004: Probe Failure HARD-GATE

**Rule**: After probe failure, the system MUST NOT retry the probe or restart dev within the same orchestration cycle. This is a non-violable constraint -- if probe has judged failure, the service has a fundamental problem, and retrying only masks the issue. The upper scheduler (e.g., CI) can distinguish retryable (exit 1) vs blocking (exit 2) to decide whether to retry in a new orchestration cycle.
**Context**: Prevents retry loops that hide real service issues. Applies to all surface types that use probe (web, api).
**Source**: feature/surface-aware-justfile BIZ-004

## Teardown

### BIZ-surface-orchestration-005: Teardown Idempotency

**Rule**: Teardown MUST be idempotent: if the PID does not exist, skip silently. If kill fails, retry once; if still failing, log process info and continue subsequent cleanup steps. Final guarantee: `.forge/test-state.json` state file is cleaned up (deleted or marked completed) regardless of teardown success.
**Context**: Ensures no orphan processes remain after test orchestration, even in failure scenarios. Dev server crashes trigger probe timeout followed by teardown.
**Source**: feature/surface-aware-justfile BIZ-005

## Naming

### BIZ-surface-orchestration-006: Surface-Key Naming Constraint

**Rule**: Surface-key values MUST match `[a-zA-Z0-9_-]` only -- no `/` or `+` characters, ensuring compatibility with just recipe name syntax. Enforcement occurs in init-justfile recipe generation.
**Context**: Prevents just recipe name injection and ensures generated recipe names are syntactically valid.
**Source**: feature/surface-aware-justfile BIZ-006
