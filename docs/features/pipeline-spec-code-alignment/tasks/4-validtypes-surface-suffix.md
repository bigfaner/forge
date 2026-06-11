---
id: "4"
title: "Add surface-suffixed type variants to ValidTypes"
priority: "P1"
estimated_time: "1h"
dependencies: [3]
surface-key: "cli"
surface-type: "cli"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 4: Add surface-suffixed type variants to ValidTypes

## Description

`autogogen.go` generates surface-suffixed task types like `test.gen-scripts.cli`, but `ValidTypes` in types.go does not include these variants. `prompt.Synthesize` at line 89 rejects them as invalid. This causes prompt synthesis to fail for surface-suffixed test tasks.

Two possible fixes (choose one):
1. Add all surface-suffixed variants to ValidTypes (static enumeration)
2. Skip validation for auto-generated types in Synthesize (dynamic approach)

Option 2 is preferred: auto-generated types follow a predictable pattern (`{category}.{subtype}.{surface-key}`), so validation can be pattern-based rather than enumeration-based.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence G1 (ValidTypes missing surface-suffixed variants, Synthesize rejects them)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 7 ValidTypes bullet
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for ValidTypes completeness

## Acceptance Criteria
- [ ] Surface-suffixed types (e.g., `test.gen-scripts.cli`) pass Synthesize validation
- [ ] All existing ValidTypes entries still work
- [ ] Existing tests pass (`go test ./...`)

## Hard Rules
- Do not hardcode surface keys into ValidTypes — use pattern matching or skip validation for generated types
- Existing non-generated type validation must remain strict

## Implementation Notes
- The autogogen.go file generates types based on surface configuration. The format is `{base-type}.{surface-key}`.
- In prompt.go Synthesize function (~line 89), the validation check should either: (a) match against a pattern for generated types, or (b) skip validation when the type prefix matches a known base type.
