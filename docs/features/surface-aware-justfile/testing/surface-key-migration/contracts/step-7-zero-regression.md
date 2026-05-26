---
journey: "surface-key-migration"
step: 7
step-action: "verify zero-regression for projects without surfaces"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-key-migration/journey.md
---

# Contract: surface-key-migration / Step 7: Verify zero-regression for projects without surfaces

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "project has no surfaces configuration in config.yaml, all migration phases completed"
- Input: "user runs the full workflow on a project with no surfaces configuration"
- Output: "all behavior is identical to the pre-feature baseline. No surface-key or surface-type fields appear in generated tasks. resolveScope() replacement falls back gracefully"
- State: "output matches pre-feature baseline exactly (diff shows no changes)"
- Side-effect: "none"

## Outcome "dead-code-and-validation"
- Preconditions: "after migration, obsolete code from the old scope system remains, or user-defined surface-keys contain invalid characters for recipe names"
- Input: "user compiles the project or runs init-justfile with surface-keys that have invalid characters"
- Output: "project compiles cleanly with all obsolete scope-related code removed, all callers updated to use surface-based APIs. For invalid surface-keys, init-justfile validates and outputs descriptive error with recovery hint"
- State: "project clean of dead code, surface-key validation enforced at recipe generation boundary"
- Side-effect: "none"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
