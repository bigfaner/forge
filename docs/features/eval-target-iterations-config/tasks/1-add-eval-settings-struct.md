---
id: "1"
title: "Add EvalSettings Go struct and Config integration"
priority: "P0"
estimated_time: "1.5h"
complexity: "medium"
dependencies: []
surface-key: ""
surface-type: "cli"
breaking: false
type: "coding.feature"
mainSession: false
---

# 1: Add EvalSettings Go struct and Config integration

## Description

Add `EvalTypeSettings` and `EvalSettings` Go structs to `forgeconfig/config.go`, and wire `Eval *EvalSettings` into the `Config` struct. The existing reflection routing (`GetConfigValue`/`SetConfigValue`) will automatically support `eval.<type>.target|iterations` keys and `forge config get eval` summary output.

Use `*int` pointer types for `Target` and `Iterations`: `nil` means not configured (fallback to rubric), non-nil overrides. The reflection router's `derefPointer` already returns `errKeyNotFound` for nil pointers, providing the correct "not configured" semantics.

## Reference Files
- `docs/proposals/eval-target-iterations-config/proposal.md` — Feasibility Assessment (Go struct design), Scope > In Scope (items 1-4, 7)
- `forge-cli/pkg/forgeconfig/config.go` — Config struct (line 278), getByPath reflection routing, formatStructSummary

## Acceptance Criteria
- [ ] `EvalTypeSettings` struct defined with `Target *int` (`yaml:"target,omitempty"`) and `Iterations *int` (`yaml:"iterations,omitempty"`)
- [ ] `EvalSettings` struct with 7 eval type fields: proposal, prd, design, ui, journey, contract, consistency (each `EvalTypeSettings`)
- [ ] `Config` struct has `Eval *EvalSettings` field with yaml tag `"eval,omitempty"`
- [ ] `forge config get eval.proposal.target` returns correct integer (e.g. `900`) when config.yaml has the value set
- [ ] `forge config get eval.proposal.target` returns exit code 1 (`errKeyNotFound`) when not configured (nil `*int`)
- [ ] `forge config set eval.proposal.target 850` correctly writes value to config.yaml and `forge config set eval.journey.iterations 5` correctly writes value to config.yaml

## Implementation Notes
- Reflection routing handles `*int` fields: `derefPointer` returns zero Value for nil → `formatValue` returns `errKeyNotFound`. No changes needed to `getByPath`/`setByPath`.
- `setByPath` initializes nil pointers via `reflect.New` before setting, so `forge config set` works on previously-unconfigured fields.
- `formatStructSummary` skips fields where `derefPointer` returns invalid Value, so `forge config get eval` only shows configured sub-types.
- The 7 eval types match rubric files: proposal, prd, design, ui, journey, contract, consistency.
