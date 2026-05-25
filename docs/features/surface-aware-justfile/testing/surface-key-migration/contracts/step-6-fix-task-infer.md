---
journey: "surface-key-migration"
step: 6
step-action: "verify quality-gate fix-task infers surface from file path"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-key-migration/journey.md
---

# Contract: surface-key-migration / Step 6: Verify quality-gate fix-task infers surface from file path

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "quality-gate creates a fix-task from a failing test file path, surface detection is available"
- Input: "quality-gate fix-task is created from a failing test file path (e.g., a test file within a configured surface directory)"
- Output: "fix-task surface-key and surface-type are inferred from the failing file path using surface detection longest-prefix-match"
- State: "fix-task created with correct surface-key and surface-type matching the failing file's surface"
- Side-effect: "surface detection CLI invoked with the failing file path"

## Outcome "inference-failure"
<!-- source: inferred -->
<!-- reasoning: When surface detection fails (not installed, config error) or returns no match for the failing file path, the fix-task must handle both cases gracefully. These are semantically similar: both result in missing surface info on the fix-task. -->
- Preconditions: "failing test file path does not match any configured surface entry, or surface detection fails or is unavailable"
- Input: "quality-gate creates fix-task but surface detection cannot resolve the file path or fails during invocation"
- Output: "fix-task created with empty or default surface-key and surface-type, error logged to stderr with recovery hint when detection fails"
- State: "fix-task exists but lacks surface-specific routing information, downstream components handle missing surface info"
- Side-effect: "none"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
