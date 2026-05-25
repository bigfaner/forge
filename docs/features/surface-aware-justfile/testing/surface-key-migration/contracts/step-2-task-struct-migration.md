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
- Preconditions: "Go data model migration is completed, task/types.go has been updated"
- Input: "user inspects task/types.go after implementation"
- Output: "Scope field replaced by SurfaceKey (string), new SurfaceType (string) field added, GetSurfaceKey() provides backward-compatible access, JSON serialization includes surfaceKey and surfaceType fields"
- State: "Task struct uses SurfaceKey and SurfaceType, Scope field removed entirely, no backward compatibility layer retained"
- Side-effect: "none"

## Outcome "legacy-scope-detected"
- Preconditions: "existing task files have scope: frontend or scope: backend in their frontmatter"
- Input: "user attempts to read the task via forge task status or similar command"
- Output: "blocking error (exit 2) to prevent silent data loss, stderr message indicates migration required: found N tasks with legacy scope field but no surface-key, run forge breakdown-tasks or forge quick-tasks to regenerate tasks"
- State: "task read blocked until migration is performed"
- Side-effect: "none"

## Outcome "migration-via-forge-task-migrate"
<!-- source: inferred -->
<!-- reasoning: Tech design defines forge task migrate subcommand to handle migration from scope to surface-key/surface-type. This is the recovery path from legacy-scope-detected outcome. Journey edge case 2b mentions forge task migrate capability. -->
- Preconditions: "legacy task files exist with scope field, forge task migrate command is available"
- Input: "user runs forge task migrate to convert scope to surface-key + surface-type"
- Output: "forge task migrate scans index.json, maps scope field values via forge surfaces CLI to surface-key + surface-type, updates index.json and frontmatter files in place"
- State: "all task files updated with surface-key and surface-type, scope field removed from frontmatter"
- Side-effect: "index.json and task frontmatter files modified on disk"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
