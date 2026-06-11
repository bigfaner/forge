---
id: "2"
title: "PRD TUI navigation template and write-prd awareness"
priority: "P1"
estimated_time: "1h"
dependencies: []
scope: "all"
breaking: false
type: "implementation"
mainSession: false
---

# 2: PRD TUI navigation template and write-prd awareness

## Description
Modify the PRD UI functions template to support TUI navigation architecture and update the write-prd SKILL.md to handle platform=tui. Currently the Navigation Architecture table only has web/mobile structures; TUI needs its own Keymap + Panel Layout + Mode structure (proposal D4).

This task is independent of the ui-design changes (Task 3) since PRD is upstream of ui-design in the pipeline.

## Reference Files
- `docs/proposals/tui-ui-design/proposal.md` — Source proposal (D4 section)
- `plugins/forge/skills/write-prd/templates/prd-ui-functions.md` — Current PRD UI functions template
- `plugins/forge/skills/write-prd/SKILL.md` — Current write-prd skill logic

## Affected Files

### Create
| File | Description |
|------|-------------|
| — | — |

### Modify
| File | Changes |
|------|---------|
| `plugins/forge/skills/write-prd/templates/prd-ui-functions.md` | Add TUI Navigation Architecture template: Keymap table, Panel Layout table, Modes table, Navigation Rules — conditionally rendered when platform=tui |
| `plugins/forge/skills/write-prd/SKILL.md` | Add platform=tui awareness so SKILL.md knows to render the TUI navigation section |

### Delete
| File | Reason |
|------|--------|
| — | — |

## Acceptance Criteria
- [ ] `prd-ui-functions.md` includes a TUI Navigation Architecture section with the structure from proposal D4: `Platform: tui`, Keymap table `[Key | Action | Context/Mode]`, Panel Layout table `[Panel | View | Position | Size Hint]`, Modes table `[Mode | Description | Default Keybindings]`, Navigation Rules
- [ ] TUI navigation section is conditionally rendered — only appears when platform=tui, does not affect web/mobile templates
- [ ] `write-prd/SKILL.md` references the TUI navigation template and triggers it when platform=tui
- [ ] Existing web/mobile PRD generation behavior unchanged

## Implementation Notes
- Study how `prd-ui-functions.md` currently handles web vs mobile navigation to follow the same conditional rendering pattern
- Proposal D4 provides the exact table structure — use it directly
- The TUI navigation structure is fundamentally different from web/mobile (keyboard-driven vs pointer-driven), so it needs its own section rather than extending the existing tables
