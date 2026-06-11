---
id: "6"
title: "forge init: TUI confirmation for surfaces"
priority: "P1"
estimated_time: "1.5h"
dependencies: ["5"]
scope: "backend"
breaking: false
type: "coding.enhancement"
mainSession: false
---

# 6: forge init: TUI confirmation for surfaces

## Description

Update the `forge init` TUI flow to display detected surfaces and allow user confirmation/editing. Show conflict signal annotations. Support add/delete operations on the surface mappings.

## Reference Files
- `proposal.md#信号冲突消歧规则` — conflict display in TUI (priority-based selection + annotation)
- `proposal.md#Success-Criteria` — TUI verification criteria (confirm button, edit entry, conflict annotation, add/delete)
- `proposal.md#Key-Scenarios` — scenarios 3-5 for TUI interactions (detection failure, multi-select, user override)

## Acceptance Criteria

- [ ] Detected surfaces displayed in TUI: scalar form shows single type, map form shows path→surface rows
- [ ] Confirm button (or Enter) writes surfaces to config and proceeds
- [ ] Edit entry: each row can enter edit mode to modify path or surface type
- [ ] Conflict annotation: conflicting signals shown as `path: surface (冲突信号: web + api，已按优先级选择 web)`
- [ ] Add: blank row input to add new mapping
- [ ] Delete: select row + press `d` to remove

## Hard Rules

- TUI must work in both scalar and map display modes
- Single-type detection should NOT show path column (scalar = just the type)

## Implementation Notes

- Modify existing init TUI flow in `forge-cli/internal/cmd/init.go`
- Use the same TUI library already used by init (check existing imports)
