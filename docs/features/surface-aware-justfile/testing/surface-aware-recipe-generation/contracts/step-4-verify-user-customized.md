---
journey: "surface-aware-recipe-generation"
step: 4
step-action: "verify user-customized protection"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-aware-recipe-generation/journey.md
---

# Contract: surface-aware-recipe-generation / Step 4: Verify user-customized protection

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "justfile exists with at least one recipe marked with # user-customized comment, user has manually modified that recipe"
- Input: "user re-runs init-justfile without --force-regenerate flag"
- Output: "recipes marked with # user-customized comment are preserved unchanged. init-justfile outputs a diff summary showing what would have changed"
- State: "user-customized recipes remain identical to their modified state, non-customized recipes are regenerated normally"
- Side-effect: "none"

## Outcome "already-exists-customized"
<!-- source: cli-required — surface rule mandates already-exists for resource creation steps -->
- Preconditions: "justfile exists with user-customized recipe that init-justfile would regenerate, user passes --force-regenerate flag"
- Input: "user re-runs init-justfile with --force-regenerate flag"
- Output: "all recipes regenerated including those previously marked # user-customized, diff summary shows overridden recipes"
- State: "all recipes regenerated from templates, # user-customized markers removed from regenerated recipes"
- Side-effect: "previous user modifications to recipes are lost"

## Journey Invariants

- init-justfile never silently overwrites a recipe marked with # user-customized without --force-regenerate
- Projects without surfaces configuration always produce output identical to the pre-feature behavior (zero regression)
- All generated recipes include dual-platform ([linux]/[windows]) variants where applicable
- cli/tui surfaces never generate run or probe recipes
- Mixed-project aggregation recipes always list services in dependency order (api before web)
