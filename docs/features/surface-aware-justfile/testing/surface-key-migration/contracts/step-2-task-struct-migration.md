---
journey: "surface-key-migration"
step: 2
step-action: "verify Task Go struct migration"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-key-migration/journey.md
---

# Contract: surface-key-migration / Step 2: Verify Task Go struct migration

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "task data model has been migrated, task file format uses new surface fields"
- Input: "user inspects task data format after implementation"
- Output: "task files use surface-key and surface-type fields instead of the old scope field, data serialization includes both new fields, the old scope field is no longer present"
- State: "task data model uses surface-key and surface-type exclusively, old scope field removed from all task file formats"
- Side-effect: "none"

## Outcome "legacy-scope-detected"
- Preconditions: "existing task files contain the old scope field in their frontmatter"
- Input: "user attempts to read task information via task status or similar command"
- Output: "blocking error (exit 2) to prevent silent data loss, error message indicates migration required with count of affected tasks and recovery command to run migration"
- State: "task read blocked until migration is performed"
- Side-effect: "none"

## Outcome "migration-via-task-migrate"
<!-- source: inferred -->
<!-- reasoning: Tech design defines a task migrate command to handle migration from scope to surface-key/surface-type. This is the recovery path from legacy-scope-detected outcome. -->
- Preconditions: "legacy task files exist with scope field, task migration command is available"
- Input: "user runs task migration command to convert scope to surface-key and surface-type"
- Output: "migration command scans task index, maps scope values to surface-key and surface-type via surface detection, updates task files and index in place"
- State: "all task files updated with surface-key and surface-type, scope field removed from frontmatter"
- Side-effect: "index.json and task frontmatter files modified on disk"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
