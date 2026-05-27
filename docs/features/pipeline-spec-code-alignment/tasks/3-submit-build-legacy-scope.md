---
id: "3"
title: "Fix submit.go scope hardcode and build.go legacy Scope field"
priority: "P0"
estimated_time: "1h"
dependencies: [2]
surface-key: "cli"
surface-type: "cli"
breaking: false
type: "coding.fix"
mainSession: false
---

# 3: Fix submit.go scope hardcode and build.go legacy Scope field

## Description

Fix two related legacy scope issues in Go code:

1. **submit.go** (~line 157): `validateQualityGate` hardcodes `scope=""` instead of using `t.SurfaceKey`. In multi-surface projects, scoped recipes (e.g., `just compile backend`) never execute during task submission. Has a TODO comment confirming this is a known defect. Fix: pass `scope=t.SurfaceKey`.

2. **build.go** (~lines 134, 307): Still populates the legacy `Scope` field in index.json from frontmatter. This causes `CheckLegacyScope` migration detection to self-sustain — as long as .md files contain `scope:` key, the legacy field is never cleaned up. Fix: stop writing `Scope` to index.json; rely on `SurfaceKey`/`SurfaceType` only.

## Reference Files
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Problem` — Evidence G9 (submit.go scope hardcode) and G11 (build.go legacy Scope)
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Proposed-Solution` — Cluster 7 subset for submit.go and build.go
- `docs/proposals/pipeline-spec-code-alignment/proposal.md#Success-Criteria` — SC for submit.go using t.SurfaceKey

## Acceptance Criteria
- [ ] `submit.go` passes `scope=t.SurfaceKey` to `validateQualityGate`
- [ ] `build.go` no longer writes `Scope` field to index.json
- [ ] Existing tests pass (`go test ./...`)

## Hard Rules
- The `Scope` field in the Go struct may still exist for reading legacy files — do not remove the struct field, only stop writing it
- Ensure `CheckLegacyScope` still works for detecting old files that do have `scope:` in frontmatter

## Implementation Notes
- submit.go has a TODO comment at the relevant line — look for it
- build.go writes the Scope field in two places (~line 134 and ~line 307)
- After this change, `CheckLegacyScope` will correctly detect and offer migration for any .md files that still have `scope:` frontmatter
