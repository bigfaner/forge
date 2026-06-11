---
journey: "surface-key-migration"
step: 3
step-action: "verify resolveScope rewrite"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-key-migration/journey.md
---

# Contract: surface-key-migration / Step 3: Verify resolveScope() rewrite

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "prompt template variable system has been updated to use surface-aware variable names"
- Input: "user invokes any prompt-based operation that uses templates with surface variables"
- Output: "template rendering substitutes surface-key variables with the task's configured surface key value and surface-type-derived arguments, no hardcoded project type references remain"
- State: "all prompt templates use surface-aware variable names, the old scope-based template resolution path no longer exists"
- Side-effect: "none"

## Outcome "cli-execution-failure"
- Preconditions: "the surface detection CLI command is not available or returns an error"
- Input: "any component that needs surface information attempts to invoke the CLI"
- Output: "component outputs error to stderr with the CLI error output and recovery hint (verify CLI is installed and at the required version), exits with exit code 1 (retryable)"
- State: "no surface-key determined, component falls back or aborts"
- Side-effect: "none"

## Outcome "render-template-variable-replaced"
<!-- source: inferred -->
<!-- reasoning: Tech design specifies all prompt templates must have old scope variables replaced with surface-aware variables. This is a direct consequence of the scope resolution rewrite. Verifying template rendering still works with new variable names. -->
- Preconditions: "template rendering system has been updated with new surface-aware variable names"
- Input: "any prompt-based operation that uses templates with the new surface key variable"
- Output: "template rendering correctly substitutes the surface key and surface type derived arguments into all prompt templates"
- State: "all prompt templates use surface-aware variable names, no hardcoded project type references remain in template logic"
- Side-effect: "none"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
