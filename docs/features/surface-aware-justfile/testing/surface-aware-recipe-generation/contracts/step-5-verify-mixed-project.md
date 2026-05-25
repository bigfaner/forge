---
journey: "surface-aware-recipe-generation"
step: 5
step-action: "verify mixed-project recipe generation"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-aware-recipe-generation/journey.md
---

# Contract: surface-aware-recipe-generation / Step 5: Verify mixed-project recipe generation

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "config.yaml has multiple surfaces configured (e.g., {admin-panel: web, payment-service: api})"
- Input: "user runs init-justfile on a project with multiple configured surfaces"
- Output: "generated justfile contains per-surface-key prefixed recipes (dev-admin-panel, dev-payment-service, probe-admin-panel, probe-payment-service, etc.) and aggregation recipes: dev (starts all services in dependency order), test (runs all test sequences). Orchestration order comment at the top of the justfile"
- State: "justfile has complete set of prefixed per-surface recipes and aggregation recipes, services listed in dependency order (api before web)"
- Side-effect: "none"

## Outcome "single-surface-fallback"
<!-- source: inferred -->
<!-- reasoning: mixed-project step tests multi-surface, but single-surface is the simpler degenerate case where no prefix is needed. Verifying the boundary between single and multi-surface behavior. -->
- Preconditions: "config.yaml has exactly one surface configured"
- Input: "user runs init-justfile on a single-surface project"
- Output: "generated justfile contains unprefixed recipes (dev, test, probe, test-teardown) matching the single surface type's rule file, no surface-key prefix in recipe names"
- State: "justfile has recipes matching the single configured surface type without key prefixes"
- Side-effect: "none"

## Journey Invariants

- init-justfile never silently overwrites a recipe marked with # user-customized without --force-regenerate
- Projects without surfaces configuration always produce output identical to the pre-feature behavior (zero regression)
- All generated recipes include dual-platform ([linux]/[windows]) variants where applicable
- cli/tui surfaces never generate run or probe recipes
- Mixed-project aggregation recipes always list services in dependency order (api before web)
- Step-specific: aggregation recipes use surface-key prefixes derived from config.yaml keys, not surface types
