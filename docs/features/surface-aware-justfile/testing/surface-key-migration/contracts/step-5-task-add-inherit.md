---
journey: "surface-key-migration"
step: 5
step-action: "verify forge task add inherits surface fields"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-key-migration/journey.md
---

# Contract: surface-key-migration / Step 5: Verify forge task add inherits surface fields

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "source task has surface-key and surface-type fields, task add command supports surface field inheritance"
- Input: "user runs task add with a source task that has surface-key/surface-type"
- Output: "new task inherits surface-key and surface-type from the source task. When no source task exists and the project has a single surface, the unique surface-type is auto-filled"
- State: "new task file created with inherited surface-key and surface-type matching the source task"
- Side-effect: "none"

## Outcome "multi-surface-ambiguous"
- Preconditions: "forge task add is called without a source task, and the project has multiple surfaces configured"
- Input: "user runs forge task add without --source-task-id on a multi-surface project"
- Output: "surface-type cannot be auto-filled due to ambiguity, command requires explicit --surface-type flag or outputs error listing available surfaces"
- State: "no task created until surface-type is explicitly specified"
- Side-effect: "none"

## Outcome "already-exists-task"
<!-- source: cli-required — surface rule mandates already-exists for resource creation steps -->
- Preconditions: "a task with the same ID or conflicting surface-key assignment already exists"
- Input: "user runs forge task add with parameters that conflict with an existing task"
- Output: "error message indicating task ID conflict or surface-key assignment conflict"
- State: "no new task created, existing task index unchanged"
- Side-effect: "none"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
