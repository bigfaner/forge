---
id: "6"
title: "Multi-platform manifest output support"
priority: "P2"
estimated_time: "30m"
dependencies: ["3"]
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 6: Multi-platform manifest output support

## Description
Modify `manifest-update-ui.md` to handle multi-platform file outputs. When a feature declares multiple platforms, the manifest must list all ui-design files and prototype directories per platform.

## Reference Files
- `docs/proposals/tui-ui-design/proposal.md` — Source proposal (D7 section)
- `plugins/forge/skills/ui-design/templates/manifest-update-ui.md` — Current manifest update template

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/ui-design/templates/manifest-update-ui.md` | Support listing multiple ui-design files (ui-design-web.md, ui-design-tui.md) and multiple prototype directories (prototype/web/, prototype/tui/) for multi-platform features |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] Single-platform feature (web only): manifest lists `ui-design.md` and `prototype/` — behavior unchanged
- [ ] Multi-platform feature (web + tui): manifest lists `ui-design-web.md`, `ui-design-tui.md`, `prototype/web/`, `prototype/tui/`
- [ ] Single TUI feature: manifest lists `ui-design-tui.md` and `prototype/`

## Implementation Notes
- Study how `manifest-update-ui.md` currently references ui-design output files — the multi-platform case adds platform suffixes to filenames and creates separate prototype directories
- Proposal D7 specifies the exact file structure for multi-platform output
- This is a small change — the template just needs to iterate over platforms instead of assuming a single output
