---
type: test.gen-contracts
category: test
identity:
  - TaskID
  - TaskType
  - FeatureSlug
context:
  - Mode
  - SurfaceKey
  - SurfaceType
  - SurfaceTypes
  - AcceptanceCriteria
---
Generate test Contract specifications for the {{.FeatureSlug}} feature.{{if .Mode}}
Mode: {{.Mode}}{{end}}
{{if .SurfaceKey}}

## Scope

{{.SurfaceKey}}{{end}}

## Discovery Strategy

Invoke the `/gen-contracts` skill to generate Contract specifications from Journey documents and code reconnaissance.

### Eval Gate by Mode

- **Breakdown mode**: An eval-journey report must exist for all Journeys before proceeding. Check for `testing/<journey>/.eval-report.md` files. If any Journey lacks an eval report or scored below target, abort this task.
- **Quick mode**: The eval-journey gate is skipped. Proceed directly to Contract generation.

## SKIP_EVAL_GATE Directive

When this task runs in Quick mode as an automated pipeline task, SKIP_EVAL_GATE=true is in effect:

- **If SKIP_EVAL_GATE=true**: Skip the eval-journey prerequisite check. Do not require `testing/<journey>/.eval-report.md` files. Proceed directly to code reconnaissance and Contract generation.
- **If SKIP_EVAL_GATE is not set** (Breakdown mode): Require eval reports for all Journeys. Abort if any Journey scored below the eval target threshold.

## Process

Follow the `/gen-contracts` skill process flow:

1. **Resolve Language & Interfaces**: Detect project language and interface types from config
2. **Read Journeys**: Enumerate Journey directories under `docs/features/{{.FeatureSlug}}/testing/` and read each `journey.md`
3. **Code Reconnaissance**: Build the Fact Table by reading source code per the reconnaissance rules
4. **Generate Contracts**: For each Journey, generate one Contract file per Step with six-dimension declarations (Preconditions, Input, Output, State, Side-effect, Invariants). Apply risk-driven Outcome density.
5. **Validate Contracts**: Schema validation for structural completeness. Retry once on failure.
6. **Write Output**: Write Contract files to `docs/features/{{.FeatureSlug}}/testing/<journey>/contracts/` and Fact Table to `.forge/fact-table.json`

## Acceptance Criteria

- [ ] At least 1 Contract file generated per Journey
- [ ] Each Contract has six-dimension declarations with semantic descriptors (no regex)
- [ ] Risk-driven Outcome density targets met per Journey risk level
- [ ] Fact Table written to `.forge/fact-table.json`
- [ ] All Contracts passed schema validation
