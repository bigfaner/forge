---
id: "3"
title: "forge surfaces CLI command"
priority: "P0"
estimated_time: "1.5h"
dependencies: ["1", "2"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 3: forge surfaces CLI command

## Description

Add independent `forge surfaces` CLI command with three sub-invocations: full listing, path query, and type list. Follows the exit code contract defined in the proposal.

## Reference Files
- `proposal.md#CLI-命令与退出码契约` — exit code table, output format for scalar vs map, gen-journeys calling contract
- `proposal.md#Config-结构` — dual-form output (scalar outputs single type, map outputs path=surface lines)
- `proposal.md#路径规范化与匹配算法` — path query delegates to matching algorithm
- `proposal.md#Success-Criteria` — CLI output verification, interfaces ignored verification

## Acceptance Criteria

- [ ] `forge surfaces` — scalar form: outputs single type (e.g., `api`, exit 0); map form: outputs `path=surface` per line (exit 0)
- [ ] `forge surfaces <path>` — returns surface type string (exit 0) or stderr error + exit 1 with manual config hint
- [ ] `forge surfaces --types` — space-separated deduplicated type list (exit 0)
- [ ] Scalar form: `forge surfaces <any-path>` always returns the single value (exit 0)
- [ ] `interfaces` field completely ignored — config with both `surfaces: {frontend: web}` and `interfaces: [api, cli]` outputs only `web`
- [ ] Command registered at top level in `root.go` (not under `forge config`)

## Hard Rules

- Do NOT register under `forge config` — independent command for separation of concerns
- Exit code 1 MUST output to stderr, not stdout
- Output format: no extra formatting, raw strings parseable by scripts

## Implementation Notes

- Follow top-level command pattern: create `forge-cli/internal/cmd/surfaces.go`, register in `root.go` init()
- Read config via `forgeconfig.ReadSurfaces()`, use `MatchSurface()` for path queries
- For `--types`: extract unique values from Surfaces map, filter unknown types
