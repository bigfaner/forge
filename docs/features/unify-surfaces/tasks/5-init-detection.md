---
id: "5"
title: "forge init: surface auto-detection logic"
priority: "P0"
estimated_time: "3h"
dependencies: ["1"]
scope: "backend"
breaking: false
type: "coding.feature"
mainSession: false
---

# 5: forge init: surface auto-detection logic

## Description

Implement surface auto-detection in `forge init`: scan project files, match dependency signals to surface types, resolve conflicts. Write detection results as `surfaces` to config. Support workspace/monorepo detection with depth limits.

## Reference Files
- `proposal.md#Technical-Feasibility` — signal table (package.json/go.mod/Cargo.toml/pyproject.toml patterns)
- `proposal.md#检测遍历策略` — traversal rules (depth limits, exclusion dirs, workspace detection, non-workspace root)
- `proposal.md#信号冲突消歧规则` — priority table (web > mobile > api > cli > tui), conflict resolution flow
- `proposal.md#Key-Scenarios` — scenarios 2-7 for detection edge cases (monorepo, failure, multi-select, Next.js, conflicts)
- `proposal.md#Success-Criteria` — detect at least 3 types, FORGE_DETECT_DEPTH validation

## Acceptance Criteria

- [ ] Detects surface types from: `package.json` (react/express/commander/blessed), `go.mod` (gin/cobra/bubbletea), `Cargo.toml` (actix/clap/ratatui), `AndroidManifest.xml`, `*.xcodeproj`, `pyproject.toml` (flask/click)
- [ ] Workspace detection: `pnpm-workspace.yaml` or `package.json#workspaces` → skip root deps, scan subdirs
- [ ] Non-workspace: root detected as `"."`, output as scalar form
- [ ] Depth limit: default 3, configurable via `FORGE_DETECT_DEPTH` (1-10, 0 is invalid with error)
- [ ] Exclusion dirs: `node_modules`, `.git`, `vendor`, `dist`, `build`, `__pycache__`, `.next`, `target`
- [ ] Signal conflict: auto-resolve via priority table (web > mobile > api > cli > tui)
- [ ] Single-type result → scalar output (`surfaces: api`); multi-type → map output
- [ ] Detection completes in <5 seconds

## Hard Rules

- `FORGE_DETECT_DEPTH=0` MUST produce an init error, not silently default
- Root deps in workspace projects MUST be skipped (root package.json is workspace config, not surface signal)
- Detection depth must NOT be unlimited — OOM/hang risk

## Implementation Notes

- Suggested location: new file `forge-cli/internal/cmd/surface_detect.go` or `forge-cli/pkg/forgeconfig/detect_surface.go`
- Keep detection logic separate from init flow for testability
- For `package.json` parsing: use `encoding/json`, extract `dependencies` + `devDependencies`
- For `go.mod` parsing: simple regex or line scan for `require` directives
