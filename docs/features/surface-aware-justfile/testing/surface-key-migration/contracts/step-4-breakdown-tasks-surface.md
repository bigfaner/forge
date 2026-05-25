---
journey: "surface-key-migration"
step: 4
step-action: "verify breakdown-tasks generates tasks with surface fields"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-key-migration/journey.md
---

# Contract: surface-key-migration / Step 4: Verify breakdown-tasks generates tasks with surface fields

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "project has configured surfaces in .forge/config.yaml, breakdown-tasks SKILL.md is updated to call forge surfaces --json"
- Input: "user runs breakdown-tasks for a project with configured surfaces"
- Output: "generated task files have surface-key and surface-type in their frontmatter, values match the configured surfaces (e.g., surface-key: admin-panel, surface-type: web)"
- State: "task files created with correct surface-key and surface-type fields derived from forge surfaces CLI output"
- Side-effect: "forge surfaces CLI invoked during task generation to resolve surface info per file path"

## Outcome "not-found-surface-for-task"
<!-- source: cli-required — surface rule mandates not-found for resource access steps -->
- Preconditions: "breakdown-tasks generates a task for a file path that does not match any configured surface entry"
- Input: "breakdown-tasks processes a file in a path outside all configured surface entries"
- Output: "task generated with empty surface-key and empty surface-type fields, or task assigned to a default surface based on project configuration"
- State: "task file created but lacks surface-specific information"
- Side-effect: "none"

## Outcome "template-variable-sync-failure"
<!-- source: inferred -->
<!-- reasoning: Journey edge case 6b describes prompt template SURFACE_KEY variable mismatch. This outcome covers the case where breakdown-tasks task template still uses old scope field instead of surface-key. Tech design lists 3 skill templates that need variable replacement. -->
- Preconditions: "some prompt or task templates still reference old frontend/backend values or use {{SCOPE}} instead of {{SURFACE_KEY}}"
- Input: "breakdown-tasks generates tasks using outdated templates"
- Output: "generated tasks have scope field instead of surface-key, or contain hardcoded frontend/backend values"
- State: "task files created with incorrect field names or values, violating the surface-key migration contract"
- Side-effect: "none"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
