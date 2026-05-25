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
- Preconditions: "prompt.go has been rewritten, resolveScope() deleted and replaced by direct SurfaceKey reading"
- Input: "user inspects prompt.go resolveScope() function (now removed)"
- Output: "resolveScope() no longer exists. renderTemplate() reads SurfaceKey directly from task struct, uses forge surfaces CLI output to determine surface-key dynamically instead of hardcoded projectType matching"
- State: "prompt.go uses SurfaceKey for template rendering, no scope-related code paths remain"
- Side-effect: "none"

## Outcome "cli-execution-failure"
- Preconditions: "forge surfaces command is not available (CLI not installed or version too old)"
- Input: "any component that calls forge surfaces attempts to invoke the CLI"
- Output: "component outputs error to stderr with the CLI output and recovery hint (check forge CLI is installed and version >= required version), exits with exit code 1 (retryable)"
- State: "no surface-key determined, component falls back or aborts"
- Side-effect: "none"

## Outcome "render-template-variable-replaced"
<!-- source: inferred -->
<!-- reasoning: Tech design specifies 16 prompt templates must have {{SCOPE}} replaced with {{SURFACE_KEY}}. This is a direct consequence of resolveScope() deletion. Verifying the template rendering still works with new variable names. -->
- Preconditions: "prompt.go renderTemplate() has been updated with new variable names"
- Input: "any prompt-based operation that uses templates with the new SURFACE_KEY variable"
- Output: "renderTemplate() correctly substitutes {{SURFACE_KEY}} with task.SurfaceKey value and {{TEST_TYPE_ARG}} with the surface type derived argument"
- State: "all 16 prompt templates use SURFACE_KEY variable, no hardcoded frontend/backend references remain in template logic"
- Side-effect: "none"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
