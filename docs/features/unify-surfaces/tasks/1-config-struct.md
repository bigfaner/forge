---
id: "1"
title: "Config struct refactor: Surfaces dual-form + remove Interfaces"
priority: "P0"
estimated_time: "2h"
dependencies: []
scope: "backend"
breaking: true
type: "coding.refactor"
mainSession: false
---

# 1: Config struct refactor: Surfaces dual-form + remove Interfaces

## Description

Replace `Interfaces []string` with `Surfaces map[string]string` in the Go Config struct. Implement custom `UnmarshalYAML`/`MarshalYAML` for scalar/map dual-form support. Add first-run auto-migration logic: detect old `interfaces` field → single-type auto-convert to scalar / multi-type error guidance.

This is the foundation task — all other tasks depend on the new `Surfaces` field being available in the Config struct.

## Reference Files
- `proposal.md#Config-结构` — Go struct declaration, dual-form YAML examples, custom marshal/unmarshal rules, omitempty rationale
- `proposal.md#Non-Functional-Requirements` — upgrade migration strategy (auto-migration for single-type, error for multi-type)
- `proposal.md#Key-Risks` — empty surfaces serialization risk, omitempty risk

## Acceptance Criteria

- [ ] `Interfaces []string` field removed from `Config` struct in `forge-cli/pkg/forgeconfig/config.go`
- [ ] `Surfaces map[string]string \`yaml:"surfaces"\`` added (no `omitempty`)
- [ ] Custom `UnmarshalYAML`: scalar `"api"` → `map[string]string{".": "api"}`; map form used as-is
- [ ] Custom `MarshalYAML`: single entry with key `"."` → scalar; otherwise → map
- [ ] Auto-migration: config has `interfaces: [api]` but no `surfaces` → auto-write `surfaces: api` (scalar) + console prompt
- [ ] Auto-migration: config has `interfaces: [web, api]` but no `surfaces` → error exit with guidance to run `forge init`
- [ ] Empty map (0 entries) serializes as `surfaces: {}`, never omitted
- [ ] Existing tests pass after struct change

## Hard Rules

- YAML tag must NOT use `omitempty` — empty `surfaces` must persist as `surfaces: {}` to avoid the silent-skip bug
- Custom marshal/unmarshal must handle: scalar string, map, nil/empty map
- Auto-migration runs on ANY forge command startup, not just init

## Implementation Notes

- Key files: `forge-cli/pkg/forgeconfig/config.go` (struct definition), `forge-cli/pkg/forgeconfig/detect.go` (ReadInterfaces → rename to ReadSurfaces)
- The `ReadInterfaces` function in `detect.go` currently returns `cfg.Interfaces` — this needs to be replaced with a `ReadSurfaces` function that returns `map[string]string`
- `uiInterfaces` map in `forge-cli/pkg/task/autogen.go` uses old naming (`"web"`, `"mobile"`) — will be cleaned in task 4
