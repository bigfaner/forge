---
feature: "surface-aware-justfile"
generated: "2026-05-26"
status: draft
---

# Business Rules: Surface-Aware Justfile

## Surface Orchestration

### BIZ-001: Surface Type Fixed Enumeration

**Rule**: Forge recognizes exactly 5 surface types: `web`, `api`, `cli`, `tui`, `mobile`. Surface-type determines the orchestration strategy (e.g., web/api require dev->probe->test->teardown; cli/tui use build->dev->test). Surface-key names are user-defined in `.forge/config.yaml` surfaces field; surface-type values are constrained to this fixed set.
**Context**: Establishes the type-system that maps user-defined surface keys to orchestration strategies. Without fixed types, skills cannot determine the correct execution sequence.
**Scope**: [CROSS]
**Source**: prd-spec.md "Surface 感知配方生成" + "Surface-key 值域统一"

### BIZ-002: Surface Orchestration Sequence Table

**Rule**: Each surface type has a fixed orchestration sequence:

| Surface | Sequence | Key Differences |
|---------|----------|-----------------|
| web | dev(background) -> probe -> test -> teardown | probe checks page root path |
| api | dev(background) -> probe -> test -> teardown | probe checks /healthz |
| cli | build -> dev -> test | no service start, no probe |
| tui | build -> dev -> test | no service start, no probe |
| mobile | test-setup -> dev -> test -> teardown | test-setup prepares emulator |

Dev server failures MUST NOT proceed to subsequent steps -- immediately teardown and exit.
**Context**: Defines the complete execution sequence per surface type, consumed by run-tests skill via surface rule files.
**Scope**: [CROSS]
**Source**: prd-spec.md "Surface 编排模式"

### BIZ-003: Probe Retry Parameters

**Rule**: Probe checks use unified retry parameters: max 3 retries, 30-second interval between retries, 90-second total timeout (max-retries x interval). Failure behavior: teardown + abort. Log format: `[probe] [retry <current>/<max>] <url> -- <reason>`.
**Context**: Ensures consistent probe behavior across all surface types that require health checks (web, api, mobile).
**Scope**: [CROSS]
**Source**: prd-spec.md "Probe 重试规格"

### BIZ-004: Probe Failure HARD-GATE

**Rule**: After probe failure, the system MUST NOT retry the probe or restart dev within the same orchestration cycle. This is a non-violable constraint -- if probe has judged failure, the service has a fundamental problem, and retrying only masks the issue. The upper scheduler (e.g., CI) can distinguish retryable (exit 1) vs blocking (exit 2) to decide whether to retry in a new orchestration cycle.
**Context**: Prevents retry loops that hide real service issues. Applies to all surface types that use probe (web, api).
**Scope**: [CROSS]
**Source**: prd-spec.md "HARD-GATE 定义"

### BIZ-005: Teardown Idempotency

**Rule**: Teardown MUST be idempotent: if the PID does not exist, skip silently. If kill fails, retry once; if still failing, log process info and continue subsequent cleanup steps. Final guarantee: `.forge/test-state.json` state file is cleaned up (deleted or marked completed) regardless of teardown success.
**Context**: Ensures no orphan processes remain after test orchestration, even in failure scenarios. Dev server crashes trigger probe timeout followed by teardown.
**Scope**: [CROSS]
**Source**: prd-spec.md "Reliability"

### BIZ-006: Surface-Key Naming Constraint

**Rule**: Surface-key values MUST match `[a-zA-Z0-9_-]` only -- no `/` or `+` characters, ensuring compatibility with just recipe name syntax. Enforcement occurs in init-justfile recipe generation.
**Context**: Prevents just recipe name injection and ensures generated recipe names are syntactically valid.
**Scope**: [CROSS]
**Source**: prd-spec.md "混合项目生成与编排流程"

## Error Handling

### BIZ-007: Surface Information Unavailable Error

**Rule**: When surface information is unavailable from both sources (task frontmatter and `forge surfaces` CLI), the system outputs an error message to stderr containing: (1) the failure reason, and (2) a recovery hint: "Please configure surfaces in `.forge/config.yaml` surfaces field, or specify surface-type in task frontmatter". Exits with code 2 (blocking).
**Context**: Dual-source fallback ensures graceful degradation. Exit code 2 (blocking) because missing surface configuration requires manual intervention.
**Scope**: [LOCAL]
**Source**: prd-spec.md "Error Handling Paths"

### BIZ-008: Unknown Surface Type Error

**Rule**: When an unknown surface type is encountered (e.g., loading a rule file for a non-existent type), output error to stderr: "Execution strategy rule file for surface type '<type>' does not exist. Supported types: web/api/cli/tui/mobile". Exit with code 2 (blocking).
**Context**: Prevents silent failures when an invalid surface type propagates through the system.
**Scope**: [LOCAL]
**Source**: prd-spec.md "Error Handling Paths"

## Compatibility

### BIZ-009: Zero Regression Guarantee

**Rule**: Projects without surface configuration MUST produce output identical to the current (pre-feature) behavior. Verified by diff output comparison.
**Context**: Ensures backward compatibility for existing Forge users who have not adopted the surface configuration model.
**Scope**: [LOCAL]
**Source**: prd-spec.md "Goals" + "Compatibility"
