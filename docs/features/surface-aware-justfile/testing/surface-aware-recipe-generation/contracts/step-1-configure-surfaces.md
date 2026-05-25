---
journey: "surface-aware-recipe-generation"
step: 1
step-action: "configure surfaces in config.yaml"
generated: "2026-05-26"
sources:
  - docs/features/surface-aware-justfile/testing/surface-aware-recipe-generation/journey.md
---

# Contract: surface-aware-recipe-generation / Step 1: Configure surfaces in config.yaml

<!-- gen-contracts: do not edit manually. Regenerate via /gen-contracts. -->

## Outcome "success"
- Preconditions: "project has .forge/config.yaml without surfaces field, or with valid surfaces field"
- Input: "user adds surfaces field to .forge/config.yaml with valid mapping of surface-key to surface-type, e.g., surfaces: {admin-panel: web}"
- Output: "config.yaml updated with surfaces field containing a valid map of surface-key to surface-type (one of web/api/cli/tui/mobile)"
- State: ".forge/config.yaml surfaces field contains valid entries ready for init-justfile consumption"
- Side-effect: "none"

## Outcome "invalid-surface-type"
- Preconditions: "user is editing .forge/config.yaml and enters an unrecognized surface-type value"
- Input: "surfaces field contains an unrecognized surface type, e.g., {my-surface: desktop}"
- Output: "descriptive error message to stderr indicating the unsupported surface type, listing the 5 supported types (web/api/cli/tui/mobile)"
- State: "config.yaml unchanged or rejected before write"
- Side-effect: "none"

## Outcome "config-validation-error"
- Preconditions: "user is editing .forge/config.yaml and enters invalid surface configuration format"
- Input: "surfaces field has invalid format (scalar value, list instead of map) or surface-key contains characters not allowed in just recipe names (e.g., my/surface or surface+key)"
- Output: "descriptive YAML parse error or validation error to stderr with recovery hint: check .forge/config.yaml surfaces field format, should be map of string to string, e.g., {admin-panel: web}. For invalid characters, suggests valid alternatives matching alphanumeric, dash, and underscore only"
- State: "config.yaml unchanged or rejected before write"
- Side-effect: "none"

## Journey Invariants

- init-justfile never silently overwrites a recipe marked with # user-customized without --force-regenerate
- Projects without surfaces configuration always produce output identical to the pre-feature behavior (zero regression)
- All generated recipes include dual-platform ([linux]/[windows]) variants where applicable
- cli/tui surfaces never generate run or probe recipes
- Mixed-project aggregation recipes always list services in dependency order (api before web)
- Step-specific: surface-type values are restricted to the fixed set (web/api/cli/tui/mobile) at configuration time
