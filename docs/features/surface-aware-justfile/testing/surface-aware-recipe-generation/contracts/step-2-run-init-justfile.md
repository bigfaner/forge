---
journey: "surface-aware-recipe-generation"
step: 2
step-action: "run init-justfile"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-aware-recipe-generation/journey.md
---

# Contract: surface-aware-recipe-generation / Step 2: Run init-justfile

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: ".forge/config.yaml exists with valid surfaces field, just binary installed and version >= 1.4.0, surface rule file exists for configured surface type"
- Input: "user executes init-justfile skill via Forge CLI"
- Output: "init-justfile detects the surface type from config.yaml, loads the corresponding surface rule file, and generates surface-specific recipes"
- State: "justfile created or updated with surface-specific recipes corresponding to the configured surface types"
- Side-effect: "justfile written to disk, surface rule files loaded from skills/init-justfile/rules/surfaces/<type>.md"

## Outcome "no-surfaces-configured"
- Preconditions: ".forge/config.yaml exists but has no surfaces field, or the field is empty"
- Input: "user executes init-justfile skill via Forge CLI"
- Output: "init-justfile generates only language-template-based recipes (compile/build/lint/fmt) with no orchestration recipes, producing output identical to the current pre-feature behavior"
- State: "justfile contains only language-template recipes, zero regression from pre-feature baseline"
- Side-effect: "none"

## Outcome "just-version-below-minimum"
- Preconditions: "installed just binary version is below 1.4.0"
- Input: "user executes init-justfile skill via Forge CLI"
- Output: "error message to stderr with the current just version and the required version (just >= 1.4.0)"
- State: "no justfile generated or modified"
- Side-effect: "process exits with exit code 2 (blocking)"

## Outcome "surface-rule-file-missing"
- Preconditions: "a supported surface type is configured (e.g., web) but the corresponding rule file (skills/init-justfile/rules/surfaces/web.md) does not exist"
- Input: "user executes init-justfile skill via Forge CLI"
- Output: "error message to stderr with the missing file path and a recovery hint (run init-justfile to regenerate rule files)"
- State: "no justfile generated or modified"
- Side-effect: "process exits with exit code 2 (blocking)"

## Journey Invariants

- init-justfile never silently overwrites a recipe marked with # user-customized without --force-regenerate
- Projects without surfaces configuration always produce output identical to the pre-feature behavior (zero regression)
- All generated recipes include dual-platform ([linux]/[windows]) variants where applicable
- cli/tui surfaces never generate run or probe recipes
- Mixed-project aggregation recipes always list services in dependency order (api before web)
- Step-specific: justfile generation is atomic — either all recipes are written or none are
