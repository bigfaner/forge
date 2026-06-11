---
status: "completed"
started: "2026-05-24 15:13"
completed: "2026-05-24 15:21"
time_spent: "~8m"
---

# Task Record: 5 forge init: surface auto-detection logic

## Summary
Implemented surface auto-detection logic in forge init: scans project files, matches dependency signals to surface types, resolves conflicts via priority table. Supports workspace/monorepo detection with configurable depth limits.

## Changes

### Files Created
- forge-cli/pkg/forgeconfig/detect_surface.go
- forge-cli/pkg/forgeconfig/detect_surface_test.go

### Files Modified
无

### Key Decisions
- Placed detection logic in forgeconfig package alongside existing config/match code for cohesion
- Special-case react-native+react conflict: suppress web signal from react when react-native present (react is shared dependency, not independent web signal)
- Depth counting starts at 0 for root's immediate children (depth 1), incrementing per level down
- Non-workspace projects scan root first, then subdirs; collapse to scalar if all paths share same type

## Test Results
- **Tests Executed**: Yes
- **Passed**: 28
- **Failed**: 0
- **Coverage**: 87.2%

## Acceptance Criteria
- [x] Detects surface types from package.json (react/express/commander/blessed), go.mod (gin/cobra/bubbletea), Cargo.toml (actix/clap/ratatui), AndroidManifest.xml, *.xcodeproj, pyproject.toml (flask/click)
- [x] Workspace detection: pnpm-workspace.yaml or package.json#workspaces skips root deps, scans subdirs
- [x] Non-workspace: root detected as '.', output as scalar form
- [x] Depth limit: default 3, configurable via FORGE_DETECT_DEPTH (1-10, 0 is invalid with error)
- [x] Exclusion dirs: node_modules, .git, vendor, dist, build, __pycache__, .next, target
- [x] Signal conflict: auto-resolve via priority table (web > mobile > api > cli > tui)
- [x] Single-type result maps to scalar output (surfaces: api); multi-type maps to map output
- [x] Detection completes in <5 seconds

## Notes
All 28 new detection tests pass. Existing forgeconfig tests (config, surfaces, match, migration) all pass unchanged. Task package tests pass.
