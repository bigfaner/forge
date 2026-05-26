---
journey: "surface-aware-recipe-generation"
step: 3
step-action: "verify generated justfile contains surface-specific recipes"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-aware-recipe-generation/journey.md
---

# Contract: surface-aware-recipe-generation / Step 3: Verify generated justfile contains surface-specific recipes

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "init-justfile has completed successfully with valid surface configuration"
- Input: "user inspects the generated justfile"
- Output: "justfile contains surface-type-specific recipes: for web/api surfaces includes dev (background start), probe (retry polling), test, test-teardown recipes; for cli/tui surfaces includes dev, test recipes only (no run, no probe). Each recipe includes [linux]/[windows] dual-platform variants"
- State: "justfile has complete set of surface-specific recipes matching the configured surface type"
- Side-effect: "none"

## Outcome "recipe-not-found"
<!-- source: cli-required — surface rule mandates not-found for resource access steps -->
- Preconditions: "init-justfile has completed but expected recipe is missing from the justfile"
- Input: "user inspects the justfile for a recipe that should exist based on the configured surface type"
- Output: "expected recipe not present in justfile, indicating init-justfile did not generate all required recipes"
- State: "justfile is incomplete, missing one or more expected surface-specific recipes"
- Side-effect: "none"

## Journey Invariants

- init-justfile never silently overwrites a recipe marked with # user-customized without --force-regenerate
- Projects without surfaces configuration always produce output identical to the pre-feature behavior (zero regression)
- All generated recipes include dual-platform ([linux]/[windows]) variants where applicable
- cli/tui surfaces never generate run or probe recipes
- Mixed-project aggregation recipes always list services in dependency order (api before web)
- Step-specific: recipe completeness is determined by the surface type's rule file requirements
