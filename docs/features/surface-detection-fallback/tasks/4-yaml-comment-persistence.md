---
id: "4"
title: "Persist source annotation as YAML comment in config"
priority: "P2"
estimated_time: "1h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 4: Persist source annotation as YAML comment in config

## Description
When writing surfaces config after inference, append a YAML comment alongside the surface value for human context (e.g., `surfaces: cli  # source: inference:cmd-dir`). Use `yaml.Node` round-trip API to preserve comments during config read/write cycles.

## Reference Files
- `proposal.md#Constraints-Dependencies` — YAML comment persistence spec, yaml.v3 Node round-trip API
- `proposal.md#Key-Risks` — risk of YAML comment stripping during config round-trip, mitigation via Node manipulation
- `proposal.md#Success-Criteria` — YAML comment persistence criterion (#15)
- `forge-cli/pkg/forgeconfig/config.go` — SurfacesMap custom YAML marshal/unmarshal (L145-197), config writing

## Acceptance Criteria
- [ ] After `forge init` infers `cli` from `cmd/` directory, config file contains `surfaces: cli  # source: inference:cmd-dir`
- [ ] Comment is present in file but ignored by YAML unmarshaler (no schema change to SurfacesMap)
- [ ] On re-detection, existing comment is read and displayed for context in the TUI prompt
- [ ] If comment is stripped during round-trip (e.g., by external editor), detection still works correctly — comment is informational only
- [ ] Integration test reads config file after `forge init` and asserts comment is present

## Hard Rules
- Use `yaml.Node` round-trip API for comment manipulation, NOT string concatenation or regex injection
- Comment is purely informational — must not affect config parsing or behavior

## Implementation Notes
- Forge uses `gopkg.in/yaml.v3` which preserves comments via `yaml.Node` API
- When writing config: unmarshal to `yaml.Node`, find surfaces node, append `HeadComment` or `LineComment` with source info, marshal back
- When reading config for re-detection: parse `yaml.Node` tree to extract comment from surfaces node for display
- Comment format: `# source: inference:cmd-dir` or `# source: dependency:cobra`
- May need to modify `SurfacesMap` marshal/unmarshal or add a separate comment-writing function called after config write
- Integration test: create temp dir with `go.mod` + `cmd/` structure, run init flow, read config file, assert `# source: inference:cmd-dir` present
