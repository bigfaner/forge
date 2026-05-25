---
journey: "surface-key-migration"
step: 1
step-action: "run forge surfaces CLI to verify surface detection"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-key-migration/journey.md
---

# Contract: surface-key-migration / Step 1: Run forge surfaces CLI to verify surface detection

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "project has .forge/config.yaml with surfaces field defining surface-key to surface-type mappings, forge surfaces CLI is installed and functional"
- Input: "user executes forge surfaces with a file path argument for a configured project path"
- Output: "CLI returns surface-key and surface-type via longest-prefix-match, e.g., forge surfaces frontend/src returns surface-key: admin-panel, surface-type: web"
- State: "surface detection verified, CLI output available for downstream consumption"
- Side-effect: "none"

## Outcome "no-match"
- Preconditions: "queried path does not match any configured surface entry in config.yaml"
- Input: "user executes forge surfaces with an unrecognized path argument"
- Output: "CLI exits with code 1, stderr contains error message with recovery hint (run forge init to configure surfaces)"
- State: "no surface information returned, downstream components handle gracefully"
- Side-effect: "none"

## Outcome "ambiguous-match"
<!-- source: inferred -->
<!-- reasoning: Journey edge case 4b describes overlapping path prefixes causing ambiguous match. Forge surfaces CLI uses longest-prefix-match; when two entries have identical prefix length, an error must be returned. This boundary is important for migration correctness. -->
- Preconditions: "file path matches multiple surface entries via longest-prefix-match with identical prefix length"
- Input: "user executes forge surfaces with an ambiguous path that has overlapping prefixes in config"
- Output: "CLI returns error indicating ambiguous configuration, lists the conflicting surface entries"
- State: "no surface information returned, user must resolve configuration ambiguity"
- Side-effect: "none"

## Outcome "not-found-cli-missing"
<!-- source: cli-required — surface rule mandates not-found for resource access steps -->
- Preconditions: "forge surfaces command is not available (CLI not installed or version too old)"
- Input: "user or downstream component attempts to invoke forge surfaces CLI"
- Output: "error to stderr with the CLI output and recovery hint (check forge CLI is installed and version >= required version), exits with exit code 1 (retryable)"
- State: "no surface detection possible, downstream components fall back or abort"
- Side-effect: "none"

## Journey Invariants

- surface-type always belongs to the fixed set (web, api, cli, tui, mobile), never user-defined
- surface-key is always user-defined and unique within a project's config.yaml
- Migration is phased: Phase 1 (data model) -> Phase 2 (upstream adapters) -> Phase 3 (downstream consumers), strict sequential dependency
- Projects without surfaces configuration produce identical output to the pre-feature baseline (zero regression guarantee)
- forge task migrate must exist before any task read operations work on old-format task files
- All 7+ components surface-key value domains are synchronized after migration
